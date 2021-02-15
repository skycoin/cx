// +build !android

package main

import (
	"os"

	cx "github.com/skycoin/cx/cx"
)

func main() {
	cx.CXLogFile(true)
	if os.Args != nil && len(os.Args) > 1 {
		Run(os.Args[1:])
	}

}
