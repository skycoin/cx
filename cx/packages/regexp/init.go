// +build regexp

package regexp

import (
	"github.com/skycoin/cx/cx/ast"
	. "github.com/skycoin/cx/cx/opcodes"
	"github.com/skycoin/cx/cx/types"
)

func RegisterPackage(prgrm *ast.CXProgram) {
	regexpPkg := ast.MakePackage("regexp")
	regexpStrct := ast.MakeStruct("Regexp")

	regexpPkg.AddStruct(regexpStrct)

	prgrm.AddPackage(regexpPkg)
	regexpStrct.AddField(prgrm, ast.MakeArgument("exp", "", 0).AddType(types.STR).AddPackage(regexpPkg))

	RegisterFunction(prgrm, "regexp.Compile", opRegexpCompile, In(ast.ConstCxArg_STR), Out(ast.Struct(prgrm, "regexp", "Regexp", "r"), ast.ConstCxArg_STR))
	RegisterFunction(prgrm, "regexp.MustCompile", opRegexpMustCompile, In(ast.ConstCxArg_STR), Out(ast.Struct(prgrm, "regexp", "Regexp", "r")))
	RegisterFunction(prgrm, "regexp.Regexp.Find", opRegexpFind, In(ast.Struct(prgrm, "regexp", "Regexp", "r"), ast.ConstCxArg_STR), Out(ast.ConstCxArg_STR))
}
