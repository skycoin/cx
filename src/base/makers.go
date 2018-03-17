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
		Packages: make([]*CXPackage, 0),
		CallStack: make([]CXCall, callStackSize, callStackSize),
		Stacks: make([]CXStack, 1, 1),
		Heap: MakeHeap(initialHeapSize),
	}
	
	newPrgrm.Stacks[0] = MakeStack(stackSize)
	newPrgrm.Stacks[0].Program = newPrgrm
	
	return newPrgrm
}

func MakePackage (name string) *CXPackage {
	return &CXPackage{
		Name: name,
		Globals: make([]*CXArgument, 0, 10),
		Imports: make([]*CXPackage, 0),
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
		MemoryType: MEM_HEAP,
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
		TotalSize: size,
		MemoryType: MEM_STACK,
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

func MakeNative (opCode int, inputs []int, outputs []int) *CXFunction {
	fn := &CXFunction{
		OpCode: opCode,
		IsNative: true,
	}

	offset := 0
	for _, typCode := range inputs {
		inp := MakeParameter("", typCode)
		inp.Offset = offset
		offset += inp.Size
		fn.Inputs = append(fn.Inputs, inp)
	}
	for _, typCode := range outputs {
		fn.Outputs = append(fn.Outputs, MakeParameter("", typCode))
		out := MakeParameter("", typCode)
		out.Offset = offset
		offset += out.Size
	}

	return fn
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

func MakeHeap (size int) CXHeap {
	return CXHeap{
		Heap: make([]byte, size, size),
		// HeapPointer: 0,
		HeapPointer: NULL_HEAP_ADDRESS_OFFSET,
	}
}

func MakeCall (op *CXFunction, ret *CXCall, mod *CXPackage, cxt *CXProgram) CXCall {
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
