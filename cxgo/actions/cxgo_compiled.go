package actions

import (
	"fmt"
	"os"
	"strconv"
	. "github.com/skycoin/cx/cx"
	"github.com/skycoin/skycoin/src/cipher/encoder"
)

var PRGRM *CXProgram
var DataOffset int = STACK_SIZE + TYPE_POINTER_SIZE // to be able to handle nil pointers

var CurrentFile string
var LineNo int = 0
var WebMode bool
var BaseOutput bool
var ReplMode bool
var HelpMode bool
var InterpretMode bool
var CompileMode bool
var ReplTargetFn string = ""
var ReplTargetStrct string = ""
var ReplTargetMod string = ""

var FoundCompileErrors bool

var InREPL bool = false

var SysInitExprs []*CXExpression



// var cxt = interpreted.MakeProgram()
//var cxt = cx0.CXT

var dStack bool = false
var InFn bool = false
//var dProgram bool = false
var tag string = ""
var asmNL = "\n"
var fileName string

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
type SelectStatement struct {
	Condition []*CXExpression
	Then      []*CXExpression
	Else      []*CXExpression
}

func ErrorHeader (currentFile string, lineNo int) string {
	return "error: " + currentFile + ":" + strconv.FormatInt(int64(lineNo), 10)
}

func DeclareGlobal(declarator *CXArgument, declaration_specifiers *CXArgument, initializer []*CXExpression, doesInitialize bool) {
	if doesInitialize {
		if pkg, err := PRGRM.GetCurrentPackage(); err == nil {
			if glbl, err := PRGRM.GetGlobal(declarator.Name); err != nil {
				expr := WritePrimary(declaration_specifiers.Type, make([]byte, declaration_specifiers.Size), true)
				exprOut := expr[0].Outputs[0]
				declaration_specifiers.Name = declarator.Name
				declaration_specifiers.Offset = exprOut.Offset
				// declaration_specifiers.Lengths = exprOut.Lengths
				declaration_specifiers.Size = exprOut.Size
				declaration_specifiers.TotalSize = exprOut.TotalSize
				declaration_specifiers.Package = exprOut.Package

				pkg.AddGlobal(declaration_specifiers)
			} else {
				if initializer[len(initializer)-1].Operator == nil {
					expr := MakeExpression(Natives[OP_IDENTITY], CurrentFile, LineNo)
					expr.Package = pkg
					declaration_specifiers.Name = declarator.Name
					declaration_specifiers.Offset = glbl.Offset
					// declaration_specifiers.Lengths = glbl.Lengths
					declaration_specifiers.Size = glbl.Size
					declaration_specifiers.TotalSize = glbl.TotalSize
					declaration_specifiers.Package = glbl.Package

					expr.AddOutput(declaration_specifiers)
					expr.AddInput(initializer[len(initializer)-1].Outputs[0])

					SysInitExprs = append(SysInitExprs, expr)
				} else {
					declaration_specifiers.Name = declarator.Name
					declaration_specifiers.Offset = glbl.Offset
					declaration_specifiers.Size = glbl.Size
					// declaration_specifiers.Lengths = glbl.Lengths
					declaration_specifiers.TotalSize = glbl.TotalSize
					declaration_specifiers.Package = glbl.Package

					expr := initializer[len(initializer)-1]
					expr.AddOutput(declaration_specifiers)

					SysInitExprs = append(SysInitExprs, initializer...)
				}
			}
		} else {
			panic(err)
		}
	} else {
		if pkg, err := PRGRM.GetCurrentPackage(); err == nil {
			if _, err := PRGRM.GetGlobal(declarator.Name); err != nil {
				expr := WritePrimary(declaration_specifiers.Type, make([]byte, declaration_specifiers.Size), true)
				exprOut := expr[0].Outputs[0]
				declaration_specifiers.Name = declarator.Name
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

func DeclareStruct(ident string, strctFlds []*CXArgument) {
	if pkg, err := PRGRM.GetCurrentPackage(); err == nil {
		if _, err := PRGRM.GetStruct(ident, pkg.Name); err != nil {
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

func DeclarePackage(ident string) {
	if pkg, err := PRGRM.GetPackage(ident); err != nil {
		pkg := MakePackage(ident)
		// pkg.AddImport(pkg)
		PRGRM.AddPackage(pkg)
	} else {
		PRGRM.SelectPackage(pkg.Name)
	}
}

func DeclareImport(ident string) {
	if pkg, err := PRGRM.GetCurrentPackage(); err == nil {
		if _, err := pkg.GetImport(ident); err != nil {
			if imp, err := PRGRM.GetPackage(ident); err == nil {
				pkg.AddImport(imp)
			} else {
				// look in the workspace
				
				panic(err)
			}
		}
	} else {
		panic(err)
	}
}

func FunctionHeader(ident string, receiver []*CXArgument, isMethod bool) *CXFunction {
	if isMethod {
		if len(receiver) > 1 {
			panic("method has multiple receivers")
		}
		if pkg, err := PRGRM.GetCurrentPackage(); err == nil {
			// fn := MakeFunction(ident)
			// pkg.AddFunction(fn)
			// fn.AddInput(receiver[0])
			// return fn

			fnName := receiver[0].CustomType.Name + "." + ident

			if fn, err := PRGRM.GetFunction(fnName, pkg.Name); err == nil {
				fn.AddInput(receiver[0])
				return fn
			} else {
				fn := MakeFunction(fnName)
				pkg.AddFunction(fn)
				fn.AddInput(receiver[0])
				return fn
			}
		} else {
			panic(err)
		}
	} else {
		if pkg, err := PRGRM.GetCurrentPackage(); err == nil {
			if fn, err := PRGRM.GetFunction(ident, pkg.Name); err == nil {
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

func DeclarationSpecifiers(declSpec *CXArgument, arraySize int, opTyp int) *CXArgument {
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

			for c := declSpec.IndirectionLevels - 1; c > 0; c-- {
				pointer = pointer.Pointee
				pointer.IndirectionLevels = c
				pointer.IsPointer = true
			}

			pointee := MakeArgument("", CurrentFile, LineNo)
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
		// byts := MakeMultiDimArray(arg.Size, arg.Lengths)

		return arg
	case DECL_SLICE:
		declSpec.DeclarationSpecifiers = append(declSpec.DeclarationSpecifiers, DECL_SLICE)

		arg := declSpec
		arg.IsSlice = true
		arg.IsReference = true
		arg.IsArray = true
		arg.PassBy = PASSBY_REFERENCE

		arg.Lengths = append([]int{0}, arg.Lengths...)
		arg.TotalSize = arg.Size
		arg.Size = TYPE_POINTER_SIZE

		return arg
	case DECL_BASIC:
		arg := declSpec
		arg.DeclarationSpecifiers = append(arg.DeclarationSpecifiers, DECL_BASIC)
		arg.TotalSize = arg.Size
		return arg
	}

	return nil
}

func DeclarationSpecifiersBasic(typ int) *CXArgument {
	arg := MakeArgument("", CurrentFile, LineNo)
	arg.AddType(TypeNames[typ])
	arg.Type = typ

	// arg.Typ = "ident"
	arg.Size = GetArgSize(typ)

	if typ == TYPE_AFF {
		// equivalent to slice of strings
		return DeclarationSpecifiers(arg, 0, DECL_SLICE)
	}

	// if typ == TYPE_STR {
	// 	return DeclarationSpecifiers(arg, 0, DECL_POINTER)
	// }
	
	return DeclarationSpecifiers(arg, 0, DECL_BASIC)
}

func DeclarationSpecifiersStruct(ident string, pkgName string, isExternal bool) *CXArgument {
	if isExternal {
		// custom type in an imported package
		if pkg, err := PRGRM.GetCurrentPackage(); err == nil {
			if imp, err := pkg.GetImport(pkgName); err == nil {
				if strct, err := PRGRM.GetStruct(ident, imp.Name); err == nil {
					arg := MakeArgument("", CurrentFile, LineNo)
					arg.AddType(ident)
					arg.CustomType = strct
					arg.Size = strct.Size
					arg.TotalSize = strct.Size

					arg.Package = pkg
					arg.DeclarationSpecifiers = append(arg.DeclarationSpecifiers, DECL_STRUCT)

					return arg
				} else {
					println(ErrorHeader(CurrentFile, LineNo), err.Error())
					os.Exit(3)
					return nil
					// panic("type '" + ident + "' does not exist")
				}
			} else {
				panic(err)
			}
		} else {
			panic(err)
		}
	} else {
		// custom type in the current package
		if pkg, err := PRGRM.GetCurrentPackage(); err == nil {
			if strct, err := PRGRM.GetStruct(ident, pkg.Name); err == nil {
				arg := MakeArgument("", CurrentFile, LineNo)
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

func StructLiteralFields(ident string) *CXExpression {
	if pkg, err := PRGRM.GetCurrentPackage(); err == nil {
		arg := MakeArgument("", CurrentFile, LineNo)
		arg.AddType(TypeNames[TYPE_IDENTIFIER])
		arg.Name = ident
		arg.Package = pkg

		// expr := &CXExpression{Outputs: []*CXArgument{arg}}
		expr := MakeExpression(nil, CurrentFile, LineNo)
		expr.Outputs = []*CXArgument{arg}
		expr.Package = pkg

		return expr
	} else {
		panic(err)
	}
}

func ArrayLiteralExpression(arrSize int, typSpec int, exprs []*CXExpression) []*CXExpression {
	var result []*CXExpression

	pkg, err := PRGRM.GetCurrentPackage()
	if err != nil {
		panic(err)
	}

	symName := MakeGenSym(LOCAL_PREFIX)

	var endPointsCounter int
	for _, expr := range exprs {
		if expr.IsArrayLiteral {
			expr.IsArrayLiteral = false

			sym := MakeArgument(symName, CurrentFile, LineNo).AddType(TypeNames[typSpec])
			sym.Package = pkg
			sym.IsShortDeclaration = true

			if sym.Type == TYPE_STR {
				sym.PassBy = PASSBY_REFERENCE
			}

			idxExpr := WritePrimary(TYPE_I32, encoder.Serialize(int32(endPointsCounter)), false)
			endPointsCounter++

			sym.Indexes = append(sym.Indexes, idxExpr[0].Outputs[0])
			sym.DereferenceOperations = append(sym.DereferenceOperations, DEREF_ARRAY)

			symExpr := MakeExpression(nil, CurrentFile, LineNo)
			symExpr.Outputs = append(symExpr.Outputs, sym)

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

	symOutput := MakeArgument(symNameOutput, CurrentFile, LineNo).AddType(TypeNames[typSpec])
	symOutput.Lengths = append(symOutput.Lengths, arrSize)
	symOutput.Package = pkg
	symOutput.IsShortDeclaration = true
	symOutput.TotalSize = symOutput.Size * TotalLength(symOutput.Lengths)

	// symOutput.DeclarationSpecifiers = append(symOutput.DeclarationSpecifiers, DECL_ARRAY)

	symInput := MakeArgument(symName, CurrentFile, LineNo).AddType(TypeNames[typSpec])
	symInput.Lengths = append(symInput.Lengths, arrSize)
	symInput.Package = pkg
	symInput.IsShortDeclaration = true
	symInput.TotalSize = symInput.Size * TotalLength(symInput.Lengths)

	symExpr := MakeExpression(Natives[OP_IDENTITY], CurrentFile, LineNo)
	symExpr.Package = pkg
	symExpr.Outputs = append(symExpr.Outputs, symOutput)
	symExpr.Inputs = append(symExpr.Inputs, symInput)

	symOutput.SynonymousTo = symInput.Name

	// marking the output so multidimensional arrays identify the expressions
	symExpr.IsArrayLiteral = true
	result = append(result, symExpr)

	return result
}

func SliceLiteralExpression(typSpec int, exprs []*CXExpression) []*CXExpression {
	var result []*CXExpression

	pkg, err := PRGRM.GetCurrentPackage()
	if err != nil {
		panic(err)
	}

	symName := MakeGenSym(LOCAL_PREFIX)

	// adding the declaration
	slcVarExpr := MakeExpression(nil, CurrentFile, LineNo)
	slcVarExpr.Package = pkg
	slcVar := MakeArgument(symName, CurrentFile, LineNo)
	slcVar = DeclarationSpecifiers(slcVar, 0, DECL_SLICE)
	slcVar.AddType(TypeNames[typSpec])
	// slcVarExpr.AddOutput(slcVar)

	// slcVar.TotalSize = slcVar.Size
	// slcVar.Size = TYPE_POINTER_SIZE
	slcVar.TotalSize = TYPE_POINTER_SIZE
	

	slcVarExpr.Outputs = append(slcVarExpr.Outputs, slcVar)
	slcVar.Package = pkg
	slcVar.IsShortDeclaration = true

	result = append(result, slcVarExpr)

	var endPointsCounter int
	for _, expr := range exprs {
		if expr.IsArrayLiteral {
			// expr.IsArrayLiteral = false

			symInp := MakeArgument(symName, CurrentFile, LineNo).AddType(TypeNames[typSpec])
			symInp.Package = pkg
			symOut := MakeArgument(symName, CurrentFile, LineNo).AddType(TypeNames[typSpec])
			symOut.Package = pkg

			// sym = DeclarationSpecifiers(sym, 0, DECL_SLICE)

			// if symInp.Type == TYPE_STR {
			// 	symInp.PassBy = PASSBY_REFERENCE
			// 	symOut.PassBy = PASSBY_REFERENCE
			// }

			endPointsCounter++

			// sym.DereferenceOperations = append(sym.DereferenceOperations, DEREF_POINTER)
			// sym.IndirectionLevels++

			symExpr := MakeExpression(nil, CurrentFile, LineNo)
			symExpr.Package = pkg
			// symExpr.Outputs = append(symExpr.Outputs, symOut)
			symExpr.AddOutput(symOut)

			if expr.Operator == nil {
				// then it's a literal
				symExpr.Operator = Natives[OP_APPEND ]

				symExpr.Inputs = nil
				symExpr.Inputs = append(symExpr.Inputs, symInp)
				symExpr.Inputs = append(symExpr.Inputs, expr.Outputs...)
			} else {
				symExpr.Operator = expr.Operator

				symExpr.Inputs = nil
				symExpr.Inputs = append(symExpr.Inputs, symInp)
				symExpr.Inputs = append(symExpr.Inputs, expr.Inputs...)

				// hack to get the correct lengths below
				expr.Outputs = append(expr.Outputs, symInp)
			}

			// result = append(result, expr)
			result = append(result, symExpr)

			// symInp.TotalSize = symInp.Size
			// symInp.Size = TYPE_POINTER_SIZE
			symInp.TotalSize = TYPE_POINTER_SIZE

			// symOut.TotalSize = symOut.Size
			// symOut.Size = TYPE_POINTER_SIZE
			symOut.TotalSize = TYPE_POINTER_SIZE
		} else {
			result = append(result, expr)
		}
	}

	symNameOutput := MakeGenSym(LOCAL_PREFIX)

	symOutput := MakeArgument(symNameOutput, CurrentFile, LineNo).AddType(TypeNames[typSpec])
	symOutput.PassBy = PASSBY_REFERENCE
	// symOutput.IsSlice = true
	symOutput.Package = pkg
	symOutput.IsShortDeclaration = true

	// symOutput.DeclarationSpecifiers = append(symOutput.DeclarationSpecifiers, DECL_ARRAY)

	symInput := MakeArgument(symName, CurrentFile, LineNo).AddType(TypeNames[typSpec])
	// symInput.IsSlice = true
	symInput.Package = pkg
	symInput.PassBy = PASSBY_REFERENCE

	// symInput.TotalSize = symInput.Size
	// symInput.Size = TYPE_POINTER_SIZE
	symInput.TotalSize = TYPE_POINTER_SIZE

	// symOutput.TotalSize = symOutput.Size
	// symOutput.Size = TYPE_POINTER_SIZE
	symOutput.TotalSize = TYPE_POINTER_SIZE

	symExpr := MakeExpression(Natives[OP_IDENTITY], CurrentFile, LineNo)
	symExpr.Package = pkg
	symExpr.Outputs = append(symExpr.Outputs, symOutput)
	symExpr.Inputs = append(symExpr.Inputs, symInput)

	symOutput.SynonymousTo = symInput.Name

	// marking the output so multidimensional arrays identify the expressions
	// symExpr.IsArrayLiteral = true
	result = append(result, symExpr)

	return result
}

func PrimaryIdentifier(ident string) []*CXExpression {
	if pkg, err := PRGRM.GetCurrentPackage(); err == nil {
		arg := MakeArgument(ident, CurrentFile, LineNo)
		arg.AddType(TypeNames[TYPE_IDENTIFIER])
		// arg.Typ = "ident"
		arg.Name = ident
		arg.Package = pkg

		// expr := &CXExpression{Outputs: []*CXArgument{arg}}
		expr := MakeExpression(nil, CurrentFile, LineNo)
		expr.Outputs = []*CXArgument{arg}
		expr.Package = pkg

		return []*CXExpression{expr}
	} else {
		panic(err)
	}
}

func PrimaryStructLiteral(ident string, strctFlds []*CXExpression) []*CXExpression {
	var result []*CXExpression
	if pkg, err := PRGRM.GetCurrentPackage(); err == nil {
		if strct, err := PRGRM.GetStruct(ident, pkg.Name); err == nil {
			for _, expr := range strctFlds {
				name := expr.Outputs[0].Name

				fld := MakeArgument(name, CurrentFile, LineNo)
				fld.Type = expr.Outputs[0].Type

				expr.IsStructLiteral = true

				expr.Outputs[0].Package = pkg
				expr.Outputs[0].Program = PRGRM

				if expr.Outputs[0].CustomType == nil {
					expr.Outputs[0].CustomType = strct
				}
				// expr.Outputs[0].CustomType = strct
				fld.CustomType = strct

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

func PrimaryStructLiteralExternal(impName string, ident string, strctFlds []*CXExpression) []*CXExpression {
	var result []*CXExpression
	if pkg, err := PRGRM.GetCurrentPackage(); err == nil {
		if _, err := pkg.GetImport(impName); err == nil {
			if strct, err := PRGRM.GetStruct(ident, impName); err == nil {
				for _, expr := range strctFlds {
					fld := MakeArgument("", CurrentFile, LineNo)
					fld.AddType(TypeNames[TYPE_IDENTIFIER])
					fld.Name = expr.Outputs[0].Name

					expr.IsStructLiteral = true

					expr.Outputs[0].Package = pkg
					expr.Outputs[0].Program = PRGRM

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

func PostfixExpressionArray(prevExprs []*CXExpression, postExprs []*CXExpression) []*CXExpression {
	prevExprs[len(prevExprs)-1].Outputs[0].IsArray = false
	pastOps := prevExprs[len(prevExprs)-1].Outputs[0].DereferenceOperations
	if len(pastOps) < 1 || pastOps[len(pastOps)-1] != DEREF_ARRAY {
		// this way we avoid calling deref_array multiple times (one for each index)
		prevExprs[len(prevExprs)-1].Outputs[0].DereferenceOperations = append(prevExprs[len(prevExprs)-1].Outputs[0].DereferenceOperations, DEREF_ARRAY)
	}

	if !prevExprs[len(prevExprs)-1].Outputs[0].IsDereferenceFirst {
		prevExprs[len(prevExprs)-1].Outputs[0].IsArrayFirst = true
	}

	if len(prevExprs[len(prevExprs)-1].Outputs[0].Fields) > 0 {
		fld := prevExprs[len(prevExprs)-1].Outputs[0].Fields[len(prevExprs[len(prevExprs)-1].Outputs[0].Fields)-1]
		fld.Indexes = append(fld.Indexes, postExprs[len(postExprs)-1].Outputs[0])
	} else {
		if len(postExprs[len(postExprs)-1].Outputs) < 1 {
			// then it's an expression (e.g. i32.add(0, 0))
			// we create a gensym for it
			idxSym := MakeArgument(MakeGenSym(LOCAL_PREFIX), CurrentFile, LineNo).AddType(TypeNames[postExprs[len(postExprs)-1].Operator.Outputs[0].Type])
			idxSym.Size = postExprs[len(postExprs)-1].Operator.Outputs[0].Size
			idxSym.TotalSize = postExprs[len(postExprs)-1].Operator.Outputs[0].Size

			idxSym.Package = postExprs[len(postExprs)-1].Package
			idxSym.IsShortDeclaration = true
			postExprs[len(postExprs)-1].Outputs = append(postExprs[len(postExprs)-1].Outputs, idxSym)

			prevExprs[len(prevExprs)-1].Outputs[0].Indexes = append(prevExprs[len(prevExprs)-1].Outputs[0].Indexes, idxSym)

			// we push the index expression
			prevExprs = append(postExprs, prevExprs...)
		} else {
			prevExprs[len(prevExprs)-1].Outputs[0].Indexes = append(prevExprs[len(prevExprs)-1].Outputs[0].Indexes, postExprs[len(postExprs)-1].Outputs[0])
		}
	}

	expr := prevExprs[len(prevExprs)-1]
	if len(expr.Inputs) < 1 {
		expr.Inputs = append(expr.Inputs, prevExprs[len(prevExprs)-1].Outputs[0])
	}

	expr.Inputs = append(expr.Inputs, postExprs[len(postExprs)-1].Outputs[0])

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
		println(ErrorHeader(CurrentFile, LineNo) + " function '" + TypeNames[typCode]+"."+opStrCode + "' does not exist")
		os.Exit(3)
		return nil
		// panic(ok)
	}
}

func PostfixExpressionEmptyFunCall(prevExprs []*CXExpression) []*CXExpression {
	if prevExprs[len(prevExprs) - 1].Outputs != nil && len(prevExprs[len(prevExprs) - 1].Outputs[0].Fields) > 0 {
		// then it's a method call or function in field
		expr := prevExprs[len(prevExprs) - 1]
		// opName := expr.Outputs[0].Fields[0].CustomType.Name + "." +
		// 	expr.Outputs[0].Fields[0].Name
		expr.IsMethodCall = true
		// method name
		expr.Operator = MakeFunction(expr.Outputs[0].Fields[0].Name)
		inp := MakeArgument(expr.Outputs[0].Name, CurrentFile, LineNo)
		inp.Package = expr.Package
		inp.Type = expr.Outputs[0].Type
		inp.CustomType = expr.Outputs[0].CustomType
		expr.Inputs = append(expr.Inputs, inp)
		// expr.Outputs[0].DereferenceOperations = expr.Outputs[0].DereferenceOperations[:len(expr.Outputs[0].DereferenceOperations)-1]
		// expr.Outputs = nil
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
	
	
	// if len(prevExprs[len(prevExprs) - 1].Outputs[0].Fields) > 0 {
	// 	// then it's a method call or function in field

	// 	expr := prevExprs[len(prevExprs) - 1]

	// 	expr.Operator = nil
		
	// 	opName := expr.Outputs[0].Fields[0].CustomType.Name + "." + expr.Outputs[0].Fields[0].Name
	// 	opPkg := expr.Outputs[0].Package
 
	// 	opName := expr.Outputs[0]
	// 	if op, err := PRGRM.GetFunction(opName, opPkg.Name); err == nil {
	// 		expr.Operator = op
	// 	} else {
	// 		panic(err)
	// 	}

	// 	expr.Outputs = nil
		
	// 	// // and we remove the last field, as it's not a field, but a method name
	// 	// expr.Outputs[0].Fields = expr.Outputs[0].Fields[:len(expr.Outputs[0].Fields)-1]
	// 	// // and also remove the field dereference operation
	// 	// expr.Outputs[0].DereferenceOperations = expr.Outputs[0].DereferenceOperations[:len(expr.Outputs[0].DereferenceOperations)-1]
	// } else if prevExprs[len(prevExprs)-1].Operator == nil {
	// 	if opCode, ok := OpCodes[prevExprs[len(prevExprs)-1].Outputs[0].Name]; ok {
	// 		if pkg, err := PRGRM.GetCurrentPackage(); err == nil {
	// 			prevExprs[0].Package = pkg
	// 		}
	// 		prevExprs[0].Outputs = nil
	// 		prevExprs[0].Operator = Natives[opCode]
	// 	}

	// 	prevExprs[0].Inputs = nil
	// }


	// return FunctionCall(prevExprs, nil)











	// if prevExprs[len(prevExprs)-1].Operator == nil {
	// 	if opCode, ok := OpCodes[prevExprs[len(prevExprs)-1].Outputs[0].Name]; ok {
	// 		if pkg, err := PRGRM.GetCurrentPackage(); err == nil {
	// 			prevExprs[0].Package = pkg
	// 		}
	// 		prevExprs[0].Outputs = nil
	// 		prevExprs[0].Operator = Natives[opCode]
	// 	}
	// }

	// prevExprs[0].Inputs = nil
	// return FunctionCall(prevExprs, nil)
}

func PostfixExpressionFunCall(prevExprs []*CXExpression, args []*CXExpression) []*CXExpression {
	if prevExprs[len(prevExprs)-1].Operator == nil {
		if opCode, ok := OpCodes[prevExprs[len(prevExprs)-1].Outputs[0].Name]; ok {
			if pkg, err := PRGRM.GetCurrentPackage(); err == nil {
				prevExprs[0].Package = pkg
			}
			prevExprs[0].Outputs = nil
			prevExprs[0].Operator = Natives[opCode]
		}
	}

	prevExprs[0].Inputs = nil

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

	exprs := append(prevExprs, expr)
	return exprs
}

func PostfixExpressionField(prevExprs []*CXExpression, ident string) {
	left := prevExprs[len(prevExprs)-1].Outputs[0]

	if left.IsRest {
		// then it can't be a package name
		// and we propagate the property to the right expression
		// right.IsRest = true
		left.DereferenceOperations = append(left.DereferenceOperations, DEREF_FIELD)
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
				left.Package = imp

				if glbl, err := imp.GetGlobal(ident); err == nil {
					// then it's a global
					prevExprs[len(prevExprs)-1].Outputs[0] = glbl
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
				if code, ok := ConstCodes[prevExprs[len(prevExprs)-1].Outputs[0].Name+"."+ident]; ok {
					constant := Constants[code]
					val := WritePrimary(constant.Type, constant.Value, false)
					prevExprs[len(prevExprs)-1].Outputs[0] = val[0].Outputs[0]
				} else if _, ok := OpCodes[prevExprs[len(prevExprs)-1].Outputs[0].Name+"."+ident]; ok {
					// then it's a native
					// TODO: we'd be referring to the function itself, not a function call
					// (functions as first-class objects)
					prevExprs[len(prevExprs)-1].Outputs[0].Name = prevExprs[len(prevExprs)-1].Outputs[0].Name + "." + ident
				} else {
					// then it's a struct
					left.IsStruct = true
					left.DereferenceOperations = append(left.DereferenceOperations, DEREF_FIELD)
					fld := MakeArgument(ident, CurrentFile, LineNo)
					fld.AddType(TypeNames[TYPE_IDENTIFIER])
					left.Fields = append(left.Fields, fld)
					
					
				}
			}
		} else {
			panic(err)
		}
	}
}

func UnaryExpression(op string, prevExprs []*CXExpression) []*CXExpression {
	exprOut := prevExprs[len(prevExprs)-1].Outputs[0]
	// exprInp := prevExprs[len(prevExprs)-1].Inputs[0]
	switch op {
	case "*":
		exprOut.DereferenceLevels++
		exprOut.DereferenceOperations = append(exprOut.DereferenceOperations, DEREF_POINTER)
		if !exprOut.IsArrayFirst {
			exprOut.IsDereferenceFirst = true
		}

		// exprOut.Outputs[0].MemoryWrite =
		// exprOut.PassBy = PASSBY_REFERENCE
		
		exprOut.IsReference = false
	case "&":
		// exprOut.IsReference = true
		// exprOut.DoesEscape = true
		exprOut.PassBy = PASSBY_REFERENCE
	case "!":
		if pkg, err := PRGRM.GetCurrentPackage(); err == nil {
			expr := MakeExpression(Natives[OP_BOOL_NOT], CurrentFile, LineNo)
			expr.Package = pkg

			expr.AddInput(prevExprs[len(prevExprs)-1].Outputs[0])

			prevExprs[len(prevExprs)-1] = expr
		} else {
			panic(err)
		}
	}
	return prevExprs
}

func ShorthandExpression(leftExprs []*CXExpression, rightExprs []*CXExpression, op int) []*CXExpression {
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

func DeclareLocal(declarator *CXArgument, declaration_specifiers *CXArgument, initializer []*CXExpression, doesInitialize bool) []*CXExpression {
	if doesInitialize {
		declaration_specifiers.IsLocalDeclaration = true

		if pkg, err := PRGRM.GetCurrentPackage(); err == nil {
			if initializer[len(initializer)-1].Operator == nil {
				// then it's a literal, e.g. var foo i32 = 10;
				expr := MakeExpression(Natives[OP_IDENTITY], CurrentFile, LineNo)
				expr.Package = pkg

				declaration_specifiers.Name = declarator.Name
				declaration_specifiers.Package = pkg
				declaration_specifiers.IsShortDeclaration = true
				// declaration_specifiers.Typ = "ident"

				expr.AddOutput(declaration_specifiers)
				expr.AddInput(initializer[len(initializer)-1].Outputs[0])

				return []*CXExpression{expr}
			} else {
				// then it's an expression (it has an operator)
				declaration_specifiers.Name = declarator.Name
				declaration_specifiers.Package = pkg
				declaration_specifiers.IsShortDeclaration = true

				expr := initializer[len(initializer)-1]
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
		if pkg, err := PRGRM.GetCurrentPackage(); err == nil {
			expr := MakeExpression(nil, CurrentFile, LineNo)
			expr.Package = pkg

			declaration_specifiers.Name = declarator.Name
			// declaration_specifiers.Typ = "ident"
			declaration_specifiers.Package = pkg
			declaration_specifiers.IsShortDeclaration = true
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

func SelectionStatement(predExprs []*CXExpression, thenExprs []*CXExpression, elseifExprs []SelectStatement, elseExprs []*CXExpression, op int) []*CXExpression {
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

func ArithmeticOperation(leftExprs []*CXExpression, rightExprs []*CXExpression, operator *CXFunction) (out []*CXExpression) {
	pkg, err := PRGRM.GetCurrentPackage()
	if err != nil {
		panic(err)
	}

	if len(leftExprs[len(leftExprs)-1].Outputs) < 1 {
		// name := MakeArgument(MakeGenSym(LOCAL_PREFIX)).AddType(TypeNames[leftExprs[len(leftExprs) - 1].Operator.Outputs[0].Type])
		name := MakeArgument(MakeGenSym(LOCAL_PREFIX), CurrentFile, LineNo).AddType(TypeNames[leftExprs[len(leftExprs)-1].Inputs[0].Type])

		name.Size = leftExprs[len(leftExprs)-1].Operator.Outputs[0].Size
		name.TotalSize = leftExprs[len(leftExprs)-1].Operator.Outputs[0].Size
		name.Type = leftExprs[len(leftExprs)-1].Operator.Outputs[0].Type
		name.Package = pkg
		name.IsShortDeclaration = true
		

		leftExprs[len(leftExprs)-1].Outputs = append(leftExprs[len(leftExprs)-1].Outputs, name)
	}

	if len(rightExprs[len(rightExprs)-1].Outputs) < 1 {
		// name := MakeArgument(MakeGenSym(LOCAL_PREFIX)).AddType(TypeNames[rightExprs[len(rightExprs) - 1].Operator.Outputs[0].Type])
		name := MakeArgument(MakeGenSym(LOCAL_PREFIX), CurrentFile, LineNo).AddType(TypeNames[rightExprs[len(rightExprs)-1].Inputs[0].Type])

		name.Size = rightExprs[len(rightExprs)-1].Operator.Outputs[0].Size
		name.TotalSize = rightExprs[len(rightExprs)-1].Operator.Outputs[0].Size
		name.Type = rightExprs[len(rightExprs)-1].Operator.Outputs[0].Type
		name.Package = pkg
		name.IsShortDeclaration = true

		rightExprs[len(rightExprs)-1].Outputs = append(rightExprs[len(rightExprs)-1].Outputs, name)
	}

	expr := MakeExpression(operator, CurrentFile, LineNo)
	expr.Package = pkg

	if len(leftExprs[len(leftExprs)-1].Outputs[0].Indexes) > 0 || leftExprs[len(leftExprs)-1].Operator != nil {
		// then it's a function call or an array access
		expr.AddInput(leftExprs[len(leftExprs)-1].Outputs[0])
		out = append(out, leftExprs...)
	} else {
		expr.Inputs = append(expr.Inputs, leftExprs[len(leftExprs)-1].Outputs[0])
	}

	if len(rightExprs[len(rightExprs)-1].Outputs[0].Indexes) > 0 || rightExprs[len(rightExprs)-1].Operator != nil {
		// then it's a function call or an array access
		expr.AddInput(rightExprs[len(rightExprs)-1].Outputs[0])
		out = append(out, rightExprs...)
	} else {
		expr.Inputs = append(expr.Inputs, rightExprs[len(rightExprs)-1].Outputs[0])
	}
	
	out = append(out, expr)

	return
}

func IsStrNil (byts []byte) bool {
	if len(byts) != 4 {
		return false
	}
	for _, byt := range byts {
		if byt != byte(0) {
			return false
		}
	}
	return true
}

// This function writes those bytes to PRGRM.Data
func WritePrimary(typ int, byts []byte, isGlobal bool) []*CXExpression {
	if pkg, err := PRGRM.GetCurrentPackage(); err == nil {
		arg := MakeArgument("", CurrentFile, LineNo)
		arg.AddType(TypeNames[typ])
		arg.Package = pkg
		arg.Program = PRGRM
		
		var size int

		size = len(byts)
		
		arg.Size = GetArgSize(typ)
		arg.TotalSize = size
		arg.Offset = DataOffset

		// arg.PassBy = PASSBY_REFERENCE
		
		if arg.Type == TYPE_STR {
			arg.PassBy = PASSBY_REFERENCE
			arg.Size = TYPE_POINTER_SIZE
			arg.TotalSize = TYPE_POINTER_SIZE
			// arg.IsPointer = true
			// arg.DereferenceOperations = append(arg.DereferenceOperations, DEREF_POINTER)
			// arg.DereferenceLevels++
		}

		for i, byt := range byts {
			PRGRM.Memory[DataOffset + i] = byt
		}
		DataOffset += size
		
		arg.PointeeSize = size

		expr := MakeExpression(nil, CurrentFile, LineNo)
		expr.Package = pkg
		expr.Outputs = append(expr.Outputs, arg)
		return []*CXExpression{expr}
	} else {
		panic(err)
	}
}

func TotalLength(lengths []int) int {
	var total int = 1
	for _, i := range lengths {
		total *= i
	}
	return total
}

func IterationExpressions(init []*CXExpression, cond []*CXExpression, incr []*CXExpression, statements []*CXExpression) []*CXExpression {
	jmpFn := Natives[OP_JMP]

	pkg, err := PRGRM.GetCurrentPackage()
	if err != nil {
		panic(err)
	}

	upExpr := MakeExpression(jmpFn, CurrentFile, LineNo)
	upExpr.Package = pkg

	trueArg := WritePrimary(TYPE_BOOL, encoder.Serialize(true), false)

	upLines := (len(statements) + len(incr) + len(cond) + 2) * -1
	downLines := 0

	upExpr.AddInput(trueArg[0].Outputs[0])
	upExpr.ThenLines = upLines
	upExpr.ElseLines = downLines

	downExpr := MakeExpression(jmpFn, CurrentFile, LineNo)
	downExpr.Package = pkg

	if len(cond[len(cond)-1].Outputs) < 1 {
		predicate := MakeArgument(MakeGenSym(LOCAL_PREFIX), CurrentFile, LineNo).AddType(TypeNames[cond[len(cond)-1].Operator.Outputs[0].Type])
		predicate.Package = pkg
		predicate.IsShortDeclaration = true
		cond[len(cond)-1].AddOutput(predicate)
		downExpr.AddInput(predicate)
	} else {
		predicate := cond[len(cond)-1].Outputs[0]
		predicate.Package = pkg
		predicate.IsShortDeclaration = true
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

func StructLiteralAssignment(to []*CXExpression, from []*CXExpression) []*CXExpression {
	for _, f := range from {
		f.Outputs[0].Name = to[0].Outputs[0].Name

		// flds := to[0].Outputs[0].Fields
		// flds = append(flds, f.Outputs[0].Fields...)
		// f.Outputs[0].Fields = flds

		if len(to[0].Outputs[0].Indexes) > 0 {
			f.Outputs[0].Lengths = to[0].Outputs[0].Lengths
			f.Outputs[0].Indexes = to[0].Outputs[0].Indexes
			f.Outputs[0].DereferenceOperations = append(f.Outputs[0].DereferenceOperations, DEREF_ARRAY)
		}

		// if !f.IsFlattened {
		// 	// if len(f.Outputs[0].Fields) > nestedLvl {
		// 	if lastCustomTypes == nil || lastCustomTypes[len(lastCustomTypes) - 1] != f.Outputs[0].CustomType {
		// 		// nestedLvl += 1
		// 		lastCustomTypes = append(lastCustomTypes, f.Outputs[0].CustomType)
		// 		nestedExprs[nestedLvl] = append(nestedExprs[len(f.Outputs[0].Fields)], f)
		// 	} else if len(f.Outputs[0].Fields) < nestedLvl {
		// 		// then current f is the field to assign to
		// 		for _, expr := range nestedExprs[nestedLvl - 1] {
		// 			// we add the previous expressions
		// 			result = append(result, expr)
		// 		}
				
		// 		for _, expr := range nestedExprs[nestedLvl] {
		// 			// we replace the field in the nested expressions for the current one
		// 			expr.Outputs[0].Fields[nestedLvl - 1] = f.Outputs[0].Fields[nestedLvl - 2]
		// 			expr.IsFlattened = true
		// 			result = append(result, expr)
		// 		}

		// 		// we reset those levels, in case we encounter another nested struct at the same level
		// 		nestedExprs[nestedLvl] = nil
		// 		nestedExprs[nestedLvl - 1] = nil
		// 		// and we decrease the nestedLvl
		// 		nestedLvl -= 1
		// 	} else {
		// 		nestedExprs[nestedLvl] = append(nestedExprs[len(f.Outputs[0].Fields)], f)
		// 	}
		// } else {
		// 	// it was already processed and assigned to its corresponding field
		// 	result = append(result, f)
		// }
		
		
		// if len(to[0].Outputs[0].Fields) > 0 {
		// 	f.Outputs[0].Fields = to[0].Outputs[0].Fields
		// 	f.Outputs[0].DereferenceOperations = append(f.Outputs[0].DereferenceOperations, DEREF_FIELD)
		// }

		// for _, fld := range to[0].Outputs[0].Fields {
		// 	f.Outputs[0].DereferenceOperations = append(f.Outputs[0].DereferenceOperations, DEREF_FIELD)
		// }


		// organizing possible nested structures
		// nestedStructs[f] = append(nestedStructs[f], f)






		// if len(f.Outputs[0].Fields) > lastFldLevel {
		// 	nested = append(nested, f)
		// }
		
		

		// if len(f.Outputs[0].Fields) == lastFldLevel {
		// 	for _, expr := range nested {
		// 		expr.Outputs[0].Fields[len(expr.Outputs[0].Fields) - 1] = f.Outputs[0].Fields[len(f.Outputs[0].Fields) - 1]
		// 		expr.Outputs[0].Fields[len(expr.Outputs[0].Fields) - 1].CustomType = f.Outputs[0].CustomType
		// 	}

		// 	nested = nil
		// 	// lastFldLevel = len(f.Outputs[0].Fields)
			
		// 	// nestedField[lastFldLevel] = f.Outputs[0].Fields[len(f.Outputs[0].Fields) - 1]
		// 	// nestedField[f] = f
		// }






		// lastFldLevel = len(f.Outputs[0].Fields)
		
		// if len(f.Outputs[0].Fields) > 1 {
		// 	// then we're dealing with nested structures
			
		// }

		f.Outputs[0].DereferenceOperations = append(f.Outputs[0].DereferenceOperations, DEREF_FIELD)
	}
	
	return from
}

func ArrayLiteralAssignment(to []*CXExpression, from []*CXExpression) []*CXExpression {
	for _, f := range from {
		f.Outputs[0].Name = to[0].Outputs[0].Name
		f.Outputs[0].DereferenceOperations = append(f.Outputs[0].DereferenceOperations, DEREF_ARRAY)
	}

	return from
}

func ShortAssign (expr *CXExpression, to []*CXExpression, from []*CXExpression, pkg *CXPackage, idx int) []*CXExpression {
	expr.AddInput(to[0].Outputs[0])
	expr.AddOutput(to[0].Outputs[0])
	expr.Package = pkg

	if from[idx].Operator == nil {
		expr.AddInput(from[idx].Outputs[0])
	} else {
		sym := MakeArgument(MakeGenSym(LOCAL_PREFIX), CurrentFile, LineNo).AddType(TypeNames[from[idx].Inputs[0].Type])
		sym.Package = pkg
		sym.IsShortDeclaration = true
		from[idx].AddOutput(sym)
		expr.AddInput(sym)
	}

	return append(from, expr)
}

func Assignment (to []*CXExpression, assignOp string, from []*CXExpression) []*CXExpression {
	idx := len(from) - 1

	if from[idx].IsArrayLiteral {
		from[0].Outputs[0].SynonymousTo = to[0].Outputs[0].Name
	}

	if pkg, err := PRGRM.GetCurrentPackage(); err == nil {

		var expr *CXExpression
		
		switch assignOp {
		case ":=":
			expr = MakeExpression(nil, CurrentFile, LineNo)
			expr.Package = pkg

			var sym *CXArgument

			if from[idx].Operator == nil {
				// then it's a literal
				sym = MakeArgument(to[0].Outputs[0].Name, CurrentFile, LineNo).AddType(TypeNames[from[idx].Outputs[0].Type])
			} else {
				sym = MakeArgument(to[0].Outputs[0].Name, CurrentFile, LineNo).AddType(TypeNames[from[idx].Inputs[0].Type])
			}
			sym.Package = pkg
			sym.IsShortDeclaration = true

			expr.AddOutput(sym)
			to = append([]*CXExpression{expr}, to...)
		case ">>=":
			expr = MakeExpression(Natives[OP_UND_BITSHR], CurrentFile, LineNo)
			return ShortAssign(expr, to, from, pkg, idx)
		case "<<=":
			expr = MakeExpression(Natives[OP_UND_BITSHL], CurrentFile, LineNo)
			return ShortAssign(expr, to, from, pkg, idx)
		case "+=":
			expr = MakeExpression(Natives[OP_UND_ADD], CurrentFile, LineNo)
			return ShortAssign(expr, to, from, pkg, idx)
		case "-=":
			expr = MakeExpression(Natives[OP_UND_SUB], CurrentFile, LineNo)
			return ShortAssign(expr, to, from, pkg, idx)
		case "*=":
			expr = MakeExpression(Natives[OP_UND_MUL], CurrentFile, LineNo)
			return ShortAssign(expr, to, from, pkg, idx)
		case "/=":
			expr = MakeExpression(Natives[OP_UND_DIV], CurrentFile, LineNo)
			return ShortAssign(expr, to, from, pkg, idx)
		case "%=":
			expr = MakeExpression(Natives[OP_UND_MOD], CurrentFile, LineNo)
			return ShortAssign(expr, to, from, pkg, idx)
		case "&=":
			expr = MakeExpression(Natives[OP_UND_BITAND], CurrentFile, LineNo)
			return ShortAssign(expr, to, from, pkg, idx)
		case "^=":
			expr = MakeExpression(Natives[OP_UND_BITXOR], CurrentFile, LineNo)
			return ShortAssign(expr, to, from, pkg, idx)
		case "|=":
			expr = MakeExpression(Natives[OP_UND_BITOR], CurrentFile, LineNo)
			return ShortAssign(expr, to, from, pkg, idx)
		}
	}
	
	if from[idx].Operator == nil {
		from[idx].Operator = Natives[OP_IDENTITY]
		to[0].Outputs[0].Size = from[idx].Outputs[0].Size
		to[0].Outputs[0].Type = from[idx].Outputs[0].Type
		to[0].Outputs[0].Lengths = from[idx].Outputs[0].Lengths
		to[0].Outputs[0].PassBy = from[idx].Outputs[0].PassBy
		to[0].Outputs[0].DoesEscape = from[idx].Outputs[0].DoesEscape
		to[0].Outputs[0].Program = PRGRM

		from[idx].Inputs = from[idx].Outputs
		from[idx].Outputs = to[len(to)-1].Outputs
		from[idx].Program = PRGRM

		return append(to[:len(to)-1], from...)
	} else {
		if from[idx].Operator.IsNative {
			// only assigning as if the operator had only one output defined
			to[0].Outputs[0].Size = Natives[from[idx].Operator.OpCode].Outputs[0].Size
			to[0].Outputs[0].Type = from[idx].Operator.Outputs[0].Type

			to[0].Outputs[0].Lengths = from[idx].Operator.Outputs[0].Lengths
			to[0].Outputs[0].DoesEscape = from[idx].Operator.Outputs[0].DoesEscape
			to[0].Outputs[0].PassBy = from[idx].Operator.Outputs[0].PassBy
			to[0].Outputs[0].Program = PRGRM
		} else {
			// we'll delegate multiple-value returns to the 'expression' grammar rule
			// only assigning as if the operator had only one output defined

			to[0].Outputs[0].Size = from[idx].Operator.Outputs[0].Size
			to[0].Outputs[0].Type = from[idx].Operator.Outputs[0].Type
			to[0].Outputs[0].Lengths = from[idx].Operator.Outputs[0].Lengths
			to[0].Outputs[0].DoesEscape = from[idx].Operator.Outputs[0].DoesEscape
			to[0].Outputs[0].PassBy = from[idx].Operator.Outputs[0].PassBy
			to[0].Outputs[0].Program = PRGRM
		}

		from[idx].Outputs = to[0].Outputs
		from[idx].Program = to[0].Program

		return append(to[:len(to)-1], from...)
	}
}

func SelectionExpressions(condExprs []*CXExpression, thenExprs []*CXExpression, elseExprs []*CXExpression) []*CXExpression {
	jmpFn := Natives[OP_JMP]
	pkg, err := PRGRM.GetCurrentPackage()
	if err != nil {
		panic(err)
	}
	ifExpr := MakeExpression(jmpFn, CurrentFile, LineNo)
	ifExpr.Package = pkg

	var predicate *CXArgument
	if condExprs[len(condExprs)-1].Operator == nil {
		// then it's a literal
		predicate = condExprs[len(condExprs)-1].Outputs[0]
	} else {
		// then it's an expression
		predicate = MakeArgument(MakeGenSym(LOCAL_PREFIX), CurrentFile, LineNo).AddType(TypeNames[condExprs[len(condExprs)-1].Operator.Outputs[0].Type])
		predicate.IsShortDeclaration = true
		condExprs[len(condExprs)-1].Outputs = append(condExprs[len(condExprs)-1].Outputs, predicate)
	}
	predicate.Package = pkg

	ifExpr.AddInput(predicate)

	thenLines := 0
	elseLines := len(thenExprs) + 1

	ifExpr.ThenLines = thenLines
	ifExpr.ElseLines = elseLines

	skipExpr := MakeExpression(jmpFn, CurrentFile, LineNo)
	skipExpr.Package = pkg

	trueArg := WritePrimary(TYPE_BOOL, encoder.Serialize(true), false)
	skipLines := len(elseExprs)

	skipExpr.AddInput(trueArg[0].Outputs[0])
	skipExpr.ThenLines = skipLines
	skipExpr.ElseLines = 0

	var exprs []*CXExpression
	if condExprs[len(condExprs)-1].Operator != nil {
		exprs = append(exprs, condExprs...)
	}
	exprs = append(exprs, ifExpr)
	exprs = append(exprs, thenExprs...)
	exprs = append(exprs, skipExpr)
	exprs = append(exprs, elseExprs...)

	return exprs
}

func GetSymType(sym *CXArgument, fn *CXFunction) int {
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

func SetFinalSize(symbols *map[string]*CXArgument, sym *CXArgument) {
	var elt *CXArgument
	var finalSize int = sym.TotalSize

	var fldIdx int
	elt = sym

	if arg, found := (*symbols)[sym.Package.Name+"."+sym.Name]; found {
		for _, op := range sym.DereferenceOperations {
			switch op {
			case DEREF_ARRAY:
				if arg.IsSlice {
					continue
				}
				var subSize int = 1
				for _, len := range elt.Lengths[:len(elt.Indexes)] {
					subSize *= len
				}
				finalSize /= subSize
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

func GetGlobalSymbol(symbols *map[string]*CXArgument, symPackage *CXPackage, symName string) {
	if _, found := (*symbols)[symPackage.Name+"."+symName]; !found {
		if glbl, err := symPackage.GetGlobal(symName); err == nil {
			(*symbols)[symPackage.Name+"."+symName] = glbl
		}
	}
}

func GiveOffset(symbols *map[string]*CXArgument, sym *CXArgument, offset *int, shouldExist bool) {
	if sym.Name != "" {
		GetGlobalSymbol(symbols, sym.Package, sym.Name)

		if arg, found := (*symbols)[sym.Package.Name+"."+sym.Name]; !found {
			panic("")
		} else {
			// if sym.IsReference {
			// 	if arg.HeapOffset < 1 && !arg.IsSlice {
			// 		// then it hasn't been assigned
			// 		// an offset of 0 is impossible because the symbol was declared before

			// 		arg.HeapOffset = *offset
			// 		// sym.HeapOffset = *offset
			// 		*offset += TYPE_POINTER_SIZE
			// 	}

			// 	// if not, then it has been assigned before
			// 	// and we just reassign it to this symbol
			// 	// we'll do this below, where we're assigning everything to sym
			// }

			// identifying customtypes of fields if they are nil
			if len(sym.Fields) > 0 {
				strct := arg.CustomType

				for c := len(sym.Fields) - 1; c >= 0; c-- {
					if sym.Fields[c].CustomType != nil {
						strct = sym.Fields[c].CustomType
					}
					if inFld, err := strct.GetField(sym.Fields[c].Name); err == nil {
						
						sym.Fields[c].CustomType = strct
						sym.Fields[c].Type = inFld.Type
						
						strct = inFld.CustomType
					}
				}
			}
			
			var isFieldPointer bool
			if len(sym.Fields) > 0 {
				var found bool

				strct := arg.CustomType
				for _, nameFld := range sym.Fields {
					if nameFld.CustomType != nil {
						strct = nameFld.CustomType
					}
					
					for _, fld := range strct.Fields {
						if nameFld.Name == fld.Name {
							if fld.IsPointer {
								sym.IsPointer = true
								isFieldPointer = true
							}
							
							found = true
							break
						}
					}

					if !found {
						panic("field '" + nameFld.Name + "' not found")
					}
				}
			}

			if sym.DereferenceLevels > 0 {
				if arg.IndirectionLevels >= sym.DereferenceLevels || isFieldPointer {
					pointer := arg

					for c := 0; c < sym.DereferenceLevels-1; c++ {
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

			// checking if it's accessing fields
			if len(sym.Fields) > 0 {
				var found bool

				strct := arg.CustomType
				for _, nameFld := range sym.Fields {
					if nameFld.CustomType != nil {
						strct = nameFld.CustomType
					}
					
					for _, fld := range strct.Fields {
						if nameFld.Name == fld.Name {
							nameFld.Lengths = fld.Lengths
							nameFld.Size = fld.Size
							nameFld.TotalSize = fld.TotalSize
							nameFld.DereferenceLevels = sym.DereferenceLevels
							nameFld.IsPointer = fld.IsPointer
							nameFld.CustomType = fld.CustomType
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

			if len(sym.Fields) > 0 {
				sym.Type = sym.Fields[len(sym.Fields) - 1].Type
				// if sym.Fields[len(sym.Fields) - 1].CustomType != nil {
				// 	// fmt.Println("huehue")
				// 	sym.Type = TYPE_CUSTOM
				// }
			} else {
				sym.Type = arg.Type
			}

			sym.IsSlice = arg.IsSlice
			sym.CustomType = arg.CustomType
			sym.Pointee = arg.Pointee
			sym.Lengths = arg.Lengths
			sym.PointeeSize = arg.PointeeSize
			sym.Package = arg.Package
			sym.Program = arg.Program
			sym.DoesEscape = arg.DoesEscape

			if sym.IsReference && !arg.IsStruct {
				sym.TotalSize = arg.TotalSize
				sym.Size = arg.Size
			} else {
				// we need to implement a more robust system, like the one in op.go
				if len(sym.Fields) > 0 {
					sym.Size = arg.Size
					sym.TotalSize = sym.Fields[len(sym.Fields)-1].TotalSize
				} else {
					sym.Size = arg.Size
					sym.TotalSize = arg.TotalSize
				}
			}
		}
	}
}

func ProcessTempVariable(expr *CXExpression) {
	if expr.Operator != nil && (expr.Operator == Natives[OP_IDENTITY] || IsUndOp(expr.Operator)) && len(expr.Outputs) > 0 && len(expr.Inputs) > 0 {
		name := expr.Outputs[0].Name
		if len(name) >= len(LOCAL_PREFIX) && name[:len(LOCAL_PREFIX)] == LOCAL_PREFIX {
			// then it's a temporary variable and it needs to adopt its input's type
			expr.Outputs[0].Type = expr.Inputs[0].Type
			expr.Outputs[0].Size = expr.Inputs[0].Size
			expr.Outputs[0].TotalSize = expr.Inputs[0].TotalSize
			expr.Outputs[0].IsShortDeclaration = true
		}
	}
}

// func ProcessShortDeclaration(expr *CXExpression) {
// 	if expr.IsShortDeclaration {
// 		expr.Outputs[0].Type = expr.Inputs[0].Type
// 		expr.Outputs[0].Size = expr.Inputs[0].Size
// 		expr.Outputs[0].TotalSize = expr.Inputs[0].TotalSize

// 		expr.Outputs[0].Lengths = expr.Inputs[0].Lengths
// 		expr.Outputs[0].Fields = expr.Inputs[0].Fields
// 	}
// }

func ProcessMethodCall(expr *CXExpression, symbols *map[string]*CXArgument, sym *CXArgument, offset *int, shouldExist bool) {
	if sym.Name != "" {
		if arg, found := (*symbols)[sym.Package.Name+"."+sym.Name]; !found {
			panic("")
		} else {
			if expr.IsMethodCall {
				if len(sym.Fields) > 0 {
					// var found bool

					strct := arg.CustomType

					if fn, err := PRGRM.GetFunction(strct.Name + "." + sym.Fields[len(sym.Fields) - 1].Name, sym.Package.Name); err == nil {
						expr.Operator = fn
					} else {
						panic("")
					}

					// expr.Operator = MakeFunction(strct.Name + "." + sym.Fields[len(sym.Fields) - 1].Name)
					sym.Fields = sym.Fields[:len(sym.Fields) - 1]
					sym.DereferenceOperations = sym.DereferenceOperations[:len(sym.DereferenceOperations) - 1]

					expr.Outputs = nil
				}
			}
		}
	}
}

func UpdateSymbolsTable(symbols *map[string]*CXArgument, sym *CXArgument, offset *int, shouldExist bool) {
	if sym.Name != "" {
		GetGlobalSymbol(symbols, sym.Package, sym.Name)

		if _, found := (*symbols)[sym.Package.Name+"."+sym.Name]; !found {
			if shouldExist {
				// it should exist. error
				println(ErrorHeader(sym.FileName, sym.FileLine) + " identifier '" + sym.Name + "' does not exist")
				os.Exit(3)
			}

			if sym.SynonymousTo != "" {
				// then the offset needs to be shared
				GetGlobalSymbol(symbols, sym.Package, sym.SynonymousTo)
				sym.Offset = (*symbols)[sym.Package.Name+"."+sym.SynonymousTo].Offset

				(*symbols)[sym.Package.Name+"."+sym.Name] = sym
			} else {
				sym.Offset = *offset
				(*symbols)[sym.Package.Name+"."+sym.Name] = sym

				if sym.IsSlice {
					*offset += sym.Size
				} else {
					*offset += sym.TotalSize
				}

				if sym.IsPointer {
					pointer := sym
					for c := 0; c < sym.IndirectionLevels-1; c++ {
						pointer = pointer.Pointee
						pointer.Offset = *offset
						*offset += pointer.TotalSize
					}
				}
			}
		}
	}
}

func GiveCustomType(symbols *map[string]*CXArgument, sym *CXArgument, offset *int, shouldExist bool) {
	if sym.Name != "" {
		
	}
}

func ProcessInputSlice(inp *CXArgument) {
	if inp.IsSlice {
		inp.DereferenceOperations = append([]int{DEREF_POINTER}, inp.DereferenceOperations...)
		inp.DereferenceLevels++
	}
}

func ProcessOutputSlice(out *CXArgument) {
	// if out.IsSlice && len(out.DereferenceOperations) < 1 {
	// 	out.PassBy = PASSBY_REFERENCE
	// }
	if out.IsSlice && len(out.DereferenceOperations) > 0 {
		out.DereferenceOperations = append([]int{DEREF_POINTER}, out.DereferenceOperations...)
		out.DereferenceLevels++
	}
}

func ProcessSliceAssignment(expr *CXExpression) {
	if expr.Operator == Natives[OP_IDENTITY] && expr.Outputs[0].IsSlice && expr.Inputs[0].IsSlice {
		expr.Inputs[0].DereferenceOperations = expr.Inputs[0].DereferenceOperations[1:]
	}
}

func CheckTypes(expr *CXExpression) {
	if expr.Operator != nil && expr.Operator.IsNative && expr.Operator.OpCode == OP_IDENTITY {
		for i, inp := range expr.Inputs {
			// if expr.Outputs[i].Type != inp.Type || expr.Outputs[i].TotalSize != inp.TotalSize {
			if expr.Outputs[i].Type != inp.Type {

				var expectedType string
				var receivedType string
				if expr.Operator.Inputs[i].CustomType != nil {
					// then it's custom type
					expectedType = expr.Outputs[i].CustomType.Name
				} else {
					// then it's native type
					expectedType = TypeNames[expr.Outputs[i].Type]
				}

				if expr.Inputs[i].CustomType != nil {
					// then it's custom type
					receivedType = inp.CustomType.Name
				} else {
					// then it's native type
					receivedType = TypeNames[inp.Type]
				}

				
				if expr.IsStructLiteral {
					println(ErrorHeader(expr.Outputs[i].FileName, expr.Outputs[i].FileLine), fmt.Sprintf("field '%s' in struct literal of type '%s' expected argument of type '%s'; '%s' was provided", expr.Outputs[i].Fields[0].Name, expr.Outputs[i].CustomType.Name, expectedType, receivedType))
				} else {
					println(ErrorHeader(expr.Outputs[i].FileName, expr.Outputs[i].FileLine), fmt.Sprintf("trying to assign argument of type '%s' to symbol '%s' of type '%s'", receivedType, expr.Outputs[i].Name, expectedType))
				}
				// os.Exit(3)

				FoundCompileErrors = true
			}
		}
	}

	// checking inputs matching operator's inputs
	if expr.Operator != nil {
		// then it's a function call and not a declaration
		for i, inp := range expr.Operator.Inputs {
			if inp.Type != expr.Inputs[i].Type && inp.Type != TYPE_UNDEFINED {
				var opName string
				if expr.Operator.IsNative {
					opName = OpNames[expr.Operator.OpCode]
				} else {
					opName = expr.Operator.Name
				}

				var expectedType string
				var receivedType string
				if expr.Operator.Inputs[i].CustomType != nil {
					// then it's custom type
					expectedType = expr.Operator.Inputs[i].CustomType.Name
				} else {
					// then it's native type
					expectedType = TypeNames[expr.Operator.Inputs[i].Type]
				}

				if expr.Inputs[i].CustomType != nil {
					// then it's custom type
					receivedType = expr.Inputs[i].CustomType.Name
				} else {
					// then it's native type
					receivedType = TypeNames[expr.Inputs[i].Type]
				}

				println(ErrorHeader(expr.Inputs[i].FileName, expr.Inputs[i].FileLine), fmt.Sprintf("function '%s' expected input argument of type '%s'; '%s' was provided", opName, expectedType, receivedType))

				FoundCompileErrors = true
				// os.Exit(3)
			}
		}
	}
}

func FunctionDeclaration(fn *CXFunction, inputs []*CXArgument, outputs []*CXArgument, exprs []*CXExpression) {
	// adding inputs, outputs
	for _, inp := range inputs {
		fn.AddInput(inp)
	}
	
	for _, out := range outputs {
		fn.AddOutput(out)
	}

	for _, out := range fn.Outputs {
		if out.IsPointer && out.Type != TYPE_STR {
			out.DoesEscape = true
		}
	}

	// getting offset to use by statements (excluding inputs, outputs and receiver)
	var offset int
	PRGRM.HeapStartsAt = DataOffset

	for i, expr := range exprs {
		if expr.Label != "" && expr.Operator == Natives[OP_JMP] {
			// then it's a goto
			for j, e := range exprs {
				if e.Label == expr.Label && i != j {
					// ElseLines is used because arg's default val is false
					expr.ThenLines = j - i - 1
				}
			}
		}

		fn.AddExpression(expr)
	}

	fn.Length = len(fn.Expressions)

	var symbols map[string]*CXArgument = make(map[string]*CXArgument, 0)
	var symbolsScope map[string]bool = make(map[string]bool, 0)

	for _, inp := range fn.Inputs {
		if inp.IsLocalDeclaration {
			symbolsScope[inp.Package.Name+"."+inp.Name] = true
		}
		inp.IsLocalDeclaration = symbolsScope[inp.Package.Name+"."+inp.Name]

		UpdateSymbolsTable(&symbols, inp, &offset, false)
		GiveOffset(&symbols, inp, &offset, false)
		SetFinalSize(&symbols, inp)

		AddPointer(fn, inp)
	}
	for _, out := range fn.Outputs {
		if out.IsLocalDeclaration {
			symbolsScope[out.Package.Name+"."+out.Name] = true
		}
		out.IsLocalDeclaration = symbolsScope[out.Package.Name+"."+out.Name]

		UpdateSymbolsTable(&symbols, out, &offset, false)
		GiveOffset(&symbols, out, &offset, false)
		SetFinalSize(&symbols, out)

		AddPointer(fn, out)
	}

	for _, expr := range fn.Expressions {
		// ProcessShortDeclaration(expr)
		for _, inp := range expr.Inputs {
			if inp.IsLocalDeclaration {
				symbolsScope[inp.Package.Name+"."+inp.Name] = true
			}
			inp.IsLocalDeclaration = symbolsScope[inp.Package.Name+"."+inp.Name]

			UpdateSymbolsTable(&symbols, inp, &offset, true)
			ProcessMethodCall(expr, &symbols, inp, &offset, true)
			GiveOffset(&symbols, inp, &offset, true)

			// we only need the input to be processed by GiveOffset to call ProcessMethodCall
			// ProcessMethodCall(expr)
			
			SetFinalSize(&symbols, inp)
			ProcessInputSlice(inp)

			for _, idx := range inp.Indexes {
				UpdateSymbolsTable(&symbols, idx, &offset, true)
				GiveOffset(&symbols, idx, &offset, true)
			}

			AddPointer(fn, inp)
		}
		for _, out := range expr.Outputs {
			if out.IsLocalDeclaration {
				symbolsScope[out.Package.Name+"."+out.Name] = true
			}

			out.IsLocalDeclaration = symbolsScope[out.Package.Name + "." + out.Name]

			if out.IsShortDeclaration {
				UpdateSymbolsTable(&symbols, out, &offset, false)
			} else {
				UpdateSymbolsTable(&symbols, out, &offset, true)
			}
			
			ProcessMethodCall(expr, &symbols, out, &offset, true)
			GiveOffset(&symbols, out, &offset, false)
			
			SetFinalSize(&symbols, out)
			ProcessOutputSlice(out)
			
			for _, idx := range out.Indexes {
				UpdateSymbolsTable(&symbols, idx, &offset, true)
				GiveOffset(&symbols, idx, &offset, true)
			}

			AddPointer(fn, out)
		}

		SetCorrectArithmeticOp(expr)
		ProcessTempVariable(expr)
		ProcessSliceAssignment(expr)
		CheckTypes(expr)

		// var affs []*CXAffordance
		// affs = expr.GetAffordances([]string{})
		// for _, aff := range affs {
		// 	fmt.Println("aff", aff.Description)
		// }
	}

	fn.Size = offset
}

func FunctionCall(exprs []*CXExpression, args []*CXExpression) []*CXExpression {
	expr := exprs[len(exprs)-1]

	if expr.Operator == nil {
		opName := expr.Outputs[0].Name
		opPkg := expr.Outputs[0].Package
		// if len(expr.Outputs[0].Fields) > 0 {
		// 	opName = expr.Outputs[0].Fields[0].Name
		// 	// it wasn't a field, but a method call. removing it as a field
		// 	// expr.Outputs[0].Fields = expr.Outputs[0].Fields[:len(expr.Outputs[0].Fields)-1]
		// 	// we remove information about the "field" (method name)
		// 	expr.AddInput(expr.Outputs[0])

		// 	expr.Outputs = expr.Outputs[:len(expr.Outputs)-1]
		// 	// expr.Inputs = expr.Inputs[:len(expr.Inputs) - 1]
		// 	// expr.AddInput(expr.Outputs[0])
		// }

		if op, err := PRGRM.GetFunction(opName, opPkg.Name); err == nil {
			expr.Operator = op
		} else {
			println(ErrorHeader(CurrentFile, LineNo), err.Error())
			os.Exit(3)
			return nil
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
					out = MakeArgument(MakeGenSym(LOCAL_PREFIX), CurrentFile, LineNo).AddType(TypeNames[inpExpr.Inputs[0].Type])
					out.Size = inpExpr.Inputs[0].Size
					out.TotalSize = inpExpr.Inputs[0].Size
					out.Type = inpExpr.Inputs[0].Type
					out.IsShortDeclaration = true
				} else {
					out = MakeArgument(MakeGenSym(LOCAL_PREFIX), CurrentFile, LineNo).AddType(TypeNames[inpExpr.Operator.Outputs[0].Type])
					out.Size = inpExpr.Operator.Outputs[0].Size
					out.TotalSize = inpExpr.Operator.Outputs[0].Size
					out.Type = inpExpr.Operator.Outputs[0].Type
					out.IsShortDeclaration = true
				}

				out.Package = inpExpr.Package
				inpExpr.AddOutput(out)
				expr.AddInput(out)
			}
			if len(inpExpr.Outputs) > 0 && inpExpr.IsArrayLiteral {
				expr.AddInput(inpExpr.Outputs[0])
			}
			nestedExprs = append(nestedExprs, inpExpr)

		}
	}

	return append(nestedExprs, exprs...)
}
