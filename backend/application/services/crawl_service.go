package services

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"backend/application/commands"
	"backend/application/queries"
	"backend/domain"
	"backend/infrastructure/persistence"

	"github.com/PuerkitoBio/goquery"
)

// CrawlService handles URL crawling business logic
type CrawlService struct{
	crawlResultRepo persistence.CrawlResultRepository
}

func NewCrawlService(repo persistence.CrawlResultRepository) *CrawlService {
	return &CrawlService{
		crawlResultRepo: repo,
	}
}

func (s *CrawlService) Crawl(cmd commands.CrawlCommand) (domain.CrawlResult, error) {
	result := domain.CrawlResult{
		HeadingCounts: make(map[string]int),
	}

	parsedURL, err := url.Parse(cmd.URL)
	if err != nil {
		result.Error = domain.ErrInvalidURLFormat.Error()
		return result, domain.ErrInvalidURLFormat
	}

	res, err := http.Get(cmd.URL)
	if err != nil {
		result.Error = fmt.Sprintf("%s: %v", domain.ErrURLFetchFailed.Error(), err)
		result.InaccessibleLinkCount = 1
		// Attempt to save even if fetch failed, to log the error
		saveErr := s.crawlResultRepo.Save(result)
		if saveErr != nil {
			fmt.Printf("Error saving crawl result (fetch failed): %v\n", saveErr)
		}
		return result, domain.ErrURLFetchFailed
	}
	defer res.Body.Close()

	if res.StatusCode >= 400 {
		result.Error = fmt.Sprintf("URL returned status code: %d", res.StatusCode)
		result.InaccessibleLinkCount = 1
		// Attempt to save even if status code is error, to log the error
		saveErr := s.crawlResultRepo.Save(result)
		if saveErr != nil {
			fmt.Printf("Error saving crawl result (status code error): %v\n", saveErr)
		}
		return result, nil // Return result with error message, but no Go error for 4xx/5xx status
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		result.Error = fmt.Sprintf("%s: %v", domain.ErrHTMLParseFailed.Error(), err)
		// Attempt to save even if HTML parse failed, to log the error
		saveErr := s.crawlResultRepo.Save(result)
		if saveErr != nil {
			fmt.Printf("Error saving crawl result (HTML parse failed): %v\n", saveErr)
		}
		return result, domain.ErrHTMLParseFailed
	}

	// HTML Version (basic inference from doctype)
	if strings.Contains(doc.Text(), "<!DOCTYPE html>") {
		result.HTMLVersion = "HTML5"
	} else if strings.Contains(doc.Text(), `<!DOCTYPE HTML PUBLIC "-//W3C//DTD HTML 4.01//EN" "http://www.w3.org/TR/html4/strict.dtd">`) {
		result.HTMLVersion = "HTML 4.01 Strict"
	} else if strings.Contains(doc.Text(), `<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Strict//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-strict.dtd">`) {
		result.HTMLVersion = "XHTML 1.0 Strict"
	} else {
		result.HTMLVersion = "Unknown/Other"
	}

	// Page Title
	result.PageTitle = doc.Find("title").Text()

	// Heading Tag Counts
	for i := 1; i <= 6; i++ {
		tag := fmt.Sprintf("h%d", i)
		result.HeadingCounts[tag] = doc.Find(tag).Length()
	}

	// Link Counts and Login Form
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if !exists || href == "" || href == "#" {
			return
		}

		linkURL, err := parsedURL.Parse(href)
		if err != nil {
			return // Malformed URL, skip
		}

		if linkURL.Host == parsedURL.Host {
			result.InternalLinkCount++
		} else {
			result.ExternalLinkCount++
		}
	})

	// Check for login form
	doc.Find("form").Each(func(i int, s *goquery.Selection) {
		if s.Find("input[type=password]").Length() > 0 {
			result.HasLoginForm = true
			return // Found a login form, no need to check further
		}
	})

	// Save the successful crawl result to the database
	saveErr := s.crawlResultRepo.Save(result)
	if saveErr != nil {
		fmt.Printf("Error saving successful crawl result: %v\n", saveErr)
		// Decide how to handle this error: return it, log it, etc.
		// For now, we'll just log and proceed with returning the result
	}

	return result, nil
}

// GetCrawlResults retrieves a paginated list of crawl results
func (s *CrawlService) GetCrawlResults(query queries.GetCrawlResultsQuery) ([]domain.CrawlResult, error) {
	// Ensure page and pageSize are valid
	if query.Page < 1 {
		query.Page = 1
	}
	if query.PageSize < 1 {
		query.PageSize = 10 // Default page size
	}

	results, err := s.crawlResultRepo.GetAll(query.Page, query.PageSize)
	if err != nil {
		return nil, fmt.Errorf("failed to get crawl results from repository: %w", err)
	}

	return results, nil
}

// DeleteCrawlResult deletes a crawl result by ID
func (s *CrawlService) DeleteCrawlResult(cmd commands.DeleteCrawlResultCommand) error {
	err := s.crawlResultRepo.Delete(cmd.ID)
	if err != nil {
		return fmt.Errorf("failed to delete crawl result: %w", err)
	}
	return nil
}
