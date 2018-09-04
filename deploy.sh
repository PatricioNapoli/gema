echo "Deploying GEMA..."

docker stack deploy --compose-file=cloud-compose.yml cloud
docker stack deploy --compose-file=elk-compose.yml elk
docker stack deploy --compose-file=sql-compose.yml sql
docker stack deploy --compose-file=nosql-compose.yml nosql
docker stack deploy --compose-file=sentry-compose.yml sentry
docker stack deploy --compose-file=metrics-compose.yml metrics
docker stack deploy --compose-file=logstash-compose.yml logstash

docker stack deploy --compose-file=portainer-compose.yml portainer
docker stack deploy --compose-file=nginx-compose.yml nginx

echo "GEMA deployed."
