package cxcore

import()

// ReadMemory ...
//TODO: Avoid all read memory commands for fixed width types (i32,f32,etc)
//TODO: Make "ReadMemoryI32", "ReadMemoryI16", etc
func ReadMemory(offset int, arg *CXArgument) []byte {
	size := GetSize(arg)
	return PROGRAM.Memory[offset : offset+size]
}
