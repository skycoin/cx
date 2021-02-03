package actions

import (
	"errors"
	"fmt"
	"os"

	. "github.com/skycoin/cx/cx"

	"github.com/jinzhu/copier"
)

// FunctionHeader takes a function name ('ident') and either creates the
// function if it's not known before or returns the already existing function
// if it is.
//
// If the function is a method (isMethod = true), then it adds the object that
// it's called on as the first argument.
//
func FunctionHeader(ident string, receiver []*CXArgument, isMethod bool) *CXFunction {
	if isMethod {
		if len(receiver) > 1 {
			panic("method has multiple receivers")
		}
		if pkg, err := PRGRM.GetCurrentPackage(); err == nil {
			fnName := receiver[0].CustomType.Name + "." + ident

			if fn, err := PRGRM.GetFunction(fnName, pkg.Name); err == nil {
				fn.AddInput(receiver[0])
				pkg.CurrentFunction = fn
				return fn
			} else {
				fn := MakeFunction(fnName, CurrentFile, LineNo)
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
				pkg.CurrentFunction = fn
				return fn
			} else {
				fn := MakeFunction(ident, CurrentFile, LineNo)
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

func isParseOp(expr *CXExpression) bool {
	if expr.Operator != nil && expr.Operator.OpCode > START_PARSE_OPS && expr.Operator.OpCode < END_PARSE_OPS {
		return true
	}
	return false
}

// CheckUndValidTypes checks if an expression with a generic operator (operators that
// accept `TYPE_UNDEFINED` arguments) is receiving arguments of valid types. For example,
// the expression `sa + sb` is not valid if they are struct instances.
func CheckUndValidTypes(expr *CXExpression) {
	if expr.Operator != nil && IsUndOpBasicTypes(expr.Operator) && !IsAllArgsBasicTypes(expr) {
		println(CompilationError(CurrentFile, LineNo), fmt.Sprintf("invalid argument types for '%s' operator", OpNames[expr.Operator.OpCode]))
	}
}

// CheckConcatStr checks if `expr`'s operator is OP_UND_ADD and if its operands are of type str.
// If this is the case, the operator is changed to OP_STR_CONCAT to concatenate the strings.
func CheckConcatStr(expr *CXExpression) {
	if expr.Operator != nil && expr.Operator.OpCode == OP_UND_ADD &&
		expr.Inputs[0].Type == TYPE_STR && expr.Inputs[1].Type == TYPE_STR {
		expr.Operator = Natives[OP_STR_CONCAT]
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

	// each element in the slice corresponds to a different scope
	var symbols *[]map[string]*CXArgument
	tmp := make([]map[string]*CXArgument, 0)
	symbols = &tmp
	*symbols = append(*symbols, make(map[string]*CXArgument, 0))

	// this variable only handles the difference between local and global scopes
	// local being function constrained variables, and global being global variables
	var symbolsScope map[string]bool = make(map[string]bool, 0)

	FunctionProcessParameters(symbols, &symbolsScope, &offset, fn, fn.Inputs)
	FunctionProcessParameters(symbols, &symbolsScope, &offset, fn, fn.Outputs)

	for i, expr := range fn.Expressions {
		if expr.ScopeOperation == SCOPE_NEW {
			*symbols = append(*symbols, make(map[string]*CXArgument, 0))
		}

		ProcessMethodCall(expr, symbols, &offset, true)
		ProcessExpressionArguments(symbols, &symbolsScope, &offset, fn, expr.Inputs, expr, true)
		ProcessExpressionArguments(symbols, &symbolsScope, &offset, fn, expr.Outputs, expr, false)

		ProcessPointerStructs(expr)

		SetCorrectArithmeticOp(expr)
		ProcessTempVariable(expr)
		ProcessSliceAssignment(expr)
		ProcessStringAssignment(expr)
		ProcessReferenceAssignment(expr)

		// process short declaration
		if len(expr.Outputs) > 0 && len(expr.Inputs) > 0 && expr.Outputs[0].IsShortDeclaration && !expr.IsStructLiteral && !isParseOp(expr) {
			if expr.IsMethodCall {
				fn.Expressions[i-1].Outputs[0].Type = fn.Expressions[i].Operator.Outputs[0].Type
				fn.Expressions[i].Outputs[0].Type = fn.Expressions[i].Operator.Outputs[0].Type
			} else {
				fn.Expressions[i-1].Outputs[0].Type = fn.Expressions[i].Inputs[0].Type
				fn.Expressions[i].Outputs[0].Type = fn.Expressions[i].Inputs[0].Type
			}
		}

		processTestExpression(expr)

		CheckTypes(expr)
		CheckUndValidTypes(expr)
		CheckConcatStr(expr)

		if expr.ScopeOperation == SCOPE_REM {
			*symbols = (*symbols)[:len(*symbols)-1]
		}
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
					out.TotalSize = GetSize(inpExpr.Inputs[0])

					out.Type = inpExpr.Inputs[0].Type
					out.PreviouslyDeclared = true
				} else {
					out = MakeArgument(MakeGenSym(LOCAL_PREFIX), CurrentFile, inpExpr.FileLine).AddType(TypeNames[inpExpr.Operator.Outputs[0].Type])
					out.DeclarationSpecifiers = inpExpr.Operator.Outputs[0].DeclarationSpecifiers

					out.CustomType = inpExpr.Operator.Outputs[0].CustomType

					if inpExpr.Operator.Outputs[0].CustomType != nil {
						if strct, err := inpExpr.Package.GetStruct(inpExpr.Operator.Outputs[0].CustomType.Name); err == nil {
							out.Size = strct.Size
							out.TotalSize = strct.Size
						}
					} else {
						out.Size = inpExpr.Operator.Outputs[0].Size
						out.TotalSize = GetSize(inpExpr.Operator.Outputs[0])
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
		return GetSize(GetAssignmentElement(expr.Inputs[0]))
	}
}

// checkSameNativeType checks if all the inputs of an expression are of the same type.
// It is used mainly to prevent implicit castings in arithmetic operations
func checkSameNativeType(expr *CXExpression) error {
	if len(expr.Inputs) < 1 {
		return errors.New("cannot perform arithmetic without operands")
	}
	var typ int = expr.Inputs[0].Type
	for _, inp := range expr.Inputs {
		if inp.Type != typ {
			return errors.New("operands are not of the same type")
		}
		typ = inp.Type
	}
	return nil
}

// isUndOpSameInputTypes checks if the received operator belongs to a list of OP_UND_***
// where its inputs' types must be of the same type
func isUndOpSameInputTypes(op *CXFunction) bool {
	switch op.OpCode {
	case
		OP_UND_EQUAL,
		OP_UND_UNEQUAL,
		OP_UND_BITAND,
		OP_UND_BITXOR,
		OP_UND_BITOR,
		OP_UND_BITCLEAR,
		OP_UND_MUL,
		OP_UND_DIV,
		OP_UND_MOD,
		OP_UND_ADD,
		OP_UND_SUB,
		OP_UND_NEG,
		OP_UND_LT,
		OP_UND_GT,
		OP_UND_LTEQ,
		OP_UND_GTEQ,
		OP_UND_BITSHL,
		OP_UND_BITSHR:
		return true
	}
	return false
}

func ProcessUndExpression(expr *CXExpression) {
	if expr.Operator != nil && isUndOpSameInputTypes(expr.Operator) {
		if err := checkSameNativeType(expr); err != nil {
			println(CompilationError(CurrentFile, LineNo), err.Error())
		}
	}
	if expr.IsUndType {
		for _, out := range expr.Outputs {
			out.Size = undOutputSize(expr)
			out.TotalSize = out.Size
		}
	}
}

func ProcessPointerStructs(expr *CXExpression) {
	for _, arg := range append(expr.Inputs, expr.Outputs...) {
		for _, fld := range arg.Fields {
			if fld.IsPointer && fld.DereferenceLevels == 0 {
				fld.DereferenceLevels++
				fld.DereferenceOperations = append(fld.DereferenceOperations, DEREF_POINTER)
			}
		}
		if arg.IsStruct && arg.IsPointer && len(arg.Fields) > 0 && arg.DereferenceLevels == 0 {
			arg.DereferenceLevels++
			arg.DereferenceOperations = append(arg.DereferenceOperations, DEREF_POINTER)
		}
	}
}

// ProcessAssertExpression checks for the special case of test calls. `assert`, `test`, `panic` are operators where
// their first input's type needs to be the same as its second input's type. This can't be handled by
// `checkSameNativeType` because these test functions' third input parameter is always a `str`.
func processTestExpression(expr *CXExpression) {
	if expr.Operator != nil {
		opCode := expr.Operator.OpCode
		if opCode == OP_ASSERT || opCode == OP_TEST || opCode == OP_PANIC {
			inp1Type := GetFormattedType(expr.Inputs[0])
			inp2Type := GetFormattedType(expr.Inputs[1])
			if inp1Type != inp2Type {
				println(CompilationError(CurrentFile, LineNo), fmt.Sprintf("first and second input arguments' types are not equal in '%s' call ('%s' != '%s')", OpNames[expr.Operator.OpCode], inp1Type, inp2Type))
			}
		}
	}
}

// checkIndexType throws an error if the type of `idx` is not `i32` or `i64`.
func checkIndexType(idx *CXArgument) {
	typ := GetFormattedType(idx)
	if typ != "i32" && typ != "i64" {
		println(CompilationError(idx.FileName, idx.FileLine), fmt.Sprintf("wrong index type; expected either 'i32' or 'i64', got '%s'", typ))
	}
}

// ProcessExpressionArguments performs a series of checks and processes to an expresion's inputs and outputs.
// Some of these checks are: checking if a an input has not been declared, assign a relative offset to the argument,
// and calculate the correct size of the argument.
func ProcessExpressionArguments(symbols *[]map[string]*CXArgument, symbolsScope *map[string]bool, offset *int, fn *CXFunction, args []*CXArgument, expr *CXExpression, isInput bool) {
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
			checkIndexType(idx)
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

// isPointerAdded checks if `sym` has already been added to `fn.ListOfPointers`.
func isPointerAdded(fn *CXFunction, sym *CXArgument) (found bool) {
	for _, ptr := range fn.ListOfPointers {
		if sym.Name == ptr.Name {
			if len(sym.Fields) == 0 && len(ptr.Fields) == 0 {
				found = true
				break
			}

			// Checking if we're referring to the same symbol and
			// same fields being accessed. For instance, `foo.bar` !=
			// `foo.car` as these will have different memory offset to
			// be considered by the garbage collector.
			if len(sym.Fields) > 0 &&
				len(sym.Fields) == len(ptr.Fields) &&
				sym.Fields[len(sym.Fields)-1].Name == ptr.Fields[len(ptr.Fields)-1].Name {
				found = true
				break
			}
		}
	}

	return found
}

// AddPointer checks if `sym` or its last field, if a struct, behaves like a
// pointer (slice, pointer, string). If this is the case, `sym` is added to
// `fn.ListOfPointers` so the CX runtime does not have to determine this.
func AddPointer(fn *CXFunction, sym *CXArgument) {
	// Ignore if it's a global variable.
	if sym.Offset > PRGRM.StackSize {
		return
	}
	// We first need to check if we're going to add `sym` with fields.
	// If `sym` has fields, then we `return` and we don't add the root `sym`.
	// If `sym` has no fields, then we check if `sym` is a pointer and
	// we add it if it is.

	// Field symbol:
	// Checking if it is a pointer candidate and if it was already
	// added to the list.
	if len(sym.Fields) > 0 {
		fld := sym.Fields[len(sym.Fields)-1]
		if IsPointer(fld) && !isPointerAdded(fn, sym) {
			fn.ListOfPointers = append(fn.ListOfPointers, sym)
		}
	}
	// Root symbol:
	// Checking if it is a pointer candidate and if it was already
	// added to the list.
	if IsPointer(sym) && !isPointerAdded(fn, sym) {
		if len(sym.Fields) > 0 {
			tmp := CXArgument{}
			copier.Copy(&tmp, sym)
			tmp.Fields = nil
			fn.ListOfPointers = append(fn.ListOfPointers, &tmp)
		} else {
			fn.ListOfPointers = append(fn.ListOfPointers, sym)
		}
	}
}

// CheckRedeclared checks if `expr` represents a variable declaration and then checks if an
// instance of that variable has already been declared.
func CheckRedeclared(symbols *[]map[string]*CXArgument, expr *CXExpression, sym *CXArgument) {
	if expr.Operator == nil && len(expr.Outputs) > 0 && len(expr.Inputs) == 0 {
		lastIdx := len(*symbols) - 1

		_, found := (*symbols)[lastIdx][sym.Package.Name+"."+sym.Name]
		if found {
			println(CompilationError(sym.FileName, sym.FileLine), fmt.Sprintf("'%s' redeclared", sym.Name))
		}
	}
}

func ProcessLocalDeclaration(symbols *[]map[string]*CXArgument, symbolsScope *map[string]bool, arg *CXArgument) {
	if arg.IsLocalDeclaration {
		(*symbolsScope)[arg.Package.Name+"."+arg.Name] = true
	}
	arg.IsLocalDeclaration = (*symbolsScope)[arg.Package.Name+"."+arg.Name]
}

func FunctionProcessParameters(symbols *[]map[string]*CXArgument, symbolsScope *map[string]bool, offset *int, fn *CXFunction, params []*CXArgument) {
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

func checkMatchParamTypes(expr *CXExpression, expected, received []*CXArgument, isInputs bool) {
	for i, inp := range expected {
		expectedType := GetFormattedType(expected[i])
		receivedType := GetFormattedType(received[i])

		if expr.IsMethodCall && expected[i].IsPointer && i == 0 {
			// if method receiver is pointer, remove *
			if expectedType[0] == '*' {
				// we need to check if it's not an `str`
				// otherwise we end up removing the `s` instead of a `*`
				expectedType = expectedType[1:]
			}
		}

		if expectedType != receivedType && inp.Type != TYPE_UNDEFINED {
			var opName string
			if expr.Operator.IsNative {
				opName = OpNames[expr.Operator.OpCode]
			} else {
				opName = expr.Operator.Name
			}

			if isInputs {
				println(CompilationError(received[i].FileName, received[i].FileLine), fmt.Sprintf("function '%s' expected input argument of type '%s'; '%s' was provided", opName, expectedType, receivedType))
			} else {
				println(CompilationError(expr.Outputs[i].FileName, expr.Outputs[i].FileLine), fmt.Sprintf("function '%s' expected receiving variable of type '%s'; '%s' was provided", opName, expectedType, receivedType))
			}

		}

		// In the case of assignment we need to check that the input's type matches the output's type.
		// FIXME: There are some expressions added by the parser where temporary variables are used.
		// These temporary variables' types are not properly being set. That's why we use !IsTempVar to
		// exclude these cases for now.
		if expr.Operator.OpCode == OP_IDENTITY && !IsTempVar(expr.Outputs[0].Name) {
			inpType := GetFormattedType(expr.Inputs[0])
			outType := GetFormattedType(expr.Outputs[0])

			// We use `isInputs` to only print the error once.
			// Otherwise we'd print the error twice: once for the input and again for the output
			if inpType != outType && isInputs {
				println(CompilationError(received[i].FileName, received[i].FileLine), fmt.Sprintf("cannot assign value of type '%s' to identifier '%s' of type '%s'", inpType, GetAssignmentElement(expr.Outputs[0]).Name, outType))
			}
		}
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

		// checking if number of expr.Outputs matches number of Operator.Outputs
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
			os.Exit(CX_COMPILATION_ERROR)
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

	// then it's a function call and not a declaration
	if expr.Operator != nil {
		// checking inputs matching operator's inputs
		checkMatchParamTypes(expr, expr.Operator.Inputs, expr.Inputs, true)

		// checking outputs matching operator's outputs
		checkMatchParamTypes(expr, expr.Operator.Outputs, expr.Outputs, false)
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

// ProcessReferenceAssignment checks if the reference of a symbol can be assigned to the expression's output.
// For example: `var foo i32; var bar i32; bar = &foo` is not valid.
func ProcessReferenceAssignment(expr *CXExpression) {
	for _, out := range expr.Outputs {
		elt := GetAssignmentElement(out)
		if elt.PassBy == PASSBY_REFERENCE &&
			!hasDeclSpec(elt, DECL_POINTER) &&
			elt.Type != TYPE_STR && !elt.IsSlice {
			println(CompilationError(CurrentFile, LineNo), "invalid reference assignment", elt.Name)
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

			// we want to pass by value if we're sending the slice as a whole (no indexing)
			// unless it's a pointer to the slice
			if assignElt.IsSlice && len(assignElt.Indexes) == 0 && !hasDeclSpec(assignElt, DECL_POINTER) {
				assignElt.PassBy = PASSBY_VALUE
			}
		}
	}
}

// lookupSymbol searches for `ident` in `symbols`, starting from the innermost scope.
func lookupSymbol(pkgName, ident string, symbols *[]map[string]*CXArgument) (*CXArgument, error) {
	fullName := pkgName + "." + ident
	for c := len(*symbols) - 1; c >= 0; c-- {
		if sym, found := (*symbols)[c][fullName]; found {
			return sym, nil
		}
	}

	// Checking if `ident` refers to a function.
	pkg, err := PRGRM.GetPackage(pkgName)
	if err != nil {
		return nil, err
	}

	notFound := errors.New("identifier '" + ident + "' does not exist")

	// We're not checking for that error
	fn, err := pkg.GetFunction(ident)
	if err != nil {
		return nil, notFound
	}
	// Then we found a function by that name. Let's create a `CXArgument` of
	// type `func` with that name.
	fnArg := MakeArgument(ident, fn.FileName, fn.FileLine).AddType(TypeNames[TYPE_FUNC])
	fnArg.Package = pkg

	return fnArg, nil
}

// UpdateSymbolsTable adds `sym` to the innermost scope (last element of slice) in `symbols`.
func UpdateSymbolsTable(symbols *[]map[string]*CXArgument, sym *CXArgument, offset *int, shouldExist bool) {
	if sym.Name != "" {
		if !sym.IsLocalDeclaration {
			GetGlobalSymbol(symbols, sym.Package, sym.Name)
		}

		lastIdx := len(*symbols) - 1
		fullName := sym.Package.Name + "." + sym.Name

		// outerSym, err := lookupSymbol(sym.Package.Name, sym.Name, symbols)
		_, err := lookupSymbol(sym.Package.Name, sym.Name, symbols)
		_, found := (*symbols)[lastIdx][fullName]

		// then it wasn't found in any scope
		if err != nil && shouldExist {
			println(CompilationError(sym.FileName, sym.FileLine), "identifier '"+sym.Name+"' does not exist")
		}

		// then it was already added in the innermost scope
		if found {
			return
		}

		// then it is a new declaration
		if !shouldExist && !found {
			// then it was declared in an outer scope
			sym.Offset = *offset
			(*symbols)[lastIdx][fullName] = sym
			*offset += GetSize(sym)
		}
	}
}

func ProcessMethodCall(expr *CXExpression, symbols *[]map[string]*CXArgument, offset *int, shouldExist bool) {
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
			// if argInp, found := (*symbols)[lastIdx][inp.Package.Name+"."+inp.Name]; !found {
			if argInp, err := lookupSymbol(inp.Package.Name, inp.Name, symbols); err != nil {
				if out == nil {
					panic("")
				}
				argOut, err := lookupSymbol(out.Package.Name, out.Name, symbols)
				if err != nil {
					println(CompilationError(out.FileName, out.FileLine), fmt.Sprintf("identifier '%s' does not exist", out.Name))
					os.Exit(CX_COMPILATION_ERROR)
				}
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
					argOut, err := lookupSymbol(out.Package.Name, out.Name, symbols)
					if err != nil {
						panic("")
					}

					strct := argOut.CustomType

					if strct == nil {
						println(CompilationError(argOut.FileName, argOut.FileLine), fmt.Sprintf("illegal method call or field access on identifier '%s' of primitive type '%s'", argOut.Name, TypeNames[argOut.Type]))
						os.Exit(CX_COMPILATION_ERROR)
					}

					expr.Inputs = append(expr.Outputs[:1], expr.Inputs...)

					expr.Outputs = expr.Outputs[:len(expr.Outputs)-1]

					if fn, err := strct.Package.GetMethod(strct.Name+"."+out.Fields[len(out.Fields)-1].Name, strct.Name); err == nil {
						expr.Operator = fn
					} else {
						panic(err)
					}

					out.Fields = out.Fields[:len(out.Fields)-1]
				}
			}
		} else {
			if out == nil {
				panic("")
			}

			argOut, err := lookupSymbol(out.Package.Name, out.Name, symbols)
			if err != nil {
				println(CompilationError(out.FileName, out.FileLine), fmt.Sprintf("identifier '%s' does not exist", out.Name))
				os.Exit(CX_COMPILATION_ERROR)
			}

			// then we found an output
			if len(out.Fields) > 0 {
				strct := argOut.CustomType

				if strct == nil {
					println(CompilationError(argOut.FileName, argOut.FileLine), fmt.Sprintf("illegal method call or field access on identifier '%s' of primitive type '%s'", argOut.Name, TypeNames[argOut.Type]))
					os.Exit(CX_COMPILATION_ERROR)
				}

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

		// checking if receiver is sent as pointer or not
		if expr.Operator.Inputs[0].IsPointer {
			expr.Inputs[0].PassBy = PASSBY_REFERENCE
		}
	}
}

func GiveOffset(symbols *[]map[string]*CXArgument, sym *CXArgument, offset *int, shouldExist bool) {
	if sym.Name != "" {
		if !sym.IsLocalDeclaration {
			GetGlobalSymbol(symbols, sym.Package, sym.Name)
		}

		arg, err := lookupSymbol(sym.Package.Name, sym.Name, symbols)
		if err == nil {
			ProcessSymbolFields(sym, arg)
			CopyArgFields(sym, arg)
		}
	}
}

func ProcessTempVariable(expr *CXExpression) {
	if expr.Operator != nil && (expr.Operator == Natives[OP_IDENTITY] || IsUndOpMimicInput(expr.Operator)) && len(expr.Outputs) > 0 && len(expr.Inputs) > 0 {
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

	if sym.FileLine != arg.FileLine {
		// FIXME Maybe we can unify this later.
		if len(sym.Fields) > 0 {
			elt := GetAssignmentElement(sym)

			declSpec := []int{}
			for c := 0; c < len(elt.DeclarationSpecifiers); c++ {
				switch elt.DeclarationSpecifiers[c] {
				case DECL_INDEXING:
					if declSpec[len(declSpec)-1] == DECL_ARRAY || declSpec[len(declSpec)-1] == DECL_SLICE {
						declSpec = declSpec[:len(declSpec)-1]
					} else {
						println(CompilationError(sym.FileName, sym.FileLine), fmt.Sprintf("invalid indexing"))
					}
				case DECL_DEREF:
					if declSpec[len(declSpec)-1] == DECL_POINTER {
						declSpec = declSpec[:len(declSpec)-1]
					} else {
						println(CompilationError(sym.FileName, sym.FileLine), fmt.Sprintf("invalid indirection"))
					}
				default:
					declSpec = append(declSpec, elt.DeclarationSpecifiers[c])
				}
			}

			elt.DeclarationSpecifiers = declSpec
		} else {
			declSpec := make([]int, len(arg.DeclarationSpecifiers))

			for i, spec := range arg.DeclarationSpecifiers {
				declSpec[i] = spec
			}

			for _, spec := range sym.DeclarationSpecifiers {
				// checking if we need to remove or add DECL_POINTERs
				// also we could be removing
				switch spec {
				case DECL_INDEXING:
					if declSpec[len(declSpec)-1] == DECL_ARRAY || declSpec[len(declSpec)-1] == DECL_SLICE {
						declSpec = declSpec[:len(declSpec)-1]
					} else {
						println(CompilationError(sym.FileName, sym.FileLine), fmt.Sprintf("invalid indexing"))
					}
				case DECL_DEREF:
					if declSpec[len(declSpec)-1] == DECL_POINTER {
						declSpec = declSpec[:len(declSpec)-1]
					} else {
						println(CompilationError(sym.FileName, sym.FileLine), fmt.Sprintf("invalid indirection"))
					}
				case DECL_POINTER:
					if sym.FileLine != arg.FileLine {
						// This function is also called so it assigns offset and other fields to signature parameters
						//
						declSpec = append(declSpec, DECL_POINTER)
					}
				}
			}

			sym.DeclarationSpecifiers = declSpec
		}
	} else {
		sym.DeclarationSpecifiers = arg.DeclarationSpecifiers
	}

	sym.IsSlice = arg.IsSlice
	sym.CustomType = arg.CustomType

	// FIXME: In other processes like ProcessSymbolFields the symbol is assigned with lengths.
	// If we already have some lengths, we skip this. This needs to be fixed in the redesign of the parser.
	if len(sym.Lengths) == 0 {
		sym.Lengths = arg.Lengths
	}

	// sym.Lengths = arg.Lengths
	sym.Package = arg.Package
	sym.DoesEscape = arg.DoesEscape
	sym.Size = arg.Size

	if arg.Type == TYPE_STR {
		sym.IsPointer = true
	}

	// Checking if it's a slice struct field. We'll do the same process as
	// below (as in the `arg.IsSlice` check), but the process differs in the
	// case of a slice struct field.
	elt := GetAssignmentElement(sym)
	if !arg.IsSlice && arg.CustomType != nil && elt.IsSlice {
		// elt.DereferenceOperations = []int{4, 4}
		for i, deref := range elt.DereferenceOperations {
			// The parser when reading `foo[5]` in postfix.go does not know if `foo`
			// is a slice or an array. At this point we now know it's a slice and we need
			// to change those dereferences to DEREF_SLICE.
			if deref == DEREF_ARRAY {
				elt.DereferenceOperations[i] = DEREF_SLICE
			}
		}
		if elt.DereferenceOperations[0] == DEREF_POINTER {
			elt.DereferenceOperations = elt.DereferenceOperations[1:]
		}
	}

	if arg.IsSlice {
		if !hasDerefOp(sym, DEREF_ARRAY) {
			// Then we're handling the slice itself, and we need to dereference it.
			sym.DereferenceOperations = append([]int{DEREF_POINTER}, sym.DereferenceOperations...)
		} else {
			for i, deref := range sym.DereferenceOperations {
				// The parser when reading `foo[5]` in postfix.go does not know if `foo`
				// is a slice or an array. At this point we now know it's a slice and we need
				// to change those dereferences to DEREF_SLICE.
				if deref == DEREF_ARRAY {
					sym.DereferenceOperations[i] = DEREF_SLICE
				}
			}
		}
	}

	if len(sym.Fields) > 0 {
		sym.Type = sym.Fields[len(sym.Fields)-1].Type
		sym.IsSlice = sym.Fields[len(sym.Fields)-1].IsSlice
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

					sym.Lengths = fld.Lengths

					// nameFld.DeclarationSpecifiers = fld.DeclarationSpecifiers
					// nameFld.DeclarationSpecifiers = append(fld.DeclarationSpecifiers, nameFld.DeclarationSpecifiers[1:]...)
					if len(nameFld.DeclarationSpecifiers) > 0 {
						nameFld.DeclarationSpecifiers = append(fld.DeclarationSpecifiers, nameFld.DeclarationSpecifiers[1:]...)
					} else {
						nameFld.DeclarationSpecifiers = fld.DeclarationSpecifiers
					}

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

				nameFld.Offset += GetSize(fld)
			}
		}
	}
}

func SetFinalSize(symbols *[]map[string]*CXArgument, sym *CXArgument) {
	var finalSize int = sym.TotalSize

	arg, err := lookupSymbol(sym.Package.Name, sym.Name, symbols)
	if err == nil {
		PreFinalSize(&finalSize, sym, arg)
		for _, fld := range sym.Fields {
			finalSize = fld.TotalSize
			PreFinalSize(&finalSize, fld, arg)
		}
	}

	sym.TotalSize = finalSize
}

// GetGlobalSymbol tries to retrieve `ident` from `symPkg`'s globals if `ident` is not found in the local scope.
func GetGlobalSymbol(symbols *[]map[string]*CXArgument, symPkg *CXPackage, ident string) {
	_, err := lookupSymbol(symPkg.Name, ident, symbols)
	if err != nil {
		if glbl, err := symPkg.GetGlobal(ident); err == nil {
			lastIdx := len(*symbols) - 1
			(*symbols)[lastIdx][symPkg.Name+"."+ident] = glbl
		}
	}
}

func PreFinalSize(finalSize *int, sym *CXArgument, arg *CXArgument) {
	idxCounter := 0
	elt := GetAssignmentElement(sym)
	for _, op := range elt.DereferenceOperations {
		if elt.IsSlice {
			continue
		}
		switch op {
		case DEREF_ARRAY:
			*finalSize /= elt.Lengths[idxCounter]
			idxCounter++
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
