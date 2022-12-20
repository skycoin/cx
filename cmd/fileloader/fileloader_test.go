package fileloader_test

import (
	"testing"

	"github.com/skycoin/cx/cmd/fileloader"
	"github.com/skycoin/cx/cx/ast"
	cxinit "github.com/skycoin/cx/cx/init"
	"github.com/skycoin/cx/cxparser/actions"
)

func TestFileLoader_LoadFiles(t *testing.T) {
	tests := []struct {
		scenario  string
		testDir   string
		wantFiles []string
	}{
		{
			scenario: "One file one package",
			testDir:  "./test_files/One_file_one_package.cx",
			wantFiles: []string{
				"test_files/One_file_one_package.cx",
			},
		},
		{
			scenario: "Has multiple packages in one file",
			testDir:  "./test_files/Has_multiple_packages_in_file.cx",
			wantFiles: []string{
				"test_files/Has_multiple_packages_in_file.cx",
			},
		},
		{
			scenario: "Has imports",
			testDir:  "./test_files/Has_Imports/",
			wantFiles: []string{
				"test_files/Has_Imports/testimport1/testimport1file1.cx",
				"test_files/Has_Imports/testimport1/testimport1file2.cx",
				"test_files/Has_Imports/testimport2/testimport2file.cx",
				"test_files/Has_Imports/testimport3/testimport3file1.cx",
				"test_files/Has_Imports/testimport4/testimport1file1.cx",
				"test_files/Has_Imports/testimport4/testimport1file2.cx",
				"test_files/Has_Imports/testmain.cx",
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.scenario, func(t *testing.T) {
			actions.AST = cxinit.MakeProgram()

			_, sourceCodes, _ := ast.ParseArgsForCX([]string{tc.testDir}, true)

			_, fileNames, err := fileloader.LoadFiles(sourceCodes)
			if err != nil {
				t.Fatal(err)
			}

			for _, wantFile := range tc.wantFiles {
				var match bool
				for _, file := range fileNames {
					if file == wantFile {
						match = true
						break
					}
				}
				if !match {
					t.Errorf("missing file: %s", wantFile)
				}
			}
		})
	}
}

func BenchmarkFileLoader_LoadFiles(b *testing.B) {
	for i := 0; i < b.N; i++ {
		actions.AST = cxinit.MakeProgram()

		_, sourceCodes, _ := ast.ParseArgsForCX([]string{"./test_files/One_file_one_package.cx"}, true)

		_, _, err := fileloader.LoadFiles(sourceCodes)
		if err != nil {
			b.Fatal(err)
		}
	}
}
