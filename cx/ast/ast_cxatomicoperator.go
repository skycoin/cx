package ast

type CXAtomicOperator struct {
	Inputs   []*CXArgument
	Outputs  []*CXArgument
	Operator *CXFunction

	Function *CXFunction
	Package  *CXPackage
	Label    string

	// used for jmp statements
	ThenLines int
	ElseLines int
}

// ----------------------------------------------------------------
//                             `CXAtomicOperator` Getters

func (op *CXAtomicOperator) GetOperatorName() string {
	if op.Operator.IsBuiltIn() {
		return OpNames[op.Operator.AtomicOPCode]
	}
	return op.Operator.Name

}

// ----------------------------------------------------------------
//                     `CXAtomicOperator` Member handling

// AddInput ...
func (op *CXAtomicOperator) AddInput(param *CXArgument) *CXAtomicOperator {
	// param.Package = op.Package
	op.Inputs = append(op.Inputs, param)
	if param.Package == nil {
		param.Package = op.Package
	}
	return op
}

// RemoveInput ...
func (op *CXAtomicOperator) RemoveInput() {
	if len(op.Inputs) > 0 {
		op.Inputs = op.Inputs[:len(op.Inputs)-1]
	}
}

// AddOutput ...
func (op *CXAtomicOperator) AddOutput(param *CXArgument) *CXAtomicOperator {
	// param.Package = op.Package
	op.Outputs = append(op.Outputs, param)
	if param.Package == nil {
		param.Package = op.Package
	}
	return op
}

// RemoveOutput ...
func (op *CXAtomicOperator) RemoveOutput() {
	if len(op.Outputs) > 0 {
		op.Outputs = op.Outputs[:len(op.Outputs)-1]
	}
}

// AddLabel ...
func (op *CXAtomicOperator) AddLabel(lbl string) *CXAtomicOperator {
	op.Label = lbl
	return op
}
