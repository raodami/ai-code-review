package github

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type PRReview struct {
	Owner       string         `json:"owner"`
	Repo        string         `json:"repo"`
	PRNumber    int            `json:"pr_number"`
	Title       string         `json:"title"`
	Description string         `json:"description"`
	Files       []FileChange   `json:"files"`
	Comments    []Comment      `json:"comments"`
	Status      string         `json:"status"`
	URL         string         `json:"url"`
}

type FileChange struct {
	Filename string      `json:"filename"`
	Status   string      `json:"status"`
	Additions int      `json:"additions"`
	Deletions int      `json:"deletions"`
	Patch    string     `json:"patch"`
}

type Comment struct {
	ID       int    `json:"id"`
	User     string `json:"user"`
	Body     string `json:"body"`
	Path     string `json:"path"`
	Line     int    `json:"line"`
}

type Client struct {
	Token string
	BaseURL string
}

func NewClient(token string) *Client {
	return &Client{
		Token: token,
		BaseURL: "https://api.github.com",
	}
}

func (c *Client) GetPullRequest(owner, repo string, prNumber int) (*PRReview, error) {
	url := fmt.Sprintf("%s/repos/%s/%s/pulls/%d", c.BaseURL, owner, repo, prNumber)
	
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("token %s", c.Token))
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("GitHub API error: %s", string(body))
	}
	
	var pr struct {
		Title       string `json:"title"`
		Body        string `json:"body"`
		URL         string `json:"url"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&pr); err != nil {
		return nil, err
	}
	
	// Get files
	filesURL := fmt.Sprintf("%s/repos/%s/%s/pulls/%d/files", c.BaseURL, owner, repo, prNumber)
	filesReq, _ := http.NewRequest("GET", filesURL, nil)
	filesReq.Header.Set("Authorization", fmt.Sprintf("token %s", c.Token))
	filesReq.Header.Set("Accept", "application/vnd.github.v3+json")
	
	filesResp, err := client.Do(filesReq)
	if err != nil {
		return nil, err
	}
	defer filesResp.Body.Close()
	
	var fileChanges []FileChange
	if err := json.NewDecoder(filesResp.Body).Decode(&fileChanges); err != nil {
		return nil, err
	}
	
	return &PRReview{
		Owner:       owner,
		Repo:        repo,
		PRNumber:    prNumber,
		Title:       pr.Title,
		Description: pr.Body,
		Files:       fileChanges,
		URL:         pr.URL,
		Status:      "open",
	}, nil
}

func extractOwnerRepo(url string) (string, string, int, error) {
	parts := strings.Split(strings.TrimSuffix(url, "/"), "/")
	if len(parts) < 5 {
		return "", "", 0, fmt.Errorf("invalid PR URL")
	}
	
	owner := parts[len(parts)-3]
	repo := parts[len(parts)-2]
	prStr := parts[len(parts)-1]
	
	var prNumber int
	fmt.Sscanf(prStr, "pull/%d", &prNumber)
	
	return owner, repo, prNumber, nil
}
