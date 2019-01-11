PKG_LIST := $(shell go list ./... | grep -v /vendor/)
GO_FILES := $(shell find . -name '*.go' | grep -v /vendor/ | grep -v _test.go)
SHORT_VERSION := $(shell echo ${VERSION} | cut -d"v" -f2)

.PHONY: all build test lint site

all: build

site_install:
	curl -LO https://github.com/gohugoio/hugo/releases/download/v0.46/hugo_0.46_Linux-64bit.tar.gz
	tar -xvzf hugo_0.46_Linux-64bit.tar.gz
	mv hugo /home/travis/bin/

site:
	cd docs && hugo

lint_install:
	curl -LO https://github.com/alecthomas/gometalinter/releases/download/v2.0.12/gometalinter-2.0.12-linux-amd64.tar.gz
	tar xf gometalinter-2.0.12-linux-amd64.tar.gz
	mv gometalinter-2.0.12-linux-amd64/* /home/travis/bin/

lint: ## Lint the files
	gometalinter --exclude=.*_test.go --concurrency=1 --deadline=1000s --line-length=100 --disable-all --enable=vet --enable=vetshadow --enable=deadcode --enable=gocyclo --enable=golint --enable=dupl --enable=ineffassign --enable=goconst --enable=gosec --enable=goimports --enable=lll --enable=misspell ./...

test: ## Run unittests
	go test -race ${PKG_LIST}

build: ## Build the binary file
	go get github.com/mitchellh/gox
	mkdir -p bin
	go mod download
	gox -osarch="darwin/386 darwin/amd64 linux/386 linux/amd64 linux/arm" -ldflags="-extldflags '-static' -X github.com/signmykeyio/signmykey/cmd.versionString=${SHORT_VERSION}" -output="bin/signmykey_{{.OS}}_{{.Arch}}"
	zip -j bin/signmykey_darwin_386.zip bin/signmykey_darwin_386
	zip -j bin/signmykey_darwin_amd64.zip bin/signmykey_darwin_amd64
	zip -j bin/signmykey_linux_386.zip bin/signmykey_linux_386
	zip -j bin/signmykey_linux_amd64.zip bin/signmykey_linux_amd64
	zip -j bin/signmykey_linux_arm.zip bin/signmykey_linux_arm

fpm_install:
	sudo apt update && sudo apt install ruby-dev build-essential rpm -y
	gem install --no-ri --no-rdoc fpm

fpm:
	cp bin/signmykey_linux_amd64 signmykey
	fpm -s dir -t deb -n signmykey -m "contact@pablo-ruth.fr" --url "https://github.com/signmykeyio/signmykey" --description "An automated SSH Certificate Authority" --category "admin" -v ${SHORT_VERSION} --prefix /usr/bin signmykey
	fpm -s dir -t rpm -n signmykey -m "contact@pablo-ruth.fr" --url "https://github.com/signmykeyio/signmykey" --description "An automated SSH Certificate Authority" --category "admin" -v ${SHORT_VERSION} --prefix /usr/bin signmykey

fpm_upload_dev:
	curl -u$(BINTRAY_USER):$(BINTRAY_TOKEN) --data '{"name":"${SHORT_VERSION}","desc":"${SHORT_VERSION}"}' https://api.bintray.com//packages/signmykeyio/signmykey-dev-deb/signmykey/versions
	curl -T signmykey_${SHORT_VERSION}_amd64.deb -u$(BINTRAY_USER):$(BINTRAY_TOKEN) "https://api.bintray.com/content/signmykeyio/signmykey-dev-deb/signmykey/${SHORT_VERSION}/pool/signmykey_${SHORT_VERSION}_amd64.deb;deb_distribution=stable;deb_component=main;deb_architecture=amd64"
	curl -u$(BINTRAY_USER):$(BINTRAY_TOKEN) --data '{"discard":true,"publish_wait_for_secs":-1,"subject":"signmykey.io"}' "https://api.bintray.com/content/signmykeyio/signmykey-dev-deb/signmykey/${SHORT_VERSION}/publish"
	curl -u$(BINTRAY_USER):$(BINTRAY_TOKEN) --data '{"name":"${SHORT_VERSION}","desc":"${SHORT_VERSION}"}' https://api.bintray.com//packages/signmykeyio/signmykey-dev-rpm/signmykey/versions
	curl -T signmykey-${SHORT_VERSION}-1.x86_64.rpm -u$(BINTRAY_USER):$(BINTRAY_TOKEN) "https://api.bintray.com/content/signmykeyio/signmykey-dev-rpm/signmykey/${SHORT_VERSION}/pool/signmykey-${SHORT_VERSION}-1.x86_64.rpm"
	curl -u$(BINTRAY_USER):$(BINTRAY_TOKEN) --data '{"discard":true,"publish_wait_for_secs":-1,"subject":"signmykey.io"}' "https://api.bintray.com/content/signmykeyio/signmykey-dev-rpm/signmykey/${SHORT_VERSION}/publish"
	
fpm_upload_tag:
	curl -u$(BINTRAY_USER):$(BINTRAY_TOKEN) --data '{"name":"${SHORT_VERSION}","desc":"${SHORT_VERSION}"}' https://api.bintray.com//packages/signmykeyio/signmykey-deb/signmykey/versions
	curl -T signmykey_${SHORT_VERSION}_amd64.deb -u$(BINTRAY_USER):$(BINTRAY_TOKEN) "https://api.bintray.com/content/signmykeyio/signmykey-deb/signmykey/${SHORT_VERSION}/pool/signmykey_${SHORT_VERSION}_amd64.deb;deb_distribution=stable;deb_component=main;deb_architecture=amd64"
	curl -u$(BINTRAY_USER):$(BINTRAY_TOKEN) --data '{"discard":true,"publish_wait_for_secs":-1,"subject":"signmykey.io"}' "https://api.bintray.com/content/signmykeyio/signmykey-deb/signmykey/${SHORT_VERSION}/publish"
	curl -u$(BINTRAY_USER):$(BINTRAY_TOKEN) --data '{"name":"${SHORT_VERSION}","desc":"${SHORT_VERSION}"}' https://api.bintray.com//packages/signmykeyio/signmykey-rpm/signmykey/versions
	curl -T signmykey-${SHORT_VERSION}-1.x86_64.rpm -u$(BINTRAY_USER):$(BINTRAY_TOKEN) "https://api.bintray.com/content/signmykeyio/signmykey-rpm/signmykey/${SHORT_VERSION}/pool/signmykey-${SHORT_VERSION}-1.x86_64.rpm"
	curl -u$(BINTRAY_USER):$(BINTRAY_TOKEN) --data '{"discard":true,"publish_wait_for_secs":-1,"subject":"signmykey.io"}' "https://api.bintray.com/content/signmykeyio/signmykey-rpm/signmykey/${SHORT_VERSION}/publish"

help: ## Display this help screen
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
