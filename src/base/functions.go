package base

import (
	"fmt"
	"github.com/skycoin/skycoin/src/cipher/encoder"
)

func addI32 (arg1 *CXArgument, arg2 *CXArgument) *CXArgument {
	if arg1.Typ.Name != "i32" || arg2.Typ.Name != "i32" {
		if arg1.Typ.Name != "i32" {
			panic(fmt.Sprintf("addI32: first argument is type '%s'; expected type 'i32'", arg1.Typ.Name))
		}
		panic(fmt.Sprintf("addI32: second argument is type '%s'; expected type 'i32'", arg1.Typ.Name))
	}

	var num1 int32
	var num2 int32
	encoder.DeserializeAtomic(*arg1.Value, &num1)
	encoder.DeserializeAtomic(*arg2.Value, &num2)

	output := encoder.SerializeAtomic(num1 + num2)

	return &CXArgument{Value: &output, Typ: MakeType("i32")}
}

func subI32 (arg1 *CXArgument, arg2 *CXArgument) *CXArgument {
	if arg1.Typ.Name != "i32" || arg2.Typ.Name != "i32" {
		panic("subI32: wrong argument type")
	}
	
	var num1 int32
	var num2 int32
	encoder.DeserializeAtomic(*arg1.Value, &num1)
	encoder.DeserializeAtomic(*arg2.Value, &num2)

	output := encoder.SerializeAtomic(num1 - num2)

	return &CXArgument{Value: &output, Typ: MakeType("i32")}
}

func mulI32 (arg1 *CXArgument, arg2 *CXArgument) *CXArgument {
	if arg1.Typ.Name != "i32" || arg2.Typ.Name != "i32" {
		panic("mulI32: wrong argument type")
	}

	var num1 int32
	var num2 int32
	encoder.DeserializeAtomic(*arg1.Value, &num1)
	encoder.DeserializeAtomic(*arg2.Value, &num2)

	output := encoder.SerializeAtomic(num1 * num2)

	return &CXArgument{Value: &output, Typ: MakeType("i32")}
}

func divI32 (arg1 *CXArgument, arg2 *CXArgument) *CXArgument {
	if arg1.Typ.Name != "i32" || arg2.Typ.Name != "i32" {
		panic("divI32: wrong argument type")
	}
	
	var num1 int32
	var num2 int32
	encoder.DeserializeAtomic(*arg1.Value, &num1)
	encoder.DeserializeAtomic(*arg2.Value, &num2)

	if num2 == 0 {
		panic("divI32: Division by 0")
	}

	output := encoder.SerializeAtomic(num1 / num2)

	return &CXArgument{Value: &output, Typ: MakeType("i32")}
}

func addI64 (arg1 *CXArgument, arg2 *CXArgument) *CXArgument {
	if arg1.Typ.Name != "i64" || arg2.Typ.Name != "i64" {
		panic("addI64: wrong argument type")
	}
	
	var num1 int64
	var num2 int64
	encoder.DeserializeAtomic(*arg1.Value, &num1)
	encoder.DeserializeAtomic(*arg2.Value, &num2)

	output := encoder.SerializeAtomic(num1 + num2)

	return &CXArgument{Value: &output, Typ: MakeType("i64")}
}

func subI64 (arg1 *CXArgument, arg2 *CXArgument) *CXArgument {
	if arg1.Typ.Name != "i64" || arg2.Typ.Name != "i64" {
		panic("subI64: wrong argument type")
	}
	
	var num1 int64
	var num2 int64
	encoder.DeserializeAtomic(*arg1.Value, &num1)
	encoder.DeserializeAtomic(*arg2.Value, &num2)

	output := encoder.SerializeAtomic(num1 - num2)

	return &CXArgument{Value: &output, Typ: MakeType("i64")}
}

func mulI64 (arg1 *CXArgument, arg2 *CXArgument) *CXArgument {
	if arg1.Typ.Name != "i64" || arg2.Typ.Name != "i64" {
		panic("mulI64: wrong argument type")
	}
	
	var num1 int64
	var num2 int64
	encoder.DeserializeAtomic(*arg1.Value, &num1)
	encoder.DeserializeAtomic(*arg2.Value, &num2)

	output := encoder.SerializeAtomic(num1 * num2)

	return &CXArgument{Value: &output, Typ: MakeType("i64")}
}

func divI64 (arg1 *CXArgument, arg2 *CXArgument) *CXArgument {
	if arg1.Typ.Name != "i64" || arg2.Typ.Name != "i64" {
		panic("divI64: wrong argument type")
	}

	var num1 int64
	var num2 int64
	encoder.DeserializeAtomic(*arg1.Value, &num1)
	encoder.DeserializeAtomic(*arg2.Value, &num2)

	if num2 == 0 {
		panic("divI64: Division by 0")
	}
	
	output := encoder.SerializeAtomic(num1 / num2)

	return &CXArgument{Value: &output, Typ: MakeType("i64")}
}

func readAByte (arr *CXArgument, idx *CXArgument) *CXArgument {
	if arr.Typ.Name != "[]byte" || idx.Typ.Name != "i32" {
		panic("readAByte: wrong argument type")
	}

	var index int32
	encoder.DeserializeAtomic(*idx.Value, &index)

	if index >= int32(len(*arr.Value)) {
		panic(fmt.Sprintf("readAByte: Index %d exceeds array of length %d", index, len(*arr.Value)))
	}

	output := make([]byte, 1)
	output[0] = (*arr.Value)[index]

	return &CXArgument{Value: &output, Typ: MakeType("byte")}
}

func writeAByte (arr *CXArgument, idx *CXArgument, val *CXArgument) *CXArgument {
	if arr.Typ.Name != "[]byte" || idx.Typ.Name != "i32" || val.Typ.Name != "byte" {
		panic("readAByte: wrong argument type")
	}

	var index int32
	encoder.DeserializeAtomic(*idx.Value, &index)
	
	if index >= int32(len(*arr.Value)) {
		panic(fmt.Sprintf("writeAByte: Index %d exceeds array of length %d", index, len(*arr.Value)))
	}

	(*arr.Value)[index] = (*val.Value)[0]

	return arr
}
