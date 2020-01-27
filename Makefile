.PHONY: build clean

# get tag of kubecm
KUBECM_VERSION=$(shell git describe --tags `git rev-list --tags --max-count=1`)

GITVERSION:=$(shell git --version | grep ^git | sed 's/^.* //g')
GITCOMMIT:=$(shell git rev-parse HEAD)

UNAME := $(shell uname)
GORELEASER_DIST=dist
BUILD_TARGET=target
BUILD_TARGET_DIR_NAME=kubecm-$(KUBECM_VERSION)
BUILD_TARGET_PKG_DIR=$(BUILD_TARGET)/kubecm-$(KUBECM_VERSION)
BUILD_TARGET_PKG_NAME=$(BUILD_TARGET)/kubecm-$(KUBECM_VERSION).tar.gz
BUILD_TARGET_PKG_FILE_PATH=$(BUILD_TARGET)/$(BUILD_TARGET_DIR_NAME)

GO_ENV=CGO_ENABLED=0
GO_MODULE=GO111MODULE=on
VERSION_PKG=github.com/sunny0826/kubecm/cmd
GO_FLAGS=-ldflags="-X ${VERSION_PKG}.kubecmVersion=$(KUBECM_VERSION) -X ${VERSION_PKG}.gitCommit=$(GITCOMMIT) -X '${VERSION_PKG}.buildDate=`date`'"
GO=env $(GO_ENV) $(GO_MODULE) go

ifeq ($(GOOS), linux)
	GO_FLAGS=-ldflags="-linkmode external -extldflags -static -X ${VERSION_PKG}.kubecmVersion=$(KUBECM_VERSION) -X ${VERSION_PKG}.gitCommit=$(GITCOMMIT) -X '${VERSION_PKG}.buildDate=`date`'"
endif

build: pre_build
	# build kubecm
	$(GO) build $(GO_FLAGS) -o $(BUILD_TARGET_PKG_DIR)/kubecm .
	# PTAH:$(BUILD_TARGET_PKG_DIR)/kubecm

pre_build:mkdir_build_target
	# clean target
	rm -rf $(BUILD_TARGET_PKG_DIR) $(BUILD_TARGET_PKG_FILE_PATH)

# create cache dir
mkdir_build_target:
ifneq ($(BUILD_TARGET_CACHE), $(wildcard $(BUILD_TARGET_CACHE)))
	mkdir -p $(BUILD_TARGET_CACHE)
endif

clean:
	$(GO) clean ./...
	rm -rf $(BUILD_TARGET)
	rm -rf $(GORELEASER_DIST)