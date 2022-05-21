package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/skycoin/cx/cmd/packageloader/encoder"
	"github.com/skycoin/cx/cmd/packageloader/loader"
)

func main() {
	loadFlag := flag.Bool("load", false, "OPTION: Load a program to the database, with a given name to load it as and path to the program")
	saveFlag := flag.Bool("save", false, "OPTION: Save a package to disk, with a given name to search on the database and a new directory path to save to")
	helpFlag := flag.Bool("help", false, "OPTION: Display this help message")
	nameFlag := flag.String("name", "", "The name of the program to load or save")
	pathFlag := flag.String("path", "", "The path to the program to load or save")
	flag.Parse()
	if *helpFlag {
		fmt.Println("Syntax: packageloader [OPTION] -path [PATH] -name [NAME]")
		flag.Usage()
		os.Exit(0)
	}
	if flag.NFlag()+flag.NArg() > 3 || flag.NFlag()+flag.NArg() < 3 {
		log.Fatal("Wrong number of arguments. Type -help for more information")
	}

	programName := *nameFlag
	path := *pathFlag
	if path[len(path)-1:] != "/" {
		path += "/"
	}
	if path[0:2] == "./" {
		path = path[2:]
	}

	if *loadFlag {
		loader.LoadPackages(programName, path)
		return
	}
	if *saveFlag {
		encoder.SavePackagesToDisk(programName, path)
		return
	}
	log.Fatal("Wrong arguments provided. Type -help for more information")
}
