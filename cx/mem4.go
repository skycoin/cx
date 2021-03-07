package cxcore

import()

/*
./cxfx/op_opengl.go:437:	obj := ReadMemory(GetFinalOffset(fp, inp1), inp1)
./cx/op_http.go:326:	// reqByts := ReadMemory(GetFinalOffset(fp, inp1), inp1)
./cx/op_http.go:493:	byts1 := ReadMemory(GetFinalOffset(fp, inp1), inp1)
./cx/fix_read3.go:110:		array := ReadMemory(offset, inp)
./cx/fix_read3.go:119:	array := ReadMemory(offset, inp)
./cx/fix_read3.go:128:	out = mustDeserializeBool(ReadMemory(offset, inp))
./cx/op_aff.go:101:	return ReadMemory(GetFinalOffset(
./cx/op_und.go:548:		obj := ReadMemory(GetFinalOffset(fp, inp2), inp2)
./cx/op_und.go:588:		obj := ReadMemory(GetFinalOffset(fp, inp3), inp3)
./cx/execute.go:291:					ReadMemory(
./cx/execute.go:424:		mem := ReadMemory(GetFinalOffset(newFP, out), out)
./cx/op_testing.go:22:		byts1 = ReadMemory(GetFinalOffset(fp, inp1), inp1)
./cx/op_testing.go:23:		byts2 = ReadMemory(GetFinalOffset(fp, inp2), inp2)
./cx/fix_read.go:11:	return mustDeserializeI8(ReadMemory(GetFinalOffset(fp, inp), inp))
./cx/fix_read.go:16:	return mustDeserializeI16(ReadMemory(GetFinalOffset(fp, inp), inp))
./cx/fix_read.go:21:	return mustDeserializeI32(ReadMemory(GetFinalOffset(fp, inp), inp))
./cx/fix_read.go:26:	return mustDeserializeI64(ReadMemory(GetFinalOffset(fp, inp), inp))
./cx/fix_read.go:31:	return mustDeserializeUI8(ReadMemory(GetFinalOffset(fp, inp), inp))
./cx/fix_read.go:36:	return mustDeserializeUI16(ReadMemory(GetFinalOffset(fp, inp), inp))
./cx/fix_read.go:41:	return mustDeserializeUI32(ReadMemory(GetFinalOffset(fp, inp), inp))
./cx/fix_read.go:46:	return mustDeserializeUI64(ReadMemory(GetFinalOffset(fp, inp), inp))
./cx/fix_read.go:51:	return mustDeserializeF32(ReadMemory(GetFinalOffset(fp, inp), inp))
./cx/fix_read.go:56:	return mustDeserializeF64(ReadMemory(GetFinalOffset(fp, inp), inp))
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
//TODO: Avoid all read memory commands for fixed width types (i32,f32,etc)
//TODO: Make "ReadMemoryI32", "ReadMemoryI16", etc
func ReadMemory(offset int, arg *CXArgument) []byte {
	size := GetSize(arg)
	return PROGRAM.Memory[offset : offset+size]
}
