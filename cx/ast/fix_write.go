package ast

import (
	"github.com/skycoin/cx/cx/constants"
	"github.com/skycoin/skycoin/src/cipher/encoder"
	"math"
)

// TODO: Move to cx/helper/serialize.go

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
	WriteObject(GetFinalOffset(fp, out), encoder.Serialize(str))
}

// FromI32 ...
func FromI32(in int32) []byte {
	var b [4]byte
    WriteMemI32(b[:4], 0, in)
	return b[:4]
}

// FromI64 ...
func FromI64(in int64) []byte {
	var b [8]byte
    WriteMemI64(b[:8], 0, in)
	return b[:8]
}

// WriteBool ...
func WriteBool(offset int, v bool) {
	WriteMemBool(PROGRAM.Memory, offset, v)
}

// WriteI8 ...
func WriteI8(offset int, v int8) {
    WriteMemI8(PROGRAM.Memory, offset, v)
}

// WriteI16 ...
func WriteI16(offset int, v int16) {
    WriteMemI16(PROGRAM.Memory, offset, v)
}

// WriteI32 ...
func WriteI32(offset int, v int32) {
    WriteMemI32(PROGRAM.Memory, offset, v)
}

// WriteI64 ...
func WriteI64(offset int, v int64) {
    WriteMemI64(PROGRAM.Memory, offset, v)
}

// WriteUI8 ...
func WriteUI8(offset int, v uint8) {
    WriteMemUI8(PROGRAM.Memory, offset, v)
}

// WriteUI16 ...
func WriteUI16(offset int, v uint16) {
    WriteMemUI16(PROGRAM.Memory, offset, v)
}

// WriteUI32 ...
func WriteUI32(offset int, v uint32) {
    WriteMemUI32(PROGRAM.Memory, offset, v)
}

// WriteUI64 ...
func WriteUI64(offset int, v uint64) {
    WriteMemUI64(PROGRAM.Memory, offset, v)
}

// WriteF32 ...
func WriteF32(offset int, v float32) {
    WriteMemF32(PROGRAM.Memory, offset, v)
}

// WriteF64 ...
func WriteF64(offset int, v float64) {
    WriteMemF64(PROGRAM.Memory, offset, v)
}

// WriteMemBool
func WriteMemBool(mem []byte, offset int, b bool) {
	v := byte(0)
	if b {
		v = 1
	}
	mem[offset] = v
}


// WriteMemI8 ...
func WriteMemI8(mem []byte, offset int, v int8) {
	mem[offset] = byte(v)
}

// WriteMemI16 ...
func WriteMemI16(mem []byte, offset int, v int16) {
	mem[offset] = byte(v)
	mem[offset+1] = byte(v >> 8)
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

// WriteMemUI8 ...
func WriteMemUI8(mem []byte, offset int, v uint8) {
	mem[offset] = v
}

// WriteMemUI16 ...
func WriteMemUI16(mem []byte, offset int, v uint16) {
	mem[offset] = byte(v)
	mem[offset+1] = byte(v >> 8)
}

// WriteMemUI32 ...
func WriteMemUI32(mem []byte, offset int, v uint32) {
	mem[offset] = byte(v)
	mem[offset+1] = byte(v >> 8)
	mem[offset+2] = byte(v >> 16)
	mem[offset+3] = byte(v >> 24)
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

// WriteMemF32 ...
func WriteMemF32(mem []byte, offset int, f float32) {
	v := math.Float32bits(f)
	mem[offset] = byte(v)
	mem[offset+1] = byte(v >> 8)
	mem[offset+2] = byte(v >> 16)
	mem[offset+3] = byte(v >> 24)
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
