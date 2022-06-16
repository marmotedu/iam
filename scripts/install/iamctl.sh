#!/usr/bin/env bash

# Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
# Use of this source code is governed by a MIT style
# license that can be found in the LICENSE file.


# The root of the build/dist directory
IAM_ROOT=$(dirname "${BASH_SOURCE[0]}")/../..
[[ -z ${COMMON_SOURCED} ]] && source ${IAM_ROOT}/scripts/install/common.sh

# 安装后打印必要的信息
function iam::iamctl::info() {
cat << EOF
iamctl test command: iamctl user list
EOF
}

# 安装
function iam::iamctl::install()
{
  pushd ${IAM_ROOT}

  # 1. 生成并安装 CA 证书和私钥
  ./scripts/gencerts.sh generate-iam-cert ${LOCAL_OUTPUT_ROOT}/cert
  iam::common::sudo "cp ${LOCAL_OUTPUT_ROOT}/cert/ca* ${IAM_CONFIG_DIR}/cert"

  ./scripts/gencerts.sh generate-iam-cert ${LOCAL_OUTPUT_ROOT}/cert admin
  #iam::common::sudo "cp ${LOCAL_OUTPUT_ROOT}/cert/admin*pem ${IAM_CONFIG_DIR}/cert"
  cert_dir=$(dirname ${CONFIG_USER_CLIENT_CERTIFICATE})
  key_dir=$(dirname ${CONFIG_USER_CLIENT_KEY})
  mkdir -p ${cert_dir} ${key_dir}
  cp ${LOCAL_OUTPUT_ROOT}/cert/admin.pem ${CONFIG_USER_CLIENT_CERTIFICATE}
  cp ${LOCAL_OUTPUT_ROOT}/cert/admin-key.pem ${CONFIG_USER_CLIENT_KEY}

  # 2. 构建 iamctl
  make build BINS=iamctl
  cp ${LOCAL_OUTPUT_ROOT}/platforms/linux/amd64/iamctl $HOME/bin/

  # 3.  生成并安装 iamctl 的配置文件（iamctl.yaml）
  mkdir -p $HOME/.iam
  ./scripts/genconfig.sh ${ENV_FILE} configs/iamctl.yaml > $HOME/.iam/iamctl.yaml
  iam::iamctl::status || return 1
  iam::iamctl::info

  iam::log::info "install iamctl successfully"
  popd
}

# 卸载
function iam::iamctl::uninstall()
{
  set +o errexit
  rm -f $HOME/bin/iamctl
  rm -f $HOME/.iam/iamctl.yaml
  #iam::common::sudo "rm -f ${IAM_CONFIG_DIR}/cert/admin*pem"
  rm -f ${CONFIG_USER_CLIENT_CERTIFICATE}
  rm -f ${CONFIG_USER_CLIENT_KEY}
  set -o errexit

  iam::log::info "uninstall iamctl successfully"
}

# 状态检查
function iam::iamctl::status()
{
  iamctl user list | grep -q admin || {
   iam::log::error "cannot list user, iamctl maybe not installed properly"
   return 1
  }

 if echo | telnet ${IAM_APISERVER_HOST} ${IAM_APISERVER_INSECURE_BIND_PORT} 2>&1|grep refused &>/dev/null;then
   iam::log::error "cannot access insecure port, iamctl maybe not startup"
   return 1
 fi
}

if [[ "$*" =~ iam::iamctl:: ]];then
  eval $*
fi
