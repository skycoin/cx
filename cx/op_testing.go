package cxcore

import (
	"fmt"
)

var assertSuccess = true

// AssertFailed ...
func AssertFailed() bool {
	return !assertSuccess
}

func assert(expr *CXExpression, fp int) (same bool) {
	inp1, inp2, inp3 := expr.Inputs[0], expr.Inputs[1], expr.Inputs[2]
	var byts1, byts2 []byte

	if inp1.Type == TYPE_STR {
		byts1 = []byte(ReadStr(fp, inp1))
		byts2 = []byte(ReadStr(fp, inp2))
	} else {
		byts1 = ReadMemory(GetFinalOffset(fp, inp1), inp1)
		byts2 = ReadMemory(GetFinalOffset(fp, inp2), inp2)
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

	message := ReadStr(fp, inp3)

	if !same {
		if message != "" {
			fmt.Printf("%s: %d: result was not equal to the expected value; %s\n", expr.FileName, expr.FileLine, message)
		} else {
			fmt.Printf("%s: %d: result was not equal to the expected value\n", expr.FileName, expr.FileLine)
		}
	}

	assertSuccess = assertSuccess && same
	return same
}

func opAssertValue(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	same := assert(expr, fp)
	WriteBool(GetFinalOffset(fp, expr.Outputs[0]), same)
}

func opTest(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	assert(expr, fp)
}

func opPanic(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	if !assert(expr, fp) {
		panic(CX_ASSERT)
	}
}

// panicIf/panicIfNot implementation
func panicIf(prgrm *CXProgram, condition bool) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()
	if ReadBool(fp, expr.Inputs[0]) == condition {
		fmt.Printf("%s : %d, %s\n", expr.FileName, expr.FileLine, ReadStr(fp, expr.Inputs[1]))
		panic(CX_ASSERT)
	}
}

// panic with CX_ASSERT exit code if condition is true
func opPanicIf(prgrm *CXProgram) {
	panicIf(prgrm, true)
}

// panic with CX_ASSERT exit code if condition is false
func opPanicIfNot(prgrm *CXProgram) {
	panicIf(prgrm, false)
}

func opStrError(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	WriteObject(GetFinalOffset(fp, expr.Outputs[0]), FromStr(ErrorString(int(ReadI32(fp, expr.Inputs[0])))))
}
