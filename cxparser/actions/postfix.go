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

		var prevExpressionOperatorOutputArg *ast.CXArgument = &ast.CXArgument{}
		if prevExpressionOperatorOutputs[0].Type == ast.TYPE_CXARGUMENT_DEPRECATE {
			prevExpressionOperatorOutputArg = prgrm.GetCXArgFromArray(ast.CXArgumentIndex(prevExpressionOperatorOutputs[0].Meta))
		} else {
			panic("type is not cx argument deprecate\n\n")
		}

		prevExpressionOperatorOutput := prevExpressionOperatorOutputArg
		out := ast.MakeArgument(genName, prevExprCXLine.FileName, prevExprCXLine.LineNumber-1).SetType(prevExpressionOperatorOutput.Type)

		out.DeclarationSpecifiers = prevExpressionOperatorOutput.DeclarationSpecifiers
		out.StructType = prevExpressionOperatorOutput.StructType
		out.Size = prevExpressionOperatorOutput.Size
		out.TotalSize = prevExpressionOperatorOutput.TotalSize
		out.Lengths = prevExpressionOperatorOutput.Lengths
		out.IsSlice = prevExpressionOperatorOutput.IsSlice
		out.PreviouslyDeclared = true
		outIdx := prgrm.AddCXArgInArray(out)

		typeSig := ast.GetCXTypeSignatureRepresentationOfCXArg_ForGlobals_CXAtomicOps(prgrm, prgrm.GetCXArgFromArray(outIdx))
		prgrm.CXAtomicOps[prevExpressionIdx].AddOutput(prgrm, typeSig)

		inp := ast.MakeArgument(genName, prevExprCXLine.FileName, prevExprCXLine.LineNumber).SetType(prevExpressionOperatorOutput.Type)

		inp.DeclarationSpecifiers = prevExpressionOperatorOutput.DeclarationSpecifiers
		inp.StructType = prevExpressionOperatorOutput.StructType
		inp.Size = prevExpressionOperatorOutput.Size
		inp.TotalSize = prevExpressionOperatorOutput.TotalSize
		inp.Lengths = prevExpressionOperatorOutput.Lengths
		inp.IsSlice = prevExpressionOperatorOutput.IsSlice
		inp.PreviouslyDeclared = true
		inpIdx := prgrm.AddCXArgInArray(inp)

		useExprCXLine := ast.MakeCXLineExpression(prgrm, prevExprCXLine.FileName, prevExprCXLine.LineNumber, prevExprCXLine.LineStr)
		useExpr := ast.MakeAtomicOperatorExpression(prgrm, nil)
		useExpressionIdx := useExpr.Index

		prgrm.CXAtomicOps[useExpressionIdx].Package = prgrm.CXAtomicOps[prevExpressionIdx].Package
		typeSig = ast.GetCXTypeSignatureRepresentationOfCXArg_ForGlobals_CXAtomicOps(prgrm, prgrm.GetCXArgFromArray(inpIdx))
		prgrm.CXAtomicOps[useExpressionIdx].AddOutput(prgrm, typeSig)
		prevExprs = append(prevExprs, *useExprCXLine, *useExpr)
	}

	prevExpression, err := prgrm.GetCXAtomicOpFromExpressions(prevExprs, len(prevExprs)-1)
	if err != nil {
		panic(err)
	}

	var prevExpressionOutputArg *ast.CXArgument = &ast.CXArgument{}
	if prevExpression.GetOutputs(prgrm)[0].Type == ast.TYPE_CXARGUMENT_DEPRECATE {
		prevExpressionOutputArg = prgrm.GetCXArgFromArray(ast.CXArgumentIndex(prevExpression.GetOutputs(prgrm)[0].Meta))
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
			// expr.AddInput(postExprs[len(postExprs)-1].ProgramOutput[0])
			indexIdx := ast.CXArgumentIndex(prgrm.CXAtomicOps[postExpressionIdx].GetOutputs(prgrm)[0].Meta)
			prgrm.CXArgs[fldIdx].Indexes = append(prgrm.CXArgs[fldIdx].Indexes, indexIdx)
		} else {
			var postExpressionOperatorOutputArg *ast.CXArgument = &ast.CXArgument{}
			if postExpressionOperatorOutputs[0].Type == ast.TYPE_CXARGUMENT_DEPRECATE {
				postExpressionOperatorOutputArg = prgrm.GetCXArgFromArray(ast.CXArgumentIndex(postExpressionOperatorOutputs[0].Meta))
			} else {
				panic("type is not cx argument deprecate\n\n")
			}

			sym := ast.MakeArgument(generateTempVarName(constants.LOCAL_PREFIX), CurrentFile, LineNo).SetType(postExpressionOperatorOutputArg.Type)
			sym.Package = prgrm.CXAtomicOps[postExpressionIdx].Package
			sym.PreviouslyDeclared = true
			symIdx := prgrm.AddCXArgInArray(sym)
			typeSig := ast.GetCXTypeSignatureRepresentationOfCXArg_ForGlobals_CXAtomicOps(prgrm, prgrm.GetCXArgFromArray(symIdx))
			prgrm.CXAtomicOps[postExpressionIdx].AddOutput(prgrm, typeSig)

			prevExprs = append(postExprs, prevExprs...)

			prgrm.CXArgs[fldIdx].Indexes = append(prgrm.CXArgs[fldIdx].Indexes, symIdx)
		}
	} else {
		if len(prgrm.CXAtomicOps[postExpressionIdx].GetOutputs(prgrm)) < 1 {
			var postExpressionOperatorOutputArg *ast.CXArgument = &ast.CXArgument{}
			if postExpressionOperatorOutputs[0].Type == ast.TYPE_CXARGUMENT_DEPRECATE {
				postExpressionOperatorOutputArg = prgrm.GetCXArgFromArray(ast.CXArgumentIndex(postExpressionOperatorOutputs[0].Meta))
			} else {
				panic("type is not cx argument deprecate\n\n")
			}

			// then it's an expression (e.g. i32.add(0, 0))
			// we create a gensym for it
			idxSym := ast.MakeArgument(generateTempVarName(constants.LOCAL_PREFIX), CurrentFile, LineNo).SetType(postExpressionOperatorOutputArg.Type)
			idxSym.Size = postExpressionOperatorOutputArg.Size
			idxSym.TotalSize = postExpressionOperatorOutputs[0].GetSize(prgrm)

			idxSym.Package = prgrm.CXAtomicOps[postExpressionIdx].Package
			idxSym.PreviouslyDeclared = true

			idxSymIdx := prgrm.AddCXArgInArray(idxSym)
			typeSig := ast.GetCXTypeSignatureRepresentationOfCXArg_ForGlobals_CXAtomicOps(prgrm, prgrm.GetCXArgFromArray(idxSymIdx))
			prgrm.CXAtomicOps[postExpressionIdx].AddOutput(prgrm, typeSig)

			var prevExpressionOutputArg *ast.CXArgument = &ast.CXArgument{}
			if prevExpression.GetOutputs(prgrm)[0].Type == ast.TYPE_CXARGUMENT_DEPRECATE {
				prevExpressionOutputArg = prgrm.GetCXArgFromArray(ast.CXArgumentIndex(prevExpression.GetOutputs(prgrm)[0].Meta))
			} else {
				panic("type is not cx argument deprecate\n\n")
			}

			prevExpressionOutput := prevExpressionOutputArg
			prevExpressionOutput.Indexes = append(prevExpressionOutput.Indexes, idxSymIdx)

			// we push the index expression
			prevExprs = append(postExprs, prevExprs...)
		} else {
			prevOutsTypeSig := prevExpression.GetOutputs(prgrm)[0]
			prevOutsIdx := ast.CXArgumentIndex(prevOutsTypeSig.Meta)

			postOutsTypeSig := prgrm.CXAtomicOps[postExpressionIdx].GetOutputs(prgrm)[0]
			postOutsIdx := ast.CXArgumentIndex(postOutsTypeSig.Meta)
			prgrm.CXArgs[prevOutsIdx].Indexes = append(prgrm.CXArgs[prevOutsIdx].Indexes, postOutsIdx)
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

	firstPrevExpressionIdx := prevExprs[0].Index

	var prevExpressionOutputArg *ast.CXArgument = &ast.CXArgument{}
	if prevExpression.Outputs != nil {

		if prevExpression.GetOutputs(prgrm)[0].Type == ast.TYPE_CXARGUMENT_DEPRECATE {
			prevExpressionOutputArg = prgrm.GetCXArgFromArray(ast.CXArgumentIndex(prevExpression.GetOutputs(prgrm)[0].Meta))
		} else {
			panic("type is not cx argument deprecate\n\n")
		}
	}

	if prevExpression.Outputs != nil && len(prevExpressionOutputArg.Fields) > 0 {
		// then it's a method call or function in field
		// prevExprs[len(prevExprs) - 1].IsMethodCall = true
		// expr.IsMethodCall = true
		// // method name
		// expr.Operator = MakeFunction(expr.ProgramOutput[0].Fields[0].Name)
		// inp := cxcore.MakeArgument(expr.ProgramOutput[0].Name, CurrentFile, LineNo)
		// inp.Package = expr.Package
		// inp.Type = expr.ProgramOutput[0].Type
		// inp.StructType = expr.ProgramOutput[0].StructType
		// expr.ProgramInput = append(expr.ProgramInput, inp)

	} else if prevExpressionOperator == nil {
		if opCode, ok := ast.OpCodes[prevExpression.GetOutputs(prgrm)[0].Name]; ok {
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

	var lastPrevExpressionOutputArg *ast.CXArgument = &ast.CXArgument{}
	if lastPrevExpression.Outputs != nil {
		if lastPrevExpression.GetOutputs(prgrm)[0].Type == ast.TYPE_CXARGUMENT_DEPRECATE {
			lastPrevExpressionOutputArg = prgrm.GetCXArgFromArray(ast.CXArgumentIndex(lastPrevExpression.GetOutputs(prgrm)[0].Meta))
		} else {
			panic("type is not cx argument deprecate\n\n")
		}
	}

	if lastPrevExpression.Outputs != nil && len(lastPrevExpressionOutputArg.Fields) > 0 {
		// then it's a method
		// prevExprs[len(prevExprs) - 1].IsMethodCall = true

	} else if lastPrevExpressionOperator == nil {
		if opCode, ok := ast.OpCodes[lastPrevExpression.GetOutputs(prgrm)[0].Name]; ok {
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

	typeSig := lastPrevExpression.GetOutputs(prgrm)[0]
	prgrm.CXAtomicOps[expressionIdx].AddInput(prgrm, typeSig)

	typeSig = ast.GetCXTypeSignatureRepresentationOfCXArg_ForGlobals_CXAtomicOps(prgrm, prgrm.GetCXArgFromArray(ast.CXArgumentIndex(valArg.Index)))
	prgrm.CXAtomicOps[expressionIdx].AddInput(prgrm, typeSig)

	typeSig = lastPrevExpression.GetOutputs(prgrm)[0]
	prgrm.CXAtomicOps[expressionIdx].AddOutput(prgrm, typeSig)

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

		var lastExpressionOperatorOutputArg *ast.CXArgument = &ast.CXArgument{}
		if lastExpressionOperatorOutputs[0].Type == ast.TYPE_CXARGUMENT_DEPRECATE {
			lastExpressionOperatorOutputArg = prgrm.GetCXArgFromArray(ast.CXArgumentIndex(lastExpressionOperatorOutputs[0].Meta))
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
		out.TotalSize = opOut.TotalSize
		// out.IsArray = opOut.IsArray
		// out.IsReference = opOut.IsReference
		out.Lengths = opOut.Lengths
		out.Package = prgrm.CXAtomicOps[lastExpressionIdx].Package
		out.PreviouslyDeclared = true

		outIdx := prgrm.AddCXArgInArray(out)
		typeSig := ast.GetCXTypeSignatureRepresentationOfCXArg_ForGlobals_CXAtomicOps(prgrm, prgrm.GetCXArgFromArray(outIdx))
		prgrm.CXAtomicOps[lastExpressionIdx].AddOutput(prgrm, typeSig)

		// we need to create an expression to hold all the modifications
		// that will take place after this if statement
		inp := ast.MakeArgument(symName, lastExprCXLine.FileName, lastExprCXLine.LineNumber)
		inp.SetType(opOut.Type)
		inp.DeclarationSpecifiers = opOut.DeclarationSpecifiers
		inp.StructType = opOut.StructType
		inp.Size = opOut.Size
		inp.TotalSize = opOut.TotalSize
		inp.Package = prgrm.CXAtomicOps[lastExpressionIdx].Package
		inpIdx := prgrm.AddCXArgInArray(inp)

		exprCXLine := ast.MakeCXLineExpression(prgrm, lastExprCXLine.FileName, lastExprCXLine.LineNumber, lastExprCXLine.LineStr)
		expr := ast.MakeAtomicOperatorExpression(prgrm, nil)

		expression, err := prgrm.GetCXAtomicOp(expr.Index)
		if err != nil {
			panic(err)
		}

		expression.Package = prgrm.CXAtomicOps[lastExpressionIdx].Package
		typeSig = ast.GetCXTypeSignatureRepresentationOfCXArg_ForGlobals_CXAtomicOps(prgrm, prgrm.GetCXArgFromArray(inpIdx))
		expression.AddOutput(prgrm, typeSig)
		prevExprs = append(prevExprs, *exprCXLine, *expr)

		lastExpr = prevExprs[len(prevExprs)-1]

		lastExpressionIdx = lastExpr.Index
	}

	leftExprOutput := prgrm.CXAtomicOps[lastExpressionIdx].GetOutputs(prgrm)[0]
	leftExprIdx := leftExprOutput.Meta

	// then left is a first (e.g first.rest) and right is a rest
	// let's check if left is a package
	pkg, err := prgrm.GetCurrentPackage()
	if err != nil {
		panic(err)
	}

	if imp, err := pkg.GetImport(prgrm, prgrm.CXArgs[leftExprIdx].Name); err == nil {
		// the external property will be propagated to the following arguments
		// this way we avoid considering these arguments as module names

		if cxpackages.IsDefaultPackage(prgrm.CXArgs[leftExprIdx].Name) {
			//TODO: constants.ConstCodes[prgrm.CXArgs[leftExprIdx].Name+"."+ident]
			//TODO: only play ConstCodes are used
			//Is used for constant declaration? But only for core packages?
			if code, ok := ConstCodes[prgrm.CXArgs[leftExprIdx].Name+"."+ident]; ok {
				constant := Constants[code]
				valArg := WritePrimary(prgrm, constant.Type, constant.Value, false)

				// TODO: to be removed and implemented with correct process to change the index
				prgrm.CXAtomicOps[lastExpressionIdx].Outputs.Fields[0].Meta = valArg.Index

				return prevExprs
			} else if _, ok := ast.OpCodes[prgrm.CXArgs[leftExprIdx].Name+"."+ident]; ok {
				// then it's a native
				// TODO: we'd be referring to the function itself, not a function call
				// (functions as first-class objects)
				leftExprOutput.Name = prgrm.CXArgs[leftExprIdx].Name + "." + ident
				prgrm.CXArgs[leftExprIdx].Name = prgrm.CXArgs[leftExprIdx].Name + "." + ident

				return prevExprs
			}
		}

		prgrm.CXArgs[leftExprIdx].Package = ast.CXPackageIndex(imp.Index)

		if glbl, err := imp.GetGlobal(prgrm, ident); err == nil {
			var glblArg *ast.CXArgument = &ast.CXArgument{}
			if glbl.Type == ast.TYPE_CXARGUMENT_DEPRECATE {
				glblArg = prgrm.GetCXArgFromArray(ast.CXArgumentIndex(glbl.Meta))
			} else {
				panic("type is not type cx argument deprecate\n\n")
			}
			// then it's a global
			// prevExprs[len(prevExprs)-1].ProgramOutput[0] = glbl
			lastExpressionOutput := prgrm.CXAtomicOps[lastExpressionIdx].GetOutputs(prgrm)[0]

			var lastExpressionOutputArg *ast.CXArgument = &ast.CXArgument{}
			if lastExpressionOutput.Type == ast.TYPE_CXARGUMENT_DEPRECATE {
				lastExpressionOutputArg = prgrm.GetCXArgFromArray(ast.CXArgumentIndex(lastExpressionOutput.Meta))
			} else {
				panic("type is not cx argument deprecate\n\n")
			}

			lastExpressionOutputArg.Name = glblArg.Name
			lastExpressionOutput.Name = glblArg.Name

			lastExpressionOutputArg.Type = glblArg.Type
			lastExpressionOutputArg.StructType = glblArg.StructType
			lastExpressionOutputArg.Size = glblArg.Size
			lastExpressionOutputArg.TotalSize = glblArg.TotalSize
			lastExpressionOutputArg.PointerTargetType = glblArg.PointerTargetType
			lastExpressionOutputArg.IsSlice = glblArg.IsSlice
			lastExpressionOutputArg.IsStruct = glblArg.IsStruct
			lastExpressionOutputArg.Package = glblArg.Package
		} else if fn, err := imp.GetFunction(prgrm, ident); err == nil {
			// then it's a function
			// not sure about this next line
			prgrm.CXAtomicOps[lastExpressionIdx].Outputs = nil
			prgrm.CXAtomicOps[lastExpressionIdx].Operator = ast.CXFunctionIndex(fn.Index)
		} else if strct, err := prgrm.GetStruct(ident, imp.Name); err == nil {
			var lastExpressionOutputArg *ast.CXArgument = &ast.CXArgument{}
			if prgrm.CXAtomicOps[lastExpressionIdx].GetOutputs(prgrm)[0].Type == ast.TYPE_CXARGUMENT_DEPRECATE {
				lastExpressionOutputArg = prgrm.GetCXArgFromArray(ast.CXArgumentIndex(prgrm.CXAtomicOps[lastExpressionIdx].GetOutputs(prgrm)[0].Meta))
			} else {
				panic("type is not cx argument deprecate\n\n")
			}

			lastExpressionOutput := lastExpressionOutputArg
			lastExpressionOutput.StructType = strct
		} else {
			// panic(err)
			fmt.Println(err)
			return nil
		}
	} else {
		// then left is not a package name
		if cxpackages.IsDefaultPackage(prgrm.CXArgs[leftExprIdx].Name) {
			println(ast.CompilationError(prgrm.CXArgs[leftExprIdx].ArgDetails.FileName, prgrm.CXArgs[leftExprIdx].ArgDetails.FileLine),
				fmt.Sprintf("identifier '%s' does not exist",
					prgrm.CXArgs[leftExprIdx].Name))
			os.Exit(constants.CX_COMPILATION_ERROR)
		}
		// then it's a struct
		prgrm.CXArgs[leftExprIdx].IsStruct = true

		leftExprField := ast.MakeArgument(ident, CurrentFile, LineNo)
		leftPkg, err := prgrm.GetPackageFromArray(prgrm.CXArgs[leftExprIdx].Package)
		if err != nil {
			panic(err)
		}
		leftExprField.SetType(types.IDENTIFIER)
		leftExprField.SetPackage(leftPkg)

		leftExprFieldIdx := prgrm.AddCXArgInArray(leftExprField)
		prgrm.CXArgs[leftExprIdx].Fields = append(prgrm.CXArgs[leftExprIdx].Fields, leftExprFieldIdx)
	}

	return prevExprs
}
