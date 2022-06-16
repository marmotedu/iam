#!/usr/bin/env bash

# Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
# Use of this source code is governed by a MIT style
# license that can be found in the LICENSE file.


# The root of the build/dist directory
IAM_ROOT=$(dirname "${BASH_SOURCE[0]}")/../..

[[ -z ${COMMON_SOURCED} ]] && source ${IAM_ROOT}/scripts/install/common.sh

# 安装后打印必要的信息
function iam::mariadb::info() {
cat << EOF
MariaDB Login: mysql -h127.0.0.1 -u${MARIADB_ADMIN_USERNAME} -p'${MARIADB_ADMIN_PASSWORD}'
EOF
}

# 安装
function iam::mariadb::install()
{
  # 1. 配置 MariaDB 10.5 Yum 源
  echo ${LINUX_PASSWORD} | sudo -S bash -c "cat << 'EOF' > /etc/yum.repos.d/mariadb-10.5.repo
# MariaDB 10.5 CentOS repository list - created 2020-10-23 01:54 UTC
# http://downloads.mariadb.org/mariadb/repositories/
[mariadb]
name = MariaDB
baseurl = https://mirrors.aliyun.com/mariadb/yum/10.5/centos8-amd64/
module_hotfixes=1
gpgkey=https://yum.mariadb.org/RPM-GPG-KEY-MariaDB
gpgcheck=0
EOF"

  # 2. 安装 MariaDB 和 MariaDB 客户端
  iam::common::sudo "yum -y install MariaDB-server MariaDB-client"

  # 3. 启动 MariaDB，并设置开机启动
  iam::common::sudo "systemctl enable mariadb"
  iam::common::sudo "systemctl start mariadb"

  # 4. 设置 root 初始密码
  iam::common::sudo "mysqladmin -u${MARIADB_ADMIN_USERNAME} password ${MARIADB_ADMIN_PASSWORD}"

  iam::mariadb::status || return 1
  iam::mariadb::info
  iam::log::info "install MariaDB successfully"
}

# 卸载
function iam::mariadb::uninstall()
{
  set +o errexit
  iam::common::sudo "systemctl stop mariadb"
  iam::common::sudo "systemctl disable mariadb"
  iam::common::sudo "yum -y remove MariaDB-server MariaDB-client"
  iam::common::sudo "rm -rf /var/lib/mysql"
  iam::common::sudo "rm -f /etc/yum.repos.d/mariadb-10.5.repo"
  set -o errexit
  iam::log::info "uninstall MariaDB successfully"
}

# 状态检查
function iam::mariadb::status()
{
  # 查看 mariadb 运行状态，如果输出中包含 active (running) 字样说明 mariadb 成功启动。
  systemctl status mariadb |grep -q 'active' || {
    iam::log::error "mariadb failed to start, maybe not installed properly"
    return 1
  }

  mysql -u${MARIADB_ADMIN_USERNAME} -p${MARIADB_ADMIN_PASSWORD} -e quit &>/dev/null || {
    iam::log::error "can not login with root, mariadb maybe not initialized properly"
    return 1
  }
}

if [[ "$*" =~ iam::mariadb:: ]];then
  eval $*
fi
