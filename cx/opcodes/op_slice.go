package opcodes

import (
	"github.com/skycoin/cx/cx/ast"
	"github.com/skycoin/cx/cx/constants"
	"github.com/skycoin/cx/cx/helper"

)

// getStrOffset ...
func getStrOffset(offset int, name string) int32 {
	if name != "" {
		// then it's not a literal
		return helper.Deserialize_i32(ast.PROGRAM.Memory[offset : offset+constants.TYPE_POINTER_SIZE])
	}
	return int32(offset) // TODO: Remove cast.
}

//TODO: Rename opSliceLen
//TODO: Rework
func opLen(inputs []ast.CXValue, outputs []ast.CXValue) {
	elt := ast.GetAssignmentElement(inputs[0].Arg)

	var sliceLen int32
	if elt.IsSlice || elt.Type == constants.TYPE_AFF { //TODO: FIX
		sliceOffset := ast.GetPointerOffset(int32(inputs[0].Offset))
		if sliceOffset > 0 {
			sliceLen = helper.Deserialize_i32(ast.GetSliceHeader(sliceOffset)[4:8])
		} else if sliceOffset < 0 {
			panic(constants.CX_RUNTIME_ERROR)
		}

		// TODO: Had to add elt.Lengths to avoid doing this for arrays, but not entirely sure why
	} else if elt.Type == constants.TYPE_STR && elt.Lengths == nil {
		var strOffset = getStrOffset(inputs[0].Offset, inputs[0].Arg.Name)
		// Checking if the string lives on the heap.
		if int(strOffset) > ast.PROGRAM.HeapStartsAt { // TODO: Remove cast.
			// Then it's on the heap and we need to consider
			// the object's header.
			strOffset += constants.OBJECT_HEADER_SIZE
		}

		sliceLen = helper.Deserialize_i32(ast.PROGRAM.Memory[strOffset:strOffset+constants.STR_HEADER_SIZE])
	} else {
		sliceLen = int32(elt.Lengths[len(elt.Indexes)])
	}

    //inputs[0].Used = int8(inputs[0].Type) // TODO: Remove hacked type check
	outputs[0].Set_i32(sliceLen)
}

//TODO: Rename OpSliceAppend
//TODO: Rework
func opAppend(inputs []ast.CXValue, outputs []ast.CXValue) {
	inp0, inp1, out0 := &inputs[0], inputs[1].Arg, &outputs[0]
	eltInp0 := ast.GetAssignmentElement(inp0.Arg)
	eltOut0 := ast.GetAssignmentElement(out0.Arg)

    if inp0.Arg.Type != inp1.Type || inp0.Arg.Type != out0.Arg.Type || !eltInp0.IsSlice || !eltOut0.IsSlice {
		panic(constants.CX_RUNTIME_INVALID_ARGUMENT)
	}

	var inputSliceLen int32
	inputSliceOffset := ast.GetPointerOffset(int32(inputs[0].Offset))
	if inputSliceOffset != 0 {
		inputSliceLen = ast.GetSliceLen(int32(inputSliceOffset))
	}

	// Preparing slice in case more memory is needed for the new element.
	outputSliceOffset := ast.SliceAppendResize(out0, inp0, inp1.Size)

	// We need to update the address of the output and input, as the final offsets
	// could be on the heap and they could have been moved by the GC.

	if inp1.Type == constants.TYPE_STR || inp1.Type == constants.TYPE_AFF {
		var obj [4]byte
		ast.WriteMemI32(obj[:], 0, int32(getStrOffset(inputs[1].Offset, inp1.Name)))
		ast.SliceAppendWrite(outputSliceOffset, obj[:], inputSliceLen)
	} else {
		obj := inputs[1].Get_bytes()
		ast.SliceAppendWrite(outputSliceOffset, obj, inputSliceLen)
	}

	//inputs[0].Used = int8(inputs[0].Type) // TODO: Remove hacked type check
	//inputs[1].Used = int8(inputs[1].Type) // TODO: Remove hacked type check
	outputs[0].SetSlice(outputSliceOffset)
}

//TODO: Rename opSliceResize
//TODO: Rework
func opResize(inputs []ast.CXValue, outputs []ast.CXValue) {
	inp0, out0 := &inputs[0], &outputs[0]
	eltInp0 := ast.GetAssignmentElement(inp0.Arg)
	eltOut0 := ast.GetAssignmentElement(out0.Arg)

    if inp0.Arg.Type != out0.Arg.Type || !eltInp0.IsSlice || !eltOut0.IsSlice {
		panic(constants.CX_RUNTIME_INVALID_ARGUMENT)
	}

	outputSliceOffset := int32(ast.SliceResize(out0, inp0, inputs[1].Get_i32(), eltInp0.TotalSize))

	//inputs[0].Used = int8(inputs[0].Type) // TODO: Remove hacked type check
    outputs[0].SetSlice(outputSliceOffset)
}

//TODO: Rename opSliceInsertElement
//TODO: Rework
func opInsert(inputs []ast.CXValue, outputs []ast.CXValue) {
	inp0, inp2, out0 := &inputs[0], inputs[2].Arg, &outputs[0]
	eltInp0 := ast.GetAssignmentElement(inp0.Arg)
	eltOut0 := ast.GetAssignmentElement(out0.Arg)

	if inp0.Type != inp2.Type || inp0.Type != out0.Type || !eltInp0.IsSlice || !eltOut0.IsSlice {
		panic(constants.CX_RUNTIME_INVALID_ARGUMENT)
	}

	index := inputs[1].Get_i32()
	var outputSliceOffset int32
	if inp2.Type == constants.TYPE_STR || inp2.Type == constants.TYPE_AFF {
		var obj [4]byte
		ast.WriteMemI32(obj[:], 0, int32(getStrOffset(inputs[2].Offset, inp2.Name)))
		outputSliceOffset = int32(ast.SliceInsert(out0, inp0, index, obj[:]))
	} else {
		obj := inputs[2].Get_bytes()
		outputSliceOffset = int32(ast.SliceInsert(out0, inp0, index, obj))
	}

	//inputs[0].Used = int8(inputs[0].Type) // TODO: Remove hacked type check
	//inputs[2].Used = int8(inputs[2].Type) // TODO: Remove hacked type check
	outputs[0].SetSlice(outputSliceOffset)
}

//TODO: Rename opSliceRemoveElement
//TODO: Rework
func opRemove(inputs []ast.CXValue, outputs []ast.CXValue) {
	inp0, out0 := &inputs[0], &outputs[0]
	eltInp0 := ast.GetAssignmentElement(inp0.Arg)
	eltOut0 := ast.GetAssignmentElement(out0.Arg)

	if inp0.Type != out0.Type || !eltInp0.IsSlice || !eltOut0.IsSlice {
		panic(constants.CX_RUNTIME_INVALID_ARGUMENT)
	}

	outputSliceOffset := int32(ast.SliceRemove(out0, inp0, inputs[1].Get_i32(), int32(eltInp0.TotalSize)))

	//inputs[0].Used = int8(inputs[0].Type) // TODO: Remove hacked type check
    outputs[0].SetSlice(outputSliceOffset)
}

//TODO: Rename opSliceCopy
func opCopy(inputs []ast.CXValue, outputs []ast.CXValue) {
	dstInput := &inputs[0]
	srcInput := &inputs[1]

    dstOffset := ast.GetSliceOffset(dstInput)
	srcOffset := ast.GetSliceOffset(srcInput)

	dstElem := ast.GetAssignmentElement(dstInput.Arg)
	srcElem := ast.GetAssignmentElement(srcInput.Arg)

	if dstInput.Type != srcInput.Type || !dstElem.IsSlice || !srcElem.IsSlice || dstElem.TotalSize != srcElem.TotalSize {
		panic(constants.CX_RUNTIME_INVALID_ARGUMENT)
	}

	var count int
	if dstInput.Type == srcInput.Type && dstOffset >= 0 && srcOffset >= 0 {
		count = copy(ast.GetSliceData(dstOffset, dstElem.TotalSize), ast.GetSliceData(srcOffset, srcElem.TotalSize))
		if count%dstElem.TotalSize != 0 {
			panic(constants.CX_RUNTIME_ERROR)
		}
	} else {
		panic(constants.CX_RUNTIME_INVALID_ARGUMENT)
	}
	
    //inputs[0].Used = int8(inputs[0].Type) // TODO: Remove hacked type check
	//inputs[1].Used = int8(inputs[1].Type) // TODO: Remove hacked type check
	outputs[0].Set_i32(int32(count/dstElem.TotalSize))
}
