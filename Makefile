GO := go
pkgs  = $(shell $(GO) list ./... | grep -v vendor)

build:
	@./scripts/build.sh

fmt:
	@go fmt $(go list ./... | grep -v _base-operator) &> /dev/null

release : build
	@echo ">>> Built release"
	@rm -rf build/craft
	@mkdir -p build/craft
	@cp -r _base-operator bin init build/craft
	@cd build ; tar  -zcf  ../craft.tar.gz craft ; cd .. 

.PHONY: build format test check_format
