package main

import (
	// "fmt"
	"github.com/skycoin/skycoin/src/cipher/encoder"
	. "github.com/skycoin/cx/cx"
)

// to decide what shorthand op to use
const (
	OP_EQUAL = iota
	OP_UNEQUAL
	
	OP_BITAND
	OP_BITXOR
	OP_BITOR
	OP_BITCLEAR

	OP_MUL
	OP_DIV
	OP_MOD
	OP_ADD
	OP_SUB
	OP_BITSHL
	OP_BITSHR
	OP_LT
	OP_GT
	OP_LTEQ
	OP_GTEQ
)

// used for selection_statement to layout its outputs
type selectStatement struct {
	Condition []*CXExpression
	Then []*CXExpression
	Else []*CXExpression
}

func DeclareGlobal (declarator *CXArgument, declaration_specifiers *CXArgument, initializer []*CXExpression, doesInitialize bool) {
	if doesInitialize {
		if pkg, err := prgrm.GetCurrentPackage(); err == nil {
			if glbl, err := prgrm.GetGlobal(declarator.Name); err != nil {
				expr := WritePrimary(declaration_specifiers.Type, make([]byte, declaration_specifiers.Size))
				exprOut := expr[0].Outputs[0]
				declaration_specifiers.Name = declarator.Name
				declaration_specifiers.MemoryRead = MEM_DATA
				declaration_specifiers.MemoryWrite = MEM_DATA
				declaration_specifiers.Offset = exprOut.Offset
				declaration_specifiers.Lengths = exprOut.Lengths
				declaration_specifiers.Size = exprOut.Size
				declaration_specifiers.TotalSize = exprOut.TotalSize
				declaration_specifiers.Package = exprOut.Package
				pkg.AddGlobal(declaration_specifiers)
			} else {
				if initializer[len(initializer) - 1].Operator == nil {
					expr := MakeExpression(Natives[OP_IDENTITY])
					expr.Package = pkg
					declaration_specifiers.Name = declarator.Name
					declaration_specifiers.MemoryRead = MEM_DATA
					declaration_specifiers.MemoryWrite = MEM_DATA
					declaration_specifiers.Offset = glbl.Offset
					declaration_specifiers.Lengths = glbl.Lengths
					declaration_specifiers.Size = glbl.Size
					declaration_specifiers.TotalSize = glbl.TotalSize
					declaration_specifiers.Package = glbl.Package
					declaration_specifiers.Value = initializer[len(initializer) - 1].Outputs[0].Value

					expr.AddOutput(declaration_specifiers)
					expr.AddInput(initializer[len(initializer) - 1].Outputs[0])

					sysInitExprs = append(sysInitExprs, expr)
				} else {
					declaration_specifiers.Name = declarator.Name
					declaration_specifiers.MemoryRead = MEM_DATA
					declaration_specifiers.MemoryWrite = MEM_DATA
					declaration_specifiers.Offset = glbl.Offset
					declaration_specifiers.Size = glbl.Size
					declaration_specifiers.Lengths = glbl.Lengths
					declaration_specifiers.TotalSize = glbl.TotalSize
					declaration_specifiers.Package = glbl.Package

					expr := initializer[len(initializer) - 1]
					expr.AddOutput(declaration_specifiers)

					sysInitExprs = append(sysInitExprs, initializer...)
				}
			}
		} else {
			panic(err)
		}
	} else {
		if pkg, err := prgrm.GetCurrentPackage(); err == nil {
			if _, err := prgrm.GetGlobal(declarator.Name); err != nil {
				expr := WritePrimary(declaration_specifiers.Type, make([]byte, declaration_specifiers.Size))
				exprOut := expr[0].Outputs[0]
				declaration_specifiers.Name = declarator.Name
				declaration_specifiers.MemoryRead = MEM_DATA
				declaration_specifiers.MemoryWrite = MEM_DATA
				declaration_specifiers.Offset = exprOut.Offset
				declaration_specifiers.Lengths = exprOut.Lengths
				declaration_specifiers.Size = exprOut.Size
				declaration_specifiers.TotalSize = exprOut.TotalSize
				declaration_specifiers.Package = exprOut.Package
				pkg.AddGlobal(declaration_specifiers)
			}
		} else {
			panic(err)
		}
	}
}

func DeclareStruct (ident string, strctFlds []*CXArgument) {
	if pkg, err := prgrm.GetCurrentPackage(); err == nil {
		if _, err := prgrm.GetStruct(ident, pkg.Name); err != nil {
			strct := MakeStruct(ident)
			pkg.AddStruct(strct)

			var size int
			for _, fld := range strctFlds {
				strct.AddField(fld)
				size += fld.TotalSize
			}
			strct.Size = size
		}
	} else {
		panic(err)
	}
}

func DeclarePackage (ident string) {
	if pkg, err := prgrm.GetPackage(ident); err != nil {
		pkg := MakePackage(ident)
		// pkg.AddImport(pkg)
		prgrm.AddPackage(pkg)
	} else {
		prgrm.SelectPackage(pkg.Name)
	}
}

func DeclareImport (ident string) {
	if pkg, err := prgrm.GetCurrentPackage(); err == nil {
		if _, err := pkg.GetImport(ident); err != nil {
			if imp, err := prgrm.GetPackage(ident); err == nil {
				pkg.AddImport(imp)
			} else {
				panic(err)
			}
		}
	} else {
		panic(err)
	}
}

func FunctionHeader (ident string, receiver []*CXArgument, isMethod bool) *CXFunction {
	if isMethod {
		if len(receiver) > 1 {
			panic("method has multiple receivers")
		}
		if pkg, err := prgrm.GetCurrentPackage(); err == nil {
			if fn, err := prgrm.GetFunction(ident, pkg.Name); err == nil {
				fn.AddInput(receiver[0])
				return fn
			} else {
				fn := MakeFunction(ident)
				pkg.AddFunction(fn)
				fn.AddInput(receiver[0])
				return fn
			}
		} else {
			panic(err)
		}
	} else {
		if pkg, err := prgrm.GetCurrentPackage(); err == nil {
			if fn, err := prgrm.GetFunction(ident, pkg.Name); err == nil {
				return fn
			} else {
				fn := MakeFunction(ident)
				pkg.AddFunction(fn)
				return fn
			}
		} else {
			panic(err)
		}
	}
}

// const (

// )

func DeclarationSpecifiers (declSpec *CXArgument, arraySize int, opTyp int) *CXArgument {
	switch opTyp {
	case DECL_POINTER:
		declSpec.DeclarationSpecifiers = append(declSpec.DeclarationSpecifiers, DECL_POINTER)
		if !declSpec.IsPointer {
			declSpec.IsPointer = true
			declSpec.PointeeSize = declSpec.Size
			declSpec.Size = TYPE_POINTER_SIZE
			declSpec.TotalSize = TYPE_POINTER_SIZE
			declSpec.IndirectionLevels++
		} else {
			pointer := declSpec

			for c := declSpec.IndirectionLevels - 1; c > 0 ; c-- {
				pointer = pointer.Pointee
				pointer.IndirectionLevels = c
				pointer.IsPointer = true
			}

			pointee := MakeArgument("")
			pointee.AddType(TypeNames[pointer.Type])
			pointee.IsPointer = true

			declSpec.IndirectionLevels++

			pointer.Size = TYPE_POINTER_SIZE
			pointer.TotalSize = TYPE_POINTER_SIZE
			pointer.Pointee = pointee
		}
		
		return declSpec
	case DECL_ARRAY:
		declSpec.DeclarationSpecifiers = append(declSpec.DeclarationSpecifiers, DECL_ARRAY)
		arg := declSpec
		arg.IsArray = true
		arg.Lengths = append([]int{arraySize}, arg.Lengths...)
		arg.TotalSize = arg.Size * TotalLength(arg.Lengths)

		// byts := make([]byte, arg.TotalSize)
		byts := MakeMultiDimArray(arg.Size, arg.Lengths)
		arg.Value = &byts
		
		return arg
	case DECL_SLICE:
		declSpec.DeclarationSpecifiers = append(declSpec.DeclarationSpecifiers, DECL_SLICE)
		arg := declSpec
		arg.IsArray = true
		arg.Lengths = append([]int{SLICE_SIZE}, arg.Lengths...)
		arg.TotalSize = arg.Size * TotalLength(arg.Lengths)
		
		// byts := make([]byte, arg.TotalSize)
		byts := MakeMultiDimArray(arg.Size, arg.Lengths)
		arg.Value = &byts
		
		return arg
	case DECL_BASIC:
		arg := declSpec
		arg.DeclarationSpecifiers = append(arg.DeclarationSpecifiers, DECL_BASIC)
		arg.TotalSize = arg.Size
		return arg
	}

	return nil
}

func DeclarationSpecifiersBasic (typ int) *CXArgument {
	arg := MakeArgument("")
	arg.AddType(TypeNames[typ])
	arg.Type = typ
	// arg.Typ = "ident"
	arg.Size = GetArgSize(typ)
	return DeclarationSpecifiers(arg, 0, DECL_BASIC)
}

func DeclarationSpecifiersStruct (ident string, pkgName string, isExternal bool) *CXArgument {
	if isExternal {
		// custom type in an imported package
		if pkg, err := prgrm.GetCurrentPackage(); err == nil {
			if imp, err := pkg.GetImport(pkgName); err == nil {
				if strct, err := prgrm.GetStruct(ident, imp.Name); err == nil {
					arg := MakeArgument("")
					// arg.AddType(TypeNames[TYPE_CUSTOM])
					// I'm not sure about the next line
					// cCX doesn't need TYPE_CUSTOM?
					arg.AddType(ident)
					arg.CustomType = strct
					arg.Size = strct.Size
					arg.TotalSize = strct.Size
					arg.Package = pkg
					arg.DeclarationSpecifiers = append(arg.DeclarationSpecifiers, DECL_STRUCT)

					return arg
				} else {
					panic("type '" + ident + "' does not exist")
				}
			} else {
				panic(err)
			}
		} else {
			panic(err)
		}
	} else {
		// custom type in the current package
		if pkg, err := prgrm.GetCurrentPackage(); err == nil {
			if strct, err := prgrm.GetStruct(ident, pkg.Name); err == nil {
				arg := MakeArgument("")
				// arg.AddType(TypeNames[TYPE_CUSTOM])
				// I'm not sure about the next line
				// cCX doesn't need TYPE_CUSTOM?
				arg.AddType(ident)
				arg.DeclarationSpecifiers = append(arg.DeclarationSpecifiers, DECL_STRUCT)
				arg.CustomType = strct
				arg.Size = strct.Size
				arg.TotalSize = strct.Size
				arg.Package = pkg

				return arg
			} else {
				panic("type '" + ident + "' does not exist")
			}
		} else {
			panic(err)
		}
	}
	
}

func StructLiteralFields (ident string) *CXExpression {
	if pkg, err := prgrm.GetCurrentPackage(); err == nil {
		arg := MakeArgument("")
		arg.AddType(TypeNames[TYPE_IDENTIFIER])
		arg.Name = ident
		arg.Package = pkg

		expr := &CXExpression{Outputs: []*CXArgument{arg}}
		expr.Package = pkg

		return expr
	} else {
		panic(err)
	}
}

func ArrayLiteralExpression (arrSize int, typSpec int, exprs []*CXExpression) []*CXExpression {
	var result []*CXExpression

	pkg, err := prgrm.GetCurrentPackage()
	if err != nil {
		panic(err)
	}

	symName := MakeGenSym(LOCAL_PREFIX)

	var endPointsCounter int
	for _, expr := range exprs {
		if expr.IsArrayLiteral {
			expr.IsArrayLiteral = false
			
			sym := MakeArgument(symName).AddType(TypeNames[typSpec])
			sym.Package = pkg

			idxExpr := WritePrimary(TYPE_I32, encoder.Serialize(int32(endPointsCounter)))
			endPointsCounter++

			sym.Indexes = append(sym.Indexes, idxExpr[0].Outputs[0])
			sym.DereferenceOperations = append(sym.DereferenceOperations, DEREF_ARRAY)

			symExpr := MakeExpression(nil)
			symExpr.Outputs = append(symExpr.Outputs, sym)

			// result = append(result, Assignment([]*CXExpression{symExpr}, []*CXExpression{expr})...)
			if expr.Operator == nil {
				// then it's a literal
				symExpr.Operator = Natives[OP_IDENTITY]
				// expr.Outputs[0].Size = symExpr.Outputs[0].Size
				// expr.Outputs[0].Lengths = symExpr.Outputs[0].Lengths

				symExpr.Inputs = expr.Outputs
				// symExpr.Outputs = expr.
			} else {
				symExpr.Operator = expr.Operator
				symExpr.Inputs = expr.Inputs

				// hack to get the correct lengths below
				expr.Outputs = append(expr.Outputs, sym)
			}
			
			// result = append(result, expr)
			result = append(result, symExpr)
			
			// sym.Lengths = append(sym.Lengths, int($2))
			sym.Lengths = append(expr.Outputs[0].Lengths, arrSize)
			sym.TotalSize = sym.Size * TotalLength(sym.Lengths)
		} else {
			result = append(result, expr)
		}
	}
	
	symNameOutput := MakeGenSym(LOCAL_PREFIX)
	
	symOutput := MakeArgument(symNameOutput).AddType(TypeNames[typSpec])
	symOutput.Lengths = append(symOutput.Lengths, arrSize)
	symOutput.Package = pkg
	symOutput.TotalSize = symOutput.Size * TotalLength(symOutput.Lengths)

	symInput := MakeArgument(symName).AddType(TypeNames[typSpec])
	symInput.Lengths = append(symInput.Lengths, arrSize)
	symInput.Package = pkg
	symInput.TotalSize = symInput.Size * TotalLength(symInput.Lengths)
	
	symExpr := MakeExpression(Natives[OP_IDENTITY])
	symExpr.Package = pkg
	symExpr.Outputs = append(symExpr.Outputs, symOutput)
	symExpr.Inputs = append(symExpr.Inputs, symInput)

	// marking the output so multidimensional arrays identify the expressions
	symExpr.IsArrayLiteral = true
	result = append(result, symExpr)

	return result
}

func PrimaryIdentifier (ident string) []*CXExpression {
	if pkg, err := prgrm.GetCurrentPackage(); err == nil {
		arg := MakeArgument(ident)
		arg.AddType(TypeNames[TYPE_IDENTIFIER])
		// arg.Typ = "ident"
		arg.Name = ident
		arg.Package = pkg

		expr := &CXExpression{Outputs: []*CXArgument{arg}}
		expr.Package = pkg

		return []*CXExpression{expr}
	} else {
		panic(err)
	}
}

func PrimaryStructLiteral (ident string, strctFlds []*CXExpression) []*CXExpression {
	var result []*CXExpression
	if pkg, err := prgrm.GetCurrentPackage(); err == nil {
		if strct, err := prgrm.GetStruct(ident, pkg.Name); err == nil {
			for _, expr := range strctFlds {
				fld := MakeArgument("")
				fld.AddType(TypeNames[TYPE_IDENTIFIER])
				fld.Name = expr.Outputs[0].Name

				expr.IsStructLiteral = true

				expr.Outputs[0].Package = pkg
				expr.Outputs[0].Program = prgrm

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

	return result
}

func PrimaryStructLiteralExternal (impName string, ident string, strctFlds []*CXExpression) []*CXExpression {
	var result []*CXExpression
	if pkg, err := prgrm.GetCurrentPackage(); err == nil {
		if _, err := pkg.GetImport(impName); err == nil {
			if strct, err := prgrm.GetStruct(ident, impName); err == nil {
				for _, expr := range strctFlds {
					fld := MakeArgument("")
					fld.AddType(TypeNames[TYPE_IDENTIFIER])
					fld.Name = expr.Outputs[0].Name

					expr.IsStructLiteral = true

					expr.Outputs[0].Package = pkg
					expr.Outputs[0].Program = prgrm

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

func PostfixExpressionArray (prevExprs []*CXExpression, postExprs []*CXExpression) []*CXExpression {
	prevExprs[len(prevExprs) - 1].Outputs[0].IsArray = false
	pastOps := prevExprs[len(prevExprs) - 1].Outputs[0].DereferenceOperations
	if len(pastOps) < 1 || pastOps[len(pastOps) - 1] != DEREF_ARRAY {
		// this way we avoid calling deref_array multiple times (one for each index)
		prevExprs[len(prevExprs) - 1].Outputs[0].DereferenceOperations = append(prevExprs[len(prevExprs) - 1].Outputs[0].DereferenceOperations, DEREF_ARRAY)
	}

	if !prevExprs[len(prevExprs) - 1].Outputs[0].IsDereferenceFirst {
		prevExprs[len(prevExprs) - 1].Outputs[0].IsArrayFirst = true
	}

	if len(prevExprs[len(prevExprs) - 1].Outputs[0].Fields) > 0 {
		fld := prevExprs[len(prevExprs) - 1].Outputs[0].Fields[len(prevExprs[len(prevExprs) - 1].Outputs[0].Fields) - 1]
		fld.Indexes = append(fld.Indexes, postExprs[len(postExprs) - 1].Outputs[0])
	} else {
		if len(postExprs[len(postExprs) - 1].Outputs) < 1 {
			// then it's an expression (e.g. i32.add(0, 0))
			// we create a gensym for it
			idxSym := MakeArgument(MakeGenSym(LOCAL_PREFIX)).AddType(TypeNames[postExprs[len(postExprs) - 1].Operator.Outputs[0].Type])
			idxSym.Size = postExprs[len(postExprs) - 1].Operator.Outputs[0].Size
			idxSym.TotalSize = postExprs[len(postExprs) - 1].Operator.Outputs[0].Size

			idxSym.Package = postExprs[len(postExprs) - 1].Package
			postExprs[len(postExprs) - 1].Outputs = append(postExprs[len(postExprs) - 1].Outputs, idxSym)

			prevExprs[len(prevExprs) - 1].Outputs[0].Indexes = append(prevExprs[len(prevExprs) - 1].Outputs[0].Indexes, idxSym)

			// we push the index expression
			prevExprs = append(postExprs, prevExprs...)
		} else {
			prevExprs[len(prevExprs) - 1].Outputs[0].Indexes = append(prevExprs[len(prevExprs) - 1].Outputs[0].Indexes, postExprs[len(postExprs) - 1].Outputs[0])
		}
	}
	
	expr := prevExprs[len(prevExprs) - 1]
	if len(expr.Inputs) < 1 {
		expr.Inputs = append(expr.Inputs, prevExprs[len(prevExprs) - 1].Outputs[0])
	}

	expr.Inputs = append(expr.Inputs, postExprs[len(postExprs) - 1].Outputs[0])

	return prevExprs
}

func PostfixExpressionNative (typCode int, opCode string) []*CXExpression {
	// these will always be native functions
	if opCode, ok := OpCodes[TypeNames[typCode] + "." + opCode]; ok {
		expr := MakeExpression(Natives[opCode])
		if pkg, err := prgrm.GetCurrentPackage(); err == nil {
			expr.Package = pkg
		} else {
			panic(err)
		}
		
		return []*CXExpression{expr}
	} else {
		panic(ok)
	}
}

func PostfixExpressionEmptyFunCall (prevExprs []*CXExpression) []*CXExpression {
	if prevExprs[len(prevExprs) - 1].Operator == nil {
		if opCode, ok := OpCodes[prevExprs[len(prevExprs) - 1].Outputs[0].Name]; ok {
			if pkg, err := prgrm.GetCurrentPackage(); err == nil {
				prevExprs[0].Package = pkg
			}
			prevExprs[0].Outputs = nil
			prevExprs[0].Operator = Natives[opCode]
		}
	}
	
	prevExprs[0].Inputs = nil
	return FunctionCall(prevExprs, nil)
}

func PostfixExpressionFunCall (prevExprs []*CXExpression, args []*CXExpression) []*CXExpression {
	if prevExprs[len(prevExprs) - 1].Operator == nil {
		if opCode, ok := OpCodes[prevExprs[len(prevExprs) - 1].Outputs[0].Name]; ok {
			if pkg, err := prgrm.GetCurrentPackage(); err == nil {
				prevExprs[0].Package = pkg
			}
			prevExprs[0].Outputs = nil
			prevExprs[0].Operator = Natives[opCode]
		}
	}

	prevExprs[0].Inputs = nil
	
	return FunctionCall(prevExprs, args)
}

func PostfixExpressionIncDec (prevExprs []*CXExpression, isInc bool) []*CXExpression {
	pkg, err := prgrm.GetCurrentPackage()
	if err != nil {
		panic(err)
	}

	var expr *CXExpression
	if isInc {
		expr = MakeExpression(Natives[OP_I32_ADD])
	} else {
		expr = MakeExpression(Natives[OP_I32_SUB])
	}
	
	val := WritePrimary(TYPE_I32, encoder.SerializeAtomic(int32(1)))

	expr.AddInput(prevExprs[len(prevExprs) - 1].Outputs[0])
	expr.AddInput(val[len(val) - 1].Outputs[0])
	expr.AddOutput(prevExprs[len(prevExprs) - 1].Outputs[0])

	expr.Package = pkg
	
	exprs := append(prevExprs, expr)
	return exprs
}

func PostfixExpressionField (prevExprs []*CXExpression, ident string) {
	left := prevExprs[len(prevExprs) - 1].Outputs[0]

	if left.IsRest {
		// then it can't be a module name
		// and we propagate the property to the right expression
		// right.IsRest = true
	} else {
		left.IsRest = true
		// then left is a first (e.g first.rest) and right is a rest
		// let's check if left is a package
		if pkg, err := prgrm.GetCurrentPackage(); err == nil {
			if imp, err := pkg.GetImport(left.Name); err == nil {
				// the external property will be propagated to the following arguments
				// this way we avoid considering these arguments as module names
				left.Package = imp

				if glbl, err := imp.GetGlobal(ident); err == nil {
					// then it's a global
					prevExprs[len(prevExprs) - 1].Outputs[0] = glbl
				} else if fn, err := prgrm.GetFunction(ident, imp.Name); err == nil {
					// then it's a function
					// not sure about this next line
					prevExprs[len(prevExprs) - 1].Outputs = nil
					prevExprs[len(prevExprs) - 1].Operator = fn
				} else if strct, err := prgrm.GetStruct(ident, imp.Name); err == nil {
					prevExprs[len(prevExprs) - 1].Outputs[0].CustomType = strct
				} else {
					panic(err)
				}
			} else {
				if code, ok := ConstCodes[prevExprs[len(prevExprs) - 1].Outputs[0].Name + "." + ident]; ok {
					constant := Constants[code]
					val := WritePrimary(constant.Type, constant.Value)
					prevExprs[len(prevExprs) - 1].Outputs[0] = val[0].Outputs[0]
				} else if _, ok := OpCodes[prevExprs[len(prevExprs) - 1].Outputs[0].Name + "." + ident]; ok {
					// then it's a native
					// TODO: we'd be referring to the function itself, not a function call
					// (functions as first-class objects)
					prevExprs[len(prevExprs) - 1].Outputs[0].Name = prevExprs[len(prevExprs) - 1].Outputs[0].Name + "." + ident
				} else {
					// then it's a struct
					left.IsStruct = true
					left.DereferenceOperations = append(left.DereferenceOperations, DEREF_FIELD)
					fld := MakeArgument(ident)
					fld.AddType(TypeNames[TYPE_IDENTIFIER])
					left.Fields = append(left.Fields, fld)
				}
			}
		} else {
			panic(err)
		}
	}
}

func UnaryExpression (op string, prevExprs []*CXExpression) []*CXExpression {
	exprOut := prevExprs[len(prevExprs) - 1].Outputs[0]
	switch op {
	case "*":
		exprOut.DereferenceLevels++
		exprOut.DereferenceOperations = append(exprOut.DereferenceOperations, DEREF_POINTER)
		if !exprOut.IsArrayFirst {
			exprOut.IsDereferenceFirst = true
		}

		exprOut.IsReference = false
	case "&":
		prevExprs[0].Outputs[0].IsReference = true
		prevExprs[0].Outputs[0].MemoryRead = MEM_STACK
		prevExprs[0].Outputs[0].MemoryWrite = MEM_HEAP
	}
	return prevExprs
}

func ShorthandExpression (leftExprs []*CXExpression, rightExprs []*CXExpression, op int) []*CXExpression {
	var operator *CXFunction
	switch op {
	case OP_EQUAL:
		operator = Natives[OP_UND_EQUAL]
	case OP_UNEQUAL:
		operator = Natives[OP_UND_UNEQUAL]
	case OP_BITAND:
		operator = Natives[OP_UND_BITAND]
	case OP_BITXOR:
		operator = Natives[OP_UND_BITXOR]
	case OP_BITOR:
		operator = Natives[OP_UND_BITOR]
	case OP_MUL:
		operator = Natives[OP_UND_MUL]
	case OP_DIV:
		operator = Natives[OP_UND_DIV]
	case OP_MOD:
		operator = Natives[OP_UND_MOD]
	case OP_ADD:
		operator = Natives[OP_UND_ADD]
	case OP_SUB:
		operator = Natives[OP_UND_SUB]
	case OP_BITSHL:
		operator = Natives[OP_UND_BITSHL]
	case OP_BITSHR:
		operator = Natives[OP_UND_BITSHR]
	case OP_BITCLEAR:
		operator = Natives[OP_UND_BITCLEAR]
	case OP_LT:
		operator = Natives[OP_UND_LT]
	case OP_GT:
		operator = Natives[OP_UND_GT]
	case OP_LTEQ:
		operator = Natives[OP_UND_LTEQ]
	case OP_GTEQ:
		operator = Natives[OP_UND_GTEQ]
	}
	
	return ArithmeticOperation(leftExprs, rightExprs, operator)
}

func DeclareLocal (declarator *CXArgument, declaration_specifiers *CXArgument, initializer []*CXExpression, doesInitialize bool) []*CXExpression {
	if doesInitialize {
		declaration_specifiers.IsLocalDeclaration = true
		
		if pkg, err := prgrm.GetCurrentPackage(); err == nil {
			if initializer[len(initializer) - 1].Operator == nil {
				// then it's a literal, e.g. var foo i32 = 10;
				expr := MakeExpression(Natives[OP_IDENTITY])
				expr.Package = pkg
				
				declaration_specifiers.Name = declarator.Name
				declaration_specifiers.Package = pkg
				// declaration_specifiers.Typ = "ident"
				
				expr.AddOutput(declaration_specifiers)
				expr.AddInput(initializer[len(initializer) - 1].Outputs[0])
				
				return []*CXExpression{expr}
			} else {
				// then it's an expression (it has an operator)
				declaration_specifiers.Name = declarator.Name
				declaration_specifiers.Package = pkg
				
				
				expr := initializer[len(initializer) - 1]
				expr.AddOutput(declaration_specifiers)
				
				// exprs := $5
				// exprs = append(exprs, expr)
				
				return initializer
			}
		} else {
			panic(err)
		}
	} else {
		declaration_specifiers.IsLocalDeclaration = true
		
		// this will tell the runtime that it's just a declaration
		if pkg, err := prgrm.GetCurrentPackage(); err == nil {
			expr := MakeExpression(nil)
			expr.Package = pkg

			declaration_specifiers.Name = declarator.Name
			// declaration_specifiers.Typ = "ident"
			declaration_specifiers.Package = pkg
			expr.AddOutput(declaration_specifiers)

			return []*CXExpression{expr}
		} else {
			panic(err)
		}
	}
}

const (
	SEL_ELSEIF = iota
	SEL_ELSEIFELSE
)

func SelectionStatement (predExprs []*CXExpression, thenExprs []*CXExpression, elseifExprs []selectStatement, elseExprs []*CXExpression, op int) []*CXExpression {
	switch op {
	case SEL_ELSEIFELSE:
		var lastElse []*CXExpression = elseExprs
		for c := len(elseifExprs) - 1; c >= 0; c-- {
			if lastElse != nil {
				lastElse = SelectionExpressions(elseifExprs[c].Condition, elseifExprs[c].Then, lastElse)
			} else {
				lastElse = SelectionExpressions(elseifExprs[c].Condition, elseifExprs[c].Then, nil)
			}
		}

		return SelectionExpressions(predExprs, thenExprs, lastElse)
	case SEL_ELSEIF:
		var lastElse []*CXExpression
		for c := len(elseifExprs) - 1; c >= 0; c-- {
			if lastElse != nil {
				lastElse = SelectionExpressions(elseifExprs[c].Condition, elseifExprs[c].Then, lastElse)
			} else {
				lastElse = SelectionExpressions(elseifExprs[c].Condition, elseifExprs[c].Then, nil)
			}
		}

		return SelectionExpressions(predExprs, thenExprs, lastElse)
	}

	panic("")
}

func ArithmeticOperation (leftExprs []*CXExpression, rightExprs []*CXExpression, operator *CXFunction) (out []*CXExpression) {
	pkg, err := prgrm.GetCurrentPackage()
	if err != nil {
		panic(err)
	}
	
	if len(leftExprs[len(leftExprs) - 1].Outputs) < 1 {
		name := MakeArgument(MakeGenSym(LOCAL_PREFIX)).AddType(TypeNames[leftExprs[len(leftExprs) - 1].Operator.Outputs[0].Type])

		name.Size = leftExprs[len(leftExprs) - 1].Operator.Outputs[0].Size
		name.TotalSize = leftExprs[len(leftExprs) - 1].Operator.Outputs[0].Size
		name.Package = pkg

		leftExprs[len(leftExprs) - 1].Outputs = append(leftExprs[len(leftExprs) - 1].Outputs, name)
	}

	if len(rightExprs[len(rightExprs) - 1].Outputs) < 1 {
		name := MakeArgument(MakeGenSym(LOCAL_PREFIX)).AddType(TypeNames[rightExprs[len(rightExprs) - 1].Operator.Outputs[0].Type])

		name.Size = rightExprs[len(rightExprs) - 1].Operator.Outputs[0].Size
		name.TotalSize = rightExprs[len(rightExprs) - 1].Operator.Outputs[0].Size
		name.Package = pkg

		rightExprs[len(rightExprs) - 1].Outputs = append(rightExprs[len(rightExprs) - 1].Outputs, name)
	}

	// var leftNestedExprs []*CXExpression
	// for _, inpExpr := range leftExprs {
		
	// }

	expr := MakeExpression(operator)
	expr.Package = pkg
	
	if leftExprs[len(leftExprs) - 1].Operator == nil {
		// then it's a literal
		expr.Inputs = append(expr.Inputs, leftExprs[len(leftExprs) - 1].Outputs[0])
	} else {
		// then it's a function call
		out = append(out, leftExprs...)
	}

	if rightExprs[len(rightExprs) - 1].Operator == nil {
		// then it's a literal
		expr.Inputs = append(expr.Inputs, rightExprs[len(rightExprs) - 1].Outputs[0])
	} else {
		// then it's a function call
		out = append(out, rightExprs...)
	}

	// out = append(out, expr)


	// var left *CXArgument
	// // var right *CXArgument
	
	// left = leftExprs[len(leftExprs) - 1].Outputs[0]
	// right = rightExprs[len(rightExprs) - 1].Outputs[0]
	
	// expr.Inputs = append(expr.Inputs, left)
	// expr.Inputs = append(expr.Inputs, right)


	
	// outName := MakeArgument(MakeGenSym(LOCAL_PREFIX)).AddType(TypeNames[left.Type])
	// // outName.Size = GetArgSize(left.Type)
	// // outName.TotalSize = GetArgSize(left.Type)
	// outName.Size = operator.Outputs[0].TotalSize
	// outName.TotalSize = operator.Outputs[0].TotalSize
	
	// outName.Package = pkg

	// expr.Outputs = append(expr.Outputs, outName)

	// out = append(leftExprs, rightExprs...)
	out = append(out, expr)

	return
}

// Primary expressions (literals) are saved in the MEM_DATA segment at compile-time
// This function writes those bytes to prgrm.Data
func WritePrimary (typ int, byts []byte) []*CXExpression {
	if pkg, err := prgrm.GetCurrentPackage(); err == nil {
		arg := MakeArgument("")
		arg.AddType(TypeNames[typ])
		arg.AddValue(&byts)
		arg.MemoryRead = MEM_DATA
		arg.MemoryWrite = MEM_DATA
		arg.Offset = dataOffset
		arg.Package = pkg
		arg.Program = prgrm
		size := len(byts)
		arg.Size = size
		arg.TotalSize = size
		arg.PointeeSize = size
		dataOffset += size
		prgrm.Data = append(prgrm.Data, Data(byts)...)
		expr := MakeExpression(nil)
		expr.Package = pkg
		expr.Outputs = append(expr.Outputs, arg)
		return []*CXExpression{expr}
	} else {
		panic(err)
	}
}

func TotalLength (lengths []int) int {
	var total int = 1
	for _, i := range lengths {
		total *= i
	}
	return total
}

func IterationExpressions (init []*CXExpression, cond []*CXExpression, incr []*CXExpression, statements []*CXExpression) []*CXExpression {
	jmpFn := Natives[OP_JMP]

	pkg, err := prgrm.GetCurrentPackage()
	if err != nil {
		panic(err)
	}
	
	upExpr := MakeExpression(jmpFn)
	upExpr.Package = pkg
	
	trueArg := WritePrimary(TYPE_BOOL, encoder.Serialize(true))

	upLines := (len(statements) + len(incr) + len(cond) + 2) * -1
	downLines := 0
	
	upExpr.AddInput(trueArg[0].Outputs[0])
	upExpr.ThenLines = upLines
	upExpr.ElseLines = downLines
	
	downExpr := MakeExpression(jmpFn)
	downExpr.Package = pkg

	if len(cond[len(cond) - 1].Outputs) < 1 {
		predicate := MakeArgument(MakeGenSym(LOCAL_PREFIX)).AddType(TypeNames[cond[len(cond) - 1].Operator.Outputs[0].Type])
		predicate.Package = pkg
		cond[len(cond) - 1].AddOutput(predicate)
		downExpr.AddInput(predicate)
	} else {
		predicate := cond[len(cond) - 1].Outputs[0]
		predicate.Package = pkg
		downExpr.AddInput(predicate)
	}

	thenLines := 0
	elseLines := len(incr) + len(statements) + 1
	
	downExpr.ThenLines = thenLines
	downExpr.ElseLines = elseLines
	
	exprs := init
	exprs = append(exprs, cond...)
	exprs = append(exprs, downExpr)
	exprs = append(exprs, statements...)
	exprs = append(exprs, incr...)
	exprs = append(exprs, upExpr)
	
	return exprs
}

func StructLiteralAssignment (to []*CXExpression, from []*CXExpression) []*CXExpression {
	pkg, err := prgrm.GetCurrentPackage()
	if err != nil {
		panic(err)
	}

	out := MakeArgument(MakeGenSym(LOCAL_PREFIX)).AddType(TypeNames[TYPE_CUSTOM])
	
	out.CustomType = from[len(from) - 1].Outputs[0].CustomType
	out.Size = from[len(from) - 1].Outputs[0].CustomType.Size
	out.TotalSize = from[len(from) - 1].Outputs[0].CustomType.Size
	out.Package = pkg
	out.Program = prgrm
	
	// before
	for _, f := range from {
		f.Outputs[0].Name = to[0].Outputs[0].Name

		if len(to[0].Outputs[0].Indexes) > 0 {
			f.Outputs[0].Lengths = to[0].Outputs[0].Lengths
			f.Outputs[0].Indexes = to[0].Outputs[0].Indexes
			f.Outputs[0].DereferenceOperations = append(f.Outputs[0].DereferenceOperations, DEREF_ARRAY)
		}

		f.Outputs[0].DereferenceOperations = append(f.Outputs[0].DereferenceOperations, DEREF_FIELD)
	}
	
	return from
}

func ArrayLiteralAssignment (to []*CXExpression, from []*CXExpression) []*CXExpression {
	for _, f := range from {
		f.Outputs[0].Name = to[0].Outputs[0].Name
		f.Outputs[0].DereferenceOperations = append(f.Outputs[0].DereferenceOperations, DEREF_ARRAY)
	}
	
	return from
}

func Assignment (to []*CXExpression, from []*CXExpression) []*CXExpression {
	idx := len(from) - 1

	if from[idx].IsArrayLiteral {
		from[0].Outputs[0].SynonymousTo = to[0].Outputs[0].Name
	}

	if glbl, err := to[0].Outputs[0].Package.GetGlobal(to[0].Outputs[0].Name); err == nil {
		for _, expr := range from {
			expr.Outputs[0].MemoryRead = glbl.MemoryRead
			expr.Outputs[0].MemoryWrite = glbl.MemoryWrite
		}
	}

	if from[idx].Operator == nil {
		from[idx].Operator = Natives[OP_IDENTITY]
		to[0].Outputs[0].Size = from[idx].Outputs[0].Size
		to[0].Outputs[0].Lengths = from[idx].Outputs[0].Lengths
		to[0].Outputs[0].Program = prgrm
		
		// // assigning .Value to field if present
		// if len(to[0].Outputs[0].Fields) > 0 {
		// 	to[0].Outputs[0].Fields[len(to[0].Outputs[0].Fields) - 1].Value = from[idx].Outputs[0].Value
		// } else {
		// 	to[0].Outputs[0].Value = from[idx].Outputs[0].Value
		// }
		
		from[idx].Inputs = from[idx].Outputs
		from[idx].Outputs = to[len(to) - 1].Outputs
		from[idx].Program = prgrm

		return append(to[:len(to) - 1], from...)
	} else {
		if from[idx].Operator.IsNative {
			// only assigning as if the operator had only one output defined
			to[0].Outputs[0].Size = Natives[from[idx].Operator.OpCode].Outputs[0].Size
			to[0].Outputs[0].Lengths = from[idx].Operator.Outputs[0].Lengths
			to[0].Outputs[0].Program = prgrm
		} else {
			// we'll delegate multiple-value returns to the 'expression' grammar rule
			// only assigning as if the operator had only one output defined
			to[0].Outputs[0].Size = from[idx].Operator.Outputs[0].Size
			to[0].Outputs[0].Lengths = from[idx].Operator.Outputs[0].Lengths
			to[0].Outputs[0].Program = prgrm
		}

		// // assigning .Value to field if present
		// if len(to[0].Outputs[0].Fields) > 0 {
		// 	to[0].Outputs[0].Fields[len(to[0].Outputs[0].Fields) - 1].Value = from[idx].Outputs[0].Value
		// } else {
		// 	to[0].Outputs[0].Value = from[idx].Outputs[0].Value
		// }
		
		from[idx].Outputs = to[0].Outputs
		from[idx].Program = to[0].Program

		if from[0].IsStructLiteral {
			from[idx].Outputs[0].MemoryRead = MEM_HEAP
		}

		return append(to[:len(to) - 1], from...)
		// return append(to, from...)
	}
}

func SelectionExpressions (condExprs []*CXExpression, thenExprs []*CXExpression, elseExprs []*CXExpression) []*CXExpression {
	jmpFn := Natives[OP_JMP]
	pkg, err := prgrm.GetCurrentPackage()
	if err != nil {
		panic(err)
	}
	ifExpr := MakeExpression(jmpFn)
	ifExpr.Package = pkg

	var predicate *CXArgument
	if condExprs[len(condExprs) - 1].Operator == nil {
		// then it's a literal
		predicate = condExprs[len(condExprs) - 1].Outputs[0]
	} else {
		// then it's an expression
		predicate = MakeArgument(MakeGenSym(LOCAL_PREFIX)).AddType(TypeNames[condExprs[len(condExprs) - 1].Operator.Outputs[0].Type])
		condExprs[len(condExprs) - 1].Outputs = append(condExprs[len(condExprs) - 1].Outputs, predicate)
	}
	predicate.Package = pkg

	ifExpr.AddInput(predicate)

	thenLines := 0
	elseLines := len(thenExprs) + 1

	ifExpr.ThenLines = thenLines
	ifExpr.ElseLines = elseLines

	skipExpr := MakeExpression(jmpFn)
	skipExpr.Package = pkg

	trueArg := WritePrimary(TYPE_BOOL, encoder.Serialize(true))
	skipLines := len(elseExprs)

	skipExpr.AddInput(trueArg[0].Outputs[0])
	skipExpr.ThenLines = skipLines
	skipExpr.ElseLines = 0

	var exprs []*CXExpression
	if condExprs[len(condExprs) - 1].Operator != nil {
		exprs = append(exprs, condExprs...)
	}
	exprs = append(exprs, ifExpr)
	exprs = append(exprs, thenExprs...)
	exprs = append(exprs, skipExpr)
	exprs = append(exprs, elseExprs...)
	
	return exprs
}

func GetSymType (sym *CXArgument, fn *CXFunction) int {
	if sym.Name == "" {
		// then literal
		return sym.Type
	}
	if glbl, err := sym.Package.GetGlobal(sym.Name); err == nil {
		// then it's a global
		return glbl.Type
	} else {
		// then it's a local
		for _, inp := range fn.Inputs {
			if inp.Name == sym.Name {
				return inp.Type
			}
		}
		for _, out := range fn.Outputs {
			if out.Name == sym.Name {
				return out.Type
			}
		}

		for _, expr := range fn.Expressions {
			for _, inp := range expr.Inputs {
				if inp.Name == sym.Name {
					return inp.Type
				}
			}
			for _, out := range expr.Outputs {
				if out.Name == sym.Name {
					return out.Type
				}
			}
		}
	}
	return TYPE_UNDEFINED
}

func SetFinalSize (symbols *map[string]*CXArgument, sym *CXArgument) {
	var elt *CXArgument
	var finalSize int = sym.TotalSize
	var fldIdx int
	elt = sym

	if arg, found := (*symbols)[sym.Package.Name + "." + sym.Name]; found {
		for _, op := range sym.DereferenceOperations {
			switch op {
			case DEREF_ARRAY:
				for i, _ := range elt.Indexes {
					var subSize int = 1
					for _, len := range elt.Lengths[i:] {
						subSize *= len
					}
					finalSize /= subSize
				}
			case DEREF_FIELD:
				elt = sym.Fields[fldIdx]
				finalSize = elt.TotalSize
				fldIdx++
			case DEREF_POINTER:
				if len(arg.DeclarationSpecifiers) > 0 {
					var subSize int
					subSize = 1
					for _, decl := range arg.DeclarationSpecifiers {
						switch decl {
						case DECL_ARRAY:
							for _, len := range arg.Lengths {
								subSize *= len
							}
						case DECL_BASIC:
							subSize = GetArgSize(elt.Type)
						case DECL_STRUCT:
							subSize = arg.CustomType.Size
						}
					}

					// finalSize = derefSize
					finalSize = subSize
				}
			}
		}
	}

	sym.TotalSize = finalSize
}

func GetGlobalSymbol (symbols *map[string]*CXArgument, symPackage *CXPackage, symName string) {
	if _, found := (*symbols)[symPackage.Name + "." + symName]; !found {
		if glbl, err := symPackage.GetGlobal(symName); err == nil {
			(*symbols)[symPackage.Name + "." + symName] = glbl
		}
	}
}

func GiveOffset (symbols *map[string]*CXArgument, sym *CXArgument, offset *int, shouldExist bool) {
	if sym.Name != "" {
		GetGlobalSymbol(symbols, sym.Package, sym.Name)

		if arg, found := (*symbols)[sym.Package.Name + "." + sym.Name]; !found {
			if shouldExist {
				// it should exist. error
				panic("identifier '" + sym.Name + "' does not exist")
			}

			if sym.SynonymousTo != "" {
				// then the offset needs to be shared
				GetGlobalSymbol(symbols, sym.Package, sym.SynonymousTo)
				sym.Offset = (*symbols)[sym.Package.Name + "." + sym.SynonymousTo].Offset
				
				(*symbols)[sym.Package.Name + "." + sym.Name] = sym
			} else {
				sym.Offset = *offset
				(*symbols)[sym.Package.Name + "." + sym.Name] = sym

				*offset += sym.TotalSize

				if sym.IsPointer {
					pointer := sym
					for c := 0; c < sym.IndirectionLevels - 1; c++ {
						pointer = pointer.Pointee
						pointer.Offset = *offset
						*offset += pointer.TotalSize
					}
				}
			}
		} else {
			if sym.IsReference {
				if arg.HeapOffset < 1 {
					// then it hasn't been assigned
					// an offset of 0 is impossible because the symbol was declared before
					arg.HeapOffset = *offset
					// sym.HeapOffset = *offset
					*offset += TYPE_POINTER_SIZE
				}
				
				// if not, then it has been assigned before
				// and we just reassign it to this symbol
				// we'll do this below, where we're assigning everything to sym
			}

			var isFieldPointer bool
			if len(sym.Fields) > 0 {
				var found bool

				strct := arg.CustomType
				for _, nameFld := range sym.Fields {
					for _, fld := range strct.Fields {
						if nameFld.Name == fld.Name {
							if fld.IsPointer {
								sym.IsPointer = true
								// sym.IndirectionLevels = fld.IndirectionLevels
								isFieldPointer = true
							}
							found = true
							if fld.CustomType != nil {
								strct = fld.CustomType
							}
							break
						}
					}
					if !found {
						panic("field '" + nameFld.Name + "' not found")
					}
				}
			}
			
			if sym.DereferenceLevels > 0 {
				if arg.IndirectionLevels >= sym.DereferenceLevels || isFieldPointer { // ||
					// 	sym.IndirectionLevels >= sym.DereferenceLevels
					// {
					pointer := arg

					for c := 0; c < sym.DereferenceLevels - 1; c++ {
						pointer = pointer.Pointee
					}

					sym.Offset = pointer.Offset
					sym.IndirectionLevels = pointer.IndirectionLevels
					sym.IsPointer = pointer.IsPointer
				} else {
					panic("invalid indirect of " + sym.Name)
				}
			} else {
				sym.Offset = arg.Offset
				sym.IsPointer = arg.IsPointer
				sym.IndirectionLevels = arg.IndirectionLevels
			}

			//if sym.IsStruct {
			// checking if it's accessing fields
			if len(sym.Fields) > 0 {
				var found bool

				strct := arg.CustomType
				for _, nameFld := range sym.Fields {
					for _, fld := range strct.Fields {
						if nameFld.Name == fld.Name {
							nameFld.Lengths = fld.Lengths
							nameFld.Size = fld.Size
							nameFld.TotalSize = fld.TotalSize
							nameFld.DereferenceLevels = sym.DereferenceLevels
							nameFld.IsPointer = fld.IsPointer
							found = true
							if fld.CustomType != nil {
								strct = fld.CustomType
							}
							break
						}
						
						nameFld.Offset += fld.TotalSize
					}
					if !found {
						panic("field '" + nameFld.Name + "' not found")
					}
				}
			}

			// sym.IsPointer = arg.IsPointer
			// sym.Typ = arg.Typ

			sym.Type = arg.Type
			sym.CustomType = arg.CustomType
			sym.Pointee = arg.Pointee
			sym.Lengths = arg.Lengths
			sym.PointeeSize = arg.PointeeSize
			sym.Package = arg.Package
			sym.Program = arg.Program
			sym.HeapOffset = arg.HeapOffset

			if !sym.IsReference {
				sym.MemoryRead = arg.MemoryRead
				sym.MemoryWrite = arg.MemoryWrite
			}

			if sym.DereferenceLevels > 0 {
				sym.MemoryRead = MEM_HEAP
				sym.MemoryWrite = MEM_HEAP
			}
			
			if sym.IsReference && !arg.IsStruct {
				// sym.Size = TYPE_POINTER_SIZE
				// sym.TotalSize = TYPE_POINTER_SIZE
				sym.TotalSize = arg.TotalSize
				
				sym.Size = arg.Size
				// sym.TotalSize = arg.TotalSize
			} else {
				// we need to implement a more robust system, like the one in op.go
				if len(sym.Fields) > 0 {
					// sym.Size = sym.Fields[len(sym.Fields) - 1].Size
					sym.Size = arg.Size
					sym.TotalSize = sym.Fields[len(sym.Fields) - 1].TotalSize
				} else {
					sym.Size = arg.Size
					sym.TotalSize = arg.TotalSize
				}
			}
		}
	}
}

func FunctionDeclaration (fn *CXFunction, inputs []*CXArgument, outputs []*CXArgument, exprs []*CXExpression) {
	// adding inputs, outputs
	for _, inp := range inputs {
		fn.AddInput(inp)
	}
	for _, out := range outputs {
		fn.AddOutput(out)
	}

	// getting offset to use by statements (excluding inputs, outputs and receiver)
	var offset int

	for _, expr := range exprs {
		fn.AddExpression(expr)
	}

	fn.Length = len(fn.Expressions)

	var symbols map[string]*CXArgument = make(map[string]*CXArgument, 0)
	var symbolsScope map[string]bool = make(map[string]bool, 0)

	for _, inp := range fn.Inputs {
		if inp.IsLocalDeclaration {
			symbolsScope[inp.Package.Name + "." + inp.Name] = true
		}
		inp.IsLocalDeclaration = symbolsScope[inp.Package.Name + "." + inp.Name]

		GiveOffset(&symbols, inp, &offset, false)
		SetFinalSize(&symbols, inp)

		AddPointer(fn, inp)
	}
	for _, out := range fn.Outputs {
		if out.IsLocalDeclaration {
			symbolsScope[out.Package.Name + "." + out.Name] = true
		}
		out.IsLocalDeclaration = symbolsScope[out.Package.Name + "." + out.Name]
		
		GiveOffset(&symbols, out, &offset, false)
		SetFinalSize(&symbols, out)
		
		AddPointer(fn, out)
	}

	for _, expr := range fn.Expressions {
		for _, inp := range expr.Inputs {
			// if inp.Name != "" {
			// 	fmt.Println("inpExprBegin", inp.Name, inp.MemoryRead, inp.MemoryWrite)
			// }
			if inp.IsLocalDeclaration {
				symbolsScope[inp.Package.Name + "." + inp.Name] = true
			}
			inp.IsLocalDeclaration = symbolsScope[inp.Package.Name + "." + inp.Name]

			GiveOffset(&symbols, inp, &offset, true)
			SetFinalSize(&symbols, inp)
			
			for _, idx := range inp.Indexes {
				GiveOffset(&symbols, idx, &offset, true)
			}

			// if _, found := symbols[inp.Package.Name + "." + inp.Name]; !found {
			// 	AddPointer(fn, inp)
			// }
			AddPointer(fn, inp)
			// if inp.Name != "" {
			// 	fmt.Println("inpExprEnd", inp.Name, inp.MemoryRead, inp.MemoryWrite)
			// }
		}
		for _, out := range expr.Outputs {
			// if out.Name != "" {
			// 	fmt.Println("outExprBegin", out.Name, out.MemoryRead, out.MemoryWrite)
			// }
			if out.IsLocalDeclaration {
				symbolsScope[out.Package.Name + "." + out.Name] = true
			}
			
			out.IsLocalDeclaration = symbolsScope[
				out.Package.Name + "." +
				out.Name]
			
			GiveOffset(&symbols, out, &offset, false)
			SetFinalSize(&symbols, out)
			for _, idx := range out.Indexes {
				GiveOffset(&symbols, idx, &offset, true)
			}

			// if _, found := symbols[out.Package.Name + "." + out.Name]; !found {
			// 	AddPointer(fn, out)
			// }
			AddPointer(fn, out)
			// if out.Name != "" {
			// 	fmt.Println("outExprEnd", out.MemoryRead, out.MemoryWrite)
			// }
			
		}
		
		SetCorrectArithmeticOp(expr)
	}

	// checking if assigning pointer to pointer
	for _, expr := range fn.Expressions {
		if expr.Operator == Natives[OP_IDENTITY] {
			for i, out := range expr.Outputs {
				if out.IsPointer && expr.Inputs[i].IsPointer {
					// we're modifying the actual pointer
					expr.Inputs[i].MemoryRead = MEM_STACK
					expr.Inputs[i].MemoryWrite = MEM_STACK
					out.MemoryRead = MEM_STACK
					out.MemoryWrite = MEM_STACK
				}
			}
		}
	}
	
	fn.Size = offset
}

func FunctionCall (exprs []*CXExpression, args []*CXExpression) []*CXExpression {
	expr := exprs[len(exprs) - 1]

	if expr.Operator == nil {
		opName := expr.Outputs[0].Name
		opPkg := expr.Outputs[0].Package
		if len(expr.Outputs[0].Fields) > 0 {
			opName = expr.Outputs[0].Fields[0].Name
			// it wasn't a field, but a method call. removing it as a field
			expr.Outputs[0].Fields = expr.Outputs[0].Fields[:len(expr.Outputs[0].Fields) - 1]
			// we remove information about the "field" (method name)
			expr.AddInput(expr.Outputs[0])
			
			expr.Outputs = expr.Outputs[:len(expr.Outputs) - 1]
			// expr.Inputs = expr.Inputs[:len(expr.Inputs) - 1]
			// expr.AddInput(expr.Outputs[0])
		}

		if op, err := prgrm.GetFunction(opName, opPkg.Name); err == nil {
			expr.Operator = op
		} else {
			panic(err)
		}
		
		expr.Outputs = nil
	}

	var nestedExprs []*CXExpression
	for _, inpExpr := range args {
		if inpExpr.Operator == nil {
			// then it's a literal
			expr.AddInput(inpExpr.Outputs[0])
		} else {
			// then it's a function call
			if len(inpExpr.Outputs) < 1 {
				var out *CXArgument
				
				if inpExpr.Operator.Outputs[0].Type == TYPE_UNDEFINED {
					// if undefined type, then adopt argument's type
					out = MakeArgument(MakeGenSym(LOCAL_PREFIX)).AddType(TypeNames[inpExpr.Inputs[0].Type])
					out.Size = inpExpr.Inputs[0].Size
					out.TotalSize = inpExpr.Inputs[0].Size
				} else {
					out = MakeArgument(MakeGenSym(LOCAL_PREFIX)).AddType(TypeNames[inpExpr.Operator.Outputs[0].Type])
					out.Size = inpExpr.Operator.Outputs[0].Size
					out.TotalSize = inpExpr.Operator.Outputs[0].Size
				}
				
				
				out.Typ = "ident"
				
				out.Package = inpExpr.Package
				inpExpr.AddOutput(out)
				expr.AddInput(out)
			}
			nestedExprs = append(nestedExprs, inpExpr)
		}
	}
	
	return append(nestedExprs, exprs...)
}
