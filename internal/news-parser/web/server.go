package web

import (
	"encoding/json"
	"io/ioutil"
	"log"
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

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("failed to read from body: ", err)
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	feed := repository.NewsFeed{}

	err = json.Unmarshal(body, &feed)
	if err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	if u, err := url.Parse(feed.URL); err != nil || u.Scheme == "" || u.Host == "" {
		log.Println("failed to parse URL: ", err)
		http.Error(w, "invalid URL", http.StatusBadRequest)
		return
	}

	if err = s.parser.AddNewsFeed(&feed); err != nil {
		log.Println("failed to add news feed: ", err)
		http.Error(w, "failed to add news feed", http.StatusInternalServerError)
		return
	}

	if _, err = w.Write([]byte("news feed succesfully added")); err != nil {
		log.Println("failed to write: ", err)
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
		log.Println("failed to get news: ", err)
		http.Error(w, "failed to get news", http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(news)
	if err != nil {
		log.Println("failed to marshal: ", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if _, err = w.Write(data); err != nil {
		log.Println("failed to write: ", err)
		http.Error(w, "failed to write", http.StatusInternalServerError)
		return
	}
}
