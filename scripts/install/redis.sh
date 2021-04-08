#!/usr/bin/env bash

# The root of the build/dist directory
IAM_ROOT=$(dirname "${BASH_SOURCE[0]}")/../..
[[ -z ${COMMON_SOURCED} ]] && source ${IAM_ROOT}/scripts/install/common.sh

# 安装后打印必要的信息
function iam::redis::info() {
cat << EOF
Redis Login: redis-cli --no-auth-warning -h ${REDIS_HOST} -p ${REDIS_PORT} -a '${REDIS_PASSWORD}'
EOF
}

# 安装
function iam::redis::install()
{
  # 1. 安装 Redis
  iam::common::sudo "yum -y install redis"

  # 2. 配置 Redis
  echo ${LINUX_PASSWORD} | sudo -S sed -i '/^daemonize/{s/no/yes/}' /etc/redis.conf
  echo ${LINUX_PASSWORD} | sudo -S sed -i '/^# bind 127.0.0.1/{s/# //}' /etc/redis.conf
  echo ${LINUX_PASSWORD} | sudo -S sed -i 's/^# requirepass.*$/requirepass '"${REDIS_PASSWORD}"'/' /etc/redis.conf
  echo ${LINUX_PASSWORD} | sudo -S sed -i '/^protected-mode/{s/yes/no/}' /etc/redis.conf

  # 3. 启动 Redis
  iam::common::sudo "redis-server /etc/redis.conf"

  # 4. 禁用防火墙
  iam::common::sudo "systemctl stop firewalld.service"
  iam::common::sudo "systemctl disable firewalld.service"

  iam::redis::status || return 1
  iam::redis::info
  iam::log::info "install Redis successfully"
}

# 卸载
function iam::redis::uninstall()
{
  set +o errexit
  iam::common::sudo "killall redis-server"
  iam::common::sudo "yum -y remove redis"
  iam::common::sudo "rm -rf /var/lib/redis"
  set -o errexit
  iam::log::info "uninstall Redis successfully"
}

# 状态检查
function iam::redis::status()
{
  if [[ -z "`pgrep redis-server`" ]];then
    iam::log::error_exit "Redis not running, maybe not installed properly"
    return 1
  fi


  redis-cli --no-auth-warning -h ${REDIS_HOST} -p ${REDIS_PORT} -a "${REDIS_PASSWORD}" --hotkeys || {
    iam::log::error "can not login with ${REDIS_USERNAME}, redis maybe not initialized properly"
    return 1
  }
}

#eval $*
if [[ "$*" =~ iam::redis:: ]];then
  eval $*
fi
