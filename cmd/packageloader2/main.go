package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/skycoin/cx/cmd/packageloader2/encoder"
	"github.com/skycoin/cx/cmd/packageloader2/graph"
	"github.com/skycoin/cx/cmd/packageloader2/loader"
	"github.com/skycoin/cx/cmd/packageloader2/tree"
)

func main() {
	redisFlag := flag.Bool("redis", false, "Use Redis Key-Value database.")
	loadFlag := flag.Bool("load", false, "OPTION: Load a program to the database, with a given name to load it as and path to the program")
	saveFlag := flag.Bool("save", false, "OPTION: Save a package to disk, with a given name to search on the database and a new directory path to save to")
	treeFlag := flag.Bool("tree", false, "OPTION: Print the import dependency tree for a given program on the database")
	graphFlag := flag.Bool("graph", false, "OPTION: Print the import dependencies for each module in a program on the database")
	nameFlag := flag.String("name", "", "The name of the program to load or save")
	pathFlag := flag.String("path", "", "The path to the program to load or save")
	flag.Parse()

	if flag.NFlag()+flag.NArg() > 4 || flag.NFlag()+flag.NArg() < 2 {
		log.Fatal("Wrong number of arguments. Type -help for more information")
	}

	var database string
	if *redisFlag {
		database = "redis"
	} else {
		database = "bolt"
	}

	programName := *nameFlag

	if *treeFlag {
		output, err := tree.GetImportTree(programName, database)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Print(output)
		return
	}

	if *graphFlag {
		output, err := graph.GetImportGraph(programName, database)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Print(output)
		return
	}

	path := *pathFlag
	if path[len(path)-1:] != "/" {
		path += "/"
	}
	if path[0:2] == "./" {
		path = path[2:]
	}

	if *loadFlag {
		_, sourceCodes, _ := loader.ParseArgsForCX([]string{path}, true)
		err := loader.LoadCXProgram(programName, sourceCodes, database)
		if err != nil {
			log.Fatal(err)
		}
		return
	}
	if *saveFlag {
		err := encoder.SavePackagesToDisk(programName, path, database)
		if err != nil {
			log.Fatal(err)
		}
		return
	}

	log.Fatal("Wrong arguments provided. Type -help for more information")
}
