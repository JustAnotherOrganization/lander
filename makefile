PACKAGE  = lander
DATE    ?= $(shell date +%Y-%m-%d_%I:%M:%S%p)
GITHASH = $(shell git rev-parse HEAD)
VERSION ?= $(shell git describe --tags --always --dirty --match=v* 2> /dev/null || \
			cat $(CURDIR)/.version 2> /dev/null || echo v0)
PKGS     = $(or $(PKG),$(shell env GO111MODULE=on $(GO) list ./...))
TESTPKGS = $(shell env GO111MODULE=on $(GO) list -f '{{ if or .TestGoFiles .XTestGoFiles }}{{ .ImportPath }}{{ end }}' $(PKGS))
BIN      = $(CURDIR)/bin
DOCKER_BUILD_CONTEXT=.
DOCKER_FILE_PATH=Dockerfile

GO      = go
TIMEOUT = 300
V = 0
Q = $(if $(filter 1,$V),,@)
M = $(shell printf "\033[34;1m▶\033[0m")

.PHONY: all
all: fmt lint build


build: pack

	$(info $(M) building executable…) @ ## Build program binary
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GO) build \
		-tags release \
		-ldflags '-X main.GitComHash=$(GITHASH) -X main.BuildStamp=$(DATE)' \
		-o $(BIN)/$(PACKAGE) cmd/main.go

docker: fmt lint build
	docker build $(DOCKER_BUILD_ARGS) -t $(PACKAGE):$(GITHASH) $(DOCKER_BUILD_CONTEXT) -f $(DOCKER_FILE_PATH)
	@DOCKER_MAJOR=$(shell docker -v | sed -e 's/.*version //' -e 's/,.*//' | cut -d\. -f1) ; \
	DOCKER_MINOR=$(shell docker -v | sed -e 's/.*version //' -e 's/,.*//' | cut -d\. -f2) ; \
	if [ $$DOCKER_MAJOR -eq 1 ] && [ $$DOCKER_MINOR -lt 10 ] ; then \
		echo docker tag -f $(PACKAGE):$(GITHASH) $(IMAGE):latest ;\
		docker tag -f $(PACKAGE):$(GITHASH) $(IMAGE):latest ;\
	else \
		echo docker tag $(PACKAGE):$(GITHASH) $(PACKAGE):latest ;\
		docker tag $(PACKAGE):$(GITHASH) $(PACKAGE):latest ; \
	fi



.PHONY: lint
lint: | $(GOLINT) ; $(info $(M) running golint…) @ ## Run golint
	$Q $(GOLINT) -set_exit_status $(PKGS)

.PHONY: fmt
fmt: ; $(info $(M) running gofmt…) @ ## Run gofmt on all source files
	$Q $(GO) fmt ./...

		#; docker build -t $(PACKAGE)

.PHONY: clean
clean: ; $(info $(M) cleaning…)	@ ## Cleanup everything
	@rm -rf $(BIN)
	@rm -rf test/tests.* test/coverage.*

.PHONY: help
help:
	@grep -E '^[ a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
	awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'

.PHONY: version
version:
	@echo $(VERSION)

pack:
	$(info $(M) binding assets) @ ## Build program binary
	go-bindata -ignore=data.go -pkg data -prefix data -o data.go ./...