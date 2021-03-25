package assert

import (
	"fmt"
	cxcore "github.com/skycoin/cx/cx"
)

func AssertOfType(arg *cxcore.CXArgument, t int) {
	if arg.Type != t {
		panic(fmt.Sprintf("Argument %s, expected type %s, got %s",
			arg.Name, cxcore.TypeNames[t], cxcore.TypeNames[arg.Type]))
	}
}

func AssertAtomic(arg *cxcore.CXArgument) {
	if arg.IsStruct || arg.IsPointer || arg.IsArray || arg.IsSlice {
		panic(fmt.Sprintf("Argument %s of type %s is not atomic", arg.Name, cxcore.TypeNames[arg.Type]))
	}
}

func AssertNotAtomic(arg *cxcore.CXArgument) {
	if !arg.IsStruct && !arg.IsPointer && !arg.IsArray && !arg.IsSlice {
		panic(fmt.Sprintf("Argument %s of type %s is atomic", arg.Name, cxcore.TypeNames[arg.Type]))
	}
}

func AssertPointer(arg *cxcore.CXArgument) {
	if !arg.IsPointer {
		panic(fmt.Sprintf("Argument %s is not a pointer", arg.Name))
	}
}

func AssertNotPointer(arg *cxcore.CXArgument) {
	if arg.IsPointer {
		panic(fmt.Sprintf("Argument %s is a pointer", arg.Name))
	}
}

func AssertSlice(arg *cxcore.CXArgument) {
	if !arg.IsSlice {
		panic(fmt.Sprintf("Argument %s is not a slice", arg.Name))
	}
}

func AssertNotSlice(arg *cxcore.CXArgument) {
	if arg.IsSlice {
		panic(fmt.Sprintf("Argument %s is a slice", arg.Name))
	}
}

func AssertStruct(arg *cxcore.CXArgument) {
	if !arg.IsStruct {
		panic(fmt.Sprintf("Argument %s is not a struct", arg.Name))
	}
}

func AssertNotStruct(arg *cxcore.CXArgument) {
	if arg.IsStruct {
		panic(fmt.Sprintf("Argument %s is a struct", arg.Name))
	}
}

func AssertArray(arg *cxcore.CXArgument) {
	if !arg.IsArray {
		panic(fmt.Sprintf("Argument %s is not a array", arg.Name))
	}
}

func AssertNotArray(arg *cxcore.CXArgument) {
	if arg.IsArray {
		panic(fmt.Sprintf("Argument %s is a array", arg.Name))
	}
}
