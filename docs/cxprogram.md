# Glossary
* [Lexer](https://en.wikipedia.org/wiki/Lexical_analysis) - takes a text [or sequence of characters] as an input and breaks it up into a list of tokens.
* [AST](https://en.wikipedia.org/wiki/Abstract_syntax_tree) - An abstract syntax tree (AST) is a way of representing the syntax of a programming language as a hierarchical tree-like structure. This structure is used for generating symbol tables for compilers and later code generation.
* [Type Signature](https://en.wikipedia.org/wiki/Type_signature) - A function signature (or type signature, or method signature) defines input and output of functions or methods. A signature can include: parameters and their types. a return value and type.

# CX Compiler specification

The cx compiler is broken down into several stages and each stages output is an input to the next stage.

A sample CX program looks like this:

``` golang
1  //main.cx
2  package main
3  
4  import package_name
5  
6  type alias i32
7  
8  type Entry struct {
9      age i32
10 }
11 
12 var a i32 = 20
13
14 func main()() {
15     var entry Entry
16     entry.age = a
17     printf("%d\n", entry.age)
18 }
```

This is a simple program that prints age of Entry.
We are going to reference this example throught the specification document.

## Stage 1

At this stage:
- The text file containing the program is is loaded into `[]byte` as `SourceCode`
- import paths are identified
- All package imports in the text file are loaded into `[]byte`
- package names are parsed from the input program and all import paths

## Stage 2

Scan - identifier, line number and filename of: `Struct`, `Type aliases`, `Enums` into a struct.

Example: in the sample program, we have a struct definition in line 8. So the parsed struct should be represented as:

`Struct -> { PackageName: "main", Name: "Entry", LineNo: 8, FileName: "main.cx" }`

## Stage 3

Global variable declarations are parsed, but not initialized. 

Example: there is a variable declaration at line 12. The parsed output should be:  
`GlobalVariable -> { PackageName: "main", Name: "a", Type: i32, LineNo: 12, FileName: "main.cx" }`.

## Stage 4

Parse function type signatures. A function's signature includes: function identifier, input parameters and return types. The only function in the above sample program is `main`. So, the compiler should parse:

`Function -> {PackageName: "main", Name: "main", Input: {}, Return: {i32}, ParentStruct: nil, LineNumber: 14, FileName: "main.cx"}`

## Stage 5

We have to initialize the global variables after we parse all functions, because global variable declaration can have function calls. At this stage, we get the `global variables` and for each global variable, we intialize it with input from the source code at the given `LineNo`.

## Stage 6

Parse function bodies.