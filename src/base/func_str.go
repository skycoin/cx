package base

import (
	"fmt"
	"bufio"
	"os"
	"strings"
	"errors"
	
	"github.com/skycoin/skycoin/src/cipher/encoder"
)

// "character arrays" (not string arrays)

func ltStr (arg1 *CXArgument, arg2 *CXArgument, expr *CXExpression, call *CXCall) error {
	if err := checkTwoTypes("str.lt", "str", "str", arg1, arg2); err == nil {
		var str1 string
		var str2 string
		encoder.DeserializeRaw(*arg1.Value, &str1)
		encoder.DeserializeRaw(*arg2.Value, &str2)

		var val []byte

		if str1 < str2 {
			val = encoder.Serialize(int32(1))
		} else {
			val = encoder.Serialize(int32(0))
		}

		assignOutput(0, val, "bool", expr, call)
		return nil
	} else {
		return err
	}
}

func gtStr (arg1 *CXArgument, arg2 *CXArgument, expr *CXExpression, call *CXCall) error {
	if err := checkTwoTypes("str.gt", "str", "str", arg1, arg2); err == nil {
		var str1 string
		var str2 string
		encoder.DeserializeRaw(*arg1.Value, &str1)
		encoder.DeserializeRaw(*arg2.Value, &str2)

		var val []byte

		if str1 > str2 {
			val = encoder.Serialize(int32(1))
		} else {
			val = encoder.Serialize(int32(0))
		}

		assignOutput(0, val, "bool", expr, call)
		return nil
	} else {
		return err
	}
}

func eqStr (arg1 *CXArgument, arg2 *CXArgument, expr *CXExpression, call *CXCall) error {
	if err := checkTwoTypes("str.eq", "str", "str", arg1, arg2); err == nil {
		var str1 string
		var str2 string
		encoder.DeserializeRaw(*arg1.Value, &str1)
		encoder.DeserializeRaw(*arg2.Value, &str2)

		var val []byte

		if str1 == str2 {
			val = encoder.Serialize(int32(1))
		} else {
			val = encoder.Serialize(int32(0))
		}

		assignOutput(0, val, "bool", expr, call)
		return nil
	} else {
		return err
	}
}

func lteqStr (arg1 *CXArgument, arg2 *CXArgument, expr *CXExpression, call *CXCall) error {
	if err := checkTwoTypes("str.lteq", "str", "str", arg1, arg2); err == nil {
		var str1 string
		var str2 string
		encoder.DeserializeRaw(*arg1.Value, &str1)
		encoder.DeserializeRaw(*arg2.Value, &str2)

		var val []byte

		if str1 <= str2 {
			val = encoder.Serialize(int32(1))
		} else {
			val = encoder.Serialize(int32(0))
		}

		assignOutput(0, val, "bool", expr, call)
		return nil
	} else {
		return err
	}
}

func gteqStr (arg1 *CXArgument, arg2 *CXArgument, expr *CXExpression, call *CXCall) error {
	if err := checkTwoTypes("str.gteq", "str", "str", arg1, arg2); err == nil {
		var str1 string
		var str2 string
		encoder.DeserializeRaw(*arg1.Value, &str1)
		encoder.DeserializeRaw(*arg2.Value, &str2)

		var val []byte

		if str1 >= str2 {
			val = encoder.Serialize(int32(1))
		} else {
			val = encoder.Serialize(int32(0))
		}

		assignOutput(0, val, "bool", expr, call)
		return nil
	} else {
		return err
	}
}

func concatStr (arg1 *CXArgument, arg2 *CXArgument, expr *CXExpression, call *CXCall) error {
	if err := checkTwoTypes("str.concat", "str", "str", arg1, arg2); err == nil {
		var str1 string
		var str2 string
		encoder.DeserializeRaw(*arg1.Value, &str1)
		encoder.DeserializeRaw(*arg2.Value, &str2)

		output := str1 + str2
		sOutput := encoder.Serialize(output)

		for _, def := range call.State {
			if def.Name == expr.OutputNames[0].Name {
				def.Value = &sOutput
				return nil
			}
		}
		call.State = append(call.State, MakeDefinition(expr.OutputNames[0].Name, &sOutput, "str"))
		assignOutput(0, sOutput, "str", expr, call)
		return nil
	} else {
		return err
	}
}

// string arrays

func readStrA (arr *CXArgument, idx *CXArgument, expr *CXExpression, call *CXCall) error {
	if err := checkTwoTypes("[]str.read", "[]str", "i32", arr, idx); err == nil {
		var index int32
		encoder.DeserializeRaw(*idx.Value, &index)

		var size int32
		encoder.DeserializeAtomic((*arr.Value)[0:4], &size)

		if index < 0 {
			return errors.New(fmt.Sprintf("[]str.read: negative index %d", index))
		}
		
		if index >= size {
			return errors.New(fmt.Sprintf("[]str.read: index %d exceeds array of length %d", index, size))
		}

		noSize := (*arr.Value)[4:]

		var offset int32
		for c := 0; c < int(index); c++ {
			var strSize int32
			encoder.DeserializeRaw(noSize[offset:offset+4], &strSize)
			offset += strSize + 4
		}

		sStrSize := noSize[offset:offset + 4]
		var strSize int32
		encoder.DeserializeRaw(sStrSize, &strSize)
		
		var value string
		encoder.DeserializeRaw(noSize[offset:offset+strSize+4], &value)
		output := encoder.Serialize(value)

		assignOutput(0, output, "str", expr, call)
		return nil
	} else {
		return err
	}
}

func writeStrA (arr *CXArgument, idx *CXArgument, val *CXArgument, expr *CXExpression, call *CXCall) error {
	if err := checkThreeTypes("[]str.write", "[]str", "i32", "str", arr, idx, val); err == nil {
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

		var array []string
		encoder.DeserializeRaw(*arr.Value, &array)

		var value string
		encoder.DeserializeRaw(*val.Value, &value)

		array[index] = value
		//*arr.Value = encoder.Serialize(array)
		sOutput := encoder.Serialize(array)
		assignOutput(0, sOutput, "[]str", expr, call)
		return nil
	} else {
		return err
	}
}

func lenStr (arr *CXArgument, expr *CXExpression, call *CXCall) error {
	if err := checkType("str.len", "str", arr); err == nil {
		size := (*arr.Value)[:4]
		assignOutput(0, size, "i32", expr, call)
		return nil
	} else {
		return err
	}
}

func lenStrA (arr *CXArgument, expr *CXExpression, call *CXCall) error {
	if err := checkType("[]str.len", "[]str", arr); err == nil {
		size := (*arr.Value)[:4]
		assignOutput(0, size, "i32", expr, call)
		return nil
	} else {
		return err
	}
}

func concatStrA (arg1 *CXArgument, arg2 *CXArgument, expr *CXExpression, call *CXCall) error {
	if err := checkTwoTypes("[]str.concat", "[]str", "[]str", arg1, arg2); err == nil {
		var slice1 []string
		var slice2 []string
		encoder.DeserializeRaw(*arg1.Value, &slice1)
		encoder.DeserializeRaw(*arg2.Value, &slice2)

		output := append(slice1, slice2...)
		sOutput := encoder.Serialize(output)

		assignOutput(0, sOutput, "[]str", expr, call)
		return nil
	} else {
		return err
	}
}

func appendStrA (arg1 *CXArgument, arg2 *CXArgument, expr *CXExpression, call *CXCall) error {
	if err := checkTwoTypes("[]str.append", "[]str", "str", arg1, arg2); err == nil {
		var slice []string
		var literal string
		encoder.DeserializeRaw(*arg1.Value, &slice)
		encoder.DeserializeRaw(*arg2.Value, &literal)

		output := append(slice, literal)
		sOutput := encoder.Serialize(output)

		//*arg1.Value = sOutput
		assignOutput(0, sOutput, "[]str", expr, call)
		return nil
	} else {
		return err
	}
}

func copyStrA (arg1 *CXArgument, arg2 *CXArgument, expr *CXExpression, call *CXCall) error {
	if err := checkTwoTypes("[]str.copy", "[]str", "[]str", arg1, arg2); err == nil {
		var slice1 []string
		var slice2 []string
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

// read string from standard input

func readStr (expr *CXExpression, call *CXCall) error {
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	text = strings.Replace(text, "\n", "", -1)
	output := encoder.Serialize(text)

	assignOutput(0, output, "str", expr, call)
	return nil
}
