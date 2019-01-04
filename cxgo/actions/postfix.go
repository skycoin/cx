package actions

import (
	"fmt"
	. "github.com/skycoin/cx/cx"
	"github.com/skycoin/skycoin/src/cipher/encoder"
	"os"
)

func PostfixExpressionArray(prevExprs []*CXExpression, postExprs []*CXExpression) []*CXExpression {
	var elt *CXArgument
	if len(prevExprs[len(prevExprs)-1].Outputs[0].Fields) > 0 {
		elt = prevExprs[len(prevExprs)-1].Outputs[0].Fields[len(prevExprs[len(prevExprs)-1].Outputs[0].Fields)-1]
	} else {
		elt = prevExprs[len(prevExprs)-1].Outputs[0]
	}

	elt.IsArray = false
	pastOps := elt.DereferenceOperations
	if len(pastOps) < 1 || pastOps[len(pastOps)-1] != DEREF_ARRAY {
		// this way we avoid calling deref_array multiple times (one for each index)
		elt.DereferenceOperations = append(elt.DereferenceOperations, DEREF_ARRAY)
	}

	if !elt.IsDereferenceFirst {
		elt.IsArrayFirst = true
	}

	if len(prevExprs[len(prevExprs)-1].Outputs[0].Fields) > 0 {
		fld := prevExprs[len(prevExprs)-1].Outputs[0].Fields[len(prevExprs[len(prevExprs)-1].Outputs[0].Fields)-1]

		if postExprs[len(postExprs)-1].Operator == nil {
			// expr.AddInput(postExprs[len(postExprs)-1].Outputs[0])
			fld.Indexes = append(fld.Indexes, postExprs[len(postExprs)-1].Outputs[0])
		} else {
			sym := MakeArgument(MakeGenSym(LOCAL_PREFIX), CurrentFile, LineNo).AddType(TypeNames[postExprs[len(postExprs)-1].Inputs[0].Type])
			sym.Package = postExprs[len(postExprs)-1].Package
			sym.PreviouslyDeclared = true
			postExprs[len(postExprs)-1].AddOutput(sym)

			prevExprs = append(postExprs, prevExprs...)

			fld.Indexes = append(fld.Indexes, sym)
			// expr.AddInput(sym)
		}

		// fld.Indexes = append(fld.Indexes, postExprs[len(postExprs)-1].Outputs[0])
	} else {
		if len(postExprs[len(postExprs)-1].Outputs) < 1 {
			// then it's an expression (e.g. i32.add(0, 0))
			// we create a gensym for it
			idxSym := MakeArgument(MakeGenSym(LOCAL_PREFIX), CurrentFile, LineNo).AddType(TypeNames[postExprs[len(postExprs)-1].Operator.Outputs[0].Type])
			idxSym.Size = postExprs[len(postExprs)-1].Operator.Outputs[0].Size
			idxSym.TotalSize = postExprs[len(postExprs)-1].Operator.Outputs[0].Size

			idxSym.Package = postExprs[len(postExprs)-1].Package
			idxSym.PreviouslyDeclared = true
			postExprs[len(postExprs)-1].Outputs = append(postExprs[len(postExprs)-1].Outputs, idxSym)

			prevExprs[len(prevExprs)-1].Outputs[0].Indexes = append(prevExprs[len(prevExprs)-1].Outputs[0].Indexes, idxSym)

			// we push the index expression
			prevExprs = append(postExprs, prevExprs...)
		} else {
			prevExprs[len(prevExprs)-1].Outputs[0].Indexes = append(prevExprs[len(prevExprs)-1].Outputs[0].Indexes, postExprs[len(postExprs)-1].Outputs[0])
		}
	}

	// expr := prevExprs[len(prevExprs)-1]
	// if len(expr.Inputs) < 1 {
	// 	expr.Inputs = append(expr.Inputs, prevExprs[len(prevExprs)-1].Outputs[0])
	// }

	return prevExprs
}

func PostfixExpressionNative(typCode int, opStrCode string) []*CXExpression {
	// these will always be native functions
	if opCode, ok := OpCodes[TypeNames[typCode]+"."+opStrCode]; ok {
		expr := MakeExpression(Natives[opCode], CurrentFile, LineNo)
		if pkg, err := PRGRM.GetCurrentPackage(); err == nil {
			expr.Package = pkg
		} else {
			panic(err)
		}

		return []*CXExpression{expr}
	} else {
		println(CompilationError(CurrentFile, LineNo) + " function '" + TypeNames[typCode] + "." + opStrCode + "' does not exist")
		return nil
		// panic(ok)
	}
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

	val := WritePrimary(TYPE_I32, encoder.SerializeAtomic(int32(1)), false)

	expr.Package = pkg

	expr.AddInput(prevExprs[len(prevExprs)-1].Outputs[0])
	expr.AddInput(val[len(val)-1].Outputs[0])
	expr.AddOutput(prevExprs[len(prevExprs)-1].Outputs[0])

	// exprs := append(prevExprs, expr)
	exprs := append([]*CXExpression{}, expr)
	return exprs
}

func PostfixExpressionField(prevExprs []*CXExpression, ident string) {
	left := prevExprs[len(prevExprs)-1].Outputs[0]

	if left.IsRest {
		// then it can't be a package name
		// and we propagate the property to the right expression
		// right.IsRest = true
		// left.DereferenceOperations = append(left.DereferenceOperations, DEREF_FIELD)
		fld := MakeArgument(ident, CurrentFile, LineNo)
		fld.AddType(TypeNames[TYPE_IDENTIFIER])
		left.Fields = append(left.Fields, fld)
	} else {
		left.IsRest = true
		// then left is a first (e.g first.rest) and right is a rest
		// let's check if left is a package
		if pkg, err := PRGRM.GetCurrentPackage(); err == nil {
			if imp, err := pkg.GetImport(left.Name); err == nil {
				// the external property will be propagated to the following arguments
				// this way we avoid considering these arguments as module names

				if IsCorePackage(left.Name) {
					if code, ok := ConstCodes[left.Name+"."+ident]; ok {
						constant := Constants[code]
						val := WritePrimary(constant.Type, constant.Value, false)
						prevExprs[len(prevExprs)-1].Outputs[0] = val[0].Outputs[0]
						return
					} else if _, ok := OpCodes[left.Name+"."+ident]; ok {
						// then it's a native
						// TODO: we'd be referring to the function itself, not a function call
						// (functions as first-class objects)
						left.Name = left.Name + "." + ident
						return
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
				} else if fn, err := PRGRM.GetFunction(ident, imp.Name); err == nil {
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
					println(CompilationError(left.FileName, left.FileLine), fmt.Sprintf("identifier '%s' does not exist", left.Name))
					os.Exit(CX_COMPILATION_ERROR)
					return
				}

				// then it's a struct
				left.IsStruct = true

				fld := MakeArgument(ident, CurrentFile, LineNo)
				fld.AddType(TypeNames[TYPE_IDENTIFIER])
				left.Fields = append(left.Fields, fld)
			}
		} else {
			panic(err)
		}
	}
}
