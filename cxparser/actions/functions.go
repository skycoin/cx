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

// FunctionHeader takes a function name ('fnName') and either creates the
// function if it's not known before or returns the already existing function
// if it is.
//
// Input arguments description:
// 	prgrm - a CXProgram that contains all the data and arrays of the program.
// 	fnName - name of the function.
//  receiver - the receiver of the method if the function is a method.
// 			   i.e. func (sampleArg i32) SomeMethod(), sampleArg is the receiver.
//  isMethod - If the function is a method (isMethod = true), then it adds the object that
// 			   it's called on as the first argument.
func FunctionHeader(prgrm *ast.CXProgram, fnName string, receiver []*ast.CXArgument, isMethod bool) ast.CXFunctionIndex {
	if isMethod {
		if len(receiver) > 1 {
			panic("method has multiple receivers")
		}
		if pkg, err := prgrm.GetCurrentPackage(); err == nil {
			fnName := receiver[0].StructType.Name + "." + fnName

			if fn, err := prgrm.GetFunction(fnName, pkg.Name); err == nil {
				fn.AddInput(prgrm, receiver[0])
				pkg.CurrentFunction = ast.CXFunctionIndex(fn.Index)
				return ast.CXFunctionIndex(fn.Index)
			} else {
				fn := ast.MakeFunction(fnName, CurrentFile, LineNo)
				fn.AddInput(prgrm, receiver[0])
				_, fnIdx := pkg.AddFunction(prgrm, fn)

				return fnIdx
			}
		} else {
			panic(err)
		}
	} else {
		if pkg, err := prgrm.GetCurrentPackage(); err == nil {
			if fn, err := prgrm.GetFunction(fnName, pkg.Name); err == nil {
				pkg.CurrentFunction = ast.CXFunctionIndex(fn.Index)
				return ast.CXFunctionIndex(fn.Index)
			} else {
				fn := ast.MakeFunction(fnName, CurrentFile, LineNo)
				_, fnIdx := pkg.AddFunction(prgrm, fn)

				return fnIdx
			}
		} else {
			panic(err)
		}
	}
}

// FunctionAddParameters adds the function's input and output parameters.
//
// Input arguments description:
// 	prgrm - a CXProgram that contains all the data and arrays of the program.
// 	fnIdx - the index of the function in the main CXFunction array.
//  inputs - parameters to be added in the function as input.
//  outputs - parameters to be added in the function as output.
func FunctionAddParameters(prgrm *ast.CXProgram, fnIdx ast.CXFunctionIndex, inputs, outputs []*ast.CXArgument) {
	fn := prgrm.GetFunctionFromArray(fnIdx)
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
}

func isParseOp(prgrm *ast.CXProgram, expr *ast.CXExpression) bool {
	exprAtomicOp, err := prgrm.GetCXAtomicOp(expr.Index)
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
	expressionIdx := expr.Index
	expressionOperator := prgrm.GetFunctionFromArray(prgrm.CXAtomicOps[expressionIdx].Operator)
	if expressionOperator != nil && ast.IsOperator(expressionOperator.AtomicOPCode) && !IsAllArgsBasicTypes(prgrm, expr) {
		println(ast.CompilationError(CurrentFile, LineNo), fmt.Sprintf("invalid argument types for '%s' operator", ast.OpNames[expressionOperator.AtomicOPCode]))
	}
}

// ProcessFunctionParameters processes the scoping, offsets, and final size
// of the parameter args of the function.
//
// Input arguments description:
// 	prgrm - a CXProgram that contains all the data and arrays of the program.
// 	symbols - a slice of string-CXArg map which corresponds to
// 			  the scoping of the  CXArguments. Each element in the slice
// 			  corresponds to a different scope. The innermost scope is
// 			  the last element of the slice.
//  symbolsScope - only handles the difference between local and global
// 				   scopes, local being function constrained variables,
// 				   and global being global variables.
//  offset - offset to use by statements (excluding inputs, outputs
// 			 and receiver).
//  fnIdx - the index of the function in the main CXFunction array.
//  params - function parameters to be processed.
func ProcessFunctionParameters(prgrm *ast.CXProgram, symbols *[]map[string]*ast.CXArgument, symbolsScope *map[string]bool, offset *types.Pointer, fnIdx ast.CXFunctionIndex, params []ast.CXArgumentIndex) {
	fn := prgrm.GetFunctionFromArray(fnIdx)

	for _, paramIdx := range params {
		ProcessLocalDeclaration(prgrm, symbolsScope, paramIdx)

		UpdateSymbolsTable(prgrm, symbols, paramIdx, offset, false)
		GiveOffset(prgrm, symbols, paramIdx, false)
		SetFinalSize(prgrm, symbols, paramIdx)

		AddPointer(prgrm, fn, paramIdx)

		param := prgrm.GetCXArgFromArray(paramIdx)
		// as these are declarations, they should not have any dereference operations
		param.DereferenceOperations = nil
	}
}

// FunctionDeclaration completes the function's composition.
//
// Input arguments description:
// 	prgrm - a CXProgram that contains all the data and arrays of the program.
// 	fnIdx - the index of the function in the main CXFunction array.
//  inputs - parameters to be added in the function as input.
//  outputs - parameters to be added in the function as output.
//  exprs - expression statements inside the function.
func FunctionDeclaration(prgrm *ast.CXProgram, fnIdx ast.CXFunctionIndex, inputs, outputs []*ast.CXArgument, exprs []ast.CXExpression) {
	// getting offset to use by statements (excluding inputs, outputs and receiver).
	var offset types.Pointer

	// symbols is a slice of string-CXArg map which corresponds to
	// the scoping of the  CXArguments. Each element in the slice
	// corresponds to a different scope. The innermost scope is
	// the last element of the slice.
	var symbols *[]map[string]*ast.CXArgument
	tmp := make([]map[string]*ast.CXArgument, 0)
	symbols = &tmp
	*symbols = append(*symbols, make(map[string]*ast.CXArgument))

	// symbolsScope only handles the difference between local and global scopes
	// local being function constrained variables, and global being global variables.
	var symbolsScope map[string]bool = make(map[string]bool)

	fn := prgrm.GetFunctionFromArray(fnIdx)

	FunctionAddParameters(prgrm, fnIdx, inputs, outputs)
	ProcessGoTos(prgrm, exprs)
	AddExprsToFunction(prgrm, fnIdx, exprs)

	ProcessFunctionParameters(prgrm, symbols, &symbolsScope, &offset, fnIdx, fn.Inputs)
	ProcessFunctionParameters(prgrm, symbols, &symbolsScope, &offset, fnIdx, fn.Outputs)

	for i, expr := range fn.Expressions {
		if expr.Type == ast.CX_LINE {
			continue
		}
		exprAtomicOp, err := prgrm.GetCXAtomicOp(expr.Index)
		if err != nil {
			panic(err)
		}

		if expr.IsScopeNew() {
			*symbols = append(*symbols, make(map[string]*ast.CXArgument))
		}

		ProcessMethodCall(prgrm, &expr, symbols)
		ProcessExpressionArguments(prgrm, symbols, &symbolsScope, &offset, fnIdx, exprAtomicOp.Inputs, &expr, true)
		ProcessExpressionArguments(prgrm, symbols, &symbolsScope, &offset, fnIdx, exprAtomicOp.Outputs, &expr, false)

		ProcessPointerStructs(prgrm, &expr)
		ProcessTempVariable(prgrm, &expr)
		ProcessSliceAssignment(prgrm, &expr)
		ProcessStringAssignment(prgrm, &expr)
		ProcessReferenceAssignment(prgrm, &expr)
		ProcessShortDeclaration(prgrm, &expr, fn.Expressions, i)
		processTestExpression(prgrm, &expr)

		CheckTypes(prgrm, fn.Expressions, i)
		CheckUndValidTypes(prgrm, &expr)
		ProcessTypedOperator(prgrm, &expr)
		if expr.IsScopeDel() {
			*symbols = (*symbols)[:len(*symbols)-1]
		}
	}

	fn.LineCount = len(fn.Expressions)
	fn.Size = offset
}

// ProcessTypedOperator gets the proper typed operator for the expression.
// i.e. add -> i32.add
//
// Input arguments description:
// 	prgrm - a CXProgram that contains all the data and arrays of the program.
// 	expr - the expression.
func ProcessTypedOperator(prgrm *ast.CXProgram, expr *ast.CXExpression) {
	expression, err := prgrm.GetCXAtomicOp(expr.Index)
	if err != nil {
		panic(err)
	}

	expressionOperator := prgrm.GetFunctionFromArray(expression.Operator)
	if expressionOperator == nil {
		return
	}
	if ast.IsOperator(expressionOperator.AtomicOPCode) {
		atomicType := prgrm.CXArgs[expression.Inputs[0]].GetType(prgrm)
		typedOp := ast.GetTypedOperator(atomicType, expressionOperator.AtomicOPCode)
		if typedOp == nil {
			return
		}
		expressionOperator.AtomicOPCode = typedOp.AtomicOPCode
	}
}

// FunctionCall completes the array of expressions for postfix function call.
//
// Input arguments description:
// 	prgrm - a CXProgram that contains all the data and arrays of the program.
// 	exprs - the array of expressions.
//  args - the input arg expressions to the function.
func FunctionCall(prgrm *ast.CXProgram, exprs []ast.CXExpression, args []ast.CXExpression) []ast.CXExpression {
	expr := exprs[len(exprs)-1]

	expression, err := prgrm.GetCXAtomicOp(expr.Index)
	if err != nil {
		panic(err)
	}
	expressionOperator := prgrm.GetFunctionFromArray(expression.Operator)
	if expressionOperator == nil {
		cxAtomicOpOutput := prgrm.GetCXArgFromArray(expression.Outputs[0])
		opName := cxAtomicOpOutput.Name
		opPkgIdx := cxAtomicOpOutput.Package

		opPkg, err := prgrm.GetPackageFromArray(opPkgIdx)
		if err != nil {
			panic(err)
		}

		if op, err := prgrm.GetFunction(opName, opPkg.Name); err == nil {
			expression.Operator = ast.CXFunctionIndex(op.Index)
		} else if cxAtomicOpOutput.Fields == nil {
			// then it's not a possible method call
			println(ast.CompilationError(CurrentFile, LineNo), err.Error())
			return nil
		} else {
			exprs[len(exprs)-1].ExpressionType = ast.CXEXPR_METHOD_CALL
		}

		if len(expression.Outputs) > 0 && cxAtomicOpOutput.Fields == nil {
			expression.Outputs = nil
		}
	}

	var nestedExprs []ast.CXExpression
	for i, inpExpr := range args {
		if inpExpr.Type == ast.CX_LINE {
			continue
		}

		inpExprAtomicOp, err := prgrm.GetCXAtomicOp(inpExpr.Index)
		if err != nil {
			panic(err)
		}
		inpExprAtomicOpOperator := prgrm.GetFunctionFromArray(inpExprAtomicOp.Operator)

		inpExprCXLine, _ := prgrm.GetPreviousCXLine(args, i)

		if inpExprAtomicOpOperator == nil {
			// then it's a literal
			expression.AddInput(prgrm, inpExprAtomicOp.Outputs[0])
		} else {
			// then it's a function call
			if len(inpExprAtomicOp.Outputs) < 1 {
				var out *ast.CXArgument

				inpExprAtomicOpOperatorOutput := prgrm.GetCXArgFromArray(inpExprAtomicOpOperator.Outputs[0])
				if inpExprAtomicOpOperatorOutput.Type == types.UNDEFINED {
					// if undefined type, then adopt argument's type
					inpExprAtomicOpInput := prgrm.GetCXArgFromArray(inpExprAtomicOp.Inputs[0])

					out = ast.MakeArgument(generateTempVarName(constants.LOCAL_PREFIX), CurrentFile, inpExprCXLine.LineNumber).SetType(inpExprAtomicOpInput.Type)
					out.StructType = inpExprAtomicOpInput.StructType
					out.Size = inpExprAtomicOpInput.Size
					out.TotalSize = ast.GetArgSize(prgrm, inpExprAtomicOpInput)
					out.Type = inpExprAtomicOpInput.Type
					out.PointerTargetType = inpExprAtomicOpInput.PointerTargetType
					out.PreviouslyDeclared = true
				} else {
					out = ast.MakeArgument(generateTempVarName(constants.LOCAL_PREFIX), CurrentFile, inpExprCXLine.LineNumber).SetType(inpExprAtomicOpOperatorOutput.Type)
					out.DeclarationSpecifiers = inpExprAtomicOpOperatorOutput.DeclarationSpecifiers

					out.StructType = inpExprAtomicOpOperatorOutput.StructType

					if inpExprAtomicOpOperatorOutput.StructType != nil {
						inpExprPkg, err := prgrm.GetPackageFromArray(inpExprAtomicOp.Package)
						if err != nil {
							panic(err)
						}
						if strct, err := inpExprPkg.GetStruct(prgrm, inpExprAtomicOpOperatorOutput.StructType.Name); err == nil {
							out.Size = strct.GetStructSize(prgrm)
							out.TotalSize = strct.GetStructSize(prgrm)
						}
					} else {
						out.Size = inpExprAtomicOpOperatorOutput.Size
						out.TotalSize = ast.GetArgSize(prgrm, inpExprAtomicOpOperatorOutput)
					}

					out.Type = inpExprAtomicOpOperatorOutput.Type
					out.PointerTargetType = inpExprAtomicOpOperatorOutput.PointerTargetType
					out.PreviouslyDeclared = true
				}

				out.Package = inpExprAtomicOp.Package
				outIdx := prgrm.AddCXArgInArray(out)

				inpExprAtomicOp.AddOutput(prgrm, outIdx)
				expression.AddInput(prgrm, outIdx)
			}
			if len(inpExprAtomicOp.Outputs) > 0 && inpExpr.IsArrayLiteral() {
				expression.AddInput(prgrm, inpExprAtomicOp.Outputs[0])
			}
			nestedExprs = append(nestedExprs, inpExpr)
		}
	}
	return append(nestedExprs, exprs...)
}

// checkSameNativeType checks if all the inputs of an expression are of the same type.
// It is used mainly to prevent implicit castings in arithmetic operations
func checkSameNativeType(prgrm *ast.CXProgram, expr *ast.CXExpression) error {
	expression, err := prgrm.GetCXAtomicOp(expr.Index)
	if err != nil {
		panic(err)
	}

	if len(expression.Inputs) < 1 {
		return errors.New("cannot perform arithmetic without operands")
	}

	typeCode := prgrm.GetCXArgFromArray(expression.Inputs[0]).Type
	if prgrm.GetCXArgFromArray(expression.Inputs[0]).Type == types.POINTER {
		typeCode = prgrm.GetCXArgFromArray(expression.Inputs[0]).PointerTargetType
	}

	for _, inpIdx := range expression.Inputs {
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

// ProcessOperatorExpression checks if the inputs of an expression
// has the same type to prevent implicit castings in arithmetic operations.
//
// Input arguments description:
// 	prgrm - a CXProgram that contains all the data and arrays of the program.
// 	expr - the operator expression.
func ProcessOperatorExpression(prgrm *ast.CXProgram, expr *ast.CXExpression) {
	expression, err := prgrm.GetCXAtomicOp(expr.Index)
	if err != nil {
		panic(err)
	}
	expressionOperator := prgrm.GetFunctionFromArray(expression.Operator)

	if expressionOperator != nil && ast.IsOperator(expressionOperator.AtomicOPCode) {
		if err := checkSameNativeType(prgrm, expr); err != nil {
			println(ast.CompilationError(CurrentFile, LineNo), err.Error())
		}
	}
	if expr.IsUndType(prgrm) {
		for _, outIdx := range expression.Outputs {
			out := prgrm.GetCXArgFromArray(outIdx)

			size := types.Pointer(1)
			if !ast.IsComparisonOperator(expressionOperator.AtomicOPCode) {
				size = ast.GetArgSize(prgrm, prgrm.GetCXArgFromArray(expression.Inputs[0]).GetAssignmentElement(prgrm))
			}
			out.Size = size
			out.TotalSize = size
		}
	}
}

// ProcessPointerStructs checks and processes if any args is a pointer struct.
//
// Input arguments description:
// 	prgrm - a CXProgram that contains all the data and arrays of the program.
// 	expr - the operator expression.
func ProcessPointerStructs(prgrm *ast.CXProgram, expr *ast.CXExpression) {
	expression, err := prgrm.GetCXAtomicOp(expr.Index)
	if err != nil {
		panic(err)
	}

	for _, argIdx := range append(expression.Inputs, expression.Outputs...) {
		arg := prgrm.GetCXArgFromArray(argIdx)
		for _, fldIdx := range arg.Fields {
			fld := prgrm.GetCXArgFromArray(fldIdx)
			if fld.IsPointer() && fld.DereferenceLevels == 0 {
				prgrm.CXArgs[fldIdx].DereferenceLevels++
				prgrm.CXArgs[fldIdx].DereferenceOperations = append(prgrm.CXArgs[fldIdx].DereferenceOperations, constants.DEREF_POINTER)
			}
		}
		if arg.IsStruct && arg.IsPointer() && len(arg.Fields) > 0 && arg.DereferenceLevels == 0 {
			prgrm.CXArgs[argIdx].DereferenceLevels++
			prgrm.CXArgs[argIdx].DereferenceOperations = append(arg.DereferenceOperations, constants.DEREF_POINTER)
		}
	}
}

// processTestExpression checks for the special case of test calls. `assert`, `test`, `panic` are operators where
// their first input's type needs to be the same as its second input's type. This can't be handled by
// `checkSameNativeType` because these test functions' third input parameter is always a `str`.
//
// Input arguments description:
// 	prgrm - a CXProgram that contains all the data and arrays of the program.
// 	expr - the operator expression.
func processTestExpression(prgrm *ast.CXProgram, expr *ast.CXExpression) {
	expression, err := prgrm.GetCXAtomicOp(expr.Index)
	if err != nil {
		panic(err)
	}
	expressionOperator := prgrm.GetFunctionFromArray(expression.Operator)

	if expressionOperator != nil {
		opCode := expressionOperator.AtomicOPCode
		if opCode == constants.OP_ASSERT || opCode == constants.OP_TEST || opCode == constants.OP_PANIC {
			inp1Type := ast.GetFormattedType(prgrm, prgrm.GetCXArgFromArray(expression.Inputs[0]))
			inp2Type := ast.GetFormattedType(prgrm, prgrm.GetCXArgFromArray(expression.Inputs[1]))
			if inp1Type != inp2Type {
				println(ast.CompilationError(CurrentFile, LineNo), fmt.Sprintf("first and second input arguments' types are not equal in '%s' call ('%s' != '%s')", ast.OpNames[expressionOperator.AtomicOPCode], inp1Type, inp2Type))
			}
		}
	}
}

// checkIndexType throws an error if the type of `idx` is not `i32` or `i64`.
func checkIndexType(prgrm *ast.CXProgram, idxIdx ast.CXArgumentIndex) {
	idx := prgrm.GetCXArgFromArray(idxIdx)

	typ := ast.GetFormattedType(prgrm, idx)
	if typ != "i32" && typ != "i64" {
		println(ast.CompilationError(idx.ArgDetails.FileName, idx.ArgDetails.FileLine), fmt.Sprintf("wrong index type; expected either 'i32' or 'i64', got '%s'", typ))
	}
}

// ProcessExpressionArguments performs a series of checks and processes to an expresion's inputs and outputs.
// Some of these checks are: checking if a an input has not been declared, assign a relative offset to the argument,
// and calculate the correct size of the argument.
//
// Input arguments description:
// 	prgrm - a CXProgram that contains all the data and arrays of the program.
// 	symbols - a slice of string-CXArg map which corresponds to
// 			  the scoping of the  CXArguments. Each element in the slice
// 			  corresponds to a different scope. The innermost scope is
// 			  the last element of the slice.
//  symbolsScope - only handles the difference between local and global
// 				   scopes, local being function constrained variables,
// 				   and global being global variables.
//  offset - offset to use by statements (excluding inputs, outputs
// 			 and receiver).
//  fnIdx - the index of the function in the main CXFunction array.
//  args - the expression arguments.
//  expr - the expression.
//  isInput - true if args are input arguments, false if they are output args.
func ProcessExpressionArguments(prgrm *ast.CXProgram, symbols *[]map[string]*ast.CXArgument, symbolsScope *map[string]bool, offset *types.Pointer, fnIdx ast.CXFunctionIndex, args []ast.CXArgumentIndex, expr *ast.CXExpression, isInput bool) {
	fn := prgrm.GetFunctionFromArray(fnIdx)

	for _, argIdx := range args {
		arg := prgrm.GetCXArgFromArray(argIdx)
		ProcessLocalDeclaration(prgrm, symbolsScope, argIdx)

		if !isInput {
			CheckRedeclared(prgrm, symbols, expr, argIdx)
		}

		if !isInput {
			ProcessOperatorExpression(prgrm, expr)
		}

		if arg.PreviouslyDeclared {
			UpdateSymbolsTable(prgrm, symbols, argIdx, offset, false)
		} else {
			UpdateSymbolsTable(prgrm, symbols, argIdx, offset, true)
		}

		if isInput {
			GiveOffset(prgrm, symbols, argIdx, true)
		} else {
			GiveOffset(prgrm, symbols, argIdx, false)
		}

		ProcessSlice(prgrm, argIdx)

		for _, idxIdx := range arg.Indexes {
			UpdateSymbolsTable(prgrm, symbols, idxIdx, offset, true)
			GiveOffset(prgrm, symbols, idxIdx, true)
			checkIndexType(prgrm, idxIdx)
		}
		for _, fldIdx := range arg.Fields {
			fld := prgrm.GetCXArgFromArray(fldIdx)
			for _, idxIdx := range fld.Indexes {
				UpdateSymbolsTable(prgrm, symbols, idxIdx, offset, true)
				GiveOffset(prgrm, symbols, idxIdx, true)
			}
		}

		SetFinalSize(prgrm, symbols, argIdx)

		AddPointer(prgrm, fn, argIdx)
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
				prgrm.CXArgs[sym.Fields[len(sym.Fields)-1]].Name == prgrm.CXArgs[ptr.Fields[len(ptr.Fields)-1]].Name {
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
//
// Input arguments description:
// 	prgrm - a CXProgram that contains all the data and arrays of the program.
// 	fn - the function the sym belongs.
//  symIdx - the index of the sym from the main CXArg array.
func AddPointer(prgrm *ast.CXProgram, fn *ast.CXFunction, symIdx ast.CXArgumentIndex) {
	sym := prgrm.GetCXArgFromArray(symIdx)

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
		field := prgrm.GetCXArgFromArray(sym.Fields[len(sym.Fields)-1])
		if field.IsPointer() && !isPointerAdded(prgrm, fn, sym) {
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
//
// Input arguments description:
// 	prgrm - a CXProgram that contains all the data and arrays of the program.
// 	symbols - a slice of string-CXArg map which corresponds to
// 			  the scoping of the  CXArguments. Each element in the slice
// 			  corresponds to a different scope. The innermost scope is
// 			  the last element of the slice.
//  expr - the expression.
//  symIdx - the index of the sym from the main CXArg array.
func CheckRedeclared(prgrm *ast.CXProgram, symbols *[]map[string]*ast.CXArgument, expr *ast.CXExpression, symIdx ast.CXArgumentIndex) {
	sym := prgrm.GetCXArgFromArray(symIdx)

	expression, err := prgrm.GetCXAtomicOp(expr.Index)
	if err != nil {
		panic(err)
	}
	expressionOperator := prgrm.GetFunctionFromArray(expression.Operator)

	if expressionOperator == nil && len(expression.Outputs) > 0 && len(expression.Inputs) == 0 {
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

// ProcessLocalDeclaration sets symbolsScope to true if the arg is a
// local declaration.
//
// Input arguments description:
// 	prgrm - a CXProgram that contains all the data and arrays of the program.
//  symbolsScope - only handles the difference between local and global
// 				   scopes, local being function constrained variables,
// 				   and global being global variables.
//  argIdx - index from the main CXArg array.
func ProcessLocalDeclaration(prgrm *ast.CXProgram, symbolsScope *map[string]bool, argIdx ast.CXArgumentIndex) {
	arg := prgrm.GetCXArgFromArray(argIdx)

	argPkg, err := prgrm.GetPackageFromArray(arg.Package)
	if err != nil {
		panic(err)
	}

	if arg.IsLocalDeclaration {
		(*symbolsScope)[argPkg.Name+"."+arg.Name] = true
	}
	arg.IsLocalDeclaration = (*symbolsScope)[argPkg.Name+"."+arg.Name]
}

// ProcessGoTos sets the ThenLines value if the expression is a goto.
//
// Input arguments description:
// 	prgrm - a CXProgram that contains all the data and arrays of the program.
//  exprs - slice of expressions.
func ProcessGoTos(prgrm *ast.CXProgram, exprs []ast.CXExpression) {
	for i, expr := range exprs {
		expression, _, _, err := prgrm.GetOperation(&expr)
		if err != nil {
			panic(err)
		}
		expressionOperator := prgrm.GetFunctionFromArray(expression.Operator)

		opGotoFn := ast.Natives[constants.OP_GOTO]
		isOpGoto := false
		if expressionOperator != nil && expressionOperator.AtomicOPCode == opGotoFn.AtomicOPCode && len(expressionOperator.Inputs) == len(opGotoFn.Inputs) && len(expressionOperator.Outputs) == len(opGotoFn.Outputs) {
			isOpGoto = true
		}
		if isOpGoto {
			// then it's a goto
			for j, e := range exprs {
				expressionToGoTo, _, _, err := prgrm.GetOperation(&e)
				if err != nil {
					panic(err)
				}

				if expressionToGoTo.Label == expression.Label && i != j {
					// ElseLines is used because arg's default val is false
					expression.ThenLines = j - i - 1
					break
				}
			}
		}
	}
}

// AddExprsToFunction add all expressions to the function.
func AddExprsToFunction(prgrm *ast.CXProgram, fnIdx ast.CXFunctionIndex, exprs []ast.CXExpression) {
	fn := prgrm.GetFunctionFromArray(fnIdx)
	for _, expr := range exprs {
		fn.AddExpression(prgrm, &expr)
	}
}

func checkMatchParamTypes(prgrm *ast.CXProgram, expr *ast.CXExpression, expected, received []ast.CXArgumentIndex, isInputs bool) {
	expression, err := prgrm.GetCXAtomicOp(expr.Index)
	if err != nil {
		panic(err)
	}
	expressionOperator := prgrm.GetFunctionFromArray(expression.Operator)

	for i, inpIdx := range expected {
		expectedArg := prgrm.GetCXArgFromArray(inpIdx)
		receivedArg := prgrm.GetCXArg(received[i])

		expectedType := ast.GetFormattedType(prgrm, expectedArg)
		receivedType := ast.GetFormattedType(prgrm, receivedArg)

		if expr.IsMethodCall() && expectedArg.IsPointer() && i == 0 {
			// if method receiver is pointer, remove *
			if expectedType[0] == '*' {
				// we need to check if it's not an `str`
				// otherwise we end up removing the `s` instead of a `*`
				expectedType = expectedType[1:]
			}
		}

		if expectedType != receivedType && expectedArg.Type != types.UNDEFINED {
			var opName string
			if expressionOperator.IsBuiltIn() {
				opName = ast.OpNames[expressionOperator.AtomicOPCode]
			} else {
				opName = expressionOperator.Name
			}

			if isInputs {
				println(ast.CompilationError(receivedArg.ArgDetails.FileName, receivedArg.ArgDetails.FileLine), fmt.Sprintf("function '%s' expected input argument of type '%s'; '%s' was provided", opName, expectedType, receivedType))
			} else {
				println(ast.CompilationError(prgrm.GetCXArgFromArray(expression.Outputs[i]).ArgDetails.FileName, prgrm.GetCXArgFromArray(expression.Outputs[i]).ArgDetails.FileLine), fmt.Sprintf("function '%s' expected receiving variable of type '%s'; '%s' was provided", opName, expectedType, receivedType))
			}

		}

		// In the case of assignment we need to check that the input's type matches the output's type.
		// FIXME: There are some expressions added by the cxgo where temporary variables are used.
		// These temporary variables' types are not properly being set. That's why we use !cxcore.IsTempVar to
		// exclude these cases for now.
		if expressionOperator.AtomicOPCode == constants.OP_IDENTITY && !IsTempVar(prgrm.GetCXArgFromArray(expression.Outputs[0]).Name) {
			inpType := ast.GetFormattedType(prgrm, prgrm.GetCXArgFromArray(expression.Inputs[0]))
			outType := ast.GetFormattedType(prgrm, prgrm.GetCXArgFromArray(expression.Outputs[0]))

			// We use `isInputs` to only print the error once.
			// Otherwise we'd print the error twice: once for the input and again for the output
			if inpType != outType && isInputs {
				println(ast.CompilationError(receivedArg.ArgDetails.FileName, receivedArg.ArgDetails.FileLine), fmt.Sprintf("cannot assign value of type '%s' to identifier '%s' of type '%s'", inpType, prgrm.GetCXArgFromArray(expression.Outputs[0]).GetAssignmentElement(prgrm).Name, outType))
			}
		}
	}
}

// CheckTypes checks if the expected types are provided and outputted correctly.
func CheckTypes(prgrm *ast.CXProgram, exprs []ast.CXExpression, currIndex int) {
	expression, err := prgrm.GetCXAtomicOpFromExpressions(exprs, currIndex)
	if err != nil {
		panic(err)
	}
	expressionOperator := prgrm.GetFunctionFromArray(expression.Operator)

	exprCXLine, _ := prgrm.GetPreviousCXLine(exprs, currIndex)

	if expressionOperator != nil {
		opName := expression.GetOperatorName(prgrm)

		// checking if number of inputs is less than the required number of inputs
		if len(expression.Inputs) != len(expressionOperator.Inputs) {
			if !(len(expressionOperator.Inputs) > 0 && prgrm.GetCXArgFromArray(expressionOperator.Inputs[len(expressionOperator.Inputs)-1]).Type != types.UNDEFINED) {
				// if the last input is of type cxcore.TYPE_UNDEFINED then it might be a variadic function, such as printf
			} else {
				// then we need to be strict in the number of inputs
				var plural1 string
				var plural2 string = "s"
				var plural3 string = "were"
				if len(expressionOperator.Inputs) > 1 {
					plural1 = "s"
				}
				if len(expression.Inputs) == 1 {
					plural2 = ""
					plural3 = "was"
				}

				println(ast.CompilationError(exprCXLine.FileName, exprCXLine.LineNumber), fmt.Sprintf("operator '%s' expects %d input%s, but %d input argument%s %s provided", opName, len(expressionOperator.Inputs), plural1, len(expression.Inputs), plural2, plural3))
				return
			}
		}

		// checking if number of expr.ProgramOutput matches number of Operator.ProgramOutput
		if len(expression.Outputs) != len(expressionOperator.Outputs) {
			var plural1 string
			var plural2 string = "s"
			var plural3 string = "were"
			if len(expressionOperator.Outputs) > 1 {
				plural1 = "s"
			}
			if len(expression.Outputs) == 1 {
				plural2 = ""
				plural3 = "was"
			}

			println(ast.CompilationError(exprCXLine.FileName, exprCXLine.LineNumber), fmt.Sprintf("operator '%s' expects to return %d output%s, but %d receiving argument%s %s provided", opName, len(expressionOperator.Outputs), plural1, len(expression.Outputs), plural2, plural3))
			os.Exit(constants.CX_COMPILATION_ERROR)
		}
	}

	if expressionOperator != nil && expressionOperator.IsBuiltIn() && expressionOperator.AtomicOPCode == constants.OP_IDENTITY {
		for i := range expression.Inputs {
			var expectedType string
			var receivedType string
			if prgrm.GetCXArgFromArray(expression.Outputs[i]).GetAssignmentElement(prgrm).StructType != nil {
				// then it's custom type
				expectedType = prgrm.GetCXArgFromArray(expression.Outputs[i]).GetAssignmentElement(prgrm).StructType.Name

			} else {
				// then it's native type
				expectedType = prgrm.GetCXArgFromArray(expression.Outputs[i]).GetAssignmentElement(prgrm).Type.Name()

				if prgrm.GetCXArgFromArray(expression.Outputs[i]).GetAssignmentElement(prgrm).Type == types.POINTER {
					expectedType = prgrm.GetCXArgFromArray(expression.Outputs[i]).GetAssignmentElement(prgrm).PointerTargetType.Name()
				}
			}

			if prgrm.GetCXArgFromArray(expression.Inputs[i]).GetAssignmentElement(prgrm).StructType != nil {
				// then it's custom type
				receivedType = prgrm.GetCXArgFromArray(expression.Inputs[i]).GetAssignmentElement(prgrm).StructType.Name
			} else {
				// then it's native type
				receivedType = prgrm.GetCXArgFromArray(expression.Inputs[i]).GetAssignmentElement(prgrm).Type.Name()

				if prgrm.GetCXArgFromArray(expression.Inputs[i]).GetAssignmentElement(prgrm).Type == types.POINTER {
					receivedType = prgrm.GetCXArgFromArray(expression.Inputs[i]).GetAssignmentElement(prgrm).PointerTargetType.Name()
				}
			}

			if receivedType != expectedType {
				if exprs[currIndex].IsStructLiteral() {
					println(ast.CompilationError(prgrm.GetCXArgFromArray(expression.Outputs[i]).ArgDetails.FileName, prgrm.GetCXArgFromArray(expression.Outputs[i]).ArgDetails.FileLine), fmt.Sprintf("field '%s' in struct literal of type '%s' expected argument of type '%s'; '%s' was provided", prgrm.GetCXArgFromArray(prgrm.GetCXArgFromArray(expression.Outputs[i]).Fields[0]).Name, prgrm.GetCXArgFromArray(expression.Outputs[i]).StructType.Name, expectedType, receivedType))
				} else {
					println(ast.CompilationError(prgrm.GetCXArgFromArray(expression.Outputs[i]).ArgDetails.FileName, prgrm.GetCXArgFromArray(expression.Outputs[i]).ArgDetails.FileLine), fmt.Sprintf("trying to assign argument of type '%s' to symbol '%s' of type '%s'", receivedType, prgrm.GetCXArgFromArray(expression.Outputs[i]).GetAssignmentElement(prgrm).Name, expectedType))
				}
			}
		}
	}

	// then it's a function call and not a declaration
	if expressionOperator != nil {
		// checking inputs matching operator's inputs
		checkMatchParamTypes(prgrm, &exprs[currIndex], expressionOperator.Inputs, expression.Inputs, true)

		// checking outputs matching operator's outputs
		checkMatchParamTypes(prgrm, &exprs[currIndex], expressionOperator.Outputs, expression.Outputs, false)
	}
}

// ProcessStringAssignment sets the args PassBy to PASSBY_VALUE if the type is string.
func ProcessStringAssignment(prgrm *ast.CXProgram, expr *ast.CXExpression) {
	expression, err := prgrm.GetCXAtomicOp(expr.Index)
	if err != nil {
		panic(err)
	}
	expressionOperator := prgrm.GetFunctionFromArray(expression.Operator)

	opIdentFn := ast.Natives[constants.OP_IDENTITY]
	isOpIdent := false
	if expressionOperator != nil && expressionOperator.AtomicOPCode == opIdentFn.AtomicOPCode && len(expressionOperator.Inputs) == len(opIdentFn.Inputs) && len(expressionOperator.Outputs) == len(opIdentFn.Outputs) {
		isOpIdent = true
	}
	if isOpIdent {
		for i, outIdx := range expression.Outputs {
			out := prgrm.GetCXArgFromArray(outIdx)
			if len(expression.Inputs) > i {
				out = out.GetAssignmentElement(prgrm)
				inp := prgrm.GetCXArgFromArray(expression.Inputs[i]).GetAssignmentElement(prgrm)

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
	expression, err := prgrm.GetCXAtomicOp(expr.Index)
	if err != nil {
		panic(err)
	}

	for _, outIdx := range expression.Outputs {
		out := prgrm.GetCXArgFromArray(outIdx)
		elt := out.GetAssignmentElement(prgrm)
		if elt.PassBy == constants.PASSBY_REFERENCE &&
			!hasDeclSpec(elt, constants.DECL_POINTER) &&
			elt.PointerTargetType != types.STR && elt.Type != types.STR && !elt.IsSlice {
			println(ast.CompilationError(CurrentFile, LineNo), "invalid reference assignment", elt.Name)
		}
	}

}

// ProcessShortDeclaration sets proper values if the expr is a short declaration.
func ProcessShortDeclaration(prgrm *ast.CXProgram, expr *ast.CXExpression, expressions []ast.CXExpression, idx int) {
	expression, err := prgrm.GetCXAtomicOp(expr.Index)
	if err != nil {
		panic(err)
	}

	// process short declaration
	if len(expression.Outputs) > 0 && len(expression.Inputs) > 0 && prgrm.GetCXArgFromArray(expression.Outputs[0]).IsShortAssignmentDeclaration && !expr.IsStructLiteral() && !isParseOp(prgrm, expr) {
		expressionOperator := prgrm.GetFunctionFromArray(expression.Operator)
		prevExpression, err := prgrm.GetPreviousCXAtomicOpFromExpressions(expressions, idx-1)
		if err != nil {
			panic(err)
		}

		var arg *ast.CXArgument
		if expr.IsMethodCall() {
			arg = prgrm.GetCXArgFromArray(expressionOperator.Outputs[0])
		} else {
			arg = prgrm.GetCXArgFromArray(expression.Inputs[0])
		}

		prevExpressionOutputIdx := prevExpression.Outputs[0]
		prgrm.CXArgs[prevExpressionOutputIdx].Type = arg.Type
		prgrm.CXArgs[prevExpressionOutputIdx].PointerTargetType = arg.PointerTargetType
		prgrm.CXArgs[prevExpressionOutputIdx].Size = arg.Size
		prgrm.CXArgs[prevExpressionOutputIdx].TotalSize = arg.TotalSize

		expressionOutputIdx := expression.Outputs[0]
		prgrm.CXArgs[expressionOutputIdx].Type = arg.Type
		prgrm.CXArgs[expressionOutputIdx].PointerTargetType = arg.PointerTargetType
		prgrm.CXArgs[expressionOutputIdx].Size = arg.Size
		prgrm.CXArgs[expressionOutputIdx].TotalSize = arg.TotalSize
	}
}

// ProcessSlice sets DereferenceOperations if the arg is a slice.
func ProcessSlice(prgrm *ast.CXProgram, inpIdx ast.CXArgumentIndex) {
	inp := prgrm.GetCXArgFromArray(inpIdx)

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

// ProcessSliceAssignent sets correct values for slice assignment expressions.
func ProcessSliceAssignment(prgrm *ast.CXProgram, expr *ast.CXExpression) {
	expression, err := prgrm.GetCXAtomicOp(expr.Index)
	if err != nil {
		panic(err)
	}
	expressionOperator := prgrm.GetFunctionFromArray(expression.Operator)

	opIdentFn := ast.Natives[constants.OP_IDENTITY]
	isOpIdent := false
	if expressionOperator != nil && expressionOperator.AtomicOPCode == opIdentFn.AtomicOPCode && len(expressionOperator.Inputs) == len(opIdentFn.Inputs) && len(expressionOperator.Outputs) == len(opIdentFn.Outputs) {
		isOpIdent = true
	}

	if isOpIdent {
		var inp *ast.CXArgument
		var out *ast.CXArgument

		inp = prgrm.GetCXArgFromArray(expression.Inputs[0]).GetAssignmentElement(prgrm)
		out = prgrm.GetCXArgFromArray(expression.Outputs[0]).GetAssignmentElement(prgrm)

		if inp.IsSlice && out.IsSlice && len(inp.Indexes) == 0 && len(out.Indexes) == 0 {
			out.PassBy = constants.PASSBY_VALUE
		}
	}
	if expressionOperator != nil && !expressionOperator.IsBuiltIn() {
		// then it's a function call
		for _, inpIdx := range expression.Inputs {
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
	fnArg := ast.MakeArgument(ident, fn.FileName, fn.FileLine).SetType(types.FUNC)
	fnArg.Package = ast.CXPackageIndex(pkg.Index)

	return fnArg, nil
}

// UpdateSymbolsTable adds `sym` to the innermost scope (last element of slice) in `symbols`.
func UpdateSymbolsTable(prgrm *ast.CXProgram, symbols *[]map[string]*ast.CXArgument, symIdx ast.CXArgumentIndex, offset *types.Pointer, shouldExist bool) {
	sym := prgrm.GetCXArgFromArray(symIdx)

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
			*offset += ast.GetArgSize(prgrm, sym)
		}
	}
}

// ProcessMethodCall completes the method call expression's inputs and outputs.
//
// Input arguments description:
// 	prgrm - a CXProgram that contains all the data and arrays of the program.
//  expr - the method call expression.
// 	symbols - a slice of string-CXArg map which corresponds to
// 			  the scoping of the  CXArguments. Each element in the slice
// 			  corresponds to a different scope. The innermost scope is
// 			  the last element of the slice.
func ProcessMethodCall(prgrm *ast.CXProgram, expr *ast.CXExpression, symbols *[]map[string]*ast.CXArgument) {
	expression, err := prgrm.GetCXAtomicOp(expr.Index)
	if err != nil {
		panic(err)
	}

	if expr.IsMethodCall() {
		var inpIdx ast.CXArgumentIndex = -1
		var outIdx ast.CXArgumentIndex = -1

		if len(expression.Inputs) > 0 && prgrm.GetCXArgFromArray(expression.Inputs[0]).Name != "" {
			inpIdx = expression.Inputs[0]
		}
		if len(expression.Outputs) > 0 && prgrm.GetCXArgFromArray(expression.Outputs[0]).Name != "" {
			outIdx = expression.Outputs[0]
		}

		if inpIdx != -1 {
			inpPkg, err := prgrm.GetPackageFromArray(prgrm.CXArgs[inpIdx].Package)
			if err != nil {
				panic(err)
			}

			if argInp, err := lookupSymbol(prgrm, inpPkg.Name, prgrm.CXArgs[inpIdx].Name, symbols); err != nil {
				if outIdx == -1 {
					panic("")
				}

				outPkg, err := prgrm.GetPackageFromArray(prgrm.CXArgs[outIdx].Package)
				if err != nil {
					panic(err)
				}

				argOut, err := lookupSymbol(prgrm, outPkg.Name, prgrm.CXArgs[outIdx].Name, symbols)
				if err != nil {
					println(ast.CompilationError(prgrm.CXArgs[outIdx].ArgDetails.FileName, prgrm.CXArgs[outIdx].ArgDetails.FileLine), fmt.Sprintf("identifier '%s' does not exist", prgrm.CXArgs[outIdx].Name))
					os.Exit(constants.CX_COMPILATION_ERROR)
				}
				// then we found an output
				if len(prgrm.CXArgs[outIdx].Fields) > 0 {
					strct := argOut.StructType
					strctPkg, err := prgrm.GetPackageFromArray(strct.Package)
					if err != nil {
						panic(err)
					}

					if fnIdx, err := strctPkg.GetMethod(prgrm, strct.Name+"."+prgrm.GetCXArgFromArray(prgrm.CXArgs[outIdx].Fields[len(prgrm.CXArgs[outIdx].Fields)-1]).Name, strct.Name); err == nil {
						prgrm.CXAtomicOps[expr.Index].Operator = fnIdx
					} else {
						panic("")
					}

					prgrm.CXAtomicOps[expr.Index].Inputs = append([]ast.CXArgumentIndex{ast.CXArgumentIndex(prgrm.CXArgs[outIdx].Index)}, prgrm.CXAtomicOps[expr.Index].Inputs...)
					prgrm.CXAtomicOps[expr.Index].Outputs = prgrm.CXAtomicOps[expr.Index].Outputs[1:]
					prgrm.CXArgs[outIdx].Fields = prgrm.CXArgs[outIdx].Fields[:len(prgrm.CXArgs[outIdx].Fields)-1]
				}
			} else {
				// then we found an input
				if len(prgrm.CXArgs[inpIdx].Fields) > 0 {
					strct := argInp.StructType

					for _, fldIdx := range prgrm.CXArgs[inpIdx].Fields {
						field := prgrm.GetCXArgFromArray(fldIdx)
						if inFld, err := strct.GetField(prgrm, field.Name); err == nil {
							if inFld.StructType != nil {
								strct = inFld.StructType
							}
						}
					}

					strctPkg, err := prgrm.GetPackageFromArray(strct.Package)
					if err != nil {
						panic(err)
					}

					if fnIdx, err := strctPkg.GetMethod(prgrm, strct.Name+"."+prgrm.GetCXArgFromArray(prgrm.CXArgs[inpIdx].Fields[len(prgrm.CXArgs[inpIdx].Fields)-1]).Name, strct.Name); err == nil {
						prgrm.CXAtomicOps[expr.Index].Operator = fnIdx
					} else {
						panic(err)
					}

					prgrm.CXArgs[inpIdx].Fields = prgrm.CXArgs[inpIdx].Fields[:len(prgrm.CXArgs[inpIdx].Fields)-1]
				} else if len(prgrm.CXArgs[outIdx].Fields) > 0 {
					outPkg, err := prgrm.GetPackageFromArray(prgrm.CXArgs[outIdx].Package)
					if err != nil {
						panic(err)
					}
					argOut, err := lookupSymbol(prgrm, outPkg.Name, prgrm.CXArgs[outIdx].Name, symbols)
					if err != nil {
						panic(err)
					}

					strct := argOut.StructType
					if strct == nil {
						println(ast.CompilationError(argOut.ArgDetails.FileName, argOut.ArgDetails.FileLine), fmt.Sprintf("illegal method call or field access on identifier '%s' of primitive type '%s'", argOut.Name, argOut.Type.Name()))
						os.Exit(constants.CX_COMPILATION_ERROR)
					}

					prgrm.CXAtomicOps[expr.Index].Inputs = append(prgrm.CXAtomicOps[expr.Index].Outputs[:1], prgrm.CXAtomicOps[expr.Index].Inputs...)
					prgrm.CXAtomicOps[expr.Index].Outputs = prgrm.CXAtomicOps[expr.Index].Outputs[:len(prgrm.CXAtomicOps[expr.Index].Outputs)-1]

					strctPkg, err := prgrm.GetPackageFromArray(strct.Package)
					if err != nil {
						panic(err)
					}

					if fnIdx, err := strctPkg.GetMethod(prgrm, strct.Name+"."+prgrm.GetCXArgFromArray(prgrm.CXArgs[outIdx].Fields[len(prgrm.CXArgs[outIdx].Fields)-1]).Name, strct.Name); err == nil {
						prgrm.CXAtomicOps[expr.Index].Operator = fnIdx
					} else {
						panic(err)
					}

					prgrm.CXArgs[outIdx].Fields = prgrm.CXArgs[outIdx].Fields[:len(prgrm.CXArgs[outIdx].Fields)-1]
				}
			}
		} else {
			if outIdx == -1 {
				panic("")
			}

			outPkg, err := prgrm.GetPackageFromArray(prgrm.CXArgs[outIdx].Package)
			if err != nil {
				panic(err)
			}

			argOut, err := lookupSymbol(prgrm, outPkg.Name, prgrm.CXArgs[outIdx].Name, symbols)
			if err != nil {
				println(ast.CompilationError(prgrm.CXArgs[outIdx].ArgDetails.FileName, prgrm.CXArgs[outIdx].ArgDetails.FileLine), fmt.Sprintf("identifier '%s' does not exist", prgrm.CXArgs[outIdx].Name))
				os.Exit(constants.CX_COMPILATION_ERROR)
			}

			// then we found an output
			if len(prgrm.CXArgs[outIdx].Fields) > 0 {
				strct := argOut.StructType

				if strct == nil {
					println(ast.CompilationError(argOut.ArgDetails.FileName, argOut.ArgDetails.FileLine), fmt.Sprintf("illegal method call or field access on identifier '%s' of primitive type '%s'", argOut.Name, argOut.Type.Name()))
					os.Exit(constants.CX_COMPILATION_ERROR)
				}

				strctPkg, err := prgrm.GetPackageFromArray(strct.Package)
				if err != nil {
					panic(err)
				}

				if fnIdx, err := strctPkg.GetMethod(prgrm, strct.Name+"."+prgrm.GetCXArgFromArray(prgrm.CXArgs[outIdx].Fields[len(prgrm.CXArgs[outIdx].Fields)-1]).Name, strct.Name); err == nil {
					prgrm.CXAtomicOps[expr.Index].Operator = fnIdx
				} else {
					panic("")
				}

				prgrm.CXAtomicOps[expr.Index].Inputs = append([]ast.CXArgumentIndex{ast.CXArgumentIndex(prgrm.CXArgs[outIdx].Index)}, prgrm.CXAtomicOps[expr.Index].Inputs...)
				prgrm.CXAtomicOps[expr.Index].Outputs = prgrm.CXAtomicOps[expr.Index].Outputs[1:]
				prgrm.CXArgs[outIdx].Fields = prgrm.CXArgs[outIdx].Fields[:len(prgrm.CXArgs[outIdx].Fields)-1]
			}
		}
		expressionOperator := prgrm.GetFunctionFromArray(expression.Operator)

		// checking if receiver is sent as pointer or not
		if prgrm.GetCXArgFromArray(expressionOperator.Inputs[0]).IsPointer() {
			prgrm.CXArgs[expression.Inputs[0]].PassBy = constants.PASSBY_REFERENCE
		}
	}
}

// GiveOffset
func GiveOffset(prgrm *ast.CXProgram, symbols *[]map[string]*ast.CXArgument, symIdx ast.CXArgumentIndex, shouldExist bool) {
	sym := prgrm.GetCXArgFromArray(symIdx)

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

// ProcessTempVariable completes the temp variable expression.
func ProcessTempVariable(prgrm *ast.CXProgram, expr *ast.CXExpression) {
	expression, err := prgrm.GetCXAtomicOp(expr.Index)
	if err != nil {
		panic(err)
	}
	expressionOperator := prgrm.GetFunctionFromArray(expression.Operator)

	opIdentFn := ast.Natives[constants.OP_IDENTITY]
	isOpIdent := false
	if expressionOperator != nil && expressionOperator.AtomicOPCode == opIdentFn.AtomicOPCode && len(expressionOperator.Inputs) == len(opIdentFn.Inputs) && len(expressionOperator.Outputs) == len(opIdentFn.Outputs) {
		isOpIdent = true
	}

	if expressionOperator != nil && (isOpIdent || ast.IsArithmeticOperator(expressionOperator.AtomicOPCode)) && len(expression.Outputs) > 0 && len(expression.Inputs) > 0 {
		outputArgIdx := expression.Outputs[0]
		name := prgrm.CXArgs[outputArgIdx].Name
		if IsTempVar(name) {
			expressionInput := prgrm.GetCXArgFromArray(expression.Inputs[0])
			// then it's a temporary variable and it needs to adopt its input's type
			prgrm.CXArgs[outputArgIdx].Type = expressionInput.Type
			prgrm.CXArgs[outputArgIdx].PointerTargetType = expressionInput.PointerTargetType
			prgrm.CXArgs[outputArgIdx].Size = expressionInput.Size
			prgrm.CXArgs[outputArgIdx].TotalSize = expressionInput.TotalSize
			prgrm.CXArgs[outputArgIdx].PreviouslyDeclared = true
		}
	}
}

// CopyArgFields copies 'arg' fields to 'sym' fields.
func CopyArgFields(prgrm *ast.CXProgram, sym *ast.CXArgument, arg *ast.CXArgument) {
	sym.Offset = arg.Offset
	sym.Type = arg.Type

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
	sym.Size = arg.Size

	// Checking if it's a slice struct field. We'll do the same process as
	// below (as in the `arg.IsSlice` check), but the process differs in the
	// case of a slice struct field.
	assignElement := sym.GetAssignmentElement(prgrm)
	if (!arg.IsSlice || hasDerefOp(sym, constants.DEREF_ARRAY)) && arg.StructType != nil && assignElement.IsSlice && assignElement != sym {
		for i, deref := range assignElement.DereferenceOperations {
			// The cxgo when reading `foo[5]` in postfix.go does not know if `foo`
			// is a slice or an array. At this point we now know it's a slice and we need
			// to change those dereferences to cxcore.DEREF_SLICE.
			if deref == constants.DEREF_ARRAY {
				assignElement.DereferenceOperations[i] = constants.DEREF_SLICE
			}
		}

		if len(assignElement.DereferenceOperations) > 0 && assignElement.DereferenceOperations[0] == constants.DEREF_POINTER {
			assignElement.DereferenceOperations = assignElement.DereferenceOperations[1:]
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

// ProcessSymbolFields copies the correct field values for the sym.Fields from their struct fields.
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
			field := prgrm.GetCXArgFromArray(fldIdx)
			if inFld, err := strct.GetField(prgrm, field.Name); err == nil {
				if inFld.StructType != nil {
					field.StructType = strct
					strct = inFld.StructType
				}
			} else {
				methodName := prgrm.GetCXArgFromArray(sym.Fields[len(sym.Fields)-1]).Name
				receiverType := strct.Name

				if methodIdx, methodErr := strctPkg.GetMethod(prgrm, receiverType+"."+methodName, receiverType); methodErr == nil {
					method := prgrm.GetFunctionFromArray(methodIdx)
					field.Type = prgrm.GetCXArgFromArray(method.Outputs[0]).Type
					field.PointerTargetType = prgrm.GetCXArgFromArray(method.Outputs[0]).PointerTargetType
				} else {
					println(ast.CompilationError(field.ArgDetails.FileName, field.ArgDetails.FileLine), err.Error())
				}
			}
		}

		strct = arg.StructType
		// then we copy all the type struct fields
		// to the respective sym.Fields
		for _, nameFldIdx := range sym.Fields {
			nameField := prgrm.GetCXArgFromArray(nameFldIdx)
			if nameField.StructType != nil {
				strct = nameField.StructType
			}

			for _, typeSignature := range strct.Fields {
				if nameField.Name == typeSignature.Name && typeSignature.Type == ast.TYPE_ATOMIC {
					nameField.Type = types.Code(typeSignature.Meta)
					nameField.StructType = nil
					nameField.Size = typeSignature.GetSize(prgrm)
					nameField.TotalSize = typeSignature.GetSize(prgrm)

					// TODO: this should not be needed.
					if len(nameField.DeclarationSpecifiers) > 0 {
						nameField.DeclarationSpecifiers = append([]int{constants.DECL_BASIC}, nameField.DeclarationSpecifiers[1:]...)
					} else {
						nameField.DeclarationSpecifiers = []int{constants.DECL_BASIC}
					}

					break
				} else if nameField.Name == typeSignature.Name && typeSignature.Type == ast.TYPE_ARRAY_ATOMIC {
					typeSignatureArray := prgrm.GetTypeSignatureArrayFromArray(typeSignature.Meta)
					nameField.Type = types.Code(typeSignatureArray.Type)
					nameField.StructType = nil
					nameField.Size = typeSignature.GetSize(prgrm)
					nameField.Lengths = []types.Pointer{typeSignature.GetArrayLength(prgrm)}
					sym.Lengths = []types.Pointer{typeSignature.GetArrayLength(prgrm)}
					nameField.TotalSize = typeSignature.GetSize(prgrm) * nameField.Lengths[0]

					// TODO: this should not be needed.
					if len(nameField.DeclarationSpecifiers) > 0 {
						nameField.DeclarationSpecifiers = append([]int{constants.DECL_BASIC, constants.DECL_ARRAY}, nameField.DeclarationSpecifiers[1:]...)
					} else {
						nameField.DeclarationSpecifiers = []int{constants.DECL_BASIC}
					}
					break
				} else if nameField.Name == typeSignature.Name && typeSignature.Type == ast.TYPE_SLICE_ATOMIC {
					nameField.Type = types.Code(typeSignature.Meta)
					nameField.StructType = nil
					nameField.Size = types.Code(typeSignature.Meta).Size()
					nameField.Lengths = []types.Pointer{0}
					sym.Lengths = []types.Pointer{0}
					nameField.TotalSize = typeSignature.GetSize(prgrm)
					nameField.IsSlice = true
					// nameField.PassBy = constants.PASSBY_REFERENCE
					// TODO: this should not be needed.
					if len(nameField.DeclarationSpecifiers) > 0 {
						nameField.DeclarationSpecifiers = append([]int{constants.DECL_BASIC, constants.DECL_SLICE}, nameField.DeclarationSpecifiers[1:]...)
					} else {
						nameField.DeclarationSpecifiers = []int{constants.DECL_BASIC}
					}

					nameField.DereferenceOperations = append([]int{constants.DEREF_POINTER}, nameField.DereferenceOperations...)
					nameField.DereferenceLevels++

					break
				}

				fieldIdx := typeSignature.Meta
				field := prgrm.CXArgs[fieldIdx]

				if nameField.Name == field.Name && typeSignature.Type != ast.TYPE_ATOMIC {
					nameField.Type = field.Type
					nameField.Lengths = field.Lengths
					nameField.Size = field.Size
					nameField.TotalSize = field.TotalSize
					nameField.DereferenceLevels = sym.DereferenceLevels
					nameField.PointerTargetType = field.PointerTargetType
					nameField.StructType = field.StructType

					sym.Lengths = field.Lengths

					if len(nameField.DeclarationSpecifiers) > 0 {
						nameField.DeclarationSpecifiers = append(field.DeclarationSpecifiers, nameField.DeclarationSpecifiers[1:]...)
					} else {
						nameField.DeclarationSpecifiers = field.DeclarationSpecifiers
					}

					if field.IsSlice {
						nameField.DereferenceOperations = append([]int{constants.DEREF_POINTER}, nameField.DereferenceOperations...)
						nameField.DereferenceLevels++
					}

					nameField.PassBy = field.PassBy
					nameField.IsSlice = field.IsSlice

					if field.Type == types.STR || field.Type == types.AFF {
						nameField.PassBy = constants.PASSBY_REFERENCE
					}

					if field.StructType != nil {
						strct = field.StructType
					}

					break
				}

				nameField.Offset += typeSignature.GetSize(prgrm)
			}
		}
	}
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

// SetFinalSize sets the finalSize of 'sym'.
func SetFinalSize(prgrm *ast.CXProgram, symbols *[]map[string]*ast.CXArgument, symIdx ast.CXArgumentIndex) {
	sym := prgrm.GetCXArgFromArray(symIdx)
	finalSize := sym.TotalSize

	symPkg, err := prgrm.GetPackageFromArray(sym.Package)
	if err != nil {
		panic(err)
	}

	arg, err := lookupSymbol(prgrm, symPkg.Name, sym.Name, symbols)
	if err == nil {
		calculateFinalSize(prgrm, &finalSize, sym, arg)
		for _, fldIdx := range sym.Fields {
			fld := prgrm.GetCXArgFromArray(fldIdx)
			finalSize = fld.TotalSize
			calculateFinalSize(prgrm, &finalSize, fld, arg)
		}
	}

	sym.TotalSize = finalSize
}

// calculateFinalSize calculates final size of 'sym'.
func calculateFinalSize(prgrm *ast.CXProgram, finalSize *types.Pointer, sym *ast.CXArgument, arg *ast.CXArgument) {
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
						subSize = arg.StructType.GetStructSize(prgrm)
					}
				}

				*finalSize = subSize
			}
		}
	}
}
