# Overview Over the CX Codebase

This file contains an overview of the CX codebase.  It is recommended reading
if you want to get involved in developing CX.

## Directory structure

### CX Source Code

* **cx/** This directory contains:
  * the main datastructures representing a CX program (`cx*.go`)
  * the runtime part of CX (`op_*.go`)
* **cxgo/** This directory contains the parser part of CX. The parser
analyzes the source code and builds a tree of the detastructures mentioned
above that represents the program.  The root of this tree is represented by an
instance of the `CXProgram` struct.
* **vendor/** This directory contains 3^rd party code that needed to be
included in the source in order for it to be built correctly. This is not of
the CX sources themselves, but rather dependencies of it.

### Tests and benchmarks

* **tests/** This directory contains product tests and regression tests for
the CX parser and interpreter.  See the file `CONTRIBUTING.md` for more
details about the tests.
* **benchmark/** This directory contains various benchmarks vs other languages
such as Go and Python.

### Tools

* **development/** contains tools for developers developing CX. See the README
                   in that directory for more details.

## Parser Workflow

The CX parser analyzes the source code of a CX program and
 works in two passes:

1. The first pass (in `cxgo/cxg0`) reads all the types, global variables and
function names and stores them.
2. The second pass (in `cxgo/`) reads the full source code of the program and
creates the object tree representing the program. This tree will be passed to
the interpreter run.

Each pass of the parser contains a lexer and a parser proper.  If you are not
familiar with these terms, we suggest that you read an introductory text about
computer compilers.

## Interpreter

The interpreter interprets the internal representation of the program. This
representation is created by a tree of the data structures in the following
subsections.


