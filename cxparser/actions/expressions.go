package actions

import (
	"fmt"
	"os"

	"github.com/skycoin/cx/cx/ast"
	"github.com/skycoin/cx/cx/constants"
	"github.com/skycoin/skycoin/src/cipher/encoder"
)

// ReturnExpressions stores the `Size` of the return arguments represented by `Expressions`.
// For example: `return foo() + bar()` is a set of 3 expressions and they represent a single return argument
type ReturnExpressions struct {
	Size        int
	Expressions []*ast.CXExpression
}

func IterationExpressions(init []*ast.CXExpression, cond []*ast.CXExpression, incr []*ast.CXExpression, statements []*ast.CXExpression) []*ast.CXExpression {
	jmpFn := ast.Natives[constants.OP_JMP]

	pkg, err := AST.GetCurrentPackage()
	if err != nil {
		panic(err)
	}

	upExpr := ast.MakeExpression(jmpFn, CurrentFile, LineNo)
	upExpr.Package = pkg

	trueArg := WritePrimary(constants.TYPE_BOOL, encoder.Serialize(true), false)

	upLines := (len(statements) + len(incr) + len(cond) + 2) * -1
	downLines := 0

	upExpr.AddInput(trueArg[0].Outputs[0])
	upExpr.ThenLines = upLines
	upExpr.ElseLines = downLines

	downExpr := ast.MakeExpression(jmpFn, CurrentFile, LineNo)
	downExpr.Package = pkg

	if len(cond[len(cond)-1].Outputs) < 1 {
		predicate := ast.MakeArgument(MakeGenSym(constants.LOCAL_PREFIX), CurrentFile, LineNo).AddType(constants.TypeNames[cond[len(cond)-1].Operator.Outputs[0].Type])
		predicate.ArgDetails.Package = pkg
		predicate.PreviouslyDeclared = true
		cond[len(cond)-1].AddOutput(predicate)
		downExpr.AddInput(predicate)
	} else {
		predicate := cond[len(cond)-1].Outputs[0]
		predicate.ArgDetails.Package = pkg
		predicate.PreviouslyDeclared = true
		downExpr.AddInput(predicate)
	}

	thenLines := 0
	elseLines := len(incr) + len(statements) + 1

	// processing possible breaks
	for i, stat := range statements {
		if stat.IsBreak() {
			stat.ThenLines = elseLines - i - 1
		}
	}

	// processing possible continues
	for i, stat := range statements {
		if stat.IsContinue() {
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

func trueJmpExpressions(opcode int) []*ast.CXExpression {
	pkg, err := AST.GetCurrentPackage()
	if err != nil {
		panic(err)
	}

	expr := ast.MakeExpression(ast.Natives[opcode], CurrentFile, LineNo)

	trueArg := WritePrimary(constants.TYPE_BOOL, encoder.Serialize(true), false)
	expr.AddInput(trueArg[0].Outputs[0])

	expr.Package = pkg

	return []*ast.CXExpression{expr}
}

func BreakExpressions() []*ast.CXExpression {
	exprs := trueJmpExpressions(constants.OP_BREAK)
	return exprs
}

func ContinueExpressions() []*ast.CXExpression {
	exprs := trueJmpExpressions(constants.OP_CONTINUE)
	return exprs
}

func SelectionExpressions(condExprs []*ast.CXExpression, thenExprs []*ast.CXExpression, elseExprs []*ast.CXExpression) []*ast.CXExpression {
	DefineNewScope(thenExprs)
	DefineNewScope(elseExprs)

	jmpFn := ast.Natives[constants.OP_JMP]
	pkg, err := AST.GetCurrentPackage()
	if err != nil {
		panic(err)
	}
	ifExpr := ast.MakeExpression(jmpFn, CurrentFile, LineNo)
	ifExpr.Package = pkg

	var predicate *ast.CXArgument
	if condExprs[len(condExprs)-1].Operator == nil && !condExprs[len(condExprs)-1].IsMethodCall() {
		// then it's a literal
		predicate = condExprs[len(condExprs)-1].Outputs[0]
	} else {
		// then it's an expression
		predicate = ast.MakeArgument(MakeGenSym(constants.LOCAL_PREFIX), CurrentFile, LineNo)
		if condExprs[len(condExprs)-1].IsMethodCall() {
			// we'll change this once we have access to method's types in
			// ProcessMethodCall
			predicate.AddType(constants.TypeNames[constants.TYPE_BOOL])
			condExprs[len(condExprs)-1].Inputs = append(condExprs[len(condExprs)-1].Outputs, condExprs[len(condExprs)-1].Inputs...)
			condExprs[len(condExprs)-1].Outputs = nil
		} else {
			predicate.AddType(constants.TypeNames[condExprs[len(condExprs)-1].Operator.Outputs[0].Type])
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

	skipExpr := ast.MakeExpression(jmpFn, CurrentFile, LineNo)
	skipExpr.Package = pkg

	trueArg := WritePrimary(constants.TYPE_BOOL, encoder.Serialize(true), false)
	skipLines := len(elseExprs)

	skipExpr.AddInput(trueArg[0].Outputs[0])
	skipExpr.ThenLines = skipLines
	skipExpr.ElseLines = 0

	var exprs []*ast.CXExpression
	if condExprs[len(condExprs)-1].Operator != nil || condExprs[len(condExprs)-1].IsMethodCall() {
		exprs = append(exprs, condExprs...)
	}
	exprs = append(exprs, ifExpr)
	exprs = append(exprs, thenExprs...)
	exprs = append(exprs, skipExpr)
	exprs = append(exprs, elseExprs...)

	return exprs
}

// resolveTypeForUnd tries to determine the type that will be returned from an expression
func resolveTypeForUnd(expr *ast.CXExpression) int {
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

// IsTempVar ...
//TODO: Delete this function; only called by next function
func IsTempVar(name string) bool {
	if len(name) >= len(constants.LOCAL_PREFIX) && name[:len(constants.LOCAL_PREFIX)] == constants.LOCAL_PREFIX {
		return true
	}
	return false
}

func OperatorExpression(leftExprs []*ast.CXExpression, rightExprs []*ast.CXExpression, opcode int) (out []*ast.CXExpression) {
	pkg, err := AST.GetCurrentPackage()
	if err != nil {
		panic(err)
	}

	if len(leftExprs[len(leftExprs)-1].Outputs) < 1 {
		name := ast.MakeArgument(MakeGenSym(constants.LOCAL_PREFIX), CurrentFile, LineNo).AddType(constants.TypeNames[resolveTypeForUnd(leftExprs[len(leftExprs)-1])])
		name.Size = leftExprs[len(leftExprs)-1].Operator.Outputs[0].Size
		name.TotalSize = ast.GetSize(leftExprs[len(leftExprs)-1].Operator.Outputs[0])
		name.Type = leftExprs[len(leftExprs)-1].Operator.Outputs[0].Type
		name.ArgDetails.Package = pkg
		name.PreviouslyDeclared = true

		leftExprs[len(leftExprs)-1].Outputs = append(leftExprs[len(leftExprs)-1].Outputs, name)
	}

	if len(rightExprs[len(rightExprs)-1].Outputs) < 1 {
		name := ast.MakeArgument(MakeGenSym(constants.LOCAL_PREFIX), CurrentFile, LineNo).AddType(constants.TypeNames[resolveTypeForUnd(rightExprs[len(rightExprs)-1])])

		name.Size = rightExprs[len(rightExprs)-1].Operator.Outputs[0].Size
		name.TotalSize = ast.GetSize(rightExprs[len(rightExprs)-1].Operator.Outputs[0])
		name.Type = rightExprs[len(rightExprs)-1].Operator.Outputs[0].Type
		name.ArgDetails.Package = pkg
		name.PreviouslyDeclared = true

		rightExprs[len(rightExprs)-1].Outputs = append(rightExprs[len(rightExprs)-1].Outputs, name)
	}

	expr := ast.MakeExpression(ast.Natives[opcode], CurrentFile, LineNo)
	// we can't know the type until we compile the full function
	expr.Package = pkg

	if len(leftExprs[len(leftExprs)-1].Outputs[0].Indexes) > 0 || leftExprs[len(leftExprs)-1].Operator != nil {
		// then it's a function call or an array access
		expr.AddInput(leftExprs[len(leftExprs)-1].Outputs[0])

		if IsTempVar(leftExprs[len(leftExprs)-1].Outputs[0].ArgDetails.Name) {
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

		if IsTempVar(rightExprs[len(rightExprs)-1].Outputs[0].ArgDetails.Name) {
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

func UnaryExpression(op string, prevExprs []*ast.CXExpression) []*ast.CXExpression {
	if len(prevExprs[len(prevExprs)-1].Outputs) == 0 {
		println(ast.CompilationError(CurrentFile, LineNo), "invalid indirection")
		// needs to be stopped immediately
		os.Exit(constants.CX_COMPILATION_ERROR)
	}

	// Some properties need to be read from the base argument
	// due to how we calculate dereferences at the moment.
	baseOut := prevExprs[len(prevExprs)-1].Outputs[0]
	exprOut := ast.GetAssignmentElement(prevExprs[len(prevExprs)-1].Outputs[0])
	switch op {
	case "*":
		exprOut.DereferenceLevels++
		exprOut.DereferenceOperations = append(exprOut.DereferenceOperations, constants.DEREF_POINTER)

		exprOut.DeclarationSpecifiers = append(exprOut.DeclarationSpecifiers, constants.DECL_DEREF)
		exprOut.IsReference = false
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
		if pkg, err := AST.GetCurrentPackage(); err == nil {
			expr := ast.MakeExpression(ast.Natives[constants.OP_BOOL_NOT], CurrentFile, LineNo)
			expr.Package = pkg

			expr.AddInput(exprOut)

			prevExprs[len(prevExprs)-1] = expr
		} else {
			panic(err)
		}
	case "-":
		if pkg, err := AST.GetCurrentPackage(); err == nil {
			expr := ast.MakeExpression(ast.Natives[constants.OP_NEG], CurrentFile, LineNo)
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
func AssociateReturnExpressions(idx int, retExprs []*ast.CXExpression) []*ast.CXExpression {
	var pkg *ast.CXPackage
	var fn *ast.CXFunction
	var err error

	pkg, err = AST.GetCurrentPackage()
	if err != nil {
		panic(err)
	}

	fn, err = pkg.GetCurrentFunction()
	if err != nil {
		panic(err)
	}

	lastExpr := retExprs[len(retExprs)-1]

	outParam := fn.Outputs[idx]

	out := ast.MakeArgument(outParam.ArgDetails.Name, CurrentFile, LineNo)
	out.AddType(constants.TypeNames[outParam.Type])
	out.CustomType = outParam.CustomType
	out.PreviouslyDeclared = true

	if lastExpr.Operator == nil {
		lastExpr.Operator = ast.Natives[constants.OP_IDENTITY]

		lastExpr.Inputs = lastExpr.Outputs
		lastExpr.Outputs = nil
		lastExpr.AddOutput(out)

		return retExprs
	} else if len(lastExpr.Outputs) > 0 {
		expr := ast.MakeExpression(ast.Natives[constants.OP_IDENTITY], CurrentFile, LineNo)
		expr.AddInput(lastExpr.Outputs[0])
		expr.AddOutput(out)

		return append(retExprs, expr)
	} else {
		lastExpr.AddOutput(out)

		return retExprs
	}
}

// AddJmpToReturnExpressions adds an jump expression that makes a function stop its execution
func AddJmpToReturnExpressions(exprs ReturnExpressions) []*ast.CXExpression {
	var pkg *ast.CXPackage
	var fn *ast.CXFunction
	var err error

	pkg, err = AST.GetCurrentPackage()
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

		println(ast.CompilationError(lastExpr.FileName, lastExpr.FileLine), fmt.Sprintf("function '%s' expects to return %d argument%s, but %d output argument%s %s provided", fn.Name, len(fn.Outputs), plural1, exprs.Size, plural2, plural3))
	}

	// expression to jump to the end of the embedding function
	expr := ast.MakeExpression(ast.Natives[constants.OP_GOTO], CurrentFile, LineNo)

	// simulating a label so it gets executed without evaluating a predicate
	expr.Label = MakeGenSym(constants.LABEL_PREFIX)
	expr.ThenLines = constants.MAX_INT32
	expr.Package = pkg

	retExprs = append(retExprs, expr)

	return retExprs
}
