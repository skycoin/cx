package base

import (
	"fmt"
	"errors"
	"time"
	"math/rand"
	"os"
	"strconv"
	"bufio"
	"strings"
	"github.com/skycoin/skycoin/src/cipher/encoder"
)

func addI32 (arg1 *CXArgument, arg2 *CXArgument, expr *CXExpression, call *CXCall) error {
	if err := checkTwoTypes("i32.add", "i32", "i32", arg1, arg2); err == nil {
		var num1 int32
		var num2 int32
		encoder.DeserializeAtomic(*arg1.Value, &num1)
		encoder.DeserializeAtomic(*arg2.Value, &num2)

		output := encoder.SerializeAtomic(int32(num1 + num2))

		assignOutput(&output, "i32", expr, call)
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

		assignOutput(&output, "i32", expr, call)
		return nil
	} else {
		return err
	}
}

func mulI32 (arg1 *CXArgument, arg2 *CXArgument, expr *CXExpression, call *CXCall) error {
	if err := checkTwoTypes("i32.mul", "i32", "i32", arg1, arg2); err == nil {
		var num1 int32
		var num2 int32
		encoder.DeserializeRaw(*arg1.Value, &num1)
		encoder.DeserializeRaw(*arg2.Value, &num2)

		output := encoder.SerializeAtomic(int32(num1 * num2))

		assignOutput(&output, "i32", expr, call)
		return nil
	} else {
		return err
	}
}

func divI32 (arg1 *CXArgument, arg2 *CXArgument, expr *CXExpression, call *CXCall) error {
	if err := checkTwoTypes("i32.div", "i32", "i32", arg1, arg2); err == nil {
		var num1 int32
		var num2 int32
		encoder.DeserializeRaw(*arg1.Value, &num1)
		encoder.DeserializeRaw(*arg2.Value, &num2)

		if num2 == int32(0) {
			return errors.New("i32.div: Division by 0")
		}

		output := encoder.SerializeAtomic(int32(num1 / num2))

		assignOutput(&output, "i32", expr, call)
		return nil
	} else {
		return err
	}
}

func modI32 (arg1 *CXArgument, arg2 *CXArgument, expr *CXExpression, call *CXCall) error {
	if err := checkTwoTypes("i32.mod", "i32", "i32", arg1, arg2); err == nil {
		var num1 int32
		var num2 int32
		encoder.DeserializeRaw(*arg1.Value, &num1)
		encoder.DeserializeRaw(*arg2.Value, &num2)

		if num2 == int32(0) {
			return errors.New("i32.mod: Division by 0")
		}

		output := encoder.Serialize(int32(num1 % num2))

		assignOutput(&output, "i32", expr, call)
		return nil
	} else {
		return err
	}
}

func andI32 (arg1 *CXArgument, arg2 *CXArgument, expr *CXExpression, call *CXCall) error {
	if err := checkTwoTypes("i32.and", "i32", "i32", arg1, arg2); err == nil {
		var num1 int32
		var num2 int32
		encoder.DeserializeRaw(*arg1.Value, &num1)
		encoder.DeserializeRaw(*arg2.Value, &num2)
		
		output := encoder.Serialize(int32(num1 & num2))

		assignOutput(&output, "i32", expr, call)
		return nil
	} else {
		return err
	}
}

func orI32 (arg1 *CXArgument, arg2 *CXArgument, expr *CXExpression, call *CXCall) error {
	if err := checkTwoTypes("i32.or", "i32", "i32", arg1, arg2); err == nil {
		var num1 int32
		var num2 int32
		encoder.DeserializeRaw(*arg1.Value, &num1)
		encoder.DeserializeRaw(*arg2.Value, &num2)

		output := encoder.Serialize(int32(num1 | num2))

		assignOutput(&output, "i32", expr, call)
		return nil
	} else {
		return err
	}
}

func xorI32 (arg1 *CXArgument, arg2 *CXArgument, expr *CXExpression, call *CXCall) error {
	if err := checkTwoTypes("i32.xor", "i32", "i32", arg1, arg2); err == nil {
		var num1 int32
		var num2 int32
		encoder.DeserializeRaw(*arg1.Value, &num1)
		encoder.DeserializeRaw(*arg2.Value, &num2)

		output := encoder.Serialize(int32(num1 ^ num2))

		assignOutput(&output, "i32", expr, call)
		return nil
	} else {
		return err
	}
}

func andNotI32 (arg1 *CXArgument, arg2 *CXArgument, expr *CXExpression, call *CXCall) error {
	if err := checkTwoTypes("i32.bitclear", "i32", "i32", arg1, arg2); err == nil {
		var num1 int32
		var num2 int32
		encoder.DeserializeRaw(*arg1.Value, &num1)
		encoder.DeserializeRaw(*arg2.Value, &num2)

		output := encoder.Serialize(int32(num1 &^ num2))

		assignOutput(&output, "i32", expr, call)
		return nil
	} else {
		return err
	}
}

func shiftLeftI32 (arg1, arg2 *CXArgument, expr *CXExpression, call *CXCall) error {
	if err := checkTwoTypes("i32.bitshl", "i32", "i32", arg1, arg2); err == nil {
		var num1 uint32
		var num2 uint32
		encoder.DeserializeRaw(*arg1.Value, &num1)
		encoder.DeserializeRaw(*arg2.Value, &num2)

		output := encoder.Serialize(int32(num1 << num2))

		assignOutput(&output, "i32", expr, call)
		return nil
	} else {
		return err
	}
}

func shiftRightI32 (arg1, arg2 *CXArgument, expr *CXExpression, call *CXCall) error {
	if err := checkTwoTypes("i32.bitshr", "i32", "i32", arg1, arg2); err == nil {
		var num1 uint32
		var num2 uint32
		encoder.DeserializeRaw(*arg1.Value, &num1)
		encoder.DeserializeRaw(*arg2.Value, &num2)

		output := encoder.Serialize(int32(num1 >> num2))

		assignOutput(&output, "i32", expr, call)
		return nil
	} else {
		return err
	}
}

func randI32 (min *CXArgument, max *CXArgument, expr *CXExpression, call *CXCall) error {
	if err := checkTwoTypes("i32.rand", "i32", "i32", min, max); err == nil {
		var minimum int32
		encoder.DeserializeRaw(*min.Value, &minimum)

		var maximum int32
		encoder.DeserializeRaw(*max.Value, &maximum)

		if minimum > maximum {
			return errors.New(fmt.Sprintf("i32.rand: min must be less than max (%d !< %d)", minimum, maximum))
		}

		rand.Seed(time.Now().UTC().UnixNano())
		output := encoder.SerializeAtomic(int32(rand.Intn(int(maximum - minimum)) + int(minimum)))

		assignOutput(&output, "i32", expr, call)
		return nil
	} else {
		return err
	}
}

func readI32A (arr *CXArgument, idx *CXArgument, expr *CXExpression, call *CXCall) error {
	if err := checkTwoTypes("[]i32.read", "[]i32", "i32", arr, idx); err == nil {
		var index int32
		encoder.DeserializeRaw(*idx.Value, &index)

		var size int32
		encoder.DeserializeAtomic((*arr.Value)[0:4], &size)

		if index < 0 {
			return errors.New(fmt.Sprintf("[]i32.read: negative index %d", index))
		}
		
		if index >= size {
			return errors.New(fmt.Sprintf("[]i32.read: index %d exceeds array of length %d", index, size))
		}

		var value int32
		encoder.DeserializeRaw((*arr.Value)[(index+1)*4:(index+2)*4], &value)
		output := encoder.Serialize(value)

		assignOutput(&output, "i32", expr, call)
		return nil
	} else {
		return err
	}
}

func writeI32A (arr *CXArgument, idx *CXArgument, val *CXArgument, expr *CXExpression, call *CXCall) error {
	if err := checkThreeTypes("[]i32.write", "[]i32", "i32", "i32", arr, idx, val); err == nil {
		var index int32
		encoder.DeserializeRaw(*idx.Value, &index)

		var size int32
		encoder.DeserializeAtomic((*arr.Value)[0:4], &size)

		if index < 0 {
			return errors.New(fmt.Sprintf("[]i32.write: negative index %d", index))
		}

		if index >= size {
			return errors.New(fmt.Sprintf("[]i32.write: index %d exceeds array of length %d", index, size))
		}

		// i := (int(index)+1)*4
		// for c := 0; c < 4; c++ {
		// 	(*arr.Value)[i + c] = (*val.Value)[c]
		// }

		offset := int(index) * 4 + 4
		firstChunk := make([]byte, offset)
		secondChunk := make([]byte, len(*arr.Value) - (offset + 4))

		copy(firstChunk, (*arr.Value)[:offset])
		copy(secondChunk, (*arr.Value)[offset + 4:])

		final := append(firstChunk, *val.Value...)
		final = append(final, secondChunk...)

		assignOutput(&final, "[]i32", expr, call)
		return nil
	} else {
		return err
	}
}

func lenI32A (arr *CXArgument, expr *CXExpression, call *CXCall) error {
	if err := checkType("[]i32.len", "[]i32", arr); err == nil {
		size := (*arr.Value)[:4]
		assignOutput(&size, "i32", expr, call)
		return nil
	} else {
		return err
	}
}

func ltI32 (arg1 *CXArgument, arg2 *CXArgument, expr *CXExpression, call *CXCall) error {
	if err := checkTwoTypes("i32.lt", "i32", "i32", arg1, arg2); err == nil {
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

		assignOutput(&val, "bool", expr, call)
		return nil
	} else {
		return err
	}
}

func gtI32 (arg1 *CXArgument, arg2 *CXArgument, expr *CXExpression, call *CXCall) error {
	if err := checkTwoTypes("i32.gt", "i32", "i32", arg1, arg2); err == nil {
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

		assignOutput(&val, "bool", expr, call)
		return nil
	} else {
		return err
	}
}

func eqI32 (arg1 *CXArgument, arg2 *CXArgument, expr *CXExpression, call *CXCall) error {
	if err := checkTwoTypes("i32.eq", "i32", "i32", arg1, arg2); err == nil {
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

		assignOutput(&val, "bool", expr, call)
		return nil
	} else {
		return err
	}
}

func lteqI32 (arg1 *CXArgument, arg2 *CXArgument, expr *CXExpression, call *CXCall) error {
	if err := checkTwoTypes("i32.lteq", "i32", "i32", arg1, arg2); err == nil {
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

		assignOutput(&val, "bool", expr, call)
		return nil
	} else {
		return err
	}
}

func gteqI32 (arg1 *CXArgument, arg2 *CXArgument, expr *CXExpression, call *CXCall) error {
	if err := checkTwoTypes("i32.gteq", "i32", "i32", arg1, arg2); err == nil {
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

		assignOutput(&val, "bool", expr, call)
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

		assignOutput(&sOutput, "[]i32", expr, call)
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

		//*arg1.Value = sOutput
		assignOutput(&sOutput, "[]i32", expr, call)
		return nil
	} else {
		return err
	}
}

func copyI32A (arg1 *CXArgument, arg2 *CXArgument, expr *CXExpression, call *CXCall) error {
	if err := checkTwoTypes("[]i32.copy", "[]i32", "[]i32", arg1, arg2); err == nil {
		copy(*arg1.Value, *arg2.Value)
		return nil
	} else {
		return err
	}
}

func readI32 (expr *CXExpression, call *CXCall) error {
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	text = strings.Replace(text, "\n", "", -1)
	num, err := strconv.ParseInt(text, 10, 32)
	if err != nil {
		return err
	}
	output := encoder.Serialize(num)

	assignOutput(&output, "i32", expr, call)
	return nil
}
