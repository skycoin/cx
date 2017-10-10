package base

import (
	"fmt"
	"errors"
	"github.com/skycoin/skycoin/src/cipher/encoder"
)

func addF32 (arg1 *CXArgument, arg2 *CXArgument, expr *CXExpression, call *CXCall) error {
	if err := checkTwoTypes("addF32", "f32", "f32", arg1, arg2); err == nil {
		
	} else {
		return err
	}
	
	var num1 float32
	var num2 float32
	encoder.DeserializeRaw(*arg1.Value, &num1)
	encoder.DeserializeRaw(*arg2.Value, &num2)

	output := encoder.Serialize(float32(num1 + num2))

	for _, def := range call.State {
		if def.Name == expr.OutputNames[0].Name {
			def.Value = &output
			return nil
		}
	}
	
	call.State = append(call.State, MakeDefinition(expr.OutputNames[0].Name, &output, "f32"))

	return nil
}

func subF32 (arg1 *CXArgument, arg2 *CXArgument, expr *CXExpression, call *CXCall) error {
	if err := checkTwoTypes("subF32", "f32", "f32", arg1, arg2); err == nil {
		var num1 float32
		var num2 float32
		encoder.DeserializeRaw(*arg1.Value, &num1)
		encoder.DeserializeRaw(*arg2.Value, &num2)

		output := encoder.Serialize(float32(num1 - num2))

		for _, def := range call.State {
			if def.Name == expr.OutputNames[0].Name {
				def.Value = &output
				return nil
			}
		}
		
		call.State = append(call.State, MakeDefinition(expr.OutputNames[0].Name, &output, "f32"))

		return nil
	} else {
		return err
	}
}

func mulF32 (arg1 *CXArgument, arg2 *CXArgument, expr *CXExpression, call *CXCall) error {
	if err := checkTwoTypes("mulF32", "f32", "f32", arg1, arg2); err == nil {
		var num1 float32
		var num2 float32
		encoder.DeserializeRaw(*arg1.Value, &num1)
		encoder.DeserializeRaw(*arg2.Value, &num2)

		output := encoder.Serialize(float32(num1 * num2))

		for _, def := range call.State {
			if def.Name == expr.OutputNames[0].Name {
				def.Value = &output
				return nil
			}
		}
		
		call.State = append(call.State, MakeDefinition(expr.OutputNames[0].Name, &output, "f32"))

		return nil
	} else {
		return err
	}
}

func divF32 (arg1 *CXArgument, arg2 *CXArgument, expr *CXExpression, call *CXCall) error {
	if err := checkTwoTypes("divF32", "f32", "f32", arg1, arg2); err == nil {
		var num1 float32
		var num2 float32
		encoder.DeserializeRaw(*arg1.Value, &num1)
		encoder.DeserializeRaw(*arg2.Value, &num2)

		if num2 == float32(0.0) {
			return errors.New("divF32: Division by 0")
		}

		output := encoder.Serialize(float32(num1 / num2))

		for _, def := range call.State {
			if def.Name == expr.OutputNames[0].Name {
				def.Value = &output
				return nil
			}
		}
		
		call.State = append(call.State, MakeDefinition(expr.OutputNames[0].Name, &output, "f32"))

		return nil
	} else {
		return err
	}	
}

func readF32A (arr *CXArgument, idx *CXArgument, expr *CXExpression, call *CXCall) error {
	if err := checkTwoTypes("readF32A", "[]f32", "i32", arr, idx); err == nil {
		var index int32
		encoder.DeserializeRaw(*idx.Value, &index)

		var size int32
		encoder.DeserializeAtomic((*arr.Value)[0:4], &size)
		
		// var array []float32
		// encoder.DeserializeRaw(*arr.Value, &array)

		if index < 0 {
			return errors.New(fmt.Sprintf("readF32A: negative index %d", index))
		}

		// if index >= int32(len(array)) {
		// 	return errors.New(fmt.Sprintf("readF32A: index %d exceeds array of length %d", index, len(array)))
		// }

		if index >= size {
			return errors.New(fmt.Sprintf("readF32A: index %d exceeds array of length %d", index, size))
		}

		var value float32
		encoder.DeserializeRaw((*arr.Value)[(index+1)*4:(index+2)*4], &value)
		sValue := encoder.Serialize(value)

		//output := encoder.Serialize(array[index])

		for _, def := range call.State {
			if def.Name == expr.OutputNames[0].Name {
				//def.Value = &output
				def.Value = &sValue
				return nil
			}
		}
		
		//call.State = append(call.State, MakeDefinition(expr.OutputNames[0].Name, &output, "f32"))
		call.State = append(call.State, MakeDefinition(expr.OutputNames[0].Name, &sValue, "f32"))

		return nil
	} else {
		return err
	}
}

func writeF32A (arr *CXArgument, idx *CXArgument, val *CXArgument, expr *CXExpression, call *CXCall) error {
	if err := checkThreeTypes("writeF32A", "[]f32", "i32", "f32", arr, idx, val); err == nil {
		var index int32
		encoder.DeserializeRaw(*idx.Value, &index)

		var value float32
		encoder.DeserializeRaw(*val.Value, &value)

		var array []float32
		encoder.DeserializeRaw(*arr.Value, &array)

		if index < 0 {
			return errors.New(fmt.Sprintf("writeF32A: negative index %d", index))
		}

		if index >= int32(len(array)) {
			return errors.New(fmt.Sprintf("writeF32A: index %d exceeds array of length %d", index, len(array)))
		}

		array[index] = value
		output := encoder.Serialize(array)

		*arr.Value = output
		return nil
	} else {
		return err
	}
}

func lenF32A (arr *CXArgument, expr *CXExpression, call *CXCall) error {
	if err := checkType("lenF32A", "[]f32", arr); err == nil {
		var array []float32
		encoder.DeserializeRaw(*arr.Value, &array)

		output := encoder.SerializeAtomic(int32(len(array)))

		for _, def := range call.State {
			if def.Name == expr.OutputNames[0].Name {
				def.Value = &output
				return nil
			}
		}
		
		call.State = append(call.State, MakeDefinition(expr.OutputNames[0].Name, &output, "i32"))
		return nil
	} else {
		return err
	}
}

func ltF32 (arg1 *CXArgument, arg2 *CXArgument, expr *CXExpression, call *CXCall) error {
	if err := checkTwoTypes("ltF32", "f32", "f32", arg1, arg2); err == nil {
		var num1 float32
		var num2 float32
		encoder.DeserializeRaw(*arg1.Value, &num1)
		encoder.DeserializeRaw(*arg2.Value, &num2)

		var val []byte

		if num1 < num2 {
			val = encoder.Serialize(int32(1))
		} else {
			val = encoder.Serialize(int32(0))
		}

		for _, def := range call.State {
			if def.Name == expr.OutputNames[0].Name {
				def.Value = &val
				return nil
			}
		}
		call.State = append(call.State, MakeDefinition(expr.OutputNames[0].Name, &val, "bool"))
		return nil
	} else {
		return err
	}
}

func gtF32 (arg1 *CXArgument, arg2 *CXArgument, expr *CXExpression, call *CXCall) error {
	if err := checkTwoTypes("gtF32", "f32", "f32", arg1, arg2); err == nil {
		var num1 float32
		var num2 float32
		encoder.DeserializeRaw(*arg1.Value, &num1)
		encoder.DeserializeRaw(*arg2.Value, &num2)

		var val []byte

		if num1 > num2 {
			val = encoder.Serialize(int32(1))
		} else {
			val = encoder.Serialize(int32(0))
		}

		for _, def := range call.State {
			if def.Name == expr.OutputNames[0].Name {
				def.Value = &val
				return nil
			}
		}
		call.State = append(call.State, MakeDefinition(expr.OutputNames[0].Name, &val, "bool"))
		return nil
	} else {
		return err
	}
}

func eqF32 (arg1 *CXArgument, arg2 *CXArgument, expr *CXExpression, call *CXCall) error {
	if err := checkTwoTypes("eqF32", "f32", "f32", arg1, arg2); err == nil {
		var num1 float32
		var num2 float32
		encoder.DeserializeRaw(*arg1.Value, &num1)
		encoder.DeserializeRaw(*arg2.Value, &num2)

		var val []byte

		if num1 == num2 {
			val = encoder.Serialize(int32(1))
		} else {
			val = encoder.Serialize(int32(0))
		}

		for _, def := range call.State {
			if def.Name == expr.OutputNames[0].Name {
				def.Value = &val
				return nil
			}
		}
		call.State = append(call.State, MakeDefinition(expr.OutputNames[0].Name, &val, "bool"))
		return nil
	} else {
		return err
	}
}

func lteqF32 (arg1 *CXArgument, arg2 *CXArgument, expr *CXExpression, call *CXCall) error {
	if err := checkTwoTypes("lteqF32", "f32", "f32", arg1, arg2); err == nil {
		var num1 float32
		var num2 float32
		encoder.DeserializeRaw(*arg1.Value, &num1)
		encoder.DeserializeRaw(*arg2.Value, &num2)

		var val []byte

		if num1 <= num2 {
			val = encoder.Serialize(int32(1))
		} else {
			val = encoder.Serialize(int32(0))
		}
		
		for _, def := range call.State {
			if def.Name == expr.OutputNames[0].Name {
				def.Value = &val
				return nil
			}
		}
		call.State = append(call.State, MakeDefinition(expr.OutputNames[0].Name, &val, "bool"))
		return nil
	} else {
		return err
	}
}

func gteqF32 (arg1 *CXArgument, arg2 *CXArgument, expr *CXExpression, call *CXCall) error {
	if err := checkTwoTypes("gteqF32", "f32", "f32", arg1, arg2); err == nil {
		var num1 float32
		var num2 float32
		encoder.DeserializeRaw(*arg1.Value, &num1)
		encoder.DeserializeRaw(*arg2.Value, &num2)

		var val []byte

		if num1 >= num2 {
			val = encoder.Serialize(int32(1))
		} else {
			val = encoder.Serialize(int32(0))
		}

		for _, def := range call.State {
			if def.Name == expr.OutputNames[0].Name {
				def.Value = &val
				return nil
			}
		}
		call.State = append(call.State, MakeDefinition(expr.OutputNames[0].Name, &val, "bool"))
		return nil
	} else {
		return err
	}
}

func concatF32A (arg1 *CXArgument, arg2 *CXArgument, expr *CXExpression, call *CXCall) error {
	if err := checkTwoTypes("[]f32.concat", "[]f32", "[]f32", arg1, arg2); err == nil {
		var slice1 []float32
		var slice2 []float32
		encoder.DeserializeRaw(*arg1.Value, &slice1)
		encoder.DeserializeRaw(*arg2.Value, &slice2)

		output := append(slice1, slice2...)
		sOutput := encoder.Serialize(output)

		for _, def := range call.State {
			if def.Name == expr.OutputNames[0].Name {
				def.Value = &sOutput
				return nil
			}
		}
		call.State = append(call.State, MakeDefinition(expr.OutputNames[0].Name, &sOutput, "[]f32"))
		return nil
	} else {
		return err
	}
}

func appendF32A (arg1 *CXArgument, arg2 *CXArgument, expr *CXExpression, call *CXCall) error {
	if err := checkTwoTypes("[]f32.append", "[]f32", "f32", arg1, arg2); err == nil {
		var slice []float32
		var literal float32
		encoder.DeserializeRaw(*arg1.Value, &slice)
		encoder.DeserializeRaw(*arg2.Value, &literal)

		output := append(slice, literal)
		sOutput := encoder.Serialize(output)

		for _, def := range call.State {
			if def.Name == expr.OutputNames[0].Name {
				def.Value = &sOutput
				return nil
			}
		}
		call.State = append(call.State, MakeDefinition(expr.OutputNames[0].Name, &sOutput, "[]f32"))
		return nil
	} else {
		return err
	}
}

func copyF32A (arg1 *CXArgument, arg2 *CXArgument, expr *CXExpression, call *CXCall) error {
	if err := checkTwoTypes("[]f32.copy", "[]f32", "[]f32", arg1, arg2); err == nil {
		var slice1 []float32
		var slice2 []float32
		encoder.DeserializeRaw(*arg1.Value, &slice1)
		encoder.DeserializeRaw(*arg2.Value, &slice2)

		copy(slice1, slice2)
		sOutput := encoder.Serialize(slice1)

		*arg1.Value = sOutput
		return nil
	} else {
		return err
	}
}
