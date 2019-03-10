#!/bin/sh

docker build -t localhost:5000/gema/agent agent/
docker stack rm agent
docker stack deploy -c agent/agent-compose.yml agent
