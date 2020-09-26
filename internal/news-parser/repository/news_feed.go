package repository

type NewsFeedRepository interface {
	Add(newsFeed *NewsFeed) error
	GetAll() ([]*NewsFeed, error)
}

type NewsFeed struct {
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

func (m *NewsFeed) IsRSS() bool {
	return m.Type == 0
}
