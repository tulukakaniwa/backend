#!/usr/bin/bash

cd /OpenROAD-flow
source setup_env.sh

cd /OpenROAD-flow/flow
make DESIGN_CONFIG=/cloud/config.mk

cp -r logs /cloud
cp -r objects /cloud
cp -r reports /cloud
cp -r results /cloud
