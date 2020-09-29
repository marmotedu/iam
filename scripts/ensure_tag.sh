#!/usr/bin/env bash

version=v`gsemver bump`
if [ -z "`git tag -l $version`" ];then
  git tag -a -m "release version $version" $version
fi
