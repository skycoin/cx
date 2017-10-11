package base

import (
	"fmt"
	"errors"
	"github.com/skycoin/skycoin/src/cipher/encoder"
)

func readByteA (arr *CXArgument, idx *CXArgument, expr *CXExpression, call *CXCall) error {
	if err := checkTwoTypes("readByteA", "[]byte", "i32", arr, idx); err == nil {
		var index int32
		encoder.DeserializeRaw(*idx.Value, &index)

		size := len(*arr.Value)

		if index < 0 {
			return errors.New(fmt.Sprintf("readByteA: negative index %d", index))
		}
		
		if index >= int32(size) {
			return errors.New(fmt.Sprintf("readByteA: index %d exceeds array of length %d", index, size))
		}

		output := make([]byte, 1)
		output[0] = (*arr.Value)[index]

		for _, def := range call.State {
			if def.Name == expr.OutputNames[0].Name {
				def.Value = &output
				return nil
			}
		}
		
		call.State = append(call.State, MakeDefinition(expr.OutputNames[0].Name, &output, "byte"))

		return nil
	} else {
		return err
	}
}

func writeByteA (arr *CXArgument, idx *CXArgument, val *CXArgument, expr *CXExpression, call *CXCall) error {
	if err := checkThreeTypes("writeByteA", "[]byte", "i32", "byte", arr, idx, val); err == nil {
		var index int32
		encoder.DeserializeRaw(*idx.Value, &index)

		if index < 0 {
			return errors.New(fmt.Sprintf("writeByteA: negative index %d", index))
		}

		size := int32(len(*arr.Value))

		if index >= size {
			return errors.New(fmt.Sprintf("writeByteA: index %d exceeds array of length %d", index, size))
		}

		(*arr.Value)[index] = (*val.Value)[0]

		return nil
	} else {
		return err
	}
}

func lenByteA (arr *CXArgument, expr *CXExpression, call *CXCall) error {
	if err := checkType("lenByteA", "[]byte", arr); err == nil {
		output := encoder.SerializeAtomic(int32(len(*arr.Value)))

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

func ltByte (arg1 *CXArgument, arg2 *CXArgument, expr *CXExpression, call *CXCall) error {
	if err := checkTwoTypes("ltByte", "byte", "byte", arg1, arg2); err == nil {
		byte1 := (*arg1.Value)[0]
		byte2 := (*arg2.Value)[0]

		var val []byte

		if byte1 < byte2 {
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

func gtByte (arg1 *CXArgument, arg2 *CXArgument, expr *CXExpression, call *CXCall) error {
	if err := checkTwoTypes("gtByte", "byte", "byte", arg1, arg2); err == nil {
		byte1 := (*arg1.Value)[0]
		byte2 := (*arg2.Value)[0]

		var val []byte

		if byte1 > byte2 {
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

func eqByte (arg1 *CXArgument, arg2 *CXArgument, expr *CXExpression, call *CXCall) error {
	if err := checkTwoTypes("eqByte", "byte", "byte", arg1, arg2); err == nil {
		byte1 := (*arg1.Value)[0]
		byte2 := (*arg2.Value)[0]

		var val []byte

		if byte1 == byte2 {
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

func lteqByte (arg1 *CXArgument, arg2 *CXArgument, expr *CXExpression, call *CXCall) error {
	if err := checkTwoTypes("lteqByte", "byte", "byte", arg1, arg2); err == nil {
		byte1 := (*arg1.Value)[0]
		byte2 := (*arg2.Value)[0]

		var val []byte

		if byte1 <= byte2 {
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

func gteqByte (arg1 *CXArgument, arg2 *CXArgument, expr *CXExpression, call *CXCall) error {
	if err := checkTwoTypes("gteqByte", "byte", "byte", arg1, arg2); err == nil {
		byte1 := (*arg1.Value)[0]
		byte2 := (*arg2.Value)[0]

		var val []byte

		if byte1 >= byte2 {
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

func concatByteA (arg1 *CXArgument, arg2 *CXArgument, expr *CXExpression, call *CXCall) error {
	if err := checkTwoTypes("[]byte.concat", "[]byte", "[]byte", arg1, arg2); err == nil {
		output := append(*arg1.Value, *arg2.Value...)

		for _, def := range call.State {
			if def.Name == expr.OutputNames[0].Name {
				def.Value = &output
				return nil
			}
		}
		call.State = append(call.State, MakeDefinition(expr.OutputNames[0].Name, &output, "[]byte"))
		return nil
	} else {
		return err
	}
}

func appendByteA (arg1 *CXArgument, arg2 *CXArgument, expr *CXExpression, call *CXCall) error {
	if err := checkTwoTypes("[]byte.append", "[]byte", "byte", arg1, arg2); err == nil {
		output := append(*arg1.Value, *arg2.Value...)

		for _, def := range call.State {
			if def.Name == expr.OutputNames[0].Name {
				def.Value = &output
				return nil
			}
		}
		call.State = append(call.State, MakeDefinition(expr.OutputNames[0].Name, &output, "[]byte"))
		return nil
	} else {
		return err
	}
}

func copyByteA (arg1 *CXArgument, arg2 *CXArgument, expr *CXExpression, call *CXCall) error {
	if err := checkTwoTypes("[]byte.copy", "[]byte", "[]byte", arg1, arg2); err == nil {
		// var slice1 []int32
		// var slice2 []int32
		// encoder.DeserializeRaw(*arg1.Value, &slice1)
		// encoder.DeserializeRaw(*arg2.Value, &slice2)

		copy(*arg1.Value, *arg2.Value)
		//sOutput := encoder.Serialize(slice1)

		//*arg1.Value = sOutput
		return nil
	} else {
		return err
	}
}
