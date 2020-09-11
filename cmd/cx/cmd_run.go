package main

import (
	"flag"
	"os"
	"path/filepath"

	"github.com/SkycoinProject/cx-chains/src/api"
	"github.com/SkycoinProject/cx-chains/src/fiber"
	"github.com/SkycoinProject/cx-chains/src/readable"
	"github.com/SkycoinProject/cx-chains/src/skycoin"
	"github.com/SkycoinProject/cx-chains/src/util/logging"
)

const (
	specFileEnv     = "CXCHAIN_SPEC"
	defaultSpecFile = "./cxchain.spec.toml"
)

// should be populated by -ldflags on compilation
var (
	version string
	commit  string
	branch  string
)

func cmdRun(args []string) {
	// env: CXCHAIN_SPEC (filepath of spec file)
	specEnv, ok := os.LookupEnv(specFileEnv)
	if !ok {
		specEnv = defaultSpecFile
	}
	specBase := filepath.Base(specEnv)
	specDir := filepath.Dir(specEnv)

	// spec file
	spec, err := fiber.NewConfig(specBase, specDir)
	if err != nil {
		errPrintf("Failed to read blockchain spec file: %v\n", err)
		os.Exit(1)
	}

	cmd := flag.NewFlagSet(args[0], flag.ExitOnError)

	// node config (from spec file)
	nodeConf := skycoin.NewNodeConfig("", spec.Node)
	nodeConf.RegisterFlags(cmd)

	// flags: parse
	parseFlagSet(cmd, args)

	// logger
	log := logging.MustGetLogger("main")

	coin := skycoin.NewCoin(skycoin.Config{
		Node: nodeConf,
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

	go func() {
		gw, ok := <-gwCh
		if !ok {
			return
		}

		// TODO
	}()

	if err := coin.Run(gwCh); err != nil {
		os.Exit(1)
	}
}
