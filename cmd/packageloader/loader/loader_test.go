package loader_test

import (
	"os"
	"testing"

	"github.com/skycoin/cx/cmd/packageloader/loader"
	"github.com/skycoin/cx/cx/ast"
)

const TEST_SRC_PATH_VALID = "test_folder/test_invalid_program/src/"
const TEST_SRC_PATH_INVALID = "test_folder/test_valid_program/src/"
const TEST_PACKAGE_FILE = "test_folder/test_various_files/package.cx"

func TestContains(t *testing.T) {
	list := []string{"a", "b", "c"}
	if !loader.Contains(list, "a") {
		t.Error("Expected true, got false")
	}
	if loader.Contains(list, "d") {
		t.Error("Expected false, got true")
	}
}

func TestGetPackageName(t *testing.T) {
	file, err := os.Open(TEST_PACKAGE_FILE)
	if err != nil {
		t.Error(err)
	}
	testPackageName, err := loader.GetPackageName(file)
	if err != nil {
		t.Error(err)
	}
	if testPackageName != "main" {
		t.Error("Wrong package name:", testPackageName)
	}
}

func TestCreateFileMap(t *testing.T) {
	_, sourceCodes, _ := ast.ParseArgsForCX([]string{TEST_SRC_PATH_VALID}, true)
	_, err := loader.CreateFileMap(sourceCodes)
	if err != nil {
		t.Error(err)
	}
	//TODO: Find a way to reliably test this function
}

func TestFileStructFromFile(t *testing.T) {
	file, err := os.Open(TEST_PACKAGE_FILE)
	if err != nil {
		t.Error(err)
	}
	testFileStruct, err := loader.FileStructFromFile(file)
	if err != nil {
		t.Error(err)
	}
	if testFileStruct.FileName != "package.cx" {
		t.Error("wrong file name:", testFileStruct.FileName)
	}
}
