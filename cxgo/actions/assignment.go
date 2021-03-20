package actions

import (
	"os"

	"github.com/skycoin/cx/cx"
)

// assignStructLiteralFields converts a struct literal to a series of struct field assignments.
// For example, `foo = Item{x: 10, y: 20}` is converted to: `foo.x = 10; foo.y = 20;`.
func assignStructLiteralFields(to []*cxcore.CXExpression, from []*cxcore.CXExpression, name string) []*cxcore.CXExpression {
	for _, f := range from {
		f.Outputs[0].Name = name

		if len(to[0].Outputs[0].Indexes) > 0 {
			f.Outputs[0].Lengths = to[0].Outputs[0].Lengths
			f.Outputs[0].Indexes = to[0].Outputs[0].Indexes
			f.Outputs[0].DereferenceOperations = append(f.Outputs[0].DereferenceOperations, cxcore.DEREF_ARRAY)
		}

		f.Outputs[0].DereferenceOperations = append(f.Outputs[0].DereferenceOperations, cxcore.DEREF_FIELD)
	}

	return from
}

// StructLiteralAssignment handles struct literals, e.g. `Item{x: 10, y: 20}`, and references to
// struct literals, e.g. `&Item{x: 10, y: 20}` in assignment expressions.
func StructLiteralAssignment(to []*cxcore.CXExpression, from []*cxcore.CXExpression) []*cxcore.CXExpression {
	lastFrom := from[len(from)-1]
	// If the last expression in `from` is declared as pointer
	// then it means the whole struct literal needs to be passed by reference.
	if !hasDeclSpec(cxcore.GetAssignmentElement(lastFrom.Outputs[0]), cxcore.DECL_POINTER) {
		return assignStructLiteralFields(to, from, to[0].Outputs[0].Name)
	} else {
		// And we also need an auxiliary variable to point to,
		// otherwise we'd be trying to assign the fields to a nil value.
		fOut := lastFrom.Outputs[0]
		auxName := cxcore.MakeGenSym(cxcore.LOCAL_PREFIX)
		aux := cxcore.MakeArgument(auxName, lastFrom.FileName, lastFrom.FileLine).AddType(cxcore.TypeNames[fOut.Type])
		aux.DeclarationSpecifiers = append(aux.DeclarationSpecifiers, cxcore.DECL_POINTER)
		aux.CustomType = fOut.CustomType
		aux.Size = fOut.Size
		aux.TotalSize = fOut.TotalSize
		aux.PreviouslyDeclared = true
		aux.Package = lastFrom.Package

		declExpr := cxcore.MakeExpression(nil, lastFrom.FileName, lastFrom.FileLine)
		declExpr.Package = lastFrom.Package
		declExpr.AddOutput(aux)

		from = assignStructLiteralFields(to, from, auxName)

		assignExpr := cxcore.MakeExpression(cxcore.Natives[cxcore.OP_IDENTITY], lastFrom.FileName, lastFrom.FileLine)
		assignExpr.Package = lastFrom.Package
		out := cxcore.MakeArgument(to[0].Outputs[0].Name, lastFrom.FileName, lastFrom.FileLine)
		out.PassBy = cxcore.PASSBY_REFERENCE
		out.Package = lastFrom.Package
		assignExpr.AddOutput(out)
		assignExpr.AddInput(aux)

		from = append([]*cxcore.CXExpression{declExpr}, from...)
		return append(from, assignExpr)
	}
}

func ArrayLiteralAssignment(to []*cxcore.CXExpression, from []*cxcore.CXExpression) []*cxcore.CXExpression {
	for _, f := range from {
		f.Outputs[0].Name = to[0].Outputs[0].Name
		f.Outputs[0].DereferenceOperations = append(f.Outputs[0].DereferenceOperations, cxcore.DEREF_ARRAY)
	}

	return from
}

func ShortAssignment(expr *cxcore.CXExpression, to []*cxcore.CXExpression, from []*cxcore.CXExpression, pkg *cxcore.CXPackage, idx int) []*cxcore.CXExpression {
	expr.AddInput(to[0].Outputs[0])
	expr.AddOutput(to[0].Outputs[0])
	expr.Package = pkg

	if from[idx].Operator == nil {
		expr.AddInput(from[idx].Outputs[0])
	} else {
		sym := cxcore.MakeArgument(cxcore.MakeGenSym(cxcore.LOCAL_PREFIX), CurrentFile, LineNo).AddType(cxcore.TypeNames[from[idx].Inputs[0].Type])
		sym.Package = pkg
		sym.PreviouslyDeclared = true
		from[idx].AddOutput(sym)
		expr.AddInput(sym)
	}

	//must check if from expression is naked previously declared variable
	if len(from) == 1 && from[0].Operator == nil && len(from[0].Outputs) > 0 && len(from[0].Inputs) == 0 {
		return []*cxcore.CXExpression{expr}
	} else {
		return append(from, expr)
	}
}

// getOutputType tries to determine what's the argument that holds the type that should be
// returned by a function call.
// This function is needed because CX has some standard library functions that return cxcore.TYPE_UNDEFINED
// arguments. In these cases, the output type depends on its input arguments' type. In the rest of
// the cases, we can simply use the function's return type.
func getOutputType(expr *cxcore.CXExpression) *cxcore.CXArgument {
	if expr.Operator.Outputs[0].Type != cxcore.TYPE_UNDEFINED {
		return expr.Operator.Outputs[0]
	}

	return expr.Inputs[0]
}

// Assignment handles assignment statements with different operators, like =, :=, +=, *=.
func Assignment(to []*cxcore.CXExpression, assignOp string, from []*cxcore.CXExpression) []*cxcore.CXExpression {
	idx := len(from) - 1

	// Checking if we're trying to assign stuff from a function call
	// And if that function call actually returns something. If not, throw an error.
	if from[idx].Operator != nil && len(from[idx].Operator.Outputs) == 0 {
		println(cxcore.CompilationError(to[0].Outputs[0].FileName, to[0].Outputs[0].FileLine), "trying to use an outputless operator in an assignment")
		os.Exit(cxcore.CX_COMPILATION_ERROR)
	}

	pkg, err := PRGRM.GetCurrentPackage()
	if err != nil {
		panic(err)
	}

	var expr *cxcore.CXExpression

	switch assignOp {
	case ":=":
		expr = cxcore.MakeExpression(nil, CurrentFile, LineNo)
		expr.Package = pkg

		var sym *cxcore.CXArgument

		if from[idx].Operator == nil {
			// then it's a literal
			sym = cxcore.MakeArgument(to[0].Outputs[0].Name, CurrentFile, LineNo).AddType(cxcore.TypeNames[from[idx].Outputs[0].Type])
		} else {
			outTypeArg := getOutputType(from[idx])

			sym = cxcore.MakeArgument(to[0].Outputs[0].Name, CurrentFile, LineNo).AddType(cxcore.TypeNames[outTypeArg.Type])

			if from[idx].IsArrayLiteral {
				sym.Size = from[idx].Inputs[0].Size
				sym.TotalSize = from[idx].Inputs[0].TotalSize
				sym.Lengths = from[idx].Inputs[0].Lengths
			}
			if outTypeArg.IsSlice {
				// if from[idx].Operator.Outputs[0].IsSlice {
				sym.Lengths = append([]int{0}, sym.Lengths...)
				sym.DeclarationSpecifiers = append(sym.DeclarationSpecifiers, cxcore.DECL_SLICE)
			}

			sym.IsSlice = outTypeArg.IsSlice
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

		to = append([]*cxcore.CXExpression{expr}, to...)
	case ">>=":
		expr = cxcore.MakeExpression(cxcore.Natives[cxcore.OP_UND_BITSHR], CurrentFile, LineNo)
		return ShortAssignment(expr, to, from, pkg, idx)
	case "<<=":
		expr = cxcore.MakeExpression(cxcore.Natives[cxcore.OP_UND_BITSHL], CurrentFile, LineNo)
		return ShortAssignment(expr, to, from, pkg, idx)
	case "+=":
		expr = cxcore.MakeExpression(cxcore.Natives[cxcore.OP_UND_ADD], CurrentFile, LineNo)
		return ShortAssignment(expr, to, from, pkg, idx)
	case "-=":
		expr = cxcore.MakeExpression(cxcore.Natives[cxcore.OP_UND_SUB], CurrentFile, LineNo)
		return ShortAssignment(expr, to, from, pkg, idx)
	case "*=":
		expr = cxcore.MakeExpression(cxcore.Natives[cxcore.OP_UND_MUL], CurrentFile, LineNo)
		return ShortAssignment(expr, to, from, pkg, idx)
	case "/=":
		expr = cxcore.MakeExpression(cxcore.Natives[cxcore.OP_UND_DIV], CurrentFile, LineNo)
		return ShortAssignment(expr, to, from, pkg, idx)
	case "%=":
		expr = cxcore.MakeExpression(cxcore.Natives[cxcore.OP_UND_MOD], CurrentFile, LineNo)
		return ShortAssignment(expr, to, from, pkg, idx)
	case "&=":
		expr = cxcore.MakeExpression(cxcore.Natives[cxcore.OP_UND_BITAND], CurrentFile, LineNo)
		return ShortAssignment(expr, to, from, pkg, idx)
	case "^=":
		expr = cxcore.MakeExpression(cxcore.Natives[cxcore.OP_UND_BITXOR], CurrentFile, LineNo)
		return ShortAssignment(expr, to, from, pkg, idx)
	case "|=":
		expr = cxcore.MakeExpression(cxcore.Natives[cxcore.OP_UND_BITOR], CurrentFile, LineNo)
		return ShortAssignment(expr, to, from, pkg, idx)
	}

	if from[idx].Operator == nil {
		from[idx].Operator = cxcore.Natives[cxcore.OP_IDENTITY]
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

			if from[idx].Operator.OpCode != cxcore.OP_IDENTITY {
				// it's a short variable declaration
				to[0].Outputs[0].Size = cxcore.Natives[from[idx].Operator.OpCode].Outputs[0].Size
				to[0].Outputs[0].Type = from[idx].Operator.Outputs[0].Type
				to[0].Outputs[0].Lengths = from[idx].Operator.Outputs[0].Lengths
			}

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
