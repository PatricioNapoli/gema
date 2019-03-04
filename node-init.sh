#!/bin/sh

apt-get update
apt-get -y install curl git fail2ban

curl -fsSL https://get.docker.com -o get-docker.sh
sh get-docker.sh

sudo sysctl -w vm.max_map_count=262144
echo 'vm.max_map_count=262144' | sudo tee --append /etc/sysctl.conf

docker swarm init
docker node update --label-add category=main pond
