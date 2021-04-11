package globals

import "github.com/skycoin/cx/cx/ast"

//What does this do?
//This is where intializers get pushed, but only used 4 times
//TODO: Get rid of this
//TODO: Move list of inits to do, to AST, not here
var SysInitExprs []*ast.CXExpression