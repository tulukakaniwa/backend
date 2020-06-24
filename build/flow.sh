#!/usr/bin/bash

cd /OpenROAD-flow
source setup_env.sh

cd /OpenROAD-flow/flow
make DESIGN_CONFIG=/cloud/config.mk

mkdir /cloud/openroad-flow
cp /cloud/repo/openroad.yml /cloud/openroad-flow
cp -r logs /cloud/openroad-flow/
cp -r objects /cloud/openroad-flow/
cp -r reports /cloud/openroad-flow/
cp -r results /cloud/openroad-flow/
