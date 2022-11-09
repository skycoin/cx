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
			Scenario:  "invalid file map",
			FilesPath: TEST_SRC_PATH_INVALID,
			ExpectedFileMap: map[string][]string{
				"main":       {"testfile.cx", "testfile2.cx"},
				"testimport": {"testimport.cx"},
			},
		},
	}

	for _, testcase := range tests {
		t.Run(testcase.Scenario, func(t *testing.T) {
			_, sourceCodes, _ := loader.ParseArgsForCX([]string{TEST_SRC_PATH_VALID}, true)
			gotFileMap, gotErr := loader.CreateFileMap(sourceCodes)
			for wantKey, wantValue := range testcase.ExpectedFileMap {
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

func TestCreateImportMap(t *testing.T) {
	tests := []struct {
		Scenario          string
		FilesPath         string
		ExpectedImportMap map[string][]string
		ExpectedErr       error
	}{
		{
			Scenario:  "valid import map",
			FilesPath: TEST_SRC_PATH_VALID,
			ExpectedImportMap: map[string][]string{
				"main":       {"testimport"},
				"testimport": {},
			},
		},
	}
	for _, testcase := range tests {
		t.Run(testcase.Scenario, func(t *testing.T) {
			_, sourceCodes, _ := loader.ParseArgsForCX([]string{TEST_SRC_PATH_VALID}, true)
			ExpectedFileMap, ExpectedErr := loader.CreateFileMap(sourceCodes)
			if ExpectedErr != nil {
				t.Error(ExpectedErr)
			}
			gotImportMap, gotErr := loader.CreateImportMap(ExpectedFileMap)

			for wantKey, wantValue := range testcase.ExpectedImportMap {
				gotValue, ok := gotImportMap[wantKey]
				if !ok {
					t.Fatalf("package %s not found in file map", wantKey)
				}
				for _, wantImport := range wantValue {
					var match bool
					var gotImportName string
					for _, gotImport := range gotValue {
						gotImportName = gotImport
						if wantImport == gotImportName {
							match = true
							break
						}
					}
					if !match {
						t.Errorf("want import %s, got %s", wantImport, gotImportName)
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
