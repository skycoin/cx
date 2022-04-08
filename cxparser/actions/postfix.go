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

	prevExprAtomicOpIdx := prevExprs[len(prevExprs)-1].Index
	prevExprAtomicOpOperator := prgrm.GetFunctionFromArray(prgrm.CXAtomicOps[prevExprAtomicOpIdx].Operator)

	prevExprCXLine, _ := prgrm.GetPreviousCXLine(prevExprs, len(prevExprs)-1)

	if prevExprAtomicOpOperator != nil && len(prgrm.CXAtomicOps[prevExprAtomicOpIdx].Outputs) == 0 {
		genName := MakeGenSym(constants.LOCAL_PREFIX)

		prevExprAtomicOpOperatorOutput := prgrm.GetCXArgFromArray(prevExprAtomicOpOperator.Outputs[0])
		out := ast.MakeArgument(genName, prevExprCXLine.FileName, prevExprCXLine.LineNumber-1).SetType(prevExprAtomicOpOperatorOutput.Type)

		out.DeclarationSpecifiers = prevExprAtomicOpOperatorOutput.DeclarationSpecifiers
		out.StructType = prevExprAtomicOpOperatorOutput.StructType
		out.Size = prevExprAtomicOpOperatorOutput.Size
		out.TotalSize = prevExprAtomicOpOperatorOutput.TotalSize
		out.Lengths = prevExprAtomicOpOperatorOutput.Lengths
		out.IsSlice = prevExprAtomicOpOperatorOutput.IsSlice
		out.PreviouslyDeclared = true
		outIdx := prgrm.AddCXArgInArray(out)

		prgrm.CXAtomicOps[prevExprAtomicOpIdx].AddOutput(prgrm, outIdx)

		inp := ast.MakeArgument(genName, prevExprCXLine.FileName, prevExprCXLine.LineNumber).SetType(prevExprAtomicOpOperatorOutput.Type)

		inp.DeclarationSpecifiers = prevExprAtomicOpOperatorOutput.DeclarationSpecifiers
		inp.StructType = prevExprAtomicOpOperatorOutput.StructType
		inp.Size = prevExprAtomicOpOperatorOutput.Size
		inp.TotalSize = prevExprAtomicOpOperatorOutput.TotalSize
		inp.Lengths = prevExprAtomicOpOperatorOutput.Lengths
		inp.IsSlice = prevExprAtomicOpOperatorOutput.IsSlice
		inp.PreviouslyDeclared = true
		inpIdx := prgrm.AddCXArgInArray(inp)

		useExprCXLine := ast.MakeCXLineExpression(prgrm, prevExprCXLine.FileName, prevExprCXLine.LineNumber, prevExprCXLine.LineStr)
		useExpr := ast.MakeAtomicOperatorExpression(prgrm, nil)
		useExprAtomicOpIdx := useExpr.Index

		prgrm.CXAtomicOps[useExprAtomicOpIdx].Package = prgrm.CXAtomicOps[prevExprAtomicOpIdx].Package
		prgrm.CXAtomicOps[useExprAtomicOpIdx].AddOutput(prgrm, inpIdx)
		prevExprs = append(prevExprs, *useExprCXLine, *useExpr)
	}

	prevExpr2AtomicOp, err := prgrm.GetCXAtomicOpFromExpressions(prevExprs, len(prevExprs)-1)
	if err != nil {
		panic(err)
	}

	if len(prgrm.GetCXArgFromArray(prevExpr2AtomicOp.Outputs[0]).Fields) > 0 {
		elt = prgrm.GetCXArgFromArray(prgrm.GetCXArgFromArray(prevExpr2AtomicOp.Outputs[0]).Fields[len(prgrm.GetCXArgFromArray(prevExpr2AtomicOp.Outputs[0]).Fields)-1])
	} else {
		elt = prgrm.GetCXArgFromArray(prevExpr2AtomicOp.Outputs[0])
	}

	// elt.IsArray = false
	elt.DereferenceOperations = append(elt.DereferenceOperations, constants.DEREF_ARRAY)
	elt.DeclarationSpecifiers = append(elt.DeclarationSpecifiers, constants.DECL_INDEXING)

	postExprsAtomicOpIdx := postExprs[len(postExprs)-1].Index
	postExprsAtomicOpOperator := prgrm.GetFunctionFromArray(prgrm.CXAtomicOps[postExprsAtomicOpIdx].Operator)

	if len(prgrm.GetCXArgFromArray(prevExpr2AtomicOp.Outputs[0]).Fields) > 0 {
		fldIdx := prgrm.GetCXArgFromArray(prevExpr2AtomicOp.Outputs[0]).Fields[len(prgrm.GetCXArgFromArray(prevExpr2AtomicOp.Outputs[0]).Fields)-1]

		if postExprsAtomicOpOperator == nil {
			// expr.AddInput(postExprs[len(postExprs)-1].ProgramOutput[0])
			indexIdx := prgrm.CXAtomicOps[postExprsAtomicOpIdx].Outputs[0]
			prgrm.CXArgs[fldIdx].Indexes = append(prgrm.CXArgs[fldIdx].Indexes, indexIdx)
		} else {
			sym := ast.MakeArgument(MakeGenSym(constants.LOCAL_PREFIX), CurrentFile, LineNo).SetType(prgrm.GetCXArgFromArray(postExprsAtomicOpOperator.Outputs[0]).Type)
			sym.Package = prgrm.CXAtomicOps[postExprsAtomicOpIdx].Package
			sym.PreviouslyDeclared = true
			symIdx := prgrm.AddCXArgInArray(sym)
			prgrm.CXAtomicOps[postExprsAtomicOpIdx].AddOutput(prgrm, symIdx)

			prevExprs = append(postExprs, prevExprs...)

			prgrm.CXArgs[fldIdx].Indexes = append(prgrm.CXArgs[fldIdx].Indexes, symIdx)
		}
	} else {
		if len(prgrm.CXAtomicOps[postExprsAtomicOpIdx].Outputs) < 1 {
			// then it's an expression (e.g. i32.add(0, 0))
			// we create a gensym for it
			idxSym := ast.MakeArgument(MakeGenSym(constants.LOCAL_PREFIX), CurrentFile, LineNo).SetType(prgrm.GetCXArgFromArray(postExprsAtomicOpOperator.Outputs[0]).Type)
			idxSym.Size = prgrm.GetCXArgFromArray(postExprsAtomicOpOperator.Outputs[0]).Size
			idxSym.TotalSize = ast.GetSize(prgrm, prgrm.GetCXArgFromArray(postExprsAtomicOpOperator.Outputs[0]))

			idxSym.Package = prgrm.CXAtomicOps[postExprsAtomicOpIdx].Package
			idxSym.PreviouslyDeclared = true

			idxSymIdx := prgrm.AddCXArgInArray(idxSym)
			prgrm.CXAtomicOps[postExprsAtomicOpIdx].Outputs = append(prgrm.CXAtomicOps[postExprsAtomicOpIdx].Outputs, idxSymIdx)

			prevExpr2AtomicOpOutput := prgrm.GetCXArgFromArray(prevExpr2AtomicOp.Outputs[0])
			prevExpr2AtomicOpOutput.Indexes = append(prevExpr2AtomicOpOutput.Indexes, idxSymIdx)

			// we push the index expression
			prevExprs = append(postExprs, prevExprs...)
		} else {
			prevOuts := prevExpr2AtomicOp.Outputs
			postOuts := prgrm.CXAtomicOps[postExprsAtomicOpIdx].Outputs
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
	cxAtomicOp, err := prgrm.GetCXAtomicOp(expr.Index)
	if err != nil {
		panic(err)
	}

	pkg, err := prgrm.GetCurrentPackage()
	if err != nil {
		panic(err)
	}
	cxAtomicOp.Package = ast.CXPackageIndex(pkg.Index)
	return []ast.CXExpression{*exprCXLine, *expr}
}

// PostfixExpressionEmptyFunCall handles postfix expression empty function calls
// or the function calls that doesn't have any input args.
//
// Input arguments description:
// 	prgrm - a CXProgram that contains all the data and arrays of the program.
// 	prevExprs - the previous expressions that compose the postfix expression.
func PostfixExpressionEmptyFunCall(prgrm *ast.CXProgram, prevExprs []ast.CXExpression) []ast.CXExpression {
	prevExprsAtomicOp, err := prgrm.GetCXAtomicOpFromExpressions(prevExprs, len(prevExprs)-1)
	if err != nil {
		panic(err)
	}
	prevExprsAtomicOpOperator := prgrm.GetFunctionFromArray(prevExprsAtomicOp.Operator)

	firstPrevExprsAtomicOp, err := prgrm.GetCXAtomicOpFromExpressions(prevExprs, 0)
	if err != nil {
		panic(err)
	}

	if prevExprsAtomicOp.Outputs != nil && len(prgrm.GetCXArgFromArray(prevExprsAtomicOp.Outputs[0]).Fields) > 0 {
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

	} else if prevExprsAtomicOpOperator == nil {
		if opCode, ok := ast.OpCodes[prgrm.GetCXArgFromArray(prevExprsAtomicOp.Outputs[0]).Name]; ok {
			if pkg, err := prgrm.GetCurrentPackage(); err == nil {
				firstPrevExprsAtomicOp.Package = ast.CXPackageIndex(pkg.Index)
			}
			firstPrevExprsAtomicOp.Outputs = nil
			opIdx := prgrm.AddNativeFunctionInArray(ast.Natives[opCode])
			firstPrevExprsAtomicOp.Operator = opIdx
		}

		firstPrevExprsAtomicOp.Inputs = nil
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
	lastPrevExprsAtomicOp, err := prgrm.GetCXAtomicOpFromExpressions(prevExprs, len(prevExprs)-1)
	if err != nil {
		panic(err)
	}
	lastPrevExprsAtomicOpOperator := prgrm.GetFunctionFromArray(lastPrevExprsAtomicOp.Operator)

	firstPrevExprsAtomicOp, err := prgrm.GetCXAtomicOpFromExpressions(prevExprs, 0)
	if err != nil {
		panic(err)
	}
	if lastPrevExprsAtomicOp.Outputs != nil && len(prgrm.GetCXArgFromArray(lastPrevExprsAtomicOp.Outputs[0]).Fields) > 0 {
		// then it's a method
		// prevExprs[len(prevExprs) - 1].IsMethodCall = true

	} else if lastPrevExprsAtomicOpOperator == nil {
		if opCode, ok := ast.OpCodes[prgrm.GetCXArgFromArray(lastPrevExprsAtomicOp.Outputs[0]).Name]; ok {
			if pkg, err := prgrm.GetCurrentPackage(); err == nil {
				firstPrevExprsAtomicOp.Package = ast.CXPackageIndex(pkg.Index)
			}
			firstPrevExprsAtomicOp.Outputs = nil
			opIdx := prgrm.AddNativeFunctionInArray(ast.Natives[opCode])
			firstPrevExprsAtomicOp.Operator = opIdx
		}

		firstPrevExprsAtomicOp.Inputs = nil
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

	cxAtomicOp, err := prgrm.GetCXAtomicOp(expr.Index)
	if err != nil {
		panic(err)
	}

	lastPrevExprsAtomicOp, err := prgrm.GetCXAtomicOpFromExpressions(prevExprs, len(prevExprs)-1)
	if err != nil {
		panic(err)
	}

	valAtomicOp, err := prgrm.GetCXAtomicOpFromExpressions(val, len(val)-1)
	if err != nil {
		panic(err)
	}

	cxAtomicOp.Package = ast.CXPackageIndex(pkg.Index)

	cxAtomicOp.AddInput(prgrm, lastPrevExprsAtomicOp.Outputs[0])
	cxAtomicOp.AddInput(prgrm, valAtomicOp.Outputs[0])
	cxAtomicOp.AddOutput(prgrm, lastPrevExprsAtomicOp.Outputs[0])

	// exprs := append(prevExprs, expr)
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

	lastExprAtomicOp, _, _, err := prgrm.GetOperation(&lastExpr)
	if err != nil {
		panic(err)
	}
	lastExprAtomicOpOperator := prgrm.GetFunctionFromArray(lastExprAtomicOp.Operator)

	lastExprCXLine, _ := prgrm.GetPreviousCXLine(prevExprs, len(prevExprs)-1)

	// Then it's a function call, e.g. foo().fld
	// and we need to create some auxiliary variables to hold the result from
	// the function call
	if lastExprAtomicOpOperator != nil {
		opOut := prgrm.GetCXArgFromArray(lastExprAtomicOpOperator.Outputs[0])
		symName := MakeGenSym(constants.LOCAL_PREFIX)

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
		out.Package = lastExprAtomicOp.Package
		out.PreviouslyDeclared = true
		out.IsInnerArg = true

		outIdx := prgrm.AddCXArgInArray(out)
		lastExprAtomicOp.Outputs = append(lastExprAtomicOp.Outputs, outIdx)

		// we need to create an expression to hold all the modifications
		// that will take place after this if statement
		inp := ast.MakeArgument(symName, lastExprCXLine.FileName, lastExprCXLine.LineNumber)
		inp.SetType(opOut.Type)
		inp.DeclarationSpecifiers = opOut.DeclarationSpecifiers
		inp.StructType = opOut.StructType
		inp.Size = opOut.Size
		inp.TotalSize = opOut.TotalSize
		inp.Package = lastExprAtomicOp.Package
		inp.IsInnerArg = true
		inpIdx := prgrm.AddCXArgInArray(inp)

		exprCXLine := ast.MakeCXLineExpression(prgrm, lastExprCXLine.FileName, lastExprCXLine.LineNumber, lastExprCXLine.LineStr)
		expr := ast.MakeAtomicOperatorExpression(prgrm, nil)

		exprAtomicOp, _, _, err := prgrm.GetOperation(expr)
		if err != nil {
			panic(err)
		}

		exprAtomicOp.Package = lastExprAtomicOp.Package
		exprAtomicOp.AddOutput(prgrm, inpIdx)
		prevExprs = append(prevExprs, *exprCXLine, *expr)

		lastExpr = prevExprs[len(prevExprs)-1]

		lastExprAtomicOp, _, _, err = prgrm.GetOperation(&lastExpr)
		if err != nil {
			panic(err)
		}
	}

	leftExprIdx := lastExprAtomicOp.Outputs[0]

	// If the left already is a rest (e.g. "var" in "pkg.var"), then
	// it can't be a package name and we propagate the property to
	//  the right side.
	if prgrm.CXArgs[leftExprIdx].IsInnerArg {
		prgrm.CXArgs[leftExprIdx].IsStruct = true
		fld := ast.MakeArgument(ident, CurrentFile, LineNo)
		leftPkg, err := prgrm.GetPackageFromArray(prgrm.CXArgs[leftExprIdx].Package)
		if err != nil {
			panic(err)
		}

		fld.SetType(types.IDENTIFIER).SetPackage(leftPkg)
		fldIdx := prgrm.AddCXArgInArray(fld)
		prgrm.CXArgs[leftExprIdx].Fields = append(prgrm.CXArgs[leftExprIdx].Fields, fldIdx)

		return prevExprs
	}

	prgrm.CXArgs[leftExprIdx].IsInnerArg = true
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

				lastExprAtomicOp.Outputs[0] = valAtomicOp.Outputs[0]

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
			lastExprAtomicOpOutput := prgrm.GetCXArgFromArray(lastExprAtomicOp.Outputs[0])
			lastExprAtomicOpOutput.Name = glbl.Name
			lastExprAtomicOpOutput.Type = glbl.Type
			lastExprAtomicOpOutput.StructType = glbl.StructType
			lastExprAtomicOpOutput.Size = glbl.Size
			lastExprAtomicOpOutput.TotalSize = glbl.TotalSize
			lastExprAtomicOpOutput.PointerTargetType = glbl.PointerTargetType
			lastExprAtomicOpOutput.IsSlice = glbl.IsSlice
			lastExprAtomicOpOutput.IsStruct = glbl.IsStruct
			lastExprAtomicOpOutput.Package = glbl.Package
		} else if fn, err := imp.GetFunction(prgrm, ident); err == nil {
			// then it's a function
			// not sure about this next line
			lastExprAtomicOp.Outputs = nil
			lastExprAtomicOp.Operator = ast.CXFunctionIndex(fn.Index)
		} else if strct, err := prgrm.GetStruct(ident, imp.Name); err == nil {
			lastExprAtomicOpOutput := prgrm.GetCXArgFromArray(lastExprAtomicOp.Outputs[0])
			lastExprAtomicOpOutput.StructType = strct
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
