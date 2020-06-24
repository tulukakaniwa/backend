#!/usr/bin/env bash


CURRENT_DIR=$(dirname "$BASH_SOURCE")
rm -rf /tmp/OpenROAD-flow
git clone --recursive https://github.com/The-OpenROAD-Project/OpenROAD-flow.git /tmp/OpenROAD-flow
cd /tmp/OpenROAD-flow
./build_openroad.sh

cd $CURRENT_DIR
docker build -t openroadcloud/flow .
