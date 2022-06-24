package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/skycoin/cx/cmd/packageloader/encoder"
	"github.com/skycoin/cx/cmd/packageloader/graph"
	"github.com/skycoin/cx/cmd/packageloader/loader"
	"github.com/skycoin/cx/cmd/packageloader/tree"
	"github.com/skycoin/cx/cx/util"
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
		_, sourceCodes, _ := ParseArgsForCX([]string{path}, true)
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

func ParseArgsForCX(args []string, alsoSubdirs bool) (cxArgs []string, sourceCode []*os.File, fileNames []string) {
	skip := false // flag for skipping arg

	for _, arg := range args {

		// skip arg if skip flag is specified
		if skip {
			skip = false
			continue
		}

		// cli flags are either "--key=value" or "-key value"
		// we have to skip both cases
		if len(arg) > 1 && arg[0] == '-' {
			if !strings.Contains(arg, "=") {
				skip = true
			}
			continue
		}

		// cli cx flags are prefixed with "++"
		if len(arg) > 2 && arg[:2] == "++" {
			cxArgs = append(cxArgs, arg)
			continue
		}

		fi, err := util.CXStatFile(arg)
		if err != nil {
			println(fmt.Sprintf("%s: source file or library not found", arg))
			os.Exit(1)
		}

		switch mode := fi.Mode(); {
		case mode.IsDir():
			var fileList []string
			var err error

			// Checking if we want to check all subdirectories.
			if alsoSubdirs {
				fileList, err = filePathWalkDir(arg)
			} else {
				fileList, err = ioReadDir(arg)
				// fileList, err = filePathWalkDir(arg)
			}

			if err != nil {
				panic(err)
			}

			for _, path := range fileList {
				file, err := util.CXOpenFile(path)

				if err != nil {
					println(fmt.Sprintf("%s: source file or library not found", arg))
					os.Exit(1)
				}

				fiName := file.Name()
				fiNameLen := len(fiName)

				if fiNameLen > 2 && fiName[fiNameLen-3:] == ".cx" {
					// only loading .cx files
					sourceCode = append(sourceCode, file)
					fileNames = append(fileNames, fiName)
				}
			}
		case mode.IsRegular():
			file, err := util.CXOpenFile(arg)

			if err != nil {
				panic(err)
			}

			fileNames = append(fileNames, file.Name())
			sourceCode = append(sourceCode, file)
		}
	}

	return cxArgs, sourceCode, fileNames
}

func filePathWalkDir(root string) ([]string, error) {
	var files []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			files = append(files, path)
		}
		return err
	})
	return files, err
}

func ioReadDir(root string) ([]string, error) {
	var files []string
	fileInfo, err := ioutil.ReadDir(root)
	if err != nil {
		return files, err
	}

	for _, file := range fileInfo {
		files = append(files, fmt.Sprintf("%s/%s", root, file.Name()))
	}
	return files, nil
}
