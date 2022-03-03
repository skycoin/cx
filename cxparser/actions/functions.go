package actions

import (
	"errors"
	"fmt"
	"os"

	"github.com/jinzhu/copier"
	"github.com/skycoin/cx/cx/ast"
	"github.com/skycoin/cx/cx/constants"
	"github.com/skycoin/cx/cx/types"
)

// FunctionHeader takes a function name ('ident') and either creates the
// function if it's not known before or returns the already existing function
// if it is.
//
// If the function is a method (isMethod = true), then it adds the object that
// it's called on as the first argument.
//
func FunctionHeader(prgrm *ast.CXProgram, ident string, receiver []*ast.CXArgument, isMethod bool) *ast.CXFunction {
	if isMethod {
		if len(receiver) > 1 {
			panic("method has multiple receivers")
		}
		if pkg, err := prgrm.GetCurrentPackage(); err == nil {
			fnName := receiver[0].StructType.Name + "." + ident

			if fn, err := prgrm.GetFunction(fnName, pkg.Name); err == nil {
				fn.AddInput(prgrm, receiver[0])
				pkg.CurrentFunction = ast.CXFunctionIndex(fn.Index)
				return fn
			} else {
				fn := ast.MakeFunction(fnName, CurrentFile, LineNo)
				_, fnIdx := pkg.AddFunction(prgrm, fn)
				newFn := prgrm.GetFunctionFromArray(fnIdx)
				newFn.AddInput(prgrm, receiver[0])

				return newFn
			}
		} else {
			panic(err)
		}
	} else {
		if pkg, err := prgrm.GetCurrentPackage(); err == nil {
			if fn, err := prgrm.GetFunction(ident, pkg.Name); err == nil {
				pkg.CurrentFunction = ast.CXFunctionIndex(fn.Index)
				return fn
			} else {
				fn := ast.MakeFunction(ident, CurrentFile, LineNo)
				_, fnIdx := pkg.AddFunction(prgrm, fn)
				newFn := prgrm.GetFunctionFromArray(fnIdx)

				return newFn
			}
		} else {
			panic(err)
		}
	}
}

func FunctionAddParameters(prgrm *ast.CXProgram, fn *ast.CXFunction, inputs, outputs []*ast.CXArgument) {
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
		fn.AddInput(prgrm, inp)
	}

	for _, out := range outputs {
		fn.AddOutput(prgrm, out)
	}

	for _, outIdx := range fn.Outputs {
		out := prgrm.GetCXArgFromArray(outIdx)
		if out.IsPointer() && out.PointerTargetType != types.STR && out.Type != types.AFF {
			out.DoesEscape = true
		}
	}
}

func isParseOp(prgrm *ast.CXProgram, expr *ast.CXExpression) bool {
	exprAtomicOp, _, _, err := prgrm.GetOperation(expr)
	if err != nil {
		panic(err)
	}
	exprAtomicOpOperator := prgrm.GetFunctionFromArray(exprAtomicOp.Operator)
	if exprAtomicOpOperator != nil && exprAtomicOpOperator.AtomicOPCode > constants.START_PARSE_OPS && exprAtomicOpOperator.AtomicOPCode < constants.END_PARSE_OPS {
		return true
	}
	return false
}

// CheckUndValidTypes checks if an expression with a generic operator (operators that
// accept `cxcore.TYPE_UNDEFINED` arguments) is receiving arguments of valid types. For example,
// the expression `sa + sb` is not valid if they are struct instances.
func CheckUndValidTypes(prgrm *ast.CXProgram, expr *ast.CXExpression) {
	exprAtomicOp, _, _, err := prgrm.GetOperation(expr)
	if err != nil {
		panic(err)
	}
	exprAtomicOpOperator := prgrm.GetFunctionFromArray(exprAtomicOp.Operator)
	if exprAtomicOpOperator != nil && ast.IsOperator(exprAtomicOpOperator.AtomicOPCode) && !IsAllArgsBasicTypes(prgrm, expr) {
		println(ast.CompilationError(CurrentFile, LineNo), fmt.Sprintf("invalid argument types for '%s' operator", ast.OpNames[exprAtomicOpOperator.AtomicOPCode]))
	}
}

func FunctionProcessParameters(prgrm *ast.CXProgram, symbols *[]map[string]*ast.CXArgument, symbolsScope *map[string]bool, offset *types.Pointer, fn *ast.CXFunction, params []*ast.CXArgument) {
	for _, param := range params {
		ProcessLocalDeclaration(prgrm, symbols, symbolsScope, param)

		UpdateSymbolsTable(prgrm, symbols, param, offset, false)
		GiveOffset(prgrm, symbols, param, offset, false)
		SetFinalSize(prgrm, symbols, param)

		AddPointer(prgrm, fn, param)

		// as these are declarations, they should not have any dereference operations
		param.DereferenceOperations = nil
	}
}

func FunctionDeclaration(prgrm *ast.CXProgram, fn *ast.CXFunction, inputs, outputs []*ast.CXArgument, exprs []ast.CXExpression) {
	//var exprs []*cxcore.CXExpression = prgrm.SysInitExprs

	FunctionAddParameters(prgrm, fn, inputs, outputs)

	// getting offset to use by statements (excluding inputs, outputs and receiver)
	var offset types.Pointer
	//TODO: Why would the heap starting position always be incrasing?
	//TODO: HeapStartsAt only increases, with every write?
	//DataOffset only increases

	ProcessGoTos(prgrm, fn, exprs)

	fn.LineCount = len(fn.Expressions)

	// each element in the slice corresponds to a different scope
	var symbols *[]map[string]*ast.CXArgument
	tmp := make([]map[string]*ast.CXArgument, 0)
	symbols = &tmp
	*symbols = append(*symbols, make(map[string]*ast.CXArgument))

	// this variable only handles the difference between local and global scopes
	// local being function constrained variables, and global being global variables
	var symbolsScope map[string]bool = make(map[string]bool)

	FunctionProcessParameters(prgrm, symbols, &symbolsScope, &offset, fn, prgrm.ConvertIndexArgsToPointerArgs(fn.Inputs))
	FunctionProcessParameters(prgrm, symbols, &symbolsScope, &offset, fn, prgrm.ConvertIndexArgsToPointerArgs(fn.Outputs))

	for i, expr := range fn.Expressions {
		if expr.Type == ast.CX_LINE {
			continue
		}
		exprAtomicOp, _, _, err := prgrm.GetOperation(&expr)
		if err != nil {
			panic(err)
		}

		if expr.IsScopeNew() {
			*symbols = append(*symbols, make(map[string]*ast.CXArgument))
		}

		ProcessMethodCall(prgrm, &expr, symbols, &offset, true)
		ProcessExpressionArguments(prgrm, symbols, &symbolsScope, &offset, fn, prgrm.ConvertIndexArgsToPointerArgs(exprAtomicOp.Inputs), &expr, true)
		ProcessExpressionArguments(prgrm, symbols, &symbolsScope, &offset, fn, prgrm.ConvertIndexArgsToPointerArgs(exprAtomicOp.Outputs), &expr, false)

		ProcessPointerStructs(prgrm, &expr)

		SetCorrectArithmeticOp(prgrm, &expr)
		ProcessTempVariable(prgrm, &expr)
		ProcessSliceAssignment(prgrm, &expr)
		ProcessStringAssignment(prgrm, &expr)
		ProcessReferenceAssignment(prgrm, &expr)

		// process short declaration
		if len(exprAtomicOp.Outputs) > 0 && len(exprAtomicOp.Inputs) > 0 && prgrm.GetCXArgFromArray(exprAtomicOp.Outputs[0]).IsShortAssignmentDeclaration && !expr.IsStructLiteral() && !isParseOp(prgrm, &expr) {
			var arg *ast.CXArgument

			exprAtomicOp, err := prgrm.GetCXAtomicOpFromExpressions(fn.Expressions, i)
			if err != nil {
				panic(err)
			}
			exprAtomicOpOperator := prgrm.GetFunctionFromArray(exprAtomicOp.Operator)

			exprBeforeAtomicOp, err := prgrm.GetPreviousCXAtomicOpFromExpressions(fn.Expressions, i-1)
			if err != nil {
				panic(err)
			}

			if expr.IsMethodCall() {
				arg = prgrm.GetCXArgFromArray(exprAtomicOpOperator.Outputs[0])
			} else {
				arg = prgrm.GetCXArgFromArray(exprAtomicOp.Inputs[0])
			}

			exprBeforeAtomicOpOutput := prgrm.GetCXArgFromArray(exprBeforeAtomicOp.Outputs[0])
			exprBeforeAtomicOpOutput.Type = arg.Type
			exprBeforeAtomicOpOutput.PointerTargetType = arg.PointerTargetType
			exprBeforeAtomicOpOutput.Size = arg.Size
			exprBeforeAtomicOpOutput.TotalSize = arg.TotalSize

			exprAtomicOpOutput := prgrm.GetCXArgFromArray(exprAtomicOp.Outputs[0])
			exprAtomicOpOutput.Type = arg.Type
			exprAtomicOpOutput.PointerTargetType = arg.PointerTargetType
			exprAtomicOpOutput.Size = arg.Size
			exprAtomicOpOutput.TotalSize = arg.TotalSize
		}

		processTestExpression(prgrm, &expr)

		CheckTypes(prgrm, fn.Expressions, i)
		CheckUndValidTypes(prgrm, &expr)

		if expr.IsScopeDel() {
			*symbols = (*symbols)[:len(*symbols)-1]
		}
	}

	fn.Size = offset

	// TODO: temporary bug fix, needs improvements
	if fn.Name == "main" {
		prgrm.CXFunctions[fn.Index] = *fn
	}
}

func FunctionCall(prgrm *ast.CXProgram, exprs []ast.CXExpression, args []ast.CXExpression) []ast.CXExpression {
	expr := exprs[len(exprs)-1]

	cxAtomicOp, _, _, err := prgrm.GetOperation(&expr)
	if err != nil {
		panic(err)
	}
	cxAtomicOpOperator := prgrm.GetFunctionFromArray(cxAtomicOp.Operator)
	if cxAtomicOpOperator == nil {
		cxAtomicOpOutput := prgrm.GetCXArgFromArray(cxAtomicOp.Outputs[0])
		opName := cxAtomicOpOutput.Name
		opPkgIdx := cxAtomicOpOutput.Package

		opPkg, err := prgrm.GetPackageFromArray(opPkgIdx)
		if err != nil {
			panic(err)
		}

		if op, err := prgrm.GetFunction(opName, opPkg.Name); err == nil {
			cxAtomicOp.Operator = ast.CXFunctionIndex(op.Index)
		} else if cxAtomicOpOutput.Fields == nil {
			// then it's not a possible method call
			println(ast.CompilationError(CurrentFile, LineNo), err.Error())
			return nil
		} else {
			exprs[len(exprs)-1].ExpressionType = ast.CXEXPR_METHOD_CALL
		}

		if len(cxAtomicOp.Outputs) > 0 && cxAtomicOpOutput.Fields == nil {
			cxAtomicOp.Outputs = nil
		}
	}

	var nestedExprs []ast.CXExpression
	for i, inpExpr := range args {
		if inpExpr.Type == ast.CX_LINE {
			continue
		}

		inpExprAtomicOp, _, _, err := prgrm.GetOperation(&inpExpr)
		if err != nil {
			panic(err)
		}
		inpExprAtomicOpOperator := prgrm.GetFunctionFromArray(inpExprAtomicOp.Operator)

		inpExprCXLine, _ := prgrm.GetPreviousCXLine(args, i)

		if inpExprAtomicOpOperator == nil {
			// then it's a literal
			cxAtomicOp.AddInput(prgrm, inpExprAtomicOp.Outputs[0])
		} else {
			// then it's a function call
			if len(inpExprAtomicOp.Outputs) < 1 {
				var out *ast.CXArgument

				inpExprAtomicOpOperatorOutput := prgrm.GetCXArgFromArray(inpExprAtomicOpOperator.Outputs[0])
				if inpExprAtomicOpOperatorOutput.Type == types.UNDEFINED {
					// if undefined type, then adopt argument's type
					inpExprAtomicOpInput := prgrm.GetCXArgFromArray(inpExprAtomicOp.Inputs[0])
					out = ast.MakeArgument(MakeGenSym(constants.LOCAL_PREFIX), CurrentFile, inpExprCXLine.LineNumber).AddType(inpExprAtomicOpInput.Type)
					out.StructType = inpExprAtomicOpInput.StructType

					out.Size = inpExprAtomicOpInput.Size
					out.TotalSize = ast.GetSize(prgrm, inpExprAtomicOpInput)

					out.Type = inpExprAtomicOpInput.Type
					out.PointerTargetType = inpExprAtomicOpInput.PointerTargetType
					out.PreviouslyDeclared = true
				} else {
					out = ast.MakeArgument(MakeGenSym(constants.LOCAL_PREFIX), CurrentFile, inpExprCXLine.LineNumber).AddType(inpExprAtomicOpOperatorOutput.Type)
					out.DeclarationSpecifiers = inpExprAtomicOpOperatorOutput.DeclarationSpecifiers

					out.StructType = inpExprAtomicOpOperatorOutput.StructType

					if inpExprAtomicOpOperatorOutput.StructType != nil {
						inpExprPkg, err := prgrm.GetPackageFromArray(inpExprAtomicOp.Package)
						if err != nil {
							panic(err)
						}
						if strct, err := inpExprPkg.GetStruct(prgrm, inpExprAtomicOpOperatorOutput.StructType.Name); err == nil {
							out.Size = strct.Size
							out.TotalSize = strct.Size
						}
					} else {
						out.Size = inpExprAtomicOpOperatorOutput.Size
						out.TotalSize = ast.GetSize(prgrm, inpExprAtomicOpOperatorOutput)
					}

					out.Type = inpExprAtomicOpOperatorOutput.Type
					out.PointerTargetType = inpExprAtomicOpOperatorOutput.PointerTargetType
					out.PreviouslyDeclared = true
				}

				out.Package = inpExprAtomicOp.Package
				outIdx := prgrm.AddCXArgInArray(out)

				inpExprAtomicOp.AddOutput(prgrm, outIdx)
				cxAtomicOp.AddInput(prgrm, outIdx)
			}
			if len(inpExprAtomicOp.Outputs) > 0 && inpExpr.IsArrayLiteral() {
				cxAtomicOp.AddInput(prgrm, inpExprAtomicOp.Outputs[0])
			}
			nestedExprs = append(nestedExprs, inpExpr)
		}
	}
	return append(nestedExprs, exprs...)
}

// checkSameNativeType checks if all the inputs of an expression are of the same type.
// It is used mainly to prevent implicit castings in arithmetic operations
func checkSameNativeType(prgrm *ast.CXProgram, expr *ast.CXExpression) error {
	cxAtomicOp, _, _, err := prgrm.GetOperation(expr)
	if err != nil {
		panic(err)
	}

	if len(cxAtomicOp.Inputs) < 1 {
		return errors.New("cannot perform arithmetic without operands")
	}

	typeCode := prgrm.GetCXArgFromArray(cxAtomicOp.Inputs[0]).Type
	if prgrm.GetCXArgFromArray(cxAtomicOp.Inputs[0]).Type == types.POINTER {
		typeCode = prgrm.GetCXArgFromArray(cxAtomicOp.Inputs[0]).PointerTargetType
	}

	for _, inpIdx := range cxAtomicOp.Inputs {
		inp := prgrm.GetCXArgFromArray(inpIdx)
		inpType := inp.Type

		if inp.Type == types.POINTER {
			inpType = inp.PointerTargetType
		}

		if inpType != typeCode {
			return errors.New("operands are not of the same type")
		}

		typeCode = inpType
	}

	return nil
}

func ProcessOperatorExpression(prgrm *ast.CXProgram, expr *ast.CXExpression) {
	cxAtomicOp, _, _, err := prgrm.GetOperation(expr)
	if err != nil {
		panic(err)
	}
	cxAtomicOpOperator := prgrm.GetFunctionFromArray(cxAtomicOp.Operator)

	if cxAtomicOpOperator != nil && ast.IsOperator(cxAtomicOpOperator.AtomicOPCode) {
		if err := checkSameNativeType(prgrm, expr); err != nil {
			println(ast.CompilationError(CurrentFile, LineNo), err.Error())
		}
	}
	if expr.IsUndType(prgrm) {
		for _, outIdx := range cxAtomicOp.Outputs {
			out := prgrm.GetCXArgFromArray(outIdx)

			size := types.Pointer(1)
			if !ast.IsComparisonOperator(cxAtomicOpOperator.AtomicOPCode) {
				size = ast.GetSize(prgrm, prgrm.GetCXArgFromArray(cxAtomicOp.Inputs[0]).GetAssignmentElement(prgrm))
			}
			out.Size = size
			out.TotalSize = size
		}
	}
}

func ProcessPointerStructs(prgrm *ast.CXProgram, expr *ast.CXExpression) {
	cxAtomicOp, _, _, err := prgrm.GetOperation(expr)
	if err != nil {
		panic(err)
	}

	for _, argIdx := range append(cxAtomicOp.Inputs, cxAtomicOp.Outputs...) {
		arg := prgrm.GetCXArgFromArray(argIdx)
		for _, fldIdx := range arg.Fields {
			fld := prgrm.GetCXArgFromArray(fldIdx)
			if fld.IsPointer() && fld.DereferenceLevels == 0 {
				fld.DereferenceLevels++
				fld.DereferenceOperations = append(fld.DereferenceOperations, constants.DEREF_POINTER)
			}
		}
		if arg.IsStruct && arg.IsPointer() && len(arg.Fields) > 0 && arg.DereferenceLevels == 0 {
			arg.DereferenceLevels++
			arg.DereferenceOperations = append(arg.DereferenceOperations, constants.DEREF_POINTER)
		}
	}
}

// ProcessAssertExpression checks for the special case of test calls. `assert`, `test`, `panic` are operators where
// their first input's type needs to be the same as its second input's type. This can't be handled by
// `checkSameNativeType` because these test functions' third input parameter is always a `str`.
func processTestExpression(prgrm *ast.CXProgram, expr *ast.CXExpression) {
	cxAtomicOp, _, _, err := prgrm.GetOperation(expr)
	if err != nil {
		panic(err)
	}
	cxAtomicOpOperator := prgrm.GetFunctionFromArray(cxAtomicOp.Operator)

	if cxAtomicOpOperator != nil {
		opCode := cxAtomicOpOperator.AtomicOPCode
		if opCode == constants.OP_ASSERT || opCode == constants.OP_TEST || opCode == constants.OP_PANIC {
			inp1Type := ast.GetFormattedType(prgrm, prgrm.GetCXArgFromArray(cxAtomicOp.Inputs[0]))
			inp2Type := ast.GetFormattedType(prgrm, prgrm.GetCXArgFromArray(cxAtomicOp.Inputs[1]))
			if inp1Type != inp2Type {
				println(ast.CompilationError(CurrentFile, LineNo), fmt.Sprintf("first and second input arguments' types are not equal in '%s' call ('%s' != '%s')", ast.OpNames[cxAtomicOpOperator.AtomicOPCode], inp1Type, inp2Type))
			}
		}
	}
}

// checkIndexType throws an error if the type of `idx` is not `i32` or `i64`.
func checkIndexType(prgrm *ast.CXProgram, idx *ast.CXArgument) {
	typ := ast.GetFormattedType(prgrm, idx)
	if typ != "i32" && typ != "i64" {
		println(ast.CompilationError(idx.ArgDetails.FileName, idx.ArgDetails.FileLine), fmt.Sprintf("wrong index type; expected either 'i32' or 'i64', got '%s'", typ))
	}
}

// ProcessExpressionArguments performs a series of checks and processes to an expresion's inputs and outputs.
// Some of these checks are: checking if a an input has not been declared, assign a relative offset to the argument,
// and calculate the correct size of the argument.
func ProcessExpressionArguments(prgrm *ast.CXProgram, symbols *[]map[string]*ast.CXArgument, symbolsScope *map[string]bool, offset *types.Pointer, fn *ast.CXFunction, args []*ast.CXArgument, expr *ast.CXExpression, isInput bool) {
	for _, arg := range args {
		ProcessLocalDeclaration(prgrm, symbols, symbolsScope, arg)

		if !isInput {
			CheckRedeclared(prgrm, symbols, expr, arg)
		}

		if !isInput {
			ProcessOperatorExpression(prgrm, expr)
		}

		if arg.PreviouslyDeclared {
			UpdateSymbolsTable(prgrm, symbols, arg, offset, false)
		} else {
			UpdateSymbolsTable(prgrm, symbols, arg, offset, true)
		}

		if isInput {
			GiveOffset(prgrm, symbols, arg, offset, true)
		} else {
			GiveOffset(prgrm, symbols, arg, offset, false)
		}

		ProcessSlice(prgrm, arg)

		for _, idxIdx := range arg.Indexes {
			idx := prgrm.GetCXArgFromArray(idxIdx)
			UpdateSymbolsTable(prgrm, symbols, idx, offset, true)
			GiveOffset(prgrm, symbols, idx, offset, true)
			checkIndexType(prgrm, idx)
		}
		for _, fldIdx := range arg.Fields {
			fld := prgrm.GetCXArgFromArray(fldIdx)
			for _, idxIdx := range fld.Indexes {
				idx := prgrm.GetCXArgFromArray(idxIdx)
				UpdateSymbolsTable(prgrm, symbols, idx, offset, true)
				GiveOffset(prgrm, symbols, idx, offset, true)
			}
		}

		SetFinalSize(prgrm, symbols, arg)

		AddPointer(prgrm, fn, arg)
	}
}

// isPointerAdded checks if `sym` has already been added to `fn.ListOfPointers`.
func isPointerAdded(prgrm *ast.CXProgram, fn *ast.CXFunction, sym *ast.CXArgument) (found bool) {
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
				prgrm.GetCXArgFromArray(sym.Fields[len(sym.Fields)-1]).Name == prgrm.GetCXArgFromArray(ptr.Fields[len(ptr.Fields)-1]).Name {
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
func AddPointer(prgrm *ast.CXProgram, fn *ast.CXFunction, sym *ast.CXArgument) {
	// Ignore if it's a global variable.
	if sym.Offset > prgrm.Stack.Size {
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
		fld := prgrm.GetCXArgFromArray(sym.Fields[len(sym.Fields)-1])
		if fld.IsPointer() && !isPointerAdded(prgrm, fn, sym) {
			fn.ListOfPointers = append(fn.ListOfPointers, sym)
		}
	}
	// Root symbol:
	// Checking if it is a pointer candidate and if it was already
	// added to the list.
	if sym.IsPointer() && !isPointerAdded(prgrm, fn, sym) {
		if len(sym.Fields) > 0 {
			tmp := ast.CXArgument{}
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
func CheckRedeclared(prgrm *ast.CXProgram, symbols *[]map[string]*ast.CXArgument, expr *ast.CXExpression, sym *ast.CXArgument) {
	cxAtomicOp, _, _, err := prgrm.GetOperation(expr)
	if err != nil {
		panic(err)
	}
	cxAtomicOpOperator := prgrm.GetFunctionFromArray(cxAtomicOp.Operator)

	if cxAtomicOpOperator == nil && len(cxAtomicOp.Outputs) > 0 && len(cxAtomicOp.Inputs) == 0 {
		lastIdx := len(*symbols) - 1

		symPkg, err := prgrm.GetPackageFromArray(sym.Package)
		if err != nil {
			panic(err)
		}

		_, found := (*symbols)[lastIdx][symPkg.Name+"."+sym.Name]
		if found {
			println(ast.CompilationError(sym.ArgDetails.FileName, sym.ArgDetails.FileLine), fmt.Sprintf("'%s' redeclared", sym.Name))
		}
	}
}

func ProcessLocalDeclaration(prgrm *ast.CXProgram, symbols *[]map[string]*ast.CXArgument, symbolsScope *map[string]bool, arg *ast.CXArgument) {
	argPkg, err := prgrm.GetPackageFromArray(arg.Package)
	if err != nil {
		panic(err)
	}

	if arg.IsLocalDeclaration {
		(*symbolsScope)[argPkg.Name+"."+arg.Name] = true
	}
	arg.IsLocalDeclaration = (*symbolsScope)[argPkg.Name+"."+arg.Name]
}

func ProcessGoTos(prgrm *ast.CXProgram, fn *ast.CXFunction, exprs []ast.CXExpression) {
	for i, expr := range exprs {
		cxAtomicOp, _, _, err := prgrm.GetOperation(&expr)
		if err != nil {
			panic(err)
		}
		cxAtomicOpOperator := prgrm.GetFunctionFromArray(cxAtomicOp.Operator)

		opGotoFn := ast.Natives[constants.OP_IDENTITY]
		if cxAtomicOpOperator != nil {
			opGotoFn.Index = cxAtomicOpOperator.Index
		}
		if cxAtomicOpOperator == opGotoFn {
			// then it's a goto
			for j, e := range exprs {
				ecxAtomicOp, _, _, err := prgrm.GetOperation(&e)
				if err != nil {
					panic(err)
				}

				if ecxAtomicOp.Label == cxAtomicOp.Label && i != j {
					// ElseLines is used because arg's default val is false
					cxAtomicOp.ThenLines = j - i - 1
					break
				}
			}
		}

		fn.AddExpression(prgrm, &expr)
	}
}

func checkMatchParamTypes(prgrm *ast.CXProgram, expr *ast.CXExpression, expected, received []*ast.CXArgument, isInputs bool) {
	cxAtomicOp, _, _, err := prgrm.GetOperation(expr)
	if err != nil {
		panic(err)
	}
	cxAtomicOpOperator := prgrm.GetFunctionFromArray(cxAtomicOp.Operator)

	for i, inp := range expected {
		expectedType := ast.GetFormattedType(prgrm, expected[i])
		receivedType := ast.GetFormattedType(prgrm, received[i])

		if expr.IsMethodCall() && expected[i].IsPointer() && i == 0 {
			// if method receiver is pointer, remove *
			if expectedType[0] == '*' {
				// we need to check if it's not an `str`
				// otherwise we end up removing the `s` instead of a `*`
				expectedType = expectedType[1:]
			}
		}

		if expectedType != receivedType && inp.Type != types.UNDEFINED {
			var opName string
			if cxAtomicOpOperator.IsBuiltIn() {
				opName = ast.OpNames[cxAtomicOpOperator.AtomicOPCode]
			} else {
				opName = cxAtomicOpOperator.Name
			}

			if isInputs {
				println(ast.CompilationError(received[i].ArgDetails.FileName, received[i].ArgDetails.FileLine), fmt.Sprintf("function '%s' expected input argument of type '%s'; '%s' was provided", opName, expectedType, receivedType))
			} else {
				println(ast.CompilationError(prgrm.GetCXArgFromArray(cxAtomicOp.Outputs[i]).ArgDetails.FileName, prgrm.GetCXArgFromArray(cxAtomicOp.Outputs[i]).ArgDetails.FileLine), fmt.Sprintf("function '%s' expected receiving variable of type '%s'; '%s' was provided", opName, expectedType, receivedType))
			}

		}

		// In the case of assignment we need to check that the input's type matches the output's type.
		// FIXME: There are some expressions added by the cxgo where temporary variables are used.
		// These temporary variables' types are not properly being set. That's why we use !cxcore.IsTempVar to
		// exclude these cases for now.
		if cxAtomicOpOperator.AtomicOPCode == constants.OP_IDENTITY && !IsTempVar(prgrm.GetCXArgFromArray(cxAtomicOp.Outputs[0]).Name) {
			inpType := ast.GetFormattedType(prgrm, prgrm.GetCXArgFromArray(cxAtomicOp.Inputs[0]))
			outType := ast.GetFormattedType(prgrm, prgrm.GetCXArgFromArray(cxAtomicOp.Outputs[0]))

			// We use `isInputs` to only print the error once.
			// Otherwise we'd print the error twice: once for the input and again for the output
			if inpType != outType && isInputs {
				println(ast.CompilationError(received[i].ArgDetails.FileName, received[i].ArgDetails.FileLine), fmt.Sprintf("cannot assign value of type '%s' to identifier '%s' of type '%s'", inpType, prgrm.GetCXArgFromArray(cxAtomicOp.Outputs[0]).GetAssignmentElement(prgrm).Name, outType))
			}
		}
	}
}

func CheckTypes(prgrm *ast.CXProgram, exprs []ast.CXExpression, currIndex int) {
	cxAtomicOp, err := prgrm.GetCXAtomicOpFromExpressions(exprs, currIndex)
	if err != nil {
		panic(err)
	}
	cxAtomicOpOperator := prgrm.GetFunctionFromArray(cxAtomicOp.Operator)

	exprCXLine, _ := prgrm.GetPreviousCXLine(exprs, currIndex)

	if cxAtomicOpOperator != nil {
		opName := cxAtomicOp.GetOperatorName(prgrm)

		// checking if number of inputs is less than the required number of inputs
		if len(cxAtomicOp.Inputs) != len(cxAtomicOpOperator.Inputs) {
			if !(len(cxAtomicOpOperator.Inputs) > 0 && prgrm.GetCXArgFromArray(cxAtomicOpOperator.Inputs[len(cxAtomicOpOperator.Inputs)-1]).Type != types.UNDEFINED) {
				// if the last input is of type cxcore.TYPE_UNDEFINED then it might be a variadic function, such as printf
			} else {
				// then we need to be strict in the number of inputs
				var plural1 string
				var plural2 string = "s"
				var plural3 string = "were"
				if len(cxAtomicOpOperator.Inputs) > 1 {
					plural1 = "s"
				}
				if len(cxAtomicOp.Inputs) == 1 {
					plural2 = ""
					plural3 = "was"
				}

				println(ast.CompilationError(exprCXLine.FileName, exprCXLine.LineNumber), fmt.Sprintf("operator '%s' expects %d input%s, but %d input argument%s %s provided", opName, len(cxAtomicOpOperator.Inputs), plural1, len(cxAtomicOp.Inputs), plural2, plural3))
				return
			}
		}

		// checking if number of expr.ProgramOutput matches number of Operator.ProgramOutput
		if len(cxAtomicOp.Outputs) != len(cxAtomicOpOperator.Outputs) {
			var plural1 string
			var plural2 string = "s"
			var plural3 string = "were"
			if len(cxAtomicOpOperator.Outputs) > 1 {
				plural1 = "s"
			}
			if len(cxAtomicOp.Outputs) == 1 {
				plural2 = ""
				plural3 = "was"
			}

			println(ast.CompilationError(exprCXLine.FileName, exprCXLine.LineNumber), fmt.Sprintf("operator '%s' expects to return %d output%s, but %d receiving argument%s %s provided", opName, len(cxAtomicOpOperator.Outputs), plural1, len(cxAtomicOp.Outputs), plural2, plural3))
			os.Exit(constants.CX_COMPILATION_ERROR)
		}
	}

	if cxAtomicOpOperator != nil && cxAtomicOpOperator.IsBuiltIn() && cxAtomicOpOperator.AtomicOPCode == constants.OP_IDENTITY {
		for i := range cxAtomicOp.Inputs {
			var expectedType string
			var receivedType string
			if prgrm.GetCXArgFromArray(cxAtomicOp.Outputs[i]).GetAssignmentElement(prgrm).StructType != nil {
				// then it's custom type
				expectedType = prgrm.GetCXArgFromArray(cxAtomicOp.Outputs[i]).GetAssignmentElement(prgrm).StructType.Name
			} else {
				// then it's native type
				expectedType = prgrm.GetCXArgFromArray(cxAtomicOp.Outputs[i]).GetAssignmentElement(prgrm).Type.Name()

				if prgrm.GetCXArgFromArray(cxAtomicOp.Outputs[i]).GetAssignmentElement(prgrm).Type == types.POINTER {
					expectedType = prgrm.GetCXArgFromArray(cxAtomicOp.Outputs[i]).GetAssignmentElement(prgrm).PointerTargetType.Name()
				}
			}

			if prgrm.GetCXArgFromArray(cxAtomicOp.Inputs[i]).GetAssignmentElement(prgrm).StructType != nil {
				// then it's custom type
				receivedType = prgrm.GetCXArgFromArray(cxAtomicOp.Inputs[i]).GetAssignmentElement(prgrm).StructType.Name
			} else {
				// then it's native type
				receivedType = prgrm.GetCXArgFromArray(cxAtomicOp.Inputs[i]).GetAssignmentElement(prgrm).Type.Name()

				if prgrm.GetCXArgFromArray(cxAtomicOp.Inputs[i]).GetAssignmentElement(prgrm).Type == types.POINTER {
					receivedType = prgrm.GetCXArgFromArray(cxAtomicOp.Inputs[i]).GetAssignmentElement(prgrm).PointerTargetType.Name()
				}
			}

			// if cxcore.GetAssignmentElement(exprs[currIndex].ProgramOutput[i]).Type != cxcore.GetAssignmentElement(inp).Type {
			if receivedType != expectedType {
				if exprs[currIndex].IsStructLiteral() {
					println(ast.CompilationError(prgrm.GetCXArgFromArray(cxAtomicOp.Outputs[i]).ArgDetails.FileName, prgrm.GetCXArgFromArray(cxAtomicOp.Outputs[i]).ArgDetails.FileLine), fmt.Sprintf("field '%s' in struct literal of type '%s' expected argument of type '%s'; '%s' was provided", prgrm.GetCXArgFromArray(prgrm.GetCXArgFromArray(cxAtomicOp.Outputs[i]).Fields[0]).Name, prgrm.GetCXArgFromArray(cxAtomicOp.Outputs[i]).StructType.Name, expectedType, receivedType))
				} else {
					println(ast.CompilationError(prgrm.GetCXArgFromArray(cxAtomicOp.Outputs[i]).ArgDetails.FileName, prgrm.GetCXArgFromArray(cxAtomicOp.Outputs[i]).ArgDetails.FileLine), fmt.Sprintf("trying to assign argument of type '%s' to symbol '%s' of type '%s'", receivedType, prgrm.GetCXArgFromArray(cxAtomicOp.Outputs[i]).GetAssignmentElement(prgrm).Name, expectedType))
				}
			}
		}
	}

	// then it's a function call and not a declaration
	if cxAtomicOpOperator != nil {
		// checking inputs matching operator's inputs
		checkMatchParamTypes(prgrm, &exprs[currIndex], prgrm.ConvertIndexArgsToPointerArgs(cxAtomicOpOperator.Inputs), prgrm.ConvertIndexArgsToPointerArgs(cxAtomicOp.Inputs), true)

		// checking outputs matching operator's outputs
		checkMatchParamTypes(prgrm, &exprs[currIndex], prgrm.ConvertIndexArgsToPointerArgs(cxAtomicOpOperator.Outputs), prgrm.ConvertIndexArgsToPointerArgs(cxAtomicOp.Outputs), false)
	}
}

func ProcessStringAssignment(prgrm *ast.CXProgram, expr *ast.CXExpression) {
	cxAtomicOp, _, _, err := prgrm.GetOperation(expr)
	if err != nil {
		panic(err)
	}
	cxAtomicOpOperator := prgrm.GetFunctionFromArray(cxAtomicOp.Operator)

	opIdentFn := ast.Natives[constants.OP_IDENTITY]
	if cxAtomicOpOperator != nil {
		opIdentFn.Index = cxAtomicOpOperator.Index
	}
	if cxAtomicOpOperator == opIdentFn {
		for i, outIdx := range cxAtomicOp.Outputs {
			out := prgrm.GetCXArgFromArray(outIdx)
			if len(cxAtomicOp.Inputs) > i {
				out = out.GetAssignmentElement(prgrm)
				inp := prgrm.GetCXArgFromArray(cxAtomicOp.Inputs[i]).GetAssignmentElement(prgrm)

				if (out.Type == types.STR || out.Type == types.AFF) && out.Name != "" &&
					(inp.Type == types.STR || inp.Type == types.AFF) && inp.Name != "" {
					out.PassBy = constants.PASSBY_VALUE
				}
			}
		}
	}
}

// ProcessReferenceAssignment checks if the reference of a symbol can be assigned to the expression's output.
// For example: `var foo i32; var bar i32; bar = &foo` is not valid.
func ProcessReferenceAssignment(prgrm *ast.CXProgram, expr *ast.CXExpression) {
	cxAtomicOp, _, _, err := prgrm.GetOperation(expr)
	if err != nil {
		panic(err)
	}

	for _, outIdx := range cxAtomicOp.Outputs {
		out := prgrm.GetCXArgFromArray(outIdx)
		elt := out.GetAssignmentElement(prgrm)
		if elt.PassBy == constants.PASSBY_REFERENCE &&
			!hasDeclSpec(elt, constants.DECL_POINTER) &&
			elt.PointerTargetType != types.STR && elt.Type != types.STR && !elt.IsSlice {
			println(ast.CompilationError(CurrentFile, LineNo), "invalid reference assignment", elt.Name)
		}
	}

}

func ProcessSlice(prgrm *ast.CXProgram, inp *ast.CXArgument) {
	var elt *ast.CXArgument

	if len(inp.Fields) > 0 {
		elt = prgrm.GetCXArgFromArray(inp.Fields[len(inp.Fields)-1])
	} else {
		elt = inp
	}

	if elt.IsSlice && len(elt.DereferenceOperations) > 0 && elt.DereferenceOperations[len(elt.DereferenceOperations)-1] == constants.DEREF_POINTER {
		elt.DereferenceOperations = elt.DereferenceOperations[:len(elt.DereferenceOperations)-1]
		return
	}

	if elt.IsSlice && len(elt.DereferenceOperations) > 0 && len(inp.Fields) == 0 {
		return
		// elt.DereferenceOperations = append([]int{cxcore.DEREF_POINTER}, elt.DereferenceOperations...)
	}
}

func ProcessSliceAssignment(prgrm *ast.CXProgram, expr *ast.CXExpression) {
	cxAtomicOp, _, _, err := prgrm.GetOperation(expr)
	if err != nil {
		panic(err)
	}
	cxAtomicOpOperator := prgrm.GetFunctionFromArray(cxAtomicOp.Operator)

	opIdentFn := ast.Natives[constants.OP_IDENTITY]
	if cxAtomicOpOperator != nil {
		opIdentFn.Index = cxAtomicOpOperator.Index
	}
	if cxAtomicOpOperator == opIdentFn {
		var inp *ast.CXArgument
		var out *ast.CXArgument

		inp = prgrm.GetCXArgFromArray(cxAtomicOp.Inputs[0]).GetAssignmentElement(prgrm)
		out = prgrm.GetCXArgFromArray(cxAtomicOp.Outputs[0]).GetAssignmentElement(prgrm)

		if inp.IsSlice && out.IsSlice && len(inp.Indexes) == 0 && len(out.Indexes) == 0 {
			out.PassBy = constants.PASSBY_VALUE
		}
	}
	if cxAtomicOpOperator != nil && !cxAtomicOpOperator.IsBuiltIn() {
		// then it's a function call
		for _, inpIdx := range cxAtomicOp.Inputs {
			inp := prgrm.GetCXArgFromArray(inpIdx)
			assignElt := inp.GetAssignmentElement(prgrm)

			// we want to pass by value if we're sending the slice as a whole (no indexing)
			// unless it's a pointer to the slice
			if assignElt.IsSlice && len(assignElt.Indexes) == 0 && !hasDeclSpec(assignElt, constants.DECL_POINTER) {
				assignElt.PassBy = constants.PASSBY_VALUE
			}
		}
	}
}

// lookupSymbol searches for `ident` in `symbols`, starting from the innermost scope.
func lookupSymbol(prgrm *ast.CXProgram, pkgName, ident string, symbols *[]map[string]*ast.CXArgument) (*ast.CXArgument, error) {
	fullName := pkgName + "." + ident
	for c := len(*symbols) - 1; c >= 0; c-- {
		if sym, found := (*symbols)[c][fullName]; found {
			return sym, nil
		}
	}

	// Checking if `ident` refers to a function.
	pkg, err := prgrm.GetPackage(pkgName)
	if err != nil {
		return nil, err
	}

	notFound := errors.New("identifier '" + ident + "' does not exist")

	// We're not checking for that error
	fn, err := pkg.GetFunction(prgrm, ident)
	if err != nil {
		return nil, errors.New(err.Error() + ":" + notFound.Error() + fmt.Sprintf("--fullName=%s", fullName))
	}
	// Then we found a function by that name. Let's create a `cxcore.CXArgument` of
	// type `func` with that name.
	fnArg := ast.MakeArgument(ident, fn.FileName, fn.FileLine).AddType(types.FUNC)
	fnArg.Package = ast.CXPackageIndex(pkg.Index)

	return fnArg, nil
}

// UpdateSymbolsTable adds `sym` to the innermost scope (last element of slice) in `symbols`.
func UpdateSymbolsTable(prgrm *ast.CXProgram, symbols *[]map[string]*ast.CXArgument, sym *ast.CXArgument, offset *types.Pointer, shouldExist bool) {
	if sym.Name != "" {
		symPkg, err := prgrm.GetPackageFromArray(sym.Package)
		if err != nil {
			panic(err)
		}

		if !sym.IsLocalDeclaration {
			GetGlobalSymbol(prgrm, symbols, symPkg, sym.Name)
		}

		lastIdx := len(*symbols) - 1
		fullName := symPkg.Name + "." + sym.Name

		// outerSym, err := lookupSymbol(sym.Package.Name, sym.Name, symbols)
		_, err = lookupSymbol(prgrm, symPkg.Name, sym.Name, symbols)
		_, found := (*symbols)[lastIdx][fullName]

		// then it wasn't found in any scope
		if err != nil && shouldExist {
			println(ast.CompilationError(sym.ArgDetails.FileName, sym.ArgDetails.FileLine), "identifier '"+sym.Name+"' does not exist")
		}

		// then it was already added in the innermost scope
		if found {
			return
		}

		// then it is a new declaration
		if !shouldExist && !found {
			// *symbols = append(*symbols, make(map[string]*ast.CXArgument))
			// lastIdx = len(*symbols) - 1

			// then it was declared in an outer scope
			sym.Offset = *offset
			(*symbols)[lastIdx][fullName] = sym
			*offset += ast.GetSize(prgrm, sym)
		}
	}
}

func ProcessMethodCall(prgrm *ast.CXProgram, expr *ast.CXExpression, symbols *[]map[string]*ast.CXArgument, offset *types.Pointer, shouldExist bool) {
	cxAtomicOp, _, _, err := prgrm.GetOperation(expr)
	if err != nil {
		panic(err)
	}

	if expr.IsMethodCall() {
		var inp *ast.CXArgument
		var out *ast.CXArgument

		if len(cxAtomicOp.Inputs) > 0 && prgrm.GetCXArgFromArray(cxAtomicOp.Inputs[0]).Name != "" {
			inp = prgrm.GetCXArgFromArray(cxAtomicOp.Inputs[0])
		}
		if len(cxAtomicOp.Outputs) > 0 && prgrm.GetCXArgFromArray(cxAtomicOp.Outputs[0]).Name != "" {
			out = prgrm.GetCXArgFromArray(cxAtomicOp.Outputs[0])
		}

		if inp != nil {
			inpPkg, err := prgrm.GetPackageFromArray(inp.Package)
			if err != nil {
				panic(err)
			}
			// if argInp, found := (*symbols)[lastIdx][inp.Package.Name+"."+inp.Name]; !found {
			if argInp, err := lookupSymbol(prgrm, inpPkg.Name, inp.Name, symbols); err != nil {
				if out == nil {
					panic("")
				}

				outPkg, err := prgrm.GetPackageFromArray(out.Package)
				if err != nil {
					panic(err)
				}

				argOut, err := lookupSymbol(prgrm, outPkg.Name, out.Name, symbols)
				if err != nil {
					println(ast.CompilationError(out.ArgDetails.FileName, out.ArgDetails.FileLine), fmt.Sprintf("identifier '%s' does not exist", out.Name))
					os.Exit(constants.CX_COMPILATION_ERROR)
				}
				// then we found an output
				if len(out.Fields) > 0 {
					strct := argOut.StructType
					strctPkg, err := prgrm.GetPackageFromArray(strct.Package)
					if err != nil {
						panic(err)
					}

					if fn, err := strctPkg.GetMethod(prgrm, strct.Name+"."+prgrm.GetCXArgFromArray(out.Fields[len(out.Fields)-1]).Name, strct.Name); err == nil {
						cxAtomicOp.Operator = ast.CXFunctionIndex(fn.Index)
					} else {
						panic("")
					}

					cxAtomicOp.Inputs = append([]ast.CXArgumentIndex{ast.CXArgumentIndex(out.Index)}, cxAtomicOp.Inputs...)

					out.Fields = out.Fields[:len(out.Fields)-1]

					cxAtomicOp.Outputs = cxAtomicOp.Outputs[1:]
				}
			} else {
				// then we found an input
				if len(inp.Fields) > 0 {
					strct := argInp.StructType

					for _, fldIdx := range inp.Fields {
						fld := prgrm.GetCXArgFromArray(fldIdx)
						if inFld, err := strct.GetField(fld.Name); err == nil {
							if inFld.StructType != nil {
								strct = inFld.StructType
							}
						}
					}

					strctPkg, err := prgrm.GetPackageFromArray(strct.Package)
					if err != nil {
						panic(err)
					}

					if fn, err := strctPkg.GetMethod(prgrm, strct.Name+"."+prgrm.GetCXArgFromArray(inp.Fields[len(inp.Fields)-1]).Name, strct.Name); err == nil {
						cxAtomicOp.Operator = ast.CXFunctionIndex(fn.Index)
					} else {
						panic(err)
					}

					inp.Fields = inp.Fields[:len(inp.Fields)-1]
				} else if len(out.Fields) > 0 {
					outPkg, err := prgrm.GetPackageFromArray(out.Package)
					if err != nil {
						panic(err)
					}
					argOut, err := lookupSymbol(prgrm, outPkg.Name, out.Name, symbols)
					if err != nil {
						panic(err)
					}

					strct := argOut.StructType

					if strct == nil {
						println(ast.CompilationError(argOut.ArgDetails.FileName, argOut.ArgDetails.FileLine), fmt.Sprintf("illegal method call or field access on identifier '%s' of primitive type '%s'", argOut.Name, argOut.Type.Name()))
						os.Exit(constants.CX_COMPILATION_ERROR)
					}

					cxAtomicOp.Inputs = append(cxAtomicOp.Outputs[:1], cxAtomicOp.Inputs...)

					cxAtomicOp.Outputs = cxAtomicOp.Outputs[:len(cxAtomicOp.Outputs)-1]

					strctPkg, err := prgrm.GetPackageFromArray(strct.Package)
					if err != nil {
						panic(err)
					}

					if fn, err := strctPkg.GetMethod(prgrm, strct.Name+"."+prgrm.GetCXArgFromArray(out.Fields[len(out.Fields)-1]).Name, strct.Name); err == nil {
						cxAtomicOp.Operator = ast.CXFunctionIndex(fn.Index)
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

			outPkg, err := prgrm.GetPackageFromArray(out.Package)
			if err != nil {
				panic(err)
			}

			argOut, err := lookupSymbol(prgrm, outPkg.Name, out.Name, symbols)
			if err != nil {
				println(ast.CompilationError(out.ArgDetails.FileName, out.ArgDetails.FileLine), fmt.Sprintf("identifier '%s' does not exist", out.Name))
				os.Exit(constants.CX_COMPILATION_ERROR)
			}

			// then we found an output
			if len(out.Fields) > 0 {
				strct := argOut.StructType

				if strct == nil {
					println(ast.CompilationError(argOut.ArgDetails.FileName, argOut.ArgDetails.FileLine), fmt.Sprintf("illegal method call or field access on identifier '%s' of primitive type '%s'", argOut.Name, argOut.Type.Name()))
					os.Exit(constants.CX_COMPILATION_ERROR)
				}

				strctPkg, err := prgrm.GetPackageFromArray(strct.Package)
				if err != nil {
					panic(err)
				}

				if fn, err := strctPkg.GetMethod(prgrm, strct.Name+"."+prgrm.GetCXArgFromArray(out.Fields[len(out.Fields)-1]).Name, strct.Name); err == nil {
					cxAtomicOp.Operator = ast.CXFunctionIndex(fn.Index)
				} else {
					panic("")
				}

				cxAtomicOp.Inputs = append([]ast.CXArgumentIndex{ast.CXArgumentIndex(out.Index)}, cxAtomicOp.Inputs...)

				out.Fields = out.Fields[:len(out.Fields)-1]

				cxAtomicOp.Outputs = cxAtomicOp.Outputs[1:]
				// expr.ProgramOutput = nil
			}
		}
		cxAtomicOpOperator := prgrm.GetFunctionFromArray(cxAtomicOp.Operator)

		// checking if receiver is sent as pointer or not
		if prgrm.GetCXArgFromArray(cxAtomicOpOperator.Inputs[0]).IsPointer() {
			prgrm.GetCXArgFromArray(cxAtomicOp.Inputs[0]).PassBy = constants.PASSBY_REFERENCE
		}
	}
}

func GiveOffset(prgrm *ast.CXProgram, symbols *[]map[string]*ast.CXArgument, sym *ast.CXArgument, offset *types.Pointer, shouldExist bool) {
	if sym.Name != "" {
		symPkg, err := prgrm.GetPackageFromArray(sym.Package)
		if err != nil {
			panic(err)
		}

		if !sym.IsLocalDeclaration {
			GetGlobalSymbol(prgrm, symbols, symPkg, sym.Name)
		}

		arg, err := lookupSymbol(prgrm, symPkg.Name, sym.Name, symbols)
		if err == nil {
			ProcessSymbolFields(prgrm, sym, arg)
			CopyArgFields(prgrm, sym, arg)
		}
	}
}

func ProcessTempVariable(prgrm *ast.CXProgram, expr *ast.CXExpression) {
	cxAtomicOp, _, _, err := prgrm.GetOperation(expr)
	if err != nil {
		panic(err)
	}
	cxAtomicOpOperator := prgrm.GetFunctionFromArray(cxAtomicOp.Operator)

	opIdentFn := ast.Natives[constants.OP_IDENTITY]
	if cxAtomicOpOperator != nil {
		opIdentFn.Index = cxAtomicOpOperator.Index
	}
	if cxAtomicOpOperator != nil && (cxAtomicOpOperator == opIdentFn || ast.IsArithmeticOperator(cxAtomicOpOperator.AtomicOPCode)) && len(cxAtomicOp.Outputs) > 0 && len(cxAtomicOp.Inputs) > 0 {
		arg := prgrm.GetCXArgFromArray(cxAtomicOp.Outputs[0])
		name := arg.Name
		if IsTempVar(name) {
			cxAtomicOpInput := prgrm.GetCXArgFromArray(cxAtomicOp.Inputs[0])
			// then it's a temporary variable and it needs to adopt its input's type
			arg.Type = cxAtomicOpInput.Type
			arg.PointerTargetType = cxAtomicOpInput.PointerTargetType
			arg.Size = cxAtomicOpInput.Size
			arg.TotalSize = cxAtomicOpInput.TotalSize
			arg.PreviouslyDeclared = true
		}
	}
}

func CopyArgFields(prgrm *ast.CXProgram, sym *ast.CXArgument, arg *ast.CXArgument) {
	sym.Offset = arg.Offset
	sym.Type = arg.Type

	// sym.IndirectionLevels = arg.IndirectionLevels

	if sym.ArgDetails.FileLine != arg.ArgDetails.FileLine {
		// FIXME Maybe we can unify this later.
		if len(sym.Fields) > 0 {
			elt := sym.GetAssignmentElement(prgrm)

			declSpec := []int{}
			for c := 0; c < len(elt.DeclarationSpecifiers); c++ {
				switch elt.DeclarationSpecifiers[c] {
				case constants.DECL_INDEXING:
					if declSpec[len(declSpec)-1] == constants.DECL_ARRAY || declSpec[len(declSpec)-1] == constants.DECL_SLICE {
						declSpec = declSpec[:len(declSpec)-1]
					} else {
						println(ast.CompilationError(sym.ArgDetails.FileName, sym.ArgDetails.FileLine), "invalid indexing")
					}
				case constants.DECL_DEREF:
					if declSpec[len(declSpec)-1] == constants.DECL_POINTER {
						declSpec = declSpec[:len(declSpec)-1]
					} else {
						println(ast.CompilationError(sym.ArgDetails.FileName, sym.ArgDetails.FileLine), "invalid indirection")
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
				// checking if we need to remove or add cxcore.DECL_POINTERs
				// also we could be removing
				switch spec {
				case constants.DECL_INDEXING:
					if declSpec[len(declSpec)-1] == constants.DECL_ARRAY || declSpec[len(declSpec)-1] == constants.DECL_SLICE {
						declSpec = declSpec[:len(declSpec)-1]
					} else {
						println(ast.CompilationError(sym.ArgDetails.FileName, sym.ArgDetails.FileLine), "invalid indexing")
					}
				case constants.DECL_DEREF:
					if declSpec[len(declSpec)-1] == constants.DECL_POINTER {
						declSpec = declSpec[:len(declSpec)-1]
					} else {
						println(ast.CompilationError(sym.ArgDetails.FileName, sym.ArgDetails.FileLine), "invalid indirection")
					}
				case constants.DECL_POINTER:
					if sym.ArgDetails.FileLine != arg.ArgDetails.FileLine {
						// This function is also called so it assigns offset and other fields to signature parameters
						//
						declSpec = append(declSpec, constants.DECL_POINTER)
					}
				}
			}

			sym.DeclarationSpecifiers = declSpec
		}
	} else {
		sym.DeclarationSpecifiers = arg.DeclarationSpecifiers
	}

	sym.IsSlice = arg.IsSlice
	sym.StructType = arg.StructType

	// FIXME: In other processes like ProcessSymbolFields the symbol is assigned with lengths.
	// If we already have some lengths, we skip this. This needs to be fixed in the redesign of the cxgo.
	if len(sym.Lengths) == 0 {
		sym.Lengths = arg.Lengths
	}

	// sym.Lengths = arg.Lengths
	sym.Package = arg.Package
	sym.DoesEscape = arg.DoesEscape
	sym.Size = arg.Size

	// Checking if it's a slice struct field. We'll do the same process as
	// below (as in the `arg.IsSlice` check), but the process differs in the
	// case of a slice struct field.
	elt := sym.GetAssignmentElement(prgrm)

	if (!arg.IsSlice || hasDerefOp(sym, constants.DEREF_ARRAY)) && arg.StructType != nil && elt.IsSlice && elt != sym {
		for i, deref := range elt.DereferenceOperations {
			// The cxgo when reading `foo[5]` in postfix.go does not know if `foo`
			// is a slice or an array. At this point we now know it's a slice and we need
			// to change those dereferences to cxcore.DEREF_SLICE.
			if deref == constants.DEREF_ARRAY {
				elt.DereferenceOperations[i] = constants.DEREF_SLICE
			}
		}

		if len(elt.DereferenceOperations) > 0 && elt.DereferenceOperations[0] == constants.DEREF_POINTER {
			elt.DereferenceOperations = elt.DereferenceOperations[1:]
		}
	}

	if arg.IsSlice {
		if !hasDerefOp(sym, constants.DEREF_ARRAY) {
			// Then we're handling the slice itself, and we need to dereference it.
			sym.DereferenceOperations = append([]int{constants.DEREF_POINTER}, sym.DereferenceOperations...)
		} else {
			for i, deref := range sym.DereferenceOperations {
				// The cxgo when reading `foo[5]` in postfix.go does not know if `foo`
				// is a slice or an array. At this point we now know it's a slice and we need
				// to change those dereferences to cxcore.DEREF_SLICE.
				if deref == constants.DEREF_ARRAY {
					sym.DereferenceOperations[i] = constants.DEREF_SLICE
				}
			}
		}
	}

	if len(sym.Fields) > 0 {
		symField := prgrm.GetCXArgFromArray(sym.Fields[len(sym.Fields)-1])
		if sym.Type == types.POINTER && arg.StructType != nil {
			sym.PointerTargetType = symField.Type
		} else {
			sym.Type = symField.Type
			sym.PointerTargetType = symField.PointerTargetType
		}

		sym.IsSlice = symField.IsSlice
	} else {
		sym.Type = arg.Type
		sym.PointerTargetType = arg.PointerTargetType
	}

	if !arg.IsStruct {
		sym.TotalSize = arg.TotalSize
	} else {
		if len(sym.Fields) > 0 {
			sym.TotalSize = prgrm.GetCXArgFromArray(sym.Fields[len(sym.Fields)-1]).TotalSize
		} else {
			sym.TotalSize = arg.TotalSize
		}
	}
}

func ProcessSymbolFields(prgrm *ast.CXProgram, sym *ast.CXArgument, arg *ast.CXArgument) {
	if len(sym.Fields) > 0 {
		if arg.StructType == nil || len(arg.StructType.Fields) == 0 {
			println(ast.CompilationError(sym.ArgDetails.FileName, sym.ArgDetails.FileLine), fmt.Sprintf("'%s' has no fields", sym.Name))
			return
		}

		// checking if fields do exist in their StructType
		// and assigning that StructType to the sym.Field
		strct := arg.StructType
		strctPkg, err := prgrm.GetPackageFromArray(strct.Package)
		if err != nil {
			panic(err)
		}

		for _, fldIdx := range sym.Fields {
			fld := prgrm.GetCXArgFromArray(fldIdx)
			if inFld, err := strct.GetField(fld.Name); err == nil {
				if inFld.StructType != nil {
					fld.StructType = strct
					strct = inFld.StructType
				}
			} else {
				methodName := prgrm.GetCXArgFromArray(sym.Fields[len(sym.Fields)-1]).Name
				receiverType := strct.Name

				if method, methodErr := strctPkg.GetMethod(prgrm, receiverType+"."+methodName, receiverType); methodErr == nil {
					fld.Type = prgrm.GetCXArgFromArray(method.Outputs[0]).Type
					fld.PointerTargetType = prgrm.GetCXArgFromArray(method.Outputs[0]).PointerTargetType
				} else {
					println(ast.CompilationError(fld.ArgDetails.FileName, fld.ArgDetails.FileLine), err.Error())
				}

			}
		}

		strct = arg.StructType
		// then we copy all the type struct fields
		// to the respective sym.Fields
		for _, nameFldIdx := range sym.Fields {
			nameFld := prgrm.GetCXArgFromArray(nameFldIdx)
			if nameFld.StructType != nil {
				strct = nameFld.StructType
			}

			for _, fld := range strct.Fields {
				if nameFld.Name == fld.Name {
					nameFld.Type = fld.Type
					nameFld.Lengths = fld.Lengths
					nameFld.Size = fld.Size
					nameFld.TotalSize = fld.TotalSize
					nameFld.DereferenceLevels = sym.DereferenceLevels
					nameFld.PointerTargetType = fld.PointerTargetType
					nameFld.StructType = fld.StructType

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
						nameFld.DereferenceOperations = append([]int{constants.DEREF_POINTER}, nameFld.DereferenceOperations...)
						nameFld.DereferenceLevels++
					}

					nameFld.PassBy = fld.PassBy
					nameFld.IsSlice = fld.IsSlice

					if fld.Type == types.STR || fld.Type == types.AFF {
						nameFld.PassBy = constants.PASSBY_REFERENCE
						// nameFld.Size = cxcore.POINTER_SIZE
						// nameFld.TotalSize = cxcore.POINTER_SIZE
					}

					if fld.StructType != nil {
						strct = fld.StructType
					}
					break
				}

				nameFld.Offset += ast.GetSize(prgrm, fld)
			}
		}
	}
}

func SetFinalSize(prgrm *ast.CXProgram, symbols *[]map[string]*ast.CXArgument, sym *ast.CXArgument) {
	finalSize := sym.TotalSize

	symPkg, err := prgrm.GetPackageFromArray(sym.Package)
	if err != nil {
		panic(err)
	}

	arg, err := lookupSymbol(prgrm, symPkg.Name, sym.Name, symbols)
	if err == nil {
		PreFinalSize(prgrm, &finalSize, sym, arg)
		for _, fldIdx := range sym.Fields {
			fld := prgrm.GetCXArgFromArray(fldIdx)
			finalSize = fld.TotalSize
			PreFinalSize(prgrm, &finalSize, fld, arg)
		}
	}

	sym.TotalSize = finalSize
}

// GetGlobalSymbol tries to retrieve `ident` from `symPkg`'s globals if `ident` is not found in the local scope.
func GetGlobalSymbol(prgrm *ast.CXProgram, symbols *[]map[string]*ast.CXArgument, symPkg *ast.CXPackage, ident string) {
	_, err := lookupSymbol(prgrm, symPkg.Name, ident, symbols)
	if err != nil {
		if glbl, err := symPkg.GetGlobal(prgrm, ident); err == nil {
			lastIdx := len(*symbols) - 1
			(*symbols)[lastIdx][symPkg.Name+"."+ident] = glbl
		}
	}
}

func PreFinalSize(prgrm *ast.CXProgram, finalSize *types.Pointer, sym *ast.CXArgument, arg *ast.CXArgument) {
	idxCounter := 0
	elt := sym.GetAssignmentElement(prgrm)
	for _, op := range elt.DereferenceOperations {
		if elt.IsSlice {
			continue
		}
		switch op {
		case constants.DEREF_ARRAY:
			*finalSize /= elt.Lengths[idxCounter]
			idxCounter++
		case constants.DEREF_POINTER:
			if len(arg.DeclarationSpecifiers) > 0 {
				subSize := types.Pointer(1)
				for _, decl := range arg.DeclarationSpecifiers {
					switch decl {
					case constants.DECL_ARRAY:
						for _, len := range arg.Lengths {
							subSize *= len
						}
					// case cxcore.DECL_SLICE:
					// 	subSize = POINTER_SIZE
					case constants.DECL_BASIC:
						subSize = sym.Type.Size()
					case constants.DECL_STRUCT:
						subSize = arg.StructType.Size
					}
				}

				*finalSize = subSize
			}
		}
	}
}
