# wechat backend


# 使用方法
- 安装etcd
```
  docker run \
  -p 2379:2379 \
  -p 2380:2380 \
  --restart=always \
  -d \
  -v /Users/puper/dev/docker/etcd/data:/etcd-data \
  --name etcd \
  quay.io/coreos/etcd:v3.4.5 \
  /usr/local/bin/etcd \
  --name s1 \
  --data-dir /etcd-data \
  --listen-client-urls http://0.0.0.0:2379 \
  --advertise-client-urls http://0.0.0.0:2379 \
  --listen-peer-urls http://0.0.0.0:2380 \
  --initial-advertise-peer-urls http://0.0.0.0:2380 \
  --initial-cluster s1=http://0.0.0.0:2380 \
  --initial-cluster-token tkn \
  --initial-cluster-state new \
  --log-level info \
  --logger zap \
  --log-outputs stderr
  ```

- 启动
    - go run main.go serve --config=config/config.toml



## api docs
* [base](docs/apis/base.md)
* [app](docs/apis/app.md)
* [accesstoken](docs/apis/accesstoken.md)
* [oauth2](docs/apis/oauth2.md)