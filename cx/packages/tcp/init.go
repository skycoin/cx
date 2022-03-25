package tcp

import (
	"github.com/skycoin/cx/cx/ast"
	"github.com/skycoin/cx/cx/opcodes"
)

func RegisterPackage(prgrm *ast.CXProgram) {

	netPkg := ast.MakePackage("tcp")
	pkgIdx := prgrm.AddPackage(netPkg)
	netPkg, _ = prgrm.GetPackageFromArray(pkgIdx)

	dialerStrct := ast.MakeStruct("Dialer")
	netPkg.AddStruct(dialerStrct)

	opcodes.RegisterFunction(prgrm, "tcp.Dial", opTCPDial, opcodes.In(ast.ConstCxArg_STR, ast.ConstCxArg_STR), opcodes.Out(ast.ConstCxArg_STR))
	opcodes.RegisterFunction(prgrm, "tcp.Listen", opTCPListen, opcodes.In(ast.ConstCxArg_STR, ast.ConstCxArg_STR), opcodes.Out(ast.ConstCxArg_STR))
	opcodes.RegisterFunction(prgrm, "tcp.Accept", opTCPAccept, opcodes.In(ast.ConstCxArg_STR, ast.ConstCxArg_STR), opcodes.Out(ast.ConstCxArg_STR))
	opcodes.RegisterFunction(prgrm, "tcp.Close", opTCPClose, nil, nil)
}
