# Glossary
* [Lexer](https://en.wikipedia.org/wiki/Lexical_analysis) - takes a text [or sequence of characters] as an input and breaks it up into a list of tokens.
* [AST](https://en.wikipedia.org/wiki/Abstract_syntax_tree) - An abstract syntax tree (AST) is a way of representing the syntax of a programming language as a hierarchical tree-like structure. This structure is used for generating symbol tables for compilers and later code generation.
* [Type Signature](https://en.wikipedia.org/wiki/Type_signature) - A function signature (or type signature, or method signature) defines input and output of functions or methods. A signature can include: parameters and their types. a return value and type.

# Lifetime of a CX program
CX command line program accepts a text input along with commands/flags.</br>
Example: `./bin/cx main.cx`, `./bin/cx --lexer main.cx`

```
//main.cx
package main

type Entry struct {
    age i32
}

func main()() {
    var entry Entry
    entry.age = 20
    printf("%d\n", entry.age)
}
```
A list of source file names are parsed, in this case `main.cx` is the only source file name. And a `*File` is grabbed for each source file names. These are an input to the CX parser, described below.


## CX Parser

At this stage, `actions.AST` is initialized and packages from `ast.PROGRAM`, which are initialized with CX program, are assigned to it.

The current CX parser has three stages: `PreliminaryStage`, `PassOne` and `PassTwo`.
<br />

## CX Lexer
CX in house `lexer` generates a chain of tokens which is used by yacc for parsing.
(found in) [lex.go](https://github.com/skycoin/cx/blob/develop/cxparser/cxparsingcompletor/lex.go )

**Example**

```
// main.cx
package main

type Entry struct {
    age i32
}

func main()() {
    var entry Entry
    entry.age = 20
    printf("%d\n", entry.age)
}
```
`cx --lexer main.cx` -> displays the lexer output.
![Sample Lexer Output](https://raw.githubusercontent.com/cbrom/cx/doc/compiler_stages/docs/images/sample-lexer-output.png)

## CX Compiler Stages

### `Preliminarystage`

This stage is an early stage regex parser that parses the following froma copy of the source code `[]*File`: 
* comments [single line and multi line] - `//` and `/* */` type comments are parsed from the source code.
* package name - if a new package name is found (example, `package utils`), it is added to `CXProgram`.
* struct types - a struct type is added to `CXPakckage`.
* global variables - variables that are in global scope are added to the package.
* package imports - package imports are checked while parsing for global variables and `ParseSourceCode` is called on the imported filepath.


At this stage the `actions.AST` has all package declarations, import chains, struct declarations, and global variables, skipping over comments. Then Passone is called with in the Preliminary stage.

### `Passone`

This stage performs a `goyacc` based partial parsing on the source codes. It is a continuation of the `Preliminary stage` that calls `actions` API. </br>
Actions called in this stage are:
* `actions.DeclareGlobal`
* `actions.DeclareStruct`
* `actions.DeclarePackage`
* `actions.DeclareImport`

In addition to this, the following components are parsed in this stage:
* function header declaration
* function type signature declaration 
* variable declarations
* expressions
