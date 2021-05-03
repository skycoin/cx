package ast

import (
	"github.com/pierrec/lz4"
)

// SerializeCXProgramV2 is using skyencoder generated
// encoder/decoder for serializing cx program.
func SerializeCXProgramV2(prgrm *CXProgram, includeDataMemory, useCompression bool) (b []byte) {
	s := SerializedCXProgram{}
	initSerialization(prgrm, &s, includeDataMemory, useCompression)

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

	if useCompression {
		// Compress using LZ4
		CompressBytesLZ4(&b)
	}

	return b
}

func DeserializeCXProgramV2(b []byte, useCompression bool) *CXProgram {
	prgrm := &CXProgram{}
	var sPrgrm SerializedCXProgram

	if useCompression {
		// Uncompress using LZ4
		UncompressBytesLZ4(&b)
	}

	DecodeSerializedCXProgram(b, &sPrgrm)
	initDeserialization(prgrm, &sPrgrm)
	return prgrm
}

func CompressBytesLZ4(data *[]byte) {
	buf := make([]byte, len(*data))
	ht := make([]int, 64<<10) // buffer for the compression table

	n, err := lz4.CompressBlock(*data, buf, ht)
	if err != nil {
		panic(err)
	}
	if n >= len(*data) {
		panic("data not compressible")
	}

	*data = buf[:n] // compressed data
}

func UncompressBytesLZ4(data *[]byte) {
	// Allocated a very large buffer for decompression.
	out := make([]byte, 1000*len(*data))
	n, err := lz4.UncompressBlock(*data, out)
	if err != nil {
		panic(err)
	}
	*data = out[:n] // uncompressed data
}
