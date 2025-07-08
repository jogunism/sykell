package main

import (
	"backend/application/services"
	"backend/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// Initialize services
	testService := services.NewTestService()
	crawlService := services.NewCrawlService()

	// Initialize handlers with their respective services
	testHandler := handlers.NewTestHandler(testService)
	crawlHandler := handlers.NewCrawlHandler(crawlService)

	// Public endpoints
	r.GET("/test", testHandler.GetTestMessage)

	// Protected endpoints (require JWT authentication)
	protected := r.Group("/")
	protected.Use(handlers.AuthMiddleware())
	{
		protected.POST("/crawl", crawlHandler.Crawl)
	}

	r.Run(":8080") // listen and serve on 0.0.0.0:8080
}
