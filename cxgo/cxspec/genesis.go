package cxspec

import "github.com/SkycoinProject/cx-chains/src/cipher"

type GenesisSpec struct {
	SpecVersion      uint64         `json:"spec_version"`
	GenesisTimestamp uint64         `json:"genesis_timestamp"`
	GenesisAddr      cipher.Address `json:"genesis_addr"` // Also the address for
	GenesisCoins     uint64         `json:"genesis_coins"`
	GenesisProgState string `json:"genesis_prog_state"`
}

type NodeSpec struct {

}

type ProtocolParams struct {
	UnconfirmedBurnFactor uint32 `json:"unconfirmed_burn_factor"` // Burn factor for an unconfirmed transaction.
	UnconfirmedMaxTransactionSize uint32 `json:"unconfirmed_max_transaction_size"` // Maximum size for an unconfirmed transaction.
	UnconfirmedMaxDropletPrecision uint8 `json:"unconfirmed_max_precision"` // Maximum number of decimals allowed for an unconfirmed transaction.

	CreateBlockBurnFactor uint32 `json:"create_block"`

}

type NodeParams struct {
	Port int `json:"port"` // Default port for wire protocol.
	WebInterfacePort int `json:"web_interface_port"` // Default port for web interface.
	DefaultConnections []string `json:"default_connections"` // Default bootstrapping nodes (trusted).
	PeerListURL string `json:"peer_list_url"` // URL pointing to a list of 'ip:port' elements (non-trusted).

}

type GovernanceSpec struct {
	SpecVersion    uint64 `json:"spec_version"`
	ChainPublicKey string `json:"chain_public_key"`
}