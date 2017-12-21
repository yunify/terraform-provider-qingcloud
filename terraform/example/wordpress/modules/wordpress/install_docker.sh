#!/bin/bash
curl -fsSL get.docker.com -o get-docker.sh 
sh get-docker.sh --mirror Aliyun
cat >> /etc/docker/daemon.json <<EOF
      {
        "registry-mirrors": ["https://registry.docker-cn.com"]
      }
EOF
systemctl start docker
