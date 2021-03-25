// +build os,mobile

package cxcore

import (
	"fmt"
	"io"
	"os"

	"golang.org/x/mobile/asset"
)

func CopyAssetsToFilesDir() bool {
	filesDir := asset.GetFilesDir()
	return CopyAssetsTo("", filesDir)
}

func CopyAssetTo(path string, dest string) bool {
	destPath := fmt.Sprintf("%s/%s", dest, path)
	fmt.Printf("Copy '%s' to '%s'\n", path, destPath)
	source, err := asset.Open(path)
	if err != nil {
		fmt.Printf("Failed to open source asset '%s', err '%v'\n", path, err)
		return false
	}
	defer source.Close()

	destination, err := os.Create(destPath)
	if err != nil {
		fmt.Printf("Failed to create destination file '%s', err '%v'\n", destPath, err)
		return false
	}
	defer destination.Close()

	if _, err := io.Copy(destination, source); err != nil {
		fmt.Printf("Failed to copy data from asset '%s' to '%s', err '%v'\n", path, destPath, err)
		return false
	}

	return true
}

func CopyAssetsTo(path string, dest string) bool {
	//fmt.Printf("COPY_ASSETS '%s' TO '%s'\n", path, dest)
	assetList := asset.GetAssetList(path)
	if len(assetList) == 0 {
		if success := CopyAssetTo(path, dest); success == false {
			fmt.Printf("ProgramError copying file '%s' to '%s'\n", path, dest)
			return false
		}
	} else {
		destPath := fmt.Sprintf("%s/%s", dest, path)
		//fmt.Printf("LOOKING FOR DIR '%s'\n", destPath)
		_, err := os.Stat(destPath)
		if os.IsNotExist(err) {
			//fmt.Printf("MKDIR '%s', '%v'\n", destPath, err)
			if err := os.Mkdir(destPath, 0766); err != nil {
				fmt.Printf("ProgramError creating dir '%s', err '%v'\n", destPath, err)
				return false
			}
		} else {
			//fmt.Printf("DEST DIR FOUND '%s', '%v'\n", destPath, err)
		}
		for _, s := range assetList {
			sourcePath := s
			//fmt.Printf("SOURCE_PATH '%s'\n", sourcePath)
			if path != "" {
				sourcePath = fmt.Sprintf("%s/%s", path, s)
				//fmt.Printf("SOURCE_PATH_EMPTY '%s' ADD '%s'\n", sourcePath, s)
			}
			if success := CopyAssetsTo(sourcePath, dest); success == false {
				fmt.Printf("ProgramError copying file or dir '%s' to '%s'\n", sourcePath, dest)
				return false
			}
		}
	}
	return true
}
