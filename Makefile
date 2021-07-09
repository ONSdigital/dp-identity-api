SHELL=bash
MAIN=dp-identity-api

BUILD=build
BUILD_ARCH=$(BUILD)/$(GOOS)-$(GOARCH)
BIN_DIR?=.

BUILD_TIME=$(shell date +%s)
GIT_COMMIT=$(shell git rev-parse HEAD)
VERSION ?= $(shell git tag --points-at HEAD | grep ^v | head -n 1)
LDFLAGS=-ldflags "-w -s -X 'main.Version=${VERSION}' -X 'main.BuildTime=$(BUILD_TIME)' -X 'main.GitCommit=$(GIT_COMMIT)'"

export GOOS?=$(shell go env GOOS)
export GOARCH?=$(shell go env GOARCH)

.PHONY: all
all: audit test build

.PHONY: audit
audit:
	go list -m all | nancy sleuth

.PHONY: build
build:
	@mkdir -p $(BUILD_ARCH)/$(BIN_DIR)
	go build $(LDFLAGS) -o $(BUILD_ARCH)/$(BIN_DIR)/$(MAIN) main.go

.PHONY: debug
debug:
	export AWS_COGNITO_USER_POOL_ID=eu-west-1_Rnma9lp2q; \
	export AWS_COGNITO_CLIENT_ID=`aws cognito-idp list-user-pool-clients --user-pool-id $$AWS_COGNITO_USER_POOL_ID --query 'UserPoolClients[0].ClientId' --output text`; \
	export AWS_COGNITO_CLIENT_SECRET=`aws cognito-idp describe-user-pool-client --user-pool-id $$AWS_COGNITO_USER_POOL_ID --client-id $$AWS_COGNITO_CLIENT_ID --query 'UserPoolClient.ClientSecret' --output text`; \
	echo AWS_COGNITO_USER_POOL_ID= $$AWS_COGNITO_USER_POOL_ID;\
	echo AWS_COGNITO_CLIENT_ID= $$AWS_COGNITO_CLIENT_ID;\
	echo AWS_COGNITO_CLIENT_SECRET= $$AWS_COGNITO_CLIENT_SECRET;\
	HUMAN_LOG=1 go run $(LDFLAGS) -race main.go
	
.PHONY: acceptance
acceptance:
	MONGODB_IMPORTS_DATABASE=test HUMAN_LOG=1 go run $(LDFLAGS) -race main.go

.PHONY: lint
lint:
	exit

.PHONY: test
test:
	go test -cover -race ./...

.PHONY: test build debug

.PHONY: test-component
test-component:
	go test -cover -race -coverpkg=github.com/ONSdigital/dp-identity-api/... -component
