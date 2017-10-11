package base

import (
	"fmt"
	"errors"
	"time"
	"math/rand"
	"github.com/skycoin/skycoin/src/cipher/encoder"
)

func addI32 (arg1 *CXArgument, arg2 *CXArgument, expr *CXExpression, call *CXCall) error {
	if err := checkTwoTypes("i32.add", "i32", "i32", arg1, arg2); err == nil {
		// num1 := *arg1.Value
		// num2 := *arg2.Value
		// plus := byte(0)
		// output := make([]byte, 4, 4)
		// for c := 0; c < 4; c++ {
		// 	res := num1[c] + num2[c] + plus
		// 	if num2[c] > 255 - num1[c] || num2[c] + plus > 255 - num1[c] {
		// 		plus = byte(1)
		// 	} else {
		// 		plus = byte(0)
		// 	}
		// 	output[c] = res
		// }

		var num1 int32
		var num2 int32
		encoder.DeserializeAtomic(*arg1.Value, &num1)
		encoder.DeserializeAtomic(*arg2.Value, &num2)

		output := encoder.SerializeAtomic(int32(num1 + num2))
		
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

func subI32 (arg1 *CXArgument, arg2 *CXArgument, expr *CXExpression, call *CXCall) error {
	if err := checkTwoTypes("i32.sub", "i32", "i32", arg1, arg2); err == nil {
		var num1 int32
		var num2 int32
		encoder.DeserializeRaw(*arg1.Value, &num1)
		encoder.DeserializeRaw(*arg2.Value, &num2)

		output := encoder.SerializeAtomic(int32(num1 - num2))

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

func mulI32 (arg1 *CXArgument, arg2 *CXArgument, expr *CXExpression, call *CXCall) error {
	if err := checkTwoTypes("mulI32", "i32", "i32", arg1, arg2); err == nil {
		var num1 int32
		var num2 int32
		encoder.DeserializeRaw(*arg1.Value, &num1)
		encoder.DeserializeRaw(*arg2.Value, &num2)

		output := encoder.SerializeAtomic(int32(num1 * num2))

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

func divI32 (arg1 *CXArgument, arg2 *CXArgument, expr *CXExpression, call *CXCall) error {
	if err := checkTwoTypes("divI32", "i32", "i32", arg1, arg2); err == nil {
		var num1 int32
		var num2 int32
		encoder.DeserializeRaw(*arg1.Value, &num1)
		encoder.DeserializeRaw(*arg2.Value, &num2)

		if num2 == int32(0) {
			return errors.New("divI32: Division by 0")
		}

		output := encoder.SerializeAtomic(int32(num1 / num2))

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

func modI32 (arg1 *CXArgument, arg2 *CXArgument, expr *CXExpression, call *CXCall) error {
	if err := checkTwoTypes("modI32", "i32", "i32", arg1, arg2); err == nil {
		var num1 int32
		var num2 int32
		encoder.DeserializeRaw(*arg1.Value, &num1)
		encoder.DeserializeRaw(*arg2.Value, &num2)

		if num2 == int32(0) {
			return errors.New("modI32: Division by 0")
		}

		output := encoder.Serialize(int32(num1 % num2))

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

func andI32 (arg1 *CXArgument, arg2 *CXArgument, expr *CXExpression, call *CXCall) error {
	if err := checkTwoTypes("andI32", "i32", "i32", arg1, arg2); err == nil {
		var num1 int32
		var num2 int32
		encoder.DeserializeRaw(*arg1.Value, &num1)
		encoder.DeserializeRaw(*arg2.Value, &num2)
		
		output := encoder.Serialize(int32(num1 & num2))

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

func orI32 (arg1 *CXArgument, arg2 *CXArgument, expr *CXExpression, call *CXCall) error {
	if err := checkTwoTypes("orI32", "i32", "i32", arg1, arg2); err == nil {
		var num1 int32
		var num2 int32
		encoder.DeserializeRaw(*arg1.Value, &num1)
		encoder.DeserializeRaw(*arg2.Value, &num2)

		output := encoder.Serialize(int32(num1 | num2))

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

func xorI32 (arg1 *CXArgument, arg2 *CXArgument, expr *CXExpression, call *CXCall) error {
	if err := checkTwoTypes("xorI32", "i32", "i32", arg1, arg2); err == nil {
		var num1 int32
		var num2 int32
		encoder.DeserializeRaw(*arg1.Value, &num1)
		encoder.DeserializeRaw(*arg2.Value, &num2)

		output := encoder.Serialize(int32(num1 ^ num2))

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

func andNotI32 (arg1 *CXArgument, arg2 *CXArgument, expr *CXExpression, call *CXCall) error {
	if err := checkTwoTypes("andNotI32", "i32", "i32", arg1, arg2); err == nil {
		var num1 int32
		var num2 int32
		encoder.DeserializeRaw(*arg1.Value, &num1)
		encoder.DeserializeRaw(*arg2.Value, &num2)

		output := encoder.Serialize(int32(num1 &^ num2))

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

func randI32 (min *CXArgument, max *CXArgument, expr *CXExpression, call *CXCall) error {
	if err := checkTwoTypes("randI32", "i32", "i32", min, max); err == nil {
		var minimum int32
		encoder.DeserializeRaw(*min.Value, &minimum)

		var maximum int32
		encoder.DeserializeRaw(*max.Value, &maximum)

		if minimum > maximum {
			return errors.New(fmt.Sprintf("randI32: min must be less than max (%d !< %d)", minimum, maximum))
		}

		rand.Seed(time.Now().UTC().UnixNano())
		output := encoder.SerializeAtomic(int32(rand.Intn(int(maximum - minimum)) + int(minimum)))

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

func readI32A (arr *CXArgument, idx *CXArgument, expr *CXExpression, call *CXCall) error {
	if err := checkTwoTypes("readI32A", "[]i32", "i32", arr, idx); err == nil {
		var index int32
		encoder.DeserializeRaw(*idx.Value, &index)

		var size int32
		encoder.DeserializeAtomic((*arr.Value)[0:4], &size)

		if index < 0 {
			return errors.New(fmt.Sprintf("readI32A: negative index %d", index))
		}
		
		if index >= size {
			return errors.New(fmt.Sprintf("readI32A: index %d exceeds array of length %d", index, size))
		}

		var value int32
		encoder.DeserializeRaw((*arr.Value)[(index+1)*4:(index+2)*4], &value)
		sValue := encoder.Serialize(value)

		for _, def := range call.State {
			if def.Name == expr.OutputNames[0].Name {
				def.Value = &sValue
				return nil
			}
		}

		call.State = append(call.State, MakeDefinition(expr.OutputNames[0].Name, &sValue, "i32"))

		return nil
	} else {
		return err
	}
}

func writeI32A (arr *CXArgument, idx *CXArgument, val *CXArgument, expr *CXExpression, call *CXCall) error {
	if err := checkThreeTypes("writeI32A", "[]i32", "i32", "i32", arr, idx, val); err == nil {
		var index int32
		encoder.DeserializeRaw(*idx.Value, &index)

		var size int32
		encoder.DeserializeAtomic((*arr.Value)[0:4], &size)
		
		if index < 0 {
			return errors.New(fmt.Sprintf("writeI32A: negative index %d", index))
		}

		if index >= size {
			return errors.New(fmt.Sprintf("writeI32A: index %d exceeds array of length %d", index, size))
		}

		i := (int(index)+1)*4
		for c := 0; c < 4; c++ {
			(*arr.Value)[i + c] = (*val.Value)[c]
		}
		
		return nil
	} else {
		return err
	}
}

func lenI32A (arr *CXArgument, expr *CXExpression, call *CXCall) error {
	if err := checkType("lenI32A", "[]i32", arr); err == nil {
		var array []int32
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

func ltI32 (arg1 *CXArgument, arg2 *CXArgument, expr *CXExpression, call *CXCall) error {
	if err := checkTwoTypes("ltI32", "i32", "i32", arg1, arg2); err == nil {
		lt := false
		for i := 3; i >= 0; i-- {
			if (*arg1.Value)[i] < (*arg2.Value)[i] {
				lt = true
				break
			}
		}
		
		val := make([]byte, 4)
		
		if lt {
			val = []byte{1, 0, 0, 0}
		} else {
			val = []byte{0, 0, 0, 0}
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

func gtI32 (arg1 *CXArgument, arg2 *CXArgument, expr *CXExpression, call *CXCall) error {
	if err := checkTwoTypes("gtI32", "i32", "i32", arg1, arg2); err == nil {
		var num1 int32
		var num2 int32
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

func eqI32 (arg1 *CXArgument, arg2 *CXArgument, expr *CXExpression, call *CXCall) error {
	if err := checkTwoTypes("eqI32", "i32", "i32", arg1, arg2); err == nil {
		equal := true
		for i, b := range *arg1.Value {
			if b != (*arg2.Value)[i] {
				equal = false
				break
			}
		}
		val := make([]byte, 4)
		
		if equal {
			val = []byte{1, 0, 0, 0}
		} else {
			val = []byte{0, 0, 0, 0}
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

func lteqI32 (arg1 *CXArgument, arg2 *CXArgument, expr *CXExpression, call *CXCall) error {
	if err := checkTwoTypes("lteqI32", "i32", "i32", arg1, arg2); err == nil {
		var num1 int32
		var num2 int32
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

func gteqI32 (arg1 *CXArgument, arg2 *CXArgument, expr *CXExpression, call *CXCall) error {
	if err := checkTwoTypes("gteqI32", "i32", "i32", arg1, arg2); err == nil {
		var num1 int32
		var num2 int32
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

func concatI32A (arg1 *CXArgument, arg2 *CXArgument, expr *CXExpression, call *CXCall) error {
	if err := checkTwoTypes("[]i32.concat", "[]i32", "[]i32", arg1, arg2); err == nil {
		var slice1 []int32
		var slice2 []int32
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
		call.State = append(call.State, MakeDefinition(expr.OutputNames[0].Name, &sOutput, "[]i32"))
		return nil
	} else {
		return err
	}
}

func appendI32A (arg1 *CXArgument, arg2 *CXArgument, expr *CXExpression, call *CXCall) error {
	if err := checkTwoTypes("[]i32.append", "[]i32", "i32", arg1, arg2); err == nil {
		var slice []int32
		var literal int32
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
		call.State = append(call.State, MakeDefinition(expr.OutputNames[0].Name, &sOutput, "[]i32"))
		return nil
	} else {
		return err
	}
}

func copyI32A (arg1 *CXArgument, arg2 *CXArgument, expr *CXExpression, call *CXCall) error {
	if err := checkTwoTypes("[]i32.copy", "[]i32", "[]i32", arg1, arg2); err == nil {
		var slice1 []int32
		var slice2 []int32
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
