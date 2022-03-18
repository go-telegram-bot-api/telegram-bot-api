# Makefile

## HELP:
.PHONY: help
## help: Show this help message.
help:
	@echo "Usage: make [target]\n"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

## :
## BUILD:

.PHONY: build
## build: Build Go code.
build:
	go build -o /dev/null ./...

## :
## DEPENDENCIES:

.PHONY: dep-clean
## dep-clean: Clean up dependency files.
dep-clean:
	@rm go.mod go.sum

.PHONY: dep-get
## dep-get: Get Go modules.
dep-get:
	go mod tidy

.PHONY: dep-init
## dep-init: Initialize Go modules.
dep-init:
	go mod init

.PHONY: dep-update
## dep-update: Update Go modules.
dep-update:
	go get -u ./...
	go mod tidy

## :
## TEST:

.PHONY: clean
## clean: Delete test output files.
clean:
	@rm coverage.out

.PHONY: test
## test: Run Go tests.
test:
	go test -coverprofile=coverage.out ./...

## :
