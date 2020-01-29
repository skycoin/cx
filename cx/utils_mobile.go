// +build base,android

package cxcore

import "golang.org/x/mobile/asset"
import "fmt"
import "os"
import "io"

/*type CXFile struct {
	name string
	impl asset.File
}

func OpenFile(path string) (CXFile, error) {
	impl, err := asset.Open(path)
	return CXFile{impl: impl, name: path}, err
}

func (file *CXFile) Close() error {
	return file.impl.Close()
}

func (file *CXFile) Impl() asset.File {
	return file.impl
}

func (file *CXFile) Name() string {
	return file.name
}*/

func CopyAssetsToFilesDir() bool {
	filesDir := asset.GetFilesDir()
	return CopyAssetsTo("", filesDir)
}

func CopyAssetTo(path string, dest string) bool {
	destPath := fmt.Sprintf("%s/%s", dest, path)
	fmt.Printf("Copy '%s' to '%s'\n", path, destPath)
	source, err := asset.Open(path)
	if err != nil {
		fmt.Printf("Failed to open source asset '%s'\n", path)
		return false
	}
	defer source.Close()

	destination, err := os.Create(destPath)
	if err != nil {
		fmt.Printf("Failed to create destination file '%s'\n", destPath)
		return false
	}
	defer destination.Close()

	if _, err := io.Copy(destination, source); err != nil {
		fmt.Printf("Failed to copy data from asset '%s' to '%s'\n", path, destPath)
		return false
	}

	return true
}

func CopyAssetsTo(path string, dest string) bool {
	//fmt.Printf("COPY_ASSETS '%s' TO '%s'\n", path, dest)
	assetList := asset.GetAssetList(path)
	if len(assetList) == 0 {
		if success := CopyAssetTo(path, dest); success == false {
			fmt.Printf("Error copying file '%s' to '%s'\n", path, dest)
			return false
		}
	} else {
		destPath := fmt.Sprintf("%s/%s", dest, path)
		//fmt.Printf("LOOKING FOR DIR '%s'\n", destPath)
		_, err := os.Stat(destPath)
		if os.IsNotExist(err) {
			//fmt.Printf("MKDIR '%s', '%v'\n", destPath, err)
			if err := os.Mkdir(destPath, 0766); err != nil {
				fmt.Printf("Error creating dir '%s'\n", destPath)
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
				fmt.Printf("Error copying file or dir '%s' to '%s'\n", sourcePath, dest)
				return false
			}
		}
	}
	return true
}
