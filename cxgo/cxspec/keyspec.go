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

// KeySpec represents a key spec.
type KeySpec struct {
	SpecEra string `json:"spec_era"`
	KeyType string `json:"key_type"` // Either "chain_key" or "genesis_key"
	PubKey string `json:"pubkey"`
	SecKey string `json:"seckey,omitempty"`
	Address string `json:"address,omitempty"`
}

// KeySpecFromSecKey generates a KeySpec from a given secret key.
func KeySpecFromSecKey(keyType string, sk cipher.SecKey, incSK, incAddr bool) KeySpec {
	if err := checkKeyType(keyType); err != nil {
		panic(err)
	}

	pk := cipher.MustPubKeyFromSecKey(sk)
	addr := cipher.AddressFromPubKey(pk)

	spec := KeySpec{
		SpecEra: Era,
		KeyType: keyType,
		PubKey:  pk.Hex(),
		SecKey:  "",
		Address: "",
	}

	if incSK {
		spec.SecKey = sk.Hex()
	}
	if incAddr {
		spec.Address = addr.String()
	}

	return spec
}

// KeySpecFromPubKey generates a KeySpec from a given public key.
func KeySpecFromPubKey(keyType string, pk cipher.PubKey, incAddr bool) KeySpec {
	if err := checkKeyType(keyType); err != nil {
		panic(err)
	}

	addr := cipher.AddressFromPubKey(pk)

	spec := KeySpec{
		SpecEra: Era,
		KeyType: keyType,
		PubKey:  pk.Hex(),
		SecKey:  "",
		Address: "",
	}

	if incAddr {
		spec.Address = addr.String()
	}

	return spec
}

func checkKeyType(keyType string) error {
	switch keyType {
	case ChainKey, GenesisKey:
		return nil
	default:
		return fmt.Errorf("invalid key type '%s'", keyType)
	}
}
