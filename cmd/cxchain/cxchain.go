package main

import (
	"flag"
	"os"

	"github.com/SkycoinProject/cx-chains/src/api"
	"github.com/SkycoinProject/cx-chains/src/cipher"
	"github.com/SkycoinProject/cx-chains/src/readable"
	"github.com/SkycoinProject/cx-chains/src/skycoin"
	"github.com/SkycoinProject/cx-chains/src/util/logging"

	"github.com/SkycoinProject/cx/cxgo/cxspec"
)

const (
	// ENV for chain spec filepath.
	specFileEnv          = "CXCHAIN_SPEC_FILEPATH"
	defaultChainSpecFile = "./skycoin.chain_spec.json"

	// ENV for the chain secret key (in hex).
	secKeyEnv = "CXCHAIN_SECRET_KEY"

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

	spec.Print()
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

// serveGateway does what it is intended to do (in the near future).
func serveGateway(gw api.Gatewayer, close chan struct{}) {
	// TODO @evanlinjin: Actually implement this.
}

func main() {
	// Parse chain spec file and secret key from envs.
	spec := parseSpecFilepathEnv() // Chain spec file (mandatory).
	chainSK := parseSecretKeyEnv() // Secret Key file (optional).

	// Node config: Init.
	conf := cxspec.BaseNodeConfig()
	ensureConfMode(&conf)

	// Node config: Populate node config based on chain spec content.
	if err := cxspec.PopulateNodeConfig(spec, &conf); err != nil {
		log.WithError(err).Fatal("Failed to parse from chain spec file.")
	}

	// Node config: If chain secret key is defined, node is block publisher.
	if !chainSK.Null() {
		conf.BlockchainSeckeyStr = chainSK.Hex()
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

	closeCh := make(chan struct{})
	defer close(closeCh)

	go func() {
		gw, ok := <-gwCh
		if !ok {
			return
		}
		serveGateway(gw, closeCh)
	}()

	if err := coin.Run(spec.RawGenesisProgState(), gwCh); err != nil {
		os.Exit(1)
	}
}
