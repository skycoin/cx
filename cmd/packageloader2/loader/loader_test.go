package loader_test

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/skycoin/cx/cmd/packageloader2/loader"
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
	tests := []struct {
		Scenario            string
		FilesPath           string
		ExpectedPackageName string
		ExpectedErr         error
	}{
		{
			Scenario:            "Test with file that has package name",
			FilesPath:           TEST_PACKAGE_FILE,
			ExpectedPackageName: "main",
			ExpectedErr:         nil,
		},
		{
			Scenario:            "Test with file that has no package name",
			FilesPath:           "test_folder/test_various_files/nopackage.cx",
			ExpectedPackageName: "",
			ExpectedErr:         errors.New("file doesn't contain a package name"),
		},
	}
	for _, testcase := range tests {
		t.Run(testcase.Scenario, func(t *testing.T) {
			file, err := os.Open(testcase.FilesPath)
			if err != nil {
				t.Error(err)
			}

			gotPackageName, gotErr := loader.GetPackageName(file)

			if gotPackageName != testcase.ExpectedPackageName {
				t.Errorf("want package name %s, got %s", testcase.ExpectedPackageName, gotPackageName)
			}

			if (gotErr != nil && testcase.ExpectedErr == nil) ||
				(gotErr == nil && testcase.ExpectedErr != nil) {
				t.Errorf("want error %v, got %v", testcase.ExpectedErr, gotErr)
			}
			if gotErr != nil && testcase.ExpectedErr != nil {
				if gotErr.Error() != testcase.ExpectedErr.Error() {
					t.Errorf("want error %v, got %v", testcase.ExpectedErr, gotErr)
				}
			}
		})
	}

}

func TestCreateFileMap(t *testing.T) {
	tests := []struct {
		Scenario        string
		FilesPath       string
		ExpectedFileMap map[string][]string
		ExpectedErr     error
	}{
		{
			Scenario:  "valid file map",
			FilesPath: TEST_SRC_PATH_VALID,
			ExpectedFileMap: map[string][]string{
				"main":       {"testfile.cx", "testfile2.cx"},
				"testimport": {"testimport.cx"},
			},
		},
		{
			Scenario:  "error in main package",
			FilesPath: TEST_SRC_PATH_INVALID,
			ExpectedFileMap: map[string][]string{
				"main": {"testfile.cx"},
			},
			ExpectedErr: errors.New("testfile2.cx: package error: package main2 found in main"),
		},
		{
			Scenario:        "multiple packages in directory",
			FilesPath:       "test_folder/test_package_error_program",
			ExpectedFileMap: map[string][]string{},
			ExpectedErr:     errors.New("testimportFile2.cx: package error: package testimport2 found in testimport"),
		},
	}

	for _, testcase := range tests {
		t.Run(testcase.Scenario, func(t *testing.T) {
			_, sourceCodes, _ := loader.ParseArgsForCX([]string{testcase.FilesPath}, true)
			gotFileMap, gotErr := loader.CreateFileMap(sourceCodes)
			for wantKey, wantValue := range testcase.ExpectedFileMap {
				gotValue, ok := gotFileMap[wantKey]
				if !ok {
					t.Errorf("package %s not found in file map", wantKey)
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
						t.Errorf("want file %s, got %s", wantFile, gotFileName)
					}
				}
			}
			if (gotErr != nil && testcase.ExpectedErr == nil) ||
				(gotErr == nil && testcase.ExpectedErr != nil) {
				t.Errorf("want error %v, got %v", testcase.ExpectedErr, gotErr)
			}

			if gotErr != nil && testcase.ExpectedErr != nil {
				if gotErr.Error() != testcase.ExpectedErr.Error() {
					t.Errorf("want error %v, got %v", testcase.ExpectedErr, gotErr)
				}
			}
		})
	}
}

// func TestCheckForDependencyLoop(t *testing.T) {
// 	tests := []struct {
// 		Scenario   string
// 		FilesPath  string
// 		ExpectsErr bool
// 	}{
// 		{
// 			Scenario:   "Test with a program that doesn't contain a dependency loop",
// 			FilesPath:  TEST_SRC_PATH_VALID,
// 			ExpectsErr: false,
// 		},
// 		{
// 			Scenario:   "Test with a program that contains a dependency loop",
// 			FilesPath:  TEST_SRC_PATH_LOOP,
// 			ExpectsErr: true,
// 		},
// 	}

// 	for _, testcase := range tests {
// 		t.Run(testcase.Scenario, func(t *testing.T) {
// 			_, sourceCodes, _ := loader.ParseArgsForCX([]string{testcase.FilesPath}, true)
// 			fileMap, err := loader.CreateFileMap(sourceCodes)
// 			if err != nil {
// 				t.Error(err)
// 			}
// 			importMap, err := loader.CreateImportMap(fileMap)
// 			if err != nil {
// 				t.Error(err)
// 			}
// 			err = loader.CheckForDependencyLoop(importMap)
// 			if (err != nil) != testcase.ExpectsErr {
// 				t.Error("Dependency check failed")
// 			}
// 		})
// 	}
// }

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

// func TestAddNewPackage(t *testing.T) {
// 	tests := []struct {
// 		Scenario  string
// 		FilesPath string
// 		Database  string
// 	}{
// 		{
// 			Scenario:  "Test adding package to Redis database",
// 			FilesPath: TEST_SRC_PATH_VALID,
// 			Database:  "redis",
// 		},
// 		{
// 			Scenario:  "Test adding package to Bolt database",
// 			FilesPath: TEST_SRC_PATH_VALID,
// 			Database:  "bolt",
// 		},
// 	}
// 	for _, testcase := range tests {
// 		t.Run(testcase.Scenario, func(t *testing.T) {
// 			_, sourceCodes, _ := loader.ParseArgsForCX([]string{testcase.FilesPath}, false)
// 			testPackageStruct := loader.PackageList{}
// 			err := loader.AddNewPackage(&testPackageStruct, "main", sourceCodes, testcase.Database)
// 			if err != nil {
// 				t.Error(err)
// 			}
// 			if len(testPackageStruct.Packages) != 1 {
// 				t.Error("Wrong number of packages added")
// 			}
// 		})
// 	}
// }

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
			_, sourceCodes, _ := loader.ParseArgsForCX([]string{testcase.FilesPath}, true)
			err := loader.LoadCXProgram("test", sourceCodes, testcase.Database)
			if err != nil {
				t.Error(err)
			}
		})
	}
}
