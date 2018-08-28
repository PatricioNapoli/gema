#!/bin/sh

if [[ -v $ENV_PROD ]]; then
    go build -o /go/bin/gema-server go/main.go
    gema-server
else
    gin --bin ../../../bin/gema-server --excludeDir static --excludeDir views --all --immediate --appPort 8080
fi