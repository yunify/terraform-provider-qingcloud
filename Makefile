GOFMT_FILES?=$$(find . -name '*.go' | grep -v vendor)
VETARGS?=-all
TEST?=$$(go list ./...)


all: build copy

build:
	go build -o terraform-provider-qingcloud

copy:
	cp terraform-provider-qingcloud $(shell dirname `which terraform`)

test: vet fmtcheck errcheck
	go test -v ./qingcloud
	TF_ACC=1 go test -v ./qingcloud -run=TestAccQingcloud -timeout=120m -parallel=4

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

dist:
	rm ./build/*
	gox -osarch="linux/amd64" -output=terraform-provider-qingcloud_{{.OS}}-{{.Arch}}
	gox -osarch="darwin/amd64" -output=terraform-provider-qingcloud_{{.OS}}-{{.Arch}}
	gox -osarch="windows/amd64" -output=terraform-provider-qingcloud_{{.OS}}-{{.Arch}}
	mkdir -p ./build
	mv terraform-provider-qingcloud_* ./build
	cd build && ls --color=no | xargs -I {} tar -czf {}.tgz {}
