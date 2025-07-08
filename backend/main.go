package main

import (
	"log"

	"backend/application/services"
	"backend/handlers"
	"backend/infrastructure/database"
	"backend/infrastructure/persistence"

	"github.com/gin-gonic/gin"
)

func main() {
	// Database connection string (replace with your actual credentials)
	dbConnStr := "admin:HyunwooCho!23$@tcp(sykell.c10yg6egqxbv.eu-central-1.rds.amazonaws.com:3306)/sykell?parseTime=true"

	// Initialize database connection
	db, err := database.InitDB(dbConnStr)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Initialize repositories
	crawlResultRepo := persistence.NewMySQLCrawlResultRepository(db)

	// Initialize services with their dependencies
	testService := services.NewTestService()
	crawlService := services.NewCrawlService(crawlResultRepo)

	// Initialize handlers with their respective services
	testHandler := handlers.NewTestHandler(testService)
	crawlHandler := handlers.NewCrawlHandler(crawlService)

	// Setup Gin router
	r := gin.Default()

	// Public endpoints
	r.GET("/test", testHandler.GetTestMessage)

	// Protected endpoints (require JWT authentication)
	protected := r.Group("/")
	protected.Use(handlers.AuthMiddleware())
	{
		protected.POST("/crawl", crawlHandler.Crawl)
		protected.GET("/crawl/list", crawlHandler.GetCrawlResults)
		protected.DELETE("/crawl/:id", crawlHandler.DeleteCrawlResult) // New endpoint for deleting crawl results
	}

	r.Run(":8080") // listen and serve on 0.0.0.0:8080
}
