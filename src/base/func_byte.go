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

		var size int32
		encoder.DeserializeAtomic((*arr.Value)[0:4], &size)

		if index < 0 {
			return errors.New(fmt.Sprintf("[]byte.read: negative index %d", index))
		}
		
		if index >= size {
			return errors.New(fmt.Sprintf("[]byte.read: index %d exceeds array of length %d", index, size))
		}
		
		var value byte
		encoder.DeserializeRaw((*arr.Value)[index+4:(index+1)+4], &value)

		output := make([]byte, 1)
		output[0] = value

		assignOutput(&output, "byte", expr, call)
		return nil
	} else {
		return err
	}
}

func writeByteA (arr *CXArgument, idx *CXArgument, val *CXArgument, expr *CXExpression, call *CXCall) error {
	if err := checkThreeTypes("writeByteA", "[]byte", "i32", "byte", arr, idx, val); err == nil {
		var index int32
		encoder.DeserializeRaw(*idx.Value, &index)

		var size int32
		encoder.DeserializeAtomic((*arr.Value)[0:4], &size)

		if index < 0 {
			return errors.New(fmt.Sprintf("[]byte.write: negative index %d", index))
		}

		if index >= size {
			return errors.New(fmt.Sprintf("[]byte.write: index %d exceeds array of length %d", index, size))
		}
		
		// (*arr.Value)[index + 4] = (*val.Value)[0]

		offset := int(index) * 1 + 4
		firstChunk := make([]byte, offset)
		secondChunk := make([]byte, len(*arr.Value) - (offset + 1))

		copy(firstChunk, (*arr.Value)[:offset])
		copy(secondChunk, (*arr.Value)[offset + 1:])

		final := append(firstChunk, *val.Value...)
		final = append(final, secondChunk...)

		assignOutput(&final, "[]byte", expr, call)
		return nil
	} else {
		return err
	}
}

func lenByteA (arr *CXArgument, expr *CXExpression, call *CXCall) error {
	if err := checkType("lenByteA", "[]byte", arr); err == nil {
		output := encoder.SerializeAtomic(int32(len((*arr.Value)[4:])))

		assignOutput(&output, "i32", expr, call)
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

		assignOutput(&val, "bool", expr, call)
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

		assignOutput(&val, "bool", expr, call)
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

		assignOutput(&val, "bool", expr, call)
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

		assignOutput(&val, "bool", expr, call)
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

		assignOutput(&val, "bool", expr, call)
		return nil
	} else {
		return err
	}
}

func concatByteA (arg1 *CXArgument, arg2 *CXArgument, expr *CXExpression, call *CXCall) error {
	if err := checkTwoTypes("[]byte.concat", "[]byte", "[]byte", arg1, arg2); err == nil {
		var slice1 []byte
		var slice2 []byte
		encoder.DeserializeRaw(*arg1.Value, &slice1)
		encoder.DeserializeRaw(*arg2.Value, &slice2)

		output := append(slice1, slice2...)
		sOutput := encoder.Serialize(output)

		assignOutput(&sOutput, "[]byte", expr, call)
		return nil
	} else {
		return err
	}
}

func appendByteA (arg1 *CXArgument, arg2 *CXArgument, expr *CXExpression, call *CXCall) error {
	if err := checkTwoTypes("[]byte.append", "[]byte", "byte", arg1, arg2); err == nil {
		var slice []byte
		encoder.DeserializeRaw(*arg1.Value, &slice)

		output := append(slice, (*arg2.Value)[0])
		sOutput := encoder.Serialize(output)

		//*arg1.Value = sOutput
		assignOutput(&sOutput, "[]byte", expr, call)
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
