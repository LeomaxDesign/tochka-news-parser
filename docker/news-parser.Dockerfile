FROM alpine:latest

WORKDIR /bin

RUN apk add --no-cache bash

COPY news-parser /bin 
COPY config.json /bin 
COPY wait-for-it.sh /bin


RUN chmod +x /bin wait-for-it.sh
CMD ["./wait-for-it.sh", "postgres:5432", "--", "news-parser"]
