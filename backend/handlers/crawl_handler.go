package handlers

import (
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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cmd := commands.CrawlCommand{
		URL: req.URL,
	}

	result, err := h.crawlService.Crawl(cmd)
	if err != nil {
		// Handle specific errors from the service layer
		if err == domain.ErrInvalidURLFormat {
			c.JSON(http.StatusBadRequest, gin.H{"error": result.Error}) // Changed to 400 Bad Request
			return
		} else if err == domain.ErrURLFetchFailed || err == domain.ErrHTMLParseFailed {
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error}) // Use result.Error for detailed message
			return
		}
		// Generic error for unexpected issues
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process crawl request"})
		return
	}

	c.JSON(http.StatusOK, result)
}

// GetCrawlResults handles the request to get a paginated list of crawl results
func (h *CrawlHandler) GetCrawlResults(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("pageSize", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 {
		pageSize = 10
	}

	query := queries.GetCrawlResultsQuery{
		Page:     page,
		PageSize: pageSize,
	}

	results, err := h.crawlService.GetCrawlResults(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, results)
}

// DeleteCrawlResult handles the request to delete a crawl result by ID
func (h *CrawlHandler) DeleteCrawlResult(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	cmd := commands.DeleteCrawlResultCommand{
		ID: id,
	}

	if err := h.crawlService.DeleteCrawlResult(cmd); err != nil {
		// Handle specific errors from the service layer
		if strings.Contains(err.Error(), "no record found") {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete crawl result"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Crawl result deleted successfully"})
}
