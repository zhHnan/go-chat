# shellcheck disable=SC2034
need_start_server_shell=(
# rpc
    user-rpc-test.sh
)

for i in ${need_start_server_shell[*]}; do
    chmod +x $i
    # shellcheck disable=SC1090
    ./$i
done

docker ps

docker exec -it etcd etcdctl get --prefix ""
