#!/bin/sh

docker build -t localhost:5000/gema/server .
docker stack rm server
docker deploy -c server-compose.yml server
