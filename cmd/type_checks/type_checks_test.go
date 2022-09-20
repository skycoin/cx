package type_checks_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/skycoin/cx/cmd/declaration_extractor"
	"github.com/skycoin/cx/cmd/type_checks"
	"github.com/skycoin/cx/cx/ast"
	cxinit "github.com/skycoin/cx/cx/init"
	cxpackages "github.com/skycoin/cx/cx/packages"
	"github.com/skycoin/cx/cx/types"
	"github.com/skycoin/cx/cxparser/actions"
)

func TestTypeCheck_ParseDeclarationSpecifier(t *testing.T) {
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
			gotDeclarationSpecifier, gotErr := type_checks.ParseDeclarationSpecifier([]byte(tc.testString), tc.fileName, tc.lineno, gotDeclarationSpecifier)
			gotDeclarationSpecifierFormattedString := ast.GetFormattedType(actions.AST, gotDeclarationSpecifier)

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

func TestTypeCheck_ParseParameterDeclaration(t *testing.T) {
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
			gotParameterDeclaration, gotErr := type_checks.ParseParameterDeclaration([]byte(tc.testString), pkg, tc.fileName, tc.lineno)
			gotParameterDeclarationFormattedName := ast.GetFormattedName(actions.AST, gotParameterDeclaration, false, pkg)
			gotParameterDeclarationFormattedType := ast.GetFormattedType(actions.AST, gotParameterDeclaration)
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

func TestTypeChecks_ParseGlobals(t *testing.T) {

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
			testDir:  "./test_files/ParseGlobals/test.cx",
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
			testDir:  "./test_files/ParseGlobals/testFile.cx",
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

			actions.AST = nil

			srcBytes, err := os.ReadFile(tc.testDir)
			if err != nil {
				t.Error(err)
			}

			ReplaceCommentsWithWhitespaces := declaration_extractor.ReplaceCommentsWithWhitespaces(srcBytes)
			ReplaceStringContentsWithWhitespaces, err := declaration_extractor.ReplaceStringContentsWithWhitespaces(ReplaceCommentsWithWhitespaces)
			if err != nil {
				t.Fatal(err)
			}

			Globals, err := declaration_extractor.ExtractGlobals(ReplaceStringContentsWithWhitespaces, tc.testDir)
			if err != nil {
				t.Fatal(err)
			}

			err = type_checks.ParseGlobals(Globals)
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

							if gotGlobal.Type == ast.TYPE_CXARGUMENT_DEPRECATE {
								gotGlobalType = ast.GetFormattedType(actions.AST, actions.AST.GetCXArg(ast.CXArgumentIndex(gotGlobal.Meta)))
							} else {
								gotGlobalType = types.Code(gotGlobal.Meta).Name()
							}

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

			actions.AST.PrintProgram()

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

func TestTypeChecks_ParseStructs(t *testing.T) {

	tests := []struct {
		scenario      string
		testDir       string
		structCXs     []ast.CXStruct
		typeSignature []ast.CXTypeSignature
	}{
		{
			scenario: "Has Structs",
			testDir:  "./test_files/test.cx",
			structCXs: []ast.CXStruct{
				{
					Name:    "CustomType",
					Index:   1,
					Package: 1,
					Fields:  []ast.CXTypeSignatureIndex{0, 1},
				},
				{
					Name:    "AnotherType",
					Index:   2,
					Package: 1,
					Fields:  []ast.CXTypeSignatureIndex{2},
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.scenario, func(t *testing.T) {

			actions.AST = nil

			srcBytes, err := os.ReadFile(tc.testDir)
			if err != nil {
				t.Error(err)
			}

			ReplaceCommentsWithWhitespaces := declaration_extractor.ReplaceCommentsWithWhitespaces(srcBytes)

			structs, err := declaration_extractor.ExtractStructs(ReplaceCommentsWithWhitespaces, tc.testDir)

			type_checks.ParseStructs(structs)

			program := actions.AST

			var ast3 string

			for _, wantStruct := range tc.structCXs {

				var match bool
				var gotStruct ast.CXStruct

				for _, pkgIdx := range program.Packages {

					pkg, err := program.GetPackageFromArray(pkgIdx)

					if err != nil {
						t.Log(err)
					}

					if cxpackages.IsDefaultPackage(pkg.Name) {
						continue
					}

					for _, structIdx := range pkg.Structs {
						gotStruct = program.CXStructs[structIdx]

						if gotStruct.Name == wantStruct.Name &&
							(gotStruct.Index == wantStruct.Index ||
								gotStruct.Package == wantStruct.Package) {
							ast3 += fmt.Sprintf("want %d. %s, got %d. %s\n", wantStruct.Index, wantStruct.Name, gotStruct.Index, gotStruct.Name)

							for _, wantFieldIdx := range wantStruct.Fields {

								for _, gotFieldIdx := range gotStruct.Fields {

									// gotField := program.GetCXTypeSignatureFromArray(gotFieldIdx)

									if gotFieldIdx == wantFieldIdx {
										match = true
										break
									}

									// if gotField.Type == ast.TYPE_CXARGUMENT_DEPRECATE {

									// } else if gotField.Type == ast.TYPE_ATOMIC {

									// } else if gotField.Type == ast.TYPE_POINTER_ATOMIC {

									// }
								}

							}

							break
						}

					}

				}

				if !match {
					t.Errorf("want struct %v, got %v", wantStruct, gotStruct)
				}
			}

		})
	}

}

func TestTypeChecks_ParseFuncHeaders(t *testing.T) {

	tests := []struct {
		scenario    string
		testDir     string
		functionCXs []ast.CXFunction
	}{
		{
			scenario: "Has funcs",
			testDir:  "./test_files/test.cx",
			functionCXs: []ast.CXFunction{
				{
					Name:     "main",
					Index:    0,
					Package:  1,
					FileName: "./test_files/test.cx",
					FileLine: 35,
				},
				{
					Name:    "add",
					Index:   1,
					Package: 1,
					Inputs: &ast.CXStruct{
						Fields: []ast.CXTypeSignatureIndex{},
					},
					Outputs: &ast.CXStruct{
						Fields: []ast.CXTypeSignatureIndex{},
					},
				},
				{
					Name:    "divide",
					Index:   2,
					Package: 1,
					Inputs: &ast.CXStruct{
						Fields: []ast.CXTypeSignatureIndex{},
					},
					Outputs: &ast.CXStruct{
						Fields: []ast.CXTypeSignatureIndex{},
					},
				},
				{
					Name:    "printer",
					Index:   3,
					Package: 1,
					Inputs: &ast.CXStruct{
						Fields: []ast.CXTypeSignatureIndex{},
					},
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.scenario, func(t *testing.T) {

			actions.AST = nil

			srcBytes, err := os.ReadFile(tc.testDir)
			if err != nil {
				t.Error(err)
			}

			ReplaceCommentsWithWhitespaces := declaration_extractor.ReplaceCommentsWithWhitespaces(srcBytes)
			ReplaceStringContentsWithWhitespaces, err := declaration_extractor.ReplaceStringContentsWithWhitespaces(ReplaceCommentsWithWhitespaces)
			if err != nil {
				t.Fatal(err)
			}

			funcs, err := declaration_extractor.ExtractFuncs(ReplaceStringContentsWithWhitespaces, tc.testDir)

			type_checks.ParseFuncHeaders(funcs)

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

							var paramMatch int = 2

							for k, wantInput := range wantFunc.GetInputs(program) {
								gotInput := gotFunc.GetInputs(program)[k]

								if gotInput != wantInput {
									paramMatch--
									break
								}
							}

							for k, wantOutput := range wantFunc.GetOutputs(program) {
								gotOutput := gotFunc.GetOutputs(program)[k]

								if gotOutput != wantOutput {
									paramMatch--
									break
								}
							}

							if paramMatch == 2 {
								match = true
							}

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
