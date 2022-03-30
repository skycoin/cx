
## Glossary
* [AST](https://en.wikipedia.org/wiki/Abstract_syntax_tree) - An abstract syntax tree (AST) is a way of representing the syntax of a programming language as a hierarchical tree-like structure. The AST contains a final parsed output that can be executed.
* [Lexer](https://en.wikipedia.org/wiki/Lexical_analysis) - takes a text [or sequence of characters] as an input and breaks it up into a list of tokens.
* [CXParser Spec](https://github.com/skycoin/cx/blob/develop/docs/cxparser_spec.md) - CX Parser Specification defines the cx parser.
---

# CX Compiler Stages

Parsing starts in [/cmd/cx/helpers.go/parseProgram()](https://github.com/skycoin/cx/blob/develop/cmd/cx/helpers.go#L27). This creates an empty AST and will call [/cxparser/cxparsing/cxparsing.go/ParseSourceCode()](https://github.com/skycoin/cx/blob/develop/cxparser/cxparsing/cxparsing.go#L30) to start parsing the source code. It will copy the CX source code into a string variable and will be passed to [/cxparser/cxparsing/utils.go/PreliminaryStage()](https://github.com/skycoin/cx/blob/develop/cxparser/cxparsing/utils.go#L21). 

`PreliminaryStage`

In [/cxparser/cxparsing/utils.go/PreliminaryStage()](https://github.com/skycoin/cx/blob/develop/cxparser/cxparsing/utils.go#L21), it will identify all the packages, structs, global variables, and package imports and add these to our AST. The comments are disregarded. After this, PassOne() will be called.

`PassOne`

PassOne or the first parsing stage, using [/cxparser/cxpartialparsing/cxpartialparsing.go/Parse()](https://github.com/skycoin/cx/blob/develop/cxparser/cxpartialparsing/cxpartialparsing.go#L33), compiles function primitives and types, along with the structure of their parameters and global variables. This stage also finalizes compilation of structs and ensures packages and their imports are correct. This uses `goyacc` for parsing and creates a chain of tokens for the parser using an in-house lexer known as `Lexer`. 

`PassTwo`

The second parsing stage, PassTwo or the second parsing stage, using [/cxparser/cxparsingcompletor/cxparsingcompletor.go/Parse()](https://github.com/skycoin/cx/blob/develop/cxparser/cxparsingcompletor/cxparsingcompletor.go#L32), fully compiles functions and all expressions. 

Finally, the `main` function and invisible `*init` functions are created, the latter of which acts as an initializer for global variables. After this, the CXProgram AST is complete and is ready to run.
