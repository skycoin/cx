package loader_test

import (
	"io/fs"
	"log"
	"os"
	"testing"

	"github.com/skycoin/cx/cmd/packageloader/loader"
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
	if !loader.Contains(list, "a") {
		t.Error("Expected true, got false")
	}
	if loader.Contains(list, "d") {
		t.Error("Expected false, got true")
	}
}

func TestGetPackageName(t *testing.T) {
	testPackageName, err := loader.GetPackageName(testFileList[0], TEST_SRC_PATH)
	if err != nil {
		t.Error(err)
	}
	if testPackageName != "main" {
		t.Error("Expected main, got", testPackageName)
	}
}

func TestGetImports(t *testing.T) {
	testImports, err := loader.GetImports(testFileList[0], TEST_SRC_PATH)
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
	tests := []struct {
		Scenario      string
		Path          string
		Files         []fs.DirEntry
		ExpectedBool  bool
		ExpectedValue string
	}{
		{
			Scenario:      "Test case for a folder with differing package names",
			Path:          TEST_SRC_PATH,
			Files:         testFileList,
			ExpectedBool:  false,
			ExpectedValue: "",
		},
		{
			Scenario:      "Test case for a folder with a single package name",
			Path:          TEST_SRC_PATH + "testimport/",
			Files:         testFileList2,
			ExpectedBool:  true,
			ExpectedValue: "testimport",
		},
	}
	for _, testcase := range tests {
		t.Run(testcase.Scenario, func(t *testing.T) {
			testSamePackage, testPackageName, _, err := loader.ComparePackageNames(testcase.Files, testcase.Path, []string{})
			if err != nil {
				t.Error(err)
			}
			if testSamePackage != testcase.ExpectedBool {
				t.Error("Expected", testcase.ExpectedBool, "got", testSamePackage)
			}
			if testSamePackage == false {
				return
			}
			if testPackageName != testcase.ExpectedValue {
				t.Error("Expected", testcase.ExpectedValue, "got", testPackageName)
			}
		})
	}
}

func TestCommentsPackage(t *testing.T) {
	files, err := os.ReadDir("test_folder/test_various_files/")
	if err != nil {
		log.Fatal(err)
	}
	for _, f := range files {
		if f.Name() == "package_comment.cx" {
			packageName, err := loader.GetPackageName(f, "test_folder/test_various_files/")
			if err != nil {
				t.Error(err)
			}
			if packageName != "main" {
				t.Error("Expected package main, got", packageName)
			}
		}
		if f.Name() == "import_comment.cx" {
			importList, err := loader.GetImports(f, "test_folder/test_various_files/")
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
	tests := []struct {
		Scenario          string
		Database          string
		WantNumberOfFiles int
	}{
		{
			Scenario:          "Test with redis database",
			Database:          "redis",
			WantNumberOfFiles: 2,
		},
		{
			Scenario:          "Test with bolt database",
			Database:          "bolt",
			WantNumberOfFiles: 2,
		},
	}
	for _, testcase := range tests {
		t.Run(testcase.Scenario, func(t *testing.T) {
			testPackage := loader.Package{}
			loader.AddFiles(&testPackage, testFileList, TEST_SRC_PATH, testcase.Database)
			if len(testPackage.Files) != testcase.WantNumberOfFiles {
				t.Error("Expected", testcase.WantNumberOfFiles, " files, got", len(testPackage.Files))
			}
		})
	}
}

func TestAddPackages(t *testing.T) {
	tests := []struct {
		Scenario             string
		Database             string
		WantNumberOfPackages int
	}{
		{
			Scenario:             "Test with redis database",
			Database:             "redis",
			WantNumberOfPackages: 2,
		},
		{
			Scenario:             "Test with bolt database",
			Database:             "bolt",
			WantNumberOfPackages: 2,
		},
	}
	for _, testcase := range tests {
		t.Run(testcase.Scenario, func(t *testing.T) {
			testPackageList := loader.PackageList{}
			_, err := loader.AddPackages(&testPackageList, TEST_SRC_PATH2, TEST_SRC_PATH2, testcase.Database, []string{})
			if err != nil {
				t.Error(err)
			}
			if len(testPackageList.Packages) != testcase.WantNumberOfPackages {
				t.Error("Expected", testcase.WantNumberOfPackages, "packages, got", len(testPackageList.Packages))
			}
		})
	}
}

func TestLoad(t *testing.T) {
	tests := []struct {
		Scenario string
		Database string
	}{
		{
			Scenario: "Test with Redis database",
			Database: "redis",
		},
		{
			Scenario: "Test with Bolt database",
			Database: "bolt",
		},
	}
	for _, testcase := range tests {
		t.Run(testcase.Scenario, func(t *testing.T) {
			loader.LoadPackages("TestValid", "test_folder/test_valid_program/", testcase.Database)
		})
	}
}
