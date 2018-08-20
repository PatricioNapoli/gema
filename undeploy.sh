echo "Deleting GEMA..."

docker stack rm nginx
docker stack rm portainer
docker stack rm architecture

echo "GEMA deleted."
