#!/bin/sh

if [ "$ENVIRONMENT" != "dev" ]; then
    mkdir /static/gema; cp -a /go/src/gema/server/static/gema/. /static/gema/$HQ_DOMAIN/
    gema-server
else
    mkdir /static/gema; cp -a /go/src/gema/server/static/gema/. /static/gema/$HQ_DOMAIN/
    gin --bin ../../../bin/gema-server -x static/ -x templates/ --all --immediate --appPort 8080
fi