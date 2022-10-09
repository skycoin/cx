package file_output

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	cxinit "github.com/skycoin/cx/cx/init"
	"github.com/skycoin/cx/cxparser/actions"

	"github.com/skycoin/cx/cmd/packageloader/bolt"
	"github.com/skycoin/cx/cmd/packageloader/loader"
	"github.com/skycoin/cx/cmd/packageloader/redis"
)

// - Adds Imports to AST
// - Returns Import Files
func GetImportFiles(packageName string, database string) (files []*os.File, err error) {

	if actions.AST == nil {
		actions.AST = cxinit.MakeProgram()
	}

	var packageList loader.PackageList
	listBytes, err := GetStructBytes(packageName, database)
	if err != nil {
		return files, err
	}
	packageList.UnmarshalBinary(listBytes)

	for _, packageString := range packageList.Packages {
		var packageStruct loader.Package
		packageBytes, err := GetStructBytes(packageString, database)
		if err != nil {
			return files, err
		}
		packageStruct.UnmarshalBinary(packageBytes)
		if err != nil {
			return files, err
		}

		for _, fileString := range packageStruct.Files {

			var fileStruct loader.File
			fileBytes, err := GetStructBytes(fileString, database)
			if err != nil {
				return files, err
			}
			fileStruct.UnmarshalBinary(fileBytes)

			fmt.Print(fileStruct.FileName)

			file, err := os.Open(fileStruct.FileName)
			if err != nil {
				return files, err
			}

			files = append(files, file)

			scanner := bufio.NewScanner(strings.NewReader(string(fileStruct.Content)))
			scanner.Split(bufio.ScanWords)

			wordBefore := ""
			for scanner.Scan() {

				if scanner.Text() != "import" {
					wordBefore = scanner.Text()
					continue
				}
				if wordBefore == "//" {
					wordBefore = scanner.Text()
					continue
				}
				if scanner.Text() == "var" || scanner.Text() == "const" || scanner.Text() == "type" || scanner.Text() == "func" {
					break
				}
				scanner.Scan()
				importString := scanner.Text()[1 : len(scanner.Text())-1]
				actions.DeclareImport(actions.AST, importString, fileString, 0)
				wordBefore = scanner.Text()
			}

		}
	}

	return files, nil
}

func GetStructBytes(structName string, database string) ([]byte, error) {
	switch database {
	case "redis":
		interfaceString, err := redis.Get(structName)
		if err != nil {
			return []byte{}, err
		}
		return []byte(interfaceString.(string)), nil
	case "bolt":
		listBytes, err := bolt.Get(structName)
		if err != nil {
			return []byte{}, err
		}
		return listBytes, nil
	}
	return []byte{}, errors.New("invalid database")
}
