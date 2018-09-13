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
* Memcached
* NextCloud
* Sentry
* Elasticsearch
* Logstash
* Kibana
* FileBeat
* MetricBeat

## Generating a self signed certificate.
`sudo openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout ./certs/domain.key -out ./certs/domain.crt`  
`sudo openssl dhparam -out ./certs/dhparam.pem 2048`

## Auth for Registry
`auth/htpasswd`  
Use htpasswd online generator, bcrypt algorithm.

## Setting up environment.
Setup preset.env into .env.  

Need to migrate dashboards from metricbeat and filebeat after deploying.  
SSH to metricbeat and filebeat containers and run `./metricbeat setup -e -strict.perms=false -E setup.kibana.host=kibana:5601`.  

Create databases gema_nc and gema_sentry.  

SSH to sentry and run ./entrypoint.sh upgrade.  

Login to cloud and run migration from admin.  

ID for Grafana dashboard: 1860.  

## Localhost resolution
Use DNSMASQ for wildcard localhost subdomain resolution, like *.localhost.
Add to /etc/dnsmasq.conf:

`address=/.localhost/127.0.0.1`

## Running

Setup sysctl map count if needed. (see node-init.sh)  
Setup main node label.  

`docker node update --label-add category=main <HOSTNAME>`  

Afterwards:  

`docker swarm init`  
`./deploy.sh`

## Chat

`mkdir -pv ./chat/volumes/app/mattermost/{data,logs,config,plugins}`  
`sudo chown -R 2000:2000 ./chat/volumes/app/mattermost/`  

## TODO

* Internal proxy auth, create dashboard for creating users. These users work for the HQ and the internal services, limit access to kibana.
* Cluster redis and postgres for data redundancy.
* GEMA dashboard, hq.geminis.io, live ops, reports generation, ML, message broadcast through WS, push notifications, service routings, let client have a user
* CORS.
* Check Performance of Proxy.
* Disable services logins.
* Disable retina detect particle js

* Code Documentation, Wiki.js vs BookStack
* Unit test GEMA
* Setup Gitlab yml CI/CD, Trigger SonarQube with Gitlab Push
* Define /health for services, (HEALTHCHECK dockerfile) GEMA Dashboard should test those for health.
* Define logstash log pattern for all services in GEMA for parsing.
* Use wikis for documentation.

* Close logstash UDP port

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
