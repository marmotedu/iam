#!/bin/bash

IAM_ROOT=$(dirname "${BASH_SOURCE[0]}")/..    
source "${IAM_ROOT}/scripts/lib/init.sh"
pushd ${IAM_ROOT}
make release
