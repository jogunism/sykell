package services

import (
	"database/sql"
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

func (s *CrawlService) Crawl(cmd commands.CrawlCommand) (domain.CrawlResult, int, error) {
	result := domain.CrawlResult{
		URL:           domain.NullString{NullString: sql.NullString{String: cmd.URL, Valid: true}},
		HeadingCounts: make(map[string]int),
	}

	parsedURL, err := url.Parse(cmd.URL)
	if err != nil {
		result.Error = domain.ErrInvalidURLFormat.Error()
		return result, 0, domain.ErrInvalidURLFormat
	}

	res, err := http.Get(cmd.URL)
	if err != nil {
		result.Error = fmt.Sprintf("%s: %v", domain.ErrURLFetchFailed.Error(), err)
		result.InaccessibleLinkCount = 0 // Set to 0 on error
		result.InternalLinkCount = 0    // Set to 0 on error
		result.ExternalLinkCount = 0    // Set to 0 on error
		id, saveErr := s.crawlResultRepo.Save(result)
		if saveErr != nil {
			fmt.Printf("Error saving crawl result (fetch failed): %v\n", saveErr)
		}
		return result, id, domain.ErrURLFetchFailed
	}
	defer res.Body.Close()

	if res.StatusCode >= 400 {
		result.Error = fmt.Sprintf("URL returned status code: %d", res.StatusCode)
		result.InaccessibleLinkCount = 0 // Set to 0 on error
		result.InternalLinkCount = 0    // Set to 0 on error
		result.ExternalLinkCount = 0    // Set to 0 on error
		id, saveErr := s.crawlResultRepo.Save(result)
		if saveErr != nil {
			fmt.Printf("Error saving crawl result (status code error): %v\n", saveErr)
		}
		return result, id, domain.ErrURLFetchFailed // Return result with error message and a Go error
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		result.Error = fmt.Sprintf("%s: %v", domain.ErrHTMLParseFailed.Error(), err)
		result.InaccessibleLinkCount = 0 // Set to 0 on error
		result.InternalLinkCount = 0    // Set to 0 on error
		result.ExternalLinkCount = 0    // Set to 0 on error
		id, saveErr := s.crawlResultRepo.Save(result)
		if saveErr != nil {
			fmt.Printf("Error saving crawl result (HTML parse failed): %v\n", saveErr)
		}
		return result, id, domain.ErrHTMLParseFailed
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
	id, saveErr := s.crawlResultRepo.Save(result)
	if saveErr != nil {
		fmt.Printf("Error saving successful crawl result: %v\n", saveErr)
		// Decide how to handle this error: return it, log it, etc.
		// For now, we'll just log and proceed with returning the result
	}

	return result, id, nil
}

type GetCrawlResultsResponse struct {
	List      []domain.CrawlResult `json:"list"`
	TotalCount int                `json:"total_count"`
}

// GetCrawlResults retrieves a paginated list of crawl results
func (s *CrawlService) GetCrawlResults(query queries.GetCrawlResultsQuery) (GetCrawlResultsResponse, error) {
	// Ensure page and pageSize are valid
	if query.CurrPage < 1 {
		query.CurrPage = 1
	}
	if query.PageSize < 1 {
		query.PageSize = 10
	}

	results, totalCount, err := s.crawlResultRepo.GetAll(query.CurrPage, query.PageSize, query.Query, query.SortingJson)
	if err != nil {
		return GetCrawlResultsResponse{}, fmt.Errorf("failed to get crawl results from repository: %w", err)
	}

	return GetCrawlResultsResponse{
		List:      results,
		TotalCount: totalCount,
	}, nil
}

// DeleteCrawlResults deletes multiple crawl results by IDs
func (s *CrawlService) DeleteCrawlResults(cmd commands.DeleteCrawlResultsCommand) error {
	if len(cmd.IDs) == 0 {
		return nil // Nothing to delete
	}
	err := s.crawlResultRepo.DeleteMany(cmd.IDs)
	if err != nil {
		return fmt.Errorf("failed to delete multiple crawl results: %w", err)
	}
	return nil
}
