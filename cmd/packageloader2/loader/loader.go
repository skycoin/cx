package loader

import (
	"bufio"
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/skycoin/cx/cmd/packageloader2/bolt"
	"github.com/skycoin/cx/cmd/packageloader2/redis"
	"github.com/skycoin/cx/cx/util"
	"golang.org/x/crypto/blake2b"
)

var SKIP_PACKAGES = []string{"al", "gl", "glfw", "time", "os", "gltext", "cx", "json", "cipher", "tcp"}
var FileHashMap = make(map[string]string)
var PackageHashMap = make(map[string]string)

func Contains(list []string, element string) bool {
	for _, elem := range list {
		if elem == element {
			return true
		}
	}
	return false
}

func createFileMap(files []*os.File) (fileMap map[string][]*os.File, err error) {
	fileMap = make(map[string][]*os.File)
	for _, file := range files {
		path := strings.Split(file.Name(), "/")
		packageName := path[len(path)-2]
		filePackageName, err := getPackageName(file)
		if err != nil {
			return fileMap, err
		}
		if packageName == "src" {
			if filePackageName != "main" {
				return fileMap, fmt.Errorf("%s: package error: package %s found in main", filepath.Base(file.Name()), filePackageName)
			}
			fileMap["main"] = append(fileMap["main"], file)
			continue
		}
		if filePackageName != packageName {
			return fileMap, fmt.Errorf("%s: package error: package %s found in %v", filepath.Base(file.Name()), filePackageName, packageName)
		}
		fileMap[packageName] = append(fileMap[packageName], file)
	}
	return fileMap, nil
}

func getPackageName(file *os.File) (string, error) {
	openFile, err := os.Open(file.Name())
	if err != nil {
		return "", err
	}
	defer openFile.Close()

	scanner := bufio.NewScanner(openFile)

	for scanner.Scan() {
		line := strings.Split(scanner.Text(), " ")
		if line[0] == "package" {
			return line[1], nil
		}
	}
	return "", errors.New("file doesn't contain a package name")
}

func getImports(file *os.File) (imports []string, err error) {
	openFile, err := os.Open(file.Name())
	if err != nil {
		return imports, err
	}
	defer openFile.Close()

	scanner := bufio.NewScanner(openFile)

	for scanner.Scan() {
		line := strings.Split(scanner.Text(), " ")
		if line[0] == "import" {
			imports = append(imports, line[1][1:len(line[1])-1])
		}
	}
	return imports, nil
}

func checkForDependencyLoop(importMap map[string][]string) (err error) {
	for packageName := range importMap {
		for _, importName := range importMap[packageName] {
			if importName == packageName {
				return errors.New("Module " + packageName + " imports itself")
			}
			if Contains(importMap[importName], packageName) {
				return errors.New("Dependency loop between modules " + packageName + " and " + importName)
			}
		}
	}
	return nil
}

func LoadCXProgram(programName string, sourceCode []*os.File, database string) (err error) {

	fileMap, err := createFileMap(sourceCode)
	if err != nil {
		return err
	}

	var packageListStruct PackageList
	importMap := make(map[string][]string)

	err = addNewPackage(&packageListStruct, "main", fileMap, importMap, database)
	if err != nil {
		return err
	}

	err = loadImportPackages(&packageListStruct, "main", fileMap, importMap, database)
	if err != nil {
		return err
	}

	err = checkForDependencyLoop(importMap)
	if err != nil {
		return err
	}

	switch database {
	case "redis":
		redis.Add(programName, packageListStruct)
	case "bolt":
		value, err := packageListStruct.MarshalBinary()
		if err != nil {
			return err
		}
		bolt.Add(programName, value)
	}

	return nil
}

// Adds a new package to the package list and the package's imports to the import map
// 1. pass the package list pointer
// 2. pass the name of the package to be imported
// 3. pass the
func addNewPackage(packageListStruct *PackageList, packageName string, fileMap map[string][]*os.File, importMap map[string][]string, database string) error {

	files, ok := fileMap[packageName]

	// Checks if package is found in the directory
	if !ok && !Contains(SKIP_PACKAGES, packageName) {
		return fmt.Errorf("import %s not found", packageName)
	}
	// Skip if the import is a built-in package
	if Contains(SKIP_PACKAGES, packageName) {
		return nil
	}

	packageStruct, err := createPackageStruct(packageName, files, Package{}, database)
	if err != nil {
		return err
	}

	importList, err := createImportList(files, []string{})
	if err != nil {
		return err
	}

	newImportList := removeDuplicates(importList)

	importMap[packageName] = newImportList

	packageListStruct.appendPackage(&packageStruct, database)

	return nil
}

func fileStructFromFile(file *os.File) (File, error) {
	path := strings.Split(file.Name(), "/")
	fileName := path[len(path)-1]
	fileInfo, err := file.Stat()
	if err != nil {
		return File{}, err
	}
	fileBytes, err := ioutil.ReadFile(file.Name())
	if err != nil {
		return File{}, err
	}
	fileHash := blake2b.Sum512(fileBytes)

	fileStruct := File{
		FileName:   fileName,
		Length:     uint32(fileInfo.Size()),
		Content:    fileBytes,
		Blake2Hash: string(fileHash[:]),
	}

	return fileStruct, nil
}

func ParseArgsForCX(args []string, alsoSubdirs bool) (cxArgs []string, sourceCode []*os.File, fileNames []string) {
	skip := false // flag for skipping arg

	for _, arg := range args {

		// skip arg if skip flag is specified
		if skip {
			skip = false
			continue
		}

		// cli flags are either "--key=value" or "-key value"
		// we have to skip both cases
		if len(arg) > 1 && arg[0] == '-' {
			if !strings.Contains(arg, "=") {
				skip = true
			}
			continue
		}

		// cli cx flags are prefixed with "++"
		if len(arg) > 2 && arg[:2] == "++" {
			cxArgs = append(cxArgs, arg)
			continue
		}

		fi, err := util.CXStatFile(arg)
		if err != nil {
			println(fmt.Sprintf("%s: source file or library not found", arg))
			os.Exit(1)
		}

		switch mode := fi.Mode(); {
		case mode.IsDir():
			var fileList []string
			var err error

			// Checking if we want to check all subdirectories.
			if alsoSubdirs {
				fileList, err = filePathWalkDir(arg)
			} else {
				fileList, err = ioReadDir(arg)
				// fileList, err = filePathWalkDir(arg)
			}

			if err != nil {
				panic(err)
			}

			for _, path := range fileList {
				file, err := util.CXOpenFile(path)

				if err != nil {
					println(fmt.Sprintf("%s: source file or library not found", arg))
					os.Exit(1)
				}

				fiName := file.Name()
				fiNameLen := len(fiName)

				if fiNameLen > 2 && fiName[fiNameLen-3:] == ".cx" {
					// only loading .cx files
					sourceCode = append(sourceCode, file)
					fileNames = append(fileNames, fiName)
				}
			}
		case mode.IsRegular():
			file, err := util.CXOpenFile(arg)

			if err != nil {
				panic(err)
			}

			fileNames = append(fileNames, file.Name())
			sourceCode = append(sourceCode, file)
		}
	}

	return cxArgs, sourceCode, fileNames
}

func filePathWalkDir(root string) ([]string, error) {
	var files []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			files = append(files, path)
		}
		return err
	})
	return files, err
}

func ioReadDir(root string) ([]string, error) {
	var files []string
	fileInfo, err := ioutil.ReadDir(root)
	if err != nil {
		return files, err
	}

	for _, file := range fileInfo {
		files = append(files, fmt.Sprintf("%s/%s", root, file.Name()))
	}
	return files, nil
}

// Creates Blake2Hash from file UUIDs
func blake2HashFromFileUUID(fileUUID []string) ([64]byte, error) {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)

	err := encoder.Encode(fileUUID)
	if err != nil {
		return [64]byte{}, err
	}
	return blake2b.Sum512(buffer.Bytes()), nil
}

// Creates Package Struct
// 1. name - name of the struct
// 2. files - file list that are part of the package
// 3. packageStruct - pass an empty package struct "Package{}"
// 4. database - "redis" or "bolt"
// This function works recursively
func createPackageStruct(name string, files []*os.File, packageStruct Package, database string) (Package, error) {

	// Base case
	if len(files) == 0 {
		return packageStruct, nil
	}

	packageStruct.PackageName = name

	// Create file struct
	currentIndex := len(files) - 1
	currentFile := files[currentIndex]
	fileStruct, err := fileStructFromFile(currentFile)
	if err != nil {
		return packageStruct, err
	}

	// Append file struct
	err = packageStruct.appendFile(&fileStruct, database)
	if err != nil {
		return packageStruct, err
	}

	// Remove file from slice
	newFiles := files[:currentIndex]

	// Generate Blake2Hash and add to package struct
	if len(newFiles) == 0 {
		hash, err := blake2HashFromFileUUID(packageStruct.Files)
		if err != nil {
			return packageStruct, err
		}

		packageStruct.Blake2Hash = string(hash[:])
	}

	return createPackageStruct(name, newFiles, packageStruct, database)
}

// Create an import list from the list of files
// 1. pass the files list to files
// 2. pass an empty string slice "[]string{}" to importList
// This function works recursively
func createImportList(files []*os.File, importList []string) ([]string, error) {
	// Base case
	if len(files) == 0 {
		return importList, nil
	}

	// Gets imports
	currentIndex := len(files) - 1
	currentFile := files[currentIndex]
	imports, err := getImports(currentFile)
	if err != nil {
		return importList, err
	}

	// Removes file from list
	importList = append(importList, imports...)
	newFiles := files[:currentIndex]

	return createImportList(newFiles, importList)
}

// Removes Duplicates from list
func removeDuplicates(list []string) []string {
	var newList []string
	for _, elem := range list {
		if !Contains(newList, elem) {
			newList = append(newList, elem)
		}
	}
	return newList
}

// Loads Import Packages
// 1. packageListStruct - package list pointer
// 2. importName - package with imports
// 3. fileMap - file map that contains the files, not an empty map
// 4. importMap - import map that contains the imports, not an empty map
// 5. database - "redis" or "bolt"
// This function works recursive loading the import packages and then loading import packages of imports
func loadImportPackages(packageListStruct *PackageList, importName string, fileMap map[string][]*os.File, importMap map[string][]string, database string) error {

	errChannel := make(chan error, len(importMap))

	var wg sync.WaitGroup
	// loops over import list
	for _, imprt := range importMap[importName] {

		wg.Add(1)

		go func(packageListStruct *PackageList, imprt string, fileMap map[string][]*os.File, importMap map[string][]string, database string) {
			defer wg.Done()
			// Add package
			err := addNewPackage(packageListStruct, imprt, fileMap, importMap, database)
			if err != nil {
				errChannel <- err
				return
			}

			// Call itself for loading imports of the import
			loadImportPackages(packageListStruct, imprt, fileMap, importMap, database)
		}(packageListStruct, imprt, fileMap, importMap, database)

	}

	wg.Wait()
	close(errChannel)
	if err := <-errChannel; err != nil {
		return err
	}
	return nil
}
