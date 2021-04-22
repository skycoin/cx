package opcodes

import (
	"github.com/skycoin/cx/cx/ast"
	"github.com/skycoin/cx/cx/constants"
)

// "fmt"

// EscapeAnalysis ...
//TODO: Comment this out
//TODO: Delete this, probably not needed
func EscapeAnalysis(input *ast.CXValue) int32 {
	heapOffset := ast.AllocateSeq(input.Arg.TotalSize + constants.OBJECT_HEADER_SIZE)

	byts := input.Get_bytes()

	// creating a header for this object
	var header = make([]byte, constants.OBJECT_HEADER_SIZE)
	ast.WriteMemI32(header, 5, int32(len(byts)))

	obj := append(header, byts...)
	ast.WriteMemory(heapOffset, obj)

	return int32(heapOffset)
}

func opIdentity(inputs []ast.CXValue, outputs []ast.CXValue) {
	out1 := outputs[0].Arg
	var elt *ast.CXArgument
	if len(out1.Fields) > 0 {
		elt = out1.Fields[len(out1.Fields)-1]
	} else {
		elt = out1
	}

	//TODO: Delete
	if elt.DoesEscape {
		outputs[0].Set_i32(EscapeAnalysis(&inputs[0]))
	} else {
		switch elt.PassBy {
		case constants.PASSBY_VALUE:
			outputs[0].Set_bytes(inputs[0].Get_bytes())
		case constants.PASSBY_REFERENCE:
			outputs[0].Set_i32(int32(inputs[0].Offset))
		}
	}

	//inputs[0].Used = int8(inputs[0].Type)
	//outputs[0].Used = int8(outputs[0].Type)
}

func opGoto(inputs []ast.CXValue, outputs []ast.CXValue) {
	call := ast.PROGRAM.GetCurrentCall()
	expr := call.Operator.Expressions[call.Line]
	call.Line = call.Line + expr.ThenLines
}

func opJmp(inputs []ast.CXValue, outputs []ast.CXValue) {
	call := ast.PROGRAM.GetCurrentCall()
	expr := inputs[0].Expr

	if inputs[0].Get_bool() {
		call.Line = call.Line + expr.ThenLines
	} else {
		call.Line = call.Line + expr.ElseLines
	}
}

func opAbsJmp(inputs []ast.CXValue, outputs []ast.CXValue) {
	call := ast.PROGRAM.GetCurrentCall()
	expr := inputs[0].Expr

	if inputs[0].Get_bool() {
		call.Line = expr.ThenLines
	} else {
		call.Line = expr.ElseLines
	}
}

func opBreak(inputs []ast.CXValue, outputs []ast.CXValue) {
	call := ast.PROGRAM.GetCurrentCall()
	expr := call.Operator.Expressions[call.Line]
	call.Line = call.Line + expr.ThenLines
}

func opContinue(inputs []ast.CXValue, outputs []ast.CXValue) {
	call := ast.PROGRAM.GetCurrentCall()
	expr := call.Operator.Expressions[call.Line]
	call.Line = call.Line + expr.ThenLines
}

func opNop(inputs []ast.CXValue, outputs []ast.CXValue) {
	// No Operation
	// Do Nothing
}
