#!/usr/bin/env bash

# Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
# Use of this source code is governed by a MIT style
# license that can be found in the LICENSE file.


mkdir -p "$DST_DIR/scripts/install"

cp -rv scripts/lib "$DST_DIR/scripts/"
cp -v scripts/install/{common.sh,environment.sh,test.sh} "$DST_DIR/scripts/install"
