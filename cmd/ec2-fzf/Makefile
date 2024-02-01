export GO111MODULE ?= on
export GOPROXY ?= https://proxy.golang.org
export GOSUMDB ?= sum.golang.org

PROJECT?=github.com/richie-tt/ec2-fzf
APP?=ec2-fzf

RELEASE?=$(shell git describe --tags 2>/dev/null || echo 0.0.1)
COMMIT?=$(shell git rev-parse --short HEAD)
BUILD_TIME?=$(shell date '+%Y-%m-%d_%H:%M:%S')
# CONTAINER_IMAGE?=docker.io/webdeva/${APP}

GOOS?=$(shell uname | tr 'A-Z' 'a-z')
# GOARCH?=amd64

LOCAL_ARCH := $(shell uname -m)
ifeq ($(LOCAL_ARCH),x86_64)
	TARGET_ARCH_LOCAL=amd64
else ifeq ($(shell echo $(LOCAL_ARCH) | head -c 5),armv8)
	TARGET_ARCH_LOCAL=arm64
else ifeq ($(shell echo $(LOCAL_ARCH) | head -c 4),armv)
	TARGET_ARCH_LOCAL=arm
else ifeq ($(shell echo $(LOCAL_ARCH) | head -c 5),arm64)
	TARGET_ARCH_LOCAL=arm64
else
	TARGET_ARCH_LOCAL=amd64
endif
export GOARCH ?= $(TARGET_ARCH_LOCAL)


clean:
	rm -f ${APP}

build: clean
	CGO_ENABLED=0 GOOS=${GOOS} GOARCH=${GOARCH} go build \
		-ldflags "-s -w \
		-X github.com/richie-tt/ec2-fzf.Release=${RELEASE} \
		-X github.com/richie-tt/ec2-fzf.Commit=${COMMIT} \
		-X github.com/richie-tt/ec2-fzf.BuildTime=${BUILD_TIME} \
		-X github.com/richie-tt/ec2-fzf.Arch=${GOARCH}" \
		-o ${APP} \
		-v

.DEFAULT_GOAL := build
