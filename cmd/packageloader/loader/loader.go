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

	"github.com/skycoin/cx/cmd/packageloader/bolt"
	"github.com/skycoin/cx/cmd/packageloader/redis"
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

func createImportMap(fileMap map[string][]*os.File) (importMap map[string][]string, err error) {
	importMap = make(map[string][]string)
	for packageName := range fileMap {
		packageImports := []string{}
		for _, file := range fileMap[packageName] {
			newImports, err := getImports(file)
			if err != nil {
				return importMap, err
			}
			packageImports = append(packageImports, newImports...)
		}

		packageImportsWithoutDuplicates := removeDuplicates(packageImports)
		importMap[packageName] = packageImportsWithoutDuplicates
	}
	return importMap, nil
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
	importMap, err := createImportMap(fileMap)
	if err != nil {
		return err
	}
	err = checkForDependencyLoop(importMap)
	if err != nil {
		return
	}

	packageListStruct := PackageList{}
	for key, files := range fileMap {
		addNewPackage(&packageListStruct, key, files, database)
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

func addNewPackage(packageListStruct *PackageList, packageName string, files []*os.File, database string) error {
	packageStruct := Package{
		PackageName: packageName,
	}

	var wg sync.WaitGroup
	for _, file := range files {
		wg.Add(1)
		errs := make(chan error)
		go func(file *os.File) {
			defer wg.Done()
			fileStruct, err := fileStructFromFile(file)
			if err != nil {
				errs <- err
			}
			err = packageStruct.appendFile(&fileStruct, database)
			if err != nil {
				errs <- err
			}
			close(errs)
		}(file)
		for err := range errs {
			if err != nil {
				return err
			}
		}
	}
	wg.Wait()

	hash, err := blake2HashFromFileUUID(packageStruct.Files)
	if err != nil {
		return err
	}

	packageStruct.Blake2Hash = string(hash[:])

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

func blake2HashFromFileUUID(fileUUID []string) ([64]byte, error) {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)

	err := encoder.Encode(fileUUID)
	if err != nil {
		return [64]byte{}, err
	}
	return blake2b.Sum512(buffer.Bytes()), nil
}

func removeDuplicates(imports []string) []string {
	var newImports []string
	for _, imprt := range imports {
		if Contains(newImports, imprt) {
			continue
		}
		newImports = append(newImports, imprt)
	}

	return newImports
}
