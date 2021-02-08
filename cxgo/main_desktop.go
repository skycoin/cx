// +build !android

package main

import (
	"os"

	. "github.com/skycoin/cx/cx"
)

func main() {
	CXLogFile(true)
	if os.Args != nil && len(os.Args) > 1 {
		Run(os.Args[1:])
	}

}
