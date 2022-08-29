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
					Length:             30,
					LineNumber:         15,
					GlobalVariableName: "fooV",
				},
				{
					PackageID:          "main",
					FileID:             "./test_files/ExtractGlobals/HasGlobals.cx",
					StartOffset:        setOffset(253, 16),
					Length:             16,
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
					Length:             56,
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
					Length:      17,
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
					Length:      11,
					LineNumber:  11,
					Type:        "",
					Value:       0,
					EnumName:    "Apples",
				},
				{
					PackageID:   "main",
					FileID:      "./test_files/ExtractEnums/HasEnums.cx",
					StartOffset: setOffset(99, 12),
					Length:      11,
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
			if err != nil {
				t.Fatal(err)
			}

			ReplaceCommentsWithWhitespaces := declaration_extractor.ReplaceCommentsWithWhitespaces(srcBytes)

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

func TestDeclarationExtractor_ExtractTypeDefinitions(t *testing.T) {

	tests := []struct {
		scenario            string
		testDir             string
		wantTypeDefinitions []declaration_extractor.TypeDefinitionDeclaration
		wantErr             error
	}{
		{
			scenario: "Has Type Definitions",
			testDir:  "./test_files/ExtractTypeDefinitions/HasTypeDefinitions.cx",
			wantTypeDefinitions: []declaration_extractor.TypeDefinitionDeclaration{
				{
					PackageID:          "main",
					FileID:             "./test_files/ExtractTypeDefinitions/HasTypeDefinitions.cx",
					StartOffset:        14,
					Length:             18,
					LineNumber:         3,
					TypeDefinitionName: "Direction",
				},
				{
					PackageID:          "main",
					FileID:             "./test_files/ExtractTypeDefinitions/HasTypeDefinitions.cx",
					StartOffset:        100,
					Length:             15,
					LineNumber:         12,
					TypeDefinitionName: "Season",
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.scenario, func(t *testing.T) {
			srcBytes, err := os.ReadFile(tc.testDir)
			if err != nil {
				t.Fatal(err)
			}

			ReplaceCommentsWithWhitespaces := declaration_extractor.ReplaceCommentsWithWhitespaces(srcBytes)

			gotTypeDefinitions, gotErr := declaration_extractor.ExtractTypeDefinitions(ReplaceCommentsWithWhitespaces, tc.testDir)

			for _, wantTypeDef := range tc.wantTypeDefinitions {

				wantTypeDef.StartOffset = setOffset(wantTypeDef.StartOffset, wantTypeDef.LineNumber)

				var gotTypeDefF declaration_extractor.TypeDefinitionDeclaration
				var match bool

				for _, gotTypeDef := range gotTypeDefinitions {

					if gotTypeDef.TypeDefinitionName == wantTypeDef.TypeDefinitionName {

						gotTypeDefF = gotTypeDef
						if gotTypeDef == wantTypeDef {
							match = true
						}

						break
					}
				}

				if !match {
					t.Errorf("want type definition %v, got %v", wantTypeDef, gotTypeDefF)
				}
			}

			if (gotErr != nil && tc.wantErr == nil) ||
				(gotErr == nil && tc.wantErr != nil) {
				t.Errorf("want error %v, got %v", tc.wantErr, gotErr)
			}

			if gotErr != nil && tc.wantErr != nil {
				if gotErr.Error() != tc.wantErr.Error() {
					t.Errorf("want error %v, got %v", tc.wantErr, gotErr)
				}
			}

		})
	}
}

func TestDeclarationExtractor_ExtractStructs(t *testing.T) {

	tests := []struct {
		scenario    string
		testDir     string
		wantStructs []declaration_extractor.StructDeclaration
		wantErr     error
	}{
		{
			scenario: "Has structs",
			testDir:  "./test_files/ExtractStructs/HasStructs.cx",
			wantStructs: []declaration_extractor.StructDeclaration{
				{
					PackageID:   "main",
					FileID:      "./test_files/ExtractStructs/HasStructs.cx",
					StartOffset: 58,
					Length:      19,
					LineNumber:  5,
					StructName:  "Point",
					StructFields: []*declaration_extractor.StructField{
						{
							StructFieldName: "x",
							StartOffset:     79,
							Length:          5,
							LineNumber:      6,
						},
						{
							StructFieldName: "y",
							StartOffset:     86,
							Length:          5,
							LineNumber:      7,
						},
					},
				},
				{
					PackageID:   "main",
					FileID:      "./test_files/ExtractStructs/HasStructs.cx",
					StartOffset: 121,
					Length:      21,
					LineNumber:  11,
					StructName:  "Strings",
					StructFields: []*declaration_extractor.StructField{
						{
							StructFieldName: "lBound",
							StartOffset:     144,
							Length:          10,
							LineNumber:      12,
						},
						{
							StructFieldName: "string",
							StartOffset:     156,
							Length:          10,
							LineNumber:      13,
						},
						{
							StructFieldName: "stringA",
							StartOffset:     168,
							Length:          13,
							LineNumber:      14,
						},
						{
							StructFieldName: "rBound",
							StartOffset:     183,
							Length:          10,
							LineNumber:      15,
						},
					},
				},
			},
		},
		{
			scenario: "Has Struct 2",
			testDir:  "./test_files/ExtractStructs/HasStructs2.cx",
			wantStructs: []declaration_extractor.StructDeclaration{
				{
					PackageID:   "main",
					FileID:      "./test_files/ExtractStructs/HasStructs2.cx",
					StartOffset: 14,
					Length:      19,
					LineNumber:  3,
					StructName:  "Point",
					StructFields: []*declaration_extractor.StructField{
						{
							StructFieldName: "x",
							StartOffset:     35,
							Length:          5,
							LineNumber:      4,
						},
						{
							StructFieldName: "y",
							StartOffset:     42,
							Length:          5,
							LineNumber:      5,
						},
					},
				},
				{
					PackageID:   "main",
					FileID:      "./test_files/ExtractStructs/HasStructs2.cx",
					StartOffset: 51,
					Length:      20,
					LineNumber:  8,
					StructName:  "Canvas",
					StructFields: []*declaration_extractor.StructField{
						{
							StructFieldName: "points",
							StartOffset:     73,
							Length:          16,
							LineNumber:      9,
						},
					},
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

			gotStructs, gotErr := declaration_extractor.ExtractStructs(ReplaceCommentsWithWhitespaces, tc.testDir)
			if err != nil {
				t.Fatal(err)
			}

			for _, wantStruct := range tc.wantStructs {

				wantStruct.StartOffset = setOffset(wantStruct.StartOffset, wantStruct.LineNumber)

				var match bool

				var gotStructF declaration_extractor.StructDeclaration

				for _, gotStruct := range gotStructs {

					if gotStruct.StructName == wantStruct.StructName {

						gotStructF = gotStruct

						if gotStruct.FileID == wantStruct.FileID &&
							gotStruct.StartOffset == wantStruct.StartOffset &&
							gotStruct.Length == wantStruct.Length &&
							gotStruct.LineNumber == wantStruct.LineNumber &&
							gotStruct.PackageID == wantStruct.PackageID {

							var fieldMatch bool = true

							for k, wantField := range wantStruct.StructFields {

								if *gotStruct.StructFields[k] != *wantField {
									fieldMatch = false
									break
								}

							}

							if fieldMatch {
								match = true
							}
							break
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

			if gotErr != tc.wantErr {
				t.Errorf("want err %v, got %v", tc.wantErr, gotErr)
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
					Length:      14,
					LineNumber:  20,
					FuncName:    "main",
				},
				{
					PackageID:   "main",
					FileID:      "./test_files/ExtractFuncs/HasFuncs.cx",
					StartOffset: setOffset(14, 3),
					Length:      55,
					LineNumber:  3,
					FuncName:    "addition",
				},
				{
					PackageID:   "main",
					FileID:      "./test_files/ExtractFuncs/HasFuncs.cx",
					StartOffset: setOffset(104, 7),
					Length:      52,
					LineNumber:  7,
					FuncName:    "minus",
				},
				{
					PackageID:   "main",
					FileID:      "./test_files/ExtractFuncs/HasFuncs.cx",
					StartOffset: setOffset(226, 15),
					Length:      31,
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
			scenario:               "Redeclared global",
			testDir:                "./test_files/ReDeclarationCheck/RedeclaredGlobal.cx",
			wantReDeclarationError: errors.New("global redeclared"),
		},
		{
			scenario:               "Redeclared enum",
			testDir:                "./test_files/ReDeclarationCheck/RedeclaredEnum.cx",
			wantReDeclarationError: errors.New("enum redeclared"),
		},
		{
			scenario:               "Redeclared struct",
			testDir:                "./test_files/ReDeclarationCheck/RedeclaredStruct.cx",
			wantReDeclarationError: errors.New("struct redeclared"),
		},
		{
			scenario:               "Redeclared struct field",
			testDir:                "./test_files/ReDeclarationCheck/RedeclaredStructField.cx",
			wantReDeclarationError: errors.New("struct field redeclared"),
		},
		{
			scenario:               "Redeclared func",
			testDir:                "./test_files/ReDeclarationCheck/RedeclaredFunc.cx",
			wantReDeclarationError: errors.New("func redeclared"),
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

			typeDefinitions, err := declaration_extractor.ExtractTypeDefinitions(ReplaceCommentsWithWhitespaces, fileName)
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

			gotReDeclarationError := declaration_extractor.ReDeclarationCheck(globals, enums, typeDefinitions, structs, funcs)

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

			typeDefinitions, err := declaration_extractor.ExtractTypeDefinitions(ReplaceCommentsWithWhitespaces, fileName)
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

			if declaration_extractor.ReDeclarationCheck(globals, enums, typeDefinitions, structs, funcs) != nil {
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

// Regex is slower because there’s a lot of conditionals and it’s a long pattern to match.
// Tokenizing is faster but isn’t as accurate.
func BenchmarkDeclarationExtractor_ExtractFuncs(b *testing.B) {
	benchmarks := []struct {
		scenario string
		testDir  string
	}{
		{scenario: "regular funcs", testDir: "./test_files/ExtractFuncs/HasFuncs.cx"},
	}

	for _, bm := range benchmarks {
		b.Run(bm.scenario, func(b *testing.B) {
			srcBytes, err := os.ReadFile(bm.testDir)
			if err != nil {
				b.Fatal(err)
			}
			for n := 0; n < b.N; n++ {
				declaration_extractor.ExtractFuncs(srcBytes, bm.testDir)
			}
		})
	}
}
