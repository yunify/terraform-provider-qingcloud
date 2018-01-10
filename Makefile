GOFMT_FILES?=$$(find . -name '*.go' | grep -v vendor)
VETARGS?=-all
TEST?=$$(go list ./... |grep -v 'vendor')
RELEASE_TAG=$$(git describe --abbrev=0 --tags)


all: build test

build: fmt
	go build -o terraform-provider-qingcloud

install:
	cp terraform-provider-qingcloud $(shell dirname `which terraform`)

test: fmtcheck
	go test -i $(TEST) || exit 1
	echo $(TEST) | \
		xargs -t -n4 go test $(TESTARGS) -timeout=30s -parallel=4

testacc: fmtcheck
	TF_ACC=1 go test $(TEST) -v $(TESTARGS) -timeout 120m

vet:
	@echo "go tool vet $(VETARGS) ."
	@go tool vet $(VETARGS) $$(ls -d */ | grep -v vendor) ; if [ $$? -eq 1 ]; then \
		echo ""; \
		echo "Vet found suspicious constructs. Please check the reported constructs"; \
		echo "and fix them if necessary before submitting the code for review."; \
		exit 1; \
	fi

fmt:
	gofmt -w $(GOFMT_FILES)

fmtcheck:
	@sh -c "'$(CURDIR)/scripts/gofmtcheck.sh'"

errcheck:
	@sh -c "'$(CURDIR)/scripts/errcheck.sh'"

dist-tools:
	@go get github.com/mitchellh/gox

dist: dist-tools
	rm -rf ./bin/*
	mkdir -p ./bin/terraform-provider-qingcloud_linux-amd64_$(RELEASE_TAG)
	mkdir -p ./bin/terraform-provider-qingcloud_darwin-amd64_$(RELEASE_TAG)
	mkdir -p ./bin/terraform-provider-qingcloud_windows-amd64_$(RELEASE_TAG)
	gox -osarch="linux/amd64" -output=./bin/terraform-provider-qingcloud_linux-amd64_$(RELEASE_TAG)/terraform-provider-qingcloud_$(RELEASE_TAG)
	gox -osarch="darwin/amd64" -output=./bin/terraform-provider-qingcloud_darwin-amd64_$(RELEASE_TAG)/terraform-provider-qingcloud_$(RELEASE_TAG)
	gox -osarch="windows/amd64" -output=./bin/terraform-provider-qingcloud_windows-amd64_$(RELEASE_TAG)/terraform-provider-qingcloud_$(RELEASE_TAG)
	cd bin && ls --color=no | xargs -I {} tar -czf {}.tgz {}

.PHONY: all build copy test vet fmt fmtcheck errcheck dist-tools dist
