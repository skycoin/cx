package compiler

import (
	"github.com/skycoin/cx/code"
	"github.com/skycoin/cx/object"
)

/*
Compiler represents Compiler for cx programming language.
*/
type Compiler struct {

	/*
		constant represents array of  Object as constant pool
	*/
	constants []object.Object

	symbolTable *SymbolTable

	scopes []CompilationScope

	scopeIndex int
}

/*
New returns Compiler Object for cx programming language.
*/
func New() *Compiler {

	mainScope := CompilationScope{
		instructions:        code.Instructions{},
		lastInstruction:     EmittedInstruction{},
		previousInstruction: EmittedInstruction{},
	}

	symbolTable := NewSymbolTable()

	for i, v := range object.Builtins {
		symbolTable.DefineBuiltin(i, v.Name)
	}

	c := &Compiler{
		constants:   []object.Object{},
		symbolTable: symbolTable,
		scopes:      []CompilationScope{mainScope},
		scopeIndex:  0,
	}

	return c
}

/*NewWithState return new state compiler object */
func NewWithState(s *SymbolTable, constants []object.Object) *Compiler {
	compiler := New()
	compiler.symbolTable = s
	compiler.constants = constants
	return compiler
}
