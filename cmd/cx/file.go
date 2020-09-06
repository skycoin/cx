package main

import (
	"fmt"
	"os"
	"strings"

	cxcore "github.com/SkycoinProject/cx/cx"
)

const (
	stdinFile  = "STDIN"
	stdoutFile = "STDOUT"
)

func openFile(filename string) (*os.File, func(), error) {
	switch filename {
	case stdinFile, "":
		return os.Stdin, func(){}, nil
	case stdoutFile:
		return nil, nil, fmt.Errorf("cannot read from %s", filename)
	default:
		f, err := cxcore.CXOpenFile(filename)
		if err != nil {
			return nil, nil, err
		}
		c := func() {
			if err := f.Close(); err != nil {
				errPrintf("Failed to close file '%s': %v\n", filename, err)
			}
		}
		return f, c, nil
	}
}

func createFile(filename string) (*os.File, func(), error) {
	switch filename {
	case stdoutFile, "":
		return os.Stdout, func(){}, nil
	case stdinFile:
		return nil, nil, fmt.Errorf("cannot write to %s", filename)
	default:
		f, err := cxcore.CXCreateFile(filename)
		if err != nil {
			return nil, nil, err
		}
		c := func() {
			if err := f.Close(); err != nil {
				errPrintf("Failed to close file '%s': %v\n", filename, err)
			}
		}
		return f, c, nil
	}
}

func errPrintf(format string, a ...interface{}) {
	if _, err := fmt.Fprintf(os.Stderr, format, a...); err != nil {
		panic(err)
	}
}

func stripFlags(args []string) []string {
	out := make([]string, 0, len(args))
	skip := false

	for _, arg := range args {

		if skip {
			skip = false
			continue
		}

		if strings.HasPrefix(arg, "-") {
			if !strings.Contains(arg, "=") {
				skip = true
			}
			continue
		}

		out = append(out, arg)
	}

	return out
}
