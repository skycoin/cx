package file_output_test

import (
	"testing"

	"github.com/skycoin/cx/cmd/packageloader2/file_output"
	"github.com/skycoin/cx/cmd/packageloader2/loader"
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
		},
		{
			scenario:    "Has Imports",
			programName: "tester",
			testDir:     "./test_files/test_tree",
			database:    "redis",
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

		})
	}
}
