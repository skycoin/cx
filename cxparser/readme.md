# CX Lexer
CX in house `lexer` generates a chain of tokens which is used by yacc for parsing.
(found in)[lex.go] (https://github.com/skycoin/cx/cxparser/cxpartialparsing/lex.go )


Parser stages 


stage 0

prilimalary stage of using regex 

Preliminarystage()
The compiler creates an AST first by using regular expression parsing (found in [utils.go](https://github.com/PratikDhanave/cx/blob/develop/cxparser/cxparsing/utils.go#L21) which structures the AST to include all package declarations, import chains, struct declarations, and globals, skipping over comments. 
This preliminary stage of parsing aids further stages since the structure of a CX repository and the names of custom types are already known. 

Step I

 `Passone`

After this preliminary stage, the first parsing stage compiles function primitives and types, along with the structure of their parameters and global variables. This stage also finalizes compilation of structs and ensures packages and their imports are correct. This uses `goyacc` for parsing and creates a chain of tokens for the parser using an in-house lexer known as `Lexer`. 

main tasks 
-the signature of functions and methods are added here
-also we figure out type of global varialble.
-imports
-type of stuct are figure out



step II

`Passtwo`

The second parsing stage fully compiles functions and all expressions. This functions similarly to the first stage, using `Lexer` and `goyacc`, but also uses `cxparser/actions` for each action the parser should take after encountering a valid syntactical production rule. These work entirely to validate program structure and to build the AST for a CX program.


