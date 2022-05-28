package loader

import (
	"io/fs"
	"log"
	"os"
	"testing"
)

var testFileList = []fs.DirEntry{}
var testFileList2 = []fs.DirEntry{}

const TEST_SRC_PATH = "testInvalidProgram/src/"
const TEST_SRC_PATH2 = "testValidProgram/src/"
const TEST_SRC_PATH3 = "testVariousFiles/src/"

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

func TestCommentsPackage(t *testing.T) {
	files, err := os.ReadDir("testVariousFiles/")
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

func TestAddFilesRedis(t *testing.T) {
	DATABASE = "redis"
	CURRENT_PATH = TEST_SRC_PATH
	testPackage := Package{}
	testPackage.addFiles(testFileList)
	if len(testPackage.Files) != 2 {
		t.Error("Expected 2 files, got", len(testPackage.Files))
	}
}

func TestAddPackagesInRedis(t *testing.T) {
	DATABASE = "redis"
	SRC_PATH = TEST_SRC_PATH2
	testPackageList := PackageList{}
	testPackageList.addPackagesIn(SRC_PATH)
	if len(testPackageList.Packages) != 2 {
		t.Error("Expected 2 packages, got", len(testPackageList.Packages))
	}
}

func TestLoadRedis(t *testing.T) {
	DATABASE = "redis"
	LoadPackages("TestValid", "testValidProgram/")
}

func TestAddFilesBolt(t *testing.T) {
	DATABASE = "bolt"
	CURRENT_PATH = TEST_SRC_PATH
	testPackage := Package{}
	testPackage.addFiles(testFileList)
	if len(testPackage.Files) != 2 {
		t.Error("Expected 2 files, got", len(testPackage.Files))
	}
}

func TestAddPackagesInBolt(t *testing.T) {
	DATABASE = "bolt"
	SRC_PATH = TEST_SRC_PATH2
	testPackageList := PackageList{}
	testPackageList.addPackagesIn(SRC_PATH)
	if len(testPackageList.Packages) != 2 {
		t.Error("Expected 2 packages, got", len(testPackageList.Packages))
	}
}

func TestLoadBolt(t *testing.T) {
	DATABASE = "bolt"
	LoadPackages("TestValid", "testValidProgram/")
}
