package opcodes

import (
	"github.com/skycoin/cx/cx/ast"
	"github.com/skycoin/cx/cx/constants"
	"github.com/skycoin/cx/cx/helper"
)

func opSerialize(inputs []ast.CXValue, outputs []ast.CXValue) {
	var slcOff int
	byts := ast.SerializeCXProgram(ast.PROGRAM, true)
	for _, b := range byts {
		slcOff = ast.WriteToSlice(slcOff, []byte{b})
	}

    outputs[0].Set_i32(int32(slcOff))
}

func opDeserialize(inputs []ast.CXValue, outputs []ast.CXValue) {
	off := inputs[0].Get_i32()

	_l := ast.PROGRAM.Memory[off+constants.OBJECT_HEADER_SIZE : off+constants.OBJECT_HEADER_SIZE+constants.SLICE_HEADER_SIZE]
	l := helper.Deserialize_i32(_l[4:8])

	ast.Deserialize(ast.PROGRAM.Memory[off+constants.OBJECT_HEADER_SIZE+constants.SLICE_HEADER_SIZE : off+constants.OBJECT_HEADER_SIZE+constants.SLICE_HEADER_SIZE+l]) // BUG : should be l * elt.TotalSize ?
}
