#!/bin/bash

# Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
# Use of this source code is governed by a MIT style
# license that can be found in the LICENSE file.

IAM_ROOT=$(dirname "${BASH_SOURCE[0]}")/..
source "${IAM_ROOT}/scripts/lib/init.sh"

if [ $# -ne 1 ];then
  iam::log::error "Usage: force_release.sh v1.0.0"
  exit 1
fi

version="$1"

set +o errexit
# 1. delete old version
git tag -d ${version}
git push origin --delete ${version}

# 2. create a new tag
git tag -a ${version} -m "release ${version}"
git push origin master
git push origin ${version}

# 3. release the new release
pushd ${IAM_ROOT}
# try to delete target github release if exist to avoid create error
iam::log::info "delete github release with tag ${version} if exist"
github-release delete  \
  --user marmotedu\
  --repo iam  \
  --tag ${version} &> /dev/null

make release
