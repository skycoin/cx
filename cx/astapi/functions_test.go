package astapi_test

import (
	"testing"

	cxast "github.com/skycoin/cx/cx/ast"
	"github.com/skycoin/cx/cx/astapi"
	cxconstants "github.com/skycoin/cx/cx/constants"
)

func TestASTAPI_Functions(t *testing.T) {
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

	t.Run("add empty function to package", func(t *testing.T) {
		err := astapi.AddEmptyFunctionToPackage(cxprogram, "main", "TestFunction")
		if err != nil {
			t.Errorf("want no error, got %v", err)
		}
	})

	t.Run("add first input to function", func(t *testing.T) {
		err := astapi.AddNativeInputToFunction(cxprogram, "main", "TestFunction", "inputOne", cxconstants.TYPE_I32)
		if err != nil {
			t.Errorf("want no error, got %v", err)
		}
	})

	t.Run("add second input to function", func(t *testing.T) {
		err := astapi.AddNativeInputToFunction(cxprogram, "main", "TestFunction", "inputTwo", cxconstants.TYPE_I32)
		if err != nil {
			t.Errorf("want no error, got %v", err)
		}
	})

	t.Run("add first output to function", func(t *testing.T) {
		err := astapi.AddNativeOutputToFunction(cxprogram, "main", "TestFunction", "outputOne", cxconstants.TYPE_I16)
		if err != nil {
			t.Errorf("want no error, got %v", err)
		}
	})

	t.Run("add second output to function", func(t *testing.T) {
		err := astapi.AddNativeOutputToFunction(cxprogram, "main", "TestFunction", "outputTwo", cxconstants.TYPE_I16)
		if err != nil {
			t.Errorf("want no error, got %v", err)
		}
	})

	t.Run("print program", func(t *testing.T) {
		cxprogram.PrintProgram()
	})

	t.Run("remove first input of function", func(t *testing.T) {
		err := astapi.RemoveFunctionInput(cxprogram, "TestFunction", "inputOne")
		if err != nil {
			t.Errorf("want no error, got %v", err)
		}
	})

	t.Run("remove first output of function", func(t *testing.T) {
		err := astapi.RemoveFunctionOutput(cxprogram, "TestFunction", "outputOne")
		if err != nil {
			t.Errorf("want no error, got %v", err)
		}
	})

	t.Run("print program", func(t *testing.T) {
		cxprogram.PrintProgram()
	})

	t.Run("remove function from package", func(t *testing.T) {
		err := astapi.RemoveFunctionFromPackage(cxprogram, "main", "TestFunction")
		if err != nil {
			t.Errorf("want no error, got %v", err)
		}
	})

	t.Run("print program", func(t *testing.T) {
		cxprogram.PrintProgram()
	})
}
