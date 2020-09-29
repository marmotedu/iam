# 创建 iam-authz-server 证书和私钥

创建证书签名请求：

``` bash
$ cd $HOME/marmotedu/work
$ cat > iam-authz-server-csr.json <<EOF
{
    "CN": "iam-authz-server",
    "key": {
        "algo": "rsa",
        "size": 2048
    },
    "hosts": [
      "127.0.0.1",
      "${IAM_AUTHZSERVER_HOST}"
    ],
    "names": [
      {
        "C": "CN",
        "ST": "BeiJing",
        "L": "BeiJing",
        "O": "iam-authz-server",
        "OU": "marmotedu"
      }
    ]
}
EOF
```
+ hosts 列表包含**所有** iam-authz-server 节点 IP；

生成证书和私钥：

``` bash
$ cd $HOME/marmotedu/work
$ cfssl gencert -ca=$HOME/marmotedu/work/ca.pem \
  -ca-key=$HOME/marmotedu/work/ca-key.pem \
  -config=$HOME/marmotedu/work/ca-config.json \
  -profile=iam iam-authz-server-csr.json | cfssljson -bare iam-authz-server
$ ls iam-authz-server*pem
iam-authz-server-key.pem  iam-authz-server.pem
```

将生成的证书和私钥分发到所有 master 节点：

``` bash
cd /opt/k8s/work
source /opt/k8s/bin/environment.sh
for node_ip in ${NODE_IPS[@]}
  do
    echo ">>> ${node_ip}"
    scp iam-authz-server*.pem root@${node_ip}:/etc/kubernetes/cert/
  done
```
