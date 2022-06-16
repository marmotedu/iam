#!/usr/bin/env bash

# Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
# Use of this source code is governed by a MIT style
# license that can be found in the LICENSE file.


# The root of the build/dist directory
IAM_ROOT=$(dirname "${BASH_SOURCE[0]}")/../..
[[ -z ${COMMON_SOURCED} ]] && source ${IAM_ROOT}/scripts/install/common.sh

# 安装后打印必要的信息
function iam::apiserver::info() {
cat << EOF
iam-apserver insecure listen on: ${IAM_APISERVER_HOST}:${IAM_APISERVER_INSECURE_BIND_PORT}
iam-apserver secure listen on: ${IAM_APISERVER_HOST}:${IAM_APISERVER_SECURE_BIND_PORT}
EOF
}

# 安装
function iam::apiserver::install()
{
  pushd ${IAM_ROOT}

  # 1. 生成 CA 证书和私钥
  ./scripts/gencerts.sh generate-iam-cert ${LOCAL_OUTPUT_ROOT}/cert
  iam::common::sudo "cp ${LOCAL_OUTPUT_ROOT}/cert/ca* ${IAM_CONFIG_DIR}/cert"

  ./scripts/gencerts.sh generate-iam-cert ${LOCAL_OUTPUT_ROOT}/cert iam-apiserver
  iam::common::sudo "cp ${LOCAL_OUTPUT_ROOT}/cert/iam-apiserver*pem ${IAM_CONFIG_DIR}/cert"

  # 2. 构建 iam-apiserver
  make build BINS=iam-apiserver
  iam::common::sudo "cp ${LOCAL_OUTPUT_ROOT}/platforms/linux/amd64/iam-apiserver ${IAM_INSTALL_DIR}/bin"

  # 3.  生成并安装 iam-apiserver 的配置文件（iam-apiserver.yaml）
  echo ${LINUX_PASSWORD} | sudo -S bash -c \
    "./scripts/genconfig.sh ${ENV_FILE} configs/iam-apiserver.yaml > ${IAM_CONFIG_DIR}/iam-apiserver.yaml"

  # 4. 创建并安装 iam-apiserver systemd unit 文件
  echo ${LINUX_PASSWORD} | sudo -S bash -c \
    "./scripts/genconfig.sh ${ENV_FILE} init/iam-apiserver.service > /etc/systemd/system/iam-apiserver.service"

  # 5. 启动 iam-apiserver 服务
  iam::common::sudo "systemctl daemon-reload"
  iam::common::sudo "systemctl restart iam-apiserver"
  iam::common::sudo "systemctl enable iam-apiserver"
  iam::apiserver::status || return 1
  iam::apiserver::info

  iam::log::info "install iam-apiserver successfully"
  popd
}

# 卸载
function iam::apiserver::uninstall()
{
  set +o errexit
  iam::common::sudo "systemctl stop iam-apiserver"
  iam::common::sudo "systemctl disable iam-apiserver"
  iam::common::sudo "rm -f ${IAM_INSTALL_DIR}/bin/iam-apiserver"
  iam::common::sudo "rm -f ${IAM_CONFIG_DIR}/iam-apiserver.yaml"
  iam::common::sudo "rm -f /etc/systemd/system/iam-apiserver.service"
  iam::common::sudo "rm -f ${IAM_CONFIG_DIR}/cert/iam-apiserver*pem"
  set -o errexit
  iam::log::info "uninstall iam-apiserver successfully"
}

# 状态检查
function iam::apiserver::status()
{
  # 查看 apiserver 运行状态，如果输出中包含 active (running) 字样说明 apiserver 成功启动。
  systemctl status iam-apiserver|grep -q 'active' || {
    iam::log::error "iam-apiserver failed to start, maybe not installed properly"
    return 1
  }

 if echo | telnet ${IAM_APISERVER_HOST} ${IAM_APISERVER_INSECURE_BIND_PORT} 2>&1|grep refused &>/dev/null;then
   iam::log::error "cannot access insecure port, iam-apiserver maybe not startup"
   return 1
 fi
}

if [[ "$*" =~ iam::apiserver:: ]];then
  eval $*
fi
