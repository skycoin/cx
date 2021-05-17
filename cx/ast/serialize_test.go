package ast_test

import (
	"bytes"
	"encoding/binary"
	"reflect"
	"strings"
	"testing"
	"time"

	cxevolvesgenerator "github.com/skycoin/cx-evolves/generator"
	"github.com/skycoin/cx/cx/ast"
	cxast "github.com/skycoin/cx/cx/ast"
	"github.com/skycoin/cx/cx/astapi"
	cxconstants "github.com/skycoin/cx/cx/constants"
	cxparsingcompletor "github.com/skycoin/cx/cxparser/cxparsingcompletor"
)

// t.Logf("Stack size=%v\n", tc.program.StackSize)
// t.Logf("dataSegment size=%v\n", tc.program.DataSegmentSize)
// t.Logf("dataSegment starts at=%v\n", tc.program.DataSegmentStartsAt)
// t.Logf("heap size=%v\n", tc.program.HeapSize)
// t.Logf("heap starts at=%v\n", tc.program.HeapStartsAt)
// t.Logf("memory len=%v\n", len(tc.program.Memory))
// t.Logf("value at memory data offset=%v\n", tc.program.Memory[tc.program.DataSegmentStartsAt:tc.program.DataSegmentStartsAt+tc.program.DataSegmentSize])

func TestSerialize_CipherEncoder(t *testing.T) {
	tests := []struct {
		scenario          string
		program           *cxast.CXProgram
		includeDataMemory bool
		useCompression    bool
	}{
		{
			scenario:          "program without literal",
			program:           generateSampleProgramFromCXEvolves(t, false),
			includeDataMemory: false,
			useCompression:    false,
		},
		{
			scenario:          "program with literal",
			program:           generateSampleProgramFromCXEvolves(t, true),
			includeDataMemory: true,
			useCompression:    false,
		},
		{
			scenario:          "program without literal, use compression",
			program:           generateSampleProgramFromCXEvolves(t, false),
			includeDataMemory: false,
			useCompression:    true,
		},
		{
			scenario:          "program with literal, use compression",
			program:           generateSampleProgramFromCXEvolves(t, true),
			includeDataMemory: true,
			useCompression:    true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.scenario, func(t *testing.T) {
			timeStart := time.Now()
			serializedBytes := cxast.SerializeCXProgram(tc.program, tc.includeDataMemory, tc.useCompression)
			t.Logf("Serialization took=%v\n", time.Since(timeStart))

			timeStart = time.Now()
			deserializedCXProgram := cxast.Deserialize(serializedBytes, tc.useCompression)
			t.Logf("Deserialization took=%v\n", time.Since(timeStart))

			if cxast.ToString(deserializedCXProgram) != cxast.ToString(tc.program) {
				t.Errorf("want same program, got different")
			}

			if tc.includeDataMemory {
				wantDataSegmentMemory := tc.program.Memory[tc.program.DataSegmentStartsAt : tc.program.DataSegmentStartsAt+tc.program.DataSegmentSize]
				gotDataSegmentMemory := deserializedCXProgram.Memory[deserializedCXProgram.DataSegmentStartsAt : deserializedCXProgram.DataSegmentStartsAt+deserializedCXProgram.DataSegmentSize]
				if !reflect.DeepEqual(wantDataSegmentMemory, gotDataSegmentMemory) {
					t.Errorf("want %v, got %v", wantDataSegmentMemory, gotDataSegmentMemory)
				}
			}
		})
	}
}

func TestSerialize_SkyEncoder(t *testing.T) {
	tests := []struct {
		scenario          string
		program           *cxast.CXProgram
		includeDataMemory bool
		useCompression    bool
	}{
		{
			scenario:          "program without literal",
			program:           generateSampleProgramFromCXEvolves(t, false),
			includeDataMemory: false,
			useCompression:    false,
		},
		{
			scenario:          "program with literal",
			program:           generateSampleProgramFromCXEvolves(t, true),
			includeDataMemory: true,
			useCompression:    false,
		},
		{
			scenario:          "program without literal, use compression",
			program:           generateSampleProgramFromCXEvolves(t, false),
			includeDataMemory: false,
			useCompression:    true,
		},
		{
			scenario:          "program with literal, use compression",
			program:           generateSampleProgramFromCXEvolves(t, true),
			includeDataMemory: true,
			useCompression:    true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.scenario, func(t *testing.T) {
			timeStart := time.Now()
			serializedBytes := cxast.SerializeCXProgramV2(tc.program, tc.includeDataMemory, tc.useCompression)
			t.Logf("Serialization took=%v\n", time.Since(timeStart))

			timeStart = time.Now()
			deserializedCXProgram := cxast.DeserializeCXProgramV2(serializedBytes, tc.useCompression)
			t.Logf("Deserialization took=%v\n", time.Since(timeStart))
			if cxast.ToString(deserializedCXProgram) != cxast.ToString(tc.program) {
				t.Errorf("want same program, got different")
			}

			if tc.includeDataMemory {
				wantDataSegmentMemory := tc.program.Memory[tc.program.DataSegmentStartsAt : tc.program.DataSegmentStartsAt+tc.program.DataSegmentSize]
				gotDataSegmentMemory := deserializedCXProgram.Memory[deserializedCXProgram.DataSegmentStartsAt : deserializedCXProgram.DataSegmentStartsAt+deserializedCXProgram.DataSegmentSize]
				if !reflect.DeepEqual(wantDataSegmentMemory, gotDataSegmentMemory) {
					t.Errorf("want %v, got %v", wantDataSegmentMemory, gotDataSegmentMemory)
				}
			}
		})
	}
}

func TestSerialize_SkyEncoder_VS_CipherEncoder(t *testing.T) {
	tests := []struct {
		scenario string
		program  *cxast.CXProgram
		wantSame bool
	}{
		{
			scenario: "Valid program",
			program:  generateSampleProgramFromCXEvolves(t, false),
		},
	}

	for _, tc := range tests {
		t.Run(tc.scenario, func(t *testing.T) {
			skyEncoderSerializedBytes := cxast.SerializeCXProgramV2(tc.program, false, false)
			cipherEncoderSerializedBytes := cxast.SerializeCXProgram(tc.program, false, false)

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

func TestCompression_LZ4(t *testing.T) {
	tests := []struct {
		scenario string
		data     []byte
	}{
		{
			scenario: "Valid data",
			data:     []byte(strings.Repeat("HelloWorld", 100)),
		},
	}

	for _, tc := range tests {
		t.Run(tc.scenario, func(t *testing.T) {
			originalLen := len(tc.data)
			t.Logf("Original length of data: %v\n", originalLen)

			ast.CompressBytesLZ4(&tc.data)
			compressedLen := len(tc.data)
			t.Logf("Compressed length of data: %v\n", compressedLen)

			ast.UncompressBytesLZ4(&tc.data)
			unCompressedLen := len(tc.data)
			t.Logf("Uncompressed length of data: %v\n", unCompressedLen)

			if compressedLen > originalLen {
				t.Errorf("want compressed length to be less than original length, got %v", compressedLen)
			}

			if originalLen != unCompressedLen {
				t.Errorf("want uncompressed length %v, got %v", originalLen, unCompressedLen)
			}
		})
	}
}

func generateSampleProgramFromCXEvolves(t *testing.T, withLiteral bool) *cxast.CXProgram {
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
	fns := cxevolvesgenerator.GetFunctionSet(functionSetNames)

	fn, _ := cxProgram.GetFunction("TestFunction", "main")
	pkg, _ := cxProgram.GetPackage("main")
	cxevolvesgenerator.GenerateRandomExpressions(fn, pkg, fns, 30)

	if withLiteral {
		buf := new(bytes.Buffer)
		var num int32 = 5
		binary.Write(buf, binary.LittleEndian, num)
		err = astapi.AddLiteralInputToExpression(cxProgram, "main", "TestFunction", buf.Bytes(), cxconstants.TYPE_I32, 2)
		if err != nil {
			t.Errorf("want no error, got %v", err)
		}

	}

	cxProgram.PrintProgram()
	return cxProgram
}
