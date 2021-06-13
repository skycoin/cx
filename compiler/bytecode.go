package compiler

import (
	"github.com/skycoin/cx/code"
	"github.com/skycoin/cx/object"
)

/*Bytecode represents bytecode */
type Bytecode struct {
	Instructions code.Instructions
	Constants    []object.Object
}
