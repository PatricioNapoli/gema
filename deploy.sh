echo "Deploying GEMA..."

docker stack deploy --compose-file=gema-compose.yml gema
docker stack deploy --compose-file=portainer-compose.yml portainer
docker stack deploy --compose-file=nginx-compose.yml nginx

echo "GEMA deployed."
