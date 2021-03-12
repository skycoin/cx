package cxcore

import (
	"github.com/skycoin/skycoin/src/cipher/encoder"
)

//Why do these functions need CXArgument as imput!?

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

// ReadData_i8 ...
func ReadData_i8(fp int, inp *CXArgument, dataType int) interface{} {
	return ReadData(fp, inp, dataType)
}

// ReadData_i16 ...
func ReadData_i16(fp int, inp *CXArgument, dataType int) interface{} {
	return ReadData(fp, inp, dataType)
}

// ReadData_i32 ...
func ReadData_i32(fp int, inp *CXArgument, dataType int) interface{} {
	return ReadData(fp, inp, dataType)
}

// ReadData_i64 ...
func ReadData_i64(fp int, inp *CXArgument, dataType int) interface{} {
	return ReadData(fp, inp, dataType)
}

// ReadData_ui8 ...
func ReadData_ui8(fp int, inp *CXArgument, dataType int) interface{} {
	return ReadData(fp, inp, dataType)
}

// ReadData_ui16 ...
func ReadData_ui16(fp int, inp *CXArgument, dataType int) interface{} {
	return ReadData(fp, inp, dataType)
}

// ReadData_ui32 ...
func ReadData_ui32(fp int, inp *CXArgument, dataType int) interface{} {
	return ReadData(fp, inp, dataType)
}

// ReadData_ui64 ...
func ReadData_ui64(fp int, inp *CXArgument, dataType int) interface{} {
	return ReadData(fp, inp, dataType)
}

// ReadData_f32 ...
func ReadData_f32(fp int, inp *CXArgument, dataType int) interface{} {
	return ReadData(fp, inp, dataType)
}

// ReadData_f64 ...
func ReadData_f64(fp int, inp *CXArgument, dataType int) interface{} {
	return ReadData(fp, inp, dataType)
}

//Note: Only called once and only by ReadData
// ReadObject ...
func ReadObject(fp int, inp *CXArgument, dataType int) interface{} {
	offset := GetFinalOffset(fp, inp)
	array := ReadMemory(offset, inp)
	return readAtomic(inp, array)
}

//Note: I modified this to crash if invalid type was used
func readAtomic(inp *CXArgument, bytes []byte) interface{} {
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
	//should this crash if it gets here?
	panic(CX_RUNTIME_INVALID_ARGUMENT) //Note: modified this so it crashes if it gets here for some reason
	return interface{}(nil)
}

// ReadSlice ...
func ReadSlice(fp int, inp *CXArgument, dataType int) interface{} {
	sliceOffset := GetSliceOffset(fp, inp)
	if sliceOffset >= 0 && (dataType < 0 || inp.Type == dataType) {
		slice := GetSliceData(sliceOffset, GetAssignmentElement(inp).Size)
		if slice != nil {
			return readAtomic(inp, slice) //readData
		}
	} else {
		panic(CX_RUNTIME_INVALID_ARGUMENT)
	}

	return interface{}(nil)
}

// ReadArray ...
func ReadArray(fp int, inp *CXArgument, dataType int) interface{} {
	offset := GetFinalOffset(fp, inp)
	if dataType < 0 || inp.Type == dataType {
		array := ReadMemory(offset, inp)
		return readAtomic(inp, array) //readData
	}
	panic(CX_RUNTIME_INVALID_ARGUMENT)
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

// second section

// ReadBool ...
func ReadBool(fp int, inp *CXArgument) (out bool) {
	offset := GetFinalOffset(fp, inp)
	out = DeserializeBool(ReadMemory(offset, inp))
	return
}

// ReadStrFromOffset ...
func ReadStrFromOffset(off int, inp *CXArgument) (out string) {
	var offset int32
	if inp.Name == "" {
		// Then it's a literal.
		offset = int32(off)
	} else {
		offset = Deserialize_i32(PROGRAM.Memory[off : off+TYPE_POINTER_SIZE])
	}

	if offset == 0 {
		// Then it's nil string.
		out = ""
		return
	}

	// We need to check if the string lives on the data segment or on the
	// heap to know if we need to take into consideration the object header's size.
	if int(offset) > PROGRAM.HeapStartsAt {
		size := Deserialize_i32(PROGRAM.Memory[offset+OBJECT_HEADER_SIZE : offset+OBJECT_HEADER_SIZE+STR_HEADER_SIZE])
		DeserializeRaw(PROGRAM.Memory[offset+OBJECT_HEADER_SIZE:offset+OBJECT_HEADER_SIZE+STR_HEADER_SIZE+size], &out)
	} else {
		size := Deserialize_i32(PROGRAM.Memory[offset : offset+STR_HEADER_SIZE])
		DeserializeRaw(PROGRAM.Memory[offset:offset+STR_HEADER_SIZE+size], &out)
	}

	return out	
}

// ReadStr ...
func ReadStr(fp int, inp *CXArgument) (out string) {
	off := GetFinalOffset(fp, inp)
	return ReadStrFromOffset(off, inp)
}

// ReadStringFromObject reads the string located at offset `off`.
func ReadStringFromObject(off int32) string {
	var plusOff int32
	if int(off) > PROGRAM.HeapStartsAt {
		// Found in heap segment.
		plusOff += OBJECT_HEADER_SIZE
	}

	size := Deserialize_i32(PROGRAM.Memory[off+plusOff : off+plusOff+STR_HEADER_SIZE])

	str := ""
	_, err := encoder.DeserializeRaw(PROGRAM.Memory[off+plusOff:off+plusOff+STR_HEADER_SIZE+size], &str)
	if err != nil {
		panic(err)
	}
	return str
}