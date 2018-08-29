#!/bin/sh

if [[ -v ENV_PROD ]]; then
    go build -o /go/bin/gema-server main.go
    gema-server
else
    gin --bin ../../../bin/gema-server -x static/ -x templates/ --all --immediate --appPort 8080
fi