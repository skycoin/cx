package astapi_test

import (
	"testing"

	cxast "github.com/skycoin/cx/cx/ast"
	"github.com/skycoin/cx/cx/astapi"
	cxconstants "github.com/skycoin/cx/cx/constants"
	parsingcompletor "github.com/skycoin/cx/cxparser/cxparsingcompletor"
)

func TestASTAPI_Arguments(t *testing.T) {
	var cxprogram *cxast.CXProgram

	// Needed for AddNativeExpressionToFunction
	// because of dependency on cxast.OpNames
	parsingcompletor.InitCXCore()

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

	t.Run("add input to function", func(t *testing.T) {
		err := astapi.AddNativeInputToFunction(cxprogram, "main", "TestFunction", "inputOne", cxconstants.TYPE_I32)
		if err != nil {
			t.Errorf("want no error, got %v", err)
		}
	})

	t.Run("add output to function", func(t *testing.T) {
		err := astapi.AddNativeOutputToFunction(cxprogram, "main", "TestFunction", "outputOne", cxconstants.TYPE_I16)
		if err != nil {
			t.Errorf("want no error, got %v", err)
		}
	})

	t.Run("add expression to function", func(t *testing.T) {
		err := astapi.AddNativeExpressionToFunction(cxprogram, "TestFunction", cxconstants.OP_ADD)
		if err != nil {
			t.Errorf("want no error, got %v", err)
		}
	})

	t.Run("add first input to expression", func(t *testing.T) {
		err := astapi.AddNativeInputToExpression(cxprogram, "main", "TestFunction", "x", cxconstants.TYPE_I16, 0)
		if err != nil {
			t.Errorf("want no error, got %v", err)
		}
	})

	t.Run("add second input to expression", func(t *testing.T) {
		err := astapi.AddNativeInputToExpression(cxprogram, "main", "TestFunction", "y", cxconstants.TYPE_I16, 0)
		if err != nil {
			t.Errorf("want no error, got %v", err)
		}
	})

	t.Run("add output to expression", func(t *testing.T) {
		err := astapi.AddNativeOutputToExpression(cxprogram, "main", "TestFunction", "z", cxconstants.TYPE_I16, 0)
		if err != nil {
			t.Errorf("want no error, got %v", err)
		}
	})

	t.Run("print program", func(t *testing.T) {
		cxprogram.PrintProgram()
	})

	t.Run("make input of an expression a pointer", func(t *testing.T) {
		err := astapi.MakeInputExpressionAPointer(cxprogram, "TestFunction", 0, 1)
		if err != nil {
			t.Errorf("want no error, got %v", err)
		}
	})

	t.Run("make output of an expression a pointer", func(t *testing.T) {
		err := astapi.MakeOutputExpressionAPointer(cxprogram, "TestFunction", 0, 0)
		if err != nil {
			t.Errorf("want no error, got %v", err)
		}
	})

	t.Run("print program", func(t *testing.T) {
		cxprogram.PrintProgram()
	})

	t.Run("remove input from an expression", func(t *testing.T) {
		err := astapi.RemoveInputFromExpression(cxprogram, "TestFunction", 0)
		if err != nil {
			t.Errorf("want no error, got %v", err)
		}
	})

	t.Run("remove output from an expression", func(t *testing.T) {
		err := astapi.RemoveOutputFromExpression(cxprogram, "TestFunction", 0)
		if err != nil {
			t.Errorf("want no error, got %v", err)
		}
	})

	t.Run("print program", func(t *testing.T) {
		cxprogram.PrintProgram()
	})

	t.Run("add global", func(t *testing.T) {
		arg := cxast.MakeGlobal("testGlobal", cxconstants.TYPE_I16, "", -1)
		cxprogram.Packages[0].AddGlobal(arg)
	})

	t.Run("get accessible i16 args", func(t *testing.T) {
		args, err := astapi.GetAccessibleArgsForFunctionByType(cxprogram, "main", "TestFunction", cxconstants.TYPE_I16)
		if err != nil {
			t.Errorf("want no error, got %v", err)
		}

		for _, arg := range args {
			t.Logf("Arg (name,type)=(%v,%v)\n", arg.ArgDetails.Name, cxconstants.TypeNames[arg.Type])
		}
	})
}
