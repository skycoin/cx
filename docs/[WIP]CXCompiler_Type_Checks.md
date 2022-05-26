# CXCompiler: TypeChecks For Signatures/Types/Assigns
Note: cant do this yes

At this stage all of the types are known
The functions in all packages that can be called
All structs are defined or known
All aliases/enums are defined and known

Parse The Function declarations (in Golang)
Check that it can be converted to New Style Type Signature Format
Check function signature return type (check if valid)
Check if the return return exist, are valid, syntactic errors
Parsing should be in golang, with a function
Check function signature input type (check if valid)
Check if the return type exists, are valid, syntactic errors

Parse The Structs declarations (in Golang)
…

Parse the Global Struct Definitions (in Golang)
…


TODO:
Once we have a parsed “New Style” Function/Struct/GlobalDeclarations Type Signature
THEN, call AST operations directly
THEN, delete or comment out goyacc code for parsing these 

Verify:
	- if a file uses PackageA.StructTypeB that PackageA is in the imports of that file, etc
Then check that the package is actually imported by the file
Etc first check that it exists