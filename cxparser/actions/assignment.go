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
	toExpression, err := prgrm.GetCXAtomicOpFromExpressions(toExprs, 0)
	if err != nil {
		panic(err)
	}

	for _, exprInfo := range fromExprs {
		if exprInfo.Type == ast.CX_LINE {
			continue
		}
		expression, err := prgrm.GetCXAtomicOp(exprInfo.Index)
		if err != nil {
			panic(err)
		}

		expressionOutputIdx := expression.Outputs[0]
		prgrm.CXArgs[expressionOutputIdx].Name = structLiteralName

		toExpressionOutput := prgrm.GetCXArgFromArray(toExpression.Outputs[0])
		if len(toExpressionOutput.Indexes) > 0 {
			prgrm.CXArgs[expressionOutputIdx].Lengths = toExpressionOutput.Lengths
			prgrm.CXArgs[expressionOutputIdx].Indexes = toExpressionOutput.Indexes
			prgrm.CXArgs[expressionOutputIdx].DereferenceOperations = append(prgrm.CXArgs[expressionOutputIdx].DereferenceOperations, constants.DEREF_ARRAY)
		}

		prgrm.CXArgs[expressionOutputIdx].DereferenceOperations = append(prgrm.CXArgs[expressionOutputIdx].DereferenceOperations, constants.DEREF_FIELD)
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

	lastFromExpression, err := prgrm.GetCXAtomicOp(lastFromExpr.Index)
	if err != nil {
		panic(err)
	}

	toExpression, err := prgrm.GetCXAtomicOpFromExpressions(toExprs, 0)
	if err != nil {
		panic(err)
	}

	lastFromCXLine, _ := prgrm.GetPreviousCXLine(fromExprs, len(fromExprs)-1)

	// If the last expression in `fromExprs` is declared as pointer
	// then it means the whole struct literal needs to be passed by reference.
	if !hasDeclSpec(prgrm.GetCXArgFromArray(lastFromExpression.Outputs[0]).GetAssignmentElement(prgrm), constants.DECL_POINTER) {
		return assignStructLiteralFields(prgrm, toExprs, fromExprs, prgrm.GetCXArgFromArray(toExpression.Outputs[0]).Name)
	} else {
		// And we also need an auxiliary variable to point to,
		// otherwise we'd be trying to assign the fields to a nil value.
		outField := prgrm.GetCXArgFromArray(lastFromExpression.Outputs[0])
		auxName := generateTempVarName(constants.LOCAL_PREFIX)
		aux := ast.MakeArgument(auxName, lastFromCXLine.FileName, lastFromCXLine.LineNumber)
		aux.SetType(outField.Type)
		aux.DeclarationSpecifiers = append(aux.DeclarationSpecifiers, constants.DECL_POINTER)
		aux.StructType = outField.StructType
		aux.Size = outField.Size
		aux.TotalSize = outField.TotalSize
		aux.PreviouslyDeclared = true
		aux.Package = lastFromExpression.Package
		auxIdx := prgrm.AddCXArgInArray(aux)

		declExprCXLine := ast.MakeCXLineExpression(prgrm, lastFromCXLine.FileName, lastFromCXLine.LineNumber, lastFromCXLine.LineStr)
		declExpr := ast.MakeAtomicOperatorExpression(prgrm, nil)
		declExpressionIdx := declExpr.Index

		prgrm.CXAtomicOps[declExpressionIdx].Package = lastFromExpression.Package
		prgrm.CXAtomicOps[declExpressionIdx].AddOutput(prgrm, auxIdx)

		fromExprs = assignStructLiteralFields(prgrm, toExprs, fromExprs, auxName)

		assignExprCXLine := ast.MakeCXLineExpression(prgrm, lastFromCXLine.FileName, lastFromCXLine.LineNumber, lastFromCXLine.LineStr)
		assignExpr := ast.MakeAtomicOperatorExpression(prgrm, ast.Natives[constants.OP_IDENTITY])
		assignExpressionIdx := assignExpr.Index

		prgrm.CXAtomicOps[assignExpressionIdx].Package = lastFromExpression.Package
		out := ast.MakeArgument(prgrm.GetCXArgFromArray(toExpression.Outputs[0]).Name, lastFromCXLine.FileName, lastFromCXLine.LineNumber)
		out.PassBy = constants.PASSBY_REFERENCE
		out.Package = lastFromExpression.Package
		outIdx := prgrm.AddCXArgInArray(out)

		prgrm.CXAtomicOps[assignExpressionIdx].AddOutput(prgrm, outIdx)
		prgrm.CXAtomicOps[assignExpressionIdx].AddInput(prgrm, auxIdx)

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
	toExpression, err := prgrm.GetCXAtomicOpFromExpressions(toExprs, 0)
	if err != nil {
		panic(err)
	}

	for _, expr := range fromExprs {
		expression, err := prgrm.GetCXAtomicOp(expr.Index)
		if err != nil {
			panic(err)
		}

		expressionOutputIdx := expression.Outputs[0]
		prgrm.CXArgs[expressionOutputIdx].Name = prgrm.GetCXArgFromArray(toExpression.Outputs[0]).Name
		prgrm.CXArgs[expressionOutputIdx].DereferenceOperations = append(prgrm.CXArgs[expressionOutputIdx].DereferenceOperations, constants.DEREF_ARRAY)
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
	expressionIdx := expr.Index

	toExpression, err := prgrm.GetCXAtomicOpFromExpressions(toExprs, 0)
	if err != nil {
		panic(err)
	}

	fromExpressionIdx := fromExprs[len(fromExprs)-1].Index

	fromExpressionOperator := prgrm.GetFunctionFromArray(prgrm.CXAtomicOps[fromExpressionIdx].Operator)

	prgrm.CXAtomicOps[expressionIdx].AddInput(prgrm, toExpression.Outputs[0])
	prgrm.CXAtomicOps[expressionIdx].AddOutput(prgrm, toExpression.Outputs[0])

	prgrm.CXAtomicOps[expressionIdx].Package = ast.CXPackageIndex(pkg.Index)

	if fromExpressionOperator == nil {
		prgrm.CXAtomicOps[expressionIdx].AddInput(prgrm, prgrm.CXAtomicOps[fromExpressionIdx].Outputs[0])
	} else {
		sym := ast.MakeArgument(generateTempVarName(constants.LOCAL_PREFIX), CurrentFile, LineNo).SetType(prgrm.GetCXArgFromArray(prgrm.CXAtomicOps[fromExpressionIdx].Inputs[0]).Type)
		sym.Package = ast.CXPackageIndex(pkg.Index)
		sym.PreviouslyDeclared = true
		symIdx := prgrm.AddCXArgInArray(sym)
		prgrm.CXAtomicOps[fromExpressionIdx].AddOutput(prgrm, symIdx)

		prgrm.CXAtomicOps[expressionIdx].AddInput(prgrm, symIdx)
	}

	//must check if from expression is naked previously declared variable
	if len(fromExprs) == 1 && fromExpressionOperator == nil && len(prgrm.CXAtomicOps[fromExpressionIdx].Outputs) > 0 && len(prgrm.CXAtomicOps[fromExpressionIdx].Inputs) == 0 {
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
	expression, err := prgrm.GetCXAtomicOp(expr.Index)
	if err != nil {
		panic(err)
	}
	expressionOperator := prgrm.GetFunctionFromArray(expression.Operator)

	if prgrm.GetCXArgFromArray(expressionOperator.Outputs[0]).Type != types.UNDEFINED {
		return prgrm.GetCXArgFromArray(expressionOperator.Outputs[0])
	}

	return prgrm.GetCXArgFromArray(expression.Inputs[0])
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
	lastFromExpressionIdx := len(fromExprs) - 1

	toExpression, err := prgrm.GetCXAtomicOpFromExpressions(toExprs, 0)
	if err != nil {
		panic(err)
	}

	fromExpressionIdx := fromExprs[lastFromExpressionIdx].Index
	fromExpressionOperator := prgrm.GetFunctionFromArray(prgrm.CXAtomicOps[fromExpressionIdx].Operator)

	// Checking if we're trying to assign stuff from a function call
	// And if that function call actually returns something. If not, throw an error.
	if fromExpressionOperator != nil && len(fromExpressionOperator.Outputs) == 0 {
		println(ast.CompilationError(prgrm.GetCXArgFromArray(toExpression.Outputs[0]).ArgDetails.FileName, prgrm.GetCXArgFromArray(toExpression.Outputs[0]).ArgDetails.FileLine), "trying to use an outputless operator in an assignment")
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
		expressionIdx := expr.Index
		prgrm.CXAtomicOps[expressionIdx].Package = ast.CXPackageIndex(pkg.Index)
		var sym *ast.CXArgument

		if fromExpressionOperator == nil {
			// then it's a literal
			sym = ast.MakeArgument(prgrm.GetCXArgFromArray(toExpression.Outputs[0]).Name, CurrentFile, LineNo).SetType(prgrm.GetCXArgFromArray(prgrm.CXAtomicOps[fromExpressionIdx].Outputs[0]).Type)
		} else {
			outTypeArg := getOutputType(prgrm, &fromExprs[lastFromExpressionIdx])

			sym = ast.MakeArgument(prgrm.GetCXArgFromArray(toExpression.Outputs[0]).Name, CurrentFile, LineNo).SetType(outTypeArg.Type)

			if fromExprs[lastFromExpressionIdx].IsArrayLiteral() {
				fromCXAtomicOpInputs := prgrm.GetCXArgFromArray(prgrm.CXAtomicOps[fromExpressionIdx].Inputs[0])
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

		prgrm.CXAtomicOps[expressionIdx].AddOutput(prgrm, symIdx)

		for _, toExpr := range toExprs {
			if toExpr.Type == ast.CX_LINE {
				continue
			}
			toExprAtomicOp, err := prgrm.GetCXAtomicOp(toExpr.Index)
			if err != nil {
				panic(err)
			}

			toExprAtomicOpOutputIdx := toExprAtomicOp.Outputs[0]
			prgrm.CXArgs[toExprAtomicOpOutputIdx].PreviouslyDeclared = true
			prgrm.CXArgs[toExprAtomicOpOutputIdx].IsShortAssignmentDeclaration = true
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

	toLastExpr, err := prgrm.GetCXAtomicOpFromExpressions(toExprs, len(toExprs)-1)
	if err != nil {
		panic(err)
	}

	if fromExpressionOperator == nil {
		opIdx := prgrm.AddNativeFunctionInArray(ast.Natives[constants.OP_IDENTITY])
		prgrm.CXAtomicOps[fromExpressionIdx].Operator = opIdx

		toExpressionOutputIdx := toExpression.Outputs[0]
		fromExpressionOutput := prgrm.GetCXArgFromArray(prgrm.CXAtomicOps[fromExpressionIdx].Outputs[0])

		prgrm.CXArgs[toExpressionOutputIdx].Size = fromExpressionOutput.Size
		prgrm.CXArgs[toExpressionOutputIdx].TotalSize = fromExpressionOutput.TotalSize
		prgrm.CXArgs[toExpressionOutputIdx].Type = fromExpressionOutput.Type
		prgrm.CXArgs[toExpressionOutputIdx].PointerTargetType = fromExpressionOutput.PointerTargetType
		prgrm.CXArgs[toExpressionOutputIdx].Lengths = fromExpressionOutput.Lengths
		prgrm.CXArgs[toExpressionOutputIdx].PassBy = fromExpressionOutput.PassBy

		if fromExprs[lastFromExpressionIdx].IsMethodCall() {
			prgrm.CXAtomicOps[fromExpressionIdx].Inputs = append(prgrm.CXAtomicOps[fromExpressionIdx].Outputs, prgrm.CXAtomicOps[fromExpressionIdx].Inputs...)
		} else {
			prgrm.CXAtomicOps[fromExpressionIdx].Inputs = prgrm.CXAtomicOps[fromExpressionIdx].Outputs
		}

		prgrm.CXAtomicOps[fromExpressionIdx].Outputs = toLastExpr.Outputs

		return append(toExprs[:len(toExprs)-1], fromExprs...)
	} else {
		fromCXAtomicOpOperatorOutput := prgrm.GetCXArgFromArray(fromExpressionOperator.Outputs[0])
		if fromExpressionOperator.IsBuiltIn() {
			// only assigning as if the operator had only one output defined

			toExpressionOutputIdx := toExpression.Outputs[0]
			if fromExpressionOperator.AtomicOPCode != constants.OP_IDENTITY {
				// it's a short variable declaration
				prgrm.CXArgs[toExpressionOutputIdx].Size = fromCXAtomicOpOperatorOutput.Size
				prgrm.CXArgs[toExpressionOutputIdx].Type = fromCXAtomicOpOperatorOutput.Type
				prgrm.CXArgs[toExpressionOutputIdx].PointerTargetType = fromCXAtomicOpOperatorOutput.PointerTargetType
				prgrm.CXArgs[toExpressionOutputIdx].Lengths = fromCXAtomicOpOperatorOutput.Lengths
			}

			prgrm.CXArgs[toExpressionOutputIdx].PassBy = fromCXAtomicOpOperatorOutput.PassBy
		} else {
			// we'll delegate multiple-value returns to the 'expression' grammar rule
			// only assigning as if the operator had only one output defined
			toExpressionOutputIdx := toExpression.Outputs[0]

			prgrm.CXArgs[toExpressionOutputIdx].Size = fromCXAtomicOpOperatorOutput.Size
			prgrm.CXArgs[toExpressionOutputIdx].Type = fromCXAtomicOpOperatorOutput.Type
			prgrm.CXArgs[toExpressionOutputIdx].PointerTargetType = fromCXAtomicOpOperatorOutput.PointerTargetType
			prgrm.CXArgs[toExpressionOutputIdx].Lengths = fromCXAtomicOpOperatorOutput.Lengths
			prgrm.CXArgs[toExpressionOutputIdx].PassBy = fromCXAtomicOpOperatorOutput.PassBy
		}

		prgrm.CXAtomicOps[fromExpressionIdx].Outputs = toLastExpr.Outputs

		return append(toExprs[:len(toExprs)-1], fromExprs...)
	}
}
