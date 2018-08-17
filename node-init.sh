apt-get update
apt-get -y install curl git fail2ban

curl -fsSL https://get.docker.com -o get-docker.sh
sh get-docker.sh

docker swarm init

docker network create --diver overlay gema
