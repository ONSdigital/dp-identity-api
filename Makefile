SHELL=bash
MAIN=dp-identity-api

BUILD=build
BUILD_ARCH=$(BUILD)/$(GOOS)-$(GOARCH)
BIN_DIR?=.

BUILD_TIME=$(shell date +%s)
GIT_COMMIT=$(shell git rev-parse HEAD)
VERSION ?= $(shell git tag --points-at HEAD | grep ^v | head -n 1)
LDFLAGS=-ldflags "-w -s -X 'main.Version=${VERSION}' -X 'main.BuildTime=$(BUILD_TIME)' -X 'main.GitCommit=$(GIT_COMMIT)'"

LOCAL_USER_POOL_ID=eu-west-2_WSD9EcAsw

export GOOS?=$(shell go env GOOS)
export GOARCH?=$(shell go env GOARCH)

.PHONY: all
all: audit lint test build

.PHONY: audit
audit:
	set -o pipefail; go list -m all | nancy sleuth

.PHONY: build
build:
	@mkdir -p $(BUILD_ARCH)/$(BIN_DIR)
	go build $(LDFLAGS) -o $(BUILD_ARCH)/$(BIN_DIR)/$(MAIN) main.go

.PHONY: debug-watch
debug-watch: 
	reflex -d none -c ./reflex

.PHONY: debug
debug:
	export AWS_COGNITO_USER_POOL_ID=$(LOCAL_USER_POOL_ID);
	export AWS_COGNITO_CLIENT_ID=${AWS_COGNITO_CLIENT_ID:?please define a valid AWS_COGNITO_CLIENT_ID in your local system, get from within pool}
	export AWS_COGNITO_CLIENT_SECRET=${AWS_COGNITO_CLIENT_SECRET:?please define a valid AWS_COGNITO_CLIENT_SECRET in your local system, get from within pool}
	echo AWS_COGNITO_USER_POOL_ID= $$AWS_COGNITO_USER_POOL_ID;\
	HUMAN_LOG=1 go run $(LDFLAGS) -race main.go
	
.PHONY: acceptance
acceptance:
	MONGODB_IMPORTS_DATABASE=test HUMAN_LOG=1 go run $(LDFLAGS) -race main.go

.PHONY: lint
lint: validate-specification
	golangci-lint run ./...

.PHONY: validate-specification
validate-specification:
	redocly lint swagger.yaml

.PHONY: test
test:
	go test -cover -race ./...

.PHONY: test build debug

.PHONY: test-component
test-component:
	go test -cover -race -coverpkg=github.com/ONSdigital/dp-identity-api/... -component

.PHONY: populate-local
populate-local:
	export AWS_COGNITO_USER_POOL_ID=$(LOCAL_USER_POOL_ID); \
	HUMAN_LOG=1 go run -race ./dummy-data/import-dummy-users/populate_dummy_data.go

.PHONY: remove-test-data
remove-test-data:
	export AWS_COGNITO_USER_POOL_ID=$(LOCAL_USER_POOL_ID); \
	HUMAN_LOG=1 go run -race ./dummy-data/delete-dummy-users/remove_dummy_data.go

.PHONY: get-jwks-keys
get-jwks-keys:
	HUMAN_LOG=1 go run ./scripts/get-jwks-keys/main.go
