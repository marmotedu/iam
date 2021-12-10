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

iam::golang::setup_env

BINS=(
  gendocs
  geniamdocs
  genman
  genyaml
)
make build -C "${IAM_ROOT}" BINS="${BINS[*]}"

iam::util::ensure-temp-dir

iam::util::gen-docs "${IAM_TEMP}"

# remove all of the old docs
iam::util::remove-gen-docs

# Copy fresh docs into the repo.
# the shopt is so that we get docs/.generated_docs from the glob.
shopt -s dotglob
cp -af "${IAM_TEMP}"/* "${IAM_ROOT}"
shopt -u dotglob
