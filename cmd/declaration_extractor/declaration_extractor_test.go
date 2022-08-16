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

func TestDeclarationExtractor_ReplaceCommentsWithWhitespaces(t *testing.T) {

	tests := []struct {
		scenario            string
		testDir             string
		wantCommentReplaced string
	}{
		{
			scenario:            "Has comments",
			testDir:             "./test_files/ReplaceCommentsWithWhitespaces/HasComments.cx",
			wantCommentReplaced: "./test_files/ReplaceCommentsWithWhitespaces/HasCommentsResult.cx",
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

func TestDeclarationExtractor_ExtractGlobals(t *testing.T) {

	tests := []struct {
		scenario    string
		testDir     string
		wantGlobals []declaration_extractor.GlobalDeclaration
		wantErr     error
	}{
		{
			scenario: "Has Globals",
			testDir:  "./test_files/ExtractGlobals/HasGlobals.cx",
			wantGlobals: []declaration_extractor.GlobalDeclaration{
				{
					PackageID:          "main",
					FileID:             "./test_files/ExtractGlobals/HasGlobals.cx",
					StartOffset:        setOffset(222, 15),
					Length:             12,
					LineNumber:         15,
					GlobalVariableName: "fooV",
				},
				{
					PackageID:          "main",
					FileID:             "./test_files/ExtractGlobals/HasGlobals.cx",
					StartOffset:        setOffset(253, 16),
					Length:             12,
					LineNumber:         16,
					GlobalVariableName: "fooA",
				},
				{
					PackageID:          "main",
					FileID:             "./test_files/ExtractGlobals/HasGlobals.cx",
					StartOffset:        setOffset(270, 17),
					Length:             12,
					LineNumber:         17,
					GlobalVariableName: "fooR",
				},
			},
			wantErr: nil,
		},
		{
			scenario: "Has Globals 2",
			testDir:  "./test_files/ExtractGlobals/HasGlobals2.cx",
			wantGlobals: []declaration_extractor.GlobalDeclaration{
				{
					PackageID:          "main",
					FileID:             "./test_files/ExtractGlobals/HasGlobals2.cx",
					StartOffset:        setOffset(153, 12),
					Length:             12,
					LineNumber:         12,
					GlobalVariableName: "fooV",
				},
			},
			wantErr: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.scenario, func(t *testing.T) {
			srcBytes, err := os.ReadFile(tc.testDir)
			ReplaceCommentsWithWhitespaces := declaration_extractor.ReplaceCommentsWithWhitespaces(srcBytes)
			fileName := tc.testDir
			if err != nil {
				t.Fatal(err)
			}

			gotGlobals, gotErr := declaration_extractor.ExtractGlobals(ReplaceCommentsWithWhitespaces, fileName)
			if gotErr != tc.wantErr {
				t.Errorf("want err %v, got %v", tc.wantErr, gotErr)
			}

			for _, wantGlobal := range tc.wantGlobals {

				var match bool = false
				var gotGlobalF declaration_extractor.GlobalDeclaration

				for _, gotGlobal := range gotGlobals {

					if gotGlobal.GlobalVariableName == wantGlobal.GlobalVariableName {

						if gotGlobal == wantGlobal {
							match = true
						}

						gotGlobalF = gotGlobal

						break
					}

				}

				if !match {
					t.Errorf("want global %v, got %v", wantGlobal, gotGlobalF)
				}

			}
		})
	}
}

func TestDeclarationExtractor_ExtractEnums(t *testing.T) {

	tests := []struct {
		scenario  string
		testDir   string
		wantEnums []declaration_extractor.EnumDeclaration
		wantErr   error
	}{
		{
			scenario: "HasEnums.cx",
			testDir:  "./test_files/ExtractEnums/HasEnums.cx",
			wantEnums: []declaration_extractor.EnumDeclaration{
				{
					PackageID:   "main",
					FileID:      "./test_files/ExtractEnums/HasEnums.cx",
					StartOffset: setOffset(26, 4),
					Length:      10,
					LineNumber:  4,
					Type:        "int",
					Value:       0,
					EnumName:    "Summer",
				},
				{
					PackageID:   "main",
					FileID:      "./test_files/ExtractEnums/HasEnums.cx",
					StartOffset: setOffset(45, 5),
					Length:      6,
					LineNumber:  5,
					Type:        "int",
					Value:       1,
					EnumName:    "Autumn",
				},
				{
					PackageID:   "main",
					FileID:      "./test_files/ExtractEnums/HasEnums.cx",
					StartOffset: setOffset(53, 6),
					Length:      6,
					LineNumber:  6,
					Type:        "int",
					Value:       2,
					EnumName:    "Winter",
				},
				{
					PackageID:   "main",
					FileID:      "./test_files/ExtractEnums/HasEnums.cx",
					StartOffset: setOffset(61, 7),
					Length:      6,
					LineNumber:  7,
					Type:        "int",
					Value:       3,
					EnumName:    "Spring",
				},
				{
					PackageID:   "main",
					FileID:      "./test_files/ExtractEnums/HasEnums.cx",
					StartOffset: setOffset(83, 11),
					Length:      6,
					LineNumber:  11,
					Type:        "",
					Value:       0,
					EnumName:    "Apples",
				},
				{
					PackageID:   "main",
					FileID:      "./test_files/ExtractEnums/HasEnums.cx",
					StartOffset: setOffset(99, 12),
					Length:      7,
					LineNumber:  12,
					Type:        "",
					Value:       1,
					EnumName:    "Oranges",
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.scenario, func(t *testing.T) {
			srcBytes, err := os.ReadFile(tc.testDir)
			ReplaceCommentsWithWhitespaces := declaration_extractor.ReplaceCommentsWithWhitespaces(srcBytes)

			if err != nil {
				t.Fatal(err)
			}

			gotEnums, gotErr := declaration_extractor.ExtractEnums(ReplaceCommentsWithWhitespaces, tc.testDir)

			for _, wantEnum := range tc.wantEnums {
				var match bool
				var gotEnumF declaration_extractor.EnumDeclaration
				for _, gotEnum := range gotEnums {
					if gotEnum.EnumName == wantEnum.EnumName {
						if gotEnum == wantEnum {
							match = true
						}
						gotEnumF = gotEnum
						break
					}
				}

				if !match {
					t.Errorf("want enum %v, got %v", wantEnum, gotEnumF)
				}
			}

			if gotErr != tc.wantErr {
				t.Errorf("want err %v, got %v", tc.wantErr, gotErr)
			}
		})
	}
}

func TestDeclarationExtractor_ExtractStructs(t *testing.T) {

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
					StartOffset: setOffset(152, 14),
					Length:      18,
					LineNumber:  14,
					StructName:  "person",
					StructFields: []*declaration_extractor.StructField{
						{
							StructFieldName: "name",
							StartOffset:     setOffset(174, 15),
							Length:          8,
							LineNumber:      15,
						},
					},
				},
				{
					PackageID:   "hello",
					FileID:      "test.cx",
					StartOffset: setOffset(230, 21),
					Length:      39,
					LineNumber:  21,
					StructName:  "animal",
					StructFields: []*declaration_extractor.StructField{
						{},
						{},
					},
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
			if err != nil {
				t.Fatal(err)
			}

			gotStructs, err := declaration_extractor.ExtractStructs(ReplaceCommentsWithWhitespaces, fileName)
			if err != nil {
				t.Fatal(err)
			}

			for _, wantStruct := range tc.wantStructs {

				var match bool

				var gotStructF declaration_extractor.StructDeclaration

				for _, gotStruct := range gotStructs {

					gotStructF = gotStruct

					if gotStruct.StructName == wantStruct.StructName {
						if gotStruct.FileID == wantStruct.FileID &&
							gotStruct.Length == wantStruct.Length &&
							gotStruct.LineNumber == wantStruct.LineNumber &&
							gotStruct.PackageID == wantStruct.PackageID {

							var fieldMatch bool = true

							for k, wantField := range wantStruct.StructFields {

								if gotStruct.StructFields[k] != wantField {
									fieldMatch = false
									break
								}

							}

							if fieldMatch {
								match = true
								break
							}
							t.Error(gotStruct)

						}
						break
					}

				}

				if !match {
					t.Errorf("want struct %v", wantStruct)
					for _, wantField := range wantStruct.StructFields {
						t.Error(wantField)
					}
					t.Errorf("got %v", gotStructF)
					for _, gotField := range gotStructF.StructFields {
						t.Error(gotField)
					}
				}

			}
		})
	}
}

func TestDeclarationExtractor_ExtractFuncs(t *testing.T) {

	tests := []struct {
		scenario  string
		testDir   string
		wantFuncs []declaration_extractor.FuncDeclaration
		wantErr   error
	}{
		{
			scenario: "Has funcs",
			testDir:  "./test_files/ExtractFuncs/HasFuncs.cx",
			wantFuncs: []declaration_extractor.FuncDeclaration{
				{
					PackageID:   "main",
					FileID:      "./test_files/ExtractFuncs/HasFuncs.cx",
					StartOffset: setOffset(322, 20),
					Length:      12,
					LineNumber:  20,
					FuncName:    "main",
				},
				{
					PackageID:   "main",
					FileID:      "./test_files/ExtractFuncs/HasFuncs.cx",
					StartOffset: setOffset(14, 3),
					Length:      53,
					LineNumber:  3,
					FuncName:    "addition",
				},
				{
					PackageID:   "main",
					FileID:      "./test_files/ExtractFuncs/HasFuncs.cx",
					StartOffset: setOffset(104, 7),
					Length:      50,
					LineNumber:  7,
					FuncName:    "minus",
				},
				{
					PackageID:   "main",
					FileID:      "./test_files/ExtractFuncs/HasFuncs.cx",
					StartOffset: setOffset(226, 15),
					Length:      29,
					LineNumber:  15,
					FuncName:    "printName",
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.scenario, func(t *testing.T) {
			srcBytes, err := os.ReadFile(tc.testDir)
			ReplaceCommentsWithWhitespaces := declaration_extractor.ReplaceCommentsWithWhitespaces(srcBytes)
			fileName := tc.testDir
			if err != nil {
				t.Fatal(err)
			}

			gotFuncs, gotErr := declaration_extractor.ExtractFuncs(ReplaceCommentsWithWhitespaces, fileName)

			for _, wantFunc := range tc.wantFuncs {
				var match bool
				var gotFuncF declaration_extractor.FuncDeclaration
				for _, gotFunc := range gotFuncs {
					if gotFunc.FuncName == wantFunc.FuncName {
						if gotFunc == wantFunc {
							match = true
						}
						gotFuncF = gotFunc
						break
					}
				}
				if !match {
					t.Errorf("want func %v, got %v", wantFunc, gotFuncF)
				}
			}

			if gotErr != tc.wantErr {
				t.Errorf("want error %v, got %v", tc.wantErr, gotErr)
			}
		})
	}
}

func TestDeclarationExtractor_ReDeclarationCheck(t *testing.T) {

	tests := []struct {
		scenario               string
		testDir                string
		wantReDeclarationError error
	}{
		{
			scenario:               "No Redeclaration",
			testDir:                "./test_files/ReDeclarationCheck/NoRedeclaration.cx",
			wantReDeclarationError: nil,
		},
		{
			scenario:               "Redeclared globals",
			testDir:                "./test_files/ReDeclarationCheck/RedeclaredGlobals.cx",
			wantReDeclarationError: errors.New("global redeclared"),
		},
		{
			scenario:               "Redeclared enums",
			testDir:                "./test_files/ReDeclarationCheck/RedeclaredEnums.cx",
			wantReDeclarationError: errors.New("enum redeclared"),
		},
		{
			scenario:               "Redeclared structs",
			testDir:                "./test_files/ReDeclarationCheck/RedeclaredStructs.cx",
			wantReDeclarationError: errors.New("struct redeclared"),
		},
		{
			scenario:               "Redeclared funcs",
			testDir:                "./test_files/ReDeclarationCheck/RedeclaredFuncs.cx",
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
			if err != nil {
				t.Fatal(err)
			}

			globals, err := declaration_extractor.ExtractGlobals(ReplaceCommentsWithWhitespaces, fileName)
			if err != nil {
				t.Fatal(err)
			}

			enums, err := declaration_extractor.ExtractEnums(ReplaceCommentsWithWhitespaces, fileName)
			if err != nil {
				t.Fatal(err)
			}

			structs, err := declaration_extractor.ExtractStructs(ReplaceCommentsWithWhitespaces, fileName)
			if err != nil {
				t.Fatal(err)
			}

			funcs, err := declaration_extractor.ExtractFuncs(ReplaceCommentsWithWhitespaces, fileName)
			if err != nil {
				t.Fatal(err)
			}

			gotReDeclarationError := declaration_extractor.ReDeclarationCheck(globals, enums, structs, funcs)

			if (gotReDeclarationError != nil && tc.wantReDeclarationError == nil) ||
				(gotReDeclarationError == nil && tc.wantReDeclarationError != nil) {
				t.Errorf("want error %v, got %v", tc.wantReDeclarationError, gotReDeclarationError)
			}

			if gotReDeclarationError != nil && tc.wantReDeclarationError != nil {
				if gotReDeclarationError.Error() != tc.wantReDeclarationError.Error() {
					t.Errorf("want error %v, got %v", tc.wantReDeclarationError, gotReDeclarationError)
				}
			}
		})
	}
}

func TestDeclarationExtractor_GetDeclarations(t *testing.T) {

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
			if err != nil {
				t.Fatal(err)
			}

			globals, err := declaration_extractor.ExtractGlobals(ReplaceCommentsWithWhitespaces, fileName)
			if err != nil {
				t.Fatal(err)
			}

			enums, err := declaration_extractor.ExtractEnums(ReplaceCommentsWithWhitespaces, fileName)
			if err != nil {
				t.Fatal(err)
			}

			structs, err := declaration_extractor.ExtractStructs(ReplaceCommentsWithWhitespaces, fileName)
			if err != nil {
				t.Fatal(err)
			}

			funcs, err := declaration_extractor.ExtractFuncs(ReplaceCommentsWithWhitespaces, fileName)
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

func TestDeclarationExtractor_ExtractAllDeclarations(t *testing.T) {

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
