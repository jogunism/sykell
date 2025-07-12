package services

import (
	"bytes"
	"database/sql"
	"fmt"
	"io/ioutil"
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

// saveFailedCrawl handles the logic for saving a crawl result when an error occurs.
func (s *CrawlService) saveFailedCrawl(result domain.CrawlResult, crawlError error, format string, args ...any) (domain.CrawlResult, int, error) {
	result.Error = fmt.Sprintf(format, args...)
	result.InaccessibleLinkCount = 0
	result.InternalLinkCount = 0
	result.ExternalLinkCount = 0
	result.PageTitle = "CRAWLING HAS FAILED"

	id, saveErr := s.crawlResultRepo.Save(result)
	if saveErr != nil {
		fmt.Printf("Error saving failed crawl result: %v (original error: %v)\n", saveErr, crawlError)
	}
	return result, id, crawlError
}

func (s *CrawlService) Crawl(cmd commands.CrawlCommand) (domain.CrawlResult, int, error) {
	result := domain.CrawlResult{
		URL:           domain.NullString{NullString: sql.NullString{String: cmd.URL, Valid: true}},
		HeadingCounts: make(map[string]int),
	}

	parsedURL, err := url.Parse(cmd.URL)
	if err != nil {
		// No need to save, as the URL is invalid from the start.
		result.Error = domain.ErrInvalidURLFormat.Error()
		return result, 0, domain.ErrInvalidURLFormat
	}

	res, err := http.Get(cmd.URL)
	if err != nil {
		return s.saveFailedCrawl(result, domain.ErrURLFetchFailed, "%s: %v", domain.ErrURLFetchFailed.Error(), err)
	}
	defer res.Body.Close()

	if res.StatusCode >= 400 {
		return s.saveFailedCrawl(result, domain.ErrURLFetchFailed, "URL returned status code: %d", res.StatusCode)
	}

	bodyBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		wrappedErr := fmt.Errorf("failed to read response body: %w", err)
		return s.saveFailedCrawl(result, wrappedErr, "Failed to read response body: %v", err)
	}

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(bodyBytes))
	if err != nil {
		return s.saveFailedCrawl(result, domain.ErrHTMLParseFailed, "%s: %v", domain.ErrHTMLParseFailed.Error(), err)
	}

	// HTML Version
	bodyString := string(bodyBytes)
	if strings.Contains(strings.ToLower(bodyString), "<!doctype html>") {
		result.HTMLVersion = "HTML5"
	} else if strings.Contains(bodyString, `<!DOCTYPE HTML PUBLIC "-//W3C//DTD HTML 4.01//EN" "http://www.w3.org/TR/html4/strict.dtd">`) {
		result.HTMLVersion = "HTML 4.01 Strict"
	} else if strings.Contains(bodyString, `<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Strict//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-strict.dtd">`) {
		result.HTMLVersion = "XHTML 1.0 Strict"
	} else {
		result.HTMLVersion = "Unknown/Other"
	}

	// Page Title
	result.PageTitle = doc.Find("title").Text()
	if result.PageTitle == "" {
		result.PageTitle = "NO TITLE"
	}

	// Count heading tags
	for i := 1; i <= 6; i++ {
		heading := fmt.Sprintf("h%d", i)
		count := doc.Find(heading).Length()
		if count > 0 {
			result.HeadingCounts[heading] = count
		}
	}

	// Count links
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		link, exists := s.Attr("href")
		if !exists {
			return
		}
		resolvedLink, err := parsedURL.Parse(link)
		if err != nil {
			result.InaccessibleLinkCount++
			return
		}
		if resolvedLink.Host == parsedURL.Host {
			result.InternalLinkCount++
		} else {
			result.ExternalLinkCount++
		}
	})

	// Save the successful result
	id, saveErr := s.crawlResultRepo.Save(result)
	if saveErr != nil {
		result.Error = fmt.Sprintf("failed to save crawl result: %v", saveErr)
		return result, 0, saveErr
	}

	return result, id, nil
}

type GetCrawlResultsResponse struct {
	List       []domain.CrawlResult `json:"list"`
	TotalCount int                  `json:"total_count"`
}

// GetCrawlResults retrieves a paginated list of crawl results
func (s *CrawlService) GetCrawlResults(query queries.GetCrawlResultsQuery) (GetCrawlResultsResponse, error) {
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
		List:       results,
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
