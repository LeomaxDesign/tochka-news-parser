package main

import (
	"log"

	_ "github.com/lib/pq"

	"github.com/LeomaxDesign/tochka-news-parser/internal/news-parser/server"
)

func main() {
	s := server.New()
	if err := s.Start(); err != nil {
		log.Fatal(err)
	}

}
