package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/SkycoinProject/cx-chains/src/cipher"

	"github.com/SkycoinProject/cx/cxgo/cxspec"
)

const (
	// ENV for chain spec filepath.
	specFileEnv          = "CXCHAIN_SPEC"
	defaultChainSpecFile = "./skycoin.chain_spec.json"

	// ENV for genesis secret key.
	genSKEnv = "CXCHAIN_GEN_SK"
)

var ErrNoSKProvided = errors.New("no secret key provided")

var globals = struct {
	spec    cxspec.ChainSpec // chain spec struct obtained from file defined in `CXCHAIN_SPEC_FILE` ENV
	specFilename string // chain spec filename
	specErr error            // error when obtaining chain spec

	genSK    cipher.SecKey // genesis secret key obtained from SK defined in `CXCHAIN_GENESIS_SK` ENV
	genSKErr error         // error when obtaining genesis secret key
}{}

// initGlobals initiates globals. This is called in main.
func initGlobals() {
	globals.spec, globals.specFilename, globals.specErr = parseSpecFileEnv()
	globals.genSK, globals.genSKErr = parseGenesisSKEnv()

	log := log.
		WithField(specFileEnv, nil).
		WithField(genSKEnv, nil)

	if globals.specErr == nil {
		log = log.WithField(specFileEnv, globals.specFilename)
	}
	if globals.genSKErr == nil {
		log = log.WithField(genSKEnv, globals.genSK.Hex())
	}

	log.Info("Environment variables:")
}

// parseSpecFileEnv parses chain spec filename from CXCHAIN_SPEC_FILEPATH env.
func parseSpecFileEnv() (cxspec.ChainSpec, string, error) {
	specFilename, ok := os.LookupEnv(specFileEnv)
	if !ok {
		specFilename = defaultChainSpecFile
	}


	spec, err := cxspec.ReadSpecFile(specFilename)
	if err != nil {
		return cxspec.ChainSpec{}, specFilename, fmt.Errorf("failed to read chain spec from %s: %w", specFilename, err)
	}

	// log.WithField("filename", specFilename).Info("Using chain spec.")
	cxspec.PopulateParamsModule(spec)
	return spec, specFilename, nil
}

// parseGenesisSKEnv parses secret key from CXCHAIN_SECRET_KEY env.
// The secret key can be null.
func parseGenesisSKEnv() (cipher.SecKey, error) {
	if skStr, ok := os.LookupEnv(genSKEnv); ok {
		sk, err := cipher.SecKeyFromHex(skStr)
		if err != nil {
			return cipher.SecKey{}, fmt.Errorf("failed to parse secret key defined in ENV '%s': %w", genSKEnv, err)
		}
		return sk, nil
	}
	return cipher.SecKey{}, ErrNoSKProvided
}