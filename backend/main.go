package main

import (
	"log"
	"os" // environment variables

	"backend/application/services"
	"backend/handlers"
	"backend/infrastructure/database"
	"backend/infrastructure/persistence"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	dbConnStr := os.Getenv("DATABASE_URL")
	if dbConnStr == "" {
		log.Fatal("DATABASE_URL environment variable not set")
	}

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

	// Configure CORS middleware
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"}, // Allow your frontend origin
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		// MaxAge: 12 * time.Hour, // Optional: How long the preflight request can be cached
	}))

	// Public endpoints
	r.GET("/test", testHandler.GetTestMessage)

	// Protected endpoints (require JWT authentication)
	protected := r.Group("/")
	protected.Use(handlers.AuthMiddleware())
	{
		protected.POST("/crawl", crawlHandler.Crawl)
		protected.GET("/crawl/list", crawlHandler.GetCrawlResults)
		protected.DELETE("/crawl", crawlHandler.DeleteCrawlResults)
	}

	r.Run(":8080") // listen and serve on 0.0.0.0:8080
}
