#!/bin/sh

if [ "$ENVIRONMENT" != "dev" ]; then
    mkdir /static/gema; cp -a /usr/local/go/src/gema/server/static/gema/. /static/gema-dash/
    gema-server
else
    mkdir /static/gema; cp -a /usr/local/go/src/gema/server/static/gema/. /static/gema-dash/
    gin --bin ../../../bin/gema-server -x static/ -x templates/ --all --immediate --appPort 8080
fi