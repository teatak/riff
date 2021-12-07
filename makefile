GOTOOLS = \
	github.com/elazarl/go-bindata-assetfs \
	github.com/go-bindata/go-bindata

# Build the project
default: tools webpack assets
	@echo "--> Running build"
	@sh -c "$(CURDIR)/scripts/build.sh"

dev: assets
	@echo "--> Running build"
	@DEV=1 sh -c "'$(CURDIR)/scripts/build.sh'"

assets:
	@echo "--> Running assets"
	@go-bindata -ignore \\.go -pkg schema -o ./schema/bindata.go ./schema/...
	@go-bindata-assetfs -ignore .DS_Store -pkg riff -o ./riff/bindata_assetfs.go  ./static/...
	@cd $(CURDIR)
	@go fmt $$(go list ./... | grep -v /vendor/)

fmt:
	@cd $(CURDIR)
	@go fmt $$(go list ./... | grep -v /vendor/)

test: tools dev
	@echo "--> Running go test"
	go list ./... | grep -v -E '^github.com/teatak/riff/(vendor|cmd/riff/vendor)' | xargs -n1 go test

webpack:
	@echo "--> Running webpack"
	@cd console && npm run product

tools:
	@echo "--> Running tools"
	@go get $(GOTOOLS)
	@go install $(GOTOOLS)
	@cd console && npm install

.PHONY: default fmt
