package ast

import (
	"github.com/skycoin/cx/cx/constants"
	"github.com/skycoin/cx/cx/helper"
	"github.com/skycoin/skycoin/src/cipher/encoder"
)

/*
cx/packages/http/op_http.go:263:	// reqByts := ReadMemory(GetFinalOffset(fp, inp1), inp1)
cx/packages/http/op_http.go:430:	byts1 := ast.ReadMemory(ast.GetFinalOffset(fp, inp1), inp1)
cx/ast/tostring.go:230:        memory : ReadMemory(offset, elt),
cx/ast/ast_cxcall.go:45:					ReadMemory(
cx/ast/read.go:15:func ReadMemory(offset int, arg *CXArgument) []byte {
cx/ast/ast_value.go:123:	return ReadMemory(value.Offset, value.Arg)
cx/execute/callback.go:50:		mem := ast.ReadMemory(ast.GetFinalOffset(newFP, out), out)
cx/opcodes/op_aff.go:103:	return ast.ReadMemory(ast.GetFinalOffset(
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
