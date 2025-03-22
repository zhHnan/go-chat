
reso_addr='registry.cn-hangzhou.aliyuncs.com/hnz-chat/im-api-dev'
tag='latest'

container_name='hnz-chat-im-api-dev'

docker stop ${container_name}
docker rm ${container_name}
docker rmi ${reso_addr}:${tag}
docker pull ${reso_addr}:${tag}

## 如果需要指定配置文件
## docker run -p 10001:8080 --network hnz-chat -v/hnz-chat/config/user-rpc:/user/conf/ --name=${container_name} -d ${reso_addr}:${tag}
#docker run -p 10000:10000 --name=${container_name} -d ${reso_addr}:${tag}
# 等待 etcd 启动（假设通过 docker-compose 管理）

# 启动新容器，配置 etcd 连接参数
docker run \
  -p 8882:8882 \
  --name=${container_name} \
  -d ${reso_addr}:${tag}

# 检查服务是否成功注册到etcd
sleep 5
echo "检查服务注册状态:"
docker exec etcd etcdctl get --prefix "im.api"