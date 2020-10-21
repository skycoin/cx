package main

import (
	"flag"
	"os"

	"github.com/SkycoinProject/cx-chains/src/util/logging"

	"github.com/SkycoinProject/cx/cxutil"
)

var log = logging.MustGetLogger("cxchain-cli")

// These values should be populated by -ldflags on compilation
var (
	version = "0.0.0"
)

var cm = cxutil.NewCommandMap(flag.CommandLine, 6, cxutil.DefaultUsageFormat("args")).
	AddSubcommand("version", func(_ []string) { cmdVersion() }).
	AddSubcommand("help", func(_ []string) { flag.CommandLine.Usage() }).
	AddSubcommand("tokenize", cmdTokenize).
	AddSubcommand("new", cmdNew).
	AddSubcommand("run", cmdRun).
	AddSubcommand("state", cmdState).
	AddSubcommand("peers", cmdPeers)

func cmdVersion() {
	cxutil.CmdPrintf(flag.CommandLine, "Version:\n  %s %s\n", os.Args[0], version)
}

func main() {
	os.Exit(cm.ParseAndRun(os.Args[1:]))
}
