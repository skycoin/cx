package actions

import (
	"github.com/skycoin/cx/cx/ast"
	"github.com/skycoin/cx/cx/constants"
	"github.com/skycoin/cx/cx/types"
	"github.com/skycoin/skycoin/src/cipher/encoder"
)

// SliceLiteralExpression handles literal expressions by converting it to a series of `append` expressions.
func SliceLiteralExpression(prgrm *ast.CXProgram, typeCode types.Code, exprs []ast.CXExpression) []ast.CXExpression {
	var result []ast.CXExpression

	pkg, err := prgrm.GetCurrentPackage()
	if err != nil {
		panic(err)
	}

	symName := MakeGenSym(constants.LOCAL_PREFIX)

	slcVarExprCXLine := ast.MakeCXLineExpression(prgrm, CurrentFile, LineNo, LineStr)
	// adding the declaration
	slcVarExpr := ast.MakeAtomicOperatorExpression(prgrm, nil)
	slcVarExprAtomicOp, _, _, err := prgrm.GetOperation(slcVarExpr)
	if err != nil {
		panic(err)
	}

	slcVarExprAtomicOp.Package = ast.CXPackageIndex(pkg.Index)
	slcVar := ast.MakeArgument(symName, CurrentFile, LineNo)
	slcVar.AddType(typeCode)
	slcVar = DeclarationSpecifiers(slcVar, []types.Pointer{0}, constants.DECL_SLICE)
	slcVar.TotalSize = types.POINTER_SIZE
	slcVar.Package = ast.CXPackageIndex(pkg.Index)
	slcVar.PreviouslyDeclared = true

	slcVarIdx := prgrm.AddCXArgInArray(slcVar)
	slcVarExprAtomicOp.AddOutput(prgrm, slcVarIdx)

	result = append(result, *slcVarExprCXLine, *slcVarExpr)

	var endPointsCounter int
	for i, expr := range exprs {
		if expr.Type == ast.CX_LINE {
			continue
		}
		exprAtomicOp, _, _, err := prgrm.GetOperation(&expr)
		if err != nil {
			panic(err)
		}
		exprAtomicOpOperator := prgrm.GetFunctionFromArray(exprAtomicOp.Operator)

		exprCXLine, _ := prgrm.GetPreviousCXLine(exprs, i)

		if expr.IsArrayLiteral() {
			symInp := ast.MakeArgument(symName, CurrentFile, LineNo).AddType(typeCode)
			symInp.Package = ast.CXPackageIndex(pkg.Index)
			symInp.TotalSize = types.POINTER_SIZE
			symInpIdx := prgrm.AddCXArgInArray(symInp)

			symOut := ast.MakeArgument(symName, CurrentFile, LineNo).AddType(typeCode)
			symOut.Package = ast.CXPackageIndex(pkg.Index)
			symOut.TotalSize = types.POINTER_SIZE
			symOutIdx := prgrm.AddCXArgInArray(symOut)

			endPointsCounter++

			symExprExprCXLine := ast.MakeCXLineExpression(prgrm, CurrentFile, LineNo, LineStr)
			symExpr := ast.MakeAtomicOperatorExpression(prgrm, nil)
			symExprAtomicOp, _, _, err := prgrm.GetOperation(symExpr)
			if err != nil {
				panic(err)
			}

			symExprAtomicOp.Package = ast.CXPackageIndex(pkg.Index)
			symExprAtomicOp.AddOutput(prgrm, symOutIdx)

			if exprAtomicOpOperator == nil {
				// then it's a literal
				opIdx := prgrm.AddFunctionInArray(ast.Natives[constants.OP_APPEND])
				symExprAtomicOp.Operator = opIdx

				symExprAtomicOp.Inputs = nil
				symExprAtomicOp.Inputs = append(symExprAtomicOp.Inputs, symInpIdx)
				symExprAtomicOp.Inputs = append(symExprAtomicOp.Inputs, exprAtomicOp.Outputs...)
			} else {
				// We need to create a temporary variable to hold the result of the
				// nested expressions. Then use that variable as part of the slice literal.
				out := ast.MakeArgument(MakeGenSym(constants.LOCAL_PREFIX), exprCXLine.FileName, exprCXLine.LineNumber)
				outArg := getOutputType(prgrm, &expr)
				out.AddType(outArg.Type)
				out.PointerTargetType = outArg.PointerTargetType
				out.StructType = outArg.StructType
				out.Size = outArg.Size
				out.TotalSize = ast.GetSize(prgrm, outArg)
				out.PreviouslyDeclared = true
				outIdx := prgrm.AddCXArgInArray(out)

				exprAtomicOp.Outputs = nil
				exprAtomicOp.AddOutput(prgrm, outIdx)

				result = append(result, expr)
				opIdx := prgrm.AddFunctionInArray(ast.Natives[constants.OP_APPEND])
				symExprAtomicOp.Operator = opIdx

				symExprAtomicOp.Inputs = nil
				symExprAtomicOp.Inputs = append(symExprAtomicOp.Inputs, symInpIdx)
				symExprAtomicOp.Inputs = append(symExprAtomicOp.Inputs, outIdx)
			}
			result = append(result, *symExprExprCXLine, *symExpr)
		} else {
			result = append(result, expr)
		}
		expr.ExpressionType = ast.CXEXPR_UNUSED
	}

	symNameOutput := MakeGenSym(constants.LOCAL_PREFIX)

	symOutput := ast.MakeArgument(symNameOutput, CurrentFile, LineNo).AddType(typeCode)
	symOutput.IsSlice = true
	symOutput.Package = ast.CXPackageIndex(pkg.Index)
	symOutput.PreviouslyDeclared = true
	symOutput.TotalSize = types.POINTER_SIZE
	symOutputIdx := prgrm.AddCXArgInArray(symOutput)

	symInput := ast.MakeArgument(symName, CurrentFile, LineNo).AddType(typeCode)
	symInput.IsSlice = true
	symInput.Package = ast.CXPackageIndex(pkg.Index)
	symInput.TotalSize = types.POINTER_SIZE
	symInputIdx := prgrm.AddCXArgInArray(symInput)

	symExprExprCXLine := ast.MakeCXLineExpression(prgrm, CurrentFile, LineNo, LineStr)
	symExpr := ast.MakeAtomicOperatorExpression(prgrm, ast.Natives[constants.OP_IDENTITY])
	symExprAtomicOp, _, _, err := prgrm.GetOperation(symExpr)
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

func PrimaryStructLiteral(prgrm *ast.CXProgram, ident string, strctFlds []ast.CXExpression) []ast.CXExpression {
	var result []ast.CXExpression

	if pkg, err := prgrm.GetCurrentPackage(); err == nil {
		if strct, err := prgrm.GetStruct(ident, pkg.Name); err == nil {
			for _, expr := range strctFlds {
				if expr.Type == ast.CX_LINE {
					continue
				}
				cxAtomicOp, _, _, err := prgrm.GetOperation(&expr)
				if err != nil {
					panic(err)
				}

				cxAtomicOpOutput := prgrm.GetCXArgFromArray(cxAtomicOp.Outputs[0])
				name := cxAtomicOpOutput.Name

				fld := ast.MakeArgument(name, CurrentFile, LineNo)
				fld.Type = cxAtomicOpOutput.Type
				fld.PointerTargetType = cxAtomicOpOutput.PointerTargetType
				expr.ExpressionType = ast.CXEXPR_STRUCT_LITERAL

				cxAtomicOpOutput.Package = ast.CXPackageIndex(pkg.Index)
				// expr.ProgramOutput[0].Program = prgrm

				if cxAtomicOpOutput.StructType == nil {
					cxAtomicOpOutput.StructType = strct
				}
				// expr.ProgramOutput[0].StructType = strct
				fld.StructType = strct

				cxAtomicOpOutput.Size = strct.Size
				cxAtomicOpOutput.TotalSize = strct.Size
				cxAtomicOpOutput.Name = ident
				fldIdx := prgrm.AddCXArgInArray(fld)
				cxAtomicOpOutput.Fields = append(cxAtomicOpOutput.Fields, fldIdx)
				result = append(result, expr)
			}
		} else {
			panic("type '" + ident + "' does not exist")
		}
	} else {
		panic(err)
	}

	return result
}

func PrimaryStructLiteralExternal(prgrm *ast.CXProgram, impName string, ident string, strctFlds []ast.CXExpression) []ast.CXExpression {
	var result []ast.CXExpression
	if pkg, err := prgrm.GetCurrentPackage(); err == nil {
		if _, err := pkg.GetImport(prgrm, impName); err == nil {
			if strct, err := prgrm.GetStruct(ident, impName); err == nil {
				for _, expr := range strctFlds {
					if expr.Type == ast.CX_LINE {
						continue
					}
					cxAtomicOp, _, _, err := prgrm.GetOperation(&expr)
					if err != nil {
						panic(err)
					}

					cxAtomicOpOutput := prgrm.GetCXArgFromArray(cxAtomicOp.Outputs[0])
					fld := ast.MakeArgument("", CurrentFile, LineNo)
					fld.AddType(types.IDENTIFIER)
					fld.Name = cxAtomicOpOutput.Name

					expr.ExpressionType = ast.CXEXPR_STRUCT_LITERAL

					cxAtomicOpOutput.Package = ast.CXPackageIndex(pkg.Index)
					// expr.ProgramOutput[0].Program = prgrm

					cxAtomicOpOutput.StructType = strct
					cxAtomicOpOutput.Size = strct.Size
					cxAtomicOpOutput.TotalSize = strct.Size
					cxAtomicOpOutput.Name = ident
					fldIdx := prgrm.AddCXArgInArray(fld)
					cxAtomicOpOutput.Fields = append(cxAtomicOpOutput.Fields, fldIdx)
					result = append(result, expr)
				}
			} else {
				panic("type '" + ident + "' does not exist")
			}
		} else {
			panic(err)
		}
	} else {
		panic(err)
	}

	return result
}

func ArrayLiteralExpression(prgrm *ast.CXProgram, arrSizes []types.Pointer, typeCode types.Code, exprs []ast.CXExpression) []ast.CXExpression {
	var result []ast.CXExpression

	pkg, err := prgrm.GetCurrentPackage()
	if err != nil {
		panic(err)
	}

	symName := MakeGenSym(constants.LOCAL_PREFIX)

	arrVarExprCXLine := ast.MakeCXLineExpression(prgrm, CurrentFile, LineNo, LineStr)
	arrVarExpr := ast.MakeAtomicOperatorExpression(prgrm, nil)
	arrVarExprAtomicOp, _, _, err := prgrm.GetOperation(arrVarExpr)
	if err != nil {
		panic(err)
	}

	arrVarExprAtomicOp.Package = ast.CXPackageIndex(pkg.Index)
	arrVar := ast.MakeArgument(symName, CurrentFile, LineNo)
	arrVar = DeclarationSpecifiers(arrVar, arrSizes, constants.DECL_ARRAY)
	arrVar.AddType(typeCode)
	arrVar.TotalSize = arrVar.Size * TotalLength(arrVar.Lengths)
	arrVar.Package = ast.CXPackageIndex(pkg.Index)
	arrVar.PreviouslyDeclared = true
	arrVarIdx := prgrm.AddCXArgInArray(arrVar)

	arrVarExprAtomicOp.AddOutput(prgrm, arrVarIdx)

	result = append(result, *arrVarExprCXLine, *arrVarExpr)

	var endPointsCounter int
	for _, expr := range exprs {
		if expr.Type == ast.CX_LINE {
			continue
		}
		exprAtomicOp, _, _, err := prgrm.GetOperation(&expr)
		if err != nil {
			panic(err)
		}
		exprAtomicOpOperator := prgrm.GetFunctionFromArray(exprAtomicOp.Operator)

		if expr.IsArrayLiteral() {
			expr.ExpressionType = ast.CXEXPR_UNUSED

			sym := ast.MakeArgument(symName, CurrentFile, LineNo).AddType(typeCode)
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
			indexIdx := idxExprAtomicOp.Outputs[0]
			sym.Indexes = append(sym.Indexes, indexIdx)
			sym.DereferenceOperations = append(sym.DereferenceOperations, constants.DEREF_ARRAY)
			symIdx := prgrm.AddCXArgInArray(sym)

			symExprCXLine := ast.MakeCXLineExpression(prgrm, CurrentFile, LineNo, LineStr)
			symExpr := ast.MakeAtomicOperatorExpression(prgrm, nil)
			symExprAtomicOp, _, _, err := prgrm.GetOperation(symExpr)
			if err != nil {
				panic(err)
			}

			symExprAtomicOp.AddOutput(prgrm, symIdx)

			if exprAtomicOpOperator == nil {
				// then it's a literal
				opIdx := prgrm.AddFunctionInArray(ast.Natives[constants.OP_IDENTITY])
				symExprAtomicOp.Operator = opIdx
				symExprAtomicOp.Inputs = exprAtomicOp.Outputs
			} else {
				symExprAtomicOp.Operator = exprAtomicOp.Operator
				symExprAtomicOp.Inputs = exprAtomicOp.Inputs

				// hack to get the correct lengths below
				exprAtomicOp.Outputs = append(exprAtomicOp.Outputs, symIdx)
			}
			result = append(result, *symExprCXLine, *symExpr)

			// sym.Lengths = append(expr.ProgramOutput[0].Lengths, arrSizes[len(arrSizes)-1])
			sym = prgrm.GetCXArgFromArray(symIdx)
			sym.Lengths = arrSizes
			sym.TotalSize = sym.Size * TotalLength(sym.Lengths)
		} else {
			result = append(result, expr)
		}
	}

	symNameOutput := MakeGenSym(constants.LOCAL_PREFIX)

	symOutput := ast.MakeArgument(symNameOutput, CurrentFile, LineNo).AddType(typeCode)
	// symOutput.Lengths = append(symOutput.Lengths, arrSizes[len(arrSizes)-1])
	symOutput.Lengths = arrSizes
	symOutput.Package = ast.CXPackageIndex(pkg.Index)
	symOutput.PreviouslyDeclared = true
	symOutput.TotalSize = symOutput.Size * TotalLength(symOutput.Lengths)

	symInput := ast.MakeArgument(symName, CurrentFile, LineNo).AddType(typeCode)
	// symInput.Lengths = append(symInput.Lengths, arrSizes[len(arrSizes)-1])
	symInput.Lengths = arrSizes
	symInput.Package = ast.CXPackageIndex(pkg.Index)
	symInput.PreviouslyDeclared = true
	symInput.TotalSize = symInput.Size * TotalLength(symInput.Lengths)

	symExprCXLine := ast.MakeCXLineExpression(prgrm, CurrentFile, LineNo, LineStr)
	symExpr := ast.MakeAtomicOperatorExpression(prgrm, ast.Natives[constants.OP_IDENTITY])
	symExprAtomicOp, _, _, err := prgrm.GetOperation(symExpr)
	if err != nil {
		panic(err)
	}
	symExprAtomicOp.Package = ast.CXPackageIndex(pkg.Index)
	symOutputIdx := prgrm.AddCXArgInArray(symOutput)
	symExprAtomicOp.Outputs = append(symExprAtomicOp.Outputs, symOutputIdx)
	symInputIdx := prgrm.AddCXArgInArray(symInput)
	symExprAtomicOp.Inputs = append(symExprAtomicOp.Inputs, symInputIdx)

	// symOutput.SynonymousTo = symInput.Name

	// marking the output so multidimensional arrays identify the expressions
	symExpr.ExpressionType = ast.CXEXPR_ARRAY_LITERAL
	result = append(result, *symExprCXLine, *symExpr)

	return result
}
