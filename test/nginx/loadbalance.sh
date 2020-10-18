#!/usr/bin/env bash

# Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
# Use of this source code is governed by a MIT style
# license that can be found in the LICENSE file.

for n in $(seq 1 1 10)
do
  nohup curl http://iam.api.marmotedu.com/healthz &>/dev/null &
done
