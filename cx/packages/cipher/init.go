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

	// PubKey
	pubkeyFld := ast.MakeArgument("PubKey", "", -1).AddType(types.UI8).AddPackage(cipherPkg)
	pubkeyFld.DeclarationSpecifiers = append(pubkeyFld.DeclarationSpecifiers, constants.DECL_ARRAY)
	// pubkeyFld.IsArray = true
	pubkeyFld.Lengths = []types.Pointer{33} // Yes, PubKey is 33 bytes long.
	pubkeyFld.TotalSize = 33                // 33 * 1 byte (ui8)

	// SecKey
	seckeyFld := ast.MakeArgument("SecKey", "", -1).AddType(types.UI8).AddPackage(cipherPkg)
	seckeyFld.DeclarationSpecifiers = append(seckeyFld.DeclarationSpecifiers, constants.DECL_ARRAY)
	// seckeyFld.IsArray = true
	seckeyFld.Lengths = []types.Pointer{32} // Yes, SecKey is 32 bytes long.
	seckeyFld.TotalSize = 33                // 33 * 1 byte (ui8)

	pubkeyStrct.AddField(pubkeyFld)
	seckeyStrct.AddField(seckeyFld)

	cipherPkg.AddStruct(pubkeyStrct)
	cipherPkg.AddStruct(seckeyStrct)

	prgrm.AddPackage(cipherPkg)

	opcodes.RegisterFunction(prgrm, "cipher.GenerateKeyPair", opCipherGenerateKeyPair, nil,
		opcodes.Out(ast.Struct(prgrm, "cipher", "PubKey", "pubKey"), ast.Struct(prgrm, "cipher", "SecKey", "sec")))
}
