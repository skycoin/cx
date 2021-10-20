package ast

import (
	"github.com/skycoin/cx/cx/types"
)

type CXValue struct {
	Arg          *CXArgument // TODO:PTR remove Arg
	Expr         *CXExpression
	Type         types.Code
	Offset       types.Pointer
	FramePointer types.Pointer // TODO:PTR remove FramePointer
}

func (value *CXValue) Get_bool() bool {
	return types.Read_bool(PROGRAM.Memory, value.Offset)
}

func (value *CXValue) Get_i8() int8 {
	return types.Read_i8(PROGRAM.Memory, value.Offset)
}

func (value *CXValue) Get_i16() int16 {
	return types.Read_i16(PROGRAM.Memory, value.Offset)
}

func (value *CXValue) Get_i32() int32 {
	return types.Read_i32(PROGRAM.Memory, value.Offset)
}

func (value *CXValue) Get_i64() int64 {
	return types.Read_i64(PROGRAM.Memory, value.Offset)
}

func (value *CXValue) Get_ui8() uint8 {
	return types.Read_ui8(PROGRAM.Memory, value.Offset)
}

func (value *CXValue) Get_ui16() uint16 {
	return types.Read_ui16(PROGRAM.Memory, value.Offset)
}

func (value *CXValue) Get_ui32() uint32 {
	return types.Read_ui32(PROGRAM.Memory, value.Offset)
}

func (value *CXValue) Get_ui64() uint64 {
	return types.Read_ui64(PROGRAM.Memory, value.Offset)
}

func (value *CXValue) Get_f32() float32 {
	return types.Read_f32(PROGRAM.Memory, value.Offset)
}

func (value *CXValue) Get_f64() float64 {
	return types.Read_f64(PROGRAM.Memory, value.Offset)
}

func (value *CXValue) Get_ptr() types.Pointer {
	return types.Read_ptr(PROGRAM.Memory, value.Offset)
}

func (value *CXValue) Get_bytes() []byte {
	return types.GetSlice_byte(PROGRAM.Memory, value.Offset, GetSize(value.Arg))
}

func (value *CXValue) Get_str() string {
	return types.Read_str(PROGRAM.Memory, value.Offset)
}

func (value *CXValue) GetSlice_i8() []int8 {
	if mem := GetSliceData(types.Read_ptr(PROGRAM.Memory, value.Offset), value.Arg.GetAssignmentElement().Size); mem != nil {
		return types.ReadSlice_i8(mem, 0)
	}
	return nil
}

func (value *CXValue) GetSlice_i16() []int16 {
	if mem := GetSliceData(types.Read_ptr(PROGRAM.Memory, value.Offset), value.Arg.GetAssignmentElement().Size); mem != nil {
		return types.ReadSlice_i16(mem, 0)
	}
	return nil
}

func (value *CXValue) GetSlice_i32() []int32 {
	if mem := GetSliceData(types.Read_ptr(PROGRAM.Memory, value.Offset), value.Arg.GetAssignmentElement().Size); mem != nil {
		return types.ReadSlice_i32(mem, 0)
	}
	return nil
}

func (value *CXValue) GetSlice_i64() []int64 {
	if mem := GetSliceData(types.Read_ptr(PROGRAM.Memory, value.Offset), value.Arg.GetAssignmentElement().Size); mem != nil {
		return types.ReadSlice_i64(mem, 0)
	}
	return nil
}

func (value *CXValue) GetSlice_ui8() []uint8 {
	if mem := GetSliceData(types.Read_ptr(PROGRAM.Memory, value.Offset), value.Arg.GetAssignmentElement().Size); mem != nil {
		return types.ReadSlice_ui8(mem, 0)
	}
	return nil
}

func (value *CXValue) GetSlice_ui16() []uint16 {
	if mem := GetSliceData(types.Read_ptr(PROGRAM.Memory, value.Offset), value.Arg.GetAssignmentElement().Size); mem != nil {
		return types.ReadSlice_ui16(mem, 0)
	}
	return nil
}

func (value *CXValue) GetSlice_ui32() []uint32 {
	if mem := GetSliceData(types.Read_ptr(PROGRAM.Memory, value.Offset), value.Arg.GetAssignmentElement().Size); mem != nil {
		return types.ReadSlice_ui32(mem, 0)
	}
	return nil
}

func (value *CXValue) GetSlice_ui64() []uint64 {
	if mem := GetSliceData(types.Read_ptr(PROGRAM.Memory, value.Offset), value.Arg.GetAssignmentElement().Size); mem != nil {
		return types.ReadSlice_ui64(mem, 0)
	}
	return nil
}

func (value *CXValue) GetSlice_f32() []float32 {
	if mem := GetSliceData(types.Read_ptr(PROGRAM.Memory, value.Offset), value.Arg.GetAssignmentElement().Size); mem != nil {
		return types.ReadSlice_f32(mem, 0)
	}
	return nil
}

func (value *CXValue) GetSlice_f64() []float64 {
	if mem := GetSliceData(types.Read_ptr(PROGRAM.Memory, value.Offset), value.Arg.GetAssignmentElement().Size); mem != nil {
		return types.ReadSlice_f64(mem, 0)
	}
	return nil
}

func (value *CXValue) GetSlice_bytes() []byte {
	return GetSliceData(types.Read_ptr(PROGRAM.Memory, value.Offset), GetSize(value.Arg))
}

func (value *CXValue) Set_bool(data bool) {
	types.Write_bool(PROGRAM.Memory, value.Offset, data)
}

func (value *CXValue) Set_i8(data int8) {
	types.Write_i8(PROGRAM.Memory, value.Offset, data)
}

func (value *CXValue) Set_i16(data int16) {
	types.Write_i16(PROGRAM.Memory, value.Offset, data)
}

func (value *CXValue) Set_i32(data int32) {
	types.Write_i32(PROGRAM.Memory, value.Offset, data)
}

func (value *CXValue) Set_i64(data int64) {
	types.Write_i64(PROGRAM.Memory, value.Offset, data)
}

func (value *CXValue) Set_ui8(data uint8) {
	types.Write_ui8(PROGRAM.Memory, value.Offset, data)
}

func (value *CXValue) Set_ui16(data uint16) {
	types.Write_ui16(PROGRAM.Memory, value.Offset, data)
}

func (value *CXValue) Set_ui32(data uint32) {
	types.Write_ui32(PROGRAM.Memory, value.Offset, data)
}

func (value *CXValue) Set_ui64(data uint64) {
	types.Write_ui64(PROGRAM.Memory, value.Offset, data)
}

func (value *CXValue) Set_f32(data float32) {
	types.Write_f32(PROGRAM.Memory, value.Offset, data)
}

func (value *CXValue) Set_f64(data float64) {
	types.Write_f64(PROGRAM.Memory, value.Offset, data)
}

func (value *CXValue) Set_ptr(data types.Pointer) {
	types.Write_ptr(PROGRAM.Memory, value.Offset, data)
}

func (value *CXValue) Set_bytes(data []byte) {
	types.WriteSlice_byte(PROGRAM.Memory, value.Offset, data)
}

func (value *CXValue) Set_str(data string) {
	types.Write_str(PROGRAM.Memory, value.Offset, data)
}
