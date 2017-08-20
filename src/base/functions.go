package base

import (
	"fmt"
	"time"
	"github.com/skycoin/skycoin/src/cipher/encoder"
)

func addI32 (arg1 *CXArgument, arg2 *CXArgument) *CXArgument {
	if arg1.Typ.Name != "i32" || arg2.Typ.Name != "i32" {
		if arg1.Typ.Name != "i32" {
			panic(fmt.Sprintf("addI32: first argument is type '%s'; expected type 'i32'", arg1.Typ.Name))
		}
		panic(fmt.Sprintf("addI32: second argument is type '%s'; expected type 'i32'", arg1.Typ.Name))
	}

	var num1 int32
	var num2 int32
	encoder.DeserializeAtomic(*arg1.Value, &num1)
	encoder.DeserializeAtomic(*arg2.Value, &num2)

	output := encoder.SerializeAtomic(num1 + num2)

	return &CXArgument{Value: &output, Typ: MakeType("i32")}
}

func subI32 (arg1 *CXArgument, arg2 *CXArgument) *CXArgument {
	if arg1.Typ.Name != "i32" || arg2.Typ.Name != "i32" {
		panic("subI32: wrong argument type")
	}
	
	var num1 int32
	var num2 int32
	encoder.DeserializeAtomic(*arg1.Value, &num1)
	encoder.DeserializeAtomic(*arg2.Value, &num2)

	output := encoder.SerializeAtomic(num1 - num2)

	return &CXArgument{Value: &output, Typ: MakeType("i32")}
}

func mulI32 (arg1 *CXArgument, arg2 *CXArgument) *CXArgument {
	if arg1.Typ.Name != "i32" || arg2.Typ.Name != "i32" {
		panic("mulI32: wrong argument type")
	}

	var num1 int32
	var num2 int32
	encoder.DeserializeAtomic(*arg1.Value, &num1)
	encoder.DeserializeAtomic(*arg2.Value, &num2)

	output := encoder.SerializeAtomic(num1 * num2)

	return &CXArgument{Value: &output, Typ: MakeType("i32")}
}

func initDef (arg1 *CXArgument) *CXArgument {
	if arg1.Typ.Name != "str" {
		panic("initDef: wrong argument type")
	}

	typName := string(*arg1.Value)

	var zeroVal []byte
	switch  typName {
	case "byte": zeroVal = []byte{byte(0)}
	case "i32": zeroVal = encoder.Serialize(int32(0))
	case "i64": zeroVal = encoder.Serialize(int64(0))
	case "f32": zeroVal = encoder.Serialize(float32(0))
	case "f64": zeroVal = encoder.Serialize(float64(0))
	case "[]byte": zeroVal = []byte{byte(0)}
	case "[]i32": zeroVal = encoder.Serialize([]int32{0})
	case "[]i64": zeroVal = encoder.Serialize([]int64{0})
	case "[]f32": zeroVal = encoder.Serialize([]float32{0})
	case "[]f64": zeroVal = encoder.Serialize([]float64{0})
	}

	return MakeArgument(&zeroVal, MakeType(typName))
}

func divI32 (arg1 *CXArgument, arg2 *CXArgument) *CXArgument {
	if arg1.Typ.Name != "i32" || arg2.Typ.Name != "i32" {
		panic("divI32: wrong argument type")
	}
	
	var num1 int32
	var num2 int32
	encoder.DeserializeAtomic(*arg1.Value, &num1)
	encoder.DeserializeAtomic(*arg2.Value, &num2)

	if num2 == int32(0) {
		panic("divI32: Division by 0")
	}

	output := encoder.SerializeAtomic(num1 / num2)

	return &CXArgument{Value: &output, Typ: MakeType("i32")}
}

func addI64 (arg1 *CXArgument, arg2 *CXArgument) *CXArgument {
	if arg1.Typ.Name != "i64" || arg2.Typ.Name != "i64" {
		panic("addI64: wrong argument type")
	}
	
	var num1 int64
	var num2 int64
	encoder.DeserializeAtomic(*arg1.Value, &num1)
	encoder.DeserializeAtomic(*arg2.Value, &num2)

	output := encoder.SerializeAtomic(num1 + num2)

	return &CXArgument{Value: &output, Typ: MakeType("i64")}
}

func subI64 (arg1 *CXArgument, arg2 *CXArgument) *CXArgument {
	if arg1.Typ.Name != "i64" || arg2.Typ.Name != "i64" {
		panic("subI64: wrong argument type")
	}
	
	var num1 int64
	var num2 int64
	encoder.DeserializeAtomic(*arg1.Value, &num1)
	encoder.DeserializeAtomic(*arg2.Value, &num2)

	output := encoder.SerializeAtomic(num1 - num2)

	return &CXArgument{Value: &output, Typ: MakeType("i64")}
}

func mulI64 (arg1 *CXArgument, arg2 *CXArgument) *CXArgument {
	if arg1.Typ.Name != "i64" || arg2.Typ.Name != "i64" {
		panic("mulI64: wrong argument type")
	}
	
	var num1 int64
	var num2 int64
	encoder.DeserializeAtomic(*arg1.Value, &num1)
	encoder.DeserializeAtomic(*arg2.Value, &num2)

	output := encoder.SerializeAtomic(num1 * num2)

	return &CXArgument{Value: &output, Typ: MakeType("i64")}
}

func divI64 (arg1 *CXArgument, arg2 *CXArgument) *CXArgument {
	if arg1.Typ.Name != "i64" || arg2.Typ.Name != "i64" {
		panic("divI64: wrong argument type")
	}

	var num1 int64
	var num2 int64
	encoder.DeserializeAtomic(*arg1.Value, &num1)
	encoder.DeserializeAtomic(*arg2.Value, &num2)

	if num2 == int64(0) {
		panic("divI64: Division by 0")
	}
	
	output := encoder.SerializeAtomic(num1 / num2)

	return &CXArgument{Value: &output, Typ: MakeType("i64")}
}

func addF32 (arg1 *CXArgument, arg2 *CXArgument) *CXArgument {
	if arg1.Typ.Name != "f32" || arg2.Typ.Name != "f32" {
		panic("addF32: wrong argument type")
	}
	
	var num1 float32
	var num2 float32
	encoder.DeserializeRaw(*arg1.Value, &num1)
	encoder.DeserializeRaw(*arg2.Value, &num2)

	output := encoder.Serialize(num1 + num2)

	return &CXArgument{Value: &output, Typ: MakeType("f32")}
}

func subF32 (arg1 *CXArgument, arg2 *CXArgument) *CXArgument {
	if arg1.Typ.Name != "f32" || arg2.Typ.Name != "f32" {
		panic("subF32: wrong argument type")
	}
	
	var num1 float32
	var num2 float32
	encoder.DeserializeRaw(*arg1.Value, &num1)
	encoder.DeserializeRaw(*arg2.Value, &num2)

	output := encoder.Serialize(num1 - num2)

	return &CXArgument{Value: &output, Typ: MakeType("f32")}
}

func mulF32 (arg1 *CXArgument, arg2 *CXArgument) *CXArgument {
	if arg1.Typ.Name != "f32" || arg2.Typ.Name != "f32" {
		panic("mulF32: wrong argument type")
	}
	
	var num1 float32
	var num2 float32
	encoder.DeserializeRaw(*arg1.Value, &num1)
	encoder.DeserializeRaw(*arg2.Value, &num2)

	output := encoder.Serialize(num1 * num2)

	return &CXArgument{Value: &output, Typ: MakeType("f32")}
}

func divF32 (arg1 *CXArgument, arg2 *CXArgument) *CXArgument {
	if arg1.Typ.Name != "f32" || arg2.Typ.Name != "f32" {
		panic("mulF32: wrong argument type")
	}
	
	var num1 float32
	var num2 float32
	encoder.DeserializeRaw(*arg1.Value, &num1)
	encoder.DeserializeRaw(*arg2.Value, &num2)

	if num2 == float32(0.0) {
		panic("divI32: Division by 0")
	}

	output := encoder.Serialize(num1 / num2)

	return &CXArgument{Value: &output, Typ: MakeType("f32")}
}

func addF64 (arg1 *CXArgument, arg2 *CXArgument) *CXArgument {
	if arg1.Typ.Name != "f64" || arg2.Typ.Name != "f64" {
		panic("addF64: wrong argument type")
	}
	
	var num1 float64
	var num2 float64
	encoder.DeserializeRaw(*arg1.Value, &num1)
	encoder.DeserializeRaw(*arg2.Value, &num2)

	output := encoder.Serialize(num1 + num2)

	return &CXArgument{Value: &output, Typ: MakeType("f64")}
}

func subF64 (arg1 *CXArgument, arg2 *CXArgument) *CXArgument {
	if arg1.Typ.Name != "f64" || arg2.Typ.Name != "f64" {
		panic("subF64: wrong argument type")
	}
	
	var num1 float64
	var num2 float64
	encoder.DeserializeRaw(*arg1.Value, &num1)
	encoder.DeserializeRaw(*arg2.Value, &num2)

	output := encoder.Serialize(num1 - num2)

	return &CXArgument{Value: &output, Typ: MakeType("f64")}
}

func mulF64 (arg1 *CXArgument, arg2 *CXArgument) *CXArgument {
	if arg1.Typ.Name != "f64" || arg2.Typ.Name != "f64" {
		panic("mulF64: wrong argument type")
	}
	
	var num1 float64
	var num2 float64
	encoder.DeserializeRaw(*arg1.Value, &num1)
	encoder.DeserializeRaw(*arg2.Value, &num2)

	output := encoder.Serialize(num1 * num2)

	return &CXArgument{Value: &output, Typ: MakeType("f64")}
}

func divF64 (arg1 *CXArgument, arg2 *CXArgument) *CXArgument {
	if arg1.Typ.Name != "f64" || arg2.Typ.Name != "f64" {
		panic("mulF64: wrong argument type")
	}
	
	var num1 float64
	var num2 float64
	encoder.DeserializeRaw(*arg1.Value, &num1)
	encoder.DeserializeRaw(*arg2.Value, &num2)

	if num2 == float64(0.0) {
		panic("divF64: Division by 0")
	}

	output := encoder.Serialize(num1 / num2)

	return &CXArgument{Value: &output, Typ: MakeType("f64")}
}

func readByteA (arr *CXArgument, idx *CXArgument) *CXArgument {
	if arr.Typ.Name != "[]byte" || idx.Typ.Name != "i32" {
		panic("readByteA: wrong argument type")
	}

	var index int32
	encoder.DeserializeAtomic(*idx.Value, &index)

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
	encoder.DeserializeAtomic(*idx.Value, &index)
	
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
	encoder.DeserializeAtomic(*idx.Value, &index)

	var array []int32
	encoder.DeserializeRaw(*arr.Value, &array)
	
	if index >= int32(len(array)) {
		panic(fmt.Sprintf("readI32A: Index %d exceeds array of length %d", index, len(array)))
	}

	output := encoder.Serialize(array[index])

	return &CXArgument{Value: &output, Typ: MakeType("i32")}
}

func writeI32A (arr *CXArgument, idx *CXArgument, val *CXArgument) *CXArgument {
	if arr.Typ.Name != "[]i32" || idx.Typ.Name != "i32" || val.Typ.Name != "i32" {
		panic("readI32A: wrong argument type")
	}

	var index int32
	encoder.DeserializeAtomic(*idx.Value, &index)

	var value int32
	encoder.DeserializeAtomic(*val.Value, &value)

	var array []int32
	encoder.DeserializeRaw(*arr.Value, &array)
	
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
	encoder.DeserializeAtomic(*idx.Value, &index)

	var array []int64
	encoder.DeserializeRaw(*arr.Value, &array)
	
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
	encoder.DeserializeAtomic(*idx.Value, &index)

	var value int64
	encoder.DeserializeAtomic(*val.Value, &value)

	var array []int64
	encoder.DeserializeRaw(*arr.Value, &array)
	
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
	encoder.DeserializeAtomic(*idx.Value, &index)

	var array []float32
	encoder.DeserializeRaw(*arr.Value, &array)
	
	if index >= int32(len(array)) {
		panic(fmt.Sprintf("readF32A: Index %d exceeds array of length %d", index, len(array)))
	}

	output := encoder.Serialize(array[index])

	return &CXArgument{Value: &output, Typ: MakeType("f32")}
}

func writeF32A (arr *CXArgument, idx *CXArgument, val *CXArgument) *CXArgument {
	if arr.Typ.Name != "[]f32" || idx.Typ.Name != "i32" || val.Typ.Name != "f32" {
		panic("readF32A: wrong argument type")
	}

	var index int32
	encoder.DeserializeAtomic(*idx.Value, &index)

	var value float32
	encoder.DeserializeAtomic(*val.Value, &value)

	var array []float32
	encoder.DeserializeRaw(*arr.Value, &array)
	
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
	encoder.DeserializeAtomic(*idx.Value, &index)

	var array []float64
	encoder.DeserializeRaw(*arr.Value, &array)
	
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
	encoder.DeserializeAtomic(*idx.Value, &index)

	var value float64
	encoder.DeserializeAtomic(*val.Value, &value)

	var array []float64
	encoder.DeserializeRaw(*arr.Value, &array)
	
	if index >= int32(len(array)) {
		panic(fmt.Sprintf("writeF64A: Index %d exceeds array of length %d", index, len(array)))
	}

	array[index] = value
	output := encoder.Serialize(array)
	arr.Value = &output

	return arr
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
	encoder.DeserializeAtomic(*arg1.Value, &num1)
	encoder.DeserializeAtomic(*arg2.Value, &num2)
	
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
	encoder.DeserializeAtomic(*arg1.Value, &num1)
	encoder.DeserializeAtomic(*arg2.Value, &num2)
	
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
	encoder.DeserializeAtomic(*arg1.Value, &num1)
	encoder.DeserializeAtomic(*arg2.Value, &num2)
	
	if num1 == num2 {
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
	encoder.DeserializeAtomic(*arg1.Value, &num1)
	encoder.DeserializeAtomic(*arg2.Value, &num2)
	
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
	encoder.DeserializeAtomic(*arg1.Value, &num1)
	encoder.DeserializeAtomic(*arg2.Value, &num2)
	
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
	encoder.DeserializeAtomic(*arg1.Value, &num1)
	encoder.DeserializeAtomic(*arg2.Value, &num2)
	
	if num1 == num2 {
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
  MetaProgramming Natives
*/

// For now, let's just add these
// Prolog
// :clauses
// :objects
// :object
// :query
// :rem
