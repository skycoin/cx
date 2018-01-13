package base

import (
	"fmt"
	//"github.com/satori/go.uuid"
	"github.com/skycoin/skycoin/src/cipher/encoder"
)

var HeapOffset int
var genSymCounter int = 0
func MakeGenSym (name string) string {
	gensym := fmt.Sprintf("%s_%d", name, genSymCounter)
	genSymCounter++
	
	return gensym
}

func MakeProgram (callStackSize int, stackSize int, initialHeapSize int) *CXProgram {
	newPrgrm := &CXProgram{
		Modules: make([]*CXModule, 0),
		CallStack: make([]CXCall, callStackSize, callStackSize),
		Stacks: make([]CXStack, 1, 1),
		Heaps: make([]Heap, 1, 1),
	}
	newPrgrm.Stacks[0] = MakeStack(stackSize)
	newPrgrm.Heaps[0] = Heap(make([]byte, initialHeapSize))
	return newPrgrm
}

func MakeModule (name string) *CXModule {
	return &CXModule{
		Name: name,
		Globals: make([]*CXArgument, 0, 10),
		Imports: make([]*CXModule, 0),
		Functions: make([]*CXFunction, 0, 10),
		Structs: make([]*CXStruct, 0),
	}
}

func MakeGlobal (name string, typ int) *CXArgument {
	size := GetArgSize(typ)
	global := &CXArgument{
		Name: name,
		Type: typ,
		Size: size,
		MemoryType: HEAP,
		Offset: HeapOffset,
	}
	HeapOffset += size
	return global
}

func MakeField (name string, typ int) *CXArgument {
	return &CXArgument{Name: name, Type: typ}
}

// Used only for native types
func MakeDefaultValue (typName string) *[]byte {
	var zeroVal []byte
	switch typName {
	case "byte": zeroVal = make([]byte, 1, 1)
	case "i64", "f64": zeroVal = make([]byte, 8, 8)
	default: zeroVal = make([]byte, 4, 4)
	}
	return &zeroVal
}

func MakeStruct (name string) *CXStruct {
	return &CXStruct{Name: name}
}

func MakeParameter (name string, typ int) *CXArgument {
	size := GetArgSize(typ)
	return &CXArgument{
		Name: name,
		Type: typ,
		Size: size,
		MemoryType: STACK,
		// this will be added in AddInput & AddOutput
		// the parent function knows how many parameters it has
		// Offset: offset,
	}
}

func MakeExpression (op *CXFunction) *CXExpression {
	return &CXExpression{Operator: op}
}

func MakeArgument (typ int) *CXArgument {
	return &CXArgument{Type: typ}
}

func MakeFunction (name string) *CXFunction {
	return &CXFunction{Name: name}
}

func MakeValue (value string) *[]byte {
	byts := encoder.Serialize(value)
	return &byts
}

func MakeStack (size int) CXStack {
	return CXStack{
		Stack: make([]byte, size, size),
		StackPointer: 0,
	}
}

func MakeCall (op *CXFunction, ret *CXCall, mod *CXModule, cxt *CXProgram) CXCall {
	return CXCall{
		Operator: op,
		Line: 0,
		FramePointer: 0}
}

func MakeAffordance (desc string, action func()) *CXAffordance {
	return &CXAffordance{
		Description: desc,
		Action: action,
	}
}

func MakeIdentityOpName (typeName string) string {
	switch typeName {
	case "str":
		return "str.id"
	case "bool":
		return "bool.id"
	case "byte":
		return "byte.id"
	case "i32":
		return "i32.id"
	case "i64":
		return "i64.id"
	case "f32":
		return "f32.id"
	case "f64":
		return "f64.id"
	case "[]bool":
		return "[]bool.id"
	case "[]byte":
		return "[]byte.id"
	case "[]str":
		return "[]str.id"
	case "[]i32":
		return "[]i32.id"
	case "[]i64":
		return "[]i64.id"
	case "[]f32":
		return "[]f32.id"
	case "[]f64":
		return "[]f64.id"
	default:
		return ""
	}
}
