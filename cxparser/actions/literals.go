package actions

import (
	"github.com/skycoin/cx/cx/ast"
	"github.com/skycoin/cx/cx/constants"
	"github.com/skycoin/cx/cx/types"
	"github.com/skycoin/skycoin/src/cipher/encoder"
)

// SliceLiteralExpression handles literal expressions by converting it to a series of `append` expressions.
//
// Input arguments description:
// 	prgrm - a CXProgram that contains all the data and arrays of the program.
//  typeCode - type code of the slice.
//  exprs - array of expressions that make up the slice literal expression.
func SliceLiteralExpression(prgrm *ast.CXProgram, typeCode types.Code, exprs []ast.CXExpression) []ast.CXExpression {
	var result []ast.CXExpression

	pkg, err := prgrm.GetCurrentPackage()
	if err != nil {
		panic(err)
	}

	symName := generateTempVarName(constants.LOCAL_PREFIX)

	slcVarExprCXLine := ast.MakeCXLineExpression(prgrm, CurrentFile, LineNo, LineStr)
	// adding the declaration
	slcVarExpr := ast.MakeAtomicOperatorExpression(prgrm, nil)
	slcVarExprAtomicOpIdx := slcVarExpr.Index

	prgrm.CXAtomicOps[slcVarExprAtomicOpIdx].Package = ast.CXPackageIndex(pkg.Index)
	slcVar := ast.MakeArgument(symName, CurrentFile, LineNo)
	slcVar.SetType(typeCode)
	slcVar = DeclarationSpecifiers(slcVar, []types.Pointer{0}, constants.DECL_SLICE)
	slcVar.TotalSize = types.POINTER_SIZE
	slcVar.Package = ast.CXPackageIndex(pkg.Index)
	slcVar.PreviouslyDeclared = true

	slcVarIdx := prgrm.AddCXArgInArray(slcVar)
	prgrm.CXAtomicOps[slcVarExprAtomicOpIdx].AddOutput(prgrm, slcVarIdx)

	result = append(result, *slcVarExprCXLine, *slcVarExpr)

	var endPointsCounter int
	for i, expr := range exprs {
		if expr.Type == ast.CX_LINE {
			continue
		}
		exprAtomicOp, err := prgrm.GetCXAtomicOp(expr.Index)
		if err != nil {
			panic(err)
		}
		exprAtomicOpOperator := prgrm.GetFunctionFromArray(exprAtomicOp.Operator)

		exprCXLine, _ := prgrm.GetPreviousCXLine(exprs, i)

		if expr.IsArrayLiteral() {
			expr.ExpressionType = ast.CXEXPR_UNUSED

			symInp := ast.MakeArgument(symName, CurrentFile, LineNo).SetType(typeCode)
			symInp.Package = ast.CXPackageIndex(pkg.Index)
			symInp.TotalSize = types.POINTER_SIZE
			symInpIdx := prgrm.AddCXArgInArray(symInp)

			symOut := ast.MakeArgument(symName, CurrentFile, LineNo).SetType(typeCode)
			symOut.Package = ast.CXPackageIndex(pkg.Index)
			symOut.TotalSize = types.POINTER_SIZE
			symOutIdx := prgrm.AddCXArgInArray(symOut)

			endPointsCounter++

			symExprExprCXLine := ast.MakeCXLineExpression(prgrm, CurrentFile, LineNo, LineStr)
			symExpr := ast.MakeAtomicOperatorExpression(prgrm, nil)
			symExprAtomicOpIdx := symExpr.Index

			prgrm.CXAtomicOps[symExprAtomicOpIdx].Package = ast.CXPackageIndex(pkg.Index)
			prgrm.CXAtomicOps[symExprAtomicOpIdx].AddOutput(prgrm, symOutIdx)

			if exprAtomicOpOperator == nil {
				// then it's a literal
				opIdx := prgrm.AddNativeFunctionInArray(ast.Natives[constants.OP_APPEND])
				prgrm.CXAtomicOps[symExprAtomicOpIdx].Operator = opIdx

				prgrm.CXAtomicOps[symExprAtomicOpIdx].Inputs = nil
				prgrm.CXAtomicOps[symExprAtomicOpIdx].Inputs = append(prgrm.CXAtomicOps[symExprAtomicOpIdx].Inputs, symInpIdx)
				prgrm.CXAtomicOps[symExprAtomicOpIdx].Inputs = append(prgrm.CXAtomicOps[symExprAtomicOpIdx].Inputs, exprAtomicOp.Outputs...)
			} else {
				// We need to create a temporary variable to hold the result of the
				// nested expressions. Then use that variable as part of the slice literal.
				out := ast.MakeArgument(generateTempVarName(constants.LOCAL_PREFIX), exprCXLine.FileName, exprCXLine.LineNumber)
				outArg := getOutputType(prgrm, &expr)
				out.SetType(outArg.Type)
				out.PointerTargetType = outArg.PointerTargetType
				out.StructType = outArg.StructType
				out.Size = outArg.Size
				out.TotalSize = ast.GetArgSize(prgrm, outArg)
				out.PreviouslyDeclared = true
				outIdx := prgrm.AddCXArgInArray(out)

				exprAtomicOp.Outputs = nil
				exprAtomicOp.AddOutput(prgrm, outIdx)

				result = append(result, expr)
				opIdx := prgrm.AddNativeFunctionInArray(ast.Natives[constants.OP_APPEND])
				prgrm.CXAtomicOps[symExprAtomicOpIdx].Operator = opIdx

				prgrm.CXAtomicOps[symExprAtomicOpIdx].Inputs = nil
				prgrm.CXAtomicOps[symExprAtomicOpIdx].Inputs = append(prgrm.CXAtomicOps[symExprAtomicOpIdx].Inputs, symInpIdx)
				prgrm.CXAtomicOps[symExprAtomicOpIdx].Inputs = append(prgrm.CXAtomicOps[symExprAtomicOpIdx].Inputs, outIdx)
			}
			result = append(result, *symExprExprCXLine, *symExpr)
		} else {
			result = append(result, expr)
		}
	}

	symNameOutput := generateTempVarName(constants.LOCAL_PREFIX)

	symOutput := ast.MakeArgument(symNameOutput, CurrentFile, LineNo)
	symOutput.SetType(typeCode)
	symOutput.IsSlice = true
	symOutput.Package = ast.CXPackageIndex(pkg.Index)
	symOutput.PreviouslyDeclared = true
	symOutput.TotalSize = types.POINTER_SIZE
	symOutputIdx := prgrm.AddCXArgInArray(symOutput)

	symInput := ast.MakeArgument(symName, CurrentFile, LineNo)
	symInput.SetType(typeCode)
	symInput.IsSlice = true
	symInput.Package = ast.CXPackageIndex(pkg.Index)
	symInput.TotalSize = types.POINTER_SIZE
	symInputIdx := prgrm.AddCXArgInArray(symInput)

	symExprExprCXLine := ast.MakeCXLineExpression(prgrm, CurrentFile, LineNo, LineStr)
	symExpr := ast.MakeAtomicOperatorExpression(prgrm, ast.Natives[constants.OP_IDENTITY])
	symExprAtomicOp, err := prgrm.GetCXAtomicOp(symExpr.Index)
	if err != nil {
		panic(err)
	}

	symExprAtomicOp.Package = ast.CXPackageIndex(pkg.Index)
	symExprAtomicOp.Outputs = append(symExprAtomicOp.Outputs, symOutputIdx)
	symExprAtomicOp.Inputs = append(symExprAtomicOp.Inputs, symInputIdx)

	// marking the output so multidimensional arrays identify the expressions
	result = append(result, *symExprExprCXLine, *symExpr)

	return result
}

// PrimaryStructLiteral creates a struct literal expression.
//
// Input arguments description:
// 	prgrm - a CXProgram that contains all the data and arrays of the program.
//  structName - name of the struct.
//  structFields - fields of the struct.
func PrimaryStructLiteral(prgrm *ast.CXProgram, structName string, structFields []ast.CXExpression) []ast.CXExpression {
	var result []ast.CXExpression

	if pkg, err := prgrm.GetCurrentPackage(); err == nil {
		if strct, err := prgrm.GetStruct(structName, pkg.Name); err == nil {
			for _, expr := range structFields {
				if expr.Type == ast.CX_LINE {
					continue
				}
				cxAtomicOp, err := prgrm.GetCXAtomicOp(expr.Index)
				if err != nil {
					panic(err)
				}

				cxAtomicOpOutputIdx := cxAtomicOp.Outputs[0]
				name := prgrm.CXArgs[cxAtomicOpOutputIdx].Name

				field := ast.MakeArgument(name, CurrentFile, LineNo)
				field.Type = prgrm.CXArgs[cxAtomicOpOutputIdx].Type
				field.PointerTargetType = prgrm.CXArgs[cxAtomicOpOutputIdx].PointerTargetType
				expr.ExpressionType = ast.CXEXPR_STRUCT_LITERAL

				prgrm.CXArgs[cxAtomicOpOutputIdx].Package = ast.CXPackageIndex(pkg.Index)

				if prgrm.CXArgs[cxAtomicOpOutputIdx].StructType == nil {
					prgrm.CXArgs[cxAtomicOpOutputIdx].StructType = strct
				}
				field.StructType = strct

				prgrm.CXArgs[cxAtomicOpOutputIdx].Size = strct.GetStructSize(prgrm)
				prgrm.CXArgs[cxAtomicOpOutputIdx].TotalSize = strct.GetStructSize(prgrm)
				prgrm.CXArgs[cxAtomicOpOutputIdx].Name = structName
				fieldIdx := prgrm.AddCXArgInArray(field)
				prgrm.CXArgs[cxAtomicOpOutputIdx].Fields = append(prgrm.CXArgs[cxAtomicOpOutputIdx].Fields, fieldIdx)

				result = append(result, expr)
			}
		} else {
			panic("type '" + structName + "' does not exist")
		}
	} else {
		panic(err)
	}

	return result
}

// PrimaryStructLiteralExternal creates a struct literal expression after a postfix expression.
//
// Input arguments description:
// 	prgrm - a CXProgram that contains all the data and arrays of the program.
//  importName - name of the import.
//  structName - name of the struct.
//  structFields - fields of the struct.
func PrimaryStructLiteralExternal(prgrm *ast.CXProgram, importName string, structName string, structFields []ast.CXExpression) []ast.CXExpression {
	var result []ast.CXExpression
	if pkg, err := prgrm.GetCurrentPackage(); err == nil {
		if _, err := pkg.GetImport(prgrm, importName); err == nil {
			if strct, err := prgrm.GetStruct(structName, importName); err == nil {
				for _, expr := range structFields {
					if expr.Type == ast.CX_LINE {
						continue
					}
					cxAtomicOp, err := prgrm.GetCXAtomicOp(expr.Index)
					if err != nil {
						panic(err)
					}

					cxAtomicOpOutputIdx := cxAtomicOp.Outputs[0]
					field := ast.MakeArgument("", CurrentFile, LineNo)
					field.SetType(types.IDENTIFIER)
					field.Name = prgrm.CXArgs[cxAtomicOpOutputIdx].Name

					expr.ExpressionType = ast.CXEXPR_STRUCT_LITERAL

					prgrm.CXArgs[cxAtomicOpOutputIdx].Package = ast.CXPackageIndex(pkg.Index)
					// expr.ProgramOutput[0].Program = prgrm

					prgrm.CXArgs[cxAtomicOpOutputIdx].StructType = strct
					prgrm.CXArgs[cxAtomicOpOutputIdx].Size = strct.GetStructSize(prgrm)
					prgrm.CXArgs[cxAtomicOpOutputIdx].TotalSize = strct.GetStructSize(prgrm)
					prgrm.CXArgs[cxAtomicOpOutputIdx].Name = structName
					fieldIdx := prgrm.AddCXArgInArray(field)
					prgrm.CXArgs[cxAtomicOpOutputIdx].Fields = append(prgrm.CXArgs[cxAtomicOpOutputIdx].Fields, fieldIdx)
					result = append(result, expr)
				}
			} else {
				panic("type '" + structName + "' does not exist")
			}
		} else {
			panic(err)
		}
	} else {
		panic(err)
	}

	return result
}

// ArrayLiteralExpression creates an array literal expression.
//
// Input arguments description:
// 	prgrm - a CXProgram that contains all the data and arrays of the program.
//  arraySizes - len(arraySizes) is the length of the array while contents of
// 				 arraySizes define the size of each dimension.
//  typeCode - type code of the array.
//  exprs - the array of expressions composing the array literal expression.
func ArrayLiteralExpression(prgrm *ast.CXProgram, arraySizes []types.Pointer, typeCode types.Code, exprs []ast.CXExpression) []ast.CXExpression {
	var result []ast.CXExpression

	pkg, err := prgrm.GetCurrentPackage()
	if err != nil {
		panic(err)
	}

	symName := generateTempVarName(constants.LOCAL_PREFIX)

	arrVarExprCXLine := ast.MakeCXLineExpression(prgrm, CurrentFile, LineNo, LineStr)
	arrVarExpr := ast.MakeAtomicOperatorExpression(prgrm, nil)
	arrVarExprAtomicOpIdx := arrVarExpr.Index

	prgrm.CXAtomicOps[arrVarExprAtomicOpIdx].Package = ast.CXPackageIndex(pkg.Index)
	arrVar := ast.MakeArgument(symName, CurrentFile, LineNo)
	arrVar = DeclarationSpecifiers(arrVar, arraySizes, constants.DECL_ARRAY)
	arrVar.SetType(typeCode)
	arrVar.TotalSize = arrVar.Size * TotalLength(arrVar.Lengths)
	arrVar.Package = ast.CXPackageIndex(pkg.Index)
	arrVar.PreviouslyDeclared = true
	arrVarIdx := prgrm.AddCXArgInArray(arrVar)

	prgrm.CXAtomicOps[arrVarExprAtomicOpIdx].AddOutput(prgrm, arrVarIdx)

	result = append(result, *arrVarExprCXLine, *arrVarExpr)

	var endPointsCounter int
	for _, expr := range exprs {
		if expr.Type == ast.CX_LINE {
			continue
		}
		expression, err := prgrm.GetCXAtomicOp(expr.Index)
		if err != nil {
			panic(err)
		}
		expressionOperator := prgrm.GetFunctionFromArray(expression.Operator)

		if expr.IsArrayLiteral() {
			expr.ExpressionType = ast.CXEXPR_UNUSED

			sym := ast.MakeArgument(symName, CurrentFile, LineNo)
			sym.SetType(typeCode)
			sym.Package = ast.CXPackageIndex(pkg.Index)
			sym.PreviouslyDeclared = true

			if sym.Type == types.STR || sym.Type == types.AFF {
				sym.PassBy = constants.PASSBY_REFERENCE
			}

			idxExpr := WritePrimary(prgrm, types.I32, encoder.Serialize(int32(endPointsCounter)), false)
			endPointsCounter++

			idxExprAtomicOp, err := prgrm.GetCXAtomicOpFromExpressions(idxExpr, 0)
			if err != nil {
				panic(err)
			}

			sym.Indexes = append(sym.Indexes, idxExprAtomicOp.Outputs[0])
			sym.DereferenceOperations = append(sym.DereferenceOperations, constants.DEREF_ARRAY)
			sym.Lengths = arraySizes
			sym.TotalSize = sym.Size * TotalLength(sym.Lengths)
			symIdx := prgrm.AddCXArgInArray(sym)

			symExprCXLine := ast.MakeCXLineExpression(prgrm, CurrentFile, LineNo, LineStr)
			symExpr := ast.MakeAtomicOperatorExpression(prgrm, nil)
			symExprAtomicOpIdx := symExpr.Index

			prgrm.CXAtomicOps[symExprAtomicOpIdx].AddOutput(prgrm, symIdx)

			if expressionOperator == nil {
				// then it's a literal
				opIdx := prgrm.AddNativeFunctionInArray(ast.Natives[constants.OP_IDENTITY])
				prgrm.CXAtomicOps[symExprAtomicOpIdx].Operator = opIdx
				prgrm.CXAtomicOps[symExprAtomicOpIdx].Inputs = expression.Outputs
			} else {
				prgrm.CXAtomicOps[symExprAtomicOpIdx].Operator = expression.Operator
				prgrm.CXAtomicOps[symExprAtomicOpIdx].Inputs = expression.Inputs

				// hack to get the correct lengths below
				expression.Outputs = append(expression.Outputs, symIdx)
			}
			result = append(result, *symExprCXLine, *symExpr)
		} else {
			result = append(result, expr)
		}
	}

	symNameOutput := generateTempVarName(constants.LOCAL_PREFIX)

	symOutput := ast.MakeArgument(symNameOutput, CurrentFile, LineNo)
	symOutput.SetType(typeCode)
	// symOutput.Lengths = append(symOutput.Lengths, arrSizes[len(arrSizes)-1])
	symOutput.Lengths = arraySizes
	symOutput.Package = ast.CXPackageIndex(pkg.Index)
	symOutput.PreviouslyDeclared = true
	symOutput.TotalSize = symOutput.Size * TotalLength(symOutput.Lengths)
	symOutputIdx := prgrm.AddCXArgInArray(symOutput)

	symInput := ast.MakeArgument(symName, CurrentFile, LineNo)
	symInput.SetType(typeCode)
	// symInput.Lengths = append(symInput.Lengths, arrSizes[len(arrSizes)-1])
	symInput.Lengths = arraySizes
	symInput.Package = ast.CXPackageIndex(pkg.Index)
	symInput.PreviouslyDeclared = true
	symInput.TotalSize = symInput.Size * TotalLength(symInput.Lengths)
	symInputIdx := prgrm.AddCXArgInArray(symInput)

	symExprCXLine := ast.MakeCXLineExpression(prgrm, CurrentFile, LineNo, LineStr)
	symExpr := ast.MakeAtomicOperatorExpression(prgrm, ast.Natives[constants.OP_IDENTITY])
	symExprAtomicOpIdx := symExpr.Index
	prgrm.CXAtomicOps[symExprAtomicOpIdx].Package = ast.CXPackageIndex(pkg.Index)
	prgrm.CXAtomicOps[symExprAtomicOpIdx].AddOutput(prgrm, symOutputIdx)
	prgrm.CXAtomicOps[symExprAtomicOpIdx].AddInput(prgrm, symInputIdx)

	// symOutput.SynonymousTo = symInput.Name

	// marking the output so multidimensional arrays identify the expressions
	symExpr.ExpressionType = ast.CXEXPR_ARRAY_LITERAL
	result = append(result, *symExprCXLine, *symExpr)

	return result
}

// StructLiteralFields creates an initial expression for the
// struct literal expression.
//
// Input arguments description:
// 	prgrm - a CXProgram that contains all the data and arrays of the program.
//	structName - name of the struct.
func StructLiteralFields(prgrm *ast.CXProgram, structName string) ast.CXExpression {
	pkg, err := prgrm.GetCurrentPackage()
	if err != nil {
		panic(err)
	}

	arg := ast.MakeArgument("", CurrentFile, LineNo)
	arg.SetType(types.IDENTIFIER)
	arg.Name = structName
	arg.Package = ast.CXPackageIndex(pkg.Index)
	argIdx := prgrm.AddCXArgInArray(arg)

	expr := ast.MakeAtomicOperatorExpression(prgrm, nil)
	expression, err := prgrm.GetCXAtomicOp(expr.Index)
	if err != nil {
		panic(err)
	}
	expression.AddOutput(prgrm, argIdx)
	expression.Package = ast.CXPackageIndex(pkg.Index)
	return *expr
}
