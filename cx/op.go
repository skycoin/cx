package cxcore

import (
	"fmt"
	"math"

	"github.com/skycoin/skycoin/src/cipher/encoder"
)

// GetSize ...
func GetSize(arg *CXArgument) int {
	if len(arg.Fields) > 0 {
		return GetSize(arg.Fields[len(arg.Fields)-1])
	}

	derefCount := len(arg.DereferenceOperations)
	if derefCount > 0 {
		deref := arg.DereferenceOperations[derefCount-1]
		if deref == DEREF_SLICE || deref == DEREF_ARRAY {
			return arg.Size
		}
	}

	for decl := range arg.DeclarationSpecifiers {
		if decl == DECL_POINTER {
			return arg.TotalSize
		}
	}

	if arg.CustomType != nil {
		return arg.CustomType.Size
	}

	return arg.TotalSize
}

// GetDerefSize ...
func GetDerefSize(arg *CXArgument) int {
	if arg.CustomType != nil {
		return arg.CustomType.Size
	}
	return arg.Size
}

// CalculateDereferences ...
func CalculateDereferences(arg *CXArgument, finalOffset *int, fp int, dbg bool) {
	var isPointer bool
	var baseOffset int
	var sizeofElement int

	idxCounter := 0
	for _, op := range arg.DereferenceOperations {
		switch op {
		case DEREF_SLICE:
			if len(arg.Indexes) == 0 {
				continue
			}

			isPointer = false
			var offset int32
			var byts []byte

			byts = PROGRAM.Memory[*finalOffset : *finalOffset+TYPE_POINTER_SIZE]

			offset = mustDeserializeI32(byts)

			*finalOffset = int(offset)

			baseOffset = *finalOffset

			*finalOffset += OBJECT_HEADER_SIZE
			*finalOffset += SLICE_HEADER_SIZE

			sizeToUse := GetDerefSize(arg)
			*finalOffset += int(ReadI32(fp, arg.Indexes[idxCounter])) * sizeToUse
			if !IsValidSliceIndex(baseOffset, *finalOffset, sizeToUse) {
				panic(CX_RUNTIME_SLICE_INDEX_OUT_OF_RANGE)
			}

			idxCounter++
		case DEREF_ARRAY:
			if len(arg.Indexes) == 0 {
				continue
			}
			var subSize = int(1)
			for _, len := range arg.Lengths[idxCounter+1:] {
				subSize *= len
			}

			sizeToUse := GetDerefSize(arg)

			baseOffset = *finalOffset
			sizeofElement = subSize * sizeToUse
			*finalOffset += int(ReadI32(fp, arg.Indexes[idxCounter])) * sizeofElement
			idxCounter++
		case DEREF_POINTER:
			isPointer = true
			var offset int32
			var byts []byte

			byts = PROGRAM.Memory[*finalOffset : *finalOffset+TYPE_POINTER_SIZE]

			offset = mustDeserializeI32(byts)
			*finalOffset = int(offset)
		}
		if dbg {
			fmt.Println("\tupdate", arg.Name, arg.DereferenceOperations, *finalOffset, PROGRAM.Memory[*finalOffset:*finalOffset+10])
		}
	}
	if dbg {
		fmt.Println("\tupdate", arg.Name, arg.DereferenceOperations, *finalOffset, PROGRAM.Memory[*finalOffset:*finalOffset+10])
	}

	// if *finalOffset >= PROGRAM.HeapStartsAt {
	if *finalOffset >= PROGRAM.HeapStartsAt && isPointer {
		// then it's an object
		*finalOffset += OBJECT_HEADER_SIZE
		if arg.IsSlice {
			*finalOffset += SLICE_HEADER_SIZE
			if !IsValidSliceIndex(baseOffset, *finalOffset, sizeofElement) {
				panic(CX_RUNTIME_SLICE_INDEX_OUT_OF_RANGE)
			}
		}
	}
}

// GetStrOffset ...
func GetStrOffset(fp int, arg *CXArgument) int {
	strOffset := GetFinalOffset(fp, arg)
	if arg.Name != "" {
		// then it's not a literal
		offset := mustDeserializeI32(PROGRAM.Memory[strOffset : strOffset+TYPE_POINTER_SIZE])
		strOffset = int(offset)
	}
	return strOffset
}

// GetFinalOffset ...
func GetFinalOffset(fp int, arg *CXArgument) int {
	// defer RuntimeError(PROGRAM)
	// var elt *CXArgument
	finalOffset := arg.Offset
	// var fldIdx int

	// elt = arg

	var dbg bool
	if arg.Name != "" {
		dbg = false
	}

	if finalOffset < PROGRAM.StackSize {
		// Then it's in the stack, not in data or heap and we need to consider the frame pointer.
		finalOffset += fp
	}

	if dbg {
		fmt.Println("(start", arg.Name, fmt.Sprintf("%s:%d", arg.FileName, arg.FileLine), arg.DereferenceOperations, finalOffset, PROGRAM.Memory[finalOffset:finalOffset+10])
	}

	// elt = arg
	CalculateDereferences(arg, &finalOffset, fp, dbg)
	for _, fld := range arg.Fields {
		// elt = fld
		finalOffset += fld.Offset
		CalculateDereferences(fld, &finalOffset, fp, dbg)
	}

	if dbg {
		fmt.Println("\t\tresult", finalOffset, PROGRAM.Memory[finalOffset:finalOffset+10], "...)")
	}

	return finalOffset
}

// ReadMemory ...
func ReadMemory(offset int, arg *CXArgument) []byte {
	size := GetSize(arg)
	return PROGRAM.Memory[offset : offset+size]
}

// ResizeMemory ...
func ResizeMemory(prgrm *CXProgram, newMemSize int, isExpand bool) {
	// We can't expand memory to a value greater than `memLimit`.
	if newMemSize > MAX_HEAP_SIZE {
		newMemSize = MAX_HEAP_SIZE
	}

	if newMemSize == prgrm.HeapSize {
		// Then we're at the limit; we can't expand anymore.
		// We can only hope that the free memory is enough for the CX program to continue running.
		return
	}

	if isExpand {
		// Adding bytes to reach a heap equal to `newMemSize`.
		prgrm.Memory = append(prgrm.Memory, make([]byte, newMemSize-prgrm.HeapSize)...)
		prgrm.HeapSize = newMemSize
	} else {
		// Removing bytes to reach a heap equal to `newMemSize`.
		prgrm.Memory = append([]byte(nil), prgrm.Memory[:prgrm.HeapStartsAt+newMemSize]...)
		prgrm.HeapSize = newMemSize
	}
}

// AllocateSeq allocates memory in the heap
func AllocateSeq(size int) (offset int) {
	// Current object trying to be allocated would use this address.
	addr := PROGRAM.HeapPointer
	// Next object to be allocated will use this address.
	newFree := addr + size

	// Checking if we can allocate the entirety of the object in the current heap.
	if newFree > PROGRAM.HeapSize {
		// It does not fit, so calling garbage collector.
		MarkAndCompact(PROGRAM)
		// Heap pointer got moved by GC and recalculate these variables based on the new pointer.
		addr = PROGRAM.HeapPointer
		newFree = addr + size

		// If the new heap pointer exceeds `MAX_HEAP_SIZE`, there's nothing left to do.
		if newFree > MAX_HEAP_SIZE {
			panic(HEAP_EXHAUSTED_ERROR)
		}

		// According to MIN_HEAP_FREE_RATIO and MAX_HEAP_FREE_RATION we can either shrink
		// or expand the heap to maintain "healthy" heap sizes. The idea is that we don't want
		// to have an absurdly amount of free heap memory, as we would be wasting resources, and we
		// don't want to have a small amount of heap memory left as we'd be calling the garbage collector
		// too frequently.

		// Calculating free heap memory percentage.
		usedPerc := float32(newFree) / float32(PROGRAM.HeapSize)
		freeMemPerc := 1.0 - usedPerc

		// Then we have less than MIN_HEAP_FREE_RATIO memory left. Expand!
		if freeMemPerc < MIN_HEAP_FREE_RATIO {
			// Calculating new heap size in order to reach MIN_HEAP_FREE_RATIO.
			newMemSize := int(float32(newFree) / (1.0 - MIN_HEAP_FREE_RATIO))
			ResizeMemory(PROGRAM, newMemSize, true)
		}

		// Then we have more than MAX_HEAP_FREE_RATIO memory left. Shrink!
		if freeMemPerc > MAX_HEAP_FREE_RATIO {
			// Calculating new heap size in order to reach MAX_HEAP_FREE_RATIO.
			newMemSize := int(float32(newFree) / (1.0 - MAX_HEAP_FREE_RATIO))

			// This check guarantees that the CX program has always at least INIT_HEAP_SIZE bytes to work with.
			// A flag could be added later to remove this, as in some cases this mechanism could not be desired.
			if newMemSize > INIT_HEAP_SIZE {
				ResizeMemory(PROGRAM, newMemSize, false)
			}
		}
	}

	PROGRAM.HeapPointer = newFree

	// Returning absolute memory address (not relative to where heap starts at).
	// Above this point we were performing all operations taking into
	// consideration only heap offsets.
	return addr + PROGRAM.HeapStartsAt
}

// WriteMemory ...
func WriteMemory(offset int, byts []byte) {
	for c := 0; c < len(byts); c++ {
		PROGRAM.Memory[offset+c] = byts[c]
	}
}

// Utilities

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

// FromStr ...
func FromStr(in string) []byte {
	return encoder.Serialize(in)
}

// FromBool ...
func FromBool(in bool) []byte {
	if in {
		return []byte{1}
	}
	return []byte{0}

}

// FromI8 ...
func FromI8(in int8) []byte {
	return encoder.SerializeAtomic(in)
}

// FromI16 ...
func FromI16(in int16) []byte {
	return encoder.SerializeAtomic(in)
}

// FromI32 ...
func FromI32(in int32) []byte {
	return encoder.SerializeAtomic(in)
}

// FromI64 ...
func FromI64(in int64) []byte {
	return encoder.SerializeAtomic(in)
}

// FromUI8 ...
func FromUI8(in uint8) []byte {
	return encoder.SerializeAtomic(in)
}

// FromUI16 ...
func FromUI16(in uint16) []byte {
	return encoder.SerializeAtomic(in)
}

// FromUI32 ...
func FromUI32(in uint32) []byte {
	return encoder.SerializeAtomic(in)
}

// FromUI64 ...
func FromUI64(in uint64) []byte {
	return encoder.SerializeAtomic(in)
}

// FromF32 ...
func FromF32(in float32) []byte {
	return FromUI32(math.Float32bits(in))
}

// FromF64 ...
func FromF64(in float64) []byte {
	return FromUI64(math.Float64bits(in))
}

// ReadData ...
func ReadData(fp int, inp *CXArgument, dataType int) interface{} {
	elt := GetAssignmentElement(inp)
	if elt.IsSlice {
		return ReadSlice(fp, inp, dataType)
	} else if elt.IsArray {
		return ReadArray(fp, inp, dataType)
	} else {
		return ReadObject(fp, inp, dataType)
	}
}

func readDataI8(bytes []byte) (out []int8) {
	count := len(bytes)
	if count > 0 {
		out = make([]int8, count)
		for i := 0; i < count; i++ {
			out[i] = mustDeserializeI8(bytes[i:])
		}
	}
	return
}

func readDataUI8(bytes []byte) (out []uint8) {
	count := len(bytes)
	if count > 0 {
		out = make([]uint8, count)
		for i := 0; i < count; i++ {
			out[i] = mustDeserializeUI8(bytes[i:])
		}
	}
	return
}

func readDataI16(bytes []byte) (out []int16) {
	count := len(bytes) / 2
	if count > 0 {
		out = make([]int16, count)
		for i := 0; i < count; i++ {
			out[i] = mustDeserializeI16(bytes[i*2:])
		}
	}
	return
}

func readDataUI16(bytes []byte) (out []uint16) {
	count := len(bytes) / 2
	if count > 0 {
		out = make([]uint16, count)
		for i := 0; i < count; i++ {
			out[i] = mustDeserializeUI16(bytes[i*2:])
		}
	}
	return
}

func readDataI32(bytes []byte) (out []int32) {
	count := len(bytes) / 4
	if count > 0 {
		out = make([]int32, count)
		for i := 0; i < count; i++ {
			out[i] = mustDeserializeI32(bytes[i*4:])
		}
	}
	return
}

func readDataUI32(bytes []byte) (out []uint32) {
	count := len(bytes) / 4
	if count > 0 {
		out = make([]uint32, count)
		for i := 0; i < count; i++ {
			out[i] = mustDeserializeUI32(bytes[i*4:])
		}
	}
	return
}

func readDataI64(bytes []byte) (out []int64) {
	count := len(bytes) / 8
	if count > 0 {
		out = make([]int64, count)
		for i := 0; i < count; i++ {
			out[i] = mustDeserializeI64(bytes[i*8:])
		}
	}
	return
}

func readDataUI64(bytes []byte) (out []uint64) {
	count := len(bytes) / 8
	if count > 0 {
		out = make([]uint64, count)
		for i := 0; i < count; i++ {
			out[i] = mustDeserializeUI64(bytes[i*8:])
		}
	}
	return
}

func readDataF32(bytes []byte) (out []float32) {
	count := len(bytes) / 4
	if count > 0 {
		out = make([]float32, count)
		for i := 0; i < count; i++ {
			out[i] = mustDeserializeF32(bytes[i*4:])
		}
	}
	return
}

func readDataF64(bytes []byte) (out []float64) {
	count := len(bytes) / 8
	if count > 0 {
		out = make([]float64, count)
		for i := 0; i < count; i++ {
			out[i] = mustDeserializeF64(bytes[i*8:])
		}
	}
	return
}

func readData(inp *CXArgument, bytes []byte) interface{} {
	switch inp.Type {
	case TYPE_I8:
		data := readDataI8(bytes)
		if len(data) > 0 {
			return interface{}(data)
		}
	case TYPE_I16:
		data := readDataI16(bytes)
		if len(data) > 0 {
			return interface{}(data)
		}
	case TYPE_I32:
		data := readDataI32(bytes)
		if len(data) > 0 {
			return interface{}(data)
		}
	case TYPE_I64:
		data := readDataI64(bytes)
		if len(data) > 0 {
			return interface{}(data)
		}
	case TYPE_UI8:
		data := readDataUI8(bytes)
		if len(data) > 0 {
			return interface{}(data)
		}
	case TYPE_UI16:
		data := readDataUI16(bytes)
		if len(data) > 0 {
			return interface{}(data)
		}
	case TYPE_UI32:
		data := readDataUI32(bytes)
		if len(data) > 0 {
			return interface{}(data)
		}
	case TYPE_UI64:
		data := readDataUI64(bytes)
		if len(data) > 0 {
			return interface{}(data)
		}
	case TYPE_F32:
		data := readDataF32(bytes)
		if len(data) > 0 {
			return interface{}(data)
		}
	case TYPE_F64:
		data := readDataF64(bytes)
		if len(data) > 0 {
			return interface{}(data)
		}
	default:
		data := readDataUI8(bytes)
		if len(data) > 0 {
			return interface{}(data)
		}
	}

	return interface{}(nil)
}

// ReadSlice ...
func ReadSlice(fp int, inp *CXArgument, dataType int) interface{} {
	sliceOffset := GetSliceOffset(fp, inp)
	if sliceOffset >= 0 && (dataType < 0 || inp.Type == dataType) {
		slice := GetSliceData(sliceOffset, GetAssignmentElement(inp).Size)
		if slice != nil {
			return readData(inp, slice)
		}
	} else {
		panic(CX_RUNTIME_INVALID_ARGUMENT)
	}

	return interface{}(nil)
}

// ReadSliceBytes ...
func ReadSliceBytes(fp int, inp *CXArgument, dataType int) []byte {
	sliceOffset := GetSliceOffset(fp, inp)
	if sliceOffset >= 0 && (dataType < 0 || inp.Type == dataType) {
		slice := GetSliceData(sliceOffset, GetAssignmentElement(inp).Size)
		return slice
	}

	panic(CX_RUNTIME_INVALID_ARGUMENT)
}

// ReadArray ...
func ReadArray(fp int, inp *CXArgument, dataType int) interface{} {
	offset := GetFinalOffset(fp, inp)
	if dataType < 0 || inp.Type == dataType {
		array := ReadMemory(offset, inp)
		return readData(inp, array)
	}
	panic(CX_RUNTIME_INVALID_ARGUMENT)
}

// ReadObject ...
func ReadObject(fp int, inp *CXArgument, dataType int) interface{} {
	offset := GetFinalOffset(fp, inp)
	array := ReadMemory(offset, inp)
	return readData(inp, array)
}

// ReadBool ...
func ReadBool(fp int, inp *CXArgument) (out bool) {
	offset := GetFinalOffset(fp, inp)
	out = mustDeserializeBool(ReadMemory(offset, inp))
	return
}

// ReadStr ...
func ReadStr(fp int, inp *CXArgument) (out string) {
	var offset int32
	off := GetFinalOffset(fp, inp)
	if inp.Name == "" {
		// Then it's a literal.
		offset = int32(off)
	} else {
		offset = mustDeserializeI32(PROGRAM.Memory[off : off+TYPE_POINTER_SIZE])
	}

	if offset == 0 {
		// Then it's nil string.
		out = ""
		return
	}

	// We need to check if the string lives on the data segment or on the
	// heap to know if we need to take into consideration the object header's size.
	if int(offset) > PROGRAM.HeapStartsAt {
		size := mustDeserializeI32(PROGRAM.Memory[offset+OBJECT_HEADER_SIZE : offset+OBJECT_HEADER_SIZE+STR_HEADER_SIZE])
		mustDeserializeRaw(PROGRAM.Memory[offset+OBJECT_HEADER_SIZE:offset+OBJECT_HEADER_SIZE+STR_HEADER_SIZE+size], &out)
	} else {
		size := mustDeserializeI32(PROGRAM.Memory[offset : offset+STR_HEADER_SIZE])
		mustDeserializeRaw(PROGRAM.Memory[offset:offset+STR_HEADER_SIZE+size], &out)
	}

	return out
}

// ReadI8 ...
func ReadI8(fp int, inp *CXArgument) int8 {
	return mustDeserializeI8(ReadMemory(GetFinalOffset(fp, inp), inp))
}

// ReadI16 ...
func ReadI16(fp int, inp *CXArgument) int16 {
	return mustDeserializeI16(ReadMemory(GetFinalOffset(fp, inp), inp))
}

// ReadI32 ...
func ReadI32(fp int, inp *CXArgument) int32 {
	return mustDeserializeI32(ReadMemory(GetFinalOffset(fp, inp), inp))
}

// ReadI64 ...
func ReadI64(fp int, inp *CXArgument) int64 {
	return mustDeserializeI64(ReadMemory(GetFinalOffset(fp, inp), inp))
}

// ReadUI8 ...
func ReadUI8(fp int, inp *CXArgument) uint8 {
	return mustDeserializeUI8(ReadMemory(GetFinalOffset(fp, inp), inp))
}

// ReadUI16 ...
func ReadUI16(fp int, inp *CXArgument) uint16 {
	return mustDeserializeUI16(ReadMemory(GetFinalOffset(fp, inp), inp))
}

// ReadUI32 ...
func ReadUI32(fp int, inp *CXArgument) uint32 {
	return mustDeserializeUI32(ReadMemory(GetFinalOffset(fp, inp), inp))
}

// ReadUI64 ...
func ReadUI64(fp int, inp *CXArgument) uint64 {
	return mustDeserializeUI64(ReadMemory(GetFinalOffset(fp, inp), inp))
}

// ReadF32 ...
func ReadF32(fp int, inp *CXArgument) float32 {
	return mustDeserializeF32(ReadMemory(GetFinalOffset(fp, inp), inp))
}

// ReadF64 ...
func ReadF64(fp int, inp *CXArgument) float64 {
	return mustDeserializeF64(ReadMemory(GetFinalOffset(fp, inp), inp))
}
