package loader_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/skycoin/cx/cmd/packageloader/loader"
)

const TEST_SRC_PATH_INVALID = "test_folder/test_invalid_program/src/"
const TEST_SRC_PATH_VALID = "test_folder/test_valid_program/src/"
const TEST_SRC_PATH_LOOP = "test_folder/test_loop_program/src/"
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
	tests := []struct {
		scenario string
		testDir  string
		fileMap  map[string][]string
	}{
		{
			scenario: "valid file map",
			testDir:  TEST_SRC_PATH_VALID,
			fileMap: map[string][]string{
				"main":       {"testfile.cx", "testfile2.cx"},
				"testimport": {"testimport.cx"},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.scenario, func(t *testing.T) {
			_, sourceCodes, _ := loader.ParseArgsForCX([]string{TEST_SRC_PATH_VALID}, true)
			gotFileMap, err := loader.CreateFileMap(sourceCodes)
			if err != nil {
				t.Error(err)
			}
			for wantKey, wantValue := range tc.fileMap {
				gotValue, ok := gotFileMap[wantKey]
				if !ok {
					t.Fatalf("package %s not found in file map", wantKey)
				}
				for _, wantFile := range wantValue {
					var match bool
					var gotFileName string
					for _, gotFile := range gotValue {
						gotFileName = filepath.Base(gotFile.Name())
						if wantFile == gotFileName {
							match = true
							break
						}
					}
					if !match {
						t.Errorf("want %s, got %s", wantFile, gotFileName)
					}
				}
			}
		})
	}
}

func TestCreateImportMap(t *testing.T) {
	_, sourceCodes, _ := loader.ParseArgsForCX([]string{TEST_SRC_PATH_VALID}, true)
	fileMap, err := loader.CreateFileMap(sourceCodes)
	if err != nil {
		t.Error(err)
	}
	_, err = loader.CreateImportMap(fileMap)
	if err != nil {
		t.Error(err)
	}
	//TODO: Find a way to reliably test this function
}

func TestCheckForDependencyLoop(t *testing.T) {
	tests := []struct {
		Scenario   string
		FilesPath  string
		ExpectsErr bool
	}{
		{
			Scenario:   "Test with a program that doesn't contain a dependency loop",
			FilesPath:  TEST_SRC_PATH_VALID,
			ExpectsErr: false,
		},
		{
			Scenario:   "Test with a program that contains a dependency loop",
			FilesPath:  TEST_SRC_PATH_LOOP,
			ExpectsErr: true,
		},
	}

	for _, testcase := range tests {
		t.Run(testcase.Scenario, func(t *testing.T) {
			_, sourceCodes, _ := loader.ParseArgsForCX([]string{testcase.FilesPath}, true)
			fileMap, err := loader.CreateFileMap(sourceCodes)
			if err != nil {
				t.Error(err)
			}
			importMap, err := loader.CreateImportMap(fileMap)
			if err != nil {
				t.Error(err)
			}
			err = loader.CheckForDependencyLoop(importMap)
			if (err != nil) != testcase.ExpectsErr {
				t.Error("Dependency check failed")
			}
		})
	}
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

func TestAddNewPackage(t *testing.T) {
	tests := []struct {
		Scenario  string
		FilesPath string
		Database  string
	}{
		{
			Scenario:  "Test adding package to Redis database",
			FilesPath: TEST_SRC_PATH_VALID,
			Database:  "redis",
		},
		{
			Scenario:  "Test adding package to Bolt database",
			FilesPath: TEST_SRC_PATH_VALID,
			Database:  "bolt",
		},
	}
	for _, testcase := range tests {
		t.Run(testcase.Scenario, func(t *testing.T) {
			_, sourceCodes, _ := loader.ParseArgsForCX([]string{testcase.FilesPath}, false)
			testPackageStruct := loader.PackageList{}
			err := loader.AddNewPackage(&testPackageStruct, "main", sourceCodes, testcase.Database)
			if err != nil {
				t.Error(err)
			}
			if len(testPackageStruct.Packages) != 1 {
				t.Error("Wrong number of packages added")
			}
		})
	}
}

func TestLoadCXProgram(t *testing.T) {
	tests := []struct {
		Scenario         string
		FilesPath        string
		Database         string
		ExpectedPackages int
	}{
		{
			Scenario:         "Test adding package to Redis database",
			FilesPath:        TEST_SRC_PATH_VALID,
			Database:         "redis",
			ExpectedPackages: 2,
		},
		{
			Scenario:         "Test adding package to Bolt database",
			FilesPath:        TEST_SRC_PATH_VALID,
			Database:         "bolt",
			ExpectedPackages: 2,
		},
	}
	for _, testcase := range tests {
		t.Run(testcase.Scenario, func(t *testing.T) {
			_, sourceCodes, _ := loader.ParseArgsForCX([]string{testcase.FilesPath}, false)
			err := loader.LoadCXProgram("test", sourceCodes, testcase.Database)
			if err != nil {
				t.Error(err)
			}
		})
	}
}
