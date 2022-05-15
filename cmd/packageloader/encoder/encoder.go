package encoder

import (
	"log"
	"os"

	"github.com/skycoin/cx/cmd/packageloader/loader"
	"github.com/skycoin/cx/cmd/packageloader/server"
)

func SavePackagesToDisk(packageName string, path string) {
	err := os.Mkdir(path, 0755)
	if err != nil {
		log.Fatal(err)
	}
	var packageList loader.PackageList
	packageList.UnmarshalBinary([]byte(server.Get(packageName).(string)))

	for _, pack := range packageList.Packages {
		var packageStruct loader.Package
		packageStruct.UnmarshalBinary([]byte(server.Get(pack).(string)))

		var filePath string
		if packageStruct.PackageName != "main" {
			continue
		}

		path = path + "src/"
		filePath = path

		err := os.Mkdir(filePath, 0755)
		if err != nil {
			log.Fatal(err)
		}

		for _, file := range packageStruct.Files {
			var fileStruct loader.File
			fileStruct.UnmarshalBinary([]byte(server.Get(file).(string)))

			err = os.WriteFile(filePath+fileStruct.FileName, fileStruct.Content, 0755)
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	for _, pack := range packageList.Packages {
		var packageStruct loader.Package
		packageStruct.UnmarshalBinary([]byte(server.Get(pack).(string)))

		var filePath string
		if packageStruct.PackageName == "main" {
			continue
		}
		filePath = path + packageStruct.PackageName + "/"

		err := os.Mkdir(filePath, 0755)
		if err != nil {
			log.Fatal(err)
		}

		for _, file := range packageStruct.Files {
			var fileStruct loader.File
			fileStruct.UnmarshalBinary([]byte(server.Get(file).(string)))
			log.Println(fileStruct.FileName)

			err = os.WriteFile(filePath+fileStruct.FileName, fileStruct.Content, 0755)
			if err != nil {
				log.Fatal(err)
			}
		}
	}

}
