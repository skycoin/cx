package base

// import (
// 	"fmt"
// 	"errors"
// 	"math"
// 	"github.com/skycoin/skycoin/src/cipher/encoder"
// )

// func addF32 (arg1 *CXArgument, arg2 *CXArgument, expr *CXExpression, call *CXCall) error {
// 	if err := checkTwoTypes("f32.add", "f32", "f32", arg1, arg2); err == nil {
		
// 	} else {
// 		return err
// 	}
	
// 	var num1 float32
// 	var num2 float32
// 	encoder.DeserializeRaw(*arg1.Value, &num1)
// 	encoder.DeserializeRaw(*arg2.Value, &num2)

// 	output := encoder.Serialize(float32(num1 + num2))

// 	assignOutput(0, output, "f32", expr, call)
// 	return nil
// }

// func subF32 (arg1 *CXArgument, arg2 *CXArgument, expr *CXExpression, call *CXCall) error {
// 	if err := checkTwoTypes("f32.sub", "f32", "f32", arg1, arg2); err == nil {
// 		var num1 float32
// 		var num2 float32
// 		encoder.DeserializeRaw(*arg1.Value, &num1)
// 		encoder.DeserializeRaw(*arg2.Value, &num2)

// 		output := encoder.Serialize(float32(num1 - num2))

// 		assignOutput(0, output, "f32", expr, call)
// 		return nil
// 	} else {
// 		return err
// 	}
// }

// func mulF32 (arg1 *CXArgument, arg2 *CXArgument, expr *CXExpression, call *CXCall) error {
// 	if err := checkTwoTypes("f32.mul", "f32", "f32", arg1, arg2); err == nil {
// 		var num1 float32
// 		var num2 float32
// 		encoder.DeserializeRaw(*arg1.Value, &num1)
// 		encoder.DeserializeRaw(*arg2.Value, &num2)

// 		output := encoder.Serialize(float32(num1 * num2))

// 		assignOutput(0, output, "f32", expr, call)
// 		return nil
// 	} else {
// 		return err
// 	}
// }

// func divF32 (arg1 *CXArgument, arg2 *CXArgument, expr *CXExpression, call *CXCall) error {
// 	if err := checkTwoTypes("f32.div", "f32", "f32", arg1, arg2); err == nil {
// 		var num1 float32
// 		var num2 float32
// 		encoder.DeserializeRaw(*arg1.Value, &num1)
// 		encoder.DeserializeRaw(*arg2.Value, &num2)

// 		if num2 == float32(0.0) {
// 			return errors.New("f32.div: Division by 0")
// 		}

// 		output := encoder.Serialize(float32(num1 / num2))

// 		assignOutput(0, output, "f32", expr, call)
// 		return nil
// 	} else {
// 		return err
// 	}	
// }

// func powF32 (arg1 *CXArgument, arg2 *CXArgument, expr *CXExpression, call *CXCall) error {
// 	if err := checkTwoTypes("f32.pow", "f32", "f32", arg1, arg2); err == nil {
// 		var num1 float32
// 		var num2 float32
// 		encoder.DeserializeRaw(*arg1.Value, &num1)
// 		encoder.DeserializeRaw(*arg2.Value, &num2)

// 		output := encoder.Serialize(float32(math.Pow(float64(num1), float64(num2))))

// 		assignOutput(0, output, "f32", expr, call)
// 		return nil
// 	} else {
// 		return err
// 	}
// }

// func absF32 (arg1 *CXArgument, expr *CXExpression, call *CXCall) error {
// 	if err := checkType("f32.abs", "f32", arg1); err == nil {
// 		var num1 float32
// 		encoder.DeserializeRaw(*arg1.Value, &num1)

// 		output := encoder.Serialize(float32(math.Abs(float64(num1))))

// 		assignOutput(0, output, "f32", expr, call)
// 		return nil
// 	} else {
// 		return err
// 	}	
// }

// func cosF32 (arg1 *CXArgument, expr *CXExpression, call *CXCall) error {
// 	if err := checkType("f32.cos", "f32", arg1); err == nil {
// 		var num1 float32
// 		encoder.DeserializeRaw(*arg1.Value, &num1)

// 		output := encoder.Serialize(float32(math.Cos(float64(num1))))

// 		assignOutput(0, output, "f32", expr, call)
// 		return nil
// 	} else {
// 		return err
// 	}	
// }

// func sinF32 (arg1 *CXArgument, expr *CXExpression, call *CXCall) error {
// 	if err := checkType("f32.sin", "f32", arg1); err == nil {
// 		var num1 float32
// 		encoder.DeserializeRaw(*arg1.Value, &num1)

// 		output := encoder.Serialize(float32(math.Sin(float64(num1))))

// 		assignOutput(0, output, "f32", expr, call)
// 		return nil
// 	} else {
// 		return err
// 	}	
// }

// func readF32A (arr *CXArgument, idx *CXArgument, expr *CXExpression, call *CXCall) error {
// 	if err := checkTwoTypes("[]f32.read", "[]f32", "i32", arr, idx); err == nil {
// 		var index int32
// 		encoder.DeserializeRaw(*idx.Value, &index)

// 		var size int32
// 		encoder.DeserializeAtomic((*arr.Value)[0:4], &size)
		
// 		if index < 0 {
// 			return errors.New(fmt.Sprintf("[]f32.read: negative index %d", index))
// 		}

// 		if index >= size {
// 			return errors.New(fmt.Sprintf("[]f32.read: index %d exceeds array of length %d", index, size))
// 		}

// 		var value float32
// 		encoder.DeserializeRaw((*arr.Value)[(index+1)*4:(index+2)*4], &value)
// 		output := encoder.Serialize(value)

// 		assignOutput(0, output, "f32", expr, call)
// 		return nil
// 	} else {
// 		return err
// 	}
// }

// func writeF32A (arr *CXArgument, idx *CXArgument, val *CXArgument, expr *CXExpression, call *CXCall) error {
// 	if err := checkThreeTypes("[]f32.write", "[]f32", "i32", "f32", arr, idx, val); err == nil {
// 		var index int32
// 		encoder.DeserializeRaw(*idx.Value, &index)

// 		var size int32
// 		encoder.DeserializeAtomic((*arr.Value)[0:4], &size)

// 		if index < 0 {
// 			return errors.New(fmt.Sprintf("[]f32.write: negative index %d", index))
// 		}

// 		if index >= size {
// 			return errors.New(fmt.Sprintf("[]f32.write: index %d exceeds array of length %d", index, size))
// 		}

// 		i := (int(index)+1)*4
// 		for c := 0; c < 4; c++ {
// 			(*arr.Value)[i + c] = (*val.Value)[c]
// 		}

		
// 		// offset := int(index) * 4 + 4
// 		// firstChunk := make([]byte, offset)
// 		// secondChunk := make([]byte, len(*arr.Value) - (offset + 4))

// 		// copy(firstChunk, (*arr.Value)[:offset])
// 		// copy(secondChunk, (*arr.Value)[offset + 4:])

// 		// final := append(firstChunk, *val.Value...)
// 		// final = append(final, secondChunk...)

// 		assignOutput(0, *arr.Value, "[]f32", expr, call)
// 		//assignOutput(0, final, "[]f32", expr, call)
		
		
// 		return nil
// 	} else {
// 		return err
// 	}
// }

// func lenF32A (arr *CXArgument, expr *CXExpression, call *CXCall) error {
// 	if err := checkType("[]f32.len", "[]f32", arr); err == nil {
// 		size := (*arr.Value)[:4]
// 		assignOutput(0, size, "i32", expr, call)
// 		return nil
// 	} else {
// 		return err
// 	}
// }

// func ltF32 (arg1 *CXArgument, arg2 *CXArgument, expr *CXExpression, call *CXCall) error {
// 	if err := checkTwoTypes("f32.lt", "f32", "f32", arg1, arg2); err == nil {
// 		var num1 float32
// 		var num2 float32
// 		encoder.DeserializeRaw(*arg1.Value, &num1)
// 		encoder.DeserializeRaw(*arg2.Value, &num2)

// 		var val []byte

// 		if num1 < num2 {
// 			val = encoder.Serialize(int32(1))
// 		} else {
// 			val = encoder.Serialize(int32(0))
// 		}

// 		assignOutput(0, val, "bool", expr, call)
// 		return nil
// 	} else {
// 		return err
// 	}
// }

// func gtF32 (arg1 *CXArgument, arg2 *CXArgument, expr *CXExpression, call *CXCall) error {
// 	if err := checkTwoTypes("f32.gt", "f32", "f32", arg1, arg2); err == nil {
// 		var num1 float32
// 		var num2 float32
// 		encoder.DeserializeRaw(*arg1.Value, &num1)
// 		encoder.DeserializeRaw(*arg2.Value, &num2)

// 		var val []byte

// 		if num1 > num2 {
// 			val = encoder.Serialize(int32(1))
// 		} else {
// 			val = encoder.Serialize(int32(0))
// 		}

// 		assignOutput(0, val, "bool", expr, call)
// 		return nil
// 	} else {
// 		return err
// 	}
// }

// func eqF32 (arg1 *CXArgument, arg2 *CXArgument, expr *CXExpression, call *CXCall) error {
// 	if err := checkTwoTypes("f32.eq", "f32", "f32", arg1, arg2); err == nil {
// 		var num1 float32
// 		var num2 float32
// 		encoder.DeserializeRaw(*arg1.Value, &num1)
// 		encoder.DeserializeRaw(*arg2.Value, &num2)

// 		var val []byte

// 		if num1 == num2 {
// 			val = encoder.Serialize(int32(1))
// 		} else {
// 			val = encoder.Serialize(int32(0))
// 		}

// 		assignOutput(0, val, "bool", expr, call)
// 		return nil
// 	} else {
// 		return err
// 	}
// }

// func uneqF32 (arg1 *CXArgument, arg2 *CXArgument, expr *CXExpression, call *CXCall) error {
// 	if err := checkTwoTypes("f32.uneq", "f32", "f32", arg1, arg2); err == nil {
// 		var num1 float32
// 		var num2 float32
// 		encoder.DeserializeRaw(*arg1.Value, &num1)
// 		encoder.DeserializeRaw(*arg2.Value, &num2)

// 		var val []byte

// 		if num1 != num2 {
// 			val = encoder.Serialize(int32(1))
// 		} else {
// 			val = encoder.Serialize(int32(0))
// 		}

// 		assignOutput(0, val, "bool", expr, call)
// 		return nil
// 	} else {
// 		return err
// 	}
// }

// func lteqF32 (arg1 *CXArgument, arg2 *CXArgument, expr *CXExpression, call *CXCall) error {
// 	if err := checkTwoTypes("f32.lteq", "f32", "f32", arg1, arg2); err == nil {
// 		var num1 float32
// 		var num2 float32
// 		encoder.DeserializeRaw(*arg1.Value, &num1)
// 		encoder.DeserializeRaw(*arg2.Value, &num2)

// 		var val []byte

// 		if num1 <= num2 {
// 			val = encoder.Serialize(int32(1))
// 		} else {
// 			val = encoder.Serialize(int32(0))
// 		}
		
// 		assignOutput(0, val, "bool", expr, call)
// 		return nil
// 	} else {
// 		return err
// 	}
// }

// func gteqF32 (arg1 *CXArgument, arg2 *CXArgument, expr *CXExpression, call *CXCall) error {
// 	if err := checkTwoTypes("f32.gteq", "f32", "f32", arg1, arg2); err == nil {
// 		var num1 float32
// 		var num2 float32
// 		encoder.DeserializeRaw(*arg1.Value, &num1)
// 		encoder.DeserializeRaw(*arg2.Value, &num2)

// 		var val []byte

// 		if num1 >= num2 {
// 			val = encoder.Serialize(int32(1))
// 		} else {
// 			val = encoder.Serialize(int32(0))
// 		}

// 		assignOutput(0, val, "bool", expr, call)
// 		return nil
// 	} else {
// 		return err
// 	}
// }

// func concatF32A (arg1 *CXArgument, arg2 *CXArgument, expr *CXExpression, call *CXCall) error {
// 	if err := checkTwoTypes("[]f32.concat", "[]f32", "[]f32", arg1, arg2); err == nil {
// 		var slice1 []float32
// 		var slice2 []float32
// 		encoder.DeserializeRaw(*arg1.Value, &slice1)
// 		encoder.DeserializeRaw(*arg2.Value, &slice2)

// 		output := append(slice1, slice2...)
// 		sOutput := encoder.Serialize(output)

// 		assignOutput(0, sOutput, "[]f32", expr, call)
// 		return nil
// 	} else {
// 		return err
// 	}
// }

// func appendF32A (arg1 *CXArgument, arg2 *CXArgument, expr *CXExpression, call *CXCall) error {
// 	if err := checkTwoTypes("[]f32.append", "[]f32", "f32", arg1, arg2); err == nil {
// 		var slice []float32
// 		var literal float32
// 		encoder.DeserializeRaw(*arg1.Value, &slice)
// 		encoder.DeserializeRaw(*arg2.Value, &literal)

// 		output := append(slice, literal)
// 		sOutput := encoder.Serialize(output)

// 		//*arg1.Value = sOutput
// 		assignOutput(0, sOutput, "[]f32", expr, call)
// 		return nil
// 	} else {
// 		return err
// 	}
// }

// func copyF32A (arg1 *CXArgument, arg2 *CXArgument, expr *CXExpression, call *CXCall) error {
// 	if err := checkTwoTypes("[]f32.copy", "[]f32", "[]f32", arg1, arg2); err == nil {
// 		copy(*arg1.Value, *arg2.Value)
// 		return nil
// 	} else {
// 		return err
// 	}
// }
