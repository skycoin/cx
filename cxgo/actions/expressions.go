package actions

import (
	"fmt"
	"os"

	"github.com/skycoin/skycoin/src/cipher/encoder"

	. "github.com/skycoin/cx/cx"
)

// ReturnExpressions stores the `Size` of the return arguments represented by `Expressions`.
// For example: `return foo() + bar()` is a set of 3 expressions and they represent a single return argument
type ReturnExpressions struct {
	Size        int
	Expressions []*CXExpression
}

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

	DefineNewScope(exprs)

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
	DefineNewScope(thenExprs)
	DefineNewScope(elseExprs)

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

// resolveTypeForUnd tries to determine the type that will be returned from an expression
func resolveTypeForUnd(expr *CXExpression) int {
	if len(expr.Inputs) > 0 {
		// it's a literal
		return expr.Inputs[0].Type
	}
	if len(expr.Outputs) > 0 {
		// it's an expression with an output
		return expr.Outputs[0].Type
	}
	if expr.Operator == nil {
		// the expression doesn't return anything
		return -1
	}
	if len(expr.Operator.Outputs) > 0 {
		// always return first output's type
		return expr.Operator.Outputs[0].Type
	}

	// error
	return -1
}

func UndefinedTypeOperation(leftExprs []*CXExpression, rightExprs []*CXExpression, operator *CXFunction) (out []*CXExpression) {
	pkg, err := PRGRM.GetCurrentPackage()
	if err != nil {
		panic(err)
	}

	if len(leftExprs[len(leftExprs)-1].Outputs) < 1 {
		name := MakeArgument(MakeGenSym(LOCAL_PREFIX), CurrentFile, LineNo).AddType(TypeNames[resolveTypeForUnd(leftExprs[len(leftExprs)-1])])
		name.Size = leftExprs[len(leftExprs)-1].Operator.Outputs[0].Size
		name.TotalSize = GetSize(leftExprs[len(leftExprs)-1].Operator.Outputs[0])
		name.Type = leftExprs[len(leftExprs)-1].Operator.Outputs[0].Type
		name.Package = pkg
		name.PreviouslyDeclared = true

		leftExprs[len(leftExprs)-1].Outputs = append(leftExprs[len(leftExprs)-1].Outputs, name)
	}

	if len(rightExprs[len(rightExprs)-1].Outputs) < 1 {
		name := MakeArgument(MakeGenSym(LOCAL_PREFIX), CurrentFile, LineNo).AddType(TypeNames[resolveTypeForUnd(rightExprs[len(rightExprs)-1])])

		name.Size = rightExprs[len(rightExprs)-1].Operator.Outputs[0].Size
		name.TotalSize = GetSize(rightExprs[len(rightExprs)-1].Operator.Outputs[0])
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
		// Handling special case of arguments being strings.
		// In this case we use `str.concat`.
		rightLen := len(rightExprs) - 1
		if rightLen >= 0 && len(rightExprs[rightLen].Outputs) > 0 && rightExprs[rightLen].Outputs[0].Type == TYPE_STR {
			operator = Natives[OP_STR_CONCAT]
		} else {
			operator = Natives[OP_UND_ADD]
		}
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
	if len(prevExprs[len(prevExprs)-1].Outputs) == 0 {
		println(CompilationError(CurrentFile, LineNo), "invalid indirection")
		// needs to be stopped immediately
		os.Exit(CX_COMPILATION_ERROR)
	}

	// Some properties need to be read from the base argument
	// due to how we calculate dereferences at the moment.
	baseOut := prevExprs[len(prevExprs)-1].Outputs[0]
	exprOut := GetAssignmentElement(prevExprs[len(prevExprs)-1].Outputs[0])
	switch op {
	case "*":
		exprOut.DereferenceLevels++
		exprOut.DereferenceOperations = append(exprOut.DereferenceOperations, DEREF_POINTER)
		if !exprOut.IsArrayFirst {
			exprOut.IsDereferenceFirst = true
		}
		exprOut.DeclarationSpecifiers = append(exprOut.DeclarationSpecifiers, DECL_DEREF)
		exprOut.IsReference = false
	case "&":
		baseOut.PassBy = PASSBY_REFERENCE
		exprOut.DeclarationSpecifiers = append(exprOut.DeclarationSpecifiers, DECL_POINTER)
		if len(baseOut.Fields) == 0 && hasDeclSpec(baseOut, DECL_INDEXING) {
			// If we're referencing an inner element, like an element of a slice (&slc[0])
			// or a field of a struct (&struct.fld) we no longer need to add
			// the OBJECT_HEADER_SIZE to the offset. The runtime uses this field to determine this.
			baseOut.IsInnerReference = true
		}
	case "!":
		if pkg, err := PRGRM.GetCurrentPackage(); err == nil {
			expr := MakeExpression(Natives[OP_BOOL_NOT], CurrentFile, LineNo)
			expr.Package = pkg

			expr.AddInput(exprOut)

			prevExprs[len(prevExprs)-1] = expr
		} else {
			panic(err)
		}
	case "-":
		if pkg, err := PRGRM.GetCurrentPackage(); err == nil {
			expr := MakeExpression(Natives[OP_UND_NEG], CurrentFile, LineNo)
			expr.Package = pkg
			expr.AddInput(exprOut)
			prevExprs[len(prevExprs)-1] = expr
		} else {
			panic(err)
		}
	}
	return prevExprs
}

// AssociateReturnExpressions associates the output of `retExprs` to the
// `idx`th output parameter of the current function.
func AssociateReturnExpressions(idx int, retExprs []*CXExpression) []*CXExpression {
	var pkg *CXPackage
	var fn *CXFunction
	var err error

	pkg, err = PRGRM.GetCurrentPackage()
	if err != nil {
		panic(err)
	}

	fn, err = pkg.GetCurrentFunction()
	if err != nil {
		panic(err)
	}

	lastExpr := retExprs[len(retExprs)-1]

	outParam := fn.Outputs[idx]

	out := MakeArgument(outParam.Name, CurrentFile, LineNo)
	out.AddType(TypeNames[outParam.Type])
	out.CustomType = outParam.CustomType
	out.PreviouslyDeclared = true

	if lastExpr.Operator == nil {
		lastExpr.Operator = Natives[OP_IDENTITY]

		lastExpr.Inputs = lastExpr.Outputs
		lastExpr.Outputs = nil
		lastExpr.AddOutput(out)

		return retExprs
	} else if len(lastExpr.Outputs) > 0 {
		expr := MakeExpression(Natives[OP_IDENTITY], CurrentFile, LineNo)
		expr.AddInput(lastExpr.Outputs[0])
		expr.AddOutput(out)

		return append(retExprs, expr)
	} else {
		lastExpr.AddOutput(out)

		return retExprs
	}
}

// AddJmpToReturnExpressions adds an jump expression that makes a function stop its execution
func AddJmpToReturnExpressions(exprs ReturnExpressions) []*CXExpression {
	var pkg *CXPackage
	var fn *CXFunction
	var err error

	pkg, err = PRGRM.GetCurrentPackage()
	if err != nil {
		panic(err)
	}

	fn, err = pkg.GetCurrentFunction()
	if err != nil {
		panic(err)
	}

	retExprs := exprs.Expressions

	if len(fn.Outputs) != exprs.Size && exprs.Expressions != nil {
		lastExpr := retExprs[len(retExprs)-1]

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

		println(CompilationError(lastExpr.FileName, lastExpr.FileLine), fmt.Sprintf("function '%s' expects to return %d argument%s, but %d output argument%s %s provided", fn.Name, len(fn.Outputs), plural1, exprs.Size, plural2, plural3))
	}

	// expression to jump to the end of the embedding function
	expr := MakeExpression(Natives[OP_JMP], CurrentFile, LineNo)

	// simulating a label so it gets executed without evaluating a predicate
	expr.Label = MakeGenSym(LABEL_PREFIX)
	expr.ThenLines = MAX_INT32
	expr.Package = pkg

	arg := MakeArgument("", CurrentFile, LineNo).AddType("bool")
	arg.Package = pkg

	expr.AddInput(arg)

	retExprs = append(retExprs, expr)

	return retExprs
}
