package export

import (
	"encoding/json"
	"fmt"
	"strings"
)

type ReviewData struct {
	Title    string        `json:"title"`
	PR       string        `json:"pull_request"`
	Summary  string        `json:"summary"`
	Issues   []interface{} `json:"issues"`
	Score    float64       `json:"score"`
}

func ExportMarkdown(data ReviewData) string {
	var sb strings.Builder
	
	sb.WriteString(fmt.Sprintf("# Code Review: %s\n\n", data.Title))
	sb.WriteString(fmt.Sprintf("**PR:** %s\n\n", data.PR))
	sb.WriteString(fmt.Sprintf("**Score:** %.1f/100\n\n", data.Score))
	sb.WriteString("---\n\n")
	
	sb.WriteString("## Summary\n\n")
	sb.WriteString(data.Summary + "\n\n")
	
	sb.WriteString("## Issues\n\n")
	for i, issue := range data.Issues {
		if m, ok := issue.(map[string]interface{}); ok {
			severity := m["severity"].(string)
			message := m["message"].(string)
			sb.WriteString(fmt.Sprintf("%d. **%s** - %s\n", i+1, strings.ToUpper(severity), message))
		}
	}
	
	return sb.String()
}

func ExportJSON(data ReviewData) string {
	bytes, _ := json.MarshalIndent(data, "", "  ")
	return string(bytes)
}

func ExportText(data ReviewData) string {
	var sb strings.Builder
	
	sb.WriteString(fmt.Sprintf("Code Review: %s\n", data.Title))
	sb.WriteString(fmt.Sprintf("PR: %s\n", data.PR))
	sb.WriteString(fmt.Sprintf("Score: %.1f/100\n\n", data.Score))
	sb.WriteString("Summary:\n")
	sb.WriteString(data.Summary + "\n\n")
	
	sb.WriteString("Issues:\n")
	for _, issue := range data.Issues {
		if m, ok := issue.(map[string]interface{}); ok {
			sb.WriteString(fmt.Sprintf("- %s\n", m["message"]))
		}
	}
	
	return sb.String()
}
