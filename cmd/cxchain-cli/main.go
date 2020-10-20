package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/SkycoinProject/cx-chains/src/util/logging"

	"github.com/SkycoinProject/cx/cxutil"
)

var log = logging.MustGetLogger("cxchain-cli")

// These values should be populated by -ldflags on compilation
var (
	version = "0.0.0"
)

var cm = cxutil.NewCommandMap(flag.CommandLine, 6, cxutil.DefaultUsageSignature()).
	Add("version", func(_ []string) { cmdVersion() }).
	Add("help", func(_ []string) { flag.CommandLine.Usage() }).
	Add("tokenize", cmdTokenize).
	Add("new", cmdNew).
	Add("run", cmdRun).
	Add("state", cmdState).
	Add("peers", cmdPeers)

func cmdVersion() {
	_, _ = fmt.Fprintf(flag.CommandLine.Output(), "%s %s\n", os.Args[0], version)
}

func main() {
	os.Exit(cm.ParseAndRun(os.Args[1:]))
}

func parseFlagSet(cmd *flag.FlagSet, args []string) {
	if err := cmd.Parse(args); err != nil {
		os.Exit(1)
	}
}
