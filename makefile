
VERSION = $(shell cat version)
GITSHA=$(shell git rev-parse HEAD)
GITBRANCH=$(shell git rev-parse --abbrev-ref HEAD)

# Build the project
default:
	@sh -c "$(CURDIR)/scripts/build.sh"

dev:
	@DEV=1 sh -c "'$(CURDIR)/scripts/build.sh'"

fmt:
	@cd $(CURDIR) ; \
	go fmt $$(go list ./... | grep -v /vendor/)

test: dev
	@echo "--> Running go test"
	@go list ./... | grep -v -E '^github.com/gimke/riff/(vendor|cmd/serf/vendor)' | xargs -n1 go test

.PHONY: default fmt
