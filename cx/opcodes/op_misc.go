package opcodes

import (
	"github.com/skycoin/cx/cx/ast"
	"github.com/skycoin/cx/cx/constants"
	"github.com/skycoin/cx/cx/types"
)

func opIdentity(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	out1 := outputs[0].Arg
	var elt *ast.CXArgument
	if len(out1.Fields) > 0 {
		elt = out1.Fields[len(out1.Fields)-1]
	} else {
		elt = out1
	}

	//TODO: Delete
	if elt.DoesEscape {
		outputs[0].Set_ptr(prgrm, types.AllocWrite_obj_data(prgrm, prgrm.Memory, inputs[0].Get_bytes(prgrm)))
	} else {
		switch elt.PassBy {
		case constants.PASSBY_VALUE:
			outputs[0].Set_bytes(prgrm, inputs[0].Get_bytes(prgrm))
		case constants.PASSBY_REFERENCE:
			outputs[0].Set_ptr(prgrm, inputs[0].Offset)
		}
	}
}

func opGoto(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	call := prgrm.GetCurrentCall()
	expr := call.Operator.Expressions[call.Line]

	cxAtomicOp, _, _, err := prgrm.GetOperation(&expr)
	if err != nil {
		panic(err)
	}

	call.Line = call.Line + cxAtomicOp.ThenLines
}

func opJmp(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	call := prgrm.GetCurrentCall()
	expr := inputs[0].Expr

	cxAtomicOp, _, _, err := prgrm.GetOperation(expr)
	if err != nil {
		panic(err)
	}

	if inputs[0].Get_bool(prgrm) {
		call.Line = call.Line + cxAtomicOp.ThenLines
	} else {
		call.Line = call.Line + cxAtomicOp.ElseLines
	}
}

func opAbsJmp(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	call := prgrm.GetCurrentCall()
	expr := inputs[0].Expr

	cxAtomicOp, _, _, err := prgrm.GetOperation(expr)
	if err != nil {
		panic(err)
	}

	if inputs[0].Get_bool(prgrm) {
		call.Line = cxAtomicOp.ThenLines
	} else {
		call.Line = cxAtomicOp.ElseLines
	}
}

func opBreak(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	call := prgrm.GetCurrentCall()
	expr := call.Operator.Expressions[call.Line]

	cxAtomicOp, _, _, err := prgrm.GetOperation(&expr)
	if err != nil {
		panic(err)
	}

	call.Line = call.Line + cxAtomicOp.ThenLines
}

func opContinue(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	call := prgrm.GetCurrentCall()
	expr := call.Operator.Expressions[call.Line]

	cxAtomicOp, _, _, err := prgrm.GetOperation(&expr)
	if err != nil {
		panic(err)
	}

	call.Line = call.Line + cxAtomicOp.ThenLines
}

func opNop(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	// No Operation
	// Do Nothing
}
