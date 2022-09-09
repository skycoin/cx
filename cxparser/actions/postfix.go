package actions

import (
	"fmt"
	"os"

	"github.com/skycoin/cx/cx/ast"
	"github.com/skycoin/cx/cx/constants"
	cxpackages "github.com/skycoin/cx/cx/packages"
	"github.com/skycoin/cx/cx/types"
)

// PostfixExpressionArray handles the postfix expression arrays.
//
// Input arguments description:
// 	prgrm - a CXProgram that contains all the data and arrays of the program.
// 	prevExprs - the previous expressions that compose the postfix expression.
//  postExprs - the array expressions.
func PostfixExpressionArray(prgrm *ast.CXProgram, prevExprs []ast.CXExpression, postExprs []ast.CXExpression) []ast.CXExpression {
	var elt *ast.CXArgument = &ast.CXArgument{}

	prevExpressionIdx := prevExprs[len(prevExprs)-1].Index
	prevExpressionOperator := prgrm.GetFunctionFromArray(prgrm.CXAtomicOps[prevExpressionIdx].Operator)

	prevExprCXLine, _ := prgrm.GetPreviousCXLine(prevExprs, len(prevExprs)-1)

	if prevExpressionOperator != nil && len(prgrm.CXAtomicOps[prevExpressionIdx].GetOutputs(prgrm)) == 0 {
		genName := generateTempVarName(constants.LOCAL_PREFIX)

		prevExpressionOperatorOutputs := prevExpressionOperator.GetOutputs(prgrm)
		prevExpressionOperatorOutputTypeSig := prgrm.GetCXTypeSignatureFromArray(prevExpressionOperatorOutputs[0])

		var prevExpressionOperatorOutputArg *ast.CXArgument = &ast.CXArgument{}
		if prevExpressionOperatorOutputTypeSig.Type == ast.TYPE_CXARGUMENT_DEPRECATE {
			prevExpressionOperatorOutputArg = prgrm.GetCXArgFromArray(ast.CXArgumentIndex(prevExpressionOperatorOutputTypeSig.Meta))
		} else {
			panic("type is not cx argument deprecate\n\n")
		}

		prevExpressionOperatorOutput := prevExpressionOperatorOutputArg
		out := ast.MakeArgument(genName, prevExprCXLine.FileName, prevExprCXLine.LineNumber-1).SetType(prevExpressionOperatorOutput.Type)

		out.DeclarationSpecifiers = prevExpressionOperatorOutput.DeclarationSpecifiers
		out.StructType = prevExpressionOperatorOutput.StructType
		out.Size = prevExpressionOperatorOutput.Size
		out.Lengths = prevExpressionOperatorOutput.Lengths
		out.IsSlice = prevExpressionOperatorOutput.IsSlice
		out.PreviouslyDeclared = true
		out.Package = prevExpressionOperatorOutput.Package
		outIdx := prgrm.AddCXArgInArray(out)

		typeSig := ast.GetCXTypeSignatureRepresentationOfCXArg(prgrm, prgrm.GetCXArgFromArray(outIdx))
		typeSigIdx := prgrm.AddCXTypeSignatureInArray(typeSig)
		prgrm.CXAtomicOps[prevExpressionIdx].AddOutput(prgrm, typeSigIdx)

		inp := ast.MakeArgument(genName, prevExprCXLine.FileName, prevExprCXLine.LineNumber).SetType(prevExpressionOperatorOutput.Type)

		inp.DeclarationSpecifiers = prevExpressionOperatorOutput.DeclarationSpecifiers
		inp.StructType = prevExpressionOperatorOutput.StructType
		inp.Size = prevExpressionOperatorOutput.Size
		inp.Lengths = prevExpressionOperatorOutput.Lengths
		inp.IsSlice = prevExpressionOperatorOutput.IsSlice
		inp.PreviouslyDeclared = true
		inpIdx := prgrm.AddCXArgInArray(inp)

		useExprCXLine := ast.MakeCXLineExpression(prgrm, prevExprCXLine.FileName, prevExprCXLine.LineNumber, prevExprCXLine.LineStr)
		useExpr := ast.MakeAtomicOperatorExpression(prgrm, nil)
		useExpressionIdx := useExpr.Index

		prgrm.CXAtomicOps[useExpressionIdx].Package = prgrm.CXAtomicOps[prevExpressionIdx].Package
		typeSig = ast.GetCXTypeSignatureRepresentationOfCXArg(prgrm, prgrm.GetCXArgFromArray(inpIdx))
		typeSigIdx = prgrm.AddCXTypeSignatureInArray(typeSig)
		prgrm.CXAtomicOps[useExpressionIdx].AddOutput(prgrm, typeSigIdx)

		prevExprs = append(prevExprs, *useExprCXLine, *useExpr)
	}

	prevExpression, err := prgrm.GetCXAtomicOpFromExpressions(prevExprs, len(prevExprs)-1)
	if err != nil {
		panic(err)
	}

	prevExpressionOutputTypeSig := prgrm.GetCXTypeSignatureFromArray(prevExpression.GetOutputs(prgrm)[0])
	var prevExpressionOutputArg *ast.CXArgument = &ast.CXArgument{}
	if prevExpressionOutputTypeSig.Type == ast.TYPE_CXARGUMENT_DEPRECATE {
		prevExpressionOutputArg = prgrm.GetCXArgFromArray(ast.CXArgumentIndex(prevExpressionOutputTypeSig.Meta))
	} else {
		panic("type is not cx argument deprecate\n\n")
	}

	prevExpressionArg := prevExpressionOutputArg
	if len(prevExpressionArg.Fields) > 0 {
		elt = prgrm.GetCXArgFromArray(prevExpressionArg.Fields[len(prevExpressionArg.Fields)-1])
	} else {
		elt = prevExpressionArg
	}

	// elt.IsArray = false
	elt.DereferenceOperations = append(elt.DereferenceOperations, constants.DEREF_ARRAY)
	elt.DeclarationSpecifiers = append(elt.DeclarationSpecifiers, constants.DECL_INDEXING)

	postExpressionIdx := postExprs[len(postExprs)-1].Index
	postExpressionOperator := prgrm.GetFunctionFromArray(prgrm.CXAtomicOps[postExpressionIdx].Operator)
	postExpressionOperatorOutputs := postExpressionOperator.GetOutputs(prgrm)

	if len(prevExpressionArg.Fields) > 0 {
		fldIdx := prevExpressionArg.Fields[len(prevExpressionArg.Fields)-1]

		if postExpressionOperator == nil {
			postExpressionOutputIndex := prgrm.CXAtomicOps[postExpressionIdx].GetOutputs(prgrm)[0]
			prgrm.CXArgs[fldIdx].Indexes = append(prgrm.CXArgs[fldIdx].Indexes, postExpressionOutputIndex)
		} else {
			postExpressionOperatorOutputTypeSig := prgrm.GetCXTypeSignatureFromArray(postExpressionOperatorOutputs[0])
			var postExpressionOperatorOutputArg *ast.CXArgument = &ast.CXArgument{}
			if postExpressionOperatorOutputTypeSig.Type == ast.TYPE_CXARGUMENT_DEPRECATE {
				postExpressionOperatorOutputArg = prgrm.GetCXArgFromArray(ast.CXArgumentIndex(postExpressionOperatorOutputTypeSig.Meta))
			} else {
				panic("type is not cx argument deprecate\n\n")
			}

			sym := ast.MakeArgument(generateTempVarName(constants.LOCAL_PREFIX), CurrentFile, LineNo).SetType(postExpressionOperatorOutputArg.Type)
			sym.Package = prgrm.CXAtomicOps[postExpressionIdx].Package
			sym.PreviouslyDeclared = true
			symIdx := prgrm.AddCXArgInArray(sym)
			typeSig := ast.GetCXTypeSignatureRepresentationOfCXArg(prgrm, prgrm.GetCXArgFromArray(symIdx))
			typeSigIdx := prgrm.AddCXTypeSignatureInArray(typeSig)
			prgrm.CXAtomicOps[postExpressionIdx].AddOutput(prgrm, typeSigIdx)

			prevExprs = append(postExprs, prevExprs...)

			prgrm.CXArgs[fldIdx].Indexes = append(prgrm.CXArgs[fldIdx].Indexes, typeSigIdx)
		}
	} else {
		if len(prgrm.CXAtomicOps[postExpressionIdx].GetOutputs(prgrm)) < 1 {
			// then it's an expression (e.g. i32.add(0, 0))
			// we create a gensym for it

			postExpressionOperatorOutputTypeSig := prgrm.GetCXTypeSignatureFromArray(postExpressionOperatorOutputs[0])
			var typeSigIdx ast.CXTypeSignatureIndex
			if postExpressionOperatorOutputTypeSig.Type == ast.TYPE_CXARGUMENT_DEPRECATE {
				postExpressionOperatorOutputArg := prgrm.GetCXArgFromArray(ast.CXArgumentIndex(postExpressionOperatorOutputTypeSig.Meta))
				idxSym := ast.MakeArgument(generateTempVarName(constants.LOCAL_PREFIX), CurrentFile, LineNo).SetType(postExpressionOperatorOutputArg.Type)
				idxSym.Size = postExpressionOperatorOutputArg.Size

				idxSym.Package = prgrm.CXAtomicOps[postExpressionIdx].Package
				idxSym.PreviouslyDeclared = true

				idxSymIdx := prgrm.AddCXArgInArray(idxSym)
				typeSig := ast.GetCXTypeSignatureRepresentationOfCXArg(prgrm, prgrm.GetCXArgFromArray(idxSymIdx))
				typeSigIdx = prgrm.AddCXTypeSignatureInArray(typeSig)
			} else if postExpressionOperatorOutputTypeSig.Type == ast.TYPE_ATOMIC {
				newTypeSig := &ast.CXTypeSignature{}
				newTypeSig.Name = generateTempVarName(constants.LOCAL_PREFIX)
				newTypeSig.Package = prgrm.CXAtomicOps[postExpressionIdx].Package
				newTypeSig.Type = postExpressionOperatorOutputTypeSig.Type
				newTypeSig.Meta = postExpressionOperatorOutputTypeSig.Meta
				newTypeSig.Offset = postExpressionOperatorOutputTypeSig.Offset
				typeSigIdx = prgrm.AddCXTypeSignatureInArray(newTypeSig)
			} else if postExpressionOperatorOutputTypeSig.Type == ast.TYPE_POINTER_ATOMIC {
				newTypeSig := &ast.CXTypeSignature{}
				newTypeSig.Name = generateTempVarName(constants.LOCAL_PREFIX)
				newTypeSig.Package = prgrm.CXAtomicOps[postExpressionIdx].Package
				newTypeSig.Type = postExpressionOperatorOutputTypeSig.Type
				newTypeSig.Meta = postExpressionOperatorOutputTypeSig.Meta
				newTypeSig.Offset = postExpressionOperatorOutputTypeSig.Offset
				typeSigIdx = prgrm.AddCXTypeSignatureInArray(newTypeSig)
			}

			prgrm.CXAtomicOps[postExpressionIdx].AddOutput(prgrm, typeSigIdx)

			prevExpressionOutputTypeSig := prgrm.GetCXTypeSignatureFromArray(prevExpression.GetOutputs(prgrm)[0])
			var prevExpressionOutputArg *ast.CXArgument = &ast.CXArgument{}
			if prevExpressionOutputTypeSig.Type == ast.TYPE_CXARGUMENT_DEPRECATE {
				prevExpressionOutputArg = prgrm.GetCXArgFromArray(ast.CXArgumentIndex(prevExpressionOutputTypeSig.Meta))
			} else {
				panic("type is not cx argument deprecate\n\n")
			}

			prevExpressionOutput := prevExpressionOutputArg
			prevExpressionOutput.Indexes = append(prevExpressionOutput.Indexes, typeSigIdx)

			// we push the index expression
			prevExprs = append(postExprs, prevExprs...)
		} else {
			prevOutsTypeSig := prgrm.GetCXTypeSignatureFromArray(prevExpression.GetOutputs(prgrm)[0])
			if prevOutsTypeSig.Type == ast.TYPE_CXARGUMENT_DEPRECATE {
				prevOutsIdx := ast.CXArgumentIndex(prevOutsTypeSig.Meta)

				postOutsIndex := prgrm.CXAtomicOps[postExpressionIdx].GetOutputs(prgrm)[0]
				prgrm.CXArgs[prevOutsIdx].Indexes = append(prgrm.CXArgs[prevOutsIdx].Indexes, postOutsIndex)

			} else if prevOutsTypeSig.Type == ast.TYPE_ATOMIC || prevOutsTypeSig.Type == ast.TYPE_POINTER_ATOMIC {
				panic("type is not cx argument deprecate\n\n")
			}

		}
	}

	return prevExprs
}

// PostfixExpressionNative handles native postfix expression function calls.
//
// Input arguments description:
// 	prgrm - a CXProgram that contains all the data and arrays of the program.
// 	typeCode - type code of the native expression.
//  opStrCode - name of the function to be called that belongs to the native
// 				expression.
func PostfixExpressionNative(prgrm *ast.CXProgram, typeCode types.Code, opStrCode string) []ast.CXExpression {
	// these will always be native functions
	opCode, ok := ast.OpCodes[typeCode.Name()+"."+opStrCode]
	if !ok {
		println(ast.CompilationError(CurrentFile, LineNo) + " function '" +
			typeCode.Name() + "." + opStrCode + "' does not exist")
		return nil
		// panic(ok)
	}

	exprCXLine := ast.MakeCXLineExpression(prgrm, CurrentFile, LineNo, LineStr)
	expr := ast.MakeAtomicOperatorExpression(prgrm, ast.Natives[opCode])
	expression, err := prgrm.GetCXAtomicOp(expr.Index)
	if err != nil {
		panic(err)
	}

	pkg, err := prgrm.GetCurrentPackage()
	if err != nil {
		panic(err)
	}
	expression.Package = ast.CXPackageIndex(pkg.Index)
	return []ast.CXExpression{*exprCXLine, *expr}
}

// PostfixExpressionEmptyFunCall handles postfix expression empty function calls
// or the function calls that doesn't have any input args.
//
// Input arguments description:
// 	prgrm - a CXProgram that contains all the data and arrays of the program.
// 	prevExprs - the previous expressions that compose the postfix expression.
func PostfixExpressionEmptyFunCall(prgrm *ast.CXProgram, prevExprs []ast.CXExpression) []ast.CXExpression {
	prevExpression, err := prgrm.GetCXAtomicOpFromExpressions(prevExprs, len(prevExprs)-1)
	if err != nil {
		panic(err)
	}

	prevExpressionOperator := prgrm.GetFunctionFromArray(prevExpression.Operator)

	var firstPrevExpressionIdx int
	for i := 0; i < len(prevExprs); i++ {
		if prevExprs[i].Type == ast.CX_ATOMIC_OPERATOR {
			firstPrevExpressionIdx = prevExprs[i].Index
			break
		}
	}

	var prevExpressionOutputArg *ast.CXArgument = &ast.CXArgument{}
	if prevExpression.Outputs != nil {
		prevExpressionOutputTypeSig := prgrm.GetCXTypeSignatureFromArray(prevExpression.GetOutputs(prgrm)[0])
		if prevExpression.Outputs != nil {
			if prevExpressionOutputTypeSig.Type == ast.TYPE_CXARGUMENT_DEPRECATE {
				prevExpressionOutputArg = prgrm.GetCXArgFromArray(ast.CXArgumentIndex(prevExpressionOutputTypeSig.Meta))
			} else {
				prevExpressionOutputArg = &ast.CXArgument{}
			}
		}
	}

	if prevExpression.Outputs != nil && len(prevExpressionOutputArg.Fields) > 0 {
		// then it's a method call or function in field

	} else if prevExpressionOperator == nil {
		prevExpressionOutputTypeSig := prgrm.GetCXTypeSignatureFromArray(prevExpression.GetOutputs(prgrm)[0])
		if opCode, ok := ast.OpCodes[prevExpressionOutputTypeSig.Name]; ok {
			if pkg, err := prgrm.GetCurrentPackage(); err == nil {
				prgrm.CXAtomicOps[firstPrevExpressionIdx].Package = ast.CXPackageIndex(pkg.Index)
			}
			prgrm.CXAtomicOps[firstPrevExpressionIdx].Outputs = nil
			opIdx := prgrm.AddNativeFunctionInArray(ast.Natives[opCode])
			prgrm.CXAtomicOps[firstPrevExpressionIdx].Operator = opIdx
		}

		prgrm.CXAtomicOps[firstPrevExpressionIdx].Inputs = nil
	}

	return FunctionCall(prgrm, prevExprs, nil)
}

// PostfixExpressionFunCall handles postfix expression function calls.
//
// Input arguments description:
// 	prgrm - a CXProgram that contains all the data and arrays of the program.
// 	prevExprs - the previous expressions that compose the postfix expression.
//  args - input args list for the function call.
func PostfixExpressionFunCall(prgrm *ast.CXProgram, prevExprs []ast.CXExpression, args []ast.CXExpression) []ast.CXExpression {
	lastPrevExpression, err := prgrm.GetCXAtomicOpFromExpressions(prevExprs, len(prevExprs)-1)
	if err != nil {
		panic(err)
	}
	lastPrevExpressionOperator := prgrm.GetFunctionFromArray(lastPrevExpression.Operator)

	firstPrevExpressionIdx := prevExprs[0].Index

	var lastPrevExpressionOutputTypeSig *ast.CXTypeSignature

	var lastPrevExpressionOutputArg *ast.CXArgument = &ast.CXArgument{}
	if lastPrevExpression.Outputs != nil {
		lastPrevExpressionOutputTypeSig = prgrm.GetCXTypeSignatureFromArray(lastPrevExpression.GetOutputs(prgrm)[0])

		if lastPrevExpressionOutputTypeSig.Type == ast.TYPE_CXARGUMENT_DEPRECATE {
			lastPrevExpressionOutputArg = prgrm.GetCXArgFromArray(ast.CXArgumentIndex(lastPrevExpressionOutputTypeSig.Meta))
		} else {
			panic("type is not cx argument deprecate\n\n")
		}
	}

	if lastPrevExpression.Outputs != nil && len(lastPrevExpressionOutputArg.Fields) > 0 {
		// then it's a method
		// prevExprs[len(prevExprs) - 1].IsMethodCall = true

	} else if lastPrevExpressionOperator == nil {
		if opCode, ok := ast.OpCodes[lastPrevExpressionOutputTypeSig.Name]; ok {
			if pkg, err := prgrm.GetCurrentPackage(); err == nil {
				prgrm.CXAtomicOps[firstPrevExpressionIdx].Package = ast.CXPackageIndex(pkg.Index)
			}
			prgrm.CXAtomicOps[firstPrevExpressionIdx].Outputs = nil
			opIdx := prgrm.AddNativeFunctionInArray(ast.Natives[opCode])
			prgrm.CXAtomicOps[firstPrevExpressionIdx].Operator = opIdx
		}

		prgrm.CXAtomicOps[firstPrevExpressionIdx].Inputs = nil
	}

	return FunctionCall(prgrm, prevExprs, args)
}

// PostfixExpressionIncDec handles the incrementing or decrementing of the
// postfix expression.
//
// Input arguments description:
// 	prgrm - a CXProgram that contains all the data and arrays of the program.
// 	prevExprs - the previous expressions that compose the postfix expression.
//  isInc - true if expression is to be incremented, false if to be decremented.
func PostfixExpressionIncDec(prgrm *ast.CXProgram, prevExprs []ast.CXExpression, isInc bool) []ast.CXExpression {
	pkg, err := prgrm.GetCurrentPackage()
	if err != nil {
		panic(err)
	}

	exprCXLine := ast.MakeCXLineExpression(prgrm, CurrentFile, LineNo, LineStr)
	var expr *ast.CXExpression
	if isInc {
		expr = ast.MakeAtomicOperatorExpression(prgrm, ast.Natives[constants.OP_ADD])
	} else {
		expr = ast.MakeAtomicOperatorExpression(prgrm, ast.Natives[constants.OP_SUB])
	}

	var valB [4]byte
	types.Write_i32(valB[:], 0, 1)
	valArg := WritePrimary(prgrm, types.I32, valB[:], false)

	lastPrevExpression, err := prgrm.GetCXAtomicOpFromExpressions(prevExprs, len(prevExprs)-1)
	if err != nil {
		panic(err)
	}

	expressionIdx := expr.Index
	prgrm.CXAtomicOps[expressionIdx].Package = ast.CXPackageIndex(pkg.Index)

	typeSigIdx := lastPrevExpression.GetOutputs(prgrm)[0]
	prgrm.CXAtomicOps[expressionIdx].AddInput(prgrm, typeSigIdx)

	typeSig := ast.GetCXTypeSignatureRepresentationOfCXArg(prgrm, prgrm.GetCXArgFromArray(ast.CXArgumentIndex(valArg.Index)))
	typeSigIdx = prgrm.AddCXTypeSignatureInArray(typeSig)
	prgrm.CXAtomicOps[expressionIdx].AddInput(prgrm, typeSigIdx)

	typeSigIdx = lastPrevExpression.GetOutputs(prgrm)[0]
	prgrm.CXAtomicOps[expressionIdx].AddOutput(prgrm, typeSigIdx)

	exprs := append([]ast.CXExpression{}, *exprCXLine, *expr)
	return exprs
}

// PostfixExpressionField handles the dot notation that can be followed by an identifier.
// Examples are: `foo.bar`, `foo().bar`, `pkg.foo`
//
// Input arguments description:
// 	prgrm - a CXProgram that contains all the data and arrays of the program.
// 	prevExprs - the previous expressions that compose the left expression
// 				i.e. 'foo' in the example above.
//  ident - field name or the identifier after the dot notation.
func PostfixExpressionField(prgrm *ast.CXProgram, prevExprs []ast.CXExpression, ident string) []ast.CXExpression {
	lastExpr := prevExprs[len(prevExprs)-1]

	lastExpressionIdx := lastExpr.Index
	lastExpressionOperator := prgrm.GetFunctionFromArray(prgrm.CXAtomicOps[lastExpressionIdx].Operator)

	lastExprCXLine, _ := prgrm.GetPreviousCXLine(prevExprs, len(prevExprs)-1)

	// Then it's a function call, e.g. foo().fld
	// and we need to create some auxiliary variables to hold the result from
	// the function call
	if lastExpressionOperator != nil {
		lastExpressionOperatorOutputs := lastExpressionOperator.GetOutputs(prgrm)
		lastExpressionOperatorOutputTypeSig := prgrm.GetCXTypeSignatureFromArray(lastExpressionOperatorOutputs[0])

		var lastExpressionOperatorOutputArg *ast.CXArgument = &ast.CXArgument{}
		if lastExpressionOperatorOutputTypeSig.Type == ast.TYPE_CXARGUMENT_DEPRECATE {
			lastExpressionOperatorOutputArg = prgrm.GetCXArgFromArray(ast.CXArgumentIndex(lastExpressionOperatorOutputTypeSig.Meta))
		} else {
			panic("type is not cx argument deprecate\n\n")
		}

		opOut := lastExpressionOperatorOutputArg
		symName := generateTempVarName(constants.LOCAL_PREFIX)

		// we associate the result of the function call to the aux variable
		out := ast.MakeArgument(symName, lastExprCXLine.FileName, lastExprCXLine.LineNumber)
		out.SetType(opOut.Type)
		out.DeclarationSpecifiers = opOut.DeclarationSpecifiers
		out.StructType = opOut.StructType
		out.Size = opOut.Size
		// out.IsArray = opOut.IsArray
		// out.IsReference = opOut.IsReference
		out.Lengths = opOut.Lengths
		out.Package = prgrm.CXAtomicOps[lastExpressionIdx].Package
		out.PreviouslyDeclared = true

		outIdx := prgrm.AddCXArgInArray(out)
		typeSig := ast.GetCXTypeSignatureRepresentationOfCXArg(prgrm, prgrm.GetCXArgFromArray(outIdx))
		typeSigIdx := prgrm.AddCXTypeSignatureInArray(typeSig)
		prgrm.CXAtomicOps[lastExpressionIdx].AddOutput(prgrm, typeSigIdx)

		// we need to create an expression to hold all the modifications
		// that will take place after this if statement
		inp := ast.MakeArgument(symName, lastExprCXLine.FileName, lastExprCXLine.LineNumber)
		inp.SetType(opOut.Type)
		inp.DeclarationSpecifiers = opOut.DeclarationSpecifiers
		inp.StructType = opOut.StructType
		inp.Size = opOut.Size
		inp.Package = prgrm.CXAtomicOps[lastExpressionIdx].Package
		inpIdx := prgrm.AddCXArgInArray(inp)

		exprCXLine := ast.MakeCXLineExpression(prgrm, lastExprCXLine.FileName, lastExprCXLine.LineNumber, lastExprCXLine.LineStr)
		expr := ast.MakeAtomicOperatorExpression(prgrm, nil)

		expression, err := prgrm.GetCXAtomicOp(expr.Index)
		if err != nil {
			panic(err)
		}

		expression.Package = prgrm.CXAtomicOps[lastExpressionIdx].Package
		typeSig = ast.GetCXTypeSignatureRepresentationOfCXArg(prgrm, prgrm.GetCXArgFromArray(inpIdx))
		typeSigIdx = prgrm.AddCXTypeSignatureInArray(typeSig)
		expression.AddOutput(prgrm, typeSigIdx)
		prevExprs = append(prevExprs, *exprCXLine, *expr)

		lastExpr = prevExprs[len(prevExprs)-1]

		lastExpressionIdx = lastExpr.Index
	}

	leftExprOutputTypeSig := prgrm.GetCXTypeSignatureFromArray(prgrm.CXAtomicOps[lastExpressionIdx].GetOutputs(prgrm)[0])
	leftExprIdx := leftExprOutputTypeSig.Meta

	// then left is a first (e.g first.rest) and right is a rest
	// let's check if left is a package
	pkg, err := prgrm.GetCurrentPackage()
	if err != nil {
		panic(err)
	}

	if imp, err := pkg.GetImport(prgrm, leftExprOutputTypeSig.Name); err == nil {
		// the external property will be propagated to the following arguments
		// this way we avoid considering these arguments as module names

		if cxpackages.IsDefaultPackage(leftExprOutputTypeSig.Name) {
			//TODO: constants.ConstCodes[leftExprOutputTypeSig.Name+"."+ident]
			//TODO: only play ConstCodes are used
			//Is used for constant declaration? But only for core packages?
			if code, ok := ConstCodes[leftExprOutputTypeSig.Name+"."+ident]; ok {
				constant := Constants[code]
				valArg := WritePrimary(prgrm, constant.Type, constant.Value, false)

				// TODO: to be removed and implemented with correct process to change the index
				lastExpressionOutputFieldTypeSig := prgrm.GetCXTypeSignatureFromArray(prgrm.CXAtomicOps[lastExpressionIdx].Outputs.Fields[0])
				if lastExpressionOutputFieldTypeSig.Type == ast.TYPE_CXARGUMENT_DEPRECATE {
					lastExpressionOutputFieldTypeSig.Meta = valArg.Index
					lastExpressionOutputFieldTypeSig.Name = valArg.Name
					lastExpressionOutputFieldTypeSig.Offset = valArg.Offset
				} else if lastExpressionOutputFieldTypeSig.Type == ast.TYPE_ATOMIC {
					// TODO: Review if this is the correct one
					typeSig := ast.GetCXTypeSignatureRepresentationOfCXArg(prgrm, valArg)
					typeSig.Offset = valArg.Offset
					typeSigIdx := prgrm.AddCXTypeSignatureInArray(typeSig)

					prgrm.CXAtomicOps[lastExpressionIdx].Outputs.Fields[0] = typeSigIdx
				} else if lastExpressionOutputFieldTypeSig.Type == ast.TYPE_POINTER_ATOMIC {
					// TODO: Review if this is the correct one
					typeSig := ast.GetCXTypeSignatureRepresentationOfCXArg(prgrm, valArg)
					typeSig.Offset = valArg.Offset
					typeSigIdx := prgrm.AddCXTypeSignatureInArray(typeSig)

					prgrm.CXAtomicOps[lastExpressionIdx].Outputs.Fields[0] = typeSigIdx
				}

				return prevExprs
			} else if _, ok := ast.OpCodes[leftExprOutputTypeSig.Name+"."+ident]; ok {
				// then it's a native
				// TODO: we'd be referring to the function itself, not a function call
				// (functions as first-class objects)
				leftExprOutputTypeSigName := leftExprOutputTypeSig.Name + "." + ident
				leftExprOutputTypeSig.Name = leftExprOutputTypeSigName
				if leftExprOutputTypeSig.Type == ast.TYPE_CXARGUMENT_DEPRECATE {
					prgrm.CXArgs[leftExprIdx].Name = leftExprOutputTypeSigName
				}

				return prevExprs
			}
		}

		leftExprOutputTypeSig.Package = ast.CXPackageIndex(imp.Index)
		if leftExprOutputTypeSig.Type == ast.TYPE_CXARGUMENT_DEPRECATE {
			prgrm.CXArgs[leftExprIdx].Package = ast.CXPackageIndex(imp.Index)
		}

		if glbl, err := imp.GetGlobal(prgrm, ident); err == nil {
			if glbl.Type == ast.TYPE_CXARGUMENT_DEPRECATE {
				glblArg := prgrm.GetCXArgFromArray(ast.CXArgumentIndex(glbl.Meta))

				// then it's a global
				// prevExprs[len(prevExprs)-1].ProgramOutput[0] = glbl
				lastExpressionOutputTypeSig := prgrm.GetCXTypeSignatureFromArray(prgrm.CXAtomicOps[lastExpressionIdx].GetOutputs(prgrm)[0])

				var lastExpressionOutputArg *ast.CXArgument = &ast.CXArgument{}
				if lastExpressionOutputTypeSig.Type == ast.TYPE_CXARGUMENT_DEPRECATE {
					lastExpressionOutputArg = prgrm.GetCXArgFromArray(ast.CXArgumentIndex(lastExpressionOutputTypeSig.Meta))
				} else {
					panic("type is not cx argument deprecate\n\n")
				}

				lastExpressionOutputArg.Name = glblArg.Name
				lastExpressionOutputTypeSig.Name = glblArg.Name

				lastExpressionOutputArg.Type = glblArg.Type
				lastExpressionOutputArg.StructType = glblArg.StructType
				lastExpressionOutputArg.Size = glblArg.Size
				lastExpressionOutputArg.PointerTargetType = glblArg.PointerTargetType
				lastExpressionOutputArg.IsSlice = glblArg.IsSlice
				lastExpressionOutputArg.IsStruct = glblArg.IsStruct
				lastExpressionOutputArg.Package = glblArg.Package
			} else if glbl.Type == ast.TYPE_ATOMIC {
				// then it's a global
				// prevExprs[len(prevExprs)-1].ProgramOutput[0] = glbl
				lastExpressionOutputTypeSig := prgrm.GetCXTypeSignatureFromArray(prgrm.CXAtomicOps[lastExpressionIdx].GetOutputs(prgrm)[0])

				var lastExpressionOutputArg *ast.CXArgument = &ast.CXArgument{}
				if lastExpressionOutputTypeSig.Type == ast.TYPE_CXARGUMENT_DEPRECATE {
					lastExpressionOutputArg = prgrm.GetCXArgFromArray(ast.CXArgumentIndex(lastExpressionOutputTypeSig.Meta))
				} else {
					panic("type is not cx argument deprecate\n\n")
				}

				lastExpressionOutputArg.Name = glbl.Name
				lastExpressionOutputTypeSig.Name = glbl.Name

				lastExpressionOutputArg.Type = types.Code(glbl.Meta)
				lastExpressionOutputArg.Offset = glbl.Offset
				lastExpressionOutputArg.Size = types.Code(glbl.Meta).Size()

				lastExpressionOutputArg.Package = glbl.Package
			} else if glbl.Type == ast.TYPE_POINTER_ATOMIC {
				// then it's a global
				// prevExprs[len(prevExprs)-1].ProgramOutput[0] = glbl
				lastExpressionOutputTypeSig := prgrm.GetCXTypeSignatureFromArray(prgrm.CXAtomicOps[lastExpressionIdx].GetOutputs(prgrm)[0])

				var lastExpressionOutputArg *ast.CXArgument = &ast.CXArgument{}
				if lastExpressionOutputTypeSig.Type == ast.TYPE_CXARGUMENT_DEPRECATE {
					lastExpressionOutputArg = prgrm.GetCXArgFromArray(ast.CXArgumentIndex(lastExpressionOutputTypeSig.Meta))
				} else {
					panic("type is not cx argument deprecate\n\n")
				}

				lastExpressionOutputArg.Name = glbl.Name
				lastExpressionOutputTypeSig.Name = glbl.Name

				lastExpressionOutputArg.Type = types.Code(glbl.Meta)
				lastExpressionOutputArg.Offset = glbl.Offset
				lastExpressionOutputArg.Size = types.POINTER.Size()

				lastExpressionOutputArg.Package = glbl.Package
			}

		} else if fn, err := imp.GetFunction(prgrm, ident); err == nil {
			// then it's a function
			// not sure about this next line
			prgrm.CXAtomicOps[lastExpressionIdx].Outputs = nil
			prgrm.CXAtomicOps[lastExpressionIdx].Operator = ast.CXFunctionIndex(fn.Index)
		} else if strct, err := prgrm.GetStruct(ident, imp.Name); err == nil {
			lastExpressionOutputTypeSig := prgrm.GetCXTypeSignatureFromArray(prgrm.CXAtomicOps[lastExpressionIdx].GetOutputs(prgrm)[0])
			var lastExpressionOutputArg *ast.CXArgument = &ast.CXArgument{}
			if lastExpressionOutputTypeSig.Type == ast.TYPE_CXARGUMENT_DEPRECATE {
				lastExpressionOutputArg = prgrm.GetCXArgFromArray(ast.CXArgumentIndex(lastExpressionOutputTypeSig.Meta))
			} else {
				panic("type is not cx argument deprecate\n\n")
			}

			lastExpressionOutputArg.StructType = strct
		} else {
			// panic(err)
			fmt.Println(err)
			return nil
		}
	} else {
		// then left is not a package name
		if cxpackages.IsDefaultPackage(leftExprOutputTypeSig.Name) {
			println(ast.CompilationError(prgrm.CXArgs[leftExprIdx].ArgDetails.FileName, prgrm.CXArgs[leftExprIdx].ArgDetails.FileLine),
				fmt.Sprintf("identifier '%s' does not exist",
					leftExprOutputTypeSig.Name))
			os.Exit(constants.CX_COMPILATION_ERROR)
		}

		// then it's a struct
		prgrm.CXArgs[leftExprIdx].IsStruct = true

		leftExprField := ast.MakeArgument(ident, CurrentFile, LineNo)
		leftPkg, err := prgrm.GetPackageFromArray(leftExprOutputTypeSig.Package)
		if err != nil {
			panic(err)
		}
		leftExprField.SetType(types.IDENTIFIER)
		leftExprField.SetPackage(leftPkg)

		if leftExprOutputTypeSig.Type == ast.TYPE_CXARGUMENT_DEPRECATE {
			leftExprFieldIdx := prgrm.AddCXArgInArray(leftExprField)
			prgrm.CXArgs[leftExprIdx].Fields = append(prgrm.CXArgs[leftExprIdx].Fields, leftExprFieldIdx)
		} else if leftExprOutputTypeSig.Type == ast.TYPE_ATOMIC || leftExprOutputTypeSig.Type == ast.TYPE_POINTER_ATOMIC {
			panic("type is not cx arg deprecate")
		}

	}

	return prevExprs
}
