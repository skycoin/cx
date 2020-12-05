package main

import (
	"flag"
	"os"

	"github.com/SkycoinProject/cx-chains/src/util/logging"

	"github.com/SkycoinProject/cx/cxutil"
)

var log = logging.MustGetLogger("cxchain-cli")

func main() {
	initGlobals()

	usageMenu := cxutil.UsageFormat(func(cmd *flag.FlagSet, subcommands []string) {
		// print: Usage
		cxutil.PrintCmdUsage(cmd, "Usage", subcommands, []string{"args"})

		// print: ENVs
		cxutil.CmdPrintf(cmd, "ENVs:")
		cxutil.CmdPrintf(cmd, "  $%s\n  \t%s", specFileEnv, "chain spec filepath")
		cxutil.CmdPrintf(cmd, "  $%s\n  \t%s", genSKEnv, "genesis secret key (hex)")

		// print: Flags
		cxutil.PrintCmdFlags(cmd, "Flags")
	})

	root := cxutil.NewCommandMap(flag.CommandLine, 7, usageMenu).
		AddSubcommand("version", func(_ []string) { cmdVersion() }).
		AddSubcommand("help", func(_ []string) { flag.CommandLine.Usage() }).
		AddSubcommand("tokenize", cmdTokenize).
		AddSubcommand("new", cmdNew).
		AddSubcommand("run", cmdRun).
		AddSubcommand("state", cmdState).
		AddSubcommand("peers", cmdPeers)

	os.Exit(root.ParseAndRun(os.Args[1:]))
}
