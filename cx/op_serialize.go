package cxcore

import (
	"github.com/skycoin/cx/cx/ast"
	"github.com/skycoin/cx/cx/constants"
)

func opSerialize(expr *ast.CXExpression, fp int) {
	inp1, out1 := expr.Inputs[0], expr.Outputs[0]
	out1Offset := GetFinalOffset(fp, out1)

	_ = inp1

	var slcOff int
	byts := SerializeCXProgram(ast.PROGRAM, true)
	for _, b := range byts {
		slcOff = WriteToSlice(slcOff, []byte{b})
	}

	WriteI32(out1Offset, int32(slcOff))
}

func opDeserialize(expr *ast.CXExpression, fp int) {
	inp := expr.Inputs[0]

	inpOffset := GetFinalOffset(fp, inp)

	off := Deserialize_i32(ast.PROGRAM.Memory[inpOffset : inpOffset+constants.TYPE_POINTER_SIZE])

	_l := ast.PROGRAM.Memory[off+constants.OBJECT_HEADER_SIZE : off+constants.OBJECT_HEADER_SIZE+constants.SLICE_HEADER_SIZE]
	l := Deserialize_i32(_l[4:8])

	Deserialize(ast.PROGRAM.Memory[off+constants.OBJECT_HEADER_SIZE+constants.SLICE_HEADER_SIZE : off+constants.OBJECT_HEADER_SIZE+constants.SLICE_HEADER_SIZE+l]) // BUG : should be l * elt.TotalSize ?
}
