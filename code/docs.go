package code



bytecode is a domain specific language for domain specific machine.



steps :

 1 : take cx programming expression 10 + 20

 2 : tokenize it and parse it using lexer and pratt parser packages.

 3 : take the resulting AST, Whose nodes are defined in our ast packages.

 4 : pass ast to the newly-built compiler, which compiles it to bycode.

 5 : take the bytecode and hand it over to the also newly-built virtual machine will 
    execute it.

 6 : make sure that virtual machine turned it into 3.
 
 

	Lexer -> Parser -> compiler -> virtual machine



	In terms of data structure

	String -> Tokens -> AST -> bytecode - > Objects



