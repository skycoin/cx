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

	fnInputs := fn.GetInputs(prgrm)
	if len(fnInputs) != len(inputs) {
		// it must be a method declaration
		// so we save the first input
		fn.Inputs.Fields = fn.Inputs.Fields[:1]
	} else {
		fn.Inputs.Fields = nil
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
func ProcessFunctionParameters(prgrm *ast.CXProgram, symbols *[]map[string]*ast.CXTypeSignature, symbolsScope *map[string]bool, offset *types.Pointer, fnIdx ast.CXFunctionIndex, params []ast.CXArgumentIndex) {
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
	var symbols *[]map[string]*ast.CXTypeSignature
	tmp := make([]map[string]*ast.CXTypeSignature, 0)
	symbols = &tmp
	*symbols = append(*symbols, make(map[string]*ast.CXTypeSignature))

	// symbolsScope only handles the difference between local and global scopes
	// local being function constrained variables, and global being global variables.
	var symbolsScope map[string]bool = make(map[string]bool)

	fn := prgrm.GetFunctionFromArray(fnIdx)

	FunctionAddParameters(prgrm, fnIdx, inputs, outputs)
	ProcessGoTos(prgrm, exprs)
	AddExprsToFunction(prgrm, fnIdx, exprs)

	var fnInputIdxs []ast.CXArgumentIndex
	for _, fnInput := range fn.GetInputs(prgrm) {
		if fnInput.Type == ast.TYPE_CXARGUMENT_DEPRECATE {
			fnInputIdxs = append(fnInputIdxs, ast.CXArgumentIndex(fnInput.Meta))
		} else {
			panic("type is not type cx argument deprecate\n\n")
		}
	}
	ProcessFunctionParameters(prgrm, symbols, &symbolsScope, &offset, fnIdx, fnInputIdxs)

	var fnOutputIdxs []ast.CXArgumentIndex
	for _, fnOutput := range fn.GetOutputs(prgrm) {
		if fnOutput.Type == ast.TYPE_CXARGUMENT_DEPRECATE {
			fnOutputIdxs = append(fnOutputIdxs, ast.CXArgumentIndex(fnOutput.Meta))
		} else {
			panic("type is not type cx argument deprecate\n\n")
		}
	}
	ProcessFunctionParameters(prgrm, symbols, &symbolsScope, &offset, fnIdx, fnOutputIdxs)

	for i, expr := range fn.Expressions {
		if expr.Type == ast.CX_LINE {
			continue
		}
		exprAtomicOp, err := prgrm.GetCXAtomicOp(expr.Index)
		if err != nil {
			panic(err)
		}

		if expr.IsScopeNew() {
			*symbols = append(*symbols, make(map[string]*ast.CXTypeSignature))
		}

		ProcessMethodCall(prgrm, &expr, symbols)

		ProcessExpressionArguments(prgrm, symbols, &symbolsScope, &offset, fnIdx, exprAtomicOp.GetInputs(prgrm), &expr, true)
		ProcessExpressionArguments(prgrm, symbols, &symbolsScope, &offset, fnIdx, exprAtomicOp.GetOutputs(prgrm), &expr, false)

		ProcessPointerStructs(prgrm, &expr)
		ProcessTempVariable(prgrm, &expr)
		ProcessSliceAssignment(prgrm, &expr)
		ProcessStringAssignment(prgrm, &expr)
		// ProcessReferenceAssignment(prgrm, &expr)
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
		atomicType := prgrm.CXArgs[expression.GetInputs(prgrm)[0].Meta].GetType(prgrm)
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
		cxAtomicOpOutput := prgrm.GetCXArgFromArray(ast.CXArgumentIndex(expression.GetOutputs(prgrm)[0].Meta))
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

		if len(expression.GetOutputs(prgrm)) > 0 && cxAtomicOpOutput.Fields == nil {
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

			typeSig := inpExprAtomicOp.GetOutputs(prgrm)[0]
			// then it's a literal
			expression.AddInput(prgrm, typeSig)
		} else {
			// then it's a function call
			if len(inpExprAtomicOp.GetOutputs(prgrm)) < 1 {
				var out *ast.CXArgument = &ast.CXArgument{}

				inpExprAtomicOpOperatorOutputs := inpExprAtomicOpOperator.GetOutputs(prgrm)
				inpExprAtomicOpOperatorOutput := prgrm.GetCXArgFromArray(ast.CXArgumentIndex(inpExprAtomicOpOperatorOutputs[0].Meta))
				if inpExprAtomicOpOperatorOutput.Type == types.UNDEFINED {
					// if undefined type, then adopt argument's type
					inpExprAtomicOpInput := prgrm.GetCXArgFromArray(ast.CXArgumentIndex(inpExprAtomicOp.GetInputs(prgrm)[0].Meta))

					out = ast.MakeArgument(generateTempVarName(constants.LOCAL_PREFIX), CurrentFile, inpExprCXLine.LineNumber).SetType(inpExprAtomicOpInput.Type)
					out.StructType = inpExprAtomicOpInput.StructType
					out.Size = inpExprAtomicOpInput.Size
					out.TotalSize = inpExprAtomicOpOperatorOutputs[0].GetSize(prgrm)
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
						out.TotalSize = inpExprAtomicOpOperatorOutputs[0].GetSize(prgrm)
					}

					out.Type = inpExprAtomicOpOperatorOutput.Type
					out.PointerTargetType = inpExprAtomicOpOperatorOutput.PointerTargetType
					out.PreviouslyDeclared = true
				}

				out.Package = inpExprAtomicOp.Package
				outIdx := prgrm.AddCXArgInArray(out)

				typeSig := ast.GetCXTypeSignatureRepresentationOfCXArg_ForGlobals_CXAtomicOps(prgrm, prgrm.GetCXArgFromArray(outIdx))
				inpExprAtomicOp.AddOutput(prgrm, typeSig)
				expression.AddInput(prgrm, typeSig)
			}
			if len(inpExprAtomicOp.GetOutputs(prgrm)) > 0 && inpExpr.IsArrayLiteral() {
				typeSig := inpExprAtomicOp.GetOutputs(prgrm)[0]
				expression.AddInput(prgrm, typeSig)
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

	if len(expression.GetInputs(prgrm)) < 1 {
		return errors.New("cannot perform arithmetic without operands")
	}

	typeCode := prgrm.GetCXArgFromArray(ast.CXArgumentIndex(expression.GetInputs(prgrm)[0].Meta)).Type
	if prgrm.GetCXArgFromArray(ast.CXArgumentIndex(expression.GetInputs(prgrm)[0].Meta)).Type == types.POINTER {
		typeCode = prgrm.GetCXArgFromArray(ast.CXArgumentIndex(expression.GetInputs(prgrm)[0].Meta)).PointerTargetType
	}

	for _, input := range expression.GetInputs(prgrm) {
		var inp *ast.CXArgument = &ast.CXArgument{}
		if input.Type == ast.TYPE_CXARGUMENT_DEPRECATE {
			inp = prgrm.GetCXArgFromArray(ast.CXArgumentIndex(input.Meta))
		} else {
			panic("type is not type cx argument deprecate\n\n")
		}

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
		for _, output := range expression.GetOutputs(prgrm) {
			var outIdx ast.CXArgumentIndex
			if output.Type == ast.TYPE_CXARGUMENT_DEPRECATE {
				outIdx = ast.CXArgumentIndex(output.Meta)
			} else {
				panic("type is not type cx argument deprecate\n\n")
			}
			out := prgrm.GetCXArgFromArray(outIdx)

			size := types.Pointer(1)
			if !ast.IsComparisonOperator(expressionOperator.AtomicOPCode) {
				size = ast.GetArgSize(prgrm, prgrm.GetCXArgFromArray(ast.CXArgumentIndex(expression.GetInputs(prgrm)[0].Meta)).GetAssignmentElement(prgrm))
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

	for _, argTypeSig := range append(expression.GetInputs(prgrm), expression.GetOutputs(prgrm)...) {
		var arg *ast.CXArgument = &ast.CXArgument{}
		if argTypeSig.Type == ast.TYPE_CXARGUMENT_DEPRECATE {
			arg = prgrm.GetCXArgFromArray(ast.CXArgumentIndex(argTypeSig.Meta))
		} else {
			panic("type is not type cx argument deprecate\n\n")
		}

		for _, fldIdx := range arg.Fields {
			fld := prgrm.GetCXArgFromArray(fldIdx)
			doesFieldHaveDerefPointer := false
			for _, deref := range fld.DereferenceOperations {
				if deref == constants.DEREF_POINTER {
					doesFieldHaveDerefPointer = true
				}
			}
			if fld.IsPointer() && !doesFieldHaveDerefPointer {
				prgrm.CXArgs[fldIdx].DereferenceOperations = append(prgrm.CXArgs[fldIdx].DereferenceOperations, constants.DEREF_POINTER)
			}
		}

		doesArgHaveDerefPointer := false
		for _, deref := range arg.DereferenceOperations {
			if deref == constants.DEREF_POINTER {
				doesArgHaveDerefPointer = true
			}
		}
		if arg.IsStruct && arg.IsPointer() && len(arg.Fields) > 0 && !doesArgHaveDerefPointer {
			prgrm.CXArgs[arg.Index].DereferenceOperations = append(arg.DereferenceOperations, constants.DEREF_POINTER)
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
			inp1Type := ast.GetFormattedType(prgrm, prgrm.GetCXArgFromArray(ast.CXArgumentIndex(expression.GetInputs(prgrm)[0].Meta)))
			inp2Type := ast.GetFormattedType(prgrm, prgrm.GetCXArgFromArray(ast.CXArgumentIndex(expression.GetInputs(prgrm)[1].Meta)))
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
func ProcessExpressionArguments(prgrm *ast.CXProgram, symbols *[]map[string]*ast.CXTypeSignature, symbolsScope *map[string]bool, offset *types.Pointer, fnIdx ast.CXFunctionIndex, args []*ast.CXTypeSignature, expr *ast.CXExpression, isInput bool) {
	fn := prgrm.GetFunctionFromArray(fnIdx)

	for _, typeSignature := range args {
		var argIdx ast.CXArgumentIndex
		if typeSignature.Type == ast.TYPE_CXARGUMENT_DEPRECATE {
			argIdx = ast.CXArgumentIndex(typeSignature.Meta)
		} else {
			panic(fmt.Sprintf("type is=%v\n\n", typeSignature.Type))
		}

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
func CheckRedeclared(prgrm *ast.CXProgram, symbols *[]map[string]*ast.CXTypeSignature, expr *ast.CXExpression, symIdx ast.CXArgumentIndex) {
	sym := prgrm.GetCXArgFromArray(symIdx)

	expression, err := prgrm.GetCXAtomicOp(expr.Index)
	if err != nil {
		panic(err)
	}
	expressionOperator := prgrm.GetFunctionFromArray(expression.Operator)

	if expressionOperator == nil && len(expression.GetOutputs(prgrm)) > 0 && len(expression.GetInputs(prgrm)) == 0 {
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
		expressionOperatorInputs := expressionOperator.GetInputs(prgrm)
		expressionOperatorOutputs := expressionOperator.GetOutputs(prgrm)

		opGotoFn := ast.Natives[constants.OP_GOTO]
		isOpGoto := false
		if expressionOperator != nil && expressionOperator.AtomicOPCode == opGotoFn.AtomicOPCode && len(expressionOperatorInputs) == len(opGotoFn.Inputs) && len(expressionOperatorOutputs) == len(opGotoFn.Outputs) {
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
				println(ast.CompilationError(prgrm.GetCXArgFromArray(ast.CXArgumentIndex(expression.GetOutputs(prgrm)[i].Meta)).ArgDetails.FileName, prgrm.GetCXArgFromArray(ast.CXArgumentIndex(expression.GetOutputs(prgrm)[i].Meta)).ArgDetails.FileLine), fmt.Sprintf("function '%s' expected receiving variable of type '%s'; '%s' was provided", opName, expectedType, receivedType))
			}

		}

		// In the case of assignment we need to check that the input's type matches the output's type.
		// FIXME: There are some expressions added by the cxgo where temporary variables are used.
		// These temporary variables' types are not properly being set. That's why we use !cxcore.IsTempVar to
		// exclude these cases for now.
		if expressionOperator.AtomicOPCode == constants.OP_IDENTITY && !IsTempVar(expression.GetOutputs(prgrm)[0].Name) {
			inpType := ast.GetFormattedType(prgrm, prgrm.GetCXArgFromArray(ast.CXArgumentIndex(expression.GetInputs(prgrm)[0].Meta)))
			outType := ast.GetFormattedType(prgrm, prgrm.GetCXArgFromArray(ast.CXArgumentIndex(expression.GetOutputs(prgrm)[0].Meta)))

			// We use `isInputs` to only print the error once.
			// Otherwise we'd print the error twice: once for the input and again for the output
			if inpType != outType && isInputs {
				println(ast.CompilationError(receivedArg.ArgDetails.FileName, receivedArg.ArgDetails.FileLine), fmt.Sprintf("cannot assign value of type '%s' to identifier '%s' of type '%s'", inpType, prgrm.GetCXArgFromArray(ast.CXArgumentIndex(expression.GetOutputs(prgrm)[0].Meta)).GetAssignmentElement(prgrm).Name, outType))
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
	expressionOperatorInputs := prgrm.ConvertIndexTypeSignaturesToPointerArgs(expressionOperator.GetInputs(prgrm))

	exprCXLine, _ := prgrm.GetPreviousCXLine(exprs, currIndex)

	if expressionOperator != nil {
		opName := expression.GetOperatorName(prgrm)

		// checking if number of inputs is less than the required number of inputs
		if len(expression.GetInputs(prgrm)) != len(expressionOperatorInputs) {
			if !(len(expressionOperatorInputs) > 0 && expressionOperatorInputs[len(expressionOperatorInputs)-1].Type != types.UNDEFINED) {
				// if the last input is of type cxcore.TYPE_UNDEFINED then it might be a variadic function, such as printf
			} else {
				// then we need to be strict in the number of inputs
				var plural1 string
				var plural2 string = "s"
				var plural3 string = "were"
				if len(expressionOperatorInputs) > 1 {
					plural1 = "s"
				}
				if len(expression.GetInputs(prgrm)) == 1 {
					plural2 = ""
					plural3 = "was"
				}

				println(ast.CompilationError(exprCXLine.FileName, exprCXLine.LineNumber), fmt.Sprintf("operator '%s' expects %d input%s, but %d input argument%s %s provided", opName, len(expressionOperatorInputs), plural1, len(expression.GetInputs(prgrm)), plural2, plural3))
				return
			}
		}
		expressionOperatorOutputs := expressionOperator.GetOutputs(prgrm)
		// checking if number of expr.ProgramOutput matches number of Operator.ProgramOutput
		if len(expression.GetOutputs(prgrm)) != len(expressionOperatorOutputs) {
			var plural1 string
			var plural2 string = "s"
			var plural3 string = "were"
			if len(expressionOperatorOutputs) > 1 {
				plural1 = "s"
			}
			if len(expression.GetOutputs(prgrm)) == 1 {
				plural2 = ""
				plural3 = "was"
			}

			println(ast.CompilationError(exprCXLine.FileName, exprCXLine.LineNumber), fmt.Sprintf("operator '%s' expects to return %d output%s, but %d receiving argument%s %s provided", opName, len(expressionOperatorOutputs), plural1, len(expression.GetOutputs(prgrm)), plural2, plural3))
			os.Exit(constants.CX_COMPILATION_ERROR)
		}
	}

	if expressionOperator != nil && expressionOperator.IsBuiltIn() && expressionOperator.AtomicOPCode == constants.OP_IDENTITY {
		for i := range expression.GetInputs(prgrm) {
			var expectedType string
			var receivedType string
			if prgrm.GetCXArgFromArray(ast.CXArgumentIndex(expression.GetOutputs(prgrm)[i].Meta)).GetAssignmentElement(prgrm).StructType != nil {
				// then it's custom type
				expectedType = prgrm.GetCXArgFromArray(ast.CXArgumentIndex(expression.GetOutputs(prgrm)[i].Meta)).GetAssignmentElement(prgrm).StructType.Name

			} else {
				// then it's native type
				expectedType = prgrm.GetCXArgFromArray(ast.CXArgumentIndex(expression.GetOutputs(prgrm)[i].Meta)).GetAssignmentElement(prgrm).Type.Name()

				if prgrm.GetCXArgFromArray(ast.CXArgumentIndex(expression.GetOutputs(prgrm)[i].Meta)).GetAssignmentElement(prgrm).Type == types.POINTER {
					expectedType = prgrm.GetCXArgFromArray(ast.CXArgumentIndex(expression.GetOutputs(prgrm)[i].Meta)).GetAssignmentElement(prgrm).PointerTargetType.Name()
				}
			}

			if prgrm.GetCXArgFromArray(ast.CXArgumentIndex(expression.GetInputs(prgrm)[i].Meta)).GetAssignmentElement(prgrm).StructType != nil {
				// then it's custom type
				receivedType = prgrm.GetCXArgFromArray(ast.CXArgumentIndex(expression.GetInputs(prgrm)[i].Meta)).GetAssignmentElement(prgrm).StructType.Name
			} else {
				// then it's native type
				receivedType = prgrm.GetCXArgFromArray(ast.CXArgumentIndex(expression.GetInputs(prgrm)[i].Meta)).GetAssignmentElement(prgrm).Type.Name()

				if prgrm.GetCXArgFromArray(ast.CXArgumentIndex(expression.GetInputs(prgrm)[i].Meta)).GetAssignmentElement(prgrm).Type == types.POINTER {
					receivedType = prgrm.GetCXArgFromArray(ast.CXArgumentIndex(expression.GetInputs(prgrm)[i].Meta)).GetAssignmentElement(prgrm).PointerTargetType.Name()
				}
			}

			if receivedType != expectedType {
				if exprs[currIndex].IsStructLiteral() {
					println(ast.CompilationError(prgrm.GetCXArgFromArray(ast.CXArgumentIndex(expression.GetOutputs(prgrm)[i].Meta)).ArgDetails.FileName, prgrm.GetCXArgFromArray(ast.CXArgumentIndex(expression.GetOutputs(prgrm)[i].Meta)).ArgDetails.FileLine), fmt.Sprintf("field '%s' in struct literal of type '%s' expected argument of type '%s'; '%s' was provided", prgrm.GetCXArgFromArray(prgrm.GetCXArgFromArray(ast.CXArgumentIndex(expression.GetOutputs(prgrm)[i].Meta)).Fields[0]).Name, prgrm.GetCXArgFromArray(ast.CXArgumentIndex(expression.GetOutputs(prgrm)[i].Meta)).StructType.Name, expectedType, receivedType))
				} else {
					println(ast.CompilationError(prgrm.GetCXArgFromArray(ast.CXArgumentIndex(expression.GetOutputs(prgrm)[i].Meta)).ArgDetails.FileName, prgrm.GetCXArgFromArray(ast.CXArgumentIndex(expression.GetOutputs(prgrm)[i].Meta)).ArgDetails.FileLine), fmt.Sprintf("trying to assign argument of type '%s' to symbol '%s' of type '%s'", receivedType, prgrm.GetCXArgFromArray(ast.CXArgumentIndex(expression.GetOutputs(prgrm)[i].Meta)).GetAssignmentElement(prgrm).Name, expectedType))
				}
			}
		}
	}

	// then it's a function call and not a declaration
	if expressionOperator != nil {
		var expressionOperatorInputsIdxs []ast.CXArgumentIndex
		for _, input := range expressionOperatorInputs {
			expressionOperatorInputsIdxs = append(expressionOperatorInputsIdxs, ast.CXArgumentIndex(input.Index))
		}

		var expressionInputsIdxs []ast.CXArgumentIndex
		for _, input := range expression.GetInputs(prgrm) {
			if input.Type == ast.TYPE_CXARGUMENT_DEPRECATE {
				expressionInputsIdxs = append(expressionInputsIdxs, ast.CXArgumentIndex(input.Meta))
			} else {
				panic("type is not type cx argument deprecate\n\n")
			}
		}
		// checking inputs matching operator's inputs
		checkMatchParamTypes(prgrm, &exprs[currIndex], expressionOperatorInputsIdxs, expressionInputsIdxs, true)

		expressionOperatorOutputs := expressionOperator.GetOutputs(prgrm)
		var expressionOperatorOutputsIdxs []ast.CXArgumentIndex
		for _, output := range expressionOperatorOutputs {
			if output.Type == ast.TYPE_CXARGUMENT_DEPRECATE {
				expressionOperatorOutputsIdxs = append(expressionOperatorOutputsIdxs, ast.CXArgumentIndex(output.Meta))
			} else {
				panic("type is not type cx argument deprecate\n\n")
			}
		}

		expressionOutputs := expression.GetOutputs(prgrm)
		var expressionOutputsIdxs []ast.CXArgumentIndex
		for _, output := range expressionOutputs {
			if output.Type == ast.TYPE_CXARGUMENT_DEPRECATE {
				expressionOutputsIdxs = append(expressionOutputsIdxs, ast.CXArgumentIndex(output.Meta))
			} else {
				panic("type is not type cx argument deprecate\n\n")
			}
		}

		// checking outputs matching operator's outputs
		checkMatchParamTypes(prgrm, &exprs[currIndex], expressionOperatorOutputsIdxs, expressionOutputsIdxs, false)
	}
}

// ProcessStringAssignment sets the args PassBy to PASSBY_VALUE if the type is string.
func ProcessStringAssignment(prgrm *ast.CXProgram, expr *ast.CXExpression) {
	expression, err := prgrm.GetCXAtomicOp(expr.Index)
	if err != nil {
		panic(err)
	}
	expressionOperator := prgrm.GetFunctionFromArray(expression.Operator)
	expressionOperatorInputs := expressionOperator.GetInputs(prgrm)
	expressionOperatorOutputs := expressionOperator.GetOutputs(prgrm)

	opIdentFn := ast.Natives[constants.OP_IDENTITY]
	isOpIdent := false
	if expressionOperator != nil && expressionOperator.AtomicOPCode == opIdentFn.AtomicOPCode && len(expressionOperatorInputs) == len(opIdentFn.Inputs) && len(expressionOperatorOutputs) == len(opIdentFn.Outputs) {
		isOpIdent = true
	}
	if isOpIdent {
		for i, output := range expression.GetOutputs(prgrm) {
			var out *ast.CXArgument = &ast.CXArgument{}
			if output.Type == ast.TYPE_CXARGUMENT_DEPRECATE {
				out = prgrm.GetCXArgFromArray(ast.CXArgumentIndex(output.Meta))
			} else {
				panic("type is not type cx argument deprecate\n\n")
			}

			if len(expression.GetInputs(prgrm)) > i {
				out = out.GetAssignmentElement(prgrm)
				inp := prgrm.GetCXArgFromArray(ast.CXArgumentIndex(expression.GetInputs(prgrm)[i].Meta)).GetAssignmentElement(prgrm)

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
// func ProcessReferenceAssignment(prgrm *ast.CXProgram, expr *ast.CXExpression) {
// 	expression, err := prgrm.GetCXAtomicOp(expr.Index)
// 	if err != nil {
// 		panic(err)
// 	}

// 	for _, outIdx := range expression.Outputs {
// 		out := prgrm.GetCXArgFromArray(outIdx)
// 		elt := out.GetAssignmentElement(prgrm)
// 		if elt.PassBy == constants.PASSBY_REFERENCE &&
// 			!hasDeclSpec(elt, constants.DECL_POINTER) &&
// 			elt.PointerTargetType != types.STR && elt.Type != types.STR && !elt.IsSlice {
// 			println(ast.CompilationError(CurrentFile, LineNo), "invalid reference assignment", elt.Name)
// 		}
// 	}

// }

// ProcessShortDeclaration sets proper values if the expr is a short declaration.
func ProcessShortDeclaration(prgrm *ast.CXProgram, expr *ast.CXExpression, expressions []ast.CXExpression, idx int) {
	if len(expressions) <= 1 {
		return
	}

	expression, err := prgrm.GetCXAtomicOp(expr.Index)
	if err != nil {
		panic(err)
	}

	// process short declaration
	if len(expression.GetOutputs(prgrm)) > 0 && len(expression.GetInputs(prgrm)) > 0 && prgrm.GetCXArgFromArray(ast.CXArgumentIndex(expression.GetOutputs(prgrm)[0].Meta)).Type == types.IDENTIFIER && !expr.IsStructLiteral() && !isParseOp(prgrm, expr) {
		expressionOperator := prgrm.GetFunctionFromArray(expression.Operator)
		expressionOperatorOutputs := expressionOperator.GetOutputs(prgrm)
		prevExpression, err := prgrm.GetPreviousCXAtomicOpFromExpressions(expressions, idx-1)
		if err != nil {
			panic(err)
		}

		var arg *ast.CXArgument = &ast.CXArgument{}
		if expr.IsMethodCall() {
			arg = prgrm.GetCXArgFromArray(ast.CXArgumentIndex(expressionOperatorOutputs[0].Meta))
		} else {
			arg = prgrm.GetCXArgFromArray(ast.CXArgumentIndex(expression.GetInputs(prgrm)[0].Meta))
		}

		prevExpressionOutput := prevExpression.GetOutputs(prgrm)[0]
		var prevExpressionOutputIdx ast.CXArgumentIndex
		if prevExpressionOutput.Type == ast.TYPE_CXARGUMENT_DEPRECATE {
			prevExpressionOutputIdx = ast.CXArgumentIndex(prevExpressionOutput.Meta)
		} else {
			panic("type is not type cx argument deprecate\n\n")
		}

		prgrm.CXArgs[prevExpressionOutputIdx].Type = arg.Type
		prgrm.CXArgs[prevExpressionOutputIdx].PointerTargetType = arg.PointerTargetType
		prgrm.CXArgs[prevExpressionOutputIdx].Size = arg.Size
		prgrm.CXArgs[prevExpressionOutputIdx].TotalSize = arg.TotalSize

		expressionOutput := expression.GetOutputs(prgrm)[0]
		var expressionOutputIdx ast.CXArgumentIndex
		if expressionOutput.Type == ast.TYPE_CXARGUMENT_DEPRECATE {
			expressionOutputIdx = ast.CXArgumentIndex(expressionOutput.Meta)
		} else {
			panic("type is not type cx argument deprecate\n\n")
		}

		prgrm.CXArgs[expressionOutputIdx].Type = arg.Type
		prgrm.CXArgs[expressionOutputIdx].PointerTargetType = arg.PointerTargetType
		prgrm.CXArgs[expressionOutputIdx].Size = arg.Size
		prgrm.CXArgs[expressionOutputIdx].TotalSize = arg.TotalSize
	}
}

// ProcessSlice sets DereferenceOperations if the arg is a slice.
func ProcessSlice(prgrm *ast.CXProgram, inpIdx ast.CXArgumentIndex) {
	inp := prgrm.GetCXArgFromArray(inpIdx)

	var elt *ast.CXArgument = &ast.CXArgument{}
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
	expressionOperatorInputs := expressionOperator.GetInputs(prgrm)
	expressionOperatorOutputs := expressionOperator.GetOutputs(prgrm)

	opIdentFn := ast.Natives[constants.OP_IDENTITY]
	isOpIdent := false
	if expressionOperator != nil && expressionOperator.AtomicOPCode == opIdentFn.AtomicOPCode && len(expressionOperatorInputs) == len(opIdentFn.Inputs) && len(expressionOperatorOutputs) == len(opIdentFn.Outputs) {
		isOpIdent = true
	}

	if isOpIdent {
		var inp *ast.CXArgument = &ast.CXArgument{}
		var out *ast.CXArgument = &ast.CXArgument{}

		inp = prgrm.GetCXArgFromArray(ast.CXArgumentIndex(expression.GetInputs(prgrm)[0].Meta)).GetAssignmentElement(prgrm)
		out = prgrm.GetCXArgFromArray(ast.CXArgumentIndex(expression.GetOutputs(prgrm)[0].Meta)).GetAssignmentElement(prgrm)

		if inp.IsSlice && out.IsSlice && len(inp.Indexes) == 0 && len(out.Indexes) == 0 {
			out.PassBy = constants.PASSBY_VALUE
		}
	}
	if expressionOperator != nil && !expressionOperator.IsBuiltIn() {
		// then it's a function call
		for _, input := range expression.GetInputs(prgrm) {
			var inp *ast.CXArgument = &ast.CXArgument{}
			if input.Type == ast.TYPE_CXARGUMENT_DEPRECATE {
				inp = prgrm.GetCXArgFromArray(ast.CXArgumentIndex(input.Meta))
			} else {
				panic("type is not type cx argument deprecate\n\n")
			}

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
func lookupSymbol(prgrm *ast.CXProgram, pkgName, ident string, symbols *[]map[string]*ast.CXTypeSignature) (*ast.CXTypeSignature, error) {
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

	fnArgIdx := prgrm.AddCXArgInArray(fnArg)
	fnArgTypeSig := &ast.CXTypeSignature{
		Name:    fnArg.Name,
		Offset:  fnArg.Offset,
		Package: fnArg.Package,
		Type:    ast.TYPE_CXARGUMENT_DEPRECATE,
		Meta:    int(fnArgIdx),
	}

	return fnArgTypeSig, nil
}

// UpdateSymbolsTable adds `sym` to the innermost scope (last element of slice) in `symbols`.
func UpdateSymbolsTable(prgrm *ast.CXProgram, symbols *[]map[string]*ast.CXTypeSignature, symIdx ast.CXArgumentIndex, offset *types.Pointer, shouldExist bool) {
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

			symTypeSig := &ast.CXTypeSignature{
				Name:    sym.Name,
				Offset:  sym.Offset,
				Package: sym.Package,
				Type:    ast.TYPE_CXARGUMENT_DEPRECATE,
				Meta:    int(symIdx),
			}
			(*symbols)[lastIdx][fullName] = symTypeSig
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
func ProcessMethodCall(prgrm *ast.CXProgram, expr *ast.CXExpression, symbols *[]map[string]*ast.CXTypeSignature) {
	expression, err := prgrm.GetCXAtomicOp(expr.Index)
	if err != nil {
		panic(err)
	}

	if expr.IsMethodCall() {
		var inpIdx ast.CXArgumentIndex = -1
		var outIdx ast.CXArgumentIndex = -1

		if len(expression.GetInputs(prgrm)) > 0 && expression.GetInputs(prgrm)[0].Name != "" {
			inpIdx = ast.CXArgumentIndex(expression.GetInputs(prgrm)[0].Meta)
		}
		if len(expression.GetOutputs(prgrm)) > 0 && expression.GetOutputs(prgrm)[0].Name != "" {
			outIdx = ast.CXArgumentIndex(expression.GetOutputs(prgrm)[0].Meta)
		}

		if inpIdx != -1 {
			inpPkg, err := prgrm.GetPackageFromArray(prgrm.CXArgs[inpIdx].Package)
			if err != nil {
				panic(err)
			}

			if argInpTypeSignature, err := lookupSymbol(prgrm, inpPkg.Name, prgrm.CXArgs[inpIdx].Name, symbols); err != nil {
				if outIdx == -1 {
					panic("")
				}

				outPkg, err := prgrm.GetPackageFromArray(prgrm.CXArgs[outIdx].Package)
				if err != nil {
					panic(err)
				}

				argOutTypeSignature, err := lookupSymbol(prgrm, outPkg.Name, prgrm.CXArgs[outIdx].Name, symbols)
				if err != nil {
					println(ast.CompilationError(prgrm.CXArgs[outIdx].ArgDetails.FileName, prgrm.CXArgs[outIdx].ArgDetails.FileLine), fmt.Sprintf("identifier '%s' does not exist", prgrm.CXArgs[outIdx].Name))
					os.Exit(constants.CX_COMPILATION_ERROR)
				}
				// then we found an output
				if len(prgrm.CXArgs[outIdx].Fields) > 0 {
					var argOut *ast.CXArgument = &ast.CXArgument{}
					if argOutTypeSignature.Type == ast.TYPE_CXARGUMENT_DEPRECATE {
						argOut = prgrm.GetCXArgFromArray(ast.CXArgumentIndex(argOutTypeSignature.Meta))
					} else {
						panic("type is cx argument deprecate\n\n")
					}

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

					typeSig := ast.GetCXTypeSignatureRepresentationOfCXArg_ForGlobals_CXAtomicOps(prgrm, &prgrm.CXArgs[outIdx])

					prgrm.CXAtomicOps[expr.Index].Inputs.Fields = append([]*ast.CXTypeSignature{typeSig}, prgrm.CXAtomicOps[expr.Index].Inputs.Fields...)
					prgrm.CXAtomicOps[expr.Index].Outputs.Fields = prgrm.CXAtomicOps[expr.Index].Outputs.Fields[1:]
					prgrm.CXArgs[outIdx].Fields = prgrm.CXArgs[outIdx].Fields[:len(prgrm.CXArgs[outIdx].Fields)-1]
				}
			} else {
				var argInp *ast.CXArgument = &ast.CXArgument{}
				if argInpTypeSignature.Type == ast.TYPE_CXARGUMENT_DEPRECATE {
					argInp = prgrm.GetCXArgFromArray(ast.CXArgumentIndex(argInpTypeSignature.Meta))
				} else {
					panic("type is cx argument deprecate\n\n")
				}

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
					argOutTypeSignature, err := lookupSymbol(prgrm, outPkg.Name, prgrm.CXArgs[outIdx].Name, symbols)
					if err != nil {
						panic(err)
					}

					var argOut *ast.CXArgument = &ast.CXArgument{}
					if argOutTypeSignature.Type == ast.TYPE_CXARGUMENT_DEPRECATE {
						argOut = prgrm.GetCXArgFromArray(ast.CXArgumentIndex(argOutTypeSignature.Meta))
					} else {
						panic("type is cx argument deprecate\n\n")
					}

					strct := argOut.StructType
					if strct == nil {
						println(ast.CompilationError(argOut.ArgDetails.FileName, argOut.ArgDetails.FileLine), fmt.Sprintf("illegal method call or field access on identifier '%s' of primitive type '%s'", argOut.Name, argOut.Type.Name()))
						os.Exit(constants.CX_COMPILATION_ERROR)
					}

					prgrm.CXAtomicOps[expr.Index].Inputs.Fields = append(prgrm.CXAtomicOps[expr.Index].Outputs.Fields[:1], prgrm.CXAtomicOps[expr.Index].Inputs.Fields...)
					prgrm.CXAtomicOps[expr.Index].Outputs.Fields = prgrm.CXAtomicOps[expr.Index].Outputs.Fields[:len(prgrm.CXAtomicOps[expr.Index].Outputs.Fields)-1]

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

			argOutTypeSignature, err := lookupSymbol(prgrm, outPkg.Name, prgrm.CXArgs[outIdx].Name, symbols)
			if err != nil {
				println(ast.CompilationError(prgrm.CXArgs[outIdx].ArgDetails.FileName, prgrm.CXArgs[outIdx].ArgDetails.FileLine), fmt.Sprintf("identifier '%s' does not exist", prgrm.CXArgs[outIdx].Name))
				os.Exit(constants.CX_COMPILATION_ERROR)
			}

			var argOut *ast.CXArgument = &ast.CXArgument{}
			if argOutTypeSignature.Type == ast.TYPE_CXARGUMENT_DEPRECATE {
				argOut = prgrm.GetCXArgFromArray(ast.CXArgumentIndex(argOutTypeSignature.Meta))
			} else {
				panic("type is cx argument deprecate\n\n")
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

				typeSig := ast.GetCXTypeSignatureRepresentationOfCXArg_ForGlobals_CXAtomicOps(prgrm, &prgrm.CXArgs[outIdx])
				newInputs := &ast.CXStruct{Fields: []*ast.CXTypeSignature{typeSig}}
				if prgrm.CXAtomicOps[expr.Index].Inputs != nil {
					for _, typeSig := range prgrm.CXAtomicOps[expr.Index].Inputs.Fields {
						newInputs.AddField_CXAtomicOps(prgrm, typeSig)
					}
				}
				prgrm.CXAtomicOps[expr.Index].Inputs = newInputs
				prgrm.CXAtomicOps[expr.Index].Outputs.Fields = prgrm.CXAtomicOps[expr.Index].Outputs.Fields[1:]
				prgrm.CXArgs[outIdx].Fields = prgrm.CXArgs[outIdx].Fields[:len(prgrm.CXArgs[outIdx].Fields)-1]
			}
		}
	}
}

// GiveOffset
func GiveOffset(prgrm *ast.CXProgram, symbols *[]map[string]*ast.CXTypeSignature, symIdx ast.CXArgumentIndex, shouldExist bool) {
	sym := prgrm.GetCXArgFromArray(symIdx)

	if sym.Name != "" {
		symPkg, err := prgrm.GetPackageFromArray(sym.Package)
		if err != nil {
			panic(err)
		}

		if !sym.IsLocalDeclaration {
			GetGlobalSymbol(prgrm, symbols, symPkg, sym.Name)
		}

		argTypeSignature, err := lookupSymbol(prgrm, symPkg.Name, sym.Name, symbols)
		if err == nil {
			var arg *ast.CXArgument = &ast.CXArgument{}
			if argTypeSignature.Type == ast.TYPE_CXARGUMENT_DEPRECATE {
				arg = prgrm.GetCXArgFromArray(ast.CXArgumentIndex(argTypeSignature.Meta))
			} else {
				panic("type is type cxargument deprecate\n\n")
			}

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
	expressionOperatorInputs := expressionOperator.GetInputs(prgrm)
	expressionOperatorOutputs := expressionOperator.GetOutputs(prgrm)

	opIdentFn := ast.Natives[constants.OP_IDENTITY]
	isOpIdent := false
	if expressionOperator != nil && expressionOperator.AtomicOPCode == opIdentFn.AtomicOPCode && len(expressionOperatorInputs) == len(opIdentFn.Inputs) && len(expressionOperatorOutputs) == len(opIdentFn.Outputs) {
		isOpIdent = true
	}

	if expressionOperator != nil && (isOpIdent || ast.IsArithmeticOperator(expressionOperator.AtomicOPCode)) && len(expression.GetOutputs(prgrm)) > 0 && len(expression.GetInputs(prgrm)) > 0 {
		outputArg := expression.GetOutputs(prgrm)[0]
		var outputArgIdx ast.CXArgumentIndex
		if outputArg.Type == ast.TYPE_CXARGUMENT_DEPRECATE {
			outputArgIdx = ast.CXArgumentIndex(outputArg.Meta)
		} else {
			panic("type is not type cx argument deprecate\n\n")
		}
		name := prgrm.CXArgs[outputArgIdx].Name
		if IsTempVar(name) {
			expressionInput := prgrm.GetCXArgFromArray(ast.CXArgumentIndex(expression.GetInputs(prgrm)[0].Meta))
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
					methodOutputs := method.GetOutputs(prgrm)
					field.Type = prgrm.GetCXArgFromArray(ast.CXArgumentIndex(methodOutputs[0].Meta)).Type
					field.PointerTargetType = prgrm.GetCXArgFromArray(ast.CXArgumentIndex(methodOutputs[0].Meta)).PointerTargetType
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
				} else if nameField.Name == typeSignature.Name && typeSignature.Type == ast.TYPE_POINTER_ATOMIC {
					nameField.Type = types.POINTER
					nameField.PointerTargetType = types.Code(typeSignature.Meta)
					nameField.StructType = nil
					nameField.Size = types.Code(typeSignature.Meta).Size()
					nameField.TotalSize = typeSignature.GetSize(prgrm)

					nameField.DereferenceOperations = append([]int{constants.DEREF_POINTER}, nameField.DereferenceOperations...)

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

					// TODO: this should not be needed.
					if len(nameField.DeclarationSpecifiers) > 0 {
						nameField.DeclarationSpecifiers = append([]int{constants.DECL_BASIC, constants.DECL_SLICE}, nameField.DeclarationSpecifiers[1:]...)
					} else {
						nameField.DeclarationSpecifiers = []int{constants.DECL_BASIC}
					}

					nameField.DereferenceOperations = append([]int{constants.DEREF_POINTER}, nameField.DereferenceOperations...)

					break
				} else if nameField.Name == typeSignature.Name && typeSignature.Type == ast.TYPE_STRUCT {
					nameField.Type = types.STRUCT
					nameField.StructType = prgrm.GetStructFromArray(ast.CXStructIndex(typeSignature.Meta))
					if nameField.StructType != nil {
						strct = nameField.StructType
					}

					// TODO: this should not be needed.
					if len(nameField.DeclarationSpecifiers) > 0 {
						nameField.DeclarationSpecifiers = append([]int{constants.DECL_STRUCT}, nameField.DeclarationSpecifiers[1:]...)
					} else {
						nameField.DeclarationSpecifiers = []int{constants.DECL_STRUCT}
					}

					break
				}

				fieldIdx := typeSignature.Meta
				field := prgrm.CXArgs[fieldIdx]

				if nameField.Name == field.Name && typeSignature.Type == ast.TYPE_CXARGUMENT_DEPRECATE {
					nameField.Type = field.Type
					nameField.Lengths = field.Lengths
					nameField.Size = field.Size
					nameField.TotalSize = field.TotalSize
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
					}

					nameField.PassBy = field.PassBy
					nameField.IsSlice = field.IsSlice

					if field.StructType != nil {
						strct = field.StructType
					}

					break
				}

				nameField.Offset += typeSignature.Offset
			}
		}
	}
}

// GetGlobalSymbol tries to retrieve `ident` from `symPkg`'s globals if `ident` is not found in the local scope.
func GetGlobalSymbol(prgrm *ast.CXProgram, symbols *[]map[string]*ast.CXTypeSignature, symPkg *ast.CXPackage, ident string) {
	_, err := lookupSymbol(prgrm, symPkg.Name, ident, symbols)
	if err != nil {
		if glbl, err := symPkg.GetGlobal(prgrm, ident); err == nil {
			lastIdx := len(*symbols) - 1
			(*symbols)[lastIdx][symPkg.Name+"."+ident] = glbl
		}
	}
}

// SetFinalSize sets the finalSize of 'sym'.
func SetFinalSize(prgrm *ast.CXProgram, symbols *[]map[string]*ast.CXTypeSignature, symIdx ast.CXArgumentIndex) {
	sym := prgrm.GetCXArgFromArray(symIdx)
	finalSize := sym.TotalSize

	symPkg, err := prgrm.GetPackageFromArray(sym.Package)
	if err != nil {
		panic(err)
	}

	argTypeSignature, err := lookupSymbol(prgrm, symPkg.Name, sym.Name, symbols)
	if err == nil {
		var arg *ast.CXArgument = &ast.CXArgument{}
		if argTypeSignature.Type == ast.TYPE_CXARGUMENT_DEPRECATE {
			arg = prgrm.GetCXArgFromArray(ast.CXArgumentIndex(argTypeSignature.Meta))
		} else {
			panic("type not cx argument deprecate\n\n")
		}

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
