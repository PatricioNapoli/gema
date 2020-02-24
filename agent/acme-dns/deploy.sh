#!/bin/sh

docker stack rm acme
docker stack deploy -c acme-dns/acme-compose.yml acme
