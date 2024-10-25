#!/bin/bash -eux

cwd=$(pwd)

pushd $cwd/dp-identity-api
  go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.61.0
  make lint
popd