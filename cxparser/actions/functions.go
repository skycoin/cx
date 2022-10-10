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

type SymbolsData struct {
	varCount     int
	symbols      []*ast.CXTypeSignature
	symbolsIndex []map[string]int
}

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
				pkg.CurrentFunction = fnIdx

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
				pkg.CurrentFunction = fnIdx

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
		fn.AddLocalVariableName(inp.Name)
		fn.AddInput(prgrm, inp)
	}

	for _, out := range outputs {
		fn.AddLocalVariableName(out.Name)
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
//  offset - offset to use by statements (excluding inputs, outputs
// 			 and receiver).
//  fnIdx - the index of the function in the main CXFunction array.
//  params - function parameters to be processed.
func ProcessFunctionParameters(prgrm *ast.CXProgram, symbolsData *SymbolsData, offset *types.Pointer, fnIdx ast.CXFunctionIndex, params []ast.CXTypeSignatureIndex) {
	fn := prgrm.GetFunctionFromArray(fnIdx)

	for _, paramIdx := range params {
		UpdateSymbolsTable(prgrm, symbolsData, paramIdx, offset, false)
		// We remove it from local var array since it is already added to the symbols
		// if fn.IsLocalVariable(typeSignatureName) {
		// 	err := fn.RemoveLocalVariableFromArray(typeSignatureName)
		// 	if err != nil {
		// 		panic(err)
		// 	}
		// }

		GiveOffset(prgrm, symbolsData, paramIdx)
		AddPointer(prgrm, fn, paramIdx)

		typeSig := prgrm.GetCXTypeSignatureFromArray(paramIdx)
		var param *ast.CXArgument
		if typeSig.Type == ast.TYPE_CXARGUMENT_DEPRECATE {
			param = prgrm.GetCXArgFromArray(ast.CXArgumentIndex(typeSig.Meta))
		} else {
			continue
		}

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

	symIndex := make([]map[string]int, 0)
	symIndex = append(symIndex, make(map[string]int))

	symbolsData := &SymbolsData{
		varCount: 1,

		// symbols is a slice of string-int map which corresponds to
		// the scoping of the local variables. Each element in the slice
		// corresponds to a different scope. The innermost scope is
		// the last element of the slice.
		symbolsIndex: symIndex,
		symbols:      make([]*ast.CXTypeSignature, 0),
	}

	pkg, err := prgrm.GetCurrentPackage()
	if err != nil {
		panic(err)
	}

	pkg.CurrentFunction = fnIdx
	fn := prgrm.GetFunctionFromArray(fnIdx)

	FunctionAddParameters(prgrm, fnIdx, inputs, outputs)
	ProcessGoTos(prgrm, exprs)
	AddExprsToFunction(prgrm, fnIdx, exprs)
	ProcessFunctionParameters(prgrm, symbolsData, &offset, fnIdx, fn.GetInputs(prgrm))
	ProcessFunctionParameters(prgrm, symbolsData, &offset, fnIdx, fn.GetOutputs(prgrm))

	for i, expr := range fn.Expressions {
		if expr.Type == ast.CX_LINE {
			continue
		}
		exprAtomicOp, err := prgrm.GetCXAtomicOp(expr.Index)
		if err != nil {
			panic(err)
		}

		if expr.IsScopeNew() {
			symbolsData.symbolsIndex = append(symbolsData.symbolsIndex, make(map[string]int))
		}

		ProcessMethodCall(prgrm, &expr, symbolsData)

		ProcessExpressionArguments(prgrm, symbolsData, &offset, fnIdx, exprAtomicOp.GetInputs(prgrm), &expr, true)
		ProcessExpressionArguments(prgrm, symbolsData, &offset, fnIdx, exprAtomicOp.GetOutputs(prgrm), &expr, false)

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
			symbolsData.symbolsIndex = (symbolsData.symbolsIndex)[:len(symbolsData.symbolsIndex)-1]
		}
	}

	fn.LineCount = len(fn.Expressions)
	fn.Size = offset

	errStr := "\n"
	errStr += "functionName=" + fn.Name + "\n"
	for i := range exprs {
		expr, _ := prgrm.GetCXAtomicOpFromExpressions(exprs, i)
		if exprs[i].Type == ast.CX_LINE {
			continue
		}

		for x, inp := range expr.GetInputs(prgrm) {
			typeSig := prgrm.GetCXTypeSignatureFromArray(inp)
			if typeSig.Type == ast.TYPE_CXARGUMENT_DEPRECATE {
				errStr += fmt.Sprintf("expr[%v] inp[%v] = %+v\n\n", i, x, prgrm.GetCXArgFromArray(ast.CXArgumentIndex(typeSig.Meta)))
			} else {
				errStr += fmt.Sprintf("expr[%v] inp[%v] = %+v\n", i, x, typeSig)
			}

		}

		for y, out := range expr.GetOutputs(prgrm) {
			typeSig := prgrm.GetCXTypeSignatureFromArray(out)

			if typeSig.Type == ast.TYPE_CXARGUMENT_DEPRECATE {
				errStr += fmt.Sprintf("expr[%v] out[%v] = %+v\n\n", i, y, prgrm.GetCXArgFromArray(ast.CXArgumentIndex(typeSig.Meta)))
			} else {
				errStr += fmt.Sprintf("expr[%v] out[%v] = %+v\n", i, y, typeSig)

			}
		}
	}
	println(errStr)
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
		expressionInputTypeSig := prgrm.GetCXTypeSignatureFromArray(expression.GetInputs(prgrm)[0])

		var atomicType types.Code
		if expressionInputTypeSig.Type == ast.TYPE_CXARGUMENT_DEPRECATE {
			atomicType = prgrm.CXArgs[expressionInputTypeSig.Meta].GetType(prgrm)
		} else if expressionInputTypeSig.Type == ast.TYPE_ATOMIC {
			atomicType = types.Code(expressionInputTypeSig.Meta)
		} else if expressionInputTypeSig.Type == ast.TYPE_POINTER_ATOMIC {
			atomicType = types.Code(expressionInputTypeSig.Meta)
		} else if expressionInputTypeSig.Type == ast.TYPE_SLICE_ATOMIC {
			sliceDetails := prgrm.GetCXTypeSignatureArrayFromArray(expressionInputTypeSig.Meta)
			atomicType = types.Code(sliceDetails.Meta)
		} else {
			panic("type is not known")
		}

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
		expressionOutputTypeSig := prgrm.GetCXTypeSignatureFromArray(expression.GetOutputs(prgrm)[0])
		var expressionOutputArg *ast.CXArgument = &ast.CXArgument{}
		if expressionOutputTypeSig.Type == ast.TYPE_CXARGUMENT_DEPRECATE {
			expressionOutputArg = prgrm.GetCXArgFromArray(ast.CXArgumentIndex(expressionOutputTypeSig.Meta))
		} else {
			expressionOutputArg = &ast.CXArgument{}
		}

		opName := expressionOutputTypeSig.Name
		opPkgIdx := expressionOutputTypeSig.Package

		opPkg, err := prgrm.GetPackageFromArray(opPkgIdx)
		if err != nil {
			panic(err)
		}

		if op, err := prgrm.GetFunction(opName, opPkg.Name); err == nil {
			expression.Operator = ast.CXFunctionIndex(op.Index)
		} else if expressionOutputArg.Fields == nil {
			// then it's not a possible method call
			println(ast.CompilationError(CurrentFile, LineNo), err.Error())
			return nil
		} else {
			exprs[len(exprs)-1].ExpressionType = ast.CXEXPR_METHOD_CALL
		}

		if len(expression.GetOutputs(prgrm)) > 0 && expressionOutputArg.Fields == nil {
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
			typeSigIdx := inpExprAtomicOp.GetOutputs(prgrm)[0]

			// typeSig := prgrm.GetCXTypeSignatureFromArray(typeSigIdx)

			// // TODO: This is still to be improved
			// // Get its CX argument if its type cx arg deprecate
			// var typeSigArg *ast.CXArgument = &ast.CXArgument{}
			// if typeSig.Type == ast.TYPE_CXARGUMENT_DEPRECATE {
			// 	typeSigArg = prgrm.GetCXArgFromArray(ast.CXArgumentIndex(typeSig.Meta))

			// 	// If an atomic type, convert cx type signature to type atomic/
			// 	if ast.IsTypeAtomic(typeSigArg) {
			// 		typeSig.Name = typeSigArg.Name
			// 		typeSig.Type = ast.TYPE_ATOMIC
			// 		typeSig.Meta = int(typeSigArg.Type)
			// 		typeSig.Offset = typeSigArg.Offset
			// 		typeSig.Package = typeSigArg.Package
			// 	}
			// } else if typeSig.Type == ast.TYPE_ATOMIC || typeSig.Type == ast.TYPE_POINTER_ATOMIC {
			// 	// do nothing
			// }

			// then it's a literal
			expression.AddInput(prgrm, typeSigIdx)
		} else {
			// then it's a function call
			if len(inpExprAtomicOp.GetOutputs(prgrm)) < 1 {
				var out *ast.CXArgument = &ast.CXArgument{}

				inpExprAtomicOpOperatorOutputs := inpExprAtomicOpOperator.GetOutputs(prgrm)
				inpExprAtomicOpOperatorOutputTypeSig := prgrm.GetCXTypeSignatureFromArray(inpExprAtomicOpOperatorOutputs[0])
				var inpExprAtomicOpOperatorOutputArg *ast.CXArgument = &ast.CXArgument{}
				if inpExprAtomicOpOperatorOutputTypeSig.Type == ast.TYPE_CXARGUMENT_DEPRECATE {
					inpExprAtomicOpOperatorOutputArg = prgrm.GetCXArgFromArray(ast.CXArgumentIndex(inpExprAtomicOpOperatorOutputTypeSig.Meta))

					var typeSigIdx ast.CXTypeSignatureIndex

					// if undefined type, then adopt argument's type
					if inpExprAtomicOpOperatorOutputArg.Type == types.UNDEFINED {
						inpExprAtomicOpInputTypeSig := prgrm.GetCXTypeSignatureFromArray(inpExprAtomicOp.GetInputs(prgrm)[0])

						var inpExprAtomicOpInputArg *ast.CXArgument = &ast.CXArgument{}
						if inpExprAtomicOpInputTypeSig.Type == ast.TYPE_CXARGUMENT_DEPRECATE {
							inpExprAtomicOpInputArg = prgrm.GetCXArgFromArray(ast.CXArgumentIndex(inpExprAtomicOpInputTypeSig.Meta))

							out = ast.MakeArgument(generateTempVarName(constants.LOCAL_PREFIX), CurrentFile, inpExprCXLine.LineNumber).SetType(inpExprAtomicOpInputArg.Type)
							out.StructType = inpExprAtomicOpInputArg.StructType
							out.Size = inpExprAtomicOpInputArg.Size
							out.Type = inpExprAtomicOpInputArg.Type
							out.PointerTargetType = inpExprAtomicOpInputArg.PointerTargetType
							out.PreviouslyDeclared = true

							out.Package = inpExprAtomicOp.Package

							typeSig := ast.GetCXTypeSignatureRepresentationOfCXArg(prgrm, out)
							typeSigIdx = prgrm.AddCXTypeSignatureInArray(typeSig)
						} else if inpExprAtomicOpInputTypeSig.Type == ast.TYPE_ATOMIC {
							outTypeSig := &ast.CXTypeSignature{}
							outTypeSig.Name = generateTempVarName(constants.LOCAL_PREFIX)
							outTypeSig.Type = ast.TYPE_ATOMIC
							outTypeSig.Meta = inpExprAtomicOpInputTypeSig.Meta
							outTypeSig.Package = inpExprAtomicOpInputTypeSig.Package
							outTypeSig.Offset = inpExprAtomicOpInputTypeSig.Offset

							typeSigIdx = prgrm.AddCXTypeSignatureInArray(outTypeSig)
						} else if inpExprAtomicOpInputTypeSig.Type == ast.TYPE_POINTER_ATOMIC {
							outTypeSig := &ast.CXTypeSignature{}
							outTypeSig.Name = generateTempVarName(constants.LOCAL_PREFIX)
							outTypeSig.Type = ast.TYPE_POINTER_ATOMIC
							outTypeSig.Meta = inpExprAtomicOpInputTypeSig.Meta
							outTypeSig.Package = inpExprAtomicOpInputTypeSig.Package
							outTypeSig.Offset = inpExprAtomicOpInputTypeSig.Offset

							typeSigIdx = prgrm.AddCXTypeSignatureInArray(outTypeSig)
						} else {
							panic("type is not known")
						}

					} else {
						out = ast.MakeArgument(generateTempVarName(constants.LOCAL_PREFIX), CurrentFile, inpExprCXLine.LineNumber).SetType(inpExprAtomicOpOperatorOutputArg.Type)
						out.DeclarationSpecifiers = inpExprAtomicOpOperatorOutputArg.DeclarationSpecifiers

						out.StructType = inpExprAtomicOpOperatorOutputArg.StructType

						if inpExprAtomicOpOperatorOutputArg.StructType != nil {
							inpExprPkg, err := prgrm.GetPackageFromArray(inpExprAtomicOp.Package)
							if err != nil {
								panic(err)
							}
							if strct, err := inpExprPkg.GetStruct(prgrm, inpExprAtomicOpOperatorOutputArg.StructType.Name); err == nil {
								out.Size = strct.GetStructSize(prgrm)
							}
						} else {
							out.Size = inpExprAtomicOpOperatorOutputArg.Size
						}

						out.Type = inpExprAtomicOpOperatorOutputArg.Type
						out.PointerTargetType = inpExprAtomicOpOperatorOutputArg.PointerTargetType
						out.PreviouslyDeclared = true
						out.Package = inpExprAtomicOp.Package

						typeSig := ast.GetCXTypeSignatureRepresentationOfCXArg(prgrm, out)
						typeSigIdx = prgrm.AddCXTypeSignatureInArray(typeSig)
					}

					inpExprAtomicOp.AddOutput(prgrm, typeSigIdx)
					expression.AddInput(prgrm, typeSigIdx)
				} else if inpExprAtomicOpOperatorOutputTypeSig.Type == ast.TYPE_ATOMIC {
					newTypeSig := *prgrm.GetCXTypeSignatureFromArray(inpExprAtomicOpOperatorOutputs[0])
					newTypeSig.Name = generateTempVarName(constants.LOCAL_PREFIX)
					newTypeSig.Package = inpExprAtomicOp.Package
					newTypeSigIdx := prgrm.AddCXTypeSignatureInArray(&newTypeSig)

					inpExprAtomicOp.AddOutput(prgrm, newTypeSigIdx)
					expression.AddInput(prgrm, newTypeSigIdx)
				} else if inpExprAtomicOpOperatorOutputTypeSig.Type == ast.TYPE_POINTER_ATOMIC {
					newTypeSig := *prgrm.GetCXTypeSignatureFromArray(inpExprAtomicOpOperatorOutputs[0])
					newTypeSig.Name = generateTempVarName(constants.LOCAL_PREFIX)
					newTypeSig.Package = inpExprAtomicOp.Package
					newTypeSigIdx := prgrm.AddCXTypeSignatureInArray(&newTypeSig)

					inpExprAtomicOp.AddOutput(prgrm, newTypeSigIdx)
					expression.AddInput(prgrm, newTypeSigIdx)
				} else {
					panic("type is not known")
				}

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

	var typeCode types.Code
	expressionInputTypeSig := prgrm.GetCXTypeSignatureFromArray(expression.GetInputs(prgrm)[0])
	if expressionInputTypeSig.Type == ast.TYPE_CXARGUMENT_DEPRECATE {
		expressionInputArg := prgrm.GetCXArgFromArray(ast.CXArgumentIndex(expressionInputTypeSig.Meta))

		typeCode = expressionInputArg.Type
		if expressionInputArg.Type == types.POINTER {
			typeCode = expressionInputArg.PointerTargetType
		}
	} else if expressionInputTypeSig.Type == ast.TYPE_ATOMIC {
		typeCode = types.Code(expressionInputTypeSig.Meta)
	} else if expressionInputTypeSig.Type == ast.TYPE_POINTER_ATOMIC {
		typeCode = types.Code(expressionInputTypeSig.Meta)
	} else if expressionInputTypeSig.Type == ast.TYPE_SLICE_ATOMIC {
		sliceDetails := prgrm.GetCXTypeSignatureArrayFromArray(expressionInputTypeSig.Meta)
		typeCode = types.Code(sliceDetails.Meta)
	} else {
		panic("type is not known")
	}

	for _, inputIdx := range expression.GetInputs(prgrm) {
		input := prgrm.GetCXTypeSignatureFromArray(inputIdx)

		var inpType types.Code
		if input.Type == ast.TYPE_CXARGUMENT_DEPRECATE {
			inp := prgrm.GetCXArgFromArray(ast.CXArgumentIndex(input.Meta))

			inpType = inp.Type
			if inp.Type == types.POINTER {
				inpType = inp.PointerTargetType
			}
		} else if input.Type == ast.TYPE_ATOMIC {
			inpType = types.Code(input.Meta)
		} else if input.Type == ast.TYPE_POINTER_ATOMIC {
			inpType = types.Code(input.Meta)
		} else if input.Type == ast.TYPE_ARRAY_ATOMIC {
			arrDetails := prgrm.GetCXTypeSignatureArrayFromArray(input.Meta)
			inpType = types.Code(arrDetails.Meta)
		} else if input.Type == ast.TYPE_SLICE_ATOMIC {
			sliceDetails := prgrm.GetCXTypeSignatureArrayFromArray(input.Meta)
			inpType = types.Code(sliceDetails.Meta)
		} else {
			panic("type is not known")
		}

		if inpType != typeCode {
			return errors.New(fmt.Sprintf("operands are not of the same type: %v!=%v", inpType, typeCode))
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
		for _, outputIdx := range expression.GetOutputs(prgrm) {
			output := prgrm.GetCXTypeSignatureFromArray(outputIdx)

			var outIdx ast.CXArgumentIndex
			if output.Type == ast.TYPE_CXARGUMENT_DEPRECATE {
				outIdx = ast.CXArgumentIndex(output.Meta)

				out := prgrm.GetCXArgFromArray(outIdx)

				size := types.Pointer(1)
				if !ast.IsComparisonOperator(expressionOperator.AtomicOPCode) {
					expressionInputTypeSig := prgrm.GetCXTypeSignatureFromArray(expression.GetInputs(prgrm)[0])

					var expressionInputArg *ast.CXArgument = &ast.CXArgument{}
					if expressionInputTypeSig.Type == ast.TYPE_CXARGUMENT_DEPRECATE {
						expressionInputArg = prgrm.GetCXArgFromArray(ast.CXArgumentIndex(expressionInputTypeSig.Meta))

						size = ast.GetArgSize(prgrm, expressionInputArg.GetAssignmentElement(prgrm))
					} else if expressionInputTypeSig.Type == ast.TYPE_ATOMIC {
						size = expressionInputTypeSig.GetSize(prgrm, false)
					} else if expressionInputTypeSig.Type == ast.TYPE_POINTER_ATOMIC {
						size = expressionInputTypeSig.GetSize(prgrm, false)
					} else if expressionInputTypeSig.Type == ast.TYPE_SLICE_ATOMIC {
						size = expressionInputTypeSig.GetSize(prgrm, false)
					} else {
						panic("type is not known")
					}
				}

				out.Size = size
			} else if output.Type == ast.TYPE_ATOMIC || output.Type == ast.TYPE_POINTER_ATOMIC {
				continue
			} else {
				panic("type is not known")
			}
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

	for _, argTypeSigIdx := range append(expression.GetInputs(prgrm), expression.GetOutputs(prgrm)...) {
		argTypeSig := prgrm.GetCXTypeSignatureFromArray(argTypeSigIdx)

		var arg *ast.CXArgument = &ast.CXArgument{}
		if argTypeSig.Type == ast.TYPE_CXARGUMENT_DEPRECATE {
			arg = prgrm.GetCXArgFromArray(ast.CXArgumentIndex(argTypeSig.Meta))
		} else if argTypeSig.Type == ast.TYPE_ATOMIC || argTypeSig.Type == ast.TYPE_POINTER_ATOMIC {
			continue
		} else if argTypeSig.Type == ast.TYPE_ARRAY_ATOMIC {
			continue
		} else if argTypeSig.Type == ast.TYPE_POINTER_ARRAY_ATOMIC {
			continue
		} else if argTypeSig.Type == ast.TYPE_SLICE_ATOMIC {
			continue
		} else if argTypeSig.Type == ast.TYPE_POINTER_SLICE_ATOMIC {
			continue
		} else {
			panic("type is not known")
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
		if arg.IsStruct() && arg.IsPointer() && len(arg.Fields) > 0 && !doesArgHaveDerefPointer {
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
			var inp1Type, inp2Type string

			expressionInputFirstTypeSig := prgrm.GetCXTypeSignatureFromArray(expression.GetInputs(prgrm)[0])
			inp1Type = ast.GetFormattedType(prgrm, expressionInputFirstTypeSig)

			expressionInputSecondTypeSig := prgrm.GetCXTypeSignatureFromArray(expression.GetInputs(prgrm)[1])
			inp2Type = ast.GetFormattedType(prgrm, expressionInputSecondTypeSig)

			if inp1Type != inp2Type {
				println(ast.CompilationError(CurrentFile, LineNo), fmt.Sprintf("first and second input arguments' types are not equal in '%s' call ('%s' != '%s')", ast.OpNames[expressionOperator.AtomicOPCode], inp1Type, inp2Type))
			}
		}
	}
}

// checkIndexType throws an error if the type of `idx` is not `i32` or `i64`.
func checkIndexType(prgrm *ast.CXProgram, idxIdx ast.CXTypeSignatureIndex) {
	idxTypeSig := prgrm.GetCXTypeSignatureFromArray(idxIdx)
	var idxType string
	var idx *ast.CXArgument = &ast.CXArgument{ArgDetails: &ast.CXArgumentDebug{}}
	if idxTypeSig.Type == ast.TYPE_CXARGUMENT_DEPRECATE {
		idx = prgrm.GetCXArg(ast.CXArgumentIndex(idxTypeSig.Meta))
	}

	idxType = ast.GetFormattedType(prgrm, idxTypeSig)

	if idxType != "i32" && idxType != "i64" {
		println(ast.CompilationError(idx.ArgDetails.FileName, idx.ArgDetails.FileLine), fmt.Sprintf("wrong index type; expected either 'i32' or 'i64', got '%s'", idxType))
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
//  offset - offset to use by statements (excluding inputs, outputs
// 			 and receiver).
//  fnIdx - the index of the function in the main CXFunction array.
//  args - the expression arguments.
//  expr - the expression.
//  isInput - true if args are input arguments, false if they are output args.
func ProcessExpressionArguments(prgrm *ast.CXProgram, symbolsData *SymbolsData, offset *types.Pointer, fnIdx ast.CXFunctionIndex, args []ast.CXTypeSignatureIndex, expr *ast.CXExpression, isInput bool) {
	fn := prgrm.GetFunctionFromArray(fnIdx)

	for _, typeSignatureIdx := range args {
		typeSignature := prgrm.GetCXTypeSignatureFromArray(typeSignatureIdx)
		var argIdx ast.CXArgumentIndex
		if typeSignature.Type == ast.TYPE_CXARGUMENT_DEPRECATE {
			argIdx = ast.CXArgumentIndex(typeSignature.Meta)
		} else {
			argIdx = -1
		}

		arg := prgrm.GetCXArgFromArray(argIdx)

		if !isInput {
			CheckRedeclared(prgrm, symbolsData, expr, typeSignatureIdx)
		}

		if !isInput {
			ProcessOperatorExpression(prgrm, expr)
		}

		// TODO: Check how to remove PreviouslyDeclared field
		// maybe its only used by temp variables
		// TempVar always have PreviouslyDeclared=true
		typeSignatureName := typeSignature.Name
		isLocalVar := fn.IsLocalVariable(typeSignatureName)
		if (arg != nil && arg.PreviouslyDeclared) || IsTempVar(typeSignatureName) || (arg == nil && isLocalVar) {
			UpdateSymbolsTable(prgrm, symbolsData, typeSignatureIdx, offset, false)
		} else {
			UpdateSymbolsTable(prgrm, symbolsData, typeSignatureIdx, offset, true)
		}

		GiveOffset(prgrm, symbolsData, typeSignatureIdx)
		ProcessSlice(prgrm, argIdx)

		if arg != nil {
			for _, idxIdx := range arg.Indexes {
				UpdateSymbolsTable(prgrm, symbolsData, idxIdx, offset, true)
				GiveOffset(prgrm, symbolsData, idxIdx)
				checkIndexType(prgrm, idxIdx)
			}
			for _, fldIdx := range arg.Fields {
				fld := prgrm.GetCXArgFromArray(fldIdx)
				for _, idxIdx := range fld.Indexes {
					UpdateSymbolsTable(prgrm, symbolsData, idxIdx, offset, true)
					GiveOffset(prgrm, symbolsData, idxIdx)
				}
			}
		}
		AddPointer(prgrm, fn, typeSignatureIdx)
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
//  typeSigIdx - the index of the type signature from the main CXTypeSignature array.
func AddPointer(prgrm *ast.CXProgram, fn *ast.CXFunction, typeSigIdx ast.CXTypeSignatureIndex) {
	typeSig := prgrm.GetCXTypeSignatureFromArray(typeSigIdx)
	var sym *ast.CXArgument = &ast.CXArgument{}
	if typeSig.Type == ast.TYPE_CXARGUMENT_DEPRECATE {
		sym = prgrm.GetCXArgFromArray(ast.CXArgumentIndex(typeSig.Meta))
	} else {
		// panic("type is not cx arg deprecate\n\n)")
		// TODO: atomic type for now doesnt have pointers
		return
	}

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
//  typeSigIdx - the index of the type signature from the main CXTypeSignature array.
func CheckRedeclared(prgrm *ast.CXProgram, symbolsData *SymbolsData, expr *ast.CXExpression, typeSigIdx ast.CXTypeSignatureIndex) {
	typeSig := prgrm.GetCXTypeSignatureFromArray(typeSigIdx)
	var arg *ast.CXArgument
	if typeSig.Type == ast.TYPE_CXARGUMENT_DEPRECATE {
		arg = prgrm.GetCXArgFromArray(ast.CXArgumentIndex(typeSig.Meta))
	} else {
		// panic("type is not cx arg deprecate\n\n")
		// TODO: temporary put empty arg
		arg = &ast.CXArgument{ArgDetails: &ast.CXArgumentDebug{}}
	}

	expression, err := prgrm.GetCXAtomicOp(expr.Index)
	if err != nil {
		panic(err)
	}
	expressionOperator := prgrm.GetFunctionFromArray(expression.Operator)

	if expressionOperator == nil && len(expression.GetOutputs(prgrm)) > 0 && len(expression.GetInputs(prgrm)) == 0 {
		lastIdx := len(symbolsData.symbolsIndex) - 1

		symPkg, err := prgrm.GetPackageFromArray(typeSig.Package)
		if err != nil {
			panic(err)
		}

		_, found := (symbolsData.symbolsIndex)[lastIdx][symPkg.Name+"."+typeSig.Name]
		if found {
			println(ast.CompilationError(arg.ArgDetails.FileName, arg.ArgDetails.FileLine), fmt.Sprintf("'%s' redeclared", typeSig.Name))
		}
	}
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

func checkMatchParamTypes(prgrm *ast.CXProgram, expr *ast.CXExpression, expected, received []ast.CXTypeSignatureIndex, isInputs bool) {
	expression, err := prgrm.GetCXAtomicOp(expr.Index)
	if err != nil {
		panic(err)
	}
	expressionOperator := prgrm.GetFunctionFromArray(expression.Operator)

	for i, expectedTypeSigIdx := range expected {
		var receivedType, expectedType string

		expectedTypeSig := prgrm.GetCXTypeSignatureFromArray(expectedTypeSigIdx)
		var expectedArg *ast.CXArgument = &ast.CXArgument{ArgDetails: &ast.CXArgumentDebug{}}
		if expectedTypeSig.Type == ast.TYPE_CXARGUMENT_DEPRECATE {
			expectedArg = prgrm.GetCXArgFromArray(ast.CXArgumentIndex(expectedTypeSig.Meta))
		}
		expectedType = ast.GetFormattedType(prgrm, expectedTypeSig)

		if expr.IsMethodCall() && expectedArg.IsPointer() && i == 0 {
			// if method receiver is pointer, remove *
			if expectedType[0] == '*' {
				// we need to check if it's not an `str`
				// otherwise we end up removing the `s` instead of a `*`
				expectedType = expectedType[1:]
			}
		}

		var receivedArg *ast.CXArgument = &ast.CXArgument{ArgDetails: &ast.CXArgumentDebug{}}
		receivedTypeSig := prgrm.GetCXTypeSignatureFromArray(received[i])
		if receivedTypeSig.Type == ast.TYPE_CXARGUMENT_DEPRECATE {
			receivedArg = prgrm.GetCXArg(ast.CXArgumentIndex(receivedTypeSig.Meta))
		}
		receivedType = ast.GetFormattedType(prgrm, receivedTypeSig)

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
				expressionOutputTypeSig := prgrm.GetCXTypeSignatureFromArray(expression.GetOutputs(prgrm)[i])
				var expressionOutputArg *ast.CXArgument = &ast.CXArgument{ArgDetails: &ast.CXArgumentDebug{}}
				if expressionOutputTypeSig.Type == ast.TYPE_CXARGUMENT_DEPRECATE {
					expressionOutputArg = prgrm.GetCXArgFromArray(ast.CXArgumentIndex(expressionOutputTypeSig.Meta))
				} else if expressionOutputTypeSig.Type == ast.TYPE_ATOMIC {
					// do nothing
				} else if expressionOutputTypeSig.Type == ast.TYPE_POINTER_ATOMIC {
					// do nothing
				} else if expressionOutputTypeSig.Type == ast.TYPE_SLICE_ATOMIC {
					// do nothing
				} else {
					panic("type is not known")
				}

				println(ast.CompilationError(expressionOutputArg.ArgDetails.FileName, expressionOutputArg.ArgDetails.FileLine), fmt.Sprintf("function '%s' expected receiving variable of type '%s'; '%s' was provided", opName, expectedType, receivedType))

			}

		}

		var expressionOutputTypeSig *ast.CXTypeSignature
		if expressionOperator.AtomicOPCode == constants.OP_IDENTITY {
			expressionOutputTypeSig = prgrm.GetCXTypeSignatureFromArray(expression.GetOutputs(prgrm)[0])
		}
		// In the case of assignment we need to check that the input's type matches the output's type.
		// FIXME: There are some expressions added by the cxgo where temporary variables are used.
		// These temporary variables' types are not properly being set. That's why we use !cxcore.IsTempVar to
		// exclude these cases for now.
		if expressionOperator.AtomicOPCode == constants.OP_IDENTITY && !IsTempVar(expressionOutputTypeSig.Name) {
			var outType, inpType string

			expressionInputTypeSig := prgrm.GetCXTypeSignatureFromArray(expression.GetInputs(prgrm)[0])
			inpType = ast.GetFormattedType(prgrm, expressionInputTypeSig)

			var expressionOutputArg *ast.CXArgument = &ast.CXArgument{}
			expressionOutputTypeSigName := expressionOutputTypeSig.Name
			if expressionOutputTypeSig.Type == ast.TYPE_CXARGUMENT_DEPRECATE {
				expressionOutputArg = prgrm.GetCXArgFromArray(ast.CXArgumentIndex(expressionOutputTypeSig.Meta))
				expressionOutputTypeSigName = expressionOutputArg.GetAssignmentElement(prgrm).Name
			}

			outType = ast.GetFormattedType(prgrm, expressionOutputTypeSig)

			// We use `isInputs` to only print the error once.
			// Otherwise we'd print the error twice: once for the input and again for the output
			if inpType != outType && isInputs {
				// println(ast.CompilationError(receivedArg.ArgDetails.FileName, receivedArg.ArgDetails.FileLine), fmt.Sprintf("cannot assign value of type '%s' to identifier '%s' of type '%s'", inpType, expressionOutputArg.GetAssignmentElement(prgrm).Name, outType))
				println(ast.CompilationError("", 0), fmt.Sprintf("cannot assign value of type '%s' to identifier '%s' of type '%s'", inpType, expressionOutputTypeSigName, outType))
			}

		}
	}
}

func checkMatchExprNumOfInputs(prgrm *ast.CXProgram, expression *ast.CXAtomicOperator, expressionOperatorInputs []ast.CXTypeSignatureIndex, opName string, exprCXLine *ast.CXLine) {
	// checking if number of inputs is not the same as the required number of inputs
	if len(expression.GetInputs(prgrm)) != len(expressionOperatorInputs) {
		expressionOperatorInputArg := &ast.CXArgument{}
		if len(expressionOperatorInputs) > 0 {
			expressionTypeSigIdx := expressionOperatorInputs[len(expressionOperatorInputs)-1]
			expressionTypeSig := prgrm.GetCXTypeSignatureFromArray(expressionTypeSigIdx)
			if expressionTypeSig.Type == ast.TYPE_CXARGUMENT_DEPRECATE {
				expressionOperatorInputArg = prgrm.GetCXArgFromArray(ast.CXArgumentIndex(expressionTypeSig.Meta))
			}
		}

		if !(len(expressionOperatorInputs) > 0 && expressionOperatorInputArg.Type != types.UNDEFINED) {
			// if the last input is of type cxcore.TYPE_UNDEFINED then it might be a variadic function, such as printf
		} else {
			// then we need to be strict in the number of inputs
			println(ast.CompilationError(exprCXLine.FileName, exprCXLine.LineNumber), fmt.Sprintf("operator '%s' expects %d input/s, but %d input argument\n were provided", opName, len(expressionOperatorInputs), len(expression.GetInputs(prgrm))))
			os.Exit(constants.CX_COMPILATION_ERROR)
		}
	}
}

func checkMatchExprNumOfOutputs(prgrm *ast.CXProgram, expression *ast.CXAtomicOperator, expressionOperatorOutputs []ast.CXTypeSignatureIndex, opName string, exprCXLine *ast.CXLine) {
	// checking if number of outputs is not the same as the required number of outputs
	if len(expression.GetOutputs(prgrm)) != len(expressionOperatorOutputs) {
		println(ast.CompilationError(exprCXLine.FileName, exprCXLine.LineNumber), fmt.Sprintf("operator '%s' expects to return %d output/s, but %d receiving argument/s were provided", opName, len(expressionOperatorOutputs), len(expression.GetOutputs(prgrm))))
		os.Exit(constants.CX_COMPILATION_ERROR)
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
		expressionOperatorInputs := expressionOperator.GetInputs(prgrm)
		expressionOperatorOutputs := expressionOperator.GetOutputs(prgrm)

		// checking if number of inputs is not the same as the required number of inputs
		checkMatchExprNumOfInputs(prgrm, expression, expressionOperatorInputs, opName, exprCXLine)
		// checking if number of outputs is not the same as the required number of outputs
		checkMatchExprNumOfOutputs(prgrm, expression, expressionOperatorOutputs, opName, exprCXLine)
	}

	if expressionOperator != nil && expressionOperator.IsBuiltIn() && expressionOperator.AtomicOPCode == constants.OP_IDENTITY {
		for i := range expression.GetInputs(prgrm) {
			var expectedType string
			var receivedType string

			expressionInputTypeSig := prgrm.GetCXTypeSignatureFromArray(expression.GetInputs(prgrm)[i])
			if expressionInputTypeSig.Type == ast.TYPE_CXARGUMENT_DEPRECATE {
				expressionInputArg := prgrm.GetCXArgFromArray(ast.CXArgumentIndex(expressionInputTypeSig.Meta))

				if expressionInputArg.GetAssignmentElement(prgrm).StructType != nil {
					// then it's custom type
					receivedType = expressionInputArg.GetAssignmentElement(prgrm).StructType.Name
				} else {
					// then it's native type
					receivedType = expressionInputArg.GetAssignmentElement(prgrm).Type.Name()

					if expressionInputArg.GetAssignmentElement(prgrm).Type == types.POINTER {
						receivedType = expressionInputArg.GetAssignmentElement(prgrm).PointerTargetType.Name()
					}
				}

			} else if expressionInputTypeSig.Type == ast.TYPE_ATOMIC {
				receivedType = types.Code(expressionInputTypeSig.Meta).Name()
			} else if expressionInputTypeSig.Type == ast.TYPE_POINTER_ATOMIC {
				receivedType = types.Code(expressionInputTypeSig.Meta).Name()
			} else if expressionInputTypeSig.Type == ast.TYPE_ARRAY_ATOMIC {
				arrDetails := prgrm.GetCXTypeSignatureArrayFromArray(expressionInputTypeSig.Meta)
				receivedType = types.Code(arrDetails.Meta).Name()
			} else if expressionInputTypeSig.Type == ast.TYPE_SLICE_ATOMIC {
				sliceDetails := prgrm.GetCXTypeSignatureArrayFromArray(expressionInputTypeSig.Meta)
				receivedType = types.Code(sliceDetails.Meta).Name()
			} else {
				panic("type is not known")
			}

			expressionOutputTypeSig := prgrm.GetCXTypeSignatureFromArray(expression.GetOutputs(prgrm)[i])
			var expressionOutputArg *ast.CXArgument = &ast.CXArgument{ArgDetails: &ast.CXArgumentDebug{}}
			if expressionOutputTypeSig.Type == ast.TYPE_CXARGUMENT_DEPRECATE {
				expressionOutputArg = prgrm.GetCXArgFromArray(ast.CXArgumentIndex(expressionOutputTypeSig.Meta))

				if expressionOutputArg.GetAssignmentElement(prgrm).StructType != nil {
					// then it's custom type
					expectedType = expressionOutputArg.GetAssignmentElement(prgrm).StructType.Name
				} else {
					// then it's native type
					expectedType = expressionOutputArg.GetAssignmentElement(prgrm).Type.Name()

					if expressionOutputArg.GetAssignmentElement(prgrm).Type == types.POINTER {
						expectedType = expressionOutputArg.GetAssignmentElement(prgrm).PointerTargetType.Name()
					}
				}
			} else if expressionOutputTypeSig.Type == ast.TYPE_ATOMIC {
				expectedType = types.Code(expressionOutputTypeSig.Meta).Name()
			} else if expressionOutputTypeSig.Type == ast.TYPE_POINTER_ATOMIC {
				expectedType = types.Code(expressionOutputTypeSig.Meta).Name()
			} else if expressionOutputTypeSig.Type == ast.TYPE_ARRAY_ATOMIC {
				arrDetails := prgrm.GetCXTypeSignatureArrayFromArray(expressionOutputTypeSig.Meta)
				expectedType = types.Code(arrDetails.Meta).Name()
			} else if expressionOutputTypeSig.Type == ast.TYPE_POINTER_ARRAY_ATOMIC {
				arrDetails := prgrm.GetCXTypeSignatureArrayFromArray(expressionOutputTypeSig.Meta)
				expectedType = types.Code(arrDetails.Meta).Name()
			} else if expressionOutputTypeSig.Type == ast.TYPE_SLICE_ATOMIC {
				arrDetails := prgrm.GetCXTypeSignatureArrayFromArray(expressionOutputTypeSig.Meta)
				expectedType = types.Code(arrDetails.Meta).Name()
			} else {
				panic("type is not known")
			}

			if receivedType != expectedType {
				if exprs[currIndex].IsStructLiteral() {
					println(ast.CompilationError(expressionOutputArg.ArgDetails.FileName, expressionOutputArg.ArgDetails.FileLine), fmt.Sprintf("field '%s' in struct literal of type '%s' expected argument of type '%s'; '%s' was provided", prgrm.GetCXArgFromArray(expressionOutputArg.Fields[0]).Name, expressionOutputArg.StructType.Name, expectedType, receivedType))
				} else {
					println(ast.CompilationError(expressionOutputArg.ArgDetails.FileName, expressionOutputArg.ArgDetails.FileLine), fmt.Sprintf("trying to assign argument of type '%s' to symbol '%s' of type '%s'", receivedType, expressionOutputTypeSig.Name, expectedType))
				}
			}

		}
	}

	// then it's a function call and not a declaration
	if expressionOperator != nil {
		// checking inputs matching operator's inputs
		checkMatchParamTypes(prgrm, &exprs[currIndex], expressionOperator.GetInputs(prgrm), expression.GetInputs(prgrm), true)

		// checking outputs matching operator's outputs
		checkMatchParamTypes(prgrm, &exprs[currIndex], expressionOperator.GetOutputs(prgrm), expression.GetOutputs(prgrm), false)
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
		for i, outputIdx := range expression.GetOutputs(prgrm) {
			output := prgrm.GetCXTypeSignatureFromArray(outputIdx)

			var out *ast.CXArgument = &ast.CXArgument{}
			if output.Type == ast.TYPE_CXARGUMENT_DEPRECATE {
				out = prgrm.GetCXArgFromArray(ast.CXArgumentIndex(output.Meta))
			} else if output.Type == ast.TYPE_ATOMIC {
				continue
			} else if output.Type == ast.TYPE_POINTER_ATOMIC {
				continue
			} else if output.Type == ast.TYPE_ARRAY_ATOMIC {
				continue
			} else if output.Type == ast.TYPE_POINTER_ARRAY_ATOMIC {
				continue
			} else if output.Type == ast.TYPE_SLICE_ATOMIC {
				continue
			} else {
				panic("type is not known")
			}

			if len(expression.GetInputs(prgrm)) > i {
				expressionInputTypeSig := prgrm.GetCXTypeSignatureFromArray(expression.GetInputs(prgrm)[i])
				var expressionInputArg *ast.CXArgument = &ast.CXArgument{}
				if expressionInputTypeSig.Type == ast.TYPE_CXARGUMENT_DEPRECATE {
					expressionInputArg = prgrm.GetCXArgFromArray(ast.CXArgumentIndex(expressionInputTypeSig.Meta))
				} else if expressionInputTypeSig.Type == ast.TYPE_ATOMIC {
					continue
				} else if expressionInputTypeSig.Type == ast.TYPE_POINTER_ATOMIC {
					continue
				} else if expressionInputTypeSig.Type == ast.TYPE_ARRAY_ATOMIC {
					continue
				} else if expressionInputTypeSig.Type == ast.TYPE_POINTER_ARRAY_ATOMIC {
					continue
				} else if expressionInputTypeSig.Type == ast.TYPE_SLICE_ATOMIC {
					continue
				} else {
					panic("type is not known")
				}

				out = out.GetAssignmentElement(prgrm)
				inp := expressionInputArg.GetAssignmentElement(prgrm)

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

	var expressionOutputArgType types.Code
	if len(expression.GetOutputs(prgrm)) > 0 {
		expressionOutputTypeSig := prgrm.GetCXTypeSignatureFromArray(expression.GetOutputs(prgrm)[0])
		if expressionOutputTypeSig.Type == ast.TYPE_CXARGUMENT_DEPRECATE {
			expressionOutputArg := prgrm.GetCXArgFromArray(ast.CXArgumentIndex(expressionOutputTypeSig.Meta))
			expressionOutputArgType = expressionOutputArg.Type
		} else {
			expressionOutputArgType = types.Code(expressionOutputTypeSig.Meta)
		}
	}

	// process short declaration
	if len(expression.GetOutputs(prgrm)) > 0 && len(expression.GetInputs(prgrm)) > 0 && expressionOutputArgType == types.IDENTIFIER && !expr.IsStructLiteral() && !isParseOp(prgrm, expr) {
		prevExpression, err := prgrm.GetPreviousCXAtomicOpFromExpressions(expressions, idx-1)
		if err != nil {
			panic(err)
		}

		var argType, argPointerTargetType types.Code
		var argSize types.Pointer
		if expr.IsMethodCall() {

			expressionOperator := prgrm.GetFunctionFromArray(expression.Operator)
			expressionOperatorOutputs := expressionOperator.GetOutputs(prgrm)
			expressionOperatorOutputTypeSig := prgrm.GetCXTypeSignatureFromArray(expressionOperatorOutputs[0])

			if expressionOperatorOutputTypeSig.Type == ast.TYPE_CXARGUMENT_DEPRECATE {
				expressionOperatorOutputArg := prgrm.GetCXArgFromArray(ast.CXArgumentIndex(expressionOperatorOutputTypeSig.Meta))
				argType = expressionOperatorOutputArg.Type
				argPointerTargetType = expressionOperatorOutputArg.PointerTargetType
				argSize = expressionOperatorOutputArg.Size
			} else if expressionOperatorOutputTypeSig.Type == ast.TYPE_ATOMIC {
				argType = types.Code(expressionOperatorOutputTypeSig.Meta)
				argSize = types.Code(expressionOperatorOutputTypeSig.Meta).Size()
			} else if expressionOperatorOutputTypeSig.Type == ast.TYPE_POINTER_ATOMIC {
				argType = types.Code(expressionOperatorOutputTypeSig.Meta)
				argSize = types.Code(expressionOperatorOutputTypeSig.Meta).Size()
			} else {
				panic("type is not known")
			}

		} else {
			expressionInputTypeSig := prgrm.GetCXTypeSignatureFromArray(expression.GetInputs(prgrm)[0])
			if expressionInputTypeSig.Type == ast.TYPE_CXARGUMENT_DEPRECATE {
				expressionInputArg := prgrm.GetCXArgFromArray(ast.CXArgumentIndex(expressionInputTypeSig.Meta))
				argType = expressionInputArg.Type
				argPointerTargetType = expressionInputArg.PointerTargetType
				argSize = expressionInputArg.Size
			} else if expressionInputTypeSig.Type == ast.TYPE_ATOMIC {
				argType = types.Code(expressionInputTypeSig.Meta)
				argSize = types.Code(expressionInputTypeSig.Meta).Size()
			} else if expressionInputTypeSig.Type == ast.TYPE_POINTER_ATOMIC {
				argType = types.Code(expressionInputTypeSig.Meta)
				argSize = types.Code(expressionInputTypeSig.Meta).Size()
			} else if expressionInputTypeSig.Type == ast.TYPE_ARRAY_ATOMIC {
				arrDetails := prgrm.GetCXTypeSignatureArrayFromArray(expressionInputTypeSig.Meta)
				argType = types.Code(arrDetails.Meta)
				argSize = types.Code(arrDetails.Meta).Size()
			} else if expressionInputTypeSig.Type == ast.TYPE_SLICE_ATOMIC {
				sliceDetails := prgrm.GetCXTypeSignatureArrayFromArray(expressionInputTypeSig.Meta)
				argType = types.Code(sliceDetails.Meta)
				argSize = types.Code(sliceDetails.Meta).Size()
			} else {
				panic("type is not known")
			}
		}

		prevExpressionOutputTypeSig := prgrm.GetCXTypeSignatureFromArray(prevExpression.GetOutputs(prgrm)[0])
		var prevExpressionOutputIdx ast.CXArgumentIndex
		if prevExpressionOutputTypeSig.Type == ast.TYPE_CXARGUMENT_DEPRECATE {
			prevExpressionOutputIdx = ast.CXArgumentIndex(prevExpressionOutputTypeSig.Meta)

			prgrm.CXArgs[prevExpressionOutputIdx].Type = argType
			prgrm.CXArgs[prevExpressionOutputIdx].PointerTargetType = argPointerTargetType
			prgrm.CXArgs[prevExpressionOutputIdx].Size = argSize
		} else if prevExpressionOutputTypeSig.Type == ast.TYPE_ATOMIC {
			prevExpressionOutputTypeSig.Meta = int(argType)
		} else if prevExpressionOutputTypeSig.Type == ast.TYPE_POINTER_ATOMIC {
			prevExpressionOutputTypeSig.Meta = int(argType)
		} else if prevExpressionOutputTypeSig.Type == ast.TYPE_ARRAY_ATOMIC {
			arrDetails := prgrm.GetCXTypeSignatureArrayFromArray(prevExpressionOutputTypeSig.Meta)
			arrDetails.Meta = int(argType)

		} else if prevExpressionOutputTypeSig.Type == ast.TYPE_SLICE_ATOMIC {
			sliceDetails := prgrm.GetCXTypeSignatureArrayFromArray(prevExpressionOutputTypeSig.Meta)
			sliceDetails.Meta = int(argType)

		} else {
			panic("type is not known")
		}

		expressionOutputTypeSig := prgrm.GetCXTypeSignatureFromArray(expression.GetOutputs(prgrm)[0])
		if expressionOutputTypeSig.Type == ast.TYPE_CXARGUMENT_DEPRECATE {
			expressionOutputIdx := ast.CXArgumentIndex(expressionOutputTypeSig.Meta)

			prgrm.CXArgs[expressionOutputIdx].Type = argType
			prgrm.CXArgs[expressionOutputIdx].PointerTargetType = argPointerTargetType
			prgrm.CXArgs[expressionOutputIdx].Size = argSize
		} else if expressionOutputTypeSig.Type == ast.TYPE_ATOMIC {
			expressionOutputTypeSig.Meta = int(argType)
		} else if expressionOutputTypeSig.Type == ast.TYPE_POINTER_ATOMIC {
			expressionOutputTypeSig.Meta = int(argType)
		} else if expressionOutputTypeSig.Type == ast.TYPE_ARRAY_ATOMIC {
			arrDetails := prgrm.GetCXTypeSignatureArrayFromArray(expressionOutputTypeSig.Meta)
			arrDetails.Meta = int(argType)
		} else if expressionOutputTypeSig.Type == ast.TYPE_SLICE_ATOMIC {
			sliceDetails := prgrm.GetCXTypeSignatureArrayFromArray(expressionOutputTypeSig.Meta)
			sliceDetails.Meta = int(argType)
		} else {
			panic("type is not known")
		}
	}
}

// ProcessSlice sets DereferenceOperations if the arg is a slice.
func ProcessSlice(prgrm *ast.CXProgram, inpIdx ast.CXArgumentIndex) {
	if inpIdx == -1 {
		return
	}

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

		expressionInputTypeSig := prgrm.GetCXTypeSignatureFromArray(expression.GetInputs(prgrm)[0])
		var expressionInputArg *ast.CXArgument = &ast.CXArgument{}
		if expressionInputTypeSig.Type == ast.TYPE_CXARGUMENT_DEPRECATE {
			expressionInputArg = prgrm.GetCXArgFromArray(ast.CXArgumentIndex(expressionInputTypeSig.Meta))
		} else {
			// panic("type is not cx argument deprecate\n\n")
			return
		}

		expressionOutputTypeSig := prgrm.GetCXTypeSignatureFromArray(expression.GetOutputs(prgrm)[0])
		var expressionOutputArg *ast.CXArgument = &ast.CXArgument{}
		if expressionOutputTypeSig.Type == ast.TYPE_CXARGUMENT_DEPRECATE {
			expressionOutputArg = prgrm.GetCXArgFromArray(ast.CXArgumentIndex(expressionOutputTypeSig.Meta))
		} else {
			// panic("type is not cx argument deprecate\n\n")
			return
		}

		inp = expressionInputArg.GetAssignmentElement(prgrm)
		out = expressionOutputArg.GetAssignmentElement(prgrm)

		if inp.IsSlice && out.IsSlice && len(inp.Indexes) == 0 && len(out.Indexes) == 0 {
			out.PassBy = constants.PASSBY_VALUE
		}
	}
	if expressionOperator != nil && !expressionOperator.IsBuiltIn() {
		// then it's a function call
		for _, inputIdx := range expression.GetInputs(prgrm) {
			input := prgrm.GetCXTypeSignatureFromArray(inputIdx)

			var inp *ast.CXArgument = &ast.CXArgument{}
			if input.Type == ast.TYPE_CXARGUMENT_DEPRECATE {
				inp = prgrm.GetCXArgFromArray(ast.CXArgumentIndex(input.Meta))
			} else if input.Type == ast.TYPE_ATOMIC {
				// Continue since atomic types doesnt include slices.
				continue
			} else if input.Type == ast.TYPE_POINTER_ATOMIC {
				// Continue since pointer atomic types doesnt include slices.
				continue
			} else if input.Type == ast.TYPE_ARRAY_ATOMIC {
				continue
			} else if input.Type == ast.TYPE_SLICE_ATOMIC {
				continue
			} else {
				panic("type is not known")
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
func lookupSymbol(prgrm *ast.CXProgram, pkgName, ident string, symbolsData *SymbolsData) (*ast.CXTypeSignature, error) {
	fullName := pkgName + "." + StripNameNumber(ident)
	for c := len(symbolsData.symbolsIndex) - 1; c >= 0; c-- {
		if symIdx, found := (symbolsData.symbolsIndex)[c][fullName]; found {
			return symbolsData.symbols[symIdx], nil
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
func UpdateSymbolsTable(prgrm *ast.CXProgram, symbolsData *SymbolsData, typeSigIdx ast.CXTypeSignatureIndex, offset *types.Pointer, shouldExist bool) {
	typeSig := prgrm.GetCXTypeSignatureFromArray(typeSigIdx)

	var sym *ast.CXArgument = &ast.CXArgument{}
	if typeSig.Type == ast.TYPE_CXARGUMENT_DEPRECATE {
		sym = prgrm.GetCXArgFromArray(ast.CXArgumentIndex(typeSig.Meta))
		if sym.Type == types.UNDEFINED && !IsTempVar(typeSig.Name) {
			return
		}
	} else {
		// panic("type is not cx arg deprecate\n\n")
		// TODO: temporary return only
		sym = nil
	}

	if typeSig.Name != "" {
		typeSigPkg, err := prgrm.GetPackageFromArray(typeSig.Package)
		if err != nil {
			panic(err)
		}

		currFn, err := prgrm.GetCurrentFunction()
		if err != nil {
			panic("error getting current function")
		}

		if !currFn.IsLocalVariable(typeSig.Name) {
			GetGlobalSymbol(prgrm, symbolsData, typeSigPkg, typeSig.Name)
		}

		lastIdx := len(symbolsData.symbolsIndex) - 1
		fullName := typeSigPkg.Name + "." + StripNameNumber(typeSig.Name)

		// outerSym, err := lookupSymbol(sym.Package.Name, sym.Name, symbols)
		_, err = lookupSymbol(prgrm, typeSigPkg.Name, typeSig.Name, symbolsData)
		_, found := (symbolsData.symbolsIndex)[lastIdx][fullName]

		// then it wasn't found in any scope
		if err != nil && shouldExist {
			// println(ast.CompilationError(sym.ArgDetails.FileName, sym.ArgDetails.FileLine), "identifier '"+typeSig.Name+"' does not exist")
			println(ast.CompilationError("", 0), "identifier '"+typeSig.Name+"' does not exist")
		}

		// then it was already added in the innermost scope
		if found {
			return
		}

		// then it is a new declaration
		if !shouldExist && !found {
			// We remove it from local var array since it is already added to the symbols
			err := currFn.RemoveLocalVariableFromArray(typeSig.Name)
			if err != nil {
				panic(err)
			}

			if !IsTempVar(typeSig.Name) {
				// Change name to n:varname format
				typeSig.Name = fmt.Sprintf("%v:%s", symbolsData.varCount, typeSig.Name)
				symbolsData.varCount++
			}

			symbolsData.symbols = append(symbolsData.symbols, typeSig)
			// add name to symbols
			symbolsData.symbolsIndex[lastIdx][fullName] = len(symbolsData.symbols) - 1

			if sym != nil {
				sym.Name = typeSig.Name
				sym.Offset = *offset
				*offset += ast.GetArgSize(prgrm, sym)
			} else {
				typeSig.Offset = *offset
				*offset += typeSig.GetSize(prgrm, true)
			}
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
func ProcessMethodCall(prgrm *ast.CXProgram, expr *ast.CXExpression, symbolsData *SymbolsData) {
	expression, err := prgrm.GetCXAtomicOp(expr.Index)
	if err != nil {
		panic(err)
	}

	if expr.IsMethodCall() {
		var inpIdx ast.CXArgumentIndex = -1
		var outIdx ast.CXArgumentIndex = -1
		if len(expression.GetInputs(prgrm)) > 0 {
			expressionInputTypeSig := prgrm.GetCXTypeSignatureFromArray(expression.GetInputs(prgrm)[0])
			if expressionInputTypeSig.Name != "" && expressionInputTypeSig.Type == ast.TYPE_CXARGUMENT_DEPRECATE {
				inpIdx = ast.CXArgumentIndex(expressionInputTypeSig.Meta)
			}
		}

		if len(expression.GetOutputs(prgrm)) > 0 {
			expressionOutputTypeSig := prgrm.GetCXTypeSignatureFromArray(expression.GetOutputs(prgrm)[0])
			if expressionOutputTypeSig.Name != "" && expressionOutputTypeSig.Type == ast.TYPE_CXARGUMENT_DEPRECATE {
				outIdx = ast.CXArgumentIndex(expressionOutputTypeSig.Meta)
			}
		}

		if inpIdx != -1 {
			inpPkg, err := prgrm.GetPackageFromArray(prgrm.CXArgs[inpIdx].Package)
			if err != nil {
				panic(err)
			}

			if argInpTypeSignature, err := lookupSymbol(prgrm, inpPkg.Name, prgrm.CXArgs[inpIdx].Name, symbolsData); err != nil {
				if outIdx == -1 {
					panic("")
				}

				outPkg, err := prgrm.GetPackageFromArray(prgrm.CXArgs[outIdx].Package)
				if err != nil {
					panic(err)
				}

				argOutTypeSignature, err := lookupSymbol(prgrm, outPkg.Name, prgrm.CXArgs[outIdx].Name, symbolsData)
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

					typeSig := ast.GetCXTypeSignatureRepresentationOfCXArg(prgrm, &prgrm.CXArgs[outIdx])
					typeSigIdx := prgrm.AddCXTypeSignatureInArray(typeSig)
					prgrm.CXAtomicOps[expr.Index].Inputs.Fields = append([]ast.CXTypeSignatureIndex{typeSigIdx}, prgrm.CXAtomicOps[expr.Index].Inputs.Fields...)
					prgrm.CXAtomicOps[expr.Index].Outputs.Fields = prgrm.CXAtomicOps[expr.Index].Outputs.Fields[1:]
					prgrm.CXArgs[outIdx].Fields = prgrm.CXArgs[outIdx].Fields[:len(prgrm.CXArgs[outIdx].Fields)-1]
				}
			} else {
				var argInp *ast.CXArgument = &ast.CXArgument{}
				if argInpTypeSignature.Type == ast.TYPE_CXARGUMENT_DEPRECATE {
					argInp = prgrm.GetCXArgFromArray(ast.CXArgumentIndex(argInpTypeSignature.Meta))
				} else if argInpTypeSignature.Type == ast.TYPE_ATOMIC || argInpTypeSignature.Type == ast.TYPE_POINTER_ATOMIC {
				} else {
					panic("type is not known")
				}

				// then we found an input
				if len(prgrm.CXArgs[inpIdx].Fields) > 0 && argInp != nil {
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
					argOutTypeSignature, err := lookupSymbol(prgrm, outPkg.Name, prgrm.CXArgs[outIdx].Name, symbolsData)
					if err != nil {
						panic(err)
					}

					var argOut *ast.CXArgument = &ast.CXArgument{StructType: nil, ArgDetails: &ast.CXArgumentDebug{}}
					var argOutType string
					if argOutTypeSignature.Type == ast.TYPE_CXARGUMENT_DEPRECATE {
						argOut = prgrm.GetCXArgFromArray(ast.CXArgumentIndex(argOutTypeSignature.Meta))
						argOutType = argOut.Type.Name()
					} else if argOutTypeSignature.Type == ast.TYPE_ATOMIC {
						argOutType = types.Code(argOutTypeSignature.Meta).Name()
					} else if argOutTypeSignature.Type == ast.TYPE_POINTER_ATOMIC {
						argOutType = types.Code(argOutTypeSignature.Meta).Name()
					} else {
						panic("type is not known")
					}

					strct := argOut.StructType
					if strct == nil {
						println(ast.CompilationError(argOut.ArgDetails.FileName, argOut.ArgDetails.FileLine), fmt.Sprintf("illegal method call or field access on identifier '%s' of primitive type '%s'", argOutTypeSignature.Name, argOutType))
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

			argOutTypeSignature, err := lookupSymbol(prgrm, outPkg.Name, prgrm.CXArgs[outIdx].Name, symbolsData)
			if err != nil {
				println(ast.CompilationError(prgrm.CXArgs[outIdx].ArgDetails.FileName, prgrm.CXArgs[outIdx].ArgDetails.FileLine), fmt.Sprintf("identifier '%s' does not exist", prgrm.CXArgs[outIdx].Name))
				os.Exit(constants.CX_COMPILATION_ERROR)
			}

			var argOut *ast.CXArgument = &ast.CXArgument{}
			if argOutTypeSignature.Type == ast.TYPE_CXARGUMENT_DEPRECATE {
				argOut = prgrm.GetCXArgFromArray(ast.CXArgumentIndex(argOutTypeSignature.Meta))

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

					typeSig := ast.GetCXTypeSignatureRepresentationOfCXArg(prgrm, &prgrm.CXArgs[outIdx])
					typeSigIdx := prgrm.AddCXTypeSignatureInArray(typeSig)
					newInputs := &ast.CXStruct{Fields: []ast.CXTypeSignatureIndex{typeSigIdx}}
					if prgrm.CXAtomicOps[expr.Index].Inputs != nil {
						for _, typeSig := range prgrm.CXAtomicOps[expr.Index].Inputs.Fields {
							newInputs.AddField_CXAtomicOps(prgrm, typeSig)
						}
					}
					prgrm.CXAtomicOps[expr.Index].Inputs = newInputs
					prgrm.CXAtomicOps[expr.Index].Outputs.Fields = prgrm.CXAtomicOps[expr.Index].Outputs.Fields[1:]
					prgrm.CXArgs[outIdx].Fields = prgrm.CXArgs[outIdx].Fields[:len(prgrm.CXArgs[outIdx].Fields)-1]
				}
			} else if argOutTypeSignature.Type == ast.TYPE_ATOMIC || argOutTypeSignature.Type == ast.TYPE_POINTER_ATOMIC {
				// TODO: improve
				// do nothing for now since len(prgrm.CXArgs[outIdx].Fields) > 0
				// of type atomics and pointers are always zero
			} else {
				panic("type is not known")
			}

		}
	}
}

// GiveOffset
func GiveOffset(prgrm *ast.CXProgram, symbolsData *SymbolsData, typeSigIdx ast.CXTypeSignatureIndex) {
	symTypeSig := prgrm.GetCXTypeSignatureFromArray(typeSigIdx)
	symTypeSigName := StripNameNumber(symTypeSig.Name)
	if symTypeSigName != "" {
		symPkg, err := prgrm.GetPackageFromArray(symTypeSig.Package)
		if err != nil {
			panic(err)
		}

		currFn, err := prgrm.GetCurrentFunction()
		if err != nil {
			// TODO: improve error handling
			panic("error getting current function")
		}

		if !currFn.IsLocalVariable(symTypeSigName) {
			GetGlobalSymbol(prgrm, symbolsData, symPkg, symTypeSigName)
		}

		argTypeSignature, err := lookupSymbol(prgrm, symPkg.Name, symTypeSigName, symbolsData)
		if err == nil {
			ProcessSymbolFields(prgrm, symTypeSig, argTypeSignature)
			CopyArgFields(prgrm, symTypeSig, argTypeSignature)
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
		outputTypeSig := prgrm.GetCXTypeSignatureFromArray(expression.GetOutputs(prgrm)[0])
		if IsTempVar(outputTypeSig.Name) {
			if outputTypeSig.Type == ast.TYPE_CXARGUMENT_DEPRECATE {
				outputArgIdx := ast.CXArgumentIndex(outputTypeSig.Meta)

				expressionInputTypeSig := prgrm.GetCXTypeSignatureFromArray(expression.GetInputs(prgrm)[0])
				if expressionInputTypeSig.Type == ast.TYPE_CXARGUMENT_DEPRECATE {
					expressionInputArg := prgrm.GetCXArgFromArray(ast.CXArgumentIndex(expressionInputTypeSig.Meta))

					expressionInput := expressionInputArg
					// then it's a temporary variable and it needs to adopt its input's type
					prgrm.CXArgs[outputArgIdx].Type = expressionInput.Type
					prgrm.CXArgs[outputArgIdx].PointerTargetType = expressionInput.PointerTargetType
					prgrm.CXArgs[outputArgIdx].Size = expressionInput.Size
					prgrm.CXArgs[outputArgIdx].PreviouslyDeclared = true
				} else if expressionInputTypeSig.Type == ast.TYPE_ATOMIC {
					prgrm.CXArgs[outputArgIdx].Type = types.Code(expressionInputTypeSig.Meta)
					prgrm.CXArgs[outputArgIdx].Size = types.Code(expressionInputTypeSig.Meta).Size()
				} else if expressionInputTypeSig.Type == ast.TYPE_POINTER_ATOMIC {
					prgrm.CXArgs[outputArgIdx].Type = types.Code(expressionInputTypeSig.Meta)
					prgrm.CXArgs[outputArgIdx].Size = types.Code(expressionInputTypeSig.Meta).Size()
				} else if expressionInputTypeSig.Type == ast.TYPE_ARRAY_ATOMIC {
					arrDetails := prgrm.GetCXTypeSignatureArrayFromArray(expressionInputTypeSig.Meta)
					prgrm.CXArgs[outputArgIdx].Type = types.Code(arrDetails.Meta)
					prgrm.CXArgs[outputArgIdx].Size = types.Code(arrDetails.Meta).Size()
				} else {
					panic("type is not known")
				}

			} else if outputTypeSig.Type == ast.TYPE_ATOMIC || outputTypeSig.Type == ast.TYPE_POINTER_ATOMIC {
				expressionInputTypeSig := prgrm.GetCXTypeSignatureFromArray(expression.GetInputs(prgrm)[0])
				if expressionInputTypeSig.Type == ast.TYPE_CXARGUMENT_DEPRECATE {
					expressionInputArg := prgrm.GetCXArgFromArray(ast.CXArgumentIndex(expressionInputTypeSig.Meta))
					outputTypeSig.Meta = int(expressionInputArg.Type)
				} else if expressionInputTypeSig.Type == ast.TYPE_ATOMIC {
					outputTypeSig.Meta = expressionInputTypeSig.Meta
				} else if expressionInputTypeSig.Type == ast.TYPE_POINTER_ATOMIC {
					outputTypeSig.Meta = expressionInputTypeSig.Meta
				} else if expressionInputTypeSig.Type == ast.TYPE_SLICE_ATOMIC {
					sliceDetails := prgrm.GetCXTypeSignatureArrayFromArray(expressionInputTypeSig.Meta)
					outputTypeSig.Meta = sliceDetails.Meta
				} else {
					panic("type is not known")
				}
			} else if outputTypeSig.Type == ast.TYPE_ARRAY_ATOMIC {
				outputArrDetails := prgrm.GetCXTypeSignatureArrayFromArray(outputTypeSig.Meta)
				expressionInputTypeSig := prgrm.GetCXTypeSignatureFromArray(expression.GetInputs(prgrm)[0])
				if expressionInputTypeSig.Type == ast.TYPE_CXARGUMENT_DEPRECATE {
					expressionInputArg := prgrm.GetCXArgFromArray(ast.CXArgumentIndex(expressionInputTypeSig.Meta))
					outputArrDetails.Meta = int(expressionInputArg.Type)
				} else if expressionInputTypeSig.Type == ast.TYPE_ATOMIC {
					outputArrDetails.Meta = expressionInputTypeSig.Meta
				} else if expressionInputTypeSig.Type == ast.TYPE_POINTER_ATOMIC {
					outputArrDetails.Meta = expressionInputTypeSig.Meta
				} else if expressionInputTypeSig.Type == ast.TYPE_ARRAY_ATOMIC {
					arrDetails := prgrm.GetCXTypeSignatureArrayFromArray(expressionInputTypeSig.Meta)
					outputArrDetails.Meta = arrDetails.Meta
				} else {
					panic("type is not known")
				}
			} else {
				panic("type is not known")
			}
		}
	}
}

// CopyArgFields copies 'arg' fields to 'sym' fields.
func CopyArgFields(prgrm *ast.CXProgram, symTypeSignature, argTypeSignature *ast.CXTypeSignature) {
	var sym *ast.CXArgument = &ast.CXArgument{}
	if symTypeSignature.Type == ast.TYPE_CXARGUMENT_DEPRECATE {
		sym = prgrm.GetCXArgFromArray(ast.CXArgumentIndex(symTypeSignature.Meta))
	} else {
		sym = nil
	}

	var arg *ast.CXArgument = &ast.CXArgument{}
	if argTypeSignature.Type == ast.TYPE_CXARGUMENT_DEPRECATE {
		arg = prgrm.GetCXArgFromArray(ast.CXArgumentIndex(argTypeSignature.Meta))
	} else {
		arg = nil
	}

	// TODO: check if this needs a change
	if sym == nil && arg == nil {
		if !(symTypeSignature.Type == ast.TYPE_SLICE_ATOMIC && argTypeSignature.Type == ast.TYPE_SLICE_ATOMIC) {
			// Only copy meta if sym and arg are not type slice atomic.
			symTypeSignature.Meta = argTypeSignature.Meta
		}

		symTypeSignature.Name = argTypeSignature.Name
		symTypeSignature.Package = argTypeSignature.Package
		symTypeSignature.Type = argTypeSignature.Type
		symTypeSignature.Offset = argTypeSignature.Offset

		return
	} else if sym != nil && ast.IsTypePointerAtomic(sym) && argTypeSignature.Type == ast.TYPE_POINTER_ATOMIC {
		symTypeSignature.Name = argTypeSignature.Name
		symTypeSignature.Package = argTypeSignature.Package
		symTypeSignature.Type = argTypeSignature.Type
		symTypeSignature.Meta = argTypeSignature.Meta
		symTypeSignature.Offset = argTypeSignature.Offset

		return

	} else if sym != nil && arg == nil && argTypeSignature.Type == ast.TYPE_POINTER_ATOMIC {
		sym.Name = argTypeSignature.Name
		sym.Package = argTypeSignature.Package
		sym.Type = types.Code(argTypeSignature.Meta)

		sym.Offset = argTypeSignature.Offset

		declSpec := []int{constants.DECL_BASIC, constants.DECL_POINTER}
		for _, spec := range sym.DeclarationSpecifiers {
			// checking if we need to remove or add cxcore.DECL_POINTERs
			// also we could be removing
			switch spec {
			case constants.DECL_INDEXING:
				println(ast.CompilationError(sym.ArgDetails.FileName, sym.ArgDetails.FileLine), "invalid indexing")
			case constants.DECL_DEREF:
				if declSpec[len(declSpec)-1] == constants.DECL_POINTER {
					declSpec = declSpec[:len(declSpec)-1]
				} else {
					println(ast.CompilationError(sym.ArgDetails.FileName, sym.ArgDetails.FileLine), "invalid indirection")
				}

			}
		}
		sym.DeclarationSpecifiers = declSpec

		return
	} else if sym != nil && arg == nil && (len(sym.DeclarationSpecifiers) > 0 || len(sym.DereferenceOperations) > 0) && argTypeSignature.Type == ast.TYPE_ARRAY_ATOMIC {
		arrDetails := prgrm.GetCXTypeSignatureArrayFromArray(argTypeSignature.Meta)
		newArrDetails := *arrDetails
		newArrDetails.Indexes = sym.Indexes
		newArrDetailsIdx := prgrm.AddCXTypeSignatureArrayInArray(&newArrDetails)

		symTypeSignature.Name = argTypeSignature.Name
		symTypeSignature.Package = argTypeSignature.Package
		symTypeSignature.Type = argTypeSignature.Type
		symTypeSignature.Meta = newArrDetailsIdx
		symTypeSignature.Offset = argTypeSignature.Offset
		symTypeSignature.PassBy = sym.PassBy

		return
	} else if sym != nil && arg == nil && (len(sym.DeclarationSpecifiers) > 0 || len(sym.DereferenceOperations) > 0) && argTypeSignature.Type == ast.TYPE_POINTER_ARRAY_ATOMIC {
		arrDetails := prgrm.GetCXTypeSignatureArrayFromArray(argTypeSignature.Meta)

		newArrDetails := *arrDetails
		newArrDetails.Indexes = sym.Indexes
		newArrDetailsIdx := prgrm.AddCXTypeSignatureArrayInArray(&newArrDetails)

		symTypeSignature.Name = argTypeSignature.Name
		symTypeSignature.Package = argTypeSignature.Package
		symTypeSignature.Type = argTypeSignature.Type
		symTypeSignature.Meta = newArrDetailsIdx
		symTypeSignature.Offset = argTypeSignature.Offset
		symTypeSignature.PassBy = sym.PassBy

		for _, decl := range sym.DeclarationSpecifiers {
			if decl == constants.DECL_DEREF {
				symTypeSignature.IsDeref = true
			}
		}

		return
	} else if sym != nil && arg == nil && argTypeSignature.Type == ast.TYPE_SLICE_ATOMIC {
		arrDetails := prgrm.GetCXTypeSignatureArrayFromArray(argTypeSignature.Meta)

		newArrDetails := *arrDetails
		newArrDetails.Indexes = sym.Indexes
		newArrDetailsIdx := prgrm.AddCXTypeSignatureArrayInArray(&newArrDetails)

		symTypeSignature.Name = argTypeSignature.Name
		symTypeSignature.Package = argTypeSignature.Package
		symTypeSignature.Type = argTypeSignature.Type
		symTypeSignature.Meta = newArrDetailsIdx
		symTypeSignature.Offset = argTypeSignature.Offset
		symTypeSignature.PassBy = sym.PassBy

		for _, decl := range sym.DeclarationSpecifiers {
			switch decl {
			case constants.DECL_DEREF, constants.DEREF_SLICE, constants.DEREF_ARRAY:
				symTypeSignature.IsDeref = true
			}
		}

		for _, deref := range sym.DereferenceOperations {
			switch deref {
			case constants.DEREF_SLICE:
				symTypeSignature.IsDeref = true
			}
		}

		return
	} else if sym != nil && arg == nil && argTypeSignature.Type == ast.TYPE_POINTER_SLICE_ATOMIC {
		arrDetails := prgrm.GetCXTypeSignatureArrayFromArray(argTypeSignature.Meta)

		newArrDetails := *arrDetails
		newArrDetails.Indexes = sym.Indexes
		newArrDetailsIdx := prgrm.AddCXTypeSignatureArrayInArray(&newArrDetails)

		symTypeSignature.Name = argTypeSignature.Name
		symTypeSignature.Package = argTypeSignature.Package
		symTypeSignature.Type = argTypeSignature.Type
		symTypeSignature.Meta = newArrDetailsIdx
		symTypeSignature.Offset = argTypeSignature.Offset
		symTypeSignature.PassBy = sym.PassBy

		for _, decl := range sym.DeclarationSpecifiers {
			switch decl {
			case constants.DECL_DEREF, constants.DEREF_SLICE, constants.DEREF_ARRAY:
				symTypeSignature.IsDeref = true
			}
		}

		for _, deref := range sym.DereferenceOperations {
			switch deref {
			case constants.DEREF_SLICE:
				symTypeSignature.IsDeref = true
			}
		}

		return
	} else if sym != nil && arg == nil && (len(sym.DeclarationSpecifiers) > 1 || len(sym.DereferenceOperations) > 1) {
		sym.Name = argTypeSignature.Name
		sym.Package = argTypeSignature.Package
		sym.Type = types.Code(argTypeSignature.Meta)
		sym.Offset = argTypeSignature.Offset

		return
	} else if sym != nil && arg == nil {
		symTypeSignature.Name = argTypeSignature.Name
		symTypeSignature.Package = argTypeSignature.Package
		symTypeSignature.Type = argTypeSignature.Type
		symTypeSignature.Meta = argTypeSignature.Meta
		symTypeSignature.Offset = argTypeSignature.Offset
		symTypeSignature.IsDeref = argTypeSignature.IsDeref
		return
	} else if sym == nil && arg != nil {
		symTypeSignature.Name = argTypeSignature.Name
		symTypeSignature.Package = argTypeSignature.Package
		symTypeSignature.Type = argTypeSignature.Type
		symTypeSignature.Meta = argTypeSignature.Meta
		symTypeSignature.Offset = argTypeSignature.Offset
		symTypeSignature.IsDeref = argTypeSignature.IsDeref

		return
	}

	// This is for test-pointers when IsTypeAtomic includes types.IDENTIFIER
	// else if sym == nil && arg != nil && symTypeSignature.PassBy != 0 && arg.Type.IsPrimitive() {
	// 	symTypeSignature.Meta = int(arg.Type)
	// 	symTypeSignature.Offset = arg.Offset

	// 	return
	// }

	sym.Name = arg.Name
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
}

// ProcessSymbolFields copies the correct field values for the sym.Fields from their struct fields.
func ProcessSymbolFields(prgrm *ast.CXProgram, symTypeSignature, argTypeSignature *ast.CXTypeSignature) {
	var arg *ast.CXArgument = &ast.CXArgument{}
	if argTypeSignature.Type == ast.TYPE_CXARGUMENT_DEPRECATE {
		arg = prgrm.GetCXArgFromArray(ast.CXArgumentIndex(argTypeSignature.Meta))
	} else {
		// panic("type is type cxargument deprecate\n\n")
		// ProcessSymbolFields is only needed for struct types
		// So we return from here
		return
	}

	var sym *ast.CXArgument = &ast.CXArgument{}
	if symTypeSignature.Type == ast.TYPE_CXARGUMENT_DEPRECATE {
		sym = prgrm.GetCXArgFromArray(ast.CXArgumentIndex(symTypeSignature.Meta))
	} else {
		// panic("type is type cxargument deprecate\n\n")
		// ProcessSymbolFields is only needed for struct types
		// So we return from here
		return
	}

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
					methodOutputTypeSig := prgrm.GetCXTypeSignatureFromArray(methodOutputs[0])

					var methodOutputArg *ast.CXArgument = &ast.CXArgument{}
					if methodOutputTypeSig.Type == ast.TYPE_CXARGUMENT_DEPRECATE {
						methodOutputArg = prgrm.GetCXArgFromArray(ast.CXArgumentIndex(methodOutputTypeSig.Meta))
					} else {
						panic("type is not cx argument deprecate\n\n")
					}

					field.Type = methodOutputArg.Type
					field.PointerTargetType = methodOutputArg.PointerTargetType
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

			for _, typeSignatureIdx := range strct.Fields {
				typeSignature := prgrm.GetCXTypeSignatureFromArray(typeSignatureIdx)

				if nameField.Name == typeSignature.Name && typeSignature.Type == ast.TYPE_ATOMIC {
					nameField.Type = types.Code(typeSignature.Meta)
					nameField.StructType = nil
					nameField.Size = typeSignature.GetSize(prgrm, false)

					// TODO: this should not be needed.
					if len(nameField.DeclarationSpecifiers) > 0 {
						nameField.DeclarationSpecifiers = append([]int{constants.DECL_BASIC}, nameField.DeclarationSpecifiers[1:]...)
					} else {
						nameField.DeclarationSpecifiers = []int{constants.DECL_BASIC}
					}

					break
				} else if nameField.Name == typeSignature.Name && typeSignature.Type == ast.TYPE_POINTER_ATOMIC {
					nameField.Type = types.Code(typeSignature.Meta)
					nameField.StructType = nil
					nameField.Size = types.Code(typeSignature.Meta).Size()

					nameField.DereferenceOperations = append([]int{constants.DECL_BASIC, constants.DEREF_POINTER}, nameField.DereferenceOperations...)

					break
				} else if nameField.Name == typeSignature.Name && typeSignature.Type == ast.TYPE_ARRAY_ATOMIC {
					typeSignatureArray := prgrm.GetCXTypeSignatureArrayFromArray(typeSignature.Meta)
					nameField.Type = types.Code(typeSignatureArray.Meta)
					nameField.StructType = nil
					nameField.Size = types.Code(typeSignatureArray.Meta).Size()
					nameField.Lengths = typeSignatureArray.Lengths
					sym.Lengths = typeSignatureArray.Lengths

					// TODO: this should not be needed.
					if len(nameField.DeclarationSpecifiers) > 0 {
						nameField.DeclarationSpecifiers = append([]int{constants.DECL_BASIC, constants.DECL_ARRAY}, nameField.DeclarationSpecifiers[1:]...)
					} else {
						nameField.DeclarationSpecifiers = []int{constants.DECL_BASIC}
					}
					break
				} else if nameField.Name == typeSignature.Name && typeSignature.Type == ast.TYPE_SLICE_ATOMIC {
					sliceDetails := prgrm.GetCXTypeSignatureArrayFromArray(typeSignature.Meta)
					nameField.Type = types.Code(sliceDetails.Meta)
					nameField.StructType = nil
					nameField.Size = nameField.Type.Size()
					nameField.Lengths = sliceDetails.Lengths
					sym.Lengths = sliceDetails.Lengths
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
func GetGlobalSymbol(prgrm *ast.CXProgram, symbolsData *SymbolsData, symPkg *ast.CXPackage, ident string) {
	_, err := lookupSymbol(prgrm, symPkg.Name, ident, symbolsData)
	if err != nil {
		if glbl, err := symPkg.GetGlobal(prgrm, ident); err == nil {
			lastIdx := len(*&symbolsData.symbolsIndex) - 1

			// add name to symbols
			symbolsData.symbols = append(symbolsData.symbols, glbl)
			(symbolsData.symbolsIndex)[lastIdx][symPkg.Name+"."+ident] = len(symbolsData.symbols) - 1
		}
	}
}
