#!/bin/sh

apt-get update
apt-get upgrade
apt-get -y install curl git fail2ban software-properties-common

curl -sL https://deb.nodesource.com/setup_11.x | sudo bash -
sudo apt-get install nodejs

curl -fsSL https://get.docker.com -o get-docker.sh
sh get-docker.sh && rm get-docker.sh

sudo sysctl -w vm.max_map_count=262144
echo 'vm.max_map_count=262144' | sudo tee --append /etc/sysctl.conf

# docker swarm init
# docker node update --label-add category=main $HOSTNAME

# Swap
dd if=/dev/zero of=/swapfile count=8192 bs=1M
chmod 600 /swapfile
mkswap /swapfile
swapon /swapfile
echo "/swapfile   none    swap    sw    0   0" >> /etc/fstab