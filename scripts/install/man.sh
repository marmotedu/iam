#!/usr/bin/env bash

# Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
# Use of this source code is governed by a MIT style
# license that can be found in the LICENSE file.


# The root of the build/dist directory
IAM_ROOT=$(dirname "${BASH_SOURCE[0]}")/../..
[[ -z ${COMMON_SOURCED} ]] && source ${IAM_ROOT}/scripts/install/common.sh

# 安装后打印必要的信息
function iam::man::info() {
cat << EOF
use: man iam-apiserver to see iam-apiserver help
EOF
}

# 安装
function iam::man::install()
{
  pushd ${IAM_ROOT}

  # 1. 生成各个组件的 man1 文件
  ${IAM_ROOT}/scripts/update-generated-docs.sh
  iam::common::sudo "cp docs/man/man1/* /usr/share/man/man1/"
  iam::man::status || return 1
  iam::man::info

  iam::log::info "install iam-apiserver successfully"
  popd
}

# 卸载
function iam::man::uninstall()
{
  set +o errexit
  iam::common::sudo "rm -f /usr/share/man/man1/iam-*"
  set -o errexit
  iam::log::info "uninstall iam man pages successfully"
}

# 状态检查
function iam::man::status()
{
  ls /usr/share/man/man1/iam-* &>/dev/null || {
    iam::log::error "iam man files not exist, maybe not installed properly"
    return 1
  }
}

if [[ "$*" =~ iam::man:: ]];then
  eval $*
fi
