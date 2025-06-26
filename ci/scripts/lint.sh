#!/bin/bash -eux

cwd=$(pwd)

npm install -g @redocly/cli

pushd $cwd/dp-identity-api
  go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.1.6
  make lint
popd
