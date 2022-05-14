package main

import (
	"log"
	"os"

	"github.com/skycoin/cx/cmd/packageloader/encoder"
	"github.com/skycoin/cx/cmd/packageloader/loader"
)

func main() {
	if len(os.Args) < 4 {
		log.Fatal("Not enough arguments given. Usage: packageloader <option> <program name> <path>")
	}
	if len(os.Args) > 4 {
		log.Fatal("Too many arguments given. Usage: packageloader <option> <program name> <path>")
	}
	programName := os.Args[2]
	var path string
	if os.Args[3][:len(os.Args[3])-1] != "/" {
		path = os.Args[3] + "/"
	} else {
		path = os.Args[3]
	}
	if os.Args[1] == "-l" || os.Args[1] == "-load" {
		loader.LoadPackages(programName, path)
	}
	if os.Args[1] == "-s" || os.Args[1] == "-save" {
		encoder.SavePackagesToDisk(programName, path)
	}
}
