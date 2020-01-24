// +build base

package cxcore

import "os"
import "fmt"

var workingDir string

func CXSetWorkingDir(dir string) {
	workingDir = dir
}

func CXOpenFile(path string) (*os.File, error) {
	return os.Open(fmt.Sprintf("%s%s", workingDir, path))
}

func CXCreateFile(path string) (*os.File, error) {
	return os.Create(fmt.Sprintf("%s%s", workingDir, path))
}

func CXRemoveFile(path string) error {
	return os.Remove(fmt.Sprintf("%s%s", workingDir, path))
}

func CXStatFile(path string) (os.FileInfo, error) {
	return os.Stat(fmt.Sprintf("%s%s", workingDir, path))
}

func CXMkdirAll(path string, perm os.FileMode) error {
	return os.MkdirAll(path, perm)
}
