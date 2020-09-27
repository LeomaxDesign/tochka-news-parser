package postgres

import (
	"context"
	"fmt"

	"github.com/LeomaxDesign/tochka-news-parser/internal/news-parser/repository"
	"github.com/jackc/pgx/v4"
)

type newsRepository struct {
	db *pgx.Conn
}

func NewNewsRepo(db *pgx.Conn) *newsRepository {
	return &newsRepository{
		db: db,
	}
}

func (r *newsRepository) GetAll(searchString string) ([]*repository.News, error) {
	var whereQuery string

	query := `SELECT 
				id,
				title,
				description,
				link,
				published,
				img
			FROM feed_news
			%[1]s
			ORDER BY id DESC`

	if searchString != "" {
		whereQuery = fmt.Sprintf("WHERE title ILIKE '%%' || '%[1]s' || '%%'", searchString)
	}

	query = fmt.Sprintf(query, whereQuery)

	rows, err := r.db.Query(context.Background(), query)
	if err != nil {
		return nil, fmt.Errorf("failed to create query: %w", err)
	}
	defer rows.Close()

	news := make([]*repository.News, 0)
	for rows.Next() {

		ni := &repository.News{}

		if err = rows.Scan(
			&ni.ID,
			&ni.Title,
			&ni.Description,
			&ni.Link,
			&ni.Published,
			&ni.Img,
		); err != nil {
			return nil, fmt.Errorf("failed to scan: %w", err)
		}

		news = append(news, ni)
	}

	return news, nil
}

func (r *newsRepository) Add(news *repository.News) error {
	var err error

	query := `INSERT INTO feed_news (feed_id, title, description, link, published, parsed, img) VALUES ($1, $2, $3, $4, $5, $6, $7)`

	if _, err = r.db.Exec(context.Background(), query, news.FeedID, news.Title, news.Description, news.Link, news.Published, news.Parsed, news.Img); err != nil {
		return fmt.Errorf("failed to exec: %w", err)
	}

	return nil
}

func (r *newsRepository) IsExists(news *repository.News) (bool, error) {
	var id int
	if err := r.db.QueryRow(context.Background(), "SELECT id FROM feed_news WHERE title = $1 OR link = $2", news.Title, news.Link).Scan(&id); err != nil {
		return false, fmt.Errorf("failed to query row: %w", err)
	}

	if id != 0 {
		return true, nil
	}

	return false, nil
}
