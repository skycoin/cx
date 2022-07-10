package type_checks_test

import (
	"os"
	"testing"

	"github.com/skycoin/cx/cmd/declaration_extraction"
	"github.com/skycoin/cx/cmd/type_checks"
	"github.com/skycoin/cx/cxparser/actions"
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

			type_checks.ParseGlobals(Globals)

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
	}{
		{
			scenario: "Has Structs",
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

			structs, err := declaration_extraction.ExtractStructs(ReplaceCommentsWithWhitespaces, tc.testDir, pkg)

			type_checks.ParseStructs(structs)

		})
	}

}

func TestTypeChecks_ParseFuncHeaders(t *testing.T) {

	tests := []struct {
		scenario string
		testDir  string
	}{
		{
			scenario: "Has funcs",
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

			funcs, err := declaration_extraction.ExtractFuncs(srcBytes, tc.testDir, pkg)

			type_checks.ParseFuncHeaders(funcs)

			actions.AST.PrintProgram()
		})
	}

}
