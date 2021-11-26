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
	slcVarExpr := ast.MakeExpression(nil, CurrentFile, LineNo)
	slcVarExpr.Package = pkg
	slcVar := ast.MakeArgument(symName, CurrentFile, LineNo)
	slcVar.AddType(typeCode)
	slcVar = DeclarationSpecifiers(slcVar, []types.Pointer{0}, constants.DECL_SLICE)

	slcVar.TotalSize = types.POINTER_SIZE

	slcVarExpr.Outputs = append(slcVarExpr.Outputs, slcVar)
	slcVar.ArgDetails.Package = pkg
	slcVar.PreviouslyDeclared = true

	result = append(result, slcVarExpr)

	var endPointsCounter int
	for _, expr := range exprs {
		if expr.IsArrayLiteral() {
			symInp := ast.MakeArgument(symName, CurrentFile, LineNo).AddType(typeCode)
			symInp.ArgDetails.Package = pkg
			symOut := ast.MakeArgument(symName, CurrentFile, LineNo).AddType(typeCode)
			symOut.ArgDetails.Package = pkg

			endPointsCounter++

			symExpr := ast.MakeExpression(nil, CurrentFile, LineNo)
			symExpr.Package = pkg
			symExpr.AddOutput(symOut)

			if expr.Operator == nil {
				// then it's a literal
				symExpr.Operator = ast.Natives[constants.OP_APPEND]

				symExpr.Inputs = nil
				symExpr.Inputs = append(symExpr.Inputs, symInp)
				symExpr.Inputs = append(symExpr.Inputs, expr.Outputs...)
			} else {
				// We need to create a temporary variable to hold the result of the
				// nested expressions. Then use that variable as part of the slice literal.
				out := ast.MakeArgument(MakeGenSym(constants.LOCAL_PREFIX), expr.FileName, expr.FileLine)
				outArg := getOutputType(expr)
				out.AddType(outArg.Type)
				out.PointerTargetType = outArg.PointerTargetType
				out.StructType = outArg.StructType
				out.Size = outArg.Size
				out.TotalSize = ast.GetSize(outArg)
				out.PreviouslyDeclared = true

				expr.Outputs = nil
				expr.AddOutput(out)
				result = append(result, expr)

				symExpr.Operator = ast.Natives[constants.OP_APPEND]

				symExpr.Inputs = nil
				symExpr.Inputs = append(symExpr.Inputs, symInp)
				symExpr.Inputs = append(symExpr.Inputs, out)
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
	symOutput.ArgDetails.Package = pkg
	symOutput.PreviouslyDeclared = true

	symInput := ast.MakeArgument(symName, CurrentFile, LineNo).AddType(typeCode)
	symInput.IsSlice = true
	symInput.ArgDetails.Package = pkg

	symInput.TotalSize = types.POINTER_SIZE
	symOutput.TotalSize = types.POINTER_SIZE

	symExpr := ast.MakeExpression(ast.Natives[constants.OP_IDENTITY], CurrentFile, LineNo)
	symExpr.Package = pkg
	symExpr.Outputs = append(symExpr.Outputs, symOutput)
	symExpr.Inputs = append(symExpr.Inputs, symInput)

	// marking the output so multidimensional arrays identify the expressions
	result = append(result, symExpr)

	return result
}

func PrimaryStructLiteral(prgrm *ast.CXProgram, ident string, strctFlds []*ast.CXExpression) []*ast.CXExpression {
	var result []*ast.CXExpression

	if pkg, err := prgrm.GetCurrentPackage(); err == nil {
		if strct, err := prgrm.GetStruct(ident, pkg.Name); err == nil {
			for _, expr := range strctFlds {
				name := expr.Outputs[0].ArgDetails.Name

				fld := ast.MakeArgument(name, CurrentFile, LineNo)
				fld.Type = expr.Outputs[0].Type
				fld.PointerTargetType = expr.Outputs[0].PointerTargetType
				expr.ExpressionType = ast.CXEXPR_STRUCT_LITERAL

				expr.Outputs[0].ArgDetails.Package = pkg
				// expr.ProgramOutput[0].Program = prgrm

				if expr.Outputs[0].StructType == nil {
					expr.Outputs[0].StructType = strct
				}
				// expr.ProgramOutput[0].StructType = strct
				fld.StructType = strct

				expr.Outputs[0].Size = strct.Size
				expr.Outputs[0].TotalSize = strct.Size
				expr.Outputs[0].ArgDetails.Name = ident
				expr.Outputs[0].Fields = append(expr.Outputs[0].Fields, fld)
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
					fld := ast.MakeArgument("", CurrentFile, LineNo)
					fld.AddType(types.IDENTIFIER)
					fld.ArgDetails.Name = expr.Outputs[0].ArgDetails.Name

					expr.ExpressionType = ast.CXEXPR_STRUCT_LITERAL

					expr.Outputs[0].ArgDetails.Package = pkg
					// expr.ProgramOutput[0].Program = prgrm

					expr.Outputs[0].StructType = strct
					expr.Outputs[0].Size = strct.Size
					expr.Outputs[0].TotalSize = strct.Size
					expr.Outputs[0].ArgDetails.Name = ident
					expr.Outputs[0].Fields = append(expr.Outputs[0].Fields, fld)
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

	arrVarExpr := ast.MakeExpression(nil, CurrentFile, LineNo)
	arrVarExpr.Package = pkg
	arrVar := ast.MakeArgument(symName, CurrentFile, LineNo)
	arrVar = DeclarationSpecifiers(arrVar, arrSizes, constants.DECL_ARRAY)
	arrVar.AddType(typeCode)
	arrVar.TotalSize = arrVar.Size * TotalLength(arrVar.Lengths)

	arrVarExpr.Outputs = append(arrVarExpr.Outputs, arrVar)
	arrVar.ArgDetails.Package = pkg
	arrVar.PreviouslyDeclared = true

	result = append(result, arrVarExpr)

	var endPointsCounter int
	for _, expr := range exprs {
		if expr.IsArrayLiteral() {
			expr.ExpressionType = ast.CXEXPR_UNUSED

			sym := ast.MakeArgument(symName, CurrentFile, LineNo).AddType(typeCode)
			sym.ArgDetails.Package = pkg
			sym.PreviouslyDeclared = true

			if sym.Type == types.STR || sym.Type == types.AFF {
				sym.PassBy = constants.PASSBY_REFERENCE
			}

			idxExpr := WritePrimary(prgrm, types.I32, encoder.Serialize(int32(endPointsCounter)), false)
			endPointsCounter++

			sym.Indexes = append(sym.Indexes, idxExpr[0].Outputs[0])
			sym.DereferenceOperations = append(sym.DereferenceOperations, constants.DEREF_ARRAY)

			symExpr := ast.MakeExpression(nil, CurrentFile, LineNo)
			symExpr.Outputs = append(symExpr.Outputs, sym)

			if expr.Operator == nil {
				// then it's a literal
				symExpr.Operator = ast.Natives[constants.OP_IDENTITY]
				symExpr.Inputs = expr.Outputs
			} else {
				symExpr.Operator = expr.Operator
				symExpr.Inputs = expr.Inputs

				// hack to get the correct lengths below
				expr.Outputs = append(expr.Outputs, sym)
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
	symOutput.ArgDetails.Package = pkg
	symOutput.PreviouslyDeclared = true
	symOutput.TotalSize = symOutput.Size * TotalLength(symOutput.Lengths)

	symInput := ast.MakeArgument(symName, CurrentFile, LineNo).AddType(typeCode)
	// symInput.Lengths = append(symInput.Lengths, arrSizes[len(arrSizes)-1])
	symInput.Lengths = arrSizes
	symInput.ArgDetails.Package = pkg
	symInput.PreviouslyDeclared = true
	symInput.TotalSize = symInput.Size * TotalLength(symInput.Lengths)

	symExpr := ast.MakeExpression(ast.Natives[constants.OP_IDENTITY], CurrentFile, LineNo)
	symExpr.Package = pkg
	symExpr.Outputs = append(symExpr.Outputs, symOutput)
	symExpr.Inputs = append(symExpr.Inputs, symInput)

	// symOutput.SynonymousTo = symInput.Name

	// marking the output so multidimensional arrays identify the expressions
	symExpr.ExpressionType = ast.CXEXPR_ARRAY_LITERAL
	result = append(result, symExpr)

	return result
}
