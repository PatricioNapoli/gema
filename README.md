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
* Prometheus
* Grafana
* cAdvisor

## Generating a self signed certificate.
```
sudo openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout ./certs/domain.key -out ./certs/domain.crt` 
sudo openssl dhparam -out ./certs/dhparam.pem 2048
```

## Auth for Registry
`auth/htpasswd`  
Use htpasswd online generator, bcrypt algorithm.

## Setting up environment.
Setup preset.env into .env.  

Need to migrate dashboards from metricbeat and filebeat after deploying.  
SSH to metricbeat and filebeat containers and run 

```
./metricbeat setup -e -strict.perms=false -E setup.kibana.host=kibana:5601
./filebeat setup -e -strict.perms=false -E setup.kibana.host=kibana:5601
```

For sending container logs to kibana: 
```
logging:
  driver: gelf
  options:
    gelf-address: udp://localhost:12201
```

Login to pgadmin and run the SQL in `databases.sql`

SSH to sentry and run `./entrypoint.sh upgrade`

Login to cloud and run migration from admin.

ID for Grafana dashboard: 1860. 

GEMA server and agent need docker build.

## Localhost resolution
Use DNSMASQ for wildcard localhost subdomain resolution, like *.localhost.
Add to /etc/dnsmasq.conf:

```
address=/.localhost/127.0.0.1
```

## Running

Setup sysctl map count if needed. (see node-init.sh)  
Setup main node label.  

```
docker node update --label-add category=main <HOSTNAME>
```

Afterwards:  

```
docker swarm init
./deploy.sh
```

## Chat

```
mkdir -pv ./chat/volumes/app/mattermost/{data,logs,config,plugins}
sudo chown -R 2000:2000 ./chat/volumes/app/mattermost/
```

## TODO

* Internal proxy auth, create dashboard for creating users, groups. These users work for the HQ and the internal services, limit access to kibana. In dash show service configs.
* Dashboard with machine learning analytics.
* Dashboard for Service Redeployments (important), would be a release manager, you can even choose image version to deploy and see current version.
* Nginx cache purge in dashboard
* Cluster redis, postgres, prometheus and elasticsearch.
* Support replication for better performance of some services: ElasticSearch, Prometheus, Redis, PGSQL, GEMA Server
* Node communication/synchronization can be achieved through DSM solution using Redis as a message broker, like events.
* Go client to use Redis as special per-node replication. (Redis multimaster?)
* GEMA dashboard, geminis.io, live ops, reports generation, ML, message broadcast through WS, push notifications, service routings, let client have a user
* CORS?
* Apache Spark?
* H20.ai
* LogSkip internal pages with yml setting
* Use Hazelcast for PGSQL and Redis cache
* Add specific versioning for every service to prevent automatic upgrades.
* Use internal docker registry?
* Support service route config modification on the fly. Maybe use JSON (like consul) when config becomes too complex.
* Add Elastic Curator.
* Filebeat for PGSQL and Redis logs.
* Speed up ElasticSearch with Redis?


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

These tests should be performed through CI in Host machine container (git pull from container?), so as to reuse it in every repo.
