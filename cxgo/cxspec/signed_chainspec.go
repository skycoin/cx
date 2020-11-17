package cxspec

import (
	"fmt"

	"github.com/SkycoinProject/cx-chains/src/cipher"
)

// SignedChainSpec contains a chain spec alongside a valid signature.
type SignedChainSpec struct {
	Spec ChainSpec `json:"spec"`
	Sig  string    `json:"sig"` // hex representation of signature
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
