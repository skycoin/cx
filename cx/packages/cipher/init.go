// +build cipher

package cipher

import (
	"github.com/skycoin/cx/cx/ast"
	"github.com/skycoin/cx/cx/constants"
	"github.com/skycoin/cx/cx/opcodes"
	"github.com/skycoin/cx/cx/types"
)

func RegisterPackage(prgrm *ast.CXProgram) {
	cipherPkg := ast.MakePackage("cipher")
	pubkeyStrct := ast.MakeStruct("PubKey")
	seckeyStrct := ast.MakeStruct("SecKey")

	pkgIdx := prgrm.AddPackage(cipherPkg)
	cPkg, _ := prgrm.GetPackageFromArray(pkgIdx)

	// PubKey
	pubkeyFld := ast.MakeArgument("PubKey", "", -1).SetType(types.UI8).SetPackage(cPkg)
	pubkeyFld.DeclarationSpecifiers = append(pubkeyFld.DeclarationSpecifiers, constants.DECL_ARRAY)
	// pubkeyFld.IsArray = true
	pubkeyFld.Lengths = []types.Pointer{33} // Yes, PubKey is 33 bytes long.
	pubkeyFld.TotalSize = 33                // 33 * 1 byte (ui8)

	// SecKey
	seckeyFld := ast.MakeArgument("SecKey", "", -1).SetType(types.UI8).SetPackage(cPkg)
	seckeyFld.DeclarationSpecifiers = append(seckeyFld.DeclarationSpecifiers, constants.DECL_ARRAY)
	// seckeyFld.IsArray = true
	seckeyFld.Lengths = []types.Pointer{32} // Yes, SecKey is 32 bytes long.
	seckeyFld.TotalSize = 33                // 33 * 1 byte (ui8)

	pubkeyStrct.AddField(prgrm, pubkeyFld.Type, pubkeyFld, nil)
	seckeyStrct.AddField(prgrm, seckeyFld.Type, seckeyFld, nil)

	cPkg.AddStruct(prgrm, pubkeyStrct)
	cPkg.AddStruct(prgrm, seckeyStrct)

	opcodes.RegisterFunction(prgrm, "cipher.GenerateKeyPair", opCipherGenerateKeyPair, nil,
		opcodes.Out(ast.MakeStructParameter(prgrm, "cipher", "PubKey", "pubKey"), ast.MakeStructParameter(prgrm, "cipher", "SecKey", "sec")))
}
