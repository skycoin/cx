package declaration_extraction_test

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/skycoin/cx/cmd/declaration_extraction"
)

func TestDeclarationExtraction_ReplaceCommentsWithWhitespaces(t *testing.T) {

	tests := []struct {
		scenario            string
		testDir             string
		wantCommentReplaced string
	}{
		{
			scenario:            "Has comments",
			testDir:             "./test_files/test.cx",
			wantCommentReplaced: "./test_files/replaceCommentResult.cx",
		},
	}

	for _, tc := range tests {
		t.Run(tc.scenario, func(t *testing.T) {
			srcBytes, err := os.ReadFile(tc.testDir)
			if err != nil {
				t.Fatal(err)
			}
			wantBytes, err := os.ReadFile(tc.wantCommentReplaced)
			if err != nil {
				t.Fatal(err)
			}
			commentReplaced := declaration_extraction.ReplaceCommentsWithWhitespaces(srcBytes)

			if len(srcBytes) != len(commentReplaced) {
				t.Errorf("Length not the same: orginal %vbytes, replaced %vbytes", len(srcBytes), len(commentReplaced))
			}

			srcLines := bytes.Count(srcBytes, []byte("\n")) + 1
			newLines := bytes.Count(commentReplaced, []byte("\n")) + 1

			if srcLines != newLines {
				t.Errorf("Lines not equal: original %vlines, new %vlines", srcLines, newLines)
			}

			if string(commentReplaced) != string(wantBytes) {
				t.Errorf("want comments replaced\n%v\ngot\n%v", string(wantBytes), string(commentReplaced))
				file, err := os.Create("gotCommentsReplaced.cx")
				if err != nil {
					t.Fatal(err)
				}
				file.Write(commentReplaced)
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
			ReplaceCommentsWithWhitespaces := declaration_extraction.ReplaceCommentsWithWhitespaces(srcBytes)

			if err != nil {
				t.Fatal(err)
			}
			pkg := declaration_extraction.ExtractPackages(ReplaceCommentsWithWhitespaces)

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
					PackageID:          "hello",
					FileID:             "test.cx",
					StartOffset:        15,
					Length:             16,
					LineNumber:         2,
					GlobalVariableName: "apple",
				},
				{
					PackageID:          "hello",
					FileID:             "test.cx",
					StartOffset:        37,
					Length:             17,
					LineNumber:         4,
					GlobalVariableName: "banana",
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.scenario, func(t *testing.T) {
			srcBytes, err := os.ReadFile(tc.testDir)
			ReplaceCommentsWithWhitespaces := declaration_extraction.ReplaceCommentsWithWhitespaces(srcBytes)
			fileName := filepath.Base(tc.testDir)
			pkg := declaration_extraction.ExtractPackages(ReplaceCommentsWithWhitespaces)
			if err != nil {
				t.Fatal(err)
			}

			globals, err := declaration_extraction.ExtractGlobals(ReplaceCommentsWithWhitespaces, fileName, pkg)
			if err != nil {
				t.Fatal(err)
			}

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
					PackageID:        "hello",
					FileID:           "test.cx",
					StartOffset:      383,
					Length:           15,
					LineNumber:       33,
					Type:             "Direction",
					Value:            0,
					EnumVariableName: "North",
				},
				{
					PackageID:        "hello",
					FileID:           "test.cx",
					StartOffset:      408,
					Length:           5,
					LineNumber:       34,
					Type:             "Direction",
					Value:            1,
					EnumVariableName: "South",
				},
				{
					PackageID:        "hello",
					FileID:           "test.cx",
					StartOffset:      416,
					Length:           4,
					LineNumber:       35,
					Type:             "Direction",
					Value:            2,
					EnumVariableName: "East",
				},
				{
					PackageID:        "hello",
					FileID:           "test.cx",
					StartOffset:      423,
					Length:           4,
					LineNumber:       36,
					Type:             "Direction",
					Value:            3,
					EnumVariableName: "West",
				},
				{
					PackageID:        "hello",
					FileID:           "test.cx",
					StartOffset:      444,
					Length:           12,
					LineNumber:       40,
					Type:             "Number",
					Value:            0,
					EnumVariableName: "First",
				},
				{
					PackageID:        "hello",
					FileID:           "test.cx",
					StartOffset:      466,
					Length:           6,
					LineNumber:       41,
					Type:             "Number",
					Value:            1,
					EnumVariableName: "Second",
				},
			},
		},
		{
			scenario: "Has enums and nested parenthesis",
			testDir:  "./test_files/enum_in_parenthesis.cx",
			wantEnums: []declaration_extraction.EnumDeclaration{
				{
					PackageID:        "hello",
					FileID:           "enum_in_parenthesis.cx",
					StartOffset:      341,
					Length:           15,
					Type:             "Direction",
					LineNumber:       33,
					Value:            0,
					EnumVariableName: "North",
				},
				{
					PackageID:        "hello",
					FileID:           "enum_in_parenthesis.cx",
					StartOffset:      366,
					Length:           5,
					LineNumber:       34,
					Type:             "Direction",
					Value:            1,
					EnumVariableName: "South",
				},
				{
					PackageID:        "hello",
					FileID:           "enum_in_parenthesis.cx",
					StartOffset:      387,
					Length:           4,
					LineNumber:       35,
					Type:             "Direction",
					Value:            2,
					EnumVariableName: "East",
				},
				{
					PackageID:        "hello",
					FileID:           "enum_in_parenthesis.cx",
					StartOffset:      394,
					Length:           4,
					LineNumber:       36,
					Type:             "Direction",
					Value:            3,
					EnumVariableName: "West",
				},
				{
					PackageID:        "hello",
					FileID:           "enum_in_parenthesis.cx",
					StartOffset:      415,
					Length:           12,
					LineNumber:       40,
					Type:             "Number",
					Value:            0,
					EnumVariableName: "First",
				},
				{
					PackageID:        "hello",
					FileID:           "enum_in_parenthesis.cx",
					StartOffset:      437,
					Length:           6,
					LineNumber:       41,
					Type:             "Number",
					Value:            1,
					EnumVariableName: "Second",
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.scenario, func(t *testing.T) {
			srcBytes, err := os.ReadFile(tc.testDir)
			ReplaceCommentsWithWhitespaces := declaration_extraction.ReplaceCommentsWithWhitespaces(srcBytes)
			fileName := filepath.Base(tc.testDir)
			pkg := declaration_extraction.ExtractPackages(ReplaceCommentsWithWhitespaces)
			if err != nil {
				t.Fatal(err)
			}

			enums, err := declaration_extraction.ExtractEnums(ReplaceCommentsWithWhitespaces, fileName, pkg)
			if err != nil {
				t.Fatal(err)
			}

			for i := range enums {
				if enums[i] != tc.wantEnums[i] {
					t.Errorf("want enums %v, got %v", tc.wantEnums[i], enums[i])
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
					PackageID:          "hello",
					FileID:             "test.cx",
					StartOffset:        171,
					Length:             18,
					LineNumber:         14,
					StructVariableName: "person",
				},
				{
					PackageID:          "hello",
					FileID:             "test.cx",
					StartOffset:        250,
					Length:             39,
					LineNumber:         21,
					StructVariableName: "animal",
				},
				{
					PackageID:          "hello",
					FileID:             "test.cx",
					StartOffset:        351,
					Length:             18,
					LineNumber:         30,
					StructVariableName: "Direction",
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.scenario, func(t *testing.T) {
			srcBytes, err := os.ReadFile(tc.testDir)
			ReplaceCommentsWithWhitespaces := declaration_extraction.ReplaceCommentsWithWhitespaces(srcBytes)
			fileName := filepath.Base(tc.testDir)
			pkg := declaration_extraction.ExtractPackages(ReplaceCommentsWithWhitespaces)
			if err != nil {
				t.Fatal(err)
			}

			structs, err := declaration_extraction.ExtractStructs(ReplaceCommentsWithWhitespaces, fileName, pkg)
			if err != nil {
				t.Fatal(err)
			}

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
			scenario: "Has funcs",
			testDir:  "./test_files/test.cx",
			wantFuncs: []declaration_extraction.FuncDeclaration{
				{
					PackageID:        "hello",
					FileID:           "test.cx",
					StartOffset:      212,
					Length:           12,
					LineNumber:       18,
					FuncVariableName: "main",
				},
				{
					PackageID:        "hello",
					FileID:           "test.cx",
					StartOffset:      321,
					Length:           19,
					LineNumber:       26,
					FuncVariableName: "functionTwo",
				},
				{
					PackageID:        "hello",
					FileID:           "test.cx",
					StartOffset:      479,
					Length:           39,
					LineNumber:       44,
					FuncVariableName: "functionWithSingleReturn",
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.scenario, func(t *testing.T) {
			srcBytes, err := os.ReadFile(tc.testDir)
			ReplaceCommentsWithWhitespaces := declaration_extraction.ReplaceCommentsWithWhitespaces(srcBytes)
			fileName := filepath.Base(tc.testDir)
			pkg := declaration_extraction.ExtractPackages(ReplaceCommentsWithWhitespaces)
			if err != nil {
				t.Fatal(err)
			}

			funcs, err := declaration_extraction.ExtractFuncs(ReplaceCommentsWithWhitespaces, fileName, pkg)
			if err != nil {
				t.Fatal(err)
			}

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
			ReplaceCommentsWithWhitespaces := declaration_extraction.ReplaceCommentsWithWhitespaces(srcBytes)
			fileName := filepath.Base(tc.testDir)
			pkg := declaration_extraction.ExtractPackages(ReplaceCommentsWithWhitespaces)
			if err != nil {
				t.Fatal(err)
			}

			globals, err := declaration_extraction.ExtractGlobals(ReplaceCommentsWithWhitespaces, fileName, pkg)
			if err != nil {
				t.Fatal(err)
			}

			enums, err := declaration_extraction.ExtractEnums(ReplaceCommentsWithWhitespaces, fileName, pkg)
			if err != nil {
				t.Fatal(err)
			}

			structs, err := declaration_extraction.ExtractStructs(ReplaceCommentsWithWhitespaces, fileName, pkg)
			if err != nil {
				t.Fatal(err)
			}

			funcs, err := declaration_extraction.ExtractFuncs(ReplaceCommentsWithWhitespaces, fileName, pkg)
			if err != nil {
				t.Fatal(err)
			}

			reDeclarationCheck := declaration_extraction.ReDeclarationCheck(globals, enums, structs, funcs)

			if errors.Is(reDeclarationCheck, tc.wantReDclr) && reDeclarationCheck != nil ||
				(reDeclarationCheck != nil && tc.wantReDclr == nil) {
				t.Errorf("want redeclaration check %v, got %v", tc.wantReDclr, reDeclarationCheck)
			}
		})
	}
}

func TestDeclarationExtraction_GetDeclarations(t *testing.T) {

	tests := []struct {
		scenario         string
		testDir          string
		wantDeclarations []string
	}{
		{
			scenario: "Has declarations",
			testDir:  "./test_files/test.cx",
			wantDeclarations: []string{
				"var apple string",
				"var banana string",
				"North Direction",
				"South",
				"East",
				"West",
				"First Number",
				"Second",
				"type person struct",
				"type animal                      struct",
				"type Direction int",
				"func main ()",
				"func functionTwo ()",
				"func functionWithSingleReturn () string",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.scenario, func(t *testing.T) {
			srcBytes, err := os.ReadFile(tc.testDir)
			ReplaceCommentsWithWhitespaces := declaration_extraction.ReplaceCommentsWithWhitespaces(srcBytes)
			fileName := filepath.Base(tc.testDir)
			pkg := declaration_extraction.ExtractPackages(ReplaceCommentsWithWhitespaces)
			if err != nil {
				t.Fatal(err)
			}

			globals, err := declaration_extraction.ExtractGlobals(ReplaceCommentsWithWhitespaces, fileName, pkg)
			if err != nil {
				t.Fatal(err)
			}

			enums, err := declaration_extraction.ExtractEnums(ReplaceCommentsWithWhitespaces, fileName, pkg)
			if err != nil {
				t.Fatal(err)
			}

			structs, err := declaration_extraction.ExtractStructs(ReplaceCommentsWithWhitespaces, fileName, pkg)
			if err != nil {
				t.Fatal(err)
			}

			funcs, err := declaration_extraction.ExtractFuncs(ReplaceCommentsWithWhitespaces, fileName, pkg)
			if err != nil {
				t.Fatal(err)
			}

			if declaration_extraction.ReDeclarationCheck(globals, enums, structs, funcs) != nil {
				t.Fatal(err)
			}

			declarations := declaration_extraction.GetDeclarations(srcBytes, globals, enums, structs, funcs)

			for i := range declarations {
				if declarations[i] != tc.wantDeclarations[i] {
					t.Errorf("want declaration %v, got %v", tc.wantDeclarations[i], declarations[i])
				}
			}

		})
	}
}

func TestDeclarationExtraction_ExtractAllDeclarations(t *testing.T) {

	tests := []struct {
		scenario    string
		testDirs    []string
		wantGlobals []declaration_extraction.GlobalDeclaration
		wantEnums   []declaration_extraction.EnumDeclaration
		wantStructs []declaration_extraction.StructDeclaration
		wantFuncs   []declaration_extraction.FuncDeclaration
		wantError   error
	}{
		{
			scenario: "Single file",
			testDirs: []string{
				"./test_files/test.cx",
			},
			wantGlobals: []declaration_extraction.GlobalDeclaration{
				{
					PackageID:          "hello",
					FileID:             "./test_files/test.cx",
					StartOffset:        15,
					Length:             16,
					LineNumber:         2,
					GlobalVariableName: "apple",
				},
				{
					PackageID:          "hello",
					FileID:             "./test_files/test.cx",
					StartOffset:        37,
					Length:             17,
					LineNumber:         4,
					GlobalVariableName: "banana",
				},
			},
			wantEnums: []declaration_extraction.EnumDeclaration{
				{
					PackageID:        "hello",
					FileID:           "./test_files/test.cx",
					StartOffset:      383,
					Length:           15,
					LineNumber:       33,
					Type:             "Direction",
					Value:            0,
					EnumVariableName: "North",
				},
				{
					PackageID:        "hello",
					FileID:           "./test_files/test.cx",
					StartOffset:      408,
					Length:           5,
					LineNumber:       34,
					Type:             "Direction",
					Value:            1,
					EnumVariableName: "South",
				},
				{
					PackageID:        "hello",
					FileID:           "./test_files/test.cx",
					StartOffset:      416,
					Length:           4,
					LineNumber:       35,
					Type:             "Direction",
					Value:            2,
					EnumVariableName: "East",
				},
				{
					PackageID:        "hello",
					FileID:           "./test_files/test.cx",
					StartOffset:      423,
					Length:           4,
					LineNumber:       36,
					Type:             "Direction",
					Value:            3,
					EnumVariableName: "West",
				},
				{
					PackageID:        "hello",
					FileID:           "./test_files/test.cx",
					StartOffset:      444,
					Length:           12,
					LineNumber:       40,
					Type:             "Number",
					Value:            0,
					EnumVariableName: "First",
				},
				{
					PackageID:        "hello",
					FileID:           "./test_files/test.cx",
					StartOffset:      466,
					Length:           6,
					LineNumber:       41,
					Type:             "Number",
					Value:            1,
					EnumVariableName: "Second",
				},
			},
			wantStructs: []declaration_extraction.StructDeclaration{
				{
					PackageID:          "hello",
					FileID:             "./test_files/test.cx",
					StartOffset:        171,
					Length:             18,
					LineNumber:         14,
					StructVariableName: "person",
				},
				{
					PackageID:          "hello",
					FileID:             "./test_files/test.cx",
					StartOffset:        250,
					Length:             39,
					LineNumber:         21,
					StructVariableName: "animal",
				},
				{
					PackageID:          "hello",
					FileID:             "./test_files/test.cx",
					StartOffset:        351,
					Length:             18,
					LineNumber:         30,
					StructVariableName: "Direction",
				},
			},
			wantFuncs: []declaration_extraction.FuncDeclaration{
				{
					PackageID:        "hello",
					FileID:           "./test_files/test.cx",
					StartOffset:      212,
					Length:           12,
					LineNumber:       18,
					FuncVariableName: "main",
				},
				{
					PackageID:        "hello",
					FileID:           "./test_files/test.cx",
					StartOffset:      321,
					Length:           19,
					LineNumber:       26,
					FuncVariableName: "functionTwo",
				},
				{
					PackageID:        "hello",
					FileID:           "./test_files/test.cx",
					StartOffset:      479,
					Length:           39,
					LineNumber:       44,
					FuncVariableName: "functionWithSingleReturn",
				},
			},
		},
		{
			scenario: "Redeclared Global",
			testDirs: []string{
				"./test_files/multiple_files/helper.cx",
				"./test_files/multiple_files/main.cx",
				"./test_files/multiple_files/utility.cx",
				"./test_files/multiple_files/worker.cx",
			},
			wantGlobals: []declaration_extraction.GlobalDeclaration{
				{
					PackageID:          "hello",
					FileID:             "./test_files/test.cx",
					StartOffset:        15,
					Length:             16,
					LineNumber:         2,
					GlobalVariableName: "apple",
				},
				{
					PackageID:          "hello",
					FileID:             "./test_files/test.cx",
					StartOffset:        37,
					Length:             17,
					LineNumber:         4,
					GlobalVariableName: "banana",
				},
			},
			wantEnums: []declaration_extraction.EnumDeclaration{
				{
					PackageID:        "hello",
					FileID:           "./test_files/test.cx",
					StartOffset:      383,
					Length:           15,
					LineNumber:       33,
					Type:             "Direction",
					Value:            0,
					EnumVariableName: "North",
				},
				{
					PackageID:        "hello",
					FileID:           "./test_files/test.cx",
					StartOffset:      408,
					Length:           5,
					LineNumber:       34,
					Type:             "Direction",
					Value:            1,
					EnumVariableName: "South",
				},
				{
					PackageID:        "hello",
					FileID:           "./test_files/test.cx",
					StartOffset:      416,
					Length:           4,
					LineNumber:       35,
					Type:             "Direction",
					Value:            2,
					EnumVariableName: "East",
				},
				{
					PackageID:        "hello",
					FileID:           "./test_files/test.cx",
					StartOffset:      423,
					Length:           4,
					LineNumber:       36,
					Type:             "Direction",
					Value:            3,
					EnumVariableName: "West",
				},
				{
					PackageID:        "hello",
					FileID:           "./test_files/test.cx",
					StartOffset:      444,
					Length:           12,
					LineNumber:       40,
					Type:             "Number",
					Value:            0,
					EnumVariableName: "First",
				},
				{
					PackageID:        "hello",
					FileID:           "./test_files/test.cx",
					StartOffset:      466,
					Length:           6,
					LineNumber:       41,
					Type:             "Number",
					Value:            1,
					EnumVariableName: "Second",
				},
			},
			wantStructs: []declaration_extraction.StructDeclaration{
				{
					PackageID:          "hello",
					FileID:             "./test_files/test.cx",
					StartOffset:        171,
					Length:             18,
					LineNumber:         14,
					StructVariableName: "person",
				},
				{
					PackageID:          "hello",
					FileID:             "./test_files/test.cx",
					StartOffset:        250,
					Length:             39,
					LineNumber:         21,
					StructVariableName: "animal",
				},
				{
					PackageID:          "hello",
					FileID:             "./test_files/test.cx",
					StartOffset:        351,
					Length:             18,
					LineNumber:         30,
					StructVariableName: "Direction",
				},
			},
			wantFuncs: []declaration_extraction.FuncDeclaration{
				{
					PackageID:        "hello",
					FileID:           "./test_files/test.cx",
					StartOffset:      212,
					Length:           12,
					LineNumber:       18,
					FuncVariableName: "main",
				},
				{
					PackageID:        "hello",
					FileID:           "./test_files/test.cx",
					StartOffset:      321,
					Length:           19,
					LineNumber:       26,
					FuncVariableName: "functionTwo",
				},
				{
					PackageID:        "hello",
					FileID:           "./test_files/test.cx",
					StartOffset:      479,
					Length:           39,
					LineNumber:       44,
					FuncVariableName: "functionWithSingleReturn",
				},
			},
			wantError: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.scenario, func(t *testing.T) {

			var files []*os.File

			for _, testDir := range tc.testDirs {

				file, err := os.Open(testDir)

				if err != nil {
					t.Fatal(err)
				}

				files = append(files, file)
			}

			Globals, Enums, Structs, Funcs, Err := declaration_extraction.ExtractAllDeclarations(files)

			for i, global := range Globals {
				if global != tc.wantGlobals[i] {
					t.Errorf("want global %v, got %v", tc.wantGlobals[i], global)
				}
			}

			for i, enum := range Enums {
				if enum != tc.wantEnums[i] {
					t.Errorf("want enum %v, got %v", tc.wantEnums[i], enum)
				}
			}

			for i, strct := range Structs {
				if strct != tc.wantStructs[i] {
					t.Errorf("want struct %v, got %v", tc.wantStructs[i], strct)
				}
			}

			for i, fun := range Funcs {
				if fun != tc.wantFuncs[i] {
					t.Errorf("want func %v, got %v", tc.wantFuncs[i], fun)
				}
			}

			if errors.Is(Err, tc.wantError) && Err != nil ||
				(Err != nil && tc.wantError == nil) {
				t.Errorf("want error %v, got %v", tc.wantError, Err)
			}

		})
	}
}
