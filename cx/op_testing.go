package cxcore

import (
	"fmt"
	"github.com/skycoin/cx/cx/ast"
	"github.com/skycoin/cx/cx/constants"
)

var assertSuccess = true

// AssertFailed ...
func AssertFailed() bool {
	return !assertSuccess
}

func assert(expr *ast.CXExpression, fp int) (same bool) {
	inp1, inp2, inp3 := expr.Inputs[0], expr.Inputs[1], expr.Inputs[2]
	var byts1, byts2 []byte

	if inp1.Type == constants.TYPE_STR {
		byts1 = []byte(ast.ReadStr(fp, inp1))
		byts2 = []byte(ast.ReadStr(fp, inp2))
	} else {
		byts1 = ast.ReadMemory(ast.GetFinalOffset(fp, inp1), inp1)
		byts2 = ast.ReadMemory(ast.GetFinalOffset(fp, inp2), inp2)
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

	message := ast.ReadStr(fp, inp3)

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

func opAssertValue(expr *ast.CXExpression, fp int) {
	same := assert(expr, fp)
	ast.WriteBool(ast.GetFinalOffset(fp, expr.Outputs[0]), same)
}

func opTest(expr *ast.CXExpression, fp int) {
	assert(expr, fp)
}

func opPanic(expr *ast.CXExpression, fp int) {
	if !assert(expr, fp) {
		panic(constants.CX_ASSERT)
	}
}

// panicIf/panicIfNot implementation
func panicIf(expr *ast.CXExpression, fp int, condition bool) {
	if ast.ReadBool(fp, expr.Inputs[0]) == condition {
		fmt.Printf("%s : %d, %s\n", expr.FileName, expr.FileLine, ast.ReadStr(fp, expr.Inputs[1]))
		panic(constants.CX_ASSERT)
	}
}

// panic with CX_ASSERT exit code if condition is true
func opPanicIf(expr *ast.CXExpression, fp int) {
	panicIf(expr, fp, true)
}

// panic with CX_ASSERT exit code if condition is false
func opPanicIfNot(expr *ast.CXExpression, fp int) {
	panicIf(expr, fp, false)
}

func opStrError(expr *ast.CXExpression, fp int) {
	ast.WriteString(fp, ast.ErrorString(int(ast.ReadI32(fp, expr.Inputs[0]))), expr.Outputs[0])
}
