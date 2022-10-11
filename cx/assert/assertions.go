package assert

import (
	"fmt"

	"github.com/skycoin/cx/cx/ast"
	"github.com/skycoin/cx/cx/types"
)

func ArgOfType(arg *ast.CXArgument, t int) {
	if arg.Type != types.Code(t) {
		panic(fmt.Sprintf("Argument %s, expected type %s, got %s",
			arg.Name, types.Code(t).Name(), arg.Type.Name()))
	}
}

func ArgAtomic(arg *ast.CXArgument) {
	if arg.IsStruct() || arg.IsString() || arg.IsSlicee() {
		panic(fmt.Sprintf("Argument %s of type %s is not atomic", arg.Name, arg.Type.Name()))
	}
}

func ArgNotAtomic(arg *ast.CXArgument) {
	if !arg.IsStruct() && !arg.IsString() && !arg.IsSlicee() {
		panic(fmt.Sprintf("Argument %s of type %s is atomic", arg.Name, arg.Type.Name()))
	}
}

func ArgPointer(arg *ast.CXArgument) {
	if !arg.IsPointer() {
		panic(fmt.Sprintf("Argument %s is not a pointer", arg.Name))
	}
}

func ArgNotPointer(arg *ast.CXArgument) {
	if arg.IsPointer() {
		panic(fmt.Sprintf("Argument %s is a pointer", arg.Name))
	}
}

func ArgSlice(arg *ast.CXArgument) {
	if !arg.IsSlicee() {
		panic(fmt.Sprintf("Argument %s is not a slice", arg.Name))
	}
}

func ArgNotSlice(arg *ast.CXArgument) {
	if arg.IsSlicee() {
		panic(fmt.Sprintf("Argument %s is a slice", arg.Name))
	}
}

func ArgStruct(arg *ast.CXArgument) {
	if !arg.IsStruct() {
		panic(fmt.Sprintf("Argument %s is not a struct", arg.Name))
	}
}

func ArgNotStruct(arg *ast.CXArgument) {
	if arg.IsStruct() {
		panic(fmt.Sprintf("Argument %s is a struct", arg.Name))
	}
}
