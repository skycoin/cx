CX Parser

The compiler creates an AST first by using regular expression parsing 
(found in /cx/cxgo/cxgo/cxparser.go) which structures the AST to include all package declarations, import chains, struct declarations, and globals, skipping over comments. This preliminary stage of parsing aids further stages since the structure of a CX repository and the names of custom types are already known.

After this preliminary stage, the first parsing stage compiles function primitives and types, along with the structure of their parameters and global variables. This stage also finalizes compilation of structs and ensures packages and their imports are correct. This uses goyacc for parsing and creates a chain of tokens for the parser using an in-house lexer known as Lexer.

The second parsing stage fully compiles functions and all expressions. This functions similarly to the first stage, using Lexer and goyacc, but also uses cxgo/actions for each action the parser should take after encountering a valid syntactical production rule. These work entirely to validate program structure and to build the AST for a CX program.

Finally, the main function and invisible *init functions are created, the latter of which acts as an initializer for global variables.

The lifetime of a CX program in depth is given here:

lexerStep0 - here, a copy of the source code is passed in as an array of strings, along with the file names. A series of regular expressions are compiled to check for various things. Finally, the source code is iterated over, and for each file, we iterate over a scanner of the file. Then, for each successful scan:
We first check if we’re in a comment or not. If we are, we continue scanning and ignore the following steps.
Next, we check if this current line contains a package declaration. If so, then we check if the program has added a package with that name. If not, then we add that package to the CXProgram (PRGM0, basically the main program). Remember, now that package is the current package.
We do the same thing for structures.
We now start a second iteration over the source code, the same as before, and for each successful scan, after checking for comments like before:
We check for import statements, and for each one, add the CXPackage to the Imports of the active CXPackage.
We check for global variables, and for each one found, add it to Globals of the active CXPackage. Additionally, we now have an inBlock counter, and when positive, it means we are inside a block, and should not add any variable declarations to globals. Things like variable type, etc will be added in later steps.
We can now do the third step, which is the first call to the parser (cxgo0). This parses function declarations and builds global variables fully, including custom types (CXStructs) and packages.
As an intermediate step, after the program has been selected for and after the program has been checked for any compile or parsing errors, os arguments are added to the program.
The OS package builtin for CX is grabbed, and for all globals in the OS package not currently added to said package, they are declared, and DeclareGlobalInPackage is called on them, after getting rid of the first two arguments, presumed to be the call to CX and the name of the program to call.
Finally, the program is fully parsed by the cxgo/parser code.






function parseProgram take source code as input 

task I
create program 
goto
actions.PRGRM = cxcore.MakeProgram()

task II 
ParseSourceCode
cxgo.ParseSourceCode(sourceCode, fileNames)

goto
ParseSourceCode 

goto
lexerStep0
// lexerStep0 performs a first pass for the CX parser. Globals, packages and
// custom types are added to `cxgo0.PRGRM0`

cxgo0.PRGRM0 has all data into Packages []*CXPackage varaiable.

