package astapi_test

import (
	"testing"

	cxast "github.com/skycoin/cx/cx/ast"
	"github.com/skycoin/cx/cx/astapi"
)

func TestASTAPI_Packages(t *testing.T) {
	var cxprogram *cxast.CXProgram
	t.Run("make program", func(t *testing.T) {
		cxprogram = cxast.MakeProgram()
	})

	t.Run("add empty package", func(t *testing.T) {
		err := astapi.AddEmptyPackage(cxprogram, "main")
		if err != nil {
			t.Errorf("want no error, got %v", err)
		}
	})

	t.Run("get packages name lists", func(t *testing.T) {
		list := astapi.GetPackagesNameList(cxprogram)
		wantLen := 1
		if len(list) != wantLen {
			t.Errorf("got len(list)=%v, want %v", len(list), wantLen)
		}
	})

	t.Run("print program", func(t *testing.T) {
		cxprogram.PrintProgram()
	})
}
