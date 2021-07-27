# CX Compiler specification

## Glossary
* [Lexer](https://en.wikipedia.org/wiki/Lexical_analysis) - takes a text [or sequence of characters] as an input and breaks it up into a list of tokens.
* [AST](https://en.wikipedia.org/wiki/Abstract_syntax_tree) - An abstract syntax tree (AST) is a way of representing the syntax of a programming language as a hierarchical tree-like structure. The AST contains a final parsed output that can be executed.
* [Type Signature](https://en.wikipedia.org/wiki/Type_signature) - A function signature (or type signature, or method signature) defines input and output of functions or methods. A signature can include: parameters and their types. a return value and type.
* [Actions](https://github.com/skycoin/cx/tree/develop/cxparser/actions) - Actions package contains functions that are used to construct the AST. These functions are called stages of the parser.
---

## Compiler Stages
* Stage 1: Parse packages and package imports.
* Stage 2: Parse structs, type aliases and enums.
* Stage 3: Parse global variables.
* Stage 4: Parse function type signatures.
* Stage 5: Set global variable construct functions.
* Stage 6: Parse function bodies.

---

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

### Stage 1: Parse Package Imports and Load Source Files

At this stage:
- The text file containing the program is is loaded into `[]byte` as `SourceCode`
- import paths are identified
- All package imports in the text file are loaded into `[]byte`
- package names are parsed from the input program and all import paths

**`source file definition`**
```
type SourceFile struct {
	FileName 		string
	SourceCode 		[]byte
}
```

**`package definition`**
```
type Package struct {
	Name 			string
	LineNo 			int
	FileName 		string
	SourceFiles		[]SourceFile
}
```

Input:
- filename: `main.cx`
- package: `package main`

Output:
`Package -> {Name: "main", LineNo: 2, FileName: "main.cx", SourceFiles: [{FileName: "main.cx", SourceCode: ...}]}`

**Actions called**
 - [`actions.DeclarePackage`](https://github.com/skycoin/cx/blob/develop/cxparser/actions/declarations.go#L212)

### Stage 2: Parse Struct and Enum types

Scan - identifier, line number and filename of: `Struct`, `Type aliases`, `Enums` into a struct.

Input: 
```
8  type Entry struct {
9      age i32
10 }
```

**`struct definition`**
``` golang
type Struct struct {
	PackageName		string
	Name		 	string
	LineNo 		 	int
	FileName 	 	string
	Fields			CXArgument
}
```

Example: in the sample program, we have a struct definition in line 8. So the parsed struct should be represented as:

Output:
`Struct -> { PackageName: "main", Name: "Entry", LineNo: 8, FileName: "main.cx", Fields: {Name: age, Type: i32} }`

**Actions called**
 - [`actions.DeclareStruct`](https://github.com/skycoin/cx/blob/develop/cxparser/actions/declarations.go#L184)

### Stage 3: Declare Global Variables

Global variable declarations are parsed, but not initialized, and thus the offset is -1 and will be set to the correct memory pointer in runtime. 

**`variable definition`**
``` golang
type CXArgument struct {
	PackageName		string
	Name		 	string
	LineNo 		 	int
	FileName 	 	string
	Type 			int //defined in constants
	Offset			int
}
```

Example: there is a variable declaration at line 12. The parsed output should be: 

Input:
``` 
12 var a i32 = 20
``` 
Output: 
`GlobalVariable -> { PackageName: "main", Name: "a", Type: i32, LineNo: 12, FileName: "main.cx", Offset: -1 }`.

**Actions called**
 - [`actions.DeclareGlobal`](https://github.com/skycoin/cx/blob/develop/cxparser/actions/declarations.go#L24)

### Stage 4: Parse Function Type Signatures

Parse function type signatures. A function's signature includes: function identifier, input parameters and return types. The only function in the above sample program is `main`. So, the compiler should parse:

**`function definition`**
``` golang
type Function struct {
	PackageName		string
	Name		 	string
	LineNo 		 	int
	FileName 	 	string
	InputParams		[]CXArgument
	ReturnParams	[]CXArgument
	Expressions		[]CXExpression
	ParentStruct 	*Struct
}
```

Input:
```
14 func main()() {
15     var entry Entry
16     entry.age = a
17     printf("%d\n", entry.age)
18 }
```

Output:
`Function -> {PackageName: "main", Name: "main", InputParams: {}, ReturnParams: {i32}, ParentStruct: nil, LineNo: 14, FileName: "main.cx"}`

**Actions called**
 - [`actions.FunctionDeclaration`](https://github.com/skycoin/cx/blob/develop/cxparser/actions/functions.go#L117)
 - [`actions.FunctionHeader`](https://github.com/skycoin/cx/blob/develop/cxparser/actions/functions.go#L21)

### Stage 5: Initialize Global Variables

At this stage, we get the `global variables` and for each global variable, we add a `construct function` to initialize it at runtime.
We have to initialize the global variables after we parse all functions, because global variable declaration can have function calls.
### Stage 6: Parse Function Bodies

Parse function bodies.

Input: 

```golang
14 //func main()() {
15     var entry Entry
16     entry.age = a
17     printf("%d\n", entry.age)
18 //}
```

Each expression within the function scope are parsed into a `Expressions` of function

**Actions called**
 - `actions.DeclarationSpecifier`
 - `actions.StructLiteralAssignment`
 - `actions.Assignment`
 - `actions.SliceLiteralExpression`
 - ...