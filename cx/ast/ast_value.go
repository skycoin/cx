package ast

import (
	"github.com/skycoin/cx/cx/constants"
	"github.com/skycoin/cx/cx/helper"
	"github.com/skycoin/skycoin/src/cipher/encoder"
)

type CXValue struct {
	Arg    *CXArgument
	Expr   *CXExpression
	Type   int
	memory []byte
	Offset int
	//size int. //unused field
	FramePointer int
}

// GetPointerOffset ...
func GetPointerOffset(pointer int32) int32 {
	return helper.Deserialize_i32(PROGRAM.Memory[pointer : pointer+constants.TYPE_POINTER_SIZE])
}


func (value *CXValue) GetSlice_i8() []int8 {
	////value.Used = TYPE_SLICE
	//value.Used = int8(value.Type) // TODO: type checking for slice is not working
	if mem := GetSliceData(GetPointerOffset(int32(value.Offset)), GetAssignmentElement(value.Arg).Size); mem != nil {
		return helper.ReadDataI8(mem)
	}
	return nil
}

func (value *CXValue) GetSlice_i16() []int16 {
	////value.Used = TYPE_SLICE
	//value.Used = int8(value.Type) // TODO: type checking for slice is not working
	if mem := GetSliceData(GetPointerOffset(int32(value.Offset)), GetAssignmentElement(value.Arg).Size); mem != nil {
		return helper.ReadDataI16(mem)
	}
	return nil
}

func (value *CXValue) GetSlice_i32() []int32 {
	////value.Used = TYPE_SLICE
	//value.Used = int8(value.Type) // TODO: type checking for slice is not working
	if mem := GetSliceData(GetPointerOffset(int32(value.Offset)), GetAssignmentElement(value.Arg).Size); mem != nil {
		return helper.ReadDataI32(mem)
	}
	return nil
}

func (value *CXValue) GetSlice_i64() []int64 {
	////value.Used = TYPE_SLICE
	//value.Used = int8(value.Type) // TODO: type checking for slice is not working
	if mem := GetSliceData(GetPointerOffset(int32(value.Offset)), GetAssignmentElement(value.Arg).Size); mem != nil {
		return helper.ReadDataI64(mem)
	}
	return nil
}

func (value *CXValue) GetSlice_ui8() []uint8 {
	////value.Used = TYPE_SLICE
	//value.Used = int8(value.Type) // TODO: type checking for slice is not working
	if mem := GetSliceData(GetPointerOffset(int32(value.Offset)), GetAssignmentElement(value.Arg).Size); mem != nil {
		return helper.ReadDataUI8(mem)
	}
	return nil
}

func (value *CXValue) GetSlice_ui16() []uint16 {
	////value.Used = TYPE_SLICE
	//value.Used = int8(value.Type) // TODO: type checking for slice is not working
	if mem := GetSliceData(GetPointerOffset(int32(value.Offset)), GetAssignmentElement(value.Arg).Size); mem != nil {
		return helper.ReadDataUI16(mem)
	}
	return nil
}

func (value *CXValue) GetSlice_ui32() []uint32 {
	////value.Used = TYPE_SLICE
	//value.Used = int8(value.Type) // TODO: type checking for slice is not working
	if mem := GetSliceData(GetPointerOffset(int32(value.Offset)), GetAssignmentElement(value.Arg).Size); mem != nil {
		return helper.ReadDataUI32(mem)
	}
	return nil
}

func (value *CXValue) GetSlice_ui64() []uint64 {
	////value.Used = TYPE_SLICE
	//value.Used = int8(value.Type) // TODO: type checking for slice is not working
	if mem := GetSliceData(GetPointerOffset(int32(value.Offset)), GetAssignmentElement(value.Arg).Size); mem != nil {
		return helper.ReadDataUI64(mem)
	}
	return nil
}

func (value *CXValue) GetSlice_f32() []float32 {
	////value.Used = TYPE_SLICE
	//value.Used = int8(value.Type) // TODO: type checking for slice is not working
	if mem := GetSliceData(GetPointerOffset(int32(value.Offset)), GetAssignmentElement(value.Arg).Size); mem != nil {
		return helper.ReadDataF32(mem)
	}
	return nil
}

func (value *CXValue) GetSlice_f64() []float64 {
	////value.Used = TYPE_SLICE
	//value.Used = int8(value.Type) // TODO: type checking for slice is not working
	if mem := GetSliceData(GetPointerOffset(int32(value.Offset)), GetAssignmentElement(value.Arg).Size); mem != nil {
		return helper.ReadDataF64(mem)
	}
	return nil
}

func (value *CXValue) SetSlice(data int32) {
	//value.Used = int8(value.Type) // TODO: type checking for slice is not working
	WriteI32(value.Offset, data)
}

func (value *CXValue) Get_bytes() []byte {
	////value.Used = TYPE_SLICE
	//value.Used = int8(value.Type) // TODO: type checking for slice is not working
	return ReadMemory(value.Offset, value.Arg)
}

func (value *CXValue) Set_bytes(data []byte) () {
	//value.Used = constants.TYPE_CUSTOM
	WriteMemory(value.Offset, data)
}

func (value *CXValue) GetSlice_bytes() []byte {
	//value.Used = int8(value.Type) // TODO: type checking for slice is not working
	return GetSliceData(GetPointerOffset(int32(value.Offset)), GetAssignmentElement(value.Arg).Size)
}

func (value *CXValue) Get_i8() int8 {
	//value.Used = constants.TYPE_I8
	return helper.Deserialize_i8(value.memory)
}

func (value *CXValue) Set_i8(data int8) {
	//value.Used = constants.TYPE_I8
	WriteI8(value.Offset, data)
}

func (value *CXValue) Get_i16() int16 {
	//value.Used = constants.TYPE_I16
	return helper.Deserialize_i16(value.memory)
}

func (value *CXValue) Set_i16(data int16) {
	//value.Used = constants.TYPE_I16
	WriteI16(value.Offset, data)
}

func (value *CXValue) Get_i32() int32 {
	//value.Used = constants.TYPE_I32
	return helper.Deserialize_i32(value.memory)
}

func (value *CXValue) Set_i32(data int32) {
	//value.Used = constants.TYPE_I32
	WriteI32(value.Offset, data)
}

func (value *CXValue) Get_i64() int64 {
	//value.Used = constants.TYPE_I64
	return helper.Deserialize_i64(value.memory)
}

func (value *CXValue) Set_i64(data int64) {
	//value.Used = constants.TYPE_I64
	WriteI64(value.Offset, data)
}

func (value *CXValue) Get_ui8() uint8 {
	//value.Used = constants.TYPE_UI8
	return helper.Deserialize_ui8(value.memory)
}

func (value *CXValue) Set_ui8(data uint8) {
	//value.Used = constants.TYPE_UI8
	WriteUI8(value.Offset, data)
}

func (value *CXValue) Get_ui16() uint16 {
	//value.Used = constants.TYPE_UI16
	return helper.Deserialize_ui16(value.memory)
}

func (value *CXValue) Set_ui16(data uint16) {
	//value.Used = constants.TYPE_UI16
	WriteUI16(value.Offset, data)
}

func (value *CXValue) Get_ui32() uint32 {
	//value.Used = constants.TYPE_UI32
	return helper.Deserialize_ui32(value.memory)
}

func (value *CXValue) Set_ui32(data uint32) {
	//value.Used = constants.TYPE_UI32
	WriteUI32(value.Offset, data)
}

func (value *CXValue) Get_ui64() uint64 {
	//value.Used = constants.TYPE_UI64
	return helper.Deserialize_ui64(value.memory)
}

func (value *CXValue) Set_ui64(data uint64) {
	//value.Used = constants.TYPE_UI64
	WriteUI64(value.Offset, data)
}

func (value *CXValue) Get_f32() float32 {
	//value.Used = constants.TYPE_F32
	return helper.Deserialize_f32(value.memory)
}

func (value *CXValue) Set_f32(data float32) {
	//value.Used = constants.TYPE_F32
	WriteF32(value.Offset, data)
}

func (value *CXValue) Get_f64() float64 {
	//value.Used = constants.TYPE_F64
	return helper.Deserialize_f64(value.memory)
}

func (value *CXValue) Set_f64(data float64) {
	//value.Used = constants.TYPE_F64
	WriteF64(value.Offset, data)
}

func (value *CXValue) Get_bool() bool {
	//value.Used = constants.TYPE_BOOL
	return helper.Deserialize_bool(value.memory)
}

func (value *CXValue) Set_bool(data bool) {
	//value.Used = constants.TYPE_BOOL
	WriteBool(value.Offset, data)
}

func (value *CXValue) Get_str() string {
	//value.Used = constants.TYPE_STR
	return ReadStrFromOffset(value.Offset, value.Arg)
}

func (value *CXValue) Set_str(data string) {
	//value.Used = constants.TYPE_STR
	WriteObject(value.Offset, encoder.Serialize(data))
}
