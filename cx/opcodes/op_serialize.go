package opcodes

import (
	"github.com/skycoin/cx/cx/ast"
	"github.com/skycoin/cx/cx/constants"
	"github.com/skycoin/cx/cx/types"
)

func opSerialize(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	var slcOff types.Pointer
	byts := ast.SerializeCXProgram(prgrm, true, false)
	for _, b := range byts {
		slcOff = ast.WriteToSlice(slcOff, []byte{b})
	}

	outputs[0].Set_ptr(slcOff)
}

func opDeserialize(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	off := inputs[0].Get_ptr()

	_l := prgrm.Memory[off+types.OBJECT_HEADER_SIZE : off+types.OBJECT_HEADER_SIZE+constants.SLICE_HEADER_SIZE]
	l := types.Read_ptr(_l, 4) // TODO:PTR remove hardcode 4

	ast.Deserialize(prgrm.Memory[off+types.OBJECT_HEADER_SIZE+constants.SLICE_HEADER_SIZE:off+types.OBJECT_HEADER_SIZE+constants.SLICE_HEADER_SIZE+l], false) // BUG : should be l * elt.TotalSize ?
}
