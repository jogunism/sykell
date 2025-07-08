package handlers

import (
	"net/http"

	"backend/application/queries"
	"backend/application/services"

	"github.com/gin-gonic/gin"
)

// TestHandler handles HTTP requests related to test endpoint
type TestHandler struct {
	testService *services.TestService
}

// NewTestAPIHandler creates a new TestAPIHandler
func NewTestHandler(testService *services.TestService) *TestHandler {
	return &TestHandler{
		testService: testService,
	}
}

// GetTestMessage handles the /test GET request
func (h *TestHandler) GetTestMessage(c *gin.Context) {
	query := queries.GetTestMessageQuery{}
	message := h.testService.GetTestMessage(query)
	c.JSON(http.StatusOK, gin.H{
		"message": message,
	})
}
