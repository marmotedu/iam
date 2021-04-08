#!/usr/bin/env bash

# Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
# Use of this source code is governed by a MIT style
# license that can be found in the LICENSE file.


#!/bin/bash

server=$1
base_dir=$PWD
interval=2
timeout=180

# 命令行参数，需要手动指定
args=""

function findPid() {
  echo -n `ps -ef -u $UID|grep "${base_dir}/${server} ${args}"|egrep -v 'grep|admin.sh'|awk '{print $2}'`
}

function waitProcess() {
  i=0
  while (($i<${timeout}))
  do
    if [ "`findPid`" == "" -a "$1" == "stop" ];then
      break
    fi

    if [ "`findPid`" != "" -a "$1" == "start" ];then
      break
    fi

    echo waiting to $1 ...

    sleep 1

    ((i++))
  done

  if [ "$i" -ge "${timeout}" -a "$1" == "stop" ];then
    echo "waiting timeout(${timeout}s), send SIGKILL signal."
    kill -9 `findPid`
  fi
  sleep 1
}


function start()
{
  if [ "`findPid`" != "" ];then
    echo "${server} already running"
    exit 0
  fi

  nohup ${base_dir}/${server} ${args} &>/dev/null &

  waitProcess start

  # check status
  if [ "`findPid`" == "" ];then
    echo "${server} start failed"
    exit 1
  fi
}

function status()
{
  if [ "`findPid`" != "" ];then
    echo ${server} is running
  else
    echo ${server} is not running
  fi
}

function stop()
{
  if [ "`findPid`" != "" ];then
    echo "send SIGKILL signal to `findPid`"
    kill `findPid`
  fi

  waitProcess stop

  if [ "`findPid`" != "" ];then
    echo "${server} stop failed after ${${timeout}}s"
    exit 1
  fi
}

case "$2" in
  'start')
    start
    ;;
  'stop')
    stop
    ;;
  'status')
    status
    ;;
  'restart')
    stop && start
    ;;
  *)
    echo "usage: $0 {start|stop|restart|status}"
    exit 0
    ;;
esac
