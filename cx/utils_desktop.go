// +build base

package cxcore

import (
	"fmt"
	"io/ioutil"
	"os"
)

var workingDir string
var logFile bool

func CXSetWorkingDir(dir string) {
	workingDir = dir
}

func CXLogFile(enable bool) {
	logFile = enable
}

func CXOpenFile(path string) (*os.File, error) {
	if logFile {
		fmt.Printf("Opening file : '%s', '%s'\n", workingDir, path)
	}

	file, err := os.Open(fmt.Sprintf("%s%s", workingDir, path))
	if logFile && err != nil {
		fmt.Printf("Failed to open file : '%s', '%s'\n", workingDir, path)
	}
	return file, err
}

func CXCreateFile(path string) (*os.File, error) {
	if logFile {
		fmt.Printf("Creating file : '%s', '%s'\n", workingDir, path)
	}

	file, err := os.Create(fmt.Sprintf("%s%s", workingDir, path))
	if logFile && err != nil {
		fmt.Printf("Failed to create file : '%s', '%s'\n", workingDir, path)
	}

	return file, err
}

func CXRemoveFile(path string) error {
	if logFile {
		fmt.Printf("Removing file : '%s', '%s'\n", workingDir, path)
	}

	err := os.Remove(fmt.Sprintf("%s%s", workingDir, path))

	if logFile && err != nil {
		fmt.Printf("Failed to remove file : '%s', '%s'\n", workingDir, path)
	}

	return err
}

func CXReadFile(path string) ([]byte, error) {
	if logFile {
		fmt.Printf("Reading file : '%s', '%s'\n", workingDir, path)
	}

	bytes, err := ioutil.ReadFile(fmt.Sprintf("%s%s", workingDir, path))

	if logFile && err != nil {
		fmt.Printf("Failed to read file : '%s', '%s'\n", workingDir, path)
	}

	return bytes, err
}

func CXStatFile(path string) (os.FileInfo, error) {
	if logFile {
		fmt.Printf("Stating file : '%s', '%s'\n", workingDir, path)
	}

	fileInfo, err := os.Stat(fmt.Sprintf("%s%s", workingDir, path))

	if logFile && err != nil {
		fmt.Printf("Failed to stat file : '%s', '%s'\n", workingDir, path)
	}

	return fileInfo, err
}

func CXMkdirAll(path string, perm os.FileMode) error {
	if logFile {
		fmt.Printf("Creating dir : '%s'\n", path)
	}

	err := os.MkdirAll(path, perm)

	if logFile && err != nil {
		fmt.Printf("Failed to create dir : '%s', '%s'\n", workingDir, path)
	}

	return err
}
