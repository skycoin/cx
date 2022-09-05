package ast

type CXAtomicOperator struct {
	Inputs   *CXStruct
	Outputs  *CXStruct
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
func (op *CXAtomicOperator) AddInput(prgrm *CXProgram, typeSignatureIdx CXTypeSignatureIndex) *CXAtomicOperator {
	if op.Inputs == nil {
		op.Inputs = &CXStruct{Package: op.Package}
	}

	typeSignature := prgrm.GetCXTypeSignatureFromArray(typeSignatureIdx)
	if typeSignature.Package == -1 {
		typeSignature.Package = op.Package
	}

	// Add Package if arg has no package
	if typeSignature.Type == TYPE_CXARGUMENT_DEPRECATE {
		arg := prgrm.GetCXArgFromArray(CXArgumentIndex(typeSignature.Meta))
		if arg.Package == -1 {
			arg.Package = op.Package
		}
	}

	op.Inputs.AddField_CXAtomicOps(prgrm, typeSignatureIdx)

	return op
}

func (op *CXAtomicOperator) GetInputs(prgrm *CXProgram) []CXTypeSignatureIndex {
	if op == nil || op.Inputs == nil {
		return []CXTypeSignatureIndex{}
	}

	return op.Inputs.Fields
}

// RemoveInput ...
// func (op *CXAtomicOperator) RemoveInput() {
// 	if len(op.Inputs) > 0 {
// 		op.Inputs = op.Inputs[:len(op.Inputs)-1]
// 	}
// }

// AddOutput ...
func (op *CXAtomicOperator) AddOutput(prgrm *CXProgram, typeSignatureIdx CXTypeSignatureIndex) *CXAtomicOperator {
	if op.Outputs == nil {
		op.Outputs = &CXStruct{Package: op.Package}
	}

	typeSignature := prgrm.GetCXTypeSignatureFromArray(typeSignatureIdx)
	if typeSignature.Package == -1 {
		typeSignature.Package = op.Package
	}

	// Add Package if arg has no package
	if typeSignature.Type == TYPE_CXARGUMENT_DEPRECATE {
		arg := prgrm.GetCXArgFromArray(CXArgumentIndex(typeSignature.Meta))
		if arg.Package == -1 {
			arg.Package = op.Package
		}
	}

	op.Outputs.AddField_CXAtomicOps(prgrm, typeSignatureIdx)
	return op
}

func (op *CXAtomicOperator) GetOutputs(prgrm *CXProgram) []CXTypeSignatureIndex {
	if op == nil || op.Outputs == nil {
		return []CXTypeSignatureIndex{}
	}

	return op.Outputs.Fields
}

// RemoveOutput ...
// func (op *CXAtomicOperator) RemoveOutput() {
// 	if len(op.Outputs) > 0 {
// 		op.Outputs = op.Outputs[:len(op.Outputs)-1]
// 	}
// }

// AddLabel ...
func (op *CXAtomicOperator) AddLabel(lbl string) *CXAtomicOperator {
	op.Label = lbl
	return op
}
