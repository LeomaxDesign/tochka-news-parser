package parser

import (
	"errors"
	"log"
	"time"

	"github.com/LeomaxDesign/tochka-news-parser/internal/news-parser/repository"

	"github.com/mmcdole/gofeed"
)

// Parser ...
type Parser struct {
	newsFeeds map[string]repository.NewsFeedModel
	repo      *repository.PostgresClient
	logger    *log.Logger
}

// New ...
func New(repo *repository.PostgresClient, logger *log.Logger) *Parser {
	return &Parser{
		newsFeeds: make(map[string]repository.NewsFeedModel),
		repo:      repo,
		logger:    logger,
	}
}

// AddNewsFeed ...
func (p *Parser) AddNewsFeed(feed *repository.NewsFeedModel) {
	if _, ok := p.newsFeeds[feed.URL]; ok {
		return
	}
	p.newsFeeds[feed.URL] = *feed
}

// Parse ...
func (p *Parser) Parse(newsFeed repository.NewsFeedModel) error {
	var (
		err  error
		news []*repository.NewsModel
	)

	if p.IsRSS(newsFeed) {
		news, err = p.ParseRSS(newsFeed)
	} else {
		news, err = p.ParseHTML(newsFeed)
	}

	if news != nil {
		p.logger.Printf("adding %d news for %s to db\n", len(news), newsFeed.URL)
		if err = p.repo.CheckIfExistAndAddNews(news); err != nil {
			return err
		}
		p.logger.Printf("%d news successfully added for %s\n", len(news), newsFeed.URL)

	}

	return nil
}

// ParseRSS ...
func (p *Parser) ParseRSS(newsFeed repository.NewsFeedModel) ([]*repository.NewsModel, error) {
	var (
		err  error
		feed *gofeed.Feed
	)

	fp := gofeed.NewParser()
	feed, err = fp.ParseURL(newsFeed.URL)
	if err != nil {
		return nil, err
	}

	if feed == nil {
		return nil, errors.New("feed is nil")
	}

	p.logger.Printf("parsed %d for %s\n", len(feed.Items), newsFeed.URL)
	news := make([]*repository.NewsModel, 0, len(feed.Items))
	for _, item := range feed.Items {
		var imgLink string
		if item.Image != nil {
			imgLink = item.Image.URL
		}
		news = append(news, &repository.NewsModel{
			FeedID:      newsFeed.ID,
			Title:       item.Title,
			Description: item.Description,
			Img:         imgLink,
			Link:        item.Link,
			Published:   item.PublishedParsed,
			Parsed:      time.Now().UTC(),
		})

	}

	return news, nil
}

// ParseHTML ...
func (p *Parser) ParseHTML(newsFeed repository.NewsFeedModel) ([]*repository.NewsModel, error) {

	return nil, nil
}

// CheckNews ...
func (p *Parser) CheckNews() {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	p.logger.Println("loading news feed")
	news, err := p.repo.GetNewsFeed()
	if err != nil {
		log.Println("error get news feed")
		return
	}

	for k := range news {
		p.AddNewsFeed(news[k])
	}

	p.logger.Printf("news feed loaded successful, total: %d", len(news))

	// for range ticker.C {
	for k := range p.newsFeeds {
		p.logger.Printf("starting parsing for: %s", p.newsFeeds[k].URL)
		if err = p.Parse(p.newsFeeds[k]); err != nil {
			p.logger.Printf("error parsing news for %s: %s\n", p.newsFeeds[k].URL, err)
		}
	}
	// }

}

func (p *Parser) check() {

}

// IsRSS ...
func (p *Parser) IsRSS(newsFeed repository.NewsFeedModel) bool {
	if newsFeed.Type == 0 {
		return true
	}

	return false
}
