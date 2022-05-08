package main

import (
	"bufio"
	"bytes"
	"encoding/gob"
	"errors"
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

var SRC_PATH = os.Args[1]
var CURRENT_PATH = ""
var IMPORTED_DIRECTORIES = []string{}
var SKIP_PACKAGES = []string{"al", "gl", "glfw", "time", "os", "gltext", "cx", "json", "cipher", "tcp"}

func contains(list []string, element string) bool {
	for _, elem := range list {
		if elem == element {
			return true
		}
	}
	return false
}

func main() {
	if SRC_PATH[len(SRC_PATH)-1:] != "/" {
		SRC_PATH += "/src/"
	} else {
		SRC_PATH += "src/"
	}
	packageList := PackageList{}

	directoryList := []string{}
	err := filepath.Walk(SRC_PATH, func(path string, info os.FileInfo, err error) error {
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

	for _, path := range directoryList {
		packageList.addPackagesIn(path)
	}
	count := 0
	for range packageList.Packages {
		count++
	}
}

func (packageList *PackageList) addPackagesIn(path string) {
	if path[len(path)-1:] != "/" {
		path += "/"
	}
	CURRENT_PATH = path
	if contains(IMPORTED_DIRECTORIES, CURRENT_PATH) {
		return
	}
	newPackage := Package{}
	imports := []string{}
	fileList := []os.FileInfo{}
	files, err := ioutil.ReadDir(CURRENT_PATH)
	if err != nil {
		log.Fatal(err)
	}
	for _, fileInfo := range files {
		if fileInfo.Name()[len(fileInfo.Name())-2:] != "cx" {
			continue
		}
		fileList = append(fileList, fileInfo)
	}
	if len(fileList) == 1 {
		packageName, err := getPackageName(fileList[0])
		if err != nil {
			log.Fatal(err)
		}
		imports, err = getImports(fileList[0], imports)
		if err != nil {
			log.Fatal(err)
		}
		newPackage.PackageName = packageName
	}
	if len(fileList) > 1 {
		samePackage := false
		packageName := ""
		samePackage, packageName, imports, err = comparePackageNames(fileList, imports)
		if err != nil {
			log.Fatal(err)
		}
		if !samePackage {
			log.Print("Files in directory " + CURRENT_PATH + " are not all in the same newPackage.\nSource of the error: " + packageName)
			log.Fatal("ErrMismatchedPackageFiles")
		}
		newPackage.PackageName = packageName
	}
	newPackage.addFiles(fileList)
	packageList.hashPackage(&newPackage)
	IMPORTED_DIRECTORIES = append(IMPORTED_DIRECTORIES, CURRENT_PATH)
	for _, path := range imports {
		if contains(SKIP_PACKAGES, path) {
			continue
		}
		packageList.addPackagesIn(SRC_PATH + path)
	}
}

// For a list of cx files, get their package names and return if they match, and add the imports
func comparePackageNames(fileList []fs.FileInfo, imports []string) (bool, string, []string, error) {
	packageName, err := getPackageName(fileList[0])
	if err != nil {
		return false, "", imports, err
	}
	imports, err = getImports(fileList[0], imports)
	if err != nil {
		return false, "", imports, err
	}
	for i := 1; i < len(fileList); i++ {
		newPackageName, err := getPackageName(fileList[i])
		if err != nil {
			return false, "", imports, err
		}
		imports, err = getImports(fileList[i], imports)
		if err != nil {
			return false, "", imports, err
		}
		if newPackageName != packageName {
			return false, newPackageName + "in" + fileList[i].Name(), imports, nil
		}
	}
	return true, packageName, imports, nil
}

// Get the package name of a cx file
func getPackageName(fileInfo fs.FileInfo) (string, error) {
	file, err := os.Open(CURRENT_PATH + fileInfo.Name())
	if err != nil {
		return "", err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanWords)
	for scanner.Text() != "package" {
		scanner.Scan()
		if scanner.Text() == "var" || scanner.Text() == "const" || scanner.Text() == "type" || scanner.Text() == "func" {
			return "", errors.New("No package name found")
		}
	}
	scanner.Scan()
	return scanner.Text(), nil
}

// Get the import names in a cx file
func getImports(fileInfo fs.FileInfo, imports []string) ([]string, error) {
	file, err := os.Open(CURRENT_PATH + fileInfo.Name())
	if err != nil {
		return imports, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		if scanner.Text() == "import" {
			scanner.Scan()
			imports = append(imports, scanner.Text()[1:len(scanner.Text())-1])
		}
		if scanner.Text() == "var" || scanner.Text() == "const" || scanner.Text() == "type" || scanner.Text() == "func" {
			break
		}
	}
	return imports, nil
}

// Add the hashes of the files in fileList to the package
func (newPackage *Package) addFiles(fileList []fs.FileInfo) {
	for _, file := range fileList {
		newFile := File{
			FileName: file.Name(),
			Length:   uint32(file.Size()),
		}
		byteArray, err := ioutil.ReadFile(CURRENT_PATH + file.Name())
		if err != nil {
			log.Fatal(err)
		}
		newFile.Content = byteArray
		h := blake2b.Sum512(byteArray)
		newFile.Blake2Hash = string(h[:])
		newPackage.hashFile(&newFile)
	}
}

// Encode a file and put it in the specified package
func (newPackage *Package) hashFile(newFile *File) error {
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
func (packageList *PackageList) hashPackage(newPackage *Package) error {
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
