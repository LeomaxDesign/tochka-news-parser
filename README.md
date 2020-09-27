# tochka-news-parser

Задача: Написать агрегатор, который принимает адрес сайта или RSS-ленту, и начинает автоматически пополнять БД новостями с этого ресурса. У пользователя должна быть возможность просмотра новостей и поиск по подстроке заголовка.


## Настройка

Все настройки производятся в  файле .env

#### .env
```
APP_ADDRESS=
APP_PORT=8000

POSTGRES_HOST=postgres
POSTGRES_PORT=5432
POSTGRES_USER=postgres
POSTGRES_PASSWORD=postgres
POSTGRES_DB=news_feed
```


## Запуск

Для запуска приложения необходимо вызвать скрипт 
```
./run.sh
```

#### run.sh
Скрипт сбилдит проект и  запустит docker-compose, который, в свою очередь, запустит 2 контейнера: `postgres (database)`, `news-parser (app)`.
```
#!/bin/sh
export GOOS=linux
export GOARCH=amd64
go build -o ./news-parser ./cmd/news-parser/main.go

docker-compose up --build
```
# API

## Add news feed

### Request:

`POST /newsfeed/add`

```json
{
  "url": "",
  "name": "",
  "type": 0,
  "frequency": 0,
  "parse_count": 0,
  "item_tag": "",
  "title_tag": "",
  "description_tag": "",
  "link_tag": "",
  "published_tag": "",
  "img_tag": ""
}
```

`frequency`: frequency to parse news (in seconds)  
`parse_count`: max items to parse for each parser call (0 - parse all available)  
`type 0|1`: 0 - rss, 1 - html  
`link_tag`: will take ***href*** attribute from the found tag  
`img_tag`: will take ***src*** attribute from the found tag  

### Response:
**Success response**  

  * **Code:** `200`<br />
    **Content:** `news feed succesfully added`

**Error response**  

  * **Code:** `400`<br />
    **Content:** `invalid request`<br />
    **Description:** something wrong with input data

  * **Code:** `400`<br />
    **Content:** `invalid URL`<br />
    **Description:** URL address in request is not valid
    
  * **Code:** `500`<br />
    **Content:** `failed to add news feed`<br />
    **Description:** error adding news feed to db

### Example:  
`POST /newsfeed/add`

**RSS**
```json
{
  "url": "https://news.yandex.ru/auto.rss",
  "name": "Auto Yandex",
  "type": 0,
  "frequency": 10,
  "parse_count": 3
}
```  
**HTML**
```json
{
  "url": "https://lenta.ru/",
  "name": "Lenta News",
  "type": 1,
  "frequency": 20,
  "parse_count": 0,
  "item_tag": ".item.article",
  "title_tag": "h3",
  "description_tag": ".rightcol",
  "link_tag": ".picture",
  "img_tag": "img"
}
```  
## Get news

### Request:

`GET /news`

**URL Params**  

***Optional:***  
`title=[string]`: find news by title

### Response
**Success response**  

  * **Code:** `200`<br />
    **Content:**  
    ```json
    [
      {
        "id": 0,
        "title": "",
        "description": "",
        "link": "",
        "img": "",
        "published": ""
      }
    ]
    ```  
**Error response**  

  * **Code:** `500`<br />
    **Content:** `failed to get news`<br />
    **Description:** failed to get news from database

  * **Code:** `500`<br />
    **Content:** `failed to marshal`<br />
    **Description:** failed to marshal response data
    
