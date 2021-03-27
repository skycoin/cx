package assert

import (
	"fmt"
	cxcore "github.com/skycoin/cx/cx"
	"github.com/skycoin/cx/cx/constants"
)

func ArgOfType(arg *cxcore.CXArgument, t int) {
	if arg.Type != t {
		panic(fmt.Sprintf("Argument %s, expected type %s, got %s",
			arg.Name, constants.TypeNames[t], constants.TypeNames[arg.Type]))
	}
}

func ArgAtomic(arg *cxcore.CXArgument) {
	if arg.IsStruct || arg.IsPointer || arg.IsArray || arg.IsSlice {
		panic(fmt.Sprintf("Argument %s of type %s is not atomic", arg.Name, constants.TypeNames[arg.Type]))
	}
}

func ArgNotAtomic(arg *cxcore.CXArgument) {
	if !arg.IsStruct && !arg.IsPointer && !arg.IsArray && !arg.IsSlice {
		panic(fmt.Sprintf("Argument %s of type %s is atomic", arg.Name, constants.TypeNames[arg.Type]))
	}
}

func ArgPointer(arg *cxcore.CXArgument) {
	if !arg.IsPointer {
		panic(fmt.Sprintf("Argument %s is not a pointer", arg.Name))
	}
}

func ArgNotPointer(arg *cxcore.CXArgument) {
	if arg.IsPointer {
		panic(fmt.Sprintf("Argument %s is a pointer", arg.Name))
	}
}

func ArgSlice(arg *cxcore.CXArgument) {
	if !arg.IsSlice {
		panic(fmt.Sprintf("Argument %s is not a slice", arg.Name))
	}
}

func ArgNotSlice(arg *cxcore.CXArgument) {
	if arg.IsSlice {
		panic(fmt.Sprintf("Argument %s is a slice", arg.Name))
	}
}

func ArgStruct(arg *cxcore.CXArgument) {
	if !arg.IsStruct {
		panic(fmt.Sprintf("Argument %s is not a struct", arg.Name))
	}
}

func ArgNotStruct(arg *cxcore.CXArgument) {
	if arg.IsStruct {
		panic(fmt.Sprintf("Argument %s is a struct", arg.Name))
	}
}

func ArgArray(arg *cxcore.CXArgument) {
	if !arg.IsArray {
		panic(fmt.Sprintf("Argument %s is not a array", arg.Name))
	}
}

func ArgNotArray(arg *cxcore.CXArgument) {
	if arg.IsArray {
		panic(fmt.Sprintf("Argument %s is a array", arg.Name))
	}
}
