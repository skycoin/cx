package main

import (
	"context"
	"flag"
	"os"
	"time"

	"github.com/skycoin/dmsg"
	cipher2 "github.com/skycoin/dmsg/cipher"

	"github.com/SkycoinProject/cx-chains/src/api"
	"github.com/SkycoinProject/cx-chains/src/cipher"
	"github.com/SkycoinProject/cx-chains/src/readable"
	"github.com/SkycoinProject/cx-chains/src/skycoin"
	"github.com/SkycoinProject/cx-chains/src/util/logging"

	"github.com/SkycoinProject/cx/cxgo/cxdmsg"
	"github.com/SkycoinProject/cx/cxgo/cxspec"
)

const (
	// ENV for chain spec filepath.
	specFileEnv          = "CXCHAIN_SPEC_FILEPATH"
	defaultChainSpecFile = "./skycoin.chain_spec.json"

	// ENV for the chain secret key (in hex).
	secKeyEnv = "CXCHAIN_SK"

	// ENVs for config modes.
	standaloneClientConfMode = "STANDALONE_CLIENT"
)

// These values should be populated by -ldflags on compilation
var (
	version  = "0.0.0"
	commit   = ""
	branch   = ""
	confMode = "" // valid values: "STANDALONE_CLIENT", ""
)

// Logger.
var log = logging.MustGetLogger("main")

// parseSpecFilepathEnv parses chain spec filename from CXCHAIN_SPEC_FILEPATH env.
func parseSpecFilepathEnv() cxspec.ChainSpec {
	specFilename, ok := os.LookupEnv(specFileEnv)
	if !ok {
		specFilename = defaultChainSpecFile
	}

	log.WithField("filename", specFilename).Info("Reading chain spec file...")

	spec, err := cxspec.ReadSpecFile(specFilename)
	if err != nil {
		log.WithError(err).Fatal("Failed to start node.")
	}

	cxspec.PopulateParamsModule(spec)
	return spec
}

// parseSecretKeyEnv parses secret key from CXCHAIN_SECRET_KEY env.
// The secret key can be null.
func parseSecretKeyEnv() cipher.SecKey {
	if skStr, ok := os.LookupEnv(secKeyEnv); ok {
		sk, err := cipher.SecKeyFromHex(skStr)
		if err != nil {
			log.WithError(err).WithField("ENV", secKeyEnv).Fatal("Provided secret key is invalid.")
		}
		return sk
	}
	return cipher.SecKey{} // return nil secret key
}

// ensureConfMode ensures 'confMode' settings are applied on the node config.
// 'confMode' is set on compile time.
func ensureConfMode(conf *skycoin.NodeConfig) {
	switch confMode {
	case "":
	case standaloneClientConfMode:
		cxspec.ApplyStandaloneClientMode(conf)
	default:
		log.Fatal("Invalid 'confMode' provided at build time. This cannot be fixed without recompiling the binary.")
	}
}

var (
	cxTrackerAddr = "127.0.0.1:9091" // cx tracker address
	dmsgDiscAddr  = "127.0.0.1:9090" // dmsg discovery address
	dmsgPort      = uint64(9090)     // dmsg listening port
)

func init() {
	cmd := flag.CommandLine
	cmd.StringVar(&cxTrackerAddr, "cx-tracker", cxTrackerAddr, "HTTP `ADDRESS` of cx tracker")
	cmd.StringVar(&dmsgDiscAddr, "dmsg-disc", dmsgDiscAddr, "HTTP `ADDRESS` of dmsg discovery")
	cmd.Uint64Var(&dmsgPort, "dmsg-port", dmsgPort, "dmsg `PORT` number to listen on")
}

func trackerUpdateLoop(nodeSK cipher.SecKey, nodeTCPAddr string, spec cxspec.ChainSpec) {
	log := logging.MustGetLogger("cx_tracker_client")

	client := cxspec.NewCXTrackerClient(log, nil, cxTrackerAddr)
	nodePK := cipher.MustPubKeyFromSecKey(nodeSK)

	block, err := spec.GenerateGenesisBlock()
	if err != nil {
		panic(err) // should not happen
	}
	hash := block.HashHeader()

	ticker := time.NewTicker(time.Second * 30)
	defer ticker.Stop()

	entry := cxspec.PeerEntry{
		PublicKey: cipher2.PubKey(nodePK),
		// LastSeen: now.Unix(),
		CXChains: map[cipher2.SHA256]cxspec.CXChainAddresses{
			cipher2.SHA256(hash): {
				DmsgAddr: dmsg.Addr{PK: cipher2.PubKey(nodePK), Port: uint16(dmsgPort)},
				TCPAddr:  nodeTCPAddr,
			},
		},
	}

	for now := range ticker.C {
		entry.LastSeen = now.Unix()

		signedEntry, err := cxspec.MakeSignedPeerEntry(entry, cipher2.SecKey(nodeSK))
		if err != nil {
			panic(err) // should not happen
		}

		if err := client.UpdatePeerEntry(context.Background(), signedEntry); err != nil {
			log.WithError(err).Warn("Failed to update peer entry in cx tracker. Retrying...")
		}
	}
}

func main() {
	// Parse chain spec file and secret key from envs.
	spec := parseSpecFilepathEnv() // Chain spec file (mandatory).
	nodeSK := parseSecretKeyEnv()  // Secret Key file (mandatory).
	nodePK := cipher.MustPubKeyFromSecKey(nodeSK)

	// Node config: Init.
	conf := cxspec.BaseNodeConfig()
	ensureConfMode(&conf)

	// Node config: Populate node config based on chain spec content.
	if err := cxspec.PopulateNodeConfig(spec, &conf); err != nil {
		log.WithError(err).Fatal("Failed to parse from chain spec file.")
	}

	// Node config: Check node secret key.
	//	- Node secret key should be defined.
	//	- If node secret key generates spec's chain pk, it is also the chain's publisher node.
	if nodeSK.Null() {
		log.Fatal("Node secret key is not defined.")
	}
	if err := nodeSK.Verify(); err != nil {
		log.WithError(err).Fatal("Failed to verify provided node secret key.")
	}
	if cipher.MustPubKeyFromSecKey(nodeSK) == spec.ProcessedChainPubKey() {
		conf.BlockchainSeckeyStr = nodeSK.Hex()
		conf.RunBlockPublisher = true
	}

	// Node config: Parse flag set.
	conf.RegisterFlags(flag.CommandLine)
	flag.Parse()

	coin := skycoin.NewCoin(skycoin.Config{
		Node: conf,
		Build: readable.BuildInfo{
			Version: version,
			Commit:  commit,
			Branch:  branch,
		},
	}, log)

	// This is, despite the name, post-processing and not "parsing".
	// Do not get confused. I did not name this function. <3 @evanlinjin
	if err := coin.ParseConfig(flag.CommandLine); err != nil {
		log.Error(err)
		os.Exit(1)
	}

	gwCh := make(chan api.Gatewayer)
	defer close(gwCh)

	// Run dmsg loop.
	go func() {
		gw, ok := <-gwCh
		if !ok {
			return
		}
		cxdmsg.ServeDmsg(
			context.Background(),
			logging.MustGetLogger("dmsgC"),
			&cxdmsg.Config{
				PK:       cipher2.PubKey(nodePK),
				SK:       cipher2.SecKey(nodeSK),
				DiscAddr: dmsgDiscAddr,
				DmsgPort: uint16(dmsgPort),
			},
			&cxdmsg.API{
				Version:   version,
				NodeConf:  conf,
				ChainSpec: spec,
				Gateway:   gw,
			},
		)
	}()

	// Run cx tracker loop.
	go trackerUpdateLoop(nodeSK, conf.Address, spec)

	if err := coin.Run(spec.RawGenesisProgState(), gwCh); err != nil {
		os.Exit(1)
	}
}
