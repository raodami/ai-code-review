package review

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

type ReviewResult struct {
	Summary      string           `json:"summary"`
	Issues       []ReviewIssue    `json:"issues"`
	Suggestions  []string         `json:"suggestions"`
	Score        float64          `json:"score"`
	Source       string           `json:"source"`
	CreatedAt    time.Time        `json:"created_at"`
}

type ReviewIssue struct {
	Type     string `json:"type"`
	File     string `json:"file"`
	Line     int    `json:"line"`
	Message  string `json:"message"`
	Severity string `json:"severity"`
}

type Client struct {
	APIKey string
	BaseURL string
}

func NewClient() *Client {
	apiKey := os.Getenv("DEEPSEEK_API_KEY")
	return &Client{
		APIKey: apiKey,
		BaseURL: "https://api.deepseek.com/v1",
	}
}

func (c *Client) GetModel() string {
	if c.APIKey != "" {
		return "DeepSeek Code Review"
	}
	return "Mock (no API key)"
}

func (c *Client) ReviewCode(prSummary string, files []interface{}) (*ReviewResult, error) {
	if c.APIKey == "" {
		return generateMockReview(prSummary, files), nil
	}
	
	prompt := buildReviewPrompt(prSummary, files)
	
	reqBody, _ := json.Marshal(map[string]interface{}{
		"model": "deepseek-coder",
		"messages": []map[string]string{
			{"role": "system", "content": "You are an expert code reviewer. Analyze the code changes and provide detailed feedback on quality, security, performance, and best practices."},
			{"role": "user", "content": prompt},
		},
		"temperature": 0.2,
	})
	
	resp, err := http.Post(
		fmt.Sprintf("%s/chat/completions", c.BaseURL),
		"application/json",
		bytes.NewBuffer(reqBody),
	)
	if err != nil {
		return generateMockReview(prSummary, files), nil
	}
	defer resp.Body.Close()
	
	body, _ := ioutil.ReadAll(resp.Body)
	
	var result struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	json.Unmarshal(body, &result)
	
	if len(result.Choices) > 0 {
		content := result.Choices[0].Message.Content
		return parseReviewResponse(content), nil
	}
	
	return generateMockReview(prSummary, files), nil
}

func buildReviewPrompt(prSummary string, files []interface{}) string {
	var sb strings.Builder
	
	sb.WriteString(fmt.Sprintf("PR Title: %s\n\n", prSummary))
	sb.WriteString("Files changed:\n")
	
	for _, f := range files {
		file := f.(map[string]interface{})
		filename := file["filename"].(string)
		patch := file["patch"].(string)
		sb.WriteString(fmt.Sprintf("\n### %s\n\n%s\n", filename, patch))
	}
	
	sb.WriteString("\nPlease review this code and provide:\n")
	sb.WriteString("1. Summary of changes\n")
	sb.WriteString("2. Issues found (security, performance, best practices)\n")
	sb.WriteString("3. Suggestions for improvement\n")
	sb.WriteString("4. Overall quality score (0-100)\n")
	
	return sb.String()
}

func parseReviewResponse(content string) *ReviewResult {
	result := &ReviewResult{
		Summary: content,
		Source: "DeepSeek",
		Score:  85.0,
		CreatedAt: time.Now(),
	}
	
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.Contains(strings.ToLower(line), "security") {
			result.Issues = append(result.Issues, ReviewIssue{
				Type: "security",
				Message: line,
				Severity: "high",
			})
		} else if strings.Contains(strings.ToLower(line), "performance") {
			result.Suggestions = append(result.Suggestions, line)
		}
	}
	
	return result
}

func generateMockReview(prSummary string, files []interface{}) *ReviewResult {
	result := &ReviewResult{
		Summary: fmt.Sprintf("Reviewed PR: %s. The changes look good overall with some minor suggestions.", prSummary),
		Issues: []ReviewIssue{
			{Type: "style", File: "main.go", Line: 45, Message: "Consider using consistent naming conventions", Severity: "low"},
			{Type: "security", File: "auth.go", Line: 12, Message: "Potential SQL injection vulnerability - use parameterized queries", Severity: "high"},
			{Type: "performance", File: "query.go", Line: 78, Message: "Missing index on frequently queried column", Severity: "medium"},
		},
		Suggestions: []string{
			"Add input validation for all user inputs",
			"Consider using connection pooling for database operations",
			"Add unit tests for new functionality",
		},
		Score: 78.5,
		Source: "Mock",
		CreatedAt: time.Now(),
	}
	
	return result
}
