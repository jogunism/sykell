package main

import (
	"log"
	"strings"

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


	// gin.SetMode(gin.ReleaseMode)


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

	// X-Forwarded-For 헤더 디버깅용 미들웨어
	r.Use(func(c *gin.Context) {
		xForwardedFor := c.GetHeader("X-Forwarded-For")
		if xForwardedFor != "" {
			ips := strings.Split(xForwardedFor, ",")
			log.Printf("X-Forwarded-For: %s (IP count: %d)", xForwardedFor, len(ips))
			
			// 30개 이상이면 경고
			if len(ips) >= 30 {
				log.Printf("WARNING: X-Forwarded-For has %d IPs, this will cause 463 error", len(ips))
			}
		}
		
		// 모든 X-Forwarded-* 헤더 로깅
		for name, values := range c.Request.Header {
			if strings.HasPrefix(name, "X-Forwarded-") {
				log.Printf("Header %s: %v", name, values)
			}
		}
		
		c.Next()
	})


	// AWS ALB 463 에러 방지를 위한 프록시 설정
	r.ForwardedByClientIP = true
	r.SetTrustedProxies([]string{
		"10.0.0.0/8",     // VPC 내부
		"172.16.0.0/12",  // ALB 대역
	})

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
		protected.POST("/crawl", crawlHandler.Crawl)
		protected.GET("/crawl/list", crawlHandler.GetCrawlResults)
		protected.DELETE("/crawl", crawlHandler.DeleteCrawlResults)
	}

	r.Run(":8080") // listen and serve on 0.0.0.0:8080
}
