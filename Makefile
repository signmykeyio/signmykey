PKG_LIST := $(shell go list ./... | grep -v /vendor/)
GO_FILES := $(shell find . -name '*.go' | grep -v /vendor/ | grep -v _test.go)
VERSION := ${VERSION}

.PHONY: all build test lint site

all: build

site_install:
	curl -LO https://github.com/gohugoio/hugo/releases/download/v0.46/hugo_0.46_Linux-64bit.tar.gz
	tar -xvzf hugo_0.46_Linux-64bit.tar.gz
	mv hugo /home/travis/bin/

site:
	cd docs && hugo

lint_install:
	curl -LO https://github.com/alecthomas/gometalinter/releases/download/v2.0.10/gometalinter-2.0.10-linux-amd64.tar.gz
	tar xf gometalinter-2.0.10-linux-amd64.tar.gz
	mv gometalinter-2.0.10-linux-amd64/* /home/travis/bin/

lint: ## Lint the files
	gometalinter --vendor --exclude=.*_test.go --concurrency=1 --deadline=1000s --line-length=100 --enable=goimports --enable=lll --enable=misspell --enable=nakedret --enable=unparam ./...

test: ## Run unittests
	go test -race ${PKG_LIST}

build: ## Build the binary file
	go get github.com/mitchellh/gox
	mkdir -p bin
	gox -ldflags="-extldflags '-static' -X github.com/signmykeyio/signmykey/cmd.versionString=$(VERSION)" -output="bin/signmykey_{{.OS}}_{{.Arch}}"

fpm_install:
	sudo apt update && sudo apt install ruby-dev build-essential rpm -y
	gem install --no-ri --no-rdoc fpm

fpm:
	fpm -s dir -t deb -n signmykey -m "contact@pablo-ruth.fr" --url "https://github.com/signmykeyio/signmykey" --description "An automated SSH Certificate Authority" --category "admin" -v $(VERSION) --prefix /usr/bin signmykey
	fpm -s dir -t rpm -n signmykey -m "contact@pablo-ruth.fr" --url "https://github.com/signmykeyio/signmykey" --description "An automated SSH Certificate Authority" --category "admin" -v $(VERSION) --prefix /usr/bin signmykey

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
