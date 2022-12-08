package file_output

import (
	"errors"

	"github.com/skycoin/cx/cmd/packageloader/bolt"
	"github.com/skycoin/cx/cmd/packageloader/loader"
	"github.com/skycoin/cx/cmd/packageloader/redis"
)

// - Adds Imports to AST
// - Returns Import Files
// - Packages must be added to AST first or an error will occur
// - Call AddPkgsToAST before calling this function
func GetImportFiles(packageName string, database string) (files []*loader.File, err error) {

	// Get package list
	var packageList loader.PackageList
	listBytes, err := GetStructBytes(packageName, database)
	if err != nil {
		return files, err
	}
	packageList.UnmarshalBinary(listBytes)

	for _, packageString := range packageList.Packages {

		//  Get package struct
		var packageStruct loader.Package
		packageBytes, err := GetStructBytes(packageString, database)
		if err != nil {
			return files, err
		}
		packageStruct.UnmarshalBinary(packageBytes)

		for _, fileString := range packageStruct.Files {

			// Get file struct
			var fileStruct loader.File
			fileBytes, err := GetStructBytes(fileString, database)
			if err != nil {
				return files, err
			}
			fileStruct.UnmarshalBinary(fileBytes)

			// Add file struct to array
			files = append(files, &fileStruct)
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
