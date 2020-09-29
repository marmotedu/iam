# 创建iam-apiserver证书和私钥

## 创建 iam-apiserver 证书和私钥

创建证书签名请求：

``` bash
$ cd $HOME/marmotedu/work
$ source $HOME/marmotedu/work/environment.sh
$ cat > iam-csr.json <<EOF
{
  "CN": "iam-apiserver",
  "hosts": [
    "127.0.0.1",
    "${IAM_APISERVER_HOST}"
  ],
  "key": {
    "algo": "rsa",
    "size": 2048
  },
  "names": [
    {
      "C": "CN",
      "ST": "BeiJing",
      "L": "BeiJing",
      "O": "iam",
      "OU": "marmotedu"
    }
  ]
}
EOF
```
+ hosts 字段指定授权使用该证书的 **IP 和域名列表**，这里列出了 iam-apiserver 节点 IP.

生成证书和私钥：

``` bash
$ cfssl gencert -ca=$HOME/marmotedu/work/ca.pem \
  -ca-key=$HOME/marmotedu/work/ca-key.pem \
  -config=$HOME/marmotedu/work/ca-config.json \
  -profile=iam iam-csr.json | cfssljson -bare iam 
$ ls iam*pem
iam-key.pem  iam.pem
```

将生成的证书和私钥文件拷贝到所有 iam-apiserver 节点：

``` bash
cd $HOME/marmotedu/work
source $HOME/marmotedu/work/environment.sh
for node_ip in ${NODE_IPS[@]}
  do
    echo ">>> ${node_ip}"
    ssh root@${node_ip} "mkdir -p /etc/iam/cert"
    scp iam-apiserverf*.pem root@${node_ip}:/etc/iam/cert/
  done
```
