package parser

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"sync"
	"time"

	"github.com/LeomaxDesign/tochka-news-parser/internal/news-parser/repository"
	"github.com/PuerkitoBio/goquery"

	"github.com/mmcdole/gofeed"
)

type Service interface {
	Parse(newsFeed *repository.NewsFeed) error
	CheckNews() error
	StartFrequencyParser(newsFeed repository.NewsFeed)
	AddNewsFeed(newsFeed *repository.NewsFeed) error
	GetNews(searchString string) ([]*repository.News, error)
}

type parser struct {
	mu           sync.Mutex
	wg           sync.WaitGroup
	newsFeeds    map[string]repository.NewsFeed
	newsFeedRepo repository.NewsFeedRepository
	newsRepo     repository.NewsRepository
}

func New(newsFeedRepo repository.NewsFeedRepository, newsRepo repository.NewsRepository) *parser {
	return &parser{
		newsFeeds:    make(map[string]repository.NewsFeed),
		newsFeedRepo: newsFeedRepo,
		newsRepo:     newsRepo,
	}
}

func (p *parser) Parse(newsFeed *repository.NewsFeed) error {
	var (
		err   error
		news  []*repository.News
		count int
	)

	if newsFeed.IsRSS() {
		if news, err = p.parseRSS(newsFeed); err != nil {
			log.Println("failed to parse rss", err)
			return err
		}

	} else {
		if news, err = p.parseHTML(newsFeed); err != nil {
			log.Println("failed to parse html", err)
			return err
		}
	}

	for _, newsItem := range news {
		exists, err := p.newsRepo.IsExists(newsItem)
		if err != nil {
			log.Println("failed to check is exists", err)
			continue
		}

		if exists {
			continue
		}

		if err = p.newsRepo.Add(newsItem); err != nil {
			log.Println("failed to insert news item", err)
		}
		count++
	}

	if count == 0 {
		log.Printf("no new news for %s\n", newsFeed.URL)
	} else {
		log.Printf("%d news added for %s\n", count, newsFeed.URL)
	}

	return nil
}

// parserSS ...
func (p *parser) parseRSS(newsFeed *repository.NewsFeed) ([]*repository.News, error) {
	var (
		err  error
		feed *gofeed.Feed
	)

	fp := gofeed.NewParser()
	feed, err = fp.ParseURL(newsFeed.URL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse url feed: %w", err)
	}

	if feed == nil {
		return nil, errors.New("feed is nil")
	}

	log.Printf("find %d news for %s\n", len(feed.Items), newsFeed.URL)
	news := make([]*repository.News, 0, len(feed.Items))
	for id, item := range feed.Items {

		if newsFeed.ParseCount > 0 && id >= newsFeed.ParseCount {
			break
		}

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

		news = append(news, &repository.News{
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
func (p *parser) parseHTML(newsFeed *repository.NewsFeed) ([]*repository.News, error) {
	resp, err := http.Get(newsFeed.URL)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to get new document from reader: %w", err)
	}

	items := doc.Find(newsFeed.ItemTag)

	log.Printf("find %d news for %s\n", items.Length(), newsFeed.URL)

	news := make([]*repository.News, 0, items.Length())
	items.EachWithBreak(func(i int, s *goquery.Selection) bool {

		if newsFeed.ParseCount > 0 && i >= newsFeed.ParseCount {
			return false
		}

		var (
			imgLink     string
			description string
		)

		imgLink = s.Find(newsFeed.ImgTag).AttrOr("src", "")

		description = s.Find(newsFeed.DescriptionTag).Text()

		if p.checkSpecialSymbols(description) {
			description = p.replaceSpecialSymbols(description)
		}

		// TODO: parse html time ?
		// published, _ = time.Parse("", s.Find(newsFeed.PublishedTag).Text())

		news = append(news, &repository.News{
			FeedID:      newsFeed.ID,
			Title:       s.Find(newsFeed.TitleTag).Text(),
			Description: description,
			Img:         imgLink,
			Link:        s.Find(newsFeed.LinkTag).AttrOr("href", ""),
			Published:   time.Now(),
			Parsed:      time.Now(),
		})

		return true
	})

	return news, nil
}

// CheckNews ...
func (p *parser) CheckNews() error {
	var err error

	log.Println("loading news feed")
	newsFeeds, err := p.newsFeedRepo.GetAll()
	if err != nil {
		return fmt.Errorf("failed to get all news feeds: %w", err)
	}

	for k := range newsFeeds {
		p.addToNewsFeedsMap(newsFeeds[k])
	}

	log.Printf("news feeds loaded successful, total: %d", len(newsFeeds))

	p.wg.Add(len(p.newsFeeds))
	for _, feed := range p.newsFeeds {
		go p.StartFrequencyParser(feed)
	}
	p.wg.Wait()

	return nil
}

func (p *parser) StartFrequencyParser(newsFeed repository.NewsFeed) {
	var err error
	log.Printf("loading parser for %s every %s", newsFeed.URL, time.Duration(newsFeed.Frequency)*time.Second)

	ticker := time.NewTicker(time.Duration(newsFeed.Frequency) * time.Second)
	defer ticker.Stop()

	p.wg.Done()

	for range ticker.C {
		log.Printf("starting parsing for: %s", newsFeed.URL)
		if err = p.Parse(&newsFeed); err != nil {
			log.Printf("error parsing news for %s: %s\n", newsFeed.URL, err)
		}

	}

}

func (p *parser) AddNewsFeed(newsFeed *repository.NewsFeed) error {
	var err error

	if err = p.newsFeedRepo.Add(newsFeed); err != nil {
		return fmt.Errorf("failed to add news repo: %w", err)
	}

	p.addToNewsFeedsMap(newsFeed)

	p.wg.Add(1)
	go p.StartFrequencyParser(*newsFeed)

	return nil
}

func (p *parser) GetNews(searchString string) ([]*repository.News, error) {
	var (
		err  error
		news []*repository.News
	)

	if news, err = p.newsRepo.GetAll(searchString); err != nil {
		return nil, fmt.Errorf("failed to get all news: %w", err)
	}

	return news, err
}

func (p *parser) checkSpecialSymbols(content string) bool {
	var specialSymbol = regexp.MustCompile(`(&)`)
	return specialSymbol.MatchString(content)
}

func (p *parser) replaceSpecialSymbols(content string) string {
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

func (p *parser) addToNewsFeedsMap(newsFeed *repository.NewsFeed) {
	defer p.mu.Unlock()

	p.mu.Lock()
	{
		if _, ok := p.newsFeeds[newsFeed.URL]; ok {
			return
		}

		p.newsFeeds[newsFeed.URL] = *newsFeed
	}
}
