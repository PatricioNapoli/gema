#!/bin/sh

if [ "$ENVIRONMENT" != "DEVELOPMENT" ]; then
    cp -a /go/src/gema/server/static/. /static/
    gema-server
else
    cp -a /go/src/gema/server/static/. /static/ 
    gin --bin ../../../bin/gema-server -x static/ -x templates/ --all --immediate --appPort 8080
fi