# CX Compiler Frontend

## Refactor Compiler Frontend
### Stages:
- [x] Stage 1: File and Package Loading
	- Load CX packages and get imports. Extracting the packages that need to be imported
	- Loading all modules/files in each package

- [x] Stage 2: Declaration Extraction
	- Function to extract all global declarations from source file
	- Function to extract all alias and enum declarations from source file
	- Function to extract all struct declarations from source files
	- Function to extract all function declarations from source file
	- Function for redeclaration checks

- [ ] Stage 3: TypeChecks For Signatures/Types/Assigns
	- Parse Packages
	- Parse Imports
	- Parse the globals
	- Parse the aliases and enums
	- Parse the structs
	- Parse function header

- [ ] Stage 4: Function Body Extraction
	- Function to extract all function bodies from source file
	- And verify () {} pairings, etc for error

Goal: Separate out compiler stage into modules. They can be in separate packages, then import and run as a library.

## Milestones:
- [ ] Stop using Goyacc to parse struct definitions; just use golang + new type format

- [ ] Only feed the function body into GOYACC for parsing
then GOYACC calls AST functions for each line
Optional: Why cant we have 8 goroutines, each parsing one function body? Are there globals in Goyacc?
Split Tokenizer, so that its only outputing data for a single function body
Tokenize all function bodies in parallel with channel+goroutine

- [ ] Parse the local variable declarations in the function bodies, without using GOYACC (using golang). Ignore assignments
Put a flag option to disable AST output for local variable declarations in GOYACC

- [ ] Only feed one line of code, at a time to GOYACC and get the AST output
Goyacc only parases individual lines of code
Split Tokenizer, so that its only outputing data for a single line of golang code
Tokenize all lines in parrallel with goroutines+channel

- [ ] Parse whole individual line of code, with golang function without GOYACC
Have option to toggle which is used, etc.
	- Extract all Global Declarations + Global Assignments
	- Extract the declaration lines
	- Extract all Function Declaration + Function Bodies
	- Check () {} to make sure they are balanced and no error
	- Extract all Struct Definitions 


