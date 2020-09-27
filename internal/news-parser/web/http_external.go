package web

import (
	"log"
	"net/http"

	"github.com/LeomaxDesign/tochka-news-parser/internal/news-parser/parser"
	"github.com/gorilla/mux"
)

type Server struct {
	router  *mux.Router
	parser  parser.Service
	address string
}

func New(parser parser.Service, address string) *Server {
	return &Server{
		router:  mux.NewRouter(),
		parser:  parser,
		address: address,
	}
}

func (s *Server) Start() error {
	var err error

	s.NewRouter()

	log.Println("Server started")
	if err = http.ListenAndServe(s.address, s.router); err != nil {
		return err
	}

	return nil
}

func (s *Server) NewRouter() {
	s.router.HandleFunc("/newsfeed/add", s.handleAddNewsFeed).Methods("POST")
	s.router.HandleFunc("/news", s.handleGetNews).Methods("GET")
	s.router.HandleFunc("/test", s.handleTest).Methods("GET")
}
