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
	var elt *ast.CXArgument

	prevExpressionIdx := prevExprs[len(prevExprs)-1].Index
	prevExpressionOperator := prgrm.GetFunctionFromArray(prgrm.CXAtomicOps[prevExpressionIdx].Operator)

	prevExprCXLine, _ := prgrm.GetPreviousCXLine(prevExprs, len(prevExprs)-1)

	if prevExpressionOperator != nil && len(prgrm.CXAtomicOps[prevExpressionIdx].Outputs) == 0 {
		genName := generateTempVarName(constants.LOCAL_PREFIX)

		prevExpressionOperatorOutputs := prevExpressionOperator.GetOutputs(prgrm)
		prevExpressionOperatorOutput := prgrm.GetCXArgFromArray(prevExpressionOperatorOutputs[0])
		out := ast.MakeArgument(genName, prevExprCXLine.FileName, prevExprCXLine.LineNumber-1).SetType(prevExpressionOperatorOutput.Type)

		out.DeclarationSpecifiers = prevExpressionOperatorOutput.DeclarationSpecifiers
		out.StructType = prevExpressionOperatorOutput.StructType
		out.Size = prevExpressionOperatorOutput.Size
		out.TotalSize = prevExpressionOperatorOutput.TotalSize
		out.Lengths = prevExpressionOperatorOutput.Lengths
		out.IsSlice = prevExpressionOperatorOutput.IsSlice
		out.PreviouslyDeclared = true
		outIdx := prgrm.AddCXArgInArray(out)

		prgrm.CXAtomicOps[prevExpressionIdx].AddOutput(prgrm, outIdx)

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
		prgrm.CXAtomicOps[useExpressionIdx].AddOutput(prgrm, inpIdx)
		prevExprs = append(prevExprs, *useExprCXLine, *useExpr)
	}

	prevExpression, err := prgrm.GetCXAtomicOpFromExpressions(prevExprs, len(prevExprs)-1)
	if err != nil {
		panic(err)
	}

	if len(prgrm.GetCXArgFromArray(prevExpression.Outputs[0]).Fields) > 0 {
		elt = prgrm.GetCXArgFromArray(prgrm.GetCXArgFromArray(prevExpression.Outputs[0]).Fields[len(prgrm.GetCXArgFromArray(prevExpression.Outputs[0]).Fields)-1])
	} else {
		elt = prgrm.GetCXArgFromArray(prevExpression.Outputs[0])
	}

	// elt.IsArray = false
	elt.DereferenceOperations = append(elt.DereferenceOperations, constants.DEREF_ARRAY)
	elt.DeclarationSpecifiers = append(elt.DeclarationSpecifiers, constants.DECL_INDEXING)

	postExpressionIdx := postExprs[len(postExprs)-1].Index
	postExpressionOperator := prgrm.GetFunctionFromArray(prgrm.CXAtomicOps[postExpressionIdx].Operator)
	postExpressionOperatorOutputs := postExpressionOperator.GetOutputs(prgrm)

	if len(prgrm.GetCXArgFromArray(prevExpression.Outputs[0]).Fields) > 0 {
		fldIdx := prgrm.GetCXArgFromArray(prevExpression.Outputs[0]).Fields[len(prgrm.GetCXArgFromArray(prevExpression.Outputs[0]).Fields)-1]

		if postExpressionOperator == nil {
			// expr.AddInput(postExprs[len(postExprs)-1].ProgramOutput[0])
			indexIdx := prgrm.CXAtomicOps[postExpressionIdx].Outputs[0]
			prgrm.CXArgs[fldIdx].Indexes = append(prgrm.CXArgs[fldIdx].Indexes, indexIdx)
		} else {

			sym := ast.MakeArgument(generateTempVarName(constants.LOCAL_PREFIX), CurrentFile, LineNo).SetType(prgrm.GetCXArgFromArray(postExpressionOperatorOutputs[0]).Type)
			sym.Package = prgrm.CXAtomicOps[postExpressionIdx].Package
			sym.PreviouslyDeclared = true
			symIdx := prgrm.AddCXArgInArray(sym)
			prgrm.CXAtomicOps[postExpressionIdx].AddOutput(prgrm, symIdx)

			prevExprs = append(postExprs, prevExprs...)

			prgrm.CXArgs[fldIdx].Indexes = append(prgrm.CXArgs[fldIdx].Indexes, symIdx)
		}
	} else {
		if len(prgrm.CXAtomicOps[postExpressionIdx].Outputs) < 1 {
			// then it's an expression (e.g. i32.add(0, 0))
			// we create a gensym for it
			idxSym := ast.MakeArgument(generateTempVarName(constants.LOCAL_PREFIX), CurrentFile, LineNo).SetType(prgrm.GetCXArgFromArray(postExpressionOperatorOutputs[0]).Type)
			idxSym.Size = prgrm.GetCXArgFromArray(postExpressionOperatorOutputs[0]).Size
			idxSym.TotalSize = ast.GetArgSize(prgrm, prgrm.GetCXArgFromArray(postExpressionOperatorOutputs[0]))

			idxSym.Package = prgrm.CXAtomicOps[postExpressionIdx].Package
			idxSym.PreviouslyDeclared = true

			idxSymIdx := prgrm.AddCXArgInArray(idxSym)
			prgrm.CXAtomicOps[postExpressionIdx].Outputs = append(prgrm.CXAtomicOps[postExpressionIdx].Outputs, idxSymIdx)

			prevExpressionOutput := prgrm.GetCXArgFromArray(prevExpression.Outputs[0])
			prevExpressionOutput.Indexes = append(prevExpressionOutput.Indexes, idxSymIdx)

			// we push the index expression
			prevExprs = append(postExprs, prevExprs...)
		} else {
			prevOuts := prevExpression.Outputs
			postOuts := prgrm.CXAtomicOps[postExpressionIdx].Outputs
			prgrm.CXArgs[prevOuts[0]].Indexes = append(prgrm.CXArgs[prevOuts[0]].Indexes, postOuts[0])
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
	if prevExpression.Outputs != nil && len(prgrm.GetCXArgFromArray(prevExpression.Outputs[0]).Fields) > 0 {
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
		if opCode, ok := ast.OpCodes[prgrm.GetCXArgFromArray(prevExpression.Outputs[0]).Name]; ok {
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
	if lastPrevExpression.Outputs != nil && len(prgrm.GetCXArgFromArray(lastPrevExpression.Outputs[0]).Fields) > 0 {
		// then it's a method
		// prevExprs[len(prevExprs) - 1].IsMethodCall = true

	} else if lastPrevExpressionOperator == nil {
		if opCode, ok := ast.OpCodes[prgrm.GetCXArgFromArray(lastPrevExpression.Outputs[0]).Name]; ok {
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
	val := WritePrimary(prgrm, types.I32, valB[:], false)

	lastPrevExpression, err := prgrm.GetCXAtomicOpFromExpressions(prevExprs, len(prevExprs)-1)
	if err != nil {
		panic(err)
	}

	valAtomicOp, err := prgrm.GetCXAtomicOpFromExpressions(val, len(val)-1)
	if err != nil {
		panic(err)
	}

	expressionIdx := expr.Index
	prgrm.CXAtomicOps[expressionIdx].Package = ast.CXPackageIndex(pkg.Index)
	prgrm.CXAtomicOps[expressionIdx].AddInput(prgrm, lastPrevExpression.Outputs[0])
	prgrm.CXAtomicOps[expressionIdx].AddInput(prgrm, valAtomicOp.Outputs[0])
	prgrm.CXAtomicOps[expressionIdx].AddOutput(prgrm, lastPrevExpression.Outputs[0])

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
		opOut := prgrm.GetCXArgFromArray(lastExpressionOperatorOutputs[0])
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
		prgrm.CXAtomicOps[lastExpressionIdx].Outputs = append(prgrm.CXAtomicOps[lastExpressionIdx].Outputs, outIdx)

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
		expression.AddOutput(prgrm, inpIdx)
		prevExprs = append(prevExprs, *exprCXLine, *expr)

		lastExpr = prevExprs[len(prevExprs)-1]

		lastExpressionIdx = lastExpr.Index
	}

	leftExprIdx := prgrm.CXAtomicOps[lastExpressionIdx].Outputs[0]

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
				val := WritePrimary(prgrm, constant.Type, constant.Value, false)
				valAtomicOp, err := prgrm.GetCXAtomicOpFromExpressions(val, 0)
				if err != nil {
					panic(err)
				}

				prgrm.CXAtomicOps[lastExpressionIdx].Outputs[0] = valAtomicOp.Outputs[0]

				return prevExprs
			} else if _, ok := ast.OpCodes[prgrm.CXArgs[leftExprIdx].Name+"."+ident]; ok {
				// then it's a native
				// TODO: we'd be referring to the function itself, not a function call
				// (functions as first-class objects)
				prgrm.CXArgs[leftExprIdx].Name = prgrm.CXArgs[leftExprIdx].Name + "." + ident

				return prevExprs
			}
		}

		prgrm.CXArgs[leftExprIdx].Package = ast.CXPackageIndex(imp.Index)

		if glbl, err := imp.GetGlobal(prgrm, ident); err == nil {
			// then it's a global
			// prevExprs[len(prevExprs)-1].ProgramOutput[0] = glbl
			lastExpressionOutput := prgrm.GetCXArgFromArray(prgrm.CXAtomicOps[lastExpressionIdx].Outputs[0])
			lastExpressionOutput.Name = glbl.Name
			lastExpressionOutput.Type = glbl.Type
			lastExpressionOutput.StructType = glbl.StructType
			lastExpressionOutput.Size = glbl.Size
			lastExpressionOutput.TotalSize = glbl.TotalSize
			lastExpressionOutput.PointerTargetType = glbl.PointerTargetType
			lastExpressionOutput.IsSlice = glbl.IsSlice
			lastExpressionOutput.IsStruct = glbl.IsStruct
			lastExpressionOutput.Package = glbl.Package
		} else if fn, err := imp.GetFunction(prgrm, ident); err == nil {
			// then it's a function
			// not sure about this next line
			prgrm.CXAtomicOps[lastExpressionIdx].Outputs = nil
			prgrm.CXAtomicOps[lastExpressionIdx].Operator = ast.CXFunctionIndex(fn.Index)
		} else if strct, err := prgrm.GetStruct(ident, imp.Name); err == nil {
			lastExpressionOutput := prgrm.GetCXArgFromArray(prgrm.CXAtomicOps[lastExpressionIdx].Outputs[0])
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
