#!/bin/sh

docker stack rm metrics
docker build -t localhost:5000/gema/node-exporter metrics/node-exporter/
docker stack deploy -c metrics/metrics-compose.yml metrics