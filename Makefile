.PHONY: build clean

# get tag of kubecm
KUBECM_VERSION=$(shell git describe --tags `git rev-list --tags --max-count=1`)
TAG=$(KUBECM_VERSION)

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

# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

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

tag:
	git tag -a $(TAG) -m "$(TAG) relsase"

push_tag:
	git push origin $(KUBECM_VERSION)

# Run go fmt against code
fmt:
	go fmt ./...

# Run go vet against code
vet:
	go vet ./...

lint: golangci
	$(GOLANGCILINT) run --disable=typecheck --timeout 10m -E golint,goimports  ./...

test: fmt vet lint
		go test -race -coverprofile=coverage.txt -covermode=atomic ./cmd/...


GOLANGCILINT_VERSION ?= v1.29.0
HOSTOS := $(shell uname -s | tr '[:upper:]' '[:lower:]')
HOSTARCH := $(shell uname -m)
ifeq ($(HOSTARCH),x86_64)
HOSTARCH := amd64
endif

golangci:
ifeq (, $(shell which golangci-lint))
	@{ \
	set -e ;\
	echo 'installing golangci-lint-$(GOLANGCILINT_VERSION)' ;\
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(GOBIN) $(GOLANGCILINT_VERSION) ;\
	echo 'Install succeed' ;\
	}
GOLANGCILINT=$(GOBIN)/golangci-lint
else
GOLANGCILINT=$(shell which golangci-lint)
endif