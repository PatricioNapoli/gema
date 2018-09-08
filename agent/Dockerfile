FROM golang:alpine
EXPOSE 80
WORKDIR /go/src/gema/agent

RUN apk add git

RUN go get -v -u github.com/go-redis/redis

COPY . .

RUN go build -v -o /go/bin/gema-agent main.go

CMD gema-agent
