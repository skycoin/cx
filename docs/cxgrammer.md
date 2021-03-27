this is first draft

cx grammer is define in [https://github.com/skycoin/cx/blob/develop/cxgo/cxgo0/cxgo0.y]

cx grammer is used to generate cx parser which is used to generate program ast.

the cx grammer is passed to yacc to generate parser for current cxgrammer.
./bin/goyacc -o cxgo/cxgo0/cxgo0.go cxgo/cxgo0/cxgo0.y


cx grammer is defined 

declarations of variable

rules of action 

functions



cx grammer consists of two things 

BNF grammar  - > action

Whenever a particular BNF grammar rule is recognized, the action describes what to do is perform on current token. 

For example

function_declaration:
                function_header function_parameters compound_statement
                {
			actions.FunctionDeclaration($1, $2, nil, $3)



if we found grammer function_declaration at the time of parsing then its corresponding action is perform 
actions.FunctionDeclaration()

all  cx action are define in cxgo/actions/


at the end of ParseSourceCode() we get complete CXProgram as object.


 