package services

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"backend/application/commands"
	"backend/domain"

	"github.com/PuerkitoBio/goquery" // crwal library
)

type CrawlService struct{}

func NewCrawlService() *CrawlService {
	return &CrawlService{}
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
		return result, domain.ErrURLFetchFailed
	}
	defer res.Body.Close()

	if res.StatusCode >= 400 {
		result.Error = fmt.Sprintf("URL returned status code: %d", res.StatusCode)
		result.InaccessibleLinkCount = 1
		return result, nil // Return result with error message, but no Go error for 4xx/5xx status
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		result.Error = fmt.Sprintf("%s: %v", domain.ErrHTMLParseFailed.Error(), err)
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

	return result, nil
}
