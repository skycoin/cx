package loader

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/skycoin/cx/cx/globals"
)

func LoadCXProgramNoSave(sourceCode []*os.File, rootDir []string) (files []*File, err error) {

	// If it's a single file program with all the packages in one file
	if len(sourceCode) == 1 && hasMultiplePkgs(sourceCode[0]) {
		file, err := fileStructFromFile(sourceCode[0])
		if err != nil {
			return files, err
		}
		files = append(files, &file)
		return files, nil
	}
	// Else it's a multiple files program

	// Gets the source files
	fileMap, err := createFileMap(sourceCode, rootDir)
	if err != nil {
		return files, err
	}

	var mx sync.Mutex
	importMap := make(map[string][]string)

	sourceCodes, ok := fileMap["main"]
	if !ok {
		return files, fmt.Errorf("main package not found")
	}

	// Dependency loops are checked before adding to DB
	// Step 9
	err = checkImports("main", sourceCodes, importMap, &mx)
	if err != nil {
		return files, err
	}

	// load the imported packages
	imprtFiles, err := loadImportPackagesNoSave("main", fileMap, importMap, &mx)
	if err != nil {
		return files, err
	}

	files = append(files, imprtFiles...)

	return files, nil
}

//	 Loads Import Packages
//		1. packageListStruct - package list pointer
//		2. importName - package with imports
//		3. fileMap - file map that contains the files
//		4. importMap - import map that contains the imports
//		5. database - "redis" or "bolt"
//
// This function works recursively, loading the import packages and then loading import packages of imports
func loadImportPackagesNoSave(importName string, fileMap map[string][]*os.File, importMap map[string][]string, mx *sync.Mutex) ([]*File, error) {

	filesChannel := make(chan []*File, len(importMap))
	errChannel := make(chan error, len(importMap))

	var wg sync.WaitGroup
	// loops over import list
	for _, imprt := range importMap[importName] {

		wg.Add(1)

		go func(imprt string, fileMap map[string][]*os.File, importMap map[string][]string, errChannel chan error, wg *sync.WaitGroup, mx *sync.Mutex) {
			defer wg.Done()
			files, ok := fileMap[imprt]

			// If package is not found check if it's in the cxpath
			if !ok && !Contains(SKIP_PACKAGES, imprt) {
				_, rootDir, sourceCode, _ := ParseArgsForCX([]string{filepath.Join(globals.SRCPATH, imprt)}, true)
				tmpMap, err := createFileMap(rootDir, sourceCode)
				if err != nil {
					errChannel <- err
					return
				}
				for k, v := range tmpMap {
					fileMap[k] = v
				}
				if strings.Contains(imprt, "/") {
					tokens := strings.Split(imprt, "/")
					imprt = tokens[len(tokens)-1]
				}
				files, ok = fileMap[imprt]
				if !ok {
					errChannel <- fmt.Errorf("package %s not found", imprt)
					return
				}
			}

			// Skip if the import is a built-in package
			if Contains(SKIP_PACKAGES, imprt) {
				return
			}

			err := checkImports(imprt, files, importMap, mx)
			if err != nil {
				errChannel <- err
				return
			}

			fileStructList, err := filesToLoaderStruct(files)
			if err != nil {
				errChannel <- err
				return
			}

			filesChannel <- fileStructList

			// Call itself for loading imports of the import
			imprtFiles, err := loadImportPackagesNoSave(imprt, fileMap, importMap, mx)
			if err != nil {
				errChannel <- err
				return
			}

			filesChannel <- imprtFiles

		}(imprt, fileMap, importMap, errChannel, &wg, mx)

	}

	wg.Wait()

	close(filesChannel)
	close(errChannel)

	var files []*File

	for err := range errChannel {
		if err != nil {
			return files, err
		}
	}

	for file := range filesChannel {
		files = append(files, file...)
	}

	return files, nil
}

func filesToLoaderStruct(files []*os.File) ([]*File, error) {
	var fileStructList []*File
	for _, file := range files {
		fileStruct, err := fileStructFromFile(file)
		if err != nil {
			return nil, err
		}
		fileStructList = append(fileStructList, &fileStruct)
	}
	return fileStructList, nil
}
