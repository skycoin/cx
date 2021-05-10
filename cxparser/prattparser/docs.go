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


	for every grammer/ token we have function assocated with to constuct AST.

	How to write Parsing


*/

package parser
