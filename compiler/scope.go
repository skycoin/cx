package compiler

import "github.com/skycoin/cx/code"

/*
CompilationScope reprsents compilation scope of instruction.
*/
type CompilationScope struct {
	instructions        code.Instructions
	lastInstruction     EmittedInstruction
	previousInstruction EmittedInstruction
}

/*
EmittedInstruction used for evaluation
*/

type EmittedInstruction struct {
	Opcode   code.Opcode
	Position int
}

/*
enterScope used for scope evaluation
*/
func (c *Compiler) enterScope() {

	scope := CompilationScope{
		instructions:        code.Instructions{},
		lastInstruction:     EmittedInstruction{},
		previousInstruction: EmittedInstruction{},
	}
	c.scopes = append(c.scopes, scope)
	c.scopeIndex++

	c.symbolTable = NewEnclosedSymbolTable(c.symbolTable)
}

/*
leaveScope used for scope evaluation
*/
func (c *Compiler) leaveScope() code.Instructions {
	instructions := c.currentInstructions()

	c.scopes = c.scopes[:len(c.scopes)-1]
	c.scopeIndex--

	c.symbolTable = c.symbolTable.Outer

	return instructions
}
