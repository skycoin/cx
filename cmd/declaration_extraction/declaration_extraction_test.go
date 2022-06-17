package declaration_extraction_test

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/skycoin/cx/cmd/declaration_extraction"
)

func setOffset(offset int, lineNumber int) int {

	var newOffset int = offset

	runtimeOS := runtime.GOOS

	if runtimeOS == "windows" {
		newOffset += lineNumber - 1
	}

	return newOffset
}

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
					StartOffset:        setOffset(14, 2),
					Length:             16,
					LineNumber:         2,
					GlobalVariableName: "apple",
				},
				{
					PackageID:          "hello",
					FileID:             "test.cx",
					StartOffset:        setOffset(34, 4),
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
					StartOffset:      setOffset(351, 33),
					Length:           15,
					LineNumber:       33,
					Type:             "Direction",
					Value:            0,
					EnumVariableName: "North",
				},
				{
					PackageID:        "hello",
					FileID:           "test.cx",
					StartOffset:      setOffset(375, 34),
					Length:           5,
					LineNumber:       34,
					Type:             "Direction",
					Value:            1,
					EnumVariableName: "South",
				},
				{
					PackageID:        "hello",
					FileID:           "test.cx",
					StartOffset:      setOffset(382, 35),
					Length:           4,
					LineNumber:       35,
					Type:             "Direction",
					Value:            2,
					EnumVariableName: "East",
				},
				{
					PackageID:        "hello",
					FileID:           "test.cx",
					StartOffset:      setOffset(388, 36),
					Length:           4,
					LineNumber:       36,
					Type:             "Direction",
					Value:            3,
					EnumVariableName: "West",
				},
				{
					PackageID:        "hello",
					FileID:           "test.cx",
					StartOffset:      setOffset(405, 40),
					Length:           12,
					LineNumber:       40,
					Type:             "Number",
					Value:            0,
					EnumVariableName: "First",
				},
				{
					PackageID:        "hello",
					FileID:           "test.cx",
					StartOffset:      setOffset(426, 41),
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
					StartOffset:      setOffset(309, 33),
					Length:           15,
					Type:             "Direction",
					LineNumber:       33,
					Value:            0,
					EnumVariableName: "North",
				},
				{
					PackageID:        "hello",
					FileID:           "enum_in_parenthesis.cx",
					StartOffset:      setOffset(333, 34),
					Length:           5,
					LineNumber:       34,
					Type:             "Direction",
					Value:            1,
					EnumVariableName: "South",
				},
				{
					PackageID:        "hello",
					FileID:           "enum_in_parenthesis.cx",
					StartOffset:      setOffset(353, 35),
					Length:           4,
					LineNumber:       35,
					Type:             "Direction",
					Value:            2,
					EnumVariableName: "East",
				},
				{
					PackageID:        "hello",
					FileID:           "enum_in_parenthesis.cx",
					StartOffset:      setOffset(359, 36),
					Length:           4,
					LineNumber:       36,
					Type:             "Direction",
					Value:            3,
					EnumVariableName: "West",
				},
				{
					PackageID:        "hello",
					FileID:           "enum_in_parenthesis.cx",
					StartOffset:      setOffset(376, 40),
					Length:           12,
					LineNumber:       40,
					Type:             "Number",
					Value:            0,
					EnumVariableName: "First",
				},
				{
					PackageID:        "hello",
					FileID:           "enum_in_parenthesis.cx",
					StartOffset:      setOffset(397, 41),
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
					StartOffset:        setOffset(158, 14),
					Length:             18,
					LineNumber:         14,
					StructVariableName: "person",
				},
				{
					PackageID:          "hello",
					FileID:             "test.cx",
					StartOffset:        setOffset(230, 21),
					Length:             39,
					LineNumber:         21,
					StructVariableName: "animal",
				},
				{
					PackageID:          "hello",
					FileID:             "test.cx",
					StartOffset:        setOffset(322, 30),
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
					StartOffset:      setOffset(195, 18),
					Length:           12,
					LineNumber:       18,
					FuncVariableName: "main",
				},
				{
					PackageID:        "hello",
					FileID:           "test.cx",
					StartOffset:      setOffset(296, 26),
					Length:           19,
					LineNumber:       26,
					FuncVariableName: "functionTwo",
				},
				{
					PackageID:        "hello",
					FileID:           "test.cx",
					StartOffset:      setOffset(436, 44),
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

			if reDeclarationCheck != nil && tc.wantReDclr == nil {
				t.Errorf("want error %v, got %v", tc.wantReDclr, reDeclarationCheck)
			}

			if reDeclarationCheck != nil && tc.wantReDclr != nil {
				if reDeclarationCheck.Error() != tc.wantReDclr.Error() {
					t.Errorf("want error %v, got %v", tc.wantReDclr, reDeclarationCheck)
				}
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
		wantGlobals int
		wantEnums   int
		wantStructs int
		wantFuncs   int
		wantError   error
	}{
		{
			scenario: "Single file",
			testDirs: []string{
				"./test_files/test.cx",
			},
			wantGlobals: 2,
			wantEnums:   6,
			wantStructs: 3,
			wantFuncs:   3,
			wantError:   nil,
		},
		{
			scenario: "Mulitple files",
			testDirs: []string{
				"./test_files/multiple_files/helper.cx",
				"./test_files/multiple_files/main.cx",
				"./test_files/multiple_files/utility.cx",
				"./test_files/multiple_files/worker.cx",
			},
			wantGlobals: 2,
			wantEnums:   3,
			wantStructs: 1,
			wantFuncs:   3,
			wantError:   nil,
		},
		{
			scenario: "Redeclared Global",
			testDirs: []string{
				"./test_files/multiple_files/helper.cx",
				"./test_files/multiple_files/main.cx",
				"./test_files/multiple_files/utility.cx",
				"./test_files/multiple_files/worker.cx",
				"./test_files/redeclaration_global.cx",
			},
			wantGlobals: 5,
			wantEnums:   9,
			wantStructs: 4,
			wantFuncs:   5,
			wantError:   errors.New("global redeclared"),
		},
		{
			scenario: "Redeclared Enum",
			testDirs: []string{
				"./test_files/multiple_files/helper.cx",
				"./test_files/multiple_files/main.cx",
				"./test_files/multiple_files/utility.cx",
				"./test_files/multiple_files/worker.cx",
				"./test_files/redeclaration_enum.cx",
			},
			wantGlobals: 4,
			wantEnums:   10,
			wantStructs: 4,
			wantFuncs:   5,
			wantError:   errors.New("enum redeclared"),
		},
		{
			scenario: "Redeclared Struct",
			testDirs: []string{
				"./test_files/multiple_files/helper.cx",
				"./test_files/multiple_files/main.cx",
				"./test_files/multiple_files/utility.cx",
				"./test_files/multiple_files/worker.cx",
				"./test_files/redeclaration_struct.cx",
			},
			wantGlobals: 4,
			wantEnums:   9,
			wantStructs: 5,
			wantFuncs:   5,
			wantError:   errors.New("struct redeclared"),
		},
		{
			scenario: "Redeclared Func",
			testDirs: []string{
				"./test_files/multiple_files/helper.cx",
				"./test_files/multiple_files/main.cx",
				"./test_files/multiple_files/utility.cx",
				"./test_files/multiple_files/worker.cx",
				"./test_files/redeclaration_func.cx",
			},
			wantGlobals: 4,
			wantEnums:   9,
			wantStructs: 4,
			wantFuncs:   6,
			wantError:   errors.New("func redeclared"),
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

			if len(Globals) == 0 && len(Enums) == 0 && len(Structs) == 0 && len(Funcs) == 0 {
				t.Error("No Declarations found")
			}

			if len(Globals) != tc.wantGlobals {
				t.Errorf("want global %v, got %v", tc.wantGlobals, len(Globals))
			}

			if len(Enums) != tc.wantEnums {
				t.Errorf("want enum %v, got %v", tc.wantEnums, len(Enums))
			}

			if len(Structs) != tc.wantStructs {
				t.Errorf("want struct %v, got %v", tc.wantStructs, len(Structs))
			}

			if len(Funcs) != tc.wantFuncs {
				t.Errorf("want func %v, got %v", tc.wantFuncs, len(Funcs))
			}

			if Err != nil && tc.wantError == nil {
				t.Errorf("want error %v, got %v", tc.wantError, Err)
			}

			if Err != nil {
				if Err.Error() != tc.wantError.Error() {
					t.Errorf("want error %v, got %v", tc.wantError, Err)
				}
			}
		})
	}
}
