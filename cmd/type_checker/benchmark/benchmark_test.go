package benchmark_test

import (
	"testing"

	"github.com/skycoin/cx/cmd/declaration_extractor"
	"github.com/skycoin/cx/cmd/packageloader2/file_output"
	"github.com/skycoin/cx/cmd/packageloader2/loader"
	"github.com/skycoin/cx/cmd/type_checker"
	"github.com/skycoin/cx/cxparser/actions"
)

func BenchmarkTypeCheckerBolt(b *testing.B) {
	for i := 0; i < b.N; i++ {
		actions.AST = nil

		_, sourceCode, _, rootDir := loader.ParseArgsForCX([]string{"./test_files/test.cx"}, true)

		err := loader.LoadCXProgram("test", sourceCode, rootDir, "bolt")
		if err != nil {
			b.Fatal(err)
		}
		files, err := file_output.GetImportFiles("test", "bolt")
		if err != nil {
			b.Fatal(err)
		}

		Imports, Globals, Enums, TypeDefinitions, Structs, Funcs, gotErr := declaration_extractor.ExtractAllDeclarations(files)
		if (Enums != nil && TypeDefinitions != nil) || gotErr != nil {
			b.Fatal(gotErr)
		}

		err = type_checker.ParseAllDeclarations(files, Imports, Globals, Structs, Funcs)
		if err != nil {
			b.Fatal(err)
		}

	}
}

func BenchmarkTypeCheckerRedis(b *testing.B) {
	for i := 0; i < b.N; i++ {
		actions.AST = nil

		_, sourceCode, _, rootDir := loader.ParseArgsForCX([]string{"./test_files/test.cx"}, true)

		err := loader.LoadCXProgram("test", sourceCode, rootDir, "redis")
		if err != nil {
			b.Fatal(err)
		}
		files, err := file_output.GetImportFiles("test", "redis")
		if err != nil {
			b.Fatal(err)
		}

		Imports, Globals, Enums, TypeDefinitions, Structs, Funcs, gotErr := declaration_extractor.ExtractAllDeclarations(files)
		if (Enums != nil && TypeDefinitions != nil) || gotErr != nil {
			b.Fatal(gotErr)
		}

		err = type_checker.ParseAllDeclarations(files, Imports, Globals, Structs, Funcs)
		if err != nil {
			b.Fatal(err)
		}

	}
}
