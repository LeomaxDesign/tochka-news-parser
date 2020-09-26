package web

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/LeomaxDesign/tochka-news-parser/internal/news-parser/repository"
)

func (s *Server) handleTest(w http.ResponseWriter, r *http.Request) {
	if _, err := w.Write([]byte("Ok")); err != nil {
		http.Error(w, "failed to write", http.StatusInternalServerError)
		return
	}
}

func (s *Server) handleAddNewsFeed(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	feed := repository.NewsFeed{}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		s.logger.Println("Error adding new feed: ", err)
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(body, &feed)
	if err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	if u, err := url.Parse(feed.URL); err != nil || u.Scheme == "" || u.Host == "" {
		http.Error(w, "failed to parse URL", http.StatusBadRequest)
		return
	}

	if err = s.parser.AddNewsFeed(&feed); err != nil {
		s.logger.Println("error adding new feed: ", err)
		http.Error(w, "failed to add new feed", http.StatusBadRequest)
		return
	}

	if _, err = w.Write([]byte("news feed succesfully added")); err != nil {
		http.Error(w, "failed to write", http.StatusInternalServerError)
		return
	}
}

func (s *Server) handleGetNews(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	searchString := r.URL.Query().Get("title")

	news, err := s.parser.GetNews(searchString)
	if err != nil {
		s.logger.Println("error get news: ", err)
		http.Error(w, "failed to get news", http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(news)
	if err != nil {
		http.Error(w, "failes to marshal", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if _, err = w.Write(data); err != nil {
		http.Error(w, "failed to write", http.StatusInternalServerError)
		return
	}
}
