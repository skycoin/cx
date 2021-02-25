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
      * [Strict Type System](#strict-type-system)
   * [CX Projects, Libraries and other Resources](#cx-projects-libraries-and-other-resources)
      * [CX Chains:](#cx-chains)
      * [CX Examples:](#cx-examples)
      * [CX Libraries:](#cx-libraries)
      * [CX Video Games:](#cx-video-games)
      * [Miscellaneous:](#miscellaneous) <!--* [CX Roadmap](#cx-roadmap) -->
   * [CX Chains (CX   Skycoin Blockchain)](#cx-chains-cx--skycoin-blockchain)
   * [Changelog](#changelog)
   * [Compiler Development](CompilerDevelopment.md)
   * [Installation](#installation)
      * [Binary Releases](#binary-releases)   <!-- * [MacOS Homebrew Install](#macos-homebrew-install) -->
      * [Compiling from Source](#compiling-from-source)
         * [Installing Go](#installing-go)
         * [Compiling CX on *nix](#compiling-cx-on-nix)
         * [Compiling CX on Windows](#compiling-cx-on-windows)
      * [Updating CX](#updating-cx)
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

Feel free to [create an issue](https://github.com/skycoin/cx/issues)
requesting a better explanation of a feature.

# CX Projects, Libraries and other Resources
[[Back to the Table of Contents] ↑](#table-of-contents)

## CX Chains:

* https://github.com/skycoin/cx-chains

## CX Examples:

* https://github.com/skycoin/cx/tree/develop/examples

## CX Libraries:

* https://github.com/skycoin/cxfx [video game library]
* https://github.com/ReewassSquared/cxm [math library]
* https://github.com/ReewassSquared/CXSL [general utilities library]
* https://github.com/ReewassSquared/CXCIPHER [crypto library]
* https://github.com/asahi3g/pumpcx [user interface library]
* https://github.com/ReewassSquared/CXML [machine learning library]

## CX Video Games:

* https://github.com/galah4d/casino-cx [slot machine]
* https://github.com/atang152/crappyBall-cx [flappy bird clone]
* https://github.com/galah4d/pacman-3d [pacman 3D clone]
* https://github.com/skycoin/cx-games/tree/master/Snake-by-Lunier [snake clone]
* https://github.com/skycoin/cx-games/tree/master/SynthCat-Brick-Breaker-by-RedCurse [brick breaker clone]
* https://github.com/skycoin/cx-games/tree/master/Pac-Man-CX-by-Galah4d [pacman 2D clone]
* https://github.com/skycoin/cx-games/tree/master/Whacky-Stack [tetris clone]
* https://github.com/skycoin/cx-games/tree/master/ridge-blaster [dig-n-rig clone]
* https://github.com/taekwondouglas/space-invaders [space invaders clone dapp using CX chains]

## Miscellaneous:

* https://github.com/skycoin/cx-website [cx.skycoin.com]

<!--# CX Roadmap

![CX Roadmap](https://raw.githubusercontent.com/skycoin/cx/master/readme-images/cx-roadmap.jpg)-->

# CX Chains (CX + Skycoin Blockchain)
[[Back to the Table of Contents] ↑](#table-of-contents)

CX Chains are Skycoin's solution for the creation of blockchain-based
programs. You can read more about them in the [CX
wiki](https://github.com/skycoin/cx/wiki/CX-Chains-Tutorial) for the latest release or in [`documentation/BLOCKCHAIN.md`](documentation/BLOCKCHAIN.md) for the `develop` branch of CX (the bleeding edge version of CX).

# Changelog
[[Back to the Table of Contents] ↑](#table-of-contents)

Check out the latest additions and bug fixes in the [changelog](CHANGELOG.md).

# Installation
[[Back to the Table of Contents] ↑](#table-of-contents)

## Binary Releases

This repository provides binary releases of the language. Check this link and download the appropriate binary release for your platfrom: https://github.com/skycoin/cx/releases

<!-- ### MacOS Homebrew Install

The simplest way to install CX on MacOS is to use the Homebrew package manager to install the prebuilt binary release. If you do not already have Homebrew installed, please visit the [Homebrew website](https://brew.sh/) for installation instructions.

Once Homebrew is installed, use the following commands to setup the Tap and then install CX.

```sh
brew tap skycoin/homebrew-skycoin
brew install skycoin-cx
```

To update use the following command:

```sh
brew update skycoin-cx 
```
-->
## Compiling from Source

If a binary release is not currently available for your platfrom or if
you want to have a nightly build of CX, you'll have to compile from
source. If you're not familiarized with Go, Git, your OS's terminal or
your OS's package manager (to name a few), we *strongly* recommend you
to try out a binary release. If you find any bugs or problems with the
binary release, submit an issue here:
https://github.com/skycoin/cx/issues, and we'll fix it for the next release.

### Installing Go

CX supports go1.15+.

[Go 1.15+ installation/setup](https://github.com/skycoin/cx/blob/develop/GO-INSTALLATION.md)

### Compiling CX on \*nix

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
Setup GOPATH for cx:

```
set GOPATH=%userprofile%/cx
```

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

# Running CX
## Hello World

Do you want to know how CX looks? This is how you print "Hello, World!"
in a terminal:

Create a `hello-world.cx` file, put these contents in the file:

```
package main

func main () {
    str.print("Hello, World!")
}
```

Then run `cx hello-world.cx`. You can also run the examples found in the `examples` directory.

## Basic Options

If you type `cx --help` or `cx -h`, you should see a text describing
CX's usage, options and more. `cx --version` provides the detail about the current cx version installed on the machine.

Some interesting options are:

* `--repl` which loads the program and makes CX run in REPL mode
(useful for debugging a program)
* `--web` which starts CX as a RESTful web service (you can send code
  to be evaluated to this endpoint: http://127.0.0.1:5336/eval)


## REPL tutorial

Once CX has been successfully installed, running `cx` should print
this in your terminal:

```
CX 0.5.13
More information about CX is available at http://cx.skycoin.com/ and https://github.com/skycoin/cx/
:func main {...
	*
```

This is the CX REPL
([read-eval-print loop](https://en.wikipedia.org/wiki/Read%E2%80%93eval%E2%80%93print_loop)),
where you can debug and modify CX programs. The CX REPL starts with a
barebones CX structure (a `main` package and a `main` function) that
you can use to start building a program.

Let's create a small program to test the REPL. First, write
`str.print("Testing the REPL")` after the `*`, and press enter. After
pressing enter you'll see the message "Testing the REPL" on the screen. If you
then write `:dp` (short for `:dProgram` or *debug program*), you
should get the current program AST printed:

```
Program
0.- Package: main
	Functions
		0.- Function: main () ()
			0.- Expression: str.print("" str)
		1.- Function: *init () ()
```

As we can see, we have a `main` package, a `main` function, and we
have a single expression: `str.print("Testing the REPL")`.

Let's now create a new function. In order to do this, we first need to
leave the `main` function. At this moment, any expression (or function
call) that we add to our program is going to be added to `main`. To
exit a function declaration, press `Ctrl+D`. The prompt (`*`) should
have changed indentation, and the REPL now shouldn't print `:func main
{...` above the prompt:

```
:func main {...
	*
*
```

Now, let's enter a function prototype (an empty function which only
specifies the name, the inputs and the outputs):

```
* func sum (num1 i32, num2 i32) (num3 i32) {}
*
```

You can check that the function was indeed added by issuing a `:dp`
command. If we want to add expressions to `sum`, we have to select it:

```
* :func sum

:func sum {...
	*
```

Notice that there's a semicolon before `func sum`. Now we can add an expression to it:

```
:func sum {...
	* num3 = num1 + num2
```

Now, exit `sum` and select `main` with the command `:func main`. Let's
add a call to `sum` and print the value that it returns when giving
the arguments 10 and 20:

```
:func main {...
	* i32.print(sum(10, 20))
30
```

# Learning CX
[[Back to the Table of Contents] ↑](#table-of-contents)

To get your hand dirty with the CX language, you can read the following reference:

## 1. CX Book

You can find the book source code and its releases in its [CX book Github
repository](https://github.com/Skycoin/cx-book).

## 2. Language Guide

You can find the language guide which is more up to date than the book in [language guide](LanguageGuide.md) section.
