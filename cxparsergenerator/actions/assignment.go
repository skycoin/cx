package actions

import (
	"os"

	"github.com/skycoin/cx/cx/ast"
	"github.com/skycoin/cx/cx/constants"
)

// assignStructLiteralFields converts a struct literal to a series of struct field assignments.
// For example, `foo = Item{x: 10, y: 20}` is converted to: `foo.x = 10; foo.y = 20;`.
func assignStructLiteralFields(to []*ast.CXExpression, from []*ast.CXExpression, name string) []*ast.CXExpression {
	for _, f := range from {
		f.Outputs[0].ArgDetails.Name = name

		if len(to[0].Outputs[0].Indexes) > 0 {
			f.Outputs[0].Lengths = to[0].Outputs[0].Lengths
			f.Outputs[0].Indexes = to[0].Outputs[0].Indexes
			f.Outputs[0].DereferenceOperations = append(f.Outputs[0].DereferenceOperations, constants.DEREF_ARRAY)
		}

		f.Outputs[0].DereferenceOperations = append(f.Outputs[0].DereferenceOperations, constants.DEREF_FIELD)
	}

	return from
}

// StructLiteralAssignment handles struct literals, e.g. `Item{x: 10, y: 20}`, and references to
// struct literals, e.g. `&Item{x: 10, y: 20}` in assignment expressions.
func StructLiteralAssignment(to []*ast.CXExpression, from []*ast.CXExpression) []*ast.CXExpression {
	lastFrom := from[len(from)-1]
	// If the last expression in `from` is declared as pointer
	// then it means the whole struct literal needs to be passed by reference.
	if !hasDeclSpec(ast.GetAssignmentElement(lastFrom.Outputs[0]), constants.DECL_POINTER) {
		return assignStructLiteralFields(to, from, to[0].Outputs[0].ArgDetails.Name)
	} else {
		// And we also need an auxiliary variable to point to,
		// otherwise we'd be trying to assign the fields to a nil value.
		fOut := lastFrom.Outputs[0]
		auxName := MakeGenSym(constants.LOCAL_PREFIX)
		aux := ast.MakeArgument(auxName, lastFrom.FileName, lastFrom.FileLine).AddType(constants.TypeNames[fOut.Type])
		aux.DeclarationSpecifiers = append(aux.DeclarationSpecifiers, constants.DECL_POINTER)
		aux.CustomType = fOut.CustomType
		aux.Size = fOut.Size
		aux.TotalSize = fOut.TotalSize
		aux.PreviouslyDeclared = true
		aux.ArgDetails.Package = lastFrom.Package

		declExpr := ast.MakeExpression(nil, lastFrom.FileName, lastFrom.FileLine)
		declExpr.Package = lastFrom.Package
		declExpr.AddOutput(aux)

		from = assignStructLiteralFields(to, from, auxName)

		assignExpr := ast.MakeExpression(ast.Natives[constants.OP_IDENTITY], lastFrom.FileName, lastFrom.FileLine)
		assignExpr.Package = lastFrom.Package
		out := ast.MakeArgument(to[0].Outputs[0].ArgDetails.Name, lastFrom.FileName, lastFrom.FileLine)
		out.PassBy = constants.PASSBY_REFERENCE
		out.ArgDetails.Package = lastFrom.Package
		assignExpr.AddOutput(out)
		assignExpr.AddInput(aux)

		from = append([]*ast.CXExpression{declExpr}, from...)
		return append(from, assignExpr)
	}
}

func ArrayLiteralAssignment(to []*ast.CXExpression, from []*ast.CXExpression) []*ast.CXExpression {
	for _, f := range from {
		f.Outputs[0].ArgDetails.Name = to[0].Outputs[0].ArgDetails.Name
		f.Outputs[0].DereferenceOperations = append(f.Outputs[0].DereferenceOperations, constants.DEREF_ARRAY)
	}

	return from
}

func ShortAssignment(expr *ast.CXExpression, to []*ast.CXExpression, from []*ast.CXExpression, pkg *ast.CXPackage, idx int) []*ast.CXExpression {
	expr.AddInput(to[0].Outputs[0])
	expr.AddOutput(to[0].Outputs[0])
	expr.Package = pkg

	if from[idx].Operator == nil {
		expr.AddInput(from[idx].Outputs[0])
	} else {
		sym := ast.MakeArgument(MakeGenSym(constants.LOCAL_PREFIX), CurrentFile, LineNo).AddType(constants.TypeNames[from[idx].Inputs[0].Type])
		sym.ArgDetails.Package = pkg
		sym.PreviouslyDeclared = true
		from[idx].AddOutput(sym)
		expr.AddInput(sym)
	}

	//must check if from expression is naked previously declared variable
	if len(from) == 1 && from[0].Operator == nil && len(from[0].Outputs) > 0 && len(from[0].Inputs) == 0 {
		return []*ast.CXExpression{expr}
	} else {
		return append(from, expr)
	}
}

// getOutputType tries to determine what's the argument that holds the type that should be
// returned by a function call.
// This function is needed because CX has some standard library functions that return cxcore.TYPE_UNDEFINED
// arguments. In these cases, the output type depends on its input arguments' type. In the rest of
// the cases, we can simply use the function's return type.
func getOutputType(expr *ast.CXExpression) *ast.CXArgument {
	if expr.Operator.Outputs[0].Type != constants.TYPE_UNDEFINED {
		return expr.Operator.Outputs[0]
	}

	return expr.Inputs[0]
}

// Assignment handles assignment statements with different operators, like =, :=, +=, *=.
func Assignment(to []*ast.CXExpression, assignOp string, from []*ast.CXExpression) []*ast.CXExpression {
	idx := len(from) - 1

	// Checking if we're trying to assign stuff from a function call
	// And if that function call actually returns something. If not, throw an error.
	if from[idx].Operator != nil && len(from[idx].Operator.Outputs) == 0 {
		println(ast.CompilationError(to[0].Outputs[0].ArgDetails.FileName, to[0].Outputs[0].ArgDetails.FileLine), "trying to use an outputless operator in an assignment")
		os.Exit(constants.CX_COMPILATION_ERROR)
	}

	pkg, err := AST.GetCurrentPackage()
	if err != nil {
		panic(err)
	}

	var expr *ast.CXExpression

	switch assignOp {
	case ":=":
		expr = ast.MakeExpression(nil, CurrentFile, LineNo)
		expr.Package = pkg

		var sym *ast.CXArgument

		if from[idx].Operator == nil {
			// then it's a literal
			sym = ast.MakeArgument(to[0].Outputs[0].ArgDetails.Name, CurrentFile, LineNo).AddType(constants.TypeNames[from[idx].Outputs[0].Type])
		} else {
			outTypeArg := getOutputType(from[idx])

			sym = ast.MakeArgument(to[0].Outputs[0].ArgDetails.Name, CurrentFile, LineNo).AddType(constants.TypeNames[outTypeArg.Type])

			if from[idx].IsArrayLiteral() {
				sym.Size = from[idx].Inputs[0].Size
				sym.TotalSize = from[idx].Inputs[0].TotalSize
				sym.Lengths = from[idx].Inputs[0].Lengths
			}
			if outTypeArg.IsSlice {
				// if from[idx].Operator.ProgramOutput[0].IsSlice {
				sym.Lengths = append([]int{0}, sym.Lengths...)
				sym.DeclarationSpecifiers = append(sym.DeclarationSpecifiers, constants.DECL_SLICE)
			}

			sym.IsSlice = outTypeArg.IsSlice
			// sym.IsSlice = from[idx].Operator.ProgramOutput[0].IsSlice
		}
		sym.ArgDetails.Package = pkg
		sym.PreviouslyDeclared = true
		sym.IsShortAssignmentDeclaration = true

		expr.AddOutput(sym)

		for _, toExpr := range to {
			toExpr.Outputs[0].PreviouslyDeclared = true
			toExpr.Outputs[0].IsShortAssignmentDeclaration = true
		}

		to = append([]*ast.CXExpression{expr}, to...)
	case ">>=":
		expr = ast.MakeExpression(ast.Natives[constants.OP_BITSHR], CurrentFile, LineNo)
		return ShortAssignment(expr, to, from, pkg, idx)
	case "<<=":
		expr = ast.MakeExpression(ast.Natives[constants.OP_BITSHL], CurrentFile, LineNo)
		return ShortAssignment(expr, to, from, pkg, idx)
	case "+=":
		expr = ast.MakeExpression(ast.Natives[constants.OP_ADD], CurrentFile, LineNo)
		return ShortAssignment(expr, to, from, pkg, idx)
	case "-=":
		expr = ast.MakeExpression(ast.Natives[constants.OP_SUB], CurrentFile, LineNo)
		return ShortAssignment(expr, to, from, pkg, idx)
	case "*=":
		expr = ast.MakeExpression(ast.Natives[constants.OP_MUL], CurrentFile, LineNo)
		return ShortAssignment(expr, to, from, pkg, idx)
	case "/=":
		expr = ast.MakeExpression(ast.Natives[constants.OP_DIV], CurrentFile, LineNo)
		return ShortAssignment(expr, to, from, pkg, idx)
	case "%=":
		expr = ast.MakeExpression(ast.Natives[constants.OP_MOD], CurrentFile, LineNo)
		return ShortAssignment(expr, to, from, pkg, idx)
	case "&=":
		expr = ast.MakeExpression(ast.Natives[constants.OP_BITAND], CurrentFile, LineNo)
		return ShortAssignment(expr, to, from, pkg, idx)
	case "^=":
		expr = ast.MakeExpression(ast.Natives[constants.OP_BITXOR], CurrentFile, LineNo)
		return ShortAssignment(expr, to, from, pkg, idx)
	case "|=":
		expr = ast.MakeExpression(ast.Natives[constants.OP_BITOR], CurrentFile, LineNo)
		return ShortAssignment(expr, to, from, pkg, idx)
	}

	if from[idx].Operator == nil {
		from[idx].Operator = ast.Natives[constants.OP_IDENTITY]
		to[0].Outputs[0].Size = from[idx].Outputs[0].Size
		to[0].Outputs[0].Type = from[idx].Outputs[0].Type
		to[0].Outputs[0].Lengths = from[idx].Outputs[0].Lengths
		to[0].Outputs[0].PassBy = from[idx].Outputs[0].PassBy
		to[0].Outputs[0].DoesEscape = from[idx].Outputs[0].DoesEscape
		// to[0].ProgramOutput[0].Program = AST

		if from[idx].IsMethodCall() {
			from[idx].Inputs = append(from[idx].Outputs, from[idx].Inputs...)
		} else {
			from[idx].Inputs = from[idx].Outputs
		}

		from[idx].Outputs = to[len(to)-1].Outputs
		// from[idx].Program = AST

		return append(to[:len(to)-1], from...)
	} else {
		if from[idx].Operator.IsBuiltin {
			// only assigning as if the operator had only one output defined

			if from[idx].Operator.OpCode != constants.OP_IDENTITY {
				// it's a short variable declaration
				to[0].Outputs[0].Size = ast.Natives[from[idx].Operator.OpCode].Outputs[0].Size
				to[0].Outputs[0].Type = from[idx].Operator.Outputs[0].Type
				to[0].Outputs[0].Lengths = from[idx].Operator.Outputs[0].Lengths
			}

			to[0].Outputs[0].DoesEscape = from[idx].Operator.Outputs[0].DoesEscape
			to[0].Outputs[0].PassBy = from[idx].Operator.Outputs[0].PassBy
			// to[0].ProgramOutput[0].Program = AST
		} else {
			// we'll delegate multiple-value returns to the 'expression' grammar rule
			// only assigning as if the operator had only one output defined

			to[0].Outputs[0].Size = from[idx].Operator.Outputs[0].Size
			to[0].Outputs[0].Type = from[idx].Operator.Outputs[0].Type
			to[0].Outputs[0].Lengths = from[idx].Operator.Outputs[0].Lengths
			to[0].Outputs[0].DoesEscape = from[idx].Operator.Outputs[0].DoesEscape
			to[0].Outputs[0].PassBy = from[idx].Operator.Outputs[0].PassBy
			// to[0].ProgramOutput[0].Program = AST
		}

		from[idx].Outputs = to[len(to)-1].Outputs
		// from[idx].Program = to[len(to) - 1].Program

		return append(to[:len(to)-1], from...)
		// return append(to, from...)
	}
}
