#!/bin/sh

if [ "$ENVIRONMENT" != "dev" ]; then
    cp -a /go/src/gema/server/static/. /static/$HQ_DOMAIN/static/
    gema-server
else
    cp -a /go/src/gema/server/static/. /static/$HQ_DOMAIN/static/
    gin --bin ../../../bin/gema-server -x static/ -x templates/ --all --immediate --appPort 8080
fi