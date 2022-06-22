package graph

import (
	"bufio"
	"errors"
	"strconv"
	"strings"

	"github.com/skycoin/cx/cmd/packageloader/bolt"
	"github.com/skycoin/cx/cmd/packageloader/loader"
	"github.com/skycoin/cx/cmd/packageloader/redis"
)

func GetImportGraph(packageName string, database string) (output string, err error) {
	var packageList loader.PackageList
	listBytes, err := GetStructBytes(packageName, database)
	if err != nil {
		return output, err
	}
	packageList.UnmarshalBinary(listBytes)

	var importMap = make(map[string][]string)
	var indexList = make([]string, 0)
	var importIndexMap = make(map[string]int)
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
		indexList = append(indexList, packageStruct.PackageName)
		importIndexMap[packageStruct.PackageName] = len(indexList) - 1
		importMap[packageStruct.PackageName] = imports
	}
	for i, module := range indexList {
		output += "id=" + strconv.Itoa(i) + ", "
		output += "module=" + module + ", "
		var tmp = "imports="
		imports := importMap[module]
		for _, importString := range imports {
			importIndex := importIndexMap[importString]
			if importIndex == 0 {
				continue
			}
			tmp += strconv.Itoa(importIndex) + ","
		}
		output += tmp + "\n"
	}
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
