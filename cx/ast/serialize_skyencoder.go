package ast

// SerializeCXProgramV2 is using skyencoder generated
// encoder/decoder for serializing cx program.
func SerializeCXProgramV2(prgrm *CXProgram, includeMemory bool) (b []byte) {
	s := SerializedCXProgram{}
	initSerialization(prgrm, &s, includeMemory)

	// serialize cx program's packages,
	// structs, functions, etc.
	serializeCXProgramElements(prgrm, &s)

	// serialize cx program's program
	serializeProgram(prgrm, &s)

	// serializing everything
	b, err := EncodeSerializedCXProgram(&s)
	if err != nil {
		panic(err)
	}

	return b
}

func DeserializeCXProgramV2(b []byte) *CXProgram {
	prgrm := &CXProgram{}
	var sPrgrm SerializedCXProgram

	DecodeSerializedCXProgram(b, &sPrgrm)
	initDeserialization(prgrm, &sPrgrm)
	return prgrm
}
