package base

import (
	"fmt"
	"errors"
	"github.com/skycoin/skycoin/src/cipher/encoder"
)

func readBoolA (arr *CXArgument, idx *CXArgument, expr *CXExpression, call *CXCall) error {
	if err := checkTwoTypes("readBoolA", "[]bool", "i32", arr, idx); err == nil {
		var index int32
		encoder.DeserializeRaw(*idx.Value, &index)

		var array []int32
		encoder.DeserializeRaw(*arr.Value, &array)

		if index < 0 {
			return errors.New(fmt.Sprintf("readBoolA: negative index %d", index))
		}
		
		if index >= int32(len(array)) {
			return errors.New(fmt.Sprintf("readBoolA: index %d exceeds array of length %d", index, len(array)))
		}

		output := encoder.Serialize(array[index])

		for _, def := range call.State {
			if def.Name == expr.OutputNames[0].Name {
				def.Value = &output
				return nil
			}
		}
		
		call.State = append(call.State, MakeDefinition(expr.OutputNames[0].Name, &output, "bool"))

		return nil
	} else {
		return err
	}
}

func writeBoolA (arr *CXArgument, idx *CXArgument, val *CXArgument, expr *CXExpression, call *CXCall) error {
	if err := checkThreeTypes("writeBoolA", "[]bool", "i32", "bool", arr, idx, val); err == nil {
		var index int32
		encoder.DeserializeRaw(*idx.Value, &index)

		var value int32
		encoder.DeserializeRaw(*val.Value, &value)
		
		var array []int32
		encoder.DeserializeRaw(*arr.Value, &array)

		if index < 0 {
			return errors.New(fmt.Sprintf("writeBoolA: negative index %d", index))
		}
		
		if index >= int32(len(*arr.Value)) {
			return errors.New(fmt.Sprintf("writeBoolA: index %d exceeds array of length %d", index, len(array)))
		}

		array[index] = value
		output := encoder.Serialize(array)

		*arr.Value = output
		return nil
	} else {
		return err
	}
}

func lenBoolA (arr *CXArgument, expr *CXExpression, call *CXCall) error {
	if err := checkType("lenBoolA", "[]bool", arr); err == nil {
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

func concatBoolA (arg1 *CXArgument, arg2 *CXArgument, expr *CXExpression, call *CXCall) error {
	if err := checkTwoTypes("[]bool.concat", "[]bool", "[]bool", arg1, arg2); err == nil {
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
		call.State = append(call.State, MakeDefinition(expr.OutputNames[0].Name, &sOutput, "[]bool"))
		return nil
	} else {
		return err
	}
}

func appendBoolA (arg1 *CXArgument, arg2 *CXArgument, expr *CXExpression, call *CXCall) error {
	if err := checkTwoTypes("[]bool.append", "[]bool", "bool", arg1, arg2); err == nil {
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
		call.State = append(call.State, MakeDefinition(expr.OutputNames[0].Name, &sOutput, "[]bool"))
		return nil
	} else {
		return err
	}
}

func copyBoolA (arg1 *CXArgument, arg2 *CXArgument, expr *CXExpression, call *CXCall) error {
	if err := checkTwoTypes("[]bool.copy", "[]bool", "[]bool", arg1, arg2); err == nil {
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
