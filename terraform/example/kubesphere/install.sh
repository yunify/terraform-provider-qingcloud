curl -O -k https://kubernetes.pek3b.qingstor.com/tools/kubekey/kk
chmod +x kk
yum install -y vim openssl socat conntrack ipset
echo -e 'yes\n' | /root/kk create cluster --with-kubesphere
