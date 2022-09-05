package type_checks_test

import (
	"os"
	"testing"

	"github.com/skycoin/cx/cmd/declaration_extractor"
	"github.com/skycoin/cx/cmd/type_checks"
	"github.com/skycoin/cx/cx/ast"
	cxpackages "github.com/skycoin/cx/cx/packages"
	"github.com/skycoin/cx/cx/types"
	"github.com/skycoin/cx/cxparser/actions"
)

func TestTypeChecks_ParseGlobals(t *testing.T) {

	tests := []struct {
		scenario      string
		testDir       string
		globalsCXArgs []*ast.CXArgument
	}{
		{
			scenario: "Has globals",
			testDir:  "./test_files/test.cx",
			globalsCXArgs: []*ast.CXArgument{
				{
					Name:      "Bool",
					Index:     0,
					Package:   1,
					Type:      types.BOOL,
					Size:      1,
					TotalSize: 1,
					Offset:    1048576,
				},
				{
					Name:      "Byte",
					Index:     2,
					Package:   1,
					Type:      types.I8,
					Size:      1,
					TotalSize: 1,
					Offset:    1048577,
				},
				{
					Name:      "I16",
					Index:     4,
					Package:   1,
					Type:      types.I16,
					Size:      2,
					TotalSize: 2,
					Offset:    1048578,
				},
				{
					Name:      "I32",
					Index:     6,
					Package:   1,
					Type:      types.I32,
					Size:      4,
					TotalSize: 4,
					Offset:    1048580,
				},
				{
					Name:      "I64",
					Index:     8,
					Package:   1,
					Type:      types.I64,
					Size:      8,
					TotalSize: 8,
					Offset:    1048584,
				},
				{
					Name:      "UByte",
					Index:     10,
					Package:   1,
					Type:      types.UI8,
					Size:      1,
					TotalSize: 1,
					Offset:    1048592,
				},
				{
					Name:      "UI16",
					Index:     12,
					Package:   1,
					Type:      types.UI16,
					Size:      2,
					TotalSize: 2,
					Offset:    1048593,
				},
				{
					Name:      "UI32",
					Index:     14,
					Package:   1,
					Type:      types.UI32,
					Size:      4,
					TotalSize: 4,
					Offset:    1048595,
				},
				{
					Name:      "UI64",
					Index:     16,
					Package:   1,
					Type:      types.UI64,
					Size:      8,
					TotalSize: 8,
					Offset:    1048599,
				},
				{
					Name:      "F32",
					Index:     18,
					Package:   1,
					Type:      types.F32,
					Size:      4,
					TotalSize: 4,
					Offset:    1048607,
				},
				{
					Name:      "F64",
					Index:     20,
					Package:   1,
					Type:      types.F64,
					Size:      8,
					TotalSize: 8,
					Offset:    1048611,
				},
				{
					Name:      "string",
					Index:     22,
					Package:   1,
					Type:      types.STR,
					Size:      8,
					TotalSize: 8,
					Offset:    1048619,
				},
				{
					Name:      "Affordance",
					Index:     24,
					Package:   1,
					Type:      types.AFF,
					TotalSize: 8,
					Offset:    1048627,
				},
				{
					Name:      "intArray",
					Index:     26,
					Package:   1,
					Type:      types.I32,
					Size:      4,
					TotalSize: 20,
					Offset:    1048635,
				},
			},
		},
		{
			scenario: "Has globals 2",
			testDir:  "./test_files/testFile.cx",
			globalsCXArgs: []*ast.CXArgument{
				{
					Name:      "number",
					Index:     0,
					Package:   1,
					Type:      types.I32,
					Size:      4,
					TotalSize: 4,
					Offset:    1048576,
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

			type_checks.ParseGlobals(Globals)

			program := actions.AST

			for _, wantGlobal := range tc.globalsCXArgs {

				for _, pkgIdx := range program.Packages {

					pkg, err := program.GetPackageFromArray(pkgIdx)

					if err != nil {
						t.Log(err)
					}

					if cxpackages.IsDefaultPackage(pkg.Name) {
						continue
					}

					var match bool
					var gotGlobal *ast.CXArgument

					for _, globalIdx := range pkg.Globals.Fields {
						gotGlobal = program.GetCXArg(ast.CXArgumentIndex(globalIdx))

						if gotGlobal.Name == wantGlobal.Name &&
							gotGlobal.Index == wantGlobal.Index &&
							gotGlobal.Package == wantGlobal.Package &&
							gotGlobal.Type == wantGlobal.Type {
							match = true
							break
						}

					}

					if !match {
						t.Errorf("want global %v, got %v", wantGlobal, gotGlobal)

					}
				}

			}
		})
	}

}

/*
func TestTypeChecks_ParseEnums(t *testing.T) {

	tests := []struct {
		scenario string
		testDir  string
	}{}

	for _, tc := range tests {
		t.Run(tc.scenario, func(t *testing.T) {

		})
	}

}

func TestTypeChecks_ParseStructs(t *testing.T) {

	tests := []struct {
		scenario  string
		testDir   string
		structCXs []ast.CXStruct
	}{
		{
			scenario: "Has Structs",
			testDir:  "./test_files/test.cx",
			structCXs: []ast.CXStruct{
				{
					Name:    "CustomType",
					Index:   1,
					Package: 1,
					Fields: []ast.CXTypeSignature{
						{
							Name:   "fieldA",
							Offset: 8,
							Type:   ast.TYPE_CXARGUMENT_DEPRECATE,
							Meta:   28,
						},
						{
							Name:   "fieldB",
							Offset: 4,
							Type:   ast.TYPE_ATOMIC,
							Meta:   4,
						},
					},
				},
				{
					Name:    "AnotherType",
					Index:   2,
					Package: 1,
					Fields: []ast.CXTypeSignature{
						{
							Name:   "name",
							Offset: 8,
							Type:   ast.TYPE_CXARGUMENT_DEPRECATE,
							Meta:   1,
						},
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
			pkg := declaration_extractor.ExtractPackages(ReplaceCommentsWithWhitespaces)

			structs, err := declaration_extractor.ExtractStructs(ReplaceCommentsWithWhitespaces, tc.testDir, pkg)

			type_checks.ParseStructs(structs)

			program := actions.AST

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
						gotStruct := program.CXStructs[structIdx]

						if gotStruct.Name == wantStruct.Name &&
							(gotStruct.Index == wantStruct.Index ||
								gotStruct.Package == wantStruct.Package) {

							for k, wantFields := range wantStruct.Fields {

								if wantFields == gotStruct.Fields[k] {
									match = true
									break
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
						Fields: []ast.CXTypeSignature{
							{
								Name: "a",
								Type: ast.TYPE_CXARGUMENT_DEPRECATE,
								Meta: 0,
							},
							{
								Name: "b",
								Type: ast.TYPE_CXARGUMENT_DEPRECATE,
								Meta: 1,
							},
						},
					},
					Outputs: &ast.CXStruct{
						Fields: []ast.CXTypeSignature{
							{
								Name: "answer",
								Type: ast.TYPE_CXARGUMENT_DEPRECATE,
								Meta: 2,
							},
						},
					},
				},
				{
					Name:    "divide",
					Index:   2,
					Package: 1,
					Inputs: &ast.CXStruct{
						Fields: []ast.CXTypeSignature{
							{
								Name: "c",
								Type: ast.TYPE_CXARGUMENT_DEPRECATE,
								Meta: 3,
							},
							{
								Name: "d",
								Type: ast.TYPE_CXARGUMENT_DEPRECATE,
								Meta: 4,
							},
						},
					},
					Outputs: &ast.CXStruct{
						Fields: []ast.CXTypeSignature{
							{
								Name: "quotient",
								Type: ast.TYPE_CXARGUMENT_DEPRECATE,
								Meta: 5,
							},
							{
								Name: "remainder",
								Type: ast.TYPE_CXARGUMENT_DEPRECATE,
								Meta: 6,
							},
						},
					},
				},
				{
					Name:    "printer",
					Index:   3,
					Package: 1,
					Inputs: &ast.CXStruct{
						Fields: []ast.CXTypeSignature{
							{
								Name: "message",
								Type: ast.TYPE_CXARGUMENT_DEPRECATE,
								Meta: 7,
							},
						},
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
			pkg := declaration_extractor.ExtractPackages(ReplaceCommentsWithWhitespaces)

			funcs, err := declaration_extractor.ExtractFuncs(srcBytes, tc.testDir, pkg)

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
*/
