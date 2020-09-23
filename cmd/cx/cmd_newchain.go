package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"os"
	"strings"

	"github.com/SkycoinProject/cx-chains/src/cipher"

	"github.com/SkycoinProject/cx/cxgo/cxspec"
)

const filePerm = 0644

func cmdNewChain(args []string) {
	cmd := flag.NewFlagSet(args[0], flag.ExitOnError)

	replace := false
	cmd.BoolVar(&replace, "replace", replace, "whether to replace output file(s)")
	cmd.BoolVar(&replace, "r", replace, "shorthand for 'replace'")

	unifyKeys := false
	cmd.BoolVar(&unifyKeys, "unify", unifyKeys, "whether to use the same keys for genesis and chain")
	cmd.BoolVar(&unifyKeys, "u", unifyKeys, "shorthand for 'unify'")

	coinName := "skycoin"
	cmd.StringVar(&coinName, "coin", coinName, "`NAME` for cx coin")
	cmd.StringVar(&coinName, "c", coinName, "shorthand for 'coin'")

	coinTicker := "SKY"
	cmd.StringVar(&coinTicker, "ticker", coinTicker, "`SYMBOL` for cx coin ticker")
	cmd.StringVar(&coinTicker, "t", coinTicker, "shorthand for 'ticker'")

	program := "STDIN"
	cmd.StringVar(&program, "program", program, "`FILE` containing genesis program source")
	cmd.StringVar(&program, "p", program, "shorthand for 'program'")

	chainSpecOut := "./{coin}.chain_spec.json"
	cmd.StringVar(&chainSpecOut, "chain-spec-output", chainSpecOut, "`FILE` for chain spec output")

	chainKeysOut := "./{coin}.chain_keys.json"
	cmd.StringVar(&chainKeysOut, "chain-keys-output", chainKeysOut, "`FILE` for chain keys output")

	genKeysOut := "./{coin}.genesis_keys.json"
	cmd.StringVar(&genKeysOut, "genesis-keys-output", genKeysOut, "`FILE` for genesis keys output")

	// Parse flags.
	parseFlagSet(cmd, args[1:])
	fillTemplate(&chainSpecOut, coinName, coinTicker)
	fillTemplate(&chainKeysOut, coinName, coinTicker)
	fillTemplate(&genKeysOut, coinName, coinTicker)

	// Check replace.
	if !replace {
		for _, name := range []string{chainSpecOut, chainKeysOut, genKeysOut} {
			if _, err := os.Stat(name); !os.IsNotExist(err) {
				errPrintf("File '%s' already exists. Replace with '--replace' flag.\n", name)
				os.Exit(1)
			}
		}
	}

	// Generate chain keys.
	chainPK, chainSK := cipher.GenerateKeyPair()

	// Generate genesis keys.
	genPK, genSK := chainPK, chainSK
	if !unifyKeys {
		genPK, genSK = cipher.GenerateKeyPair()
	}
	genAddr := cipher.AddressFromPubKey(genPK)

	// Generate and write chain spec file.
	cSpec, err := cxspec.New(coinName, coinTicker, chainSK, genAddr, nil)
	if err != nil {
		errPrintf("Failed to generate chain spec: %v\n", err)
		os.Exit(1)
	}
	cSpecB, err := json.MarshalIndent(cSpec, "", "\t")
	if err != nil {
		errPrintf("Failed to encode chain spec to json: %v\n", err)
		os.Exit(1)
	}
	if err := ioutil.WriteFile(chainSpecOut, cSpecB, filePerm); err != nil {
		errPrintf("Failed to write chain spec to file '%s': %v\n", chainSpecOut, err)
		os.Exit(1)
	}

	// Write chain keys file.
	cKeys := cxspec.KeySpecFromSecKey(cxspec.ChainKey, chainSK, true, true)
	cKeysB, err := json.MarshalIndent(cKeys, "", "\t")
	if err != nil {
		errPrintf("Failed to encode chain keys to json: %v\n", err)
		os.Exit(1)
	}
	if err := ioutil.WriteFile(chainKeysOut, cKeysB, filePerm); err != nil {
		errPrintf("Failed to write chain keys to file '%s': %v\n", chainKeysOut, err)
		os.Exit(1)
	}

	// Write genesis keys to file.
	gKeys := cxspec.KeySpecFromSecKey(cxspec.GenesisKey, genSK, true, true)
	gKeysB, err := json.MarshalIndent(gKeys, "", "\t")
	if err != nil {
		errPrintf("Failed to encode genesis keys to json: %v\n", err)
		os.Exit(1)
	}
	if err := ioutil.WriteFile(genKeysOut, gKeysB, filePerm); err != nil {
		errPrintf("Failed to write genesis keys to file '%s': %v\n", genKeysOut, err)
		os.Exit(1)
	}
}

func fillTemplate(s *string, coinName, coinTicker string) {
	*s = strings.ReplaceAll(*s, "{coin}", coinName)
	*s = strings.ReplaceAll(*s, "{ticker}", coinTicker)
}
