package cxcore

import (
	"github.com/skycoin/cx/cx/ast"
	"github.com/skycoin/cx/cx/constants"
	"github.com/skycoin/cx/cx/helper"
)

// "fmt"

// EscapeAnalysis ...
func EscapeAnalysis(fp int, inpOffset, outOffset int, arg *ast.CXArgument) {
	heapOffset := ast.AllocateSeq(arg.TotalSize + constants.OBJECT_HEADER_SIZE)

	byts := ast.ReadMemory(inpOffset, arg)

	// creating a header for this object
	var header = make([]byte, constants.OBJECT_HEADER_SIZE)
	ast.WriteMemI32(header, 5, int32(len(byts)))

	obj := append(header, byts...)
	ast.WriteMemory(heapOffset, obj)

	ast.WriteI32(outOffset, int32(heapOffset))
}

func opIdentity(expr *ast.CXExpression, fp int) {
	inp1, out1 := expr.Inputs[0], expr.Outputs[0]
	inp1Offset := ast.GetFinalOffset(fp, inp1)
	out1Offset := ast.GetFinalOffset(fp, out1)

	var elt *ast.CXArgument
	if len(out1.Fields) > 0 {
		elt = out1.Fields[len(out1.Fields)-1]
	} else {
		elt = out1
	}

	if elt.DoesEscape {
		EscapeAnalysis(fp, inp1Offset, out1Offset, inp1)
	} else {
		switch elt.PassBy {
		case constants.PASSBY_VALUE:
			ast.WriteMemory(out1Offset, ast.ReadMemory(inp1Offset, inp1))
		case constants.PASSBY_REFERENCE:
			ast.WriteI32(out1Offset, int32(inp1Offset))
		}
	}
}

func opJmp(expr *ast.CXExpression, fp int) {
	call := ast.PROGRAM.GetCurrentCall()
	inp1 := expr.Inputs[0]
	var predicate bool

	if expr.Label != "" {
		// then it's a goto
		call.Line = call.Line + expr.ThenLines
	} else {
		inp1Offset := ast.GetFinalOffset(fp, inp1)

		predicateB := ast.PROGRAM.Memory[inp1Offset : inp1Offset+ast.GetSize(inp1)]
		predicate = helper.DeserializeBool(predicateB)

		if predicate {
			call.Line = call.Line + expr.ThenLines
		} else {
			call.Line = call.Line + expr.ElseLines
		}
	}
}
