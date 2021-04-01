package cxcore

import (
	
    "github.com/skycoin/cx/cx/ast"
	"github.com/skycoin/cx/cx/constants"
	//"github.com/skycoin/cx/cx/helper"
)

// "fmt"

// EscapeAnalysis ...
func EscapeAnalysis(inpOffset, outOffset int, arg *ast.CXArgument) {
	heapOffset := ast.AllocateSeq(arg.TotalSize + constants.OBJECT_HEADER_SIZE)

	byts := ast.ReadMemory(inpOffset, arg)

	// creating a header for this object
	var header = make([]byte, constants.OBJECT_HEADER_SIZE)
	ast.WriteMemI32(header, 5, int32(len(byts)))

	obj := append(header, byts...)
	ast.WriteMemory(heapOffset, obj)

	ast.WriteI32(outOffset, int32(heapOffset))
}

func opIdentity(inputs []ast.CXValue, outputs []ast.CXValue) {
    inp1Offset := inputs[0].Offset
	out1Offset := outputs[0].Offset

    out1 := outputs[0].Arg
	var elt *ast.CXArgument
	if len(out1.Fields) > 0 {
		elt = out1.Fields[len(out1.Fields)-1]
	} else {
		elt = out1
	}

	if elt.DoesEscape {
		EscapeAnalysis(inp1Offset, out1Offset, inputs[0].Arg)
	} else {
		switch elt.PassBy {
		case constants.PASSBY_VALUE:
			ast.WriteMemory(out1Offset, ast.ReadMemory(inp1Offset, inputs[0].Arg))
		case constants.PASSBY_REFERENCE:
			ast.WriteI32(out1Offset, int32(inp1Offset))
		}
	}
}

func opGoto(inputs []ast.CXValue, outputs []ast.CXValue) {
	call := ast.PROGRAM.GetCurrentCall()
	expr := call.Operator.Expressions[call.Line]
	call.Line = call.Line + expr.ThenLines
}

func opJmp(inputs []ast.CXValue, outputs []ast.CXValue) {
	call := ast.PROGRAM.GetCurrentCall()
	expr := call.Operator.Expressions[call.Line]
	if inputs[0].Get_bool() {
		call.Line = call.Line + expr.ThenLines
	} else {
		call.Line = call.Line + expr.ElseLines
	}
}
