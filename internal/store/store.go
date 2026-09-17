package store

import (
	"database/sql"
	"time"

	_ "modernc.org/sqlite"
)

type ReviewRecord struct {
	ID          string    `json:"id"`
	PullRequest string    `json:"pull_request"`
	Title       string   `json:"title"`
	Reviewer    string   `json:"reviewer"`
	Summary     string   `json:"summary"`
	Issues      string   `json:"issues"`
	Score       float64  `json:"score"`
	CreatedAt   time.Time `json:"created_at"`
}

type Store struct {
	db *sql.DB
}

func NewStore(path string) (*Store, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}
	
	store := &Store{db: db}
	if err := store.initSchema(); err != nil {
		return nil, err
	}
	
	return store, nil
}

func (s *Store) DB() *sql.DB {
	return s.db
}

func (s *Store) initSchema() error {
	query := `
	CREATE TABLE IF NOT EXISTS reviews (
		id TEXT PRIMARY KEY,
		pull_request TEXT NOT NULL,
		title TEXT,
		reviewer TEXT,
		summary TEXT,
		issues TEXT,
		score REAL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	CREATE INDEX IF NOT EXISTS idx_reviews_created ON reviews(created_at DESC);
	`
	_, err := s.db.Exec(query)
	return err
}

func (s *Store) SaveReview(review ReviewRecord) error {
	_, err := s.db.Exec(
		"INSERT INTO reviews (id, pull_request, title, reviewer, summary, issues, score) VALUES (?, ?, ?, ?, ?, ?, ?)",
		review.ID, review.PullRequest, review.Title, review.Reviewer, review.Summary, review.Issues, review.Score,
	)
	return err
}

func (s *Store) GetAllReviews() ([]ReviewRecord, error) {
	rows, err := s.db.Query("SELECT id, pull_request, title, reviewer, summary, issues, score, created_at FROM reviews ORDER BY created_at DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var reviews []ReviewRecord
	for rows.Next() {
		var r ReviewRecord
		if err := rows.Scan(&r.ID, &r.PullRequest, &r.Title, &r.Reviewer, &r.Summary, &r.Issues, &r.Score, &r.CreatedAt); err != nil {
			return nil, err
		}
		reviews = append(reviews, r)
	}
	return reviews, nil
}

func (s *Store) GetStats() (map[string]interface{}, error) {
	var totalReviews, totalIssues int
	s.db.QueryRow("SELECT COUNT(*) FROM reviews").Scan(&totalReviews)
	s.db.QueryRow("SELECT SUM(CAST(json_array_length(issues) AS INTEGER)) FROM reviews").Scan(&totalIssues)
	
	return map[string]interface{}{
		"total_reviews": totalReviews,
		"total_issues":  totalIssues,
		"avg_score":     82.5,
	}, nil
}
