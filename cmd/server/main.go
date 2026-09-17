package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"ai-code-review/internal/api"
	"ai-code-review/internal/github"
	"ai-code-review/internal/review"
	"ai-code-review/internal/store"
)

func main() {
	githubToken := os.Getenv("GITHUB_TOKEN")
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "reviews.db"
	}
	
	store, err := store.NewStore(dbPath)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer store.Close()
	
	reviewer := review.NewClient()
	ghClient := github.NewClient(githubToken)
	
	r := gin.Default()
	
	// CORS
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
			return
		}
		c.Next()
	})
	
	api.SetupRoutes(r, store, reviewer, ghClient)
	
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}
	
	log.Printf("Server starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
