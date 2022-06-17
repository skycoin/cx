package opcodes

import (
	"fmt"

	"github.com/skycoin/cx/cx/ast"
	"github.com/skycoin/cx/cx/constants"
	"github.com/skycoin/cx/cx/types"
)

var assertSuccess = true

// AssertFailed ...
func AssertFailed() bool {
	return !assertSuccess
}

//TODO: Rework
func assert(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) (same bool) {
	var byts1, byts2 []byte
	var inp0 *ast.CXArgument
	if inputs[0].TypeSignature.Type == ast.TYPE_CXARGUMENT_DEPRECATE {
		inp0 = prgrm.GetCXArgFromArray(ast.CXArgumentIndex(inputs[0].TypeSignature.Meta))
	} else {
		panic("type is not type cx argument deprecate\n\n")
	}
	if inp0.Type == types.STR || inp0.PointerTargetType == types.STR {
		byts1 = []byte(inputs[0].Get_str(prgrm))
		byts2 = []byte(inputs[1].Get_str(prgrm))
	} else {
		byts1 = inputs[0].Get_bytes(prgrm)
		byts2 = inputs[1].Get_bytes(prgrm)
	}

	same = true

	if len(byts1) != len(byts2) {
		same = false
		fmt.Println("byts1", byts1)
		fmt.Println("byts2", byts2)
	}

	if same {
		for i, byt := range byts1 {
			if byt != byts2[i] {
				same = false
				fmt.Println("byts1", byts1)
				fmt.Println("byts2", byts2)
				break
			}
		}
	}

	message := inputs[2].Get_str(prgrm)

	if !same {
		call := prgrm.GetCurrentCall()
		// expr := call.Operator.Expressions[call.Line]
		cxLine, _ := prgrm.GetPreviousCXLine(call.Operator.Expressions, call.Line)

		if message != "" {
			fmt.Printf("%s: %d: result was not equal to the expected value; %s\n", cxLine.FileName, cxLine.LineNumber, message)
		} else {
			fmt.Printf("%s: %d: result was not equal to the expected value\n", cxLine.FileName, cxLine.LineNumber)
		}
	}

	assertSuccess = assertSuccess && same
	return same
}

func opAssertValue(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	same := assert(prgrm, inputs, outputs)
	outputs[0].Set_bool(prgrm, same)
}

func opTest(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	assert(prgrm, inputs, outputs)
}

func opPanic(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	if !assert(prgrm, inputs, outputs) {
		panic(constants.CX_ASSERT)
	}
}

// panicIf/panicIfNot implementation
func panicIf(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue, condition bool) {
	str := inputs[1].Get_str(prgrm)
	if inputs[0].Get_bool(prgrm) == condition {
		call := prgrm.GetCurrentCall()
		// expr := call.Operator.Expressions[call.Line]
		cxLine, _ := prgrm.GetPreviousCXLine(call.Operator.Expressions, call.Line)

		fmt.Printf("%s : %d, %s\n", cxLine.FileName, cxLine.LineNumber, str)
		panic(constants.CX_ASSERT)
	}
}

// panic with CX_ASSERT exit code if condition is true
func opPanicIf(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	panicIf(prgrm, inputs, outputs, true)
}

// panic with CX_ASSERT exit code if condition is false
func opPanicIfNot(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	panicIf(prgrm, inputs, outputs, false)
}

func opStrError(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outputs[0].Set_str(prgrm, ast.ErrorString(int(inputs[0].Get_i32(prgrm))))
}
