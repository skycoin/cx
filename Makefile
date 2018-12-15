.DEFAULT_GOAL := help
.PHONY: build-parser build test update-golden-files
.PHONY: install-deps-Darwin install-deps-Linux install-deps install

PKG_NAMES_LINUX := 'glade xvfb libxinerama-dev libxcursor-dev libxrandr-dev libgl1-mesa-dev libxi-dev gir1.2-gtk-3.0 libgtk2.0-dev libperl-dev libcairo2-dev libpango1.0-dev libgtk-3-dev gtk+3.0 libglib2.0-dev'
PKG_NAMES_MACOS := 'gtk gtk-mac-integration gtk+3 glade'
UNAME_S := $(shell uname -s)
INSTALL_DEPS := "build-deps-$(UNAME_S)"

ifeq ($(UNAME_S), Linux)
  DISPLAY       := ':99.0'
  GTK_VERSION   := $(pkg-config --modversion gtk+-3.0 | tr . _| cut -d '_' -f 1-2)
  GLIB_VERSION  := $(pkg-config --modversion glib-2.0)
  CAIRO_VERSION := $(pkg-config --modversion cairo)
  PANGO_VERSION := $(pkg-config --modversion pango)
endif

configure: ## Configure the system to build and run CX
	export PATH="${GOPATH}/bin:${PATH}"

build-parser: ## Generate lexer and parser for CX grammar
	nex -e cxgo/cxgo0/cxgo0.nex
	goyacc -o cxgo/cxgo0/cxgo0.go cxgo/cxgo0/cxgo0.y
	nex -e cxgo/cxgo.nex
	goyacc -o cxgo/cxgo.go cxgo/cxgo.y

build: configure build-parser ## Build CX from sources
	go build -tags full -i -o ${GOPATH}/bin/cx ./cxgo/

install-deps-Linux:
	sudo apt-get update -qq
	sudo apt-get install -y $(PKG_NAMES_LINUX) --no-install-recommends
	export DISPLAY=$(DISPLAY)
	sudo /usr/bin/Xvfb $(DISPLAY) 2>1 > /dev/null &
	export GTK_VERSION=$(GTK_VERSION)
	export Glib_VERSION=$(GLIB_VERSION)
	export Cairo_VERSION=$(CAIRO_VERSION)
	export Pango_VERSION=$(PANGO_VERSION)
	echo "GTK version ${GTK_VERSION} (Glib ${Glib_VERSION}, Cairo ${Cairo_VERSION}, Pango ${Pango_VERSION})"

install-deps-Darwin:
	brew install $(PKG_NAMES_MACOS)

install-deps: $(INSTALL_DEPS)
	go test -race -tags full -i github.com/skycoin/cx/cxgo/

install: configure install-deps ## Install CX from sources. Build dependencies
	source ./cx.sh
	go install -tags full -i -o ${GOPATH}/bin/cx ./cxgo/

test: build ## Run CX test suite.
	go test -race -tags full github.com/skycoin/cx/cxgo/
	cx ./tests/main.cx ++wdir=./tests ++disable-tests=gui,issue

update-golden-files: ## Update golden files used in CX test suite
	# TODO: Implement

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

