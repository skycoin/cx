package interpreted

import (
	"fmt"
	"errors"
	"math"
	"github.com/skycoin/skycoin/src/cipher/encoder"
)

func addF64 (arg1 *CXArgument, arg2 *CXArgument, expr *CXExpression, call *CXCall) error {
	if err := checkTwoTypes("f64.add", "f64", "f64", arg1, arg2); err == nil {
		var num1 float64
		var num2 float64
		encoder.DeserializeRaw(*arg1.Value, &num1)
		encoder.DeserializeRaw(*arg2.Value, &num2)

		output := encoder.Serialize(float64(num1 + num2))

		assignOutput(0, output, "f64", expr, call)
		return nil
	} else {
		return err
	}
}

func subF64 (arg1 *CXArgument, arg2 *CXArgument, expr *CXExpression, call *CXCall) error {
	if err := checkTwoTypes("f64.sub", "f64", "f64", arg1, arg2); err == nil {
		var num1 float64
		var num2 float64
		encoder.DeserializeRaw(*arg1.Value, &num1)
		encoder.DeserializeRaw(*arg2.Value, &num2)

		output := encoder.Serialize(float64(num1 - num2))

		assignOutput(0, output, "f64", expr, call)
		return nil
	} else {
		return err
	}
}

func mulF64 (arg1 *CXArgument, arg2 *CXArgument, expr *CXExpression, call *CXCall) error {
	if err := checkTwoTypes("f64.mul", "f64", "f64", arg1, arg2); err == nil {
		var num1 float64
		var num2 float64
		encoder.DeserializeRaw(*arg1.Value, &num1)
		encoder.DeserializeRaw(*arg2.Value, &num2)

		output := encoder.Serialize(float64(num1 * num2))

		assignOutput(0, output, "f64", expr, call)
		return nil
	} else {
		return err
	}
}

func divF64 (arg1 *CXArgument, arg2 *CXArgument, expr *CXExpression, call *CXCall) error {
	if err := checkTwoTypes("f64.div", "f64", "f64", arg1, arg2); err == nil {
		var num1 float64
		var num2 float64
		encoder.DeserializeRaw(*arg1.Value, &num1)
		encoder.DeserializeRaw(*arg2.Value, &num2)

		if num2 == float64(0.0) {
			return errors.New("f64.div: Division by 0")
		}

		output := encoder.Serialize(float64(num1 / num2))

		assignOutput(0, output, "f64", expr, call)
		return nil
	} else {
		return err
	}
}

func powF64 (arg1 *CXArgument, arg2 *CXArgument, expr *CXExpression, call *CXCall) error {
	if err := checkTwoTypes("f64.pow", "f64", "f64", arg1, arg2); err == nil {
		var num1 float64
		var num2 float64
		encoder.DeserializeRaw(*arg1.Value, &num1)
		encoder.DeserializeRaw(*arg2.Value, &num2)

		output := encoder.Serialize(math.Pow(num1, num2))

		assignOutput(0, output, "f64", expr, call)
		return nil
	} else {
		return err
	}
}

func absF64 (arg1 *CXArgument, expr *CXExpression, call *CXCall) error {
	if err := checkType("f64.abs", "f64", arg1); err == nil {
		var num1 float64
		encoder.DeserializeRaw(*arg1.Value, &num1)

		output := encoder.Serialize(math.Abs(num1))

		assignOutput(0, output, "f64", expr, call)
		return nil
	} else {
		return err
	}	
}

func cosF64 (arg1 *CXArgument, expr *CXExpression, call *CXCall) error {
	if err := checkType("f64.cos", "f64", arg1); err == nil {
		var num1 float64
		encoder.DeserializeRaw(*arg1.Value, &num1)

		output := encoder.Serialize(math.Cos(num1))

		assignOutput(0, output, "f64", expr, call)
		return nil
	} else {
		return err
	}	
}

func sinF64 (arg1 *CXArgument, expr *CXExpression, call *CXCall) error {
	if err := checkType("f64.sin", "f64", arg1); err == nil {
		var num1 float64
		encoder.DeserializeRaw(*arg1.Value, &num1)

		output := encoder.Serialize(math.Sin(num1))

		assignOutput(0, output, "f64", expr, call)
		return nil
	} else {
		return err
	}	
}

func readF64A (arr *CXArgument, idx *CXArgument, expr *CXExpression, call *CXCall) error {
	if err := checkTwoTypes("[]f64.read", "[]f64", "i32", arr, idx); err == nil {
		var index int32
		encoder.DeserializeRaw(*idx.Value, &index)

		var size int32
		encoder.DeserializeAtomic((*arr.Value)[0:4], &size)

		if index < 0 {
			return errors.New(fmt.Sprintf("[]f64.read: negative index %d", index))
		}

		if index >= size {
			return errors.New(fmt.Sprintf("[]f64.read: index %d exceeds array of length %d", index, size))
		}

		var value float64
		encoder.DeserializeRaw((*arr.Value)[((index)*8)+4:((index+1)*8)+4], &value)
		output := encoder.Serialize(value)

		assignOutput(0, output, "f64", expr, call)
		return nil
	} else {
		return err
	}
}

func writeF64A (arr *CXArgument, idx *CXArgument, val *CXArgument, expr *CXExpression, call *CXCall) error {
	if err := checkThreeTypes("[]f64.write", "[]f64", "i32", "f64", arr, idx, val); err == nil {
		var index int32
		encoder.DeserializeRaw(*idx.Value, &index)

		var size int32
		encoder.DeserializeAtomic((*arr.Value)[0:4], &size)

		if index < 0 {
			return errors.New(fmt.Sprintf("[]f64.write: negative index %d", index))
		}

		if index >= size {
			return errors.New(fmt.Sprintf("[]f64.write: index %d exceeds array of length %d", index, size))
		}

		// i := (int(index)*8)+4
		// for c := 0; c < 8; c++ {
		// 	(*arr.Value)[i + c] = (*val.Value)[c]
		// }

		offset := int(index) * 8 + 4
		firstChunk := make([]byte, offset)
		secondChunk := make([]byte, len(*arr.Value) - (offset + 8))

		copy(firstChunk, (*arr.Value)[:offset])
		copy(secondChunk, (*arr.Value)[offset + 8:])

		final := append(firstChunk, *val.Value...)
		final = append(final, secondChunk...)

		assignOutput(0, final, "[]f64", expr, call)

		return nil
	} else {
		return err
	}
}

func lenF64A (arr *CXArgument, expr *CXExpression, call *CXCall) error {
	if err := checkType("[]f64.len", "[]f64", arr); err == nil {
		size := (*arr.Value)[:4]
		assignOutput(0, size, "i32", expr, call)
		return nil
	} else {
		return err
	}
}

func ltF64 (arg1 *CXArgument, arg2 *CXArgument, expr *CXExpression, call *CXCall) error {
	if err := checkTwoTypes("f64.lt", "f64", "f64", arg1, arg2); err == nil {
		var num1 float64
		var num2 float64
		encoder.DeserializeRaw(*arg1.Value, &num1)
		encoder.DeserializeRaw(*arg2.Value, &num2)

		var val []byte

		if num1 < num2 {
			val = encoder.Serialize(int32(1))
		} else {
			val = encoder.Serialize(int32(0))
		}

		assignOutput(0, val, "bool", expr, call)
		return nil
	} else {
		return err
	}
}

func gtF64 (arg1 *CXArgument, arg2 *CXArgument, expr *CXExpression, call *CXCall) error {
	if err := checkTwoTypes("f64.gt", "f64", "f64", arg1, arg2); err == nil {
		var num1 float64
		var num2 float64
		encoder.DeserializeRaw(*arg1.Value, &num1)
		encoder.DeserializeRaw(*arg2.Value, &num2)

		var val []byte

		if num1 > num2 {
			val = encoder.Serialize(int32(1))
		} else {
			val = encoder.Serialize(int32(0))
		}
		
		assignOutput(0, val, "bool", expr, call)
		return nil
	} else {
		return err
	}
}

func eqF64 (arg1 *CXArgument, arg2 *CXArgument, expr *CXExpression, call *CXCall) error {
	if err := checkTwoTypes("f64.eq", "f64", "f64", arg1, arg2); err == nil {
		var num1 float64
		var num2 float64
		encoder.DeserializeRaw(*arg1.Value, &num1)
		encoder.DeserializeRaw(*arg2.Value, &num2)

		var val []byte

		if num1 == num2 {
			val = encoder.Serialize(int32(1))
		} else {
			val = encoder.Serialize(int32(0))
		}

		assignOutput(0, val, "bool", expr, call)
		return nil
	} else {
		return err
	}
}

func uneqF64 (arg1 *CXArgument, arg2 *CXArgument, expr *CXExpression, call *CXCall) error {
	if err := checkTwoTypes("f64.uneq", "f64", "f64", arg1, arg2); err == nil {
		var num1 float64
		var num2 float64
		encoder.DeserializeRaw(*arg1.Value, &num1)
		encoder.DeserializeRaw(*arg2.Value, &num2)

		var val []byte

		if num1 != num2 {
			val = encoder.Serialize(int32(1))
		} else {
			val = encoder.Serialize(int32(0))
		}

		assignOutput(0, val, "bool", expr, call)
		return nil
	} else {
		return err
	}
}

func lteqF64 (arg1 *CXArgument, arg2 *CXArgument, expr *CXExpression, call *CXCall) error {
	if err := checkTwoTypes("f64.lteq", "f64", "f64", arg1, arg2); err == nil {
		var num1 float64
		var num2 float64
		encoder.DeserializeRaw(*arg1.Value, &num1)
		encoder.DeserializeRaw(*arg2.Value, &num2)

		var val []byte

		if num1 <= num2 {
			val = encoder.Serialize(int32(1))
		} else {
			val = encoder.Serialize(int32(0))
		}

		assignOutput(0, val, "bool", expr, call)
		return nil
	} else {
		return err
	}
}

func gteqF64 (arg1 *CXArgument, arg2 *CXArgument, expr *CXExpression, call *CXCall) error {
	if err := checkTwoTypes("f64.gteq", "f64", "f64", arg1, arg2); err == nil {
		var num1 float64
		var num2 float64
		encoder.DeserializeRaw(*arg1.Value, &num1)
		encoder.DeserializeRaw(*arg2.Value, &num2)

		var val []byte

		if num1 >= num2 {
			val = encoder.Serialize(int32(1))
		} else {
			val = encoder.Serialize(int32(0))
		}

		assignOutput(0, val, "bool", expr, call)
		return nil
	} else {
		return err
	}
}

func concatF64A (arg1 *CXArgument, arg2 *CXArgument, expr *CXExpression, call *CXCall) error {
	if err := checkTwoTypes("[]f64.concat", "[]f64", "[]f64", arg1, arg2); err == nil {
		var slice1 []float64
		var slice2 []float64
		encoder.DeserializeRaw(*arg1.Value, &slice1)
		encoder.DeserializeRaw(*arg2.Value, &slice2)

		output := append(slice1, slice2...)
		sOutput := encoder.Serialize(output)

		assignOutput(0, sOutput, "[]f64", expr, call)
		return nil
	} else {
		return err
	}
}

func appendF64A (arg1 *CXArgument, arg2 *CXArgument, expr *CXExpression, call *CXCall) error {
	if err := checkTwoTypes("[]f64.append", "[]f64", "f64", arg1, arg2); err == nil {
		var slice []float64
		var literal float64
		encoder.DeserializeRaw(*arg1.Value, &slice)
		encoder.DeserializeRaw(*arg2.Value, &literal)

		output := append(slice, literal)
		sOutput := encoder.Serialize(output)

		//*arg1.Value = sOutput
		assignOutput(0, sOutput, "[]f64", expr, call)
		return nil
	} else {
		return err
	}
}

func copyF64A (arg1 *CXArgument, arg2 *CXArgument, expr *CXExpression, call *CXCall) error {
	if err := checkTwoTypes("[]f64.copy", "[]f64", "[]f64", arg1, arg2); err == nil {
		copy(*arg1.Value, *arg2.Value)
		return nil
	} else {
		return err
	}
}
