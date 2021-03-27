package ast

import (
	"github.com/skycoin/cx/cx/constants"
	"github.com/skycoin/cx/cx/helper"
	"github.com/skycoin/skycoin/src/cipher/encoder"
)

/*
./cxfx/op_opengl.go:437:	obj := ReadMemory(GetFinalOffset(fp, inp1), inp1)
./cx/op_http.go:326:	// reqByts := ReadMemory(GetFinalOffset(fp, inp1), inp1)
./cx/op_http.go:493:	byts1 := ReadMemory(GetFinalOffset(fp, inp1), inp1)
./cx/fix_read3.go:110:		array := ReadMemory(offset, inp)
./cx/fix_read3.go:119:	array := ReadMemory(offset, inp)
./cx/fix_read3.go:128:	out = DeserializeBool(ReadMemory(offset, inp))
./cx/op_aff.go:101:	return ReadMemory(GetFinalOffset(
./cx/op_und.go:548:		obj := ReadMemory(GetFinalOffset(fp, inp2), inp2)
./cx/op_und.go:588:		obj := ReadMemory(GetFinalOffset(fp, inp3), inp3)
./cx/execute.go:291:					ReadMemory(
./cx/execute.go:424:		mem := ReadMemory(GetFinalOffset(newFP, out), out)
./cx/op_testing.go:22:		byts1 = ReadMemory(GetFinalOffset(fp, inp1), inp1)
./cx/op_testing.go:23:		byts2 = ReadMemory(GetFinalOffset(fp, inp2), inp2)
./cx/fix_read.go:11:	return Deserialize_i8(ReadMemory(GetFinalOffset(fp, inp), inp))
./cx/fix_read.go:16:	return Deserialize_i16(ReadMemory(GetFinalOffset(fp, inp), inp))
./cx/fix_read.go:21:	return Deserialize_i32(ReadMemory(GetFinalOffset(fp, inp), inp))
./cx/fix_read.go:26:	return Deserialize_i64(ReadMemory(GetFinalOffset(fp, inp), inp))
./cx/fix_read.go:31:	return Deserialize_ui8(ReadMemory(GetFinalOffset(fp, inp), inp))
./cx/fix_read.go:36:	return Deserialize_ui16(ReadMemory(GetFinalOffset(fp, inp), inp))
./cx/fix_read.go:41:	return Deserialize_ui32(ReadMemory(GetFinalOffset(fp, inp), inp))
./cx/fix_read.go:46:	return Deserialize_ui64(ReadMemory(GetFinalOffset(fp, inp), inp))
./cx/fix_read.go:51:	return Deserialize_f32(ReadMemory(GetFinalOffset(fp, inp), inp))
./cx/fix_read.go:56:	return Deserialize_f64(ReadMemory(GetFinalOffset(fp, inp), inp))
./cx/op_misc.go:9:	byts := ReadMemory(inpOffset, arg)
./cx/op_misc.go:41:			WriteMemory(out1Offset, ReadMemory(inp1Offset, inp1))
./cx/op.go:183:// ReadMemory ...
./cx/op.go:185://TODO: Make "ReadMemoryI32", "ReadMemoryI16", etc
./cx/op.go:186:func ReadMemory(offset int, arg *CXArgument) []byte {
./cx/fix_readmemory.go:5:// ReadMemory ...
./cx/fix_readmemory.go:7://TODO: Make "ReadMemoryI32", "ReadMemoryI16", etc
./cx/fix_readmemory.go:8:func ReadMemory(offset int, arg *CXArgument) []byte {
*/
// ReadMemory ...
//TODO: DELETE THIS FUNCTION
//TODO: Avoid all read memory commands for fixed width types (i32,f32,etc)
//TODO: Make "ReadMemoryI32", "ReadMemoryI16", etc
func ReadMemory(offset int, arg *CXArgument) []byte {
	size := GetSize(arg)
	return PROGRAM.Memory[offset : offset+size]
}

// ReadStr ...
func ReadStr(fp int, inp *CXArgument) (out string) {
	off := GetFinalOffset(fp, inp)
	return ReadStrFromOffset(off, inp)
}

// ReadStrFromOffset ...
func ReadStrFromOffset(off int, inp *CXArgument) (out string) {
	var offset int32
	if inp.Name == "" {
		// Then it's a literal.
		offset = int32(off)
	} else {
		offset = helper.Deserialize_i32(PROGRAM.Memory[off : off+constants.TYPE_POINTER_SIZE])
	}

	if offset == 0 {
		// Then it's nil string.
		out = ""
		return
	}

	// We need to check if the string lives on the data segment or on the
	// heap to know if we need to take into consideration the object header's size.
	if int(offset) > PROGRAM.HeapStartsAt {
		size := helper.Deserialize_i32(PROGRAM.Memory[offset+constants.OBJECT_HEADER_SIZE : offset+constants.OBJECT_HEADER_SIZE+constants.STR_HEADER_SIZE])
		helper.DeserializeRaw(PROGRAM.Memory[offset+constants.OBJECT_HEADER_SIZE:offset+constants.OBJECT_HEADER_SIZE+constants.STR_HEADER_SIZE+size], &out)
	} else {
		size := helper.Deserialize_i32(PROGRAM.Memory[offset : offset+constants.STR_HEADER_SIZE])
		helper.DeserializeRaw(PROGRAM.Memory[offset:offset+constants.STR_HEADER_SIZE+size], &out)
	}

	return out
}

// ReadStringFromObject reads the string located at offset `off`.
func ReadStringFromObject(off int32) string {
	var plusOff int32
	if int(off) > PROGRAM.HeapStartsAt {
		// Found in heap segment.
		plusOff += constants.OBJECT_HEADER_SIZE
	}

	size := helper.Deserialize_i32(PROGRAM.Memory[off+plusOff : off+plusOff+constants.STR_HEADER_SIZE])

	str := ""
	_, err := encoder.DeserializeRaw(PROGRAM.Memory[off+plusOff:off+plusOff+constants.STR_HEADER_SIZE+size], &str)
	if err != nil {
		panic(err)
	}
	return str
}
