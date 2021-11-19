# Systemd 配置、安装和启动

- [Systemd 配置、安装和启动](#systemd-配置安装和启动)
	- [1. 前置操作（需要 root 权限）](#前置操作需要-root-权限)
	- [1. 创建 iam-apiserver systemd unit 模板文件](#创建-iam-apiserver-systemd-unit-模板文件)
	- [3. 创建 iam-authz-server systemd unit 模板文件](#创建-iam-authz-server-systemd-unit-模板文件)
	- [4. 复制 systemd unit 模板文件到 sysmted 配置目录(需要有root权限)](#复制-systemd-unit-模板文件到-sysmted-配置目录需要有root权限)
	- [5. 启动 systemd 服务](#启动-systemd-服务)

## 1. 前置操作（需要 root 权限）

1. 根据注释配置 `environment.sh`

2. 创建 data 目录 

```
mkdir -p ${IAM_DATA_DIR}/{iam-apiserver,iam-authz-server}
```

3. 创建 bin 目录，并将 `iam-apiserver` 和 `iam-authz-server` 可执行文件复制过去

```bash
source ./environment.sh
mkdir -p ${IAM_INSTALL_DIR}/bin
cp iam-apiserver iam-authz-server ${IAM_INSTALL_DIR}/bin
```

4. 将 `iam-apiserver` 和 `iam-authz-server` 配置文件拷贝到 `${IAM_CONFIG_DIR}` 目录下

```bash
mkdir -p ${IAM_CONFIG_DIR}
cp iam-apiserver.yaml iam-authz-server.yaml ${IAM_CONFIG_DIR}
```

## 2. 创建 iam-apiserver systemd unit 模板文件

执行如下 shell 脚本生成 `iam-apiserver.service.template`

```bash
source ./environment.sh
cat > iam-apiserver.service.template <<EOF
[Unit]
Description=IAM APIServer
Documentation=https://github.com/marmotedu/iam/blob/master/init/README.md

[Service]
WorkingDirectory=${IAM_DATA_DIR}/iam-apiserver
ExecStart=${IAM_INSTALL_DIR}/bin/iam-apiserver --apiconfig=${IAM_CONFIG_DIR}/iam-apiserver.yaml
Restart=always
RestartSec=5
StartLimitInterval=0

[Install]
WantedBy=multi-user.target
EOF
```

## 3. 创建 iam-authz-server systemd unit 模板文件

执行如下 shell 脚本生成 `iam-authz-server.service.template`

```bash
source ./environment.sh
cat > iam-authz-server.service.template <<EOF
[Unit]
Description=IAM AuthzServer
Documentation=https://github.com/marmotedu/iam/blob/master/init/README.md

[Service]
WorkingDirectory=${IAM_DATA_DIR}/iam-authz-server
ExecStart=${IAM_INSTALL_DIR}/bin/iam-authz-server --authzconfig=${IAM_CONFIG_DIR}/iam-authz-server.yaml
Restart=always
RestartSec=5
StartLimitInterval=0

[Install]
WantedBy=multi-user.target
EOF
```

## 4. 创建 iam-pump systemd unit 模板文件

执行如下 shell 脚本生成 `iam-pump.service.template`

```bash
source ./environment.sh
cat > iam-pump.service.template <<EOF
[Unit]
Description=IAM Pump Server
Documentation=https://github.com/marmotedu/iam/blob/master/init/README.md

[Service]
WorkingDirectory=${IAM_DATA_DIR}/iam-pump
ExecStart=${IAM_INSTALL_DIR}/bin/iam-pump --authzconfig=${IAM_CONFIG_DIR}/iam-pump.yaml
Restart=always
RestartSec=5
StartLimitInterval=0

[Install]
WantedBy=multi-user.target
EOF
```

## 5. 创建 iam-watcher systemd unit 模板文件

执行如下 shell 脚本生成 `iam-watcher.service.template`

```bash
source ./environment.sh
cat > iam-watcher.service.template <<EOF
[Unit]
Description=IAM Watcher Server
Documentation=https://github.com/marmotedu/iam/blob/master/init/README.md

[Service]
WorkingDirectory=${IAM_DATA_DIR}/iam-watcher
ExecStart=${IAM_INSTALL_DIR}/bin/iam-watcher --authzconfig=${IAM_CONFIG_DIR}/iam-watcher.yaml
Restart=always
RestartSec=5
StartLimitInterval=0

[Install]
WantedBy=multi-user.target
EOF
```

## 6. 复制 systemd unit 模板文件到 sysmted 配置目录(需要有root权限)

```bash
cp iam-apiserver.service.template /etc/systemd/system/iam-apiserver.service
cp iam-authz-server.service.template /etc/systemd/system/iam-authz-server.service
cp iam-pump.service.template /etc/systemd/system/iam-pump.service
cp iam-watcher.service.template /etc/systemd/system/iam-watcher.service
```

## 7. 启动 systemd 服务

```bash
systemctl daemon-reload && systemctl enable iam-apiserver && systemctl restart iam-apiserver
systemctl daemon-reload && systemctl enable iam-authz-server && systemctl restart iam-authz-server
systemctl daemon-reload && systemctl enable iam-pump && systemctl restart iam-pump
systemctl daemon-reload && systemctl enable iam-watcher && systemctl restart iam-watcher
```
