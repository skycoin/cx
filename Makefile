export GO111MODULE=on

.DEFAULT_GOAL := help
.PHONY: build-parser build build-full test test-full
.PHONY: install-deps install install-full
.PHONY: dep

PWD := $(shell pwd)

UNAME_S := $(shell uname -s)

CXVERSION := $(shell $(PWD)/bin/cx --version 2> /dev/null)

ifneq (,$(findstring Linux, $(UNAME_S)))
PLATFORM := LINUX
SUBSYSTEM := LINUX
PACKAGES := PGK_NAMES_LINUX
DISPLAY  := :99.0
endif

ifneq (,$(findstring Darwin, $(UNAME_S)))
PLATFORM := MACOS
SUBSYSTEM := MACOS
PACKAGES := PKG_NAMES_MACOS
endif

ifneq (,$(findstring MINGW, $(UNAME_S)))
PLATFORM := WINDOWS
SUBSYSTEM := MINGW
PACKAGES := PKG_NAMES_WINDOWS
endif

ifneq (,$(findstring MSYS, $(UNAME_S)))
PLATFORM := WINDOWS
SUBSYSTEM := MSYS
PACKAGES := PKG_NAMES_WINDOWS
endif

ifeq ($(PLATFORM), WINDOWS)
GOPATH := $(subst \,/,${GOPATH})
HOME := $(subst \,/,${HOME})
CXPATH := $(subst, \,/, ${CXPATH})
endif

GLOBAL_GOPATH := $(GOPATH)
LOCAL_GOPATH  := $(HOME)/go

ifdef GLOBAL_GOPATH
  GOPATH := $(GLOBAL_GOPATH)
else
  GOPATH := $(LOCAL_GOPATH)
endif

GOLANGCI_LINT_VERSION ?= latest

GO_OPTS ?= -mod=vendor

ifdef CXPATH
	CX_PATH := $(CXPATH)
else
	CX_PATH := $(HOME)/cx
endif

ifeq ($(UNAME_S), Linux)
endif

build:  ## Build CX from sources
	go build $(GO_OPTS) -tags="base" -i -o ./bin/cx github.com/skycoin/cx/cxgo/
	chmod +x ./bin/cx

clean: ## Removes binaries.
	rm -r ./bin/cx

build-full: install-full  ## Build CX from sources with all build tags
	go build $(GO_OPTS) -tags="base cxfx" -i -o ./bin/cx github.com/skycoin/cx/cxgo/
	chmod +x ./bin/cx

build-android: install-full install-mobile
	# TODO @evanlinjin: We should switch this to use 'github.com/SkycoinProject/gomobile' once it can build.
	go get $(GO_OPTS) -u golang.org/x/mobile/cmd/gomobile

token-fuzzer:
	go build $(GO_OPTS) -i -o ./bin/cx-token-fuzzer $(PWD)/development/token-fuzzer/main.go
	chmod +x ${GOPATH}/bin/cx-token-fuzzer

build-parser: install-deps ## Generate lexer and parser for CX grammar
	./bin/goyacc -o cxgo/cxgo0/cxgo0.go cxgo/cxgo0/cxgo0.y
	./bin/goyacc -o cxgo/parser/cxgo.go cxgo/parser/cxgo.y

install: install-deps build configure-workspace ## Install CX from sources. Build dependencies
	@echo 'NOTE:\tWe recommend you to test your CX installation by running "cx ./tests"'
	./bin/cx -v

install-full: install-deps configure-workspace

install-deps:
	@echo "Installing go package dependencies"
	go get $(GO_OPTS) -u modernc.org/goyacc

test:  ## Run CX test suite.
ifndef CXVERSION
	@echo "cx not found in $(PWD)/bin, please run make install first"
else
	go test $(GO_OPTS) -race -tags base github.com/skycoin/cx/cxgo/
	./bin/cx ./lib/args.cx ./tests/main.cx ++wdir=./tests ++disable-tests=gui,issue
endif

test-full: build ## Run CX test suite with all build tags
	go test $(GO_OPTS) -race -tags="base cxfx" github.com/skycoin/cx/cxgo/
	./bin/cx ./lib/args.cx ./tests/main.cx ++wdir=./tests ++disable-tests=gui,issue

configure-workspace: ## Configure CX workspace environment
	mkdir -p $(CX_PATH)/src $(CX_PATH)/bin $(CX_PATH)/pkg
	@echo "NOTE:\tCX workspace at $(CX_PATH)"

format: ## Formats the code. Must have goimports installed (use make install-linters).
	goimports -w -local github.com/skycoin/cx ./cx
	goimports -w -local github.com/skycoin/cx ./cxfx
	goimports -w -local github.com/skycoin/cx ./cxgo

dep: ## Update go vendor
	go mod $(GO_OPTS) vendor
	go mod $(GO_OPTS) verify
	go mod $(GO_OPTS) tidy

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
