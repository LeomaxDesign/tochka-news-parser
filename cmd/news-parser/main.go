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
	filename = ".env"
	filepath = "."
)

func main() {
	var err error

	viper.SetConfigType("env")
	viper.SetConfigName(filename)
	viper.AddConfigPath(filepath)
	if err = viper.ReadInConfig(); err != nil {
		log.Fatal("error loading config ", err)
	}

	repo := repository.New(
		viper.GetString(`POSTGRES_HOST`),
		viper.GetString(`POSTGRES_USER`),
		viper.GetString(`POSTGRES_PASSWORD`),
		viper.GetString(`POSTGRES_DB`),
		viper.GetInt(`POSTGRES_PORT`),
	)

	if err = repo.Connect(); err != nil {
		log.Fatal("failed connect to db:", err)
	}
	defer repo.Disconnect()

	newsFeedRepo := postgres.NewNewsFeedRepo(repo.DB)
	newsRepo := postgres.NewNewsRepo(repo.DB)

	parser := parser.New(newsFeedRepo, newsRepo)

	if err = parser.CheckNews(); err != nil {
		log.Fatal("failed to check news:", err)
	}

	address := viper.GetString(`APP_ADDRESS`) + ":" + viper.GetString(`APP_PORT`)

	server := web.New(parser, address)
	if err := server.Start(); err != nil {
		log.Fatal(err)
	}

}
