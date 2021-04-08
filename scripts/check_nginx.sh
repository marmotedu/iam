#!/usr/bin/env bash

# Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
# Use of this source code is governed by a MIT style
# license that can be found in the LICENSE file.

for port in 80
do
  if echo |telnet 127.0.0.1 $port 2>&1|grep refused &>/dev/null;then
    exit 1
  fi
done
