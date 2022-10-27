package type_checker_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/skycoin/cx/cmd/declaration_extractor"
	"github.com/skycoin/cx/cmd/packageloader/file_output"
	"github.com/skycoin/cx/cmd/packageloader/loader"
	"github.com/skycoin/cx/cmd/type_checker"
	"github.com/skycoin/cx/cx/ast"
	cxinit "github.com/skycoin/cx/cx/init"
	cxpackages "github.com/skycoin/cx/cx/packages"
	"github.com/skycoin/cx/cx/types"
	"github.com/skycoin/cx/cxparser/actions"
)

func TestTypeChecker_ParseDeclarationSpecifier(t *testing.T) {
	tests := []struct {
		scenario                                string
		testString                              string
		fileName                                string
		lineno                                  int
		strctName                               string
		pkgName                                 string
		wantDeclarationSpecifierFormattedString string
		wantErr                                 error
	}{
		{
			scenario:                                "Has Type Specifier",
			testString:                              "str",
			fileName:                                "./myFile",
			lineno:                                  4,
			wantDeclarationSpecifierFormattedString: "str",
			wantErr:                                 nil,
		},
		{
			scenario:                                "Has Indentifier",
			testString:                              "Animal",
			fileName:                                "./testFile",
			lineno:                                  10,
			strctName:                               "Animal",
			pkgName:                                 "",
			wantDeclarationSpecifierFormattedString: "Animal",
			wantErr:                                 nil,
		},
		{
			scenario:                                "Has External Indentifier",
			testString:                              "tester.Direction",
			fileName:                                "./myFile",
			lineno:                                  15,
			strctName:                               "Direction",
			pkgName:                                 "tester",
			wantDeclarationSpecifierFormattedString: "Direction",
			wantErr:                                 nil,
		},
		{
			scenario:                                "Has External Type Specifier",
			testString:                              "i32.counter",
			fileName:                                "./myFile",
			lineno:                                  23,
			strctName:                               "counter",
			pkgName:                                 "i32",
			wantDeclarationSpecifierFormattedString: "counter",
			wantErr:                                 nil,
		},
		{
			scenario:                                "Has Array ",
			testString:                              "[5]i64",
			fileName:                                "./testFile",
			lineno:                                  67,
			wantDeclarationSpecifierFormattedString: "[5]i64",
			wantErr:                                 nil,
		},
		{
			scenario:                                "Has Slice",
			testString:                              "[]str",
			fileName:                                "./myFile",
			lineno:                                  45,
			wantDeclarationSpecifierFormattedString: "[]str",
			wantErr:                                 nil,
		},
		{
			scenario:                                "Has Pointer",
			testString:                              "*ui8",
			fileName:                                "./testFile",
			lineno:                                  23,
			wantDeclarationSpecifierFormattedString: "*ui8",
			wantErr:                                 nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.scenario, func(t *testing.T) {
			actions.AST = cxinit.MakeProgram()

			if tc.strctName != "" {

				pkg := ast.MakePackage(tc.pkgName)
				pkgIdx := actions.AST.AddPackage(pkg)
				strct := ast.MakeStruct(tc.strctName)
				strct.Package = ast.CXPackageIndex(pkgIdx)
				pkg = pkg.AddStruct(actions.AST, strct)

			}

			if tc.pkgName != "" {
				actions.DeclareImport(actions.AST, tc.pkgName, tc.fileName, tc.lineno)
			}

			var gotDeclarationSpecifier *ast.CXArgument
			gotDeclarationSpecifier, gotErr := type_checker.ParseDeclarationSpecifier([]byte(tc.testString), tc.fileName, tc.lineno, gotDeclarationSpecifier)
			gotTypeSignature := ast.GetCXTypeSignatureRepresentationOfCXArg(actions.AST, gotDeclarationSpecifier)
			gotDeclarationSpecifierFormattedString := ast.GetFormattedType(actions.AST, gotTypeSignature)

			if gotDeclarationSpecifierFormattedString != tc.wantDeclarationSpecifierFormattedString {
				t.Errorf("want %s, got %s", tc.wantDeclarationSpecifierFormattedString, gotDeclarationSpecifierFormattedString)
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

func TestTypeChecker_ParseParameterDeclaration(t *testing.T) {
	tests := []struct {
		scenario                                string
		testString                              string
		fileName                                string
		lineno                                  int
		strctName                               string
		pkgName                                 string
		wantParameterDeclarationFormattedString string
		wantErr                                 error
	}{
		{
			scenario:                                "Has Type Specifier",
			testString:                              "name str",
			fileName:                                "./myFile",
			lineno:                                  4,
			wantParameterDeclarationFormattedString: "name str",
			wantErr:                                 nil,
		},
		{
			scenario:                                "Has Indentifier",
			testString:                              "cat Animal",
			fileName:                                "./testFile",
			lineno:                                  10,
			strctName:                               "Animal",
			pkgName:                                 "",
			wantParameterDeclarationFormattedString: "cat Animal",
			wantErr:                                 nil,
		},
		{
			scenario:                                "Has External Indentifier",
			testString:                              "North tester.Direction",
			fileName:                                "./myFile",
			lineno:                                  15,
			strctName:                               "Direction",
			pkgName:                                 "tester",
			wantParameterDeclarationFormattedString: "North Direction",
			wantErr:                                 nil,
		},
		{
			scenario:                                "Has External Type Specifier",
			testString:                              "clock i32.counter",
			fileName:                                "./myFile",
			lineno:                                  23,
			strctName:                               "counter",
			pkgName:                                 "i32",
			wantParameterDeclarationFormattedString: "clock counter",
			wantErr:                                 nil,
		},
		{
			scenario:                                "Has Array ",
			testString:                              "lottery [5]i64",
			fileName:                                "./testFile",
			lineno:                                  67,
			wantParameterDeclarationFormattedString: "lottery [5]i64",
			wantErr:                                 nil,
		},
		{
			scenario:                                "Has Slice",
			testString:                              "months []str",
			fileName:                                "./myFile",
			lineno:                                  45,
			wantParameterDeclarationFormattedString: "months []str",
			wantErr:                                 nil,
		},
		{
			scenario:                                "Has Pointer",
			testString:                              "link *ui8",
			fileName:                                "./testFile",
			lineno:                                  23,
			wantParameterDeclarationFormattedString: "link *ui8",
			wantErr:                                 nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.scenario, func(t *testing.T) {
			actions.AST = cxinit.MakeProgram()

			var pkg *ast.CXPackage
			pkg = ast.MakePackage(tc.pkgName)
			pkgIdx := actions.AST.AddPackage(pkg)

			if tc.strctName != "" {
				strct := ast.MakeStruct(tc.strctName)
				strct.Package = ast.CXPackageIndex(pkgIdx)
				pkg = pkg.AddStruct(actions.AST, strct)
			}

			if tc.pkgName != "" {
				actions.DeclareImport(actions.AST, tc.pkgName, tc.fileName, tc.lineno)
			}

			var gotParameterDeclaration *ast.CXArgument
			gotParameterDeclaration, gotErr := type_checker.ParseParameterDeclaration([]byte(tc.testString), pkg, tc.fileName, tc.lineno)
			gotParameterDeclarationFormattedName := ast.GetFormattedName(actions.AST, gotParameterDeclaration, false, pkg)
			gotTypeSignature := ast.GetCXTypeSignatureRepresentationOfCXArg(actions.AST, gotParameterDeclaration)
			gotParameterDeclarationFormattedType := ast.GetFormattedType(actions.AST, gotTypeSignature)
			gotParameterDeclarationFormattedString := gotParameterDeclarationFormattedName + " " + gotParameterDeclarationFormattedType

			if gotParameterDeclarationFormattedString != tc.wantParameterDeclarationFormattedString {
				t.Errorf("want %s, got %s", tc.wantParameterDeclarationFormattedString, gotParameterDeclarationFormattedString)
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

func TestTypeChecker_ParseImports(t *testing.T) {

	tests := []struct {
		scenario string
		testDir  string
		imports  []string
	}{
		{
			scenario: "Has Imports",
			testDir:  "./test_files/ParseImports/HasImports.cx",
			imports: []string{
				"myImport",
				"anotherImport",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.scenario, func(t *testing.T) {

			actions.AST = cxinit.MakeProgram()

			srcBytes, err := os.ReadFile(tc.testDir)
			if err != nil {
				t.Error(err)
			}

			ReplaceCommentsWithWhitespaces := declaration_extractor.ReplaceCommentsWithWhitespaces(srcBytes)

			Imports, err := declaration_extractor.ExtractImports(ReplaceCommentsWithWhitespaces, tc.testDir)
			if err != nil {
				t.Fatal(err)
			}

			err = type_checker.ParseImports(Imports)

			program := actions.AST

			for _, wantImport := range tc.imports {

				var match bool
				var gotImport string
				for _, pkgIdx := range program.Packages {
					pkg, err := program.GetPackageFromArray(pkgIdx)
					if err != nil {
						panic(err)
					}

					if cxpackages.IsDefaultPackage(pkg.Name) {
						continue
					}

					for _, impIdx := range pkg.Imports {
						impPkg, err := program.GetPackageFromArray(impIdx)
						if err != nil {
							panic(err)
						}

						gotImport = impPkg.Name

						if gotImport == wantImport {
							match = true
							break
						}
					}
				}

				if !match {
					t.Errorf("want import %s, got %s", wantImport, gotImport)
				}
			}

		})
	}
}

func TestTypeChecker_ParseGlobals(t *testing.T) {

	type GlobalTypeSignature struct {
		Package string
		Index   int
		Name    string
		Type    string
	}

	tests := []struct {
		scenario             string
		testDir              string
		globalTypeSignatures []GlobalTypeSignature
	}{
		{
			scenario: "Has globals",
			testDir:  "./test_files/ParseGlobals/HasGlobals",
			globalTypeSignatures: []GlobalTypeSignature{
				{
					Name:    "Bool",
					Index:   0,
					Package: "main",
					Type:    "bool",
				},
				{
					Name:    "Byte",
					Index:   1,
					Package: "main",
					Type:    "i8",
				},
				{
					Name:    "I16",
					Index:   2,
					Package: "main",
					Type:    "i16",
				},
				{
					Name:    "I32",
					Index:   3,
					Package: "main",
					Type:    "i32",
				},
				{
					Name:    "I64",
					Index:   4,
					Package: "main",
					Type:    "i64",
				},
				{
					Name:    "UByte",
					Index:   5,
					Package: "main",
					Type:    "ui8",
				},
				{
					Name:    "UI16",
					Index:   6,
					Package: "main",
					Type:    "ui16",
				},
				{
					Name:    "UI32",
					Index:   7,
					Package: "main",
					Type:    "ui32",
				},
				{
					Name:    "UI64",
					Index:   8,
					Package: "main",
					Type:    "ui64",
				},
				{
					Name:    "F32",
					Index:   9,
					Package: "main",
					Type:    "f32",
				},
				{
					Name:    "F64",
					Index:   10,
					Package: "main",
					Type:    "f64",
				},
				{
					Name:    "string",
					Index:   11,
					Package: "main",
					Type:    "str",
				},
				{
					Name:    "Affordance",
					Index:   12,
					Package: "main",
					Type:    "[]aff",
				},
				{
					Name:    "intArray",
					Index:   13,
					Package: "main",
					Type:    "[5]i32",
				},
			},
		},
		{
			scenario: "Has globals 2",
			testDir:  "./test_files/ParseGlobals/HasGlobals2",
			globalTypeSignatures: []GlobalTypeSignature{
				{
					Name:    "number",
					Index:   0,
					Package: "main",
					Type:    "i32",
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.scenario, func(t *testing.T) {

			actions.AST = cxinit.MakeProgram()

			_, sourceCode, _ := loader.ParseArgsForCX([]string{tc.testDir}, true)

			err := loader.LoadCXProgram("test", sourceCode, "bolt")
			if err != nil {
				t.Fatal(err)
			}

			files, err := file_output.GetImportFiles("test", "bolt")
			if err != nil {
				t.Fatal(err)
			}

			ReplaceCommentsWithWhitespaces := declaration_extractor.ReplaceCommentsWithWhitespaces(files[0].Content)
			ReplaceStringContentsWithWhitespaces, err := declaration_extractor.ReplaceStringContentsWithWhitespaces(ReplaceCommentsWithWhitespaces)
			if err != nil {
				t.Fatal(err)
			}

			Globals, err := declaration_extractor.ExtractGlobals(ReplaceStringContentsWithWhitespaces, files[0].FileName)
			if err != nil {
				t.Fatal(err)
			}

			err = type_checker.ParseGlobals(files, Globals)
			if err != nil {
				t.Fatal(err)
			}

			program := actions.AST

			for _, wantGlobal := range tc.globalTypeSignatures {

				for _, pkgIdx := range program.Packages {

					pkg, err := program.GetPackageFromArray(pkgIdx)
					gotPkgName := pkg.Name

					if err != nil {
						t.Log(err)
					}

					if cxpackages.IsDefaultPackage(pkg.Name) {
						continue
					}

					var match bool
					var gotGlobal *ast.CXTypeSignature
					var gotGlobalType string

					for _, globalIdx := range pkg.Globals.Fields {
						gotGlobal = program.GetCXTypeSignatureFromArray(globalIdx)

						if gotGlobal.Name == wantGlobal.Name {

							gotGlobalType = ast.GetFormattedType(actions.AST, gotGlobal)

							if int(gotGlobal.Index) == wantGlobal.Index &&
								gotPkgName == wantGlobal.Package &&
								gotGlobalType == wantGlobal.Type {
								match = true
							}

							break
						}

					}

					if !match {

						if gotGlobal.Type == ast.TYPE_CXARGUMENT_DEPRECATE {
							t.Errorf("want global %v %v %v %v, got %v %v %v %v", wantGlobal.Name, wantGlobal.Index, wantGlobal.Package, wantGlobal.Type, program.GetCXArg(ast.CXArgumentIndex(gotGlobal.Meta)).Name, gotGlobal.Index, gotPkgName, gotGlobalType)
						} else {
							t.Errorf("want global %v %v %v %v, got %v %v %v %v", wantGlobal.Name, wantGlobal.Index, wantGlobal.Package, wantGlobal.Type, gotGlobal.Name, gotGlobal.Index, gotPkgName, gotGlobalType)
						}
					}
				}

			}

		})
	}

}

// func TestTypeChecks_ParseEnums(t *testing.T) {

// 	tests := []struct {
// 		scenario string
// 		testDir  string
// 	}{}

// 	for _, tc := range tests {
// 		t.Run(tc.scenario, func(t *testing.T) {

// 		})
// 	}

// }

func TestTypeChecker_ParseStructs(t *testing.T) {

	type StructFieldTypeSignature struct {
		Index int
		Name  string
		Type  string
	}

	type StructTypeSignature struct {
		Package string
		Index   int
		Name    string
		Type    string
		Fields  []StructFieldTypeSignature
	}

	tests := []struct {
		scenario             string
		testDir              string
		structTypeSignatures []StructTypeSignature
	}{
		{
			scenario: "Has Structs",
			testDir:  "./test_files/ParseStructs/HasStructs",
			structTypeSignatures: []StructTypeSignature{
				{
					Name:    "CustomType",
					Index:   1,
					Package: "main",
					Fields: []StructFieldTypeSignature{
						{
							Index: 0,
							Name:  "fieldA",
							Type:  "str",
						},
						{
							Index: 1,
							Name:  "fieldB",
							Type:  "i32",
						},
					},
				},
				{
					Name:    "AnotherType",
					Index:   2,
					Package: "main",
					Fields: []StructFieldTypeSignature{
						{
							Index: 0,
							Name:  "name",
							Type:  "str",
						},
					},
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.scenario, func(t *testing.T) {

			actions.AST = cxinit.MakeProgram()

			_, sourceCode, _ := loader.ParseArgsForCX([]string{tc.testDir}, true)

			err := loader.LoadCXProgram("test", sourceCode, "bolt")
			if err != nil {
				t.Fatal(err)
			}

			files, err := file_output.GetImportFiles("test", "bolt")
			if err != nil {
				t.Fatal(err)
			}

			ReplaceCommentsWithWhitespaces := declaration_extractor.ReplaceCommentsWithWhitespaces(files[0].Content)

			structs, err := declaration_extractor.ExtractStructs(ReplaceCommentsWithWhitespaces, files[0].FileName)

			type_checker.ParseStructs(files, structs)

			program := actions.AST

			for _, wantStruct := range tc.structTypeSignatures {

				var match bool
				var ast3 string
				var gotStruct ast.CXStruct

				for _, pkgIdx := range program.Packages {

					pkg, err := program.GetPackageFromArray(pkgIdx)
					gotPkgName := pkg.Name

					if err != nil {
						t.Log(err)
					}

					if cxpackages.IsDefaultPackage(pkg.Name) {
						continue
					}

					for _, structIdx := range pkg.Structs {
						gotStruct = program.CXStructs[structIdx]

						if gotStruct.Name == wantStruct.Name &&
							gotPkgName == wantStruct.Package {
							ast3 += fmt.Sprintf("want %s %d. %s, got %s %d. %s\n", wantStruct.Package, wantStruct.Index, wantStruct.Name, gotPkgName, gotStruct.Index, gotStruct.Name)

							if gotStruct.Index != wantStruct.Index || len(gotStruct.Fields) != len(wantStruct.Fields) {
								ast3 += fmt.Sprintf("want %d fields, got %d fields", len(wantStruct.Fields), len(gotStruct.Fields))
								break
							}

							for _, wantField := range wantStruct.Fields {

								for _, gotFieldIdx := range gotStruct.Fields {

									gotField := program.GetCXTypeSignatureFromArray(gotFieldIdx)
									var gotFieldType string
									var gotFieldIndex int

									if gotField.Name == wantField.Name {

										gotFieldType = ast.GetFormattedType(actions.AST, gotField)

										ast3 += fmt.Sprintf("want field %d. %s %s, got %d. %s %s\n", wantField.Index, wantField.Name, wantField.Type, gotFieldIndex, gotField.Name, gotFieldType)

										if gotFieldIndex == wantField.Index &&
											gotFieldType == wantField.Type {
											match = true
											break
										}
									}

								}

							}

							break
						}

					}

				}

				if !match {
					t.Error(ast3)
				}

			}

		})
	}

}

func TestTypeChecker_ParseFuncHeaders(t *testing.T) {

	tests := []struct {
		scenario    string
		testDir     string
		functionCXs []ast.CXFunction
	}{
		{
			scenario: "Has funcs",
			testDir:  "./test_files/ParseFuncs/HasFuncs",
			functionCXs: []ast.CXFunction{
				{
					Name:    "main",
					Index:   0,
					Package: 1,
				},
				{
					Name:    "add",
					Index:   1,
					Package: 1,
				},
				{
					Name:    "divide",
					Index:   2,
					Package: 1,
				},
				{
					Name:    "printer",
					Index:   3,
					Package: 1,
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.scenario, func(t *testing.T) {

			actions.AST = cxinit.MakeProgram()

			_, sourceCode, _ := loader.ParseArgsForCX([]string{tc.testDir}, true)

			err := loader.LoadCXProgram("test", sourceCode, "bolt")
			if err != nil {
				t.Fatal(err)
			}

			files, err := file_output.GetImportFiles("test", "bolt")
			if err != nil {
				t.Fatal(err)
			}

			ReplaceCommentsWithWhitespaces := declaration_extractor.ReplaceCommentsWithWhitespaces(files[0].Content)
			ReplaceStringContentsWithWhitespaces, err := declaration_extractor.ReplaceStringContentsWithWhitespaces(ReplaceCommentsWithWhitespaces)
			if err != nil {
				t.Fatal(err)
			}

			funcs, err := declaration_extractor.ExtractFuncs(ReplaceStringContentsWithWhitespaces, files[0].FileName)

			type_checker.ParseFuncHeaders(files, funcs)

			program := actions.AST

			for _, wantFunc := range tc.functionCXs {

				var match bool
				var gotFunc *ast.CXFunction

				for _, pkgIdx := range program.Packages {

					pkg, err := program.GetPackageFromArray(pkgIdx)

					if err != nil {
						t.Log(err)
					}

					for _, funcIdx := range pkg.Functions {

						gotFunc = program.GetFunctionFromArray(funcIdx)

						if gotFunc.Name == wantFunc.Name &&
							gotFunc.Index == wantFunc.Index &&
							gotFunc.Package == wantFunc.Package {

							match = true

							break

						}
					}

				}

				if !match {
					t.Errorf("want func \n%v\n\tInputs: %v\n\tOutputs: %v\n, got \n%v\n\tInputs: %v\n\tOutputs: %v\n", wantFunc, wantFunc.Inputs, wantFunc.Outputs, gotFunc, gotFunc.Inputs, gotFunc.Outputs)
				}
			}

		})
	}

}

func TestTypeChecker_ParseAllDeclarations(t *testing.T) {

	type Global struct {
		Name string
		Type string
	}

	type StructField struct {
		Name string
		Type string
	}

	type Struct struct {
		Name   string
		Fields []StructField
	}

	type Func struct {
		Name         string
		RecieverName string
		RecieverType string
	}

	type Package struct {
		Name    string
		Imports []string
		Globals []Global
		Structs []Struct
		Funcs   []Func
	}

	tests := []struct {
		scenario    string
		testDir     string
		wantProgram []Package
	}{
		{
			scenario: "Has Declarations",
			testDir:  "./test_files/ParseAllDeclarations/HasDeclarations",
			wantProgram: []Package{
				{
					Name: "main",
					Imports: []string{
						"helper",
						"config",
					},
					Globals: []Global{
						{
							Name: "Bool",
							Type: "bool",
						},
						{
							Name: "Byte",
							Type: "i8",
						},
						{
							Name: "I16",
							Type: "i16",
						},
						{
							Name: "I32",
							Type: "i32",
						},
						{
							Name: "I64",
							Type: "i64",
						},
						{
							Name: "UByte",
							Type: "ui8",
						},
						{
							Name: "UI16",
							Type: "ui16",
						},
						{
							Name: "UI32",
							Type: "ui32",
						},
						{
							Name: "UI64",
							Type: "ui64",
						},
						{
							Name: "F32",
							Type: "f32",
						},
						{
							Name: "F64",
							Type: "f64",
						},
						{
							Name: "string",
							Type: "str",
						},
						{
							Name: "Affordance",
							Type: "[]aff",
						},
						{
							Name: "intArray",
							Type: "[5]i32",
						},
					},
					Structs: []Struct{
						{
							Name: "CustomType",
							Fields: []StructField{
								{
									Name: "fieldA",
									Type: "str",
								},
								{
									Name: "fieldB",
									Type: "i32",
								},
							},
						},
						{
							Name: "AnotherType",
							Fields: []StructField{
								{
									Name: "name",
									Type: "str",
								},
							},
						},
					},
					Funcs: []Func{
						{
							Name: "main",
						},
						{
							Name: "add",
						},
						{
							Name: "divide",
						},
						{
							Name: "printer",
						},
						{
							Name:         "setFieldA",
							RecieverName: "customType",
							RecieverType: "CustomType",
						},
					},
				},
			},
		},
		{
			scenario: "Has Multiple Files",
			testDir:  "./test_files/ParseAllDeclarations/HasMultipleFiles",
			wantProgram: []Package{
				{
					Name: "main",
					Imports: []string{
						"math",
					},
					Globals: []Global{},
					Structs: []Struct{},
					Funcs: []Func{
						{
							Name: "main",
						},
					},
				},
				{
					Name:    "math",
					Globals: []Global{},
					Structs: []Struct{},
					Funcs: []Func{
						{
							Name: "double",
						},
					},
				},
			},
		},
		{
			scenario: "Has Imports",
			testDir:  "./test_files/ParseAllDeclarations/HasImports",
			wantProgram: []Package{
				{
					Name: "main",
					Globals: []Global{
						{
							Name: "dog",
							Type: "Animal",
						},
					},
					Structs: []Struct{},
					Funcs: []Func{
						{
							Name: "main",
						},
					},
				},
				{
					Name:    "helper",
					Globals: []Global{},
					Structs: []Struct{
						{
							Name: "Animal",
							Fields: []StructField{
								{
									Name: "sound",
									Type: "str",
								},
							},
						},
					},
					Funcs: []Func{
						{
							Name:         "Speak",
							RecieverName: "a",
							RecieverType: "Animal",
						},
					},
				},
				{
					Name: "config",
					Globals: []Global{
						{
							Name: "Name",
							Type: "str",
						},
						{
							Name: "Apple",
							Type: "i32",
						},
					},
					Structs: []Struct{},
					Funcs:   []Func{},
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.scenario, func(t *testing.T) {

			actions.AST = nil

			_, sourceCode, _ := loader.ParseArgsForCX([]string{tc.testDir}, true)

			err := loader.LoadCXProgram("test", sourceCode, "bolt")
			if err != nil {
				t.Fatal(err)
			}

			files, err := file_output.GetImportFiles("test", "bolt")
			if err != nil {
				t.Fatal(err)
			}

			Imports, Globals, Enums, TypeDefinitions, Structs, Funcs, gotErr := declaration_extractor.ExtractAllDeclarations(files)
			if (Enums != nil && TypeDefinitions != nil) || gotErr != nil {
				t.Fatal(gotErr)
			}

			err = type_checker.ParseAllDeclarations(files, Imports, Globals, Structs, Funcs)
			if err != nil {
				t.Fatal(err)
			}

			program := actions.AST

			for _, wantPkg := range tc.wantProgram {
				gotPkg, err := program.GetPackage(wantPkg.Name)
				if err != nil {
					t.Fatal(err)
				}

				if cxpackages.IsDefaultPackage(gotPkg.Name) {
					continue
				}

				// Globals
				for _, wantGlobal := range wantPkg.Globals {
					var gotGlobalName string
					var gotGlobalType string
					var match bool
					for _, globalIdx := range gotPkg.Globals.Fields {
						gotGlobal := program.GetCXTypeSignatureFromArray(globalIdx)
						gotGlobalName = gotGlobal.Name

						if gotGlobalName == wantGlobal.Name {

							gotGlobalType = ast.GetFormattedType(program, gotGlobal)

							if gotGlobalType == wantGlobal.Type {
								match = true
							}

							break
						}

					}

					if !match {
						t.Errorf("want global %s %s %s, got %s %s %s", wantPkg.Name, wantGlobal.Name, wantGlobal.Type, gotPkg.Name, gotGlobalName, gotGlobalType)
					}
				}

				// Structs
				for _, wantStruct := range wantPkg.Structs {
					var match bool
					var fields string
					var gotStruct ast.CXStruct
					for _, structIdx := range gotPkg.Structs {
						gotStruct = program.CXStructs[structIdx]

						if gotStruct.Name == wantStruct.Name {

							var fieldMatch int
							for _, wantField := range wantStruct.Fields {

								for _, gotFieldIdx := range gotStruct.Fields {

									gotField := program.GetCXTypeSignatureFromArray(gotFieldIdx)
									var gotFieldType string

									gotFieldType = ast.GetFormattedType(program, gotField)

									if gotField.Name == wantField.Name {

										if gotFieldType == wantField.Type {
											fieldMatch++
											break
										}
									}

									fields += fmt.Sprintf("want field %s %s, got %s %s\n", wantField.Name, wantField.Type, gotField.Name, gotFieldType)

								}

								if fieldMatch == len(wantStruct.Fields) {
									match = true
								}

							}

							break
						}

					}
					if !match {
						t.Errorf("want %s %s, got %s %s\n%s", wantPkg.Name, wantStruct.Name, gotPkg.Name, gotStruct.Name, fields)
					}

				}

				// Funcs
				for _, wantFunc := range wantPkg.Funcs {

					var match bool

					var wantFuncName string = wantFunc.Name
					var wantFuncReciever string

					if wantFunc.RecieverName != "" {
						wantFuncName = fmt.Sprintf("%s.%s", wantFunc.RecieverType, wantFunc.Name)
						wantFuncReciever = fmt.Sprintf("%s *%s", wantFunc.RecieverName, wantFunc.RecieverType)
					}

					var gotFuncName string
					var gotFuncReciever string
					for _, gotFuncIdx := range gotPkg.Functions {
						gotFuncReciever = ""
						gotFunc := program.GetFunctionFromArray(gotFuncIdx)
						gotFuncName = gotFunc.Name
						gotFuncInputs := gotFunc.GetInputs(program)
						if len(gotFuncInputs) != 0 {
							gotRecieverTypeSignatureIdx := gotFuncInputs[0]
							gotRecieverTypeSignature := program.GetCXTypeSignatureFromArray(gotRecieverTypeSignatureIdx)
							var gotRecieverName string
							var gotRecieverType string

							if gotRecieverTypeSignature.Type == ast.TYPE_CXARGUMENT_DEPRECATE {
								gotRecieverArg := program.CXArgs[gotRecieverTypeSignature.Meta]
								gotRecieverName = ast.GetFormattedName(program, &gotRecieverArg, false, gotPkg)
								gotRecieverType = ast.GetFormattedType(program, gotRecieverTypeSignature)
							} else if gotRecieverTypeSignature.Type == ast.TYPE_ATOMIC {
								gotRecieverType = types.Code(gotRecieverTypeSignature.Meta).Name()
							} else if gotRecieverTypeSignature.Type == ast.TYPE_POINTER_ATOMIC {
								gotRecieverType = "*" + types.Code(gotRecieverTypeSignature.Meta).Name()
							} else {
								gotRecieverType = "type is not known"
							}

							gotFuncReciever = fmt.Sprintf("%s %s", gotRecieverName, gotRecieverType)
						}

						if gotFunc.Name == wantFuncName {
							if gotFuncReciever == wantFuncReciever {
								match = true
							}
							break
						}

					}

					if !match {
						t.Errorf("want func %s (%s) (), got %s (%s) ()", wantFuncName, wantFuncReciever, gotFuncName, gotFuncReciever)

					}
				}

			}

		})
	}
}
