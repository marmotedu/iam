#!/usr/bin/env bash

for n in $(seq 1 1 10)
do
  nohup curl http://iam.api.marmotedu.com/healthz &>/dev/null &
done
