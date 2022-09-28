package opcodes

import (
	"github.com/skycoin/cx/cx/ast"
	"github.com/skycoin/cx/cx/constants"
	"github.com/skycoin/cx/cx/types"
)

//TODO: Rework
func opSliceLen(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	var inp0 *ast.CXArgument
	var sliceLen types.Pointer
	if inputs[0].TypeSignature.Type == ast.TYPE_CXARGUMENT_DEPRECATE {
		inp0 = prgrm.GetCXArgFromArray(ast.CXArgumentIndex(inputs[0].TypeSignature.Meta))

		elt := inp0.GetAssignmentElement(prgrm)

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
	} else if inputs[0].TypeSignature.Type == ast.TYPE_ARRAY_ATOMIC {
		arrDetails := prgrm.GetCXTypeSignatureArrayFromArray(inputs[0].TypeSignature.Meta)
		sliceLen = arrDetails.Lengths[len(arrDetails.Indexes)]
	} else if inputs[0].TypeSignature.Type == ast.TYPE_SLICE_ATOMIC {
		sliceOffset := types.Read_ptr(prgrm.Memory, inputs[0].Offset)
		if sliceOffset > 0 && sliceOffset.IsValid() {
			sliceLen = ast.GetSliceLen(prgrm, sliceOffset)
		} else if sliceOffset < 0 {
			panic(constants.CX_RUNTIME_ERROR)
		}
	} else {
		panic("type is not known\n\n")
	}

	outputs[0].Set_i32(prgrm, types.Cast_ptr_to_i32(sliceLen)) // TODO:PTR remove hardcode i32, should use ptr alias.
}

//TODO: Rework
func opSliceAppend(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	var inp0, out0 *ast.CXArgument
	var inp0Type, inp1Type, out0Type types.Code
	var input0IsSlice, outputIsSlice bool
	if inputs[0].TypeSignature.Type == ast.TYPE_CXARGUMENT_DEPRECATE {
		inp0 = prgrm.GetCXArgFromArray(ast.CXArgumentIndex(inputs[0].TypeSignature.Meta))
		inp0Type = inp0.Type
		if inp0.Type == types.POINTER {
			inp0Type = inp0.PointerTargetType
		}

		eltInp0 := inp0.GetAssignmentElement(prgrm)
		input0IsSlice = eltInp0.IsSlice
	} else if inputs[0].TypeSignature.Type == ast.TYPE_SLICE_ATOMIC {
		sliceDetails := prgrm.GetCXTypeSignatureArrayFromArray(inputs[0].TypeSignature.Meta)
		inp0Type = types.Code(sliceDetails.Type)

		input0IsSlice = true
	} else {
		panic("type is not type cx argument deprecate\n\n")
	}

	if inputs[1].TypeSignature.Type == ast.TYPE_CXARGUMENT_DEPRECATE {
		inp1 := prgrm.GetCXArgFromArray(ast.CXArgumentIndex(inputs[1].TypeSignature.Meta))
		inp1Type = inp1.Type
		if inp1.Type == types.POINTER {
			inp1Type = inp1.PointerTargetType
		}
	} else if inputs[1].TypeSignature.Type == ast.TYPE_ATOMIC || inputs[1].TypeSignature.Type == ast.TYPE_POINTER_ATOMIC {
		inp1Type = types.Code(inputs[1].TypeSignature.Meta)
	} else if inputs[1].TypeSignature.Type == ast.TYPE_SLICE_ATOMIC {
		sliceDetails := prgrm.GetCXTypeSignatureArrayFromArray(inputs[1].TypeSignature.Meta)

		inp1Type = types.Code(sliceDetails.Type)
	} else {
		panic("type is not known")
	}

	if outputs[0].TypeSignature.Type == ast.TYPE_CXARGUMENT_DEPRECATE {
		out0 = prgrm.GetCXArgFromArray(ast.CXArgumentIndex(outputs[0].TypeSignature.Meta))

		out0 := out0.GetAssignmentElement(prgrm)

		out0Type = out0.Type
		if out0.Type == types.POINTER {
			out0Type = out0.PointerTargetType
		}

		outputIsSlice = out0.IsSlice

	} else if outputs[0].TypeSignature.Type == ast.TYPE_SLICE_ATOMIC {
		sliceDetails := prgrm.GetCXTypeSignatureArrayFromArray(outputs[0].TypeSignature.Meta)
		out0Type = types.Code(sliceDetails.Type)
		outputIsSlice = true
	} else {
		panic("type is not type cx argument deprecate\n\n")
	}

	sliceInputs := inputs[1:]
	sliceInputsLen := types.Cast_int_to_ptr(len(sliceInputs))

	if inp0Type != inp1Type || inp0Type != out0Type || !input0IsSlice || !outputIsSlice {
		panic(constants.CX_RUNTIME_INVALID_ARGUMENT)
	}

	var inputSliceLen types.Pointer
	inputSliceOffset := types.Read_ptr(prgrm.Memory, inputs[0].Offset)
	if inputSliceOffset != 0 && inputSliceOffset.IsValid() {
		inputSliceLen = ast.GetSliceLen(prgrm, inputSliceOffset)
	}

	// Preparing slice in case more memory is needed for the new element.
	outputSliceOffset := ast.SliceAppendResize(prgrm, inputs[0].FramePointer, outputs[0].TypeSignature, inputs[0].TypeSignature, ast.GetDerefSizeSlice(prgrm, inputs[0].TypeSignature), sliceInputsLen)

	// We need to update the address of the output and input, as the final offsets
	// could be on the heap and they could have been moved by the GC.

	for i, input := range sliceInputs {
		var inp *ast.CXArgument
		var inpType types.Code
		if input.TypeSignature.Type == ast.TYPE_CXARGUMENT_DEPRECATE {
			inp = prgrm.GetCXArgFromArray(ast.CXArgumentIndex(input.TypeSignature.Meta))
			inpType = inp.Type
			if inp.Type == types.POINTER {
				inpType = inp.PointerTargetType
			}
		} else if input.TypeSignature.Type == ast.TYPE_ATOMIC || input.TypeSignature.Type == ast.TYPE_POINTER_ATOMIC {
			inpType = types.Code(input.TypeSignature.Meta)
		} else if input.TypeSignature.Type == ast.TYPE_SLICE_ATOMIC {
			sliceDetails := prgrm.GetCXTypeSignatureArrayFromArray(input.TypeSignature.Meta)
			inpType = types.Code(sliceDetails.Type)
		} else {
			panic("type is not known")
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
	var inpSize types.Pointer
	var inpType, outType types.Code
	var inpIsSlice, outIsSlice bool
	if inputs[0].TypeSignature.Type == ast.TYPE_CXARGUMENT_DEPRECATE {
		inp0 := prgrm.GetCXArgFromArray(ast.CXArgumentIndex(inputs[0].TypeSignature.Meta))
		eltInp0 := inp0.GetAssignmentElement(prgrm)
		inpSize = eltInp0.Size
	} else if inputs[0].TypeSignature.Type == ast.TYPE_SLICE_ATOMIC {
		sliceDetails := prgrm.GetCXTypeSignatureArrayFromArray(inputs[0].TypeSignature.Meta)
		inpSize = types.Code(sliceDetails.Type).Size()

		inpType = types.Code(sliceDetails.Type)
		inpIsSlice = true
	} else {
		panic("type is not type cx argument deprecate\n\n")
	}

	if outputs[0].TypeSignature.Type == ast.TYPE_CXARGUMENT_DEPRECATE {
		out0 := prgrm.GetCXArgFromArray(ast.CXArgumentIndex(outputs[0].TypeSignature.Meta))
		outType = out0.Type
		outIsSlice = out0.GetAssignmentElement(prgrm).IsSlice
	} else if outputs[0].TypeSignature.Type == ast.TYPE_SLICE_ATOMIC {
		sliceDetails := prgrm.GetCXTypeSignatureArrayFromArray(outputs[0].TypeSignature.Meta)

		outType = types.Code(sliceDetails.Type)
		outIsSlice = true
	} else {
		panic("type is not type cx argument deprecate\n\n")
	}

	fp := inputs[0].FramePointer

	if inpType != outType || !inpIsSlice || !outIsSlice {
		panic(constants.CX_RUNTIME_INVALID_ARGUMENT)
	}

	outputSliceOffset := ast.SliceResize(prgrm, fp, outputs[0].TypeSignature, inputs[0].TypeSignature, types.Cast_i32_to_ptr(inputs[1].Get_i32(prgrm)), inpSize)

	outputs[0].Set_ptr(prgrm, outputSliceOffset)
}

//TODO: Rework
func opSliceInsertElement(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	var inp0Type, inp2Type, out0Type types.Code
	var inpIsSlice, outIsSlice bool
	if inputs[0].TypeSignature.Type == ast.TYPE_CXARGUMENT_DEPRECATE {
		inp0 := prgrm.GetCXArgFromArray(ast.CXArgumentIndex(inputs[0].TypeSignature.Meta))

		inp0Type = inp0.Type
		if inp0.Type == types.POINTER {
			inp0Type = inp0.PointerTargetType
		}

		inpIsSlice = inp0.IsSlice
	} else if inputs[0].TypeSignature.Type == ast.TYPE_ATOMIC || inputs[0].TypeSignature.Type == ast.TYPE_POINTER_ATOMIC {
		inp0Type = types.Code(inputs[0].TypeSignature.Meta)

		inpIsSlice = false
	} else if inputs[0].TypeSignature.Type == ast.TYPE_SLICE_ATOMIC {
		sliceDetails := prgrm.GetCXTypeSignatureArrayFromArray(inputs[0].TypeSignature.Meta)

		inp0Type = types.Code(sliceDetails.Type)
		inpIsSlice = true
	} else {
		panic("type is not known")
	}

	if inputs[2].TypeSignature.Type == ast.TYPE_CXARGUMENT_DEPRECATE {
		inp2 := prgrm.GetCXArgFromArray(ast.CXArgumentIndex(inputs[2].TypeSignature.Meta))

		inp2Type = inp2.Type
		if inp2.Type == types.POINTER {
			inp2Type = inp2.PointerTargetType
		}
	} else if inputs[2].TypeSignature.Type == ast.TYPE_ATOMIC || inputs[2].TypeSignature.Type == ast.TYPE_POINTER_ATOMIC {
		inp2Type = types.Code(inputs[2].TypeSignature.Meta)
	} else {
		panic("type is not known")
	}

	if outputs[0].TypeSignature.Type == ast.TYPE_CXARGUMENT_DEPRECATE {
		out0 := prgrm.GetCXArgFromArray(ast.CXArgumentIndex(outputs[0].TypeSignature.Meta))

		out0Type = out0.Type
		if out0.Type == types.POINTER {
			out0Type = out0.PointerTargetType
		}

		outIsSlice = out0.IsSlice
	} else if outputs[0].TypeSignature.Type == ast.TYPE_ATOMIC || outputs[0].TypeSignature.Type == ast.TYPE_POINTER_ATOMIC {
		out0Type = types.Code(outputs[0].TypeSignature.Meta)

		outIsSlice = false
	} else if outputs[0].TypeSignature.Type == ast.TYPE_SLICE_ATOMIC {
		sliceDetails := prgrm.GetCXTypeSignatureArrayFromArray(outputs[0].TypeSignature.Meta)

		out0Type = types.Code(sliceDetails.Type)
		outIsSlice = true
	} else {
		panic("type is not known")
	}

	fp := inputs[0].FramePointer

	if inp0Type != inp2Type || inp0Type != out0Type || !inpIsSlice || !outIsSlice {
		panic(constants.CX_RUNTIME_INVALID_ARGUMENT)
	}

	index := types.Cast_i32_to_ptr(inputs[1].Get_i32(prgrm))

	var outputSliceOffset types.Pointer
	if inp2Type == types.STR || inp2Type == types.AFF {
		var obj [types.POINTER_SIZE]byte
		types.Write_ptr(obj[:], 0, types.Read_ptr(prgrm.Memory, inputs[2].Offset))
		outputSliceOffset = ast.SliceInsert(prgrm, fp, outputs[0].TypeSignature, inputs[0].TypeSignature, index, obj[:])
	} else {
		obj := inputs[2].Get_bytes(prgrm)
		outputSliceOffset = ast.SliceInsert(prgrm, fp, outputs[0].TypeSignature, inputs[0].TypeSignature, index, obj)
	}

	outputs[0].Set_ptr(prgrm, outputSliceOffset)
}

//TODO: Rework
func opSliceRemoveElement(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	var inpType, outType types.Code
	var inpIsSlice, outIsSlice bool
	var inpSize types.Pointer

	if inputs[0].TypeSignature.Type == ast.TYPE_CXARGUMENT_DEPRECATE {
		inp0 := prgrm.GetCXArgFromArray(ast.CXArgumentIndex(inputs[0].TypeSignature.Meta))

		inpType = inp0.Type
		inpIsSlice = inp0.GetAssignmentElement(prgrm).IsSlice
		inpSize = inp0.GetAssignmentElement(prgrm).Size
	} else if inputs[0].TypeSignature.Type == ast.TYPE_SLICE_ATOMIC {
		sliceDetails := prgrm.GetCXTypeSignatureArrayFromArray(inputs[0].TypeSignature.Meta)

		inpType = types.Code(sliceDetails.Type)
		inpIsSlice = true
		inpSize = inpType.Size()
	} else {
		panic("type is not type cx argument deprecate\n\n")
	}

	if outputs[0].TypeSignature.Type == ast.TYPE_CXARGUMENT_DEPRECATE {
		out0 := prgrm.GetCXArgFromArray(ast.CXArgumentIndex(outputs[0].TypeSignature.Meta))

		outType = out0.Type
		outIsSlice = out0.GetAssignmentElement(prgrm).IsSlice
	} else if outputs[0].TypeSignature.Type == ast.TYPE_SLICE_ATOMIC {
		sliceDetails := prgrm.GetCXTypeSignatureArrayFromArray(outputs[0].TypeSignature.Meta)

		outType = types.Code(sliceDetails.Type)
		outIsSlice = true
	} else {
		panic("type is not type cx argument deprecate\n\n")
	}

	fp := inputs[0].FramePointer

	if inpType != outType || !inpIsSlice || !outIsSlice {
		panic(constants.CX_RUNTIME_INVALID_ARGUMENT)
	}

	outputSliceOffset := ast.SliceRemove(prgrm, fp, outputs[0].TypeSignature, inputs[0].TypeSignature, types.Cast_i32_to_ptr(inputs[1].Get_i32(prgrm)), inpSize)

	outputs[0].Set_ptr(prgrm, outputSliceOffset)
}

func opSliceCopy(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	var dstInput *ast.CXArgument
	if inputs[0].TypeSignature.Type == ast.TYPE_CXARGUMENT_DEPRECATE {
		dstInput = prgrm.GetCXArgFromArray(ast.CXArgumentIndex(inputs[0].TypeSignature.Meta))
	} else {
		panic("type is not type cx argument deprecate\n\n")
	}

	var srcInput *ast.CXArgument
	if inputs[1].TypeSignature.Type == ast.TYPE_CXARGUMENT_DEPRECATE {
		srcInput = prgrm.GetCXArgFromArray(ast.CXArgumentIndex(inputs[1].TypeSignature.Meta))
	} else {
		panic("type is not type cx argument deprecate\n\n")
	}

	fp := inputs[0].FramePointer

	dstOffset := ast.GetSliceOffset(prgrm, fp, inputs[0].TypeSignature)
	srcOffset := ast.GetSliceOffset(prgrm, fp, inputs[1].TypeSignature)

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
