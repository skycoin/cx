package types

import (
	"math"
	//"fmt"
)

const (
	CX_SUCCESS = iota //zero can be success
	CX_COMPILATION_ERROR
	CX_PANIC // 2
	CX_INTERNAL_ERROR
	CX_ASSERT
	CX_RUNTIME_ERROR
	CX_RUNTIME_STACK_OVERFLOW_ERROR
	CX_RUNTIME_HEAP_EXHAUSTED_ERROR
	CX_RUNTIME_INVALID_ARGUMENT
	CX_RUNTIME_SLICE_INDEX_OUT_OF_RANGE
	CX_RUNTIME_NOT_IMPLEMENTED
	CX_RUNTIME_INVALID_CAST
)

const (
	MAX_INT   = int(MAX_UINT >> 1)
	MAX_INT32 = int(MAX_UINT32 >> 1)
	MIN_INT32 = -MAX_INT32 - 1

	MAX_UINT   = ^uint(0)
	MAX_UINT8  = ^uint8(0)
	MAX_UINT16 = ^uint16(0)
	MAX_UINT32 = ^uint32(0)
	MAX_UINT64 = ^uint64(0)
)

const (
	BOOL_SIZE = Pointer(1)

	I8_SIZE  = Pointer(1)
	I16_SIZE = Pointer(2)
	I32_SIZE = Pointer(4)
	I64_SIZE = Pointer(8)

	UI8_SIZE  = Pointer(1)
	UI16_SIZE = Pointer(2)
	UI32_SIZE = Pointer(4)
	UI64_SIZE = Pointer(8)

	F32_SIZE = Pointer(4)
	F64_SIZE = Pointer(8)

	ARRAY_SIZE = POINTER_SIZE
	SLICE_SIZE = POINTER_SIZE

	STR_SIZE = POINTER_SIZE
)

type Code int

const (
	UNUSED Code = iota //reserve zero value, if this value appears, program should crash; should be assert

	BOOL

	I8
	I16
	I32
	I64

	UI8
	UI16
	UI32
	UI64

	F32
	F64

	POINTER
	STR

	ARRAY
	SLICE
	STRUCT

	FUNC
	AFF

	UNDEFINED
	IDENTIFIER

	COUNT
)

func (t Code) Name() string {
	return definitions[t].name
}

func (t Code) Size() Pointer {
	//	panicIf(!definitions[t].size.IsValid(), CX_RUNTIME_INVALID_ARGUMENT)
	return definitions[t].size
}

func (t Code) IsPrimitive() bool {
	return definitions[t].isPrimitive
}

type Type struct {
	name        string
	size        Pointer
	isPrimitive bool
}

var definitions []Type = []Type{
	{"UNUSED", InvalidPointer, false},

	{"bool", BOOL_SIZE, true},

	{"i8", I8_SIZE, true},
	{"i16", I16_SIZE, true},
	{"i32", I32_SIZE, true},
	{"i64", I64_SIZE, true},

	{"ui8", UI8_SIZE, true},
	{"ui16", UI16_SIZE, true},
	{"ui32", UI32_SIZE, true},
	{"ui64", UI64_SIZE, true},

	{"f32", F32_SIZE, true},
	{"f64", F64_SIZE, true},

	{"ptr", POINTER_SIZE, false},
	{"str", STR_SIZE, true}, // TODO:PTR check why str needs to be a primitive type or we need to have both isPrimitive && isAtomic.
	{"array", ARRAY_SIZE, false},
	{"slice", SLICE_SIZE, false},
	{"struct", InvalidPointer, false},
	{"func", InvalidPointer, false},
	{"aff", InvalidPointer, false},
	{"und", InvalidPointer, false},
	{"ident", POINTER_SIZE, false}, // TODO:PTR use InvalidPointer to track addressing issues.
}

type AllocatorHandler func(Pointer) Pointer

var Allocator AllocatorHandler

func debugPanicIf(condition bool, error int) {
	if condition {
		panic(error)
	}
}

func panicIf(condition bool, error int) {
	if condition {
		panic(error)
	}
}

func (pointer Pointer) IsValid() bool {
	return pointer != InvalidPointer
}

func (pointer *Pointer) Add(value Pointer) bool {
	if pointer.IsValid() && value.IsValid() {
		*pointer += value
		return true
	}

	*pointer = InvalidPointer
	return false
}

func Cast_sint_to_sptr(value []int) []Pointer {
	l := len(value)
	if l == 0 {
		return nil
	}

	sptr := make([]Pointer, l)
	for i, k := range value {
		sptr[i] = Cast_int_to_ptr(k)
	}
	return sptr
}

func Read_bool(memory []byte, offset Pointer) bool {
	memory = memory[offset:]
	return memory[0] != 0
}

func Read_i8(memory []byte, offset Pointer) int8 {
	memory = memory[offset:]
	return int8(memory[0])
}

func Read_i16(memory []byte, offset Pointer) int16 {
	memory = memory[offset:]
	return int16(memory[0]) | int16(memory[1])<<8
}

func Read_i32(memory []byte, offset Pointer) int32 {
	memory = memory[offset:]
	return int32(memory[0]) | int32(memory[1])<<8 | int32(memory[2])<<16 | int32(memory[3])<<24
}

func Read_i64(memory []byte, offset Pointer) int64 {
	memory = memory[offset:]
	return int64(memory[0]) | int64(memory[1])<<8 | int64(memory[2])<<16 | int64(memory[3])<<24 |
		int64(memory[4])<<32 | int64(memory[5])<<40 | int64(memory[6])<<48 | int64(memory[7])<<56
}

func Read_ui8(memory []byte, offset Pointer) uint8 {
	memory = memory[offset:]
	return uint8(memory[0])
}

func Read_ui16(memory []byte, offset Pointer) uint16 {
	memory = memory[offset:]
	return uint16(memory[0]) | uint16(memory[1])<<8
}

func Read_ui32(memory []byte, offset Pointer) uint32 {
	memory = memory[offset:]
	return uint32(memory[0]) | uint32(memory[1])<<8 | uint32(memory[2])<<16 | uint32(memory[3])<<24
}

func Read_ui64(memory []byte, offset Pointer) uint64 {
	memory = memory[offset:]
	return uint64(memory[0]) | uint64(memory[1])<<8 | uint64(memory[2])<<16 | uint64(memory[3])<<24 |
		uint64(memory[4])<<32 | uint64(memory[5])<<40 | uint64(memory[6])<<48 | uint64(memory[7])<<56
}

func Read_f32(memory []byte, offset Pointer) float32 {
	return math.Float32frombits(Read_ui32(memory, offset))
}

func Read_f64(memory []byte, offset Pointer) float64 {
	return math.Float64frombits(Read_ui64(memory, offset))
}

func ReadSlice_i8(memory []byte, offset Pointer) (out []int8) {
	count := Cast_int_to_ptr(len(memory))
	if count > 0 {
		out = make([]int8, count)
		for i := Pointer(0); i < count; i++ {
			out[i] = Read_i8(memory, i)
		}
	}
	return
}

func ReadSlice_i16(memory []byte, offset Pointer) (out []int16) {
	count := Cast_int_to_ptr(len(memory) / 2)
	if count > 0 {
		out = make([]int16, count)
		for i := Pointer(0); i < count; i++ {
			out[i] = Read_i16(memory, i*2)
		}
	}
	return
}

func ReadSlice_i32(memory []byte, offset Pointer) (out []int32) {
	count := Cast_int_to_ptr(len(memory) / 4)
	if count > 0 {
		out = make([]int32, count)
		for i := Pointer(0); i < count; i++ {
			out[i] = Read_i32(memory, i*4)
		}
	}
	return
}

func ReadSlice_i64(memory []byte, offset Pointer) (out []int64) {
	count := Cast_int_to_ptr(len(memory) / 8)
	if count > 0 {
		out = make([]int64, count)
		for i := Pointer(0); i < count; i++ {
			out[i] = Read_i64(memory, i*8)
		}
	}
	return
}

func ReadSlice_ui8(memory []byte, offset Pointer) (out []uint8) {
	count := Cast_int_to_ptr(len(memory))
	if count > 0 {
		out = make([]uint8, count)
		for i := Pointer(0); i < count; i++ {
			out[i] = Read_ui8(memory, i)
		}
	}
	return
}

func ReadSlice_ui16(memory []byte, offset Pointer) (out []uint16) {
	count := Cast_int_to_ptr(len(memory) / 2)
	if count > 0 {
		out = make([]uint16, count)
		for i := Pointer(0); i < count; i++ {
			out[i] = Read_ui16(memory, i*2)
		}
	}
	return
}

func ReadSlice_ui32(memory []byte, offset Pointer) (out []uint32) {
	count := Cast_int_to_ptr(len(memory) / 4)
	if count > 0 {
		out = make([]uint32, count)
		for i := Pointer(0); i < count; i++ {
			out[i] = Read_ui32(memory, i*4)
		}
	}
	return
}

func ReadSlice_ui64(memory []byte, offset Pointer) (out []uint64) {
	count := Cast_int_to_ptr(len(memory) / 8)
	if count > 0 {
		out = make([]uint64, count)
		for i := Pointer(0); i < count; i++ {
			out[i] = Read_ui64(memory, i*8)
		}
	}
	return
}

func ReadSlice_f32(memory []byte, offset Pointer) (out []float32) {
	count := Cast_int_to_ptr(len(memory) / 4)
	if count > 0 {
		out = make([]float32, count)
		for i := Pointer(0); i < count; i++ {
			out[i] = Read_f32(memory, i*4)
		}
	}
	return
}

func ReadSlice_f64(memory []byte, offset Pointer) (out []float64) {
	count := Cast_int_to_ptr(len(memory) / 8)
	if count > 0 {
		out = make([]float64, count)
		for i := Pointer(0); i < count; i++ {
			out[i] = Read_f64(memory, i*8)
		}
	}
	return
}

func GetSlice_byte(memory []byte, offset Pointer, size Pointer) []byte {
	return memory[offset : offset+size]
}

func Write_bool(memory []byte, offset Pointer, value bool) {
	if value {
		memory[offset] = 1
	} else {
		memory[offset] = 0
	}
}

func Write_i8(mem []byte, offset Pointer, v int8) {
	mem[offset] = byte(v)
}

func Write_i16(mem []byte, offset Pointer, v int16) {
	mem[offset] = byte(v)
	mem[offset+1] = byte(v >> 8)
}

func Write_i32(mem []byte, offset Pointer, v int32) {
	mem[offset] = byte(v)
	mem[offset+1] = byte(v >> 8)
	mem[offset+2] = byte(v >> 16)
	mem[offset+3] = byte(v >> 24)
}

func Write_i64(mem []byte, offset Pointer, v int64) {
	mem[offset] = byte(v)
	mem[offset+1] = byte(v >> 8)
	mem[offset+2] = byte(v >> 16)
	mem[offset+3] = byte(v >> 24)
	mem[offset+4] = byte(v >> 32)
	mem[offset+5] = byte(v >> 40)
	mem[offset+6] = byte(v >> 48)
	mem[offset+7] = byte(v >> 56)
}

func Write_ui8(mem []byte, offset Pointer, v uint8) {
	mem[offset] = v
}

func Write_ui16(mem []byte, offset Pointer, v uint16) {
	mem[offset] = byte(v)
	mem[offset+1] = byte(v >> 8)
}

func Write_ui32(memory []byte, offset Pointer, value uint32) {
	memory[offset] = byte(value)
	memory[offset+1] = byte(value >> 8)
	memory[offset+2] = byte(value >> 16)
	memory[offset+3] = byte(value >> 24)
}

func Write_ui64(memory []byte, offset Pointer, value uint64) {
	memory[offset] = byte(value)
	memory[offset+1] = byte(value >> 8)
	memory[offset+2] = byte(value >> 16)
	memory[offset+3] = byte(value >> 24)
	memory[offset+4] = byte(value >> 32)
	memory[offset+5] = byte(value >> 40)
	memory[offset+6] = byte(value >> 48)
	memory[offset+7] = byte(value >> 56)
}

func Write_f32(mem []byte, offset Pointer, f float32) {
	v := math.Float32bits(f)
	mem[offset] = byte(v)
	mem[offset+1] = byte(v >> 8)
	mem[offset+2] = byte(v >> 16)
	mem[offset+3] = byte(v >> 24)
}

func Write_f64(mem []byte, offset Pointer, f float64) {
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

func WriteSlice_byte(memory []byte, offset Pointer, byts []byte) {
	// TODO:PTR use copy()
	count := Cast_int_to_ptr(len(byts))
	for c := Pointer(0); c < count; c++ {
		memory[offset+c] = byts[c]
	}
}
