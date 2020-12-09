package actions

import (
	"fmt"
	"os"

	. "github.com/skycoin/cx/cx"
)

// PostfixExpressionArray...
//
func PostfixExpressionArray(prevExprs []*CXExpression, postExprs []*CXExpression) []*CXExpression {
	var elt *CXArgument
	prevExpr := prevExprs[len(prevExprs)-1]

	if prevExpr.Operator != nil && len(prevExpr.Outputs) == 0 {
		genName := MakeGenSym(LOCAL_PREFIX)

		out := MakeArgument(genName, prevExpr.FileName, prevExpr.FileLine-1).AddType(TypeNames[prevExpr.Operator.Outputs[0].Type])

		out.DeclarationSpecifiers = prevExpr.Operator.Outputs[0].DeclarationSpecifiers
		out.CustomType = prevExpr.Operator.Outputs[0].CustomType
		out.Size = prevExpr.Operator.Outputs[0].Size
		out.TotalSize = prevExpr.Operator.Outputs[0].TotalSize
		out.Lengths = prevExpr.Operator.Outputs[0].Lengths
		out.IsSlice = prevExpr.Operator.Outputs[0].IsSlice
		out.PreviouslyDeclared = true

		prevExpr.AddOutput(out)

		inp := MakeArgument(genName, prevExpr.FileName, prevExpr.FileLine).AddType(TypeNames[prevExpr.Operator.Outputs[0].Type])

		inp.DeclarationSpecifiers = prevExpr.Operator.Outputs[0].DeclarationSpecifiers
		inp.CustomType = prevExpr.Operator.Outputs[0].CustomType
		inp.Size = prevExpr.Operator.Outputs[0].Size
		inp.TotalSize = prevExpr.Operator.Outputs[0].TotalSize
		inp.Lengths = prevExpr.Operator.Outputs[0].Lengths
		inp.IsSlice = prevExpr.Operator.Outputs[0].IsSlice
		inp.PreviouslyDeclared = true

		useExpr := MakeExpression(nil, prevExpr.FileName, prevExpr.FileLine)
		useExpr.Package = prevExpr.Package
		useExpr.AddOutput(inp)

		prevExprs = append(prevExprs, useExpr)
	}

	prevExpr = prevExprs[len(prevExprs)-1]

	if len(prevExpr.Outputs[0].Fields) > 0 {
		elt = prevExpr.Outputs[0].Fields[len(prevExpr.Outputs[0].Fields)-1]
	} else {
		elt = prevExpr.Outputs[0]
	}

	elt.IsArray = false
	elt.DereferenceOperations = append(elt.DereferenceOperations, DEREF_ARRAY)
	elt.DeclarationSpecifiers = append(elt.DeclarationSpecifiers, DECL_INDEXING)

	if !elt.IsDereferenceFirst {
		elt.IsArrayFirst = true
	}

	if len(prevExprs[len(prevExprs)-1].Outputs[0].Fields) > 0 {
		fld := prevExprs[len(prevExprs)-1].Outputs[0].Fields[len(prevExprs[len(prevExprs)-1].Outputs[0].Fields)-1]

		if postExprs[len(postExprs)-1].Operator == nil {
			// expr.AddInput(postExprs[len(postExprs)-1].Outputs[0])
			fld.Indexes = append(fld.Indexes, postExprs[len(postExprs)-1].Outputs[0])
		} else {
			sym := MakeArgument(MakeGenSym(LOCAL_PREFIX), CurrentFile, LineNo).AddType(TypeNames[postExprs[len(postExprs)-1].Operator.Outputs[0].Type])
			sym.Package = postExprs[len(postExprs)-1].Package
			sym.PreviouslyDeclared = true
			postExprs[len(postExprs)-1].AddOutput(sym)

			prevExprs = append(postExprs, prevExprs...)

			fld.Indexes = append(fld.Indexes, sym)
			// expr.AddInput(sym)
		}
	} else {
		if len(postExprs[len(postExprs)-1].Outputs) < 1 {
			// then it's an expression (e.g. i32.add(0, 0))
			// we create a gensym for it
			idxSym := MakeArgument(MakeGenSym(LOCAL_PREFIX), CurrentFile, LineNo).AddType(TypeNames[postExprs[len(postExprs)-1].Operator.Outputs[0].Type])
			idxSym.Size = postExprs[len(postExprs)-1].Operator.Outputs[0].Size
			idxSym.TotalSize = GetSize(postExprs[len(postExprs)-1].Operator.Outputs[0])

			idxSym.Package = postExprs[len(postExprs)-1].Package
			idxSym.PreviouslyDeclared = true
			postExprs[len(postExprs)-1].Outputs = append(postExprs[len(postExprs)-1].Outputs, idxSym)

			prevExprs[len(prevExprs)-1].Outputs[0].Indexes = append(prevExprs[len(prevExprs)-1].Outputs[0].Indexes, idxSym)

			// we push the index expression
			prevExprs = append(postExprs, prevExprs...)
		} else {
			prevOuts := prevExprs[len(prevExprs)-1].Outputs
			postOuts := postExprs[len(postExprs)-1].Outputs
			prevOuts[0].Indexes = append(prevOuts[0].Indexes, postOuts[0])
		}
	}

	return prevExprs
}

func PostfixExpressionNative(typCode int, opStrCode string) []*CXExpression {
	// these will always be native functions
	opCode, ok := OpCodes[TypeNames[typCode]+"."+opStrCode]
	if !ok {
		println(CompilationError(CurrentFile, LineNo) + " function '" +
			TypeNames[typCode] + "." + opStrCode + "' does not exist")
		return nil
		// panic(ok)
	}

	expr := MakeExpression(Natives[opCode], CurrentFile, LineNo)
	pkg, err := PRGRM.GetCurrentPackage()
	if err != nil {
		panic(err)
	}
	expr.Package = pkg

	return []*CXExpression{expr}
}

func PostfixExpressionEmptyFunCall(prevExprs []*CXExpression) []*CXExpression {
	if prevExprs[len(prevExprs)-1].Outputs != nil && len(prevExprs[len(prevExprs)-1].Outputs[0].Fields) > 0 {
		// then it's a method call or function in field
		// prevExprs[len(prevExprs) - 1].IsMethodCall = true
		// expr.IsMethodCall = true
		// // method name
		// expr.Operator = MakeFunction(expr.Outputs[0].Fields[0].Name)
		// inp := MakeArgument(expr.Outputs[0].Name, CurrentFile, LineNo)
		// inp.Package = expr.Package
		// inp.Type = expr.Outputs[0].Type
		// inp.CustomType = expr.Outputs[0].CustomType
		// expr.Inputs = append(expr.Inputs, inp)

	} else if prevExprs[len(prevExprs)-1].Operator == nil {
		if opCode, ok := OpCodes[prevExprs[len(prevExprs)-1].Outputs[0].Name]; ok {
			if pkg, err := PRGRM.GetCurrentPackage(); err == nil {
				prevExprs[0].Package = pkg
			}
			prevExprs[0].Outputs = nil
			prevExprs[0].Operator = Natives[opCode]
		}

		prevExprs[0].Inputs = nil
	}

	return FunctionCall(prevExprs, nil)
}

func PostfixExpressionFunCall(prevExprs []*CXExpression, args []*CXExpression) []*CXExpression {
	if prevExprs[len(prevExprs)-1].Outputs != nil && len(prevExprs[len(prevExprs)-1].Outputs[0].Fields) > 0 {
		// then it's a method
		// prevExprs[len(prevExprs) - 1].IsMethodCall = true

	} else if prevExprs[len(prevExprs)-1].Operator == nil {
		if opCode, ok := OpCodes[prevExprs[len(prevExprs)-1].Outputs[0].Name]; ok {
			if pkg, err := PRGRM.GetCurrentPackage(); err == nil {
				prevExprs[0].Package = pkg
			}
			prevExprs[0].Outputs = nil
			prevExprs[0].Operator = Natives[opCode]
		}

		prevExprs[0].Inputs = nil
	}

	return FunctionCall(prevExprs, args)
}

func PostfixExpressionIncDec(prevExprs []*CXExpression, isInc bool) []*CXExpression {
	pkg, err := PRGRM.GetCurrentPackage()
	if err != nil {
		panic(err)
	}

	var expr *CXExpression
	if isInc {
		expr = MakeExpression(Natives[OP_I32_ADD], CurrentFile, LineNo)
	} else {
		expr = MakeExpression(Natives[OP_I32_SUB], CurrentFile, LineNo)
	}

	var valB [4]byte
	WriteMemI32(valB[:], 0, int32(1))
	val := WritePrimary(TYPE_I32, valB[:], false)

	expr.Package = pkg

	expr.AddInput(prevExprs[len(prevExprs)-1].Outputs[0])
	expr.AddInput(val[len(val)-1].Outputs[0])
	expr.AddOutput(prevExprs[len(prevExprs)-1].Outputs[0])

	// exprs := append(prevExprs, expr)
	exprs := append([]*CXExpression{}, expr)
	return exprs
}

// PostfixExpressionField handles the dot notation that can follow an identifier.
// Examples are: `foo.bar`, `foo().bar`, `pkg.foo`
func PostfixExpressionField(prevExprs []*CXExpression, ident string) []*CXExpression {
	lastExpr := prevExprs[len(prevExprs)-1]

	// Then it's a function call, e.g. foo().fld
	// and we need to create some auxiliary variables to hold the result from
	// the function call
	if lastExpr.Operator != nil {
		opOut := lastExpr.Operator.Outputs[0]
		symName := MakeGenSym(LOCAL_PREFIX)

		// we associate the result of the function call to the aux variable
		out := MakeArgument(symName, lastExpr.FileName, lastExpr.FileLine).AddType(TypeNames[opOut.Type])
		out.DeclarationSpecifiers = opOut.DeclarationSpecifiers
		out.CustomType = opOut.CustomType
		out.Size = opOut.Size
		out.TotalSize = opOut.TotalSize
		out.IsArray = opOut.IsArray
		out.IsReference = opOut.IsReference
		out.Lengths = opOut.Lengths
		out.Package = lastExpr.Package
		out.PreviouslyDeclared = true
		out.IsRest = true

		lastExpr.Outputs = append(lastExpr.Outputs, out)

		// we need to create an expression to hold all the modifications
		// that will take place after this if statement
		inp := MakeArgument(symName, lastExpr.FileName, lastExpr.FileLine).AddType(TypeNames[opOut.Type])
		inp.DeclarationSpecifiers = opOut.DeclarationSpecifiers
		inp.CustomType = opOut.CustomType
		inp.Size = opOut.Size
		inp.TotalSize = opOut.TotalSize
		inp.Package = lastExpr.Package
		inp.IsRest = true

		expr := MakeExpression(nil, lastExpr.FileName, lastExpr.FileLine)
		expr.Package = lastExpr.Package
		expr.AddOutput(inp)

		prevExprs = append(prevExprs, expr)

		lastExpr = prevExprs[len(prevExprs)-1]
	}

	left := lastExpr.Outputs[0]

	// If the left already is a rest (e.g. "var" in "pkg.var"), then
	// it can't be a package name and we propagate the property to
	//  the right side.
	if left.IsRest {
		// right.IsRest = true
		// left.DereferenceOperations = append(left.DereferenceOperations, DEREF_FIELD)
		left.IsStruct = true
		fld := MakeArgument(ident, CurrentFile, LineNo)
		fld.AddType(TypeNames[TYPE_IDENTIFIER]).AddPackage(left.Package)
		left.Fields = append(left.Fields, fld)
		return prevExprs
	}

	left.IsRest = true
	// then left is a first (e.g first.rest) and right is a rest
	// let's check if left is a package
	pkg, err := PRGRM.GetCurrentPackage()
	if err != nil {
		panic(err)
	}

	if imp, err := pkg.GetImport(left.Name); err == nil {
		// the external property will be propagated to the following arguments
		// this way we avoid considering these arguments as module names

		if IsCorePackage(left.Name) {
			if code, ok := ConstCodes[left.Name+"."+ident]; ok {
				constant := Constants[code]
				val := WritePrimary(constant.Type, constant.Value, false)
				prevExprs[len(prevExprs)-1].Outputs[0] = val[0].Outputs[0]
				return prevExprs
			} else if _, ok := OpCodes[left.Name+"."+ident]; ok {
				// then it's a native
				// TODO: we'd be referring to the function itself, not a function call
				// (functions as first-class objects)
				left.Name = left.Name + "." + ident
				return prevExprs
			}
		}

		left.Package = imp

		if glbl, err := imp.GetGlobal(ident); err == nil {
			// then it's a global
			// prevExprs[len(prevExprs)-1].Outputs[0] = glbl
			prevExprs[len(prevExprs)-1].Outputs[0].Name = glbl.Name
			prevExprs[len(prevExprs)-1].Outputs[0].Type = glbl.Type
			prevExprs[len(prevExprs)-1].Outputs[0].CustomType = glbl.CustomType
			prevExprs[len(prevExprs)-1].Outputs[0].Size = glbl.Size
			prevExprs[len(prevExprs)-1].Outputs[0].TotalSize = glbl.TotalSize
			prevExprs[len(prevExprs)-1].Outputs[0].IsPointer = glbl.IsPointer
			prevExprs[len(prevExprs)-1].Outputs[0].IsSlice = glbl.IsSlice
			prevExprs[len(prevExprs)-1].Outputs[0].IsStruct = glbl.IsStruct
			prevExprs[len(prevExprs)-1].Outputs[0].Package = glbl.Package
		} else if fn, err := imp.GetFunction(ident); err == nil {
			// then it's a function
			// not sure about this next line
			prevExprs[len(prevExprs)-1].Outputs = nil
			prevExprs[len(prevExprs)-1].Operator = fn
		} else if strct, err := PRGRM.GetStruct(ident, imp.Name); err == nil {
			prevExprs[len(prevExprs)-1].Outputs[0].CustomType = strct
		} else {
			panic(err)
		}
	} else {
		// then left is not a package name
		if IsCorePackage(left.Name) {
			println(CompilationError(left.FileName, left.FileLine),
				fmt.Sprintf("identifier '%s' does not exist",
					left.Name))
			os.Exit(CX_COMPILATION_ERROR)
		}
		// then it's a struct
		left.IsStruct = true

		fld := MakeArgument(ident, CurrentFile, LineNo)
		fld.AddType(TypeNames[TYPE_IDENTIFIER]).AddPackage(left.Package)

		left.Fields = append(left.Fields, fld)
	}

	return prevExprs
}
