FROM golang:alpine
EXPOSE 8080
WORKDIR /go/src/gema/server

VOLUME /go/src/gema/server
VOLUME /go/bin/static
VOLUME /go/bin/templates

RUN apk add git

RUN go get -v -u github.com/kataras/iris
RUN go get -v -u github.com/go-pg/pg
RUN go get -v -u github.com/codegangsta/gin
RUN go get -v -u golang.org/x/crypto/bcrypt
RUN go get -v github.com/iris-contrib/middleware/raven

COPY . .

COPY static /go/bin/static
COPY templates /go/bin/templates

ENTRYPOINT ["./entrypoint.sh"]
