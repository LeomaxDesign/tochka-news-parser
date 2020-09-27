FROM alpine:latest

WORKDIR /bin

RUN apk add --no-cache bash

ARG APP_PORT

COPY news-parser /bin 
COPY .env /bin  
COPY ./docker/wait-for-it.sh /bin 


RUN chmod +x /bin wait-for-it.sh
CMD ["./wait-for-it.sh", "postgres:5432", "--", "news-parser"]
