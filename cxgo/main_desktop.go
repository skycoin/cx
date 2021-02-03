// +build !android

package main

import (
	"os"

	. "github.com/skycoin/cx/cx"
)

func main() {
	CXLogFile(true)
	Run(os.Args[1:])
}
