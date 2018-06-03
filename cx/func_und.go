package base

import (
	// "fmt"
	"github.com/skycoin/skycoin/src/cipher/encoder"
)

func gtUnd (arg1, arg2 *CXArgument, expr *CXExpression, call *CXCall) error {
	var val []byte

	switch arg1.Type {
	case TYPE_I32:
		var num1 int32
		var num2 int32
		encoder.DeserializeAtomic(*arg1.Value, &num1)
		encoder.DeserializeAtomic(*arg2.Value, &num2)

		if num1 > num2 {
			val = encoder.Serialize(true)
		} else {
			val = encoder.Serialize(false)
		}
	case TYPE_I64:
		var num1 int64
		var num2 int64
		encoder.DeserializeRaw(*arg1.Value, &num1)
		encoder.DeserializeRaw(*arg2.Value, &num2)

		if num1 > num2 {
			val = encoder.Serialize(true)
		} else {
			val = encoder.Serialize(false)
		}
	case TYPE_F32: 
		var num1 float32
		var num2 float32
		encoder.DeserializeRaw(*arg1.Value, &num1)
		encoder.DeserializeRaw(*arg2.Value, &num2)

		if num1 > num2 {
			val = encoder.Serialize(true)
		} else {
			val = encoder.Serialize(false)
		}
	case TYPE_F64:
		var num1 float64
		var num2 float64
		encoder.DeserializeRaw(*arg1.Value, &num1)
		encoder.DeserializeRaw(*arg2.Value, &num2)

		if num1 > num2 {
			val = encoder.Serialize(true)
		} else {
			val = encoder.Serialize(false)
		}
	}
	
	

	assignOutput(0, val, "bool", expr, call)
	return nil
}

func lenUnd (arg1 *CXArgument, expr *CXExpression, call *CXCall) error {
	val := encoder.Serialize(int32(arg1.Lengths[len(arg1.Lengths) - 1]))
	assignOutput(0, val, "i32", expr, call)
	return nil
}

