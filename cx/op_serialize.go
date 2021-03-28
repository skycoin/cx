package cxcore

import (
	"github.com/skycoin/cx/cx/ast"
	"github.com/skycoin/cx/cx/constants"
	"github.com/skycoin/cx/cx/helper"
)

func opSerialize(expr *ast.CXExpression, fp int) {
	inp1, out1 := expr.Inputs[0], expr.Outputs[0]
	out1Offset := ast.GetFinalOffset(fp, out1)

	_ = inp1

	var slcOff int
	byts := ast.SerializeCXProgram(ast.PROGRAM, true)
	for _, b := range byts {
		slcOff = ast.WriteToSlice(slcOff, []byte{b})
	}

	ast.WriteI32(out1Offset, int32(slcOff))
}

func opDeserialize(expr *ast.CXExpression, fp int) {
	inp := expr.Inputs[0]

	inpOffset := ast.GetFinalOffset(fp, inp)

	off := helper.Deserialize_i32(ast.PROGRAM.Memory[inpOffset : inpOffset+constants.TYPE_POINTER_SIZE])

	_l := ast.PROGRAM.Memory[off+constants.OBJECT_HEADER_SIZE : off+constants.OBJECT_HEADER_SIZE+constants.SLICE_HEADER_SIZE]
	l := helper.Deserialize_i32(_l[4:8])

	ast.Deserialize(ast.PROGRAM.Memory[off+constants.OBJECT_HEADER_SIZE+constants.SLICE_HEADER_SIZE : off+constants.OBJECT_HEADER_SIZE+constants.SLICE_HEADER_SIZE+l]) // BUG : should be l * elt.TotalSize ?
}
