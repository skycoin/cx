package main

import (
	"context"
	"encoding/hex"
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
	trackerAddr  = "http://127.0.0.1:9091"           // cx tracker address
	dmsgDiscAddr = cxdmsg.DefaultDiscAddr     // dmsg discovery address
	dmsgPort     = uint64(cxdmsg.DefaultPort) // dmsg listening port
)

func init() {
	cmd := flag.CommandLine
	cmd.StringVar(&trackerAddr, "cx-tracker", trackerAddr, "HTTP `ADDRESS` of cx tracker")
	cmd.StringVar(&dmsgDiscAddr, "dmsg-disc", dmsgDiscAddr, "HTTP `ADDRESS` of dmsg discovery")
	cmd.Uint64Var(&dmsgPort, "dmsg-port", dmsgPort, "dmsg `PORT` number to listen on")
}

func trackerUpdateLoop(nodeSK cipher.SecKey, nodeTCPAddr string, spec cxspec.ChainSpec) {
	log := logging.MustGetLogger("cx_tracker_client")

	client := cxspec.NewCXTrackerClient(log, nil, trackerAddr)
	nodePK := cipher.MustPubKeyFromSecKey(nodeSK)

	block, err := spec.GenerateGenesisBlock()
	if err != nil {
		panic(err) // should not happen
	}
	hash := block.HashHeader()

	// If publisher, ensure spec is registered.
	if isPub := nodePK == spec.ProcessedChainPubKey(); isPub {
		signedSpec, err := cxspec.MakeSignedChainSpec(spec, nodeSK)
		if err != nil {
			panic(err) // should not happen
		}
		for {
			if err := client.PostSpec(context.Background(), signedSpec); err != nil {
				log.WithError(err).Error("Failed to post spec, retrying again...")
				time.Sleep(time.Second * 10)
				continue
			}
			break
		}
	}

	// Prepare ticker for cx tracker peer entry update loop.
	ticker := time.NewTicker(time.Second * 30)
	defer ticker.Stop()

	// All peers need to update entry.
	entry := cxspec.PeerEntry{
		PublicKey: cipher2.PubKey(nodePK),
		CXChains: map[string]cxspec.CXChainAddresses{
			hex.EncodeToString(hash[:]): {
				DmsgAddr: dmsg.Addr{PK: cipher2.PubKey(nodePK), Port: uint16(dmsgPort)},
				TCPAddr:  nodeTCPAddr,
			},
		},
	}
	updateEntry := func(now int64) {
		entry.LastSeen = now
		signedEntry, err := cxspec.MakeSignedPeerEntry(entry, cipher2.SecKey(nodeSK))
		if err != nil {
			panic(err) // should not happen
		}
		if err := client.UpdatePeerEntry(context.Background(), signedEntry); err != nil {
			log.WithError(err).Warn("Failed to update peer entry in cx tracker. Retrying...")
		}
	}
	updateEntry(time.Now().Unix())
	for now := range ticker.C {
		updateEntry(now.Unix())
	}
}

func main() {
	// Parse chain spec file and secret key from envs.
	spec := parseSpecFilepathEnv() // Chain spec file (mandatory).
	nodeSK := parseSecretKeyEnv()  // Secret Key file (mandatory).
	var nodePK cipher.PubKey

	// Node config: Init.
	conf := cxspec.BaseNodeConfig()
	ensureConfMode(&conf)

	// Node config: Populate node config based on chain spec content.
	if err := cxspec.PopulateNodeConfig(trackerAddr, spec, &conf); err != nil {
		log.WithError(err).Fatal("Failed to parse from chain spec file.")
	}

	// Node config: Check node secret key.
	//	- If node secret key is null, randomly generate one.
	//	- If node secret key generates spec's chain pk, it is also the chain's publisher node.
	if nodeSK.Null() {
		nodePK, nodeSK = cipher.GenerateKeyPair()
		log.WithField("node_pk", nodePK.Hex()).
			Warn("Node secret key is not defined. Random key pair generated.")
	}
	if err := nodeSK.Verify(); err != nil {
		log.WithError(err).Fatal("Failed to verify provided node secret key.")
	}
	if nodePK = cipher.MustPubKeyFromSecKey(nodeSK); nodePK == spec.ProcessedChainPubKey() {
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

	go func() {
		gw, ok := <-gwCh
		if !ok {
			return
		}

		// Run cx tracker loop.
		go trackerUpdateLoop(nodeSK, gw.DaemonConfig().Address, spec)

		// Run dmsg loop.
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

	if err := coin.Run(spec.RawGenesisProgState(), gwCh); err != nil {
		os.Exit(1)
	}
}
