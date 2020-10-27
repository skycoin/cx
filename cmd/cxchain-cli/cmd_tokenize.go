package main

import (
	"flag"
	"os"

	"github.com/SkycoinProject/cx/cxgo/parser"
	"github.com/SkycoinProject/cx/cxutil"
)

func cmdTokenize(args []string) {
	cmd := flag.NewFlagSet("cxchain-cli tokenize", flag.ExitOnError)

	cmd.Usage = func() {
		usage := cxutil.DefaultUsageFormat("flags")
		usage(cmd, nil)
	}

	// flag: output, o
	out := stdoutFile
	cmd.StringVar(&out, "output", out, "`FILE` to use as compile output")
	cmd.StringVar(&out, "o", out, "shorthand for 'output'")

	// flag: input, i
	in := stdinFile
	cmd.StringVar(&in, "input", in, "`FILE` to use as compile input")
	cmd.StringVar(&in, "i", in, "shorthand for 'input'")

	// parse:
	if err := cmd.Parse(args); err != nil {
		os.Exit(1)
	}

	inF, closeIn, err := openFile(in)
	if err != nil {
		errPrintf("Failed to open input file '%s': %v\n", in, err)
		os.Exit(1)
	}
	defer closeIn()

	outF, closeOut, err := createFile(out)
	if err != nil {
		errPrintf("Failed to create output file '%s': %v\n", out, err)
		os.Exit(1)
	}
	defer closeOut()

	parser.Tokenize(inF, outF)
}
