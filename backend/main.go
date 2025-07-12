package main

import (
	"log"

	"backend/application/services"
	"backend/handlers"
	"backend/infrastructure/database"
	"backend/infrastructure/persistence"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
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

	// Configure CORS middleware
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://sykell-alb-1699323201.eu-central-1.elb.amazonaws.com"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		// MaxAge: 12 * time.Hour, // Optional: How long the preflight request can be cached
	}))

	// Public endpoints
	api := r.Group("/api")
	api.GET("/test", testHandler.GetTestMessage)

	// Protected endpoints (require JWT authentication)
	protected := api.Group("")
	protected.Use(handlers.AuthMiddleware())
	{
		protected.GET("/crawl/list", crawlHandler.GetCrawlResults)
		protected.DELETE("/crawl", crawlHandler.DeleteCrawlResults)
		protected.POST("/crawl", crawlHandler.Crawl)
	}

	r.Run(":8080") // listen and serve on 0.0.0.0:8080
}
