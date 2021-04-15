package ast

import "github.com/skycoin/skycoin/src/cipher/encoder"

// SerializeCXProgramV2 translates cx program to slice of bytes that we can save.
// These slice of bytes can then be deserialize in the future and
// be translated back to cx program.
func SerializeCXProgramV2(prgrm *CXProgram, includeMemory bool) (b []byte) {
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
	b = encoder.Serialize(s)
	b = append(b, s.Strings...)
	b = append(b, s.Memory...)

	return b
}
