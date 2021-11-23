package ast

type CXAtomicOperator struct {
	Inputs  []*CXArgument
	Outputs []*CXArgument
	Opcode  int

	// used for jmp statements
	ThenLines int
	ElseLines int
}
