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


## Ensure $GOBIN is set.
GOLANGCI_LINT_VERSION ?= latest

GO_OPTS ?= GOBIN=$(GOBIN)

ifdef CXPATH
	CX_PATH := $(CXPATH)
else
	CX_PATH := $(HOME)/cx
endif

ifeq ($(UNAME_S), Linux)
endif

configure-workspace: ## Configure CX workspace environment
	mkdir -p $(CX_PATH)/src $(CX_PATH)/bin $(CX_PATH)/pkg
	@echo "NOTE:\tCX workspace at $(CX_PATH)"

build-parser: install-deps ## Generate lexer and parser for CX grammar
	$(GOBIN)/goyacc -o cxgo/cxgo0/cxgo0.go cxgo/cxgo0/cxgo0.y
	$(GOBIN)/goyacc -o cxgo/parser/cxgo.go cxgo/parser/cxgo.y

build:  ## Build CX from sources
	$(GO_OPTS) go build -tags="base" -i -o $(GOBIN)/cx github.com/skycoin/cx/cxgo/
	chmod +x $(GOBIN)/cx

build-full: install-full  ## Build CX from sources with all build tags
	$(GO_OPTS) go build -tags="base cxfx" -i -o $(GOBIN)/cx github.com/skycoin/cx/cxgo/
	chmod +x $(GOBIN)/cx

build-android: install-full install-mobile 
	# TODO @evanlinjin: We should switch this to use 'github.com/SkycoinProject/gomobile' once it can build.
	$(GO_OPTS) go get -u golang.org/x/mobile/cmd/gomobile

install-deps:
	@echo "Installing go package dependencies"
	$(GO_OPTS) go get -u modernc.org/goyacc

install: install-deps build configure-workspace ## Install CX from sources. Build dependencies
	@echo 'NOTE:\tWe recommend you to test your CX installation by running "cx ./tests"'
	$(GOBIN)/cx -v

install-full: install-deps configure-workspace

install-mobile:
	$(GO_OPTS) go get golang.org/x/mobile/gl # TODO @evanlinjin: This is a library. needed?

clean: ## Removes binaries. 
	rm -r $(GOBIN)/cx

token-fuzzer:
	$(GO_OPTS) go build -i -o $(GOBIN)/cx-token-fuzzer $(PWD)/development/token-fuzzer/main.go
	chmod +x ${GOPATH}/bin/cx-token-fuzzer

test: #build ## Run CX test suite.
ifndef CXVERSION
	@echo "cx not found in $(PWD)/bin, please run make install first"
else	
	$(GO_OPTS) go test -race -tags base github.com/skycoin/cx/cxgo/
	$(GOBIN)/cx ./lib/args.cx ./tests/main.cx ++wdir=./tests ++disable-tests=gui,issue
endif

test-full: build ## Run CX test suite with all build tags
	$(GO_OPTS) go test -race -tags="base cxfx" github.com/skycoin/cx/cxgo/
	$(GOBIN)/cx ./lib/args.cx ./tests/main.cx ++wdir=./tests ++disable-tests=gui,issue

check: test ## Perform self-tests

format: ## Formats the code. Must have goimports installed (use make install-linters).
	goimports -w -local github.com/skycoin/cx ./cx
	goimports -w -local github.com/skycoin/cx ./cxfx
	goimports -w -local github.com/skycoin/cx ./cxgo

dep: ## Update go vendor
	$(GO_OPTS) go mod vendor
	$(GO_OPTS) go mod verify
	$(GO_OPTS) go mod tidy

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
