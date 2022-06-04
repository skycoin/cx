package declaration_extraction_test

import (
	"errors"
	"os"
	"testing"

	"test.com/declaration_extraction"
)

func TestDeclarationExtraction_ExtractGlobal(t *testing.T) {

	tests := []struct {
		scenario    string
		testDir     string
		wantGlobals []declaration_extraction.GlobalDeclaration
	}{
		{
			scenario: "Has globals",
			testDir:  "./test.cx",
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
			srcBytes, err := os.Open(tc.testDir)
			if err != nil {
				t.Fatal(err)
			}
			globals, err := declaration_extraction.ExtractGlobals(srcBytes)
			if err != nil {
				t.Fatal(err)
			}
			for i := range globals {
				if globals[i] != tc.wantGlobals[i] {
					t.Errorf("want globals %+v, got %+v", tc.wantGlobals[i], globals[i])
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
			testDir:  "./test.cx",
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
			testDir:  "./testInPrts.cx",
			wantEnums: []declaration_extraction.EnumDeclaration{
				{
					PackageID:   "hello",
					FileID:      "testInPrts.cx",
					StartOffset: 1,
					Length:      15,
					Name:        "North",
					Type:        "Direction",
					Value:       0,
				},
				{
					PackageID:   "hello",
					FileID:      "testInPrts.cx",
					StartOffset: 1,
					Length:      5,
					Name:        "South",
					Type:        "Direction",
					Value:       1,
				},
				{
					PackageID:   "hello",
					FileID:      "testInPrts.cx",
					StartOffset: 1,
					Length:      4,
					Name:        "East",
					Type:        "Direction",
					Value:       2,
				},
				{
					PackageID:   "hello",
					FileID:      "testInPrts.cx",
					StartOffset: 1,
					Length:      4,
					Name:        "West",
					Type:        "Direction",
					Value:       3,
				},
				{
					PackageID:   "hello",
					FileID:      "testInPrts.cx",
					StartOffset: 1,
					Length:      12,
					Name:        "First",
					Type:        "Number",
					Value:       0,
				},
				{
					PackageID:   "hello",
					FileID:      "testInPrts.cx",
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
			srcBytes, err := os.Open(tc.testDir)
			if err != nil {
				t.Fatal(err)
			}
			enums, err := declaration_extraction.ExtractEnums(srcBytes)
			if err != nil {
				t.Fatal(err)
			}
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
			testDir:  "./test.cx",
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
			srcBytes, err := os.Open(tc.testDir)
			if err != nil {
				t.Fatal(err)
			}
			structs, err := declaration_extraction.ExtractStructs(srcBytes)
			if err != nil {
				t.Fatal(err)
			}
			for i := range structs {
				if structs[i] != tc.wantStructs[i] {
					t.Errorf("want structs %+v, got %+v", tc.wantStructs[i], structs[i])
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
			testDir:  "./test.cx",
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
			srcBytes, err := os.Open(tc.testDir)
			if err != nil {
				t.Fatal(err)
			}
			funcs, err := declaration_extraction.ExtractFuncs(srcBytes)
			if err != nil {
				t.Fatal(err)
			}
			for i := range funcs {
				if funcs[i] != tc.wantFuncs[i] {
					t.Errorf("want funcs %+v, got %+v", tc.wantFuncs[i], funcs[i])
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
			testDir:    "./test.cx",
			wantReDclr: nil,
		},
		{
			scenario:   "Redeclared globals",
			testDir:    "./reDclrGlbl.cx",
			wantReDclr: errors.New("global redeclared"),
		},
		{
			scenario:   "Redeclared enums",
			testDir:    "./reDclrEnum.cx",
			wantReDclr: errors.New("enum redeclared"),
		},
		{
			scenario:   "Redeclared structs",
			testDir:    "./reDclrStrct.cx",
			wantReDclr: errors.New("struct redeclared"),
		},
		{
			scenario:   "Redeclared funcs",
			testDir:    "./reDclrFunc.cx",
			wantReDclr: errors.New("func redeclared"),
		},
	}

	// Problems
	//Is there a way to open the file once and use it for all?
	// How to compare error nil to error nil?

	for _, tc := range tests {
		t.Run(tc.scenario, func(t *testing.T) {
			srcBytes, err := os.Open(tc.testDir)
			if err != nil {
				t.Fatal(err)
			}

			globals, err := declaration_extraction.ExtractGlobals(srcBytes)
			if err != nil {
				t.Fatal(err)
			}

			srcBytes, err = os.Open(tc.testDir)
			if err != nil {
				t.Fatal(err)
			}

			enums, err := declaration_extraction.ExtractEnums(srcBytes)
			if err != nil {
				t.Fatal(err)
			}

			srcBytes, err = os.Open(tc.testDir)
			if err != nil {
				t.Fatal(err)
			}
			structs, err := declaration_extraction.ExtractStructs(srcBytes)
			if err != nil {
				t.Fatal(err)
			}

			srcBytes, err = os.Open(tc.testDir)
			if err != nil {
				t.Fatal(err)
			}
			funcs, err := declaration_extraction.ExtractFuncs(srcBytes)
			if err != nil {
				t.Fatal(err)
			}

			reDclr := declaration_extraction.ReDeclarationCheck(globals, enums, structs, funcs)
			if errors.Is(reDclr, tc.wantReDclr) {
				t.Errorf("want %+v, got %+v", tc.wantReDclr, reDclr)
			}

		})
	}
}
