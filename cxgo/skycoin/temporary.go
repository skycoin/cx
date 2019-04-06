package skycoin

import (
	"flag"
	_ "net/http/pprof"
	"os"

	"github.com/skycoin/skycoin/src/readable"
	"github.com/skycoin/skycoin/src/skycoin"
	"github.com/skycoin/skycoin/src/util/logging"
)

var (
	logger = logging.MustGetLogger("main")

	// CoinName name of coin
	CoinName = "cxcoin"

	// GenesisSignatureStr hex string of genesis signature
	GenesisSignatureStr = "a214e0361ff99d80d2f9d646b25f93b8d1d2deb9f7bae0ff908d2302193d8cc31b8388b7bd38c019304b932bfd570444dbe8561aa9d47da021fd31a70146defd01"
	// GenesisAddressStr genesis address string
	GenesisAddressStr = "23v7mT1uLpViNKZHh9aww4VChxizqKsNq4E"
	// BlockchainPubkeyStr pubic key string
	BlockchainPubkeyStr = "02583e5ebbf85522474e0f17e681e62ca37910db6b8792763af4e97663c31a7984"
	// BlockchainSeckeyStr empty private key string
	BlockchainSeckeyStr = ""

	// GenesisTimestamp genesis block create unix time
	GenesisTimestamp uint64 = 1426562704
	// GenesisCoinVolume represents the coin capacity
	GenesisCoinVolume uint64 = 100000000000000

	// DefaultConnections the default trust node addresses
	DefaultConnections = []string{
	}

	nodeConfig = skycoin.NewNodeConfig(ConfigMode, skycoin.NodeParameters{
		CoinName:                       CoinName,
		GenesisSignatureStr:            GenesisSignatureStr,
		GenesisAddressStr:              GenesisAddressStr,
		GenesisCoinVolume:              GenesisCoinVolume,
		GenesisTimestamp:               GenesisTimestamp,
		BlockchainPubkeyStr:            BlockchainPubkeyStr,
		BlockchainSeckeyStr:            BlockchainSeckeyStr,
		DefaultConnections:             DefaultConnections,
		PeerListURL:                    "https://127.0.0.1/peers.txt",
		Port:                           6000,
		WebInterfacePort:               6420,
		DataDirectory:                  "$HOME/.cxcoin",
		UnconfirmedBurnFactor:          2,
		UnconfirmedMaxTransactionSize:  65535,
		UnconfirmedMaxDropletPrecision: 3,
		CreateBlockBurnFactor:          2,
		CreateBlockMaxTransactionSize:  65535,
		CreateBlockMaxDropletPrecision: 3,
		MaxBlockSize:                   65535,
	})

	parseFlags = true
)
