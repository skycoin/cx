package cxcore

import()


func readData(inp *CXArgument, bytes []byte) interface{} {
	switch inp.Type {
	case TYPE_I8:
		data := readDataI8(bytes)
		if len(data) > 0 {
			return interface{}(data)
		}
	case TYPE_I16:
		data := readDataI16(bytes)
		if len(data) > 0 {
			return interface{}(data)
		}
	case TYPE_I32:
		data := readDataI32(bytes)
		if len(data) > 0 {
			return interface{}(data)
		}
	case TYPE_I64:
		data := readDataI64(bytes)
		if len(data) > 0 {
			return interface{}(data)
		}
	case TYPE_UI8:
		data := readDataUI8(bytes)
		if len(data) > 0 {
			return interface{}(data)
		}
	case TYPE_UI16:
		data := readDataUI16(bytes)
		if len(data) > 0 {
			return interface{}(data)
		}
	case TYPE_UI32:
		data := readDataUI32(bytes)
		if len(data) > 0 {
			return interface{}(data)
		}
	case TYPE_UI64:
		data := readDataUI64(bytes)
		if len(data) > 0 {
			return interface{}(data)
		}
	case TYPE_F32:
		data := readDataF32(bytes)
		if len(data) > 0 {
			return interface{}(data)
		}
	case TYPE_F64:
		data := readDataF64(bytes)
		if len(data) > 0 {
			return interface{}(data)
		}
	default:
		data := readDataUI8(bytes)
		if len(data) > 0 {
			return interface{}(data)
		}
	}

	return interface{}(nil)
}

// ReadSlice ...
func ReadSlice(fp int, inp *CXArgument, dataType int) interface{} {
	sliceOffset := GetSliceOffset(fp, inp)
	if sliceOffset >= 0 && (dataType < 0 || inp.Type == dataType) {
		slice := GetSliceData(sliceOffset, GetAssignmentElement(inp).Size)
		if slice != nil {
			return readData(inp, slice)
		}
	} else {
		panic(CX_RUNTIME_INVALID_ARGUMENT)
	}

	return interface{}(nil)
}

// ReadSliceBytes ...
func ReadSliceBytes(fp int, inp *CXArgument, dataType int) []byte {
	sliceOffset := GetSliceOffset(fp, inp)
	if sliceOffset >= 0 && (dataType < 0 || inp.Type == dataType) {
		slice := GetSliceData(sliceOffset, GetAssignmentElement(inp).Size)
		return slice
	}

	panic(CX_RUNTIME_INVALID_ARGUMENT)
}

// ReadArray ...
func ReadArray(fp int, inp *CXArgument, dataType int) interface{} {
	offset := GetFinalOffset(fp, inp)
	if dataType < 0 || inp.Type == dataType {
		array := ReadMemory(offset, inp)
		return readData(inp, array)
	}
	panic(CX_RUNTIME_INVALID_ARGUMENT)
}

// ReadObject ...
func ReadObject(fp int, inp *CXArgument, dataType int) interface{} {
	offset := GetFinalOffset(fp, inp)
	array := ReadMemory(offset, inp)
	return readData(inp, array)
}

// second section

// ReadBool ...
func ReadBool(fp int, inp *CXArgument) (out bool) {
	offset := GetFinalOffset(fp, inp)
	out = mustDeserializeBool(ReadMemory(offset, inp))
	return
}

// ReadStr ...
func ReadStr(fp int, inp *CXArgument) (out string) {
	var offset int32
	off := GetFinalOffset(fp, inp)
	if inp.Name == "" {
		// Then it's a literal.
		offset = int32(off)
	} else {
		offset = mustDeserializeI32(PROGRAM.Memory[off : off+TYPE_POINTER_SIZE])
	}

	if offset == 0 {
		// Then it's nil string.
		out = ""
		return
	}

	// We need to check if the string lives on the data segment or on the
	// heap to know if we need to take into consideration the object header's size.
	if int(offset) > PROGRAM.HeapStartsAt {
		size := mustDeserializeI32(PROGRAM.Memory[offset+OBJECT_HEADER_SIZE : offset+OBJECT_HEADER_SIZE+STR_HEADER_SIZE])
		mustDeserializeRaw(PROGRAM.Memory[offset+OBJECT_HEADER_SIZE:offset+OBJECT_HEADER_SIZE+STR_HEADER_SIZE+size], &out)
	} else {
		size := mustDeserializeI32(PROGRAM.Memory[offset : offset+STR_HEADER_SIZE])
		mustDeserializeRaw(PROGRAM.Memory[offset:offset+STR_HEADER_SIZE+size], &out)
	}

	return out
}
