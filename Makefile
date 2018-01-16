PKG_LIST := $(shell go list ./... | grep -v /vendor/)
GO_FILES := $(shell find . -name '*.go' | grep -v /vendor/ | grep -v _test.go)
VERSION := ${VERSION}

.PHONY: all build test lint

all: build

lint_install:
	apt-get update -y
	apt-get install -y curl lbzip2
	curl -LO https://github.com/alecthomas/gometalinter/releases/download/v2.0.2/gometalinter-v2.0.2-linux-amd64.tar.bz2
	tar xf gometalinter-v2.0.2-linux-amd64.tar.bz2
	mv gometalinter-v2.0.2-linux-amd64/gometalinter gometalinter-v2.0.2-linux-amd64/linters/* /usr/local/go/bin/

lint: ## Lint the files
	gometalinter --vendor --exclude=.*_test.go --concurrency=4 --deadline=120s --line-length=100 --enable=goimports --enable=lll --enable=misspell --enable=nakedret --enable=unparam ./...

test: ## Run unittests
	go test -short ${PKG_LIST}

build: ## Build the binary file
	go build -ldflags "-w -s -extldflags '-static' -X gitlab.com/signmykey/signmykey/cmd.versionString=$(VERSION)"

help: ## Display this help screen
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
