FROM golang:alpine
EXPOSE 80
WORKDIR /go/src/gema/backup

COPY . .

RUN go build -v -o /go/bin/gema-backup main.go

CMD gema-backup
