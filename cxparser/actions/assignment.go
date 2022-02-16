package actions

import (
	"os"

	"github.com/skycoin/cx/cx/ast"
	"github.com/skycoin/cx/cx/constants"
	"github.com/skycoin/cx/cx/types"
)

// assignStructLiteralFields converts a struct literal to a series of struct field assignments.
// For example, `foo = Item{x: 10, y: 20}` is converted to: `foo.x = 10; foo.y = 20;`.
func assignStructLiteralFields(prgrm *ast.CXProgram, to []ast.CXExpression, from []ast.CXExpression, name string) []ast.CXExpression {
	toCXAtomicOp, err := prgrm.GetCXAtomicOpFromExpressions(to, 0)
	if err != nil {
		panic(err)
	}

	for _, f := range from {
		if f.Type == ast.CX_LINE {
			continue
		}
		cxAtomicOp, _, _, err := prgrm.GetOperation(&f)
		if err != nil {
			panic(err)
		}

		cxAtomicOp.Outputs[0].Name = name

		if len(toCXAtomicOp.Outputs[0].Indexes) > 0 {
			cxAtomicOp.Outputs[0].Lengths = toCXAtomicOp.Outputs[0].Lengths
			cxAtomicOp.Outputs[0].Indexes = toCXAtomicOp.Outputs[0].Indexes
			cxAtomicOp.Outputs[0].DereferenceOperations = append(cxAtomicOp.Outputs[0].DereferenceOperations, constants.DEREF_ARRAY)
		}

		cxAtomicOp.Outputs[0].DereferenceOperations = append(cxAtomicOp.Outputs[0].DereferenceOperations, constants.DEREF_FIELD)
	}

	return from
}

// StructLiteralAssignment handles struct literals, e.g. `Item{x: 10, y: 20}`, and references to
// struct literals, e.g. `&Item{x: 10, y: 20}` in assignment expressions.
func StructLiteralAssignment(prgrm *ast.CXProgram, to []ast.CXExpression, from []ast.CXExpression) []ast.CXExpression {
	lastFrom := from[len(from)-1]

	lastFromAtomicOp, _, _, err := prgrm.GetOperation(&lastFrom)
	if err != nil {
		panic(err)
	}

	toCXAtomicOp, err := prgrm.GetCXAtomicOpFromExpressions(to, 0)
	if err != nil {
		panic(err)
	}

	lastFromCXLine, _ := prgrm.GetPreviousCXLine(from, len(from)-1)

	// If the last expression in `from` is declared as pointer
	// then it means the whole struct literal needs to be passed by reference.
	if !hasDeclSpec(lastFromAtomicOp.Outputs[0].GetAssignmentElement(), constants.DECL_POINTER) {
		return assignStructLiteralFields(prgrm, to, from, toCXAtomicOp.Outputs[0].Name)
	} else {
		// And we also need an auxiliary variable to point to,
		// otherwise we'd be trying to assign the fields to a nil value.
		fOut := lastFromAtomicOp.Outputs[0]
		auxName := MakeGenSym(constants.LOCAL_PREFIX)
		aux := ast.MakeArgument(auxName, lastFromCXLine.FileName, lastFromCXLine.LineNumber).AddType(fOut.Type)
		aux.DeclarationSpecifiers = append(aux.DeclarationSpecifiers, constants.DECL_POINTER)
		aux.StructType = fOut.StructType
		aux.Size = fOut.Size
		aux.TotalSize = fOut.TotalSize
		aux.PreviouslyDeclared = true
		aux.Package = lastFromAtomicOp.Package

		declExprCXLine := ast.MakeCXLineExpression(prgrm, lastFromCXLine.FileName, lastFromCXLine.LineNumber, lastFromCXLine.LineStr)
		declExpr := ast.MakeAtomicOperatorExpression(prgrm, nil)
		declExprAtomicOp, _, _, err := prgrm.GetOperation(declExpr)
		if err != nil {
			panic(err)
		}
		declExprAtomicOp.Package = lastFromAtomicOp.Package
		declExprAtomicOp.AddOutput(aux)
		from = assignStructLiteralFields(prgrm, to, from, auxName)

		assignExprCXLine := ast.MakeCXLineExpression(prgrm, lastFromCXLine.FileName, lastFromCXLine.LineNumber, lastFromCXLine.LineStr)
		assignExpr := ast.MakeAtomicOperatorExpression(prgrm, ast.Natives[constants.OP_IDENTITY])
		assignExprAtomicOp, _, _, err := prgrm.GetOperation(assignExpr)
		if err != nil {
			panic(err)
		}

		assignExprAtomicOp.Package = lastFromAtomicOp.Package
		out := ast.MakeArgument(toCXAtomicOp.Outputs[0].Name, lastFromCXLine.FileName, lastFromCXLine.LineNumber)
		out.PassBy = constants.PASSBY_REFERENCE
		out.Package = lastFromAtomicOp.Package
		assignExprAtomicOp.AddOutput(out)
		assignExprAtomicOp.AddInput(aux)

		from = append([]ast.CXExpression{*declExprCXLine, *declExpr}, from...)
		return append(from, *assignExprCXLine, *assignExpr)
	}
}

func ArrayLiteralAssignment(prgrm *ast.CXProgram, to []ast.CXExpression, from []ast.CXExpression) []ast.CXExpression {
	toCXAtomicOp, err := prgrm.GetCXAtomicOpFromExpressions(to, 0)
	if err != nil {
		panic(err)
	}

	for _, f := range from {
		cxAtomicOp, _, _, err := prgrm.GetOperation(&f)
		if err != nil {
			panic(err)
		}

		cxAtomicOp.Outputs[0].Name = toCXAtomicOp.Outputs[0].Name
		cxAtomicOp.Outputs[0].DereferenceOperations = append(cxAtomicOp.Outputs[0].DereferenceOperations, constants.DEREF_ARRAY)
	}

	return from
}

func ShortAssignment(prgrm *ast.CXProgram, expr *ast.CXExpression, exprCXLine *ast.CXExpression, to []ast.CXExpression, from []ast.CXExpression, pkg *ast.CXPackage, idx int) []ast.CXExpression {
	cxAtomicOp, _, _, err := prgrm.GetOperation(expr)
	if err != nil {
		panic(err)
	}

	toCXAtomicOp, err := prgrm.GetCXAtomicOpFromExpressions(to, 0)
	if err != nil {
		panic(err)
	}

	fromCXAtomicOp, err := prgrm.GetCXAtomicOpFromExpressions(from, idx)
	if err != nil {
		panic(err)
	}

	cxAtomicOp.AddInput(toCXAtomicOp.Outputs[0])
	cxAtomicOp.AddOutput(toCXAtomicOp.Outputs[0])
	cxAtomicOp.Package = ast.CXPackageIndex(pkg.Index)

	if fromCXAtomicOp.Operator == nil {
		cxAtomicOp.AddInput(fromCXAtomicOp.Outputs[0])
	} else {
		sym := ast.MakeArgument(MakeGenSym(constants.LOCAL_PREFIX), CurrentFile, LineNo).AddType(fromCXAtomicOp.Inputs[0].Type)
		sym.Package = ast.CXPackageIndex(pkg.Index)
		sym.PreviouslyDeclared = true
		fromCXAtomicOp.AddOutput(sym)
		cxAtomicOp.AddInput(sym)
	}

	//must check if from expression is naked previously declared variable
	if len(from) == 1 && fromCXAtomicOp.Operator == nil && len(fromCXAtomicOp.Outputs) > 0 && len(fromCXAtomicOp.Inputs) == 0 {
		return []ast.CXExpression{*exprCXLine, *expr}
	} else {
		return append(from, *exprCXLine, *expr)
	}
}

// getOutputType tries to determine what's the argument that holds the type that should be
// returned by a function call.
// This function is needed because CX has some standard library functions that return cxcore.TYPE_UNDEFINED
// arguments. In these cases, the output type depends on its input arguments' type. In the rest of
// the cases, we can simply use the function's return type.
func getOutputType(prgrm *ast.CXProgram, expr *ast.CXExpression) *ast.CXArgument {
	cxAtomicOp, _, _, err := prgrm.GetOperation(expr)
	if err != nil {
		panic(err)
	}

	if prgrm.GetCXArgFromArray(cxAtomicOp.Operator.Outputs[0]).Type != types.UNDEFINED {
		return prgrm.GetCXArgFromArray(cxAtomicOp.Operator.Outputs[0])
	}

	return cxAtomicOp.Inputs[0]
}

// Assignment handles assignment statements with different operators, like =, :=, +=, *=.
func Assignment(prgrm *ast.CXProgram, to []ast.CXExpression, assignOp string, from []ast.CXExpression) []ast.CXExpression {
	idx := len(from) - 1

	toCXAtomicOp, err := prgrm.GetCXAtomicOpFromExpressions(to, 0)
	if err != nil {
		panic(err)
	}

	fromCXAtomicOp, err := prgrm.GetCXAtomicOpFromExpressions(from, idx)
	if err != nil {
		panic(err)
	}

	// Checking if we're trying to assign stuff from a function call
	// And if that function call actually returns something. If not, throw an error.
	if fromCXAtomicOp.Operator != nil && len(fromCXAtomicOp.Operator.Outputs) == 0 {
		println(ast.CompilationError(toCXAtomicOp.Outputs[0].ArgDetails.FileName, toCXAtomicOp.Outputs[0].ArgDetails.FileLine), "trying to use an outputless operator in an assignment")
		os.Exit(constants.CX_COMPILATION_ERROR)
	}

	pkg, err := prgrm.GetCurrentPackage()
	if err != nil {
		panic(err)
	}

	var expr *ast.CXExpression

	exprCXLine := ast.MakeCXLineExpression(prgrm, CurrentFile, LineNo, LineStr)

	switch assignOp {
	case ":=":
		expr = ast.MakeAtomicOperatorExpression(prgrm, nil)
		cxAtomicOp, _, _, err := prgrm.GetOperation(expr)
		if err != nil {
			panic(err)
		}
		cxAtomicOp.Package = ast.CXPackageIndex(pkg.Index)
		var sym *ast.CXArgument

		if fromCXAtomicOp.Operator == nil {
			// then it's a literal
			sym = ast.MakeArgument(toCXAtomicOp.Outputs[0].Name, CurrentFile, LineNo).AddType(fromCXAtomicOp.Outputs[0].Type)
		} else {
			outTypeArg := getOutputType(prgrm, &from[idx])

			sym = ast.MakeArgument(toCXAtomicOp.Outputs[0].Name, CurrentFile, LineNo).AddType(outTypeArg.Type)

			if from[idx].IsArrayLiteral() {
				sym.Size = fromCXAtomicOp.Inputs[0].Size
				sym.TotalSize = fromCXAtomicOp.Inputs[0].TotalSize
				sym.Lengths = fromCXAtomicOp.Inputs[0].Lengths
			}
			if outTypeArg.IsSlice {
				// if from[idx].Operator.ProgramOutput[0].IsSlice {
				sym.Lengths = append([]types.Pointer{0}, sym.Lengths...)
				sym.DeclarationSpecifiers = append(sym.DeclarationSpecifiers, constants.DECL_SLICE)
			}

			sym.IsSlice = outTypeArg.IsSlice
			// sym.IsSlice = from[idx].Operator.ProgramOutput[0].IsSlice
		}
		sym.Package = ast.CXPackageIndex(pkg.Index)
		sym.PreviouslyDeclared = true
		sym.IsShortAssignmentDeclaration = true

		cxAtomicOp.AddOutput(sym)
		for _, toExpr := range to {
			if toExpr.Type == ast.CX_LINE {
				continue
			}
			toExprAtomicOp, _, _, err := prgrm.GetOperation(&toExpr)
			if err != nil {
				panic(err)
			}

			toExprAtomicOp.Outputs[0].PreviouslyDeclared = true
			toExprAtomicOp.Outputs[0].IsShortAssignmentDeclaration = true
		}

		to = append([]ast.CXExpression{*exprCXLine, *expr}, to...)
	case ">>=":
		expr = ast.MakeAtomicOperatorExpression(prgrm, ast.Natives[constants.OP_BITSHR])
		return ShortAssignment(prgrm, expr, exprCXLine, to, from, pkg, idx)
	case "<<=":
		expr = ast.MakeAtomicOperatorExpression(prgrm, ast.Natives[constants.OP_BITSHL])
		return ShortAssignment(prgrm, expr, exprCXLine, to, from, pkg, idx)
	case "+=":
		expr = ast.MakeAtomicOperatorExpression(prgrm, ast.Natives[constants.OP_ADD])
		return ShortAssignment(prgrm, expr, exprCXLine, to, from, pkg, idx)
	case "-=":
		expr = ast.MakeAtomicOperatorExpression(prgrm, ast.Natives[constants.OP_SUB])
		return ShortAssignment(prgrm, expr, exprCXLine, to, from, pkg, idx)
	case "*=":
		expr = ast.MakeAtomicOperatorExpression(prgrm, ast.Natives[constants.OP_MUL])
		return ShortAssignment(prgrm, expr, exprCXLine, to, from, pkg, idx)
	case "/=":
		expr = ast.MakeAtomicOperatorExpression(prgrm, ast.Natives[constants.OP_DIV])
		return ShortAssignment(prgrm, expr, exprCXLine, to, from, pkg, idx)
	case "%=":
		expr = ast.MakeAtomicOperatorExpression(prgrm, ast.Natives[constants.OP_MOD])
		return ShortAssignment(prgrm, expr, exprCXLine, to, from, pkg, idx)
	case "&=":
		expr = ast.MakeAtomicOperatorExpression(prgrm, ast.Natives[constants.OP_BITAND])
		return ShortAssignment(prgrm, expr, exprCXLine, to, from, pkg, idx)
	case "^=":
		expr = ast.MakeAtomicOperatorExpression(prgrm, ast.Natives[constants.OP_BITXOR])
		return ShortAssignment(prgrm, expr, exprCXLine, to, from, pkg, idx)
	case "|=":
		expr = ast.MakeAtomicOperatorExpression(prgrm, ast.Natives[constants.OP_BITOR])
		return ShortAssignment(prgrm, expr, exprCXLine, to, from, pkg, idx)
	}

	toLastExprAtomicOp, err := prgrm.GetCXAtomicOpFromExpressions(to, len(to)-1)
	if err != nil {
		panic(err)
	}

	if fromCXAtomicOp.Operator == nil {
		fromCXAtomicOp.Operator = ast.Natives[constants.OP_IDENTITY]
		toCXAtomicOp.Outputs[0].Size = fromCXAtomicOp.Outputs[0].Size
		toCXAtomicOp.Outputs[0].TotalSize = fromCXAtomicOp.Outputs[0].TotalSize
		toCXAtomicOp.Outputs[0].Type = fromCXAtomicOp.Outputs[0].Type
		toCXAtomicOp.Outputs[0].PointerTargetType = fromCXAtomicOp.Outputs[0].PointerTargetType
		toCXAtomicOp.Outputs[0].Lengths = fromCXAtomicOp.Outputs[0].Lengths
		toCXAtomicOp.Outputs[0].PassBy = fromCXAtomicOp.Outputs[0].PassBy
		toCXAtomicOp.Outputs[0].DoesEscape = fromCXAtomicOp.Outputs[0].DoesEscape
		// toCXAtomicOp.ProgramOutput[0].Program = prgrm

		if from[idx].IsMethodCall() {
			fromCXAtomicOp.Inputs = append(fromCXAtomicOp.Outputs, fromCXAtomicOp.Inputs...)
		} else {
			fromCXAtomicOp.Inputs = fromCXAtomicOp.Outputs
		}

		fromCXAtomicOp.Outputs = toLastExprAtomicOp.Outputs
		// from[idx].Program = prgrm

		// TODO: temporary bug fix, needs improvements
		prgrm.CXAtomicOps[from[idx].Index] = *fromCXAtomicOp

		return append(to[:len(to)-1], from...)
	} else {

		fromCXAtomicOpOperatorOutput := prgrm.GetCXArgFromArray(fromCXAtomicOp.Operator.Outputs[0])
		if fromCXAtomicOp.Operator.IsBuiltIn() {
			// only assigning as if the operator had only one output defined

			if fromCXAtomicOp.Operator.AtomicOPCode != constants.OP_IDENTITY {
				// it's a short variable declaration
				toCXAtomicOp.Outputs[0].Size = fromCXAtomicOpOperatorOutput.Size
				toCXAtomicOp.Outputs[0].Type = fromCXAtomicOpOperatorOutput.Type
				toCXAtomicOp.Outputs[0].PointerTargetType = fromCXAtomicOpOperatorOutput.PointerTargetType
				toCXAtomicOp.Outputs[0].Lengths = fromCXAtomicOpOperatorOutput.Lengths
			}

			toCXAtomicOp.Outputs[0].DoesEscape = fromCXAtomicOpOperatorOutput.DoesEscape
			toCXAtomicOp.Outputs[0].PassBy = fromCXAtomicOpOperatorOutput.PassBy
			// toCXAtomicOp.ProgramOutput[0].Program = prgrm
		} else {
			// we'll delegate multiple-value returns to the 'expression' grammar rule
			// only assigning as if the operator had only one output defined

			toCXAtomicOp.Outputs[0].Size = fromCXAtomicOpOperatorOutput.Size
			toCXAtomicOp.Outputs[0].Type = fromCXAtomicOpOperatorOutput.Type
			toCXAtomicOp.Outputs[0].PointerTargetType = fromCXAtomicOpOperatorOutput.PointerTargetType
			toCXAtomicOp.Outputs[0].Lengths = fromCXAtomicOpOperatorOutput.Lengths
			toCXAtomicOp.Outputs[0].DoesEscape = fromCXAtomicOpOperatorOutput.DoesEscape
			toCXAtomicOp.Outputs[0].PassBy = fromCXAtomicOpOperatorOutput.PassBy
			// toCXAtomicOp.ProgramOutput[0].Program = prgrm
		}

		fromCXAtomicOp.Outputs = toLastExprAtomicOp.Outputs

		// TODO: temporary bug fix, needs improvements
		prgrm.CXAtomicOps[from[idx].Index] = *fromCXAtomicOp
		return append(to[:len(to)-1], from...)
		// return append(to, from...)
	}
}
