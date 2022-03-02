package ast

type CXAtomicOperator struct {
	Inputs   []CXArgumentIndex
	Outputs  []CXArgumentIndex
	Operator CXFunctionIndex

	Function CXFunctionIndex
	Package  CXPackageIndex
	Label    string

	// used for jmp statements
	ThenLines int
	ElseLines int
}

// ----------------------------------------------------------------
//                             `CXAtomicOperator` Getters

func (op *CXAtomicOperator) GetOperatorName(prgrm *CXProgram) string {
	opOperator := prgrm.GetFunctionFromArray(op.Operator)
	if opOperator.IsBuiltIn() {
		return OpNames[opOperator.AtomicOPCode]
	}
	return opOperator.Name

}

// ----------------------------------------------------------------
//                     `CXAtomicOperator` Member handling

// AddInput ...
func (op *CXAtomicOperator) AddInput(prgrm *CXProgram, paramIdx CXArgumentIndex) *CXAtomicOperator {
	if prgrm.GetCXArgFromArray(paramIdx).Package == -1 {
		prgrm.GetCXArgFromArray(paramIdx).Package = op.Package
	}
	op.Inputs = append(op.Inputs, paramIdx)

	return op
}

// RemoveInput ...
func (op *CXAtomicOperator) RemoveInput() {
	if len(op.Inputs) > 0 {
		op.Inputs = op.Inputs[:len(op.Inputs)-1]
	}
}

// AddOutput ...
func (op *CXAtomicOperator) AddOutput(prgrm *CXProgram, paramIdx CXArgumentIndex) *CXAtomicOperator {
	if prgrm.GetCXArgFromArray(paramIdx).Package == -1 {
		prgrm.GetCXArgFromArray(paramIdx).Package = op.Package
	}

	op.Outputs = append(op.Outputs, paramIdx)

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
