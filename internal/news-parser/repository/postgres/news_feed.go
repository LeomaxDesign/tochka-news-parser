package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/LeomaxDesign/tochka-news-parser/internal/news-parser/repository"
	"github.com/jackc/pgx/v4"
)

type newsFeedRepository struct {
	db *pgx.Conn
}

func NewNewsFeedRepo(db *pgx.Conn) *newsFeedRepository {
	return &newsFeedRepository{
		db: db,
	}
}

func (r *newsFeedRepository) Add(newsFeed *repository.NewsFeed) error {
	if err := r.db.QueryRow(context.Background(), `
		INSERT INTO news_feeds 
				(url, title, type, frequency, parse_count, item_tag, title_tag, description_tag, link_tag, published_tag, img_tag) 
			VALUES 
				($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id`,
		newsFeed.URL,
		newsFeed.Title,
		newsFeed.Type,
		newsFeed.Frequency,
		newsFeed.ParseCount,
		newsFeed.ItemTag,
		newsFeed.TitleTag,
		newsFeed.DescriptionTag,
		newsFeed.LinkTag,
		newsFeed.PublishedTag,
		newsFeed.ImgTag,
	).Scan(&newsFeed.ID); err != nil {
		return fmt.Errorf("failed to insert: %w", err)
	}

	return nil
}

func (r *newsFeedRepository) GetAll() ([]*repository.NewsFeed, error) {
	var itemTag, titleTag, descriptionTag, linkTag, publishedTag, imgTag sql.NullString

	rows, err := r.db.Query(context.Background(), `
		SELECT 
			id, 
			title, 
			url, 
			type, 
			frequency, 
			parse_count, 
			item_tag, 
			title_tag, 
			description_tag, 
			link_tag, 
			published_tag, 
			img_tag 
		FROM news_feeds`)
	if err != nil {
		return nil, fmt.Errorf("failed to create query: %w", err)
	}
	defer rows.Close()

	newsFeed := make([]*repository.NewsFeed, 0)
	for rows.Next() {

		nf := &repository.NewsFeed{}

		if err = rows.Scan(
			&nf.ID,
			&nf.Title,
			&nf.URL,
			&nf.Type,
			&nf.Frequency,
			&nf.ParseCount,
			&itemTag,
			&titleTag,
			&descriptionTag,
			&linkTag,
			&publishedTag,
			&imgTag,
		); err != nil {
			return nil, fmt.Errorf("failed to scan: %w", err)
		}

		nf.ItemTag = itemTag.String
		nf.TitleTag = titleTag.String
		nf.DescriptionTag = descriptionTag.String
		nf.LinkTag = linkTag.String
		nf.PublishedTag = publishedTag.String
		nf.ImgTag = imgTag.String

		newsFeed = append(newsFeed, nf)

	}
	return newsFeed, nil
}
