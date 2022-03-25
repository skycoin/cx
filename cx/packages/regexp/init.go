// +build regexp

package regexp

import (
	"github.com/skycoin/cx/cx/ast"
	"github.com/skycoin/cx/cx/opcodes"
	"github.com/skycoin/cx/cx/types"
)

func RegisterPackage(prgrm *ast.CXProgram) {
	regexpPkg := ast.MakePackage("regexp")
	pkgIdx := prgrm.AddPackage(regexpPkg)
	regexpPkg, _ = prgrm.GetPackageFromArray(pkgIdx)

	regexpStrct := ast.MakeStruct("Regexp")
	regexpStrct.AddField(prgrm, ast.MakeArgument("exp", "", 0).AddType(types.STR).AddPackage(regexpPkg))
	regexpPkg.AddStruct(regexpStrct)

	opcodes.RegisterFunction(prgrm, "regexp.Compile", opRegexpCompile, opcodes.In(ast.ConstCxArg_STR), opcodes.Out(ast.Struct(prgrm, "regexp", "Regexp", "r"), ast.ConstCxArg_STR))
	opcodes.RegisterFunction(prgrm, "regexp.MustCompile", opRegexpMustCompile, opcodes.In(ast.ConstCxArg_STR), opcodes.Out(ast.Struct(prgrm, "regexp", "Regexp", "r")))
	opcodes.RegisterFunction(prgrm, "regexp.Regexp.Find", opRegexpFind, opcodes.In(ast.Struct(prgrm, "regexp", "Regexp", "r"), ast.ConstCxArg_STR), opcodes.Out(ast.ConstCxArg_STR))
}
