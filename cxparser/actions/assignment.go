package actions

import (
	"os"

	"github.com/skycoin/cx/cx/ast"
	"github.com/skycoin/cx/cx/constants"
	"github.com/skycoin/cx/cx/types"
)

// assignStructLiteralFields converts a struct literal to a series of struct field assignments.
// For example, `foo = Item{x: 10, y: 20}` is converted to: `foo.x = 10; foo.y = 20;`.
//
// Input arguments description:
// 	prgrm - a CXProgram that contains all the data and array of the program.
// 	toExprs - toExprs are the array of expressions that contains the data needed
// 			  to construct the series of struct field assignments.
// 	fromExprs - fromExprs are the array of expressions that will contain the
// 			    series of struct field assignments.
// 	structLiteralName - name of the struct, in the example above this is "foo".
func assignStructLiteralFields(prgrm *ast.CXProgram, toExprs []ast.CXExpression, fromExprs []ast.CXExpression, structLiteralName string) []ast.CXExpression {
	toCXAtomicOp, err := prgrm.GetCXAtomicOpFromExpressions(toExprs, 0)
	if err != nil {
		panic(err)
	}

	for _, expr := range fromExprs {
		if expr.Type == ast.CX_LINE {
			continue
		}
		cxAtomicOp, err := prgrm.GetCXAtomicOp(expr.Index)
		if err != nil {
			panic(err)
		}

		cxAtomicOpOutputIdx := cxAtomicOp.Outputs[0]
		prgrm.CXArgs[cxAtomicOpOutputIdx].Name = structLiteralName

		toCXAtomicOpOutput := prgrm.GetCXArgFromArray(toCXAtomicOp.Outputs[0])
		if len(toCXAtomicOpOutput.Indexes) > 0 {
			prgrm.CXArgs[cxAtomicOpOutputIdx].Lengths = toCXAtomicOpOutput.Lengths
			prgrm.CXArgs[cxAtomicOpOutputIdx].Indexes = toCXAtomicOpOutput.Indexes
			prgrm.CXArgs[cxAtomicOpOutputIdx].DereferenceOperations = append(prgrm.CXArgs[cxAtomicOpOutputIdx].DereferenceOperations, constants.DEREF_ARRAY)
		}

		prgrm.CXArgs[cxAtomicOpOutputIdx].DereferenceOperations = append(prgrm.CXArgs[cxAtomicOpOutputIdx].DereferenceOperations, constants.DEREF_FIELD)
	}

	return fromExprs
}

// StructLiteralAssignment handles struct literals, e.g. `Item{x: 10, y: 20}`, and references to
// struct literals, e.g. `&Item{x: 10, y: 20}` in assignment expressions.
//
// Input arguments description:
// 	prgrm - a CXProgram that contains all the data and array of the program.
// 	toExprs - toExprs are the array of expressions that contains the data needed
// 			  to construct the series of struct field assignments.
// 	fromExprs - fromExprs are the array of expressions that will contain the
// 				series of struct field assignments.
func StructLiteralAssignment(prgrm *ast.CXProgram, toExprs []ast.CXExpression, fromExprs []ast.CXExpression) []ast.CXExpression {
	lastFromExpr := fromExprs[len(fromExprs)-1]

	lastFromAtomicOp, _, _, err := prgrm.GetOperation(&lastFromExpr)
	if err != nil {
		panic(err)
	}

	toCXAtomicOp, err := prgrm.GetCXAtomicOpFromExpressions(toExprs, 0)
	if err != nil {
		panic(err)
	}

	lastFromCXLine, _ := prgrm.GetPreviousCXLine(fromExprs, len(fromExprs)-1)

	// If the last expression in `fromExprs` is declared as pointer
	// then it means the whole struct literal needs to be passed by reference.
	if !hasDeclSpec(prgrm.GetCXArgFromArray(lastFromAtomicOp.Outputs[0]).GetAssignmentElement(prgrm), constants.DECL_POINTER) {
		return assignStructLiteralFields(prgrm, toExprs, fromExprs, prgrm.GetCXArgFromArray(toCXAtomicOp.Outputs[0]).Name)
	} else {
		// And we also need an auxiliary variable to point to,
		// otherwise we'd be trying to assign the fields to a nil value.
		outField := prgrm.GetCXArgFromArray(lastFromAtomicOp.Outputs[0])
		auxName := MakeGenSym(constants.LOCAL_PREFIX)
		aux := ast.MakeArgument(auxName, lastFromCXLine.FileName, lastFromCXLine.LineNumber)
		aux.SetType(outField.Type)
		aux.DeclarationSpecifiers = append(aux.DeclarationSpecifiers, constants.DECL_POINTER)
		aux.StructType = outField.StructType
		aux.Size = outField.Size
		aux.TotalSize = outField.TotalSize
		aux.PreviouslyDeclared = true
		aux.Package = lastFromAtomicOp.Package
		auxIdx := prgrm.AddCXArgInArray(aux)

		declExprCXLine := ast.MakeCXLineExpression(prgrm, lastFromCXLine.FileName, lastFromCXLine.LineNumber, lastFromCXLine.LineStr)
		declExpr := ast.MakeAtomicOperatorExpression(prgrm, nil)
		declExprAtomicOpIdx := declExpr.Index

		prgrm.CXAtomicOps[declExprAtomicOpIdx].Package = lastFromAtomicOp.Package
		prgrm.CXAtomicOps[declExprAtomicOpIdx].AddOutput(prgrm, auxIdx)

		fromExprs = assignStructLiteralFields(prgrm, toExprs, fromExprs, auxName)

		assignExprCXLine := ast.MakeCXLineExpression(prgrm, lastFromCXLine.FileName, lastFromCXLine.LineNumber, lastFromCXLine.LineStr)
		assignExpr := ast.MakeAtomicOperatorExpression(prgrm, ast.Natives[constants.OP_IDENTITY])
		assignExprAtomicOpIdx := assignExpr.Index

		prgrm.CXAtomicOps[assignExprAtomicOpIdx].Package = lastFromAtomicOp.Package
		out := ast.MakeArgument(prgrm.GetCXArgFromArray(toCXAtomicOp.Outputs[0]).Name, lastFromCXLine.FileName, lastFromCXLine.LineNumber)
		out.PassBy = constants.PASSBY_REFERENCE
		out.Package = lastFromAtomicOp.Package
		outIdx := prgrm.AddCXArgInArray(out)

		prgrm.CXAtomicOps[assignExprAtomicOpIdx].AddOutput(prgrm, outIdx)
		prgrm.CXAtomicOps[assignExprAtomicOpIdx].AddInput(prgrm, auxIdx)

		fromExprs = append([]ast.CXExpression{*declExprCXLine, *declExpr}, fromExprs...)
		return append(fromExprs, *assignExprCXLine, *assignExpr)
	}
}

// ArrayLiteralAssignment handles array literals.
//
// Input arguments description:
// 	prgrm - a CXProgram that contains all the data and array of the program.
// 	toExprs - toExprs are the array of expressions that contains the data needed
// 			  to construct the series of array literals.
// 	fromExprs - fromExprs are the array of expressions that will contain the
// 				series of array literal assignments.
func ArrayLiteralAssignment(prgrm *ast.CXProgram, toExprs []ast.CXExpression, fromExprs []ast.CXExpression) []ast.CXExpression {
	toCXAtomicOp, err := prgrm.GetCXAtomicOpFromExpressions(toExprs, 0)
	if err != nil {
		panic(err)
	}

	for _, expr := range fromExprs {
		cxAtomicOp, err := prgrm.GetCXAtomicOp(expr.Index)
		if err != nil {
			panic(err)
		}

		cxAtomicOpOutputIdx := cxAtomicOp.Outputs[0]
		prgrm.CXArgs[cxAtomicOpOutputIdx].Name = prgrm.GetCXArgFromArray(toCXAtomicOp.Outputs[0]).Name
		prgrm.CXArgs[cxAtomicOpOutputIdx].DereferenceOperations = append(prgrm.CXArgs[cxAtomicOpOutputIdx].DereferenceOperations, constants.DEREF_ARRAY)
	}

	return fromExprs
}

// ShortAssignment handles short assignments for ">>=","<<=",
// "+=","-=","*=","/=","%=","&=","^=", and "|=" operators.
//
// Input arguments description:
// 	prgrm - a CXProgram that contains all the data and array of the program.
// 	expr - the expression for the short assignment.
// 	exprCXLine - the cx line or the line string of the short assignment expression.
// 	toExprs - Contains the output cx arg to be added to the expression.
// 	fromExprs - Contains the output cx arg to be added to the expression.
// 	pkg - the package the expression belongs.
func ShortAssignment(prgrm *ast.CXProgram, expr *ast.CXExpression, exprCXLine *ast.CXExpression, toExprs []ast.CXExpression, fromExprs []ast.CXExpression, pkg *ast.CXPackage) []ast.CXExpression {
	cxAtomicOp, err := prgrm.GetCXAtomicOp(expr.Index)
	if err != nil {
		panic(err)
	}

	toCXAtomicOp, err := prgrm.GetCXAtomicOpFromExpressions(toExprs, 0)
	if err != nil {
		panic(err)
	}

	fromCXAtomicOp, err := prgrm.GetCXAtomicOpFromExpressions(fromExprs, len(fromExprs)-1)
	if err != nil {
		panic(err)
	}
	fromCXAtomicOpOperator := prgrm.GetFunctionFromArray(fromCXAtomicOp.Operator)

	cxAtomicOp.AddInput(prgrm, toCXAtomicOp.Outputs[0])
	cxAtomicOp.AddOutput(prgrm, toCXAtomicOp.Outputs[0])

	cxAtomicOp.Package = ast.CXPackageIndex(pkg.Index)

	if fromCXAtomicOpOperator == nil {
		cxAtomicOp.AddInput(prgrm, fromCXAtomicOp.Outputs[0])
	} else {
		sym := ast.MakeArgument(MakeGenSym(constants.LOCAL_PREFIX), CurrentFile, LineNo).SetType(prgrm.GetCXArgFromArray(fromCXAtomicOp.Inputs[0]).Type)
		sym.Package = ast.CXPackageIndex(pkg.Index)
		sym.PreviouslyDeclared = true
		symIdx := prgrm.AddCXArgInArray(sym)
		fromCXAtomicOp.AddOutput(prgrm, symIdx)

		cxAtomicOp.AddInput(prgrm, symIdx)
	}

	//must check if from expression is naked previously declared variable
	if len(fromExprs) == 1 && fromCXAtomicOpOperator == nil && len(fromCXAtomicOp.Outputs) > 0 && len(fromCXAtomicOp.Inputs) == 0 {
		return []ast.CXExpression{*exprCXLine, *expr}
	} else {
		return append(fromExprs, *exprCXLine, *expr)
	}
}

// getOutputType tries to determine what's the argument that holds the type that should be
// returned by a function call.
// This function is needed because CX has some standard library functions that return cxcore.TYPE_UNDEFINED
// arguments. In these cases, the output type depends on its input arguments' type. In the rest of
// the cases, we can simply use the function's return type.
func getOutputType(prgrm *ast.CXProgram, expr *ast.CXExpression) *ast.CXArgument {
	cxAtomicOp, err := prgrm.GetCXAtomicOp(expr.Index)
	if err != nil {
		panic(err)
	}
	cxAtomicOpOperator := prgrm.GetFunctionFromArray(cxAtomicOp.Operator)

	if prgrm.GetCXArgFromArray(cxAtomicOpOperator.Outputs[0]).Type != types.UNDEFINED {
		return prgrm.GetCXArgFromArray(cxAtomicOpOperator.Outputs[0])
	}

	return prgrm.GetCXArgFromArray(cxAtomicOp.Inputs[0])
}

// Assignment handles assignment statements with different operators,
// like =, :=, +=, *=.
//
// Input arguments description:
// 	prgrm - a CXProgram that contains all the data and array of the program.
// 	toExprs, fromExprs - array of expressions where the assingment
// 			  		  	 expression will be added.
// 	assignOp - the assignment operator, "=", ":=", ">>=","<<=",
// 			   "+=","-=","*=","/=","%=","&=","^=", and "|=".
func Assignment(prgrm *ast.CXProgram, toExprs []ast.CXExpression, assignOp string, fromExprs []ast.CXExpression) []ast.CXExpression {
	idx := len(fromExprs) - 1

	toCXAtomicOp, err := prgrm.GetCXAtomicOpFromExpressions(toExprs, 0)
	if err != nil {
		panic(err)
	}

	fromCXAtomicOpIdx := fromExprs[idx].Index
	fromCXAtomicOpOperator := prgrm.GetFunctionFromArray(prgrm.CXAtomicOps[fromCXAtomicOpIdx].Operator)

	// Checking if we're trying to assign stuff from a function call
	// And if that function call actually returns something. If not, throw an error.
	if fromCXAtomicOpOperator != nil && len(fromCXAtomicOpOperator.Outputs) == 0 {
		println(ast.CompilationError(prgrm.GetCXArgFromArray(toCXAtomicOp.Outputs[0]).ArgDetails.FileName, prgrm.GetCXArgFromArray(toCXAtomicOp.Outputs[0]).ArgDetails.FileLine), "trying to use an outputless operator in an assignment")
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
		cxAtomicOp, err := prgrm.GetCXAtomicOp(expr.Index)
		if err != nil {
			panic(err)
		}
		cxAtomicOp.Package = ast.CXPackageIndex(pkg.Index)
		var sym *ast.CXArgument

		if fromCXAtomicOpOperator == nil {
			// then it's a literal
			sym = ast.MakeArgument(prgrm.GetCXArgFromArray(toCXAtomicOp.Outputs[0]).Name, CurrentFile, LineNo).SetType(prgrm.GetCXArgFromArray(prgrm.CXAtomicOps[fromCXAtomicOpIdx].Outputs[0]).Type)
		} else {
			outTypeArg := getOutputType(prgrm, &fromExprs[idx])

			sym = ast.MakeArgument(prgrm.GetCXArgFromArray(toCXAtomicOp.Outputs[0]).Name, CurrentFile, LineNo).SetType(outTypeArg.Type)

			if fromExprs[idx].IsArrayLiteral() {
				fromCXAtomicOpInputs := prgrm.GetCXArgFromArray(prgrm.CXAtomicOps[fromCXAtomicOpIdx].Inputs[0])
				sym.Size = fromCXAtomicOpInputs.Size
				sym.TotalSize = fromCXAtomicOpInputs.TotalSize
				sym.Lengths = fromCXAtomicOpInputs.Lengths
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
		symIdx := prgrm.AddCXArgInArray(sym)

		cxAtomicOp.AddOutput(prgrm, symIdx)

		for _, toExpr := range toExprs {
			if toExpr.Type == ast.CX_LINE {
				continue
			}
			toExprAtomicOp, _, _, err := prgrm.GetOperation(&toExpr)
			if err != nil {
				panic(err)
			}

			toExprAtomicOpOutput := prgrm.GetCXArgFromArray(toExprAtomicOp.Outputs[0])
			toExprAtomicOpOutput.PreviouslyDeclared = true
			toExprAtomicOpOutput.IsShortAssignmentDeclaration = true
		}

		toExprs = append([]ast.CXExpression{*exprCXLine, *expr}, toExprs...)
	case ">>=":
		expr = ast.MakeAtomicOperatorExpression(prgrm, ast.Natives[constants.OP_BITSHR])
		return ShortAssignment(prgrm, expr, exprCXLine, toExprs, fromExprs, pkg)
	case "<<=":
		expr = ast.MakeAtomicOperatorExpression(prgrm, ast.Natives[constants.OP_BITSHL])
		return ShortAssignment(prgrm, expr, exprCXLine, toExprs, fromExprs, pkg)
	case "+=":
		expr = ast.MakeAtomicOperatorExpression(prgrm, ast.Natives[constants.OP_ADD])
		return ShortAssignment(prgrm, expr, exprCXLine, toExprs, fromExprs, pkg)
	case "-=":
		expr = ast.MakeAtomicOperatorExpression(prgrm, ast.Natives[constants.OP_SUB])
		return ShortAssignment(prgrm, expr, exprCXLine, toExprs, fromExprs, pkg)
	case "*=":
		expr = ast.MakeAtomicOperatorExpression(prgrm, ast.Natives[constants.OP_MUL])
		return ShortAssignment(prgrm, expr, exprCXLine, toExprs, fromExprs, pkg)
	case "/=":
		expr = ast.MakeAtomicOperatorExpression(prgrm, ast.Natives[constants.OP_DIV])
		return ShortAssignment(prgrm, expr, exprCXLine, toExprs, fromExprs, pkg)
	case "%=":
		expr = ast.MakeAtomicOperatorExpression(prgrm, ast.Natives[constants.OP_MOD])
		return ShortAssignment(prgrm, expr, exprCXLine, toExprs, fromExprs, pkg)
	case "&=":
		expr = ast.MakeAtomicOperatorExpression(prgrm, ast.Natives[constants.OP_BITAND])
		return ShortAssignment(prgrm, expr, exprCXLine, toExprs, fromExprs, pkg)
	case "^=":
		expr = ast.MakeAtomicOperatorExpression(prgrm, ast.Natives[constants.OP_BITXOR])
		return ShortAssignment(prgrm, expr, exprCXLine, toExprs, fromExprs, pkg)
	case "|=":
		expr = ast.MakeAtomicOperatorExpression(prgrm, ast.Natives[constants.OP_BITOR])
		return ShortAssignment(prgrm, expr, exprCXLine, toExprs, fromExprs, pkg)
	}

	toLastExprAtomicOp, err := prgrm.GetCXAtomicOpFromExpressions(toExprs, len(toExprs)-1)
	if err != nil {
		panic(err)
	}

	if fromCXAtomicOpOperator == nil {
		opIdx := prgrm.AddNativeFunctionInArray(ast.Natives[constants.OP_IDENTITY])
		prgrm.CXAtomicOps[fromCXAtomicOpIdx].Operator = opIdx

		toCXAtomicOpOutput := prgrm.GetCXArgFromArray(toCXAtomicOp.Outputs[0])
		fromCXAtomicOpOutput := prgrm.GetCXArgFromArray(prgrm.CXAtomicOps[fromCXAtomicOpIdx].Outputs[0])

		toCXAtomicOpOutput.Size = fromCXAtomicOpOutput.Size
		toCXAtomicOpOutput.TotalSize = fromCXAtomicOpOutput.TotalSize
		toCXAtomicOpOutput.Type = fromCXAtomicOpOutput.Type
		toCXAtomicOpOutput.PointerTargetType = fromCXAtomicOpOutput.PointerTargetType
		toCXAtomicOpOutput.Lengths = fromCXAtomicOpOutput.Lengths
		toCXAtomicOpOutput.PassBy = fromCXAtomicOpOutput.PassBy
		toCXAtomicOpOutput.DoesEscape = fromCXAtomicOpOutput.DoesEscape
		// toCXAtomicOp.ProgramOutput[0].Program = prgrm

		if fromExprs[idx].IsMethodCall() {
			prgrm.CXAtomicOps[fromCXAtomicOpIdx].Inputs = append(prgrm.CXAtomicOps[fromCXAtomicOpIdx].Outputs, prgrm.CXAtomicOps[fromCXAtomicOpIdx].Inputs...)
		} else {
			prgrm.CXAtomicOps[fromCXAtomicOpIdx].Inputs = prgrm.CXAtomicOps[fromCXAtomicOpIdx].Outputs
		}

		prgrm.CXAtomicOps[fromCXAtomicOpIdx].Outputs = toLastExprAtomicOp.Outputs

		return append(toExprs[:len(toExprs)-1], fromExprs...)
	} else {
		fromCXAtomicOpOperatorOutput := prgrm.GetCXArgFromArray(fromCXAtomicOpOperator.Outputs[0])
		if fromCXAtomicOpOperator.IsBuiltIn() {
			// only assigning as if the operator had only one output defined

			toCXAtomicOpOutput := prgrm.GetCXArgFromArray(toCXAtomicOp.Outputs[0])
			if fromCXAtomicOpOperator.AtomicOPCode != constants.OP_IDENTITY {
				// it's a short variable declaration
				toCXAtomicOpOutput.Size = fromCXAtomicOpOperatorOutput.Size
				toCXAtomicOpOutput.Type = fromCXAtomicOpOperatorOutput.Type
				toCXAtomicOpOutput.PointerTargetType = fromCXAtomicOpOperatorOutput.PointerTargetType
				toCXAtomicOpOutput.Lengths = fromCXAtomicOpOperatorOutput.Lengths
			}

			toCXAtomicOpOutput.DoesEscape = fromCXAtomicOpOperatorOutput.DoesEscape
			toCXAtomicOpOutput.PassBy = fromCXAtomicOpOperatorOutput.PassBy
		} else {
			// we'll delegate multiple-value returns to the 'expression' grammar rule
			// only assigning as if the operator had only one output defined
			toCXAtomicOpOutput := prgrm.GetCXArgFromArray(toCXAtomicOp.Outputs[0])

			toCXAtomicOpOutput.Size = fromCXAtomicOpOperatorOutput.Size
			toCXAtomicOpOutput.Type = fromCXAtomicOpOperatorOutput.Type
			toCXAtomicOpOutput.PointerTargetType = fromCXAtomicOpOperatorOutput.PointerTargetType
			toCXAtomicOpOutput.Lengths = fromCXAtomicOpOperatorOutput.Lengths
			toCXAtomicOpOutput.DoesEscape = fromCXAtomicOpOperatorOutput.DoesEscape
			toCXAtomicOpOutput.PassBy = fromCXAtomicOpOperatorOutput.PassBy
		}

		prgrm.CXAtomicOps[fromCXAtomicOpIdx].Outputs = toLastExprAtomicOp.Outputs

		return append(toExprs[:len(toExprs)-1], fromExprs...)
	}
}
