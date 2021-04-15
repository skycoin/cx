package ast

// SerializeCXProgramV3 translates cx program to slice of bytes that we can save.
// These slice of bytes can then be deserialize in the future and
// be translated back to cx program.
func SerializeCXProgramV3(prgrm *CXProgram, includeMemory bool) (b []byte) {
	s := SerializedCXProgram{}
	initSerialization(prgrm, &s, includeMemory)

	// serialize cx program's packages,
	// structs, functions, etc.
	serializeCXProgramElements(prgrm, &s)

	// serialize cx program's program
	serializeProgram(prgrm, &s)

	// assign cx program's offsets
	assignSerializedCXProgramOffset(&s)

	// serializing everything
	b, err := EncodeSerializedCXProgram(&s)
	if err != nil {
		panic(err)
	}

	return b
}

func DeserializeCXProgramV3(b []byte) *CXProgram {
	prgrm := &CXProgram{}
	var sPrgrm SerializedCXProgram

	DecodeSerializedCXProgram(b, &sPrgrm)
	initDeserialization(prgrm, &sPrgrm)
	return prgrm
}
