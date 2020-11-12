// +build base

package cxcore

import (
	"github.com/SkycoinProject/skycoin/src/cipher"

	. "github.com/SkycoinProject/cx/cx"
)

func init() {
	cipherPkg := MakePackage("cipher")
	pubkeyStrct := MakeStruct("PubKey")
	seckeyStrct := MakeStruct("SecKey")

	// PubKey
	pubkeyFld := MakeArgument("PubKey", "", -1).AddType(TypeNames[TYPE_UI8]).AddPackage(cipherPkg)
	pubkeyFld.DeclarationSpecifiers = append(pubkeyFld.DeclarationSpecifiers, DECL_ARRAY)
	pubkeyFld.IsArray = true
	pubkeyFld.Lengths = []int{33} // Yes, PubKey is 33 bytes long.
	pubkeyFld.TotalSize = 33      // 33 * 1 byte (ui8)

	// SecKey
	seckeyFld := MakeArgument("SecKey", "", -1).AddType(TypeNames[TYPE_UI8]).AddPackage(cipherPkg)
	seckeyFld.DeclarationSpecifiers = append(seckeyFld.DeclarationSpecifiers, DECL_ARRAY)
	seckeyFld.IsArray = true
	seckeyFld.Lengths = []int{32} // Yes, SecKey is 32 bytes long.
	seckeyFld.TotalSize = 33      // 33 * 1 byte (ui8)

	pubkeyStrct.AddField(pubkeyFld)
	seckeyStrct.AddField(seckeyFld)

	cipherPkg.AddStruct(pubkeyStrct)
	cipherPkg.AddStruct(seckeyStrct)

	PROGRAM.AddPackage(cipherPkg)
}

// opCipherGenerateKeyPair generates a PubKey and a SecKey.
func opCipherGenerateKeyPair(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()
	out1, out2 := expr.Outputs[0], expr.Outputs[1]

	pubKey, secKey := cipher.GenerateKeyPair()

	bPubKey := make([]byte, len(pubKey))
	bSecKey := make([]byte, len(secKey))

	// Copying bytes
	for i, byt := range pubKey {
		bPubKey[i] = byt
	}
	for i, byt := range secKey {
		bSecKey[i] = byt
	}

	WriteMemory(GetFinalOffset(fp, out1), bPubKey)
	WriteMemory(GetFinalOffset(fp, out2), bSecKey)
}
