package file_output_test

import (
	"testing"

	"github.com/skycoin/cx/cmd/packageloader/file_output"
	"github.com/skycoin/cx/cmd/packageloader/loader"
	"github.com/skycoin/cx/cxparser/actions"
)

func TestFile_Output_GetImportFiles(t *testing.T) {

	type pkg struct {
		pkgName string
		imports []string
	}
	tests := []struct {
		scenario    string
		programName string
		testDir     string
		database    string
		files       []loader.File
		packages    []pkg
	}{
		{
			scenario:    "Has Imports",
			programName: "tester",
			testDir:     "./test_files/test_valid_program",
			database:    "bolt",
			files: []loader.File{
				{
					FileName: "testimport.cx",
				},
				{
					FileName: "testfile.cx",
				},
				{
					FileName: "testfile.cx",
				},
			},
			packages: []pkg{
				{
					pkgName: "main",
					imports: []string{
						"os",
						"testimport1",
						"testimport2",
					},
				},
			},
		},
		{
			scenario:    "Has Imports",
			programName: "tester",
			testDir:     "./test_files/test_tree",
			database:    "bolt",
			files: []loader.File{
				{
					FileName: "testimport1file1.cx",
				},
				{
					FileName: "testimport1file2.cx",
				},
				{
					FileName: "testimport2file.cx",
				},
				{
					FileName: "testimport3file1.cx",
				},
				{
					FileName: "testmain.cx",
				},
			},
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

			for _, wantFile := range tc.files {
				var match bool
				var wantFileName string = wantFile.FileName
				var gotFileName string
				for _, gotFile := range files {
					gotFileName = gotFile.FileName
					if gotFileName == wantFileName {
						match = true
						break
					}
				}

				if !match {
					t.Errorf("want file %s, got %s", wantFileName, gotFileName)
				}
			}

			for _, wantPackage := range tc.packages {
				var match bool
				var wantPackageName string = wantPackage.pkgName
				var gotPackageName string

				for _, gotPackageIdx := range actions.AST.Packages {
					gotPackage, err := actions.AST.GetPackageFromArray(gotPackageIdx)
					if err != nil {
						t.Fatal(err)
					}
					gotPackageName = gotPackage.Name
					if gotPackageName == wantPackageName {
						match = true

						for _, wantImport := range wantPackage.imports {
							var impMatch bool
							var gotImpName string
							for _, gotImportIdx := range gotPackage.Imports {
								impPkg, err := prgrm.GetPackageFromArray(impIdx)
								if err != nil {
									t.Fatal(err)
								}
								*ast += fmt.Sprintf("\t\t%d.- Import: %s\n", count, impPkg.Name)								if gotImpor == wantImport {
								
							}
						}
	
						break
					}
				}

				if !match {
					t.Errorf("want package %s, got %s", wantPackageName, gotPackageName)
				}
			}

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
