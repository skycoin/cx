package util

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

var workingDir string
var logFile bool

// CXSetWorkingDir ...
// TODO @evanlinjin: This should be in a module named 'util'.
func CXSetWorkingDir(dir string) {
	workingDir = dir
}

// CXLogFile ...
// TODO @evanlinjin: This should be in a module named 'util'.
func CXLogFile(enable bool) {
	logFile = enable
}

// CXOpenFile ...
// TODO @evanlinjin: This should be in a module named 'util'.
func CXOpenFile(filename string) (*os.File, error) {
	filename = filepath.Join(workingDir, filename)

	if logFile {
		fmt.Printf("CXOpenFile: Opening '%s'\n", filename)
	}

	file, err := os.Open(filename)
	if logFile && err != nil {
		fmt.Printf("CXOpenFile: Failed to open '%s': %v\n", filename, err)
	}
	return file, err
}

// CXCreateFile ...
// TODO @evanlinjin: This should be in a module named 'util'.
func CXCreateFile(filename string) (*os.File, error) {
	filename = filepath.Join(workingDir, filename)

	if logFile {
		fmt.Printf("Creating file : '%s', '%s'\n", workingDir, filename)
	}

	file, err := os.Create(filepath.Join(workingDir, filename))
	if logFile && err != nil {
		fmt.Printf("Failed to create file : '%s', '%s', err '%v'\n", workingDir, filename, err)
	}

	return file, err
}

// CXRemoveFile ...
// TODO @evanlinjin: This should be in a module named 'util'.
func CXRemoveFile(path string) error {
	if logFile {
		fmt.Printf("Removing file : '%s', '%s'\n", workingDir, path)
	}

	err := os.Remove(fmt.Sprintf("%s%s", workingDir, path))

	if logFile && err != nil {
		fmt.Printf("Failed to remove file : '%s', '%s', err '%v'\n", workingDir, path, err)
	}

	return err
}

// CXReadFile ...
// TODO @evanlinjin: This should be in a module named 'util'.
func CXReadFile(path string) ([]byte, error) {
	if logFile {
		fmt.Printf("Reading file : '%s', '%s'\n", workingDir, path)
	}

	bytes, err := ioutil.ReadFile(fmt.Sprintf("%s%s", workingDir, path))

	if logFile && err != nil {
		fmt.Printf("Failed to read file : '%s', '%s', err '%v'\n", workingDir, path, err)
	}

	return bytes, err
}

// CXStatFile ...
// TODO @evanlinjin: This should be in a module named 'util'.
func CXStatFile(path string) (os.FileInfo, error) {
	if logFile {
		fmt.Printf("Stating file : '%s', '%s'\n", workingDir, path)
	}

	fileInfo, err := os.Stat(fmt.Sprintf("%s%s", workingDir, path))

	if logFile && err != nil {
		fmt.Printf("Failed to stat file : '%s', '%s', err '%v'\n", workingDir, path, err)
	}

	return fileInfo, err
}

// CXMkdirAll ...
// TODO @evanlinjin: This should be in a module named 'util'.
func CXMkdirAll(path string, perm os.FileMode) error {
	if logFile {
		fmt.Printf("Creating dir : '%s'\n", path)
	}

	err := os.MkdirAll(fmt.Sprintf("%s%s", workingDir, path), perm)

	if logFile && err != nil {
		fmt.Printf("Failed to create dir : '%s', '%s', err '%v'\n", workingDir, path, err)
	}

	return err
}
