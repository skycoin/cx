package declaration_extraction_test

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/skycoin/cx/cmd/declaration_extraction"
)

func TestDeclarationExtraction_RemoveComment(t *testing.T) {

	tests := []struct {
		scenario           string
		testDir            string
		wantCommentRemoved string
	}{
		{
			scenario:           "Has comments",
			testDir:            "./test_files/test.cx",
			wantCommentRemoved: "./test_files/removeCommentResult.cx",
		},
	}

	for _, tc := range tests {
		t.Run(tc.scenario, func(t *testing.T) {
			srcBytes, err := os.ReadFile(tc.testDir)
			if err != nil {
				t.Fatal(err)
			}
			wantBytes, err := os.ReadFile(tc.wantCommentRemoved)
			if err != nil {
				t.Fatal(err)
			}
			commentRemoved := declaration_extraction.RemoveComment(srcBytes)

			if string(commentRemoved) != string(wantBytes) {
				t.Errorf("want removed comments %v, got %v", string(wantBytes), string(commentRemoved))
			}
		})
	}
}

func TestDeclarationExtraction_ExtractPackages(t *testing.T) {

	tests := []struct {
		scenario    string
		testDir     string
		wantPackage string
	}{
		{
			scenario:    "Has package",
			testDir:     "./test_files/test.cx",
			wantPackage: "hello",
		},
	}

	for _, tc := range tests {
		t.Run(tc.scenario, func(t *testing.T) {
			srcBytes, err := os.ReadFile(tc.testDir)
			if err != nil {
				t.Fatal(err)
			}
			pkg := declaration_extraction.ExtractPackages(srcBytes)

			if pkg != tc.wantPackage {
				t.Errorf("want packages %v, got %v", tc.wantPackage, pkg)
			}

		})
	}
}

func TestDeclarationExtraction_ExtractGlobal(t *testing.T) {

	tests := []struct {
		scenario    string
		testDir     string
		wantGlobals []declaration_extraction.GlobalDeclaration
	}{
		{
			scenario: "Has globals",
			testDir:  "./test_files/test.cx",
			wantGlobals: []declaration_extraction.GlobalDeclaration{
				{
					PackageID:   "hello",
					FileID:      "test.cx",
					StartOffset: 0,
					Length:      16,
					Name:        "apple",
				},
				{
					PackageID:   "hello",
					FileID:      "test.cx",
					StartOffset: 2,
					Length:      17,
					Name:        "banana",
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.scenario, func(t *testing.T) {
			srcBytes, err := os.ReadFile(tc.testDir)
			removeComment := declaration_extraction.RemoveComment(srcBytes)
			fileName := filepath.Base(tc.testDir)
			pkg := declaration_extraction.ExtractPackages(removeComment)
			if err != nil {
				t.Fatal(err)
			}
			globals := declaration_extraction.ExtractGlobals(removeComment, fileName, pkg)

			for i := range globals {
				if globals[i] != tc.wantGlobals[i] {
					t.Errorf("want globals %v, got %v", tc.wantGlobals[i], globals[i])
				}
			}
		})
	}
}

func TestDeclarationExtraction_ExtractEnums(t *testing.T) {

	tests := []struct {
		scenario  string
		testDir   string
		wantEnums []declaration_extraction.EnumDeclaration
	}{
		{
			scenario: "Has enums",
			testDir:  "./test_files/test.cx",
			wantEnums: []declaration_extraction.EnumDeclaration{
				{
					PackageID:   "hello",
					FileID:      "test.cx",
					StartOffset: 1,
					Length:      15,
					Name:        "North",
					Type:        "Direction",
					Value:       0,
				},
				{
					PackageID:   "hello",
					FileID:      "test.cx",
					StartOffset: 1,
					Length:      5,
					Name:        "South",
					Type:        "Direction",
					Value:       1,
				},
				{
					PackageID:   "hello",
					FileID:      "test.cx",
					StartOffset: 1,
					Length:      4,
					Name:        "East",
					Type:        "Direction",
					Value:       2,
				},
				{
					PackageID:   "hello",
					FileID:      "test.cx",
					StartOffset: 1,
					Length:      4,
					Name:        "West",
					Type:        "Direction",
					Value:       3,
				},
				{
					PackageID:   "hello",
					FileID:      "test.cx",
					StartOffset: 1,
					Length:      12,
					Name:        "First",
					Type:        "Number",
					Value:       0,
				},
				{
					PackageID:   "hello",
					FileID:      "test.cx",
					StartOffset: 1,
					Length:      6,
					Name:        "Second",
					Type:        "Number",
					Value:       1,
				},
			},
		},
		{
			scenario: "Has enums and nested parenthesis",
			testDir:  "./test_files/enum_in_parenthesis.cx",
			wantEnums: []declaration_extraction.EnumDeclaration{
				{
					PackageID:   "hello",
					FileID:      "enum_in_parenthesis.cx",
					StartOffset: 1,
					Length:      15,
					Name:        "North",
					Type:        "Direction",
					Value:       0,
				},
				{
					PackageID:   "hello",
					FileID:      "enum_in_parenthesis.cx",
					StartOffset: 1,
					Length:      5,
					Name:        "South",
					Type:        "Direction",
					Value:       1,
				},
				{
					PackageID:   "hello",
					FileID:      "enum_in_parenthesis.cx",
					StartOffset: 1,
					Length:      4,
					Name:        "East",
					Type:        "Direction",
					Value:       2,
				},
				{
					PackageID:   "hello",
					FileID:      "enum_in_parenthesis.cx",
					StartOffset: 1,
					Length:      4,
					Name:        "West",
					Type:        "Direction",
					Value:       3,
				},
				{
					PackageID:   "hello",
					FileID:      "enum_in_parenthesis.cx",
					StartOffset: 1,
					Length:      12,
					Name:        "First",
					Type:        "Number",
					Value:       0,
				},
				{
					PackageID:   "hello",
					FileID:      "enum_in_parenthesis.cx",
					StartOffset: 1,
					Length:      6,
					Name:        "Second",
					Type:        "Number",
					Value:       1,
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.scenario, func(t *testing.T) {
			srcBytes, err := os.ReadFile(tc.testDir)
			removeComment := declaration_extraction.RemoveComment(srcBytes)
			fileName := filepath.Base(tc.testDir)
			pkg := declaration_extraction.ExtractPackages(removeComment)
			if err != nil {
				t.Fatal(err)
			}
			enums := declaration_extraction.ExtractEnums(removeComment, fileName, pkg)

			for i := range enums {
				if enums[i] != tc.wantEnums[i] {
					t.Errorf("want enums %+v, got %+v", tc.wantEnums[i], enums[i])
				}
			}
		})
	}
}

func TestDeclarationExtraction_ExtractStructs(t *testing.T) {

	tests := []struct {
		scenario    string
		testDir     string
		wantStructs []declaration_extraction.StructDeclaration
	}{
		{
			scenario: "Has structs",
			testDir:  "./test_files/test.cx",
			wantStructs: []declaration_extraction.StructDeclaration{
				{
					PackageID:   "hello",
					FileID:      "test.cx",
					StartOffset: 0,
					Length:      18,
					Name:        "person",
				},
				{
					PackageID:   "hello",
					FileID:      "test.cx",
					StartOffset: 1,
					Length:      18,
					Name:        "animal",
				},
				{
					PackageID:   "hello",
					FileID:      "test.cx",
					StartOffset: 0,
					Length:      18,
					Name:        "Direction",
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.scenario, func(t *testing.T) {
			srcBytes, err := os.ReadFile(tc.testDir)
			removeComment := declaration_extraction.RemoveComment(srcBytes)
			fileName := filepath.Base(tc.testDir)
			pkg := declaration_extraction.ExtractPackages(removeComment)
			if err != nil {
				t.Fatal(err)
			}
			structs := declaration_extraction.ExtractStructs(removeComment, fileName, pkg)

			for i := range structs {
				if structs[i] != tc.wantStructs[i] {
					t.Errorf("want structs %v, got %v", tc.wantStructs[i], structs[i])
				}
			}
		})
	}
}

func TestDeclarationExtraction_ExtractFuncs(t *testing.T) {

	tests := []struct {
		scenario  string
		testDir   string
		wantFuncs []declaration_extraction.FuncDeclaration
	}{
		{
			scenario: "Has structs",
			testDir:  "./test_files/test.cx",
			wantFuncs: []declaration_extraction.FuncDeclaration{
				{
					PackageID:   "hello",
					FileID:      "test.cx",
					StartOffset: 0,
					Length:      12,
					Name:        "main",
				},
				{
					PackageID:   "hello",
					FileID:      "test.cx",
					StartOffset: 1,
					Length:      19,
					Name:        "functionTwo",
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.scenario, func(t *testing.T) {
			srcBytes, err := os.ReadFile(tc.testDir)
			removeComment := declaration_extraction.RemoveComment(srcBytes)
			fileName := filepath.Base(tc.testDir)
			pkg := declaration_extraction.ExtractPackages(removeComment)
			if err != nil {
				t.Fatal(err)
			}
			funcs := declaration_extraction.ExtractFuncs(removeComment, fileName, pkg)

			for i := range funcs {
				if funcs[i] != tc.wantFuncs[i] {
					t.Errorf("want funcs %v, got %v", tc.wantFuncs[i], funcs[i])
				}
			}
		})
	}
}

func TestDeclarationExtraction_ReDeclarationCheck(t *testing.T) {

	tests := []struct {
		scenario   string
		testDir    string
		wantReDclr error
	}{
		{
			scenario:   "No Redeclarations",
			testDir:    "./test_files/test.cx",
			wantReDclr: nil,
		},
		{
			scenario:   "Redeclared globals",
			testDir:    "./test_files/redeclaration_global.cx",
			wantReDclr: errors.New("global redeclared"),
		},
		{
			scenario:   "Redeclared enums",
			testDir:    "./test_files/redeclaration_enum.cx",
			wantReDclr: errors.New("enum redeclared"),
		},
		{
			scenario:   "Redeclared structs",
			testDir:    "./test_files/redeclaration_struct.cx",
			wantReDclr: errors.New("struct redeclared"),
		},
		{
			scenario:   "Redeclared funcs",
			testDir:    "./test_files/redeclaration_func.cx",
			wantReDclr: errors.New("func redeclared"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.scenario, func(t *testing.T) {
			srcBytes, err := os.ReadFile(tc.testDir)
			removeComment := declaration_extraction.RemoveComment(srcBytes)
			fileName := filepath.Base(tc.testDir)
			pkg := declaration_extraction.ExtractPackages(removeComment)
			if err != nil {
				t.Fatal(err)
			}
			globals := declaration_extraction.ExtractGlobals(removeComment, fileName, pkg)
			enums := declaration_extraction.ExtractEnums(removeComment, fileName, pkg)
			structs := declaration_extraction.ExtractStructs(removeComment, fileName, pkg)
			funcs := declaration_extraction.ExtractFuncs(removeComment, fileName, pkg)

			reDeclarationCheck := declaration_extraction.ReDeclarationCheck(globals, enums, structs, funcs)

			if errors.Is(reDeclarationCheck, tc.wantReDclr) && reDeclarationCheck != nil ||
				(reDeclarationCheck != nil && tc.wantReDclr == nil) {
				t.Errorf("want redeclaration check %v, got %v", tc.wantReDclr, reDeclarationCheck)
			}
		})
	}
}
