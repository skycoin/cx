package cxcore

import (
	"github.com/skycoin/cx/cx/ast"
	"github.com/skycoin/cx/cx/constants"
)

//TODO: Rename opSliceLen
func opLen(expr *ast.CXExpression, fp int) {
	inp1, out1 := expr.Inputs[0], expr.Outputs[0]
	elt := ast.GetAssignmentElement(inp1)

	if elt.IsSlice || elt.Type == constants.TYPE_AFF {
		var sliceOffset = ast.GetSliceOffset(fp, inp1)
		if sliceOffset > 0 {
			sliceLen := ast.GetSliceHeader(sliceOffset)[4:8]
			ast.WriteMemory(ast.GetFinalOffset(fp, out1), sliceLen)
		} else if sliceOffset == 0 {
			ast.WriteI32(ast.GetFinalOffset(fp, out1), 0)
		} else {
			panic(constants.CX_RUNTIME_ERROR)
		}

		// TODO: Had to add elt.Lengths to avoid doing this for arrays, but not entirely sure why
	} else if elt.Type == constants.TYPE_STR && elt.Lengths == nil {
		var strOffset = ast.GetStrOffset(fp, inp1)
		// Checking if the string lives on the heap.
		if strOffset > ast.PROGRAM.HeapStartsAt {
			// Then it's on the heap and we need to consider
			// the object's header.
			strOffset += constants.OBJECT_HEADER_SIZE
		}

		ast.WriteMemory(ast.GetFinalOffset(fp, out1), ast.PROGRAM.Memory[strOffset:strOffset+constants.STR_HEADER_SIZE])
	} else {
		outV0 := int32(elt.Lengths[len(elt.Indexes)])
		ast.WriteI32(ast.GetFinalOffset(fp, out1), outV0)
	}
}

//TODO: Rename OpSliceAppend
func opAppend(expr *ast.CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]

	eltInp1 := ast.GetAssignmentElement(inp1)
	eltOut1 := ast.GetAssignmentElement(out1)
	if inp1.Type != inp2.Type || inp1.Type != out1.Type || !eltInp1.IsSlice || !eltOut1.IsSlice {
		panic(constants.CX_RUNTIME_INVALID_ARGUMENT)
	}

	var inputSliceLen int32
	inputSliceOffset := ast.GetSliceOffset(fp, inp1)
	if inputSliceOffset != 0 {
		inputSliceLen = ast.GetSliceLen(inputSliceOffset)
	}

	// Preparing slice in case more memory is needed for the new element.
	outputSliceOffset := ast.SliceAppendResize(fp, out1, inp1, inp2.Size)

	// We need to update the address of the output and input, as the final offsets
	// could be on the heap and they could have been moved by the GC.
	outputSlicePointer := ast.GetFinalOffset(fp, out1)

	if inp2.Type == constants.TYPE_STR || inp2.Type == constants.TYPE_AFF {
		var obj [4]byte
		ast.WriteMemI32(obj[:], 0, int32(ast.GetStrOffset(fp, inp2)))
		ast.SliceAppendWrite(outputSliceOffset, obj[:], inputSliceLen)
	} else {
		obj := ast.ReadMemory(ast.GetFinalOffset(fp, inp2), inp2)
		ast.SliceAppendWrite(outputSliceOffset, obj, inputSliceLen)
	}

	ast.WriteI32(outputSlicePointer, outputSliceOffset)
}

//TODO: Rename opSliceResize
func opResize(expr *ast.CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]

	if inp1.Type != out1.Type || !ast.GetAssignmentElement(inp1).IsSlice || !ast.GetAssignmentElement(out1).IsSlice {
		panic(constants.CX_RUNTIME_INVALID_ARGUMENT)
	}

	outputSliceOffset := int32(ast.SliceResize(fp, out1, inp1, ast.ReadI32(fp, inp2), ast.GetAssignmentElement(inp1).TotalSize))
	outputSlicePointer := ast.GetFinalOffset(fp, out1)
	ast.WriteI32(outputSlicePointer, outputSliceOffset)
}

//TODO: Rename opSliceInsertElement
func opInsert(expr *ast.CXExpression, fp int) {
	inp1, inp2, inp3, out1 := expr.Inputs[0], expr.Inputs[1], expr.Inputs[2], expr.Outputs[0]

	if inp1.Type != inp3.Type || inp1.Type != out1.Type || !ast.GetAssignmentElement(inp1).IsSlice || !ast.GetAssignmentElement(out1).IsSlice {
		panic(constants.CX_RUNTIME_INVALID_ARGUMENT)
	}

	outputSlicePointer := ast.GetFinalOffset(fp, out1)

	if inp3.Type == constants.TYPE_STR || inp3.Type == constants.TYPE_AFF {
		var obj [4]byte
		ast.WriteMemI32(obj[:], 0, int32(ast.GetStrOffset(fp, inp3)))
		outputSliceOffset := int32(ast.SliceInsert(fp, out1, inp1, ast.ReadI32(fp, inp2), obj[:]))
		ast.WriteI32(outputSlicePointer, outputSliceOffset)
	} else {
		obj := ast.ReadMemory(ast.GetFinalOffset(fp, inp3), inp3)
		outputSliceOffset := int32(ast.SliceInsert(fp, out1, inp1, ast.ReadI32(fp, inp2), obj))
		ast.WriteI32(outputSlicePointer, outputSliceOffset)
	}
}

//TODO: Rename opSliceRemoveElement
func opRemove(expr *ast.CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]

	if inp1.Type != out1.Type || !ast.GetAssignmentElement(inp1).IsSlice || !ast.GetAssignmentElement(out1).IsSlice {
		panic(constants.CX_RUNTIME_INVALID_ARGUMENT)
	}

	outputSlicePointer := ast.GetFinalOffset(fp, out1)
	outputSliceOffset := int32(ast.SliceRemove(fp, out1, inp1, ast.ReadI32(fp, inp2), int32(ast.GetAssignmentElement(inp1).TotalSize)))
	ast.WriteI32(outputSlicePointer, outputSliceOffset)
}

//TODO: Rename opSliceCopy
func opCopy(expr *ast.CXExpression, fp int) {
	dstInput := expr.Inputs[0]
	srcInput := expr.Inputs[1]
	dstOffset := ast.GetSliceOffset(fp, dstInput)
	srcOffset := ast.GetSliceOffset(fp, srcInput)

	dstElem := ast.GetAssignmentElement(dstInput)
	srcElem := ast.GetAssignmentElement(srcInput)

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
	ast.WriteI32(ast.GetFinalOffset(fp, expr.Outputs[0]), int32(count/dstElem.TotalSize))
}
