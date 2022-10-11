package file_output_test

import (
	"testing"

	"github.com/skycoin/cx/cmd/packageloader/file_output"
	"github.com/skycoin/cx/cmd/packageloader/loader"
)

func TestFile_Output_GetImportFiles(t *testing.T) {
	tests := []struct {
		scenario    string
		programName string
		testDir     string
		database    string
	}{
		// {
		// 	scenario:    "Has Imports",
		// 	programName: "tester",
		// 	testDir:     "./test_files/test_valid_program",
		// 	database:    "bolt",
		// },
		{
			scenario:    "Has Imports",
			programName: "tester",
			testDir:     "./test_files/test_tree",
			database:    "bolt",
		},
	}

	for _, tc := range tests {
		t.Run(tc.scenario, func(t *testing.T) {
			_, sourceCodes, _ := loader.ParseArgsForCX([]string{tc.testDir}, true)
			err := loader.LoadCXProgram(tc.programName, sourceCodes, tc.database)
			if err != nil {
				t.Fatal(err)
			}

			files, err := file_output.GetImportFiles(tc.programName, tc.database)
			if err != nil {
				t.Fatal(err)
			}

			t.Error(len(files))
		})
	}
}
