package base

import (
	"fmt"
	"time"
	"errors"
	"math/rand"
	"github.com/skycoin/skycoin/src/cipher/encoder"
)

func checkType (fnName string, typ string, arg *CXArgument) error {
	if arg.Typ.Name != typ {
		return errors.New(fmt.Sprintf("%s: argument 1 is type '%s'; expected type '%s'", fnName, arg.Typ.Name, typ))
	}
	return nil
}

func checkTwoTypes (fnName string, typ1 string, typ2 string, arg1 *CXArgument, arg2 *CXArgument) error {
	if arg1.Typ.Name != typ1 || arg2.Typ.Name != typ2 {
		if arg1.Typ.Name != typ1 {
			return errors.New(fmt.Sprintf("%s: argument 1 is type '%s'; expected type '%s'", fnName, arg1.Typ.Name, typ1))
		}
		return errors.New(fmt.Sprintf("%s: argument 2 is type '%s'; expected type '%s'", fnName, arg2.Typ.Name, typ2))
	}
	return nil
}

func checkThreeTypes (fnName string, typ1 string, typ2 string, typ3 string, arg1 *CXArgument, arg2 *CXArgument, arg3 *CXArgument) error {
	if arg1.Typ.Name != typ1 || arg2.Typ.Name != typ2 || arg3.Typ.Name != typ3 {
		if arg1.Typ.Name != typ1 {
			return errors.New(fmt.Sprintf("%s: argument 1 is type '%s'; expected type '%s'", fnName, arg1.Typ.Name, typ1))
		} else if arg2.Typ.Name != typ2 {
			return errors.New(fmt.Sprintf("%s: argument 2 is type '%s'; expected type '%s'", fnName, arg2.Typ.Name, typ2))
		}
		return errors.New(fmt.Sprintf("%s: argument 3 is type '%s'; expected type '%s'", fnName, arg3.Typ.Name, typ3))
	}
	return nil
}

func addI32 (arg1 *CXArgument, arg2 *CXArgument) (*CXArgument, error) {
	if err := checkTwoTypes("addI32", "i32", "i32", arg1, arg2); err == nil {
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
	if err := checkTwoTypes("subI32", "i32", "i32", arg1, arg2); err == nil {
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
	if err := checkTwoTypes("mulI32", "i32", "i32", arg1, arg2); err == nil {
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
	if err := checkTwoTypes("divI32", "i32", "i32", arg1, arg2); err == nil {
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
	if err := checkTwoTypes("addI64", "i64", "i64", arg1, arg2); err == nil {
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
	if err := checkTwoTypes("subI64", "i64", "i64", arg1, arg2); err == nil {
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
	if err := checkTwoTypes("mulI64", "i64", "i64", arg1, arg2); err == nil {
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
	if err := checkTwoTypes("divI64", "i64", "i64", arg1, arg2); err == nil {
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
	if err := checkTwoTypes("addF32", "f32", "f32", arg1, arg2); err == nil {
		
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
	if err := checkTwoTypes("subF32", "f32", "f32", arg1, arg2); err == nil {
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
	if err := checkTwoTypes("mulF32", "f32", "f32", arg1, arg2); err == nil {
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
	if err := checkTwoTypes("divF32", "f32", "f32", arg1, arg2); err == nil {
		var num1 float32
		var num2 float32
		encoder.DeserializeRaw(*arg1.Value, &num1)
		encoder.DeserializeRaw(*arg2.Value, &num2)

		if num2 == float32(0.0) {
			return nil, errors.New("divF32: Division by 0")
		}

		output := encoder.Serialize(float32(num1 / num2))

		return &CXArgument{Value: &output, Typ: MakeType("f32")}, nil
	} else {
		return nil, err
	}	
}

func addF64 (arg1 *CXArgument, arg2 *CXArgument) (*CXArgument, error) {
	if err := checkTwoTypes("addF64", "f64", "f64", arg1, arg2); err == nil {
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
	if err := checkTwoTypes("subF64", "f64", "f64", arg1, arg2); err == nil {
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
	if err := checkTwoTypes("mulF64", "f64", "f64", arg1, arg2); err == nil {
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
	if err := checkTwoTypes("divF64", "f64", "f64", arg1, arg2); err == nil {
		var num1 float64
		var num2 float64
		encoder.DeserializeRaw(*arg1.Value, &num1)
		encoder.DeserializeRaw(*arg2.Value, &num2)

		if num2 == float64(0.0) {
			return nil, errors.New("divF64: Division by 0")
		}

		output := encoder.Serialize(float64(num1 / num2))

		return &CXArgument{Value: &output, Typ: MakeType("f64")}, nil
	} else {
		return nil, err
	}
}

func modI32 (arg1 *CXArgument, arg2 *CXArgument) (*CXArgument, error) {
	if err := checkTwoTypes("modI32", "i32", "i32", arg1, arg2); err == nil {
		var num1 int32
		var num2 int32
		encoder.DeserializeRaw(*arg1.Value, &num1)
		encoder.DeserializeRaw(*arg2.Value, &num2)

		if num2 == int32(0) {
			return nil, errors.New("modI32: Division by 0")
		}

		output := encoder.Serialize(int32(num1 % num2))

		return &CXArgument{Value: &output, Typ: MakeType("i32")}, nil
	} else {
		return nil, err
	}
}

func modI64 (arg1 *CXArgument, arg2 *CXArgument) (*CXArgument, error) {
	if err := checkTwoTypes("modI64", "i64", "i64", arg1, arg2); err == nil {
		var num1 int64
		var num2 int64
		encoder.DeserializeRaw(*arg1.Value, &num1)
		encoder.DeserializeRaw(*arg2.Value, &num2)

		if num2 == int64(0) {
			return nil, errors.New("modI64: Division by 0")
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
	if err := checkTwoTypes("andI32", "i32", "i32", arg1, arg2); err == nil {
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
	if err := checkTwoTypes("orI32", "i32", "i32", arg1, arg2); err == nil {
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
	if err := checkTwoTypes("xorI32", "i32", "i32", arg1, arg2); err == nil {
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
	if err := checkTwoTypes("andNotI32", "i32", "i32", arg1, arg2); err == nil {
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
	if err := checkTwoTypes("andI64", "i64", "i64", arg1, arg2); err == nil {
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
	if err := checkTwoTypes("orI64", "i64", "i64", arg1, arg2); err == nil {
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
	if err := checkTwoTypes("xorI64", "i64", "i64", arg1, arg2); err == nil {
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
	if err := checkTwoTypes("andNotI64", "i64", "i64", arg1, arg2); err == nil {
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

func readBoolA (arr *CXArgument, idx *CXArgument) (*CXArgument, error) {
	if err := checkTwoTypes("readBoolA", "[]bool", "i32", arr, idx); err == nil {
		var index int32
		encoder.DeserializeRaw(*idx.Value, &index)

		var array []int32
		encoder.DeserializeRaw(*arr.Value, &array)

		if index < 0 {
			return nil, errors.New(fmt.Sprintf("readBoolA: negative index %d", index))
		}
		
		if index >= int32(len(array)) {
			return nil, errors.New(fmt.Sprintf("readBoolA: index %d exceeds array of length %d", index, len(array)))
		}

		output := encoder.Serialize(array[index])

		return &CXArgument{Value: &output, Typ: MakeType("bool")}, nil
	} else {
		return nil, err
	}
}

func writeBoolA (arr *CXArgument, idx *CXArgument, val *CXArgument) (*CXArgument, error) {
	if err := checkThreeTypes("writeBoolA", "[]bool", "i32", "bool", arr, idx, val); err == nil {
		var index int32
		encoder.DeserializeRaw(*idx.Value, &index)

		var value int32
		encoder.DeserializeRaw(*val.Value, &value)
		
		var array []int32
		encoder.DeserializeRaw(*arr.Value, &array)

		if index < 0 {
			return nil, errors.New(fmt.Sprintf("writeBoolA: negative index %d", index))
		}
		
		if index >= int32(len(*arr.Value)) {
			return nil, errors.New(fmt.Sprintf("writeBoolA: index %d exceeds array of length %d", index, len(array)))
		}

		array[index] = value
		output := encoder.Serialize(array)
		arr.Value = &output

		return arr, nil
	} else {
		return nil, err
	}
}

func readByteA (arr *CXArgument, idx *CXArgument) (*CXArgument, error) {
	if err := checkTwoTypes("readByteA", "[]byte", "i32", arr, idx); err == nil {
		var index int32
		encoder.DeserializeRaw(*idx.Value, &index)

		if index < 0 {
			return nil, errors.New(fmt.Sprintf("readByteA: negative index %d", index))
		}
		
		if index >= int32(len(*arr.Value)) {
			return nil, errors.New(fmt.Sprintf("readByteA: index %d exceeds array of length %d", index, len(*arr.Value)))
		}

		output := make([]byte, 1)
		output[0] = (*arr.Value)[index]

		return &CXArgument{Value: &output, Typ: MakeType("byte")}, nil
	} else {
		return nil, err
	}
}

func writeByteA (arr *CXArgument, idx *CXArgument, val *CXArgument) (*CXArgument, error) {
	if err := checkThreeTypes("writeByteA", "[]byte", "i32", "byte", arr, idx, val); err == nil {
		var index int32
		encoder.DeserializeRaw(*idx.Value, &index)

		if index < 0 {
			return nil, errors.New(fmt.Sprintf("writeByteA: negative index %d", index))
		}
		
		if index >= int32(len(*arr.Value)) {
			return nil, errors.New(fmt.Sprintf("writeByteA: index %d exceeds array of length %d", index, len(*arr.Value)))
		}

		(*arr.Value)[index] = (*val.Value)[0]

		return arr, nil
	} else {
		return nil, err
	}
}

func readI32A (arr *CXArgument, idx *CXArgument) (*CXArgument, error) {
	if err := checkTwoTypes("readI32A", "[]i32", "i32", arr, idx); err == nil {
		var index int32
		encoder.DeserializeRaw(*idx.Value, &index)

		var array []int32
		encoder.DeserializeRaw(*arr.Value, &array)

		if index < 0 {
			return nil, errors.New(fmt.Sprintf("readI32A: negative index %d", index))
		}
		
		if index >= int32(len(array)) {
			return nil, errors.New(fmt.Sprintf("readI32A: index %d exceeds array of length %d", index, len(array)))
		}

		output := encoder.Serialize(array[index])

		return &CXArgument{Value: &output, Typ: MakeType("i32")}, nil
	} else {
		return nil, err
	}
}

func writeI32A (arr *CXArgument, idx *CXArgument, val *CXArgument) (*CXArgument, error) {
	if err := checkThreeTypes("writeI32A", "[]i32", "i32", "i32", arr, idx, val); err == nil {
		var index int32
		encoder.DeserializeRaw(*idx.Value, &index)

		var value int32
		encoder.DeserializeRaw(*val.Value, &value)

		var array []int32
		encoder.DeserializeRaw(*arr.Value, &array)

		if index < 0 {
			return nil, errors.New(fmt.Sprintf("writeI32A: negative index %d", index))
		}

		if index >= int32(len(array)) {
			return nil, errors.New(fmt.Sprintf("writeI32A: index %d exceeds array of length %d", index, len(array)))
		}

		array[index] = value
		output := encoder.Serialize(array)
		arr.Value = &output

		return arr, nil
	} else {
		return nil, err
	}
}

func readI64A (arr *CXArgument, idx *CXArgument) (*CXArgument, error) {
	if err := checkTwoTypes("readI64A", "[]i64", "i32", arr, idx); err == nil {
		var index int32
		encoder.DeserializeRaw(*idx.Value, &index)

		var array []int64
		encoder.DeserializeRaw(*arr.Value, &array)

		if index < 0 {
			return nil, errors.New(fmt.Sprintf("readI64A: negative index %d", index))
		}

		if index >= int32(len(array)) {
			return nil, errors.New(fmt.Sprintf("readI64A: index %d exceeds array of length %d", index, len(array)))
		}

		output := encoder.Serialize(array[index])

		return &CXArgument{Value: &output, Typ: MakeType("i64")}, nil
	} else {
		return nil, err
	}
}

func writeI64A (arr *CXArgument, idx *CXArgument, val *CXArgument) (*CXArgument, error) {
	if err := checkThreeTypes("writeI64A", "[]i64", "i32", "i64", arr, idx, val); err == nil {
		var index int32
		encoder.DeserializeRaw(*idx.Value, &index)

		var value int64
		encoder.DeserializeRaw(*val.Value, &value)

		var array []int64
		encoder.DeserializeRaw(*arr.Value, &array)

		if index < 0 {
			return nil, errors.New(fmt.Sprintf("writeI64A: negative index %d", index))
		}

		if index >= int32(len(array)) {
			return nil, errors.New(fmt.Sprintf("writeI64A: index %d exceeds array of length %d", index, len(array)))
		}

		array[index] = value
		output := encoder.Serialize(array)
		arr.Value = &output

		return arr, nil
	} else {
		return nil, err
	}
}

func readF32A (arr *CXArgument, idx *CXArgument) (*CXArgument, error) {
	if err := checkTwoTypes("readF32A", "[]f32", "i32", arr, idx); err == nil {
		var index int32
		encoder.DeserializeRaw(*idx.Value, &index)

		var array []float32
		encoder.DeserializeRaw(*arr.Value, &array)

		if index < 0 {
			return nil, errors.New(fmt.Sprintf("readF32A: negative index %d", index))
		}

		if index >= int32(len(array)) {
			return nil, errors.New(fmt.Sprintf("readF32A: index %d exceeds array of length %d", index, len(array)))
		}

		output := encoder.Serialize(array[index])

		return &CXArgument{Value: &output, Typ: MakeType("f32")}, nil
	} else {
		return nil, err
	}
}

func writeF32A (arr *CXArgument, idx *CXArgument, val *CXArgument) (*CXArgument, error) {
	if err := checkThreeTypes("writeF32A", "[]f32", "i32", "f32", arr, idx, val); err == nil {
		var index int32
		encoder.DeserializeRaw(*idx.Value, &index)

		var value float32
		encoder.DeserializeRaw(*val.Value, &value)

		var array []float32
		encoder.DeserializeRaw(*arr.Value, &array)

		if index < 0 {
			return nil, errors.New(fmt.Sprintf("writeF32A: negative index %d", index))
		}

		if index >= int32(len(array)) {
			return nil, errors.New(fmt.Sprintf("writeF32A: index %d exceeds array of length %d", index, len(array)))
		}

		array[index] = value
		output := encoder.Serialize(array)
		arr.Value = &output

		return arr, nil
	} else {
		return nil, err
	}
}

func readF64A (arr *CXArgument, idx *CXArgument) (*CXArgument, error) {
	if err := checkTwoTypes("readF64A", "[]f64", "i32", arr, idx); err == nil {
		var index int32
		encoder.DeserializeRaw(*idx.Value, &index)

		var array []float64
		encoder.DeserializeRaw(*arr.Value, &array)

		if index < 0 {
			return nil, errors.New(fmt.Sprintf("readF64A: negative index %d", index))
		}

		if index >= int32(len(array)) {
			return nil, errors.New(fmt.Sprintf("readF64A: index %d exceeds array of length %d", index, len(array)))
		}

		output := encoder.Serialize(array[index])

		return &CXArgument{Value: &output, Typ: MakeType("f64")}, nil
	} else {
		return nil, err
	}
}

func writeF64A (arr *CXArgument, idx *CXArgument, val *CXArgument) (*CXArgument, error) {
	if err := checkThreeTypes("writeF64A", "[]f64", "i32", "f64", arr, idx, val); err == nil {
		var index int32
		encoder.DeserializeRaw(*idx.Value, &index)

		var value float64
		encoder.DeserializeRaw(*val.Value, &value)

		var array []float64
		encoder.DeserializeRaw(*arr.Value, &array)

		if index < 0 {
			return nil, errors.New(fmt.Sprintf("writeF64A: negative index %d", index))
		}

		if index >= int32(len(array)) {
			return nil, errors.New(fmt.Sprintf("writeF64A: index %d exceeds array of length %d", index, len(array)))
		}

		array[index] = value
		output := encoder.Serialize(array)
		arr.Value = &output

		return arr, nil
	} else {
		return nil, err
	}
}

func lenBoolA (arr *CXArgument) (*CXArgument, error) {
	if err := checkType("lenBoolA", "[]bool", arr); err == nil {
		var array []int32
		encoder.DeserializeRaw(*arr.Value, &array)

		output := encoder.SerializeAtomic(int32(len(array)))
		return &CXArgument{Value: &output, Typ: MakeType("i32")}, nil
	} else {
		return nil, err
	}
}

func lenByteA (arr *CXArgument) (*CXArgument, error) {
	if err := checkType("lenByteA", "[]byte", arr); err == nil {
		output := encoder.SerializeAtomic(int32(len(*arr.Value)))
		
		return &CXArgument{Value: &output, Typ: MakeType("i32")}, nil
	} else {
		return nil, err
	}
}

func lenI32A (arr *CXArgument) (*CXArgument, error) {
	if err := checkType("lenI32A", "[]i32", arr); err == nil {
		var array []int32
		encoder.DeserializeRaw(*arr.Value, &array)

		output := encoder.SerializeAtomic(int32(len(array)))
		return &CXArgument{Value: &output, Typ: MakeType("i32")}, nil
	} else {
		return nil, err
	}
}

func lenI64A (arr *CXArgument) (*CXArgument, error) {
	if err := checkType("lenI64A", "[]i64", arr); err == nil {
		var array []int64
		encoder.DeserializeRaw(*arr.Value, &array)

		output := encoder.SerializeAtomic(int32(len(array)))
		return &CXArgument{Value: &output, Typ: MakeType("i32")}, nil
	} else {
		return nil, err
	}
}

func lenF32A (arr *CXArgument) (*CXArgument, error) {
	if err := checkType("lenF32A", "[]f32", arr); err == nil {
		var array []float32
		encoder.DeserializeRaw(*arr.Value, &array)

		output := encoder.SerializeAtomic(int32(len(array)))
		return &CXArgument{Value: &output, Typ: MakeType("i32")}, nil
	} else {
		return nil, err
	}
}

func lenF64A (arr *CXArgument) (*CXArgument, error) {
	if err := checkType("lenF64A", "[]f64", arr); err == nil {
		var array []float64
		encoder.DeserializeRaw(*arr.Value, &array)

		output := encoder.SerializeAtomic(int32(len(array)))
		return &CXArgument{Value: &output, Typ: MakeType("i32")}, nil
	} else {
		return nil, err
	}
}

/*
  Logical Operators
*/

func and (arg1 *CXArgument, arg2 *CXArgument) (*CXArgument, error) {
	if err := checkTwoTypes("and", "bool", "bool", arg1, arg2); err == nil {
		var c1 int32
		var c2 int32
		encoder.DeserializeRaw(*arg1.Value, &c1)
		encoder.DeserializeRaw(*arg2.Value, &c2)

		if c1 == 1 && c2 == 1 {
			val := encoder.Serialize(int32(1))
			return &CXArgument{Value: &val, Typ: MakeType("bool")}, nil
		} else {
			val := encoder.Serialize(int32(0))
			return &CXArgument{Value: &val, Typ: MakeType("bool")}, nil
		}
	} else {
		return nil, err
	}
}

func or (arg1 *CXArgument, arg2 *CXArgument) (*CXArgument, error) {
	if err := checkTwoTypes("or", "bool", "bool", arg1, arg2); err == nil {
		var c1 int32
		var c2 int32
		encoder.DeserializeRaw(*arg1.Value, &c1)
		encoder.DeserializeRaw(*arg2.Value, &c2)

		if c1 == 1 || c2 == 1 {
			val := encoder.Serialize(int32(1))
			return &CXArgument{Value: &val, Typ: MakeType("bool")}, nil
		} else {
			val := encoder.Serialize(int32(0))
			return &CXArgument{Value: &val, Typ: MakeType("bool")}, nil
		}
	} else {
		return nil, err
	}
}

func not (arg1 *CXArgument) (*CXArgument, error) {
	if err := checkType("not", "bool", arg1); err == nil {
		var c1 int32
		encoder.DeserializeRaw(*arg1.Value, &c1)

		if c1 == 0 {
			val := encoder.Serialize(int32(1))
			return &CXArgument{Value: &val, Typ: MakeType("bool")}, nil
		} else {
			val := encoder.Serialize(int32(0))
			return &CXArgument{Value: &val, Typ: MakeType("bool")}, nil
		}
	} else {
		return nil, err
	}
}

/*
  Relational Operators
*/

func ltI32 (arg1 *CXArgument, arg2 *CXArgument) (*CXArgument, error) {
	if err := checkTwoTypes("ltI32", "i32", "i32", arg1, arg2); err == nil {
		var num1 int32
		var num2 int32
		encoder.DeserializeRaw(*arg1.Value, &num1)
		encoder.DeserializeRaw(*arg2.Value, &num2)
		
		if num1 < num2 {
			val := encoder.Serialize(int32(1))
			return &CXArgument{Value: &val, Typ: MakeType("bool")}, nil
		} else {
			val := encoder.Serialize(int32(0))
			return &CXArgument{Value: &val, Typ: MakeType("bool")}, nil
		}
	} else {
		return nil, err
	}
}

func gtI32 (arg1 *CXArgument, arg2 *CXArgument) (*CXArgument, error) {
	if err := checkTwoTypes("gtI32", "i32", "i32", arg1, arg2); err == nil {
		var num1 int32
		var num2 int32
		encoder.DeserializeRaw(*arg1.Value, &num1)
		encoder.DeserializeRaw(*arg2.Value, &num2)
		
		if num1 > num2 {
			val := encoder.Serialize(int32(1))
			return &CXArgument{Value: &val, Typ: MakeType("bool")}, nil
		} else {
			val := encoder.Serialize(int32(0))
			return &CXArgument{Value: &val, Typ: MakeType("bool")}, nil
		}
	} else {
		return nil, err
	}
}

func eqI32 (arg1 *CXArgument, arg2 *CXArgument) (*CXArgument, error) {
	if err := checkTwoTypes("eqI32", "i32", "i32", arg1, arg2); err == nil {
		var num1 int32
		var num2 int32
		encoder.DeserializeRaw(*arg1.Value, &num1)
		encoder.DeserializeRaw(*arg2.Value, &num2)
		
		if num1 == num2 {
			val := encoder.Serialize(int32(1))
			return &CXArgument{Value: &val, Typ: MakeType("bool")}, nil
		} else {
			val := encoder.Serialize(int32(0))
			return &CXArgument{Value: &val, Typ: MakeType("bool")}, nil
		}
	} else {
		return nil, err
	}
}

func lteqI32 (arg1 *CXArgument, arg2 *CXArgument) (*CXArgument, error) {
	if err := checkTwoTypes("lteqI32", "i32", "i32", arg1, arg2); err == nil {
		var num1 int32
		var num2 int32
		encoder.DeserializeRaw(*arg1.Value, &num1)
		encoder.DeserializeRaw(*arg2.Value, &num2)
		
		if num1 <= num2 {
			val := encoder.Serialize(int32(1))
			return &CXArgument{Value: &val, Typ: MakeType("bool")}, nil
		} else {
			val := encoder.Serialize(int32(0))
			return &CXArgument{Value: &val, Typ: MakeType("bool")}, nil
		}
	} else {
		return nil, err
	}
}

func gteqI32 (arg1 *CXArgument, arg2 *CXArgument) (*CXArgument, error) {
	if err := checkTwoTypes("gteqI32", "i32", "i32", arg1, arg2); err == nil {
		var num1 int32
		var num2 int32
		encoder.DeserializeRaw(*arg1.Value, &num1)
		encoder.DeserializeRaw(*arg2.Value, &num2)
		
		if num1 >= num2 {
			val := encoder.Serialize(int32(1))
			return &CXArgument{Value: &val, Typ: MakeType("bool")}, nil
		} else {
			val := encoder.Serialize(int32(0))
			return &CXArgument{Value: &val, Typ: MakeType("bool")}, nil
		}
	} else {
		return nil, err
	}
}

func ltI64 (arg1 *CXArgument, arg2 *CXArgument) (*CXArgument, error) {
	if err := checkTwoTypes("ltI64", "i64", "i64", arg1, arg2); err == nil {
		var num1 int64
		var num2 int64
		encoder.DeserializeRaw(*arg1.Value, &num1)
		encoder.DeserializeRaw(*arg2.Value, &num2)
		
		if num1 < num2 {
			val := encoder.Serialize(int32(1))
			return &CXArgument{Value: &val, Typ: MakeType("bool")}, nil
		} else {
			val := encoder.Serialize(int32(0))
			return &CXArgument{Value: &val, Typ: MakeType("bool")}, nil
		}
	} else {
		return nil, err
	}
}

func gtI64 (arg1 *CXArgument, arg2 *CXArgument) (*CXArgument, error) {
	if err := checkTwoTypes("gtI64", "i64", "i64", arg1, arg2); err == nil {
		var num1 int64
		var num2 int64
		encoder.DeserializeRaw(*arg1.Value, &num1)
		encoder.DeserializeRaw(*arg2.Value, &num2)
		
		if num1 > num2 {
			val := encoder.Serialize(int32(1))
			return &CXArgument{Value: &val, Typ: MakeType("bool")}, nil
		} else {
			val := encoder.Serialize(int32(0))
			return &CXArgument{Value: &val, Typ: MakeType("bool")}, nil
		}
	} else {
		return nil, err
	}
}

func eqI64 (arg1 *CXArgument, arg2 *CXArgument) (*CXArgument, error) {
	if err := checkTwoTypes("eqI64", "i64", "i64", arg1, arg2); err == nil {
		var num1 int64
		var num2 int64
		encoder.DeserializeRaw(*arg1.Value, &num1)
		encoder.DeserializeRaw(*arg2.Value, &num2)
		
		if num1 == num2 {
			val := encoder.Serialize(int32(1))
			return &CXArgument{Value: &val, Typ: MakeType("bool")}, nil
		} else {
			val := encoder.Serialize(int32(0))
			return &CXArgument{Value: &val, Typ: MakeType("bool")}, nil
		}
	} else {
		return nil, err
	}
}

func lteqI64 (arg1 *CXArgument, arg2 *CXArgument) (*CXArgument, error) {
	if err := checkTwoTypes("lteqI64", "i64", "i64", arg1, arg2); err == nil {
		var num1 int64
		var num2 int64
		encoder.DeserializeRaw(*arg1.Value, &num1)
		encoder.DeserializeRaw(*arg2.Value, &num2)
		
		if num1 <= num2 {
			val := encoder.Serialize(int32(1))
			return &CXArgument{Value: &val, Typ: MakeType("bool")}, nil
		} else {
			val := encoder.Serialize(int32(0))
			return &CXArgument{Value: &val, Typ: MakeType("bool")}, nil
		}
	} else {
		return nil, err
	}
}

func gteqI64 (arg1 *CXArgument, arg2 *CXArgument) (*CXArgument, error) {
	if err := checkTwoTypes("gteqI64", "i64", "i64", arg1, arg2); err == nil {
		var num1 int64
		var num2 int64
		encoder.DeserializeRaw(*arg1.Value, &num1)
		encoder.DeserializeRaw(*arg2.Value, &num2)
		
		if num1 >= num2 {
			val := encoder.Serialize(int32(1))
			return &CXArgument{Value: &val, Typ: MakeType("bool")}, nil
		} else {
			val := encoder.Serialize(int32(0))
			return &CXArgument{Value: &val, Typ: MakeType("bool")}, nil
		}
	} else {
		return nil, err
	}
}

func ltF32 (arg1 *CXArgument, arg2 *CXArgument) (*CXArgument, error) {
	if err := checkTwoTypes("ltF32", "f32", "f32", arg1, arg2); err == nil {
		var num1 float32
		var num2 float32
		encoder.DeserializeRaw(*arg1.Value, &num1)
		encoder.DeserializeRaw(*arg2.Value, &num2)
		
		if num1 < num2 {
			val := encoder.Serialize(int32(1))
			return &CXArgument{Value: &val, Typ: MakeType("bool")}, nil
		} else {
			val := encoder.Serialize(int32(0))
			return &CXArgument{Value: &val, Typ: MakeType("bool")}, nil
		}
	} else {
		return nil, err
	}
}

func gtF32 (arg1 *CXArgument, arg2 *CXArgument) (*CXArgument, error) {
	if err := checkTwoTypes("gtF32", "f32", "f32", arg1, arg2); err == nil {
		var num1 float32
		var num2 float32
		encoder.DeserializeRaw(*arg1.Value, &num1)
		encoder.DeserializeRaw(*arg2.Value, &num2)
		
		if num1 > num2 {
			val := encoder.Serialize(int32(1))
			return &CXArgument{Value: &val, Typ: MakeType("bool")}, nil
		} else {
			val := encoder.Serialize(int32(0))
			return &CXArgument{Value: &val, Typ: MakeType("bool")}, nil
		}
	} else {
		return nil, err
	}
}

func eqF32 (arg1 *CXArgument, arg2 *CXArgument) (*CXArgument, error) {
	if err := checkTwoTypes("eqF32", "f32", "f32", arg1, arg2); err == nil {
		var num1 float32
		var num2 float32
		encoder.DeserializeRaw(*arg1.Value, &num1)
		encoder.DeserializeRaw(*arg2.Value, &num2)
		
		if num1 == num2 {
			val := encoder.Serialize(int32(1))
			return &CXArgument{Value: &val, Typ: MakeType("bool")}, nil
		} else {
			val := encoder.Serialize(int32(0))
			return &CXArgument{Value: &val, Typ: MakeType("bool")}, nil
		}
	} else {
		return nil, err
	}
}

func lteqF32 (arg1 *CXArgument, arg2 *CXArgument) (*CXArgument, error) {
	if err := checkTwoTypes("lteqF32", "f32", "f32", arg1, arg2); err == nil {
		var num1 float32
		var num2 float32
		encoder.DeserializeRaw(*arg1.Value, &num1)
		encoder.DeserializeRaw(*arg2.Value, &num2)
		
		if num1 <= num2 {
			val := encoder.Serialize(int32(1))
			return &CXArgument{Value: &val, Typ: MakeType("bool")}, nil
		} else {
			val := encoder.Serialize(int32(0))
			return &CXArgument{Value: &val, Typ: MakeType("bool")}, nil
		}
	} else {
		return nil, err
	}
}

func gteqF32 (arg1 *CXArgument, arg2 *CXArgument) (*CXArgument, error) {
	if err := checkTwoTypes("gteqF32", "f32", "f32", arg1, arg2); err == nil {
		var num1 float32
		var num2 float32
		encoder.DeserializeRaw(*arg1.Value, &num1)
		encoder.DeserializeRaw(*arg2.Value, &num2)
		
		if num1 >= num2 {
			val := encoder.Serialize(int32(1))
			return &CXArgument{Value: &val, Typ: MakeType("bool")}, nil
		} else {
			val := encoder.Serialize(int32(0))
			return &CXArgument{Value: &val, Typ: MakeType("bool")}, nil
		}
	} else {
		return nil, err
	}
}

func ltF64 (arg1 *CXArgument, arg2 *CXArgument) (*CXArgument, error) {
	if err := checkTwoTypes("ltF64", "f64", "f64", arg1, arg2); err == nil {
		var num1 float64
		var num2 float64
		encoder.DeserializeRaw(*arg1.Value, &num1)
		encoder.DeserializeRaw(*arg2.Value, &num2)
		
		if num1 < num2 {
			val := encoder.Serialize(int32(1))
			return &CXArgument{Value: &val, Typ: MakeType("bool")}, nil
		} else {
			val := encoder.Serialize(int32(0))
			return &CXArgument{Value: &val, Typ: MakeType("bool")}, nil
		}
	} else {
		return nil, err
	}
}

func gtF64 (arg1 *CXArgument, arg2 *CXArgument) (*CXArgument, error) {
	if err := checkTwoTypes("gtF64", "f64", "f64", arg1, arg2); err == nil {
		var num1 float64
		var num2 float64
		encoder.DeserializeRaw(*arg1.Value, &num1)
		encoder.DeserializeRaw(*arg2.Value, &num2)
		
		if num1 > num2 {
			val := encoder.Serialize(int32(1))
			return &CXArgument{Value: &val, Typ: MakeType("bool")}, nil
		} else {
			val := encoder.Serialize(int32(0))
			return &CXArgument{Value: &val, Typ: MakeType("bool")}, nil
		}
	} else {
		return nil, err
	}
}

func eqF64 (arg1 *CXArgument, arg2 *CXArgument) (*CXArgument, error) {
	if err := checkTwoTypes("eqF64", "f64", "f64", arg1, arg2); err == nil {
		var num1 float64
		var num2 float64
		encoder.DeserializeRaw(*arg1.Value, &num1)
		encoder.DeserializeRaw(*arg2.Value, &num2)
		
		if num1 == num2 {
			val := encoder.Serialize(int32(1))
			return &CXArgument{Value: &val, Typ: MakeType("bool")}, nil
		} else {
			val := encoder.Serialize(int32(0))
			return &CXArgument{Value: &val, Typ: MakeType("bool")}, nil
		}
	} else {
		return nil, err
	}
}

func lteqF64 (arg1 *CXArgument, arg2 *CXArgument) (*CXArgument, error) {
	if err := checkTwoTypes("lteqF64", "f64", "f64", arg1, arg2); err == nil {
		var num1 float64
		var num2 float64
		encoder.DeserializeRaw(*arg1.Value, &num1)
		encoder.DeserializeRaw(*arg2.Value, &num2)
		
		if num1 <= num2 {
			val := encoder.Serialize(int32(1))
			return &CXArgument{Value: &val, Typ: MakeType("bool")}, nil
		} else {
			val := encoder.Serialize(int32(0))
			return &CXArgument{Value: &val, Typ: MakeType("bool")}, nil
		}
	} else {
		return nil, err
	}
}

func gteqF64 (arg1 *CXArgument, arg2 *CXArgument) (*CXArgument, error) {
	if err := checkTwoTypes("gteqF64", "f64", "f64", arg1, arg2); err == nil {
		var num1 float64
		var num2 float64
		encoder.DeserializeRaw(*arg1.Value, &num1)
		encoder.DeserializeRaw(*arg2.Value, &num2)
		
		if num1 >= num2 {
			val := encoder.Serialize(int32(1))
			return &CXArgument{Value: &val, Typ: MakeType("bool")}, nil
		} else {
			val := encoder.Serialize(int32(0))
			return &CXArgument{Value: &val, Typ: MakeType("bool")}, nil
		}
	} else {
		return nil, err
	}
}

func ltStr (arg1 *CXArgument, arg2 *CXArgument) (*CXArgument, error) {
	if err := checkTwoTypes("ltStr", "str", "str", arg1, arg2); err == nil {
		str1 := string(*arg1.Value)
		str2 := string(*arg2.Value)
		
		if str1 < str2 {
			val := encoder.Serialize(int32(1))
			return &CXArgument{Value: &val, Typ: MakeType("bool")}, nil
		} else {
			val := encoder.Serialize(int32(0))
			return &CXArgument{Value: &val, Typ: MakeType("bool")}, nil
		}
	} else {
		return nil, err
	}
}

func gtStr (arg1 *CXArgument, arg2 *CXArgument) (*CXArgument, error) {
	if err := checkTwoTypes("gtStr", "str", "str", arg1, arg2); err == nil {
		str1 := string(*arg1.Value)
		str2 := string(*arg2.Value)
		
		if str1 > str2 {
			val := encoder.Serialize(int32(1))
			return &CXArgument{Value: &val, Typ: MakeType("bool")}, nil
		} else {
			val := encoder.Serialize(int32(0))
			return &CXArgument{Value: &val, Typ: MakeType("bool")}, nil
		}
	} else {
		return nil, err
	}
}

func eqStr (arg1 *CXArgument, arg2 *CXArgument) (*CXArgument, error) {
	if err := checkTwoTypes("eqStr", "str", "str", arg1, arg2); err == nil {
		str1 := string(*arg1.Value)
		str2 := string(*arg2.Value)
		
		if str1 == str2 {
			val := encoder.Serialize(int32(1))
			return &CXArgument{Value: &val, Typ: MakeType("bool")}, nil
		} else {
			val := encoder.Serialize(int32(0))
			return &CXArgument{Value: &val, Typ: MakeType("bool")}, nil
		}
	} else {
		return nil, err
	}
}

func lteqStr (arg1 *CXArgument, arg2 *CXArgument) (*CXArgument, error) {
	if err := checkTwoTypes("lteqStr", "str", "str", arg1, arg2); err == nil {
		str1 := string(*arg1.Value)
		str2 := string(*arg2.Value)
		
		if str1 <= str2 {
			val := encoder.Serialize(int32(1))
			return &CXArgument{Value: &val, Typ: MakeType("bool")}, nil
		} else {
			val := encoder.Serialize(int32(0))
			return &CXArgument{Value: &val, Typ: MakeType("bool")}, nil
		}
	} else {
		return nil, err
	}
}

func gteqStr (arg1 *CXArgument, arg2 *CXArgument) (*CXArgument, error) {
	if err := checkTwoTypes("gteqStr", "str", "str", arg1, arg2); err == nil {
		str1 := string(*arg1.Value)
		str2 := string(*arg2.Value)
		
		if str1 >= str2 {
			val := encoder.Serialize(int32(1))
			return &CXArgument{Value: &val, Typ: MakeType("bool")}, nil
		} else {
			val := encoder.Serialize(int32(0))
			return &CXArgument{Value: &val, Typ: MakeType("bool")}, nil
		}
	} else {
		return nil, err
	}
}

func ltByte (arg1 *CXArgument, arg2 *CXArgument) (*CXArgument, error) {
	if err := checkTwoTypes("ltByte", "byte", "byte", arg1, arg2); err == nil {
		byte1 := (*arg1.Value)[0]
		byte2 := (*arg2.Value)[0]
		
		if byte1 < byte2 {
			val := encoder.Serialize(int32(1))
			return &CXArgument{Value: &val, Typ: MakeType("bool")}, nil
		} else {
			val := encoder.Serialize(int32(0))
			return &CXArgument{Value: &val, Typ: MakeType("bool")}, nil
		}
	} else {
		return nil, err
	}
}

func gtByte (arg1 *CXArgument, arg2 *CXArgument) (*CXArgument, error) {
	if err := checkTwoTypes("gtByte", "byte", "byte", arg1, arg2); err == nil {
		byte1 := (*arg1.Value)[0]
		byte2 := (*arg2.Value)[0]
		
		if byte1 > byte2 {
			val := encoder.Serialize(int32(1))
			return &CXArgument{Value: &val, Typ: MakeType("bool")}, nil
		} else {
			val := encoder.Serialize(int32(0))
			return &CXArgument{Value: &val, Typ: MakeType("bool")}, nil
		}
	} else {
		return nil, err
	}
}

func eqByte (arg1 *CXArgument, arg2 *CXArgument) (*CXArgument, error) {
	if err := checkTwoTypes("eqByte", "byte", "byte", arg1, arg2); err == nil {
		byte1 := (*arg1.Value)[0]
		byte2 := (*arg2.Value)[0]
		
		if byte1 == byte2 {
			val := encoder.Serialize(int32(1))
			return &CXArgument{Value: &val, Typ: MakeType("bool")}, nil
		} else {
			val := encoder.Serialize(int32(0))
			return &CXArgument{Value: &val, Typ: MakeType("bool")}, nil
		}
	} else {
		return nil, err
	}
}

func lteqByte (arg1 *CXArgument, arg2 *CXArgument) (*CXArgument, error) {
	if err := checkTwoTypes("lteqByte", "byte", "byte", arg1, arg2); err == nil {
		byte1 := (*arg1.Value)[0]
		byte2 := (*arg2.Value)[0]
		
		if byte1 <= byte2 {
			val := encoder.Serialize(int32(1))
			return &CXArgument{Value: &val, Typ: MakeType("bool")}, nil
		} else {
			val := encoder.Serialize(int32(0))
			return &CXArgument{Value: &val, Typ: MakeType("bool")}, nil
		}
	} else {
		return nil, err
	}
}

func gteqByte (arg1 *CXArgument, arg2 *CXArgument) (*CXArgument, error) {
	if err := checkTwoTypes("gteqByte", "byte", "byte", arg1, arg2); err == nil {
		byte1 := (*arg1.Value)[0]
		byte2 := (*arg2.Value)[0]
		
		if byte1 >= byte2 {
			val := encoder.Serialize(int32(1))
			return &CXArgument{Value: &val, Typ: MakeType("bool")}, nil
		} else {
			val := encoder.Serialize(int32(0))
			return &CXArgument{Value: &val, Typ: MakeType("bool")}, nil
		}
	} else {
		return nil, err
	}
}



/*
  Cast functions
*/

func castToStr (arg *CXArgument) (*CXArgument, error) {
	strTyp := MakeType("str")
	switch arg.Typ.Name {
	case "[]byte":
		newArg := MakeArgument(arg.Value, strTyp)
		return newArg, nil
	default:
		return nil, errors.New(fmt.Sprintf("castToStr: type '%s' can't be casted to type 'str'", arg.Typ.Name))
	}
}

func castToByteA (arg *CXArgument) (*CXArgument, error) {
	byteATyp := MakeType("[]byte")
	switch arg.Typ.Name {
	case "str":
		return MakeArgument(arg.Value, byteATyp), nil
	case "[]byte":
		return arg, nil
	case "[]i32":
		var val []int32
		encoder.DeserializeRaw(*arg.Value, &val)

		output := make([]byte, len(val))
		for i, n := range val {
			output[i] = byte(n)
		}
		
		newVal := output
		return MakeArgument(&newVal, byteATyp), nil
	case "[]i64":
		var val []int64
		encoder.DeserializeRaw(*arg.Value, &val)

		output := make([]byte, len(val))
		for i, n := range val {
			output[i] = byte(n)
		}
		
		newVal := output
		return MakeArgument(&newVal, byteATyp), nil
	case "[]f32":
		var val []float32
		encoder.DeserializeRaw(*arg.Value, &val)

		output := make([]byte, len(val))
		for i, n := range val {
			output[i] = byte(n)
		}
		
		newVal := output
		return MakeArgument(&newVal, byteATyp), nil
	case "[]f64":
		var val []float64
		encoder.DeserializeRaw(*arg.Value, &val)

		output := make([]byte, len(val))
		for i, n := range val {
			output[i] = byte(n)
		}
		
		newVal := output
		return MakeArgument(&newVal, byteATyp), nil
	default:
		return nil, errors.New(fmt.Sprintf("castToByteA: type '%s' can't be casted to type '[]byte'", arg.Typ.Name))
	}
}

func castToByte (arg *CXArgument) (*CXArgument, error) {
	byteTyp := MakeType("byte")
	switch arg.Typ.Name {
	case "byte":
		return arg, nil
	case "i32":
		var val int32
		encoder.DeserializeRaw(*arg.Value, &val)
		newVal := []byte{byte(val)}
		return MakeArgument(&newVal, byteTyp), nil
	case "i64":
		var val int64
		encoder.DeserializeRaw(*arg.Value, &val)
		newVal := []byte{byte(val)}
		return MakeArgument(&newVal, byteTyp), nil
	case "f32":
		var val float32
		encoder.DeserializeRaw(*arg.Value, &val)
		newVal := []byte{byte(val)}
		return MakeArgument(&newVal, byteTyp), nil
	case "f64":
		var val float64
		encoder.DeserializeRaw(*arg.Value, &val)
		newVal := []byte{byte(val)}
		return MakeArgument(&newVal, byteTyp), nil
	default:
		return nil, errors.New(fmt.Sprintf("castToByte: type '%s' can't be casted to type 'byte'", arg.Typ.Name))
	}
}

func castToI32 (arg *CXArgument) (*CXArgument, error) {
	i32Typ := MakeType("i32")
	switch arg.Typ.Name {
	case "byte":
		val := (*arg.Value)[0]
		newVal := encoder.Serialize(int32(val))
		return MakeArgument(&newVal, i32Typ), nil
	case "i32":
		return arg, nil
	case "i64":
		var val int64
		encoder.DeserializeRaw(*arg.Value, &val)
		newVal := encoder.Serialize(int32(val))
		return MakeArgument(&newVal, i32Typ), nil
	case "f32":
		var val float32
		encoder.DeserializeRaw(*arg.Value, &val)
		newVal := encoder.Serialize(int32(val))
		return MakeArgument(&newVal, i32Typ), nil
	case "f64":
		var val float64
		encoder.DeserializeRaw(*arg.Value, &val)
		newVal := encoder.Serialize(int32(val))
		return MakeArgument(&newVal, i32Typ), nil
	default:
		return nil, errors.New(fmt.Sprintf("castToI32: type '%s' can't be casted to type 'i32'", arg.Typ.Name))
	}
}

func castToI64 (arg *CXArgument) (*CXArgument, error) {
	i64Typ := MakeType("i64")
	switch arg.Typ.Name {
	case "byte":
		val := (*arg.Value)[0]
		newVal := encoder.Serialize(int64(val))
		return MakeArgument(&newVal, i64Typ), nil
	case "i32":
		var val int32
		encoder.DeserializeRaw(*arg.Value, &val)
		newVal := encoder.Serialize(int64(val))
		return MakeArgument(&newVal, i64Typ), nil
	case "i64":
		return arg, nil
	case "f32":
		var val float32
		encoder.DeserializeRaw(*arg.Value, &val)
		newVal := encoder.Serialize(int64(val))
		return MakeArgument(&newVal, i64Typ), nil
	case "f64":
		var val float64
		encoder.DeserializeRaw(*arg.Value, &val)
		newVal := encoder.Serialize(int64(val))
		return MakeArgument(&newVal, i64Typ), nil
	default:
		return nil, errors.New(fmt.Sprintf("castToI64: type '%s' can't be casted to type 'i64'", arg.Typ.Name))
	}
}

func castToF32 (arg *CXArgument) (*CXArgument, error) {
	f32Typ := MakeType("f32")
	switch arg.Typ.Name {
	case "byte":
		val := (*arg.Value)[0]
		newVal := encoder.Serialize(float32(val))
		return MakeArgument(&newVal, f32Typ), nil
	case "i32":
		var val int32
		encoder.DeserializeRaw(*arg.Value, &val)
		newVal := encoder.Serialize(float32(val))
		return MakeArgument(&newVal, f32Typ), nil
	case "i64":
		var val int64
		encoder.DeserializeRaw(*arg.Value, &val)
		newVal := encoder.Serialize(float32(val))
		return MakeArgument(&newVal, f32Typ), nil
	case "f32":
		return arg, nil
	case "f64":
		var val float64
		encoder.DeserializeRaw(*arg.Value, &val)
		newVal := encoder.Serialize(float32(val))
		return MakeArgument(&newVal, f32Typ), nil
	default:
		return nil, errors.New(fmt.Sprintf("castToF32: type '%s' can't be casted to type 'f32'", arg.Typ.Name))
	}
}

func castToF64 (arg *CXArgument) (*CXArgument, error) {
	f64Typ := MakeType("f64")
	switch arg.Typ.Name {
	case "byte":
		val := (*arg.Value)[0]
		newVal := encoder.Serialize(float64(val))
		return MakeArgument(&newVal, f64Typ), nil
	case "i32":
		var val int32
		encoder.DeserializeRaw(*arg.Value, &val)
		newVal := encoder.Serialize(float64(val))
		return MakeArgument(&newVal, f64Typ), nil
	case "i64":
		var val int64
		encoder.DeserializeRaw(*arg.Value, &val)
		newVal := encoder.Serialize(float64(val))
		return MakeArgument(&newVal, f64Typ), nil
	case "f32":
		var val float32
		encoder.DeserializeRaw(*arg.Value, &val)
		newVal := encoder.Serialize(float64(val))
		return MakeArgument(&newVal, f64Typ), nil
	case "f64":
		return arg, nil
	default:
		return nil, errors.New(fmt.Sprintf("castToF64: type '%s' can't be casted to type 'f64'", arg.Typ.Name))
	}
}

func castToI32A (arg *CXArgument) (*CXArgument, error) {
	i32ATyp := MakeType("[]i32")
	switch arg.Typ.Name {
	case "[]byte":
		val := *arg.Value

		output := make([]int32, len(val))
		for i, n := range val {
			output[i] = int32(n)
		}
		
		newVal := encoder.Serialize(output)
		newArg := MakeArgument(&newVal, i32ATyp)
		return newArg, nil
	case "[]i32":
		return arg, nil
	case "[]i64":
		var val []int64
		encoder.DeserializeRaw(*arg.Value, &val)

		output := make([]int32, len(val))
		for i, n := range val {
			output[i] = int32(n)
		}
		
		newVal := encoder.Serialize(output)
		newArg := MakeArgument(&newVal, i32ATyp)
		return newArg, nil
	case "[]f32":
		var val []float32
		encoder.DeserializeRaw(*arg.Value, &val)

		output := make([]int32, len(val))
		for i, n := range val {
			output[i] = int32(n)
		}
		
		newVal := encoder.Serialize(output)
		newArg := MakeArgument(&newVal, i32ATyp)
		return newArg, nil
	case "[]f64":
		var val []float64
		encoder.DeserializeRaw(*arg.Value, &val)

		output := make([]int32, len(val))
		for i, n := range val {
			output[i] = int32(n)
		}
		
		newVal := encoder.Serialize(output)
		newArg := MakeArgument(&newVal, i32ATyp)
		return newArg, nil
	default:
		return nil, errors.New(fmt.Sprintf("castToI32A: type '%s' can't be casted to type '[]i32'", arg.Typ.Name))
	}
}

func castToI64A (arg *CXArgument) (*CXArgument, error) {
	i64ATyp := MakeType("[]i64")
	switch arg.Typ.Name {
	case "[]byte":
		val := *arg.Value

		output := make([]int64, len(val))
		for i, n := range val {
			output[i] = int64(n)
		}
		
		newVal := encoder.Serialize(output)
		return MakeArgument(&newVal, i64ATyp), nil
	case "[]i32":
		var val []int32
		encoder.DeserializeRaw(*arg.Value, &val)

		output := make([]int64, len(val))
		for i, n := range val {
			output[i] = int64(n)
		}
		
		newVal := encoder.Serialize(output)
		return MakeArgument(&newVal, i64ATyp), nil
	case "[]i64":
		return arg, nil
	case "[]f32":
		var val []float32
		encoder.DeserializeRaw(*arg.Value, &val)

		output := make([]int64, len(val))
		for i, n := range val {
			output[i] = int64(n)
		}
		
		newVal := encoder.Serialize(output)
		return MakeArgument(&newVal, i64ATyp), nil
	case "[]f64":
		var val []float64
		encoder.DeserializeRaw(*arg.Value, &val)

		output := make([]int64, len(val))
		for i, n := range val {
			output[i] = int64(n)
		}
		
		newVal := encoder.Serialize(output)
		return MakeArgument(&newVal, i64ATyp), nil
	default:
		return nil, errors.New(fmt.Sprintf("castToI64A: type '%s' can't be casted to type '[]i64'", arg.Typ.Name))
	}
}

func castToF32A (arg *CXArgument) (*CXArgument, error) {
	f32ATyp := MakeType("[]f32")
	switch arg.Typ.Name {
	case "[]byte":
		val := *arg.Value

		output := make([]float32, len(val))
		for i, n := range val {
			output[i] = float32(n)
		}
		
		newVal := encoder.Serialize(output)
		return MakeArgument(&newVal, f32ATyp), nil
	case "[]i32":
		var val []int32
		encoder.DeserializeRaw(*arg.Value, &val)

		output := make([]float32, len(val))
		for i, n := range val {
			output[i] = float32(n)
		}
		
		newVal := encoder.Serialize(output)
		return MakeArgument(&newVal, f32ATyp), nil
	case "[]i64":
		var val []int64
		encoder.DeserializeRaw(*arg.Value, &val)

		output := make([]float32, len(val))
		for i, n := range val {
			output[i] = float32(n)
		}
		
		newVal := encoder.Serialize(output)
		return MakeArgument(&newVal, f32ATyp), nil
	case "[]f32":
		return arg, nil
	case "[]f64":
		var val []float64
		encoder.DeserializeRaw(*arg.Value, &val)

		output := make([]float32, len(val))
		for i, n := range val {
			output[i] = float32(n)
		}
		
		newVal := encoder.Serialize(output)
		return MakeArgument(&newVal, f32ATyp), nil
	default:
		return nil, errors.New(fmt.Sprintf("castToF32A: type '%s' can't be casted to type '[]f32'", arg.Typ.Name))
	}
}

func castToF64A (arg *CXArgument) (*CXArgument, error) {
	f64ATyp := MakeType("[]f64")
	switch arg.Typ.Name {
	case "[]byte":
		val := *arg.Value

		output := make([]float64, len(val))
		for i, n := range val {
			output[i] = float64(n)
		}
		
		newVal := encoder.Serialize(output)
		return MakeArgument(&newVal, f64ATyp), nil
	case "[]i32":
		var val []int32
		encoder.DeserializeRaw(*arg.Value, &val)

		output := make([]float64, len(val))
		for i, n := range val {
			output[i] = float64(n)
		}
		
		newVal := encoder.Serialize(output)
		return MakeArgument(&newVal, f64ATyp), nil
	case "[]i64":
		var val []int64
		encoder.DeserializeRaw(*arg.Value, &val)

		output := make([]float64, len(val))
		for i, n := range val {
			output[i] = float64(n)
		}
		
		newVal := encoder.Serialize(output)
		return MakeArgument(&newVal, f64ATyp), nil
	case "[]f32":
		var val []float32
		encoder.DeserializeRaw(*arg.Value, &val)

		output := make([]float64, len(val))
		for i, n := range val {
			output[i] = float64(n)
		}
		
		newVal := encoder.Serialize(output)
		return MakeArgument(&newVal, f64ATyp), nil
	case "[]f64":
		return arg, nil
	default:
		return nil, errors.New(fmt.Sprintf("castToF64A: type '%s' can't be casted to type '[]f64'", arg.Typ.Name))
	}
}

// goTo increments/decrements the call.Line to the desired expression line.
// Used for if/else and loop statements.
func goTo (call *CXCall, predicate *CXArgument, thenLine *CXArgument, elseLine *CXArgument) (*CXArgument, error) {
	if err := checkThreeTypes("goTo", "bool", "i32", "i32", predicate, thenLine, elseLine); err == nil {
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
			call.Line = call.Line + int(elseLineNo) - 1
		} else {
			call.Line = call.Line + int(thenLineNo) - 1
		}
		
		if isFalse {
			val := encoder.Serialize(int32(0))
			return MakeArgument(&val, MakeType("bool")), nil
		} else {
			val := encoder.Serialize(int32(1))
			return MakeArgument(&val, MakeType("bool")), nil
		}
	} else {
		return nil, err
	}
}

/*
  Time functions
*/

func sleep (ms *CXArgument) (*CXArgument, error) {
	if err := checkType("sleep", "i32", ms); err == nil {
		var duration int32
		encoder.DeserializeRaw(*ms.Value, &duration)

		time.Sleep(time.Duration(duration) * time.Millisecond)

		return ms, nil
	} else {
		return nil, err
	}
}


/*
  Prolog functions
*/

func setClauses (clss *CXArgument, mod *CXModule) (*CXArgument, error) {
	if err := checkType("setClauses", "str", clss); err == nil {
		clauses := string(*clss.Value)
		mod.AddClauses(clauses)

		return clss, nil
	} else {
		return nil, err
	}
}

func addObject (obj *CXArgument, mod *CXModule) (*CXArgument, error) {
	if err := checkType("addObject", "str", obj); err == nil {
		mod.AddObject(MakeObject(string(*obj.Value)))

		return obj, nil
	} else {
		return nil, err
	}
}

func setQuery (qry *CXArgument, mod *CXModule) (*CXArgument, error) {
	if err := checkType("setQuery", "str", qry); err == nil {
		query := string(*qry.Value)
		mod.AddQuery(query)

		return qry, nil
	} else {
		return nil, err
	}
}

func remObject (obj *CXArgument, mod *CXModule) (*CXArgument, error) {
	if err := checkType("remObject", "str", obj); err == nil {
		object := string(*obj.Value)
		mod.RemoveObject(object)

		return obj, nil
	} else {
		return nil, err
	}
}

func remObjects (mod *CXModule) (*CXArgument, error) {
	mod.RemoveObjects()

	success := encoder.Serialize(int32(1))
	return MakeArgument(&success, MakeType("bool")), nil
}

/*
  Meta-programming functions
*/

func remArg (tag *CXArgument, caller *CXFunction) (*CXArgument, error) {
	if err := checkType("remArg", "str", tag); err == nil {
		for _, expr := range caller.Expressions {
			if expr.Tag == string(*tag.Value) {
				expr.RemoveArgument()
				val := encoder.Serialize(int32(0))
				return MakeArgument(&val, MakeType("bool")), nil
			}
		}
	} else {
		return nil, err
	}
	return nil, errors.New(fmt.Sprintf("remArg: no expression with tag '%s' was found", string(*tag.Value)))
}

// func addArg (tag *CXArgument, ident *CXArgument, caller *CXFunction) (*CXArgument, error) {
// 	if err := checkTwoTypes("addArg", "str", "str", tag, ident); err == nil {
// 		for _, expr := range caller.Expressions {
// 			if expr.Tag == string(*tag.Value) {
// 				expr.AddArgument(MakeArgument(ident.Value, MakeType("ident")))
// 				val := encoder.Serialize(int32(0))
// 				return MakeArgument(&val, MakeType("bool")), nil
// 			}
// 		}
// 	} else {
// 		return nil, err
// 	}
// 	return nil, errors.New(fmt.Sprintf("remArg: no expression with tag '%s' was found", string(*tag.Value)))
// }

func addExpr (tag *CXArgument, fnName *CXArgument, caller *CXFunction, line int) (*CXArgument, error) {
	if err := checkType("addExpr", "str", fnName); err == nil {
		mod := caller.Module
		
		opName := string(*fnName.Value)
		if fn, err := mod.Context.GetFunction(opName, mod.Name); err == nil {
			expr := MakeExpression(fn)
			expr.AddTag(string(*tag.Value))
			
			caller.AddExpression(expr)

			//caller.Expressions = append(caller.Expressions[:line], expr, caller.Expressions[line:(len(caller.Expressions)-2)]...)

			// re-indexing expression line numbers
			// for i, expr := range caller.Expressions {
			// 	expr.Line = i
			// }
			
			val := encoder.Serialize(int32(0))
			return MakeArgument(&val, MakeType("bool")), nil
		} else {
			val := encoder.Serialize(int32(1))
			return MakeArgument(&val, MakeType("bool")), nil
		}
	} else {
		return nil, err
	}
}

func remExpr (tag *CXArgument, caller *CXFunction) (*CXArgument, error) {
	if err := checkType("remExpr", "str", tag); err == nil {
		for i, expr := range caller.Expressions {
			if expr.Tag == string(*tag.Value) {
				caller.RemoveExpression(i)
				val := encoder.Serialize(int32(0))
				return MakeArgument(&val, MakeType("bool")), nil
			}
		}
	} else {
		return nil, err
	}
	return nil, errors.New(fmt.Sprintf("remExpr: no expression with tag '%s' was found", string(*tag.Value)))
}

//func affFn (filter *CXArgument, )

func affExpr (tag *CXArgument, filter *CXArgument, idx *CXArgument, caller *CXFunction) (*CXArgument, error) {
	if err := checkThreeTypes("affExpr", "str", "str", "i32", tag, filter, idx); err == nil {
		var index int32
		encoder.DeserializeRaw(*idx.Value, &index)

		if index == -1 {
			for _, expr := range caller.Expressions {
				if expr.Tag == string(*tag.Value) {
					affs := FilterAffordances(expr.GetAffordances(), string(*filter.Value))
					PrintAffordances(affs)
					val := encoder.Serialize(int32(len(affs)))
					return MakeArgument(&val, MakeType("i32")), nil
				}
			}
		} else if index < -1 {
			for _, expr := range caller.Expressions {
				if expr.Tag == string(*tag.Value) {
					affs := FilterAffordances(expr.GetAffordances(), string(*filter.Value))
					val := encoder.Serialize(int32(len(affs)))
					return MakeArgument(&val, MakeType("i32")), nil
				}
			}
		} else {
			for _, expr := range caller.Expressions {
				if expr.Tag == string(*tag.Value) {
					affs := FilterAffordances(expr.GetAffordances(), string(*filter.Value))
					affs[index].ApplyAffordance()
					val := encoder.Serialize(int32(len(affs)))
					return MakeArgument(&val, MakeType("i32")), nil
				}
			}
		}
	} else {
		return nil, err
	}
	return nil, errors.New(fmt.Sprintf("affArg: no expression with tag '%s' was found", string(*tag.Value)))
}

func initDef (arg1 *CXArgument) (*CXArgument, error) {
	if err := checkType("initDef", "str", arg1); err == nil {
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

		return MakeArgument(&zeroVal, MakeType(typName)), nil
	} else {
		return nil, err
	}
}

func randI32 (min *CXArgument, max *CXArgument) (*CXArgument, error) {
	if err := checkTwoTypes("randI32", "i32", "i32", min, max); err == nil {
		var minimum int32
		encoder.DeserializeRaw(*min.Value, &minimum)

		var maximum int32
		encoder.DeserializeRaw(*max.Value, &maximum)

		if minimum > maximum {
			return nil, errors.New(fmt.Sprintf("randI32: min must be less than max (%d !< %d)", minimum, maximum))
		}

		rand.Seed(time.Now().UTC().UnixNano())
		output := encoder.SerializeAtomic(int32(rand.Intn(int(maximum - minimum)) + int(minimum)))
		return MakeArgument(&output, MakeType("i32")), nil
	} else {
		return nil, err
	}
}

func randI64 (min *CXArgument, max *CXArgument) (*CXArgument, error) {
	if err := checkTwoTypes("randI64", "i64", "i64", min, max); err == nil {
		var minimum int64
		encoder.DeserializeRaw(*min.Value, &minimum)

		var maximum int64
		encoder.DeserializeRaw(*max.Value, &maximum)

		if minimum > maximum {
			return nil, errors.New(fmt.Sprintf("randI64: min must be less than max (%d !< %d)", minimum, maximum))
		}

		rand.Seed(time.Now().UTC().UnixNano())
		output := encoder.SerializeAtomic(int32(rand.Intn(int(maximum - minimum)) + int(minimum)))
		return MakeArgument(&output, MakeType("i64")), nil
	} else {
		return nil, err
	}
}

/*
  Make Array
*/

func makeArray (typ string, size *CXArgument) (*CXArgument, error) {
	if err := checkType("makeArray", "i32", size); err == nil {
		var len int32
		encoder.DeserializeRaw(*size.Value, &len)

		switch typ {
		case "[]bool":
			arr := make([]bool, len)
			val := encoder.Serialize(arr)
			return MakeArgument(&val, MakeType(typ)), nil
		case "[]byte":
			arr := make([]byte, len)
			return MakeArgument(&arr, MakeType(typ)), nil
		case "[]i32":
			arr := make([]int32, len)
			val := encoder.Serialize(arr)
			return MakeArgument(&val, MakeType(typ)), nil
		case "[]i64":
			arr := make([]int64, len)
			val := encoder.Serialize(arr)
			return MakeArgument(&val, MakeType(typ)), nil
		case "[]f32":
			arr := make([]float32, len)
			val := encoder.Serialize(arr)
			return MakeArgument(&val, MakeType(typ)), nil
		case "[]f64":
			arr := make([]float64, len)
			val := encoder.Serialize(arr)
			return MakeArgument(&val, MakeType(typ)), nil
		case "default":
			return nil, errors.New(fmt.Sprintf("makeArray: argument 1 is type '%s'; expected type 'i32'", size.Typ.Name))
		}
		return nil, nil
	} else {
		return nil, err
	}
}
