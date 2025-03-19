#!/bin/bash -eux

cwd=$(pwd)

npm install -g @redocly/cli

pushd $cwd/dp-identity-api
  go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.64.8
  make lint
popd
