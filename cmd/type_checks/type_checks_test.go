package type_checks_test

import (
	"os"
	"testing"

	"github.com/skycoin/cx/cmd/declaration_extraction"
	"github.com/skycoin/cx/cmd/type_checks"
)

func TestTypeChecks_ParseGlobals(t *testing.T) {

	tests := []struct {
		scenario string
		testDir  string
	}{
		{
			scenario: "Has globals",
			testDir:  "./test_files/test.cx",
		},
	}

	for _, tc := range tests {
		t.Run(tc.scenario, func(t *testing.T) {

			srcBytes, err := os.ReadFile(tc.testDir)
			if err != nil {
				t.Error(err)
			}

			ReplaceCommentsWithWhitespaces := declaration_extraction.ReplaceCommentsWithWhitespaces(srcBytes)
			pkg := declaration_extraction.ExtractPackages(ReplaceCommentsWithWhitespaces)

			Globals, err := declaration_extraction.ExtractGlobals(ReplaceCommentsWithWhitespaces, tc.testDir, pkg)

			t.Error(Globals)

			type_checks.ParseGlobals(Globals)

			// cxpartialparsing.Program.PrintProgram()

		})
	}

}

func TestTypeChecks_ParseEnums(t *testing.T) {

	tests := []struct {
		scenario string
		testDir  string
	}{}

	for _, tc := range tests {
		t.Run(tc.scenario, func(t *testing.T) {

		})
	}

}

func TestTypeChecks_ParseStructs(t *testing.T) {

	tests := []struct {
		scenario string
		testDir  string
	}{}

	for _, tc := range tests {
		t.Run(tc.scenario, func(t *testing.T) {

		})
	}

}

func TestTypeChecks_ParseFuncs(t *testing.T) {

	tests := []struct {
		scenario string
		testDir  string
	}{}

	for _, tc := range tests {
		t.Run(tc.scenario, func(t *testing.T) {

		})
	}

}
