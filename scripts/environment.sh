# Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
# Use of this source code is governed by a MIT style
# license that can be found in the LICENSE file.

#!/usr/bin/bash

: << EOF
环境变量名和configs/*.yaml中的配置项一一对应，例如：configs/iam-apiserver.yaml中有如下配置：
insecure:
    bind-address: 127.0.0.1

则环境变量名称为：IAM_APISERVER_INSECURE_BIND_ADDRESS。环境变量的具体含义也请参考configs/*.yaml中的注释。
EOF

# iam 配置
export IAM_DATA_DIR="/data/iam" # iam 各组件数据目录
export IAM_INSTALL_DIR="/opt/iam" # iam 安装文件存放目录
export IAM_CONFIG_DIR="/etc/iam" # iam 配置文件存放目录
export IAM_LOG_DIR="/var/log/iam" # iam 日志文件存放目录
export CA_FILE=${IAM_CONFIG_DIR}/cert/ca.pem # CA

# MySQL 配置信息
export MYSQL_HOST=127.0.0.1:3306        # MySQL主机地址
export MYSQL_USERNAME=iam               # 登陆用户
export MYSQL_PASSWORD=iam1234           # 登陆密码
export MYSQL_DATABASE=iam               # iam数据库

# Redis 配置信息
export REDIS_HOST=127.0.0.1             # Redis 主机地址
export REDIS_PORT=6379                  # Redis 监听端口
export REDIS_PASSWORD=iam1234 # Redis 密码

# MongoDB 配置
MONGODB_HOST=127.0.0.1
MONGODB_PORT=27017
MONGODB_USERNAME=iam
MONGODB_PASSWORD=iam1234

# iam-apiserver 配置
IAM_APISERVER_HOST=127.0.0.1 # iam-apiserver 部署机器 IP 地址
IAM_APISERVER_GRPC_BIND_ADDRESS=0.0.0.0
IAM_APISERVER_GRPC_BIND_PORT=8081
IAM_APISERVER_INSECURE_BIND_ADDRESS=127.0.0.1
IAM_APISERVER_INSECURE_BIND_PORT=8080
IAM_APISERVER_SECURE_BIND_ADDRESS=0.0.0.0
IAM_APISERVER_SECURE_BIND_PORT=8443
IAM_APISERVER_SECURE_TLS_CERT_KEY_CERT_FILE=${IAM_CONFIG_DIR}/cert/iam-apiserver.pem
IAM_APISERVER_SECURE_TLS_CERT_KEY_PRIVATE_KEY_FILE=${IAM_CONFIG_DIR}/cert/iam-apiserver-key.pem

# iam-authz-server 配置
IAM_AUTHZ_SERVER_HOST=127.0.0.1 # iam-authz-server 部署机器 IP 地址
IAM_AUTHZ_SERVER_INSECURE_BIND_ADDRESS=127.0.0.1
IAM_AUTHZ_SERVER_INSECURE_BIND_PORT=9090
IAM_AUTHZ_SERVER_SECURE_BIND_ADDRESS=0.0.0.0
IAM_AUTHZ_SERVER_SECURE_BIND_PORT=9443
IAM_AUTHZ_SERVER_SECURE_TLS_CERT_KEY_CERT_FILE=${IAM_CONFIG_DIR}/cert/iam-authz-server.pem
IAM_AUTHZ_SERVER_SECURE_TLS_CERT_KEY_PRIVATE_KEY_FILE=${IAM_CONFIG_DIR}/cert/iam-authz-server-key.pem
IAM_AUTHZ_SERVER_CLIENT_CA_FILE=${CA_FILE}
IAM_AUTHZ_SERVER_RPCSERVER=${IAM_APISERVER_HOST}:${IAM_APISERVER_GRPC_BIND_PORT}

# iam-pump 配置
IAM_PUMP_HOST=127.0.0.1 # iam-pump 部署机器 IP 地址
IAM_PUMP_COLLECTION_NAME=iam_analytics
IAM_PUMP_MONGO_URL=mongodb://${MONGODB_USERNAME}:${MONGODB_PASSWORD}@${MONGODB_HOST}:${MONGODB_PORT}/${IAM_PUMP_COLLECTION_NAME}?authSource=admin

# iamctl 配置
CONFIG_USER_USERNAME=admin
CONFIG_USER_PASSWORD=Admin@2020
CONFIG_USER_CLIENT_CERTIFICATE=${IAM_CONFIG_DIR}/cert/admin.pem
CONFIG_USER_CLIENT_KEY=${IAM_CONFIG_DIR}/cert/admin-key.pem
CONFIG_SERVER_ADDRESS=${IAM_APISERVER_HOST}:${IAM_APISERVER_SECURE_BIND_PORT}
CONFIG_SERVER_CERTIFICATE_AUTHORITY=${CA_FILE}
