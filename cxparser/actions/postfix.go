package actions

import (
	"fmt"
	"os"

	"github.com/skycoin/cx/cx/ast"
	"github.com/skycoin/cx/cx/constants"
	"github.com/skycoin/cx/cx/types"
	constants2 "github.com/skycoin/cx/cxparser/constants"
)

// PostfixExpressionArray...
//
func PostfixExpressionArray(prgrm *ast.CXProgram, prevExprs []*ast.CXExpression, postExprs []*ast.CXExpression) []*ast.CXExpression {
	var elt *ast.CXArgument

	prevExprAtomicOp, err := prgrm.GetCXAtomicOpFromExpressions(prevExprs, len(prevExprs)-1)
	if err != nil {
		panic(err)
	}

	prevExprCXLine, _ := prgrm.GetPreviousCXLine(prevExprs, len(prevExprs)-1)

	if prevExprAtomicOp.Operator != nil && len(prevExprAtomicOp.Outputs) == 0 {
		genName := MakeGenSym(constants.LOCAL_PREFIX)

		prevExprAtomicOpOperatorOutput := prgrm.GetCXArgFromArray(prevExprAtomicOp.Operator.Outputs[0])
		out := ast.MakeArgument(genName, prevExprCXLine.FileName, prevExprCXLine.LineNumber-1).AddType(prevExprAtomicOpOperatorOutput.Type)

		out.DeclarationSpecifiers = prevExprAtomicOpOperatorOutput.DeclarationSpecifiers
		out.StructType = prevExprAtomicOpOperatorOutput.StructType
		out.Size = prevExprAtomicOpOperatorOutput.Size
		out.TotalSize = prevExprAtomicOpOperatorOutput.TotalSize
		out.Lengths = prevExprAtomicOpOperatorOutput.Lengths
		out.IsSlice = prevExprAtomicOpOperatorOutput.IsSlice
		out.PreviouslyDeclared = true

		prevExprAtomicOp.AddOutput(out)

		inp := ast.MakeArgument(genName, prevExprCXLine.FileName, prevExprCXLine.LineNumber).AddType(prevExprAtomicOpOperatorOutput.Type)

		inp.DeclarationSpecifiers = prevExprAtomicOpOperatorOutput.DeclarationSpecifiers
		inp.StructType = prevExprAtomicOpOperatorOutput.StructType
		inp.Size = prevExprAtomicOpOperatorOutput.Size
		inp.TotalSize = prevExprAtomicOpOperatorOutput.TotalSize
		inp.Lengths = prevExprAtomicOpOperatorOutput.Lengths
		inp.IsSlice = prevExprAtomicOpOperatorOutput.IsSlice
		inp.PreviouslyDeclared = true

		useExprCXLine := ast.MakeCXLineExpression(prgrm, prevExprCXLine.FileName, prevExprCXLine.LineNumber, prevExprCXLine.LineStr)
		useExpr := ast.MakeAtomicOperatorExpression(prgrm, nil)
		useExprAtomicOp, _, _, err := prgrm.GetOperation(useExpr)
		if err != nil {
			panic(err)
		}

		useExprAtomicOp.Package = prevExprAtomicOp.Package
		useExprAtomicOp.AddOutput(inp)
		prevExprs = append(prevExprs, useExprCXLine, useExpr)
	}

	prevExpr2AtomicOp, err := prgrm.GetCXAtomicOpFromExpressions(prevExprs, len(prevExprs)-1)
	if err != nil {
		panic(err)
	}

	if len(prevExpr2AtomicOp.Outputs[0].Fields) > 0 {
		elt = prevExpr2AtomicOp.Outputs[0].Fields[len(prevExpr2AtomicOp.Outputs[0].Fields)-1]
	} else {
		elt = prevExpr2AtomicOp.Outputs[0]
	}

	// elt.IsArray = false
	elt.DereferenceOperations = append(elt.DereferenceOperations, constants.DEREF_ARRAY)
	elt.DeclarationSpecifiers = append(elt.DeclarationSpecifiers, constants.DECL_INDEXING)

	postExprsAtomicOp, err := prgrm.GetCXAtomicOpFromExpressions(postExprs, len(postExprs)-1)
	if err != nil {
		panic(err)
	}

	if len(prevExpr2AtomicOp.Outputs[0].Fields) > 0 {
		fld := prevExpr2AtomicOp.Outputs[0].Fields[len(prevExpr2AtomicOp.Outputs[0].Fields)-1]

		if postExprsAtomicOp.Operator == nil {
			// expr.AddInput(postExprs[len(postExprs)-1].ProgramOutput[0])
			fld.Indexes = append(fld.Indexes, postExprsAtomicOp.Outputs[0])
		} else {
			sym := ast.MakeArgument(MakeGenSym(constants.LOCAL_PREFIX), CurrentFile, LineNo).AddType(prgrm.GetCXArgFromArray(postExprsAtomicOp.Operator.Outputs[0]).Type)
			sym.Package = postExprsAtomicOp.Package
			sym.PreviouslyDeclared = true
			postExprsAtomicOp.AddOutput(sym)

			prevExprs = append(postExprs, prevExprs...)

			fld.Indexes = append(fld.Indexes, sym)
			// expr.AddInput(sym)
		}
	} else {
		if len(postExprsAtomicOp.Outputs) < 1 {
			// then it's an expression (e.g. i32.add(0, 0))
			// we create a gensym for it
			idxSym := ast.MakeArgument(MakeGenSym(constants.LOCAL_PREFIX), CurrentFile, LineNo).AddType(prgrm.GetCXArgFromArray(postExprsAtomicOp.Operator.Outputs[0]).Type)
			idxSym.Size = prgrm.GetCXArgFromArray(postExprsAtomicOp.Operator.Outputs[0]).Size
			idxSym.TotalSize = ast.GetSize(prgrm.GetCXArgFromArray(postExprsAtomicOp.Operator.Outputs[0]))

			idxSym.Package = postExprsAtomicOp.Package
			idxSym.PreviouslyDeclared = true
			postExprsAtomicOp.Outputs = append(postExprsAtomicOp.Outputs, idxSym)

			prevExpr2AtomicOp.Outputs[0].Indexes = append(prevExpr2AtomicOp.Outputs[0].Indexes, idxSym)

			// we push the index expression
			prevExprs = append(postExprs, prevExprs...)
		} else {
			prevOuts := prevExpr2AtomicOp.Outputs
			postOuts := postExprsAtomicOp.Outputs
			prevOuts[0].Indexes = append(prevOuts[0].Indexes, postOuts[0])
		}
	}

	return prevExprs
}

func PostfixExpressionNative(prgrm *ast.CXProgram, typeCode types.Code, opStrCode string) []*ast.CXExpression {
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
	cxAtomicOp, _, _, err := prgrm.GetOperation(expr)
	if err != nil {
		panic(err)
	}

	pkg, err := prgrm.GetCurrentPackage()
	if err != nil {
		panic(err)
	}
	cxAtomicOp.Package = ast.CXPackageIndex(pkg.Index)
	return []*ast.CXExpression{exprCXLine, expr}
}

func PostfixExpressionEmptyFunCall(prgrm *ast.CXProgram, prevExprs []*ast.CXExpression) []*ast.CXExpression {
	prevExprsAtomicOp, err := prgrm.GetCXAtomicOpFromExpressions(prevExprs, len(prevExprs)-1)
	if err != nil {
		panic(err)
	}

	firstPrevExprsAtomicOp, err := prgrm.GetCXAtomicOpFromExpressions(prevExprs, 0)
	if err != nil {
		panic(err)
	}

	if prevExprsAtomicOp.Outputs != nil && len(prevExprsAtomicOp.Outputs[0].Fields) > 0 {
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

	} else if prevExprsAtomicOp.Operator == nil {
		if opCode, ok := ast.OpCodes[prevExprsAtomicOp.Outputs[0].Name]; ok {
			if pkg, err := prgrm.GetCurrentPackage(); err == nil {
				firstPrevExprsAtomicOp.Package = ast.CXPackageIndex(pkg.Index)
			}
			firstPrevExprsAtomicOp.Outputs = nil
			firstPrevExprsAtomicOp.Operator = ast.Natives[opCode]
		}

		firstPrevExprsAtomicOp.Inputs = nil
	}

	return FunctionCall(prgrm, prevExprs, nil)
}

func PostfixExpressionFunCall(prgrm *ast.CXProgram, prevExprs []*ast.CXExpression, args []*ast.CXExpression) []*ast.CXExpression {
	lastPrevExprsAtomicOp, err := prgrm.GetCXAtomicOpFromExpressions(prevExprs, len(prevExprs)-1)
	if err != nil {
		panic(err)
	}

	firstPrevExprsAtomicOp, err := prgrm.GetCXAtomicOpFromExpressions(prevExprs, 0)
	if err != nil {
		panic(err)
	}

	if lastPrevExprsAtomicOp.Outputs != nil && len(lastPrevExprsAtomicOp.Outputs[0].Fields) > 0 {
		// then it's a method
		// prevExprs[len(prevExprs) - 1].IsMethodCall = true

	} else if lastPrevExprsAtomicOp.Operator == nil {
		if opCode, ok := ast.OpCodes[lastPrevExprsAtomicOp.Outputs[0].Name]; ok {
			if pkg, err := prgrm.GetCurrentPackage(); err == nil {
				firstPrevExprsAtomicOp.Package = ast.CXPackageIndex(pkg.Index)
			}
			firstPrevExprsAtomicOp.Outputs = nil
			firstPrevExprsAtomicOp.Operator = ast.Natives[opCode]
		}

		firstPrevExprsAtomicOp.Inputs = nil
	}

	return FunctionCall(prgrm, prevExprs, args)
}

func PostfixExpressionIncDec(prgrm *ast.CXProgram, prevExprs []*ast.CXExpression, isInc bool) []*ast.CXExpression {
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

	cxAtomicOp, _, _, err := prgrm.GetOperation(expr)
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

	cxAtomicOp.AddInput(lastPrevExprsAtomicOp.Outputs[0])
	cxAtomicOp.AddInput(valAtomicOp.Outputs[0])
	cxAtomicOp.AddOutput(lastPrevExprsAtomicOp.Outputs[0])

	// exprs := append(prevExprs, expr)
	exprs := append([]*ast.CXExpression{}, exprCXLine, expr)
	return exprs
}

// PostfixExpressionField handles the dot notation that can follow an identifier.
// Examples are: `foo.bar`, `foo().bar`, `pkg.foo`
func PostfixExpressionField(prgrm *ast.CXProgram, prevExprs []*ast.CXExpression, ident string) []*ast.CXExpression {
	lastExpr := prevExprs[len(prevExprs)-1]

	lastExprAtomicOp, _, _, err := prgrm.GetOperation(lastExpr)
	if err != nil {
		panic(err)
	}

	lastExprCXLine, _ := prgrm.GetPreviousCXLine(prevExprs, len(prevExprs)-1)

	// Then it's a function call, e.g. foo().fld
	// and we need to create some auxiliary variables to hold the result from
	// the function call
	if lastExprAtomicOp.Operator != nil {
		opOut := prgrm.GetCXArgFromArray(lastExprAtomicOp.Operator.Outputs[0])
		symName := MakeGenSym(constants.LOCAL_PREFIX)

		// we associate the result of the function call to the aux variable
		out := ast.MakeArgument(symName, lastExprCXLine.FileName, lastExprCXLine.LineNumber).AddType(opOut.Type)
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

		lastExprAtomicOp.Outputs = append(lastExprAtomicOp.Outputs, out)

		// we need to create an expression to hold all the modifications
		// that will take place after this if statement
		inp := ast.MakeArgument(symName, lastExprCXLine.FileName, lastExprCXLine.LineNumber).AddType(opOut.Type)
		inp.DeclarationSpecifiers = opOut.DeclarationSpecifiers
		inp.StructType = opOut.StructType
		inp.Size = opOut.Size
		inp.TotalSize = opOut.TotalSize
		inp.Package = lastExprAtomicOp.Package
		inp.IsInnerArg = true

		exprCXLine := ast.MakeCXLineExpression(prgrm, lastExprCXLine.FileName, lastExprCXLine.LineNumber, lastExprCXLine.LineStr)
		expr := ast.MakeAtomicOperatorExpression(prgrm, nil)

		exprAtomicOp, _, _, err := prgrm.GetOperation(expr)
		if err != nil {
			panic(err)
		}

		exprAtomicOp.Package = lastExprAtomicOp.Package
		exprAtomicOp.AddOutput(inp)
		prevExprs = append(prevExprs, exprCXLine, expr)

		lastExpr = prevExprs[len(prevExprs)-1]

		lastExprAtomicOp, _, _, err = prgrm.GetOperation(lastExpr)
		if err != nil {
			panic(err)
		}
	}

	left := lastExprAtomicOp.Outputs[0]

	// If the left already is a rest (e.g. "var" in "pkg.var"), then
	// it can't be a package name and we propagate the property to
	//  the right side.
	if left.IsInnerArg {
		// right.IsInnerArg = true
		// left.DereferenceOperations = append(left.DereferenceOperations, cxcore.DEREF_FIELD)
		left.IsStruct = true
		fld := ast.MakeArgument(ident, CurrentFile, LineNo)
		leftPkg, err := prgrm.GetPackageFromArray(left.Package)
		if err != nil {
			panic(err)
		}

		fld.AddType(types.IDENTIFIER).AddPackage(leftPkg)
		left.Fields = append(left.Fields, fld)
		return prevExprs
	}

	left.IsInnerArg = true
	// then left is a first (e.g first.rest) and right is a rest
	// let's check if left is a package
	pkg, err := prgrm.GetCurrentPackage()
	if err != nil {
		panic(err)
	}

	if imp, err := pkg.GetImport(prgrm, left.Name); err == nil {
		// the external property will be propagated to the following arguments
		// this way we avoid considering these arguments as module names

		if constants2.IsCorePackage(left.Name) {

			//TODO: constants.ConstCodes[left.Name+"."+ident]
			//TODO: only play ConstCodes are used
			//Is used for constant declaration? But only for core packages?
			if code, ok := ConstCodes[left.Name+"."+ident]; ok {
				constant := Constants[code]
				val := WritePrimary(prgrm, constant.Type, constant.Value, false)
				valAtomicOp, err := prgrm.GetCXAtomicOpFromExpressions(val, 0)
				if err != nil {
					panic(err)
				}

				lastExprAtomicOp.Outputs[0] = valAtomicOp.Outputs[0]
				return prevExprs
			} else if _, ok := ast.OpCodes[left.Name+"."+ident]; ok {
				// then it's a native
				// TODO: we'd be referring to the function itself, not a function call
				// (functions as first-class objects)
				left.Name = left.Name + "." + ident
				return prevExprs
			}
		}

		left.Package = ast.CXPackageIndex(imp.Index)

		if glbl, err := imp.GetGlobal(prgrm, ident); err == nil {
			// then it's a global
			// prevExprs[len(prevExprs)-1].ProgramOutput[0] = glbl
			lastExprAtomicOp.Outputs[0].Name = glbl.Name
			lastExprAtomicOp.Outputs[0].Type = glbl.Type
			lastExprAtomicOp.Outputs[0].StructType = glbl.StructType
			lastExprAtomicOp.Outputs[0].Size = glbl.Size
			lastExprAtomicOp.Outputs[0].TotalSize = glbl.TotalSize
			lastExprAtomicOp.Outputs[0].PointerTargetType = glbl.PointerTargetType
			lastExprAtomicOp.Outputs[0].IsSlice = glbl.IsSlice
			lastExprAtomicOp.Outputs[0].IsStruct = glbl.IsStruct
			lastExprAtomicOp.Outputs[0].Package = glbl.Package
		} else if fn, err := imp.GetFunction(prgrm, ident); err == nil {
			// then it's a function
			// not sure about this next line
			lastExprAtomicOp.Outputs = nil
			lastExprAtomicOp.Operator = fn
		} else if strct, err := prgrm.GetStruct(ident, imp.Name); err == nil {
			lastExprAtomicOp.Outputs[0].StructType = strct
		} else {
			// panic(err)
			fmt.Println(err)
			return nil
		}
	} else {
		// then left is not a package name
		if constants2.IsCorePackage(left.Name) {
			println(ast.CompilationError(left.ArgDetails.FileName, left.ArgDetails.FileLine),
				fmt.Sprintf("identifier '%s' does not exist",
					left.Name))
			os.Exit(constants.CX_COMPILATION_ERROR)
		}
		// then it's a struct
		left.IsStruct = true

		fld := ast.MakeArgument(ident, CurrentFile, LineNo)
		leftPkg, err := prgrm.GetPackageFromArray(left.Package)
		if err != nil {
			panic(err)
		}
		fld.AddType(types.IDENTIFIER).AddPackage(leftPkg)

		left.Fields = append(left.Fields, fld)
	}

	return prevExprs
}
