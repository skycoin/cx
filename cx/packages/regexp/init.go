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
	regexpArg := ast.MakeArgument("exp", "", 0).SetType(types.STR).SetPackage(regexpPkg)
	regexpStrct.AddField(prgrm, regexpArg.Type, regexpArg)
	regexpPkg.AddStruct(prgrm, regexpStrct)

	opcodes.RegisterFunction(prgrm, "regexp.Compile", opRegexpCompile, opcodes.In(ast.ConstCxArg_STR), opcodes.Out(ast.MakeStructParameter(prgrm, "regexp", "Regexp", "r"), ast.ConstCxArg_STR))
	opcodes.RegisterFunction(prgrm, "regexp.MustCompile", opRegexpMustCompile, opcodes.In(ast.ConstCxArg_STR), opcodes.Out(ast.MakeStructParameter(prgrm, "regexp", "Regexp", "r")))
	opcodes.RegisterFunction(prgrm, "regexp.Regexp.Find", opRegexpFind, opcodes.In(ast.MakeStructParameter(prgrm, "regexp", "Regexp", "r"), ast.ConstCxArg_STR), opcodes.Out(ast.ConstCxArg_STR))
}
