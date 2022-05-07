package main

import (
	"bufio"
	"bytes"
	"encoding/gob"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"golang.org/x/crypto/blake2b"
)

type File struct {
	FileName   string
	Length     uint32
	Content    []byte
	Blake2Hash string
}

type Package struct {
	PackageName string
	Files       []string
}

type PackageList struct {
	Packages []string
}

func main() {
	directoryList := []string{}
	err := filepath.Walk(os.Args[1], func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			directoryList = append(directoryList, path)
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	packageList := PackageList{}
	var importList []string

	// For each directory, create a package and get the files in that directory
	for _, path := range directoryList {
		var fileList = []fs.FileInfo{}
		files, err := ioutil.ReadDir(path)
		if err != nil {
			log.Fatal(err)
		}
		for _, file := range files {
			if file.Name()[len(file.Name())-2:] != "cx" {
				continue
			}
			fileList = append(fileList, file)
		}

		newPackage := Package{}

		// For each file in the directory, add it to the package
		for i := 1; i < len(fileList); i++ {
			samePackage, packageName, err := comparePackages(path+"/"+fileList[i].Name(), path+"/"+fileList[i-1].Name())
			if err != nil {
				log.Fatal(err)
			}

			if !samePackage {
				log.Print("Files in directory " + path + " are not all in the same newPackage.\nSource of the error: " + fileList[i].Name())
				log.Fatal("ErrMismatchedPackageFiles")
			}
			// Once files are taken care of, add imports to the package
			newImports, err := getImports(path + "/" + fileList[i].Name())
			if err != nil {
				log.Fatal(err)
			}
			importList = append(importList, newImports...)
			newFile := File{
				FileName: fileList[i].Name(),
				Length:   uint32(fileList[i].Size()),
			}
			h := blake2b.Sum512(newFile.Content)
			newFile.Blake2Hash = string(h[:])
			newFile.Content, err = ioutil.ReadFile(path + "/" + fileList[i].Name())
			if err != nil {
				log.Fatal(err)
			}
			newPackage.PackageName = packageName
			newPackage.hashFile(newFile)
		}
		packageList.hashPackage(newPackage)
	}
}

// Encode a file and put it in the specified package
func (newPackage Package) hashFile(newFile File) error {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)

	err := encoder.Encode(newFile)
	if err != nil {
		return err
	}
	h := blake2b.Sum512(buffer.Bytes())
	newPackage.Files = append(newPackage.Files, string(h[:]))
	return nil
}

// Encode a package and put it in the specified package list
func (packageList PackageList) hashPackage(newPackage Package) error {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)

	err := encoder.Encode(newPackage)
	if err != nil {
		return err
	}
	h := blake2b.Sum512(buffer.Bytes())
	packageList.Packages = append(packageList.Packages, string(h[:]))
	return nil
}

func comparePackages(filepath1 string, filepath2 string) (bool, string, error) {
	file1, err := os.Open(filepath1)
	if err != nil {
		return false, "", err
	}
	defer file1.Close()
	file2, err := os.Open(filepath2)
	if err != nil {
		return false, "", err
	}
	defer file2.Close()
	scanner1 := bufio.NewScanner(file1)
	scanner2 := bufio.NewScanner(file2)
	scanner1.Split(bufio.ScanWords)
	scanner2.Split(bufio.ScanWords)
	for scanner1.Text() != "package" {
		scanner1.Scan()
	}
	for scanner2.Text() != "package" {
		scanner2.Scan()
	}
	scanner1.Scan()
	scanner2.Scan()
	return scanner1.Text() == scanner2.Text(), scanner1.Text(), nil
}

func getImports(filepath string) (imports []string, err error) {
	file, err := os.Open(filepath)
	if err != nil {
		return []string{}, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		if scanner.Text() == "import" {
			scanner.Scan()
			imports = append(imports, scanner.Text())
		}
	}
	return imports, nil
}
