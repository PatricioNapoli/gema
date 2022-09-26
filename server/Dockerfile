FROM golang:alpine
EXPOSE 8080
WORKDIR /usr/local/go/src/gema/server

VOLUME /usr/local/go/src/gema/server
VOLUME /usr/local/go/bin/static
VOLUME /usr/local/go/bin/templates

RUN apk add git build-base

ENV GO111MODULE=on

RUN go get -v -u github.com/kataras/iris@latest
RUN go get -v -u github.com/go-pg/pg
RUN go get -v -u github.com/codegangsta/gin
RUN go get -v -u golang.org/x/crypto/bcrypt
RUN go get -v -u github.com/iris-contrib/middleware/raven
RUN go get -v -u go.elastic.co/apm
RUN go get -v -u github.com/go-redis/redis
RUN go get -v -u github.com/gorilla/websocket

COPY . .

COPY templates /usr/local/go/bin/templates

RUN go build -mod=mod -v -o /usr/local/go/bin/gema-server main.go
ENTRYPOINT ["./entrypoint.sh"]
