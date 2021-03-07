package cxcore

import()

// ReadData ...
func ReadData(fp int, inp *CXArgument, dataType int) interface{} {
	elt := GetAssignmentElement(inp)
	if elt.IsSlice {
		return ReadSlice(fp, inp, dataType)
	} else if elt.IsArray {
		return ReadArray(fp, inp, dataType)
	} else {
		return ReadObject(fp, inp, dataType)
	}
}

func readDataI8(bytes []byte) (out []int8) {
	count := len(bytes)
	if count > 0 {
		out = make([]int8, count)
		for i := 0; i < count; i++ {
			out[i] = mustDeserializeI8(bytes[i:])
		}
	}
	return
}

func readDataUI8(bytes []byte) (out []uint8) {
	count := len(bytes)
	if count > 0 {
		out = make([]uint8, count)
		for i := 0; i < count; i++ {
			out[i] = mustDeserializeUI8(bytes[i:])
		}
	}
	return
}

func readDataI16(bytes []byte) (out []int16) {
	count := len(bytes) / 2
	if count > 0 {
		out = make([]int16, count)
		for i := 0; i < count; i++ {
			out[i] = mustDeserializeI16(bytes[i*2:])
		}
	}
	return
}

func readDataUI16(bytes []byte) (out []uint16) {
	count := len(bytes) / 2
	if count > 0 {
		out = make([]uint16, count)
		for i := 0; i < count; i++ {
			out[i] = mustDeserializeUI16(bytes[i*2:])
		}
	}
	return
}

func readDataI32(bytes []byte) (out []int32) {
	count := len(bytes) / 4
	if count > 0 {
		out = make([]int32, count)
		for i := 0; i < count; i++ {
			out[i] = mustDeserializeI32(bytes[i*4:])
		}
	}
	return
}

func readDataUI32(bytes []byte) (out []uint32) {
	count := len(bytes) / 4
	if count > 0 {
		out = make([]uint32, count)
		for i := 0; i < count; i++ {
			out[i] = mustDeserializeUI32(bytes[i*4:])
		}
	}
	return
}

func readDataI64(bytes []byte) (out []int64) {
	count := len(bytes) / 8
	if count > 0 {
		out = make([]int64, count)
		for i := 0; i < count; i++ {
			out[i] = mustDeserializeI64(bytes[i*8:])
		}
	}
	return
}

func readDataUI64(bytes []byte) (out []uint64) {
	count := len(bytes) / 8
	if count > 0 {
		out = make([]uint64, count)
		for i := 0; i < count; i++ {
			out[i] = mustDeserializeUI64(bytes[i*8:])
		}
	}
	return
}

func readDataF32(bytes []byte) (out []float32) {
	count := len(bytes) / 4
	if count > 0 {
		out = make([]float32, count)
		for i := 0; i < count; i++ {
			out[i] = mustDeserializeF32(bytes[i*4:])
		}
	}
	return
}

func readDataF64(bytes []byte) (out []float64) {
	count := len(bytes) / 8
	if count > 0 {
		out = make([]float64, count)
		for i := 0; i < count; i++ {
			out[i] = mustDeserializeF64(bytes[i*8:])
		}
	}
	return
}

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

// ReadI8 ...
func ReadI8(fp int, inp *CXArgument) int8 {
	return mustDeserializeI8(ReadMemory(GetFinalOffset(fp, inp), inp))
}

// ReadI16 ...
func ReadI16(fp int, inp *CXArgument) int16 {
	return mustDeserializeI16(ReadMemory(GetFinalOffset(fp, inp), inp))
}

// ReadI32 ...
func ReadI32(fp int, inp *CXArgument) int32 {
	return mustDeserializeI32(ReadMemory(GetFinalOffset(fp, inp), inp))
}

// ReadI64 ...
func ReadI64(fp int, inp *CXArgument) int64 {
	return mustDeserializeI64(ReadMemory(GetFinalOffset(fp, inp), inp))
}

// ReadUI8 ...
func ReadUI8(fp int, inp *CXArgument) uint8 {
	return mustDeserializeUI8(ReadMemory(GetFinalOffset(fp, inp), inp))
}

// ReadUI16 ...
func ReadUI16(fp int, inp *CXArgument) uint16 {
	return mustDeserializeUI16(ReadMemory(GetFinalOffset(fp, inp), inp))
}

// ReadUI32 ...
func ReadUI32(fp int, inp *CXArgument) uint32 {
	return mustDeserializeUI32(ReadMemory(GetFinalOffset(fp, inp), inp))
}

// ReadUI64 ...
func ReadUI64(fp int, inp *CXArgument) uint64 {
	return mustDeserializeUI64(ReadMemory(GetFinalOffset(fp, inp), inp))
}

// ReadF32 ...
func ReadF32(fp int, inp *CXArgument) float32 {
	return mustDeserializeF32(ReadMemory(GetFinalOffset(fp, inp), inp))
}

// ReadF64 ...
func ReadF64(fp int, inp *CXArgument) float64 {
	return mustDeserializeF64(ReadMemory(GetFinalOffset(fp, inp), inp))
}