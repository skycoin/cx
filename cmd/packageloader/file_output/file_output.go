package file_output

import (
	"errors"
	"os"

	"github.com/skycoin/cx/cmd/packageloader/bolt"
	"github.com/skycoin/cx/cmd/packageloader/loader"
	"github.com/skycoin/cx/cmd/packageloader/redis"
)

func GetImportFiles(packageName string, database string) (files []*os.File, err error) {
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
