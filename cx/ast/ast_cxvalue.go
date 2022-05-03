package ast

import (
	"github.com/skycoin/cx/cx/types"
)

type CXValue struct {
	Arg          *CXArgument // TODO:PTR remove Arg
	Expr         *CXExpression
	Type         types.Code
	Offset       types.Pointer
	Size         types.Pointer
	FramePointer types.Pointer // TODO:PTR remove FramePointer
}

func (value *CXValue) Get_bool(prgrm *CXProgram) bool {
	return types.Read_bool(prgrm.Memory, value.Offset)
}

func (value *CXValue) Get_i8(prgrm *CXProgram) int8 {
	return types.Read_i8(prgrm.Memory, value.Offset)
}

func (value *CXValue) Get_i16(prgrm *CXProgram) int16 {
	return types.Read_i16(prgrm.Memory, value.Offset)
}

func (value *CXValue) Get_i32(prgrm *CXProgram) int32 {
	return types.Read_i32(prgrm.Memory, value.Offset)
}

func (value *CXValue) Get_i64(prgrm *CXProgram) int64 {
	return types.Read_i64(prgrm.Memory, value.Offset)
}

func (value *CXValue) Get_ui8(prgrm *CXProgram) uint8 {
	return types.Read_ui8(prgrm.Memory, value.Offset)
}

func (value *CXValue) Get_ui16(prgrm *CXProgram) uint16 {
	return types.Read_ui16(prgrm.Memory, value.Offset)
}

func (value *CXValue) Get_ui32(prgrm *CXProgram) uint32 {
	return types.Read_ui32(prgrm.Memory, value.Offset)
}

func (value *CXValue) Get_ui64(prgrm *CXProgram) uint64 {
	return types.Read_ui64(prgrm.Memory, value.Offset)
}

func (value *CXValue) Get_f32(prgrm *CXProgram) float32 {
	return types.Read_f32(prgrm.Memory, value.Offset)
}

func (value *CXValue) Get_f64(prgrm *CXProgram) float64 {
	return types.Read_f64(prgrm.Memory, value.Offset)
}

func (value *CXValue) Get_ptr(prgrm *CXProgram) types.Pointer {
	return types.Read_ptr(prgrm.Memory, value.Offset)
}

func (value *CXValue) Get_bytes(prgrm *CXProgram) []byte {
	return types.GetSlice_byte(prgrm.Memory, value.Offset, value.Size)
}

func (value *CXValue) Get_str(prgrm *CXProgram) string {
	return types.Read_str(prgrm.Memory, value.Offset)
}

func (value *CXValue) GetSlice_i8(prgrm *CXProgram) []int8 {
	if mem := GetSliceData(prgrm, types.Read_ptr(prgrm.Memory, value.Offset), value.Size); mem != nil {
		return types.ReadSlice_i8(mem, 0)
	}
	return nil
}

func (value *CXValue) GetSlice_i16(prgrm *CXProgram) []int16 {
	if mem := GetSliceData(prgrm, types.Read_ptr(prgrm.Memory, value.Offset), value.Size); mem != nil {
		return types.ReadSlice_i16(mem, 0)
	}
	return nil
}

func (value *CXValue) GetSlice_i32(prgrm *CXProgram) []int32 {
	if mem := GetSliceData(prgrm, types.Read_ptr(prgrm.Memory, value.Offset), value.Size); mem != nil {
		return types.ReadSlice_i32(mem, 0)
	}
	return nil
}

func (value *CXValue) GetSlice_i64(prgrm *CXProgram) []int64 {
	if mem := GetSliceData(prgrm, types.Read_ptr(prgrm.Memory, value.Offset), value.Size); mem != nil {
		return types.ReadSlice_i64(mem, 0)
	}
	return nil
}

func (value *CXValue) GetSlice_ui8(prgrm *CXProgram) []uint8 {
	if mem := GetSliceData(prgrm, types.Read_ptr(prgrm.Memory, value.Offset), value.Size); mem != nil {
		return types.ReadSlice_ui8(mem, 0)
	}
	return nil
}

func (value *CXValue) GetSlice_ui16(prgrm *CXProgram) []uint16 {
	if mem := GetSliceData(prgrm, types.Read_ptr(prgrm.Memory, value.Offset), value.Size); mem != nil {
		return types.ReadSlice_ui16(mem, 0)
	}
	return nil
}

func (value *CXValue) GetSlice_ui32(prgrm *CXProgram) []uint32 {
	if mem := GetSliceData(prgrm, types.Read_ptr(prgrm.Memory, value.Offset), value.Size); mem != nil {
		return types.ReadSlice_ui32(mem, 0)
	}
	return nil
}

func (value *CXValue) GetSlice_ui64(prgrm *CXProgram) []uint64 {
	if mem := GetSliceData(prgrm, types.Read_ptr(prgrm.Memory, value.Offset), value.Size); mem != nil {
		return types.ReadSlice_ui64(mem, 0)
	}
	return nil
}

func (value *CXValue) GetSlice_f32(prgrm *CXProgram) []float32 {
	if mem := GetSliceData(prgrm, types.Read_ptr(prgrm.Memory, value.Offset), value.Size); mem != nil {
		return types.ReadSlice_f32(mem, 0)
	}
	return nil
}

func (value *CXValue) GetSlice_f64(prgrm *CXProgram) []float64 {
	if mem := GetSliceData(prgrm, types.Read_ptr(prgrm.Memory, value.Offset), value.Size); mem != nil {
		return types.ReadSlice_f64(mem, 0)
	}
	return nil
}

func (value *CXValue) GetSlice_bytes(prgrm *CXProgram) []byte {
	return GetSliceData(prgrm, types.Read_ptr(prgrm.Memory, value.Offset), value.Size)
}

func (value *CXValue) Set_bool(prgrm *CXProgram, data bool) {
	types.Write_bool(prgrm.Memory, value.Offset, data)
}

func (value *CXValue) Set_i8(prgrm *CXProgram, data int8) {
	types.Write_i8(prgrm.Memory, value.Offset, data)
}

func (value *CXValue) Set_i16(prgrm *CXProgram, data int16) {
	types.Write_i16(prgrm.Memory, value.Offset, data)
}

func (value *CXValue) Set_i32(prgrm *CXProgram, data int32) {
	types.Write_i32(prgrm.Memory, value.Offset, data)
}

func (value *CXValue) Set_i64(prgrm *CXProgram, data int64) {
	types.Write_i64(prgrm.Memory, value.Offset, data)
}

func (value *CXValue) Set_ui8(prgrm *CXProgram, data uint8) {
	types.Write_ui8(prgrm.Memory, value.Offset, data)
}

func (value *CXValue) Set_ui16(prgrm *CXProgram, data uint16) {
	types.Write_ui16(prgrm.Memory, value.Offset, data)
}

func (value *CXValue) Set_ui32(prgrm *CXProgram, data uint32) {
	types.Write_ui32(prgrm.Memory, value.Offset, data)
}

func (value *CXValue) Set_ui64(prgrm *CXProgram, data uint64) {
	types.Write_ui64(prgrm.Memory, value.Offset, data)
}

func (value *CXValue) Set_f32(prgrm *CXProgram, data float32) {
	types.Write_f32(prgrm.Memory, value.Offset, data)
}

func (value *CXValue) Set_f64(prgrm *CXProgram, data float64) {
	types.Write_f64(prgrm.Memory, value.Offset, data)
}

func (value *CXValue) Set_ptr(prgrm *CXProgram, data types.Pointer) {
	types.Write_ptr(prgrm.Memory, value.Offset, data)
}

func (value *CXValue) Set_bytes(prgrm *CXProgram, data []byte) {
	types.WriteSlice_byte(prgrm.Memory, value.Offset, data)
}

func (value *CXValue) Set_str(prgrm *CXProgram, data string) {
	types.Write_str(prgrm, prgrm.Memory, value.Offset, data)
}
