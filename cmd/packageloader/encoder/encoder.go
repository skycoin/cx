package encoder

import (
	"os"

	"github.com/skycoin/cx/cmd/packageloader/bolt"
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
	var listBytes []byte
	switch DATABASE {
	case "redis":
		interfaceString, err := redis.Get(packageName)
		if err != nil {
			return err
		}
		listBytes = []byte(interfaceString.(string))
	case "bolt":
		listBytes, err = bolt.Get(packageName)
		if err != nil {
			return err
		}
	}
	packageList.UnmarshalBinary(listBytes)
	for _, pack := range packageList.Packages {
		var packageStruct loader.Package
		var packageBytes []byte
		switch DATABASE {
		case "redis":
			interfaceString, err := redis.Get(pack)
			if err != nil {
				return err
			}
			packageBytes = []byte(interfaceString.(string))
		case "bolt":
			packageBytes, err = bolt.Get(pack)
			if err != nil {
				return err
			}
		}
		packageStruct.UnmarshalBinary(packageBytes)

		if packageStruct.PackageName != "main" {
			continue
		}

		path = path + "src/"
		var filePath = path

		SaveFilesToDisk(packageStruct, filePath)
	}

	for _, pack := range packageList.Packages {
		var packageStruct loader.Package
		var packageBytes []byte
		switch DATABASE {
		case "redis":
			interfaceString, err := redis.Get(pack)
			if err != nil {
				return err
			}
			packageBytes = []byte(interfaceString.(string))
		case "bolt":
			packageBytes, err = bolt.Get(pack)
			if err != nil {
				return err
			}
		}
		packageStruct.UnmarshalBinary(packageBytes)

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
		var fileBytes []byte
		switch DATABASE {
		case "redis":
			interfaceString, err := redis.Get(file)
			if err != nil {
				return err
			}
			fileBytes = []byte(interfaceString.(string))
		case "bolt":
			fileBytes, err = bolt.Get(file)
			if err != nil {
				return err
			}
		}
		fileStruct.UnmarshalBinary(fileBytes)

		err = os.WriteFile(filePath+fileStruct.FileName, fileStruct.Content, 0755)
		if err != nil {
			return err
		}
	}
	return nil
}
