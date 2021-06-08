package vm

import (
	"fmt"

	"github.com/skycoin/cx/object"
)

/*

	LastPoppedStackElem return Last Popped Stack object.Object.
*/
func (vm *VM) LastPoppedStackElem() object.Object {
	return vm.stack[vm.sp]
}

/*

	push pushes Object onto stack.
*/
func (vm *VM) push(o object.Object) error {

	if vm.sp >= StackSize {
		return fmt.Errorf("stack overflow")
	}

	vm.stack[vm.sp] = o
	vm.sp++

	return nil
}

/*

	pop pops Object onto stack.
*/
func (vm *VM) pop() object.Object {

	o := vm.stack[vm.sp-1]

	vm.sp--

	return o
}
