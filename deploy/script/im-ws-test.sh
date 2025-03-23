
reso_addr='registry.cn-hangzhou.aliyuncs.com/hnz-chat/im-ws-dev'
tag='latest'

pod_idb="192.168.232.12"

container_name='hnz-chat-im-ws-dev'

docker stop ${container_name}
docker rm ${container_name}
docker rmi ${reso_addr}:${tag}
docker pull ${reso_addr}:${tag}

# 启动新容器，配置 etcd 连接参数
docker run \
  -p 10090:10090 \
  -e POD_IP=${pod_idb} \
  --name=${container_name} \
  -d ${reso_addr}:${tag}

# 检查服务是否成功注册到etcd
sleep 5
echo "检查服务注册状态:"
docker exec etcd etcdctl get --prefix "im.ws"