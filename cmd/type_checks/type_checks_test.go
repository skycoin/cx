package type_checks_test

import (
	"os"
	"testing"

	"github.com/skycoin/cx/cmd/declaration_extraction"
	"github.com/skycoin/cx/cmd/type_checks"
	"github.com/skycoin/cx/cx/ast"
	"github.com/skycoin/cx/cx/packages"
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
					Size:      18446744073709551615,
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
	}

	for _, tc := range tests {
		t.Run(tc.scenario, func(t *testing.T) {

			srcBytes, err := os.ReadFile(tc.testDir)
			if err != nil {
				t.Error(err)
			}

			ReplaceCommentsWithWhitespaces := declaration_extraction.ReplaceCommentsWithWhitespaces(srcBytes)
			pkg := declaration_extraction.ExtractPackages(ReplaceCommentsWithWhitespaces)

			Globals, err := declaration_extraction.ExtractGlobals(ReplaceCommentsWithWhitespaces, tc.testDir, pkg)

			type_checks.ParseGlobals(Globals)

			var i int

			program := actions.AST

			for _, pkgIdx := range program.Packages {
				pkg, err := program.GetPackageFromArray(pkgIdx)
				if err != nil {
					panic(err)
				}

				if packages.IsDefaultPackage(pkg.Name) {
					continue
				}

				for _, globalIdx := range pkg.Globals.Fields {
					global := program.GetCXArg(ast.CXArgumentIndex(globalIdx.Meta))
					testGlobal := tc.globalsCXArgs[i]

					if global.Name != testGlobal.Name ||
						global.Index != testGlobal.Index ||
						global.Package != testGlobal.Package ||
						global.Type != testGlobal.Type ||
						global.Size != testGlobal.Size ||
						global.Offset != testGlobal.Offset {
						t.Errorf("want global %v, got %v", testGlobal, global)
					}

					i++
				}
			}

		})
	}

}

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
		scenario string
		testDir  string
	}{
		{
			scenario: "Has Structs",
			testDir:  "./test_files/test.cx",
		},
	}

	for _, tc := range tests {
		t.Run(tc.scenario, func(t *testing.T) {

			srcBytes, err := os.ReadFile(tc.testDir)
			if err != nil {
				t.Error(err)
			}

			ReplaceCommentsWithWhitespaces := declaration_extraction.ReplaceCommentsWithWhitespaces(srcBytes)
			pkg := declaration_extraction.ExtractPackages(ReplaceCommentsWithWhitespaces)

			structs, err := declaration_extraction.ExtractStructs(ReplaceCommentsWithWhitespaces, tc.testDir, pkg)

			type_checks.ParseStructs(structs)

		})
	}

}

func TestTypeChecks_ParseFuncHeaders(t *testing.T) {

	tests := []struct {
		scenario string
		testDir  string
	}{
		{
			scenario: "Has funcs",
			testDir:  "./test_files/test.cx",
		},
	}

	for _, tc := range tests {
		t.Run(tc.scenario, func(t *testing.T) {

			srcBytes, err := os.ReadFile(tc.testDir)
			if err != nil {
				t.Error(err)
			}

			ReplaceCommentsWithWhitespaces := declaration_extraction.ReplaceCommentsWithWhitespaces(srcBytes)
			pkg := declaration_extraction.ExtractPackages(ReplaceCommentsWithWhitespaces)

			funcs, err := declaration_extraction.ExtractFuncs(srcBytes, tc.testDir, pkg)

			type_checks.ParseFuncHeaders(funcs)

		})
	}

}
