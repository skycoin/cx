.DEFAULT_GOAL := help
.PHONY: build-parser build test update-golden-files
.PHONY: install-deps-Darwin install-deps-Linux install-deps install

PWD := $(shell pwd)
PKG_NAMES_LINUX := glade xvfb libxinerama-dev libxcursor-dev libxrandr-dev libgl1-mesa-dev libxi-dev gir1.2-gtk-3.0 libgtk2.0-dev libperl-dev libcairo2-dev libpango1.0-dev libgtk-3-dev gtk+3.0 libglib2.0-dev
PKG_NAMES_MACOS := gtk gtk-mac-integration gtk+3 glade
UNAME_S := $(shell uname -s)
INSTALL_DEPS := install-deps-$(UNAME_S)

ifeq ($(UNAME_S), Linux)
  DISPLAY       := :99.0
endif

configure: ## Configure the system to build and run CX
	if [ -z "${GOPATH+x}" ]; then echo "NOTE:\tGOPATH not set" ; export GOPATH="${HOME}/go"; export PATH="${GOPATH}/bin:${PATH}" ; fi
	mkdir -p ${GOPATH}
	echo "GOPATH=${GOPATH}"
	if [ ! -d ${GOPATH}/src/github.com/skycoin/cx ]; then mkdir -p ${GOPATH}/src/github.com/skycoin ; ln -s $(PWD) ${GOPATH}/src/github.com/skycoin/cx ; fi

configure-workspace: ## Configure CX workspace environment
	if [ -z "${CXPATH+x}" ]; then export CX_PATH="${HOME}/cx" ; else export CX_PATH=${CXPATH} ; fi
	mkdir -p ${CX_PATH}/{,src,bin,pkg}
	echo "NOTE:\tCX workspace at ${CX_PATH}"

build-parser: configure ## Generate lexer and parser for CX grammar
	nex -e cxgo/cxgo0/cxgo0.nex
	goyacc -o cxgo/cxgo0/cxgo0.go cxgo/cxgo0/cxgo0.y
	nex -e cxgo/cxgo.nex
	goyacc -o cxgo/cxgo.go cxgo/cxgo.y

build: configure build-parser ## Build CX from sources
	go build -tags full -i -o ${GOPATH}/bin/cx github.com/skycoin/cx/cxgo/
	chmod +x ${GOPATH}/bin/cx

install-deps-Linux:
	echo 'Installing dependencies for $(UNAME_S)'
	sudo apt-get update -qq
	sudo apt-get install -y $(PKG_NAMES_LINUX) --no-install-recommends
	export DISPLAY=$(DISPLAY)
	sudo /usr/bin/Xvfb ${DISPLAY} 2>1 > /dev/null &
	export GTK_VERSION="$(shell pkg-config --modversion gtk+-3.0 | tr . _| cut -d '_' -f 1-2)"
	export Glib_VERSION="$(shell pkg-config --modversion glib-2.0)"
	export Cairo_VERSION="$(shell pkg-config --modversion cairo)"
	export Pango_VERSION="$(shell pkg-config --modversion pango)"
	echo "GTK version ${GTK_VERSION} (Glib ${Glib_VERSION}, Cairo ${Cairo_VERSION}, Pango ${Pango_VERSION})"

install-deps-Darwin:
	echo 'Installing dependencies for $(UNAME_S)'
	brew install $(PKG_NAMES_MACOS)

install-deps: configure $(INSTALL_DEPS)
	echo "Installing go package dependencies"
	go get github.com/skycoin/skycoin/...
	go get github.com/go-gl/gl/v2.1/gl
	go get github.com/go-gl/glfw/v3.2/glfw
	go get github.com/go-gl/gltext
	go get github.com/blynn/nex
	go get github.com/cznic/goyacc
#	go get github.com/skycoin/cx/...

install: install-deps build configure-workspace ## Install CX from sources. Build dependencies
	echo 'NOTE:\tWe recommend you to test your CX installation by running "cx ${GOPATH}/src/github.com/skycoin/cx/tests"'
	cx -v

test: build ## Run CX test suite.
	go test -race -tags full github.com/skycoin/cx/cxgo/
	cx ./tests/main.cx ++wdir=./tests ++disable-tests=gui,issue

update-golden-files: ## Update golden files used in CX test suite
	# TODO: Implement

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

