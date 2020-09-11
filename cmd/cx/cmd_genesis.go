package main

import (
	"flag"

	"github.com/SkycoinProject/cx-chains/src/cipher"
	"github.com/SkycoinProject/cx-chains/src/coin"
)

func cmdGenesis(args []string) {
	cmd := flag.NewFlagSet(args[0], flag.ExitOnError)

	parseFlagSet(cmd, args)
}



func GenerateGenesisSpec(addr cipher.Address, coins uint64, timestamp uint64, progState []byte) ([]byte, error) {
	block, err := coin.NewGenesisBlock(addr, coins, timestamp, progState)
	if err != nil {
		return nil, err
	}
	block.Body.Hash()

}