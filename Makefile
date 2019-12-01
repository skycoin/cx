.DEFAULT_GOAL := help
.PHONY: build-parser build build-full test test-full update-golden-files
.PHONY: install-gfx-deps install-gfx-deps-Darwin install-gfx-deps-Linux install-deps install

PWD := $(shell pwd)
# PKG_NAMES_LINUX := glade xvfb libxinerama-dev libxcursor-dev libxrandr-dev libgl1-mesa-dev libxi-dev gir1.2-gtk-3.0 libgtk2.0-dev libperl-dev libcairo2-dev libpango1.0-dev libgtk-3-dev gtk+3.0 libglib2.0-dev
PKG_NAMES_LINUX := glade xvfb libxinerama-dev libxcursor-dev libxrandr-dev libgl1-mesa-dev libxi-dev libperl-dev libcairo2-dev libpango1.0-dev libglib2.0-dev
# PKG_NAMES_MACOS := gtk gtk-mac-integration gtk+3 glade
UNAME_S := $(shell uname -s)
INSTALL_GFX_DEPS := install-gfx-deps-$(UNAME_S)

GLOBAL_GOPATH := ${GOPATH}
LOCAL_GOPATH  := ${HOME}/go
CXPATH        := ${CXPATH}

ifdef GLOBAL_GOPATH
  GOPATH := $(GLOBAL_GOPATH)
else
  GOPATH := $(LOCAL_GOPATH)
endif

ifdef CXPATH
	CX_PATH := $(CXPATH)
else
	CX_PATH := ${HOME}/cx
endif

ifeq ($(UNAME_S), Linux)
  DISPLAY       := :99.0
endif

configure: ## Configure the system to build and run CX
	@if [ -z "$(GLOBAL_GOPATH)" ]; then echo "NOTE:\tGOPATH not set" ; export GOPATH="$(LOCAL_GOPATH)"; export PATH="$(LOCAL_GOPATH)/bin:${PATH}" ; fi
	@echo "GOPATH=$(GOPATH)"
	@mkdir -p $(GOPATH)/src/github.com/SkycoinProject
	@if [ ! -e $(GOPATH)/src/github.com/SkycoinProject/cx ]; then mkdir -p $(GOPATH)/src/github.com/SkycoinProject ; ln -s $(PWD) $(GOPATH)/src/github.com/SkycoinProject/cx ; fi

configure-workspace: ## Configure CX workspace environment
	mkdir -p $(CX_PATH)/{,src,bin,pkg}
	@echo "NOTE:\tCX workspace at $(CX_PATH)"

build-parser: configure install-deps ## Generate lexer and parser for CX grammar
	nex -e cxgo/cxgo0/cxgo0.nex
	goyacc -o cxgo/cxgo0/cxgo0.go cxgo/cxgo0/cxgo0.y
	nex -e cxgo/parser/cxgo.nex
	goyacc -o cxgo/parser/cxgo.go cxgo/parser/cxgo.y

build: configure build-parser ## Build CX from sources
	go build -tags base -i -o $(GOPATH)/bin/cx github.com/SkycoinProject/cx/cxgo/
	chmod +x $(GOPATH)/bin/cx

build-full: configure build-parser ## Build CX from sources with all build tags
	go build -tags="base opengl" -i -o $(GOPATH)/bin/cx github.com/SkycoinProject/cx/cxgo/
	chmod +x $(GOPATH)/bin/cx

install-gfx-deps-Linux:
	@echo 'Installing dependencies for $(UNAME_S)'
	sudo apt-get update -qq
	sudo apt-get install -y $(PKG_NAMES_LINUX) --no-install-recommends
	export DISPLAY=$(DISPLAY)
	sudo /usr/bin/Xvfb ${DISPLAY} 2>1 > /dev/null &
	export GTK_VERSION="$(shell pkg-config --modversion gtk+-3.0 | tr . _| cut -d '_' -f 1-2)"
	export Glib_VERSION="$(shell pkg-config --modversion glib-2.0)"
	export Cairo_VERSION="$(shell pkg-config --modversion cairo)"
	export Pango_VERSION="$(shell pkg-config --modversion pango)"

install-gfx-deps-Darwin:
	# echo 'Installing dependencies for $(UNAME_S)'
	# brew install $(PKG_NAMES_MACOS)

install-deps: configure
	@echo "Installing go package dependencies"
	go get github.com/blynn/nex
	go get github.com/cznic/goyacc

install-gfx-deps: configure $(INSTALL_GFX_DEPS)
	go get github.com/go-gl/gl/v2.1/gl
	go get github.com/go-gl/glfw/v3.2/glfw
	go get github.com/go-gl/gltext

install: install-deps build configure-workspace ## Install CX from sources. Build dependencies
	@echo 'NOTE:\tWe recommend you to test your CX installation by running "cx ${GOPATH}/src/github.com/SkycoinProject/cx/tests"'
	cx -v

install-linters: ## Install linters
	go get -u github.com/FiloSottile/vendorcheck
	# For some reason this install method is not recommended, see https://github.com/golangci/golangci-lint#install
	# However, they suggest `curl ... | bash` which we should not do
	go get -u github.com/golangci/golangci-lint/cmd/golangci-lint

lint: ## Run linters. Use make install-linters first.
	vendorcheck ./...
	golangci-lint run -c .golangci.yml ./cx

test: build ## Run CX test suite.
	go test -race -tags base github.com/SkycoinProject/cx/cxgo/
	cx ./tests/main.cx ++wdir=./tests ++disable-tests=gui,issue

test-full: build ## Run CX test suite with all build tags
	go test -race -tags="base opengl" github.com/SkycoinProject/cx/cxgo/
	cx ./tests/main.cx ++wdir=./tests ++disable-tests=gui,issue

update-golden-files: build ## Update golden files used in CX test suite
	ls -1 tests/ | grep '.cx$$' | while read -r NAME; do echo "Processing $$NAME"; cx -t -co tests/testdata/tokens/$${NAME}.txt tests/$$NAME || true ; done

check-golden-files: update-golden-files ## Ensure golden files are up to date
	if [ "$(shell git diff tests/testdata | wc -l | tr -d ' ')" != "0" ] ; then echo 'Changes detected. Goden files not up to date' ; exit 2 ; fi

check: check-golden-files test ## Perform self-tests

format: ## Formats the code. Must have goimports installed (use make install-linters).
	goimports -w -local github.com/SkycoinProject/cx ./cx
	goimports -w -local github.com/SkycoinProject/cx ./cxgo/actions

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

