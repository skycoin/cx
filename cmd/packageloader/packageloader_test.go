package main

import (
	"io/ioutil"
	"log"
	"os"
	"testing"
)

var testFileList = []os.FileInfo{}
var testFileList2 = []os.FileInfo{}

const TEST_SRC_PATH = "./test/src/"

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
	if testPackageName != "testone" {
		t.Error("Expected testone, got", testPackageName)
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
	testImports := []string{}
	testSamePackage, _, testImports, err := comparePackageNames(testFileList, testImports)
	if err != nil {
		t.Error(err)
	}
	if testSamePackage {
		t.Error("Expected false, got true")
	}
}

func TestComparePackageNamesTrue(t *testing.T) {
	CURRENT_PATH = TEST_SRC_PATH + "testimport/"
	testImports := []string{}
	testSamePackage, testPackageName, testImports, err := comparePackageNames(testFileList2, testImports)
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
