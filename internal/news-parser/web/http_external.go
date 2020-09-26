package web

import (
	"log"
	"net/http"

	"github.com/LeomaxDesign/tochka-news-parser/internal/news-parser/parser"
	"github.com/gorilla/mux"
)

// Server ...
type Server struct {
	router  *mux.Router
	parser  parser.Service
	address string
}

// New ...
func New(parser parser.Service, address string) *Server {
	return &Server{
		router:  mux.NewRouter(),
		parser:  parser,
		address: address,
	}
}

// Start ...
func (s *Server) Start() error {
	var err error

	s.NewRouter()

	log.Println("Server started")
	if err = http.ListenAndServe(s.address, s.router); err != nil {
		return err
	}

	return nil
}

// NewRouter ...
func (s *Server) NewRouter() {
	s.router.HandleFunc("/feed/add", s.handleAddNewsFeed).Methods("POST")
	s.router.HandleFunc("/news", s.handleGetNews).Methods("GET")
	s.router.HandleFunc("/test", s.handleTest).Methods("GET")
}
