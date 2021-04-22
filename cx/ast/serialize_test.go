package ast_test

import (
	"fmt"
	"testing"

	cxevolves "github.com/skycoin/cx-evolves/evolve"
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
			program:  generateSampleProgramFromCXEvolves(t),
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
			program:  generateSampleProgramFromCXEvolves(t),
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

func TestSerialize_SkyEncoder_VS_CipherEncoder(t *testing.T) {
	// prgrm := readCXProgramFromTestData(t, "serialized_1")
	prgrm := generateSampleProgramFromCXEvolves(t)
	// prgrm.PrintProgram()
	tests := []struct {
		scenario string
		program  *cxast.CXProgram
		wantSame bool
	}{
		{
			scenario: "Valid program",
			program:  prgrm,
		},
	}

	for _, tc := range tests {
		t.Run(tc.scenario, func(t *testing.T) {
			skyEncoderSerializedBytes := cxast.SerializeCXProgramV2(tc.program, false)
			cipherEncoderSerializedBytes := cxast.SerializeCXProgram(tc.program, false)

			// Check byte per byte
			var diffCount int = 0
			var pos []int
			for i := range skyEncoderSerializedBytes {
				if skyEncoderSerializedBytes[i] != cipherEncoderSerializedBytes[i] {
					pos = append(pos, i)
					diffCount++
				}
			}
			t.Logf("There were %v indexes that have different values. \nThese indexes are %v", diffCount, pos)
		})
	}
}

func generateSampleProgramFromCXEvolves(t *testing.T) *cxast.CXProgram {
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

	err = astapi.AddNativeOutputToFunction(cxProgram, "main", "TestFunction", "outputOne", cxconstants.TYPE_I32)
	if err != nil {
		t.Errorf("want no error, got %v", err)
	}
	functionSetNames := []string{"i32.add", "i32.mul", "i32.sub", "i32.eq", "i32.uneq", "i32.gt", "i32.gteq", "i32.lt", "i32.lteq", "bool.not", "bool.or", "bool.and", "bool.uneq", "bool.eq", "i32.neg", "i32.abs", "i32.bitand", "i32.bitor", "i32.bitxor", "i32.bitclear", "i32.bitshl", "i32.bitshr", "i32.max", "i32.min", "i32.rand"}
	fns := cxevolves.GetFunctionSet(functionSetNames)

	fn, _ := cxProgram.GetFunction("TestFunction", "main")
	fmt.Printf("got func=%v\n", fn)
	pkg, _ := cxProgram.GetPackage("main")
	fmt.Printf("got pkg=%v\n", pkg)
	fmt.Printf("len expr:=%v\n", len(fn.Expressions))
	cxevolves.GenerateRandomExpressions(fn, pkg, fns, 30)

	cxProgram.PrintProgram()
	return cxProgram
}
