echo "Deploying GEMA..."

docker stack rm nginx
docker stack rm portainer
docker stack rm architecture

echo "Waiting for Docker cleanup.."
sleep 10

docker stack deploy --compose-file=architecture-compose.yml architecture
docker stack deploy --compose-file=portainer-compose.yml portainer
docker stack deploy --compose-file=nginx-compose.yml nginx

echo "GEMA deployed."
