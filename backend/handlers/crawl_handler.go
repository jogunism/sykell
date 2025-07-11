package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"backend/application/commands"
	"backend/application/queries"
	"backend/application/services"
	"backend/domain"

	"github.com/gin-gonic/gin"
)

// CrawlHandler handles HTTP requests related to URL crawling
type CrawlHandler struct {
	crawlService *services.CrawlService
}

// NewCrawlHandler creates a new CrawlHandler
func NewCrawlHandler(crawlService *services.CrawlService) *CrawlHandler {
	return &CrawlHandler{
		crawlService: crawlService,
	}
}

// Crawl handles the URL crawling request
func (h *CrawlHandler) Crawl(c *gin.Context) {
	var req domain.CrawlRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	cmd := commands.CrawlCommand{
		URL: req.URL,
	}

	result, id, err := h.crawlService.Crawl(cmd)
	if err != nil {
		fmt.Println(">>>>> ", result.Error)
		switch err {
			case domain.ErrInvalidURLFormat:
				c.JSON(http.StatusBadRequest, gin.H{"message": result.Error})
				return
			case domain.ErrURLFetchFailed, domain.ErrHTMLParseFailed:
				c.JSON(http.StatusInternalServerError, gin.H{"message": result.Error})
				return
		}

		// Generic error for unexpected issues
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to process crawl request"})
		
		return
	}

	// fmt.Println(result)
	c.JSON(http.StatusOK, gin.H{"message": "Crawl the url successfully", "id": id})
}

// GetCrawlResults handles the request to get a paginated list of crawl results
func (h *CrawlHandler) GetCrawlResults(c *gin.Context) {
	currPageStr := c.DefaultQuery("currPage", "1")
	pageSizeStr := c.DefaultQuery("pageSize", "10")
	queryStr := c.DefaultQuery("query", "")
	sortingJsonStr := c.DefaultQuery("sorting", "") // Get sorting JSON string

	currPage, err := strconv.Atoi(currPageStr)
	if err != nil || currPage < 1 {
		currPage = 1
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 {
		pageSize = 10
	}

	query := queries.GetCrawlResultsQuery{
		CurrPage:   currPage,
		PageSize:   pageSize,
		Query:      queryStr,
		SortingJson: sortingJsonStr, // Assign sorting JSON string
	}

	response, err := h.crawlService.GetCrawlResults(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"list": response.List, "total_count": response.TotalCount})
}

// DeleteCrawlResults handles the request to delete multiple crawl results by IDs
func (h *CrawlHandler) DeleteCrawlResults(c *gin.Context) {
	var req domain.DeleteCrawlResultsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	cmd := commands.DeleteCrawlResultsCommand{
		IDs: req.IDs,
	}

	if err := h.crawlService.DeleteCrawlResults(cmd); err != nil {
		// Handle specific errors from the service layer
		if strings.Contains(err.Error(), "no records found") {
			c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to delete multiple crawl results"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Crawl results deleted successfully"})
}
