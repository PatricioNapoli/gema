# GEMA - Geminis Architecture

Project for centralizing Geminis architecture backend, as well as LiveOps to support the everquest of integrating new applications into our baseline.

This project aims to swarm the following services:

* GEMA Websocket Server
* GEMA LiveOps Server Dashboard
* GEMA FedAuth Layer
* GEMA Service Proxy Layer
* GEMA Service Discovery Agent

Using the following environment:

* Docker Registry
* Portainer
* NGINX
* PostgreSQL
* PGAdmin4
* Redis
* Elasticsearch
* Logstash
* Kibana
* FileBeat
* MetricBeat

## Generating a self signed certificate.
`sudo openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout ./certs/domain.key -out ./certs/domain.crt`
`sudo openssl dhparam -out ./certs/dhparam.pem 2048`

## Auth (deprecated)
`auth/htpasswd`  
Use htpasswd online generator, bcrypt algorithm.

## Setting up env variables
Setup preset.env into .env.  

Need to migrate dashboards from metricbeat and filebeat after deploying.  
SSH to metricbeat and filebeat containers and run `./metricbeat setup -E setup.kibana.host=kibana:5601`.

## Localhost resolution
Add to /etc/hosts:  

```
127.0.0.1 hq.localhost
127.0.0.1 portainer.localhost
127.0.0.1 pgadmin.localhost
127.0.0.1 registry.localhost
127.0.0.1 kibana.localhost
```

## Running

Setup sysctl map count if needed. (see node-init.sh)  

`docker swarm init`  
`./deploy.sh`  

## TODO

* Use unix sockets for DB 
* Internal proxy auth 
* Internal proxy discovery agent, options: https, port, websocket support, require auth, domain
* Cluster redis and postgres for data redundancy
* Investigate Swarm behaviour when deploying NGINX proxy, always use exposed master node.
* Stg environment considerations 
* GEMA dashboard, hq.geminis.io, live ops, reports generation, ML, message broadcast through WS, push notifications, service routings, let client have a user
* Use swarm secrets
* Use swarm configs
* Unit test GEMA
* Tweak bash for deployment
* Investigate service updating, (dynamic configs?)
* Setup Gitlab yml CI/CD
* Go Expvar
* Go APM
* Sentry Docker
* OwnCloud Docker
* Filebeat Syslog

---

### Data Analytics
Genetic algorithm on top of models, supervising:  
* Hidden markov model
* Recurrent Neural network (RNN, LSTM) / GAN / CNN, KERAS
* ARIMA-GARCH (error varies over time)

---

### CI/CD
Setup for CI:

* Code quality
* Test Coverage
* SAST & DAST
* Performance Test
* Project badges

These tests should be performed through CI in Host machine, so as to use it in every repo.
