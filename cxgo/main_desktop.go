// +build !android

package main

import (
	"os"
)

func main() {
	Run(os.Args[1:])
}
