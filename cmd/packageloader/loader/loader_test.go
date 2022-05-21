package loader

import (
	"io/ioutil"
	"log"
	"os"
	"testing"
)

var testFileList = []os.FileInfo{}
var testFileList2 = []os.FileInfo{}

const TEST_SRC_PATH = "testInvalidProgram/src/"
const TEST_SRC_PATH2 = "testValidProgram/src/"
const TEST_SRC_PATH3 = "testVariousFiles/src/"

func init() {

	files, err := ioutil.ReadDir(TEST_SRC_PATH)
	if err != nil {
		log.Fatal(err)
	}
	for _, fileInfo := range files {
		if fileInfo.Name()[len(fileInfo.Name())-2:] != "cx" {
			continue
		}
		testFileList = append(testFileList, fileInfo)
	}

	files, err = ioutil.ReadDir(TEST_SRC_PATH + "testimport/")
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

func TestComparePackageNamesFalse(t *testing.T) {
	CURRENT_PATH = TEST_SRC_PATH
	testSamePackage, _, _, err := comparePackageNames(testFileList, []string{})
	if err != nil {
		t.Error(err)
	}
	if testSamePackage {
		t.Error("Expected false, got true")
	}
}

func TestComparePackageNamesTrue(t *testing.T) {
	CURRENT_PATH = TEST_SRC_PATH + "testimport/"
	testSamePackage, testPackageName, _, err := comparePackageNames(testFileList2, []string{})
	if err != nil {
		t.Error(err)
	}
	if !testSamePackage {
		t.Error("Expected true, got false")
	}
	if testPackageName != "testimport" {
		t.Error("Expected testimport, got", testPackageName)
	}
}

func TestAddFiles(t *testing.T) {
	CURRENT_PATH = TEST_SRC_PATH
	testPackage := Package{}
	testPackage.addFiles(testFileList)
	if len(testPackage.Files) != 2 {
		t.Error("Expected 2 files, got", len(testPackage.Files))
	}
}

func TestAddPackagesIn(t *testing.T) {
	SRC_PATH = TEST_SRC_PATH2
	testPackageList := PackageList{}
	testPackageList.addPackagesIn(SRC_PATH)
	if len(testPackageList.Packages) != 2 {
		t.Error("Expected 2 packages, got", len(testPackageList.Packages))
	}
}

func TestLoad(t *testing.T) {
	LoadPackages("TestValid", "testValidProgram/")
}

func TestCommentsPackage(t *testing.T) {
	files, err := ioutil.ReadDir("testVariousFiles/")
	if err != nil {
		log.Fatal(err)
	}
	CURRENT_PATH = "testVariousFiles/"
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
