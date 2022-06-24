package loader

import (
	"bufio"
	"errors"
	"io/ioutil"
	"os"
	"strings"
	"sync"

	"github.com/skycoin/cx/cmd/packageloader/bolt"
	"github.com/skycoin/cx/cmd/packageloader/redis"
	"golang.org/x/crypto/blake2b"
)

var SKIP_PACKAGES = []string{"al", "gl", "glfw", "time", "os", "gltext", "cx", "json", "cipher", "tcp"}

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
			if filePackageName == "main" {
				fileMap["main"] = append(fileMap["main"], file)
				continue
			}
		}
		if filePackageName == packageName {
			fileMap[packageName] = append(fileMap[packageName], file)
		}
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

func LoadCXProgram(programName string, sourceCode []*os.File, database string) (err error) {
	fileMap, err := createFileMap(sourceCode)
	if err != nil {
		return err
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
