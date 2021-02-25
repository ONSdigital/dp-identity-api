#!/bin/bash -eux

cwd=$(pwd)

pushd dp-identity-api
  make test
popd
