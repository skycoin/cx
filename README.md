![cx logo](https://user-images.githubusercontent.com/26845312/32426758-2a4bbb00-c282-11e7-858e-a1eaf3ea92f3.png)

# CX Programming Language
 
[![Build Status](https://travis-ci.com/skycoin/cx.svg?branch=develop)](https://travis-ci.com/skycoin/cx) [![Build status](https://ci.appveyor.com/api/projects/status/y04pofhhfmpw8vef/branch/master?svg=true)](https://ci.appveyor.com/project/skycoin/cx/branch/master)

CX is a general purpose, interpreted and compiled programming
language, with a very strict type system and a syntax
similar to Golang's. CX provides a new programming paradigm based on
the concept of affordances.

## Table of Contents

   * [CX Programming Language](#cx-programming-language-1)
   * [CX Chains (CX Skycoin Blockchain)](#cx-chains-cx--skycoin-blockchain)
   * [Compiler Development](CompilerDevelopment.md)
   * [Installation](#installation)
      * [Binary Releases](#binary-releases)  
      * [Compiling from Source](#compiling-from-source)
         * [Installing Go](#installing-go)
         * [Compiling CX on *nix](#compiling-cx-on-nix)
         * [Compiling CX on Windows](#compiling-cx-on-windows)
      * [Updating CX](#updating-cx)
   * [Resources and libraries](#resources-and-libraries)
   * [Running CX](#running-cx)
      * [Hello World](#hello-world)
      * [Basic Options](#other-options)
         * [Running CX Programs](#running-cx-programs)
      * [REPL tutorial](#cx-repl)
   * [Learning CX](#learning-cx)

# CX Programming Language
[[Back to the Table of Contents] ↑](#table-of-contents)

CX is a general purpose, interpreted and compiled programming
language, with a very strict type system and a syntax
similar to Golang's. CX provides a new programming paradigm based on
the concept of affordances, where the user can ask the programming
language at runtime what can be done with a CX object (functions,
expressions, packages, etc.), and interactively or automatically choose
one of the affordances to be applied. This paradigm has the main objective
of providing an additional security layer for decentralized,
blockchain-based applications, but can also be used for general
purpose programming. 

## Installation

CX requires a Golang version of `1.15` or higher. 

### Binary Releases

You can find binary releases for most major systems on the [release page](https://github.com/skycoin/cx/releases). 

### Compiling on Linux

If you are using a `apt` compatible system, install the dependencies with"

```
sudo apt-get update -qq

sudo apt-get install -y glade xvfb libxinerama-dev libxcursor-dev libxrandr-dev libgl1-mesa-dev libxi-dev libperl-dev libcairo2-dev libpango1.0-dev libglib2.0-dev libopenal-dev libxxf86vm-dev --no-install-recommends
```

Download CX's repository using Go:

```
go get github.com/skycoin/cx
```

Get required Go dependencies with:

```
go get -u golang.org/x/mobile/cmd/gomobile
go get -u modernc.org/goyacc
go get golang.org/x/mobile/gl 
```

Navigate to CX's repository via:

```
cd $GOPATH/src/github.com/skycoin/cx
```

Build CX's binary and install by running:

```
make install
```

`make install` builds a CX binary and installs it into `$HOME/cx/bin`. Add the CX binary path to your operating system's `$PATH`. For example, in Linux:

```
export PATH=$PATH:$HOME/cx/bin
```

You should test your installation by running:

```
make test
```

If you intend to develop games with CX, then run:

```
make test-full
```

### Compiling on MacOS

Download CX's repository using Go:

```
go get github.com/skycoin/cx
```

Navigate to CX's repository via:

```
cd $GOPATH/src/github.com/skycoin/cx
```

Build CX's binary and install by running:

```
make install
```

`make install` builds a CX binary and installs it into `$HOME/cx/bin`. Add the CX binary path to your operating system's `$PATH`. For example, in Linux:

```
export PATH=$PATH:$HOME/cx/bin
```

You should test your installation by running:

```
make test
```

If you intend to develop games with CX, then run:

```
make test-full
```

### Compiling CX on Windows

Compiling CX on Windows requires a recent version of Git to be installed. 

Before compiling CX, install dependencies with `pacman`:

```
pacman -Sy

pacman -S mingw-w64-x86_64-openal

if [ ! -a /mingw64/lib/libOpenAL32.a]; then ln -s /mingw64/lib/libopenal.a /mingw64/lib/libOpenAL32.a; fi

if [ ! -a /mingw64/lib/libOpenAL32.dll.a]; then ln -s /mingw64/lib/libopenal.dll.a /mingw64/lib/libOpenAL32.dll.a; fi
```

You can compile CX by running: 

```
cx-setup.bat
```

Test your installation by running:

```
cx tests\main.cx ++wdir=tests ++disable-tests=issue
```

## Updating CX

You can update your CX installation by running:

```
make install
```

Or on Windows:

```
cx-setup.bat
```

## Resources and libraries

If you are interested in learning more about CX, please refer to the [resources documentation](docs/cx-resources.md). 

If you want to get started with some basic example programs and tutorials check out the [tutorials section](docs/cx-tutorials.md). 

The docs also provide a high level [overview](docs/overview.md) over the language. 