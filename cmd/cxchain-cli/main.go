package main

import (
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/SkycoinProject/cx-chains/src/util/logging"
)

var log = logging.MustGetLogger("cxchain-cli")

// These values should be populated by -ldflags on compilation
var (
	version = "0.0.0"
)

const cmdCount = 5

var (
	cmdMap  = make(map[string]func(), cmdCount)
	cmdList = make([]string, 0, cmdCount)
)

func init() {
	add := func(name string, fn func()) {
		cmdMap[name] = fn
		cmdList = append(cmdList, name)
	}
	add("version", cmdVersion)
	add("help", cmdHelp)
	add("tokenize", func() { cmdTokenize(os.Args[1:]) })
	add("new", func() { cmdNew(os.Args[1:]) })
	add("run", func() { cmdRun(os.Args[1:]) })
	add("state", func() { cmdState(os.Args[1:]) })

	sort.Strings(cmdList)
}

func cmdVersion() {
	_, _ = fmt.Fprintf(os.Stderr, "%s %s\n", os.Args[0], version)
}

func cmdHelp() {
	_, _ = fmt.Fprintf(os.Stderr, "usage: %s [%s] [args...]\n",
		os.Args[0],
		strings.Join(cmdList, "|"))
	flag.CommandLine.PrintDefaults()
}

func main() {
	flag.CommandLine.Usage = cmdHelp
	flag.Parse()

	subCmd, ok := cmdMap[os.Args[1]]
	if !ok {
		_, _ = fmt.Fprintf(os.Stderr, "subcommand '%s' does not exist\n", os.Args[1])
		flag.CommandLine.Usage()
		return
	}
	subCmd()
}

func parseFlagSet(cmd *flag.FlagSet, args []string) {
	if err := cmd.Parse(args); err != nil {
		os.Exit(1)
	}
}
