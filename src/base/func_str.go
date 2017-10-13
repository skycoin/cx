package base

import (
	"github.com/skycoin/skycoin/src/cipher/encoder"
)

func ltStr (arg1 *CXArgument, arg2 *CXArgument, expr *CXExpression, call *CXCall) error {
	if err := checkTwoTypes("ltStr", "str", "str", arg1, arg2); err == nil {
		str1 := string(*arg1.Value)
		str2 := string(*arg2.Value)

		var val []byte

		if str1 < str2 {
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

func gtStr (arg1 *CXArgument, arg2 *CXArgument, expr *CXExpression, call *CXCall) error {
	if err := checkTwoTypes("gtStr", "str", "str", arg1, arg2); err == nil {
		str1 := string(*arg1.Value)
		str2 := string(*arg2.Value)

		var val []byte

		if str1 > str2 {
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

func eqStr (arg1 *CXArgument, arg2 *CXArgument, expr *CXExpression, call *CXCall) error {
	if err := checkTwoTypes("eqStr", "str", "str", arg1, arg2); err == nil {
		str1 := string(*arg1.Value)
		str2 := string(*arg2.Value)

		var val []byte

		if str1 == str2 {
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

func lteqStr (arg1 *CXArgument, arg2 *CXArgument, expr *CXExpression, call *CXCall) error {
	if err := checkTwoTypes("lteqStr", "str", "str", arg1, arg2); err == nil {
		str1 := string(*arg1.Value)
		str2 := string(*arg2.Value)

		var val []byte

		if str1 <= str2 {
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

func gteqStr (arg1 *CXArgument, arg2 *CXArgument, expr *CXExpression, call *CXCall) error {
	if err := checkTwoTypes("gteqStr", "str", "str", arg1, arg2); err == nil {
		str1 := string(*arg1.Value)
		str2 := string(*arg2.Value)

		var val []byte

		if str1 >= str2 {
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

func concatStr (arg1 *CXArgument, arg2 *CXArgument, expr *CXExpression, call *CXCall) error {
	if err := checkTwoTypes("str.concat", "str", "str", arg1, arg2); err == nil {
		// var slice1 []int32
		// var slice2 []int32
		// encoder.DeserializeRaw(*arg1.Value, &slice1)
		// encoder.DeserializeRaw(*arg2.Value, &slice2)

		output := append(*arg1.Value, *arg2.Value...)
		sOutput := encoder.Serialize(output)

		for _, def := range call.State {
			if def.Name == expr.OutputNames[0].Name {
				def.Value = &sOutput
				return nil
			}
		}
		call.State = append(call.State, MakeDefinition(expr.OutputNames[0].Name, &sOutput, "str"))
		assignOutput(&sOutput, "str", expr, call)
		return nil
	} else {
		return err
	}
}
