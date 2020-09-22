package chainspec

import (
	"fmt"

	"github.com/SkycoinProject/cx-chains/src/cipher"
)

// For KeySpec.KeyType
const (
	ChainKey   = "chain_key"
	GenesisKey = "genesis_key"
)

type KeySpec struct {
	KeyType string `json:"key_type"` // Either "chain_key" or "genesis_key"
	PubKey string `json:"pubkey"`
	SecKey string `json:"seckey,omitempty"`
	Address string `json:"address,omitempty"`
}

func KeySpecFromSecKey(keyType string, sk cipher.SecKey, incSK, incAddr bool) KeySpec {
	return KeySpec{}
}

func checkKeyType(keyType string) error {
	switch keyType {
	case ChainKey, GenesisKey:
		return nil
	default:
		return fmt.Errorf("")
	}
}