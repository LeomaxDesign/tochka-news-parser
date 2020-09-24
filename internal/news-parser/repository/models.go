package repository

import (
	"database/sql"
	"time"
)

// PostgresClient ...
type PostgresClient struct {
	host     string
	username string
	password string
	dbname   string
	port     int
	db       *sql.DB
}

// NewsFeedModel ...
type NewsFeedModel struct {
	ID         int
	URL        string `json:"url"`
	Name       string
	Type       int `json:"type"`
	Frequency  int `json:"frequency"`
	ParseCount int `json:"parse_count"`
}

// NewsModel ...
type NewsModel struct {
	ID          int
	FeedID      int
	Title       string
	Description string
	Link        string
	Img         string
	Published   *time.Time
	Parsed      time.Time
}
