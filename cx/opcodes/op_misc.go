package opcodes

import (
	"github.com/skycoin/cx/cx/ast"
	"github.com/skycoin/cx/cx/constants"
	"github.com/skycoin/cx/cx/types"
)

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
		outputs[0].Set_ptr(types.AllocWrite_obj_data(ast.PROGRAM.Memory, inputs[0].Get_bytes()))
	} else {
		switch elt.PassBy {
		case constants.PASSBY_VALUE:
			outputs[0].Set_bytes(inputs[0].Get_bytes())
		case constants.PASSBY_REFERENCE:
			outputs[0].Set_ptr(inputs[0].Offset)
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
