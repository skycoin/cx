export GO111MODULE=on


.DEFAULT_GOAL := help
.PHONY: build-parser build test test-parser build-core
.PHONY: install
.PHONY: dep

PWD := $(shell pwd)
PTR := ptr32
##PTR := ptr64

UNAME_S := $(shell uname -s)


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
CXEXE := cx.exe
else
CXEXE := cx
endif

CXVERSION := $(shell $(PWD)/bin/$(CXEXE) --version 2> /dev/null)

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

build: ## Build CX from sources
	$(GO_OPTS) go build -tags="$(PTR) cipher cxfx cxos http regexp" -o ./bin/$(CXEXE) github.com/skycoin/cx/cmd/cx
	chmod +x ./bin/$(CXEXE)

build-debug: ## Build CX from sources
	$(GO_OPTS) go build -gcflags="all=-N -l" -tags="$(PTR) cipher cxfx cxos http regexp" -o ./bin/$(CXEXE) github.com/skycoin/cx/cmd/cx
	chmod +x ./bin/$(CXEXE)

build-core: ## Build CX with CXFX support. Done via satisfying 'cxfx' build tag.
	$(GO_OPTS) go build -tags="$(PTR) base" -o ./bin/$(CXEXE) github.com/skycoin/cx/cmd/cx
	chmod +x ./bin/$(CXEXE)

build-tests: build-debug
	go build -gcflags="all=-N -l" -o ./bin/cxtests github.com/skycoin/cx/cmd/cxtest/
	chmod +x  ./bin/cxtests

clean: ## Removes binaries.
	rm -rf ./bin/cx
	rm -rf ./bin/$(CXEXE)
	rm -rf ./bin/cxtests

install: configure-workspace ## Install CX from sources. Build dependencies
	@echo 'NOTE:\tWe recommend you to test your CX installation by running "$(CXEXE) ./tests"'
	./bin/$(CXEXE) -v

test-parser: build-parser build test

test:  ## Run CX test suite.
ifndef CXVERSION
	@echo "$(CXEXE) not found in $(PWD)/bin, please run make install first"
else
	# go test $(GO_OPTS) -race -tags base github.com/skycoin/cx/cxgo/
	go run -mod=vendor ./cmd/cxtest/main.go --cxpath=$(PWD)/bin/$(CXEXE) --wdir=./tests --log=fail,stderr --disable-tests=gui,issue
endif

test-all:  ## Run CX test suite.
ifndef CXVERSION
	@echo "$(CXEXE) not found in $(PWD)/bin, please run make install first"
else
	# go test $(GO_OPTS) -race -tags base github.com/skycoin/cx/cxgo/
	go run -mod=vendor ./cmd/cxtest/main.go --cxpath=$(PWD)/bin/$(CXEXE) --wdir=./tests --log=fail,stderr
endif

build-goyacc: ## Builds goyacc into /bin/goyacc
	go build -o ./bin/goyacc ./cmd/goyacc/main.go

build-parser: ## Generate lexer and parser for CX grammar
	#go build -o ./bin/goyacc ./cmd/goyacc/main.go
	./bin/goyacc -o cxparser/cxpartialparsing/cxpartialparsing.go cxparser/cxpartialparsing/cxpartialparsing.y
	./bin/goyacc -o cxparser/cxparsingcompletor/parsingcompletor.go cxparser/cxparsingcompletor/parsingcompletor.y

token-fuzzer:
	go build $(GO_OPTS) -o ./bin/cx-token-fuzzer $(PWD)/development/token-fuzzer/main.go
	chmod +x ./bin/cx-token-fuzzer

configure-workspace: ## Configure CX workspace environment
	mkdir -p $(CX_PATH)/src $(CX_PATH)/bin $(CX_PATH)/pkg
	@echo "NOTE:\tCX workspace at $(CX_PATH)"

format: ## Formats the code. Must have goimports installed (use make install-linters).
	goimports -w -local github.com/skycoin/cx ./cmd
	goimports -w -local github.com/skycoin/cx ./cx
	goimports -w -local github.com/skycoin/cx ./cxfx
	goimports -w -local github.com/skycoin/cx ./cxgo

dep: ## Update go vendor
	go mod vendor
	go mod verify
	go mod tidy

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

benchmark: #Benchmark
ifndef CXVERSION
	@echo "$(CXEXE) not found in $(PWD)/bin, please run make install first"
else
	mkdir -p $(PWD)/cmd/cxbenchmark/bin/
	rm -f $(PWD)/cmd/cxbenchmark/bin/$(CXEXE)
	cp $(PWD)/bin/$(CXEXE) $(PWD)/cmd/cxbenchmark/bin/
	go test $(PWD)/cmd/cxbenchmark -run BenchmarkCX -tags ptr32 -bench=.
endif