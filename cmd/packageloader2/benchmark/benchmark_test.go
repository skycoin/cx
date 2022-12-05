package benchmark_test

import (
	"testing"

	"github.com/skycoin/cx/cmd/packageloader2/file_output"
	"github.com/skycoin/cx/cmd/packageloader2/loader"
	"github.com/skycoin/cx/cxparser/actions"
)

func BenchmarkPackageloaderBolt(b *testing.B) {
	for i := 0; i < b.N; i++ {
		actions.AST = nil

		_, sourceCode, _, rootDir := loader.ParseArgsForCX([]string{"./test_files/test.cx"}, true)

		err := loader.LoadCXProgram("test", sourceCode, rootDir, "bolt")
		if err != nil {
			b.Fatal(err)
		}

		_, err = file_output.GetImportFiles("test", "bolt")
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkPackageloaderRedis(b *testing.B) {
	for i := 0; i < b.N; i++ {
		actions.AST = nil

		_, sourceCode, _, rootDir := loader.ParseArgsForCX([]string{"./test_files/test.cx"}, true)

		err := loader.LoadCXProgram("test", sourceCode, rootDir, "redis")
		if err != nil {
			b.Fatal(err)
		}

		_, err = file_output.GetImportFiles("test", "redis")
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkPackageloaderNoSave(b *testing.B) {
	for i := 0; i < b.N; i++ {
		actions.AST = nil

		_, sourceCode, _, rootDir := loader.ParseArgsForCX([]string{"./test_files/test.cx"}, true)

		_, err := loader.LoadCXProgramNoSave(sourceCode, rootDir)
		if err != nil {
			b.Fatal(err)
		}
	}
}
