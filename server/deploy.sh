#!/bin/sh

docker build -t localhost:5000/gema/server .
docker stack rm server
docker stack deploy -c server-compose.yml server
