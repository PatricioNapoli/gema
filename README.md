# Geminis Architecture

Project for centralizing Geminis architecture backend, as well as LiveOps to support the everquest of integrating new applications into our baseline.

This project aims to swarm the following services:

* Websocket server
* HTTP Server LiveOps
* Docker Registry
* Portainer
* Prometheus
* Grafana
* NGINX
* NGINX Amplify
* PostgreSQL
* PGAdmin4
* Redis
* Kafka
* Elasticsearch
* Logstash
* Kibana

## Generating a self signed certificate to test.
`sudo openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout ./certs/domain.key -out ./certs/domain.crt`
`sudo openssl dhparam -out ./certs/dhparam.pem 2048`

## Auth
`auth/htpasswd`
Use htpasswd online generator, bcrypt algorithm.

## Env vars
Setup preset.env into .env.

## Localhost resolution
Add to /etc/hosts:

127.0.0.1 portainer.localhost
127.0.0.1 pgadmin.localhost
127.0.0.1 registry.localhost
127.0.0.1 prometheus.localhost
127.0.0.1 grafana.localhost
127.0.0.1 kibana.localhost

TODO: use unix sockets for DB, internal proxy auth, internal proxy discovery agent
