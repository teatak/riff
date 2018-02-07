GOTOOLS = \
	github.com/elazarl/go-bindata-assetfs/... \
	github.com/jteeuwen/go-bindata/...

VERSION = $(shell cat version)
GITSHA=$(shell git rev-parse HEAD)
GITBRANCH=$(shell git rev-parse --abbrev-ref HEAD)

# Build the project
default: tools webpack assets
	@echo "--> Running build"
	@sh -c "$(CURDIR)/scripts/build.sh"

dev: webpack assets
	@echo "--> Running build"
	@DEV=1 sh -c "'$(CURDIR)/scripts/build.sh'"

fmt:
	@cd $(CURDIR) ; \
	go fmt $$(go list ./... | grep -v /vendor/)

test: tools dev
	@echo "--> Running go test"
	go list ./... | grep -v -E '^github.com/gimke/riff/(vendor|cmd/riff/vendor)' | xargs -n1 go test

webpack:
	@echo "--> Running webpack"
	@npm run product

assets:
	@echo "--> Running assets"
	@go-bindata-assetfs -ignore .DS_Store -pkg riff ./static/...
	@mv bindata_assetfs.go riff/
	@cd $(CURDIR) ; \
	go fmt $$(go list ./... | grep -v /vendor/)

tools:
	@echo "--> Running tools"
	@go get -u -v $(GOTOOLS)

.PHONY: default fmt
