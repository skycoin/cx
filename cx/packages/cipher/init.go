// +build cipher

package cipher

import (
	"github.com/skycoin/cx/cx/ast"
	"github.com/skycoin/cx/cx/constants"
	"github.com/skycoin/cx/cx/opcodes"
)

func RegisterPackage() {
	cipherPkg := ast.MakePackage("cipher")
	pubkeyStrct := ast.MakeStruct("PubKey")
	seckeyStrct := ast.MakeStruct("SecKey")

	// PubKey
	pubkeyFld := ast.MakeArgument("PubKey", "", -1).AddType(constants.TypeNames[constants.TYPE_UI8]).AddPackage(cipherPkg)
	pubkeyFld.DeclarationSpecifiers = append(pubkeyFld.DeclarationSpecifiers, constants.DECL_ARRAY)
	// pubkeyFld.IsArray = true
	pubkeyFld.Lengths = []int{33} // Yes, PubKey is 33 bytes long.
	pubkeyFld.TotalSize = 33      // 33 * 1 byte (ui8)

	// SecKey
	seckeyFld := ast.MakeArgument("SecKey", "", -1).AddType(constants.TypeNames[constants.TYPE_UI8]).AddPackage(cipherPkg)
	seckeyFld.DeclarationSpecifiers = append(seckeyFld.DeclarationSpecifiers, constants.DECL_ARRAY)
	// seckeyFld.IsArray = true
	seckeyFld.Lengths = []int{32} // Yes, SecKey is 32 bytes long.
	seckeyFld.TotalSize = 33      // 33 * 1 byte (ui8)

	pubkeyStrct.AddField(pubkeyFld)
	seckeyStrct.AddField(seckeyFld)

	cipherPkg.AddStruct(pubkeyStrct)
	cipherPkg.AddStruct(seckeyStrct)

	ast.PROGRAM.AddPackage(cipherPkg)

	opcodes.RegisterFunction("cipher.GenerateKeyPair", opCipherGenerateKeyPair, nil,
		opcodes.Out(opcodes.Struct("cipher", "PubKey", "pubKey"), opcodes.Struct("cipher", "SecKey", "sec")))
}
