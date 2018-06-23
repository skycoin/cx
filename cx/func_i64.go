package base

import (
	"errors"
	"fmt"
	"math"
	"math/rand"

	"github.com/skycoin/skycoin/src/cipher/encoder"
)

func addI64(arg1 *CXArgument, arg2 *CXArgument, expr *CXExpression, call *CXCall) error {
	if err := checkTwoTypes("i64.add", "i64", "i64", arg1, arg2); err == nil {
		var num1 int64
		var num2 int64
		encoder.DeserializeRaw(*arg1.Value, &num1)
		encoder.DeserializeRaw(*arg2.Value, &num2)

		output := encoder.SerializeAtomic(int64(num1 + num2))

		assignOutput(0, output, "i64", expr, call)
		return nil
	} else {
		return err
	}
}

func subI64(arg1 *CXArgument, arg2 *CXArgument, expr *CXExpression, call *CXCall) error {
	if err := checkTwoTypes("i64.sub", "i64", "i64", arg1, arg2); err == nil {
		var num1 int64
		var num2 int64
		encoder.DeserializeRaw(*arg1.Value, &num1)
		encoder.DeserializeRaw(*arg2.Value, &num2)

		output := encoder.SerializeAtomic(int64(num1 - num2))

		assignOutput(0, output, "i64", expr, call)
		return nil
	} else {
		return err
	}
}

func mulI64(arg1 *CXArgument, arg2 *CXArgument, expr *CXExpression, call *CXCall) error {
	if err := checkTwoTypes("i64mul", "i64", "i64", arg1, arg2); err == nil {
		var num1 int64
		var num2 int64
		encoder.DeserializeRaw(*arg1.Value, &num1)
		encoder.DeserializeRaw(*arg2.Value, &num2)

		output := encoder.SerializeAtomic(int64(num1 * num2))

		assignOutput(0, output, "i64", expr, call)
		return nil
	} else {
		return err
	}
}

func divI64(arg1 *CXArgument, arg2 *CXArgument, expr *CXExpression, call *CXCall) error {
	if err := checkTwoTypes("i64.div", "i64", "i64", arg1, arg2); err == nil {
		var num1 int64
		var num2 int64
		encoder.DeserializeRaw(*arg1.Value, &num1)
		encoder.DeserializeRaw(*arg2.Value, &num2)

		if num2 == int64(0) {
			return errors.New("i64.div: Division by 0")
		}

		output := encoder.SerializeAtomic(int64(num1 / num2))

		assignOutput(0, output, "i64", expr, call)
		return nil
	} else {
		return err
	}
}

func powI64(arg1 *CXArgument, arg2 *CXArgument, expr *CXExpression, call *CXCall) error {
	if err := checkTwoTypes("i64.pow", "i64", "i64", arg1, arg2); err == nil {
		var num1 int64
		var num2 int64
		encoder.DeserializeRaw(*arg1.Value, &num1)
		encoder.DeserializeRaw(*arg2.Value, &num2)

		output := encoder.Serialize(int64(math.Pow(float64(num1), float64(num2))))

		assignOutput(0, output, "i64", expr, call)
		return nil
	} else {
		return err
	}
}

func sqrtI64(arg1 *CXArgument, expr *CXExpression, call *CXCall) error {
	if err := checkType("i64.sqrt", "i64", arg1); err == nil {
		var num1 int64
		encoder.DeserializeRaw(*arg1.Value, &num1)

		output := encoder.Serialize(int64(math.Sqrt(float64(num1))))

		assignOutput(0, output, "i64", expr, call)
		return nil
	} else {
		return err
	}
}

func absI64(arg1 *CXArgument, expr *CXExpression, call *CXCall) error {
	if err := checkType("i64.abs", "i64", arg1); err == nil {
		var num1 int64
		encoder.DeserializeRaw(*arg1.Value, &num1)

		output := encoder.Serialize(int64(math.Abs(float64(num1))))

		assignOutput(0, output, "i64", expr, call)
		return nil
	} else {
		return err
	}
}

func modI64(arg1 *CXArgument, arg2 *CXArgument, expr *CXExpression, call *CXCall) error {
	if err := checkTwoTypes("i64.mod", "i64", "i64", arg1, arg2); err == nil {
		var num1 int64
		var num2 int64
		encoder.DeserializeRaw(*arg1.Value, &num1)
		encoder.DeserializeRaw(*arg2.Value, &num2)

		if num2 == int64(0) {
			return errors.New("i64.mod: Division by 0")
		}

		output := encoder.Serialize(int64(num1 % num2))

		assignOutput(0, output, "i64", expr, call)
		return nil
	} else {
		return err
	}
}

func andI64(arg1 *CXArgument, arg2 *CXArgument, expr *CXExpression, call *CXCall) error {
	if err := checkTwoTypes("i64.and", "i64", "i64", arg1, arg2); err == nil {
		var num1 int64
		var num2 int64
		encoder.DeserializeRaw(*arg1.Value, &num1)
		encoder.DeserializeRaw(*arg2.Value, &num2)

		output := encoder.Serialize(int64(num1 & num2))

		assignOutput(0, output, "i64", expr, call)
		return nil
	} else {
		return err
	}
}

func orI64(arg1 *CXArgument, arg2 *CXArgument, expr *CXExpression, call *CXCall) error {
	if err := checkTwoTypes("i64.or", "i64", "i64", arg1, arg2); err == nil {
		var num1 int64
		var num2 int64
		encoder.DeserializeRaw(*arg1.Value, &num1)
		encoder.DeserializeRaw(*arg2.Value, &num2)

		output := encoder.Serialize(int64(num1 | num2))

		assignOutput(0, output, "i64", expr, call)
		return nil
	} else {
		return err
	}
}

func xorI64(arg1 *CXArgument, arg2 *CXArgument, expr *CXExpression, call *CXCall) error {
	if err := checkTwoTypes("i64.xor", "i64", "i64", arg1, arg2); err == nil {
		var num1 int64
		var num2 int64
		encoder.DeserializeRaw(*arg1.Value, &num1)
		encoder.DeserializeRaw(*arg2.Value, &num2)

		output := encoder.Serialize(int64(num1 ^ num2))

		assignOutput(0, output, "i64", expr, call)
		return nil
	} else {
		return err
	}
}

func andNotI64(arg1 *CXArgument, arg2 *CXArgument, expr *CXExpression, call *CXCall) error {
	if err := checkTwoTypes("i64.bitclear", "i64", "i64", arg1, arg2); err == nil {
		var num1 int64
		var num2 int64
		encoder.DeserializeRaw(*arg1.Value, &num1)
		encoder.DeserializeRaw(*arg2.Value, &num2)

		output := encoder.Serialize(int64(num1 &^ num2))

		assignOutput(0, output, "i64", expr, call)
		return nil
	} else {
		return err
	}
}

func shiftLeftI64(arg1, arg2 *CXArgument, expr *CXExpression, call *CXCall) error {
	if err := checkTwoTypes("i64.bitshl", "i64", "i64", arg1, arg2); err == nil {
		var num1 uint64
		var num2 uint64
		encoder.DeserializeRaw(*arg1.Value, &num1)
		encoder.DeserializeRaw(*arg2.Value, &num2)

		output := encoder.Serialize(int64(num1 << num2))

		assignOutput(0, output, "i64", expr, call)
		return nil
	} else {
		return err
	}
}

func shiftRightI64(arg1, arg2 *CXArgument, expr *CXExpression, call *CXCall) error {
	if err := checkTwoTypes("i64.bitshr", "i64", "i64", arg1, arg2); err == nil {
		var num1 uint64
		var num2 uint64
		encoder.DeserializeRaw(*arg1.Value, &num1)
		encoder.DeserializeRaw(*arg2.Value, &num2)

		output := encoder.Serialize(int64(num1 >> num2))

		assignOutput(0, output, "i64", expr, call)
		return nil
	} else {
		return err
	}
}

func randI64(min *CXArgument, max *CXArgument, expr *CXExpression, call *CXCall) error {
	if err := checkTwoTypes("i64.rand", "i64", "i64", min, max); err == nil {
		var minimum int64
		encoder.DeserializeRaw(*min.Value, &minimum)

		var maximum int64
		encoder.DeserializeRaw(*max.Value, &maximum)

		if minimum > maximum {
			return errors.New(fmt.Sprintf("i64.rand: min must be less than max (%d !< %d)", minimum, maximum))
		}

		output := encoder.Serialize(int64(rand.Intn(int(maximum-minimum)) + int(minimum)))

		assignOutput(0, output, "i64", expr, call)
		return nil
	} else {
		return err
	}
}

func readI64A(arr *CXArgument, idx *CXArgument, expr *CXExpression, call *CXCall) error {
	if err := checkTwoTypes("[]i64.read", "[]i64", "i32", arr, idx); err == nil {
		var index int32
		encoder.DeserializeRaw(*idx.Value, &index)

		var size int32
		encoder.DeserializeAtomic((*arr.Value)[0:4], &size)

		if index < 0 {
			return errors.New(fmt.Sprintf("[]i64.read: negative index %d", index))
		}

		if index >= size {
			return errors.New(fmt.Sprintf("[]i64.read: index %d exceeds array of length %d", index, size))
		}

		var value int64
		//encoder.DeserializeRaw((*arr.Value)[(index+1)*4:(index+2)*4], &value)
		encoder.DeserializeRaw((*arr.Value)[((index)*8)+4:((index+1)*8)+4], &value)
		output := encoder.Serialize(value)

		assignOutput(0, output, "i64", expr, call)
		return nil
	} else {
		return err
	}
}

func writeI64A(arr *CXArgument, idx *CXArgument, val *CXArgument, expr *CXExpression, call *CXCall) error {
	if err := checkThreeTypes("[]i64.write", "[]i64", "i32", "i64", arr, idx, val); err == nil {
		var index int32
		encoder.DeserializeRaw(*idx.Value, &index)

		var size int32
		encoder.DeserializeAtomic((*arr.Value)[0:4], &size)

		if index < 0 {
			return errors.New(fmt.Sprintf("[]i64.write: negative index %d", index))
		}

		if index >= size {
			return errors.New(fmt.Sprintf("[]i64.write: index %d exceeds array of length %d", index, size))
		}

		// i := (int(index)*8)+4
		// for c := 0; c < 8; c++ {
		// 	(*arr.Value)[i + c] = (*val.Value)[c]
		// }

		offset := int(index)*8 + 4
		firstChunk := make([]byte, offset)
		secondChunk := make([]byte, len(*arr.Value)-(offset+8))

		copy(firstChunk, (*arr.Value)[:offset])
		copy(secondChunk, (*arr.Value)[offset+8:])

		final := append(firstChunk, *val.Value...)
		final = append(final, secondChunk...)

		assignOutput(0, final, "[]i64", expr, call)

		return nil
	} else {
		return err
	}
}

func lenI64A(arr *CXArgument, expr *CXExpression, call *CXCall) error {
	if err := checkType("[]i64.len", "[]i64", arr); err == nil {
		size := (*arr.Value)[:4]
		assignOutput(0, size, "i32", expr, call)
		return nil
	} else {
		return err
	}
}

func ltI64(arg1 *CXArgument, arg2 *CXArgument, expr *CXExpression, call *CXCall) error {
	if err := checkTwoTypes("i64.lt", "i64", "i64", arg1, arg2); err == nil {
		var num1 int64
		var num2 int64
		encoder.DeserializeRaw(*arg1.Value, &num1)
		encoder.DeserializeRaw(*arg2.Value, &num2)

		var val []byte

		if num1 < num2 {
			val = encoder.Serialize(true)
		} else {
			val = encoder.Serialize(false)
		}

		assignOutput(0, val, "bool", expr, call)
		return nil
	} else {
		return err
	}
}

func gtI64(arg1 *CXArgument, arg2 *CXArgument, expr *CXExpression, call *CXCall) error {
	if err := checkTwoTypes("i64.gt", "i64", "i64", arg1, arg2); err == nil {
		var num1 int64
		var num2 int64
		encoder.DeserializeRaw(*arg1.Value, &num1)
		encoder.DeserializeRaw(*arg2.Value, &num2)

		var val []byte

		if num1 > num2 {
			val = encoder.Serialize(true)
		} else {
			val = encoder.Serialize(false)
		}

		assignOutput(0, val, "bool", expr, call)
		return nil
	} else {
		return err
	}
}

func eqI64(arg1 *CXArgument, arg2 *CXArgument, expr *CXExpression, call *CXCall) error {
	if err := checkTwoTypes("i64.eq", "i64", "i64", arg1, arg2); err == nil {
		var num1 int64
		var num2 int64
		encoder.DeserializeRaw(*arg1.Value, &num1)
		encoder.DeserializeRaw(*arg2.Value, &num2)

		var val []byte

		if num1 == num2 {
			val = encoder.Serialize(true)
		} else {
			val = encoder.Serialize(false)
		}

		assignOutput(0, val, "bool", expr, call)
		return nil
	} else {
		return err
	}
}

func uneqI64(arg1 *CXArgument, arg2 *CXArgument, expr *CXExpression, call *CXCall) error {
	if err := checkTwoTypes("i64.uneq", "i64", "i64", arg1, arg2); err == nil {
		var num1 int64
		var num2 int64
		encoder.DeserializeRaw(*arg1.Value, &num1)
		encoder.DeserializeRaw(*arg2.Value, &num2)

		var val []byte

		if num1 != num2 {
			val = encoder.Serialize(true)
		} else {
			val = encoder.Serialize(false)
		}

		assignOutput(0, val, "bool", expr, call)
		return nil
	} else {
		return err
	}
}

func lteqI64(arg1 *CXArgument, arg2 *CXArgument, expr *CXExpression, call *CXCall) error {
	if err := checkTwoTypes("i64.lteq", "i64", "i64", arg1, arg2); err == nil {
		var num1 int64
		var num2 int64
		encoder.DeserializeRaw(*arg1.Value, &num1)
		encoder.DeserializeRaw(*arg2.Value, &num2)

		var val []byte

		if num1 <= num2 {
			val = encoder.Serialize(true)
		} else {
			val = encoder.Serialize(false)
		}

		assignOutput(0, val, "bool", expr, call)
		return nil
	} else {
		return err
	}
}

func gteqI64(arg1 *CXArgument, arg2 *CXArgument, expr *CXExpression, call *CXCall) error {
	if err := checkTwoTypes("i64.gteq", "i64", "i64", arg1, arg2); err == nil {
		var num1 int64
		var num2 int64
		encoder.DeserializeRaw(*arg1.Value, &num1)
		encoder.DeserializeRaw(*arg2.Value, &num2)

		var val []byte

		if num1 >= num2 {
			val = encoder.Serialize(true)
		} else {
			val = encoder.Serialize(false)
		}

		assignOutput(0, val, "bool", expr, call)
		return nil
	} else {
		return err
	}
}

func concatI64A(arg1 *CXArgument, arg2 *CXArgument, expr *CXExpression, call *CXCall) error {
	if err := checkTwoTypes("[]i64.concat", "[]i64", "[]i64", arg1, arg2); err == nil {
		var slice1 []int64
		var slice2 []int64
		encoder.DeserializeRaw(*arg1.Value, &slice1)
		encoder.DeserializeRaw(*arg2.Value, &slice2)

		output := append(slice1, slice2...)
		sOutput := encoder.Serialize(output)

		assignOutput(0, sOutput, "[]i64", expr, call)
		return nil
	} else {
		return err
	}
}

func appendI64A(arg1 *CXArgument, arg2 *CXArgument, expr *CXExpression, call *CXCall) error {
	if err := checkTwoTypes("[]i64.append", "[]i64", "i64", arg1, arg2); err == nil {
		var slice []int64
		var literal int64
		encoder.DeserializeRaw(*arg1.Value, &slice)
		encoder.DeserializeRaw(*arg2.Value, &literal)

		output := append(slice, literal)
		sOutput := encoder.Serialize(output)

		//*arg1.Value = sOutput
		assignOutput(0, sOutput, "[]i64", expr, call)
		return nil
	} else {
		return err
	}
}

func copyI64A(arg1 *CXArgument, arg2 *CXArgument, expr *CXExpression, call *CXCall) error {
	if err := checkTwoTypes("[]i64.copy", "[]i64", "[]i64", arg1, arg2); err == nil {
		copy(*arg1.Value, *arg2.Value)
		return nil
	} else {
		return err
	}
}
