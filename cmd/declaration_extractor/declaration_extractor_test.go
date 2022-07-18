package declaration_extractor_test

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/skycoin/cx/cmd/declaration_extractor"
)

//Sets the offset for windows or other os
func setOffset(offset int, lineNumber int) int {

	//Input offset is the offset for linux/mac
	var newOffset int = offset

	runtimeOS := runtime.GOOS

	//Determines runtime os and sets the offset accordingly
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
		{
			scenario:            "Has quoted single line comment ",
			testDir:             "./test_files/test_2.cx",
			wantCommentReplaced: "./test_files/replaceCommentResult_2.cx",
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
			commentReplaced := declaration_extractor.ReplaceCommentsWithWhitespaces(srcBytes)

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
			ReplaceCommentsWithWhitespaces := declaration_extractor.ReplaceCommentsWithWhitespaces(srcBytes)

			if err != nil {
				t.Fatal(err)
			}
			pkg := declaration_extractor.ExtractPackages(ReplaceCommentsWithWhitespaces)

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
		wantGlobals []declaration_extractor.GlobalDeclaration
	}{
		{
			scenario: "Has globals",
			testDir:  "./test_files/test.cx",
			wantGlobals: []declaration_extractor.GlobalDeclaration{
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
		{
			scenario: "Has Globals 2",
			testDir:  "./test_files/test_2.cx",
			wantGlobals: []declaration_extractor.GlobalDeclaration{
				{
					PackageID:          "test_2",
					FileID:             "test_2.cx",
					StartOffset:        setOffset(35, 7),
					Length:             17,
					LineNumber:         7,
					GlobalVariableName: "global",
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.scenario, func(t *testing.T) {
			srcBytes, err := os.ReadFile(tc.testDir)
			ReplaceCommentsWithWhitespaces := declaration_extractor.ReplaceCommentsWithWhitespaces(srcBytes)
			fileName := filepath.Base(tc.testDir)
			pkg := declaration_extractor.ExtractPackages(ReplaceCommentsWithWhitespaces)
			if err != nil {
				t.Fatal(err)
			}

			globals, err := declaration_extractor.ExtractGlobals(ReplaceCommentsWithWhitespaces, fileName, pkg)
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
		wantEnums []declaration_extractor.EnumDeclaration
	}{
		{
			scenario: "Has enums",
			testDir:  "./test_files/test.cx",
			wantEnums: []declaration_extractor.EnumDeclaration{
				{
					PackageID:   "hello",
					FileID:      "test.cx",
					StartOffset: setOffset(351, 33),
					Length:      15,
					LineNumber:  33,
					Type:        "Direction",
					Value:       0,
					EnumName:    "North",
				},
				{
					PackageID:   "hello",
					FileID:      "test.cx",
					StartOffset: setOffset(375, 34),
					Length:      5,
					LineNumber:  34,
					Type:        "Direction",
					Value:       1,
					EnumName:    "South",
				},
				{
					PackageID:   "hello",
					FileID:      "test.cx",
					StartOffset: setOffset(382, 35),
					Length:      4,
					LineNumber:  35,
					Type:        "Direction",
					Value:       2,
					EnumName:    "East",
				},
				{
					PackageID:   "hello",
					FileID:      "test.cx",
					StartOffset: setOffset(388, 36),
					Length:      4,
					LineNumber:  36,
					Type:        "Direction",
					Value:       3,
					EnumName:    "West",
				},
				{
					PackageID:   "hello",
					FileID:      "test.cx",
					StartOffset: setOffset(405, 40),
					Length:      12,
					LineNumber:  40,
					Type:        "Number",
					Value:       0,
					EnumName:    "First",
				},
				{
					PackageID:   "hello",
					FileID:      "test.cx",
					StartOffset: setOffset(426, 41),
					Length:      6,
					LineNumber:  41,
					Type:        "Number",
					Value:       1,
					EnumName:    "Second",
				},
			},
		},
		{
			scenario: "Has enums and nested parenthesis",
			testDir:  "./test_files/enum_in_parenthesis.cx",
			wantEnums: []declaration_extractor.EnumDeclaration{
				{
					PackageID:   "hello",
					FileID:      "enum_in_parenthesis.cx",
					StartOffset: setOffset(309, 33),
					Length:      15,
					Type:        "Direction",
					LineNumber:  33,
					Value:       0,
					EnumName:    "North",
				},
				{
					PackageID:   "hello",
					FileID:      "enum_in_parenthesis.cx",
					StartOffset: setOffset(333, 34),
					Length:      5,
					LineNumber:  34,
					Type:        "Direction",
					Value:       1,
					EnumName:    "South",
				},
				{
					PackageID:   "hello",
					FileID:      "enum_in_parenthesis.cx",
					StartOffset: setOffset(353, 35),
					Length:      4,
					LineNumber:  35,
					Type:        "Direction",
					Value:       2,
					EnumName:    "East",
				},
				{
					PackageID:   "hello",
					FileID:      "enum_in_parenthesis.cx",
					StartOffset: setOffset(359, 36),
					Length:      4,
					LineNumber:  36,
					Type:        "Direction",
					Value:       3,
					EnumName:    "West",
				},
				{
					PackageID:   "hello",
					FileID:      "enum_in_parenthesis.cx",
					StartOffset: setOffset(376, 40),
					Length:      12,
					LineNumber:  40,
					Type:        "Number",
					Value:       0,
					EnumName:    "First",
				},
				{
					PackageID:   "hello",
					FileID:      "enum_in_parenthesis.cx",
					StartOffset: setOffset(397, 41),
					Length:      6,
					LineNumber:  41,
					Type:        "Number",
					Value:       1,
					EnumName:    "Second",
				},
			},
		},
		{
			scenario: "Has Enums 2",
			testDir:  "./test_files/test_2.cx",
			wantEnums: []declaration_extractor.EnumDeclaration{
				{
					PackageID:   "test_2",
					FileID:      "test_2.cx",
					StartOffset: setOffset(161, 21),
					Length:      13,
					LineNumber:  21,
					Type:        "string",
					Value:       0,
					EnumName:    "Spring",
				},
				{
					PackageID:   "test_2",
					FileID:      "test_2.cx",
					StartOffset: setOffset(183, 22),
					Length:      6,
					LineNumber:  22,
					Type:        "string",
					Value:       1,
					EnumName:    "Summer",
				},
				{
					PackageID:   "test_2",
					FileID:      "test_2.cx",
					StartOffset: setOffset(191, 23),
					Length:      6,
					LineNumber:  23,
					Type:        "string",
					Value:       2,
					EnumName:    "Autumn",
				},
				{
					PackageID:   "test_2",
					FileID:      "test_2.cx",
					StartOffset: setOffset(199, 24),
					Length:      6,
					LineNumber:  24,
					Type:        "string",
					Value:       3,
					EnumName:    "Winter",
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.scenario, func(t *testing.T) {
			srcBytes, err := os.ReadFile(tc.testDir)
			ReplaceCommentsWithWhitespaces := declaration_extractor.ReplaceCommentsWithWhitespaces(srcBytes)
			fileName := filepath.Base(tc.testDir)
			pkg := declaration_extractor.ExtractPackages(ReplaceCommentsWithWhitespaces)
			if err != nil {
				t.Fatal(err)
			}

			enums, err := declaration_extractor.ExtractEnums(ReplaceCommentsWithWhitespaces, fileName, pkg)
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
		wantStructs []declaration_extractor.StructDeclaration
	}{
		{
			scenario: "Has structs",
			testDir:  "./test_files/test.cx",
			wantStructs: []declaration_extractor.StructDeclaration{
				{
					PackageID:   "hello",
					FileID:      "test.cx",
					StartOffset: setOffset(158, 14),
					Length:      18,
					LineNumber:  14,
					StructName:  "person",
				},
				{
					PackageID:   "hello",
					FileID:      "test.cx",
					StartOffset: setOffset(230, 21),
					Length:      39,
					LineNumber:  21,
					StructName:  "animal",
				},
			},
		},
		{
			scenario: "Has Struct 2",
			testDir:  "./test_files/test_2.cx",
			wantStructs: []declaration_extractor.StructDeclaration{
				{
					PackageID:   "test_2",
					FileID:      "test_2.cx",
					StartOffset: setOffset(89, 15),
					Length:      18,
					LineNumber:  15,
					StructName:  "object",
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.scenario, func(t *testing.T) {
			srcBytes, err := os.ReadFile(tc.testDir)
			ReplaceCommentsWithWhitespaces := declaration_extractor.ReplaceCommentsWithWhitespaces(srcBytes)
			fileName := filepath.Base(tc.testDir)
			pkg := declaration_extractor.ExtractPackages(ReplaceCommentsWithWhitespaces)
			if err != nil {
				t.Fatal(err)
			}

			gotStructs, err := declaration_extractor.ExtractStructs(ReplaceCommentsWithWhitespaces, fileName, pkg)
			if err != nil {
				t.Fatal(err)
			}

			for _, wantStruct := range tc.wantStructs {

				var match bool

				var gotStruct declaration_extractor.StructDeclaration

				for _, gotStruct := range gotStructs {

					if gotStruct.StructName == wantStruct.StructName {
						match = true
						break
					}

				}

				if !match {
					t.Errorf("want struct: %v\n\t%v\ngot:%v\n\t%v", wantStruct, wantStruct.StructFields, gotStruct, gotStruct.StructFields)
				}

			}
		})
	}
}

func TestDeclarationExtraction_ExtractFuncs(t *testing.T) {

	tests := []struct {
		scenario  string
		testDir   string
		wantFuncs []declaration_extractor.FuncDeclaration
	}{
		{
			scenario: "Has funcs",
			testDir:  "./test_files/test.cx",
			wantFuncs: []declaration_extractor.FuncDeclaration{
				{
					PackageID:   "hello",
					FileID:      "test.cx",
					StartOffset: setOffset(195, 18),
					Length:      12,
					LineNumber:  18,
					FuncName:    "main",
				},
				{
					PackageID:   "hello",
					FileID:      "test.cx",
					StartOffset: setOffset(296, 26),
					Length:      19,
					LineNumber:  26,
					FuncName:    "functionTwo",
				},
				{
					PackageID:   "hello",
					FileID:      "test.cx",
					StartOffset: setOffset(436, 44),
					Length:      39,
					LineNumber:  44,
					FuncName:    "functionWithSingleReturn",
				},
			},
		},
		{
			scenario: "test_2",
			testDir:  "./test_files/test_2.cx",
			wantFuncs: []declaration_extractor.FuncDeclaration{
				{
					PackageID:   "test_2",
					FileID:      "test_2.cx",
					StartOffset: setOffset(209, 27),
					Length:      46,
					LineNumber:  27,
					FuncName:    "add",
				},
				{
					PackageID:   "test_2",
					FileID:      "test_2.cx",
					StartOffset: setOffset(299, 32),
					Length:      11,
					LineNumber:  32,
					FuncName:    "main",
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.scenario, func(t *testing.T) {
			srcBytes, err := os.ReadFile(tc.testDir)
			ReplaceCommentsWithWhitespaces := declaration_extractor.ReplaceCommentsWithWhitespaces(srcBytes)
			fileName := filepath.Base(tc.testDir)
			pkg := declaration_extractor.ExtractPackages(ReplaceCommentsWithWhitespaces)
			if err != nil {
				t.Fatal(err)
			}

			funcs, err := declaration_extractor.ExtractFuncs(ReplaceCommentsWithWhitespaces, fileName, pkg)
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
		scenario               string
		testDir                string
		wantReDeclarationError error
	}{
		{
			scenario:               "No Redeclarations",
			testDir:                "./test_files/test.cx",
			wantReDeclarationError: nil,
		},
		{
			scenario:               "Redeclared globals",
			testDir:                "./test_files/redeclaration_global.cx",
			wantReDeclarationError: errors.New("global redeclared"),
		},
		{
			scenario:               "Redeclared enums",
			testDir:                "./test_files/redeclaration_enum.cx",
			wantReDeclarationError: errors.New("enum redeclared"),
		},
		{
			scenario:               "Redeclared structs",
			testDir:                "./test_files/redeclaration_struct.cx",
			wantReDeclarationError: errors.New("struct redeclared"),
		},
		{
			scenario:               "Redeclared funcs",
			testDir:                "./test_files/redeclaration_func.cx",
			wantReDeclarationError: errors.New("func redeclared"),
		},
		{
			scenario:               "No Redeclarations 2",
			testDir:                "./test_files/test_2.cx",
			wantReDeclarationError: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.scenario, func(t *testing.T) {
			srcBytes, err := os.ReadFile(tc.testDir)
			ReplaceCommentsWithWhitespaces := declaration_extractor.ReplaceCommentsWithWhitespaces(srcBytes)
			fileName := filepath.Base(tc.testDir)
			pkg := declaration_extractor.ExtractPackages(ReplaceCommentsWithWhitespaces)
			if err != nil {
				t.Fatal(err)
			}

			globals, err := declaration_extractor.ExtractGlobals(ReplaceCommentsWithWhitespaces, fileName, pkg)
			if err != nil {
				t.Fatal(err)
			}

			enums, err := declaration_extractor.ExtractEnums(ReplaceCommentsWithWhitespaces, fileName, pkg)
			if err != nil {
				t.Fatal(err)
			}

			structs, err := declaration_extractor.ExtractStructs(ReplaceCommentsWithWhitespaces, fileName, pkg)
			if err != nil {
				t.Fatal(err)
			}

			funcs, err := declaration_extractor.ExtractFuncs(ReplaceCommentsWithWhitespaces, fileName, pkg)
			if err != nil {
				t.Fatal(err)
			}

			reDeclarationError := declaration_extractor.ReDeclarationCheck(globals, enums, structs, funcs)

			if reDeclarationError != nil && tc.wantReDeclarationError == nil {
				t.Errorf("want error %v, got %v", tc.wantReDeclarationError, reDeclarationError)
			}

			if reDeclarationError != nil && tc.wantReDeclarationError != nil {
				if reDeclarationError.Error() != tc.wantReDeclarationError.Error() {
					t.Errorf("want error %v, got %v", tc.wantReDeclarationError, reDeclarationError)
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
		{
			scenario: "Has declarations 2",
			testDir:  "./test_files/test_2.cx",
			wantDeclarations: []string{
				"var global string",
				"Spring string",
				"Summer",
				"Autumn",
				"Winter",
				"type object struct",
				"func add(obj *object, name string, number int)",
				"func main()",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.scenario, func(t *testing.T) {
			srcBytes, err := os.ReadFile(tc.testDir)
			ReplaceCommentsWithWhitespaces := declaration_extractor.ReplaceCommentsWithWhitespaces(srcBytes)
			fileName := filepath.Base(tc.testDir)
			pkg := declaration_extractor.ExtractPackages(ReplaceCommentsWithWhitespaces)
			if err != nil {
				t.Fatal(err)
			}

			globals, err := declaration_extractor.ExtractGlobals(ReplaceCommentsWithWhitespaces, fileName, pkg)
			if err != nil {
				t.Fatal(err)
			}

			enums, err := declaration_extractor.ExtractEnums(ReplaceCommentsWithWhitespaces, fileName, pkg)
			if err != nil {
				t.Fatal(err)
			}

			structs, err := declaration_extractor.ExtractStructs(ReplaceCommentsWithWhitespaces, fileName, pkg)
			if err != nil {
				t.Fatal(err)
			}

			funcs, err := declaration_extractor.ExtractFuncs(ReplaceCommentsWithWhitespaces, fileName, pkg)
			if err != nil {
				t.Fatal(err)
			}

			if declaration_extractor.ReDeclarationCheck(globals, enums, structs, funcs) != nil {
				t.Fatal(err)
			}

			declarations := declaration_extractor.GetDeclarations(srcBytes, globals, enums, structs, funcs)

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
			scenario: "Single file 2",
			testDirs: []string{
				"./test_files/test_2.cx",
			},
			wantGlobals: 1,
			wantEnums:   4,
			wantStructs: 1,
			wantFuncs:   2,
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

			Globals, Enums, Structs, Funcs, Err := declaration_extractor.ExtractAllDeclarations(files)

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
