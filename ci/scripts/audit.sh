#!/bin/bash -eux

export cwd=$(pwd)

pushd $cwd/dp-identity-api
  make audit
popd