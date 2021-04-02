
Parser stages 

prilimalary stage of using regex 

before lexerstep0 cxgo0.Parse() function call
This preliminary stage of parsing we figure out package ,stucts, and globals.


Step I

 lexerstep0 cxgo0.Parse() function call

the first parsing stage compiles function primitives and types, along with the structure of their parameters and global variables. This stage also finalizes compilation of structs and ensures packages and their imports are correct. This uses goyacc for parsing and creates a chain of tokens for the parser using an in-house lexer known as Lexer.

main tasks 
-the signature of functions and methods are added here
-also we figure out type of global varialble.
-imports
-type of stuct are figure out



step II

cxgo.Parse(cxgo.NewLexer(b) call

The second parsing stage fully compiles functions and all expressions. This functions similarly to the first stage, using Lexer and goyacc, but also uses cxgo/actions for each action the parser should take after encountering a valid syntactical production rule. These work entirely to validate program structure and to build the AST for a CX program.


