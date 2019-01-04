package actions

import (
	. "github.com/skycoin/cx/cx"
)

func StructLiteralAssignment(to []*CXExpression, from []*CXExpression) []*CXExpression {
	for _, f := range from {
		f.Outputs[0].Name = to[0].Outputs[0].Name

		if len(to[0].Outputs[0].Indexes) > 0 {
			f.Outputs[0].Lengths = to[0].Outputs[0].Lengths
			f.Outputs[0].Indexes = to[0].Outputs[0].Indexes
			f.Outputs[0].DereferenceOperations = append(f.Outputs[0].DereferenceOperations, DEREF_ARRAY)
		}

		f.Outputs[0].DereferenceOperations = append(f.Outputs[0].DereferenceOperations, DEREF_FIELD)
	}

	return from
}

func ArrayLiteralAssignment(to []*CXExpression, from []*CXExpression) []*CXExpression {
	for _, f := range from {
		f.Outputs[0].Name = to[0].Outputs[0].Name
		f.Outputs[0].DereferenceOperations = append(f.Outputs[0].DereferenceOperations, DEREF_ARRAY)
	}

	return from
}

func ShortAssignment(expr *CXExpression, to []*CXExpression, from []*CXExpression, pkg *CXPackage, idx int) []*CXExpression {
	expr.AddInput(to[0].Outputs[0])
	expr.AddOutput(to[0].Outputs[0])
	expr.Package = pkg

	if from[idx].Operator == nil {
		expr.AddInput(from[idx].Outputs[0])
	} else {
		sym := MakeArgument(MakeGenSym(LOCAL_PREFIX), CurrentFile, LineNo).AddType(TypeNames[from[idx].Inputs[0].Type])
		sym.Package = pkg
		sym.PreviouslyDeclared = true
		from[idx].AddOutput(sym)
		expr.AddInput(sym)
	}

	return append(from, expr)
}

func Assignment(to []*CXExpression, assignOp string, from []*CXExpression) []*CXExpression {
	idx := len(from) - 1

	if pkg, err := PRGRM.GetCurrentPackage(); err == nil {

		var expr *CXExpression

		switch assignOp {
		case ":=":
			expr = MakeExpression(nil, CurrentFile, LineNo)
			expr.Package = pkg

			var sym *CXArgument

			if from[idx].Operator == nil {
				// then it's a literal
				sym = MakeArgument(to[0].Outputs[0].Name, CurrentFile, LineNo).AddType(TypeNames[from[idx].Outputs[0].Type])
			} else {
				sym = MakeArgument(to[0].Outputs[0].Name, CurrentFile, LineNo).AddType(TypeNames[from[idx].Inputs[0].Type])
				// sym = MakeArgument(to[0].Outputs[0].Name, CurrentFile, LineNo).AddType(TypeNames[from[idx].Operator.Outputs[0].Type])

				if from[idx].IsArrayLiteral {
					sym.Size = from[idx].Inputs[0].Size
					sym.TotalSize = from[idx].Inputs[0].TotalSize
					sym.Lengths = from[idx].Inputs[0].Lengths
				}
				if from[idx].Inputs[0].IsSlice {
					// if from[idx].Operator.Outputs[0].IsSlice {
					sym.Lengths = append([]int{0}, sym.Lengths...)
				}

				sym.IsSlice = from[idx].Inputs[0].IsSlice
				// sym.IsSlice = from[idx].Operator.Outputs[0].IsSlice
			}
			sym.Package = pkg
			sym.PreviouslyDeclared = true
			sym.IsShortDeclaration = true

			expr.AddOutput(sym)

			for _, toExpr := range to {
				toExpr.Outputs[0].PreviouslyDeclared = true
				toExpr.Outputs[0].IsShortDeclaration = true
			}

			to = append([]*CXExpression{expr}, to...)
		case ">>=":
			expr = MakeExpression(Natives[OP_UND_BITSHR], CurrentFile, LineNo)
			return ShortAssignment(expr, to, from, pkg, idx)
		case "<<=":
			expr = MakeExpression(Natives[OP_UND_BITSHL], CurrentFile, LineNo)
			return ShortAssignment(expr, to, from, pkg, idx)
		case "+=":
			expr = MakeExpression(Natives[OP_UND_ADD], CurrentFile, LineNo)
			return ShortAssignment(expr, to, from, pkg, idx)
		case "-=":
			expr = MakeExpression(Natives[OP_UND_SUB], CurrentFile, LineNo)
			return ShortAssignment(expr, to, from, pkg, idx)
		case "*=":
			expr = MakeExpression(Natives[OP_UND_MUL], CurrentFile, LineNo)
			return ShortAssignment(expr, to, from, pkg, idx)
		case "/=":
			expr = MakeExpression(Natives[OP_UND_DIV], CurrentFile, LineNo)
			return ShortAssignment(expr, to, from, pkg, idx)
		case "%=":
			expr = MakeExpression(Natives[OP_UND_MOD], CurrentFile, LineNo)
			return ShortAssignment(expr, to, from, pkg, idx)
		case "&=":
			expr = MakeExpression(Natives[OP_UND_BITAND], CurrentFile, LineNo)
			return ShortAssignment(expr, to, from, pkg, idx)
		case "^=":
			expr = MakeExpression(Natives[OP_UND_BITXOR], CurrentFile, LineNo)
			return ShortAssignment(expr, to, from, pkg, idx)
		case "|=":
			expr = MakeExpression(Natives[OP_UND_BITOR], CurrentFile, LineNo)
			return ShortAssignment(expr, to, from, pkg, idx)
		}
	}

	if from[idx].Operator == nil {
		from[idx].Operator = Natives[OP_IDENTITY]
		to[0].Outputs[0].Size = from[idx].Outputs[0].Size
		to[0].Outputs[0].Type = from[idx].Outputs[0].Type
		to[0].Outputs[0].Lengths = from[idx].Outputs[0].Lengths
		to[0].Outputs[0].PassBy = from[idx].Outputs[0].PassBy
		to[0].Outputs[0].DoesEscape = from[idx].Outputs[0].DoesEscape
		// to[0].Outputs[0].Program = PRGRM

		if from[idx].IsMethodCall {
			from[idx].Inputs = append(from[idx].Outputs, from[idx].Inputs...)
		} else {
			from[idx].Inputs = from[idx].Outputs
		}

		from[idx].Outputs = to[len(to)-1].Outputs
		// from[idx].Program = PRGRM

		return append(to[:len(to)-1], from...)
	} else {
		if from[idx].Operator.IsNative {
			// only assigning as if the operator had only one output defined

			if from[idx].Operator.OpCode != OP_IDENTITY {
				// it's a short variable declaration
				to[0].Outputs[0].Size = Natives[from[idx].Operator.OpCode].Outputs[0].Size
				to[0].Outputs[0].Type = from[idx].Operator.Outputs[0].Type
				to[0].Outputs[0].Lengths = from[idx].Operator.Outputs[0].Lengths
			}

			// to[0].Outputs[0].Type = from[idx].Operator.Outputs[0].Type
			// to[0].Outputs[0].Lengths = from[idx].Operator.Outputs[0].Lengths
			// to[0].Outputs[0].Size = Natives[from[idx].Operator.OpCode].Outputs[0].Size

			to[0].Outputs[0].DoesEscape = from[idx].Operator.Outputs[0].DoesEscape
			to[0].Outputs[0].PassBy = from[idx].Operator.Outputs[0].PassBy
			// to[0].Outputs[0].Program = PRGRM
		} else {
			// we'll delegate multiple-value returns to the 'expression' grammar rule
			// only assigning as if the operator had only one output defined

			to[0].Outputs[0].Size = from[idx].Operator.Outputs[0].Size
			to[0].Outputs[0].Type = from[idx].Operator.Outputs[0].Type
			to[0].Outputs[0].Lengths = from[idx].Operator.Outputs[0].Lengths
			to[0].Outputs[0].DoesEscape = from[idx].Operator.Outputs[0].DoesEscape
			to[0].Outputs[0].PassBy = from[idx].Operator.Outputs[0].PassBy
			// to[0].Outputs[0].Program = PRGRM
		}

		from[idx].Outputs = to[len(to)-1].Outputs
		// from[idx].Program = to[len(to) - 1].Program

		return append(to[:len(to)-1], from...)
		// return append(to, from...)
	}
}
