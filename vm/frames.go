package vm

import (
	"github.com/skycoin/cx/code"
	"github.com/skycoin/cx/object"
)

/*
	Frame represents frame of function on stack of frames.
*/
type Frame struct {
	cl          *object.Closure
	ip          int
	basePointer int
}

/*
	NewFrame returns Frame Object.
*/

func NewFrame(cl *object.Closure, basePointer int) *Frame {

	f := &Frame{
		cl:          cl,
		ip:          -1,
		basePointer: basePointer,
	}

	return f
}

/*
	Instructions returns code.Instructions of current frame.
*/
func (f *Frame) Instructions() code.Instructions {
	return f.cl.Fn.Instructions
}

/*
	currentFrame returns currentFrame of vm.
*/
func (vm *VM) currentFrame() *Frame {
	return vm.frames[vm.framesIndex-1]
}

/*
	currentFrame returns currentFrame of vm.
*/
func (vm *VM) pushFrame(f *Frame) {
	vm.frames[vm.framesIndex] = f
	vm.framesIndex++
}

/*
	currentFrame returns currentFrame of vm.
*/
func (vm *VM) popFrame() *Frame {
	vm.framesIndex--
	return vm.frames[vm.framesIndex]
}
