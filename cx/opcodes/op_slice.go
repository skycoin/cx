package opcodes

import (
	"github.com/skycoin/cx/cx/ast"
	"github.com/skycoin/cx/cx/constants"
	"github.com/skycoin/cx/cx/types"
)

//TODO: Rework
func opSliceLen(inputs []ast.CXValue, outputs []ast.CXValue) {
	elt := inputs[0].Arg.GetAssignmentElement()

	var sliceLen types.Pointer
	if elt.IsSlice || elt.Type == types.AFF { //TODO: FIX
		sliceOffset := types.Read_ptr(ast.PROGRAM.Memory, inputs[0].Offset)
		if sliceOffset > 0 && sliceOffset.IsValid() {
			sliceLen = ast.GetSliceLen(sliceOffset)
		} else if sliceOffset < 0 {
			panic(constants.CX_RUNTIME_ERROR)
		}
		// TODO: Had to add elt.Lengths to avoid doing this for arrays, but not entirely sure why
	} else if elt.Type == types.STR && elt.Lengths == nil {
		sliceLen = types.Read_str_size(ast.PROGRAM.Memory, inputs[0].Offset)
	} else {
		sliceLen = elt.Lengths[len(elt.Indexes)]
	}

	outputs[0].Set_i32(types.Cast_ptr_to_i32(sliceLen)) // TODO:PTR remove hardcode i32, should use ptr alias.
}

//TODO: Rework
func opSliceAppend(inputs []ast.CXValue, outputs []ast.CXValue) {
	inp0, inp1, out0 := inputs[0].Arg, inputs[1].Arg, outputs[0].Arg
	sliceInputs := inputs[1:]
	sliceInputsLen := types.Cast_int_to_ptr(len(sliceInputs))

	eltInp0 := inp0.GetAssignmentElement()
	eltOut0 := out0.GetAssignmentElement()

	if inp0.Type != inp1.Type || inp0.Type != out0.Type || !eltInp0.IsSlice || !eltOut0.IsSlice {
		panic(constants.CX_RUNTIME_INVALID_ARGUMENT)
	}

	var inputSliceLen types.Pointer
	inputSliceOffset := types.Read_ptr(ast.PROGRAM.Memory, inputs[0].Offset)
	if inputSliceOffset != 0 && inputSliceOffset.IsValid() {
		inputSliceLen = ast.GetSliceLen(inputSliceOffset)
	}

	// Preparing slice in case more memory is needed for the new element.
	outputSliceOffset := ast.SliceAppendResize(inputs[0].FramePointer, out0, inp0, ast.GetDerefSizeSlice(eltInp0), sliceInputsLen)

	// We need to update the address of the output and input, as the final offsets
	// could be on the heap and they could have been moved by the GC.

	for i, input := range sliceInputs {
		inp := input.Arg
		if inp0.Type != inp.Type {
			panic(constants.CX_RUNTIME_INVALID_ARGUMENT)
		}
		if inp.Type == types.STR || inp.Type == types.AFF {
			var obj [types.POINTER_SIZE]byte
			types.Write_ptr(obj[:], 0, types.Read_ptr(ast.PROGRAM.Memory, input.Offset))
			ast.SliceAppendWrite(outputSliceOffset, obj[:], inputSliceLen+types.Cast_int_to_ptr(i))
		} else {
			obj := inputs[1].Get_bytes()
			ast.SliceAppendWrite(outputSliceOffset, obj, inputSliceLen+types.Cast_int_to_ptr(i))
		}
	}
	outputs[0].Set_ptr(outputSliceOffset)
}

//TODO: Rework
func opSliceResize(inputs []ast.CXValue, outputs []ast.CXValue) {
	inp0, out0 := inputs[0].Arg, outputs[0].Arg
	fp := inputs[0].FramePointer

	if inp0.Type != out0.Type || !inp0.GetAssignmentElement().IsSlice || !out0.GetAssignmentElement().IsSlice {
		panic(constants.CX_RUNTIME_INVALID_ARGUMENT)
	}

	eltInp0 := inp0.GetAssignmentElement()

	outputSliceOffset := ast.SliceResize(fp, out0, inp0, types.Cast_i32_to_ptr(inputs[1].Get_i32()), eltInp0.Size)

	outputs[0].Set_ptr(outputSliceOffset)
}

//TODO: Rework
func opSliceInsertElement(inputs []ast.CXValue, outputs []ast.CXValue) {
	inp0, inp2, out0 := inputs[0].Arg, inputs[2].Arg, outputs[0].Arg
	fp := inputs[0].FramePointer

	if inp0.Type != inp2.Type || inp0.Type != out0.Type || !inp0.GetAssignmentElement().IsSlice || !out0.GetAssignmentElement().IsSlice {
		panic(constants.CX_RUNTIME_INVALID_ARGUMENT)
	}

	index := types.Cast_i32_to_ptr(inputs[1].Get_i32())
	var outputSliceOffset types.Pointer
	if inp2.Type == types.STR || inp2.Type == types.AFF {
		var obj [types.POINTER_SIZE]byte
		types.Write_ptr(obj[:], 0, types.Read_ptr(ast.PROGRAM.Memory, inputs[2].Offset))
		outputSliceOffset = ast.SliceInsert(fp, out0, inp0, index, obj[:])
	} else {
		obj := inputs[2].Get_bytes()
		outputSliceOffset = ast.SliceInsert(fp, out0, inp0, index, obj)
	}

	outputs[0].Set_ptr(outputSliceOffset)
}

//TODO: Rework
func opSliceRemoveElement(inputs []ast.CXValue, outputs []ast.CXValue) {
	inp0, out0 := inputs[0].Arg, outputs[0].Arg
	fp := inputs[0].FramePointer

	if inp0.Type != out0.Type || !inp0.GetAssignmentElement().IsSlice || !out0.GetAssignmentElement().IsSlice {
		panic(constants.CX_RUNTIME_INVALID_ARGUMENT)
	}

	outputSliceOffset := ast.SliceRemove(fp, out0, inp0, types.Cast_i32_to_ptr(inputs[1].Get_i32()), inp0.GetAssignmentElement().Size)

	outputs[0].Set_ptr(outputSliceOffset)
}

func opSliceCopy(inputs []ast.CXValue, outputs []ast.CXValue) {
	dstInput := inputs[0].Arg
	srcInput := inputs[1].Arg
	fp := inputs[0].FramePointer

	dstOffset := ast.GetSliceOffset(fp, dstInput)
	srcOffset := ast.GetSliceOffset(fp, srcInput)

	dstElem := dstInput.GetAssignmentElement()
	srcElem := srcInput.GetAssignmentElement()

	if dstInput.Type != srcInput.Type || !dstElem.IsSlice || !srcElem.IsSlice || dstElem.Size != srcElem.Size {
		panic(constants.CX_RUNTIME_INVALID_ARGUMENT)
	}

	dst := ast.GetSliceData(dstOffset, dstElem.Size)
	src := ast.GetSliceData(srcOffset, srcElem.Size)

	var count types.Pointer
	if dstInput.Type == srcInput.Type && dstOffset >= 0 && srcOffset >= 0 {
		count = types.Cast_int_to_ptr(copy(dst, src))
		if count%dstElem.Size != 0 {
			panic(constants.CX_RUNTIME_ERROR)
		}
	} else {
		panic(constants.CX_RUNTIME_INVALID_ARGUMENT)
	}

	outputs[0].Set_i32(types.Cast_ptr_to_i32(count / dstElem.Size)) // TODO:PTR use ptr instead of i32
}
