FROM golang:alpine
EXPOSE 8080
WORKDIR /go/src/gema/agent

VOLUME /configs

RUN apk add git

RUN go get -v -u github.com/go-redis/redis

COPY . .

ENTRYPOINT ["./entrypoint.sh"]
