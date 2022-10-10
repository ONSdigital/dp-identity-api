# see main.go for all env var options/defaults

SHELL=bash
APP=collections-migration

ENV?=develop
TIMEOUT?=30s

host_num?=publishing_mount 1
host_bin=bin-$(APP)

GOOS?=$(shell go env GOOS)
GOARCH?=$(shell go env GOARCH)

BUILD=build
BUILD_ARCH=$(BUILD)/$(GOOS)-$(GOARCH)
BUILD_SCRIPT=$(BUILD_ARCH)/$(APP).sh

.PHONY: all
all: clean
	make GOOS=linux build deploy
	make GOOS=linux build deploy

.PHONY: pre-build
pre-build:
	mkdir -p $(BUILD_ARCH)

.PHONY: build
build: pre-build
	go build -o $(BUILD_ARCH)/$(APP) main.go

.PHONY: run
run:
	go run -race main.go

.PHONY: deploy
deploy: build
	dp scp $(ENV) $(host_num) -r -- $(BUILD_ARCH) $(host_bin)

.PHONY: clean
clean: clean-deploy clean-build

.PHONY: clean-build
clean-build:
	[[ ! -d $(BUILD) ]] || rm -r $(BUILD)

.PHONY: clean-deploy
clean-deploy:
	dp ssh $(ENV) $(host_num) -- bash -c '"[[ ! -d $(host_bin) ]] || rm -r $(host_bin)"'

