#!/usr/bin/env bash

# Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
# Use of this source code is governed by a MIT style
# license that can be found in the LICENSE file.

# /etc/keepalived/keepalived_notify.sh
log_file=/var/log/keepalived.log

iam::keepalived::mail() {
  # 这里可以添加email逻辑，当keepalived变动时及时告警
  :
}
iam::keepalived::log() {
    echo "[`date '+%Y-%m-%d %T'`] $1" >> ${log_file}
}

[ ! -d /var/keepalived/ ] && mkdir -p /var/keepalived/

case "$1" in
    "MASTER" )
        iam::keepalived::log "notify_master"
        ;;
    "BACKUP" )
        iam::keepalived::log "notify_backup"
        ;;
    "FAULT" )
        iam::keepalived::log "notify_fault"
        ;;
    "STOP" )
        iam::keepalived::log "notify_stop"
        ;;
    *)
        iam::keepalived::log "keepalived_notify.sh: state error!"
        ;;
esac
