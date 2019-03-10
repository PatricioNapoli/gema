#!/bin/sh

docker stack rm chat
docker build -t localhost:5000/gema/chat chat/
docker stack deploy -c chat/chat-compose.yml chat