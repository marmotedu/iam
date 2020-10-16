#!/usr/bin/env bash

for port in 80
do    
  if echo |telnet 127.0.0.1 $port 2>&1|grep refused &>/dev/null;then    
    exit 1    
  fi
done
