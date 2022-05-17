package actions

import (
	"fmt"
	"os"

	"github.com/skycoin/cx/cx/ast"
	"github.com/skycoin/cx/cx/constants"
	"github.com/skycoin/cx/cx/types"
	"github.com/skycoin/skycoin/src/cipher/encoder"
)

// ReturnExpressions stores the `Size` of the return arguments represented by `Expressions`.
// For example: `return foo() + bar()` is a set of 3 expressions and they represent a single return argument
type ReturnExpressions struct {
	Size        int
	Expressions []ast.CXExpression
}

// IterationExpressions creates series of expressions that will create
// a for loop condition.
//
// Input arguments description:
// 	prgrm - a CXProgram that contains all the data and arrays of the program.
// 	initializeExprs - contains the initialization of the variables used
// 					  for the condition of the iteration expression.
//  conditionExprs - contains the condition of the iteration expression.
//  incrementExprs - contains the increment expr for the variable used
// 					 in the condition.
//  statementExprs - contains the statements inside the iteration expression.
func IterationExpressions(prgrm *ast.CXProgram,
	initializeExprs []ast.CXExpression,
	conditionExprs []ast.CXExpression,
	incrementExprs []ast.CXExpression,
	statementExprs []ast.CXExpression) []ast.CXExpression {
	jmpFn := ast.Natives[constants.OP_JMP]

	pkg, err := prgrm.GetCurrentPackage()
	if err != nil {
		panic(err)
	}

	upExprCXLine := ast.MakeCXLineExpression(prgrm, CurrentFile, LineNo, LineStr)
	upExpr := ast.MakeAtomicOperatorExpression(prgrm, jmpFn)
	upExpressionIdx := upExpr.Index

	prgrm.CXAtomicOps[upExpressionIdx].Package = ast.CXPackageIndex(pkg.Index)

	// -2 for the cx line expression addition for up and down expr
	upLines := ((len(statementExprs) + len(incrementExprs) + len(conditionExprs) + 2) * -1) - 2
	downLines := 0

	trueArg := WritePrimary(prgrm, types.BOOL, encoder.Serialize(true), false)
	prgrm.CXAtomicOps[upExpressionIdx].AddInput(prgrm, ast.CXArgumentIndex(trueArg.Index))
	prgrm.CXAtomicOps[upExpressionIdx].ThenLines = upLines
	prgrm.CXAtomicOps[upExpressionIdx].ElseLines = downLines

	downExprCXLine := ast.MakeCXLineExpression(prgrm, CurrentFile, LineNo, LineStr)
	downExpr := ast.MakeAtomicOperatorExpression(prgrm, jmpFn)
	downExpressionIdx := downExpr.Index

	prgrm.CXAtomicOps[downExpressionIdx].Package = ast.CXPackageIndex(pkg.Index)

	lastCondExpressionIdx := conditionExprs[len(conditionExprs)-1].Index
	lastCondExpressionOperator := prgrm.GetFunctionFromArray(prgrm.CXAtomicOps[lastCondExpressionIdx].Operator)
	lastCondExpressionOperatorOutputs := lastCondExpressionOperator.GetOutputs(prgrm)
	if len(prgrm.CXAtomicOps[lastCondExpressionIdx].Outputs) < 1 {
		predicate := ast.MakeArgument(generateTempVarName(constants.LOCAL_PREFIX), CurrentFile, LineNo).SetType(prgrm.GetCXArgFromArray(lastCondExpressionOperatorOutputs[0]).Type)
		predicate.Package = ast.CXPackageIndex(pkg.Index)
		predicate.PreviouslyDeclared = true
		predicateIdx := prgrm.AddCXArgInArray(predicate)

		prgrm.CXAtomicOps[lastCondExpressionIdx].AddOutput(prgrm, predicateIdx)
		prgrm.CXAtomicOps[downExpressionIdx].AddInput(prgrm, predicateIdx)
	} else {
		predicateIdx := prgrm.CXAtomicOps[lastCondExpressionIdx].Outputs[0]
		prgrm.CXArgs[predicateIdx].Package = ast.CXPackageIndex(pkg.Index)
		prgrm.CXArgs[predicateIdx].PreviouslyDeclared = true

		prgrm.CXAtomicOps[downExpressionIdx].AddInput(prgrm, predicateIdx)
	}

	thenLines := 0
	// + 1 for the cx line expression addition for the down expr jmp
	elseLines := len(incrementExprs) + len(statementExprs) + 1 + 1

	// processing possible breaks
	for i, statementExpr := range statementExprs {
		if statementExpr.IsBreak(prgrm) {
			statementExpressionIdx := statementExpr.Index
			prgrm.CXAtomicOps[statementExpressionIdx].ThenLines = elseLines - i - 1
		}
	}

	// processing possible continues
	for i, statementExpr := range statementExprs {
		if statementExpr.IsContinue(prgrm) {
			statementExpressionIdx := statementExpr.Index
			prgrm.CXAtomicOps[statementExpressionIdx].ThenLines = len(statementExprs) - i - 1
		}
	}

	prgrm.CXAtomicOps[downExpressionIdx].ThenLines = thenLines
	prgrm.CXAtomicOps[downExpressionIdx].ElseLines = elseLines

	exprs := initializeExprs
	exprs = append(exprs, conditionExprs...)
	exprs = append(exprs, *downExprCXLine, *downExpr)
	exprs = append(exprs, statementExprs...)
	exprs = append(exprs, incrementExprs...)
	exprs = append(exprs, *upExprCXLine, *upExpr)

	DefineNewScope(exprs)

	return exprs
}

// trueJmpExpressions makes a true jump expression for BREAKS and CONTINUES.
//
// Input arguments description:
// 	prgrm - a CXProgram that contains all the data and arrays of the program.
// 	opcode - the opcode of the expression. Currently used as either a
//  		 BREAK or CONTINUE expression.
func trueJmpExpressions(prgrm *ast.CXProgram, opcode int) []ast.CXExpression {
	pkg, err := prgrm.GetCurrentPackage()
	if err != nil {
		panic(err)
	}

	exprCXLine := ast.MakeCXLineExpression(prgrm, CurrentFile, LineNo, LineStr)
	expr := ast.MakeAtomicOperatorExpression(prgrm, ast.Natives[opcode])
	expressionIdx := expr.Index

	trueArg := WritePrimary(prgrm, types.BOOL, encoder.Serialize(true), false)
	prgrm.CXAtomicOps[expressionIdx].AddInput(prgrm, ast.CXArgumentIndex(trueArg.Index))
	prgrm.CXAtomicOps[expressionIdx].Package = ast.CXPackageIndex(pkg.Index)

	return []ast.CXExpression{*exprCXLine, *expr}
}

// BreakExpressions makes a true jump expression for BREAKS.
//
// Input arguments description:
// 	prgrm - a CXProgram that contains all the data and arrays of the program.
func BreakExpressions(prgrm *ast.CXProgram) []ast.CXExpression {
	exprs := trueJmpExpressions(prgrm, constants.OP_BREAK)
	return exprs
}

// ContinueExpressions makes a true jump expression for CONTINUES.
//
// Input arguments description:
// 	prgrm - a CXProgram that contains all the data and arrays of the program.
func ContinueExpressions(prgrm *ast.CXProgram) []ast.CXExpression {
	exprs := trueJmpExpressions(prgrm, constants.OP_CONTINUE)
	return exprs
}

// SelectionExpressions creates series of expressions that will create
// an if else condition.
//
// Input arguments description:
// 	prgrm - a CXProgram that contains all the data and arrays of the program.
//  conditionExprs - contains the condition of the selection expression.
// 	thenExprs - contains the statements if condition is true.
//  elseExprs - contains the statements if condition is false.
func SelectionExpressions(prgrm *ast.CXProgram, conditionExprs []ast.CXExpression, thenExprs []ast.CXExpression, elseExprs []ast.CXExpression) []ast.CXExpression {
	DefineNewScope(thenExprs)
	DefineNewScope(elseExprs)

	jmpFn := ast.Natives[constants.OP_JMP]
	pkg, err := prgrm.GetCurrentPackage()
	if err != nil {
		panic(err)
	}

	ifExprCXLine := ast.MakeCXLineExpression(prgrm, CurrentFile, LineNo, LineStr)
	ifExpr := ast.MakeAtomicOperatorExpression(prgrm, jmpFn)
	ifExpressionIdx := ifExpr.Index
	prgrm.CXAtomicOps[ifExpressionIdx].Package = ast.CXPackageIndex(pkg.Index)

	lastCondExpressionIdx := conditionExprs[len(conditionExprs)-1].Index
	lastCondExpressionOperator := prgrm.GetFunctionFromArray(prgrm.CXAtomicOps[lastCondExpressionIdx].Operator)

	var predicateIdx ast.CXArgumentIndex
	if lastCondExpressionOperator == nil && !conditionExprs[len(conditionExprs)-1].IsMethodCall() {
		// then it's a literal
		predicateIdx = prgrm.CXAtomicOps[lastCondExpressionIdx].Outputs[0]
	} else {
		// then it's an expression
		predicate := ast.MakeArgument(generateTempVarName(constants.LOCAL_PREFIX), CurrentFile, LineNo)
		if conditionExprs[len(conditionExprs)-1].IsMethodCall() {
			// we'll change this once we have access to method's types in
			// ProcessMethodCall
			predicate.SetType(types.BOOL)
			prgrm.CXAtomicOps[lastCondExpressionIdx].Inputs = append(prgrm.CXAtomicOps[lastCondExpressionIdx].Outputs, prgrm.CXAtomicOps[lastCondExpressionIdx].Inputs...)
			prgrm.CXAtomicOps[lastCondExpressionIdx].Outputs = nil
		} else {
			lastCondExpressionOperatorOutputs := lastCondExpressionOperator.GetOutputs(prgrm)
			predicate.SetType(prgrm.GetCXArgFromArray(lastCondExpressionOperatorOutputs[0]).Type)
		}
		predicate.PreviouslyDeclared = true

		predicateIdx = prgrm.AddCXArgInArray(predicate)
		prgrm.CXAtomicOps[lastCondExpressionIdx].AddOutput(prgrm, predicateIdx)
	}
	prgrm.CXAtomicOps[ifExpressionIdx].AddInput(prgrm, predicateIdx)

	thenLines := 0
	// + 1 for cx line expression addition
	elseLines := len(thenExprs) + 1 + 1

	prgrm.CXAtomicOps[ifExpressionIdx].ThenLines = thenLines
	prgrm.CXAtomicOps[ifExpressionIdx].ElseLines = elseLines

	skipExprCXLine := ast.MakeCXLineExpression(prgrm, CurrentFile, LineNo, LineStr)
	skipExpr := ast.MakeAtomicOperatorExpression(prgrm, jmpFn)
	skipExpressionIdx := skipExpr.Index

	prgrm.CXAtomicOps[skipExpressionIdx].Package = ast.CXPackageIndex(pkg.Index)

	skipLines := len(elseExprs)

	trueArg := WritePrimary(prgrm, types.BOOL, encoder.Serialize(true), false)
	prgrm.CXAtomicOps[skipExpressionIdx].AddInput(prgrm, ast.CXArgumentIndex(trueArg.Index))
	prgrm.CXAtomicOps[skipExpressionIdx].ThenLines = skipLines
	prgrm.CXAtomicOps[skipExpressionIdx].ElseLines = 0

	var exprs []ast.CXExpression
	if lastCondExpressionOperator != nil || conditionExprs[len(conditionExprs)-1].IsMethodCall() {
		exprs = append(exprs, conditionExprs...)
	}

	exprs = append(exprs, *ifExprCXLine, *ifExpr)
	exprs = append(exprs, thenExprs...)
	exprs = append(exprs, *skipExprCXLine, *skipExpr)
	exprs = append(exprs, elseExprs...)

	return exprs
}

// resolveTypeForUnd tries to determine the type that will be returned from an expression
func resolveTypeForUnd(prgrm *ast.CXProgram, expr *ast.CXExpression) types.Code {
	expression, err := prgrm.GetCXAtomicOp(expr.Index)
	if err != nil {
		panic(err)
	}

	expressionOperator := prgrm.GetFunctionFromArray(expression.Operator)

	if len(expression.Inputs) > 0 {
		// it's a literal
		return prgrm.GetCXArgFromArray(expression.Inputs[0]).Type
	}
	if len(expression.Outputs) > 0 {
		// it's an expression with an output
		return prgrm.GetCXArgFromArray(expression.Outputs[0]).Type
	}
	if expressionOperator == nil {
		// the expression doesn't return anything
		return -1
	}
	expressionOperatorOutputs := expressionOperator.GetOutputs(prgrm)
	if len(expressionOperatorOutputs) > 0 {
		// always return first output's type
		return prgrm.GetCXArgFromArray(expressionOperatorOutputs[0]).Type
	}

	// error
	return -1
}

// IsTempVar ...
//TODO: Delete this function; only called by next function
func IsTempVar(name string) bool {
	if len(name) >= len(constants.LOCAL_PREFIX) && name[:len(constants.LOCAL_PREFIX)] == constants.LOCAL_PREFIX {
		return true
	}
	return false
}

// OperatorExpression creates an expression for an expression OP expression instance.
// i.e. ((expr+expr) * (expr-expr))
//
// Input arguments description:
// 	prgrm - a CXProgram that contains all the data and arrays of the program.
//  leftExprs - exprs on the left side of the operator expression
// 	rightExprs - exprs on the right side of the operator expression
// 	opcode - the opcode of the operator.
func OperatorExpression(prgrm *ast.CXProgram, leftExprs []ast.CXExpression, rightExprs []ast.CXExpression, opcode int) (out []ast.CXExpression) {
	pkg, err := prgrm.GetCurrentPackage()
	if err != nil {
		panic(err)
	}

	lastLeftExpressionIdx := leftExprs[len(leftExprs)-1].Index
	lastLeftExpressionOperator := prgrm.GetFunctionFromArray(prgrm.CXAtomicOps[lastLeftExpressionIdx].Operator)
	lastLeftExpressionOperatorOutputs := lastLeftExpressionOperator.GetOutputs(prgrm)
	if len(prgrm.CXAtomicOps[lastLeftExpressionIdx].Outputs) < 1 {
		lastLeftExpressionOperatorOutput := prgrm.GetCXArgFromArray(lastLeftExpressionOperatorOutputs[0])

		out := ast.MakeArgument(generateTempVarName(constants.LOCAL_PREFIX), CurrentFile, LineNo)
		out.SetType(resolveTypeForUnd(prgrm, &leftExprs[len(leftExprs)-1]))
		out.Size = lastLeftExpressionOperatorOutput.Size
		out.TotalSize = ast.GetArgSize(prgrm, lastLeftExpressionOperatorOutput)
		out.Type = lastLeftExpressionOperatorOutput.Type
		out.PointerTargetType = lastLeftExpressionOperatorOutput.PointerTargetType
		out.Package = ast.CXPackageIndex(pkg.Index)
		out.PreviouslyDeclared = true

		outIdx := prgrm.AddCXArgInArray(out)
		prgrm.CXAtomicOps[lastLeftExpressionIdx].Outputs = append(prgrm.CXAtomicOps[lastLeftExpressionIdx].Outputs, outIdx)
	}

	lastRightExpressionIdx := rightExprs[len(rightExprs)-1].Index
	lastRightExpressionOperator := prgrm.GetFunctionFromArray(prgrm.CXAtomicOps[lastRightExpressionIdx].Operator)
	lastRightExpressionOperatorOutputs := lastRightExpressionOperator.GetOutputs(prgrm)
	if len(prgrm.CXAtomicOps[lastRightExpressionIdx].Outputs) < 1 {
		lastRightExpressionOperatorOutput := prgrm.GetCXArgFromArray(lastRightExpressionOperatorOutputs[0])

		out := ast.MakeArgument(generateTempVarName(constants.LOCAL_PREFIX), CurrentFile, LineNo)
		out.SetType(resolveTypeForUnd(prgrm, &rightExprs[len(rightExprs)-1]))
		out.Size = lastRightExpressionOperatorOutput.Size
		out.TotalSize = ast.GetArgSize(prgrm, lastRightExpressionOperatorOutput)
		out.Type = lastRightExpressionOperatorOutput.Type
		out.PointerTargetType = lastRightExpressionOperatorOutput.PointerTargetType
		out.Package = ast.CXPackageIndex(pkg.Index)
		out.PreviouslyDeclared = true

		outIdx := prgrm.AddCXArgInArray(out)
		prgrm.CXAtomicOps[lastRightExpressionIdx].Outputs = append(prgrm.CXAtomicOps[lastRightExpressionIdx].Outputs, outIdx)
	}

	exprCXLine := ast.MakeCXLineExpression(prgrm, CurrentFile, LineNo, LineStr)
	expr := ast.MakeAtomicOperatorExpression(prgrm, ast.Natives[opcode])
	expressionIdx := expr.Index

	// we can't know the type until we compile the full function
	prgrm.CXAtomicOps[expressionIdx].Package = ast.CXPackageIndex(pkg.Index)

	if len(prgrm.GetCXArgFromArray(prgrm.CXAtomicOps[lastLeftExpressionIdx].Outputs[0]).Indexes) > 0 || lastLeftExpressionOperator != nil {
		// then it's a function call or an array access
		prgrm.CXAtomicOps[expressionIdx].AddInput(prgrm, prgrm.CXAtomicOps[lastLeftExpressionIdx].Outputs[0])

		if IsTempVar(prgrm.GetCXArgFromArray(prgrm.CXAtomicOps[lastLeftExpressionIdx].Outputs[0]).Name) {
			out = append(out, leftExprs...)
		} else {
			out = append(out, leftExprs[:len(leftExprs)-1]...)
		}
	} else {
		prgrm.CXAtomicOps[expressionIdx].Inputs = append(prgrm.CXAtomicOps[expressionIdx].Inputs, prgrm.CXAtomicOps[lastLeftExpressionIdx].Outputs[0])
	}

	if len(prgrm.GetCXArgFromArray(prgrm.CXAtomicOps[lastRightExpressionIdx].Outputs[0]).Indexes) > 0 || lastRightExpressionOperator != nil {
		// then it's a function call or an array access
		prgrm.CXAtomicOps[expressionIdx].AddInput(prgrm, prgrm.CXAtomicOps[lastRightExpressionIdx].Outputs[0])

		if IsTempVar(prgrm.GetCXArgFromArray(prgrm.CXAtomicOps[lastRightExpressionIdx].Outputs[0]).Name) {
			out = append(out, rightExprs...)
		} else {
			out = append(out, rightExprs[:len(rightExprs)-1]...)
		}
	} else {
		prgrm.CXAtomicOps[expressionIdx].Inputs = append(prgrm.CXAtomicOps[expressionIdx].Inputs, prgrm.CXAtomicOps[lastRightExpressionIdx].Outputs[0])
	}

	out = append(out, *exprCXLine, *expr)

	return
}

// UnaryExpression creates an expression for unary operations,
// operations that only need one input, '*', '&', '!', and '-'.
//
// Input arguments description:
// 	prgrm - a CXProgram that contains all the data and arrays of the program.
//  op - the unary operator, '*', '&', '!', and '-'.
//  prevExprs - the array of expressions the unary expression belongs.
func UnaryExpression(prgrm *ast.CXProgram, op string, prevExprs []ast.CXExpression) []ast.CXExpression {
	lastPrevExpression, err := prgrm.GetCXAtomicOpFromExpressions(prevExprs, len(prevExprs)-1)
	if err != nil {
		panic(err)
	}

	if len(lastPrevExpression.Outputs) == 0 {
		println(ast.CompilationError(CurrentFile, LineNo), "invalid indirection")
		// needs to be stopped immediately
		os.Exit(constants.CX_COMPILATION_ERROR)
	}

	// Some properties need to be read from the base argument
	// due to how we calculate dereferences at the moment.
	lastPrevExpressionOutput := prgrm.GetCXArgFromArray(lastPrevExpression.Outputs[0])
	baseOut := lastPrevExpressionOutput
	exprOut := lastPrevExpressionOutput.GetAssignmentElement(prgrm)
	switch op {
	case "*":
		exprOut.DereferenceOperations = append(exprOut.DereferenceOperations, constants.DEREF_POINTER)
		exprOut.DeclarationSpecifiers = append(exprOut.DeclarationSpecifiers, constants.DECL_DEREF)
	case "&":
		baseOut.PassBy = constants.PASSBY_REFERENCE

		// panic(fmt.Sprintf("passby ref baseOut=%+v\n\nexprOut=%+v\n\n", baseOut, exprOut))
		exprOut.DeclarationSpecifiers = append(exprOut.DeclarationSpecifiers, constants.DECL_POINTER)
		if len(baseOut.Fields) == 0 && hasDeclSpec(baseOut, constants.DECL_INDEXING) {
			// If we're referencing an inner element, like an element of a slice (&slc[0])
			// or a field of a struct (&struct.fld) we no longer need to add
			// the OBJECT_HEADER_SIZE to the offset. The runtime uses this field to determine this.
			baseOut.IsInnerReference = true
		}
	case "!":
		if pkg, err := prgrm.GetCurrentPackage(); err == nil {
			expr := ast.MakeAtomicOperatorExpression(prgrm, ast.Natives[constants.OP_BOOL_NOT])
			expression, err := prgrm.GetCXAtomicOp(expr.Index)
			if err != nil {
				panic(err)
			}
			expression.Package = ast.CXPackageIndex(pkg.Index)

			expression.AddInput(prgrm, ast.CXArgumentIndex(exprOut.Index))
			prevExprs[len(prevExprs)-1] = *expr
		} else {
			panic(err)
		}
	case "-":
		if pkg, err := prgrm.GetCurrentPackage(); err == nil {
			expr := ast.MakeAtomicOperatorExpression(prgrm, ast.Natives[constants.OP_NEG])
			expression, err := prgrm.GetCXAtomicOp(expr.Index)
			if err != nil {
				panic(err)
			}
			expression.Package = ast.CXPackageIndex(pkg.Index)
			expression.AddInput(prgrm, ast.CXArgumentIndex(exprOut.Index))
			prevExprs[len(prevExprs)-1] = *expr
		} else {
			panic(err)
		}
	}

	return prevExprs
}

// AssociateReturnExpressions associate the output of `retExprs` to the
// `idx`th output parameter of the current function.
//
// Input arguments description:
// 	prgrm - a CXProgram that contains all the data and arrays of the program.
//  idx - referring to the nth output parameter of the current function.
// 	retExprs - the return expressions to be associated or linked to the
// 			   output parameter of the current function.
func AssociateReturnExpressions(prgrm *ast.CXProgram, idx int, retExprs []ast.CXExpression) []ast.CXExpression {
	var pkg *ast.CXPackage
	var fn *ast.CXFunction
	var err error

	pkg, err = prgrm.GetCurrentPackage()
	if err != nil {
		panic(err)
	}

	fn, err = pkg.GetCurrentFunction(prgrm)
	if err != nil {
		panic(err)
	}

	lastExpr := retExprs[len(retExprs)-1]

	fnOutputs := fn.GetOutputs(prgrm)
	outParam := prgrm.GetCXArgFromArray(fnOutputs[idx])

	out := ast.MakeArgument(outParam.Name, CurrentFile, LineNo)
	out.SetType(outParam.Type)
	out.StructType = outParam.StructType
	out.PreviouslyDeclared = true
	outIdx := prgrm.AddCXArgInArray(out)

	lastExpression, err := prgrm.GetCXAtomicOp(lastExpr.Index)
	if err != nil {
		panic(err)
	}

	lastExpressionOperator := prgrm.GetFunctionFromArray(lastExpression.Operator)
	if lastExpressionOperator == nil {
		opIdx := prgrm.AddNativeFunctionInArray(ast.Natives[constants.OP_IDENTITY])
		lastExpression.Operator = opIdx

		lastExpression.Inputs = lastExpression.Outputs
		lastExpression.Outputs = nil
		lastExpression.AddOutput(prgrm, outIdx)
		return retExprs
	} else if len(lastExpression.Outputs) > 0 {
		exprCXLine := ast.MakeCXLineExpression(prgrm, CurrentFile, LineNo, LineStr)
		expr := ast.MakeAtomicOperatorExpression(prgrm, ast.Natives[constants.OP_IDENTITY])
		expression, err := prgrm.GetCXAtomicOp(expr.Index)
		if err != nil {
			panic(err)
		}
		expression.AddInput(prgrm, lastExpression.Outputs[0])
		expression.AddOutput(prgrm, outIdx)

		return append(retExprs, *exprCXLine, *expr)
	} else {
		lastExpression.AddOutput(prgrm, outIdx)
		return retExprs
	}
}

// AddJmpToReturnExpressions adds an jump expression that makes a function stop its execution.
//
// Input arguments description:
// 	prgrm - a CXProgram that contains all the data and arrays of the program.
//  exprs - contains the array of return expressions where the jump expr will be added.
func AddJmpToReturnExpressions(prgrm *ast.CXProgram, exprs ReturnExpressions) []ast.CXExpression {
	var pkg *ast.CXPackage
	var fn *ast.CXFunction
	var err error

	pkg, err = prgrm.GetCurrentPackage()
	if err != nil {
		panic(err)
	}

	fn, err = pkg.GetCurrentFunction(prgrm)
	if err != nil {
		panic(err)
	}

	retExprs := exprs.Expressions
	fnOutputs := fn.GetOutputs(prgrm)
	if len(fnOutputs) != exprs.Size && exprs.Expressions != nil {
		// lastExpr := retExprs[len(retExprs)-1]
		lastExprCXLine, _ := prgrm.GetPreviousCXLine(retExprs, len(retExprs)-1)

		var plural1 string
		var plural2 string = "s"
		var plural3 string = "were"
		if len(fnOutputs) > 1 {
			plural1 = "s"
		}
		if exprs.Size == 1 {
			plural2 = ""
			plural3 = "was"
		}

		println(ast.CompilationError(lastExprCXLine.FileName, lastExprCXLine.LineNumber), fmt.Sprintf("function '%s' expects to return %d argument%s, but %d output argument%s %s provided", fn.Name, len(fnOutputs), plural1, exprs.Size, plural2, plural3))
	}

	exprCXLine := ast.MakeCXLineExpression(prgrm, CurrentFile, LineNo, LineStr)
	// expression to jump to the end of the embedding function
	expr := ast.MakeAtomicOperatorExpression(prgrm, ast.Natives[constants.OP_GOTO])
	expression, err := prgrm.GetCXAtomicOp(expr.Index)
	if err != nil {
		panic(err)
	}

	// simulating a label so it gets executed without evaluating a predicate
	expression.Label = generateTempVarName(constants.LABEL_PREFIX)
	expression.ThenLines = types.MAX_INT32
	expression.Package = ast.CXPackageIndex(pkg.Index)

	retExprs = append(retExprs, *exprCXLine, *expr)

	return retExprs
}
