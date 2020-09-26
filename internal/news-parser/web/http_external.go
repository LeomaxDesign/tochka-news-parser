package web

import (
	"log"
	"net/http"
	"os"

	"github.com/LeomaxDesign/tochka-news-parser/internal/news-parser/parser"
	"github.com/gorilla/mux"
)

// Server ...
type Server struct {
	router *mux.Router
	logger *log.Logger
	parser parser.Service
}

// New ...
func New(parser parser.Service) *Server {
	return &Server{
		router: mux.NewRouter(),
		parser: parser,
		logger: log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile),
	}
}

// Start ...
func (s *Server) Start() error {
	var err error

	s.NewRouter()

	s.logger.Println("Server started")
	if err = http.ListenAndServe(":8000", s.router); err != nil {
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
