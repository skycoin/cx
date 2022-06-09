package loader

import (
	"bufio"
	"errors"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/skycoin/cx/cmd/packageloader/bolt"
	"github.com/skycoin/cx/cmd/packageloader/redis"
	"golang.org/x/crypto/blake2b"
)

var DATABASE = "bolt"
var SRC_PATH string
var CURRENT_PATH string
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

func LoadPackages(programName string, path string) error {
	SRC_PATH = path + "src/"

	packageList := PackageList{}

	directoryList := []string{}
	err := filepath.WalkDir(SRC_PATH, func(path string, info fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			directoryList = append(directoryList, path)
		}
		return nil
	})
	if err != nil {
		return err
	}

	for _, path := range directoryList {
		packageList.addPackagesIn(path)
	}
	switch DATABASE {
	case "redis":
		redis.Add(programName, packageList)
	case "bolt":
		value, err := packageList.MarshalBinary()
		if err != nil {
			return err
		}
		bolt.Add(programName, value)
	}
	return nil
}

func (packageList *PackageList) addPackagesIn(path string) error {
	if path[len(path)-1:] != "/" {
		path += "/"
	}
	CURRENT_PATH = path
	if contains(IMPORTED_DIRECTORIES, CURRENT_PATH) {
		return nil
	}
	newPackage := Package{}
	imports := []string{}
	fileList := []os.DirEntry{}
	files, err := os.ReadDir(CURRENT_PATH)
	if err != nil {
		return err
	}
	for _, dirEntry := range files {
		if dirEntry.Name()[len(dirEntry.Name())-2:] != "cx" {
			continue
		}
		fileList = append(fileList, dirEntry)
	}
	if len(fileList) == 1 {
		packageName, err := getPackageName(fileList[0])
		if err != nil {
			return err
		}
		imports, err = getImports(fileList[0], imports)
		if err != nil {
			return err
		}
		newPackage.PackageName = packageName
	}
	if len(fileList) > 1 {
		samePackage := false
		packageName := ""
		samePackage, packageName, imports, err = comparePackageNames(fileList, imports)
		if err != nil {
			return err
		}
		if !samePackage {
			log.Print("Files in directory " + CURRENT_PATH + " are not all in the same newPackage.\nSource of the error: " + packageName)
			return errors.New("ErrMismatchedPackageFiles")
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
	return nil
}

// For a list of cx files, get their package names and return if they match, and add the imports
func comparePackageNames(fileList []fs.DirEntry, imports []string) (bool, string, []string, error) {
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
func getPackageName(dirEntry fs.DirEntry) (string, error) {
	file, err := os.Open(CURRENT_PATH + dirEntry.Name())
	if err != nil {
		return "", err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanWords)

	wordBefore := ""
	for scanner.Scan() {
		if scanner.Text() != "package" {
			wordBefore = scanner.Text()
			continue
		}
		if wordBefore == "//" {
			wordBefore = scanner.Text()
			continue
		}
		if scanner.Text() == "import" || scanner.Text() == "var" || scanner.Text() == "const" || scanner.Text() == "type" || scanner.Text() == "func" {
			return "", errors.New("no package name found")
		}
		break
	}
	scanner.Scan()
	return scanner.Text(), nil
}

// Get the import names in a cx file
func getImports(dirEntry fs.DirEntry, imports []string) ([]string, error) {
	file, err := os.Open(CURRENT_PATH + dirEntry.Name())
	if err != nil {
		return imports, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanWords)

	wordBefore := ""
	for scanner.Scan() {
		if scanner.Text() != "import" {
			wordBefore = scanner.Text()
			continue
		}
		if wordBefore == "//" {
			wordBefore = scanner.Text()
			continue
		}
		if scanner.Text() == "var" || scanner.Text() == "const" || scanner.Text() == "type" || scanner.Text() == "func" {
			break
		}
		scanner.Scan()
		imports = append(imports, scanner.Text()[1:len(scanner.Text())-1])
		wordBefore = scanner.Text()
	}
	return imports, nil
}

// Add the hashes of the files in fileList to the package
func (newPackage *Package) addFiles(fileList []fs.DirEntry) error {
	for _, file := range fileList {
		fileInfo, err := file.Info()
		if err != nil {
			return err
		}

		newFile := File{
			FileName: file.Name(),
			Length:   uint32(fileInfo.Size()),
		}
		byteArray, err := ioutil.ReadFile(CURRENT_PATH + file.Name())
		if err != nil {
			return err
		}
		newFile.Content = byteArray
		h := blake2b.Sum512(byteArray)
		newFile.Blake2Hash = string(h[:])
		newPackage.hashFile(&newFile)
	}
	return nil
}
