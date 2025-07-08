package handlers

import (
	"net/http"

	"backend/application/commands"
	"backend/application/services"
	"backend/domain"

	"github.com/gin-gonic/gin"
)

// CrawlHandler handles HTTP requests related to URL crawling
type CrawlHandler struct {
	crawlService *services.CrawlService
}

// NewCrawlAPIHandler creates a new CrawlAPIHandler
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
		if err == domain.ErrInvalidURLFormat || err == domain.ErrURLFetchFailed || err == domain.ErrHTMLParseFailed {
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error}) // Use result.Error for detailed message
			return
		}
		// Generic error for unexpected issues
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process crawl request"})
		return
	}

	c.JSON(http.StatusOK, result)
}