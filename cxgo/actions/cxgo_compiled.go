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
var IdeMode bool
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

var dStack bool = false
var InFn bool = false
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
	FoundCompileErrors = true
	return "error: " + currentFile + ":" + strconv.FormatInt(int64(lineNo), 10)
}

// this function adds the roots (pointers) for some GC algorithms
func AddPointer(fn *CXFunction, sym *CXArgument) {
	if sym.IsPointer && sym.Name != "" {
		var found bool
		for _, ptr := range fn.ListOfPointers {
			if sym.Name == ptr.Name {
				found = true
				break
			}
		}
		if !found {
			fn.ListOfPointers = append(fn.ListOfPointers, sym)
		}
	}
}

func DeclareGlobal(declarator *CXArgument, declaration_specifiers *CXArgument, initializer []*CXExpression, doesInitialize bool) {
	if pkg, err := PRGRM.GetCurrentPackage(); err == nil {
		declaration_specifiers.Package = pkg
		
		if glbl, err := PRGRM.GetGlobal(declarator.Name); err == nil {
			// then it is already defined

			if glbl.Offset < 0 || glbl.Size == 0 || glbl.TotalSize == 0 {
				// then it was only added a reference to the symbol
				var offExpr []*CXExpression
				if declaration_specifiers.IsSlice {
					offExpr = WritePrimary(declaration_specifiers.Type, make([]byte, declaration_specifiers.Size), true)
				} else {
					offExpr = WritePrimary(declaration_specifiers.Type, make([]byte, declaration_specifiers.TotalSize), true)
				}

				glbl.Offset = offExpr[0].Outputs[0].Offset
				glbl.PassBy = offExpr[0].Outputs[0].PassBy
			}

			if doesInitialize {
				// then we just re-assign offsets
				if initializer[len(initializer)-1].Operator == nil {
					// then it's a literal
					declaration_specifiers.Name = glbl.Name
					declaration_specifiers.Offset = glbl.Offset
					declaration_specifiers.PassBy = glbl.PassBy

					*glbl = *declaration_specifiers

					initializer[len(initializer) - 1].AddInput(initializer[len(initializer) - 1].Outputs[0])
					initializer[len(initializer) - 1].Outputs = nil
					initializer[len(initializer) - 1].AddOutput(glbl)
					initializer[len(initializer) - 1].Operator = Natives[OP_IDENTITY]

					SysInitExprs = append(SysInitExprs, initializer...)
				} else {
					// then it's an expression
					declaration_specifiers.Name = glbl.Name
					declaration_specifiers.Offset = glbl.Offset
					declaration_specifiers.PassBy = glbl.PassBy

					*glbl = *declaration_specifiers
					
					if initializer[len(initializer) - 1].IsStructLiteral {
						initializer = StructLiteralAssignment([]*CXExpression{&CXExpression{Outputs: []*CXArgument{glbl}}}, initializer)
					} else {
						initializer[len(initializer) - 1].Outputs = nil
						initializer[len(initializer) - 1].AddOutput(glbl)
					}
					
					SysInitExprs = append(SysInitExprs, initializer...)
				}
			} else {
				// we keep the last value for now
				declaration_specifiers.Name = glbl.Name
				declaration_specifiers.Offset = glbl.Offset
				declaration_specifiers.PassBy = glbl.PassBy
				*glbl = *declaration_specifiers
			}
		} else {
			// then it hasn't been defined
			var offExpr []*CXExpression
			if declaration_specifiers.IsSlice {
				offExpr = WritePrimary(declaration_specifiers.Type, make([]byte, declaration_specifiers.Size), true)
			} else {
				offExpr = WritePrimary(declaration_specifiers.Type, make([]byte, declaration_specifiers.TotalSize), true)
			}
			if doesInitialize {
				if initializer[len(initializer)-1].Operator == nil {
					// then it's a literal

					declaration_specifiers.Name = declarator.Name
					declaration_specifiers.Offset = offExpr[0].Outputs[0].Offset
					declaration_specifiers.Size = offExpr[0].Outputs[0].Size
					declaration_specifiers.TotalSize = offExpr[0].Outputs[0].TotalSize
					declaration_specifiers.Package = pkg

					initializer[len(initializer) - 1].Operator = Natives[OP_IDENTITY]
					initializer[len(initializer) - 1].AddInput(initializer[len(initializer) - 1].Outputs[0])
					initializer[len(initializer) - 1].Outputs = nil
					initializer[len(initializer) - 1].AddOutput(declaration_specifiers)
					
					pkg.AddGlobal(declaration_specifiers)

					SysInitExprs = append(SysInitExprs, initializer...)
				} else {
					// then it's an expression
					declaration_specifiers.Name = declarator.Name
					declaration_specifiers.Offset = offExpr[0].Outputs[0].Offset
					declaration_specifiers.Size = offExpr[0].Outputs[0].Size
					declaration_specifiers.TotalSize = offExpr[0].Outputs[0].TotalSize
					declaration_specifiers.Package = pkg

					if initializer[len(initializer) - 1].IsStructLiteral {
						initializer = StructLiteralAssignment([]*CXExpression{&CXExpression{Outputs: []*CXArgument{declaration_specifiers}}}, initializer)
					} else {
						initializer[len(initializer) - 1].Outputs = nil
						initializer[len(initializer) - 1].AddOutput(declaration_specifiers)
					}

					pkg.AddGlobal(declaration_specifiers)
					SysInitExprs = append(SysInitExprs, initializer...)
				}
			} else {
				// offExpr := WritePrimary(declaration_specifiers.Type, make([]byte, declaration_specifiers.Size), true)
				// exprOut := expr[0].Outputs[0]

				declaration_specifiers.Name = declarator.Name
				declaration_specifiers.Offset = offExpr[0].Outputs[0].Offset
				declaration_specifiers.Size = offExpr[0].Outputs[0].Size
				declaration_specifiers.TotalSize = offExpr[0].Outputs[0].TotalSize
				declaration_specifiers.Package = pkg
				
				pkg.AddGlobal(declaration_specifiers)
			}
		}
	} else {
		panic(err)
	}
}

func DeclareStruct (ident string, strctFlds []*CXArgument) {
	if pkg, err := PRGRM.GetCurrentPackage(); err == nil {
		if strct, err := PRGRM.GetStruct(ident, pkg.Name); err == nil {
			strct.Fields = nil
			strct.Size = 0

			// var size int
			for _, fld := range strctFlds {
				strct.AddField(fld)
				// size += fld.TotalSize
			}
			// strct.Size = size
		} else {
			panic(err)
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
		PRGRM.SelectPackage(pkg.Name)
	} else {
		PRGRM.SelectPackage(pkg.Name)
	}
}

func AffordanceStructs (pkg *CXPackage) {
	// Argument type
	argStrct := MakeStruct("Argument")
	// argStrct.Size = GetArgSize(TYPE_STR) + GetArgSize(TYPE_STR)
	
	argFldName := MakeField("Name", TYPE_STR, "", 0)
	argFldName.TotalSize = GetArgSize(TYPE_STR)
	argFldIndex := MakeField("Index", TYPE_I32, "", 0)
	argFldIndex.TotalSize = GetArgSize(TYPE_I32)
	argFldType := MakeField("Type", TYPE_STR, "", 0)
	argFldType.TotalSize = GetArgSize(TYPE_STR)
	
	argStrct.AddField(argFldName)
	argStrct.AddField(argFldIndex)
	argStrct.AddField(argFldType)
	
	pkg.AddStruct(argStrct)

	// Expression type
	exprStrct := MakeStruct("Expression")
	// exprStrct.Size = GetArgSize(TYPE_STR)
	
	exprFldOperator := MakeField("Operator", TYPE_STR, "", 0)
	
	exprStrct.AddField(exprFldOperator)
	
	pkg.AddStruct(exprStrct)

	// Function type
	fnStrct := MakeStruct("Function")
	// fnStrct.Size = GetArgSize(TYPE_STR) + GetArgSize(TYPE_STR) + GetArgSize(TYPE_STR)
	
	fnFldName := MakeField("Name", TYPE_STR, "", 0)
	fnFldName.TotalSize = GetArgSize(TYPE_STR)
	
	fnFldInpSig := MakeField("InputSignature", TYPE_STR, "", 0)
	fnFldInpSig.Size = GetArgSize(TYPE_STR)
	fnFldInpSig = DeclarationSpecifiers(fnFldInpSig, 0, DECL_SLICE)

	fnFldOutSig := MakeField("OutputSignature", TYPE_STR, "", 0)
	fnFldOutSig.Size = GetArgSize(TYPE_STR)
	fnFldOutSig = DeclarationSpecifiers(fnFldOutSig, 0, DECL_SLICE)

	fnStrct.AddField(fnFldName)
	fnStrct.AddField(fnFldInpSig)

	fnStrct.AddField(fnFldOutSig)
	
	pkg.AddStruct(fnStrct)
	
	// Structure type
	strctStrct := MakeStruct("Structure")
	// strctStrct.Size = GetArgSize(TYPE_STR)
	
	strctFldName := MakeField("Name", TYPE_STR, "", 0)
	strctFldName.TotalSize = GetArgSize(TYPE_STR)
	
	strctStrct.AddField(strctFldName)
	
	pkg.AddStruct(strctStrct)
	
	// Package type
	pkgStrct := MakeStruct("Structure")
	// pkgStrct.Size = GetArgSize(TYPE_STR)
	
	pkgFldName := MakeField("Name", TYPE_STR, "", 0)
	
	pkgStrct.AddField(pkgFldName)
	
	pkg.AddStruct(pkgStrct)

	// Caller type
	callStrct := MakeStruct("Caller")
	// callStrct.Size = GetArgSize(TYPE_STR) + GetArgSize(TYPE_I32)
	
	callFldFnName := MakeField("FnName", TYPE_STR, "", 0)
	callFldFnName.TotalSize = GetArgSize(TYPE_STR)
	callFldFnSize := MakeField("FnSize", TYPE_I32, "", 0)
	callFldFnSize.TotalSize = GetArgSize(TYPE_I32)
	
	callStrct.AddField(callFldFnName)
	callStrct.AddField(callFldFnSize)

	pkg.AddStruct(callStrct)
	
	// Program type
	prgrmStrct := MakeStruct("Program")
	// prgrmStrct.Size = GetArgSize(TYPE_I32) + GetArgSize(TYPE_I64)
	
	prgrmFldCallCounter := MakeField("CallCounter", TYPE_I32, "", 0)
	prgrmFldCallCounter.TotalSize = GetArgSize(TYPE_I32)
	prgrmFldFreeHeap := MakeField("HeapUsed", TYPE_I64, "", 0)
	prgrmFldFreeHeap.TotalSize = GetArgSize(TYPE_I64)

	// prgrmFldCaller := MakeField("Caller", TYPE_CUSTOM, "", 0)
	prgrmFldCaller := DeclarationSpecifiersStruct(callStrct.Name, callStrct.Package.Name, false)
	prgrmFldCaller.Name = "Caller"

	prgrmStrct.AddField(prgrmFldCallCounter)
	prgrmStrct.AddField(prgrmFldFreeHeap)
	prgrmStrct.AddField(prgrmFldCaller)
	
	pkg.AddStruct(prgrmStrct)
}

func DeclareImport(ident string, currentFile string, lineNo int) {
	if pkg, err := PRGRM.GetCurrentPackage(); err == nil {
		if _, err := pkg.GetImport(ident); err != nil {
			
			if imp, err := PRGRM.GetPackage(ident); err == nil {
				pkg.AddImport(imp)
			} else {
				// TODO look in the workspace
				if IsCorePackage(ident) {
					imp := MakePackage(ident)
					pkg.AddImport(imp)
					PRGRM.AddPackage(imp)
					PRGRM.CurrentPackage = pkg

					if ident == "aff" {
						AffordanceStructs(imp)
					}
				} else {
					println(ErrorHeader(currentFile, lineNo), err.Error())
				}
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
			declSpec.Size = TYPE_POINTER_SIZE
			declSpec.TotalSize = TYPE_POINTER_SIZE
			declSpec.IndirectionLevels++
		} else {
			pointer := declSpec

			for c := declSpec.IndirectionLevels - 1; c > 0; c-- {
				pointer.IndirectionLevels = c
				pointer.IsPointer = true
			}

			declSpec.IndirectionLevels++

			pointer.Size = TYPE_POINTER_SIZE
			pointer.TotalSize = TYPE_POINTER_SIZE
		}

		return declSpec
	case DECL_ARRAY:
		declSpec.DeclarationSpecifiers = append(declSpec.DeclarationSpecifiers, DECL_ARRAY)
		arg := declSpec
		arg.IsArray = true
		arg.Lengths = append([]int{arraySize}, arg.Lengths...)
		arg.TotalSize = arg.Size * TotalLength(arg.Lengths)

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

	arg.Size = GetArgSize(typ)

	if typ == TYPE_AFF {
		// equivalent to slice of strings
		return DeclarationSpecifiers(arg, 0, DECL_SLICE)
	}
	
	return DeclarationSpecifiers(arg, 0, DECL_BASIC)
}

func DeclarationSpecifiersStruct (ident string, pkgName string, isExternal bool) *CXArgument {
	if isExternal {
		// custom type in an imported package
		if pkg, err := PRGRM.GetCurrentPackage(); err == nil {
			if imp, err := pkg.GetImport(pkgName); err == nil {
				if strct, err := PRGRM.GetStruct(ident, imp.Name); err == nil {
					arg := MakeArgument("", CurrentFile, LineNo)
					arg.Type = TYPE_CUSTOM
					arg.CustomType = strct
					arg.Size = strct.Size
					arg.TotalSize = strct.Size

					arg.Package = pkg
					arg.DeclarationSpecifiers = append(arg.DeclarationSpecifiers, DECL_STRUCT)

					return arg
				} else {
					println(ErrorHeader(CurrentFile, LineNo), err.Error())
					return nil
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
				arg.Type = TYPE_CUSTOM
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

	arrVarExpr := MakeExpression(nil, CurrentFile, LineNo)
	arrVarExpr.Package = pkg
	arrVar := MakeArgument(symName, CurrentFile, LineNo)
	arrVar = DeclarationSpecifiers(arrVar, arrSize, DECL_ARRAY)
	arrVar.AddType(TypeNames[typSpec])
	arrVar.TotalSize = arrVar.Size * TotalLength(arrVar.Lengths)

	arrVarExpr.Outputs = append(arrVarExpr.Outputs, arrVar)
	arrVar.Package = pkg
	arrVar.PreviouslyDeclared = true

	result = append(result, arrVarExpr)

	var endPointsCounter int
	for _, expr := range exprs {
		if expr.IsArrayLiteral {
			expr.IsArrayLiteral = false

			sym := MakeArgument(symName, CurrentFile, LineNo).AddType(TypeNames[typSpec])
			sym.Package = pkg
			sym.PreviouslyDeclared = true

			if sym.Type == TYPE_STR || sym.Type == TYPE_AFF {
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
				symExpr.Inputs = expr.Outputs
			} else {
				symExpr.Operator = expr.Operator
				symExpr.Inputs = expr.Inputs

				// hack to get the correct lengths below
				expr.Outputs = append(expr.Outputs, sym)
			}
			
			result = append(result, symExpr)

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
	symOutput.PreviouslyDeclared = true
	symOutput.TotalSize = symOutput.Size * TotalLength(symOutput.Lengths)

	symInput := MakeArgument(symName, CurrentFile, LineNo).AddType(TypeNames[typSpec])
	symInput.Lengths = append(symInput.Lengths, arrSize)
	symInput.Package = pkg
	symInput.PreviouslyDeclared = true
	symInput.TotalSize = symInput.Size * TotalLength(symInput.Lengths)

	symExpr := MakeExpression(Natives[OP_IDENTITY], CurrentFile, LineNo)
	symExpr.Package = pkg
	symExpr.Outputs = append(symExpr.Outputs, symOutput)
	symExpr.Inputs = append(symExpr.Inputs, symInput)

	// symOutput.SynonymousTo = symInput.Name

	// marking the output so multidimensional arrays identify the expressions
	symExpr.IsArrayLiteral = true
	result = append(result, symExpr)

	return result
}

func SliceLiteralExpression (typSpec int, exprs []*CXExpression) []*CXExpression {
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

	// slcVar.IsSlice = true

	slcVar.TotalSize = TYPE_POINTER_SIZE
	

	slcVarExpr.Outputs = append(slcVarExpr.Outputs, slcVar)
	slcVar.Package = pkg
	slcVar.PreviouslyDeclared = true

	result = append(result, slcVarExpr)

	var endPointsCounter int
	for _, expr := range exprs {
		if expr.IsArrayLiteral {
			symInp := MakeArgument(symName, CurrentFile, LineNo).AddType(TypeNames[typSpec])
			symInp.Package = pkg
			symOut := MakeArgument(symName, CurrentFile, LineNo).AddType(TypeNames[typSpec])
			symOut.Package = pkg

			// symOut.IsSlice = true
			// symInp.IsSlice = true

			endPointsCounter++

			symExpr := MakeExpression(nil, CurrentFile, LineNo)
			symExpr.Package = pkg
			// symExpr.Outputs = append(symExpr.Outputs, symOut)
			symExpr.AddOutput(symOut)

			if expr.Operator == nil {
				// then it's a literal
				symExpr.Operator = Natives[OP_APPEND]

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

			symInp.TotalSize = TYPE_POINTER_SIZE
			symOut.TotalSize = TYPE_POINTER_SIZE
		} else {
			result = append(result, expr)
		}
	}

	symNameOutput := MakeGenSym(LOCAL_PREFIX)

	symOutput := MakeArgument(symNameOutput, CurrentFile, LineNo).AddType(TypeNames[typSpec])
	// symOutput.PassBy = PASSBY_REFERENCE
	symOutput.IsSlice = true
	symOutput.Package = pkg
	symOutput.PreviouslyDeclared = true

	// symOutput.DeclarationSpecifiers = append(symOutput.DeclarationSpecifiers, DECL_ARRAY)
	

	symInput := MakeArgument(symName, CurrentFile, LineNo).AddType(TypeNames[typSpec])
	// symInput.DereferenceOperations = append(symInput.DereferenceOperations, DEREF_POINTER)
	symInput.IsSlice = true
	symInput.Package = pkg
	// symInput.PassBy = PASSBY_REFERENCE

	symInput.TotalSize = TYPE_POINTER_SIZE
	symOutput.TotalSize = TYPE_POINTER_SIZE

	symExpr := MakeExpression(Natives[OP_IDENTITY], CurrentFile, LineNo)
	symExpr.Package = pkg
	symExpr.Outputs = append(symExpr.Outputs, symOutput)
	symExpr.Inputs = append(symExpr.Inputs, symInput)

	// symExpr.IsArrayLiteral = true

	// symOutput.SynonymousTo = symInput.Name

	// marking the output so multidimensional arrays identify the expressions
	result = append(result, symExpr)

	return result
}

func PrimaryIdentifier (ident string) []*CXExpression {
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

func PrimaryStructLiteral (ident string, strctFlds []*CXExpression) []*CXExpression {
	var result []*CXExpression
	if pkg, err := PRGRM.GetCurrentPackage(); err == nil {
		if strct, err := PRGRM.GetStruct(ident, pkg.Name); err == nil {
			for _, expr := range strctFlds {
				name := expr.Outputs[0].Name

				fld := MakeArgument(name, CurrentFile, LineNo)
				fld.Type = expr.Outputs[0].Type

				expr.IsStructLiteral = true

				expr.Outputs[0].Package = pkg
				// expr.Outputs[0].Program = PRGRM

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

func PrimaryStructLiteralExternal (impName string, ident string, strctFlds []*CXExpression) []*CXExpression {
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
					// expr.Outputs[0].Program = PRGRM

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
	var elt *CXArgument
	if len(prevExprs[len(prevExprs)-1].Outputs[0].Fields) > 0 {
		elt = prevExprs[len(prevExprs)-1].Outputs[0].Fields[len(prevExprs[len(prevExprs)-1].Outputs[0].Fields) - 1]
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

func PostfixExpressionNative (typCode int, opStrCode string) []*CXExpression {
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
		return nil
		// panic(ok)
	}
}

func PostfixExpressionEmptyFunCall (prevExprs []*CXExpression) []*CXExpression {
	if prevExprs[len(prevExprs) - 1].Outputs != nil && len(prevExprs[len(prevExprs) - 1].Outputs[0].Fields) > 0 {
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

func PostfixExpressionFunCall (prevExprs []*CXExpression, args []*CXExpression) []*CXExpression {
	if prevExprs[len(prevExprs) - 1].Outputs != nil && len(prevExprs[len(prevExprs) - 1].Outputs[0].Fields) > 0 {
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

func PostfixExpressionIncDec (prevExprs []*CXExpression, isInc bool) []*CXExpression {
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

func PostfixExpressionField (prevExprs []*CXExpression, ident string) {
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
					} else if _, ok := OpCodes[left.Name+"."+ident]; ok {
						// then it's a native
						// TODO: we'd be referring to the function itself, not a function call
						// (functions as first-class objects)
						left.Name = left.Name + "." + ident
					}
					return
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
					println(ErrorHeader(left.FileName, left.FileLine), fmt.Sprintf("identifier '%s' does not exist", left.Name))
					os.Exit(3)
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

		exprOut.IsReference = false
	case "&":
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

func DeclareLocal (declarator *CXArgument, declaration_specifiers *CXArgument, initializer []*CXExpression, doesInitialize bool) []*CXExpression {
	if doesInitialize {
		declaration_specifiers.IsLocalDeclaration = true

		if pkg, err := PRGRM.GetCurrentPackage(); err == nil {
			if initializer[len(initializer)-1].Operator == nil {
				// then it's a literal, e.g. var foo i32 = 10;
				expr := MakeExpression(Natives[OP_IDENTITY], CurrentFile, LineNo)
				expr.Package = pkg

				declaration_specifiers.Name = declarator.Name
				declaration_specifiers.Package = pkg
				declaration_specifiers.PreviouslyDeclared = true

				expr.AddOutput(declaration_specifiers)
				expr.AddInput(initializer[len(initializer)-1].Outputs[0])

				return []*CXExpression{expr}
			} else {
				// then it's an expression (it has an operator)
				declaration_specifiers.Name = declarator.Name
				declaration_specifiers.Package = pkg
				declaration_specifiers.PreviouslyDeclared = true

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
			expr := MakeExpression(nil, declarator.FileName, declarator.FileLine)
			expr.Package = pkg

			declaration_specifiers.Name = declarator.Name
			declaration_specifiers.Package = pkg
			declaration_specifiers.PreviouslyDeclared = true
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
		name := MakeArgument(MakeGenSym(LOCAL_PREFIX), CurrentFile, LineNo).AddType(TypeNames[leftExprs[len(leftExprs)-1].Inputs[0].Type])
		
		name.Size = leftExprs[len(leftExprs)-1].Operator.Outputs[0].Size
		name.TotalSize = leftExprs[len(leftExprs)-1].Operator.Outputs[0].Size
		name.Type = leftExprs[len(leftExprs)-1].Operator.Outputs[0].Type
		name.Package = pkg
		name.PreviouslyDeclared = true
		

		leftExprs[len(leftExprs)-1].Outputs = append(leftExprs[len(leftExprs)-1].Outputs, name)
	}

	if len(rightExprs[len(rightExprs)-1].Outputs) < 1 {
		name := MakeArgument(MakeGenSym(LOCAL_PREFIX), CurrentFile, LineNo).AddType(TypeNames[rightExprs[len(rightExprs)-1].Inputs[0].Type])

		name.Size = rightExprs[len(rightExprs)-1].Operator.Outputs[0].Size
		name.TotalSize = rightExprs[len(rightExprs)-1].Operator.Outputs[0].Size
		name.Type = rightExprs[len(rightExprs)-1].Operator.Outputs[0].Type
		name.Package = pkg
		name.PreviouslyDeclared = true

		rightExprs[len(rightExprs)-1].Outputs = append(rightExprs[len(rightExprs)-1].Outputs, name)
	}

	expr := MakeExpression(operator, CurrentFile, LineNo)
	// we can't know the type until we compile the full function
	expr.IsUndType = true
	expr.Package = pkg

	if len(leftExprs[len(leftExprs)-1].Outputs[0].Indexes) > 0 || leftExprs[len(leftExprs)-1].Operator != nil {
		// then it's a function call or an array access
		expr.AddInput(leftExprs[len(leftExprs)-1].Outputs[0])
		
		if IsTempVar(leftExprs[len(leftExprs)-1].Outputs[0].Name) {
			out = append(out, leftExprs...)
		} else {
			out = append(out, leftExprs[:len(leftExprs) - 1]...)
		}
	} else {
		expr.Inputs = append(expr.Inputs, leftExprs[len(leftExprs)-1].Outputs[0])
	}
	
	if len(rightExprs[len(rightExprs)-1].Outputs[0].Indexes) > 0 || rightExprs[len(rightExprs)-1].Operator != nil {
		// then it's a function call or an array access
		expr.AddInput(rightExprs[len(rightExprs)-1].Outputs[0])
		
		if IsTempVar(rightExprs[len(rightExprs)-1].Outputs[0].Name) {
			out = append(out, rightExprs...)
		} else {
			out = append(out, rightExprs[:len(rightExprs) - 1]...)
		}
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
		// arg.Program = PRGRM
		
		var size int

		size = len(byts)

		arg.Size = GetArgSize(typ)
		arg.TotalSize = size
		arg.Offset = DataOffset

		if arg.Type == TYPE_STR || arg.Type == TYPE_AFF {
			arg.PassBy = PASSBY_REFERENCE
			arg.Size = TYPE_POINTER_SIZE
			arg.TotalSize = TYPE_POINTER_SIZE
		}

		for i, byt := range byts {
			PRGRM.Memory[DataOffset + i] = byt
		}
		DataOffset += size
		
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
		predicate.PreviouslyDeclared = true
		cond[len(cond)-1].AddOutput(predicate)
		downExpr.AddInput(predicate)
	} else {
		predicate := cond[len(cond)-1].Outputs[0]
		predicate.Package = pkg
		predicate.PreviouslyDeclared = true
		downExpr.AddInput(predicate)
	}

	thenLines := 0
	elseLines := len(incr) + len(statements) + 1

	// processing possible breaks
	for i, stat := range statements {
		if stat.IsBreak {
			stat.ThenLines = elseLines - i - 1
		}
	}

	// processing possible continues
	for i, stat := range statements {
		if stat.IsContinue {
			stat.ThenLines = len(statements) - i - 1
		}
	}

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

		if len(to[0].Outputs[0].Indexes) > 0 {
			f.Outputs[0].Lengths = to[0].Outputs[0].Lengths
			f.Outputs[0].Indexes = to[0].Outputs[0].Indexes
			f.Outputs[0].DereferenceOperations = append(f.Outputs[0].DereferenceOperations, DEREF_ARRAY)
		}

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
		sym.PreviouslyDeclared = true
		from[idx].AddOutput(sym)
		expr.AddInput(sym)
	}

	return append(from, expr)
}

func Assignment (to []*CXExpression, assignOp string, from []*CXExpression) []*CXExpression {
	idx := len(from) - 1

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
				// sym = MakeArgument(to[0].Outputs[0].Name, CurrentFile, LineNo).AddType(TypeNames[from[idx].Operator.Outputs[0].Type])
				
				if from[idx].IsArrayLiteral {
					sym.Size = from[idx].Inputs[0].Size
					sym.TotalSize = from[idx].Inputs[0].TotalSize
					sym.Lengths = from[idx].Inputs[0].Lengths
				}
				if from[idx].Inputs[0].IsSlice {
				// if from[idx].Operator.Outputs[0].IsSlice {
					sym.Lengths = append([]int{0}, sym.Lengths...)
				}
				
				sym.IsSlice = from[idx].Inputs[0].IsSlice
				// sym.IsSlice = from[idx].Operator.Outputs[0].IsSlice
			}
			sym.Package = pkg
			sym.PreviouslyDeclared = true
			sym.IsShortDeclaration = true

			expr.AddOutput(sym)

			for _, toExpr := range to {
				toExpr.Outputs[0].PreviouslyDeclared = true
				toExpr.Outputs[0].IsShortDeclaration = true
			}
			
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
		// to[0].Outputs[0].Program = PRGRM

		if from[idx].IsMethodCall {
			from[idx].Inputs = append(from[idx].Outputs, from[idx].Inputs...)
		} else {
			from[idx].Inputs = from[idx].Outputs
		}

		from[idx].Outputs = to[len(to)-1].Outputs
		// from[idx].Program = PRGRM

		return append(to[:len(to)-1], from...)
	} else {
		if from[idx].Operator.IsNative {
			// only assigning as if the operator had only one output defined

			if from[idx].Operator.OpCode != OP_IDENTITY {
				// it's a short variable declaration
				to[0].Outputs[0].Size = Natives[from[idx].Operator.OpCode].Outputs[0].Size
				to[0].Outputs[0].Type = from[idx].Operator.Outputs[0].Type
				to[0].Outputs[0].Lengths = from[idx].Operator.Outputs[0].Lengths
			}

			// to[0].Outputs[0].Type = from[idx].Operator.Outputs[0].Type
			// to[0].Outputs[0].Lengths = from[idx].Operator.Outputs[0].Lengths
			// to[0].Outputs[0].Size = Natives[from[idx].Operator.OpCode].Outputs[0].Size

			to[0].Outputs[0].DoesEscape = from[idx].Operator.Outputs[0].DoesEscape
			to[0].Outputs[0].PassBy = from[idx].Operator.Outputs[0].PassBy
			// to[0].Outputs[0].Program = PRGRM
		} else {
			// we'll delegate multiple-value returns to the 'expression' grammar rule
			// only assigning as if the operator had only one output defined

			to[0].Outputs[0].Size = from[idx].Operator.Outputs[0].Size
			to[0].Outputs[0].Type = from[idx].Operator.Outputs[0].Type
			to[0].Outputs[0].Lengths = from[idx].Operator.Outputs[0].Lengths
			to[0].Outputs[0].DoesEscape = from[idx].Operator.Outputs[0].DoesEscape
			to[0].Outputs[0].PassBy = from[idx].Operator.Outputs[0].PassBy
			// to[0].Outputs[0].Program = PRGRM
		}

		from[idx].Outputs = to[len(to) - 1].Outputs
		// from[idx].Program = to[len(to) - 1].Program

		return append(to[:len(to)-1], from...)
		// return append(to, from...)
	}
}

func trueJmpExpressions () []*CXExpression {
	pkg, err := PRGRM.GetCurrentPackage()
	if err != nil {
		panic(err)
	}
	
	expr := MakeExpression(Natives[OP_JMP], CurrentFile, LineNo)

	trueArg := WritePrimary(TYPE_BOOL, encoder.Serialize(true), false)
	expr.AddInput(trueArg[0].Outputs[0])
	
	expr.Package = pkg

	return []*CXExpression{expr}
}

func BreakExpressions () []*CXExpression {
	exprs := trueJmpExpressions()
	exprs[0].IsBreak = true
	return exprs
}

func ContinueExpressions () []*CXExpression {
	exprs := trueJmpExpressions()
	exprs[0].IsContinue = true
	return exprs
}

func SelectionExpressions (condExprs []*CXExpression, thenExprs []*CXExpression, elseExprs []*CXExpression) []*CXExpression {
	jmpFn := Natives[OP_JMP]
	pkg, err := PRGRM.GetCurrentPackage()
	if err != nil {
		panic(err)
	}
	ifExpr := MakeExpression(jmpFn, CurrentFile, LineNo)
	ifExpr.Package = pkg

	var predicate *CXArgument
	if condExprs[len(condExprs)-1].Operator == nil && !condExprs[len(condExprs)-1].IsMethodCall {
		// then it's a literal
		predicate = condExprs[len(condExprs)-1].Outputs[0]
	} else {
		// then it's an expression
		predicate = MakeArgument(MakeGenSym(LOCAL_PREFIX), CurrentFile, LineNo)
		if condExprs[len(condExprs)-1].IsMethodCall {
			// we'll change this once we have access to method's types in
			// ProcessMethodCall
			predicate.AddType(TypeNames[TYPE_BOOL])
			condExprs[len(condExprs)-1].Inputs = append(condExprs[len(condExprs)-1].Outputs, condExprs[len(condExprs)-1].Inputs...)
			condExprs[len(condExprs)-1].Outputs = nil
		} else {
			predicate.AddType(TypeNames[condExprs[len(condExprs)-1].Operator.Outputs[0].Type])
		}
		predicate.PreviouslyDeclared = true
		condExprs[len(condExprs)-1].Outputs = append(condExprs[len(condExprs)-1].Outputs, predicate)
	}
	// predicate.Package = pkg

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
	if condExprs[len(condExprs)-1].Operator != nil || condExprs[len(condExprs)-1].IsMethodCall {
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

func PreFinalSize (finalSize *int, sym *CXArgument, arg *CXArgument) {
	for _, op := range sym.DereferenceOperations {
		switch op {
		case DEREF_ARRAY:
			if GetAssignmentElement(sym).IsSlice {
				continue
			}
			var subSize int = 1
			
			for _, len := range GetAssignmentElement(sym).Lengths[:len(GetAssignmentElement(sym).Indexes)] {
				subSize *= len
			}
			*finalSize /= subSize
		// case DEREF_FIELD:
		// 	elt = sym.Fields[fldIdx]
		// 	finalSize = elt.TotalSize
		// 	fldIdx++
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
					// case DECL_SLICE:
					// 	subSize = TYPE_POINTER_SIZE
					case DECL_BASIC:
						subSize = GetArgSize(sym.Type)
					case DECL_STRUCT:
						subSize = arg.CustomType.Size
					}
				}
				
				*finalSize = subSize
			}
		}
	}
}

func SetFinalSize (symbols *map[string]*CXArgument, sym *CXArgument) {
	// var elt *CXArgument
	var finalSize int = sym.TotalSize

	// var fldIdx int
	// elt = sym
	if arg, found := (*symbols)[sym.Package.Name+"."+sym.Name]; found {
		PreFinalSize(&finalSize, sym, arg)
		for _, fld := range sym.Fields {
			finalSize = fld.TotalSize
			PreFinalSize(&finalSize, fld, arg)
		}
	}
	sym.TotalSize = finalSize
}

func GetGlobalSymbol(symbols *map[string]*CXArgument, symPackage *CXPackage, symName string) {
	if _, found := (*symbols)[symPackage.Name + "." + symName]; !found {
		if glbl, err := symPackage.GetGlobal(symName); err == nil {
			(*symbols)[symPackage.Name + "." + symName] = glbl
		}
	}
}

func ProcessDereferenceLevels () {
	// handles the levels of dereferencing
	// the symbol in the end adopts the fields of its pointee
	// might not be needed
	// if sym.DereferenceLevels > 0 {
	// 	if arg.IndirectionLevels >= sym.DereferenceLevels || isFieldPointer {
	// 		pointer := arg

	// 		for c := 0; c < sym.DereferenceLevels-1; c++ {
	// 			pointer = pointer.Pointee
	// 		}

	// 		sym.Offset = pointer.Offset
	// 		sym.IndirectionLevels = pointer.IndirectionLevels
	// 		sym.IsPointer = pointer.IsPointer
	// 	} else {
	// 		panic("invalid indirect of " + sym.Name)
	// 	}
	// } else {
	// 	sym.Offset = arg.Offset
	// 	sym.IsPointer = arg.IsPointer
	// 	sym.IndirectionLevels = arg.IndirectionLevels
	// }
}

func ProcessSymbolFields (sym *CXArgument, arg *CXArgument) {
	if len(sym.Fields) > 0 {
		if arg.CustomType == nil || len(arg.CustomType.Fields) == 0 {
			println(ErrorHeader(sym.FileName, sym.FileLine), fmt.Sprintf("'%s' has no fields", sym.Name))
			os.Exit(3)
		}
		
		// checking if fields do exist in their CustomType
		// and assigning that CustomType to the sym.Field
		strct := arg.CustomType

		// methodName := arg.CustomType.Fields[len(arg.CustomType.Fields) - 1].Name
		
		for _, fld := range sym.Fields {
			if inFld, err := strct.GetField(fld.Name); err == nil {
				if inFld.CustomType != nil {
					fld.CustomType = strct
					strct = inFld.CustomType
				}
			} else {
				methodName := sym.Fields[len(sym.Fields) - 1].Name
				receiverType := strct.Name

				if method, methodErr := strct.Package.GetMethod(receiverType + "." + methodName, receiverType); methodErr == nil {
					fld.Type = method.Outputs[0].Type
				} else {
					println(ErrorHeader(fld.FileName, fld.FileLine), err.Error())
				}
				
				
			}
		}

		strct = arg.CustomType
		// then we copy all the type struct fields
		// to the respective sym.Fields
		for _, nameFld := range sym.Fields {
			if nameFld.CustomType != nil {
				strct = nameFld.CustomType
			}
			
			for _, fld := range strct.Fields {
				if nameFld.Name == fld.Name {
					nameFld.Type = fld.Type
					nameFld.Lengths = fld.Lengths
					nameFld.Size = fld.Size
					nameFld.TotalSize = fld.TotalSize
					nameFld.DereferenceLevels = sym.DereferenceLevels
					nameFld.IsPointer = fld.IsPointer
					nameFld.CustomType = fld.CustomType
					
					// sym.DereferenceOperations = append(sym.DereferenceOperations, DEREF_FIELD)
					
					if fld.IsSlice {
						nameFld.DereferenceOperations = append([]int{DEREF_POINTER}, nameFld.DereferenceOperations...)
						nameFld.DereferenceLevels++
					}

					nameFld.PassBy = fld.PassBy
					nameFld.IsSlice = fld.IsSlice
					
					if fld.Type == TYPE_STR || fld.Type == TYPE_AFF {
						nameFld.PassBy = PASSBY_REFERENCE
						// nameFld.Size = TYPE_POINTER_SIZE
						// nameFld.TotalSize = TYPE_POINTER_SIZE
					}
					
					if fld.CustomType != nil {
						strct = fld.CustomType
					}
					break
				}

				nameFld.Offset += fld.TotalSize
			}
		}
	}
}

func CopyArgFields (sym *CXArgument, arg *CXArgument) {
	sym.Offset = arg.Offset
	sym.IsPointer = arg.IsPointer
	sym.IndirectionLevels = arg.IndirectionLevels

	sym.IsSlice = arg.IsSlice
	sym.CustomType = arg.CustomType

	sym.Lengths = arg.Lengths
	sym.Package = arg.Package
	// sym.Program = arg.Program
	sym.DoesEscape = arg.DoesEscape
	sym.Size = arg.Size

	if arg.Type == TYPE_STR {
		sym.IsPointer = true
	}

	if arg.IsSlice {
		sym.DereferenceOperations = append([]int{DEREF_POINTER}, sym.DereferenceOperations...)
		sym.DereferenceLevels++
	}

	// assignElt := GetAssignmentElement(sym)
	
	// if assignElt.IsSlice && len(assignElt.Indexes) > 0 {
	// 	sym.DereferenceOperations = append([]int{DEREF_POINTER}, sym.DereferenceOperations...)
	// 	sym.DereferenceLevels++
	// }
	
	// if (GetAssignmentElement(sym).IsSlice || (arg.IsSlice && len(sym.Fields) == 0)) &&
	// 	(len(GetAssignmentElement(sym).Fields) > 0 || len(GetAssignmentElement(sym).Indexes) > 0) {
	// 	sym.DereferenceOperations = append([]int{DEREF_POINTER}, sym.DereferenceOperations...)
	// 	sym.DereferenceLevels++
	// }

	if len(sym.Fields) > 0 {
		sym.Type = sym.Fields[len(sym.Fields) - 1].Type
		// sym.IsSlice = sym.Fields[len(sym.Fields) - 1].IsSlice
	} else {
		sym.Type = arg.Type
	}

	if sym.IsReference && !arg.IsStruct {
		sym.TotalSize = arg.TotalSize
	} else {
		if len(sym.Fields) > 0 {
			sym.TotalSize = sym.Fields[len(sym.Fields)-1].TotalSize
		} else {
			sym.TotalSize = arg.TotalSize
		}
	}
}

func GiveOffset (symbols *map[string]*CXArgument, sym *CXArgument, offset *int, shouldExist bool) {
	if sym.Name != "" {
		if !sym.IsLocalDeclaration {
			GetGlobalSymbol(symbols, sym.Package, sym.Name)
		}

		if arg, found := (*symbols)[sym.Package.Name+"."+sym.Name]; found {
			// ProcessDereferenceLevels()
			ProcessSymbolFields(sym, arg)
			CopyArgFields(sym, arg)
		}
	}
}

func ProcessTempVariable (expr *CXExpression) {
	if expr.Operator != nil && (expr.Operator == Natives[OP_IDENTITY] || IsUndOp(expr.Operator)) && len(expr.Outputs) > 0 && len(expr.Inputs) > 0 {
		name := expr.Outputs[0].Name
		arg := expr.Outputs[0]
		if IsTempVar(name) {
			// then it's a temporary variable and it needs to adopt its input's type
			arg.Type = expr.Inputs[0].Type
			arg.Size = expr.Inputs[0].Size
			arg.TotalSize = expr.Inputs[0].TotalSize
			arg.PreviouslyDeclared = true
		}
	}
}

// func ProcessShortDeclaration (expr *CXExpression) {
// 	// if len(expr.Outputs) > 0 && len(expr.Inputs) > 0 && expr.Outputs[0].PreviouslyDeclared && (expr.Operator == nil || expr.Operator.OpCode == OP_IDENTITY) {
// 	if len(expr.Inputs) > 0 && len(expr.Outputs) > 0 && expr.Outputs[0].PreviouslyDeclared {
// 		expr.Outputs[0].Type = expr.Inputs[0].Type
// 		expr.Outputs[0].Size = expr.Inputs[0].Size
// 		expr.Outputs[0].TotalSize = expr.Inputs[0].TotalSize
		
// 		expr.Outputs[0].Lengths = expr.Inputs[0].Lengths
// 		expr.Outputs[0].Fields = expr.Inputs[0].Fields
		
// 		// if expr.Operator != nil && expr.Operator.OpCode == OP_IDENTITY {
// 		// 	expr.Outputs[0].Type = expr.Inputs[0].Type
// 		// 	expr.Outputs[0].Size = expr.Inputs[0].Size
// 		// 	expr.Outputs[0].TotalSize = expr.Inputs[0].TotalSize

// 		// 	expr.Outputs[0].Lengths = expr.Inputs[0].Lengths
// 		// 	expr.Outputs[0].Fields = expr.Inputs[0].Fields
// 		// }
// 		//  else {
// 		// 	expr.Outputs[0].Type = expr.Inputs[0].Type
// 		// 	expr.Outputs[0].Size = expr.Inputs[0].Size
// 		// 	expr.Outputs[0].TotalSize = expr.Inputs[0].TotalSize

// 		// 	expr.Outputs[0].Lengths = expr.Inputs[0].Lengths
// 		// 	expr.Outputs[0].Fields = expr.Inputs[0].Fields
// 		// }
// 	}
// }

func ProcessMethodCall (expr *CXExpression, symbols *map[string]*CXArgument, offset *int, shouldExist bool) {
	if expr.IsMethodCall {
		var inp *CXArgument
		var out *CXArgument

		if len(expr.Inputs) > 0 && expr.Inputs[0].Name != "" {
			inp = expr.Inputs[0]
		}
		if len(expr.Outputs) > 0 && expr.Outputs[0].Name != "" {
			out = expr.Outputs[0]
		}
		
		if inp != nil {
			if argInp, found := (*symbols)[inp.Package.Name+"."+inp.Name]; !found {
				if out != nil {
					if argOut, found := (*symbols)[out.Package.Name+"."+out.Name]; !found {
						panic("")
					} else {
						// then we found an output
						if len(out.Fields) > 0 {
							strct := argOut.CustomType

							if fn, err := strct.Package.GetMethod(strct.Name + "." + out.Fields[len(out.Fields) - 1].Name, strct.Name); err == nil {
								expr.Operator = fn
							} else {
								panic("")
							}

							expr.Inputs = append([]*CXArgument{out}, expr.Inputs...)

							out.Fields = out.Fields[:len(out.Fields) - 1]
							
							expr.Outputs = expr.Outputs[1:]
						}
					}
				} else {
					panic("")
				}
			} else {
				// then we found an input

				if len(inp.Fields) > 0 {
					strct := argInp.CustomType

					for _, fld := range inp.Fields {
						if inFld, err := strct.GetField(fld.Name); err == nil {
							if inFld.CustomType != nil {
								strct = inFld.CustomType
							}
						}
					}

					if fn, err := strct.Package.GetMethod(strct.Name + "." + inp.Fields[len(inp.Fields) - 1].Name, strct.Name); err == nil {
						expr.Operator = fn
					} else {
						panic(err)
					}
					
					inp.Fields = inp.Fields[:len(inp.Fields) - 1]
				} else if len(out.Fields) > 0 {
					if argOut, found := (*symbols)[out.Package.Name + "." + out.Name]; found {
						strct := argOut.CustomType

						expr.Inputs = append(expr.Outputs[:1], expr.Inputs...)
						
						expr.Outputs = expr.Outputs[:len(expr.Outputs) - 1]
						
						if fn, err := strct.Package.GetMethod(strct.Name + "." + out.Fields[len(out.Fields) - 1].Name, strct.Name); err == nil {
							expr.Operator = fn
						} else {
							panic(err)
						}
						
						out.Fields = out.Fields[:len(out.Fields) - 1]
					} else {
						panic("")
					}
				}
			}
		} else {
			if out != nil {
				if argOut, found := (*symbols)[out.Package.Name+"."+out.Name]; !found {
					panic("")
				} else {
					// then we found an output
					if len(out.Fields) > 0 {
						strct := argOut.CustomType

						if fn, err := strct.Package.GetMethod(strct.Name + "." + out.Fields[len(out.Fields) - 1].Name, strct.Name); err == nil {
							expr.Operator = fn
						} else {
							panic("")
						}

						expr.Inputs = append([]*CXArgument{out}, expr.Inputs...)

						
						out.Fields = out.Fields[:len(out.Fields) - 1]
						
						expr.Outputs = expr.Outputs[1:]
						// expr.Outputs = nil
					}
				}
			} else {
				panic("")
			}
		}

		// checking if receiver is sent as pointer or not
		if expr.Operator.Inputs[0].IsPointer {
			expr.Inputs[0].PassBy = PASSBY_REFERENCE
		}
	}
}

func UpdateSymbolsTable(symbols *map[string]*CXArgument, sym *CXArgument, offset *int, shouldExist bool) {
	if sym.Name != "" {
		if !sym.IsLocalDeclaration {
			GetGlobalSymbol(symbols, sym.Package, sym.Name)
		}
		
		if _, found := (*symbols)[sym.Package.Name+"."+sym.Name]; !found {
			if shouldExist {
				// it should exist. error
				println(ErrorHeader(sym.FileName, sym.FileLine) + " identifier '" + sym.Name + "' does not exist")
				os.Exit(3)
			}

			sym.Offset = *offset
			(*symbols)[sym.Package.Name+"."+sym.Name] = sym

			if sym.IsSlice {
				*offset += sym.Size
			} else {
				*offset += sym.TotalSize
			}
		}
	}
}

func GiveCustomType(symbols *map[string]*CXArgument, sym *CXArgument, offset *int, shouldExist bool) {
	if sym.Name != "" {
		
	}
}

func ProcessSlice (inp *CXArgument) {
	var elt *CXArgument

	if len(inp.Fields) > 0 {
		elt = inp.Fields[len(inp.Fields) - 1]
	} else {
		elt = inp
	}

	// elt.IsPointer = true

	if elt.IsSlice && len(elt.DereferenceOperations) > 0 && elt.DereferenceOperations[len(elt.DereferenceOperations) - 1] == DEREF_POINTER {
		elt.DereferenceOperations = elt.DereferenceOperations[:len(elt.DereferenceOperations) - 1]
	} else if elt.IsSlice && len(elt.DereferenceOperations) > 0 && len(inp.Fields) == 0 {
		// elt.DereferenceOperations = append([]int{DEREF_POINTER}, elt.DereferenceOperations...)
	}
}

func ProcessSliceAssignment (expr *CXExpression) {
	if expr.Operator == Natives[OP_IDENTITY] {
		var inp *CXArgument
		var out *CXArgument

		inp = GetAssignmentElement(expr.Inputs[0])
		out = GetAssignmentElement(expr.Outputs[0])

		if inp.IsSlice && out.IsSlice && len(inp.Indexes) == 0 && len(out.Indexes) == 0 {
			out.PassBy = PASSBY_VALUE
		}
		
		// if len(expr.Inputs[0].Fields) > 0 {
		// 	inp = expr.Inputs[0].Fields[len(expr.Inputs[0].Fields) - 1]
		// } else {
		// 	inp = expr.Inputs[0]
		// }
		// if len(expr.Outputs[0].Fields) > 0 {
		// 	out = expr.Outputs[0].Fields[len(expr.Outputs[0].Fields) - 1]
		// } else {
		// 	out = expr.Outputs[0]
		// }
		
		// if out.IsSlice && inp.IsSlice {
		// 	inp.DereferenceOperations = append([]int{DEREF_POINTER}, inp.DereferenceOperations...)
		// 	out.PassBy = PASSBY_REFERENCE
		// }
	}
	if expr.Operator != nil && !expr.Operator.IsNative {
		// then it's a function call
		for _, inp := range expr.Inputs {
			assignElt := GetAssignmentElement(inp)
			
			if assignElt.IsSlice && len(assignElt.Indexes) == 0 {
				assignElt.PassBy = PASSBY_VALUE
			}
		}
	}
}

func ProcessStringAssignment (expr *CXExpression) {
	if expr.Operator == Natives[OP_IDENTITY] {
		for i, out := range expr.Outputs {
			if len(expr.Inputs) > i {
				out = GetAssignmentElement(out)
				inp := GetAssignmentElement(expr.Inputs[i])

				if (out.Type == TYPE_STR || out.Type == TYPE_AFF) && out.Name != "" &&
					(inp.Type == TYPE_STR || inp.Type == TYPE_AFF) && inp.Name != "" {
					out.PassBy = PASSBY_VALUE
				}
			}
		}
	}
}

func CheckTypes(expr *CXExpression) {
	if expr.Operator != nil {
		var opName string
		if expr.Operator.IsNative {
			opName = OpNames[expr.Operator.OpCode]
		} else {
			opName = expr.Operator.Name
		}

		// checking if number of inputs is less than the required number of inputs
		if len(expr.Inputs) != len(expr.Operator.Inputs) {
			if !(len(expr.Operator.Inputs) > 0 && expr.Operator.Inputs[len(expr.Operator.Inputs) - 1].Type != TYPE_UNDEFINED) {
				// if the last input is of type TYPE_UNDEFINED then it might be a variadic function, such as printf
			} else {
				// then we need to be strict in the number of inputs
				var plural1 string
				var plural2 string = "s"
				var plural3 string = "were"
				if len(expr.Operator.Inputs) > 1 {
					plural1 = "s"
				}
				if len(expr.Inputs) == 1 {
					plural2 = ""
					plural3 = "was"
				}

				println(ErrorHeader(expr.FileName, expr.FileLine), fmt.Sprintf("operator '%s' expects %d input%s, but %d input argument%s %s provided", opName, len(expr.Operator.Inputs), plural1, len(expr.Inputs), plural2, plural3))
				os.Exit(3)
			}
		}

		// checking if number of expr.Outputs match number of Operator.Outputs
		if len(expr.Outputs) != len(expr.Operator.Outputs) {
			var plural1 string
			var plural2 string = "s"
			var plural3 string = "were"
			if len(expr.Operator.Outputs) > 1 {
				plural1 = "s"
			}
			if len(expr.Outputs) == 1 {
				plural2 = ""
				plural3 = "was"
			}
			println(ErrorHeader(expr.FileName, expr.FileLine), fmt.Sprintf("operator '%s' expects to return %d output%s, but %d receiving argument%s %s provided", opName, len(expr.Operator.Outputs), plural1, len(expr.Outputs), plural2, plural3)) 
		}
	}

	if expr.Operator != nil && expr.Operator.IsNative && expr.Operator.OpCode == OP_IDENTITY {
		for i, _ := range expr.Inputs {
			var expectedType string
			var receivedType string
			if GetAssignmentElement(expr.Outputs[i]).CustomType != nil {
				// then it's custom type
				expectedType = GetAssignmentElement(expr.Outputs[i]).CustomType.Name
			} else {
				// then it's native type
				expectedType = TypeNames[GetAssignmentElement(expr.Outputs[i]).Type]
			}
			

			if GetAssignmentElement(expr.Inputs[i]).CustomType != nil {
				// then it's custom type
				receivedType = GetAssignmentElement(expr.Inputs[i]).CustomType.Name
			} else {
				// then it's native type
				receivedType = TypeNames[GetAssignmentElement(expr.Inputs[i]).Type]
			}

			// if GetAssignmentElement(expr.Outputs[i]).Type != GetAssignmentElement(inp).Type {
			if receivedType != expectedType {
				if expr.IsStructLiteral {
					println(ErrorHeader(expr.Outputs[i].FileName, expr.Outputs[i].FileLine), fmt.Sprintf("field '%s' in struct literal of type '%s' expected argument of type '%s'; '%s' was provided", expr.Outputs[i].Fields[0].Name, expr.Outputs[i].CustomType.Name, expectedType, receivedType))
				} else {
					println(ErrorHeader(expr.Outputs[i].FileName, expr.Outputs[i].FileLine), fmt.Sprintf("trying to assign argument of type '%s' to symbol '%s' of type '%s'", receivedType, GetAssignmentElement(expr.Outputs[i]).Name, expectedType))
				}
			}
		}
	}

	// checking inputs matching operator's inputs
	if expr.Operator != nil {
		// then it's a function call and not a declaration
		for i, inp := range expr.Operator.Inputs {

			var expectedType string
			var receivedType string
			if expr.Operator.Inputs[i].CustomType != nil {
				// then it's custom type
				expectedType = expr.Operator.Inputs[i].CustomType.Name
			} else {
				// then it's native type
				expectedType = TypeNames[expr.Operator.Inputs[i].Type]
			}

			if GetAssignmentElement(expr.Inputs[i]).CustomType != nil {
				// then it's custom type
				receivedType = GetAssignmentElement(expr.Inputs[i]).CustomType.Name
			} else {
				// then it's native type
				receivedType = TypeNames[GetAssignmentElement(expr.Inputs[i]).Type]
			}
			
			// if inp.Type != expr.Inputs[i].Type && inp.Type != TYPE_UNDEFINED {
			if expectedType != receivedType && inp.Type != TYPE_UNDEFINED {
				var opName string
				if expr.Operator.IsNative {
					opName = OpNames[expr.Operator.OpCode]
				} else {
					opName = expr.Operator.Name
				}

				println(ErrorHeader(expr.Inputs[i].FileName, expr.Inputs[i].FileLine), fmt.Sprintf("function '%s' expected input argument of type '%s'; '%s' was provided", opName, expectedType, receivedType))
			}
		}
	}
}

func FunctionAddParameters (fn *CXFunction, inputs, outputs []*CXArgument) {
	if len(fn.Inputs) != len(inputs) {
		// it must be a method declaration
		// so we save the first input
		fn.Inputs = fn.Inputs[:1]
	} else {
		fn.Inputs = nil
	}

	// we need to wipe the inputs recognized in the first pass
	// as these don't have all the fields correctly
	fn.Outputs = nil
	
	for _, inp := range inputs {
		fn.AddInput(inp)
	}
	
	for _, out := range outputs {
		fn.AddOutput(out)
	}

	for _, out := range fn.Outputs {
		if out.IsPointer && out.Type != TYPE_STR && out.Type != TYPE_AFF {
			out.DoesEscape = true
		}
	}
}

func ProcessGoTos (fn *CXFunction, exprs []*CXExpression) {
	for i, expr := range exprs {
		if expr.Label != "" && expr.Operator == Natives[OP_JMP] {
			// then it's a goto
			for j, e := range exprs {
				if e.Label == expr.Label && i != j {
					// ElseLines is used because arg's default val is false
					expr.ThenLines = j - i - 1
					break
				}
			}
		}

		fn.AddExpression(expr)
	}
}

func FunctionProcessParameters (symbols *map[string]*CXArgument, symbolsScope *map[string]bool, offset *int, fn *CXFunction, params []*CXArgument) {
	for _, param := range params {
		ProcessLocalDeclaration(symbols, symbolsScope, param)

		UpdateSymbolsTable(symbols, param, offset, false)
		GiveOffset(symbols, param, offset, false)
		SetFinalSize(symbols, param)

		AddPointer(fn, param)

		// as these are declarations, they should not have any dereference operations
		param.DereferenceOperations = nil
	}
}

func ProcessLocalDeclaration (symbols *map[string]*CXArgument, symbolsScope *map[string]bool, arg *CXArgument) {
	if arg.IsLocalDeclaration {
		(*symbolsScope)[arg.Package.Name+"."+arg.Name] = true
	}
	arg.IsLocalDeclaration = (*symbolsScope)[arg.Package.Name+"."+arg.Name]
}

func CheckRedeclared (symbols *map[string]*CXArgument, expr *CXExpression, sym *CXArgument) {
	if expr.Operator == nil && len(expr.Outputs) > 0 && len(expr.Inputs) == 0 {
		if _, found := (*symbols)[sym.Package.Name+"."+sym.Name]; found {
			println(ErrorHeader(sym.FileName, sym.FileLine), fmt.Sprintf("'%s' redeclared", sym.Name))
		}
	}
}

func ProcessExpressionArguments (symbols *map[string]*CXArgument, symbolsScope *map[string]bool, offset *int, fn *CXFunction, args []*CXArgument, expr *CXExpression, isInput bool) {
	for _, arg := range args {		
		ProcessLocalDeclaration(symbols, symbolsScope, arg)

		if !isInput {
			CheckRedeclared(symbols, expr, arg)
		}
		
		if !isInput {
			ProcessUndExpression(expr)
		}

		if arg.PreviouslyDeclared {
			UpdateSymbolsTable(symbols, arg, offset, false)
		} else {
			UpdateSymbolsTable(symbols, arg, offset, true)
		}

		if isInput {
			GiveOffset(symbols, arg, offset, true)
		} else {
			GiveOffset(symbols, arg, offset, false)
		}

		ProcessSlice(arg)
		
		for _, idx := range arg.Indexes {
			UpdateSymbolsTable(symbols, idx, offset, true)
			GiveOffset(symbols, idx, offset, true)
		}
		for _, fld := range arg.Fields {
			for _, idx := range fld.Indexes {
				UpdateSymbolsTable(symbols, idx, offset, true)
				GiveOffset(symbols, idx, offset, true)
			}
		}

		SetFinalSize(symbols, arg)

		AddPointer(fn, arg)
	}
}

// func ProcessMethodCalls (exprs []*CXExpression, symbols *map[string]*CXArgument) []*CXExpression {
// 	for _, expr := range exprs {
// 		if expr.IsMethodCall {
// 			out := MakeArgument(MakeGenSym(LOCAL_PREFIX), expr.FileName, expr.FileLine)
// 			out.PreviouslyDeclared = true
			
// 			newExpr := MakeExpression(Natives[OP_IDENTITY], expr.FileName, expr.FileLine)
// 			newExpr.AddOutput(out)
// 			newExpr.AddInput(exprs[len(exprs) - 1].Outputs[0])

// 			out.Package = expr.Package
// 			expr.Outputs[0].Package = expr.Package

// 			// out.Fields = expr.Outputs[0].Fields[len(expr.Outputs[0].Fields) - 1 : ]
// 			// expr.Outputs[0].Fields = expr.Outputs[0].Fields[:len(expr.Outputs[0].Fields) - 1]
// 			// expr.Inputs = []*CXArgument{out}

// 			// nestedExprs = append(nestedExprs, newExpr)
// 		}
// 	}

// 	return exprs
// }

func ProcessPointerStructs (expr *CXExpression) {
	for _, arg := range append(expr.Inputs, expr.Outputs...) {
		if arg.IsStruct && arg.IsPointer && len(arg.Fields) > 0 && arg.DereferenceLevels == 0 {
			arg.DereferenceLevels++
			arg.DereferenceOperations = append(arg.DereferenceOperations, DEREF_POINTER)
		}
	}
}

// Depending on the operator, we're going to return the input's size or a prefixed size (like a Boolean)
func undOutputSize (expr *CXExpression) int {
	switch expr.Operator.OpCode {
	case OP_UND_EQUAL, OP_UND_UNEQUAL, OP_UND_LT, OP_UND_GT, OP_UND_LTEQ, OP_UND_GTEQ:
		// the result is a Boolean for any of these
		return 1
	default:
		return GetAssignmentElement(expr.Inputs[0]).Size
	}
}

func ProcessUndExpression (expr *CXExpression) {
	if expr.IsUndType {
		for _, out := range expr.Outputs {
			out.Size = undOutputSize(expr)
			out.TotalSize = out.Size
		}
	}
}

func FunctionDeclaration (fn *CXFunction, inputs, outputs []*CXArgument, exprs []*CXExpression) {
	if FoundCompileErrors {
		os.Exit(3)
	}

	FunctionAddParameters(fn, inputs, outputs)

	// getting offset to use by statements (excluding inputs, outputs and receiver)
	var offset int
	PRGRM.HeapStartsAt = DataOffset

	ProcessGoTos(fn, exprs)

	fn.Length = len(fn.Expressions)

	var symbols map[string]*CXArgument = make(map[string]*CXArgument, 0)
	var symbolsScope map[string]bool = make(map[string]bool, 0)

	FunctionProcessParameters(&symbols, &symbolsScope, &offset, fn, fn.Inputs)
	FunctionProcessParameters(&symbols, &symbolsScope, &offset, fn, fn.Outputs)

	// fn.Expressions = ProcessMethodCalls(fn.Expressions, &symbols)

	for i, expr := range fn.Expressions {
		// ProcessShortDeclaration(expr)
		
		ProcessMethodCall(expr, &symbols, &offset, true)
		ProcessExpressionArguments(&symbols, &symbolsScope, &offset, fn, expr.Inputs, expr, true)
		ProcessExpressionArguments(&symbols, &symbolsScope, &offset, fn, expr.Outputs, expr, false)

		ProcessPointerStructs(expr)
		
		SetCorrectArithmeticOp(expr)
		ProcessTempVariable(expr)
		ProcessSliceAssignment(expr)
		ProcessStringAssignment(expr)

		// process short declaration
		if len(expr.Outputs) > 0 && len(expr.Inputs) > 0 && expr.Outputs[0].IsShortDeclaration {
			fn.Expressions[i - 1].Outputs[0].Type = fn.Expressions[i].Inputs[0].Type
			fn.Expressions[i].Outputs[0].Type = fn.Expressions[i].Inputs[0].Type
		}

		CheckTypes(expr)
	}

	fn.Size = offset
}

func FunctionCall (exprs []*CXExpression, args []*CXExpression) []*CXExpression {
	expr := exprs[len(exprs)-1]
	
	if expr.Operator == nil {
		opName := expr.Outputs[0].Name
		opPkg := expr.Outputs[0].Package

		if op, err := PRGRM.GetFunction(opName, opPkg.Name); err == nil {
			expr.Operator = op
		} else if expr.Outputs[0].Fields == nil {
			// then it's not a possible method call
			println(ErrorHeader(CurrentFile, LineNo), err.Error())
			os.Exit(3)
			return nil
		} else {
			expr.IsMethodCall = true
		}

		if len(expr.Outputs) > 0 && expr.Outputs[0].Fields == nil {
			expr.Outputs = nil
		}
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
					out = MakeArgument(MakeGenSym(LOCAL_PREFIX), CurrentFile, inpExpr.FileLine).AddType(TypeNames[inpExpr.Inputs[0].Type])
					out.CustomType = inpExpr.Inputs[0].CustomType

					out.Size = inpExpr.Inputs[0].Size
					out.TotalSize = inpExpr.Inputs[0].Size
					
					out.Type = inpExpr.Inputs[0].Type
					out.PreviouslyDeclared = true
				} else {
					out = MakeArgument(MakeGenSym(LOCAL_PREFIX), CurrentFile, inpExpr.FileLine).AddType(TypeNames[inpExpr.Operator.Outputs[0].Type])
					

					out.CustomType = inpExpr.Operator.Outputs[0].CustomType
					

					if inpExpr.Operator.Outputs[0].CustomType != nil {
						if strct, err := inpExpr.Package.GetStruct(inpExpr.Operator.Outputs[0].CustomType.Name); err == nil {
							out.Size = strct.Size
							out.TotalSize = strct.Size
						}
					} else {
						out.Size = inpExpr.Operator.Outputs[0].Size
						out.TotalSize = inpExpr.Operator.Outputs[0].Size
					}

					out.Type = inpExpr.Operator.Outputs[0].Type
					out.PreviouslyDeclared = true
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
