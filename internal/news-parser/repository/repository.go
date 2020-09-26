package repository

import (
	"database/sql"
	"fmt"
	"strconv"
)

type PostgresClient struct {
	host     string
	username string
	password string
	dbname   string
	port     int
	DB       *sql.DB
}

type PostgresRepository struct {
	pc           *PostgresClient
	newsFeedRepo *NewsFeedRepository
	newsRepo     *NewsRepository
}

func New(host, username, password, dbname string, port int) *PostgresClient {
	return &PostgresClient{
		host:     host,
		username: username,
		password: password,
		dbname:   dbname,
		port:     port,
	}
}

func (pc *PostgresClient) Connect() error {
	var err error

	pc.DB, err = sql.Open("postgres", "host="+pc.host+" user="+pc.username+" dbname="+pc.dbname+" password="+pc.password+" port="+strconv.Itoa(pc.port)+" sslmode=disable")
	if err != nil {
		return fmt.Errorf("failed to open sql: %w", err)
	}

	if err = pc.DB.Ping(); err != nil {
		return fmt.Errorf("failed to ping sql: %w", err)
	}

	return nil
}

func (pc *PostgresClient) Disconnect() {
	pc.DB.Close()
}

// // CheckIfExistAndAddNews ...
// func (pc *PostgresClient) CheckIfExistAndAddNews(news []*NewsModel) (int, error) {
// 	var (
// 		err   error
// 		count int
// 	)
// 	query := `INSERT INTO feed_news (feed_id, title, description, link, published, parsed, img) VALUES ($1, $2, $3, $4, $5, $6, $7)`

// 	for _, item := range news {
// 		if pc.newsItemExists(item) {
// 			continue
// 		}

// 		if _, err = pc.db.Exec(query, item.FeedID, item.Title, item.Description, item.Link, item.Published, item.Parsed, item.Img); err != nil {
// 			return count, err
// 		}
// 		count++
// 	}

// 	return count, nil
// }
