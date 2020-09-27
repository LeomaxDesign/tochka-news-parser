package repository

import (
	"time"
)

type NewsRepository interface {
	GetAll(searchString string) ([]*News, error)
	Add(news *News) error
	IsExists(news *News) (bool, error)
}

type News struct {
	ID          int       `json:"id"`
	NewsFeedID  int       `json:"-"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Link        string    `json:"link"`
	Img         string    `json:"img"`
	Published   time.Time `json:"published"`
	Parsed      time.Time `json:"-"`
}
