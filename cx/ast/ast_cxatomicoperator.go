package ast

type CXAtomicOperator struct {
	Inputs  []*CXArgument
	Outputs []*CXArgument
	Opcode  int
}
