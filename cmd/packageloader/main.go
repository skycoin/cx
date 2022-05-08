package main

import (
	"log"
	"os"

	"github.com/skycoin/cx/cmd/packageloader/loader"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("No arguments given. Usage: packageloader <path>")
	}
	if len(os.Args) > 2 {
		log.Fatal("Too many arguments given. Usage: packageloader <path>")
	}
	loader.LoadPackages(os.Args[1])
}
