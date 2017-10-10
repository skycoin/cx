package base

import (
	"fmt"
	"errors"
	"time"
	"math/rand"
	"github.com/skycoin/skycoin/src/cipher/encoder"
)

func addI64 (arg1 *CXArgument, arg2 *CXArgument, expr *CXExpression, call *CXCall) error {
	if err := checkTwoTypes("addI64", "i64", "i64", arg1, arg2); err == nil {
		var num1 int64
		var num2 int64
		encoder.DeserializeRaw(*arg1.Value, &num1)
		encoder.DeserializeRaw(*arg2.Value, &num2)

		output := encoder.SerializeAtomic(int64(num1 + num2))

		for _, def := range call.State {
			if def.Name == expr.OutputNames[0].Name {
				def.Value = &output
				return nil
			}
		}
		
		call.State = append(call.State, MakeDefinition(expr.OutputNames[0].Name, &output, "i64"))

		return nil
	} else {
		return err
	}
}

func subI64 (arg1 *CXArgument, arg2 *CXArgument, expr *CXExpression, call *CXCall) error {
	if err := checkTwoTypes("subI64", "i64", "i64", arg1, arg2); err == nil {
		var num1 int64
		var num2 int64
		encoder.DeserializeRaw(*arg1.Value, &num1)
		encoder.DeserializeRaw(*arg2.Value, &num2)

		output := encoder.SerializeAtomic(int64(num1 - num2))

		for _, def := range call.State {
			if def.Name == expr.OutputNames[0].Name {
				def.Value = &output
				return nil
			}
		}
		
		call.State = append(call.State, MakeDefinition(expr.OutputNames[0].Name, &output, "i64"))

		return nil
	} else {
		return err
	}
}

func mulI64 (arg1 *CXArgument, arg2 *CXArgument, expr *CXExpression, call *CXCall) error {
	if err := checkTwoTypes("mulI64", "i64", "i64", arg1, arg2); err == nil {
		var num1 int64
		var num2 int64
		encoder.DeserializeRaw(*arg1.Value, &num1)
		encoder.DeserializeRaw(*arg2.Value, &num2)

		output := encoder.SerializeAtomic(int64(num1 * num2))

		for _, def := range call.State {
			if def.Name == expr.OutputNames[0].Name {
				def.Value = &output
				return nil
			}
		}
		
		call.State = append(call.State, MakeDefinition(expr.OutputNames[0].Name, &output, "i64"))

		return nil
	} else {
		return err
	}
}

func divI64 (arg1 *CXArgument, arg2 *CXArgument, expr *CXExpression, call *CXCall) error {
	if err := checkTwoTypes("divI64", "i64", "i64", arg1, arg2); err == nil {
		var num1 int64
		var num2 int64
		encoder.DeserializeRaw(*arg1.Value, &num1)
		encoder.DeserializeRaw(*arg2.Value, &num2)

		if num2 == int64(0) {
			return errors.New("divI64: Division by 0")
		}
		
		output := encoder.SerializeAtomic(int64(num1 / num2))

		for _, def := range call.State {
			if def.Name == expr.OutputNames[0].Name {
				def.Value = &output
				return nil
			}
		}
		
		call.State = append(call.State, MakeDefinition(expr.OutputNames[0].Name, &output, "i64"))

		return nil
	} else {
		return err
	}
}

func modI64 (arg1 *CXArgument, arg2 *CXArgument, expr *CXExpression, call *CXCall) error {
	if err := checkTwoTypes("modI64", "i64", "i64", arg1, arg2); err == nil {
		var num1 int64
		var num2 int64
		encoder.DeserializeRaw(*arg1.Value, &num1)
		encoder.DeserializeRaw(*arg2.Value, &num2)

		if num2 == int64(0) {
			return errors.New("modI64: Division by 0")
		}

		output := encoder.Serialize(int64(num1 % num2))

		for _, def := range call.State {
			if def.Name == expr.OutputNames[0].Name {
				def.Value = &output
				return nil
			}
		}
		
		call.State = append(call.State, MakeDefinition(expr.OutputNames[0].Name, &output, "i64"))

		return nil
	} else {
		return err
	}
}

func andI64 (arg1 *CXArgument, arg2 *CXArgument, expr *CXExpression, call *CXCall) error {
	if err := checkTwoTypes("andI64", "i64", "i64", arg1, arg2); err == nil {
		var num1 int64
		var num2 int64
		encoder.DeserializeRaw(*arg1.Value, &num1)
		encoder.DeserializeRaw(*arg2.Value, &num2)

		output := encoder.Serialize(int64(num1 & num2))

		for _, def := range call.State {
			if def.Name == expr.OutputNames[0].Name {
				def.Value = &output
				return nil
			}
		}
		
		call.State = append(call.State, MakeDefinition(expr.OutputNames[0].Name, &output, "i64"))

		return nil
	} else {
		return err
	}
}

func orI64 (arg1 *CXArgument, arg2 *CXArgument, expr *CXExpression, call *CXCall) error {
	if err := checkTwoTypes("orI64", "i64", "i64", arg1, arg2); err == nil {
		var num1 int64
		var num2 int64
		encoder.DeserializeRaw(*arg1.Value, &num1)
		encoder.DeserializeRaw(*arg2.Value, &num2)

		output := encoder.Serialize(int64(num1 | num2))

		for _, def := range call.State {
			if def.Name == expr.OutputNames[0].Name {
				def.Value = &output
				return nil
			}
		}
		
		call.State = append(call.State, MakeDefinition(expr.OutputNames[0].Name, &output, "i64"))

		return nil
	} else {
		return err
	}
}

func xorI64 (arg1 *CXArgument, arg2 *CXArgument, expr *CXExpression, call *CXCall) error {
	if err := checkTwoTypes("xorI64", "i64", "i64", arg1, arg2); err == nil {
		var num1 int64
		var num2 int64
		encoder.DeserializeRaw(*arg1.Value, &num1)
		encoder.DeserializeRaw(*arg2.Value, &num2)

		output := encoder.Serialize(int64(num1 ^ num2))

		for _, def := range call.State {
			if def.Name == expr.OutputNames[0].Name {
				def.Value = &output
				return nil
			}
		}
		
		call.State = append(call.State, MakeDefinition(expr.OutputNames[0].Name, &output, "i64"))

		return nil
	} else {
		return err
	}
}

func andNotI64 (arg1 *CXArgument, arg2 *CXArgument, expr *CXExpression, call *CXCall) error {
	if err := checkTwoTypes("andNotI64", "i64", "i64", arg1, arg2); err == nil {
		var num1 int64
		var num2 int64
		encoder.DeserializeRaw(*arg1.Value, &num1)
		encoder.DeserializeRaw(*arg2.Value, &num2)

		output := encoder.Serialize(int64(num1 &^ num2))

		for _, def := range call.State {
			if def.Name == expr.OutputNames[0].Name {
				def.Value = &output
				return nil
			}
		}
		
		call.State = append(call.State, MakeDefinition(expr.OutputNames[0].Name, &output, "i64"))

		return nil
	} else {
		return err
	}
}

func randI64 (min *CXArgument, max *CXArgument, expr *CXExpression, call *CXCall) error {
	if err := checkTwoTypes("randI64", "i64", "i64", min, max); err == nil {
		var minimum int64
		encoder.DeserializeRaw(*min.Value, &minimum)

		var maximum int64
		encoder.DeserializeRaw(*max.Value, &maximum)

		if minimum > maximum {
			return errors.New(fmt.Sprintf("randI64: min must be less than max (%d !< %d)", minimum, maximum))
		}

		rand.Seed(time.Now().UTC().UnixNano())
		output := encoder.SerializeAtomic(int32(rand.Intn(int(maximum - minimum)) + int(minimum)))

		for _, def := range call.State {
			if def.Name == expr.OutputNames[0].Name {
				def.Value = &output
				return nil
			}
		}
		
		call.State = append(call.State, MakeDefinition(expr.OutputNames[0].Name, &output, "i64"))
		return nil
	} else {
		return err
	}
}

func readI64A (arr *CXArgument, idx *CXArgument, expr *CXExpression, call *CXCall) error {
	if err := checkTwoTypes("readI64A", "[]i64", "i32", arr, idx); err == nil {
		var index int32
		encoder.DeserializeRaw(*idx.Value, &index)

		var array []int64
		encoder.DeserializeRaw(*arr.Value, &array)

		if index < 0 {
			return errors.New(fmt.Sprintf("readI64A: negative index %d", index))
		}

		if index >= int32(len(array)) {
			return errors.New(fmt.Sprintf("readI64A: index %d exceeds array of length %d", index, len(array)))
		}

		output := encoder.Serialize(array[index])

		for _, def := range call.State {
			if def.Name == expr.OutputNames[0].Name {
				def.Value = &output
				return nil
			}
		}
		
		call.State = append(call.State, MakeDefinition(expr.OutputNames[0].Name, &output, "i64"))

		return nil
	} else {
		return err
	}
}

func writeI64A (arr *CXArgument, idx *CXArgument, val *CXArgument, expr *CXExpression, call *CXCall) error {
	if err := checkThreeTypes("writeI64A", "[]i64", "i32", "i64", arr, idx, val); err == nil {
		var index int32
		encoder.DeserializeRaw(*idx.Value, &index)

		var value int64
		encoder.DeserializeRaw(*val.Value, &value)

		var array []int64
		encoder.DeserializeRaw(*arr.Value, &array)

		if index < 0 {
			return errors.New(fmt.Sprintf("writeI64A: negative index %d", index))
		}

		if index >= int32(len(array)) {
			return errors.New(fmt.Sprintf("writeI64A: index %d exceeds array of length %d", index, len(array)))
		}

		array[index] = value
		output := encoder.Serialize(array)

		*arr.Value = output
		return nil
	} else {
		return err
	}
}

func lenI64A (arr *CXArgument, expr *CXExpression, call *CXCall) error {
	if err := checkType("lenI64A", "[]i64", arr); err == nil {
		var array []int64
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

func ltI64 (arg1 *CXArgument, arg2 *CXArgument, expr *CXExpression, call *CXCall) error {
	if err := checkTwoTypes("ltI64", "i64", "i64", arg1, arg2); err == nil {
		var num1 int64
		var num2 int64
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

func gtI64 (arg1 *CXArgument, arg2 *CXArgument, expr *CXExpression, call *CXCall) error {
	if err := checkTwoTypes("gtI64", "i64", "i64", arg1, arg2); err == nil {
		var num1 int64
		var num2 int64
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

func eqI64 (arg1 *CXArgument, arg2 *CXArgument, expr *CXExpression, call *CXCall) error {
	if err := checkTwoTypes("eqI64", "i64", "i64", arg1, arg2); err == nil {
		var num1 int64
		var num2 int64
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

func lteqI64 (arg1 *CXArgument, arg2 *CXArgument, expr *CXExpression, call *CXCall) error {
	if err := checkTwoTypes("lteqI64", "i64", "i64", arg1, arg2); err == nil {
		var num1 int64
		var num2 int64
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

func gteqI64 (arg1 *CXArgument, arg2 *CXArgument, expr *CXExpression, call *CXCall) error {
	if err := checkTwoTypes("gteqI64", "i64", "i64", arg1, arg2); err == nil {
		var num1 int64
		var num2 int64
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

func concatI64A (arg1 *CXArgument, arg2 *CXArgument, expr *CXExpression, call *CXCall) error {
	if err := checkTwoTypes("[]i64.concat", "[]i64", "[]i64", arg1, arg2); err == nil {
		var slice1 []int64
		var slice2 []int64
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
		call.State = append(call.State, MakeDefinition(expr.OutputNames[0].Name, &sOutput, "[]i64"))
		return nil
	} else {
		return err
	}
}

func appendI64A (arg1 *CXArgument, arg2 *CXArgument, expr *CXExpression, call *CXCall) error {
	if err := checkTwoTypes("[]i64.append", "[]i64", "i64", arg1, arg2); err == nil {
		var slice []int64
		var literal int64
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
		call.State = append(call.State, MakeDefinition(expr.OutputNames[0].Name, &sOutput, "[]i64"))
		return nil
	} else {
		return err
	}
}

func copyI64A (arg1 *CXArgument, arg2 *CXArgument, expr *CXExpression, call *CXCall) error {
	if err := checkTwoTypes("[]i64.copy", "[]i64", "[]i64", arg1, arg2); err == nil {
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
