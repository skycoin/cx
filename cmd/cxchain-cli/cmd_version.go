package main

import (
	"flag"
	"os"

	"github.com/SkycoinProject/cx/cxutil"
)

// These values should be populated by -ldflags on compilation
var (
	version = "0.0.0"
)

func cmdVersion() {
	cxutil.CmdPrintf(flag.CommandLine, "Version:\n  %s %s\n", os.Args[0], version)
}
