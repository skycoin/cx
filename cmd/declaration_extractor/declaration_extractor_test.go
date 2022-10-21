package declaration_extractor_test

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/skycoin/cx/cmd/declaration_extractor"
	"github.com/skycoin/cx/cmd/packageloader/file_output"
	"github.com/skycoin/cx/cmd/packageloader/loader"
	"github.com/skycoin/cx/cxparser/actions"
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

func TestDeclarationExtractor_ReplaceStringContentsWithWhitespaces(t *testing.T) {

	tests := []struct {
		scenario                   string
		testDir                    string
		wantStringContentsReplaced string
		wantErr                    error
	}{
		{
			scenario:                   "Has String",
			testDir:                    "./test_files/ReplaceStringContentsWithWhitespaces/HasString.cx",
			wantStringContentsReplaced: "./test_files/ReplaceStringContentsWithWhitespaces/HasString.cxStringContentsReplaced.cx",
			wantErr:                    nil,
		},
		{
			scenario:                   "Syntax Error",
			testDir:                    "./test_files/ReplaceStringContentsWithWhitespaces/SyntaxError.cx",
			wantStringContentsReplaced: "./test_files/ReplaceStringContentsWithWhitespaces/SyntaxError.cxStringContentsReplaced.cx",
			wantErr:                    errors.New("9: syntax error: quote not terminated"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.scenario, func(t *testing.T) {
			srcBytes, err := os.ReadFile(tc.testDir)
			if err != nil {
				t.Fatal(err)
			}
			wantBytes, err := os.ReadFile(tc.wantStringContentsReplaced)
			if err != nil {
				t.Fatal(err)
			}
			stringContentsReplaced, gotErr := declaration_extractor.ReplaceStringContentsWithWhitespaces(srcBytes)

			if len(srcBytes) != len(stringContentsReplaced) {
				t.Errorf("Length not the same: orginal %vbytes, replaced %vbytes", len(srcBytes), len(stringContentsReplaced))
			}

			srcLines := bytes.Count(srcBytes, []byte("\n")) + 1
			newLines := bytes.Count(stringContentsReplaced, []byte("\n")) + 1

			if srcLines != newLines {
				t.Errorf("Lines not equal: original %vlines, new %vlines", srcLines, newLines)
			}

			if string(stringContentsReplaced) != string(wantBytes) {
				t.Errorf("want string contents replaced\n%v\ngot\n%v", string(wantBytes), string(stringContentsReplaced))
				file, err := os.Create(tc.testDir + "gotStringContentsReplaced.cx")
				if err != nil {
					t.Fatal(err)
				}
				file.Write(stringContentsReplaced)
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

func TestDeclarationExtractor_ExtractImports(t *testing.T) {

	tests := []struct {
		scenario    string
		testDir     string
		wantImports []declaration_extractor.ImportDeclaration
		wantErr     error
	}{
		{
			scenario: "Has Imports",
			testDir:  "./test_files/ExtractImports/HasImports.cx",
			wantImports: []declaration_extractor.ImportDeclaration{
				{
					PackageID:  "bar",
					FileID:     "./test_files/ExtractImports/HasImports.cx",
					LineNumber: 34,
					ImportName: "foo",
				},
				{
					PackageID:  "bar",
					FileID:     "./test_files/ExtractImports/HasImports.cx",
					LineNumber: 35,
					ImportName: "nosal",
				},
				{
					PackageID:  "main",
					FileID:     "./test_files/ExtractImports/HasImports.cx",
					LineNumber: 57,
					ImportName: "foo",
				},
				{
					PackageID:  "main",
					FileID:     "./test_files/ExtractImports/HasImports.cx",
					LineNumber: 58,
					ImportName: "bar",
				},
			},
		},
		{
			scenario: "Has Imports",
			testDir:  "./test_files/ExtractImports/PackageError.cx",
			wantImports: []declaration_extractor.ImportDeclaration{
				{
					PackageID:  "bar",
					FileID:     "./test_files/ExtractImports/PackageError.cx",
					LineNumber: 34,
					ImportName: "foo",
				},
				{
					PackageID:  "bar",
					FileID:     "./test_files/ExtractImports/PackageError.cx",
					LineNumber: 35,
					ImportName: "nosal",
				},
			},
			wantErr: errors.New("PackageError.cx:54: syntax error: package declaration"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.scenario, func(t *testing.T) {
			srcBytes, err := os.ReadFile(tc.testDir)
			if err != nil {
				t.Fatal(err)
			}

			ReplaceCommentsWithWhitespaces := declaration_extractor.ReplaceCommentsWithWhitespaces(srcBytes)
			if err != nil {
				t.Fatal(filepath.Base(tc.testDir), err)
			}

			gotImports, gotErr := declaration_extractor.ExtractImports(ReplaceCommentsWithWhitespaces, tc.testDir)

			for _, wantImport := range tc.wantImports {

				var match bool = false
				var gotImportF declaration_extractor.ImportDeclaration

				for _, gotImport := range gotImports {
					gotImportF = gotImport
					if gotImport.ImportName == wantImport.ImportName &&
						gotImport.PackageID == wantImport.PackageID {
						if gotImport == wantImport {
							match = true
						}
						break
					}
				}

				if !match {
					t.Errorf("want Import %v, got %v", wantImport, gotImportF)
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
					StartOffset:        222,
					Length:             30,
					LineNumber:         15,
					GlobalVariableName: "fooV",
				},
				{
					PackageID:          "main",
					FileID:             "./test_files/ExtractGlobals/HasGlobals.cx",
					StartOffset:        253,
					Length:             16,
					LineNumber:         16,
					GlobalVariableName: "fooA",
				},
				{
					PackageID:          "main",
					FileID:             "./test_files/ExtractGlobals/HasGlobals.cx",
					StartOffset:        270,
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
					StartOffset:        153,
					Length:             56,
					LineNumber:         12,
					GlobalVariableName: "fooV",
				},
			},
			wantErr: nil,
		},
		{
			scenario:    "Package Error",
			testDir:     "./test_files/ExtractGlobals/PackageError.cx",
			wantGlobals: []declaration_extractor.GlobalDeclaration{},
			wantErr:     errors.New("PackageError.cx:10: syntax error: package declaration"),
		},
		{
			scenario: "Syntax Error",
			testDir:  "./test_files/ExtractGlobals/SyntaxError.cx",
			wantGlobals: []declaration_extractor.GlobalDeclaration{
				{
					PackageID:          "main",
					FileID:             "./test_files/ExtractGlobals/SyntaxError.cx",
					StartOffset:        153,
					Length:             56,
					LineNumber:         12,
					GlobalVariableName: "fooV",
				},
			},
			wantErr: errors.New("SyntaxError.cx:23: syntax error: global declaration"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.scenario, func(t *testing.T) {
			srcBytes, err := os.ReadFile(tc.testDir)
			if err != nil {
				t.Fatal(err)
			}

			ReplaceCommentsWithWhitespaces := declaration_extractor.ReplaceCommentsWithWhitespaces(srcBytes)
			ReplaceStringContentsWithWhitespaces, err := declaration_extractor.ReplaceStringContentsWithWhitespaces(ReplaceCommentsWithWhitespaces)
			if err != nil {
				t.Fatal(filepath.Base(tc.testDir), err)
			}

			gotGlobals, gotErr := declaration_extractor.ExtractGlobals(ReplaceStringContentsWithWhitespaces, tc.testDir)

			for _, wantGlobal := range tc.wantGlobals {

				wantGlobal.StartOffset = setOffset(wantGlobal.StartOffset, wantGlobal.LineNumber)

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

func TestDeclarationExtractor_ExtractEnums(t *testing.T) {

	tests := []struct {
		scenario  string
		testDir   string
		wantEnums []declaration_extractor.EnumDeclaration
		wantErr   error
	}{
		{
			scenario: "Has Enums",
			testDir:  "./test_files/ExtractEnums/HasEnums.cx",
			wantEnums: []declaration_extractor.EnumDeclaration{
				{
					PackageID:   "main",
					FileID:      "./test_files/ExtractEnums/HasEnums.cx",
					StartOffset: 26,
					Length:      17,
					LineNumber:  4,
					Type:        "int",
					Value:       0,
					EnumName:    "Summer",
				},
				{
					PackageID:   "main",
					FileID:      "./test_files/ExtractEnums/HasEnums.cx",
					StartOffset: 45,
					Length:      6,
					LineNumber:  5,
					Type:        "int",
					Value:       1,
					EnumName:    "Autumn",
				},
				{
					PackageID:   "main",
					FileID:      "./test_files/ExtractEnums/HasEnums.cx",
					StartOffset: 53,
					Length:      6,
					LineNumber:  6,
					Type:        "int",
					Value:       2,
					EnumName:    "Winter",
				},
				{
					PackageID:   "main",
					FileID:      "./test_files/ExtractEnums/HasEnums.cx",
					StartOffset: 61,
					Length:      6,
					LineNumber:  7,
					Type:        "int",
					Value:       3,
					EnumName:    "Spring",
				},
				{
					PackageID:   "main",
					FileID:      "./test_files/ExtractEnums/HasEnums.cx",
					StartOffset: 83,
					Length:      11,
					LineNumber:  11,
					Type:        "",
					Value:       0,
					EnumName:    "Apples",
				},
				{
					PackageID:   "main",
					FileID:      "./test_files/ExtractEnums/HasEnums.cx",
					StartOffset: 99,
					Length:      11,
					LineNumber:  12,
					Type:        "",
					Value:       1,
					EnumName:    "Oranges",
				},
			},
		},
		{
			scenario:  "Package Error",
			testDir:   "./test_files/ExtractEnums/PackageError.cx",
			wantEnums: []declaration_extractor.EnumDeclaration{},
			wantErr:   errors.New("PackageError.cx:1: syntax error: package declaration"),
		},
		{
			scenario: "Syntax Error",
			testDir:  "./test_files/ExtractEnums/SyntaxError.cx",
			wantEnums: []declaration_extractor.EnumDeclaration{
				{
					PackageID:   "main",
					FileID:      "./test_files/ExtractEnums/SyntaxError.cx",
					StartOffset: 26,
					Length:      17,
					LineNumber:  4,
					Type:        "int",
					Value:       0,
					EnumName:    "Summer",
				},
				{
					PackageID:   "main",
					FileID:      "./test_files/ExtractEnums/SyntaxError.cx",
					StartOffset: 45,
					Length:      6,
					LineNumber:  5,
					Type:        "int",
					Value:       1,
					EnumName:    "Autumn",
				},
				{
					PackageID:   "main",
					FileID:      "./test_files/ExtractEnums/SyntaxError.cx",
					StartOffset: 53,
					Length:      6,
					LineNumber:  6,
					Type:        "int",
					Value:       2,
					EnumName:    "Winter",
				},
				{
					PackageID:   "main",
					FileID:      "./test_files/ExtractEnums/SyntaxError.cx",
					StartOffset: 61,
					Length:      6,
					LineNumber:  7,
					Type:        "int",
					Value:       3,
					EnumName:    "Spring",
				},
				{
					PackageID:   "main",
					FileID:      "./test_files/ExtractEnums/SyntaxError.cx",
					StartOffset: 83,
					Length:      11,
					LineNumber:  11,
					Type:        "",
					Value:       0,
					EnumName:    "Apples",
				},
				{
					PackageID:   "main",
					FileID:      "./test_files/ExtractEnums/SyntaxError.cx",
					StartOffset: 99,
					Length:      11,
					LineNumber:  12,
					Type:        "",
					Value:       1,
					EnumName:    "Oranges",
				},
			},
			wantErr: errors.New("SyntaxError.cx:13: syntax error: enum declaration"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.scenario, func(t *testing.T) {
			srcBytes, err := os.ReadFile(tc.testDir)
			if err != nil {
				t.Fatal(err)
			}

			ReplaceCommentsWithWhitespaces := declaration_extractor.ReplaceCommentsWithWhitespaces(srcBytes)
			ReplaceStringContentsWithWhitespaces, err := declaration_extractor.ReplaceStringContentsWithWhitespaces(ReplaceCommentsWithWhitespaces)
			if err != nil {
				t.Fatal(filepath.Base(tc.testDir), err)
			}

			gotEnums, gotErr := declaration_extractor.ExtractEnums(ReplaceStringContentsWithWhitespaces, tc.testDir)

			for _, wantEnum := range tc.wantEnums {

				wantEnum.StartOffset = setOffset(wantEnum.StartOffset, wantEnum.LineNumber)

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
		{
			scenario:            "Package Error",
			testDir:             "./test_files/ExtractTypeDefinitions/PackageError.cx",
			wantTypeDefinitions: []declaration_extractor.TypeDefinitionDeclaration{},
			wantErr:             errors.New("PackageError.cx:1: syntax error: package declaration"),
		},
		{
			scenario: "Syntax Error",
			testDir:  "./test_files/ExtractTypeDefinitions/SyntaxError.cx",
			wantTypeDefinitions: []declaration_extractor.TypeDefinitionDeclaration{
				{
					PackageID:          "main",
					FileID:             "./test_files/ExtractTypeDefinitions/SyntaxError.cx",
					StartOffset:        14,
					Length:             18,
					LineNumber:         3,
					TypeDefinitionName: "Direction",
				},
				{
					PackageID:          "main",
					FileID:             "./test_files/ExtractTypeDefinitions/SyntaxError.cx",
					StartOffset:        100,
					Length:             15,
					LineNumber:         12,
					TypeDefinitionName: "Season",
				},
			},
			wantErr: errors.New("SyntaxError.cx:21: syntax error: type definition declaration"),
		},
		{
			scenario: "Syntax Error 2",
			testDir:  "./test_files/ExtractTypeDefinitions/SyntaxError2.cx",
			wantTypeDefinitions: []declaration_extractor.TypeDefinitionDeclaration{
				{
					PackageID:          "main",
					FileID:             "./test_files/ExtractTypeDefinitions/SyntaxError2.cx",
					StartOffset:        14,
					Length:             18,
					LineNumber:         3,
					TypeDefinitionName: "Direction",
				},
				{
					PackageID:          "main",
					FileID:             "./test_files/ExtractTypeDefinitions/SyntaxError2.cx",
					StartOffset:        100,
					Length:             15,
					LineNumber:         12,
					TypeDefinitionName: "Season",
				},
			},
			wantErr: errors.New("SyntaxError2.cx:21: syntax error: type definition declaration"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.scenario, func(t *testing.T) {
			srcBytes, err := os.ReadFile(tc.testDir)
			if err != nil {
				t.Fatal(err)
			}

			ReplaceCommentsWithWhitespaces := declaration_extractor.ReplaceCommentsWithWhitespaces(srcBytes)
			ReplaceStringContentsWithWhitespaces, err := declaration_extractor.ReplaceStringContentsWithWhitespaces(ReplaceCommentsWithWhitespaces)
			if err != nil {
				t.Fatal(filepath.Base(tc.testDir), err)
			}

			gotTypeDefinitions, gotErr := declaration_extractor.ExtractTypeDefinitions(ReplaceStringContentsWithWhitespaces, tc.testDir)

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
					Length:      17,
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
					Length:      19,
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
					Length:      17,
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
					Length:      18,
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
		{
			scenario:    "Package Error",
			testDir:     "./test_files/ExtractStructs/PackageError.cx",
			wantStructs: []declaration_extractor.StructDeclaration{},
			wantErr:     errors.New("PackageError.cx:2: syntax error: package declaration"),
		},
		{
			scenario: "Syntax Error",
			testDir:  "./test_files/ExtractStructs/SyntaxError.cx",
			wantStructs: []declaration_extractor.StructDeclaration{
				{
					PackageID:   "main",
					FileID:      "./test_files/ExtractStructs/SyntaxError.cx",
					StartOffset: 14,
					Length:      17,
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
			},
			wantErr: errors.New("SyntaxError.cx:8: syntax error: struct declaration"),
		},
		{
			scenario: "Syntax Error 2",
			testDir:  "./test_files/ExtractStructs/SyntaxError2.cx",
			wantStructs: []declaration_extractor.StructDeclaration{
				{
					PackageID:   "main",
					FileID:      "./test_files/ExtractStructs/SyntaxError2.cx",
					StartOffset: 58,
					Length:      17,
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
			},
			wantErr: errors.New("SyntaxError2.cx:16: syntax error:struct field"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.scenario, func(t *testing.T) {
			srcBytes, err := os.ReadFile(tc.testDir)
			if err != nil {
				t.Fatal(err)
			}
			ReplaceCommentsWithWhitespaces := declaration_extractor.ReplaceCommentsWithWhitespaces(srcBytes)
			ReplaceStringContentsWithWhitespaces, err := declaration_extractor.ReplaceStringContentsWithWhitespaces(ReplaceCommentsWithWhitespaces)
			if err != nil {
				t.Fatal(filepath.Base(tc.testDir), err)
			}

			gotStructs, gotErr := declaration_extractor.ExtractStructs(ReplaceStringContentsWithWhitespaces, tc.testDir)

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
					StartOffset: 322,
					Length:      12,
					LineNumber:  20,
					FuncName:    "main",
				},
				{
					PackageID:   "main",
					FileID:      "./test_files/ExtractFuncs/HasFuncs.cx",
					StartOffset: 14,
					Length:      53,
					LineNumber:  3,
					FuncName:    "addition",
				},
				{
					PackageID:   "main",
					FileID:      "./test_files/ExtractFuncs/HasFuncs.cx",
					StartOffset: 104,
					Length:      50,
					LineNumber:  7,
					FuncName:    "minus",
				},
				{
					PackageID:   "main",
					FileID:      "./test_files/ExtractFuncs/HasFuncs.cx",
					StartOffset: 226,
					Length:      29,
					LineNumber:  15,
					FuncName:    "printName",
				},
			},
		},
		{
			scenario:  "Package Error",
			testDir:   "./test_files/ExtractFuncs/PackageError.cx",
			wantFuncs: []declaration_extractor.FuncDeclaration{},
			wantErr:   errors.New("PackageError.cx:1: syntax error: package declaration"),
		},
		{
			scenario: "Syntax Error",
			testDir:  "./test_files/ExtractFuncs/SyntaxError.cx",
			wantFuncs: []declaration_extractor.FuncDeclaration{
				{
					PackageID:   "main",
					FileID:      "./test_files/ExtractFuncs/SyntaxError.cx",
					StartOffset: 322,
					Length:      12,
					LineNumber:  20,
					FuncName:    "main",
				},
				{
					PackageID:   "main",
					FileID:      "./test_files/ExtractFuncs/SyntaxError.cx",
					StartOffset: 14,
					Length:      53,
					LineNumber:  3,
					FuncName:    "addition",
				},
				{
					PackageID:   "main",
					FileID:      "./test_files/ExtractFuncs/SyntaxError.cx",
					StartOffset: 104,
					Length:      50,
					LineNumber:  7,
					FuncName:    "minus",
				},
				{
					PackageID:   "main",
					FileID:      "./test_files/ExtractFuncs/SyntaxError.cx",
					StartOffset: 226,
					Length:      29,
					LineNumber:  15,
					FuncName:    "printName",
				},
			},
			wantErr: errors.New("SyntaxError.cx:31: syntax error: func declaration"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.scenario, func(t *testing.T) {
			srcBytes, err := os.ReadFile(tc.testDir)
			if err != nil {
				t.Fatal(err)
			}

			ReplaceCommentsWithWhitespaces := declaration_extractor.ReplaceCommentsWithWhitespaces(srcBytes)
			ReplaceStringContentsWithWhitespaces, err := declaration_extractor.ReplaceStringContentsWithWhitespaces(ReplaceCommentsWithWhitespaces)
			if err != nil {
				t.Fatal(filepath.Base(tc.testDir), err)
			}

			gotFuncs, gotErr := declaration_extractor.ExtractFuncs(ReplaceStringContentsWithWhitespaces, tc.testDir)

			for _, wantFunc := range tc.wantFuncs {

				wantFunc.StartOffset = setOffset(wantFunc.StartOffset, wantFunc.LineNumber)

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
			scenario:               "Redeclared import",
			testDir:                "./test_files/ReDeclarationCheck/RedeclaredImport.cx",
			wantReDeclarationError: errors.New("RedeclaredImport.cx:4: redeclaration error: import: testimport"),
		},
		{
			scenario:               "Redeclared global",
			testDir:                "./test_files/ReDeclarationCheck/RedeclaredGlobal.cx",
			wantReDeclarationError: errors.New("RedeclaredGlobal.cx:5: redeclaration error: global: banana"),
		},
		{
			scenario:               "Redeclared enum",
			testDir:                "./test_files/ReDeclarationCheck/RedeclaredEnum.cx",
			wantReDeclarationError: errors.New("RedeclaredEnum.cx:8: redeclaration error: enum: myEnum"),
		},
		{
			scenario:               "Redeclared type definition",
			testDir:                "./test_files/ReDeclarationCheck/RedeclaredTypeDefinition.cx",
			wantReDeclarationError: errors.New("RedeclaredTypeDefinition.cx:21: redeclaration error: type definition: Season"),
		},
		{
			scenario:               "Redeclared struct",
			testDir:                "./test_files/ReDeclarationCheck/RedeclaredStruct.cx",
			wantReDeclarationError: errors.New("RedeclaredStruct.cx:15: redeclaration error: struct: Animal"),
		},
		{
			scenario:               "Redeclared struct field",
			testDir:                "./test_files/ReDeclarationCheck/RedeclaredStructField.cx",
			wantReDeclarationError: errors.New("RedeclaredStructField.cx:13: redeclaration error: struct field: name"),
		},
		{
			scenario:               "Redeclared func",
			testDir:                "./test_files/ReDeclarationCheck/RedeclaredFunc.cx",
			wantReDeclarationError: errors.New("RedeclaredFunc.cx:23: redeclaration error: func: add"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.scenario, func(t *testing.T) {
			srcBytes, err := os.ReadFile(tc.testDir)
			if err != nil {
				t.Fatal(err)
			}

			ReplaceCommentsWithWhitespaces := declaration_extractor.ReplaceCommentsWithWhitespaces(srcBytes)
			ReplaceStringContentsWithWhitespaces, err := declaration_extractor.ReplaceStringContentsWithWhitespaces(ReplaceCommentsWithWhitespaces)
			if err != nil {
				t.Fatal(filepath.Base(tc.testDir), err)
			}

			imports, err := declaration_extractor.ExtractImports(ReplaceCommentsWithWhitespaces, tc.testDir)
			if err != nil {
				t.Fatal(err)
			}

			globals, err := declaration_extractor.ExtractGlobals(ReplaceStringContentsWithWhitespaces, tc.testDir)
			if err != nil {
				t.Fatal(err)
			}

			enums, err := declaration_extractor.ExtractEnums(ReplaceStringContentsWithWhitespaces, tc.testDir)
			if err != nil {
				t.Fatal(err)
			}

			typeDefinitions, err := declaration_extractor.ExtractTypeDefinitions(ReplaceStringContentsWithWhitespaces, tc.testDir)
			if err != nil {
				t.Fatal(err)
			}

			structs, err := declaration_extractor.ExtractStructs(ReplaceStringContentsWithWhitespaces, tc.testDir)
			if err != nil {
				t.Fatal(err)
			}

			funcs, err := declaration_extractor.ExtractFuncs(ReplaceStringContentsWithWhitespaces, tc.testDir)
			if err != nil {
				t.Fatal(err)
			}

			gotReDeclarationError := declaration_extractor.ReDeclarationCheck(imports, globals, enums, typeDefinitions, structs, funcs)

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
			testDir:  "./test_files/GetDeclarations/HasDeclarations.cx",
			wantDeclarations: []string{
				"var number i32",
				`var string str = "Hello World"`,
				"Apple fruit = iota",
				"Orange",
				"Banana",
				"type fruit i64",
				"type Blender struct",
				"func (b *Blender) blend ()",
				"func main ()",
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
			ReplaceStringContentsWithWhitespaces, err := declaration_extractor.ReplaceStringContentsWithWhitespaces(ReplaceCommentsWithWhitespaces)
			if err != nil {
				t.Fatal(filepath.Base(tc.testDir), err)
			}

			imports, err := declaration_extractor.ExtractImports(ReplaceCommentsWithWhitespaces, tc.testDir)
			if err != nil {
				t.Fatal(err)
			}

			globals, err := declaration_extractor.ExtractGlobals(ReplaceStringContentsWithWhitespaces, tc.testDir)
			if err != nil {
				t.Fatal(err)
			}

			enums, err := declaration_extractor.ExtractEnums(ReplaceStringContentsWithWhitespaces, tc.testDir)
			if err != nil {
				t.Fatal(err)
			}

			typeDefinitions, err := declaration_extractor.ExtractTypeDefinitions(ReplaceStringContentsWithWhitespaces, tc.testDir)
			if err != nil {
				t.Fatal(err)
			}

			structs, err := declaration_extractor.ExtractStructs(ReplaceStringContentsWithWhitespaces, tc.testDir)
			if err != nil {
				t.Fatal(err)
			}

			funcs, err := declaration_extractor.ExtractFuncs(ReplaceStringContentsWithWhitespaces, tc.testDir)
			if err != nil {
				t.Fatal(err)
			}

			if declaration_extractor.ReDeclarationCheck(imports, globals, enums, typeDefinitions, structs, funcs) != nil {
				t.Fatal(err)
			}

			declarations := declaration_extractor.GetDeclarations(ReplaceCommentsWithWhitespaces, globals, enums, typeDefinitions, structs, funcs)

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
		scenario            string
		testDir             string
		wantImports         int
		wantGlobals         int
		wantEnums           int
		wantTypeDefinitions int
		wantStructs         int
		wantFuncs           int
		wantError           error
	}{
		{
			scenario:            "Single file",
			testDir:             "./test_files/ExtractAllDeclarations/SingleFile",
			wantGlobals:         2,
			wantEnums:           3,
			wantTypeDefinitions: 1,
			wantStructs:         1,
			wantFuncs:           2,
			wantError:           nil,
		},
		{
			scenario:            "Multiple files",
			testDir:             "./test_files/ExtractAllDeclarations/MultipleFiles",
			wantImports:         2,
			wantGlobals:         3,
			wantEnums:           12,
			wantTypeDefinitions: 3,
			wantStructs:         3,
			wantFuncs:           5,
			wantError:           nil,
		},
		{
			scenario:            "Global Syntax Error",
			testDir:             "./test_files/ExtractAllDeclarations/GlobalSyntaxError",
			wantImports:         2,
			wantGlobals:         0,
			wantEnums:           12,
			wantTypeDefinitions: 3,
			wantStructs:         3,
			wantFuncs:           5,
			wantError:           errors.New("main.cx:9: syntax error: global declaration"),
		},
		{
			scenario:            "Enum Syntax Error",
			testDir:             "./test_files/ExtractAllDeclarations/EnumSyntaxError",
			wantImports:         2,
			wantGlobals:         3,
			wantEnums:           0,
			wantTypeDefinitions: 3,
			wantStructs:         3,
			wantFuncs:           5,
			wantError:           errors.New("program.cx:27: syntax error: enum declaration"),
		},
		{
			scenario:            "Type Definition Syntax Error",
			testDir:             "./test_files/ExtractAllDeclarations/TypeDefinitionSyntaxError",
			wantImports:         2,
			wantGlobals:         3,
			wantEnums:           12,
			wantTypeDefinitions: 0,
			wantStructs:         3,
			wantFuncs:           5,
			wantError:           errors.New("program.cx:21: syntax error: type definition declaration"),
		},
		{
			scenario:            "Struct Syntax Error",
			testDir:             "./test_files/ExtractAllDeclarations/StructSyntaxError",
			wantImports:         2,
			wantGlobals:         3,
			wantEnums:           12,
			wantTypeDefinitions: 3,
			wantStructs:         0,
			wantFuncs:           5,
			wantError:           errors.New("helper.cx:20: syntax error: struct declaration"),
		},
		{
			scenario:            "Func Syntax Error",
			testDir:             "./test_files/ExtractAllDeclarations/FuncSyntaxError",
			wantImports:         2,
			wantGlobals:         3,
			wantEnums:           12,
			wantTypeDefinitions: 3,
			wantStructs:         3,
			wantFuncs:           2,
			wantError:           errors.New("helper.cx:34: syntax error: func declaration"),
		},
		{
		  scenario:            "Redeclaration Error",
			testDir:             "./test_files/ExtractAllDeclarations/RedeclarationError",
			wantImports:         2,
			wantGlobals:         4,
			wantEnums:           12,
			wantTypeDefinitions: 3,
			wantStructs:         3,
			wantFuncs:           5,
			wantError:           errors.New("main.cx:13: redeclaration error: global: number"),
		},
		{
			scenario:            "String Syntax Error",
			testDir:             "./test_files/ExtractAllDeclarations/StringSyntaxError",
			wantImports:         2,
			wantGlobals:         3,
			wantEnums:           0,
			wantTypeDefinitions: 0,
			wantStructs:         3,
			wantFuncs:           4,
			wantError:           errors.New("program.cx:41: syntax error: quote not terminated"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.scenario, func(t *testing.T) {

			actions.AST = nil

			_, sourceCodes, _ := loader.ParseArgsForCX([]string{tc.testDir}, true)
			err := loader.LoadCXProgram("mypkg1", sourceCodes, "bolt")
			if err != nil {
				t.Fatal(err)
			}

			err = file_output.AddPkgsToAST("mypkg1", "bolt")
			if err != nil {
				t.Fatal(err)
			}

			files, err := file_output.GetImportFiles("mypkg1", "bolt")
			if err != nil {
				t.Fatal(err)
			}

			Imports, Globals, Enums, TypeDefinitions, Structs, Funcs, gotErr := declaration_extractor.ExtractAllDeclarations(files)

			if len(Globals) == 0 && len(Enums) == 0 && len(Structs) == 0 && len(Funcs) == 0 {
				t.Error("No Declarations found")
			}

			if len(Imports) != tc.wantImports {
				t.Errorf("want import %v, got %v", tc.wantImports, len(Imports))
			}

			if len(Globals) != tc.wantGlobals {
				t.Errorf("want global %v, got %v", tc.wantGlobals, len(Globals))
			}

			if len(Enums) != tc.wantEnums {
				t.Errorf("want enum %v, got %v", tc.wantEnums, len(Enums))
			}

			if len(TypeDefinitions) != tc.wantTypeDefinitions {
				t.Errorf("want type definition %v, got %v", tc.wantTypeDefinitions, len(TypeDefinitions))
			}

			if len(Structs) != tc.wantStructs {
				t.Errorf("want struct %v, got %v", tc.wantStructs, len(Structs))
			}

			if len(Funcs) != tc.wantFuncs {
				t.Errorf("want func %v, got %v", tc.wantFuncs, len(Funcs))
			}

			if (gotErr != nil && tc.wantError == nil) ||
				(gotErr == nil && tc.wantError != nil) {
				t.Errorf("want error %v, got %v", tc.wantError, gotErr)
			}

			if gotErr != nil && tc.wantError != nil {
				if gotErr.Error() != tc.wantError.Error() {
					t.Errorf("want error %v, got %v", tc.wantError, gotErr)
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
