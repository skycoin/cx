package encoder

import (
	"os"

	"github.com/skycoin/cx/cmd/packageloader/loader"
	"github.com/skycoin/cx/cmd/packageloader/redis"
)

var DATABASE = "bolt"

func SavePackagesToDisk(packageName string, path string) error {
	err := os.Mkdir(path, 0755)
	if err != nil {
		return err
	}
	var packageList loader.PackageList
	listString, err := redis.Get(packageName)
	if err != nil {
		return err
	}
	packageList.UnmarshalBinary([]byte(listString.(string)))

	for _, pack := range packageList.Packages {
		var packageStruct loader.Package
		packageString, err := redis.Get(pack)
		if err != nil {
			return err
		}
		packageStruct.UnmarshalBinary([]byte(packageString.(string)))

		if packageStruct.PackageName != "main" {
			continue
		}

		path = path + "src/"
		var filePath = path

		SaveFilesToDisk(packageStruct, filePath)
	}

	for _, pack := range packageList.Packages {
		var packageStruct loader.Package
		packageString, err := redis.Get(pack)
		if err != nil {
			return err
		}
		packageStruct.UnmarshalBinary([]byte(packageString.(string)))

		if packageStruct.PackageName == "main" {
			continue
		}
		var filePath = path + packageStruct.PackageName + "/"
		SaveFilesToDisk(packageStruct, filePath)
	}
	return nil
}

func SaveFilesToDisk(packageStruct loader.Package, filePath string) error {
	err := os.Mkdir(filePath, 0755)
	if err != nil {
		return err
	}
	for _, file := range packageStruct.Files {
		var fileStruct loader.File
		fileString, err := redis.Get(file)
		if err != nil {
			return err
		}
		fileStruct.UnmarshalBinary([]byte(fileString.(string)))

		err = os.WriteFile(filePath+fileStruct.FileName, fileStruct.Content, 0755)
		if err != nil {
			return err
		}
	}
	return nil
}
