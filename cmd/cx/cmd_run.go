package main

import (
	"flag"
	"fmt"
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
	keyEnv = "CXCHAIN_SECRET_KEY"

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

func cmdRun(args []string) {
	// logger
	log := logging.MustGetLogger("main")

	// Prepare ENV: CXCHAIN_SPEC_FILEPATH (filepath of chain spec file).
	specFilename, ok := os.LookupEnv(specFileEnv)
	if !ok {
		specFilename = defaultChainSpecFile
	}
	fmt.Printf("Reading chain spec file from '%s'...\n", specFilename)
	spec, err := cxspec.ReadSpecFile(specFilename)
	if err != nil {
		errPrintf("Failed to start node: %v\n", err)
		os.Exit(1)
	}
	spec.Print()
	cxspec.PopulateParamsModule(spec)

	// Prepare ENV CXCHAIN_SECRET_KEY (chain private key in hex).
	var chainSK cipher.SecKey
	if skStr, ok := os.LookupEnv(keyEnv); ok {
		sk, err := cipher.SecKeyFromHex(skStr)
		if err != nil {
			errPrintf("Provided secret key '%s' in ENV '%s' is invalid: %v\n", skStr, keyEnv, err)
			os.Exit(1)
		}
		chainSK = sk
	}

	// Prepare node config.
	conf := cxspec.BaseNodeConfig()
	ensureConfMode(&conf)

	if err := cxspec.PopulateNodeConfig(spec, &conf); err != nil {
		errPrintf("Failed to parse from chain spec file: %v\n", err)
		os.Exit(1)
	}

	if !chainSK.Null() {
		conf.BlockchainSeckeyStr = chainSK.Hex()
		conf.RunBlockPublisher = true
	}

	// Parse flag set into node config.
	cmd := flag.NewFlagSet(args[0], flag.ExitOnError)

	conf.RegisterFlags(cmd)
	parseFlagSet(cmd, args[1:])



	coin := skycoin.NewCoin(skycoin.Config{
		Node: conf,
		Build: readable.BuildInfo{
			Version: version,
			Commit:  commit,
			Branch:  branch,
		},
	}, log)

	if err := coin.ParseConfig(cmd); err != nil {
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

	if err := coin.Run(gwCh); err != nil {
		os.Exit(1)
	}
}

func ensureConfMode(conf *skycoin.NodeConfig) {
	switch confMode {
	case "":
	case standaloneClientConfMode:
		cxspec.ApplyStandaloneClientMode(conf)
	default:
		panic("Invalid 'confMode' provided at build time. This cannot be fixed without recompiling the binary.")
	}
}

func serveGateway(gw api.Gatewayer, close chan struct{}) {
	// TODO
}
