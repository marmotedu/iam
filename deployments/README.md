# iam应用容器化部署指南

## 依赖检查

- Kubernetes: `>= 1.16.0-0`
- Helm: `>= 3.0`

假设 iam项目根目录路径为 `IAM_ROOT`

进入iam项目根目录

```bash
$ cd ${IAM_ROOT}
```

## 容器化安装

具体安装步骤如下：

1) 生成配置文件

```bash
$ export IAM_APISERVER_INSECURE_BIND_ADDRESS=0.0.0.0
$ export IAM_AUTHZ_SERVER_INSECURE_BIND_ADDRESS=0.0.0.0
$ ./scripts/genconfig.sh scripts/install/environment.sh configs/iam-apiserver.yaml > deployments/iam/configs/iam-apiserver.yaml
$ ./scripts/genconfig.sh scripts/install/environment.sh configs/iam-authz-server.yaml > deployments/iam/configs/iam-authz-server.yaml
$ ./scripts/genconfig.sh scripts/install/environment.sh configs/iam-pump.yaml > deployments/iam/configs/iam-pump.yaml
```

```bash
$ kubectl -n iam create configmap iam --from-file=/etc/iam/
$ kubectl create configmap iam-cert --from-file=/etc/iam/cert
```

2) 使用Helm模板生成部署yaml文件: `iam.yaml`

```bash
$ helm template deployments/iam > deployments/iam.yaml
```

3) 安装iam应用

```bash
$ kubectl -n iam apply -f deployments/iam.yaml
```

4) 检查安装是否成功

**检查iam-apiserver**

```bash
$ export IAM_APISERVER_HOST=x.x.x.x
$ export IAM_APISERVER_INSECURE_BIND_PORT=30080
$ ./scripts/install/test.sh iam::test::apiserver
```

**检查iam-authz-server**

```bash
$ export IAM_APISERVER_HOST=x.x.x.x
$ export IAM_APISERVER_INSECURE_BIND_PORT=30080
$ ./scripts/install/test.sh iam::test::authzserver
```

**检查iam-pump**

```bash
$ export IAM_APISERVER_HOST=x.x.x.x
$ export IAM_APISERVER_INSECURE_BIND_PORT=30080
$ ./scripts/install/test.sh iam::test::pump
```

## Helm安装

1) 生成配置文件

```bash
$ export IAM_APISERVER_INSECURE_BIND_ADDRESS=0.0.0.0
$ export IAM_AUTHZ_SERVER_INSECURE_BIND_ADDRESS=0.0.0.0
$ ./scripts/genconfig.sh scripts/install/environment.sh configs/iam-apiserver.yaml > deployments/iam/configs/iam-apiserver.yaml
$ ./scripts/genconfig.sh scripts/install/environment.sh configs/iam-authz-server.yaml > deployments/iam/configs/iam-authz-server.yaml
$ ./scripts/genconfig.sh scripts/install/environment.sh configs/iam-pump.yaml > deployments/iam/configs/iam-pump.yaml
```

2) Helm install

```bash
$ helm install iam deployments/iam
```

3) 测试
