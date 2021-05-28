#!/usr/bin/env bash

# Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
# Use of this source code is governed by a MIT style
# license that can be found in the LICENSE file.


# The root of the build/dist directory
IAM_ROOT=$(dirname "${BASH_SOURCE[0]}")/../..
[[ -z ${COMMON_SOURCED} ]] && source ${IAM_ROOT}/scripts/install/common.sh

INSECURE_APISERVER=${IAM_APISERVER_HOST}:${IAM_APISERVER_INSECURE_BIND_POR}
INSECURE_AUTHZSERVER=${IAM_AUTHZ_SERVER_HOST}:${IAM_AUTHZ_SERVER_INSECURE_BIND_PORT}
CCURL="curl -s -XPOST -H'Content-Type: application/json'" # Create
UCURL="curl -s -XPUT -H'Content-Type: application/json'" # Update
RCURL="curl -s -XGET" # Retrieve
DCURL="curl -s -XDELETE" # Delete


iam::test::login()
{
  ${CCURL} -d'{"username":"admin","password":"Admin@2021"}' http://${INSECURE_APISERVER}/login
}

iam::test::user()
{
  token="-HAuthorization: Bearer $(iam::test::login)"

  # 1. 如果有colin、mark、john用户先清空
  ${DCURL} ${token} http://${INSECURE_APISERVER}/v1/users/colin
  ${DCURL} ${token} http://${INSECURE_APISERVER}/v1/users/mark
  ${DCURL} ${token} http://${INSECURE_APISERVER}/v1/users/john

  # 2. 创建colin、mark、john用户
  ${CCURL} ${token} -d'{"password":"User@2021","metadata":{"name":"colin"},"nickname":"colin","email":"colin@foxmail.com","phone":"1812884xxxx"}' http://${INSECURE_APISERVER}/v1/users

  # 3. 列出所有用户
  ${RCURL} ${token} '"http://${INSECURE_APISERVER}/v1/users?offset=0&limit=10"'

  # 4. 获取colin用户的详细信息
  ${RCURL} ${token} http://${INSECURE_APISERVER}/v1/users/colin

  # 5. 修改colin用户
  ${UCURL} ${token} -d'{"metadata":{"name":"colin"},"nickname":"colin","email":"colin_modified@foxmail.com","phone":"1812884xxxx"}' http://${INSECURE_APISERVER}/v1/users

  # 6. 删除colin用户
  ${DCURL} ${token} http://${INSECURE_APISERVER}/v1/users/colin

  # 7. 批量删除用户
  ${DCURL} ${token} "http://${INSECURE_APISERVER}/v1/users?name=mark&name=john"
  iam::log::info "congratulations, /v1/user test passed!"
}

iam::test::secret()
{
  token="-HAuthorization: Bearer $(iam::test::login)"

  # 1. 如果有secret0密钥先清空
  ${DCURL} ${token} http://${INSECURE_APISERVER}/v1/secrets/secret0

  # 2. 创建secret0密钥
  ${CCURL} ${token} -d'{"metadata":{"name":"secret0"},"expires":0,"description":"admin secret"}' http://${INSECURE_APISERVER}/v1/secrets

  # 3. 列出所有密钥
  ${RCURL} ${token} http://${INSECURE_APISERVER}/v1/secrets

  # 4. 获取secret0密钥的详细信息
  ${RCURL} ${token} http://${INSECURE_APISERVER}/v1/secrets/secret0

  # 5. 修改secret0密钥
  ${UCURL} ${token} -d'{"metadata":{"name":"secret0"},"expires":0,"description":"admin secret(modified)"}' http://${INSECURE_APISERVER}/v1/secrets

  # 6. 删除secret0密钥
  ${DCURL} ${token} http://${INSECURE_APISERVER}/v1/secrets/secret0
  iam::log::info "congratulations, /v1/secret test passed!"
}

iam::test::policy()
{
  token="-HAuthorization: Bearer $(iam::test::login)"

  # 1. 如果有policy0策略先清空
  ${DCURL} ${token} http://${INSECURE_APISERVER}/v1/policies/policy0

  # 2. 创建policy0策略
  ${CCURL} ${token} -d'{"metadata":{"name":"policy0"},"policy":{"description":"One policy to rule them all.","subjects":["users:<peter|ken>","users:maria","groups:admins"],"actions":["delete","<create|update>"],"effect":"allow","resources":["resources:articles:<.*>","resources:printer"],"conditions":{"remoteIP":{"type":"CIDRCondition","options":{"cidr":"192.168.0.1/16"}}}}}' http://${INSECURE_APISERVER}/v1/policies

  # 3. 列出所有策略
  ${RCURL} ${token} http://${INSECURE_APISERVER}/v1/policies

  # 4. 获取policy0策略的详细信息
  ${RCURL} ${token} http://${INSECURE_APISERVER}/v1/policies/policy0

  # 5. 修改policy0策略
  ${UCURL} ${token} -d'{"metadata":{"name":"policy0"},"policy":{"description":"One policy to rule them all(modified).","subjects":["users:<peter|ken>","users:maria","groups:admins"],"actions":["delete","<create|update>"],"effect":"allow","resources":["resources:articles:<.*>","resources:printer"],"conditions":{"remoteIP":{"type":"CIDRCondition","options":{"cidr":"192.168.0.1/16"}}}}}' http://${INSECURE_APISERVER}/v1/policies

  # 6. 删除policy0策略
  ${DCURL} ${token} http://${INSECURE_APISERVER}/v1/policies/policy0
  iam::log::info "congratulations, /v1/policy test passed!"
}

iam::test::apiserver()
{
  iam::test::user
  iam::test::secret
  iam::test::policy
  iam::log::info "congratulations, iam-apiserver test passed!"
}

iam::test::authz()
{
    $CCURL -H"'Authorization: Bearer $token'" -d'{"subject":"users:peter","action":"delete","resource":"resources:articles:ladon-introduction","context":{"remoteIP":"193.168.0.5"}}' http://${IAM_AUTHZSERVER_INSECURE_ADDRESS}/v1/authz
  iam::log::info "congratulations, /v1/authz test passed!"
}

iam::test::authzserver()
{
  iam::test::authz || return 1
  iam::log::info "congratulations, iam-authz-server test passed!"
}

iam::test::pump()
{
  ${RCURL} http://${IAM_PUMP_HOST}:7070/healthz | egrep -q 'status.*ok' || {
    iam::log::error "cannot access iam-pump healthz api, iam-pump maybe down"
    return 1
  }
  iam::log::info "congratulations, iam-pump test passed!"
}

iam::test::iamctl()
{
  iamctl user list | egrep -q admin || {
    iam::log::error "iamctl cannot list users from iam-apiserver"
    return 1
  }
  iam::log::info "congratulations, iamctl test passed!"
}

function iam::test::man()
{
  man iam-apiserver | grep -q 'IAM API Server' || {
    iam::log::error "iam man page not installed or may not installed properly"
    return 1
  }
  iam::log::info "congratulations, man test passed!"
}

iam::test::test()
{
  iam::test::apiserver || return 1
  iam::test::authzserver || return 1
  iam::test::pump || return 1
  iam::test::iamctl || return 1
  iam::test::man || return 1

  iam::log::info "congratulations, all test passed!"
}

if [[ "$*" =~ iam::test:: ]];then
  eval $*  
fi
