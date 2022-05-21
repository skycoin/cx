package encoder

import (
	"log"
	"os"

	"github.com/skycoin/cx/cmd/packageloader/database"
	"github.com/skycoin/cx/cmd/packageloader/loader"
)

func SavePackagesToDisk(packageName string, path string) {
	err := os.Mkdir(path, 0755)
	if err != nil {
		log.Fatal(err)
	}
	var packageList loader.PackageList
	packageList.UnmarshalBinary([]byte(database.Get(packageName).(string)))

	for _, pack := range packageList.Packages {
		var packageStruct loader.Package
		packageStruct.UnmarshalBinary([]byte(database.Get(pack).(string)))

		if packageStruct.PackageName != "main" {
			continue
		}

		path = path + "src/"
		var filePath = path

		SaveFilesToDisk(packageStruct, filePath)
	}

	for _, pack := range packageList.Packages {
		var packageStruct loader.Package
		packageStruct.UnmarshalBinary([]byte(database.Get(pack).(string)))

		if packageStruct.PackageName == "main" {
			continue
		}
		var filePath = path + packageStruct.PackageName + "/"
		SaveFilesToDisk(packageStruct, filePath)
	}

}

func SaveFilesToDisk(packageStruct loader.Package, filePath string) {
	err := os.Mkdir(filePath, 0755)
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range packageStruct.Files {
		var fileStruct loader.File
		fileStruct.UnmarshalBinary([]byte(database.Get(file).(string)))

		err = os.WriteFile(filePath+fileStruct.FileName, fileStruct.Content, 0755)
		if err != nil {
			log.Fatal(err)
		}
	}
}
