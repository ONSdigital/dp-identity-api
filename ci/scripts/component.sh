#!/bin/bash -eux

export cwd=$(pwd)

pushd $cwd/dp-identity-api
  make test-component
popd