package benchmark_test

import (
	"testing"

	"github.com/skycoin/cx/cmd/packageloader/file_output"
	"github.com/skycoin/cx/cmd/packageloader/loader"
	"github.com/skycoin/cx/cxparser/actions"
)

func BenchmarkPackageloader(b *testing.B) {
	for i := 0; i < b.N; i++ {
		actions.AST = nil

		_, sourceCode, _ := loader.ParseArgsForCX([]string{"./test_files/test.cx"}, true)

		err := loader.LoadCXProgram("test", sourceCode, "bolt")
		if err != nil {
			b.Fatal(err)
		}

		_, err = file_output.GetImportFiles("test", "bolt")
		if err != nil {
			b.Fatal(err)
		}
	}
}
