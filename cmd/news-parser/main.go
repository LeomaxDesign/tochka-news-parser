package main

import (
	"log"

	"github.com/LeomaxDesign/tochka-news-parser/internal/news-parser/parser"
	"github.com/LeomaxDesign/tochka-news-parser/internal/news-parser/repository"
	"github.com/LeomaxDesign/tochka-news-parser/internal/news-parser/repository/postgres"
	"github.com/LeomaxDesign/tochka-news-parser/internal/news-parser/web"

	"github.com/spf13/viper"
)

const (
	filename = "config"
	filepath = "."
)

func main() {
	var err error

	viper.SetConfigName(filename)
	viper.AddConfigPath(filepath)
	if err = viper.ReadInConfig(); err != nil {
		log.Fatal("error loading config ", err)
	}

	repo := repository.New(
		viper.GetString(`postgres.host`),
		viper.GetString(`postgres.user`),
		viper.GetString(`postgres.password`),
		viper.GetString(`postgres.dbname`),
		viper.GetInt(`postgres.port`),
	)

	if err = repo.Connect(); err != nil {
		log.Fatal("error connecting to database:", err)
	}
	defer repo.Disconnect()

	newsFeedRepo := postgres.NewNewsFeedRepo(repo.DB)
	newsRepo := postgres.NewNewsRepo(repo.DB)

	parser := parser.New(newsFeedRepo, newsRepo)

	if err = parser.CheckNews(); err != nil {
		log.Fatal("error connecting to database:", err)
	}

	server := web.New(parser, viper.GetString(`address`))
	if err := server.Start(); err != nil {
		log.Fatal(err)
	}

}
