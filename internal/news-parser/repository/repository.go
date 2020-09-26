package repository

import (
	"database/sql"
	"fmt"
	"strconv"
)

// New ...
func New(host, username, password, dbname string, port int) *PostgresClient {
	return &PostgresClient{
		host:     host,
		username: username,
		password: password,
		dbname:   dbname,
		port:     port,
	}
}

// Connect ...
func (pc *PostgresClient) Connect() error {
	var err error

	pc.db, err = sql.Open("postgres", "host="+pc.host+" user="+pc.username+" dbname="+pc.dbname+" password="+pc.password+" port="+strconv.Itoa(pc.port)+" sslmode=disable")
	if err != nil {
		return err
	}

	if err = pc.db.Ping(); err != nil {
		return err
	}

	return nil
}

// Disconnect ...
func (pc *PostgresClient) Disconnect() {
	pc.db.Close()
}

// AddNewsFeed ...
func (pc *PostgresClient) AddNewsFeed(newsFeed *NewsFeedModel) error {
	if _, err := pc.db.Exec("INSERT INTO feed (url, title, type, frequency, parse_count, item_tag, title_tag, description_tag, link_tag, published_tag, img_tag) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)",
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
	); err != nil {
		return err
	}

	return nil
}

// GetNewsFeed ...
func (pc *PostgresClient) GetNewsFeed() ([]*NewsFeedModel, error) {
	var itemTag, titleTag, descriptionTag, linkTag, publishedTag, imgTag sql.NullString

	rows, err := pc.db.Query(`SELECT 
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
							FROM feed`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	newsFeed := make([]*NewsFeedModel, 0)
	for rows.Next() {

		nf := &NewsFeedModel{}

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
			return nil, err
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

// CheckIfExistAndAddNews ...
func (pc *PostgresClient) CheckIfExistAndAddNews(news []*NewsModel) (int, error) {
	var (
		err   error
		count int
	)
	query := `INSERT INTO feed_news (feed_id, title, description, link, published, parsed, img) VALUES ($1, $2, $3, $4, $5, $6, $7)`

	for _, item := range news {
		if pc.newsItemExists(item) {
			continue
		}

		if _, err = pc.db.Exec(query, item.FeedID, item.Title, item.Description, item.Link, item.Published, item.Parsed, item.Img); err != nil {
			return count, err
		}
		count++
	}

	return count, nil
}

func (pc *PostgresClient) newsItemExists(news *NewsModel) bool {
	var id int
	if err := pc.db.QueryRow("SELECT id FROM feed_news WHERE title = $1", news.Title).Scan(&id); err != nil {
		return false
	}

	if id != 0 {
		return true
	}

	return false

}

// IsRSS ...
func (pc *PostgresClient) IsRSS(newsFeed *NewsFeedModel) bool {
	if newsFeed.Type == 0 {
		return true
	}

	return false
}

// GetNews ...
func (pc *PostgresClient) GetNews(searchString string) ([]*NewsModel, error) {
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

	rows, err := pc.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	news := make([]*NewsModel, 0)
	for rows.Next() {

		ni := &NewsModel{}

		if err = rows.Scan(
			&ni.ID,
			&ni.Title,
			&ni.Description,
			&ni.Link,
			&ni.Published,
			&ni.Img,
		); err != nil {
			return nil, err
		}

		news = append(news, ni)
	}

	return news, nil
}
