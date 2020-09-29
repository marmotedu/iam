#!/usr/bin/env bash

# Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
# Use of this source code is governed by a MIT style
# license that can be found in the LICENSE file.

# 本脚本功能：根据 scripts/environment.sh 配置，生成 IAM 组件 YAML 配置文件。
# 示例：genconfig.sh scripts/environment.sh configs/iam-apiserver.yaml

env_file="$1"
template_file="$2"

IAM_ROOT=$(dirname "${BASH_SOURCE[0]}")/..

source "${IAM_ROOT}/scripts/lib/init.sh"

if [ $# -ne 2 ];then
    iam::log::error "Usage: genconfig.sh scripts/environment.sh configs/iam-apiserver.yaml"
    exit 1
fi

source "${env_file}"

declare -A envs

set +u
for env in $(sed -n 's/^[^#].*${\(.*\)}.*/\1/p' ${template_file})
do
    if [ -z "$(eval echo \$${env})" ];then
        iam::log::error "environment variable '${env}' not set"
        missing=true
    fi
done

if [ "${missing}" ];then
    iam::log::error 'You may run `source scripts/environment.sh` to set these environment'
    exit 1
fi

eval "cat << EOF
$(cat ${template_file})
EOF"
