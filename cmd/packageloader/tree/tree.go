package tree

import (
	"bufio"
	"errors"
	"strings"

	"github.com/skycoin/cx/cmd/packageloader/bolt"
	"github.com/skycoin/cx/cmd/packageloader/loader"
	"github.com/skycoin/cx/cmd/packageloader/redis"
)

func GetImportTree(packageName string, database string) (output string, err error) {
	var packageList loader.PackageList
	listBytes, err := GetStructBytes(packageName, database)
	if err != nil {
		return output, err
	}
	packageList.UnmarshalBinary(listBytes)

	var mainImports []string
	var alreadyPrinted []string
	var hasPackages bool
	for i, packageString := range packageList.Packages {
		var imports []string
		var packageStruct loader.Package
		packageBytes, err := GetStructBytes(packageString, database)
		if err != nil {
			return output, err
		}
		packageStruct.UnmarshalBinary(packageBytes)

		for _, fileString := range packageStruct.Files {
			var fileStruct loader.File
			fileBytes, err := GetStructBytes(fileString, database)
			if err != nil {
				return output, err
			}
			fileStruct.UnmarshalBinary(fileBytes)

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
				if !loader.Contains(imports, importString) {
					imports = append(imports, importString)
				}
				wordBefore = scanner.Text()
			}
		}

		if packageStruct.PackageName == "main" {
			mainImports = imports
		} else {
			hasPackages = true
			if i == len(packageList.Packages)-1 {
				output += "`--" + packageStruct.PackageName + "\n"
				break
			}
			output += "|--" + packageStruct.PackageName + "\n"
			alreadyPrinted = append(alreadyPrinted, packageStruct.PackageName)
			output += "|  "
			for i, importString := range imports {
				if i == len(imports)-1 {
					output += "`--" + importString + "\n"
					break
				}
				output += "|--" + importString + "\n"
			}
		}
	}
	for i, importString := range mainImports {
		if !loader.Contains(alreadyPrinted, importString) {
			if !hasPackages {
				if i == len(mainImports)-1 {
					output = "`--" + importString + "\n" + output
					break
				}
			}
			output = "|--" + importString + "\n" + output
		}
	}
	output = "main\n" + output
	return output, nil
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
