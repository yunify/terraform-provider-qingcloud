#!/bin/bash 
curl -fsSL get.docker.com -o get-docker.sh 
sh get-docker.sh --mirror Aliyun 
mkdir /etc/docker
cat >> /etc/docker/daemon.json <<EOF
{
  "registry-mirrors": ["https://registry.docker-cn.com"] 
}
EOF
systemctl restart docker
