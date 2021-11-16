#!/usr/bin/env bash

# Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
# Use of this source code is governed by a MIT style
# license that can be found in the LICENSE file.

# This file is not intended to be run automatically. It is meant to be run
# immediately before exporting docs. We do not want to check these documents in
# by default.

set -o errexit
set -o nounset
set -o pipefail

IAM_ROOT=$(dirname "${BASH_SOURCE[0]}")/..
source "${IAM_ROOT}/scripts/lib/init.sh"

COMPONENTS=(iam-apiserver iam-authz-server iam-pump iam-watcher)
KINDS=(deployment service)

for component in ${COMPONENTS[@]}
do
  truncate -s 0 ${IAM_ROOT}/deployments/${component}.yaml

  for kind in ${KINDS[@]}
  do
    echo -e "---\n# Source: deployments/${component}-${kind}.yaml" >> ${IAM_ROOT}/deployments/${component}.yaml
    sed '/^#\|^$/d' ${IAM_ROOT}/deployments/${component}-${kind}.yaml >> ${IAM_ROOT}/deployments/${component}.yaml
  done

  iam::log::info "generate ${IAM_ROOT}/deployments/${component}.yaml success"
done
