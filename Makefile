.DEFAULT_GOAL := help
.PHONY: build-parser build build-full test test-full update-golden-files
.PHONY: install-gfx-deps install-gfx-deps-LINUX install-gfx-deps-MSYS install-gfx-deps-MINGW install-gfx-deps-MACOS install-deps install install-full

PWD := $(shell pwd)

#PKG_NAMES_LINUX := glade xvfb libxinerama-dev libxcursor-dev libxrandr-dev libgl1-mesa-dev libxi-dev gir1.2-gtk-3.0 libgtk2.0-dev libperl-dev libcairo2-dev libpango1.0-dev libgtk-3-dev gtk+3.0 libglib2.0-dev
PKG_NAMES_LINUX := glade xvfb libxinerama-dev libxcursor-dev libxrandr-dev libgl1-mesa-dev libxi-dev libperl-dev libcairo2-dev libpango1.0-dev libglib2.0-dev libopenal-dev
#PKG_NAMES_MACOS := gtk gtk-mac-integration gtk+3 glade
PKG_NAMES_WINDOWS := mingw-w64-x86_64-openal

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

#ifneq (,$(findstring CYGWIN, $(UNAME_S)))
#PLATFORM := WINDOWS
#SUBSYSTEM := CYGWIN
#endif

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

INSTALL_GFX_DEPS := install-gfx-deps-$(SUBSYSTEM)

GLOBAL_GOPATH := $(GOPATH)
LOCAL_GOPATH  := $(HOME)/go

ifdef GLOBAL_GOPATH
  GOPATH := $(GLOBAL_GOPATH)
else
  GOPATH := $(LOCAL_GOPATH)
endif

ifdef CXPATH
	CX_PATH := $(CXPATH)
else
	CX_PATH := $(HOME)/cx
endif

ifeq ($(UNAME_S), Linux)
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
	go build -tags="base" -i -o $(GOPATH)/bin/cx github.com/SkycoinProject/cx/cxgo/
	chmod +x $(GOPATH)/bin/cx

build-full: configure build-parser ## Build CX from sources with all build tags
	go build -tags="base cxfx" -i -o $(GOPATH)/bin/cx github.com/SkycoinProject/cx/cxgo/
	chmod +x $(GOPATH)/bin/cx

build-android: configure build-parser
#go get github.com/SkycoinProject/gltext
	git clone https://github.com/SkycoinProject/gomobile $(GOPATH)/src/golang.org/x/mobile 2> /dev/null || true
	cd $(GOPATH)/src/golang.org/x/mobile/; git pull origin master; go get ./cmd/gomobile
	cp -R $(GOPATH)/src/github.com/SkycoinProject/cxfx/resources/fonts/ $(GOPATH)/src/github.com/SkycoinProject/cx/cxgo/assets/cxfx/resources/fonts/
	cp -R $(GOPATH)/src/github.com/SkycoinProject/cxfx/resources/shaders/ $(GOPATH)/src/github.com/SkycoinProject/cx/cxgo/assets/cxfx/resources/shaders/
	cp -R $(GOPATH)/src/github.com/SkycoinProject/cxfx/tutorials/ $(GOPATH)/src/github.com/SkycoinProject/cx/cxgo/assets/cxfx/tutorials/
	cp -R $(GOPATH)/src/github.com/SkycoinProject/cxfx/src/ $(GOPATH)/src/github.com/SkycoinProject/cx/cxgo/assets/cxfx/src/
	cp -R $(GOPATH)/src/github.com/SkycoinProject/cxfx/games/skylight/src/ $(GOPATH)/src/github.com/SkycoinProject/cx/cxgo/assets/cxfx/games/skylight/src/
	gomobile install -tags="base cxfx mobile android_gles31" -target=android $(GOPATH)/src/github.com/SkycoinProject/cx/cxgo/
#-DHOST=armv7a-linux-androideabi29

install-gfx-deps-LINUX:
	@echo 'Installing dependencies for $(UNAME_S)'
	sudo apt-get update -qq
	sudo apt-get install -y $(PKG_NAMES_LINUX) --no-install-recommends
#	export DISPLAY=$(DISPLAY)
#	sudo /usr/bin/Xvfb ${DISPLAY} 2>1 > /dev/null &
#	export GTK_VERSION="$(shell pkg-config --modversion gtk+-3.0 | tr . _| cut -d '_' -f 1-2)"
#	export Glib_VERSION="$(shell pkg-config --modversion glib-2.0)"
#	export Cairo_VERSION="$(shell pkg-config --modversion cairo)"
#	export Pango_VERSION="$(shell pkg-config --modversion pango)"

install-gfx-deps-MSYS:
	@echo 'Installing dependencies for $(UNAME_S)'
	pacman -Sy
	pacman -S $(PKG_NAMES_WINDOWS)
	ln -s /mingw64/lib/libopenal.a /mingw64/lib/libOpenAL32.a
	ln -s /mingw64/lib/libopenal.dll.a /mingw64/lib/libOpenAL32.dll.a

install-gfx-deps-MINGW: install-gfx-deps-MSYS

install-gfx-deps-MACOS:
	@echo 'Installing dependencies for $(UNAME_S)'
#brew install $(PKG_NAMES_MACOS)

install-deps: configure
	@echo "Installing go package dependencies"
	go get github.com/SkycoinProject/nex
	go get github.com/cznic/goyacc

install-gfx-deps: configure $(INSTALL_GFX_DEPS)
	go get github.com/SkycoinProject/gltext
	go get github.com/go-gl/gl/v3.2-compatibility/gl
	go get github.com/go-gl/glfw/v3.3/glfw
	go get golang.org/x/mobile/exp/audio/al
	go get github.com/mjibson/go-dsp/wav

install: install-deps build configure-workspace ## Install CX from sources. Build dependencies
	@echo 'NOTE:\tWe recommend you to test your CX installation by running "cx $(GOPATH)/src/github.com/SkycoinProject/cx/tests"'
	cx -v

install-full: install-gfx-deps install

install-mobile:
	go get golang.org/x/mobile/gl

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
	cx ./lib/args.cx ./tests/main.cx ++wdir=./tests ++disable-tests=gui,issue

test-full: build ## Run CX test suite with all build tags
	go test -race -tags="base cxfx" github.com/SkycoinProject/cx/cxgo/
	cx ./lib/args.cx ./tests/main.cx ++wdir=./tests ++disable-tests=gui,issue

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
