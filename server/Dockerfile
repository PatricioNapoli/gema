FROM golang:alpine
EXPOSE 8080
WORKDIR /go/src/gema/server

VOLUME /go/src/gema/server
VOLUME /go/bin/static
VOLUME /go/bin/views

RUN apk add git

COPY deps.sh deps.sh
RUN ./deps.sh

COPY . .

COPY static /go/bin/static
COPY views /go/bin/views

ENTRYPOINT ["./entrypoint.sh"]
