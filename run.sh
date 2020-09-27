#!/bin/sh
export GOOS=linux
export GOARCH=amd64
go build -mod vendor -o ./news-parser ./cmd/news-parser/main.go

docker-compose up --build
