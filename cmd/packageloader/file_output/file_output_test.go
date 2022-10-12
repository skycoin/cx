package file_output_test

import (
	"testing"

	"github.com/skycoin/cx/cmd/packageloader/file_output"
	"github.com/skycoin/cx/cmd/packageloader/loader"
	"github.com/skycoin/cx/cxparser/actions"
)

func TestFile_Output_GetImportFiles(t *testing.T) {
	tests := []struct {
		scenario    string
		programName string
		testDir     string
		database    string
		files       []*loader.File
	}{
		{
			scenario:    "Has Imports",
			programName: "tester",
			testDir:     "./test_files/test_valid_program",
			database:    "bolt",
			files: []*loader.File{
				{
					FileName: "",
				},
			},
		},
		{
			scenario:    "Has Imports",
			programName: "tester",
			testDir:     "./test_files/test_tree",
			database:    "bolt",
		},
	}

	for _, tc := range tests {
		t.Run(tc.scenario, func(t *testing.T) {

			actions.AST = nil

			_, sourceCodes, _ := loader.ParseArgsForCX([]string{tc.testDir}, true)
			err := loader.LoadCXProgram(tc.programName, sourceCodes, tc.database)
			if err != nil {
				t.Fatal(err)
			}

			err = file_output.AddPkgsToAST(tc.programName, tc.database)
			if err != nil {
				t.Fatal(err)
			}

			files, err := file_output.GetImportFiles(tc.programName, tc.database)
			if err != nil {
				t.Fatal(err)
			}

			t.Error(files)

		})
	}
}

func TestFile_Output_AddPkgsToAST(t *testing.T) {
	tests := []struct {
		scenario    string
		testDir     string
		programName string
		database    string
		packages    []string
	}{
		{
			scenario:    "Has Packages",
			testDir:     "./test_files/test_valid_program",
			programName: "MyPkg",
			database:    "bolt",
			packages: []string{
				"main",
				"testimport",
			},
		},
		{
			scenario:    "Has Packages 2",
			testDir:     "./test_files/test_tree",
			programName: "Pkg2",
			database:    "bolt",
			packages: []string{
				"main",
				"testimport2",
				"testimport2",
				"testimport3",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.programName, func(t *testing.T) {
			actions.AST = nil

			_, sourceCodes, _ := loader.ParseArgsForCX([]string{tc.testDir}, true)
			err := loader.LoadCXProgram(tc.programName, sourceCodes, tc.database)
			if err != nil {
				t.Fatal(err)
			}

			err = file_output.AddPkgsToAST(tc.programName, tc.database)
			if err != nil {
				t.Fatal(err)
			}

			for _, wantPackage := range tc.packages {

				var match bool
				var gotPackageName string

				for _, gotPackageIdx := range actions.AST.Packages {
					gotPackage, err := actions.AST.GetPackageFromArray(gotPackageIdx)
					gotPackageName = gotPackage.Name
					if err != nil {
						t.Fatal(err)
					}

					if gotPackageName == wantPackage {
						match = true
						break
					}
				}

				if !match {
					t.Errorf("want package %s, got %s", wantPackage, gotPackageName)
				}
			}
		})
	}
}
