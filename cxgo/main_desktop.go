// +build !android

package main

import (
	. "github.com/SkycoinProject/cx/cx"
	"os"
)

func main() {
	CXLogFile(true)
	Run(os.Args[1:])
}
