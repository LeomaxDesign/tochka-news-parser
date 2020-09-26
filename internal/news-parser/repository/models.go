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
	ID             int
	URL            string `json:"url"`
	Title          string `json:"title"`
	Type           int    `json:"type"`
	Frequency      int    `json:"frequency"`
	ParseCount     int    `json:"parse_count"`
	ItemTag        string `json:"item_tag"`
	TitleTag       string `json:"title_tag"`
	DescriptionTag string `json:"description_tag"`
	LinkTag        string `json:"link_tag"`
	PublishedTag   string `json:"published_tag"`
	ImgTag         string `json:"img_tag"`
}

// NewsModel ...
type NewsModel struct {
	ID          int       `json:"id"`
	FeedID      int       `json:"-"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Link        string    `json:"link"`
	Img         string    `json:"img"`
	Published   time.Time `json:"published"`
	Parsed      time.Time `json:"-"`
}
