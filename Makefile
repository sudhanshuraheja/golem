SHELL := /bin/bash

GO ?= go
BIN ?= bin/golem
BIN_DIR := $(dir $(BIN))
VERSION_FILE := golem/version.go
COVER_PROFILE := coverage.out
COVER_HTML := coverage.html
COVER_FUNC := coverage.func

.PHONY: help version build install test lint fmt deps clean

help:
	@echo "Common targets:"
	@echo "  make build    Build the binary to $(BIN)"
	@echo "  make install  Install the binary into GOPATH/bin"
	@echo "  make test     Run tests with race detector and coverage reports"
	@echo "  make lint     Run go vet"
	@echo "  make fmt      Run go fmt on the module"
	@echo "  make deps     Upgrade module dependencies"
	@echo "  make clean    Remove build and coverage artifacts"

version:
	./version.sh

$(BIN): version
	@mkdir -p $(BIN_DIR)
	$(GO) build -o $(BIN) ./...

build: $(BIN)

install: version
	$(GO) install ./...

test: version
	$(GO) test ./... -timeout 15s -race -cover -coverprofile=$(COVER_PROFILE)
	$(GO) tool cover -html=$(COVER_PROFILE) -o $(COVER_HTML)
	$(GO) tool cover -func=$(COVER_PROFILE) -o $(COVER_FUNC)

lint: version
	$(GO) vet ./...

fmt:
	$(GO) fmt ./...

deps:
	$(GO) get -u ./...
	$(GO) mod tidy

clean:
	rm -f $(BIN) $(COVER_PROFILE) $(COVER_HTML) $(COVER_FUNC)
