package encoder

import (
	"log"
	"os"

	"github.com/skycoin/cx/cmd/packageloader/loader"
	"github.com/skycoin/cx/cmd/packageloader/server"
)

func SavePackagesToDisk(packageName string, path string) {
	var packageList loader.PackageList
	packageList.UnmarshalBinary([]byte(server.Get(packageName).(string)))

	for _, pack := range packageList.Packages {
		var packageStruct loader.Package
		packageStruct.UnmarshalBinary([]byte(server.Get(pack).(string)))

		var filePath string
		if packageStruct.PackageName == "main" {
			filePath = path + "src/"
			path = path + "src/"
		} else {
			filePath = path + packageStruct.PackageName + "/"
		}

		err := os.Mkdir(filePath, 0755)
		if err != nil {
			log.Fatal(err)
		}

		for _, file := range packageStruct.Files {
			fileStruct := server.Get(file).(loader.File)
			file, err := os.Create(path + fileStruct.FileName)
			if err != nil {
				log.Fatal(err)
			}
			defer file.Close()
			err = os.WriteFile(filePath+fileStruct.FileName, fileStruct.Content, 0755)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}
