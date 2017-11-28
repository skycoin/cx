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

		var size int32
		encoder.DeserializeAtomic((*arr.Value)[0:4], &size)

		if index < 0 {
			return errors.New(fmt.Sprintf("[]bool.read: negative index %d", index))
		}
		
		if index >= size {
			return errors.New(fmt.Sprintf("[]bool.read: index %d exceeds array of length %d", index, size))
		}

		var value int32
		encoder.DeserializeRaw((*arr.Value)[(index+1)*4:(index+2)*4], &value)
		output := encoder.Serialize(value)

		assignOutput(0, output, "bool", expr, call)
		return nil
	} else {
		return err
	}
}

func writeBoolA (arr *CXArgument, idx *CXArgument, val *CXArgument, expr *CXExpression, call *CXCall) error {
	if err := checkThreeTypes("writeBoolA", "[]bool", "i32", "bool", arr, idx, val); err == nil {
		var index int32
		encoder.DeserializeRaw(*idx.Value, &index)

		var size int32
		encoder.DeserializeAtomic((*arr.Value)[0:4], &size)

		if index < 0 {
			return errors.New(fmt.Sprintf("[]bool.write: negative index %d", index))
		}
		
		if index >= size {
			return errors.New(fmt.Sprintf("[]bool.write: index %d exceeds array of length %d", index, size))
		}

		// i := (int(index)+1) * 4
		// for c := 0; c < 4; c++ {
		// 	(*arr.Value)[i + c] = (*val.Value)[c]
		// }

		offset := int(index) * 4 + 4
		firstChunk := make([]byte, offset)
		secondChunk := make([]byte, len(*arr.Value) - offset)

		copy(firstChunk, (*arr.Value)[:offset])
		copy(secondChunk, (*arr.Value)[offset + 4:])

		final := append(firstChunk, *val.Value...)
		final = append(final, secondChunk...)

		assignOutput(0, final, "[]bool", expr, call)
		return nil
	} else {
		return err
	}
}

func lenBoolA (arr *CXArgument, expr *CXExpression, call *CXCall) error {
	if err := checkType("lenBoolA", "[]bool", arr); err == nil {
		size := (*arr.Value)[:4]
		assignOutput(0, size, "i32", expr, call)
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

		assignOutput(0, sOutput, "[]bool", expr, call)
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

		assignOutput(0, sOutput, "[]bool", expr, call)
		return nil
	} else {
		return err
	}
}

func copyBoolA (arg1 *CXArgument, arg2 *CXArgument, expr *CXExpression, call *CXCall) error {
	if err := checkTwoTypes("[]bool.copy", "[]bool", "[]bool", arg1, arg2); err == nil {
		copy(*arg1.Value, *arg2.Value)
		return nil
	} else {
		return err
	}
}
