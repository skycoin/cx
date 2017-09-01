package base

import (
	"fmt"
	"time"
	"errors"
	"math/rand"
	"github.com/skycoin/skycoin/src/cipher/encoder"
)

func checkTwoTypes (fnName string, typ string, arg1 *CXArgument, arg2 *CXArgument) error {
	if arg1.Typ.Name != typ || arg2.Typ.Name != typ {
		if arg1.Typ.Name != typ {
			return errors.New(fmt.Sprintf("%s: first argument is type '%s'; expected type '%s'", fnName, arg1.Typ.Name, typ))
		}
		return errors.New(fmt.Sprintf("%s: second argument is type '%s'; expected type '%s'", fnName, arg1.Typ.Name, typ))
	}
	return nil
}

func addI32 (arg1 *CXArgument, arg2 *CXArgument) (*CXArgument, error) {
	if err := checkTwoTypes("addI32", "i32", arg1, arg2); err == nil {
		var num1 int32
		var num2 int32
		encoder.DeserializeRaw(*arg1.Value, &num1)
		encoder.DeserializeRaw(*arg2.Value, &num2)

		output := encoder.SerializeAtomic(int32(num1 + num2))

		return &CXArgument{Value: &output, Typ: MakeType("i32")}, nil
	} else {
		return nil, err
	}
}

func subI32 (arg1 *CXArgument, arg2 *CXArgument) (*CXArgument, error) {
	if err := checkTwoTypes("subI32", "i32", arg1, arg2); err == nil {
		var num1 int32
		var num2 int32
		encoder.DeserializeRaw(*arg1.Value, &num1)
		encoder.DeserializeRaw(*arg2.Value, &num2)

		output := encoder.SerializeAtomic(int32(num1 - num2))

		return &CXArgument{Value: &output, Typ: MakeType("i32")}, nil
	} else {
		return nil, err
	}
}

func mulI32 (arg1 *CXArgument, arg2 *CXArgument) (*CXArgument, error) {
	if err := checkTwoTypes("mulI32", "i32", arg1, arg2); err == nil {
		var num1 int32
		var num2 int32
		encoder.DeserializeRaw(*arg1.Value, &num1)
		encoder.DeserializeRaw(*arg2.Value, &num2)

		output := encoder.SerializeAtomic(int32(num1 * num2))

		return &CXArgument{Value: &output, Typ: MakeType("i32")}, nil
	} else {
		return nil, err
	}
}

func divI32 (arg1 *CXArgument, arg2 *CXArgument) (*CXArgument, error) {
	if err := checkTwoTypes("divI32", "i32", arg1, arg2); err == nil {
		var num1 int32
		var num2 int32
		encoder.DeserializeRaw(*arg1.Value, &num1)
		encoder.DeserializeRaw(*arg2.Value, &num2)

		if num2 == int32(0) {
			return nil, errors.New("divI32: Division by 0")
		}

		output := encoder.SerializeAtomic(int32(num1 / num2))

		return &CXArgument{Value: &output, Typ: MakeType("i32")}, nil
	} else {
		return nil, err
	}
}

func addI64 (arg1 *CXArgument, arg2 *CXArgument) (*CXArgument, error) {
	if err := checkTwoTypes("addI64", "i64", arg1, arg2); err == nil {
		var num1 int64
		var num2 int64
		encoder.DeserializeRaw(*arg1.Value, &num1)
		encoder.DeserializeRaw(*arg2.Value, &num2)

		output := encoder.SerializeAtomic(int64(num1 + num2))

		return &CXArgument{Value: &output, Typ: MakeType("i64")}, nil
	} else {
		return nil, err
	}
}

func subI64 (arg1 *CXArgument, arg2 *CXArgument) (*CXArgument, error) {
	if err := checkTwoTypes("subI64", "i64", arg1, arg2); err == nil {
		var num1 int64
		var num2 int64
		encoder.DeserializeRaw(*arg1.Value, &num1)
		encoder.DeserializeRaw(*arg2.Value, &num2)

		output := encoder.SerializeAtomic(int64(num1 - num2))

		return &CXArgument{Value: &output, Typ: MakeType("i64")}, nil
	} else {
		return nil, err
	}
}

func mulI64 (arg1 *CXArgument, arg2 *CXArgument) (*CXArgument, error) {
	if err := checkTwoTypes("mulI64", "i64", arg1, arg2); err == nil {
		var num1 int64
		var num2 int64
		encoder.DeserializeRaw(*arg1.Value, &num1)
		encoder.DeserializeRaw(*arg2.Value, &num2)

		output := encoder.SerializeAtomic(int64(num1 * num2))

		return &CXArgument{Value: &output, Typ: MakeType("i64")}, nil
	} else {
		return nil, err
	}
}

func divI64 (arg1 *CXArgument, arg2 *CXArgument) (*CXArgument, error) {
	if err := checkTwoTypes("divI64", "i64", arg1, arg2); err == nil {
		var num1 int64
		var num2 int64
		encoder.DeserializeRaw(*arg1.Value, &num1)
		encoder.DeserializeRaw(*arg2.Value, &num2)

		if num2 == int64(0) {
			return nil, errors.New("divI64: Division by 0")
		}
		
		output := encoder.SerializeAtomic(int64(num1 / num2))

		return &CXArgument{Value: &output, Typ: MakeType("i64")}, nil
	} else {
		return nil, err
	}
}

func addF32 (arg1 *CXArgument, arg2 *CXArgument) (*CXArgument, error) {
	if err := checkTwoTypes("addF32", "f32", arg1, arg2); err == nil {
		
	} else {
		return nil, err
	}
	
	var num1 float32
	var num2 float32
	encoder.DeserializeRaw(*arg1.Value, &num1)
	encoder.DeserializeRaw(*arg2.Value, &num2)

	output := encoder.Serialize(float32(num1 + num2))

	return &CXArgument{Value: &output, Typ: MakeType("f32")}, nil
}

func subF32 (arg1 *CXArgument, arg2 *CXArgument) (*CXArgument, error) {
	if err := checkTwoTypes("subF32", "f32", arg1, arg2); err == nil {
		var num1 float32
	var num2 float32
	encoder.DeserializeRaw(*arg1.Value, &num1)
	encoder.DeserializeRaw(*arg2.Value, &num2)

	output := encoder.Serialize(float32(num1 - num2))

	return &CXArgument{Value: &output, Typ: MakeType("f32")}, nil
	} else {
		return nil, err
	}
}

func mulF32 (arg1 *CXArgument, arg2 *CXArgument) (*CXArgument, error) {
	if err := checkTwoTypes("mulF32", "f32", arg1, arg2); err == nil {
		var num1 float32
		var num2 float32
		encoder.DeserializeRaw(*arg1.Value, &num1)
		encoder.DeserializeRaw(*arg2.Value, &num2)

		output := encoder.Serialize(float32(num1 * num2))

		return &CXArgument{Value: &output, Typ: MakeType("f32")}, nil
	} else {
		return nil, err
	}
}

func divF32 (arg1 *CXArgument, arg2 *CXArgument) (*CXArgument, error) {
	if err := checkTwoTypes("divF32", "f32", arg1, arg2); err == nil {
		var num1 float32
		var num2 float32
		encoder.DeserializeRaw(*arg1.Value, &num1)
		encoder.DeserializeRaw(*arg2.Value, &num2)

		if num2 == float32(0.0) {
			panic("divF32: Division by 0")
		}

		output := encoder.Serialize(float32(num1 / num2))

		return &CXArgument{Value: &output, Typ: MakeType("f32")}, nil
	} else {
		return nil, err
	}	
}

func addF64 (arg1 *CXArgument, arg2 *CXArgument) (*CXArgument, error) {
	if err := checkTwoTypes("addF64", "f64", arg1, arg2); err == nil {
		var num1 float64
		var num2 float64
		encoder.DeserializeRaw(*arg1.Value, &num1)
		encoder.DeserializeRaw(*arg2.Value, &num2)

		output := encoder.Serialize(float64(num1 + num2))

		return &CXArgument{Value: &output, Typ: MakeType("f64")}, nil
	} else {
		return nil, err
	}
}

func subF64 (arg1 *CXArgument, arg2 *CXArgument) (*CXArgument, error) {
	if err := checkTwoTypes("subF64", "f64", arg1, arg2); err == nil {
		var num1 float64
		var num2 float64
		encoder.DeserializeRaw(*arg1.Value, &num1)
		encoder.DeserializeRaw(*arg2.Value, &num2)

		output := encoder.Serialize(float64(num1 - num2))

		return &CXArgument{Value: &output, Typ: MakeType("f64")}, nil
	} else {
		return nil, err
	}
}

func mulF64 (arg1 *CXArgument, arg2 *CXArgument) (*CXArgument, error) {
	if err := checkTwoTypes("mulF64", "f64", arg1, arg2); err == nil {
		var num1 float64
		var num2 float64
		encoder.DeserializeRaw(*arg1.Value, &num1)
		encoder.DeserializeRaw(*arg2.Value, &num2)

		output := encoder.Serialize(float64(num1 * num2))

		return &CXArgument{Value: &output, Typ: MakeType("f64")}, nil
	} else {
		return nil, err
	}
}

func divF64 (arg1 *CXArgument, arg2 *CXArgument) (*CXArgument, error) {
	if err := checkTwoTypes("divF64", "f64", arg1, arg2); err == nil {
		var num1 float64
		var num2 float64
		encoder.DeserializeRaw(*arg1.Value, &num1)
		encoder.DeserializeRaw(*arg2.Value, &num2)

		if num2 == float64(0.0) {
			panic("divF64: Division by 0")
		}

		output := encoder.Serialize(float64(num1 / num2))

		return &CXArgument{Value: &output, Typ: MakeType("f64")}, nil
	} else {
		return nil, err
	}
}

func modI32 (arg1 *CXArgument, arg2 *CXArgument) (*CXArgument, error) {
	if err := checkTwoTypes("modI32", "i32", arg1, arg2); err == nil {
		var num1 int32
		var num2 int32
		encoder.DeserializeRaw(*arg1.Value, &num1)
		encoder.DeserializeRaw(*arg2.Value, &num2)

		if num2 == int32(0) {
			panic("modI32: Division by 0")
		}

		output := encoder.Serialize(int32(num1 % num2))

		return &CXArgument{Value: &output, Typ: MakeType("i32")}, nil
	} else {
		return nil, err
	}
}

func modI64 (arg1 *CXArgument, arg2 *CXArgument) (*CXArgument, error) {
	if err := checkTwoTypes("modI64", "i64", arg1, arg2); err == nil {
		var num1 int64
		var num2 int64
		encoder.DeserializeRaw(*arg1.Value, &num1)
		encoder.DeserializeRaw(*arg2.Value, &num2)

		if num2 == int64(0) {
			panic("modI64: Division by 0")
		}

		output := encoder.Serialize(int64(num1 % num2))

		return &CXArgument{Value: &output, Typ: MakeType("i64")}, nil
	} else {
		return nil, err
	}
}

/*
  Bitwise operators
*/

func andI32 (arg1 *CXArgument, arg2 *CXArgument) (*CXArgument, error) {
	if err := checkTwoTypes("andI32", "i32", arg1, arg2); err == nil {
		var num1 int32
		var num2 int32
		encoder.DeserializeRaw(*arg1.Value, &num1)
		encoder.DeserializeRaw(*arg2.Value, &num2)

		output := encoder.Serialize(int32(num1 & num2))

		return &CXArgument{Value: &output, Typ: MakeType("i32")}, nil
	} else {
		return nil, err
	}
}

func orI32 (arg1 *CXArgument, arg2 *CXArgument) (*CXArgument, error) {
	if err := checkTwoTypes("orI32", "i32", arg1, arg2); err == nil {
		var num1 int32
		var num2 int32
		encoder.DeserializeRaw(*arg1.Value, &num1)
		encoder.DeserializeRaw(*arg2.Value, &num2)

		output := encoder.Serialize(int32(num1 | num2))

		return &CXArgument{Value: &output, Typ: MakeType("i32")}, nil
	} else {
		return nil, err
	}
}

func xorI32 (arg1 *CXArgument, arg2 *CXArgument) (*CXArgument, error) {
	if err := checkTwoTypes("xorI32", "i32", arg1, arg2); err == nil {
		var num1 int32
		var num2 int32
		encoder.DeserializeRaw(*arg1.Value, &num1)
		encoder.DeserializeRaw(*arg2.Value, &num2)

		output := encoder.Serialize(int32(num1 ^ num2))

		return &CXArgument{Value: &output, Typ: MakeType("i32")}, nil
	} else {
		return nil, err
	}
}

func andNotI32 (arg1 *CXArgument, arg2 *CXArgument) (*CXArgument, error) {
	if err := checkTwoTypes("andNotI32", "i32", arg1, arg2); err == nil {
		var num1 int32
		var num2 int32
		encoder.DeserializeRaw(*arg1.Value, &num1)
		encoder.DeserializeRaw(*arg2.Value, &num2)

		output := encoder.Serialize(int32(num1 &^ num2))

		return &CXArgument{Value: &output, Typ: MakeType("i32")}, nil
	} else {
		return nil, err
	}
}

func andI64 (arg1 *CXArgument, arg2 *CXArgument) (*CXArgument, error) {
	if err := checkTwoTypes("andI64", "i64", arg1, arg2); err == nil {
		var num1 int64
		var num2 int64
		encoder.DeserializeRaw(*arg1.Value, &num1)
		encoder.DeserializeRaw(*arg2.Value, &num2)

		output := encoder.Serialize(int64(num1 & num2))

		return &CXArgument{Value: &output, Typ: MakeType("i64")}, nil
	} else {
		return nil, err
	}
}

func orI64 (arg1 *CXArgument, arg2 *CXArgument) (*CXArgument, error) {
	if err := checkTwoTypes("orI64", "i64", arg1, arg2); err == nil {
		var num1 int64
		var num2 int64
		encoder.DeserializeRaw(*arg1.Value, &num1)
		encoder.DeserializeRaw(*arg2.Value, &num2)

		output := encoder.Serialize(int64(num1 | num2))

		return &CXArgument{Value: &output, Typ: MakeType("i64")}, nil
	} else {
		return nil, err
	}
}

func xorI64 (arg1 *CXArgument, arg2 *CXArgument) (*CXArgument, error) {
	if err := checkTwoTypes("xorI64", "i64", arg1, arg2); err == nil {
		var num1 int64
		var num2 int64
		encoder.DeserializeRaw(*arg1.Value, &num1)
		encoder.DeserializeRaw(*arg2.Value, &num2)

		output := encoder.Serialize(int64(num1 ^ num2))

		return &CXArgument{Value: &output, Typ: MakeType("i64")}, nil
	} else {
		return nil, err
	}
}

func andNotI64 (arg1 *CXArgument, arg2 *CXArgument) (*CXArgument, error) {
	if err := checkTwoTypes("andNotI64", "i64", arg1, arg2); err == nil {
		var num1 int64
		var num2 int64
		encoder.DeserializeRaw(*arg1.Value, &num1)
		encoder.DeserializeRaw(*arg2.Value, &num2)

		output := encoder.Serialize(int64(num1 &^ num2))

		return &CXArgument{Value: &output, Typ: MakeType("i64")}, nil
	} else {
		return nil, err
	}
}

/*
  Array functions
*/

func readBoolA (arr *CXArgument, idx *CXArgument) *CXArgument {
	if arr.Typ.Name != "[]bool" || idx.Typ.Name != "i32" {
		panic("readBoolA: wrong argument type")
	}

	var index int32
	encoder.DeserializeRaw(*idx.Value, &index)

	var array []int32
	encoder.DeserializeRaw(*arr.Value, &array)

	if index < 0 {
		panic(fmt.Sprintf("readBoolA: Negative index %d", index))
	}
	
	if index >= int32(len(array)) {
		panic(fmt.Sprintf("readBoolA: Index %d exceeds array of length %d", index, len(array)))
	}

	output := encoder.Serialize(array[index])

	return &CXArgument{Value: &output, Typ: MakeType("bool")}
}

func writeBoolA (arr *CXArgument, idx *CXArgument, val *CXArgument) *CXArgument {
	if arr.Typ.Name != "[]bool" || idx.Typ.Name != "i32" || val.Typ.Name != "bool" {
		panic("readBoolA: wrong argument type")
	}

	var index int32
	encoder.DeserializeRaw(*idx.Value, &index)

	var value int32
	encoder.DeserializeRaw(*val.Value, &value)
	
	var array []int32
	encoder.DeserializeRaw(*arr.Value, &array)

	if index < 0 {
		panic(fmt.Sprintf("writeBoolA: Negative index %d", index))
	}
	
	if index >= int32(len(*arr.Value)) {
		panic(fmt.Sprintf("writeBoolA: Index %d exceeds array of length %d", index, len(array)))
	}

	array[index] = value
	output := encoder.Serialize(array)
	arr.Value = &output

	return arr
}

func readByteA (arr *CXArgument, idx *CXArgument) *CXArgument {
	if arr.Typ.Name != "[]byte" || idx.Typ.Name != "i32" {
		panic("readByteA: wrong argument type")
	}

	var index int32
	encoder.DeserializeRaw(*idx.Value, &index)

	if index < 0 {
		panic(fmt.Sprintf("readByteA: Negative index %d", index))
	}
	
	if index >= int32(len(*arr.Value)) {
		panic(fmt.Sprintf("readByteA: Index %d exceeds array of length %d", index, len(*arr.Value)))
	}

	output := make([]byte, 1)
	output[0] = (*arr.Value)[index]

	return &CXArgument{Value: &output, Typ: MakeType("byte")}
}

func writeByteA (arr *CXArgument, idx *CXArgument, val *CXArgument) *CXArgument {
	if arr.Typ.Name != "[]byte" || idx.Typ.Name != "i32" || val.Typ.Name != "byte" {
		panic("readByteA: wrong argument type")
	}

	var index int32
	encoder.DeserializeRaw(*idx.Value, &index)

	if index < 0 {
		panic(fmt.Sprintf("writeByteA: Negative index %d", index))
	}
	
	if index >= int32(len(*arr.Value)) {
		panic(fmt.Sprintf("writeByteA: Index %d exceeds array of length %d", index, len(*arr.Value)))
	}

	(*arr.Value)[index] = (*val.Value)[0]

	return arr
}

func readI32A (arr *CXArgument, idx *CXArgument) *CXArgument {
	if arr.Typ.Name != "[]i32" || idx.Typ.Name != "i32" {
		panic("readI32A: wrong argument type")
	}

	var index int32
	encoder.DeserializeRaw(*idx.Value, &index)

	var array []int32
	encoder.DeserializeRaw(*arr.Value, &array)

	if index < 0 {
		panic(fmt.Sprintf("readI32A: Negative index %d", index))
	}
	
	if index >= int32(len(array)) {
		panic(fmt.Sprintf("readI32A: Index %d exceeds array of length %d", index, len(array)))
	}

	output := encoder.Serialize(array[index])

	return &CXArgument{Value: &output, Typ: MakeType("i32")}
}

func writeI32A (arr *CXArgument, idx *CXArgument, val *CXArgument) *CXArgument {
	if arr.Typ.Name != "[]i32" || idx.Typ.Name != "i32" || val.Typ.Name != "i32" {
		panic("writeI32A: wrong argument type")
	}

	var index int32
	encoder.DeserializeRaw(*idx.Value, &index)

	var value int32
	encoder.DeserializeRaw(*val.Value, &value)

	var array []int32
	encoder.DeserializeRaw(*arr.Value, &array)

	if index < 0 {
		panic(fmt.Sprintf("writeI32A: Negative index %d", index))
	}

	if index >= int32(len(array)) {
		panic(fmt.Sprintf("writeI32A: Index %d exceeds array of length %d", index, len(array)))
	}

	array[index] = value
	output := encoder.Serialize(array)
	arr.Value = &output

	return arr
}

func readI64A (arr *CXArgument, idx *CXArgument) *CXArgument {
	if arr.Typ.Name != "[]i64" || idx.Typ.Name != "i32" {
		panic("readI64A: wrong argument type")
	}

	var index int32
	encoder.DeserializeRaw(*idx.Value, &index)

	var array []int64
	encoder.DeserializeRaw(*arr.Value, &array)

	if index < 0 {
		panic(fmt.Sprintf("readI64A: Negative index %d", index))
	}

	if index >= int32(len(array)) {
		panic(fmt.Sprintf("readI64A: Index %d exceeds array of length %d", index, len(array)))
	}

	output := encoder.Serialize(array[index])

	return &CXArgument{Value: &output, Typ: MakeType("i64")}
}

func writeI64A (arr *CXArgument, idx *CXArgument, val *CXArgument) *CXArgument {
	if arr.Typ.Name != "[]i64" || idx.Typ.Name != "i32" || val.Typ.Name != "i64" {
		panic("readI64A: wrong argument type")
	}

	var index int32
	encoder.DeserializeRaw(*idx.Value, &index)

	var value int64
	encoder.DeserializeRaw(*val.Value, &value)

	var array []int64
	encoder.DeserializeRaw(*arr.Value, &array)

	if index < 0 {
		panic(fmt.Sprintf("writeI64A: Negative index %d", index))
	}

	if index >= int32(len(array)) {
		panic(fmt.Sprintf("writeI64A: Index %d exceeds array of length %d", index, len(array)))
	}

	array[index] = value
	output := encoder.Serialize(array)
	arr.Value = &output

	return arr
}

func readF32A (arr *CXArgument, idx *CXArgument) *CXArgument {
	if arr.Typ.Name != "[]f32" || idx.Typ.Name != "i32" {
		panic("readF32A: wrong argument type")
	}

	var index int32
	encoder.DeserializeRaw(*idx.Value, &index)

	var array []float32
	encoder.DeserializeRaw(*arr.Value, &array)

	if index < 0 {
		panic(fmt.Sprintf("readF32A: Negative index %d", index))
	}

	if index >= int32(len(array)) {
		panic(fmt.Sprintf("readF32A: Index %d exceeds array of length %d", index, len(array)))
	}

	output := encoder.Serialize(array[index])

	return &CXArgument{Value: &output, Typ: MakeType("f32")}
}

func writeF32A (arr *CXArgument, idx *CXArgument, val *CXArgument) *CXArgument {
	if arr.Typ.Name != "[]f32" || idx.Typ.Name != "i32" || val.Typ.Name != "f32" {
		panic("writeF32A: wrong argument type")
	}

	var index int32
	encoder.DeserializeRaw(*idx.Value, &index)

	var value float32
	encoder.DeserializeRaw(*val.Value, &value)

	var array []float32
	encoder.DeserializeRaw(*arr.Value, &array)

	if index < 0 {
		panic(fmt.Sprintf("writeF32A: Negative index %d", index))
	}

	if index >= int32(len(array)) {
		panic(fmt.Sprintf("writeF32A: Index %d exceeds array of length %d", index, len(array)))
	}

	array[index] = value
	output := encoder.Serialize(array)
	arr.Value = &output

	return arr
}

func readF64A (arr *CXArgument, idx *CXArgument) *CXArgument {
	if arr.Typ.Name != "[]f64" || idx.Typ.Name != "i32" {
		panic("readF64A: wrong argument type")
	}

	var index int32
	encoder.DeserializeRaw(*idx.Value, &index)

	var array []float64
	encoder.DeserializeRaw(*arr.Value, &array)

	if index < 0 {
		panic(fmt.Sprintf("readF64A: Negative index %d", index))
	}

	if index >= int32(len(array)) {
		panic(fmt.Sprintf("readF64A: Index %d exceeds array of length %d", index, len(array)))
	}

	output := encoder.Serialize(array[index])

	return &CXArgument{Value: &output, Typ: MakeType("f64")}
}

func writeF64A (arr *CXArgument, idx *CXArgument, val *CXArgument) *CXArgument {
	if arr.Typ.Name != "[]f64" || idx.Typ.Name != "i32" || val.Typ.Name != "f64" {
		panic("readF64A: wrong argument type")
	}

	var index int32
	encoder.DeserializeRaw(*idx.Value, &index)

	var value float64
	encoder.DeserializeRaw(*val.Value, &value)

	var array []float64
	encoder.DeserializeRaw(*arr.Value, &array)

	if index < 0 {
		panic(fmt.Sprintf("writeF64A: Negative index %d", index))
	}

	if index >= int32(len(array)) {
		panic(fmt.Sprintf("writeF64A: Index %d exceeds array of length %d", index, len(array)))
	}

	array[index] = value
	output := encoder.Serialize(array)
	arr.Value = &output

	return arr
}

func lenBoolA (arr *CXArgument) *CXArgument {
	if arr.Typ.Name != "[]bool" {
		panic("lenBoolA: wrong argument type")
	}

	var array []int32
	encoder.DeserializeRaw(*arr.Value, &array)

	output := encoder.SerializeAtomic(int32(len(array)))
	return &CXArgument{Value: &output, Typ: MakeType("i32")}
}

func lenByteA (arr *CXArgument) *CXArgument {
	if arr.Typ.Name != "[]byte" {
		panic("lenByteA: wrong argument type")
	}

	output := encoder.SerializeAtomic(int32(len(*arr.Value)))
	
	return &CXArgument{Value: &output, Typ: MakeType("i32")}
}

func lenI32A (arr *CXArgument) *CXArgument {
	if arr.Typ.Name != "[]i32" {
		panic("lenI32A: wrong argument type")
	}

	var array []int32
	encoder.DeserializeRaw(*arr.Value, &array)

	output := encoder.SerializeAtomic(int32(len(array)))
	return &CXArgument{Value: &output, Typ: MakeType("i32")}
}

func lenI64A (arr *CXArgument) *CXArgument {
	if arr.Typ.Name != "[]i64" {
		panic("lenI64A: wrong argument type")
	}

	var array []int64
	encoder.DeserializeRaw(*arr.Value, &array)

	output := encoder.SerializeAtomic(int32(len(array)))
	return &CXArgument{Value: &output, Typ: MakeType("i32")}
}

func lenF32A (arr *CXArgument) *CXArgument {
	if arr.Typ.Name != "[]f32" {
		panic("lenF32A: wrong argument type")
	}

	var array []float32
	encoder.DeserializeRaw(*arr.Value, &array)

	output := encoder.SerializeAtomic(int32(len(array)))
	return &CXArgument{Value: &output, Typ: MakeType("i32")}
}

func lenF64A (arr *CXArgument) *CXArgument {
	if arr.Typ.Name != "[]f64" {
		panic("lenF64A: wrong argument type")
	}

	var array []float64
	encoder.DeserializeRaw(*arr.Value, &array)

	output := encoder.SerializeAtomic(int32(len(array)))
	return &CXArgument{Value: &output, Typ: MakeType("i32")}
}

/*
  Logical Operators
*/

func and (arg1 *CXArgument, arg2 *CXArgument) *CXArgument {
	if arg1.Typ.Name != "bool" || arg2.Typ.Name != "bool" {
		panic("and: wrong argument type")
	}

	var c1 int32
	var c2 int32
	encoder.DeserializeRaw(*arg1.Value, &c1)
	encoder.DeserializeRaw(*arg2.Value, &c2)

	if c1 == 1 && c2 == 1 {
		val := encoder.Serialize(int32(1))
		return &CXArgument{Value: &val, Typ: MakeType("bool")}
	} else {
		val := encoder.Serialize(int32(0))
		return &CXArgument{Value: &val, Typ: MakeType("bool")}
	}
}

func or (arg1 *CXArgument, arg2 *CXArgument) *CXArgument {
	if arg1.Typ.Name != "bool" || arg2.Typ.Name != "bool" {
		panic("or: wrong argument type")
	}

	var c1 int32
	var c2 int32
	encoder.DeserializeRaw(*arg1.Value, &c1)
	encoder.DeserializeRaw(*arg2.Value, &c2)

	if c1 == 1 || c2 == 1 {
		val := encoder.Serialize(int32(1))
		return &CXArgument{Value: &val, Typ: MakeType("bool")}
	} else {
		val := encoder.Serialize(int32(0))
		return &CXArgument{Value: &val, Typ: MakeType("bool")}
	}
}

func not (arg1 *CXArgument) *CXArgument {
	if arg1.Typ.Name != "bool" {
		panic("not: wrong argument type")
	}

	var c1 int32
	encoder.DeserializeRaw(*arg1.Value, &c1)

	if c1 == 0 {
		val := encoder.Serialize(int32(1))
		return &CXArgument{Value: &val, Typ: MakeType("bool")}
	} else {
		val := encoder.Serialize(int32(0))
		return &CXArgument{Value: &val, Typ: MakeType("bool")}
	}
}

/*
  Relational Operators
*/

func ltI32 (arg1 *CXArgument, arg2 *CXArgument) *CXArgument {
	if arg1.Typ.Name != "i32" || arg2.Typ.Name != "i32" {
		panic("ltI32: wrong argument type")
	}

	var num1 int32
	var num2 int32
	encoder.DeserializeRaw(*arg1.Value, &num1)
	encoder.DeserializeRaw(*arg2.Value, &num2)
	
	if num1 < num2 {
		val := encoder.Serialize(int32(1))
		return &CXArgument{Value: &val, Typ: MakeType("bool")}
	} else {
		val := encoder.Serialize(int32(0))
		return &CXArgument{Value: &val, Typ: MakeType("bool")}
	}
}

func gtI32 (arg1 *CXArgument, arg2 *CXArgument) *CXArgument {
	if arg1.Typ.Name != "i32" || arg2.Typ.Name != "i32" {
		panic("gtI32: wrong argument type")
	}

	var num1 int32
	var num2 int32
	encoder.DeserializeRaw(*arg1.Value, &num1)
	encoder.DeserializeRaw(*arg2.Value, &num2)
	
	if num1 > num2 {
		val := encoder.Serialize(int32(1))
		return &CXArgument{Value: &val, Typ: MakeType("bool")}
	} else {
		val := encoder.Serialize(int32(0))
		return &CXArgument{Value: &val, Typ: MakeType("bool")}
	}
}

func eqI32 (arg1 *CXArgument, arg2 *CXArgument) *CXArgument {
	if arg1.Typ.Name != "i32" || arg2.Typ.Name != "i32" {
		panic("eqI32: wrong argument type")
	}

	var num1 int32
	var num2 int32
	encoder.DeserializeRaw(*arg1.Value, &num1)
	encoder.DeserializeRaw(*arg2.Value, &num2)
	
	if num1 == num2 {
		val := encoder.Serialize(int32(1))
		return &CXArgument{Value: &val, Typ: MakeType("bool")}
	} else {
		val := encoder.Serialize(int32(0))
		return &CXArgument{Value: &val, Typ: MakeType("bool")}
	}
}

func lteqI32 (arg1 *CXArgument, arg2 *CXArgument) *CXArgument {
	if arg1.Typ.Name != "i32" || arg2.Typ.Name != "i32" {
		panic("lteqI32: wrong argument type")
	}

	var num1 int32
	var num2 int32
	encoder.DeserializeRaw(*arg1.Value, &num1)
	encoder.DeserializeRaw(*arg2.Value, &num2)
	
	if num1 <= num2 {
		val := encoder.Serialize(int32(1))
		return &CXArgument{Value: &val, Typ: MakeType("bool")}
	} else {
		val := encoder.Serialize(int32(0))
		return &CXArgument{Value: &val, Typ: MakeType("bool")}
	}
}

func gteqI32 (arg1 *CXArgument, arg2 *CXArgument) *CXArgument {
	if arg1.Typ.Name != "i32" || arg2.Typ.Name != "i32" {
		panic("lteqI32: wrong argument type")
	}

	var num1 int32
	var num2 int32
	encoder.DeserializeRaw(*arg1.Value, &num1)
	encoder.DeserializeRaw(*arg2.Value, &num2)
	
	if num1 >= num2 {
		val := encoder.Serialize(int32(1))
		return &CXArgument{Value: &val, Typ: MakeType("bool")}
	} else {
		val := encoder.Serialize(int32(0))
		return &CXArgument{Value: &val, Typ: MakeType("bool")}
	}
}

func ltI64 (arg1 *CXArgument, arg2 *CXArgument) *CXArgument {
	if arg1.Typ.Name != "i64" || arg2.Typ.Name != "i64" {
		panic("ltI64: wrong argument type")
	}

	var num1 int64
	var num2 int64
	encoder.DeserializeRaw(*arg1.Value, &num1)
	encoder.DeserializeRaw(*arg2.Value, &num2)
	
	if num1 < num2 {
		val := encoder.Serialize(int32(1))
		return &CXArgument{Value: &val, Typ: MakeType("bool")}
	} else {
		val := encoder.Serialize(int32(0))
		return &CXArgument{Value: &val, Typ: MakeType("bool")}
	}
}

func gtI64 (arg1 *CXArgument, arg2 *CXArgument) *CXArgument {
	if arg1.Typ.Name != "i64" || arg2.Typ.Name != "i64" {
		panic("gtI64: wrong argument type")
	}

	var num1 int64
	var num2 int64
	encoder.DeserializeRaw(*arg1.Value, &num1)
	encoder.DeserializeRaw(*arg2.Value, &num2)
	
	if num1 > num2 {
		val := encoder.Serialize(int32(1))
		return &CXArgument{Value: &val, Typ: MakeType("bool")}
	} else {
		val := encoder.Serialize(int32(0))
		return &CXArgument{Value: &val, Typ: MakeType("bool")}
	}
}

func eqI64 (arg1 *CXArgument, arg2 *CXArgument) *CXArgument {
	if arg1.Typ.Name != "i64" || arg2.Typ.Name != "i64" {
		panic("eqI64: wrong argument type")
	}

	var num1 int64
	var num2 int64
	encoder.DeserializeRaw(*arg1.Value, &num1)
	encoder.DeserializeRaw(*arg2.Value, &num2)
	
	if num1 == num2 {
		val := encoder.Serialize(int32(1))
		return &CXArgument{Value: &val, Typ: MakeType("bool")}
	} else {
		val := encoder.Serialize(int32(0))
		return &CXArgument{Value: &val, Typ: MakeType("bool")}
	}
}

func lteqI64 (arg1 *CXArgument, arg2 *CXArgument) *CXArgument {
	if arg1.Typ.Name != "i64" || arg2.Typ.Name != "i64" {
		panic("lteqI64: wrong argument type")
	}

	var num1 int64
	var num2 int64
	encoder.DeserializeRaw(*arg1.Value, &num1)
	encoder.DeserializeRaw(*arg2.Value, &num2)
	
	if num1 <= num2 {
		val := encoder.Serialize(int32(1))
		return &CXArgument{Value: &val, Typ: MakeType("bool")}
	} else {
		val := encoder.Serialize(int32(0))
		return &CXArgument{Value: &val, Typ: MakeType("bool")}
	}
}

func gteqI64 (arg1 *CXArgument, arg2 *CXArgument) *CXArgument {
	if arg1.Typ.Name != "i64" || arg2.Typ.Name != "i64" {
		panic("gteqI64: wrong argument type")
	}

	var num1 int64
	var num2 int64
	encoder.DeserializeRaw(*arg1.Value, &num1)
	encoder.DeserializeRaw(*arg2.Value, &num2)
	
	if num1 >= num2 {
		val := encoder.Serialize(int32(1))
		return &CXArgument{Value: &val, Typ: MakeType("bool")}
	} else {
		val := encoder.Serialize(int32(0))
		return &CXArgument{Value: &val, Typ: MakeType("bool")}
	}
}

func ltF32 (arg1 *CXArgument, arg2 *CXArgument) *CXArgument {
	if arg1.Typ.Name != "f32" || arg2.Typ.Name != "f32" {
		panic("ltF32: wrong argument type")
	}

	var num1 float32
	var num2 float32
	encoder.DeserializeRaw(*arg1.Value, &num1)
	encoder.DeserializeRaw(*arg2.Value, &num2)
	
	if num1 < num2 {
		val := encoder.Serialize(int32(1))
		return &CXArgument{Value: &val, Typ: MakeType("bool")}
	} else {
		val := encoder.Serialize(int32(0))
		return &CXArgument{Value: &val, Typ: MakeType("bool")}
	}
}

func gtF32 (arg1 *CXArgument, arg2 *CXArgument) *CXArgument {
	if arg1.Typ.Name != "f32" || arg2.Typ.Name != "f32" {
		panic("gtF32: wrong argument type")
	}

	var num1 float32
	var num2 float32
	encoder.DeserializeRaw(*arg1.Value, &num1)
	encoder.DeserializeRaw(*arg2.Value, &num2)
	
	if num1 > num2 {
		val := encoder.Serialize(int32(1))
		return &CXArgument{Value: &val, Typ: MakeType("bool")}
	} else {
		val := encoder.Serialize(int32(0))
		return &CXArgument{Value: &val, Typ: MakeType("bool")}
	}
}

func eqF32 (arg1 *CXArgument, arg2 *CXArgument) *CXArgument {
	if arg1.Typ.Name != "f32" || arg2.Typ.Name != "f32" {
		panic("eqF32: wrong argument type")
	}

	var num1 float32
	var num2 float32
	encoder.DeserializeRaw(*arg1.Value, &num1)
	encoder.DeserializeRaw(*arg2.Value, &num2)
	
	if num1 == num2 {
		val := encoder.Serialize(int32(1))
		return &CXArgument{Value: &val, Typ: MakeType("bool")}
	} else {
		val := encoder.Serialize(int32(0))
		return &CXArgument{Value: &val, Typ: MakeType("bool")}
	}
}

func lteqF32 (arg1 *CXArgument, arg2 *CXArgument) *CXArgument {
	if arg1.Typ.Name != "f32" || arg2.Typ.Name != "f32" {
		panic("lteqF32: wrong argument type")
	}

	var num1 float32
	var num2 float32
	encoder.DeserializeRaw(*arg1.Value, &num1)
	encoder.DeserializeRaw(*arg2.Value, &num2)
	
	if num1 <= num2 {
		val := encoder.Serialize(int32(1))
		return &CXArgument{Value: &val, Typ: MakeType("bool")}
	} else {
		val := encoder.Serialize(int32(0))
		return &CXArgument{Value: &val, Typ: MakeType("bool")}
	}
}

func gteqF32 (arg1 *CXArgument, arg2 *CXArgument) *CXArgument {
	if arg1.Typ.Name != "f32" || arg2.Typ.Name != "f32" {
		panic("gteqF32: wrong argument type")
	}

	var num1 float32
	var num2 float32
	encoder.DeserializeRaw(*arg1.Value, &num1)
	encoder.DeserializeRaw(*arg2.Value, &num2)
	
	if num1 >= num2 {
		val := encoder.Serialize(int32(1))
		return &CXArgument{Value: &val, Typ: MakeType("bool")}
	} else {
		val := encoder.Serialize(int32(0))
		return &CXArgument{Value: &val, Typ: MakeType("bool")}
	}
}












func ltF64 (arg1 *CXArgument, arg2 *CXArgument) *CXArgument {
	if arg1.Typ.Name != "f64" || arg2.Typ.Name != "f64" {
		panic("ltF64: wrong argument type")
	}

	var num1 float64
	var num2 float64
	encoder.DeserializeRaw(*arg1.Value, &num1)
	encoder.DeserializeRaw(*arg2.Value, &num2)
	
	if num1 < num2 {
		val := encoder.Serialize(int32(1))
		return &CXArgument{Value: &val, Typ: MakeType("bool")}
	} else {
		val := encoder.Serialize(int32(0))
		return &CXArgument{Value: &val, Typ: MakeType("bool")}
	}
}

func gtF64 (arg1 *CXArgument, arg2 *CXArgument) *CXArgument {
	if arg1.Typ.Name != "f64" || arg2.Typ.Name != "f64" {
		panic("gtF64: wrong argument type")
	}

	var num1 float64
	var num2 float64
	encoder.DeserializeRaw(*arg1.Value, &num1)
	encoder.DeserializeRaw(*arg2.Value, &num2)
	
	if num1 > num2 {
		val := encoder.Serialize(int32(1))
		return &CXArgument{Value: &val, Typ: MakeType("bool")}
	} else {
		val := encoder.Serialize(int32(0))
		return &CXArgument{Value: &val, Typ: MakeType("bool")}
	}
}

func eqF64 (arg1 *CXArgument, arg2 *CXArgument) *CXArgument {
	if arg1.Typ.Name != "f64" || arg2.Typ.Name != "f64" {
		panic("eqF64: wrong argument type")
	}

	var num1 float64
	var num2 float64
	encoder.DeserializeRaw(*arg1.Value, &num1)
	encoder.DeserializeRaw(*arg2.Value, &num2)
	
	if num1 == num2 {
		val := encoder.Serialize(int32(1))
		return &CXArgument{Value: &val, Typ: MakeType("bool")}
	} else {
		val := encoder.Serialize(int32(0))
		return &CXArgument{Value: &val, Typ: MakeType("bool")}
	}
}

func lteqF64 (arg1 *CXArgument, arg2 *CXArgument) *CXArgument {
	if arg1.Typ.Name != "f64" || arg2.Typ.Name != "f64" {
		panic("lteqF64: wrong argument type")
	}

	var num1 float64
	var num2 float64
	encoder.DeserializeRaw(*arg1.Value, &num1)
	encoder.DeserializeRaw(*arg2.Value, &num2)
	
	if num1 <= num2 {
		val := encoder.Serialize(int32(1))
		return &CXArgument{Value: &val, Typ: MakeType("bool")}
	} else {
		val := encoder.Serialize(int32(0))
		return &CXArgument{Value: &val, Typ: MakeType("bool")}
	}
}

func gteqF64 (arg1 *CXArgument, arg2 *CXArgument) *CXArgument {
	if arg1.Typ.Name != "f64" || arg2.Typ.Name != "f64" {
		panic("gteqF64: wrong argument type")
	}

	var num1 float64
	var num2 float64
	encoder.DeserializeRaw(*arg1.Value, &num1)
	encoder.DeserializeRaw(*arg2.Value, &num2)
	
	if num1 >= num2 {
		val := encoder.Serialize(int32(1))
		return &CXArgument{Value: &val, Typ: MakeType("bool")}
	} else {
		val := encoder.Serialize(int32(0))
		return &CXArgument{Value: &val, Typ: MakeType("bool")}
	}
}

func eqStr (arg1 *CXArgument, arg2 *CXArgument) *CXArgument {
	if arg1.Typ.Name != "str" || arg2.Typ.Name != "str" {
		panic("eqStr: wrong argument type")
	}

	str1 := string(*arg1.Value)
	str2 := string(*arg2.Value)
	
	if str1 == str2 {
		val := encoder.Serialize(int32(1))
		return &CXArgument{Value: &val, Typ: MakeType("bool")}
	} else {
		val := encoder.Serialize(int32(0))
		return &CXArgument{Value: &val, Typ: MakeType("bool")}
	}
}

/*
  Cast functions
*/

func castToStr (arg *CXArgument) *CXArgument {
	if arg.Typ.Name != "[]byte" {
		panic("castToStr: wrong argument type")
	}
	strTyp := MakeType("str")
	switch arg.Typ.Name {
	case "[]byte":
		newArg := MakeArgument(arg.Value, strTyp)
		return newArg
	default:
		panic(fmt.Sprintf("Type '%s' can't be casted to type 'str'", arg.Typ.Name))
	}
}

func byteAtoStr (arg *CXArgument) *CXArgument {
	if arg.Typ.Name != "[]byte" {
		panic("byteAtoStr: wrong argument type")
	}
	
	strTyp := MakeType("str")
	newArg := MakeArgument(arg.Value, strTyp)
		return newArg
}

func castToI32 (arg *CXArgument) *CXArgument {
	if arg.Typ.Name != "i32" && arg.Typ.Name != "i64" && arg.Typ.Name != "f32" && arg.Typ.Name != "f64" {
		panic("castToI32: wrong argument type")
	}
	i32Typ := MakeType("i32")
	switch arg.Typ.Name {
	case "i32":
		return arg
	case "i64":
		var val int64
		encoder.DeserializeRaw(*arg.Value, &val)
		newVal := encoder.Serialize(int32(val))
		newArg := MakeArgument(&newVal, i32Typ)
		return newArg
	case "f32":
		var val float32
		encoder.DeserializeRaw(*arg.Value, &val)
		newVal := encoder.Serialize(int32(val))
		newArg := MakeArgument(&newVal, i32Typ)
		return newArg
	case "f64":
		var val float64
		encoder.DeserializeRaw(*arg.Value, &val)
		newVal := encoder.Serialize(int32(val))
		newArg := MakeArgument(&newVal, i32Typ)
		return newArg
	default:
		panic(fmt.Sprintf("Type '%s' can't be casted to type 'i32'", arg.Typ.Name))
	}
}

func castToI64 (arg *CXArgument) *CXArgument {
	if arg.Typ.Name != "i32" && arg.Typ.Name != "i64" && arg.Typ.Name != "f32" && arg.Typ.Name != "f64" {
		panic("castToI64: wrong argument type")
	}
	i64Typ := MakeType("i64")
	switch arg.Typ.Name {
	case "i32":
		var val int32
		encoder.DeserializeRaw(*arg.Value, &val)
		newVal := encoder.Serialize(int64(val))
		newArg := MakeArgument(&newVal, i64Typ)
		return newArg
	case "i64":
		return arg
	case "f32":
		var val float32
		encoder.DeserializeRaw(*arg.Value, &val)
		newVal := encoder.Serialize(int64(val))
		newArg := MakeArgument(&newVal, i64Typ)
		return newArg
	case "f64":
		var val float64
		encoder.DeserializeRaw(*arg.Value, &val)
		newVal := encoder.Serialize(int64(val))
		newArg := MakeArgument(&newVal, i64Typ)
		return newArg
	default:
		panic(fmt.Sprintf("Type '%s' can't be casted to type 'i64'", arg.Typ.Name))
	}
}

func castToF32 (arg *CXArgument) *CXArgument {
	if arg.Typ.Name != "i32" && arg.Typ.Name != "i64" && arg.Typ.Name != "f32" && arg.Typ.Name != "f64" {
		panic("castToF32: wrong argument type")
	}
	f32Typ := MakeType("f32")
	switch arg.Typ.Name {
	case "i32":
		var val int32
		encoder.DeserializeRaw(*arg.Value, &val)
		newVal := encoder.Serialize(float32(val))
		newArg := MakeArgument(&newVal, f32Typ)
		return newArg
	case "i64":
		var val int64
		encoder.DeserializeRaw(*arg.Value, &val)
		newVal := encoder.Serialize(float32(val))
		newArg := MakeArgument(&newVal, f32Typ)
		return newArg
	case "f32":
		return arg
	case "f64":
		var val float64
		encoder.DeserializeRaw(*arg.Value, &val)
		newVal := encoder.Serialize(float32(val))
		newArg := MakeArgument(&newVal, f32Typ)
		return newArg
	default:
		panic(fmt.Sprintf("Type '%s' can't be casted to type 'f32'", arg.Typ.Name))
	}
}

func castToF64 (arg *CXArgument) *CXArgument {
	if arg.Typ.Name != "i32" && arg.Typ.Name != "i64" && arg.Typ.Name != "f32" && arg.Typ.Name != "f64" {
		panic("castToF64: wrong argument type")
	}
	
	f64Typ := MakeType("f64")
	switch arg.Typ.Name {
	case "i32":
		var val int32
		encoder.DeserializeRaw(*arg.Value, &val)
		newVal := encoder.Serialize(float64(val))
		newArg := MakeArgument(&newVal, f64Typ)
		return newArg
	case "i64":
		var val int64
		encoder.DeserializeRaw(*arg.Value, &val)
		newVal := encoder.Serialize(float64(val))
		newArg := MakeArgument(&newVal, f64Typ)
		return newArg
	case "f32":
		var val float32
		encoder.DeserializeRaw(*arg.Value, &val)
		newVal := encoder.Serialize(float64(val))
		newArg := MakeArgument(&newVal, f64Typ)
		return newArg
	case "f64":
		return arg
	default:
		panic(fmt.Sprintf("Type '%s' can't be casted to type 'f64'", arg.Typ.Name))
	}
}

func castToI32A (arg *CXArgument) *CXArgument {
	if arg.Typ.Name != "[]i32" && arg.Typ.Name != "[]i64" && arg.Typ.Name != "[]f32" && arg.Typ.Name != "[]f64" {
		panic("castToI32A: wrong argument type")
	}
	i32ATyp := MakeType("[]i32")
	switch arg.Typ.Name {
	case "[]i32":
		return arg
	case "[]i64":
		var val []int64
		encoder.DeserializeRaw(*arg.Value, &val)

		output := make([]int32, len(val))
		for i, n := range val {
			output[i] = int32(n)
		}
		
		newVal := encoder.Serialize(output)
		newArg := MakeArgument(&newVal, i32ATyp)
		return newArg
	case "[]f32":
		var val []float32
		encoder.DeserializeRaw(*arg.Value, &val)

		output := make([]int32, len(val))
		for i, n := range val {
			output[i] = int32(n)
		}
		
		newVal := encoder.Serialize(output)
		newArg := MakeArgument(&newVal, i32ATyp)
		return newArg
	case "[]f64":
		var val []float64
		encoder.DeserializeRaw(*arg.Value, &val)

		output := make([]int32, len(val))
		for i, n := range val {
			output[i] = int32(n)
		}
		
		newVal := encoder.Serialize(output)
		newArg := MakeArgument(&newVal, i32ATyp)
		return newArg
	default:
		panic(fmt.Sprintf("Type '%s' can't be casted to type '[]i32'", arg.Typ.Name))
	}
}

func castToI64A (arg *CXArgument) *CXArgument {
	if arg.Typ.Name != "[]i32" && arg.Typ.Name != "[]i64" && arg.Typ.Name != "[]f32" && arg.Typ.Name != "[]f64" {
		panic("castToI64A: wrong argument type")
	}
	i64ATyp := MakeType("[]i64")
	switch arg.Typ.Name {
	case "[]i32":
		var val []int32
		encoder.DeserializeRaw(*arg.Value, &val)

		output := make([]int64, len(val))
		for i, n := range val {
			output[i] = int64(n)
		}
		
		newVal := encoder.Serialize(output)
		newArg := MakeArgument(&newVal, i64ATyp)
		return newArg
	case "[]i64":
		return arg
	case "[]f32":
		var val []float32
		encoder.DeserializeRaw(*arg.Value, &val)

		output := make([]int64, len(val))
		for i, n := range val {
			output[i] = int64(n)
		}
		
		newVal := encoder.Serialize(output)
		newArg := MakeArgument(&newVal, i64ATyp)
		return newArg
	case "[]f64":
		var val []float64
		encoder.DeserializeRaw(*arg.Value, &val)

		output := make([]int64, len(val))
		for i, n := range val {
			output[i] = int64(n)
		}
		
		newVal := encoder.Serialize(output)
		newArg := MakeArgument(&newVal, i64ATyp)
		return newArg
	default:
		panic(fmt.Sprintf("Type '%s' can't be casted to type '[]i64'", arg.Typ.Name))
	}
}

func castToF32A (arg *CXArgument) *CXArgument {
	if arg.Typ.Name != "[]i32" && arg.Typ.Name != "[]i64" && arg.Typ.Name != "[]f32" && arg.Typ.Name != "[]f64" {
		panic("castToF32A: wrong argument type")
	}
	f32ATyp := MakeType("[]f32")
	switch arg.Typ.Name {
	case "[]i32":
		var val []int32
		encoder.DeserializeRaw(*arg.Value, &val)

		output := make([]float32, len(val))
		for i, n := range val {
			output[i] = float32(n)
		}
		
		newVal := encoder.Serialize(output)
		newArg := MakeArgument(&newVal, f32ATyp)
		return newArg
	case "[]i64":
		var val []int64
		encoder.DeserializeRaw(*arg.Value, &val)

		output := make([]float32, len(val))
		for i, n := range val {
			output[i] = float32(n)
		}
		
		newVal := encoder.Serialize(output)
		newArg := MakeArgument(&newVal, f32ATyp)
		return newArg
	case "[]f32":
		return arg
	case "[]f64":
		var val []float64
		encoder.DeserializeRaw(*arg.Value, &val)

		output := make([]float32, len(val))
		for i, n := range val {
			output[i] = float32(n)
		}
		
		newVal := encoder.Serialize(output)
		newArg := MakeArgument(&newVal, f32ATyp)
		return newArg
	default:
		panic(fmt.Sprintf("Type '%s' can't be casted to type '[]f32'", arg.Typ.Name))
	}
}

func castToF64A (arg *CXArgument) *CXArgument {
	if arg.Typ.Name != "[]i32" && arg.Typ.Name != "[]i64" && arg.Typ.Name != "[]f32" && arg.Typ.Name != "[]f64" {
		panic("castToF64A: wrong argument type")
	}
	f64ATyp := MakeType("[]f64")
	switch arg.Typ.Name {
	case "[]i32":
		var val []int32
		encoder.DeserializeRaw(*arg.Value, &val)

		output := make([]float64, len(val))
		for i, n := range val {
			output[i] = float64(n)
		}
		
		newVal := encoder.Serialize(output)
		newArg := MakeArgument(&newVal, f64ATyp)
		return newArg
	case "[]i64":
		var val []int64
		encoder.DeserializeRaw(*arg.Value, &val)

		output := make([]float64, len(val))
		for i, n := range val {
			output[i] = float64(n)
		}
		
		newVal := encoder.Serialize(output)
		newArg := MakeArgument(&newVal, f64ATyp)
		return newArg
	case "[]f32":
		var val []float32
		encoder.DeserializeRaw(*arg.Value, &val)

		output := make([]float64, len(val))
		for i, n := range val {
			output[i] = float64(n)
		}
		
		newVal := encoder.Serialize(output)
		newArg := MakeArgument(&newVal, f64ATyp)
		return newArg
	case "[]f64":
		return arg
	default:
		panic(fmt.Sprintf("Type '%s' can't be casted to type '[]f64'", arg.Typ.Name))
	}
}

// goTo increments/decrements the call.Line to the desired expression line.
// Used for if/else and loop statements.
func goTo (call *CXCall, predicate *CXArgument, thenLine *CXArgument, elseLine *CXArgument) *CXArgument {
	if predicate.Typ.Name != "bool" || thenLine.Typ.Name != "i32" || elseLine.Typ.Name != "i32" {
		panic("goTo: wrong argument type")
	}
	var isFalse bool

	var pred int32
	encoder.DeserializeRaw(*predicate.Value, &pred)

	if pred == 0 {
		isFalse = true
	} else {
		isFalse = false
	}

	var thenLineNo int32
	var elseLineNo int32

	encoder.DeserializeRaw(*thenLine.Value, &thenLineNo)
	encoder.DeserializeRaw(*elseLine.Value, &elseLineNo)

	if isFalse {
		//call.Line = int(elseLineNo) - 1
		call.Line = call.Line + int(elseLineNo) - 1
	} else {
		//call.Line = int(thenLineNo) - 1
		call.Line = call.Line + int(thenLineNo) - 1
	}
	
	if isFalse {
		val := encoder.Serialize(int32(0))
		return MakeArgument(&val, MakeType("bool"))
	} else {
		val := encoder.Serialize(int32(1))
		return MakeArgument(&val, MakeType("bool"))
	}
}

/*
  Time functions
*/

func sleep (ms *CXArgument) *CXArgument {
	if ms.Typ.Name != "i32" {
		panic("sleep: wrong argument type")
	}
	
	var duration int32
	encoder.DeserializeRaw(*ms.Value, &duration)

	time.Sleep(time.Duration(duration) * time.Millisecond)

	return ms
}


/*
  Prolog functions
*/

func setClauses (clss *CXArgument, mod *CXModule) *CXArgument {
	if clss.Typ.Name != "str" {
		panic("setClauses: wrong argument type")
	}

	clauses := string(*clss.Value)
	mod.AddClauses(clauses)

	return clss
}

func addObject (obj *CXArgument, mod *CXModule) *CXArgument {
	if obj.Typ.Name != "str" {
		panic("addObject: wrong argument type")
	}
	
	object := string(*obj.Value)
	mod.AddObject(object)

	return obj
}

func setQuery (qry *CXArgument, mod *CXModule) *CXArgument {
	if qry.Typ.Name != "str" {
		panic("setQuery: wrong argument type")
	}
	
	query := string(*qry.Value)
	mod.AddQuery(query)

	return qry
}

func remObject (obj *CXArgument, mod *CXModule) *CXArgument {
	if obj.Typ.Name != "str" {
		panic("remObject: wrong argument type")
	}
	
	object := string(*obj.Value)
	mod.RemoveObject(object)

	return obj
}

func remObjects (mod *CXModule) *CXArgument {
	mod.RemoveObjects()

	success := encoder.Serialize(int32(1))
	return MakeArgument(&success, MakeType("bool"))
}

func remExpr (fnName *CXArgument, ln *CXArgument, mod *CXModule) *CXArgument {
	if fnName.Typ.Name != "str" {
		panic("remExpr: wrong argument type")
	}

	var line int32
	encoder.DeserializeRaw(*ln.Value, &line)
	
	name := string(*fnName.Value)
	if fn, err := mod.Context.GetFunction(name, mod.Name); err == nil {
		fn.RemoveExpression(int(line))
		val := encoder.Serialize(int32(1))
		return MakeArgument(&val, MakeType("bool"))
	} else {
		val := encoder.Serialize(int32(0))
		return MakeArgument(&val, MakeType("bool"))
	}
}

func remArg (fnName *CXArgument, mod *CXModule) *CXArgument {
	if fnName.Typ.Name != "str" {
		panic("remArg: wrong argument type")
	}
	
	name := string(*fnName.Value)
	if fn, err := mod.Context.GetFunction(name, mod.Name); err == nil {
		if expr, err := fn.GetCurrentExpression(); err == nil {
			expr.RemoveArgument()
		}
		
		val := encoder.Serialize(int32(1))
		return MakeArgument(&val, MakeType("bool"))
	} else {
		val := encoder.Serialize(int32(0))
		return MakeArgument(&val, MakeType("bool"))
	}
}

func addExpr (fnName *CXArgument, caller *CXFunction) *CXArgument {
	if fnName.Typ.Name != "str" {
		panic("addExpr: wrong argument type")
	}

	mod := caller.Module
	
	opName := string(*fnName.Value)
	if fn, err := mod.Context.GetFunction(opName, mod.Name); err == nil {
		expr := MakeExpression(fn)
		caller.AddExpression(expr)
		
		val := encoder.Serialize(int32(1))
		return MakeArgument(&val, MakeType("bool"))
	} else {
		val := encoder.Serialize(int32(0))
		return MakeArgument(&val, MakeType("bool"))
	}
}

func exprAff (fltr *CXArgument, op *CXFunction) *CXArgument {
	filter := string(*fltr.Value)
	
	if expr, err := op.GetExpression(len(op.Expressions)-1); err == nil {
		affs := FilterAffordances(expr.GetAffordances(), filter)
		affs[0].ApplyAffordance()

		
		
		val := encoder.Serialize(int32(1))
		return MakeArgument(&val, MakeType("bool"))
	} else {
		val := encoder.Serialize(int32(0))
		return MakeArgument(&val, MakeType("bool"))
	}
}

func initDef (arg1 *CXArgument) *CXArgument {
	if arg1.Typ.Name != "str" {
		panic("initDef: wrong argument type")
	}

	typName := string(*arg1.Value)

	var zeroVal []byte
	switch  typName {
	case "bool": zeroVal = encoder.Serialize(int32(0))
	case "byte": zeroVal = []byte{byte(0)}
	case "i32": zeroVal = encoder.Serialize(int32(0))
	case "i64": zeroVal = encoder.Serialize(int64(0))
	case "f32": zeroVal = encoder.Serialize(float32(0))
	case "f64": zeroVal = encoder.Serialize(float64(0))
	case "[]bool": zeroVal = encoder.Serialize([]int32{0})
	case "[]byte": zeroVal = []byte{byte(0)}
	case "[]i32": zeroVal = encoder.Serialize([]int32{0})
	case "[]i64": zeroVal = encoder.Serialize([]int64{0})
	case "[]f32": zeroVal = encoder.Serialize([]float32{0})
	case "[]f64": zeroVal = encoder.Serialize([]float64{0})
	}

	return MakeArgument(&zeroVal, MakeType(typName))
}

// func randSeed (seed *CXArgument) *CXArgument {
// 	if seed.Typ.Name != "i64" {
// 		panic("randSeed: wrong argument type")
// 	}

	
// }

func randI32 (min *CXArgument, max *CXArgument) *CXArgument {
	if min.Typ.Name != "i32" || max.Typ.Name != "i32" {
		panic("randI32: wrong argument type")
	}

	//const n int = 312
	//const m int = 156





	
	var minimum int32
	encoder.DeserializeRaw(*min.Value, &minimum)

	var maximum int32
	encoder.DeserializeRaw(*max.Value, &maximum)

	if minimum > maximum {
		panic(fmt.Sprintf("random: min must be less than max (%d !< %d)", minimum, maximum))
	}

	rand.Seed(time.Now().UTC().UnixNano())
	output := encoder.SerializeAtomic(int32(rand.Intn(int(maximum - minimum)) + int(minimum)))
	return MakeArgument(&output, MakeType("i32"))
}

