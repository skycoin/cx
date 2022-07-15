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

		expressionOutputTypeSigIdx := expression.GetOutputs(prgrm)[0]
		expressionOutput := prgrm.GetCXTypeSignatureFromArray(expressionOutputTypeSigIdx)

		var expressionOutputIdx ast.CXArgumentIndex
		if expressionOutput.Type == ast.TYPE_CXARGUMENT_DEPRECATE {
			expressionOutputIdx = ast.CXArgumentIndex(expressionOutput.Meta)
		} else {
			panic("type is not type cx argument deprecate\n\n")
		}
		prgrm.CXArgs[expressionOutputIdx].Name = structLiteralName
		expressionOutput.Name = structLiteralName

		toExpressionOutputTypeSig := prgrm.GetCXTypeSignatureFromArray(toExpression.GetOutputs(prgrm)[0])
		var toExpressionOutput *ast.CXArgument = &ast.CXArgument{}
		if toExpressionOutputTypeSig.Type == ast.TYPE_CXARGUMENT_DEPRECATE {
			toExpressionOutput = prgrm.GetCXArgFromArray(ast.CXArgumentIndex(toExpressionOutputTypeSig.Meta))
		} else {
			panic("type is not type cx argument deprecate\n\n")
		}

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

	lastFromExpressionOutputTypeSig := prgrm.GetCXTypeSignatureFromArray(lastFromExpression.GetOutputs(prgrm)[0])
	var lastFromExpressionOutputArg *ast.CXArgument = &ast.CXArgument{}
	if lastFromExpressionOutputTypeSig.Type == ast.TYPE_CXARGUMENT_DEPRECATE {
		lastFromExpressionOutputArg = prgrm.GetCXArgFromArray(ast.CXArgumentIndex(lastFromExpressionOutputTypeSig.Meta))
	} else {
		panic("type is not cx argument deprecate\n\n")
	}

	toExpressionOutputTypeSig := prgrm.GetCXTypeSignatureFromArray(toExpression.GetOutputs(prgrm)[0])

	// If the last expression in `fromExprs` is declared as pointer
	// then it means the whole struct literal needs to be passed by reference.
	if !hasDeclSpec(lastFromExpressionOutputArg.GetAssignmentElement(prgrm), constants.DECL_POINTER) {
		return assignStructLiteralFields(prgrm, toExprs, fromExprs, toExpressionOutputTypeSig.Name)
	} else {
		// And we also need an auxiliary variable to point to,
		// otherwise we'd be trying to assign the fields to a nil value.
		outField := lastFromExpressionOutputArg
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

		typeSig := ast.GetCXTypeSignatureRepresentationOfCXArg_ForGlobals_CXAtomicOps(prgrm, prgrm.GetCXArgFromArray(auxIdx))
		typeSigIdx := prgrm.AddCXTypeSignatureInArray(typeSig)
		prgrm.CXAtomicOps[declExpressionIdx].AddOutput(prgrm, typeSigIdx)

		fromExprs = assignStructLiteralFields(prgrm, toExprs, fromExprs, auxName)

		assignExprCXLine := ast.MakeCXLineExpression(prgrm, lastFromCXLine.FileName, lastFromCXLine.LineNumber, lastFromCXLine.LineStr)
		assignExpr := ast.MakeAtomicOperatorExpression(prgrm, ast.Natives[constants.OP_IDENTITY])
		assignExpressionIdx := assignExpr.Index

		prgrm.CXAtomicOps[assignExpressionIdx].Package = lastFromExpression.Package
		out := ast.MakeArgument(toExpressionOutputTypeSig.Name, lastFromCXLine.FileName, lastFromCXLine.LineNumber)
		out.Package = lastFromExpression.Package
		outIdx := prgrm.AddCXArgInArray(out)

		outTypeSig := ast.GetCXTypeSignatureRepresentationOfCXArg_ForGlobals_CXAtomicOps(prgrm, prgrm.GetCXArgFromArray(outIdx))
		outTypeSigIdx := prgrm.AddCXTypeSignatureInArray(outTypeSig)
		prgrm.CXAtomicOps[assignExpressionIdx].AddOutput(prgrm, outTypeSigIdx)

		auxTypeSig := ast.GetCXTypeSignatureRepresentationOfCXArg_ForGlobals_CXAtomicOps(prgrm, prgrm.GetCXArgFromArray(auxIdx))
		auxTypeSigIdx := prgrm.AddCXTypeSignatureInArray(auxTypeSig)
		prgrm.CXAtomicOps[assignExpressionIdx].AddInput(prgrm, auxTypeSigIdx)

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

	toExpressionOutput := prgrm.GetCXTypeSignatureFromArray(toExpression.GetOutputs(prgrm)[0])

	for _, expr := range fromExprs {
		expression, err := prgrm.GetCXAtomicOp(expr.Index)
		if err != nil {
			panic(err)
		}

		expressionOutput := prgrm.GetCXTypeSignatureFromArray(expression.GetOutputs(prgrm)[0])
		var expressionOutputIdx ast.CXArgumentIndex
		if expressionOutput.Type == ast.TYPE_CXARGUMENT_DEPRECATE {
			expressionOutputIdx = ast.CXArgumentIndex(expressionOutput.Meta)
		} else {
			panic("type is not type cx argument deprecate\n\n")
		}

		prgrm.CXArgs[expressionOutputIdx].Name = toExpressionOutput.Name
		expressionOutput.Name = toExpressionOutput.Name
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

	typeSig := toExpression.GetOutputs(prgrm)[0]

	prgrm.CXAtomicOps[expressionIdx].AddInput(prgrm, typeSig)
	prgrm.CXAtomicOps[expressionIdx].AddOutput(prgrm, typeSig)

	prgrm.CXAtomicOps[expressionIdx].Package = ast.CXPackageIndex(pkg.Index)

	if fromExpressionOperator == nil {
		typeSig := prgrm.CXAtomicOps[fromExpressionIdx].GetOutputs(prgrm)[0]
		prgrm.CXAtomicOps[expressionIdx].AddInput(prgrm, typeSig)
	} else {
		var fromExpressionInputArg *ast.CXArgument = &ast.CXArgument{}
		fromExpressionInputTypeSig := prgrm.GetCXTypeSignatureFromArray(prgrm.CXAtomicOps[fromExpressionIdx].GetInputs(prgrm)[0])
		if fromExpressionInputTypeSig.Type == ast.TYPE_CXARGUMENT_DEPRECATE {
			fromExpressionInputArg = prgrm.GetCXArgFromArray(ast.CXArgumentIndex(fromExpressionInputTypeSig.Meta))
		} else {
			panic("type is not cx argument deprecate\n\n")
		}

		sym := ast.MakeArgument(generateTempVarName(constants.LOCAL_PREFIX), CurrentFile, LineNo).SetType(fromExpressionInputArg.Type)
		sym.Package = ast.CXPackageIndex(pkg.Index)
		sym.PreviouslyDeclared = true
		symIdx := prgrm.AddCXArgInArray(sym)
		typeSig := ast.GetCXTypeSignatureRepresentationOfCXArg_ForGlobals_CXAtomicOps(prgrm, prgrm.GetCXArgFromArray(symIdx))
		typeSigIdx := prgrm.AddCXTypeSignatureInArray(typeSig)
		prgrm.CXAtomicOps[fromExpressionIdx].AddOutput(prgrm, typeSigIdx)
		prgrm.CXAtomicOps[expressionIdx].AddInput(prgrm, typeSigIdx)
	}

	//must check if from expression is naked previously declared variable
	if len(fromExprs) == 1 && fromExpressionOperator == nil && len(prgrm.CXAtomicOps[fromExpressionIdx].GetOutputs(prgrm)) > 0 && len(prgrm.CXAtomicOps[fromExpressionIdx].GetInputs(prgrm)) == 0 {
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
	expressionOperatorOutputs := expressionOperator.GetOutputs(prgrm)
	expressionOperatorOutputTypeSig := prgrm.GetCXTypeSignatureFromArray(expressionOperatorOutputs[0])
	var expressionOperatorOutputArg *ast.CXArgument = &ast.CXArgument{}
	if expressionOperatorOutputTypeSig.Type == ast.TYPE_CXARGUMENT_DEPRECATE {
		expressionOperatorOutputArg = prgrm.GetCXArgFromArray(ast.CXArgumentIndex(expressionOperatorOutputTypeSig.Meta))
	} else {
		panic("type is not cx argument deprecate\n\n")
	}

	if expressionOperatorOutputArg.Type != types.UNDEFINED {
		return expressionOperatorOutputArg
	}

	expressionInputTypeSig := prgrm.GetCXTypeSignatureFromArray(expression.GetInputs(prgrm)[0])
	var expressionInputArg *ast.CXArgument = &ast.CXArgument{}
	if expressionInputTypeSig.Type == ast.TYPE_CXARGUMENT_DEPRECATE {
		expressionInputArg = prgrm.GetCXArgFromArray(ast.CXArgumentIndex(expressionInputTypeSig.Meta))
	} else {
		panic("type is not cx argument deprecate\n\n")
	}

	return expressionInputArg
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
	fromExpressionOperatorOutputs := fromExpressionOperator.GetOutputs(prgrm)
	// Checking if we're trying to assign stuff from a function call
	// And if that function call actually returns something. If not, throw an error.
	if fromExpressionOperator != nil && len(fromExpressionOperatorOutputs) == 0 {
		toExpressionOutputs := toExpression.GetOutputs(prgrm)
		toExpressionOutputTypeSig := prgrm.GetCXTypeSignatureFromArray(toExpressionOutputs[0])

		var toExpressionOutputArg *ast.CXArgument = &ast.CXArgument{}
		if toExpressionOutputTypeSig.Type == ast.TYPE_CXARGUMENT_DEPRECATE {
			toExpressionOutputArg = prgrm.GetCXArgFromArray(ast.CXArgumentIndex(toExpressionOutputTypeSig.Meta))
		} else {
			panic("type is not cx argument deprecate\n\n")
		}

		println(ast.CompilationError(toExpressionOutputArg.ArgDetails.FileName, toExpressionOutputArg.ArgDetails.FileLine), "trying to use an outputless operator in an assignment")
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
		var sym *ast.CXArgument = &ast.CXArgument{}

		toExpressionOutputTypeSig := prgrm.GetCXTypeSignatureFromArray(toExpression.GetOutputs(prgrm)[0])

		if fromExpressionOperator == nil {
			fromExpressionOutputTypeSig := prgrm.GetCXTypeSignatureFromArray(prgrm.CXAtomicOps[fromExpressionIdx].GetOutputs(prgrm)[0])
			var fromExpressionOutputArg *ast.CXArgument = &ast.CXArgument{}
			if fromExpressionOutputTypeSig.Type == ast.TYPE_CXARGUMENT_DEPRECATE {
				fromExpressionOutputArg = prgrm.GetCXArgFromArray(ast.CXArgumentIndex(fromExpressionOutputTypeSig.Meta))
			} else {
				panic("type is not cx argument deprecate\n\n")
			}

			// then it's a literal
			sym = ast.MakeArgument(toExpressionOutputTypeSig.Name, CurrentFile, LineNo).SetType(fromExpressionOutputArg.Type)
		} else {
			outTypeArg := getOutputType(prgrm, &fromExprs[lastFromExpressionIdx])

			sym = ast.MakeArgument(toExpressionOutputTypeSig.Name, CurrentFile, LineNo).SetType(outTypeArg.Type)

			if fromExprs[lastFromExpressionIdx].IsArrayLiteral() {
				fromExpressionInputTypeSig := prgrm.GetCXTypeSignatureFromArray(prgrm.CXAtomicOps[fromExpressionIdx].GetInputs(prgrm)[0])
				var fromExpressionInputArg *ast.CXArgument = &ast.CXArgument{}
				if fromExpressionInputTypeSig.Type == ast.TYPE_CXARGUMENT_DEPRECATE {
					fromExpressionInputArg = prgrm.GetCXArgFromArray(ast.CXArgumentIndex(fromExpressionInputTypeSig.Meta))
				} else {
					panic("type is not cx argument deprecate\n\n")
				}

				fromCXAtomicOpInputs := fromExpressionInputArg
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
		symIdx := prgrm.AddCXArgInArray(sym)

		typeSig := ast.GetCXTypeSignatureRepresentationOfCXArg_ForGlobals_CXAtomicOps(prgrm, prgrm.GetCXArgFromArray(symIdx))
		typeSigIdx := prgrm.AddCXTypeSignatureInArray(typeSig)
		prgrm.CXAtomicOps[expressionIdx].AddOutput(prgrm, typeSigIdx)

		for _, toExpr := range toExprs {
			if toExpr.Type == ast.CX_LINE {
				continue
			}
			toExprAtomicOp, err := prgrm.GetCXAtomicOp(toExpr.Index)
			if err != nil {
				panic(err)
			}

			toExprAtomicOpOutputTypeSig := prgrm.GetCXTypeSignatureFromArray(toExprAtomicOp.GetOutputs(prgrm)[0])
			var toExprAtomicOpOutputIdx ast.CXArgumentIndex
			if toExprAtomicOpOutputTypeSig.Type == ast.TYPE_CXARGUMENT_DEPRECATE {
				toExprAtomicOpOutputIdx = ast.CXArgumentIndex(toExprAtomicOpOutputTypeSig.Meta)
			} else {
				panic("type is not type cx argument deprecate\n\n")
			}

			prgrm.CXArgs[toExprAtomicOpOutputIdx].PreviouslyDeclared = true

			prgrm.CXArgs[toExprAtomicOpOutputIdx].Type = sym.Type
			prgrm.CXArgs[toExprAtomicOpOutputIdx].PointerTargetType = sym.PointerTargetType
			prgrm.CXArgs[toExprAtomicOpOutputIdx].Size = sym.Size
			prgrm.CXArgs[toExprAtomicOpOutputIdx].TotalSize = sym.TotalSize
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

		toExpressionOutputTypeSig := prgrm.GetCXTypeSignatureFromArray(toExpression.GetOutputs(prgrm)[0])
		var toExpressionOutputIdx ast.CXArgumentIndex
		if toExpressionOutputTypeSig.Type == ast.TYPE_CXARGUMENT_DEPRECATE {
			toExpressionOutputIdx = ast.CXArgumentIndex(toExpressionOutputTypeSig.Meta)
		} else {
			panic("type is not type cx argument deprecate\n\n")
		}

		fromExpressionOutputTypeSig := prgrm.GetCXTypeSignatureFromArray(prgrm.CXAtomicOps[fromExpressionIdx].GetOutputs(prgrm)[0])
		var fromExpressionOutputArg *ast.CXArgument = &ast.CXArgument{}
		if fromExpressionOutputTypeSig.Type == ast.TYPE_CXARGUMENT_DEPRECATE {
			fromExpressionOutputArg = prgrm.GetCXArgFromArray(ast.CXArgumentIndex(fromExpressionOutputTypeSig.Meta))
		} else {
			panic("type is not type cx argument deprecate\n\n")
		}

		// if fromExpressionOutputArg type is atomic then change the type signature to type atomic
		if !fromExpressionOutputArg.IsSlice && len(fromExpressionOutputArg.Lengths) == 0 && fromExpressionOutputArg.Type.IsPrimitive() {
			toExpressionOutputTypeSig.Type = ast.TYPE_ATOMIC
			toExpressionOutputTypeSig.Meta = int(fromExpressionOutputArg.Type)
			toExpressionOutputTypeSig.Offset = prgrm.CXArgs[toExpressionOutputIdx].Offset

			fromExpressionOutputTypeSig.Type = ast.TYPE_ATOMIC
			fromExpressionOutputTypeSig.Meta = int(fromExpressionOutputArg.Type)
			fromExpressionOutputTypeSig.Offset = fromExpressionOutputArg.Offset
			fromExpressionOutputTypeSig.Package = fromExpressionOutputArg.Package
		} else {
			prgrm.CXArgs[toExpressionOutputIdx].Size = fromExpressionOutputArg.Size
			prgrm.CXArgs[toExpressionOutputIdx].TotalSize = fromExpressionOutputArg.TotalSize
			prgrm.CXArgs[toExpressionOutputIdx].Type = fromExpressionOutputArg.Type
			prgrm.CXArgs[toExpressionOutputIdx].PointerTargetType = fromExpressionOutputArg.PointerTargetType
			prgrm.CXArgs[toExpressionOutputIdx].Lengths = fromExpressionOutputArg.Lengths
			prgrm.CXArgs[toExpressionOutputIdx].PassBy = fromExpressionOutputArg.PassBy
		}

		if fromExprs[lastFromExpressionIdx].IsMethodCall() {
			newInputs := prgrm.CXAtomicOps[fromExpressionIdx].Outputs
			if prgrm.CXAtomicOps[fromExpressionIdx].Inputs != nil {
				for _, typeSig := range prgrm.CXAtomicOps[fromExpressionIdx].Inputs.Fields {
					newInputs.AddField_CXAtomicOps(prgrm, typeSig)
				}
			}
			prgrm.CXAtomicOps[fromExpressionIdx].Inputs = newInputs
		} else {
			prgrm.CXAtomicOps[fromExpressionIdx].Inputs = prgrm.CXAtomicOps[fromExpressionIdx].Outputs
		}

		prgrm.CXAtomicOps[fromExpressionIdx].Outputs = toLastExpr.Outputs

		return append(toExprs[:len(toExprs)-1], fromExprs...)
	} else {
		fromExpressionOperatorOutputs := fromExpressionOperator.GetOutputs(prgrm)
		fromExpressionOperatorOutputTypeSig := prgrm.GetCXTypeSignatureFromArray(fromExpressionOperatorOutputs[0])
		var fromCXAtomicOpOperatorOutput *ast.CXArgument = &ast.CXArgument{}
		if fromExpressionOperatorOutputTypeSig.Type == ast.TYPE_CXARGUMENT_DEPRECATE {
			fromCXAtomicOpOperatorOutput = prgrm.GetCXArgFromArray(ast.CXArgumentIndex(fromExpressionOperatorOutputTypeSig.Meta))

			if fromExpressionOperator.IsBuiltIn() {
				// only assigning as if the operator had only one output defined

				toExpressionOutputTypeSig := prgrm.GetCXTypeSignatureFromArray(toExpression.GetOutputs(prgrm)[0])

				var toExpressionOutputIdx ast.CXArgumentIndex
				if toExpressionOutputTypeSig.Type == ast.TYPE_CXARGUMENT_DEPRECATE {
					toExpressionOutputIdx = ast.CXArgumentIndex(toExpressionOutputTypeSig.Meta)
				} else {
					panic("type is not type cx argument deprecate\n\n")
				}
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
				toExpressionOutputTypeSig := prgrm.GetCXTypeSignatureFromArray(toExpression.GetOutputs(prgrm)[0])

				var toExpressionOutputIdx ast.CXArgumentIndex
				if toExpressionOutputTypeSig.Type == ast.TYPE_CXARGUMENT_DEPRECATE {
					toExpressionOutputIdx = ast.CXArgumentIndex(toExpressionOutputTypeSig.Meta)
				} else {
					panic("type is not type cx argument deprecate\n\n")
				}

				prgrm.CXArgs[toExpressionOutputIdx].Size = fromCXAtomicOpOperatorOutput.Size
				prgrm.CXArgs[toExpressionOutputIdx].Type = fromCXAtomicOpOperatorOutput.Type
				prgrm.CXArgs[toExpressionOutputIdx].PointerTargetType = fromCXAtomicOpOperatorOutput.PointerTargetType
				prgrm.CXArgs[toExpressionOutputIdx].Lengths = fromCXAtomicOpOperatorOutput.Lengths
				prgrm.CXArgs[toExpressionOutputIdx].PassBy = fromCXAtomicOpOperatorOutput.PassBy
			}
		} else {
			// panic("type is not cx argument deprecate\n\n")

			if fromExpressionOperator.IsBuiltIn() {
				// only assigning as if the operator had only one output defined

				toExpressionOutputTypeSig := prgrm.GetCXTypeSignatureFromArray(toExpression.GetOutputs(prgrm)[0])
				if fromExpressionOperator.AtomicOPCode != constants.OP_IDENTITY {
					// it's a short variable declaration

					toExpressionOutputTypeSig.Type = ast.TYPE_ATOMIC
					toExpressionOutputTypeSig.Meta = fromExpressionOperatorOutputTypeSig.Meta
				}

			} else {
				// we'll delegate multiple-value returns to the 'expression' grammar rule
				// only assigning as if the operator had only one output defined
				toExpressionOutputTypeSig := prgrm.GetCXTypeSignatureFromArray(toExpression.GetOutputs(prgrm)[0])

				toExpressionOutputTypeSig.Type = ast.TYPE_ATOMIC
				toExpressionOutputTypeSig.Meta = fromExpressionOperatorOutputTypeSig.Meta
			}
		}

		prgrm.CXAtomicOps[fromExpressionIdx].Outputs = toLastExpr.Outputs

		return append(toExprs[:len(toExprs)-1], fromExprs...)
	}
}
