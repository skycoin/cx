package loader_test

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
	"testing"

	"github.com/skycoin/cx/cmd/packageloader2/loader"
)

const TEST_SRC_PATH_INVALID = "test_folder/test_invalid_program/src/"
const TEST_SRC_PATH_VALID = "test_folder/test_valid_program/src/"
const TEST_SRC_PATH_LOOP = "test_folder/test_loop_program/src/"
const TEST_PACKAGE_FILE = "test_folder/test_various_files/package.cx"

func TestContains(t *testing.T) {
	tests := []struct {
		Scenario string
		List     []string
		Element  string
		WantBool bool
	}{
		{
			Scenario: "Testing if list does contain the element",
			List:     []string{"a", "b", "c"},
			Element:  "a",
			WantBool: true,
		},
		{
			Scenario: "Testing if list doesn't contain the element",
			List:     []string{"a", "b", "c"},
			Element:  "d",
			WantBool: false,
		},
	}
	for _, testcase := range tests {
		t.Run(testcase.Scenario, func(t *testing.T) {
			gotBool := loader.Contains(testcase.List, testcase.Element)
			if gotBool != testcase.WantBool {
				t.Errorf("want bool %v, got %v", gotBool, testcase.WantBool)
			}
		})
	}
}

func TestRemoveDuplicates(t *testing.T) {
	tests := []struct {
		Scenario string
		List     []string
		WantList []string
	}{
		{
			Scenario: "Testing list without duplicates",
			List:     []string{"a", "b", "c"},
			WantList: []string{"a", "b", "c"},
		},
		{
			Scenario: "Testing list with duplicates",
			List:     []string{"a", "b", "b", "c", "c"},
			WantList: []string{"a", "b", "c"},
		},
	}
	for _, testcase := range tests {
		t.Run(testcase.Scenario, func(t *testing.T) {
			gotList := loader.RemoveDuplicates(testcase.List)
			if !reflect.DeepEqual(gotList, testcase.WantList) {
				t.Errorf("want list %v, got %v", testcase.WantList, gotList)
			}
		})
	}
}

func TestGetPackageName(t *testing.T) {
	tests := []struct {
		Scenario        string
		FilesPath       string
		WantPackageName string
		WantErr         error
	}{
		{
			Scenario:        "Test with file that has package name",
			FilesPath:       TEST_PACKAGE_FILE,
			WantPackageName: "main",
			WantErr:         nil,
		},
		{
			Scenario:        "Test with file that has no package name",
			FilesPath:       "test_folder/test_various_files/nopackage.cx",
			WantPackageName: "",
			WantErr:         errors.New("file doesn't contain a package name"),
		},
	}
	for _, testcase := range tests {
		t.Run(testcase.Scenario, func(t *testing.T) {
			file, err := os.Open(testcase.FilesPath)
			if err != nil {
				t.Error(err)
			}

			gotPackageName, gotErr := loader.GetPackageName(file)

			if gotPackageName != testcase.WantPackageName {
				t.Errorf("want package name %s, got %s", testcase.WantPackageName, gotPackageName)
			}

			if (gotErr != nil && testcase.WantErr == nil) ||
				(gotErr == nil && testcase.WantErr != nil) {
				t.Errorf("want error %v, got %v", testcase.WantErr, gotErr)
			}
			if gotErr != nil && testcase.WantErr != nil {
				if gotErr.Error() != testcase.WantErr.Error() {
					t.Errorf("want error %v, got %v", testcase.WantErr, gotErr)
				}
			}
		})
	}

}

func TestCreateFileMap(t *testing.T) {
	tests := []struct {
		Scenario    string
		FilesPath   string
		WantFileMap map[string][]string
		WantErr     error
	}{
		{
			Scenario:  "valid file map",
			FilesPath: TEST_SRC_PATH_VALID,
			WantFileMap: map[string][]string{
				"main":       {"testfile.cx", "testfile2.cx"},
				"testimport": {"testimport.cx"},
			},
		},
		{
			Scenario:    "no package name",
			FilesPath:   "test_folder/test_various_files/nopackage.cx",
			WantFileMap: map[string][]string{},
			WantErr:     errors.New("file doesn't contain a package name"),
		},
		{
			Scenario:  "error in main package",
			FilesPath: TEST_SRC_PATH_INVALID,
			WantFileMap: map[string][]string{
				"main": {"testfile.cx"},
			},
			WantErr: errors.New("testfile2.cx: package error: package main2 found in main"),
		},
		{
			Scenario:    "multiple packages in directory",
			FilesPath:   "test_folder/test_package_error_program",
			WantFileMap: map[string][]string{},
			WantErr:     errors.New("testimportFile2.cx: package error: package testimport2 found in testimport"),
		},
	}

	for _, testcase := range tests {
		t.Run(testcase.Scenario, func(t *testing.T) {
			_, sourceCodes, _, rootDir := loader.ParseArgsForCX([]string{testcase.FilesPath}, true)
			gotFileMap, gotErr := loader.CreateFileMap(sourceCodes, rootDir)
			for wantKey, wantValue := range testcase.WantFileMap {
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
			if (gotErr != nil && testcase.WantErr == nil) ||
				(gotErr == nil && testcase.WantErr != nil) {
				t.Errorf("want error %v, got %v", testcase.WantErr, gotErr)
			}

			if gotErr != nil && testcase.WantErr != nil {
				if gotErr.Error() != testcase.WantErr.Error() {
					t.Errorf("want error %v, got %v", testcase.WantErr, gotErr)
				}
			}
		})
	}
}

func TestFileStructFromFile(t *testing.T) {
	tests := []struct {
		Scenario       string
		FilesPath      string
		WantFileStruct loader.File
		WantErr        error
	}{
		{
			Scenario:  "Testing package file with no errors",
			FilesPath: TEST_PACKAGE_FILE,
			WantFileStruct: loader.File{
				FileName: "package.cx",
			},
			WantErr: nil,
		},
	}
	for _, testcase := range tests {
		t.Run(testcase.Scenario, func(t *testing.T) {
			file, err := os.Open(testcase.FilesPath)
			if err != nil {
				t.Error(err)
			}
			gotFileStruct, gotErr := loader.FileStructFromFile(file)

			if gotFileStruct.FileName != testcase.WantFileStruct.FileName {
				t.Errorf("want file %s, got %s", testcase.WantFileStruct.FileName, gotFileStruct.FileName)
			}

			if (gotErr != nil && testcase.WantErr == nil) ||
				(gotErr == nil && testcase.WantErr != nil) {
				t.Errorf("want error %v, got %v", testcase.WantErr, gotErr)
			}

			if gotErr != nil && testcase.WantErr != nil {
				if gotErr.Error() != testcase.WantErr.Error() {
					t.Errorf("want error %v, got %v", testcase.WantErr, gotErr)
				}
			}
		})
	}

}

func TestAddNewPackage(t *testing.T) {
	tests := []struct {
		Scenario             string
		FilesPath            string
		Database             string
		WantNumberOfPackages int
		WantErr              error
	}{
		{
			Scenario:             "Test adding package to Redis database",
			FilesPath:            TEST_SRC_PATH_VALID,
			Database:             "redis",
			WantNumberOfPackages: 1,

			WantErr: nil,
		},
		{
			Scenario:             "Test adding package to Bolt database",
			FilesPath:            TEST_SRC_PATH_VALID,
			Database:             "bolt",
			WantNumberOfPackages: 1,

			WantErr: nil,
		},
	}
	for _, testcase := range tests {
		t.Run(testcase.Scenario, func(t *testing.T) {
			_, sourceCodes, _, rootDir := loader.ParseArgsForCX([]string{testcase.FilesPath}, false)
			fileMap, err := loader.CreateFileMap(sourceCodes, rootDir)
			if err != nil {
				t.Error(err)
			}

			testPackageStruct := loader.PackageList{}

			files, ok := fileMap["main"]
			if !ok {
				t.Error("package main not found")
			}

			gotErr := loader.AddNewPackage(&testPackageStruct, "main", files, testcase.Database)

			gotNumberOfPackages := len(testPackageStruct.Packages)
			if gotNumberOfPackages != testcase.WantNumberOfPackages {
				t.Errorf("want %d packages, got %d packages", testcase.WantNumberOfPackages, gotNumberOfPackages)
			}

			if (gotErr != nil && testcase.WantErr == nil) ||
				(gotErr == nil && testcase.WantErr != nil) {
				t.Errorf("want error %v, got %v", testcase.WantErr, gotErr)
			}
			if gotErr != nil && testcase.WantErr != nil {
				if gotErr.Error() != testcase.WantErr.Error() {
					t.Errorf("want error %v, got %v", testcase.WantErr, gotErr)
				}
			}
		})
	}
}

func TestLoadImportPackages(t *testing.T) {
	tests := []struct {
		Scenario             string
		FilesPath            string
		Database             string
		WantNumberOfPackages int
		WantImportMap        map[string][]string
		WantErr              error
	}{
		{
			Scenario:             "Test adding package to Redis database",
			FilesPath:            TEST_SRC_PATH_VALID,
			Database:             "redis",
			WantNumberOfPackages: 2,
			WantImportMap: map[string][]string{
				"main":       {"os", "testimport"},
				"testimport": {},
			},
			WantErr: nil,
		},
		{
			Scenario:             "Test adding package to Bolt database",
			FilesPath:            TEST_SRC_PATH_VALID,
			Database:             "bolt",
			WantNumberOfPackages: 2,
			WantImportMap: map[string][]string{
				"main":       {"os", "testimport"},
				"testimport": {},
			},
			WantErr: nil,
		},
	}

	for _, testcase := range tests {
		t.Run(testcase.Scenario, func(t *testing.T) {
			_, sourceCodes, _, rootDir := loader.ParseArgsForCX([]string{testcase.FilesPath}, true)
			fileMap, err := loader.CreateFileMap(sourceCodes, rootDir)
			if err != nil {
				t.Error(err)
			}
			var mx sync.Mutex
			gotImportMap := make(map[string][]string)
			testPackageStruct := loader.PackageList{}

			files, ok := fileMap["main"]
			if !ok {
				t.Error("package main not found")
			}

			err = loader.AddNewPackage(&testPackageStruct, "main", files, testcase.Database)
			if err != nil {
				t.Error(err)
			}

			err = loader.CheckImports("main", files, gotImportMap, &mx)
			if err != nil {
				t.Error(err)
			}

			gotErr := loader.LoadImportPackages(&testPackageStruct, "main", fileMap, gotImportMap, testcase.Database, &mx)

			gotNumberOfPackages := len(testPackageStruct.Packages)
			if gotNumberOfPackages != testcase.WantNumberOfPackages {
				t.Errorf("want %d packages, got %d packages", testcase.WantNumberOfPackages, gotNumberOfPackages)
			}

			if fmt.Sprint(gotImportMap) != fmt.Sprint(testcase.WantImportMap) {
				t.Errorf("want import map %v, got %v", testcase.WantImportMap, gotImportMap)
			}

			if (gotErr != nil && testcase.WantErr == nil) ||
				(gotErr == nil && testcase.WantErr != nil) {
				t.Errorf("want error %v, got %v", testcase.WantErr, gotErr)
			}
			if gotErr != nil && testcase.WantErr != nil {
				if gotErr.Error() != testcase.WantErr.Error() {
					t.Errorf("want error %v, got %v", testcase.WantErr, gotErr)
				}
			}
		})
	}
}

func TestCheckForDependencyLoop(t *testing.T) {
	tests := []struct {
		Scenario  string
		ImportMap map[string][]string
		WantErr   error
	}{
		{
			Scenario: "Testing with a program with no dependency loop",
			ImportMap: map[string][]string{
				"main": {"testimport"},
			},
			WantErr: nil,
		},
		{
			Scenario: "Test with a program that contains a self dependency loop",
			ImportMap: map[string][]string{
				"main": {"main"},
			},
			WantErr: errors.New("Module main imports itself"),
		},
		{
			Scenario: "Test with a program that contains a dependency loop between modules",
			ImportMap: map[string][]string{
				"main":       {"testimport"},
				"testimport": {"main"},
			},
			WantErr: errors.New("Dependency loop between modules main and testimport"),
		},
	}

	for _, testcase := range tests {
		t.Run(testcase.Scenario, func(t *testing.T) {
			gotErr := loader.CheckForDependencyLoop(testcase.ImportMap, "main")

			if (gotErr != nil && testcase.WantErr == nil) ||
				(gotErr == nil && testcase.WantErr != nil) {
				t.Errorf("want error %v, got %v", testcase.WantErr, gotErr)
			}

			if gotErr != nil && testcase.WantErr != nil {
				if gotErr.Error() != testcase.WantErr.Error() &&
					!(strings.Contains(gotErr.Error(), "between") && strings.Contains(testcase.WantErr.Error(), "between")) {
					t.Errorf("want error %v, got %v", testcase.WantErr, gotErr)
				}
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
			_, sourceCodes, _, rootDir := loader.ParseArgsForCX([]string{testcase.FilesPath}, true)
			err := loader.LoadCXProgram("test", sourceCodes, rootDir, testcase.Database)
			if err != nil {
				t.Error(err)
			}
		})
	}
}
