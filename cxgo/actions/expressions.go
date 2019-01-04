package actions

import (
	. "github.com/skycoin/cx/cx"
	"github.com/skycoin/skycoin/src/cipher/encoder"
)

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

func trueJmpExpressions() []*CXExpression {
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

func BreakExpressions() []*CXExpression {
	exprs := trueJmpExpressions()
	exprs[0].IsBreak = true
	return exprs
}

func ContinueExpressions() []*CXExpression {
	exprs := trueJmpExpressions()
	exprs[0].IsContinue = true
	return exprs
}

func SelectionExpressions(condExprs []*CXExpression, thenExprs []*CXExpression, elseExprs []*CXExpression) []*CXExpression {
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

func UndefinedTypeOperation(leftExprs []*CXExpression, rightExprs []*CXExpression, operator *CXFunction) (out []*CXExpression) {
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
			out = append(out, leftExprs[:len(leftExprs)-1]...)
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
			out = append(out, rightExprs[:len(rightExprs)-1]...)
		}
	} else {
		expr.Inputs = append(expr.Inputs, rightExprs[len(rightExprs)-1].Outputs[0])
	}

	out = append(out, expr)

	return
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

	return UndefinedTypeOperation(leftExprs, rightExprs, operator)
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
