package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

func main() {
	if len(os.Args) < 2 || strings.HasPrefix(os.Args[1], "-") {
		fmt.Println("A subcommand was expected.")
		cmdHelp()
		os.Exit(1)
	}

	switch c := os.Args[1]; c {
	case "version", "v":
		cmdVersion()

	case "help", "h":
		cmdHelp()

	case "tokenize", "t":
		cmdTokenize(os.Args[1:])

	case "newchain", "n":
		cmdNewChain(os.Args[1:])

	default:
		fmt.Printf("Subcommand '%s' is not found.\n", c)
		cmdHelp()
		os.Exit(1)
	}
}

func parseFlagSet(cmd *flag.FlagSet, args []string) {
	if err := cmd.Parse(args); err != nil {
		fmt.Println(err.Error())
		fmt.Println()
		fmt.Println("TODO: Help menu goes here.")

		if err == flag.ErrHelp {
			os.Exit(0)
		}

		os.Exit(1)
	}
}
