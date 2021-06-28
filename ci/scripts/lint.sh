#!/bin/bash -eux

cwd=$(pwd)

pushd $cwd/dp-identity-api
  make lint
popd