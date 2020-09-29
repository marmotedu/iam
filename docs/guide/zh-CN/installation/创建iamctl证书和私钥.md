## 创建 iamctl admin 证书和私钥

iamctl 使用 https 协议与 iam-apiserver 进行安全通信，iam-apiserver 对 iamctl 请求包含的证书进行认证和授权。

iamctl 后续用于集群管理，所以这里创建具有**最高权限**的 admin 证书。

创建证书签名请求：

``` bash
$ cd $HOME/marmotedu/work
$ cat > admin-csr.json <<EOF
{
  "CN": "admin",
  "hosts": [],
  "key": {
    "algo": "rsa",
    "size": 2048
  },
  "names": [
    {
      "C": "CN",
      "ST": "BeiJing",
      "L": "BeiJing",
      "O": "iamctl",
      "OU": "marmotedu"
    }
  ]
}
EOF
```
+ 该证书只会被 iamctl 当做 client 证书使用，所以 `hosts` 字段为空；

生成证书和私钥：

``` bash
$ cd $HOME/marmotedu/work
$ cfssl gencert -ca=$HOME/marmotedu/work/ca.pem \
  -ca-key=$HOME/marmotedu/work/ca-key.pem \
  -config=$HOME/marmotedu/work/ca-config.json \
  -profile=iam admin-csr.json | cfssljson -bare admin
$ ls admin*
admin.csr  admin-csr.json  admin-key.pem  admin.pem
```
+ 忽略警告消息 `[WARNING] This certificate lacks a "hosts" field.`；

## 创建 iamconfig 文件


示例配置如下：

```yaml
apiVersion: v1
user:
  #token:
  username: admin                                                     # 用户名
  password: Admin@2020                                                # 密码
#secret-id: 
#secret-key: 
  client-certificate: /home/colin/marmotedu/work/admin.pem
  client-key: /home/colin/marmotedu/work/admin-key.pem
  #client-certificate-data:
  #client-key-data:
server:
  address: https://127.0.0.1:8443    # iam api-server 地址
  timeout: 10s                                                        # 请求 api-server 超时时间
  #max-retries:                                                       # 最大重试次数
  #retry-interval:                                                    # 重试间隔
  #tls-server-name:
  #insecure-skip-tls-verify:
  certificate-authority: /home/colin/marmotedu/work/ca.pem
  #certificate-authority-data:
```

## 创建 kubeconfig 文件

kubectl 使用 kubeconfig 文件访问 apiserver，该文件包含 kube-apiserver 的地址和认证信息（CA 证书和客户端证书）：

``` bash
cd /opt/k8s/work
source /opt/k8s/bin/environment.sh

# 设置集群参数
kubectl config set-cluster kubernetes \
  --certificate-authority=/opt/k8s/work/ca.pem \
  --embed-certs=true \
  --server=https://${NODE_IPS[0]}:6443 \
  --kubeconfig=kubectl.kubeconfig

# 设置客户端认证参数
kubectl config set-credentials admin \
  --client-certificate=/opt/k8s/work/admin.pem \
  --client-key=/opt/k8s/work/admin-key.pem \
  --embed-certs=true \
  --kubeconfig=kubectl.kubeconfig

# 设置上下文参数
kubectl config set-context kubernetes \
  --cluster=kubernetes \
  --user=admin \
  --kubeconfig=kubectl.kubeconfig

# 设置默认上下文
kubectl config use-context kubernetes --kubeconfig=kubectl.kubeconfig
```
+ `--certificate-authority`：验证 kube-apiserver 证书的根证书；
+ `--client-certificate`、`--client-key`：刚生成的 `admin` 证书和私钥，与 kube-apiserver https 通信时使用；
+ `--embed-certs=true`：将 ca.pem 和 admin.pem 证书内容嵌入到生成的 kubectl.kubeconfig 文件中(否则，写入的是证书文件路径，后续拷贝 kubeconfig 到其它机器时，还需要单独拷贝证书文件，不方便。)；
+ `--server`：指定 kube-apiserver 的地址，这里指向第一个节点上的服务；

## 分发 kubeconfig 文件

分发到所有使用 `kubectl` 命令的节点：

``` bash
cd /opt/k8s/work
source /opt/k8s/bin/environment.sh
for node_ip in ${NODE_IPS[@]}
  do
    echo ">>> ${node_ip}"
    ssh root@${node_ip} "mkdir -p ~/.kube"
    scp kubectl.kubeconfig root@${node_ip}:~/.kube/config
  done
```
