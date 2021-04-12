// +build cipher

package cipher

import (
	"github.com/skycoin/cx/cx/ast"
	"github.com/skycoin/skycoin/src/cipher"
)

// opCipherGenerateKeyPair generates a PubKey and a SecKey.
func opCipherGenerateKeyPair(inputs []ast.CXValue, outputs []ast.CXValue) {
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

    outputs[0].Set_bytes(bPubKey)
    outputs[1].Set_bytes(bSecKey)
}
