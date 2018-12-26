# Overview Over the CX Codebase

This file contains an overview of the CX codebase.  It is recommended reading
if you want to get involved in developing CX itself.

## Directory structure

The top directory contains informational files, installation scripts and
subdirectories.  We will try to describe then grouped in some form of logical
order rather than alphabetical order.

The informational files are:

* **README.md** is an overall introduction and description of the CX project.
* **CHANGELOG.md** lists changes between different versions of CX.
* **CONTRIBUTING.md** lists various ways that you can contribute to CX
* **OVERVIEW.md** This file

The tools and installation scripts are:

* **cx-setup.bat** The setup and installation script for CX on Windows
* **cx.sh** The setup and installation script for CX on Linux
* **Makefile** The standard way of building CX, using the `make` tool.

### CX Source Code

* **cx/** This directory contains two main parts:
  * the main datastructures representing a CX program (`cx*.go`)
  * the runtime part of CX (`op_*.go`)
* **cxgo/** This directory contains the parser part of CX. The parser
analyzes the source code and builds a tree of the detastructures mentioned
above that represents the program.  The root of this tree is represented by an
instance of the `CXProgram` struct.
* **vendor/** This directory contains 3^rd party code that needed to be
included in the source in order for it to be built correctly. This is not of
the CX sources themselves, but rather dependencies of it.

### Tests and Similar

* **tests/** This directory contains product tests and regression tests for
the CX parser and interpreter.  See the file `CONTRIBUTING.md` for more
details about the tests.
* **benchmark/** This directory contains various benchmarks vs other languages
such as Go and Python.

### Documentation 

* **documentation/** contains documentation about the internals
                     of the CX codebase. This file (OVERVIEW.md) is in that
                     directory. 
* **book/** This directory contains a full CX tutorial book written in the
            LaTeX typesetting system.  Please note that the contents of the
            book is created on a site called Overleaf, which provides a good
            authoring environment for LaTeX files.  Only occasionally are
            snapshots taken from Overleaf and stored into git as a backup.
            If you want to help out with the book, please ask the team for
            access to Overleaf. 
* **README-images/** stores images used in the `README.md` file.

### Tools

* **development/** contains tools for developers developing CX. See the README
                   in that directory for more details.
* **gui-helper/**  FIXME: What is this?
* **ide/**         contains a CX IDE (Integrated Development
                   Environment). This will become a graphic development tool
                   for CX.
* **object-explorer/** contains a tool for querying live objects on the heap.

### Example Code

* **examples/** contains many small and smallish exemples of how to write CX
programs. 
* **cx-games/** contains a number of games written in CX. These games can be
used for testing CX or as templates for other CX programs.

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

### CXProgram

### CXPackage

### CXStruct

### CXFunction

### CXArgument

### CXExpression

## Opcodes

## Serialization and Deserialization

## Affordances

