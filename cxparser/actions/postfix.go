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

		useExprCXLine := ast.MakeCXLineExpression(prgrm, prevExprCXLine.FileName, prevExprCXLine.LineNumber, prevExprCXLine.LineStr)
		useExpr := ast.MakeAtomicOperatorExpression(prgrm, nil)
		useExpressionIdx := useExpr.Index
		prgrm.CXAtomicOps[useExpressionIdx].Package = prgrm.CXAtomicOps[prevExpressionIdx].Package

		prevExpressionOperatorOutputs := prevExpressionOperator.GetOutputs(prgrm)
		prevExpressionOperatorOutputTypeSig := prgrm.GetCXTypeSignatureFromArray(prevExpressionOperatorOutputs[0])
		if prevExpressionOperatorOutputTypeSig.Type == ast.TYPE_CXARGUMENT_DEPRECATE {
			prevExpressionOperatorOutputArg := prgrm.GetCXArgFromArray(ast.CXArgumentIndex(prevExpressionOperatorOutputTypeSig.Meta))

			prevExpressionOperatorOutput := prevExpressionOperatorOutputArg
			out := ast.MakeArgument(genName, prevExprCXLine.FileName, prevExprCXLine.LineNumber-1).SetType(prevExpressionOperatorOutput.Type)

			out.DeclarationSpecifiers = prevExpressionOperatorOutput.DeclarationSpecifiers
			out.StructType = prevExpressionOperatorOutput.StructType
			out.Size = prevExpressionOperatorOutput.Size
			out.Lengths = prevExpressionOperatorOutput.Lengths
			out.PreviouslyDeclared = true
			out.Package = prevExpressionOperatorOutput.Package

			typeSig := ast.GetCXTypeSignatureRepresentationOfCXArg(prgrm, out)
			typeSigIdx := prgrm.AddCXTypeSignatureInArray(typeSig)
			prgrm.CXAtomicOps[prevExpressionIdx].AddOutput(prgrm, typeSigIdx)

			inp := ast.MakeArgument(genName, prevExprCXLine.FileName, prevExprCXLine.LineNumber).SetType(prevExpressionOperatorOutput.Type)

			inp.DeclarationSpecifiers = prevExpressionOperatorOutput.DeclarationSpecifiers
			inp.StructType = prevExpressionOperatorOutput.StructType
			inp.Size = prevExpressionOperatorOutput.Size
			inp.Lengths = prevExpressionOperatorOutput.Lengths
			inp.PreviouslyDeclared = true

			typeSig = ast.GetCXTypeSignatureRepresentationOfCXArg(prgrm, inp)
			typeSigIdx = prgrm.AddCXTypeSignatureInArray(typeSig)
			prgrm.CXAtomicOps[useExpressionIdx].AddOutput(prgrm, typeSigIdx)
		} else if prevExpressionOperatorOutputTypeSig.Type == ast.TYPE_SLICE_ATOMIC {
			sliceDetails := prgrm.GetCXTypeSignatureArrayFromArray(prevExpressionOperatorOutputTypeSig.Meta)

			newTypeSigSlice := ast.CXTypeSignature_Array{}
			newTypeSigSlice.Meta = sliceDetails.Meta
			newTypeSigSlice.Lengths = sliceDetails.Lengths
			newTypeSigSliceIdx := prgrm.AddCXTypeSignatureArrayInArray(&newTypeSigSlice)

			newTypeSig := ast.CXTypeSignature{}
			newTypeSig.Name = genName
			newTypeSig.Package = prevExpressionOperatorOutputTypeSig.Package
			newTypeSig.Type = prevExpressionOperatorOutputTypeSig.Type
			newTypeSig.Meta = newTypeSigSliceIdx

			outNewTypeSigIdx := prgrm.AddCXTypeSignatureInArray(&newTypeSig)
			prgrm.CXAtomicOps[prevExpressionIdx].AddOutput(prgrm, outNewTypeSigIdx)

			inpNewTypeSigSlice := ast.CXTypeSignature_Array{}
			inpNewTypeSigSlice.Meta = sliceDetails.Meta
			inpNewTypeSigSlice.Lengths = sliceDetails.Lengths
			inpNewTypeSigSliceIdx := prgrm.AddCXTypeSignatureArrayInArray(&inpNewTypeSigSlice)

			inpNewTypeSig := ast.CXTypeSignature{}
			inpNewTypeSig.Name = genName
			inpNewTypeSig.Package = prevExpressionOperatorOutputTypeSig.Package
			inpNewTypeSig.Type = prevExpressionOperatorOutputTypeSig.Type
			inpNewTypeSig.Meta = inpNewTypeSigSliceIdx
			inpNewTypeSigIdx := prgrm.AddCXTypeSignatureInArray(&inpNewTypeSig)
			prgrm.CXAtomicOps[useExpressionIdx].AddOutput(prgrm, inpNewTypeSigIdx)
		} else {
			panic("type is not known\n")
		}

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

		if len(prevExpressionOutputArg.Fields) > 0 {
			elt = prgrm.GetCXArgFromArray(prevExpressionOutputArg.Fields[len(prevExpressionOutputArg.Fields)-1])
		} else {
			elt = prevExpressionOutputArg
		}

		// elt.IsArray = false
		elt.DereferenceOperations = append(elt.DereferenceOperations, constants.DEREF_ARRAY)
		elt.DeclarationSpecifiers = append(elt.DeclarationSpecifiers, constants.DECL_INDEXING)

	} else if prevExpressionOutputTypeSig.Type == ast.TYPE_ATOMIC || prevExpressionOutputTypeSig.Type == ast.TYPE_POINTER_ATOMIC {
		// TODO: implement correct way when it is possible
		// temporarily return back to cx arg deprecate
		arg := ast.MakeArgument(prevExpressionOutputTypeSig.Name, CurrentFile, LineNo) // fix: line numbers in errors sometimes report +1 or -1. Issue #195
		arg.SetType(types.IDENTIFIER)
		arg.Name = prevExpressionOutputTypeSig.Name
		arg.Package = prevExpressionOutputTypeSig.Package
		arg.PassBy = prevExpressionOutputTypeSig.PassBy

		argIdx := prgrm.AddCXArgInArray(arg)
		prevExpressionOutputTypeSig.Type = ast.TYPE_CXARGUMENT_DEPRECATE
		prevExpressionOutputTypeSig.Meta = int(argIdx)

		prevExpressionOutputArg = prgrm.GetCXArgFromArray(ast.CXArgumentIndex(prevExpressionOutputTypeSig.Meta))

		elt = prevExpressionOutputArg
		// elt.IsArray = false
		elt.DereferenceOperations = append(elt.DereferenceOperations, constants.DEREF_ARRAY)
		elt.DeclarationSpecifiers = append(elt.DeclarationSpecifiers, constants.DECL_INDEXING)
	} else if prevExpressionOutputTypeSig.Type == ast.TYPE_SLICE_ATOMIC {
		prevExpressionOutputTypeSig.IsDeref = true
		prevExpressionOutputArg = &ast.CXArgument{ArgDetails: &ast.CXArgumentDebug{}}
	} else {
		panic("type is not known")
	}

	postExpressionIdx := postExprs[len(postExprs)-1].Index
	postExpressionOperator := prgrm.GetFunctionFromArray(prgrm.CXAtomicOps[postExpressionIdx].Operator)
	postExpressionOperatorOutputs := postExpressionOperator.GetOutputs(prgrm)

	if len(prevExpressionOutputArg.Fields) > 0 {
		fldIdx := prevExpressionOutputArg.Fields[len(prevExpressionOutputArg.Fields)-1]

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

			typeSig := ast.GetCXTypeSignatureRepresentationOfCXArg(prgrm, sym)
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

				typeSig := ast.GetCXTypeSignatureRepresentationOfCXArg(prgrm, idxSym)
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
			} else {
				panic("type is not known")
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
			} else if prevOutsTypeSig.Type == ast.TYPE_SLICE_ATOMIC {
				sliceDetails := prgrm.GetCXTypeSignatureArrayFromArray(prevOutsTypeSig.Meta)

				postOutsIndex := prgrm.CXAtomicOps[postExpressionIdx].GetOutputs(prgrm)[0]
				sliceDetails.Indexes = append(sliceDetails.Indexes, postOutsIndex)
			} else {
				panic("type is not known")
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
		} else if lastPrevExpressionOutputTypeSig.Type == ast.TYPE_ATOMIC || lastPrevExpressionOutputTypeSig.Type == ast.TYPE_POINTER_ATOMIC {
			// do nothing
		} else {
			panic("type is not known")
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

	typeSig := ast.GetCXTypeSignatureRepresentationOfCXArg(prgrm, valArg)
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

		symName := generateTempVarName(constants.LOCAL_PREFIX)
		if lastExpressionOperatorOutputTypeSig.Type == ast.TYPE_CXARGUMENT_DEPRECATE {
			lastExpressionOperatorOutputArg := prgrm.GetCXArgFromArray(ast.CXArgumentIndex(lastExpressionOperatorOutputTypeSig.Meta))

			opOut := lastExpressionOperatorOutputArg

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

			typeSig := ast.GetCXTypeSignatureRepresentationOfCXArg(prgrm, out)
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

			exprCXLine := ast.MakeCXLineExpression(prgrm, lastExprCXLine.FileName, lastExprCXLine.LineNumber, lastExprCXLine.LineStr)
			expr := ast.MakeAtomicOperatorExpression(prgrm, nil)

			expression, err := prgrm.GetCXAtomicOp(expr.Index)
			if err != nil {
				panic(err)
			}

			expression.Package = prgrm.CXAtomicOps[lastExpressionIdx].Package
			typeSig = ast.GetCXTypeSignatureRepresentationOfCXArg(prgrm, inp)
			typeSigIdx = prgrm.AddCXTypeSignatureInArray(typeSig)
			expression.AddOutput(prgrm, typeSigIdx)
			prevExprs = append(prevExprs, *exprCXLine, *expr)

			lastExpr = prevExprs[len(prevExprs)-1]

			lastExpressionIdx = lastExpr.Index
		} else if lastExpressionOperatorOutputTypeSig.Type == ast.TYPE_STRUCT {
			var newStructDetails ast.CXTypeSignature_Struct
			structDetails := prgrm.GetCXTypeSignatureStructFromArray(lastExpressionOperatorOutputTypeSig.Meta)
			newStructDetails = *structDetails
			outStructDetailsIdx := prgrm.AddCXTypeSignatureStructInArray(&newStructDetails)

			outTypeSig := *lastExpressionOperatorOutputTypeSig
			outTypeSig.Name = symName
			outTypeSig.Meta = outStructDetailsIdx
			outTypeSigIdx := prgrm.AddCXTypeSignatureInArray(&outTypeSig)
			prgrm.CXAtomicOps[lastExpressionIdx].AddOutput(prgrm, outTypeSigIdx)

			// we need to create an expression to hold all the modifications
			// that will take place after this if statement
			inpTypeSig := *lastExpressionOperatorOutputTypeSig
			inpTypeSig.Name = symName
			inpStructDetailsIdx := prgrm.AddCXTypeSignatureStructInArray(&newStructDetails)
			inpTypeSig.Meta = inpStructDetailsIdx
			inpTypeSigIdx := prgrm.AddCXTypeSignatureInArray(&inpTypeSig)

			exprCXLine := ast.MakeCXLineExpression(prgrm, lastExprCXLine.FileName, lastExprCXLine.LineNumber, lastExprCXLine.LineStr)
			expr := ast.MakeAtomicOperatorExpression(prgrm, nil)

			expression, err := prgrm.GetCXAtomicOp(expr.Index)
			if err != nil {
				panic(err)
			}

			expression.Package = prgrm.CXAtomicOps[lastExpressionIdx].Package
			expression.AddOutput(prgrm, inpTypeSigIdx)
			prevExprs = append(prevExprs, *exprCXLine, *expr)

			lastExpr = prevExprs[len(prevExprs)-1]

			lastExpressionIdx = lastExpr.Index
		} else {
			panic("type is not cx argument deprecate\n\n")
		}

	}

	leftExprOutputTypeSig := prgrm.GetCXTypeSignatureFromArray(prgrm.CXAtomicOps[lastExpressionIdx].GetOutputs(prgrm)[0])

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
				} else {
					panic("type is not known")
				}

				return prevExprs
			} else if _, ok := ast.OpCodes[leftExprOutputTypeSig.Name+"."+ident]; ok {
				// then it's a native
				// TODO: we'd be referring to the function itself, not a function call
				// (functions as first-class objects)
				leftExprOutputTypeSigName := leftExprOutputTypeSig.Name + "." + ident
				leftExprOutputTypeSig.Name = leftExprOutputTypeSigName
				if leftExprOutputTypeSig.Type == ast.TYPE_CXARGUMENT_DEPRECATE {
					leftExprIdx := leftExprOutputTypeSig.Meta
					prgrm.CXArgs[leftExprIdx].Name = leftExprOutputTypeSigName
				}

				return prevExprs
			}
		}

		leftExprOutputTypeSig.Package = ast.CXPackageIndex(imp.Index)
		if leftExprOutputTypeSig.Type == ast.TYPE_CXARGUMENT_DEPRECATE {
			leftExprIdx := leftExprOutputTypeSig.Meta
			prgrm.CXArgs[leftExprIdx].Package = ast.CXPackageIndex(imp.Index)
		}

		if glbl, err := imp.GetGlobal(prgrm, ident); err == nil {
			// then it's a global
			// prevExprs[len(prevExprs)-1].ProgramOutput[0] = glbl
			if glbl.Type == ast.TYPE_CXARGUMENT_DEPRECATE {
				glblArg := prgrm.GetCXArgFromArray(ast.CXArgumentIndex(glbl.Meta))

				// then it's a global
				// prevExprs[len(prevExprs)-1].ProgramOutput[0] = glbl
				lastExpressionOutputTypeSig := prgrm.GetCXTypeSignatureFromArray(prgrm.CXAtomicOps[lastExpressionIdx].GetOutputs(prgrm)[0])

				var lastExpressionOutputArg *ast.CXArgument = &ast.CXArgument{}
				if lastExpressionOutputTypeSig.Type == ast.TYPE_CXARGUMENT_DEPRECATE {
					lastExpressionOutputArg = prgrm.GetCXArgFromArray(ast.CXArgumentIndex(lastExpressionOutputTypeSig.Meta))
					lastExpressionOutputArg.Name = glblArg.Name
					lastExpressionOutputTypeSig.Name = glblArg.Name

					lastExpressionOutputArg.Type = glblArg.Type
					lastExpressionOutputArg.StructType = glblArg.StructType
					lastExpressionOutputArg.Size = glblArg.Size
					lastExpressionOutputArg.PointerTargetType = glblArg.PointerTargetType
					lastExpressionOutputArg.Package = glblArg.Package
				} else if lastExpressionOutputTypeSig.Type == ast.TYPE_ATOMIC || lastExpressionOutputTypeSig.Type == ast.TYPE_POINTER_ATOMIC {
					lastExpressionOutputTypeSig.Name = glbl.Name
					lastExpressionOutputTypeSig.Type = glbl.Type
					lastExpressionOutputTypeSig.Meta = glbl.Meta
					lastExpressionOutputTypeSig.Offset = glbl.Offset
				} else {
					panic("type is not known")
				}

			} else if glbl.Type == ast.TYPE_ATOMIC || glbl.Type == ast.TYPE_POINTER_ATOMIC {
				lastExpressionOutputTypeSig := prgrm.GetCXTypeSignatureFromArray(prgrm.CXAtomicOps[lastExpressionIdx].GetOutputs(prgrm)[0])
				lastExpressionOutputTypeSig.Name = glbl.Name
				lastExpressionOutputTypeSig.Type = glbl.Type
				lastExpressionOutputTypeSig.Meta = glbl.Meta
				lastExpressionOutputTypeSig.Offset = glbl.Offset
			} else {
				panic("type is not known")
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
			leftExprIdx := leftExprOutputTypeSig.Meta
			println(ast.CompilationError(prgrm.CXArgs[leftExprIdx].ArgDetails.FileName, prgrm.CXArgs[leftExprIdx].ArgDetails.FileLine),
				fmt.Sprintf("identifier '%s' does not exist",
					leftExprOutputTypeSig.Name))
			os.Exit(constants.CX_COMPILATION_ERROR)
		}

		// then it's a struct
		leftExprField := ast.MakeArgument(ident, CurrentFile, LineNo)
		leftPkg, err := prgrm.GetPackageFromArray(leftExprOutputTypeSig.Package)
		if err != nil {
			panic(err)
		}
		leftExprField.SetType(types.IDENTIFIER)
		leftExprField.SetPackage(leftPkg)
		leftExprFieldIdx := prgrm.AddCXArgInArray(leftExprField)

		var leftExprIdx int
		if leftExprOutputTypeSig.Type == ast.TYPE_CXARGUMENT_DEPRECATE {
			leftExprIdx = leftExprOutputTypeSig.Meta

			prgrm.CXArgs[leftExprIdx].Fields = append(prgrm.CXArgs[leftExprIdx].Fields, leftExprFieldIdx)
		} else if leftExprOutputTypeSig.Type == ast.TYPE_ATOMIC || leftExprOutputTypeSig.Type == ast.TYPE_POINTER_ATOMIC {
			// TODO: implement correct way when it is possible
			// temporarily return back to cx arg deprecate
			// since the arg will become a struct
			arg := ast.MakeArgument(leftExprOutputTypeSig.Name, CurrentFile, LineNo) // fix: line numbers in errors sometimes report +1 or -1. Issue #195
			arg.SetType(types.IDENTIFIER)
			arg.Name = leftExprOutputTypeSig.Name
			arg.Package = leftExprOutputTypeSig.Package
			arg.PassBy = leftExprOutputTypeSig.PassBy

			argIdx := prgrm.AddCXArgInArray(arg)
			leftExprOutputTypeSig.Type = ast.TYPE_CXARGUMENT_DEPRECATE
			leftExprOutputTypeSig.Meta = int(argIdx)

			leftExprIdx = int(argIdx)
			prgrm.CXArgs[leftExprIdx].Fields = append(prgrm.CXArgs[leftExprIdx].Fields, leftExprFieldIdx)
		} else if leftExprOutputTypeSig.Type == ast.TYPE_STRUCT {
			structDetails := prgrm.GetCXTypeSignatureStructFromArray(leftExprOutputTypeSig.Meta)
			structDetails.Fields = append(structDetails.Fields, leftExprFieldIdx)
		} else {
			panic("type is not known")
		}
	}

	return prevExprs
}
