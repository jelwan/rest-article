SHELL 		:= /bin/bash
PACKAGE		:= $(shell go list)
APP_NAME	?= $(shell basename $(PACKAGE))
DOCKER_TAG  ?= $(shell echo $(PACKAGE) | awk -F '/' '{print $$1 "/" $$1}')
DATE		?= $(shell date +%FT%T%z)
VERSION		?= $(shell git describe --tags --exact-match 2> /dev/null || \
					git symbolic-ref -q --short HEAD 2> /dev/null || \
					cat $(CURDIR)/.version 2> /dev/null || echo v0)

GOPATH 	= $(shell echo $${GOPATH} | awk -F : '{print $$NF}')
BASE 	= $(GOPATH)/src/$(PACKAGE)
PKGS 	= $(or $(PKG),$(shell cd $(BASE) && env GOPATH=$(GOPATH) $(GO) list ./...))

# Tools
GO      	= go
GOFMT   	= gofmt
TIMEOUT 	= 180

.PHONY: status

status:
	@echo "PACKAGE: $(PACKAGE)"
	@echo "APP_NAME: $(APP_NAME)"
	@echo "DOCKER_TAG: 	$(DOCKER_TAG)"
	@echo "DATE: $(DATE)"
	@echo "VERSION: $(VERSION)"
	@echo "GOPATH: $(GOPATH)"
	@echo "BASE: $(BASE)"
	@echo "PKGS: $(PKGS)"
	@echo "GO: $(GO)"
	@echo "GOFMT: $(GOFMT)"
	@echo "TIMEOUT: $(TIMEOUT)"

.PHONY: docker

docker:	docker		## Build docker image
	docker build \
		--pull \
		--network=host \
		--build-arg PACKAGE=$(PACKAGE) \
		-t $(DOCKER_TAG):$(VERSION) \
		--label "COMMIT_SHA=$$(git rev-parse HEAD)" \

