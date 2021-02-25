![cx logo](https://user-images.githubusercontent.com/26845312/32426758-2a4bbb00-c282-11e7-858e-a1eaf3ea92f3.png)

# CX Programming Language

[![Build Status](https://travis-ci.com/skycoin/cx.svg?branch=develop)](https://travis-ci.com/skycoin/cx) [![Build status](https://ci.appveyor.com/api/projects/status/y04pofhhfmpw8vef/branch/master?svg=true)](https://ci.appveyor.com/project/skycoin/cx/branch/master)

CX is a general purpose, interpreted and compiled programming
language, with a very strict type system and a syntax
similar to Golang's. CX provides a new programming paradigm based on
the concept of affordances.

## Table of Contents
=================

   * [CX Programming Language](#cx-programming-language)
   * [Table of Contents](#table-of-contents)
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
      
# CX Programming Language

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

# Installation

CX requires a Golang version of `1.15` or higher. 

## Binary Releases

This repository provides binary releases of the language. Check this link and download the appropriate binary release for your platfrom: https://github.com/skycoin/cx/releases

## Compiling from Source

If a binary release is not currently available for your platfrom or if
you want to have a nightly build of CX, you'll have to compile from
source. If you're not familiarized with Go, Git, your OS's terminal or
your OS's package manager (to name a few), we *strongly* recommend you
to try out a binary release. If you find any bugs or problems with the
binary release, submit an issue here:
https://github.com/skycoin/cx/issues, and we'll fix it for the next release.

### Compiling on Linux

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

Requires installation of GIT from https://git-scm.com/downloads prior to compile.
An installation script is also provided for Windows named `cx-setup.bat`. You can compile CX on Windows by running:

```
cx-setup.bat
```

You should test your installation by running:

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