package opcodes

import (
	"github.com/skycoin/cx/cx/ast"
	"github.com/skycoin/cx/cx/constants"
	"github.com/skycoin/cx/cx/types"
)

//TODO: Rework
func opSliceLen(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	elt := inputs[0].Arg.GetAssignmentElement(prgrm)

	var sliceLen types.Pointer
	if elt.IsSlice || elt.Type == types.AFF { //TODO: FIX
		sliceOffset := types.Read_ptr(prgrm.Memory, inputs[0].Offset)
		if sliceOffset > 0 && sliceOffset.IsValid() {
			sliceLen = ast.GetSliceLen(prgrm, sliceOffset)
		} else if sliceOffset < 0 {
			panic(constants.CX_RUNTIME_ERROR)
		}
		// TODO: Had to add elt.Lengths to avoid doing this for arrays, but not entirely sure why
	} else if (elt.PointerTargetType == types.STR || elt.Type == types.STR) && elt.Lengths == nil {
		sliceLen = types.Read_str_size(prgrm.Memory, inputs[0].Offset)
	} else {
		sliceLen = elt.Lengths[len(elt.Indexes)]
	}

	outputs[0].Set_i32(prgrm, types.Cast_ptr_to_i32(sliceLen)) // TODO:PTR remove hardcode i32, should use ptr alias.
}

//TODO: Rework
func opSliceAppend(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	inp0, inp1, out0 := inputs[0].Arg, inputs[1].Arg, outputs[0].Arg
	sliceInputs := inputs[1:]
	sliceInputsLen := types.Cast_int_to_ptr(len(sliceInputs))

	eltInp0 := inp0.GetAssignmentElement(prgrm)
	eltOut0 := out0.GetAssignmentElement(prgrm)

	inp0Type := inp0.Type
	inp1Type := inp1.Type
	out0Type := out0.Type
	if inp0.Type == types.POINTER {
		inp0Type = inp0.PointerTargetType
	}

	if inp1.Type == types.POINTER {
		inp1Type = inp1.PointerTargetType
	}

	if out0.Type == types.POINTER {
		out0Type = out0.PointerTargetType
	}

	if inp0Type != inp1Type || inp0Type != out0Type || !eltInp0.IsSlice || !eltOut0.IsSlice {
		panic(constants.CX_RUNTIME_INVALID_ARGUMENT)
	}

	var inputSliceLen types.Pointer
	inputSliceOffset := types.Read_ptr(prgrm.Memory, inputs[0].Offset)
	if inputSliceOffset != 0 && inputSliceOffset.IsValid() {
		inputSliceLen = ast.GetSliceLen(prgrm, inputSliceOffset)
	}

	// Preparing slice in case more memory is needed for the new element.
	outputSliceOffset := ast.SliceAppendResize(prgrm, inputs[0].FramePointer, out0, inp0, ast.GetDerefSizeSlice(prgrm, eltInp0), sliceInputsLen)

	// We need to update the address of the output and input, as the final offsets
	// could be on the heap and they could have been moved by the GC.

	for i, input := range sliceInputs {
		inp := input.Arg
		inpType := inp.Type
		if inp.Type == types.POINTER {
			inpType = inp.PointerTargetType
		}

		if inp0Type != inpType {
			panic(constants.CX_RUNTIME_INVALID_ARGUMENT)
		}
		if inpType == types.STR || inpType == types.AFF {
			var obj [types.POINTER_SIZE]byte
			types.Write_ptr(obj[:], 0, types.Read_ptr(prgrm.Memory, input.Offset))
			ast.SliceAppendWrite(prgrm, outputSliceOffset, obj[:], inputSliceLen+types.Cast_int_to_ptr(i))
		} else {
			obj := inputs[1].Get_bytes(prgrm)
			ast.SliceAppendWrite(prgrm, outputSliceOffset, obj, inputSliceLen+types.Cast_int_to_ptr(i))
		}
	}
	outputs[0].Set_ptr(prgrm, outputSliceOffset)
}

//TODO: Rework
func opSliceResize(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	inp0, out0 := inputs[0].Arg, outputs[0].Arg
	fp := inputs[0].FramePointer

	if inp0.Type != out0.Type || !inp0.GetAssignmentElement(prgrm).IsSlice || !out0.GetAssignmentElement(prgrm).IsSlice {
		panic(constants.CX_RUNTIME_INVALID_ARGUMENT)
	}

	eltInp0 := inp0.GetAssignmentElement(prgrm)

	outputSliceOffset := ast.SliceResize(prgrm, fp, out0, inp0, types.Cast_i32_to_ptr(inputs[1].Get_i32(prgrm)), eltInp0.Size)

	outputs[0].Set_ptr(prgrm, outputSliceOffset)
}

//TODO: Rework
func opSliceInsertElement(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	inp0, inp2, out0 := inputs[0].Arg, inputs[2].Arg, outputs[0].Arg
	fp := inputs[0].FramePointer

	inp0Type := inp0.Type
	if inp0.Type == types.POINTER {
		inp0Type = inp0.PointerTargetType
	}

	inp2Type := inp2.Type
	if inp2.Type == types.POINTER {
		inp2Type = inp2.PointerTargetType
	}

	out0Type := out0.Type
	if out0.Type == types.POINTER {
		out0Type = out0.PointerTargetType
	}

	if inp0Type != inp2Type || inp0Type != out0Type || !inp0.GetAssignmentElement(prgrm).IsSlice || !out0.GetAssignmentElement(prgrm).IsSlice {
		panic(constants.CX_RUNTIME_INVALID_ARGUMENT)
	}

	index := types.Cast_i32_to_ptr(inputs[1].Get_i32(prgrm))

	var outputSliceOffset types.Pointer
	if inp2Type == types.STR || inp2Type == types.AFF {
		var obj [types.POINTER_SIZE]byte
		types.Write_ptr(obj[:], 0, types.Read_ptr(prgrm.Memory, inputs[2].Offset))
		outputSliceOffset = ast.SliceInsert(prgrm, fp, out0, inp0, index, obj[:])
	} else {
		obj := inputs[2].Get_bytes(prgrm)
		outputSliceOffset = ast.SliceInsert(prgrm, fp, out0, inp0, index, obj)
	}

	outputs[0].Set_ptr(prgrm, outputSliceOffset)
}

//TODO: Rework
func opSliceRemoveElement(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	inp0, out0 := inputs[0].Arg, outputs[0].Arg
	fp := inputs[0].FramePointer

	if inp0.Type != out0.Type || !inp0.GetAssignmentElement(prgrm).IsSlice || !out0.GetAssignmentElement(prgrm).IsSlice {
		panic(constants.CX_RUNTIME_INVALID_ARGUMENT)
	}

	outputSliceOffset := ast.SliceRemove(prgrm, fp, out0, inp0, types.Cast_i32_to_ptr(inputs[1].Get_i32(prgrm)), inp0.GetAssignmentElement(prgrm).Size)

	outputs[0].Set_ptr(prgrm, outputSliceOffset)
}

func opSliceCopy(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	dstInput := inputs[0].Arg
	srcInput := inputs[1].Arg
	fp := inputs[0].FramePointer

	dstOffset := ast.GetSliceOffset(prgrm, fp, dstInput)
	srcOffset := ast.GetSliceOffset(prgrm, fp, srcInput)

	dstElem := dstInput.GetAssignmentElement(prgrm)
	srcElem := srcInput.GetAssignmentElement(prgrm)

	if dstInput.Type != srcInput.Type || !dstElem.IsSlice || !srcElem.IsSlice || dstElem.Size != srcElem.Size {
		panic(constants.CX_RUNTIME_INVALID_ARGUMENT)
	}

	dst := ast.GetSliceData(prgrm, dstOffset, dstElem.Size)
	src := ast.GetSliceData(prgrm, srcOffset, srcElem.Size)

	var count types.Pointer
	if dstInput.Type == srcInput.Type && dstOffset >= 0 && srcOffset >= 0 {
		count = types.Cast_int_to_ptr(copy(dst, src))
		if count%dstElem.Size != 0 {
			panic(constants.CX_RUNTIME_ERROR)
		}
	} else {
		panic(constants.CX_RUNTIME_INVALID_ARGUMENT)
	}

	outputs[0].Set_i32(prgrm, types.Cast_ptr_to_i32(count/dstElem.Size)) // TODO:PTR use ptr instead of i32
}
