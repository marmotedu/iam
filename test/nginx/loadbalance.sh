#!/usr/bin/env bash

for domain in iam.api.marmotedu.com iam.authz.marmotedu.com
do
  for n in $(seq 1 1 10)
  do
    echo $domain
    nohup curl http://${domain}/healthz &>/dev/null &
  done
done
