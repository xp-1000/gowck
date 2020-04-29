export GO111MODULE=on

# Project related variables.
VERSION := $(shell git describe --tags --always --dirty)
BUILD := $(shell git rev-parse --short HEAD)
PROJECTNAME := $(shell basename "$(PWD)")

# Go related variables.
GOBASE := $(shell pwd)
GOBIN := $(GOBASE)/bin
NUM_CORES ?= $(shell getconf _NPROCESSORS_ONLN)

# Make is verbose in Linux. Make it silent.
MAKEFLAGS += --silent

# Use linker flags to provide version/build settings
LDFLAGS=-ldflags "-X=main.Version=$(VERSION) -X=main.Build=$(BUILD)"

.PHONY: install
## install: Compile binary and install it in ./bin
install:
	@echo ' > Build binary'
	GOBIN=$(GOBIN) go build -trimpath $(LDFLAGS) -o $(GOBIN)/$(PROJECTNAME)

.PHONY: clean
## clean: Clean build files. Runs `go clean` internally.
clean:
	@echo ' > Cleaning build files'
	@-rm $(GOBIN)/$(PROJECTNAME) 2> /dev/null
	@echo ' > Run go clean'
	GOBIN=$(GOBIN) go clean

.PHONY: check
## check: Run lint, vet and test
check: lint test

## test: Run go test
.PHONY: test
test:
	@echo ' > Run go test'
	CGO_ENABLED=0 go test -p $(NUM_CORES) ./...

## lint: Linting go files
.PHONY: lint
lint:
	@which golangci-lint 1>/dev/null || (echo " > Installing golangci-lint" && GO111MODULE=on go get github.com/golangci/golangci-lint/cmd/golangci-lint)
	@echo ' > Linting code'
	CGO_ENABLED=0 GOGC=40 golangci-lint run -E misspell -E golint -E gofmt --deadline 5m
 
## fmt: Format go files
.PHONY: fmt
fmt:
	@echo ' > Format code'
	CGO_ENABLED=0 go fmt ./...

## tidy: Clean unused modules
.PHONY: tidy
tidy:
	@echo ' > Run go mod tiny'
	go mod tidy

.PHONY: help
all: help
help: Makefile
	@echo
	@echo " Choose a command run in "$(PROJECTNAME)":"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo
