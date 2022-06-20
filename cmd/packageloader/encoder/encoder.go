package encoder

import (
	"os"

	"github.com/skycoin/cx/cmd/packageloader/bolt"
	"github.com/skycoin/cx/cmd/packageloader/loader"
	"github.com/skycoin/cx/cmd/packageloader/redis"
)

func SavePackagesToDisk(packageName string, path string, database string) error {
	err := os.Mkdir(path, 0755)
	if err != nil {
		return err
	}

	var packageList loader.PackageList
	var listBytes []byte
	switch database {
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
		switch database {
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
		err := packageStruct.UnmarshalBinary(packageBytes)
		if err != nil {
			return err
		}

		var filePath string
		if packageStruct.PackageName == "main" {
			filePath = path + "src/"
		} else {
			filePath = path + "src/" + packageStruct.PackageName + "/"
		}

		err = SaveFilesToDisk(packageStruct, filePath, database)
		if err != nil {
			return err
		}
	}
	return nil
}

func SaveFilesToDisk(packageStruct loader.Package, filePath string, database string) error {
	err := os.Mkdir(filePath, 0755)
	if err != nil {
		return err
	}
	for _, file := range packageStruct.Files {
		var fileStruct loader.File
		var fileBytes []byte
		switch database {
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
