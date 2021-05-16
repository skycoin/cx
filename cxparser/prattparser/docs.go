/*

	Pratt Parsing : Notes

	Parsing is the process by which a compiler turns a sequence of tokens into a tree representation:


                            	Add
                 	Parser     / \
	 "1 + 2 * 3"    ------->   1  Mul
                              / \
                             2   3


	whats Important

	Nest the nodes of AST correctly.

	How It works


	1 + 2 + 3


	recursive parsing

	func parseExpression() {
    ...
    loop {
        ...
        parseExpression()
        ...
    }

	}

	pratt parsing

	func parseExpression() {


		prefixfn = GetprefixParseFns()

		leftExp := prefixfn()

		loop {

			infixfn := GetinfixParseFns(p)

			leftExp = infixfn(leftExp){
					parseExpression
			}
    }
	}



	for every token we have prefixfn and infixfn to evaluate operator precedence

	base on this precedence we construct ast recursively.


	In the pratt parsing we construct ast as well as evaluate operator precedence.


	for every grammer/ token we have functions assocated with to constuct AST node.

	for example : parseIntegerLiteral ()

	How to write Parsing



References :

i) https://engineering.desmos.com/articles/pratt-parser/
ii)https://matklad.github.io/2020/04/13/simple-but-powerful-pratt-parsing.html
iii) http://craftinginterpreters.com/parsing-expressions.html
iv) https://www.oilshell.org/blog/2017/03/31.html
v) http://journal.stuffwithstuff.com/2011/03/19/pratt-parsers-expression-parsing-made-easy/
vi) https://en.wikipedia.org/wiki/Operator-precedence_parser
vii) https://www.oilshell.org/blog/2016/11/01.html
viii) https://www.npmjs.com/package/pratt-parser
ix) https://eli.thegreenplace.net/2010/01/02/top-down-operator-precedence-parsing
x) https://dev.to/jrop/pratt-parsing
xi) https://github.com/jrop/pratt-calculator
xii) https://crockford.com/javascript/tdop/tdop.html
xiii) https://elementpath.readthedocs.io/en/latest/pratt_api.html

*/

package parser
