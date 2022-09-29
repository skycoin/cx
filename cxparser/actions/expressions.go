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
	typeSig := ast.GetCXTypeSignatureRepresentationOfCXArg(prgrm, trueArg)
	typeSigIdx := prgrm.AddCXTypeSignatureInArray(typeSig)
	prgrm.CXAtomicOps[upExpressionIdx].AddInput(prgrm, typeSigIdx)
	prgrm.CXAtomicOps[upExpressionIdx].ThenLines = upLines
	prgrm.CXAtomicOps[upExpressionIdx].ElseLines = downLines

	downExprCXLine := ast.MakeCXLineExpression(prgrm, CurrentFile, LineNo, LineStr)
	downExpr := ast.MakeAtomicOperatorExpression(prgrm, jmpFn)
	downExpressionIdx := downExpr.Index

	prgrm.CXAtomicOps[downExpressionIdx].Package = ast.CXPackageIndex(pkg.Index)

	lastCondExpressionIdx := conditionExprs[len(conditionExprs)-1].Index
	lastCondExpressionOperator := prgrm.GetFunctionFromArray(prgrm.CXAtomicOps[lastCondExpressionIdx].Operator)
	lastCondExpressionOperatorOutputs := lastCondExpressionOperator.GetOutputs(prgrm)
	if len(prgrm.CXAtomicOps[lastCondExpressionIdx].GetOutputs(prgrm)) < 1 {
		lastCondExpressionOperatorOutputTypeSig := prgrm.GetCXTypeSignatureFromArray(lastCondExpressionOperatorOutputs[0])

		var lastCondExpressionOperatorOutputType types.Code
		var typeSigIdx ast.CXTypeSignatureIndex
		if lastCondExpressionOperatorOutputTypeSig.Type == ast.TYPE_CXARGUMENT_DEPRECATE {
			lastCondExpressionOperatorOutputArg := prgrm.GetCXArgFromArray(ast.CXArgumentIndex(lastCondExpressionOperatorOutputTypeSig.Meta))
			lastCondExpressionOperatorOutputType = lastCondExpressionOperatorOutputArg.Type
			predicate := ast.MakeArgument(generateTempVarName(constants.LOCAL_PREFIX), CurrentFile, LineNo).SetType(lastCondExpressionOperatorOutputType)
			predicate.Package = ast.CXPackageIndex(pkg.Index)
			predicate.PreviouslyDeclared = true

			typeSig := ast.GetCXTypeSignatureRepresentationOfCXArg(prgrm, predicate)
			typeSigIdx = prgrm.AddCXTypeSignatureInArray(typeSig)
		} else if lastCondExpressionOperatorOutputTypeSig.Type == ast.TYPE_ATOMIC {
			var newTypeSig ast.CXTypeSignature
			newTypeSig = *lastCondExpressionOperatorOutputTypeSig
			newTypeSig.Name = generateTempVarName(constants.LOCAL_PREFIX)
			newTypeSig.Package = ast.CXPackageIndex(pkg.Index)
			newTypeSig.Offset = types.Pointer(0)
			typeSigIdx = prgrm.AddCXTypeSignatureInArray(&newTypeSig)
		} else if lastCondExpressionOperatorOutputTypeSig.Type == ast.TYPE_POINTER_ATOMIC {
			var newTypeSig ast.CXTypeSignature
			newTypeSig = *lastCondExpressionOperatorOutputTypeSig
			newTypeSig.Name = generateTempVarName(constants.LOCAL_PREFIX)
			newTypeSig.Package = ast.CXPackageIndex(pkg.Index)
			newTypeSig.Offset = types.Pointer(0)
			typeSigIdx = prgrm.AddCXTypeSignatureInArray(&newTypeSig)
		} else {
			panic("type is not known")
		}

		prgrm.CXAtomicOps[lastCondExpressionIdx].AddOutput(prgrm, typeSigIdx)
		prgrm.CXAtomicOps[downExpressionIdx].AddInput(prgrm, typeSigIdx)
	} else {
		predicateTypeSigIdx := prgrm.CXAtomicOps[lastCondExpressionIdx].GetOutputs(prgrm)[0]
		predicateTypeSig := prgrm.GetCXTypeSignatureFromArray(predicateTypeSigIdx)
		var predicateIdx ast.CXArgumentIndex
		if predicateTypeSig.Type == ast.TYPE_CXARGUMENT_DEPRECATE {
			predicateIdx = ast.CXArgumentIndex(predicateTypeSig.Meta)
			prgrm.CXArgs[predicateIdx].Package = ast.CXPackageIndex(pkg.Index)
			prgrm.CXArgs[predicateIdx].PreviouslyDeclared = true
		} else if predicateTypeSig.Type == ast.TYPE_ATOMIC || predicateTypeSig.Type == ast.TYPE_POINTER_ATOMIC {
			// do nothing
		} else {
			panic("type is not known")
		}

		prgrm.CXAtomicOps[downExpressionIdx].AddInput(prgrm, predicateTypeSigIdx)
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
	typeSig := ast.GetCXTypeSignatureRepresentationOfCXArg(prgrm, trueArg)
	typeSigIdx := prgrm.AddCXTypeSignatureInArray(typeSig)
	prgrm.CXAtomicOps[expressionIdx].AddInput(prgrm, typeSigIdx)
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

	var predicateTypeSigIdx ast.CXTypeSignatureIndex
	if lastCondExpressionOperator == nil && !conditionExprs[len(conditionExprs)-1].IsMethodCall() {
		// then it's a literal
		predicateTypeSigIdx = prgrm.CXAtomicOps[lastCondExpressionIdx].GetOutputs(prgrm)[0]
	} else {
		// then it's an expression
		predicate := ast.MakeArgument(generateTempVarName(constants.LOCAL_PREFIX), CurrentFile, LineNo)
		if conditionExprs[len(conditionExprs)-1].IsMethodCall() {
			// we'll change this once we have access to method's types in
			// ProcessMethodCall
			predicate.SetType(types.BOOL)
			prgrm.CXAtomicOps[lastCondExpressionIdx].Inputs.Fields = append(prgrm.CXAtomicOps[lastCondExpressionIdx].Outputs.Fields, prgrm.CXAtomicOps[lastCondExpressionIdx].Inputs.Fields...)
			prgrm.CXAtomicOps[lastCondExpressionIdx].Outputs = nil
		} else {
			lastCondExpressionOperatorOutputs := lastCondExpressionOperator.GetOutputs(prgrm)
			lastCondExpressionOperatorOutputTypeSig := prgrm.GetCXTypeSignatureFromArray(lastCondExpressionOperatorOutputs[0])

			var lastCondExpressionOperatorOutputArgType types.Code
			if lastCondExpressionOperatorOutputTypeSig.Type == ast.TYPE_CXARGUMENT_DEPRECATE {
				lastCondExpressionOperatorOutputArg := prgrm.GetCXArgFromArray(ast.CXArgumentIndex(lastCondExpressionOperatorOutputTypeSig.Meta))
				lastCondExpressionOperatorOutputArgType = lastCondExpressionOperatorOutputArg.Type
			} else if lastCondExpressionOperatorOutputTypeSig.Type == ast.TYPE_ATOMIC {
				lastCondExpressionOperatorOutputArgType = types.Code(lastCondExpressionOperatorOutputTypeSig.Meta)
			} else if lastCondExpressionOperatorOutputTypeSig.Type == ast.TYPE_POINTER_ATOMIC {
				lastCondExpressionOperatorOutputArgType = types.Code(lastCondExpressionOperatorOutputTypeSig.Meta)
			} else {
				panic("type is not known")
			}
			predicate.SetType(lastCondExpressionOperatorOutputArgType)
		}
		predicate.PreviouslyDeclared = true
		predicate.Package = ast.CXPackageIndex(pkg.Index)

		predicateTypeSig := ast.GetCXTypeSignatureRepresentationOfCXArg(prgrm, predicate)
		predicateTypeSigIdx = prgrm.AddCXTypeSignatureInArray(predicateTypeSig)
		prgrm.CXAtomicOps[lastCondExpressionIdx].AddOutput(prgrm, predicateTypeSigIdx)
	}
	prgrm.CXAtomicOps[ifExpressionIdx].AddInput(prgrm, predicateTypeSigIdx)

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
	typeSig := ast.GetCXTypeSignatureRepresentationOfCXArg(prgrm, trueArg)
	typeSigIdx := prgrm.AddCXTypeSignatureInArray(typeSig)
	prgrm.CXAtomicOps[skipExpressionIdx].AddInput(prgrm, typeSigIdx)
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

	expressionInputs := expression.GetInputs(prgrm)
	expressionOutputs := expression.GetOutputs(prgrm)
	if len(expressionInputs) > 0 {
		expressionInputTypeSig := prgrm.GetCXTypeSignatureFromArray(expressionInputs[0])

		var expressionInputArgType types.Code
		if expressionInputTypeSig.Type == ast.TYPE_CXARGUMENT_DEPRECATE {
			expressionInputArg := prgrm.GetCXArgFromArray(ast.CXArgumentIndex(expressionInputTypeSig.Meta))
			expressionInputArgType = expressionInputArg.Type
		} else if expressionInputTypeSig.Type == ast.TYPE_ATOMIC {
			expressionInputArgType = types.Code(expressionInputTypeSig.Meta)
		} else if expressionInputTypeSig.Type == ast.TYPE_POINTER_ATOMIC {
			expressionInputArgType = types.Code(expressionInputTypeSig.Meta)
		} else {
			panic("type is not known")
		}

		// it's a literal
		return expressionInputArgType
	}

	if len(expressionOutputs) > 0 {
		expressionOutputTypeSig := prgrm.GetCXTypeSignatureFromArray(expressionOutputs[0])
		var expressionOutputArgType types.Code
		if expressionOutputTypeSig.Type == ast.TYPE_CXARGUMENT_DEPRECATE {
			expressionOutputArg := prgrm.GetCXArgFromArray(ast.CXArgumentIndex(expressionOutputTypeSig.Meta))
			expressionOutputArgType = expressionOutputArg.Type
		} else if expressionOutputTypeSig.Type == ast.TYPE_ATOMIC {
			expressionOutputArgType = types.Code(expressionOutputTypeSig.Meta)
		} else if expressionOutputTypeSig.Type == ast.TYPE_POINTER_ATOMIC {
			expressionOutputArgType = types.Code(expressionOutputTypeSig.Meta)
		} else {
			panic("type is not known")
		}

		// it's an expression with an output
		return expressionOutputArgType
	}

	if expressionOperator == nil {
		// the expression doesn't return anything
		return -1
	}
	expressionOperatorOutputs := expressionOperator.GetOutputs(prgrm)
	if len(expressionOperatorOutputs) > 0 {
		expressionOperatorOutputTypeSig := prgrm.GetCXTypeSignatureFromArray(expressionOperatorOutputs[0])

		var expressionOperatorOutputArgType types.Code
		if expressionOperatorOutputTypeSig.Type == ast.TYPE_CXARGUMENT_DEPRECATE {
			expressionOperatorOutputArg := prgrm.GetCXArgFromArray(ast.CXArgumentIndex(expressionOperatorOutputTypeSig.Meta))
			expressionOperatorOutputArgType = expressionOperatorOutputArg.Type
		} else if expressionOperatorOutputTypeSig.Type == ast.TYPE_ATOMIC {
			expressionOperatorOutputArgType = types.Code(expressionOperatorOutputTypeSig.Meta)
		} else if expressionOperatorOutputTypeSig.Type == ast.TYPE_POINTER_ATOMIC {
			expressionOperatorOutputArgType = types.Code(expressionOperatorOutputTypeSig.Meta)
		} else {
			panic("type is not known")
		}

		// always return first output's type
		return expressionOperatorOutputArgType
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
	if len(prgrm.CXAtomicOps[lastLeftExpressionIdx].GetOutputs(prgrm)) < 1 {
		lastLeftExpressionOperatorOutputTypeSig := prgrm.GetCXTypeSignatureFromArray(lastLeftExpressionOperatorOutputs[0])

		var typeSigIdx ast.CXTypeSignatureIndex
		if lastLeftExpressionOperatorOutputTypeSig.Type == ast.TYPE_CXARGUMENT_DEPRECATE {
			lastLeftExpressionOperatorOutputArg := prgrm.GetCXArgFromArray(ast.CXArgumentIndex(lastLeftExpressionOperatorOutputTypeSig.Meta))

			out := ast.MakeArgument(generateTempVarName(constants.LOCAL_PREFIX), CurrentFile, LineNo)
			out.SetType(resolveTypeForUnd(prgrm, &leftExprs[len(leftExprs)-1]))
			out.Size = lastLeftExpressionOperatorOutputArg.Size
			out.Type = lastLeftExpressionOperatorOutputArg.Type
			out.PointerTargetType = lastLeftExpressionOperatorOutputArg.PointerTargetType
			out.Package = ast.CXPackageIndex(pkg.Index)
			out.PreviouslyDeclared = true

			typeSig := ast.GetCXTypeSignatureRepresentationOfCXArg(prgrm, out)
			typeSigIdx = prgrm.AddCXTypeSignatureInArray(typeSig)
		} else if lastLeftExpressionOperatorOutputTypeSig.Type == ast.TYPE_ATOMIC {
			typeSig := &ast.CXTypeSignature{}
			typeSig.Name = generateTempVarName(constants.LOCAL_PREFIX)
			typeSig.Package = lastLeftExpressionOperatorOutputTypeSig.Package
			typeSig.Type = ast.TYPE_ATOMIC
			typeSig.Meta = lastLeftExpressionOperatorOutputTypeSig.Meta
			typeSig.Offset = lastLeftExpressionOperatorOutputTypeSig.Offset

			typeSigIdx = prgrm.AddCXTypeSignatureInArray(typeSig)
		} else if lastLeftExpressionOperatorOutputTypeSig.Type == ast.TYPE_POINTER_ATOMIC {
			typeSig := &ast.CXTypeSignature{}
			typeSig.Name = generateTempVarName(constants.LOCAL_PREFIX)
			typeSig.Package = lastLeftExpressionOperatorOutputTypeSig.Package
			typeSig.Type = ast.TYPE_POINTER_ATOMIC
			typeSig.Meta = lastLeftExpressionOperatorOutputTypeSig.Meta
			typeSig.Offset = lastLeftExpressionOperatorOutputTypeSig.Offset

			typeSigIdx = prgrm.AddCXTypeSignatureInArray(typeSig)
		} else {
			panic("type is not known")
		}

		prgrm.CXAtomicOps[lastLeftExpressionIdx].AddOutput(prgrm, typeSigIdx)
	}

	lastRightExpressionIdx := rightExprs[len(rightExprs)-1].Index
	lastRightExpressionOperator := prgrm.GetFunctionFromArray(prgrm.CXAtomicOps[lastRightExpressionIdx].Operator)
	lastRightExpressionOperatorOutputs := lastRightExpressionOperator.GetOutputs(prgrm)
	if len(prgrm.CXAtomicOps[lastRightExpressionIdx].GetOutputs(prgrm)) < 1 {
		lastRightExpressionOperatorOutputTypeSig := prgrm.GetCXTypeSignatureFromArray(lastRightExpressionOperatorOutputs[0])
		var lastRightExpressionOperatorOutputArg *ast.CXArgument = &ast.CXArgument{}

		var typeSigIdx ast.CXTypeSignatureIndex
		if lastRightExpressionOperatorOutputTypeSig.Type == ast.TYPE_CXARGUMENT_DEPRECATE {
			lastRightExpressionOperatorOutputArg = prgrm.GetCXArgFromArray(ast.CXArgumentIndex(lastRightExpressionOperatorOutputTypeSig.Meta))

			out := ast.MakeArgument(generateTempVarName(constants.LOCAL_PREFIX), CurrentFile, LineNo)
			out.SetType(resolveTypeForUnd(prgrm, &rightExprs[len(rightExprs)-1]))
			out.Size = lastRightExpressionOperatorOutputArg.Size
			out.Type = lastRightExpressionOperatorOutputArg.Type
			out.PointerTargetType = lastRightExpressionOperatorOutputArg.PointerTargetType
			out.Package = ast.CXPackageIndex(pkg.Index)
			out.PreviouslyDeclared = true

			typeSig := ast.GetCXTypeSignatureRepresentationOfCXArg(prgrm, out)
			typeSigIdx = prgrm.AddCXTypeSignatureInArray(typeSig)
		} else if lastRightExpressionOperatorOutputTypeSig.Type == ast.TYPE_ATOMIC {
			typeSig := &ast.CXTypeSignature{}
			typeSig.Name = generateTempVarName(constants.LOCAL_PREFIX)
			typeSig.Package = lastRightExpressionOperatorOutputTypeSig.Package
			typeSig.Type = ast.TYPE_ATOMIC
			typeSig.Meta = lastRightExpressionOperatorOutputTypeSig.Meta
			typeSig.Offset = lastRightExpressionOperatorOutputTypeSig.Offset

			typeSigIdx = prgrm.AddCXTypeSignatureInArray(typeSig)
		} else if lastRightExpressionOperatorOutputTypeSig.Type == ast.TYPE_POINTER_ATOMIC {
			typeSig := &ast.CXTypeSignature{}
			typeSig.Name = generateTempVarName(constants.LOCAL_PREFIX)
			typeSig.Package = lastRightExpressionOperatorOutputTypeSig.Package
			typeSig.Type = ast.TYPE_POINTER_ATOMIC
			typeSig.Meta = lastRightExpressionOperatorOutputTypeSig.Meta
			typeSig.Offset = lastRightExpressionOperatorOutputTypeSig.Offset

			typeSigIdx = prgrm.AddCXTypeSignatureInArray(typeSig)
		} else {
			panic("type is not known")
		}

		prgrm.CXAtomicOps[lastRightExpressionIdx].AddOutput(prgrm, typeSigIdx)
	}

	exprCXLine := ast.MakeCXLineExpression(prgrm, CurrentFile, LineNo, LineStr)
	expr := ast.MakeAtomicOperatorExpression(prgrm, ast.Natives[opcode])
	expressionIdx := expr.Index

	// we can't know the type until we compile the full function
	prgrm.CXAtomicOps[expressionIdx].Package = ast.CXPackageIndex(pkg.Index)

	lastLeftExpressionOutputTypeSig := prgrm.GetCXTypeSignatureFromArray(prgrm.CXAtomicOps[lastLeftExpressionIdx].GetOutputs(prgrm)[0])
	var lastLeftExpressionOutputArg *ast.CXArgument = &ast.CXArgument{}
	if lastLeftExpressionOutputTypeSig.Type == ast.TYPE_CXARGUMENT_DEPRECATE {
		lastLeftExpressionOutputArg = prgrm.GetCXArgFromArray(ast.CXArgumentIndex(lastLeftExpressionOutputTypeSig.Meta))
	} else if lastLeftExpressionOutputTypeSig.Type == ast.TYPE_ATOMIC || lastLeftExpressionOutputTypeSig.Type == ast.TYPE_POINTER_ATOMIC {
		lastLeftExpressionOutputArg = &ast.CXArgument{}
	} else {
		panic("type is not known")
	}

	if len(lastLeftExpressionOutputArg.Indexes) > 0 || lastLeftExpressionOperator != nil {
		typeSigIdx := prgrm.CXAtomicOps[lastLeftExpressionIdx].GetOutputs(prgrm)[0]
		// then it's a function call or an array access
		prgrm.CXAtomicOps[expressionIdx].AddInput(prgrm, typeSigIdx)

		if IsTempVar(lastLeftExpressionOutputTypeSig.Name) {
			out = append(out, leftExprs...)
		} else {
			out = append(out, leftExprs[:len(leftExprs)-1]...)
		}
	} else {
		typeSigIdx := prgrm.CXAtomicOps[lastLeftExpressionIdx].GetOutputs(prgrm)[0]
		prgrm.CXAtomicOps[expressionIdx].AddInput(prgrm, typeSigIdx)
	}

	lastRightExpressionOutputTypeSig := prgrm.GetCXTypeSignatureFromArray(prgrm.CXAtomicOps[lastRightExpressionIdx].GetOutputs(prgrm)[0])
	var lastRightExpressionOutputArg *ast.CXArgument = &ast.CXArgument{}
	if lastRightExpressionOutputTypeSig.Type == ast.TYPE_CXARGUMENT_DEPRECATE {
		lastRightExpressionOutputArg = prgrm.GetCXArgFromArray(ast.CXArgumentIndex(lastRightExpressionOutputTypeSig.Meta))
	} else if lastRightExpressionOutputTypeSig.Type == ast.TYPE_ATOMIC || lastRightExpressionOutputTypeSig.Type == ast.TYPE_POINTER_ATOMIC {
		lastRightExpressionOutputArg = &ast.CXArgument{}
	} else {
		panic("type is not known")
	}

	if len(lastRightExpressionOutputArg.Indexes) > 0 || lastRightExpressionOperator != nil {
		typeSigIdx := prgrm.CXAtomicOps[lastRightExpressionIdx].GetOutputs(prgrm)[0]
		// then it's a function call or an array access
		prgrm.CXAtomicOps[expressionIdx].AddInput(prgrm, typeSigIdx)

		if IsTempVar(lastRightExpressionOutputTypeSig.Name) {
			out = append(out, rightExprs...)
		} else {
			out = append(out, rightExprs[:len(rightExprs)-1]...)
		}
	} else {
		typeSigIdx := prgrm.CXAtomicOps[lastRightExpressionIdx].GetOutputs(prgrm)[0]
		prgrm.CXAtomicOps[expressionIdx].AddInput(prgrm, typeSigIdx)
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

	if len(lastPrevExpression.GetOutputs(prgrm)) == 0 {
		println(ast.CompilationError(CurrentFile, LineNo), "invalid indirection")
		// needs to be stopped immediately
		os.Exit(constants.CX_COMPILATION_ERROR)
	}

	lastPrevExpressionOutputTypeSig := prgrm.GetCXTypeSignatureFromArray(lastPrevExpression.GetOutputs(prgrm)[0])
	var lastPrevExpressionOutputArg *ast.CXArgument = &ast.CXArgument{}
	if lastPrevExpressionOutputTypeSig.Type == ast.TYPE_CXARGUMENT_DEPRECATE {
		lastPrevExpressionOutputArg = prgrm.GetCXArgFromArray(ast.CXArgumentIndex(lastPrevExpressionOutputTypeSig.Meta))

		// Some properties need to be read from the base argument
		// due to how we calculate dereferences at the moment.
		lastPrevExpressionOutput := lastPrevExpressionOutputArg
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

				typeSig := ast.GetCXTypeSignatureRepresentationOfCXArg(prgrm, prgrm.GetCXArgFromArray(ast.CXArgumentIndex(exprOut.Index)))
				typeSigIdx := prgrm.AddCXTypeSignatureInArray(typeSig)
				expression.AddInput(prgrm, typeSigIdx)
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
				typeSig := ast.GetCXTypeSignatureRepresentationOfCXArg(prgrm, prgrm.GetCXArgFromArray(ast.CXArgumentIndex(exprOut.Index)))
				typeSigIdx := prgrm.AddCXTypeSignatureInArray(typeSig)
				expression.AddInput(prgrm, typeSigIdx)
				prevExprs[len(prevExprs)-1] = *expr
			} else {
				panic(err)
			}
		}
	} else if lastPrevExpressionOutputTypeSig.Type == ast.TYPE_ATOMIC {
		switch op {
		case "*":
			// TODO: what is the most efficient alternative for this
			// exprOut.DereferenceOperations = append(exprOut.DereferenceOperations, constants.DEREF_POINTER)
			// exprOut.DeclarationSpecifiers = append(exprOut.DeclarationSpecifiers, constants.DECL_DEREF)

			// a pointer type that we want to get its value
			lastPrevExpressionOutputTypeSig.IsDeref = true
		case "&":
			// TODO: what is the most efficient alternative for this
			// baseOut.PassBy = constants.PASSBY_REFERENCE

			// // panic(fmt.Sprintf("passby ref baseOut=%+v\n\nexprOut=%+v\n\n", baseOut, exprOut))
			// exprOut.DeclarationSpecifiers = append(exprOut.DeclarationSpecifiers, constants.DECL_POINTER)
			// if len(baseOut.Fields) == 0 && hasDeclSpec(baseOut, constants.DECL_INDEXING) {
			// 	// If we're referencing an inner element, like an element of a slice (&slc[0])
			// 	// or a field of a struct (&struct.fld) we no longer need to add
			// 	// the OBJECT_HEADER_SIZE to the offset. The runtime uses this field to determine this.
			// 	baseOut.IsInnerReference = true
			// }
			lastPrevExpressionOutputTypeSig.PassBy = constants.PASSBY_REFERENCE
		case "!":
			if pkg, err := prgrm.GetCurrentPackage(); err == nil {
				expr := ast.MakeAtomicOperatorExpression(prgrm, ast.Natives[constants.OP_BOOL_NOT])
				expression, err := prgrm.GetCXAtomicOp(expr.Index)
				if err != nil {
					panic(err)
				}
				expression.Package = ast.CXPackageIndex(pkg.Index)

				typeSigIdx := lastPrevExpressionOutputTypeSig.Index
				expression.AddInput(prgrm, typeSigIdx)
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

				typeSigIdx := lastPrevExpressionOutputTypeSig.Index
				expression.AddInput(prgrm, typeSigIdx)
				prevExprs[len(prevExprs)-1] = *expr
			} else {
				panic(err)
			}
		}

	} else {
		panic("type is not cx argument deprecate\n\n")
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
	fnOutputTypeSig := prgrm.GetCXTypeSignatureFromArray(fnOutputs[idx])
	var outParam *ast.CXArgument = &ast.CXArgument{}
	var typeSigIdx ast.CXTypeSignatureIndex
	if fnOutputTypeSig.Type == ast.TYPE_CXARGUMENT_DEPRECATE {
		outParam = prgrm.GetCXArgFromArray(ast.CXArgumentIndex(fnOutputTypeSig.Meta))

		out := ast.MakeArgument(outParam.Name, CurrentFile, LineNo)
		out.SetType(outParam.Type)
		out.StructType = outParam.StructType
		out.PreviouslyDeclared = true

		typeSig := ast.GetCXTypeSignatureRepresentationOfCXArg(prgrm, out)
		typeSigIdx = prgrm.AddCXTypeSignatureInArray(typeSig)
	} else if fnOutputTypeSig.Type == ast.TYPE_ATOMIC {
		newCXTypeSig := &ast.CXTypeSignature{}
		newCXTypeSig.Name = fnOutputTypeSig.Name
		newCXTypeSig.Package = fnOutputTypeSig.Package
		newCXTypeSig.Type = ast.TYPE_ATOMIC
		newCXTypeSig.Meta = int(fnOutputTypeSig.Type)
		newCXTypeSig.Offset = fnOutputTypeSig.Offset
		typeSigIdx = prgrm.AddCXTypeSignatureInArray(newCXTypeSig)
	} else if fnOutputTypeSig.Type == ast.TYPE_POINTER_ATOMIC {
		newCXTypeSig := &ast.CXTypeSignature{}
		newCXTypeSig.Name = fnOutputTypeSig.Name
		newCXTypeSig.Package = fnOutputTypeSig.Package
		newCXTypeSig.Type = ast.TYPE_POINTER_ATOMIC
		newCXTypeSig.Meta = int(fnOutputTypeSig.Type)
		newCXTypeSig.Offset = fnOutputTypeSig.Offset
		typeSigIdx = prgrm.AddCXTypeSignatureInArray(newCXTypeSig)
	} else {
		panic("type is not known")
	}

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

		lastExpression.AddOutput(prgrm, typeSigIdx)
		return retExprs
	} else if len(lastExpression.GetOutputs(prgrm)) > 0 {
		exprCXLine := ast.MakeCXLineExpression(prgrm, CurrentFile, LineNo, LineStr)
		expr := ast.MakeAtomicOperatorExpression(prgrm, ast.Natives[constants.OP_IDENTITY])
		expression, err := prgrm.GetCXAtomicOp(expr.Index)
		if err != nil {
			panic(err)
		}

		exprOutTypeSig := lastExpression.GetOutputs(prgrm)[0]
		expression.AddInput(prgrm, exprOutTypeSig)
		expression.AddOutput(prgrm, typeSigIdx)

		return append(retExprs, *exprCXLine, *expr)
	} else {
		lastExpression.AddOutput(prgrm, typeSigIdx)
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
