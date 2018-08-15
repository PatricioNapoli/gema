apt-get update
apt-get install curl

curl -fsSL https://get.docker.com -o get-docker.sh
sh get-docker.sh

docker swarm init
