package cxspec

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/skycoin/dmsg"
	"github.com/skycoin/dmsg/cipher"
)

// CXChainAddresses contains the addresses of a cx node within cx tracker.
type CXChainAddresses struct {
	DmsgAddr dmsg.Addr `json:"dmsg_addr"`
	TCPAddr  string    `json:"tcp_addr,omitempty"`
}

// PeerEntry represents a peer entry in cx tracker.
type PeerEntry struct {
	PublicKey cipher.PubKey                      `json:"public_key"`
	LastSeen  int64                              `json:"last_seen"`
	CXChains  map[cipher.SHA256]CXChainAddresses `json:"cx_chains"`
}

// Check ensures all fields of PeerEntry conforms to all rules.
func (pe *PeerEntry) Check() error {
	const lastSeenTolerance = int64(time.Minute)

	if pe.PublicKey.Null() {
		return fmt.Errorf("field 'public_key' cannot have value '%s'", pe.PublicKey)
	}

	now := time.Now().Unix()

	if pe.LastSeen < now-lastSeenTolerance || pe.LastSeen > now+lastSeenTolerance {
		return fmt.Errorf("field 'last_seen' is invalid '%d'", pe.LastSeen)
	}

	return nil
}

// Hash hashes the PeerEntry.
func (pe *PeerEntry) Hash() cipher.SHA256 {
	b, err := json.Marshal(pe)
	if err != nil {
		panic(err)
	}
	return cipher.SumSHA256(b)
}

// SignedPeerEntry contains a chain spec alongside a valid signature.
type SignedPeerEntry struct {
	Entry PeerEntry  `json:"entry"`
	Sig   cipher.Sig `json:"sig"`
}

// MakeSignedPeerEntry generates a signed peer entry from a PeerEntry and secret
// key. It checks the PeerEntry's validity before signing.
func MakeSignedPeerEntry(entry PeerEntry, sk cipher.SecKey) (SignedPeerEntry, error) {
	if err := entry.Check(); err != nil {
		return SignedPeerEntry{}, err
	}

	b, err := json.Marshal(entry)
	if err != nil {
		panic(err)
	}

	sig, err := cipher.SignPayload(b, sk)
	if err != nil {
		return SignedPeerEntry{}, err
	}

	signedEntry := SignedPeerEntry{
		Entry: entry,
		Sig:   sig,
	}

	return signedEntry, nil
}
