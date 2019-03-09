#!/bin/sh

docker build -t localhost:5000/gema/server server/
docker stack rm server
docker stack deploy -c server/server-compose.yml server
