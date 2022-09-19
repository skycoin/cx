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
	slcVar.Package = ast.CXPackageIndex(pkg.Index)
	slcVar.PreviouslyDeclared = true

	typeSig := ast.GetCXTypeSignatureRepresentationOfCXArg(prgrm, slcVar)
	typeSigIdx := prgrm.AddCXTypeSignatureInArray(typeSig)
	prgrm.CXAtomicOps[slcVarExprAtomicOpIdx].AddOutput(prgrm, typeSigIdx)

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

			symOut := ast.MakeArgument(symName, CurrentFile, LineNo).SetType(typeCode)
			symOut.Package = ast.CXPackageIndex(pkg.Index)

			endPointsCounter++

			symExprExprCXLine := ast.MakeCXLineExpression(prgrm, CurrentFile, LineNo, LineStr)
			symExpr := ast.MakeAtomicOperatorExpression(prgrm, nil)
			symExprAtomicOpIdx := symExpr.Index

			prgrm.CXAtomicOps[symExprAtomicOpIdx].Package = ast.CXPackageIndex(pkg.Index)
			typeSig := ast.GetCXTypeSignatureRepresentationOfCXArg(prgrm, symOut)
			typeSigIdx := prgrm.AddCXTypeSignatureInArray(typeSig)
			prgrm.CXAtomicOps[symExprAtomicOpIdx].AddOutput(prgrm, typeSigIdx)

			if exprAtomicOpOperator == nil {
				// then it's a literal
				opIdx := prgrm.AddNativeFunctionInArray(ast.Natives[constants.OP_APPEND])
				prgrm.CXAtomicOps[symExprAtomicOpIdx].Operator = opIdx

				prgrm.CXAtomicOps[symExprAtomicOpIdx].Inputs = nil
				typeSig := ast.GetCXTypeSignatureRepresentationOfCXArg(prgrm, symInp)
				typeSigIdx := prgrm.AddCXTypeSignatureInArray(typeSig)
				prgrm.CXAtomicOps[symExprAtomicOpIdx].AddInput(prgrm, typeSigIdx)
				prgrm.CXAtomicOps[symExprAtomicOpIdx].Inputs.Fields = append(prgrm.CXAtomicOps[symExprAtomicOpIdx].Inputs.Fields, exprAtomicOp.Outputs.Fields...)
			} else {
				// We need to create a temporary variable to hold the result of the
				// nested expressions. Then use that variable as part of the slice literal.

				outTypeSig := getOutputType(prgrm, &expr)
				var outArg *ast.CXArgument
				var outTypeSigIdx ast.CXTypeSignatureIndex
				if outTypeSig.Type == ast.TYPE_CXARGUMENT_DEPRECATE {
					outArg = prgrm.GetCXArgFromArray(ast.CXArgumentIndex(outTypeSig.Meta))

					out := ast.MakeArgument(generateTempVarName(constants.LOCAL_PREFIX), exprCXLine.FileName, exprCXLine.LineNumber)
					out.SetType(outArg.Type)
					out.PointerTargetType = outArg.PointerTargetType
					out.StructType = outArg.StructType
					out.Size = outArg.Size
					out.PreviouslyDeclared = true
					out.Package = ast.CXPackageIndex(pkg.Index)

					typeSig := ast.GetCXTypeSignatureRepresentationOfCXArg(prgrm, out)
					outTypeSigIdx = prgrm.AddCXTypeSignatureInArray(typeSig)
				} else if outTypeSig.Type == ast.TYPE_ATOMIC || outTypeSig.Type == ast.TYPE_POINTER_ATOMIC {
					var newTypeSig ast.CXTypeSignature
					newTypeSig = *outTypeSig
					newTypeSig.Name = generateTempVarName(constants.LOCAL_PREFIX)
					newTypeSig.Package = ast.CXPackageIndex(pkg.Index)
					outTypeSigIdx = prgrm.AddCXTypeSignatureInArray(&newTypeSig)
				} else {
					panic("type is not known")
				}

				exprAtomicOp.Outputs = nil

				exprAtomicOp.AddOutput(prgrm, outTypeSigIdx)

				result = append(result, expr)
				opIdx := prgrm.AddNativeFunctionInArray(ast.Natives[constants.OP_APPEND])
				prgrm.CXAtomicOps[symExprAtomicOpIdx].Operator = opIdx

				prgrm.CXAtomicOps[symExprAtomicOpIdx].Inputs = nil
				typeSig = ast.GetCXTypeSignatureRepresentationOfCXArg(prgrm, symInp)
				typeSigIdx = prgrm.AddCXTypeSignatureInArray(typeSig)
				prgrm.CXAtomicOps[symExprAtomicOpIdx].AddInput(prgrm, typeSigIdx)
				prgrm.CXAtomicOps[symExprAtomicOpIdx].AddInput(prgrm, outTypeSigIdx)
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

	symInput := ast.MakeArgument(symName, CurrentFile, LineNo)
	symInput.SetType(typeCode)
	symInput.IsSlice = true
	symInput.Package = ast.CXPackageIndex(pkg.Index)

	symExprExprCXLine := ast.MakeCXLineExpression(prgrm, CurrentFile, LineNo, LineStr)
	symExpr := ast.MakeAtomicOperatorExpression(prgrm, ast.Natives[constants.OP_IDENTITY])
	symExprAtomicOp, err := prgrm.GetCXAtomicOp(symExpr.Index)
	if err != nil {
		panic(err)
	}

	symExprAtomicOp.Package = ast.CXPackageIndex(pkg.Index)
	typeSig = ast.GetCXTypeSignatureRepresentationOfCXArg(prgrm, symOutput)
	typeSigIdx = prgrm.AddCXTypeSignatureInArray(typeSig)
	symExprAtomicOp.AddOutput(prgrm, typeSigIdx)
	typeSig = ast.GetCXTypeSignatureRepresentationOfCXArg(prgrm, symInput)
	typeSigIdx = prgrm.AddCXTypeSignatureInArray(typeSig)
	symExprAtomicOp.AddInput(prgrm, typeSigIdx)

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

				cxAtomicOpOutputTypeSig := prgrm.GetCXTypeSignatureFromArray(cxAtomicOp.GetOutputs(prgrm)[0])
				if cxAtomicOpOutputTypeSig.Type == ast.TYPE_CXARGUMENT_DEPRECATE {
					cxAtomicOpOutputIdx := cxAtomicOpOutputTypeSig.Meta
					name := prgrm.CXArgs[cxAtomicOpOutputIdx].Name

					field := ast.MakeArgument(name, CurrentFile, LineNo)
					field.Type = prgrm.CXArgs[cxAtomicOpOutputIdx].Type
					field.PointerTargetType = prgrm.CXArgs[cxAtomicOpOutputIdx].PointerTargetType
					field.StructType = strct

					prgrm.CXArgs[cxAtomicOpOutputIdx].Package = ast.CXPackageIndex(pkg.Index)

					if prgrm.CXArgs[cxAtomicOpOutputIdx].StructType == nil {
						prgrm.CXArgs[cxAtomicOpOutputIdx].StructType = strct
					}

					prgrm.CXArgs[cxAtomicOpOutputIdx].Size = strct.GetStructSize(prgrm)
					prgrm.CXArgs[cxAtomicOpOutputIdx].Name = structName
					cxAtomicOpOutputTypeSig.Name = structName

					fieldIdx := prgrm.AddCXArgInArray(field)
					prgrm.CXArgs[cxAtomicOpOutputIdx].Fields = append(prgrm.CXArgs[cxAtomicOpOutputIdx].Fields, fieldIdx)
				} else if cxAtomicOpOutputTypeSig.Type == ast.TYPE_ATOMIC {
					// TODO: give proper change when we implement type_structs
					// Looks like we have to convert the arg to type cx arg deprecate again
					newCXArg := &ast.CXArgument{ArgDetails: &ast.CXArgumentDebug{}}
					newCXArg.Type = types.Code(cxAtomicOpOutputTypeSig.Meta)
					newCXArg.Package = ast.CXPackageIndex(pkg.Index)
					newCXArg.Offset = cxAtomicOpOutputTypeSig.Offset

					newCXArg.StructType = strct
					newCXArg.Size = strct.GetStructSize(prgrm)
					newCXArg.Name = structName

					field := ast.MakeArgument(cxAtomicOpOutputTypeSig.Name, CurrentFile, LineNo)
					field.Type = types.Code(cxAtomicOpOutputTypeSig.Meta)
					field.StructType = strct
					fieldIdx := prgrm.AddCXArgInArray(field)

					newCXArg.Fields = append(newCXArg.Fields, fieldIdx)
					newCXArgIdx := prgrm.AddCXArgInArray(newCXArg)

					cxAtomicOpOutputTypeSig.Name = structName
					cxAtomicOpOutputTypeSig.Type = ast.TYPE_CXARGUMENT_DEPRECATE
					cxAtomicOpOutputTypeSig.Meta = int(newCXArgIdx)
				} else if cxAtomicOpOutputTypeSig.Type == ast.TYPE_POINTER_ATOMIC {
					// TODO: give proper change when we implement type_structs
					// Looks like we have to convert the arg to type cx arg deprecate again
					newCXArg := &ast.CXArgument{ArgDetails: &ast.CXArgumentDebug{}}
					newCXArg.Type = types.Code(cxAtomicOpOutputTypeSig.Meta)
					newCXArg.Package = ast.CXPackageIndex(pkg.Index)
					newCXArg.Offset = cxAtomicOpOutputTypeSig.Offset
					newCXArg.DeclarationSpecifiers = []int{constants.DECL_BASIC, constants.DEREF_POINTER}

					newCXArg.StructType = strct
					newCXArg.Size = strct.GetStructSize(prgrm)
					newCXArg.Name = structName

					field := ast.MakeArgument(cxAtomicOpOutputTypeSig.Name, CurrentFile, LineNo)
					field.Type = types.Code(cxAtomicOpOutputTypeSig.Meta)
					field.StructType = strct
					fieldIdx := prgrm.AddCXArgInArray(field)

					newCXArg.Fields = append(newCXArg.Fields, fieldIdx)
					newCXArgIdx := prgrm.AddCXArgInArray(newCXArg)

					cxAtomicOpOutputTypeSig.Name = structName
					cxAtomicOpOutputTypeSig.Type = ast.TYPE_CXARGUMENT_DEPRECATE
					cxAtomicOpOutputTypeSig.Meta = int(newCXArgIdx)
				} else {
					panic("type is not known")
				}

				expr.ExpressionType = ast.CXEXPR_STRUCT_LITERAL
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

					cxAtomicOpOutputTypeSig := prgrm.GetCXTypeSignatureFromArray(cxAtomicOp.GetOutputs(prgrm)[0])
					if cxAtomicOpOutputTypeSig.Type == ast.TYPE_CXARGUMENT_DEPRECATE {
						cxAtomicOpOutputIdx := cxAtomicOpOutputTypeSig.Meta
						field := ast.MakeArgument("", CurrentFile, LineNo)
						field.SetType(types.IDENTIFIER)
						field.Name = prgrm.CXArgs[cxAtomicOpOutputIdx].Name

						expr.ExpressionType = ast.CXEXPR_STRUCT_LITERAL

						prgrm.CXArgs[cxAtomicOpOutputIdx].Package = ast.CXPackageIndex(pkg.Index)
						// expr.ProgramOutput[0].Program = prgrm

						prgrm.CXArgs[cxAtomicOpOutputIdx].StructType = strct
						prgrm.CXArgs[cxAtomicOpOutputIdx].Size = strct.GetStructSize(prgrm)
						prgrm.CXArgs[cxAtomicOpOutputIdx].Name = structName
						cxAtomicOpOutputTypeSig.Name = structName

						fieldIdx := prgrm.AddCXArgInArray(field)
						prgrm.CXArgs[cxAtomicOpOutputIdx].Fields = append(prgrm.CXArgs[cxAtomicOpOutputIdx].Fields, fieldIdx)
					} else if cxAtomicOpOutputTypeSig.Type == ast.TYPE_ATOMIC || cxAtomicOpOutputTypeSig.Type == ast.TYPE_POINTER_ATOMIC {
						panic("type signature is type atomic")
					} else {
						panic("type is not known")
					}

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
	arrVar.SetType(typeCode)
	arrVar = DeclarationSpecifiers(arrVar, arraySizes, constants.DECL_ARRAY)
	arrVar.Package = ast.CXPackageIndex(pkg.Index)
	arrVar.PreviouslyDeclared = true

	typeSig := ast.GetCXTypeSignatureRepresentationOfCXArg(prgrm, arrVar)
	typeSigIdx := prgrm.AddCXTypeSignatureInArray(typeSig)
	prgrm.CXAtomicOps[arrVarExprAtomicOpIdx].AddOutput(prgrm, typeSigIdx)

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

			idxArg := WritePrimary(prgrm, types.I32, encoder.Serialize(int32(endPointsCounter)), false)
			// index is always i32
			// hence, always an atomic type

			indexTypeSignature := &ast.CXTypeSignature{
				Name:    idxArg.Name,
				Type:    ast.TYPE_ATOMIC,
				Meta:    int(idxArg.Type),
				Offset:  idxArg.Offset, // important for this
				Package: idxArg.Package,
			}

			indexTypeSignatureIdx := prgrm.AddCXTypeSignatureInArray(indexTypeSignature)
			endPointsCounter++

			sym.Indexes = append(sym.Indexes, indexTypeSignatureIdx)
			sym.DereferenceOperations = append(sym.DereferenceOperations, constants.DEREF_ARRAY)
			sym.Lengths = arraySizes

			symExprCXLine := ast.MakeCXLineExpression(prgrm, CurrentFile, LineNo, LineStr)
			symExpr := ast.MakeAtomicOperatorExpression(prgrm, nil)
			symExprAtomicOpIdx := symExpr.Index

			typeSig := ast.GetCXTypeSignatureRepresentationOfCXArg(prgrm, sym)
			typeSigIdx := prgrm.AddCXTypeSignatureInArray(typeSig)
			prgrm.CXAtomicOps[symExprAtomicOpIdx].AddOutput(prgrm, typeSigIdx)

			if expressionOperator == nil {
				// then it's a literal
				opIdx := prgrm.AddNativeFunctionInArray(ast.Natives[constants.OP_IDENTITY])
				prgrm.CXAtomicOps[symExprAtomicOpIdx].Operator = opIdx
				prgrm.CXAtomicOps[symExprAtomicOpIdx].Inputs = expression.Outputs
			} else {
				prgrm.CXAtomicOps[symExprAtomicOpIdx].Operator = expression.Operator
				prgrm.CXAtomicOps[symExprAtomicOpIdx].Inputs = expression.Inputs

				typeSig := ast.GetCXTypeSignatureRepresentationOfCXArg(prgrm, sym)
				typeSigIdx := prgrm.AddCXTypeSignatureInArray(typeSig)
				// hack to get the correct lengths below
				expression.AddOutput(prgrm, typeSigIdx)
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

	symInput := ast.MakeArgument(symName, CurrentFile, LineNo)
	symInput.SetType(typeCode)
	// symInput.Lengths = append(symInput.Lengths, arrSizes[len(arrSizes)-1])
	symInput.Lengths = arraySizes
	symInput.Package = ast.CXPackageIndex(pkg.Index)
	symInput.PreviouslyDeclared = true

	symExprCXLine := ast.MakeCXLineExpression(prgrm, CurrentFile, LineNo, LineStr)
	symExpr := ast.MakeAtomicOperatorExpression(prgrm, ast.Natives[constants.OP_IDENTITY])
	symExprAtomicOpIdx := symExpr.Index
	prgrm.CXAtomicOps[symExprAtomicOpIdx].Package = ast.CXPackageIndex(pkg.Index)

	typeSig = ast.GetCXTypeSignatureRepresentationOfCXArg(prgrm, symOutput)
	typeSigIdx = prgrm.AddCXTypeSignatureInArray(typeSig)
	prgrm.CXAtomicOps[symExprAtomicOpIdx].AddOutput(prgrm, typeSigIdx)

	typeSig = ast.GetCXTypeSignatureRepresentationOfCXArg(prgrm, symInput)
	typeSigIdx = prgrm.AddCXTypeSignatureInArray(typeSig)
	prgrm.CXAtomicOps[symExprAtomicOpIdx].AddInput(prgrm, typeSigIdx)

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

	expr := ast.MakeAtomicOperatorExpression(prgrm, nil)
	expression, err := prgrm.GetCXAtomicOp(expr.Index)
	if err != nil {
		panic(err)
	}

	typeSig := ast.GetCXTypeSignatureRepresentationOfCXArg(prgrm, arg)
	typeSigIdx := prgrm.AddCXTypeSignatureInArray(typeSig)
	expression.AddOutput(prgrm, typeSigIdx)

	expression.Package = ast.CXPackageIndex(pkg.Index)
	return *expr
}
