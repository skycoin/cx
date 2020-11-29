package cxspec

import (
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
		return fmt.Errorf("field 'last_seen' is invalid '%s'", pe.LastSeen)
	}

	return nil
}
