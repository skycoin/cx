package cxspec

import (
	"fmt"

	"github.com/SkycoinProject/cx-chains/src/cipher"
)

// SignedChainSpec contains a chain spec alongside a valid signature.
type SignedChainSpec struct {
	Spec        ChainSpec `json:"spec"`
	GenesisHash string    `json:"genesis_hash,omitempty"`
	Sig         string    `json:"sig"` // hex representation of signature
}

// MakeSignedChainSpec generates a signed spec from a ChainSpec and secret key.
// Note that the secret key needs to be able to generate the ChainSpec's public
// key to be valid.
func MakeSignedChainSpec(spec ChainSpec, sk cipher.SecKey) (SignedChainSpec, error) {
	genesis, err := spec.GenerateGenesisBlock()
	if err != nil {
		return SignedChainSpec{}, fmt.Errorf("chain spec failed to generate genesis block: %w", err)
	}

	pk, err := cipher.PubKeyFromSecKey(sk)
	if err != nil {
		return SignedChainSpec{}, err
	}

	if pk != spec.ProcessedChainPubKey() {
		return SignedChainSpec{}, fmt.Errorf("provided sk does not generate chain pk '%s'", spec.ChainPubKey)
	}

	sig, err := cipher.SignHash(spec.SpecHash(), sk)
	if err != nil {
		return SignedChainSpec{}, err
	}

	signedSpec := SignedChainSpec{
		Spec:        spec,
		GenesisHash: genesis.HashHeader().Hex(),
		Sig:         sig.Hex(),
	}

	return signedSpec, nil
}

// Verify checks the following:
// - Spec is of right era, has valid chain pk, and generates valid genesis block.
// - Signature is valid
func (ss *SignedChainSpec) Verify() error {
	if era := ss.Spec.SpecEra; era != Era {
		return fmt.Errorf("unexpected chain spec era '%s' (expected '%s')",
			era, Era)
	}

	if _, err := ss.Spec.GenerateGenesisBlock(); err != nil {
		return fmt.Errorf("chain spec failed to generate genesis block: %w", err)
	}

	sig, err := cipher.SigFromHex(ss.Sig)
	if err != nil {
		return fmt.Errorf("failed to decode spec signature: %w", err)
	}

	pk := ss.Spec.ProcessedChainPubKey()
	hash := ss.Spec.SpecHash()

	if err := cipher.VerifyPubKeySignedHash(pk, sig, hash); err != nil {
		return fmt.Errorf("failed to verify spec signature: %w", err)
	}

	return nil
}
