package vm

import (
	"github.com/skycoin/cx/compiler"
	"github.com/skycoin/cx/object"
)

type VM struct {
	constants []object.Object

	stack []object.Object
	sp    int // Always points to the next value. Top of stack is stack[sp-1]

	globals []object.Object

	frames      []*Frame
	framesIndex int
}

func New(bytecode *compiler.Bytecode) *VM {

	mainFn := &object.CompiledFunction{Instructions: bytecode.Instructions}

	mainClosure := &object.Closure{Fn: mainFn}

	mainFrame := NewFrame(mainClosure, 0)

	frames := make([]*Frame, MaxFrames)

	frames[0] = mainFrame

	return &VM{
		constants: bytecode.Constants,

		stack: make([]object.Object, StackSize),
		sp:    0,

		globals: make([]object.Object, GlobalsSize),

		frames:      frames,
		framesIndex: 1,
	}
}
