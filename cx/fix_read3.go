package cxcore

// ReadStr ...
/*
func ReadStr(fp int, inp *CXArgument) (out string) {
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
*/


//TODO:
// -ReadStr is weird
// -ReadStrFromOffset is weird

