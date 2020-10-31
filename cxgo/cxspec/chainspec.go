package cxspec

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/SkycoinProject/cx-chains/src/cipher"
	"github.com/SkycoinProject/cx-chains/src/coin"
)

const (
	Era = "cx_alpha"
)

var progSEnc = base64.StdEncoding

// ProtocolParams defines the coin's consensus parameters.
type ProtocolParams struct {
	UnconfirmedBurnFactor          uint32 `json:"unconfirmed_burn_factor"`           // Burn factor for an unconfirmed transaction.
	UnconfirmedMaxTransactionSize  uint32 `json:"unconfirmed_max_transaction_size"`  // Maximum size for an unconfirmed transaction.
	UnconfirmedMaxDropletPrecision uint8  `json:"unconfirmed_max_droplet_precision"` // Maximum number of decimals allowed for an unconfirmed transaction.

	CreateBlockBurnFactor          uint32 `json:"create_block_burn_factor"`           // Burn factor to transactions when publishing blocks.
	CreateBlockMaxTransactionSize  uint32 `json:"create_block_max_transaction_size"`  // Maximum size of a transaction when publishing blocks.
	CreateBlockMaxDropletPrecision uint8  `json:"create_block_max_droplet_precision"` // Maximum number of decimals allowed for a transaction when publishing blocks.

	MaxBlockTransactionSize uint32 `json:"max_block_transaction_size"` // Maximum total size of transactions when publishing a block.
}

// DefaultProtocolParams returns default values for ProtocolParams.
func DefaultProtocolParams() ProtocolParams {
	return ProtocolParams{
		UnconfirmedBurnFactor:          10,
		UnconfirmedMaxTransactionSize:  5 * 1024 * 1024,
		UnconfirmedMaxDropletPrecision: 3,
		CreateBlockBurnFactor:          10,
		CreateBlockMaxTransactionSize:  5 * 1024 * 1024,
		CreateBlockMaxDropletPrecision: 3,
		MaxBlockTransactionSize:        5 * 1024 * 1024,
	}
}

// NodeParams defines the coin's default node parameters.
// TODO @evanlinjin: In the future, we may use the same network for different cx-chains.
// TODO: If that ever comes to light, we can remove these.
type NodeParams struct {
	Port               int      `json:"port"`                // Default port for wire protocol.
	WebInterfacePort   int      `json:"web_interface_port"`  // Default port for web interface.
	DefaultConnections []string `json:"default_connections"` // Default bootstrapping nodes (trusted).
	PeerListURL        string   `json:"peer_list_url"`       // URL pointing to a list of 'ip:port' elements (non-trusted).

	/* Parameters for user-generated transactions. */
	UserBurnFactor          uint64 `json:"user_burn_factor"`           // Inverse fraction of coin hours that must be burned (used when creating transactions).
	UserMaxTransactionSize  uint32 `json:"user_max_transaction_size"`  // Maximum size of user-created transactions (typically equal to the max size of a block).
	UserMaxDropletPrecision uint64 `json:"user_max_droplet_precision"` // Decimal precision of droplets (smallest coin unit).
}

// DefaultNodeParams returns the default values for NodeParams.
func DefaultNodeParams() NodeParams {
	return NodeParams{
		Port:             6001,
		WebInterfacePort: 6421,
		DefaultConnections: []string{
			"127.0.0.1:6001",
		},
		PeerListURL:             "https://127.0.0.1/peers.txt",
		UserBurnFactor:          10,
		UserMaxTransactionSize:  32 * 1024,
		UserMaxDropletPrecision: 3,
	}
}

// ChainSpec is...
// Functions:
// - Generate genesis block hash.
// - All checks.
type ChainSpec struct {
	SpecEra string `json:"spec_era"`

	ChainPubKey string `json:"chain_pubkey"` // Blockchain public key.

	Protocol ProtocolParams `json:"protocol"` // Params that define the transaction protocol.
	Node     NodeParams     `json:"node"`     // Default params for a node of given coin (this may be removed in future eras).

	/* Identity Params */
	CoinName        string `json:"coin_name"`         // Coin display name (e.g. Skycoin).
	CoinTicker      string `json:"coin_ticker"`       // Coin price ticker (e.g. SKY).
	CoinHoursName   string `json:"coin_hours_name"`   // Coin hours display name (e.g. Skycoin Coin Hours).
	CoinHoursTicker string `json:"coin_hours_ticker"` // Coin hours ticker (e.g SCH).

	/* Genesis Params */
	GenesisAddr       string `json:"genesis_address"`       // Genesis address (base58 representation).
	GenesisSig        string `json:"genesis_signature"`     // Genesis signature (hex representation).
	GenesisCoinVolume uint64 `json:"genesis_coin_volume"`   // Genesis coin volume.
	GenesisProgState  string `json:"genesis_program_state"` // Initial program state on genesis addr (hex representation).
	GenesisTimestamp  uint64 `json:"genesis_timestamp"`     // Timestamp of genesis block (in seconds, UTC time).

	/* Distribution Params */
	// TODO @evanlinjin: Figure out if these are needed for the time being.
	MaxCoinSupply uint64 `json:"max_coin_supply"` // Maximum coin supply.
	// InitialUnlockedCount      uint64   `json:"initial_unlocked_count"`       // Initial number of unlocked addresses.
	// UnlockAddressRate         uint64   `json:"unlock_address_rate"`          // Number of addresses to unlock per time interval.
	// UnlockAddressTimeInterval uint64   `json:"unlock_address_time_interval"` // Time interval (in seconds) in which addresses are unlocked. Once the InitialUnlockedCount is exhausted, UnlockAddressRate addresses will be unlocked per UnlockTimeInterval.
	// DistributionAddresses     []string `json:"distribution_addresses"`       // Addresses that receive coins.

	/* post-processed params */
	chainPK      cipher.PubKey
	genAddr      cipher.Address
	genSig       cipher.Sig
	genProgState []byte
	// distAddresses []cipher.Address // TODO @evanlinjin: May not be needed.
}

// New generates a new chain spec.
func New(coin, ticker string, chainSK cipher.SecKey, genesisAddr cipher.Address, genesisProgState []byte) (*ChainSpec, error) {
	coin = strings.ToLower(strings.Replace(coin, " ", "", -1))
	ticker = strings.ToUpper(strings.Replace(ticker, " ", "", -1))

	spec := &ChainSpec{
		SpecEra:     Era,
		ChainPubKey: "", // ChainPubKey is generated at a later step via generateAndSignGenesisBlock
		Protocol:    DefaultProtocolParams(),
		Node:        DefaultNodeParams(),

		CoinName:        coin,
		CoinTicker:      ticker,
		CoinHoursName:   fmt.Sprintf("%s coin hours", coin),
		CoinHoursTicker: fmt.Sprintf("%sCH", ticker),

		GenesisAddr:       genesisAddr.String(),
		GenesisSig:        "", // GenesisSig is generated at a later step via generateAndSignGenesisBlock
		GenesisCoinVolume: 100e12,
		GenesisProgState:  progSEnc.EncodeToString(genesisProgState),
		GenesisTimestamp:  uint64(time.Now().UTC().Unix()),

		MaxCoinSupply: 1e8,
	}

	// Fill post-processed fields.
	if err := postProcess(spec, true); err != nil {
		return nil, err
	}

	// Generate genesis signature.
	if _, err := generateAndSignGenesisBlock(spec, chainSK); err != nil {
		return nil, err
	}

	return spec, nil
}

func (cs *ChainSpec) RawGenesisProgState() []byte {
	b, err := progSEnc.DecodeString(cs.GenesisProgState)
	if err != nil {
		panic(err)
	}
	return b
}

// GenerateGenesisBlock generates a genesis block from the chain spec and verifies it.
// It returns an error if anything fails.
func (cs *ChainSpec) GenerateGenesisBlock() (*coin.Block, error) {
	if err := postProcess(cs, false); err != nil {
		return nil, err
	}

	block, err := generateGenesisBlock(cs)
	if err != nil {
		return nil, err
	}

	if err := cipher.VerifyPubKeySignedHash(cs.chainPK, cs.genSig, block.HashHeader()); err != nil {
		return nil, err
	}

	return block, nil
}

// GenerateAndSignGenesisBlock is used to generate and sign a genesis block.
// Associated fields in chain spec are also updated.
func (cs *ChainSpec) GenerateAndSignGenesisBlock(sk cipher.SecKey) (*coin.Block, error) {
	if err := postProcess(cs, true); err != nil {
		return nil, err
	}

	return generateAndSignGenesisBlock(cs, sk)
}

// Sign signs the genesis block and fills the chain spec with the resultant data.
// The fields changed are ChainPubKey and GenesisSig (and also the post-processed fields).
func (cs *ChainSpec) Sign(sk cipher.SecKey) error {
	if err := postProcess(cs, true); err != nil {
		return err
	}

	_, err := generateAndSignGenesisBlock(cs, sk)
	return err
}

func (cs ChainSpec) ProcessedChainPubKey() cipher.PubKey {
	return cs.chainPK
}

func (cs *ChainSpec) Print() {
	b, err := json.MarshalIndent(cs, "", "\t")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))
}

// SpecHash returns the hashed spec object.
func (cs *ChainSpec) SpecHash() cipher.SHA256 {
	b, err := json.Marshal(cs)
	if err != nil {
		panic(err)
	}
	return cipher.SumSHA256(b)
}

// postProcess fills post-process fields of chain spec.
// The 'allowEmpty'
func postProcess(cs *ChainSpec, allowEmpty bool) error {
	wrapErr := func(name string, err error) error {
		return fmt.Errorf("chain spec: failed to post-process '%s': %w", name, err)
	}
	var err error

	if !allowEmpty || cs.ChainPubKey != "" {
		if cs.chainPK, err = cipher.PubKeyFromHex(cs.ChainPubKey); err != nil {
			return wrapErr("chain_pubkey", err)
		}
	}
	if cs.genAddr, err = cipher.DecodeBase58Address(cs.GenesisAddr); err != nil {
		return wrapErr("genesis_address", err)
	}
	if !allowEmpty || cs.GenesisSig != "" {
		if cs.genSig, err = cipher.SigFromHex(cs.GenesisSig); err != nil {
			return wrapErr("genesis_signature", err)
		}
	}
	if cs.genProgState, err = progSEnc.DecodeString(cs.GenesisProgState); err != nil {
		return wrapErr("genesis_prog_state", err)
	}
	return nil
}

// generateGenesisBlock generates a genesis block from the chain spec with no checks.
func generateGenesisBlock(cs *ChainSpec) (*coin.Block, error) {
	block, err := coin.NewGenesisBlock(cs.genAddr, cs.GenesisCoinVolume, cs.GenesisTimestamp, cs.genProgState)
	if err != nil {
		return nil, fmt.Errorf("chain spec: %w", err)
	}
	return block, nil
}

// generateAndSignGenesisBlock generates and signs the genesis block (using specified fields from chain spec).
// Hence, ChainPubKey and GenesisSig fields are also filled.
func generateAndSignGenesisBlock(cs *ChainSpec, sk cipher.SecKey) (*coin.Block, error) {
	block, err := generateGenesisBlock(cs)
	if err != nil {
		return nil, err
	}

	pk, err := cipher.PubKeyFromSecKey(sk)
	if err != nil {
		return nil, err
	}

	cs.chainPK = pk
	cs.ChainPubKey = pk.Hex()

	blockSig, err := cipher.SignHash(block.HashHeader(), sk)
	if err != nil {
		return nil, fmt.Errorf("failed to sign genesis block: %w", err)
	}

	cs.genSig = blockSig
	cs.GenesisSig = blockSig.Hex()

	return block, nil
}
