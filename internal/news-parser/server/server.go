package server

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/LeomaxDesign/tochka-news-parser/internal/news-parser/parser"
	"github.com/LeomaxDesign/tochka-news-parser/internal/news-parser/repository"

	"github.com/gorilla/mux"
)

// Server ...
type Server struct {
	router *mux.Router
	repo   *repository.PostgresClient
	logger *log.Logger
	parser *parser.Parser
}

// New ...
func New() *Server {
	return &Server{
		router: mux.NewRouter(),
		repo:   repository.New("localhost", "postgres", "postgres", "news_feed", 5432),
		logger: log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile),
	}
}

// Start ...
func (s *Server) Start() error {
	var err error

	s.NewRouter()

	if err = s.repo.Connect(); err != nil {
		return err
	}
	defer s.repo.Disconnect()
	s.logger.Println("Database connected")

	s.parser = parser.New(s.repo, s.logger)

	go s.parser.CheckNews()

	s.logger.Println("Server started")
	if err = http.ListenAndServe(":8000", s.router); err != nil {
		return err
	}

	return nil
}

// NewRouter ...
func (s *Server) NewRouter() {
	s.router.HandleFunc("/feed/add", s.handleAddFeed).Methods("POST")
	s.router.HandleFunc("/news", s.handleGetNews).Methods("GET")
	s.router.HandleFunc("/test", s.handleTest).Methods("GET")
}

func (s *Server) handleTest(rw http.ResponseWriter, r *http.Request) {
	rw.Write([]byte("Ok"))
}

func (s *Server) handleAddFeed(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	feed := repository.NewsFeedModel{}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		s.logger.Println("Error adding new feed: ", err)
		http.Error(w, "Request is not valid", http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(body, &feed)
	if err != nil {
		http.Error(w, "Request is not valid", http.StatusBadRequest)
		return
	}

	if _, err = url.Parse(feed.URL); err != nil {
		http.Error(w, "Link is not valid", http.StatusBadRequest)
		return
	}

	if err = s.repo.AddNewsFeed(feed); err != nil {
		s.logger.Println("Error adding new feed: ", err)
		http.Error(w, "Error adding new feed", http.StatusBadRequest)
		return
	}

	s.parser.AddNewsFeed(&feed)

	w.Write([]byte("new feed succesfully added"))
	w.WriteHeader(http.StatusOK)
}

func (s *Server) handleGetNews(rw http.ResponseWriter, r *http.Request) {

	// feed, err := s.parser.Parse("https://www.liga.net/biz/all/rss.xml")
	// if err != nil {
	// 	log.Println("ERROR: ", err)
	// }

	// data, err := json.Marshal(feed.Items[0])
	// if err != nil {
	// 	log.Println("ERROR: ", err)
	// }
	// rw.Write(data)
	// return
}
