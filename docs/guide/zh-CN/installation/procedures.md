# 安装步骤

[部署架构](./architecture.md)

## 1. 需求检查 & 依赖安装

请参考：[需求检查](./requirement.md)

## 2. 代码包下载

```bash
git clone https://github.com/marmotedu/iam
```

## 3. 编译

```bash
cd iam
make
```

## 4. 生成配置

1. 准备 `environment.sh`

```bash
cat scripts/environment.sh
#!/usr/bin/bash

# iam 各组件数据目录
export IAM_DATA_DIR="/data/iam"

# iam 配置存放目录
export IAM_CONFIG_DIR="/opt/iam"

# MySQL配置信息
export MYSQL_HOST=127.0.0.1
export MYSQL_PORT=3306
export MYSQL_USERNAME=iam
export MYSQL_PASSWORD=iam1234

# Redis配置信息
export REDIS_HOST=127.0.0.1
export REDIS_PORT=6379
export REDIS_PASSWORD=redis


# 端口指定
# iam-apiserver insecure(http)端口
IAM_APISERVER_HOST=127.0.0.1
IAM_APISERVER_HTTP_PORT=8080
# iam-apiserver secure(https)端口
IAM_APISERVER_HTTPS_PORT=8443
# iam-apiserver grpc端口
IAM_APISERVER_GRPC_PORT=8081

IAM_AUTHZSERVER_HOST=127.0.0.1
# iam-authz-server insecure(http)端口
IAM_AUTHZSERVER_HTTP_PORT=9090
# iam-authz-server secure(https)端口
IAM_AUTHZSERVER_HTTPS_PORT=9443
```

> 根据需要进行配置，通常需要修改 Mysql 和 Redis 服务的地址和端口

2. 创建配置文件存放目录 & 载入环境变量

```bash
mkdir $HOME/.iam
source scripts/environment.sh
```

3. 生成 `iam-apiserver.yaml` 配置文件

```bash
cat > $HOME/.iam/iam-apiserver.yaml << EOF
# iam-apiserver 全配置

# RESTful服务配置
server:
    max-ping-count: 3                                                # http 服务启动后，自检尝试次数
    middlewares: recovery,secure,options,nocache,cors,requestid,dump # 加载的 gin 中间件
    healthz: true                                                    # 是否开启健康检查，如果开启会安装 /healthz 路由
    mode: debug                                                      # server mode: release, debug, test

# GRPC服务配置
grpc:
  bind-address:  0.0.0.0                                              # grpc 服务绑定地址
  bind-port: ${IAM_APISERVER_GRPC_PORT}                               # grpc 服务绑定端口

# HTTP配置
insecure:
    bind-address: 0.0.0.0
    bind-port: ${IAM_APISERVER_HTTP_PORT}

# HTTPS配置
secure:
    bind-address: 0.0.0.0
    bind-port: ${IAM_APISERVER_HTTPS_PORT}
    tls:
        cert-dir: config
        pair-name: cccc
        cert-key:
            cert-file: conf/server.crt
            private-key-file: conf/server.key

# Mysql数据库相关配置
mysql:
  host: ${MYSQL_HOST}:${MYSQL_PORT}                                   # mysql 机器 ip 和端口
  username: ${MYSQL_USERNAME}                                         # mysql 用户名(建议授权最小权限集)
  password: ${MYSQL_PASSWORD}                                         # mysql 用户密码
  database: iam                                                       # iam 系统所用的数据库名
  max-idle-connections: 100
  max-open-connections: 100
  max-connection-life-time: 10s
  log-mode: true

# Redis配置
redis:
  host: ${REDIS_HOST}                                                 # redis 地址
  port: ${REDIS_PORT}                                                 # redis 端口
  password: ${REDIS_PASSWORD}                                         # redis 密码
  #master-name:
  #username:
  #database:
  #optimisation-max-idle:
  #ooptimisation-max-active:
  #timeout:
  #enable-cluster:
  #ssl-insecure-skip-verify:

# JWT配置
jwt:
  realm: iam jwt                                                      # jwt 标识
  key: dfVpOK8LZeJLZHYmHdb1VdyRrACKpqoo                               # 服务端密钥
  timeout: 24h                                                        # token 过期时间(小时)
  max-refresh: 24h                                                    # token 更新时间(小时)

log:
    level: info                                                       # debug, info, warn, error, dpanic, panic, fatal
    format: console                                                   # console, json
    disable-color: false                                              # 是否开启颜色输出，true:是，false:否
    enable-caller: true
    output-paths: /tmp/iam.log,stdout                                 # 多个输出，逗号分开。stdout：标准输出，
    #error-output-paths

feature:
  enable-metrics: true                                                # 开启 metrics, router:  /metrics
  profiling: true                                                     # 开启 profiling, router 'host:port/debug/pprof/'
EOF
```
4. 生成 i`am-authz-config.yaml` 配置文件

```bash
cat > $HOME/.iam/iam-authz-server.yaml << EOF
# iam-authz-server 全配置

# IAM rpc 服务地址
rpcserver: ${IAM_APISERVER_HOST}:${IAM_APISERVER_GRPC_PORT}

# RESTful服务配置
server:
  max-ping-count: 3                                                   # http 服务启动后，自检尝试次数
  middlewares: recovery,secure,options,nocache,cors,requestid,dump    # 加载的 gin 中间件
  healthz: true                                                       # 是否开启健康检查，如果开启会安装 /healthz 路由
  mode: debug                                                         # server mode: release, debug, test

# HTTP配置
insecure:
    bind-address: 0.0.0.0
    bind-port: ${IAM_AUTHZSERVER_HTTP_PORT}

# HTTPS配置
secure:
    bind-address: 0.0.0.0
    bind-port: ${IAM_AUTHZSERVER_HTTPS_PORT}
    tls:
        cert-dir: config
        pair-name: cccc
        cert-key:
            cert-file: conf/server.crt
            private-key-file: conf/server.key

# Redis配置
redis:
  host: ${REDIS_HOST}                                                 # redis地址
  port: ${REDIS_PORT}                                                 # redis端口
  password: ${REDIS_PASSWORD}                                         # redis密码
  #master-name:
  #username:
  #database:
  #optimisation-max-idle:
  #ooptimisation-max-active:
  #timeout:
  #enable-cluster:
  #ssl-insecure-skip-verify:

log:
    level: info                                                       # debug, info, warn, error, dpanic, panic, fatal
    format: console                                                   # console, json
    disable-color: false                                              # 是否开启颜色输出，true:是，false:否
    enable-caller: true
    output-paths: /tmp/iam.log,stdout                                 # 多个输出，逗号分开。stdout：标准输出，
    #error-output-paths

feature:
  enable-metrics: true                                                # 开启 metrics, router:  /metrics
  profiling: true                                                     # 开启 profiling, router 'host:port/debug/pprof/'
EOF
```

5. 生成 `config` 配置文件

```bash
cat > $HOME/.iam/config << EOF
apiVersion: v1
user:
  #token:
  username: admin                                                     # 用户名
  password: Admin@2020                                                # 密码
  #secret-id:
  #secret-key:
  #client-certificate:
  #client-certificate-data:
  #client-key:
  #client-key-data:
server:
  address: http://${IAM_APISERVER_HOST}:${IAM_APISERVER_HTTP_PORT}    # iam api-server 地址
  timeout: 10s                                                        # 请求 api-server 超时时间
  #max-retries:                                                       # 最大重试次数
  #retry-interval:                                                    # 重试间隔
  #tls-server-name:
  #insecure-skip-tls-verify:
  #certificate-authority:
  #certificate-authority-data:
EOF
```

### 4.1 初始化数据库

```bash
source environment.sh
mysql -h ${MYSQL_HOST} -P${MYSQL_PORT} -u${MYSQL_USERNAME} -p${MYSQL_PASSWORD}

MySQL [(none)]> source scripts/iam.20200711.sql
```

> 初始化数据库后，会自动创建一个：admin (用户名)，Admin@2020 (密码)的管理员账户

## 5. 部署 (需要有 root 权限)

1. 安装可执行文件到指定目录
```bash
source scripts/environment.sh

cp _output/platforms/linux/amd64/* ${IAM_CONFIG_DIR}/{bin}/
```

2. 启动 `iam-apiserver`

```bash
${IAM_CONFIG_DIR}/{bin}/iam-apiserver
```
3. 启动 `iam-authz-server`

```bash
${IAM_CONFIG_DIR}/{bin}/iam-authz-server
```

## 6. 测试

1. Health测试

```bash
# 请求
curl --request GET --url http://106.52.131.123:8080/healthz

# 返回
"ok"
```
2. 创建用户

```bash
# 请求
curl --request POST \
  --url http://konglingfei.com:8080/v1/users \
  --header 'content-type: application/json' \
  --data '{
	"metadata":{
		"name":"foo"
	},
	"nickname":"foo",
	"password":"Foo@2020",
	"email":"foo@foxmail.com",
	"phone":"1812884xxxx"
}'

# 返回
{
  "metadata": {
    "id": 30,
    "name": "foo",
    "createdAt": "2020-07-12 09:17:53",
    "updatedAt": "2020-07-12 09:17:53"
  },
  "nickname": "foo",
  "password": "$2a$10$l3Fso2F/8o60IOVXyaKgV.pr55tFSlx9kIY007CJ0c/hgRYhBwJAa",
  "email": "foo@foxmail.com",
  "phone": "1812884xxxx"
}
```

3. 登陆管理员账号

```bash
 # 请求
 curl --request POST \
  --url http://konglingfei.com:8080/login \
  --header 'content-type: application/json' \
  --data '{
	"name":"admin",
	"password":"Admin@2020"
}'

 # 返回
 {
  "expire": "2020-07-13T08:53:18+08:00",
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1OTQ2MDE1OTgsImlkZW50aXR5IjoiYWRtaW4iLCJpc3MiOiJpYW0tYXBpc2VydmVyIiwib3JpZ19pYXQiOjE1OTQ1MTUxOTgsInN1YiI6InVzZXIgb2YgaWFtLWFwaXNlcnZlciJ9.6mzvO32148Xax3y7Tj-9WspML0ATgGeUDZqr9x2OLI4"
}
```

> 返回的 token 用户设置后面请求的 `Authorization: Bearer`

4. 创建密钥

```bash
# 请求
curl --request POST \
  --url http://konglingfei.com:8080/v1/secrets \
  --header 'authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1OTQ2MDE1OTgsImlkZW50aXR5IjoiYWRtaW4iLCJpc3MiOiJpYW0tYXBpc2VydmVyIiwib3JpZ19pYXQiOjE1OTQ1MTUxOTgsInN1YiI6InVzZXIgb2YgaWFtLWFwaXNlcnZlciJ9.6mzvO32148Xax3y7Tj-9WspML0ATgGeUDZqr9x2OLI4' \
  --header 'content-type: application/json' \
  --data '{
	"metadata":{
		"name":"secret1"
	},
	"expires":0,
	"description":"admin secret"
}'

# 返回
{
  "metadata": {
    "id": 18,
    "name": "secret1",
    "createdAt": "2020-07-12 09:17:53",
    "updatedAt": "2020-07-12 09:17:53"
  },
  "username": "admin",
  "secretID": "LRhchcQgd6dzu5DhwERr3fFDsB8d35yqhLJL",
  "secretKey": "NX6hYwvdbSgq8nsJK8aEoRrMV38aqLFU",
  "expires": 0,
  "description": "admin secret"
}
```

4. 创建授权策略

```bash
# 请求
curl --request POST \
  --url http://konglingfei.com:8080/v1/policies \
  --header 'authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1OTQ2MDE1OTgsImlkZW50aXR5IjoiYWRtaW4iLCJpc3MiOiJpYW0tYXBpc2VydmVyIiwib3JpZ19pYXQiOjE1OTQ1MTUxOTgsInN1YiI6InVzZXIgb2YgaWFtLWFwaXNlcnZlciJ9.6mzvO32148Xax3y7Tj-9WspML0ATgGeUDZqr9x2OLI4' \
  --header 'content-type: application/json' \
  --data '{
	"metadata":{
		"name":"policy1"
	},
  "description": "One policy to rule them all.",
  "subjects": [
    "users:<peter|ken>",
    "users:maria",
    "groups:admins"
  ],
  "actions": [
    "delete",
    "<create|update>"
  ],
  "effect": "allow",
  "resources": [
    "resources:articles:<.*>",
    "resources:printer"
  ],
  "conditions": {
    "remoteIP": {
      "type": "CIDRCondition",
      "options": {
        "cidr": "192.168.0.1/16"
      }
    }
  }
}'

# 返回
{
  "metadata": {
    "id": 35,
    "name": "policy1",
    "createdAt": "2020-07-12 10:03:18",
    "updatedAt": "2020-07-12 10:03:18"
  },
  "username": "admin",
  "policy": {
    "id": "",
    "description": "One policy to rule them all.",
    "subjects": [
      "users:<peter|ken>",
      "users:maria",
      "groups:admins"
    ],
    "effect": "allow",
    "resources": [
      "resources:articles:<.*>",
      "resources:printer"
    ],
    "actions": [
      "delete",
      "<create|update>"
    ],
    "conditions": {
      "remoteIP": {
        "type": "CIDRCondition",
        "options": {
          "cidr": "192.168.0.1/16"
        }
      }
    },
    "meta": null
  }
}
```

5. 生成 Token

```bash
iamctl sign LRhchcQgd6dzu5DhwERr3fFDsB8d35yqhLJL NX6hYwvdbSgq8nsJK8aEoRrMV38aqLFU --timeout=2h --algorithm=HS256
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1OTQ1MjY4MDIsImlhdCI6MTU5NDUxOTYwMiwianRpIjoiTFJoY2hjUWdkNmR6dTVEaHdFUnIzZkZEc0I4ZDM1eXFoTEpMIn0.vYyPLSILfubHwtvCM9YSmO0PonFB-OeUWJhucZHBA7c
```

> 用第3步生成的 secretID 和 secretKey 作为输入参数

6. 请求 iam-authz-server进行权限验证

```bash
# 请求
curl --request POST \
  --url http://konglingfei.com:9090/v1/authz \
  --header 'authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1OTQ1MjY4MDIsImlhdCI6MTU5NDUxOTYwMiwianRpIjoiTFJoY2hjUWdkNmR6dTVEaHdFUnIzZkZEc0I4ZDM1eXFoTEpMIn0.vYyPLSILfubHwtvCM9YSmO0PonFB-OeUWJhucZHBA7c' \
  --header 'content-type: application/json' \
  --data '{
  "subject": "users:peter",
  "action": "delete",
  "resource": "resources:articles:ladon-introduction",
  "context": {
    "remoteIP": "192.168.0.5"
  }
}'

# 返回
{
  "allowed": true
}
```
