#!/usr/bin/env bash

mkdir -p "$DST_DIR/scripts/install"

cp -rv scripts/lib "$DST_DIR/scripts/"
cp -v scripts/install/{common.sh,environment.sh,test.sh} "$DST_DIR/scripts/install"
