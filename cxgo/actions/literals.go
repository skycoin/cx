package actions

import (
	. "github.com/skycoin/cx/cx"
	"github.com/skycoin/skycoin/src/cipher/encoder"
)

func SliceLiteralExpression(typSpec int, exprs []*CXExpression) []*CXExpression {
	var result []*CXExpression

	pkg, err := PRGRM.GetCurrentPackage()
	if err != nil {
		panic(err)
	}

	symName := MakeGenSym(LOCAL_PREFIX)

	// adding the declaration
	slcVarExpr := MakeExpression(nil, CurrentFile, LineNo)
	slcVarExpr.Package = pkg
	slcVar := MakeArgument(symName, CurrentFile, LineNo)
	slcVar = DeclarationSpecifiers(slcVar, 0, DECL_SLICE)
	slcVar.AddType(TypeNames[typSpec])

	// slcVar.IsSlice = true

	slcVar.TotalSize = TYPE_POINTER_SIZE

	slcVarExpr.Outputs = append(slcVarExpr.Outputs, slcVar)
	slcVar.Package = pkg
	slcVar.PreviouslyDeclared = true

	result = append(result, slcVarExpr)

	var endPointsCounter int
	for _, expr := range exprs {
		if expr.IsArrayLiteral {
			symInp := MakeArgument(symName, CurrentFile, LineNo).AddType(TypeNames[typSpec])
			symInp.Package = pkg
			symOut := MakeArgument(symName, CurrentFile, LineNo).AddType(TypeNames[typSpec])
			symOut.Package = pkg

			// symOut.IsSlice = true
			// symInp.IsSlice = true

			endPointsCounter++

			symExpr := MakeExpression(nil, CurrentFile, LineNo)
			symExpr.Package = pkg
			// symExpr.Outputs = append(symExpr.Outputs, symOut)
			symExpr.AddOutput(symOut)

			if expr.Operator == nil {
				// then it's a literal
				symExpr.Operator = Natives[OP_APPEND]

				symExpr.Inputs = nil
				symExpr.Inputs = append(symExpr.Inputs, symInp)
				symExpr.Inputs = append(symExpr.Inputs, expr.Outputs...)
			} else {
				symExpr.Operator = expr.Operator

				symExpr.Inputs = nil
				symExpr.Inputs = append(symExpr.Inputs, symInp)
				symExpr.Inputs = append(symExpr.Inputs, expr.Inputs...)

				// hack to get the correct lengths below
				expr.Outputs = append(expr.Outputs, symInp)
			}

			// result = append(result, expr)
			result = append(result, symExpr)

			symInp.TotalSize = TYPE_POINTER_SIZE
			symOut.TotalSize = TYPE_POINTER_SIZE
		} else {
			result = append(result, expr)
		}
	}

	symNameOutput := MakeGenSym(LOCAL_PREFIX)

	symOutput := MakeArgument(symNameOutput, CurrentFile, LineNo).AddType(TypeNames[typSpec])
	// symOutput.PassBy = PASSBY_REFERENCE
	symOutput.IsSlice = true
	symOutput.Package = pkg
	symOutput.PreviouslyDeclared = true

	// symOutput.DeclarationSpecifiers = append(symOutput.DeclarationSpecifiers, DECL_ARRAY)

	symInput := MakeArgument(symName, CurrentFile, LineNo).AddType(TypeNames[typSpec])
	// symInput.DereferenceOperations = append(symInput.DereferenceOperations, DEREF_POINTER)
	symInput.IsSlice = true
	symInput.Package = pkg
	// symInput.PassBy = PASSBY_REFERENCE

	symInput.TotalSize = TYPE_POINTER_SIZE
	symOutput.TotalSize = TYPE_POINTER_SIZE

	symExpr := MakeExpression(Natives[OP_IDENTITY], CurrentFile, LineNo)
	symExpr.Package = pkg
	symExpr.Outputs = append(symExpr.Outputs, symOutput)
	symExpr.Inputs = append(symExpr.Inputs, symInput)

	// symExpr.IsArrayLiteral = true

	// symOutput.SynonymousTo = symInput.Name

	// marking the output so multidimensional arrays identify the expressions
	result = append(result, symExpr)

	return result
}

func PrimaryStructLiteral(ident string, strctFlds []*CXExpression) []*CXExpression {
	var result []*CXExpression

	if pkg, err := PRGRM.GetCurrentPackage(); err == nil {
		if strct, err := PRGRM.GetStruct(ident, pkg.Name); err == nil {
			for _, expr := range strctFlds {
				name := expr.Outputs[0].Name

				fld := MakeArgument(name, CurrentFile, LineNo)
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

func PrimaryStructLiteralExternal(impName string, ident string, strctFlds []*CXExpression) []*CXExpression {
	var result []*CXExpression
	if pkg, err := PRGRM.GetCurrentPackage(); err == nil {
		if _, err := pkg.GetImport(impName); err == nil {
			if strct, err := PRGRM.GetStruct(ident, impName); err == nil {
				for _, expr := range strctFlds {
					fld := MakeArgument("", CurrentFile, LineNo)
					fld.AddType(TypeNames[TYPE_IDENTIFIER])
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

func ArrayLiteralExpression(arrSize int, typSpec int, exprs []*CXExpression) []*CXExpression {
	var result []*CXExpression

	pkg, err := PRGRM.GetCurrentPackage()
	if err != nil {
		panic(err)
	}

	symName := MakeGenSym(LOCAL_PREFIX)

	arrVarExpr := MakeExpression(nil, CurrentFile, LineNo)
	arrVarExpr.Package = pkg
	arrVar := MakeArgument(symName, CurrentFile, LineNo)
	arrVar = DeclarationSpecifiers(arrVar, arrSize, DECL_ARRAY)
	arrVar.AddType(TypeNames[typSpec])
	arrVar.TotalSize = arrVar.Size * TotalLength(arrVar.Lengths)

	arrVarExpr.Outputs = append(arrVarExpr.Outputs, arrVar)
	arrVar.Package = pkg
	arrVar.PreviouslyDeclared = true

	result = append(result, arrVarExpr)

	var endPointsCounter int
	for _, expr := range exprs {
		if expr.IsArrayLiteral {
			expr.IsArrayLiteral = false

			sym := MakeArgument(symName, CurrentFile, LineNo).AddType(TypeNames[typSpec])
			sym.Package = pkg
			sym.PreviouslyDeclared = true

			if sym.Type == TYPE_STR || sym.Type == TYPE_AFF {
				sym.PassBy = PASSBY_REFERENCE
			}

			idxExpr := WritePrimary(TYPE_I32, encoder.Serialize(int32(endPointsCounter)), false)
			endPointsCounter++

			sym.Indexes = append(sym.Indexes, idxExpr[0].Outputs[0])
			sym.DereferenceOperations = append(sym.DereferenceOperations, DEREF_ARRAY)

			symExpr := MakeExpression(nil, CurrentFile, LineNo)
			symExpr.Outputs = append(symExpr.Outputs, sym)

			if expr.Operator == nil {
				// then it's a literal
				symExpr.Operator = Natives[OP_IDENTITY]
				symExpr.Inputs = expr.Outputs
			} else {
				symExpr.Operator = expr.Operator
				symExpr.Inputs = expr.Inputs

				// hack to get the correct lengths below
				expr.Outputs = append(expr.Outputs, sym)
			}

			result = append(result, symExpr)

			sym.Lengths = append(expr.Outputs[0].Lengths, arrSize)
			sym.TotalSize = sym.Size * TotalLength(sym.Lengths)
		} else {
			result = append(result, expr)
		}
	}

	symNameOutput := MakeGenSym(LOCAL_PREFIX)

	symOutput := MakeArgument(symNameOutput, CurrentFile, LineNo).AddType(TypeNames[typSpec])
	symOutput.Lengths = append(symOutput.Lengths, arrSize)
	symOutput.Package = pkg
	symOutput.PreviouslyDeclared = true
	symOutput.TotalSize = symOutput.Size * TotalLength(symOutput.Lengths)

	symInput := MakeArgument(symName, CurrentFile, LineNo).AddType(TypeNames[typSpec])
	symInput.Lengths = append(symInput.Lengths, arrSize)
	symInput.Package = pkg
	symInput.PreviouslyDeclared = true
	symInput.TotalSize = symInput.Size * TotalLength(symInput.Lengths)

	symExpr := MakeExpression(Natives[OP_IDENTITY], CurrentFile, LineNo)
	symExpr.Package = pkg
	symExpr.Outputs = append(symExpr.Outputs, symOutput)
	symExpr.Inputs = append(symExpr.Inputs, symInput)

	// symOutput.SynonymousTo = symInput.Name

	// marking the output so multidimensional arrays identify the expressions
	symExpr.IsArrayLiteral = true
	result = append(result, symExpr)

	return result
}
