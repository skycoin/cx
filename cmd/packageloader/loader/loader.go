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

var SKIP_PACKAGES = []string{"al", "gl", "glfw", "time", "os", "gltext", "cx", "json", "cipher", "tcp"}

func contains(list []string, element string) bool {
	for _, elem := range list {
		if elem == element {
			return true
		}
	}
	return false
}

func LoadPackages(programName string, path string, database string) error {
	srcPath := path + "src/"

	packageList := PackageList{}

	directoryList := []string{}
	err := filepath.WalkDir(srcPath, func(path string, info fs.DirEntry, err error) error {
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
	importedDirectories := []string{}
	for _, path := range directoryList {
		importedDirectories, err = addPackages(&packageList, srcPath, path, database, importedDirectories)
		if err != nil {
			return err
		}
	}
	switch database {
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

func addPackages(packageList *PackageList, srcPath string, packagePath string, database string, importedDirectories []string) ([]string, error) {
	if packagePath[len(packagePath)-1:] != "/" {
		packagePath += "/"
	}
	if contains(importedDirectories, packagePath) {
		return importedDirectories, nil
	}

	newPackage := Package{}
	imports := []string{}
	fileList := []os.DirEntry{}
	files, err := os.ReadDir(packagePath)
	if err != nil {
		return importedDirectories, err
	}

	for _, dirEntry := range files {
		if dirEntry.Name()[len(dirEntry.Name())-2:] != "cx" {
			continue
		}
		fileList = append(fileList, dirEntry)
	}

	if len(fileList) == 1 {
		packageName, err := getPackageName(fileList[0], packagePath)
		if err != nil {
			return importedDirectories, err
		}
		newImports, err := getImports(fileList[0], packagePath)
		imports = append(imports, newImports...)
		if err != nil {
			return importedDirectories, err
		}
		newPackage.PackageName = packageName
	}
	if len(fileList) > 1 {
		samePackage := false
		packageName := ""
		samePackage, packageName, imports, err = comparePackageNames(fileList, packagePath, imports)
		if err != nil {
			return importedDirectories, err
		}
		if !samePackage {
			log.Print("Files in directory " + packagePath + " are not all in the same newPackage.\nSource of the error: " + packageName)
			return importedDirectories, errors.New("ErrMismatchedPackageFiles")
		}
		newPackage.PackageName = packageName
	}

	addFiles(&newPackage, fileList, packagePath, database)
	packageList.addPackage(&newPackage, database)
	importedDirectories = append(importedDirectories, packagePath)

	for _, importName := range imports {
		importPath := srcPath + importName + "/"
		if contains(SKIP_PACKAGES, importName) || contains(importedDirectories, importPath) {
			continue
		}
		return addPackages(packageList, srcPath, importPath, database, importedDirectories)
	}
	return importedDirectories, nil
}

// For a list of cx files, get their package names and return if they match, and add the imports
func comparePackageNames(fileList []fs.DirEntry, packagePath string, imports []string) (bool, string, []string, error) {
	packageName, err := getPackageName(fileList[0], packagePath)
	if err != nil {
		return false, "", imports, err
	}
	newImports, err := getImports(fileList[0], packagePath)
	imports = append(imports, newImports...)
	if err != nil {
		return false, "", imports, err
	}
	for i := 1; i < len(fileList); i++ {
		newPackageName, err := getPackageName(fileList[i], packagePath)
		if err != nil {
			return false, "", imports, err
		}
		newImports, err := getImports(fileList[i], packagePath)
		imports = append(imports, newImports...)
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
func getPackageName(dirEntry fs.DirEntry, packagePath string) (string, error) {
	file, err := os.Open(packagePath + dirEntry.Name())
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
func getImports(dirEntry fs.DirEntry, importPath string) (imports []string, err error) {
	file, err := os.Open(importPath + dirEntry.Name())
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
func addFiles(newPackage *Package, fileList []fs.DirEntry, packagePath string, database string) error {
	for _, file := range fileList {
		fileInfo, err := file.Info()
		if err != nil {
			return err
		}

		newFile := File{
			FileName: file.Name(),
			Length:   uint32(fileInfo.Size()),
		}
		byteArray, err := ioutil.ReadFile(packagePath + file.Name())
		if err != nil {
			return err
		}
		newFile.Content = byteArray
		h := blake2b.Sum512(byteArray)
		newFile.Blake2Hash = string(h[:])
		newPackage.addFile(&newFile, database)
	}
	return nil
}
