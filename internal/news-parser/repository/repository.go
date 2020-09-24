package repository

import (
	"database/sql"
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
func (pc *PostgresClient) AddNewsFeed(newsFeed NewsFeedModel) error {
	if _, err := pc.db.Exec("INSERT INTO feed (url, name, type, frequency, parse_count) VALUES ($1, $2, $3, $4, $5)", newsFeed.URL, newsFeed.Name, newsFeed.Type, newsFeed.Frequency, newsFeed.ParseCount); err != nil {
		return err
	}

	return nil
}

// GetNewsFeed ...
func (pc *PostgresClient) GetNewsFeed() ([]*NewsFeedModel, error) {

	rows, err := pc.db.Query("SELECT id, name, url, type, frequency, parse_count FROM feed")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	newsFeed := make([]*NewsFeedModel, 0)
	for rows.Next() {

		nf := &NewsFeedModel{}

		if err = rows.Scan(
			&nf.ID,
			&nf.Name,
			&nf.URL,
			&nf.Type,
			&nf.Frequency,
			&nf.ParseCount,
		); err != nil {
			return nil, err
		}

		newsFeed = append(newsFeed, nf)

	}
	return newsFeed, nil
}

// CheckIfExistAndAddNews ...
func (pc *PostgresClient) CheckIfExistAndAddNews(news []*NewsModel) error {
	var err error
	query := `INSERT INTO feed_news (feed_id, title, description, link, published, parsed, img) VALUES ($1, $2, $3, $4, $5, $6, $7)`

	for _, item := range news {
		if pc.newsItemExists(item) {
			continue
		}

		if _, err = pc.db.Exec(query, item.FeedID, item.Title, item.Description, item.Link, item.Published, item.Parsed, item.Img); err != nil {
			return err
		}

	}

	return nil
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
