package encoder

import (
	"log"
	"os"

	"github.com/skycoin/cx/cmd/packageloader/loader"
)

func SavePackagesToDisk(packageName string, path string) {
	packageList := loader.PackageListLookup[packageName]
	for _, pack := range packageList.Packages {
		packageStruct := loader.PackageLookup[pack]
		var filePath string
		if packageStruct.PackageName == "main" {
			filePath = path + "src/"
			err := os.Mkdir(filePath, 0755)
			if err != nil {
				log.Fatal(err)
			}
			path = path + "src/"
		} else {
			filePath = path + packageStruct.PackageName + "/"
			err := os.Mkdir(filePath, 0755)
			if err != nil {
				log.Fatal(err)
			}
		}
		for _, file := range packageStruct.Files {
			fileStruct := loader.FileLookup[file]
			file, err := os.Create(path + fileStruct.FileName)
			if err != nil {
				log.Fatal(err)
			}
			defer file.Close()
			err = os.WriteFile(filePath+fileStruct.FileName, fileStruct.Content, 0644)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}
