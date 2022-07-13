package opcodes

import (
	"github.com/skycoin/cx/cx/ast"
	"github.com/skycoin/cx/cx/constants"
	"github.com/skycoin/cx/cx/types"
)

func opIdentity(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	var out1 *ast.CXArgument
	if outputs[0].TypeSignature.Type == ast.TYPE_CXARGUMENT_DEPRECATE {
		out1 = prgrm.GetCXArgFromArray(ast.CXArgumentIndex(outputs[0].TypeSignature.Meta))

		var elt *ast.CXArgument
		if len(out1.Fields) > 0 {
			elt = prgrm.GetCXArgFromArray(out1.Fields[len(out1.Fields)-1])
		} else {
			elt = out1
		}

		if elt.IsPointer() && elt.PointerTargetType != types.STR && elt.Type != types.AFF || elt.PassBy == constants.PASSBY_REFERENCE {
			// outputs[0].Set_ptr(prgrm, types.AllocWrite_obj_data(prgrm, prgrm.Memory, inputs[0].Get_bytes(prgrm)))
			outputs[0].Set_ptr(prgrm, inputs[0].Offset)
		} else {
			// Pass by value
			outputs[0].Set_bytes(prgrm, inputs[0].Get_bytes(prgrm))
		}
	} else {
		// panic("type is not type cx argument deprecate\n\n")
		// TODO: type atomic for now so automatically
		// pass by value
		outputs[0].Set_bytes(prgrm, inputs[0].Get_bytes(prgrm))
	}

}

func opGoto(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	call := prgrm.GetCurrentCall()
	expr := call.Operator.Expressions[call.Line]

	cxAtomicOp, err := prgrm.GetCXAtomicOp(expr.Index)
	if err != nil {
		panic(err)
	}

	call.Line = call.Line + cxAtomicOp.ThenLines
}

func opJmp(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	call := prgrm.GetCurrentCall()
	expr := inputs[0].Expr

	cxAtomicOp, err := prgrm.GetCXAtomicOp(expr.Index)
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

	cxAtomicOp, err := prgrm.GetCXAtomicOp(expr.Index)
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

	cxAtomicOp, err := prgrm.GetCXAtomicOp(expr.Index)
	if err != nil {
		panic(err)
	}

	call.Line = call.Line + cxAtomicOp.ThenLines
}

func opContinue(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	call := prgrm.GetCurrentCall()
	expr := call.Operator.Expressions[call.Line]

	cxAtomicOp, err := prgrm.GetCXAtomicOp(expr.Index)
	if err != nil {
		panic(err)
	}

	call.Line = call.Line + cxAtomicOp.ThenLines
}

func opNop(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	// No Operation
	// Do Nothing
}
