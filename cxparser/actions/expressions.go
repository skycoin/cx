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

func IterationExpressions(prgrm *ast.CXProgram, init []ast.CXExpression, cond []ast.CXExpression, incr []ast.CXExpression, statements []ast.CXExpression) []ast.CXExpression {
	jmpFn := ast.Natives[constants.OP_JMP]

	pkg, err := prgrm.GetCurrentPackage()
	if err != nil {
		panic(err)
	}

	upExprCXLine := ast.MakeCXLineExpression(prgrm, CurrentFile, LineNo, LineStr)
	upExpr := ast.MakeAtomicOperatorExpression(prgrm, jmpFn)
	upExprAtomicOpIdx := upExpr.Index

	prgrm.CXAtomicOps[upExprAtomicOpIdx].Package = ast.CXPackageIndex(pkg.Index)

	trueArg := WritePrimary(prgrm, types.BOOL, encoder.Serialize(true), false)
	trueArgAtomicOp, err := prgrm.GetCXAtomicOpFromExpressions(trueArg, 0)
	if err != nil {
		panic(err)
	}

	// -2 for the cx line expression addition for up and down expr
	upLines := ((len(statements) + len(incr) + len(cond) + 2) * -1) - 2
	downLines := 0

	prgrm.CXAtomicOps[upExprAtomicOpIdx].AddInput(prgrm, trueArgAtomicOp.Outputs[0])
	prgrm.CXAtomicOps[upExprAtomicOpIdx].ThenLines = upLines
	prgrm.CXAtomicOps[upExprAtomicOpIdx].ElseLines = downLines

	downExprCXLine := ast.MakeCXLineExpression(prgrm, CurrentFile, LineNo, LineStr)
	downExpr := ast.MakeAtomicOperatorExpression(prgrm, jmpFn)
	downExprAtomicOp, _, _, err := prgrm.GetOperation(downExpr)
	if err != nil {
		panic(err)
	}
	downExprAtomicOp.Package = ast.CXPackageIndex(pkg.Index)

	lastCondAtomicOp, err := prgrm.GetCXAtomicOpFromExpressions(cond, len(cond)-1)
	if err != nil {
		panic(err)
	}
	lastCondAtomicOpOperator := prgrm.GetFunctionFromArray(lastCondAtomicOp.Operator)

	if len(lastCondAtomicOp.Outputs) < 1 {
		predicate := ast.MakeArgument(MakeGenSym(constants.LOCAL_PREFIX), CurrentFile, LineNo).SetType(prgrm.GetCXArgFromArray(lastCondAtomicOpOperator.Outputs[0]).Type)
		predicate.Package = ast.CXPackageIndex(pkg.Index)
		predicate.PreviouslyDeclared = true
		predicateIdx := prgrm.AddCXArgInArray(predicate)

		lastCondAtomicOp.AddOutput(prgrm, predicateIdx)
		downExprAtomicOp.AddInput(prgrm, predicateIdx)
	} else {
		predicate := prgrm.GetCXArgFromArray(lastCondAtomicOp.Outputs[0])
		predicate.Package = ast.CXPackageIndex(pkg.Index)
		predicate.PreviouslyDeclared = true
		predicateIdx := lastCondAtomicOp.Outputs[0]

		downExprAtomicOp.AddInput(prgrm, predicateIdx)
	}

	thenLines := 0
	// + 1 for the cx line expression addition for the down expr jmp
	elseLines := len(incr) + len(statements) + 1 + 1

	// processing possible breaks
	for i, stat := range statements {
		if stat.IsBreak(prgrm) {
			statAtomicOp, _, _, err := prgrm.GetOperation(&stat)
			if err != nil {
				panic(err)
			}

			statAtomicOp.ThenLines = elseLines - i - 1
		}
	}

	// processing possible continues
	for i, stat := range statements {
		if stat.IsContinue(prgrm) {
			statAtomicOp, _, _, err := prgrm.GetOperation(&stat)
			if err != nil {
				panic(err)
			}

			statAtomicOp.ThenLines = len(statements) - i - 1
		}
	}

	downExprAtomicOp.ThenLines = thenLines
	downExprAtomicOp.ElseLines = elseLines

	exprs := init
	exprs = append(exprs, cond...)
	exprs = append(exprs, *downExprCXLine, *downExpr)
	exprs = append(exprs, statements...)
	exprs = append(exprs, incr...)
	exprs = append(exprs, *upExprCXLine, *upExpr)

	DefineNewScope(exprs)

	return exprs
}

func trueJmpExpressions(prgrm *ast.CXProgram, opcode int) []ast.CXExpression {
	pkg, err := prgrm.GetCurrentPackage()
	if err != nil {
		panic(err)
	}

	exprCXLine := ast.MakeCXLineExpression(prgrm, CurrentFile, LineNo, LineStr)
	expr := ast.MakeAtomicOperatorExpression(prgrm, ast.Natives[opcode])
	exprAtomicOp, _, _, err := prgrm.GetOperation(expr)
	if err != nil {
		panic(err)
	}

	trueArg := WritePrimary(prgrm, types.BOOL, encoder.Serialize(true), false)
	trueArgAtomicOp, err := prgrm.GetCXAtomicOpFromExpressions(trueArg, 0)
	if err != nil {
		panic(err)
	}

	exprAtomicOp.AddInput(prgrm, trueArgAtomicOp.Outputs[0])
	exprAtomicOp.Package = ast.CXPackageIndex(pkg.Index)

	return []ast.CXExpression{*exprCXLine, *expr}
}

func BreakExpressions(prgrm *ast.CXProgram) []ast.CXExpression {
	exprs := trueJmpExpressions(prgrm, constants.OP_BREAK)
	return exprs
}

func ContinueExpressions(prgrm *ast.CXProgram) []ast.CXExpression {
	exprs := trueJmpExpressions(prgrm, constants.OP_CONTINUE)
	return exprs
}

func SelectionExpressions(prgrm *ast.CXProgram, condExprs []ast.CXExpression, thenExprs []ast.CXExpression, elseExprs []ast.CXExpression) []ast.CXExpression {
	DefineNewScope(thenExprs)
	DefineNewScope(elseExprs)

	jmpFn := ast.Natives[constants.OP_JMP]
	pkg, err := prgrm.GetCurrentPackage()
	if err != nil {
		panic(err)
	}

	ifExprCXLine := ast.MakeCXLineExpression(prgrm, CurrentFile, LineNo, LineStr)
	ifExpr := ast.MakeAtomicOperatorExpression(prgrm, jmpFn)
	ifExprAtomicOp, _, _, err := prgrm.GetOperation(ifExpr)
	if err != nil {
		panic(err)
	}
	ifExprAtomicOp.Package = ast.CXPackageIndex(pkg.Index)

	lastCondExprsAtomicOp, err := prgrm.GetCXAtomicOpFromExpressions(condExprs, len(condExprs)-1)
	if err != nil {
		panic(err)
	}
	lastCondExprsAtomicOpOperator := prgrm.GetFunctionFromArray(lastCondExprsAtomicOp.Operator)

	var predicateIdx ast.CXArgumentIndex
	if lastCondExprsAtomicOpOperator == nil && !condExprs[len(condExprs)-1].IsMethodCall() {
		// then it's a literal
		predicateIdx = lastCondExprsAtomicOp.Outputs[0]
	} else {
		// then it's an expression
		predicate := ast.MakeArgument(MakeGenSym(constants.LOCAL_PREFIX), CurrentFile, LineNo)
		if condExprs[len(condExprs)-1].IsMethodCall() {
			// we'll change this once we have access to method's types in
			// ProcessMethodCall
			predicate.SetType(types.BOOL)
			lastCondExprsAtomicOp.Inputs = append(lastCondExprsAtomicOp.Outputs, lastCondExprsAtomicOp.Inputs...)
			lastCondExprsAtomicOp.Outputs = nil
		} else {
			predicate.SetType(prgrm.GetCXArgFromArray(lastCondExprsAtomicOpOperator.Outputs[0]).Type)
		}
		predicate.PreviouslyDeclared = true

		predicateIdx = prgrm.AddCXArgInArray(predicate)
		lastCondExprsAtomicOp.AddOutput(prgrm, predicateIdx)
	}
	ifExprAtomicOp.AddInput(prgrm, predicateIdx)

	thenLines := 0
	// + 1 for cx line expression addition
	elseLines := len(thenExprs) + 1 + 1

	ifExprAtomicOp.ThenLines = thenLines
	ifExprAtomicOp.ElseLines = elseLines

	skipExprCXLine := ast.MakeCXLineExpression(prgrm, CurrentFile, LineNo, LineStr)
	skipExpr := ast.MakeAtomicOperatorExpression(prgrm, jmpFn)
	skipExprAtomicOpIdx := skipExpr.Index

	prgrm.CXAtomicOps[skipExprAtomicOpIdx].Package = ast.CXPackageIndex(pkg.Index)

	trueArg := WritePrimary(prgrm, types.BOOL, encoder.Serialize(true), false)
	trueArgAtomicOp, err := prgrm.GetCXAtomicOpFromExpressions(trueArg, 0)
	if err != nil {
		panic(err)
	}
	skipLines := len(elseExprs)

	prgrm.CXAtomicOps[skipExprAtomicOpIdx].AddInput(prgrm, trueArgAtomicOp.Outputs[0])
	prgrm.CXAtomicOps[skipExprAtomicOpIdx].ThenLines = skipLines
	prgrm.CXAtomicOps[skipExprAtomicOpIdx].ElseLines = 0

	var exprs []ast.CXExpression
	if lastCondExprsAtomicOpOperator != nil || condExprs[len(condExprs)-1].IsMethodCall() {
		exprs = append(exprs, condExprs...)
	}

	exprs = append(exprs, *ifExprCXLine, *ifExpr)
	exprs = append(exprs, thenExprs...)
	exprs = append(exprs, *skipExprCXLine, *skipExpr)
	exprs = append(exprs, elseExprs...)

	return exprs
}

// resolveTypeForUnd tries to determine the type that will be returned from an expression
func resolveTypeForUnd(prgrm *ast.CXProgram, expr *ast.CXExpression) types.Code {
	cxAtomicOp, _, _, err := prgrm.GetOperation(expr)
	if err != nil {
		panic(err)
	}

	cxAtomicOpOperator := prgrm.GetFunctionFromArray(cxAtomicOp.Operator)

	if len(cxAtomicOp.Inputs) > 0 {
		// it's a literal
		return prgrm.GetCXArgFromArray(cxAtomicOp.Inputs[0]).Type
	}
	if len(cxAtomicOp.Outputs) > 0 {
		// it's an expression with an output
		return prgrm.GetCXArgFromArray(cxAtomicOp.Outputs[0]).Type
	}
	if cxAtomicOpOperator == nil {
		// the expression doesn't return anything
		return -1
	}
	if len(cxAtomicOpOperator.Outputs) > 0 {
		// always return first output's type
		return prgrm.GetCXArgFromArray(cxAtomicOpOperator.Outputs[0]).Type
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

func OperatorExpression(prgrm *ast.CXProgram, leftExprs []ast.CXExpression, rightExprs []ast.CXExpression, opcode int) (out []ast.CXExpression) {
	pkg, err := prgrm.GetCurrentPackage()
	if err != nil {
		panic(err)
	}

	lastLeftExprsAtomicOp, err := prgrm.GetCXAtomicOpFromExpressions(leftExprs, len(leftExprs)-1)
	if err != nil {
		panic(err)
	}
	lastLeftExprsAtomicOpOperator := prgrm.GetFunctionFromArray(lastLeftExprsAtomicOp.Operator)

	if len(lastLeftExprsAtomicOp.Outputs) < 1 {
		lastLeftExprsAtomicOpOperatorOutput := prgrm.GetCXArgFromArray(lastLeftExprsAtomicOpOperator.Outputs[0])
		name := ast.MakeArgument(MakeGenSym(constants.LOCAL_PREFIX), CurrentFile, LineNo).SetType(resolveTypeForUnd(prgrm, &leftExprs[len(leftExprs)-1]))
		name.Size = lastLeftExprsAtomicOpOperatorOutput.Size
		name.TotalSize = ast.GetSize(prgrm, lastLeftExprsAtomicOpOperatorOutput)
		name.Type = lastLeftExprsAtomicOpOperatorOutput.Type
		name.PointerTargetType = lastLeftExprsAtomicOpOperatorOutput.PointerTargetType
		name.Package = ast.CXPackageIndex(pkg.Index)
		name.PreviouslyDeclared = true

		nameIdx := prgrm.AddCXArgInArray(name)
		lastLeftExprsAtomicOp.Outputs = append(lastLeftExprsAtomicOp.Outputs, nameIdx)
	}

	lastRightExprsAtomicOp, err := prgrm.GetCXAtomicOpFromExpressions(rightExprs, len(rightExprs)-1)
	if err != nil {
		panic(err)
	}
	lastRightExprsAtomicOpOperator := prgrm.GetFunctionFromArray(lastRightExprsAtomicOp.Operator)

	if len(lastRightExprsAtomicOp.Outputs) < 1 {
		name := ast.MakeArgument(MakeGenSym(constants.LOCAL_PREFIX), CurrentFile, LineNo).SetType(resolveTypeForUnd(prgrm, &rightExprs[len(rightExprs)-1]))

		lastRightExprsAtomicOpOperatorOutput := prgrm.GetCXArgFromArray(lastRightExprsAtomicOpOperator.Outputs[0])
		name.Size = lastRightExprsAtomicOpOperatorOutput.Size
		name.TotalSize = ast.GetSize(prgrm, lastRightExprsAtomicOpOperatorOutput)
		name.Type = lastRightExprsAtomicOpOperatorOutput.Type
		name.PointerTargetType = lastRightExprsAtomicOpOperatorOutput.PointerTargetType
		name.Package = ast.CXPackageIndex(pkg.Index)
		name.PreviouslyDeclared = true

		nameIdx := prgrm.AddCXArgInArray(name)
		lastRightExprsAtomicOp.Outputs = append(lastRightExprsAtomicOp.Outputs, nameIdx)
	}

	exprCXLine := ast.MakeCXLineExpression(prgrm, CurrentFile, LineNo, LineStr)
	expr := ast.MakeAtomicOperatorExpression(prgrm, ast.Natives[opcode])
	cxAtomicOp, _, _, err := prgrm.GetOperation(expr)
	if err != nil {
		panic(err)
	}

	// we can't know the type until we compile the full function
	cxAtomicOp.Package = ast.CXPackageIndex(pkg.Index)

	if len(prgrm.GetCXArgFromArray(lastLeftExprsAtomicOp.Outputs[0]).Indexes) > 0 || lastLeftExprsAtomicOpOperator != nil {
		// then it's a function call or an array access
		cxAtomicOp.AddInput(prgrm, lastLeftExprsAtomicOp.Outputs[0])

		if IsTempVar(prgrm.GetCXArgFromArray(lastLeftExprsAtomicOp.Outputs[0]).Name) {
			out = append(out, leftExprs...)
		} else {
			out = append(out, leftExprs[:len(leftExprs)-1]...)
		}
	} else {
		cxAtomicOp.Inputs = append(cxAtomicOp.Inputs, lastLeftExprsAtomicOp.Outputs[0])
	}

	if len(prgrm.GetCXArgFromArray(lastRightExprsAtomicOp.Outputs[0]).Indexes) > 0 || lastRightExprsAtomicOpOperator != nil {
		// then it's a function call or an array access
		cxAtomicOp.AddInput(prgrm, lastRightExprsAtomicOp.Outputs[0])

		if IsTempVar(prgrm.GetCXArgFromArray(lastRightExprsAtomicOp.Outputs[0]).Name) {
			out = append(out, rightExprs...)
		} else {
			out = append(out, rightExprs[:len(rightExprs)-1]...)
		}
	} else {
		cxAtomicOp.Inputs = append(cxAtomicOp.Inputs, lastRightExprsAtomicOp.Outputs[0])
	}

	out = append(out, *exprCXLine, *expr)

	return
}

func UnaryExpression(prgrm *ast.CXProgram, op string, prevExprs []ast.CXExpression) []ast.CXExpression {
	lastPrevExprsAtomicOp, err := prgrm.GetCXAtomicOpFromExpressions(prevExprs, len(prevExprs)-1)
	if err != nil {
		panic(err)
	}

	if len(lastPrevExprsAtomicOp.Outputs) == 0 {
		println(ast.CompilationError(CurrentFile, LineNo), "invalid indirection")
		// needs to be stopped immediately
		os.Exit(constants.CX_COMPILATION_ERROR)
	}

	// Some properties need to be read from the base argument
	// due to how we calculate dereferences at the moment.
	lastPrevExprsAtomicOpOutput := prgrm.GetCXArgFromArray(lastPrevExprsAtomicOp.Outputs[0])
	baseOut := lastPrevExprsAtomicOpOutput
	exprOut := lastPrevExprsAtomicOpOutput.GetAssignmentElement(prgrm)
	switch op {
	case "*":
		exprOut.DereferenceLevels++
		exprOut.DereferenceOperations = append(exprOut.DereferenceOperations, constants.DEREF_POINTER)

		exprOut.DeclarationSpecifiers = append(exprOut.DeclarationSpecifiers, constants.DECL_DEREF)
		// exprOut.IsReference = false
	case "&":
		baseOut.PassBy = constants.PASSBY_REFERENCE
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
			cxAtomicOp, _, _, err := prgrm.GetOperation(expr)
			if err != nil {
				panic(err)
			}
			cxAtomicOp.Package = ast.CXPackageIndex(pkg.Index)

			cxAtomicOp.AddInput(prgrm, ast.CXArgumentIndex(exprOut.Index))
			prevExprs[len(prevExprs)-1] = *expr
		} else {
			panic(err)
		}
	case "-":
		if pkg, err := prgrm.GetCurrentPackage(); err == nil {
			expr := ast.MakeAtomicOperatorExpression(prgrm, ast.Natives[constants.OP_NEG])
			cxAtomicOp, _, _, err := prgrm.GetOperation(expr)
			if err != nil {
				panic(err)
			}
			cxAtomicOp.Package = ast.CXPackageIndex(pkg.Index)
			cxAtomicOp.AddInput(prgrm, ast.CXArgumentIndex(exprOut.Index))
			prevExprs[len(prevExprs)-1] = *expr
		} else {
			panic(err)
		}
	}

	return prevExprs
}

// AssociateReturnExpressions associates the output of `retExprs` to the
// `idx`th output parameter of the current function.
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

	outParam := prgrm.GetCXArgFromArray(fn.Outputs[idx])

	out := ast.MakeArgument(outParam.Name, CurrentFile, LineNo)
	out.SetType(outParam.Type)
	out.StructType = outParam.StructType
	out.PreviouslyDeclared = true
	outIdx := prgrm.AddCXArgInArray(out)

	lastExprAtomicOp, _, _, err := prgrm.GetOperation(&lastExpr)
	if err != nil {
		panic(err)
	}

	lastExprAtomicOpOperator := prgrm.GetFunctionFromArray(lastExprAtomicOp.Operator)
	if lastExprAtomicOpOperator == nil {
		opIdx := prgrm.AddNativeFunctionInArray(ast.Natives[constants.OP_IDENTITY])
		lastExprAtomicOp.Operator = opIdx

		lastExprAtomicOp.Inputs = lastExprAtomicOp.Outputs
		lastExprAtomicOp.Outputs = nil
		lastExprAtomicOp.AddOutput(prgrm, outIdx)
		return retExprs
	} else if len(lastExprAtomicOp.Outputs) > 0 {
		exprCXLine := ast.MakeCXLineExpression(prgrm, CurrentFile, LineNo, LineStr)
		expr := ast.MakeAtomicOperatorExpression(prgrm, ast.Natives[constants.OP_IDENTITY])
		cxAtomicOp, _, _, err := prgrm.GetOperation(expr)
		if err != nil {
			panic(err)
		}
		cxAtomicOp.AddInput(prgrm, lastExprAtomicOp.Outputs[0])
		cxAtomicOp.AddOutput(prgrm, outIdx)

		return append(retExprs, *exprCXLine, *expr)
	} else {
		lastExprAtomicOp.AddOutput(prgrm, outIdx)
		return retExprs
	}
}

// AddJmpToReturnExpressions adds an jump expression that makes a function stop its execution
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

	if len(fn.Outputs) != exprs.Size && exprs.Expressions != nil {
		// lastExpr := retExprs[len(retExprs)-1]
		lastExprCXLine, _ := prgrm.GetPreviousCXLine(retExprs, len(retExprs)-1)

		var plural1 string
		var plural2 string = "s"
		var plural3 string = "were"
		if len(fn.Outputs) > 1 {
			plural1 = "s"
		}
		if exprs.Size == 1 {
			plural2 = ""
			plural3 = "was"
		}

		println(ast.CompilationError(lastExprCXLine.FileName, lastExprCXLine.LineNumber), fmt.Sprintf("function '%s' expects to return %d argument%s, but %d output argument%s %s provided", fn.Name, len(fn.Outputs), plural1, exprs.Size, plural2, plural3))
	}

	exprCXLine := ast.MakeCXLineExpression(prgrm, CurrentFile, LineNo, LineStr)
	// expression to jump to the end of the embedding function
	expr := ast.MakeAtomicOperatorExpression(prgrm, ast.Natives[constants.OP_GOTO])
	cxAtomicOp, _, _, err := prgrm.GetOperation(expr)
	if err != nil {
		panic(err)
	}

	// simulating a label so it gets executed without evaluating a predicate
	cxAtomicOp.Label = MakeGenSym(constants.LABEL_PREFIX)
	cxAtomicOp.ThenLines = types.MAX_INT32
	cxAtomicOp.Package = ast.CXPackageIndex(pkg.Index)

	retExprs = append(retExprs, *exprCXLine, *expr)

	return retExprs
}
