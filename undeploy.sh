echo "Deleting GEMA..."

docker stack rm nginx
docker stack rm portainer
docker stack rm gema

echo "GEMA deleted."
