package cxspec

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"github.com/SkycoinProject/cx-chains/src/api"
	"github.com/SkycoinProject/cx-chains/src/kvstorage"
	"github.com/SkycoinProject/cx-chains/src/params"
	"github.com/SkycoinProject/cx-chains/src/readable"
	"github.com/SkycoinProject/cx-chains/src/skycoin"
	"github.com/SkycoinProject/cx-chains/src/wallet"
)

// PopulateParamsModule populates the params module within cx chain.
func PopulateParamsModule(cs ChainSpec) {
	// TODO @evanlinjin: Figure out distribution.
	params.MainNetDistribution = params.Distribution{
		MaxCoinSupply:        cs.MaxCoinSupply,
		InitialUnlockedCount: 1,
		UnlockAddressRate:    0,
		UnlockTimeInterval:   0,
		Addresses:            []string{cs.GenesisAddr},
	}
	params.UserVerifyTxn = params.VerifyTxn{
		BurnFactor:          uint32(cs.Node.UserBurnFactor),
		MaxTransactionSize:  cs.Node.UserMaxTransactionSize,
		MaxDropletPrecision: uint8(cs.Node.UserMaxDropletPrecision),
	}
	params.InitFromEnv()
}

// PopulateNodeConfig populates the node config with values from cx chain spec.
func PopulateNodeConfig(spec ChainSpec, conf *skycoin.NodeConfig) error {
	if spec.SpecEra != Era {
		return fmt.Errorf("unsupported spec era '%s'", spec.SpecEra)
	}

	// conf := DefaultNodeConfig(defaultSpecFilename)
	conf.CoinName = spec.CoinName
	conf.PeerListURL = spec.Node.PeerListURL
	conf.Port = spec.Node.Port
	conf.WebInterfacePort = spec.Node.WebInterfacePort
	conf.UnconfirmedVerifyTxn = params.VerifyTxn{
		BurnFactor:          spec.Protocol.UnconfirmedBurnFactor,
		MaxTransactionSize:  spec.Protocol.UnconfirmedMaxTransactionSize,
		MaxDropletPrecision: spec.Protocol.UnconfirmedMaxDropletPrecision,
	}
	conf.CreateBlockVerifyTxn = params.VerifyTxn{
		BurnFactor:          spec.Protocol.CreateBlockBurnFactor,
		MaxTransactionSize:  spec.Protocol.CreateBlockMaxTransactionSize,
		MaxDropletPrecision: spec.Protocol.CreateBlockMaxDropletPrecision,
	}
	conf.MaxBlockTransactionsSize = spec.Protocol.MaxBlockTransactionSize
	conf.GenesisSignatureStr = spec.GenesisSig
	conf.GenesisAddressStr = spec.GenesisAddr
	conf.BlockchainPubkeyStr = spec.ChainPubKey
	conf.GenesisTimestamp = spec.GenesisTimestamp
	conf.GenesisCoinVolume = spec.GenesisCoinVolume
	conf.DefaultConnections = spec.Node.DefaultConnections
	conf.Fiber = readable.FiberConfig{
		Name:            spec.CoinName,
		DisplayName:     spec.CoinName,
		Ticker:          spec.CoinTicker,
		CoinHoursName:   spec.CoinHoursName,
		CoinHoursTicker: spec.CoinHoursTicker,
		ExplorerURL:     "", // TODO @evanlinjin: CX Chain explorer?
	}

	if conf.DataDirectory == "" {
		conf.DataDirectory = "$HOME/.cxchain/" + spec.CoinName
	} else {
		conf.DataDirectory = strings.ReplaceAll(conf.DataDirectory, "{coin}", spec.CoinName)
	}

	return nil
}

// ReadSpecFile reads chain spec from given filename.
func ReadSpecFile(filename string) (ChainSpec, error) {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return ChainSpec{}, fmt.Errorf("failed to read chain spec file '%s': %w", filename, err)
	}
	var spec ChainSpec
	if err := json.Unmarshal(b, &spec); err != nil {
		return ChainSpec{}, fmt.Errorf("chain spec file '%s' is ill-formed: %w", filename, err)
	}
	if _, err := spec.GenerateGenesisBlock(); err != nil {
		return ChainSpec{}, fmt.Errorf("chain spec file '%s' cannot generate genesis block: %w", filename, err)
	}
	return spec, nil
}

func ReadKeysFile(filename string) (KeySpec, error) {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return KeySpec{}, fmt.Errorf("failed to read key file '%s': %w", filename, err)
	}
	var spec KeySpec
	if err := json.Unmarshal(b, &spec); err != nil {
		return KeySpec{}, fmt.Errorf("key file '%s' is ill-formed: %w", filename, err)
	}
	return spec, nil
}

// BaseNodeConfig returns the base node config.
// Fields which requires values from the cx chain spec are left blank.
func BaseNodeConfig() skycoin.NodeConfig {
	conf := skycoin.NodeConfig{
		CoinName:                          "", // populate with cx spec
		DisablePEX:                        false,
		DownloadPeerList:                  false,
		PeerListURL:                       "", // populate with cx spec
		DisableOutgoingConnections:        false,
		DisableIncomingConnections:        false,
		DisableNetworking:                 false,
		EnableGUI:                         false,
		DisableCSRF:                       false,
		DisableHeaderCheck:                false,
		DisableCSP:                        false,
		EnabledAPISets:                    "",
		DisabledAPISets:                   "",
		EnableAllAPISets:                  false,
		HostWhitelist:                     "",
		LocalhostOnly:                     true,
		Address:                           "",
		Port:                              0, // populate with cx spec
		MaxConnections:                    128,
		MaxOutgoingConnections:            8,
		MaxDefaultPeerOutgoingConnections: 1,
		OutgoingConnectionsRate:           time.Second * 5,
		MaxOutgoingMessageLength:          5243081 * 2, // TODO @evanlinjin: Find a way to regulate this with cx txns (originally 256 * 1024).
		MaxIncomingMessageLength:          5243081 * 4, // TODO @evanlinjin: Find a way to regulate this with cx txns (originally 1024 * 1024).
		PeerlistSize:                      65535,
		WebInterface:                      true,
		WebInterfacePort:                  0, // populate with cx spec
		WebInterfaceAddr:                  "127.0.0.1",
		WebInterfaceCert:                  "",
		WebInterfaceKey:                   "",
		WebInterfaceHTTPS:                 false,
		WebInterfaceUsername:              "",
		WebInterfacePassword:              "",
		WebInterfacePlaintextAuth:         false,
		LaunchBrowser:                     false,
		DataDirectory:                     "$HOME/.cxchain/{coin}", // populate with cx spec
		GUIDirectory:                      "./src/gui/static/",
		HTTPReadTimeout:                   time.Second * 10,
		HTTPWriteTimeout:                  time.Second * 60,
		HTTPIdleTimeout:                   time.Second * 120,
		UserAgentRemark:                   "",
		ColorLog:                          true,
		LogLevel:                          "INFO",
		DisablePingPong:                   false,
		VerifyDB:                          false,
		ResetCorruptDB:                    false,
		UnconfirmedVerifyTxn:              params.VerifyTxn{}, // populate with cx spec
		CreateBlockVerifyTxn:              params.VerifyTxn{}, // populate with cx spec
		MaxBlockTransactionsSize:          0,                  // populate with cx spec
		WalletDirectory:                   "",
		WalletCryptoType:                  string(wallet.CryptoTypeScryptChacha20poly1305),
		KVStorageDirectory:                "",
		EnabledStorageTypes: []kvstorage.Type{
			kvstorage.TypeTxIDNotes,
			kvstorage.TypeGeneral,
		},
		DisableDefaultPeers: false,
		CustomPeersFile:     "",
		RunBlockPublisher:   false,
		ProfileCPU:          false,
		ProfileCPUFile:      "cpu.prof",
		HTTPProf:            false,
		HTTPProfHost:        "localhost:6060",
		DBPath:              "",
		DBReadOnly:          false,
		LogToFile:           false,
		Version:             false,
		GenesisSignatureStr: "", // populate with cx spec
		GenesisAddressStr:   "", // populate with cx spec
		BlockchainPubkeyStr: "", // populate with cx spec
		BlockchainSeckeyStr: "",
		GenesisTimestamp:    0,                      // populate with cx spec
		GenesisCoinVolume:   0,                      // populate with cx spec
		DefaultConnections:  nil,                    // populate with cx spec
		Fiber:               readable.FiberConfig{}, // populate with cx spec
	}

	return conf
}

// ApplyStandaloneClientMode alters a node config for standalone node use-cases.
func ApplyStandaloneClientMode(conf *skycoin.NodeConfig) {
	conf.EnableAllAPISets = true
	conf.EnabledAPISets = api.EndpointsInsecureWalletSeed
	conf.EnableGUI = true
	conf.LaunchBrowser = true
	conf.DisableCSRF = false
	conf.DisableHeaderCheck = false
	conf.DisableCSP = false
	conf.DownloadPeerList = true
	conf.WebInterface = true
	conf.LogToFile = false
	conf.ResetCorruptDB = true
	conf.WebInterfacePort = 0 // randomize web interface port
}
