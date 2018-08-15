# Geminis Architecture

Project for centralizing Geminis architecture backend, as well as LiveOps to support the everquest of integrating new applications into our baseline.

This project aims to swarm the following services:

* Docker Registry
* Portainer
* Netdata
* NGINX
* NGINX Amplify
* C++ Websocket
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

##
`auth/htpasswd`
`user password`
