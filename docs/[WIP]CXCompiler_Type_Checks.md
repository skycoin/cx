# CXCompiler: TypeChecks For Signatures/Types/Assigns
At this stage all of the types are known. The functions in all packages that can be called, all structs are defined or known, all aliases/enums are defined and known.

## Goals:
- [x] Function to parse the globals
- [ ] Function to parse the aliases and enums
- [x] Function to parse the structs
- [x] Function to parse function header
	- Check that it can be converted to New Style Type Signature Format
	- Check function signature return type (check if valid)
	- Check if the return return exist, are valid, syntactic errors
- [ ] Verify
	- if a file uses PackageA.StructTypeB that PackageA is in the imports of that file, etc
	- Then check that the package is actually imported by the file(first check that it exists)

Note: 
Call AST API directly to build the AST. Then, delete or comment out goyacc code for parsing these 

## Progress
- Parsing data types only support the pattern/order of *[5]i32 other patterns/orders like [4]\*str aren't supported
- Performance of parsing data types is not optimal yet can be improved with simpler regex or better algorithm for identifying and removing parts from back to front 
- Parsing enums/aliases and type definitions aren't implemented yet.
- Type signature checks are done by the Actions API and panics if there's any error
- Imports verification is done by the Actions API and panics if there's any error