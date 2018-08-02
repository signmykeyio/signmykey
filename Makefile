PKG_LIST := $(shell go list ./... | grep -v /vendor/)
GO_FILES := $(shell find . -name '*.go' | grep -v /vendor/ | grep -v _test.go)
VERSION := ${VERSION}

.PHONY: all build test lint

all: build

lint_install:
	apt-get update -y
	apt-get install -y curl
	curl -LO https://github.com/alecthomas/gometalinter/releases/download/v2.0.5/gometalinter-2.0.5-linux-amd64.tar.gz
	tar xf gometalinter-2.0.5-linux-amd64.tar.gz
	mv gometalinter-2.0.5-linux-amd64/* /usr/local/go/bin/

lint: ## Lint the files
	gometalinter --vendor --exclude=.*_test.go --concurrency=1 --deadline=1000s --line-length=100 --enable=goimports --enable=lll --enable=misspell --enable=nakedret --enable=unparam ./...

test: ## Run unittests
	go test -short ${PKG_LIST}

build: ## Build the binary file
	go build -ldflags "-w -s -extldflags '-static' -X gitlab.com/signmykey/signmykey/cmd.versionString=$(VERSION)"

fpm_install:
	apt-get update -y && apt-get install ruby ruby-dev rubygems build-essential -y
	gem install --no-ri --no-rdoc fpm

fpm:
	fpm -s dir -t deb -n signmykey -m "contact@pablo-ruth.fr" --url "https://gitlab.com/signmykey/signmykey" --description "A light command to sign ssh keys with signmykey-server" --category "admin" -v $(VERSION) --prefix /usr/bin signmykey

fpm_upload_dev:
	@curl -u $(APTLY_USER):$(APTLY_PASSWORD) -X POST -F file=@signmykey_$(VERSION)_amd64.deb https://apt.signmykey.io/api/files/signmykey_$(VERSION)
	@curl -u $(APTLY_USER):$(APTLY_PASSWORD) -X POST https://apt.signmykey.io/api/repos/signmykey-dev/file/signmykey_$(VERSION)
	@curl -u $(APTLY_USER):$(APTLY_PASSWORD) -X PUT -H 'Content-Type: application/json' --data '{"Signing": {"Skip": true}}' https://apt.signmykey.io/api/publish/signmykey-dev/xenial
	@curl -f -o /dev/null --silent --head https://apt.signmykey.io/signmykey-dev/pool/main/s/signmykey/signmykey_$(VERSION)_amd64.deb
	
fpm_upload_tag:
	@curl -u $(APTLY_USER):$(APTLY_PASSWORD) -X POST -F file=@signmykey_$(VERSION)_amd64.deb https://apt.signmykey.io/api/files/signmykey_$(VERSION)
	@curl -u $(APTLY_USER):$(APTLY_PASSWORD) -X POST https://apt.signmykey.io/api/repos/signmykey/file/signmykey_$(VERSION)
	@curl -u $(APTLY_USER):$(APTLY_PASSWORD) -X PUT -H 'Content-Type: application/json' --data '{"Signing": {"Skip": true}}' https://apt.signmykey.io/api/publish/signmykey/xenial
	@curl -f -o /dev/null --silent --head https://apt.signmykey.io/signmykey/pool/main/s/signmykey/signmykey_$(VERSION)_amd64.deb

help: ## Display this help screen
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
