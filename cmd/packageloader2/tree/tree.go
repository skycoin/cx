package tree

import (
	"bufio"
	"errors"
	"strings"

	"github.com/skycoin/cx/cmd/packageloader2/bolt"
	"github.com/skycoin/cx/cmd/packageloader2/loader"
	"github.com/skycoin/cx/cmd/packageloader2/redis"
)

func GetImportTree(packageName string, database string) (output string, err error) {
	var packageList loader.PackageList
	listBytes, err := GetStructBytes(packageName, database)
	if err != nil {
		return output, err
	}
	packageList.UnmarshalBinary(listBytes)

	var importMap = make(map[string][]string)
	for _, packageString := range packageList.Packages {
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
		importMap[packageStruct.PackageName] = imports
	}
	output += "main\n"
	var layers = make(map[int]bool)
	for i, mainImport := range importMap["main"] {
		if i == len(importMap["main"])-1 {
			output += "`--" + mainImport + "\n"
			layers[0] = false
		} else {
			output += "|--" + mainImport + "\n"
			layers[0] = true
		}
		output += AddDependenciesRecur(importMap, importMap[mainImport], 1, layers)
	}
	return output, nil
}

func AddDependenciesRecur(importMap map[string][]string, dependencies []string, depth int, layers map[int]bool) (output string) {
	if len(dependencies) == 0 {
		return output
	}
	for i, dependencyImport := range dependencies {
		for l := 0; l < depth; l++ {
			if layers[l] {
				output += "|  "
			} else {
				output += "   "
			}
		}
		if i == len(dependencies)-1 {
			output += "`--" + dependencyImport + "\n"
			layers[depth] = false
		} else {
			output += "|--" + dependencyImport + "\n"
			layers[depth] = true
		}
		output += AddDependenciesRecur(importMap, importMap[dependencyImport], depth+1, layers)
	}
	return output
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
