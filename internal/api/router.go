package api

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"ai-code-review/internal/export"
	"ai-code-review/internal/github"
	"ai-code-review/internal/review"
	"ai-code-review/internal/store"
)

func SetupRoutes(r *gin.Engine, s *store.Store, reviewer *review.Client, ghClient *github.Client) {
	r.POST("/api/review", func(c *gin.Context) {
		var req struct {
			PRURL string `json:"pr_url" binding:"required"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}
		
		// Extract PR info from URL
		parts := strings.Split(req.PRURL, "/")
		if len(parts) < 6 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid PR URL format"})
			return
		}
		
		owner := parts[3]
		repo := parts[4]
		var prNumber int
		fmt.Sscanf(parts[5], "pull/%d", &prNumber)
		
		if prNumber == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid PR number"})
			return
		}
		
		// Fetch PR from GitHub
		pr, err := ghClient.GetPullRequest(owner, repo, prNumber)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to fetch PR: %v", err)})
			return
		}
		
		// Generate review
		result, err := reviewer.ReviewCode(pr.Title, pr.Files)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate review"})
			return
		}
		
		// Save to database
		reviewRecord := store.ReviewRecord{
			ID:          uuid.New().String(),
			PullRequest: pr.URL,
			Title:       pr.Title,
			Reviewer:    result.Source,
			Summary:     result.Summary,
			Issues:      toJSON(result.Issues),
			Score:       result.Score,
			CreatedAt:   time.Now(),
		}
		
		if err := s.SaveReview(reviewRecord); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save review"})
			return
		}
		
		c.JSON(http.StatusOK, gin.H{
			"id": reviewRecord.ID,
			"result": result,
		})
	})
	
	r.GET("/api/review/:id", func(c *gin.Context) {
		id := c.Param("id")
		// TODO: implement GetReviewByID
		c.JSON(http.StatusOK, gin.H{"id": id})
	})
	
	r.GET("/api/history", func(c *gin.Context) {
		reviews, err := s.GetAllReviews()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch history"})
			return
		}
		c.JSON(http.StatusOK, reviews)
	})
	
	r.GET("/api/stats", func(c *gin.Context) {
		stats, err := s.GetStats()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get stats"})
			return
		}
		c.JSON(http.StatusOK, stats)
	})
	
	r.GET("/api/export/:id/:format", func(c *gin.Context) {
		id := c.Param("id")
		format := c.Param("format")
		// TODO: implement export with actual review data
		data := export.ReviewData{
			Title: "Sample Review",
			Score: 85.0,
		}
		
		var content string
		switch format {
		case "markdown":
			content = export.ExportMarkdown(data)
		case "json":
			content = export.ExportJSON(data)
		default:
			content = export.ExportText(data)
		}
		
		c.Header("Content-Disposition", fmt.Sprintf(`attachment; filename="review.%s"`, format))
		c.Data(http.StatusOK, "text/plain", []byte(content))
	})
}

func toJSON(v interface{}) string {
	b, _ := json.Marshal(v)
	return string(b)
}
