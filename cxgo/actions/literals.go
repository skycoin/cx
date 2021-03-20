package actions

import (
	"github.com/skycoin/skycoin/src/cipher/encoder"

	"github.com/skycoin/cx/cx"
)

// SliceLiteralExpression handles literal expressions by converting it to a series of `append` expressions.
func SliceLiteralExpression(typSpec int, exprs []*cxcore.CXExpression) []*cxcore.CXExpression {
	var result []*cxcore.CXExpression

	pkg, err := PRGRM.GetCurrentPackage()
	if err != nil {
		panic(err)
	}

	symName := cxcore.MakeGenSym(cxcore.LOCAL_PREFIX)

	// adding the declaration
	slcVarExpr := cxcore.MakeExpression(nil, CurrentFile, LineNo)
	slcVarExpr.Package = pkg
	slcVar := cxcore.MakeArgument(symName, CurrentFile, LineNo)
	slcVar.AddType(cxcore.TypeNames[typSpec])
	slcVar = DeclarationSpecifiers(slcVar, []int{0}, cxcore.DECL_SLICE)

	slcVar.TotalSize = cxcore.TYPE_POINTER_SIZE

	slcVarExpr.Outputs = append(slcVarExpr.Outputs, slcVar)
	slcVar.Package = pkg
	slcVar.PreviouslyDeclared = true

	result = append(result, slcVarExpr)

	var endPointsCounter int
	for _, expr := range exprs {
		if expr.IsArrayLiteral {
			symInp := cxcore.MakeArgument(symName, CurrentFile, LineNo).AddType(cxcore.TypeNames[typSpec])
			symInp.Package = pkg
			symOut := cxcore.MakeArgument(symName, CurrentFile, LineNo).AddType(cxcore.TypeNames[typSpec])
			symOut.Package = pkg

			endPointsCounter++

			symExpr := cxcore.MakeExpression(nil, CurrentFile, LineNo)
			symExpr.Package = pkg
			symExpr.AddOutput(symOut)

			if expr.Operator == nil {
				// then it's a literal
				symExpr.Operator = cxcore.Natives[cxcore.OP_APPEND]

				symExpr.Inputs = nil
				symExpr.Inputs = append(symExpr.Inputs, symInp)
				symExpr.Inputs = append(symExpr.Inputs, expr.Outputs...)
			} else {
				// We need to create a temporary variable to hold the result of the
				// nested expressions. Then use that variable as part of the slice literal.
				out := cxcore.MakeArgument(cxcore.MakeGenSym(cxcore.LOCAL_PREFIX), expr.FileName, expr.FileLine)
				outArg := getOutputType(expr)
				out.AddType(cxcore.TypeNames[outArg.Type])
				out.CustomType = outArg.CustomType
				out.Size = outArg.Size
				out.TotalSize = cxcore.GetSize(outArg)
				out.PreviouslyDeclared = true

				expr.Outputs = nil
				expr.AddOutput(out)
				result = append(result, expr)

				symExpr.Operator = cxcore.Natives[cxcore.OP_APPEND]

				symExpr.Inputs = nil
				symExpr.Inputs = append(symExpr.Inputs, symInp)
				symExpr.Inputs = append(symExpr.Inputs, out)
			}

			result = append(result, symExpr)

			symInp.TotalSize = cxcore.TYPE_POINTER_SIZE
			symOut.TotalSize = cxcore.TYPE_POINTER_SIZE
		} else {
			result = append(result, expr)
		}
		expr.IsArrayLiteral = false
	}

	symNameOutput := cxcore.MakeGenSym(cxcore.LOCAL_PREFIX)

	symOutput := cxcore.MakeArgument(symNameOutput, CurrentFile, LineNo).AddType(cxcore.TypeNames[typSpec])
	symOutput.IsSlice = true
	symOutput.Package = pkg
	symOutput.PreviouslyDeclared = true

	symInput := cxcore.MakeArgument(symName, CurrentFile, LineNo).AddType(cxcore.TypeNames[typSpec])
	symInput.IsSlice = true
	symInput.Package = pkg

	symInput.TotalSize = cxcore.TYPE_POINTER_SIZE
	symOutput.TotalSize = cxcore.TYPE_POINTER_SIZE

	symExpr := cxcore.MakeExpression(cxcore.Natives[cxcore.OP_IDENTITY], CurrentFile, LineNo)
	symExpr.Package = pkg
	symExpr.Outputs = append(symExpr.Outputs, symOutput)
	symExpr.Inputs = append(symExpr.Inputs, symInput)

	// marking the output so multidimensional arrays identify the expressions
	result = append(result, symExpr)

	return result
}

func PrimaryStructLiteral(ident string, strctFlds []*cxcore.CXExpression) []*cxcore.CXExpression {
	var result []*cxcore.CXExpression

	if pkg, err := PRGRM.GetCurrentPackage(); err == nil {
		if strct, err := PRGRM.GetStruct(ident, pkg.Name); err == nil {
			for _, expr := range strctFlds {
				name := expr.Outputs[0].Name

				fld := cxcore.MakeArgument(name, CurrentFile, LineNo)
				fld.Type = expr.Outputs[0].Type

				expr.IsStructLiteral = true

				expr.Outputs[0].Package = pkg
				// expr.Outputs[0].Program = PRGRM

				if expr.Outputs[0].CustomType == nil {
					expr.Outputs[0].CustomType = strct
				}
				// expr.Outputs[0].CustomType = strct
				fld.CustomType = strct

				expr.Outputs[0].Size = strct.Size
				expr.Outputs[0].TotalSize = strct.Size
				expr.Outputs[0].Name = ident
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

func PrimaryStructLiteralExternal(impName string, ident string, strctFlds []*cxcore.CXExpression) []*cxcore.CXExpression {
	var result []*cxcore.CXExpression
	if pkg, err := PRGRM.GetCurrentPackage(); err == nil {
		if _, err := pkg.GetImport(impName); err == nil {
			if strct, err := PRGRM.GetStruct(ident, impName); err == nil {
				for _, expr := range strctFlds {
					fld := cxcore.MakeArgument("", CurrentFile, LineNo)
					fld.AddType(cxcore.TypeNames[cxcore.TYPE_IDENTIFIER])
					fld.Name = expr.Outputs[0].Name

					expr.IsStructLiteral = true

					expr.Outputs[0].Package = pkg
					// expr.Outputs[0].Program = PRGRM

					expr.Outputs[0].CustomType = strct
					expr.Outputs[0].Size = strct.Size
					expr.Outputs[0].TotalSize = strct.Size
					expr.Outputs[0].Name = ident
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

func ArrayLiteralExpression(arrSizes []int, typSpec int, exprs []*cxcore.CXExpression) []*cxcore.CXExpression {
	var result []*cxcore.CXExpression

	pkg, err := PRGRM.GetCurrentPackage()
	if err != nil {
		panic(err)
	}

	symName := cxcore.MakeGenSym(cxcore.LOCAL_PREFIX)

	arrVarExpr := cxcore.MakeExpression(nil, CurrentFile, LineNo)
	arrVarExpr.Package = pkg
	arrVar := cxcore.MakeArgument(symName, CurrentFile, LineNo)
	arrVar = DeclarationSpecifiers(arrVar, arrSizes, cxcore.DECL_ARRAY)
	arrVar.AddType(cxcore.TypeNames[typSpec])
	arrVar.TotalSize = arrVar.Size * TotalLength(arrVar.Lengths)

	arrVarExpr.Outputs = append(arrVarExpr.Outputs, arrVar)
	arrVar.Package = pkg
	arrVar.PreviouslyDeclared = true

	result = append(result, arrVarExpr)

	var endPointsCounter int
	for _, expr := range exprs {
		if expr.IsArrayLiteral {
			expr.IsArrayLiteral = false

			sym := cxcore.MakeArgument(symName, CurrentFile, LineNo).AddType(cxcore.TypeNames[typSpec])
			sym.Package = pkg
			sym.PreviouslyDeclared = true

			if sym.Type == cxcore.TYPE_STR || sym.Type == cxcore.TYPE_AFF {
				sym.PassBy = cxcore.PASSBY_REFERENCE
			}

			idxExpr := WritePrimary(cxcore.TYPE_I32, encoder.Serialize(int32(endPointsCounter)), false)
			endPointsCounter++

			sym.Indexes = append(sym.Indexes, idxExpr[0].Outputs[0])
			sym.DereferenceOperations = append(sym.DereferenceOperations, cxcore.DEREF_ARRAY)

			symExpr := cxcore.MakeExpression(nil, CurrentFile, LineNo)
			symExpr.Outputs = append(symExpr.Outputs, sym)

			if expr.Operator == nil {
				// then it's a literal
				symExpr.Operator = cxcore.Natives[cxcore.OP_IDENTITY]
				symExpr.Inputs = expr.Outputs
			} else {
				symExpr.Operator = expr.Operator
				symExpr.Inputs = expr.Inputs

				// hack to get the correct lengths below
				expr.Outputs = append(expr.Outputs, sym)
			}

			result = append(result, symExpr)

			// sym.Lengths = append(expr.Outputs[0].Lengths, arrSizes[len(arrSizes)-1])
			sym.Lengths = arrSizes
			sym.TotalSize = sym.Size * TotalLength(sym.Lengths)
		} else {
			result = append(result, expr)
		}
	}

	symNameOutput := cxcore.MakeGenSym(cxcore.LOCAL_PREFIX)

	symOutput := cxcore.MakeArgument(symNameOutput, CurrentFile, LineNo).AddType(cxcore.TypeNames[typSpec])
	// symOutput.Lengths = append(symOutput.Lengths, arrSizes[len(arrSizes)-1])
	symOutput.Lengths = arrSizes
	symOutput.Package = pkg
	symOutput.PreviouslyDeclared = true
	symOutput.TotalSize = symOutput.Size * TotalLength(symOutput.Lengths)

	symInput := cxcore.MakeArgument(symName, CurrentFile, LineNo).AddType(cxcore.TypeNames[typSpec])
	// symInput.Lengths = append(symInput.Lengths, arrSizes[len(arrSizes)-1])
	symInput.Lengths = arrSizes
	symInput.Package = pkg
	symInput.PreviouslyDeclared = true
	symInput.TotalSize = symInput.Size * TotalLength(symInput.Lengths)

	symExpr := cxcore.MakeExpression(cxcore.Natives[cxcore.OP_IDENTITY], CurrentFile, LineNo)
	symExpr.Package = pkg
	symExpr.Outputs = append(symExpr.Outputs, symOutput)
	symExpr.Inputs = append(symExpr.Inputs, symInput)

	// symOutput.SynonymousTo = symInput.Name

	// marking the output so multidimensional arrays identify the expressions
	symExpr.IsArrayLiteral = true
	result = append(result, symExpr)

	return result
}
