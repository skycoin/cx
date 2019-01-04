package actions

import (
	"fmt"
	. "github.com/skycoin/cx/cx"
)

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

func FunctionAddParameters(fn *CXFunction, inputs, outputs []*CXArgument) {
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

func FunctionDeclaration(fn *CXFunction, inputs, outputs []*CXArgument, exprs []*CXExpression) {
	if FoundCompileErrors {
		return
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
		if len(expr.Outputs) > 0 && len(expr.Inputs) > 0 && expr.Outputs[0].IsShortDeclaration && !expr.IsStructLiteral {
			fn.Expressions[i-1].Outputs[0].Type = fn.Expressions[i].Inputs[0].Type
			fn.Expressions[i].Outputs[0].Type = fn.Expressions[i].Inputs[0].Type
		}

		CheckTypes(expr)
	}

	fn.Size = offset
}

func FunctionCall(exprs []*CXExpression, args []*CXExpression) []*CXExpression {
	expr := exprs[len(exprs)-1]

	if expr.Operator == nil {
		opName := expr.Outputs[0].Name
		opPkg := expr.Outputs[0].Package

		if op, err := PRGRM.GetFunction(opName, opPkg.Name); err == nil {
			expr.Operator = op
		} else if expr.Outputs[0].Fields == nil {
			// then it's not a possible method call
			println(CompilationError(CurrentFile, LineNo), err.Error())
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

// Depending on the operator, we're going to return the input's size or a prefixed size (like a Boolean)
func undOutputSize(expr *CXExpression) int {
	switch expr.Operator.OpCode {
	case OP_UND_EQUAL, OP_UND_UNEQUAL, OP_UND_LT, OP_UND_GT, OP_UND_LTEQ, OP_UND_GTEQ:
		// the result is a Boolean for any of these
		return 1
	default:
		return GetAssignmentElement(expr.Inputs[0]).Size
	}
}

func ProcessUndExpression(expr *CXExpression) {
	if expr.IsUndType {
		for _, out := range expr.Outputs {
			out.Size = undOutputSize(expr)
			out.TotalSize = out.Size
		}
	}
}

func ProcessPointerStructs(expr *CXExpression) {
	for _, arg := range append(expr.Inputs, expr.Outputs...) {
		if arg.IsStruct && arg.IsPointer && len(arg.Fields) > 0 && arg.DereferenceLevels == 0 {
			arg.DereferenceLevels++
			arg.DereferenceOperations = append(arg.DereferenceOperations, DEREF_POINTER)
		}
	}
}

func ProcessExpressionArguments(symbols *map[string]*CXArgument, symbolsScope *map[string]bool, offset *int, fn *CXFunction, args []*CXArgument, expr *CXExpression, isInput bool) {
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

func CheckRedeclared(symbols *map[string]*CXArgument, expr *CXExpression, sym *CXArgument) {
	if expr.Operator == nil && len(expr.Outputs) > 0 && len(expr.Inputs) == 0 {
		if _, found := (*symbols)[sym.Package.Name+"."+sym.Name]; found {
			println(CompilationError(sym.FileName, sym.FileLine), fmt.Sprintf("'%s' redeclared", sym.Name))
		}
	}
}

func ProcessLocalDeclaration(symbols *map[string]*CXArgument, symbolsScope *map[string]bool, arg *CXArgument) {
	if arg.IsLocalDeclaration {
		(*symbolsScope)[arg.Package.Name+"."+arg.Name] = true
	}
	arg.IsLocalDeclaration = (*symbolsScope)[arg.Package.Name+"."+arg.Name]
}

func FunctionProcessParameters(symbols *map[string]*CXArgument, symbolsScope *map[string]bool, offset *int, fn *CXFunction, params []*CXArgument) {
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

func ProcessGoTos(fn *CXFunction, exprs []*CXExpression) {
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

func CheckTypes(expr *CXExpression) {
	if expr.Operator != nil {
		opName := ExprOpName(expr)

		// checking if number of inputs is less than the required number of inputs
		if len(expr.Inputs) != len(expr.Operator.Inputs) {
			if !(len(expr.Operator.Inputs) > 0 && expr.Operator.Inputs[len(expr.Operator.Inputs)-1].Type != TYPE_UNDEFINED) {
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

				println(CompilationError(expr.FileName, expr.FileLine), fmt.Sprintf("operator '%s' expects %d input%s, but %d input argument%s %s provided", opName, len(expr.Operator.Inputs), plural1, len(expr.Inputs), plural2, plural3))
				return
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
			println(CompilationError(expr.FileName, expr.FileLine), fmt.Sprintf("operator '%s' expects to return %d output%s, but %d receiving argument%s %s provided", opName, len(expr.Operator.Outputs), plural1, len(expr.Outputs), plural2, plural3))
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
					println(CompilationError(expr.Outputs[i].FileName, expr.Outputs[i].FileLine), fmt.Sprintf("field '%s' in struct literal of type '%s' expected argument of type '%s'; '%s' was provided", expr.Outputs[i].Fields[0].Name, expr.Outputs[i].CustomType.Name, expectedType, receivedType))
				} else {
					println(CompilationError(expr.Outputs[i].FileName, expr.Outputs[i].FileLine), fmt.Sprintf("trying to assign argument of type '%s' to symbol '%s' of type '%s'", receivedType, GetAssignmentElement(expr.Outputs[i]).Name, expectedType))
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

				println(CompilationError(expr.Inputs[i].FileName, expr.Inputs[i].FileLine), fmt.Sprintf("function '%s' expected input argument of type '%s'; '%s' was provided", opName, expectedType, receivedType))
			}
		}
	}
}

func ProcessStringAssignment(expr *CXExpression) {
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

func ProcessSlice(inp *CXArgument) {
	var elt *CXArgument

	if len(inp.Fields) > 0 {
		elt = inp.Fields[len(inp.Fields)-1]
	} else {
		elt = inp
	}

	// elt.IsPointer = true

	if elt.IsSlice && len(elt.DereferenceOperations) > 0 && elt.DereferenceOperations[len(elt.DereferenceOperations)-1] == DEREF_POINTER {
		elt.DereferenceOperations = elt.DereferenceOperations[:len(elt.DereferenceOperations)-1]
	} else if elt.IsSlice && len(elt.DereferenceOperations) > 0 && len(inp.Fields) == 0 {
		// elt.DereferenceOperations = append([]int{DEREF_POINTER}, elt.DereferenceOperations...)
	}
}

func ProcessSliceAssignment(expr *CXExpression) {
	if expr.Operator == Natives[OP_IDENTITY] {
		var inp *CXArgument
		var out *CXArgument

		inp = GetAssignmentElement(expr.Inputs[0])
		out = GetAssignmentElement(expr.Outputs[0])

		if inp.IsSlice && out.IsSlice && len(inp.Indexes) == 0 && len(out.Indexes) == 0 {
			out.PassBy = PASSBY_VALUE
		}
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

func UpdateSymbolsTable(symbols *map[string]*CXArgument, sym *CXArgument, offset *int, shouldExist bool) {
	if sym.Name != "" {
		if !sym.IsLocalDeclaration {
			GetGlobalSymbol(symbols, sym.Package, sym.Name)
		}

		if _, found := (*symbols)[sym.Package.Name+"."+sym.Name]; !found {
			if shouldExist {
				// it should exist. error
				println(CompilationError(sym.FileName, sym.FileLine) + " identifier '" + sym.Name + "' does not exist")
				return
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

func ProcessMethodCall(expr *CXExpression, symbols *map[string]*CXArgument, offset *int, shouldExist bool) {
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

							if fn, err := strct.Package.GetMethod(strct.Name+"."+out.Fields[len(out.Fields)-1].Name, strct.Name); err == nil {
								expr.Operator = fn
							} else {
								panic("")
							}

							expr.Inputs = append([]*CXArgument{out}, expr.Inputs...)

							out.Fields = out.Fields[:len(out.Fields)-1]

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

					if fn, err := strct.Package.GetMethod(strct.Name+"."+inp.Fields[len(inp.Fields)-1].Name, strct.Name); err == nil {
						expr.Operator = fn
					} else {
						panic(err)
					}

					inp.Fields = inp.Fields[:len(inp.Fields)-1]
				} else if len(out.Fields) > 0 {
					if argOut, found := (*symbols)[out.Package.Name+"."+out.Name]; found {
						strct := argOut.CustomType

						expr.Inputs = append(expr.Outputs[:1], expr.Inputs...)

						expr.Outputs = expr.Outputs[:len(expr.Outputs)-1]

						if fn, err := strct.Package.GetMethod(strct.Name+"."+out.Fields[len(out.Fields)-1].Name, strct.Name); err == nil {
							expr.Operator = fn
						} else {
							panic(err)
						}

						out.Fields = out.Fields[:len(out.Fields)-1]
					} else {
						panic("")
					}
				}
			}
		} else {
			if out != nil {
				if argOut, found := (*symbols)[out.Package.Name+"."+out.Name]; !found {
					println(CompilationError(out.FileName, out.FileLine), fmt.Sprintf("'%s.%s' not found", out.Package.Name, out.Name))
					return
				} else {
					// then we found an output
					if len(out.Fields) > 0 {
						strct := argOut.CustomType

						if fn, err := strct.Package.GetMethod(strct.Name+"."+out.Fields[len(out.Fields)-1].Name, strct.Name); err == nil {
							expr.Operator = fn
						} else {
							panic("")
						}

						expr.Inputs = append([]*CXArgument{out}, expr.Inputs...)

						out.Fields = out.Fields[:len(out.Fields)-1]

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

func GiveOffset(symbols *map[string]*CXArgument, sym *CXArgument, offset *int, shouldExist bool) {
	if sym.Name != "" {
		if !sym.IsLocalDeclaration {
			GetGlobalSymbol(symbols, sym.Package, sym.Name)
		}

		if arg, found := (*symbols)[sym.Package.Name+"."+sym.Name]; found {
			ProcessSymbolFields(sym, arg)
			CopyArgFields(sym, arg)
		}
	}
}

func ProcessTempVariable(expr *CXExpression) {
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

func CopyArgFields(sym *CXArgument, arg *CXArgument) {
	sym.Offset = arg.Offset
	sym.IsPointer = arg.IsPointer
	sym.IndirectionLevels = arg.IndirectionLevels

	sym.IsSlice = arg.IsSlice
	sym.CustomType = arg.CustomType

	sym.Lengths = arg.Lengths
	sym.Package = arg.Package
	sym.DoesEscape = arg.DoesEscape
	sym.Size = arg.Size

	if arg.Type == TYPE_STR {
		sym.IsPointer = true
	}

	if arg.IsSlice {
		sym.DereferenceOperations = append([]int{DEREF_POINTER}, sym.DereferenceOperations...)
		sym.DereferenceLevels++
	}

	if len(sym.Fields) > 0 {
		sym.Type = sym.Fields[len(sym.Fields)-1].Type
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

func ProcessSymbolFields(sym *CXArgument, arg *CXArgument) {
	if len(sym.Fields) > 0 {
		if arg.CustomType == nil || len(arg.CustomType.Fields) == 0 {
			println(CompilationError(sym.FileName, sym.FileLine), fmt.Sprintf("'%s' has no fields", sym.Name))
			return
		}

		// checking if fields do exist in their CustomType
		// and assigning that CustomType to the sym.Field
		strct := arg.CustomType

		for _, fld := range sym.Fields {
			if inFld, err := strct.GetField(fld.Name); err == nil {
				if inFld.CustomType != nil {
					fld.CustomType = strct
					strct = inFld.CustomType
				}
			} else {
				methodName := sym.Fields[len(sym.Fields)-1].Name
				receiverType := strct.Name

				if method, methodErr := strct.Package.GetMethod(receiverType+"."+methodName, receiverType); methodErr == nil {
					fld.Type = method.Outputs[0].Type
				} else {
					println(CompilationError(fld.FileName, fld.FileLine), err.Error())
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

func SetFinalSize(symbols *map[string]*CXArgument, sym *CXArgument) {
	var finalSize int = sym.TotalSize

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
	if _, found := (*symbols)[symPackage.Name+"."+symName]; !found {
		if glbl, err := symPackage.GetGlobal(symName); err == nil {
			(*symbols)[symPackage.Name+"."+symName] = glbl
		}
	}
}

func PreFinalSize(finalSize *int, sym *CXArgument, arg *CXArgument) {
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
