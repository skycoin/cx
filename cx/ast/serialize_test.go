package ast_test

import (
	"testing"

	cxast "github.com/skycoin/cx/cx/ast"
	"github.com/skycoin/cx/cx/astapi"
	cxconstants "github.com/skycoin/cx/cx/constants"
	cxparsingcompletor "github.com/skycoin/cx/cxparser/cxparsingcompletor"
)

func TestSerialize_CipherEncoder(t *testing.T) {
	tests := []struct {
		scenario string
		program  *cxast.CXProgram
		wantErr  error
	}{
		{
			scenario: "Valid program",
			program:  generateSampleProgram(t),
		},
	}

	for _, tc := range tests {
		t.Run(tc.scenario, func(t *testing.T) {
			serializedBytes := cxast.SerializeCXProgram(tc.program, false)
			deserializedCXProgram := cxast.Deserialize(serializedBytes)

			if cxast.ToString(deserializedCXProgram) != cxast.ToString(tc.program) {
				t.Errorf("want same program, got different")
			}
		})
	}
}

func TestSerialize_SkyEncoder(t *testing.T) {
	tests := []struct {
		scenario string
		program  *cxast.CXProgram
		wantSame bool
	}{
		{
			scenario: "Valid program",
			program:  generateSampleProgram(t),
		},
	}

	for _, tc := range tests {
		t.Run(tc.scenario, func(t *testing.T) {
			serializedBytes := cxast.SerializeCXProgramV2(tc.program, false)
			deserializedCXProgram := cxast.DeserializeCXProgramV2(serializedBytes)
			if cxast.ToString(deserializedCXProgram) != cxast.ToString(tc.program) {
				t.Errorf("want same program, got different")
			}
		})
	}
}

func generateSampleProgram(t *testing.T) *cxast.CXProgram {
	var cxProgram *cxast.CXProgram

	// Needed for AddNativeExpressionToFunction
	// because of dependency on cxast.OpNames
	cxparsingcompletor.InitCXCore()
	cxProgram = cxast.MakeProgram()

	err := astapi.AddEmptyPackage(cxProgram, "main")
	if err != nil {
		t.Errorf("want no error, got %v", err)
	}

	err = astapi.AddEmptyFunctionToPackage(cxProgram, "main", "TestFunction")
	if err != nil {
		t.Errorf("want no error, got %v", err)
	}

	err = astapi.AddNativeInputToFunction(cxProgram, "main", "TestFunction", "inputOne", cxconstants.TYPE_I32)
	if err != nil {
		t.Errorf("want no error, got %v", err)
	}

	err = astapi.AddNativeOutputToFunction(cxProgram, "main", "TestFunction", "outputOne", cxconstants.TYPE_I16)
	if err != nil {
		t.Errorf("want no error, got %v", err)
	}

	err = astapi.AddNativeExpressionToFunction(cxProgram, "TestFunction", cxconstants.OP_ADD)
	if err != nil {
		t.Errorf("want no error, got %v", err)
	}

	err = astapi.AddNativeInputToExpression(cxProgram, "main", "TestFunction", "x", cxconstants.TYPE_I16, 0)
	if err != nil {
		t.Errorf("want no error, got %v", err)
	}

	err = astapi.AddNativeInputToExpression(cxProgram, "main", "TestFunction", "y", cxconstants.TYPE_I16, 0)
	if err != nil {
		t.Errorf("want no error, got %v", err)
	}

	err = astapi.AddNativeOutputToExpression(cxProgram, "main", "TestFunction", "z", cxconstants.TYPE_I16, 0)
	if err != nil {
		t.Errorf("want no error, got %v", err)
	}

	err = astapi.AddNativeExpressionToFunction(cxProgram, "TestFunction", cxconstants.OP_SUB)
	if err != nil {
		t.Errorf("want no error, got %v", err)
	}

	err = astapi.AddNativeInputToExpression(cxProgram, "main", "TestFunction", "x", cxconstants.TYPE_I16, 1)
	if err != nil {
		t.Errorf("want no error, got %v", err)
	}

	err = astapi.AddNativeInputToExpression(cxProgram, "main", "TestFunction", "y", cxconstants.TYPE_I16, 1)
	if err != nil {
		t.Errorf("want no error, got %v", err)
	}

	err = astapi.AddNativeOutputToExpression(cxProgram, "main", "TestFunction", "z", cxconstants.TYPE_I16, 1)
	if err != nil {
		t.Errorf("want no error, got %v", err)
	}

	return cxProgram
}
