package file_output

import (
	"bufio"
	"errors"
	"strings"

	"github.com/skycoin/cx/cx/ast"
	cxinit "github.com/skycoin/cx/cx/init"
	"github.com/skycoin/cx/cxparser/actions"

	"github.com/skycoin/cx/cmd/packageloader/bolt"
	"github.com/skycoin/cx/cmd/packageloader/loader"
	"github.com/skycoin/cx/cmd/packageloader/redis"
)

// - Adds Imports to AST
// - Returns Import Files
func GetImportFiles(packageName string, database string) (files []*loader.File, err error) {

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

		actions.AST.SelectPackage(packageStruct.PackageName)

		for _, fileString := range packageStruct.Files {

			var fileStruct loader.File
			fileBytes, err := GetStructBytes(fileString, database)
			if err != nil {
				return files, err
			}
			fileStruct.UnmarshalBinary(fileBytes)

			files = append(files, &fileStruct)

			scanner := bufio.NewScanner(strings.NewReader(string(fileStruct.Content)))
			scanner.Split(bufio.ScanWords)

			var lineno int
			wordBefore := ""
			for scanner.Scan() {

				if strings.Contains(scanner.Text(), "\n") {
					lineno++
				}
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

				actions.DeclareImport(actions.AST, importString, fileStruct.FileName, lineno)
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

func AddPkgsToAST(packageName string, database string) (err error) {
	if actions.AST == nil {
		actions.AST = cxinit.MakeProgram()
	}

	var packageList loader.PackageList
	listBytes, err := GetStructBytes(packageName, database)
	if err != nil {
		return err
	}
	packageList.UnmarshalBinary(listBytes)

	for _, packageString := range packageList.Packages {
		var packageStruct loader.Package
		packageBytes, err := GetStructBytes(packageString, database)
		if err != nil {
			return err
		}
		packageStruct.UnmarshalBinary(packageBytes)

		if pkg, err := actions.AST.GetPackage(packageStruct.PackageName); err != nil {
			pkg = ast.MakePackage(packageStruct.PackageName)
			pkgIdx := actions.AST.AddPackage(pkg)
			pkg, err = actions.AST.GetPackageFromArray(pkgIdx)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
