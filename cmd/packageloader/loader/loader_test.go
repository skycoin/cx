package loader

import (
	"io/fs"
	"log"
	"os"
	"testing"
)

var testFileList = []fs.DirEntry{}
var testFileList2 = []fs.DirEntry{}

const TEST_SRC_PATH = "test_folder/test_invalid_program/src/"
const TEST_SRC_PATH2 = "test_folder/test_valid_program/src/"
const TEST_SRC_PATH3 = "test_folder/test_various_files/src/"

func init() {
	files, err := os.ReadDir(TEST_SRC_PATH)
	if err != nil {
		log.Fatal(err)
	}
	for _, fileInfo := range files {
		if fileInfo.Name()[len(fileInfo.Name())-2:] != "cx" {
			continue
		}
		testFileList = append(testFileList, fileInfo)
	}

	files, err = os.ReadDir(TEST_SRC_PATH + "testimport/")
	if err != nil {
		log.Fatal(err)
	}
	for _, fileInfo := range files {
		if fileInfo.Name()[len(fileInfo.Name())-2:] != "cx" {
			continue
		}
		testFileList2 = append(testFileList2, fileInfo)
	}

}

func TestContains(t *testing.T) {
	list := []string{"a", "b", "c"}
	if !contains(list, "a") {
		t.Error("Expected true, got false")
	}
	if contains(list, "d") {
		t.Error("Expected false, got true")
	}
}

func TestGetPackageName(t *testing.T) {
	CURRENT_PATH = TEST_SRC_PATH
	testPackageName, err := getPackageName(testFileList[0])
	if err != nil {
		t.Error(err)
	}
	if testPackageName != "main" {
		t.Error("Expected main, got", testPackageName)
	}
}

func TestGetImports(t *testing.T) {
	CURRENT_PATH = TEST_SRC_PATH
	testImports := []string{}
	testImports, err := getImports(testFileList[0], testImports)
	if err != nil {
		t.Error(err)
	}
	if len(testImports) != 1 {
		t.Error("Expected 1 import, got", len(testImports))
	}
	if testImports[0] != "testimport" {
		t.Error("Expected testimport, got", testImports[0])
	}
}

func TestComparePackageNames(t *testing.T) {
	testValues := []struct {
		Path     string
		Files    []fs.DirEntry
		Expected bool
		ExpValue string
	}{{TEST_SRC_PATH, testFileList, false, ""}, {TEST_SRC_PATH + "testimport/", testFileList2, true, "testimport"}}
	for _, testcase := range testValues {
		CURRENT_PATH = testcase.Path
		testSamePackage, testPackageName, _, err := comparePackageNames(testcase.Files, []string{})
		if err != nil {
			t.Error(err)
		}
		if testSamePackage != testcase.Expected {
			t.Error("Expected", testcase.Expected, "got", testSamePackage)
		}
		if testSamePackage == false {
			return
		}
		if testPackageName != testcase.ExpValue {
			t.Error("Expected", testcase.ExpValue, "got", testPackageName)
		}
	}
}

func TestCommentsPackage(t *testing.T) {
	CURRENT_PATH = "test_folder/test_various_files/"
	files, err := os.ReadDir(CURRENT_PATH)
	if err != nil {
		log.Fatal(err)
	}
	for _, f := range files {
		if f.Name() == "package_comment.cx" {
			packageName, err := getPackageName(f)
			if err != nil {
				t.Error(err)
			}
			if packageName != "main" {
				t.Error("Expected package main, got", packageName)
			}
		}
		if f.Name() == "import_comment.cx" {
			importList := []string{}
			importList, err := getImports(f, importList)
			if err != nil {
				t.Error(err)
			}
			if len(importList) != 1 || importList[0] != "cx" {
				t.Error("Expected import cx, got imports", importList)
			}
		}
	}
}

func TestAddFiles(t *testing.T) {
	for _, v := range []string{"redis", "bolt"} {
		DATABASE = v
		CURRENT_PATH = TEST_SRC_PATH
		testPackage := Package{}
		testPackage.addFiles(testFileList)
		if len(testPackage.Files) != 2 {
			t.Error("Expected 2 files, got", len(testPackage.Files))
		}
	}
}

func TestAddPackagesIn(t *testing.T) {
	for _, v := range []string{"redis", "bolt"} {
		DATABASE = v
		SRC_PATH = TEST_SRC_PATH2
		IMPORTED_DIRECTORIES = []string{}
		testPackageList := PackageList{}
		testPackageList.addPackagesIn(SRC_PATH)
		if len(testPackageList.Packages) != 2 {
			t.Error("Expected 2 packages, got", len(testPackageList.Packages))
		}
	}
}

func TestLoad(t *testing.T) {
	for _, v := range []string{"redis", "bolt"} {
		DATABASE = v
		IMPORTED_DIRECTORIES = []string{}
		LoadPackages("TestValid", "test_folder/test_valid_program/")
	}
}
