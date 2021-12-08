package actions

import (
	"github.com/skycoin/cx/cx/ast"
	"github.com/skycoin/cx/cx/constants"
	"github.com/skycoin/cx/cx/types"
	"github.com/skycoin/skycoin/src/cipher/encoder"
)

// SliceLiteralExpression handles literal expressions by converting it to a series of `append` expressions.
func SliceLiteralExpression(prgrm *ast.CXProgram, typeCode types.Code, exprs []*ast.CXExpression) []*ast.CXExpression {
	var result []*ast.CXExpression

	pkg, err := prgrm.GetCurrentPackage()
	if err != nil {
		panic(err)
	}

	symName := MakeGenSym(constants.LOCAL_PREFIX)

	// adding the declaration
	slcVarExpr := ast.MakeExpression(prgrm, nil, CurrentFile, LineNo)
	slcVarExprAtomicOp, _, _, err := prgrm.GetOperation(slcVarExpr)
	if err != nil {
		panic(err)
	}

	slcVarExprAtomicOp.Package = pkg
	slcVar := ast.MakeArgument(symName, CurrentFile, LineNo)
	slcVar.AddType(typeCode)
	slcVar = DeclarationSpecifiers(slcVar, []types.Pointer{0}, constants.DECL_SLICE)

	slcVar.TotalSize = types.POINTER_SIZE

	slcVarExprAtomicOp.AddOutput(slcVar)
	slcVar.Package = pkg
	slcVar.PreviouslyDeclared = true

	result = append(result, slcVarExpr)

	var endPointsCounter int
	for _, expr := range exprs {
		exprAtomicOp, _, _, err := prgrm.GetOperation(expr)
		if err != nil {
			panic(err)
		}

		if expr.IsArrayLiteral() {
			symInp := ast.MakeArgument(symName, CurrentFile, LineNo).AddType(typeCode)
			symInp.Package = pkg
			symOut := ast.MakeArgument(symName, CurrentFile, LineNo).AddType(typeCode)
			symOut.Package = pkg

			endPointsCounter++

			symExpr := ast.MakeExpression(prgrm, nil, CurrentFile, LineNo)
			symExprAtomicOp, _, _, err := prgrm.GetOperation(symExpr)
			if err != nil {
				panic(err)
			}

			symExprAtomicOp.Package = pkg
			symExprAtomicOp.AddOutput(symOut)

			if exprAtomicOp.Operator == nil {
				// then it's a literal
				symExprAtomicOp.Operator = ast.Natives[constants.OP_APPEND]

				symExprAtomicOp.Inputs = nil
				symExprAtomicOp.Inputs = append(symExprAtomicOp.Inputs, symInp)
				symExprAtomicOp.Inputs = append(symExprAtomicOp.Inputs, exprAtomicOp.Outputs...)
			} else {
				// We need to create a temporary variable to hold the result of the
				// nested expressions. Then use that variable as part of the slice literal.
				out := ast.MakeArgument(MakeGenSym(constants.LOCAL_PREFIX), expr.FileName, expr.FileLine)
				outArg := getOutputType(prgrm, expr)
				out.AddType(outArg.Type)
				out.PointerTargetType = outArg.PointerTargetType
				out.StructType = outArg.StructType
				out.Size = outArg.Size
				out.TotalSize = ast.GetSize(outArg)
				out.PreviouslyDeclared = true

				exprAtomicOp.Outputs = nil
				exprAtomicOp.AddOutput(out)
				result = append(result, expr)

				symExprAtomicOp.Operator = ast.Natives[constants.OP_APPEND]

				symExprAtomicOp.Inputs = nil
				symExprAtomicOp.Inputs = append(symExprAtomicOp.Inputs, symInp)
				symExprAtomicOp.Inputs = append(symExprAtomicOp.Inputs, out)
			}

			result = append(result, symExpr)

			symInp.TotalSize = types.POINTER_SIZE
			symOut.TotalSize = types.POINTER_SIZE
		} else {
			result = append(result, expr)
		}
		expr.ExpressionType = ast.CXEXPR_UNUSED
	}

	symNameOutput := MakeGenSym(constants.LOCAL_PREFIX)

	symOutput := ast.MakeArgument(symNameOutput, CurrentFile, LineNo).AddType(typeCode)
	symOutput.IsSlice = true
	symOutput.Package = pkg
	symOutput.PreviouslyDeclared = true

	symInput := ast.MakeArgument(symName, CurrentFile, LineNo).AddType(typeCode)
	symInput.IsSlice = true
	symInput.Package = pkg

	symInput.TotalSize = types.POINTER_SIZE
	symOutput.TotalSize = types.POINTER_SIZE

	symExpr := ast.MakeExpression(prgrm, ast.Natives[constants.OP_IDENTITY], CurrentFile, LineNo)
	symExprAtomicOp, _, _, err := prgrm.GetOperation(symExpr)
	if err != nil {
		panic(err)
	}

	symExprAtomicOp.Package = pkg
	symExprAtomicOp.Outputs = append(symExprAtomicOp.Outputs, symOutput)
	symExprAtomicOp.Inputs = append(symExprAtomicOp.Inputs, symInput)

	// marking the output so multidimensional arrays identify the expressions
	result = append(result, symExpr)

	return result
}

func PrimaryStructLiteral(prgrm *ast.CXProgram, ident string, strctFlds []*ast.CXExpression) []*ast.CXExpression {
	var result []*ast.CXExpression

	if pkg, err := prgrm.GetCurrentPackage(); err == nil {
		if strct, err := prgrm.GetStruct(ident, pkg.Name); err == nil {
			for _, expr := range strctFlds {
				cxAtomicOp, _, _, err := prgrm.GetOperation(expr)
				if err != nil {
					panic(err)
				}
				name := cxAtomicOp.Outputs[0].Name

				fld := ast.MakeArgument(name, CurrentFile, LineNo)
				fld.Type = cxAtomicOp.Outputs[0].Type
				fld.PointerTargetType = cxAtomicOp.Outputs[0].PointerTargetType
				expr.ExpressionType = ast.CXEXPR_STRUCT_LITERAL

				cxAtomicOp.Outputs[0].Package = pkg
				// expr.ProgramOutput[0].Program = prgrm

				if cxAtomicOp.Outputs[0].StructType == nil {
					cxAtomicOp.Outputs[0].StructType = strct
				}
				// expr.ProgramOutput[0].StructType = strct
				fld.StructType = strct

				cxAtomicOp.Outputs[0].Size = strct.Size
				cxAtomicOp.Outputs[0].TotalSize = strct.Size
				cxAtomicOp.Outputs[0].Name = ident
				cxAtomicOp.Outputs[0].Fields = append(cxAtomicOp.Outputs[0].Fields, fld)
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

func PrimaryStructLiteralExternal(prgrm *ast.CXProgram, impName string, ident string, strctFlds []*ast.CXExpression) []*ast.CXExpression {
	var result []*ast.CXExpression
	if pkg, err := prgrm.GetCurrentPackage(); err == nil {
		if _, err := pkg.GetImport(impName); err == nil {
			if strct, err := prgrm.GetStruct(ident, impName); err == nil {
				for _, expr := range strctFlds {
					cxAtomicOp, _, _, err := prgrm.GetOperation(expr)
					if err != nil {
						panic(err)
					}

					fld := ast.MakeArgument("", CurrentFile, LineNo)
					fld.AddType(types.IDENTIFIER)
					fld.Name = cxAtomicOp.Outputs[0].Name

					expr.ExpressionType = ast.CXEXPR_STRUCT_LITERAL

					cxAtomicOp.Outputs[0].Package = pkg
					// expr.ProgramOutput[0].Program = prgrm

					cxAtomicOp.Outputs[0].StructType = strct
					cxAtomicOp.Outputs[0].Size = strct.Size
					cxAtomicOp.Outputs[0].TotalSize = strct.Size
					cxAtomicOp.Outputs[0].Name = ident
					cxAtomicOp.Outputs[0].Fields = append(cxAtomicOp.Outputs[0].Fields, fld)
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

func ArrayLiteralExpression(prgrm *ast.CXProgram, arrSizes []types.Pointer, typeCode types.Code, exprs []*ast.CXExpression) []*ast.CXExpression {
	var result []*ast.CXExpression

	pkg, err := prgrm.GetCurrentPackage()
	if err != nil {
		panic(err)
	}

	symName := MakeGenSym(constants.LOCAL_PREFIX)

	arrVarExpr := ast.MakeExpression(prgrm, nil, CurrentFile, LineNo)
	arrVarExprAtomicOp, _, _, err := prgrm.GetOperation(arrVarExpr)
	if err != nil {
		panic(err)
	}

	arrVarExprAtomicOp.Package = pkg
	arrVar := ast.MakeArgument(symName, CurrentFile, LineNo)
	arrVar = DeclarationSpecifiers(arrVar, arrSizes, constants.DECL_ARRAY)
	arrVar.AddType(typeCode)
	arrVar.TotalSize = arrVar.Size * TotalLength(arrVar.Lengths)

	arrVarExprAtomicOp.AddOutput(arrVar)
	arrVar.Package = pkg
	arrVar.PreviouslyDeclared = true

	result = append(result, arrVarExpr)

	var endPointsCounter int
	for _, expr := range exprs {
		exprAtomicOp, _, _, err := prgrm.GetOperation(expr)
		if err != nil {
			panic(err)
		}

		if expr.IsArrayLiteral() {
			expr.ExpressionType = ast.CXEXPR_UNUSED

			sym := ast.MakeArgument(symName, CurrentFile, LineNo).AddType(typeCode)
			sym.Package = pkg
			sym.PreviouslyDeclared = true

			if sym.Type == types.STR || sym.Type == types.AFF {
				sym.PassBy = constants.PASSBY_REFERENCE
			}

			idxExpr := WritePrimary(prgrm, types.I32, encoder.Serialize(int32(endPointsCounter)), false)
			endPointsCounter++

			idxExprAtomicOp, _, _, err := prgrm.GetOperation(idxExpr[0])
			if err != nil {
				panic(err)
			}

			sym.Indexes = append(sym.Indexes, idxExprAtomicOp.Outputs[0])
			sym.DereferenceOperations = append(sym.DereferenceOperations, constants.DEREF_ARRAY)

			symExpr := ast.MakeExpression(prgrm, nil, CurrentFile, LineNo)
			symExprAtomicOp, _, _, err := prgrm.GetOperation(symExpr)
			if err != nil {
				panic(err)
			}

			symExprAtomicOp.AddOutput(sym)

			if exprAtomicOp.Operator == nil {
				// then it's a literal
				symExprAtomicOp.Operator = ast.Natives[constants.OP_IDENTITY]
				symExprAtomicOp.Inputs = exprAtomicOp.Outputs
			} else {
				symExprAtomicOp.Operator = exprAtomicOp.Operator
				symExprAtomicOp.Inputs = exprAtomicOp.Inputs

				// hack to get the correct lengths below
				exprAtomicOp.Outputs = append(exprAtomicOp.Outputs, sym)
			}

			result = append(result, symExpr)

			// sym.Lengths = append(expr.ProgramOutput[0].Lengths, arrSizes[len(arrSizes)-1])
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
	symOutput.Package = pkg
	symOutput.PreviouslyDeclared = true
	symOutput.TotalSize = symOutput.Size * TotalLength(symOutput.Lengths)

	symInput := ast.MakeArgument(symName, CurrentFile, LineNo).AddType(typeCode)
	// symInput.Lengths = append(symInput.Lengths, arrSizes[len(arrSizes)-1])
	symInput.Lengths = arrSizes
	symInput.Package = pkg
	symInput.PreviouslyDeclared = true
	symInput.TotalSize = symInput.Size * TotalLength(symInput.Lengths)

	symExpr := ast.MakeExpression(prgrm, ast.Natives[constants.OP_IDENTITY], CurrentFile, LineNo)
	symExprAtomicOp, _, _, err := prgrm.GetOperation(symExpr)
	if err != nil {
		panic(err)
	}
	symExprAtomicOp.Package = pkg
	symExprAtomicOp.Outputs = append(symExprAtomicOp.Outputs, symOutput)
	symExprAtomicOp.Inputs = append(symExprAtomicOp.Inputs, symInput)

	// symOutput.SynonymousTo = symInput.Name

	// marking the output so multidimensional arrays identify the expressions
	symExpr.ExpressionType = ast.CXEXPR_ARRAY_LITERAL
	result = append(result, symExpr)

	return result
}
