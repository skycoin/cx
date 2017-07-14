package base

import (
	"github.com/skycoin/skycoin/src/cipher/encoder"
)

func addI32 (arg1 *cxArgument, arg2 *cxArgument) *cxArgument {
	var num1 int32
	var num2 int32
	encoder.DeserializeAtomic(*arg1.Value, &num1)
	encoder.DeserializeAtomic(*arg2.Value, &num2)

	output := encoder.SerializeAtomic(num1 + num2)

	return &cxArgument{Value: &output, Typ: MakeType("i32")}
}

func subI32 (arg1 *cxArgument, arg2 *cxArgument) *cxArgument {
	var num1 int32
	var num2 int32
	encoder.DeserializeAtomic(*arg1.Value, &num1)
	encoder.DeserializeAtomic(*arg2.Value, &num2)

	output := encoder.SerializeAtomic(num1 - num2)

	return &cxArgument{Value: &output, Typ: MakeType("i32")}
}

func mulI32 (arg1 *cxArgument, arg2 *cxArgument) *cxArgument {
	var num1 int32
	var num2 int32
	encoder.DeserializeAtomic(*arg1.Value, &num1)
	encoder.DeserializeAtomic(*arg2.Value, &num2)

	output := encoder.SerializeAtomic(num1 * num2)

	return &cxArgument{Value: &output, Typ: MakeType("i32")}
}
