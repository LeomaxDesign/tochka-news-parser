FROM alpine:latest

WORKDIR /bin

RUN apk add --no-cache bash

ARG APP_PORT

COPY news-parser /bin 
COPY config.json /bin 
COPY wait-for-it.sh /bin


EXPOSE ${APP_PORT}

RUN chmod +x /bin wait-for-it.sh
CMD ["./wait-for-it.sh", "postgres:5432", "--", "news-parser"]
