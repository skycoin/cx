// +build regexp

package regexp

import (
	"github.com/skycoin/cx/cx/ast"
    "github.com/skycoin/cx/cx/types"
    . "github.com/skycoin/cx/cx/opcodes"
)

 func RegisterPackage() {
	regexpPkg := ast.MakePackage("regexp")
	regexpStrct := ast.MakeStruct("Regexp")

	regexpStrct.AddField(ast.MakeArgument("exp", "", 0).AddType(types.STR).AddPackage(regexpPkg))

	regexpPkg.AddStruct(regexpStrct)

	ast.PROGRAM.AddPackage(regexpPkg)


	RegisterFunction("regexp.Compile", opRegexpCompile, In(ast.ConstCxArg_STR), Out(ast.Struct("regexp", "Regexp", "r"), ast.ConstCxArg_STR))
	RegisterFunction("regexp.MustCompile", opRegexpMustCompile, In(ast.ConstCxArg_STR), Out(ast.Struct("regexp", "Regexp", "r")))
	RegisterFunction("regexp.Regexp.Find", opRegexpFind, In(ast.Struct("regexp", "Regexp", "r"), ast.ConstCxArg_STR), Out(ast.ConstCxArg_STR))
}
