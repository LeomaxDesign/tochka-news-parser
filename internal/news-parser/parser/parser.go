package parser

import (
	"errors"
	"log"
	"net/http"
	"regexp"
	"sync"
	"time"

	"github.com/LeomaxDesign/tochka-news-parser/internal/news-parser/repository"
	"github.com/PuerkitoBio/goquery"

	"github.com/mmcdole/gofeed"
)

// Parser ...
type Parser struct {
	mu        sync.Mutex
	Wg        sync.WaitGroup
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
	p.mu.Lock()
	{
		if _, ok := p.newsFeeds[feed.URL]; ok {
			return
		}
		p.newsFeeds[feed.URL] = *feed
	}
	p.mu.Unlock()
}

// Parse ...
func (p *Parser) Parse(newsFeed *repository.NewsFeedModel) error {
	var (
		err   error
		news  []*repository.NewsModel
		count int
	)

	if p.repo.IsRSS(newsFeed) {
		news, err = p.ParseRSS(newsFeed)
	} else {
		news, err = p.ParseHTML(newsFeed)
	}

	if news != nil {
		if count, err = p.repo.CheckIfExistAndAddNews(news); err != nil {
			return err
		}

		if count == 0 {
			p.logger.Printf("no new news for %s\n", newsFeed.URL)
		} else {
			p.logger.Printf("%d news added for %s\n", count, newsFeed.URL)
		}

	}

	return nil
}

// ParseRSS ...
func (p *Parser) ParseRSS(newsFeed *repository.NewsFeedModel) ([]*repository.NewsModel, error) {
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

	p.logger.Printf("find %d news for %s\n", len(feed.Items), newsFeed.URL)
	news := make([]*repository.NewsModel, 0, len(feed.Items))
	for _, item := range feed.Items {
		var imgLink string

		if item.Image != nil {
			imgLink = item.Image.URL
		}

		if imgLink == "" && item.Enclosures != nil && item.Enclosures[0] != nil {
			imgLink = item.Enclosures[0].URL
		}

		if p.checkSpecialSymbols(item.Description) {
			item.Description = p.replaceSpecialSymbols(item.Description)
		}

		news = append(news, &repository.NewsModel{
			FeedID:      newsFeed.ID,
			Title:       item.Title,
			Description: item.Description,
			Img:         imgLink,
			Link:        item.Link,
			Published:   *item.PublishedParsed,
			Parsed:      time.Now(),
		})

	}

	return news, nil
}

// ParseHTML ...
func (p *Parser) ParseHTML(newsFeed *repository.NewsFeedModel) ([]*repository.NewsModel, error) {
	resp, err := http.Get(newsFeed.URL)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromReader(resp.Body)

	items := doc.Find(newsFeed.ItemTag)
	news := make([]*repository.NewsModel, 0, items.Length())

	p.logger.Printf("find %d news for %s\n", items.Length(), newsFeed.URL)
	items.Each(func(i int, s *goquery.Selection) {
		var (
			imgLink     string
			description string
		)

		imgLink = s.Find(newsFeed.ImgTag).AttrOr("src", "")

		description = s.Find(newsFeed.DescriptionTag).Text()

		if p.checkSpecialSymbols(description) {
			description = p.replaceSpecialSymbols(description)
		}

		// published, _ = time.Parse("", s.Find(newsFeed.PublishedTag).Text())

		news = append(news, &repository.NewsModel{
			FeedID:      newsFeed.ID,
			Title:       s.Find(newsFeed.TitleTag).Text(),
			Description: description,
			Img:         imgLink,
			Link:        s.Find(newsFeed.LinkTag).AttrOr("href", ""),
			Published:   time.Now(),
			Parsed:      time.Now(),
		})

	})

	return news, nil
}

// CheckNews ...
func (p *Parser) CheckNews() error {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	p.logger.Println("loading news feed")
	news, err := p.repo.GetNewsFeed()
	if err != nil {
		p.logger.Println("error get news feed")
		return err
	}

	for k := range news {
		p.AddNewsFeed(news[k])
	}

	p.logger.Printf("news feeds loaded successful, total: %d", len(news))

	p.Wg.Add(len(p.newsFeeds))
	for _, feed := range p.newsFeeds {
		go p.StartFrequencyParser(feed)
	}
	p.Wg.Wait()

	return nil
}

// StartFrequencyParser ...
func (p *Parser) StartFrequencyParser(newsFeed repository.NewsFeedModel) {
	var err error
	p.logger.Printf("loading parser for %s every %s", newsFeed.URL, time.Duration(newsFeed.Frequency)*time.Second)

	ticker := time.NewTicker(time.Duration(newsFeed.Frequency) * time.Second)
	defer ticker.Stop()

	p.Wg.Done()

	for range ticker.C {
		p.logger.Printf("starting parsing for: %s", newsFeed.URL)
		if err = p.Parse(&newsFeed); err != nil {
			p.logger.Printf("error parsing news for %s: %s\n", newsFeed.URL, err)
		}

	}

}

func (p *Parser) checkSpecialSymbols(content string) bool {
	var specialSymbol = regexp.MustCompile(`(&)`)
	return specialSymbol.MatchString(content)
}

func (p *Parser) replaceSpecialSymbols(content string) string {
	specSymb := map[string]string{
		"&amp;":    `&`,
		"&lt;":     `<`,
		"&gt;":     `>`,
		"&nbsp;":   ` `,
		"&sect;":   `§`,
		"&copy;":   `©`,
		"&reg;":    `®`,
		"&deg;":    `°`,
		"&laquo;":  `«`,
		"&raquo;":  `»`,
		"&middot;": `·`,
		"&trade;":  `™`,
		"&plusmn;": `±`,
		"&quot;":   `"`,
		"&hellip;": `…`,
		"&ndash;":  `–`,
		"&mdash;":  `—`,
		"&lsquo;":  `‘`,
		"&rsquo;":  `’`,
		"&sbquo;":  `‚`,
		"&ldquo;":  `“`,
		"&rdquo;":  `”`,
		"&bdquo;":  `„`,
		"&prime;":  `′`,
		"&Prime;":  `″`,
	}
	for idx, symb := range specSymb {
		re := regexp.MustCompile(idx)
		content = re.ReplaceAllLiteralString(content, symb)
	}
	return content
}
