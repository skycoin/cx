package ast

import (
	"github.com/skycoin/cx/cx/constants"
	"github.com/skycoin/skycoin/src/cipher/encoder"
	"math"
)

// WriteMemory ...
func WriteMemory(offset int, byts []byte) {
	for c := 0; c < len(byts); c++ {
		PROGRAM.Memory[offset+c] = byts[c]
	}
}

// Utilities

// WriteObjectRef
// WARNING, is using heap variables?
//Is this "Write object ot heap?"
func WriteObjectData(obj []byte) int {
	size := len(obj) + constants.OBJECT_HEADER_SIZE
	heapOffset := AllocateSeq(size)
	WriteI32(heapOffset, int32(size))
	WriteMemory(heapOffset +constants.OBJECT_HEADER_SIZE, obj)
	return heapOffset
}

// WriteObject ...
func WriteObject(out1Offset int, obj []byte) {
	heapOffset := WriteObjectData(obj)
	WriteI32(out1Offset, int32(heapOffset))
}

// WriteStringData writes `str` to the heap as an object and returns its absolute offset.
func WriteStringData(str string) int {
	return WriteObjectData(encoder.Serialize(str))
}

// WriteString writes the string `str` on memory, starting at byte number `fp`.
func WriteString(fp int, str string, out *CXArgument) {
	WriteObject(GetOffset_str(fp, out), encoder.Serialize(str))
}


// WriteBool ...
func WriteBool(offset int, b bool) {
	v := byte(0)
	if b {
		v = 1
	}
	PROGRAM.Memory[offset] = v
}

// WriteI8 ...
func WriteI8(offset int, v int8) {
	PROGRAM.Memory[offset] = byte(v)
}

// WriteMemI8 ...
func WriteMemI8(mem []byte, offset int, v int8) {
	mem[offset] = byte(v)
}

// WriteI16 ...
func WriteI16(offset int, v int16) {
	PROGRAM.Memory[offset] = byte(v)
	PROGRAM.Memory[offset+1] = byte(v >> 8)
}

// WriteMemI16 ...
func WriteMemI16(mem []byte, offset int, v int16) {
	mem[offset] = byte(v)
	mem[offset+1] = byte(v >> 8)
}

// WriteI32 ...
func WriteI32(offset int, v int32) {
	PROGRAM.Memory[offset] = byte(v)
	PROGRAM.Memory[offset+1] = byte(v >> 8)
	PROGRAM.Memory[offset+2] = byte(v >> 16)
	PROGRAM.Memory[offset+3] = byte(v >> 24)
}

// WriteMemI32 ...
//TODO: This is an atomic type, do atomic write (not byte by byte)
//TODO: Fixed size type, doesnt need []byte slice passed in
func WriteMemI32(mem []byte, offset int, v int32) {
	mem[offset] = byte(v)
	mem[offset+1] = byte(v >> 8)
	mem[offset+2] = byte(v >> 16)
	mem[offset+3] = byte(v >> 24)
}

// WriteI64 ...
func WriteI64(offset int, v int64) {
	PROGRAM.Memory[offset] = byte(v)
	PROGRAM.Memory[offset+1] = byte(v >> 8)
	PROGRAM.Memory[offset+2] = byte(v >> 16)
	PROGRAM.Memory[offset+3] = byte(v >> 24)
	PROGRAM.Memory[offset+4] = byte(v >> 32)
	PROGRAM.Memory[offset+5] = byte(v >> 40)
	PROGRAM.Memory[offset+6] = byte(v >> 48)
	PROGRAM.Memory[offset+7] = byte(v >> 56)
}

// WriteMemI64 ...
//TODO: This is an atomic type, wtf, do atomic write
//TODO: Fixed size type, doesnt need []byte slice passed in
func WriteMemI64(mem []byte, offset int, v int64) {
	mem[offset] = byte(v)
	mem[offset+1] = byte(v >> 8)
	mem[offset+2] = byte(v >> 16)
	mem[offset+3] = byte(v >> 24)
	mem[offset+4] = byte(v >> 32)
	mem[offset+5] = byte(v >> 40)
	mem[offset+6] = byte(v >> 48)
	mem[offset+7] = byte(v >> 56)
}

// WriteUI8 ...
func WriteUI8(offset int, v uint8) {
	PROGRAM.Memory[offset] = v
}

// WriteMemUI8 ...
func WriteMemUI8(mem []byte, offset int, v uint8) {
	mem[offset] = v
}

// WriteUI16 ...
func WriteUI16(offset int, v uint16) {
	PROGRAM.Memory[offset] = byte(v)
	PROGRAM.Memory[offset+1] = byte(v >> 8)
}

// WriteMemUI16 ...
func WriteMemUI16(mem []byte, offset int, v uint16) {
	mem[offset] = byte(v)
	mem[offset+1] = byte(v >> 8)
}

// WriteUI32 ...
func WriteUI32(offset int, v uint32) {
	PROGRAM.Memory[offset] = byte(v)
	PROGRAM.Memory[offset+1] = byte(v >> 8)
	PROGRAM.Memory[offset+2] = byte(v >> 16)
	PROGRAM.Memory[offset+3] = byte(v >> 24)
}

// WriteMemUI32 ...
func WriteMemUI32(mem []byte, offset int, v uint32) {
	mem[offset] = byte(v)
	mem[offset+1] = byte(v >> 8)
	mem[offset+2] = byte(v >> 16)
	mem[offset+3] = byte(v >> 24)
}

// WriteUI64 ...
func WriteUI64(offset int, v uint64) {
	PROGRAM.Memory[offset] = byte(v)
	PROGRAM.Memory[offset+1] = byte(v >> 8)
	PROGRAM.Memory[offset+2] = byte(v >> 16)
	PROGRAM.Memory[offset+3] = byte(v >> 24)
	PROGRAM.Memory[offset+4] = byte(v >> 32)
	PROGRAM.Memory[offset+5] = byte(v >> 40)
	PROGRAM.Memory[offset+6] = byte(v >> 48)
	PROGRAM.Memory[offset+7] = byte(v >> 56)
}

// WriteMemUI64 ...
func WriteMemUI64(mem []byte, offset int, v uint64) {
	mem[offset] = byte(v)
	mem[offset+1] = byte(v >> 8)
	mem[offset+2] = byte(v >> 16)
	mem[offset+3] = byte(v >> 24)
	mem[offset+4] = byte(v >> 32)
	mem[offset+5] = byte(v >> 40)
	mem[offset+6] = byte(v >> 48)
	mem[offset+7] = byte(v >> 56)
}

// WriteF32 ...
func WriteF32(offset int, f float32) {
	v := math.Float32bits(f)
	PROGRAM.Memory[offset] = byte(v)
	PROGRAM.Memory[offset+1] = byte(v >> 8)
	PROGRAM.Memory[offset+2] = byte(v >> 16)
	PROGRAM.Memory[offset+3] = byte(v >> 24)
}

// WriteMemF32 ...
func WriteMemF32(mem []byte, offset int, f float32) {
	v := math.Float32bits(f)
	mem[offset] = byte(v)
	mem[offset+1] = byte(v >> 8)
	mem[offset+2] = byte(v >> 16)
	mem[offset+3] = byte(v >> 24)
}

// WriteF64 ...
func WriteF64(offset int, f float64) {
	v := math.Float64bits(f)
	PROGRAM.Memory[offset] = byte(v)
	PROGRAM.Memory[offset+1] = byte(v >> 8)
	PROGRAM.Memory[offset+2] = byte(v >> 16)
	PROGRAM.Memory[offset+3] = byte(v >> 24)
	PROGRAM.Memory[offset+4] = byte(v >> 32)
	PROGRAM.Memory[offset+5] = byte(v >> 40)
	PROGRAM.Memory[offset+6] = byte(v >> 48)
	PROGRAM.Memory[offset+7] = byte(v >> 56)
}

// WriteMemF64 ...
func WriteMemF64(mem []byte, offset int, f float64) {
	v := math.Float64bits(f)
	mem[offset] = byte(v)
	mem[offset+1] = byte(v >> 8)
	mem[offset+2] = byte(v >> 16)
	mem[offset+3] = byte(v >> 24)
	mem[offset+4] = byte(v >> 32)
	mem[offset+5] = byte(v >> 40)
	mem[offset+6] = byte(v >> 48)
	mem[offset+7] = byte(v >> 56)
}
