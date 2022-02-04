package generator

import (
	"bytes"
	"encoding/binary"
	"math/rand"
	"strconv"
	"testing"

	"github.com/skycoin/cx/cx/ast"
	"github.com/skycoin/cx/cx/types"

	"github.com/jinzhu/copier"
	cxast "github.com/skycoin/cx/cx/ast"
	"github.com/skycoin/cx/cx/astapi"
	cxconstants "github.com/skycoin/cx/cx/constants"
	cxparsingcompletor "github.com/skycoin/cx/cxparser/cxparsingcompletor"
)

func GenerateRandomExpressions(prgrm *cxast.CXProgram, inputFn *cxast.CXFunction, inputPkg *cxast.CXPackage, fns []*cxast.CXFunction, numExprs int) {
	preExistingExpressions := len(inputFn.Expressions)
	// Checking if we need to add more expressions.
	for i := 0; i < numExprs-preExistingExpressions; i++ {
		op := getRandFn(fns)
		// Last expression output must be the same as function output.
		if i == (numExprs-preExistingExpressions)-1 && len(op.Outputs) > 0 && len(inputFn.Outputs) > 0 {
			for len(op.Outputs) == 0 || op.Outputs[0].Type != inputFn.Outputs[0].Type {
				op = getRandFn(fns)
			}
		}

		exprCXLine := cxast.MakeCXLineExpression(prgrm, "", -1, "")
		expr := cxast.MakeAtomicOperatorExpression(prgrm, op)
		cxAtomicOp, _, _, err := prgrm.GetOperation(expr)
		if err != nil {
			panic(err)
		}

		cxAtomicOp.Package = cxast.CXPackageIndex(inputPkg.Index)
		cxAtomicOp.Function = cxast.CXFunctionIndex(inputFn.Index)
		for c := 0; c < len(op.Inputs); c++ {
			cxAtomicOp.Inputs = append(cxAtomicOp.Inputs, getRandInp(prgrm, expr))
		}

		// if operator is jmp, add then and else lines
		if IsJumpOperator(op.AtomicOPCode) {
			lineNumOptions := numExprs - len(inputFn.Expressions)
			if lineNumOptions < 0 {
				lineNumOptions = (lineNumOptions * -1) - 2
			}
			randThenLineIndex := 0
			if lineNumOptions > 0 {
				randThenLineIndex = rand.Intn(lineNumOptions)
			}

			cxAtomicOp.ThenLines = 1
			cxAtomicOp.ElseLines = randThenLineIndex
		}

		// We need to add the expression at this point, so we
		// can consider this expression's output as a
		// possibility to assign stuff.
		inputFn.Expressions = append(inputFn.Expressions, exprCXLine, expr)

		// Adding last expression, so output must be fn's output.
		if i == numExprs-preExistingExpressions-1 {
			cxAtomicOp.Outputs = append(cxAtomicOp.Outputs, inputFn.Outputs[0])
		} else {
			for c := 0; c < len(op.Outputs); c++ {
				cxAtomicOp.Outputs = append(cxAtomicOp.Outputs, getRandOut(prgrm, expr))
			}
		}
	}
	inputFn.Size = calcFnSize(prgrm, inputFn)
	inputFn.LineCount = numExprs
}

func IsJumpOperator(opCode int) bool {
	switch opCode {
	case cxconstants.OP_JMP,
		cxconstants.OP_ABS_JMP,
		cxconstants.OP_JMP_EQ,
		cxconstants.OP_JMP_UNEQ,
		cxconstants.OP_JMP_GT,
		cxconstants.OP_JMP_GTEQ,
		cxconstants.OP_JMP_LT,
		cxconstants.OP_JMP_LTEQ,
		cxconstants.OP_JMP_ZERO,
		cxconstants.OP_JMP_NOT_ZERO:
		return true
	default:
		return false
	}
}

func getRandFn(fnSet []*cxast.CXFunction) *cxast.CXFunction {
	return fnSet[rand.Intn(len(fnSet))]
}

func calcFnSize(prgrm *cxast.CXProgram, fn *cxast.CXFunction) (size types.Pointer) {
	for _, arg := range fn.Inputs {
		size += arg.TotalSize
	}
	for _, arg := range fn.Outputs {
		size += arg.TotalSize
	}
	for _, expr := range fn.Expressions {
		if expr.Type == ast.CX_LINE {
			continue
		}
		cxAtomicOp, _, _, err := prgrm.GetOperation(expr)
		if err != nil {
			panic(err)
		}

		// TODO: We're only considering one output per operator.
		/// Not because of practicality, but because multiple returns in CX are currently buggy anyway.
		if len(cxAtomicOp.Operator.Outputs) > 0 {
			size += cxAtomicOp.Operator.Outputs[0].TotalSize
		}
	}

	return size
}

func getRandInp(prgrm *cxast.CXProgram, expr *cxast.CXExpression) *cxast.CXArgument {
	var rndExprIdx int
	var argToCopy *cxast.CXArgument
	var arg cxast.CXArgument

	cxAtomicOp, _, _, err := prgrm.GetOperation(expr)
	if err != nil {
		panic(err)
	}

	cxAtomicOpFunction, err := prgrm.GetFunctionFromArray(cxAtomicOp.Function)
	if err != nil {
		panic(err)
	}

	// Find available arg options.
	optionsFromInputs, optionsFromExpressions := findArgOptions(prgrm, expr, cxAtomicOp.Operator.Inputs[0].Type)
	lengthOfOptions := len(optionsFromInputs) + len(optionsFromExpressions)

	// if no options available or if operator is jump, add new i32_LT expression.
	if lengthOfOptions == 0 || cxAtomicOp.Operator.AtomicOPCode == cxconstants.OP_JMP {
		// TODO: improve process when there's OP_JMP
		return addNewExpression(prgrm, expr, cxast.OpCodes["i32.lt"])
	}

	rndExprIdx = rand.Intn(lengthOfOptions)
	gotOptionsFromFunctionInputs := rndExprIdx < len(optionsFromInputs)

	if gotOptionsFromFunctionInputs {
		argToCopy = cxAtomicOpFunction.Inputs[optionsFromInputs[rndExprIdx]]
	} else {
		rndExprIdx -= len(optionsFromInputs)
		cxAtomicOp1, err := prgrm.GetCXAtomicOpFromExpressions(cxAtomicOpFunction.Expressions, optionsFromExpressions[rndExprIdx])
		if err != nil {
			panic(err)
		}

		argToCopy = cxAtomicOp1.Operator.Outputs[0]
	}

	// Making a copy of the argument
	err = copier.Copy(&arg, argToCopy)
	if err != nil {
		panic(err)
	}

	if !gotOptionsFromFunctionInputs {
		determineExpressionOffset(prgrm, &arg, expr, optionsFromExpressions[rndExprIdx])
		arg.Name = strconv.Itoa(optionsFromExpressions[rndExprIdx])
	}
	arg.Package = cxAtomicOpFunction.Package

	return &arg
}

func addNewExpression(prgrm *cxast.CXProgram, expr *cxast.CXExpression, expressionType int) *cxast.CXArgument {
	var rndExprIdx int
	var argToAdd *cxast.CXArgument

	exprCXAtomicOp, _, _, err := prgrm.GetOperation(expr)
	if err != nil {
		panic(err)
	}

	exprCXAtomicOpFunction, err := prgrm.GetFunctionFromArray(exprCXAtomicOp.Function)
	if err != nil {
		panic(err)
	}

	newExprCXLine := cxast.MakeCXLineExpression(prgrm, "", -1, "")
	newExpr := cxast.MakeAtomicOperatorExpression(prgrm, cxast.Natives[expressionType])
	newExprCXAtomicOp, _, _, err := prgrm.GetOperation(newExpr)
	if err != nil {
		panic(err)
	}

	newExprCXAtomicOp.Operator.Name = cxast.OpNames[expressionType]

	// Add expression's inputs
	for i := 0; i < 2; i++ {
		optionsFromInputs, optionsFromExpressions := findArgOptions(prgrm, expr, newExprCXAtomicOp.Operator.Inputs[0].Type)
		rndExprIdx = rand.Intn(len(optionsFromInputs) + len(optionsFromExpressions))
		if rndExprIdx < len(optionsFromInputs) {
			argToAdd = exprCXAtomicOpFunction.Inputs[optionsFromInputs[rndExprIdx]]
		} else {
			rndExprIdx -= len(optionsFromInputs)
			argToAddCXAtomicOp, err := prgrm.GetCXAtomicOpFromExpressions(exprCXAtomicOpFunction.Expressions, optionsFromExpressions[rndExprIdx])
			if err != nil {
				panic(err)
			}
			argToAdd = argToAddCXAtomicOp.Outputs[0]
		}
		newExprCXAtomicOp.AddInput(argToAdd)
	}

	// Add expression's output
	argOutName := strconv.Itoa(len(exprCXAtomicOpFunction.Expressions))
	argOut := cxast.MakeField(argOutName, types.BOOL, "", -1)
	argOut.AddType(types.Code(types.BOOL))
	argOut.Package = exprCXAtomicOpFunction.Package
	newExprCXAtomicOp.AddOutput(argOut)
	exprCXAtomicOpFunction.AddExpression(prgrm, newExprCXLine)
	exprCXAtomicOpFunction.AddExpression(prgrm, newExpr)

	determineExpressionOffset(prgrm, argOut, expr, len(exprCXAtomicOpFunction.Expressions))

	return argOut
}

func findArgOptions(prgrm *cxast.CXProgram, expr *cxast.CXExpression, argTypeToFind types.Code) ([]int, []int) {
	var optionsFromInputs []int
	var optionsFromExpressions []int

	cxAtomicOp, _, _, err := prgrm.GetOperation(expr)
	if err != nil {
		panic(err)
	}

	cxAtomicOpFunction, err := prgrm.GetFunctionFromArray(cxAtomicOp.Function)
	if err != nil {
		panic(err)
	}
	// loop in inputs
	for i, inp := range cxAtomicOpFunction.Inputs {
		if inp.Type == argTypeToFind && inp.Name != "" {
			// add index to options from inputs
			optionsFromInputs = append(optionsFromInputs, i)
		}
	}

	// loop in expression outputs
	for i, exp := range cxAtomicOpFunction.Expressions {
		expCXAtomicOp, _, _, err := prgrm.GetOperation(exp)
		if err != nil {
			panic(err)
		}

		if len(expCXAtomicOp.Outputs) > 0 && expCXAtomicOp.Outputs[0].Type == argTypeToFind && expCXAtomicOp.Outputs[0].Name != "" {
			// add index to options from inputs
			optionsFromExpressions = append(optionsFromExpressions, i)
		}
	}
	return optionsFromInputs, optionsFromExpressions
}

func getRandOut(prgrm *cxast.CXProgram, expr *cxast.CXExpression) *cxast.CXArgument {
	var arg cxast.CXArgument
	var optionsFromExpressions []int

	cxAtomicOp, _, _, err := prgrm.GetOperation(expr)
	if err != nil {
		panic(err)
	}

	cxAtomicOpFunction, err := prgrm.GetFunctionFromArray(cxAtomicOp.Function)
	if err != nil {
		panic(err)
	}

	for i, exp := range cxAtomicOpFunction.Expressions {
		if exp.Type == ast.CX_LINE {
			continue
		}
		expCXAtomicOp, _, _, err := prgrm.GetOperation(exp)
		if err != nil {
			panic(err)
		}

		if len(expCXAtomicOp.Operator.Outputs) > 0 && expCXAtomicOp.Operator.Outputs[0].Type == cxAtomicOp.Operator.Outputs[0].Type {
			optionsFromExpressions = append(optionsFromExpressions, i)
		}
	}

	rndExprIdx := rand.Intn(len(optionsFromExpressions))

	copyCXAtomicOp, err := prgrm.GetCXAtomicOpFromExpressions(cxAtomicOpFunction.Expressions, optionsFromExpressions[rndExprIdx])
	if err != nil {
		panic(err)
	}

	// Making a copy of the argument
	err = copier.Copy(&arg, copyCXAtomicOp.Operator.Outputs[0])
	if err != nil {
		panic(err)
	}

	determineExpressionOffset(prgrm, &arg, expr, optionsFromExpressions[rndExprIdx])
	arg.Package = cxAtomicOpFunction.Package
	arg.Name = strconv.Itoa(optionsFromExpressions[rndExprIdx])

	return &arg
}

func determineExpressionOffset(prgrm *cxast.CXProgram, arg *cxast.CXArgument, expr *cxast.CXExpression, indexOfSelectedOption int) {
	cxAtomicOp, _, _, err := prgrm.GetOperation(expr)
	if err != nil {
		panic(err)
	}

	cxAtomicOpFunction, err := prgrm.GetFunctionFromArray(cxAtomicOp.Function)
	if err != nil {
		panic(err)
	}

	// Determining the offset where the expression should be writing to.
	for c := 0; c < len(cxAtomicOpFunction.Inputs); c++ {
		arg.Offset += cxAtomicOpFunction.Inputs[c].TotalSize
	}
	for c := 0; c < len(cxAtomicOpFunction.Outputs); c++ {
		arg.Offset += cxAtomicOpFunction.Outputs[c].TotalSize
	}
	for c := 0; c < indexOfSelectedOption; c++ {
		cxAtomicOp1, err := prgrm.GetCXAtomicOpFromExpressions(cxAtomicOpFunction.Expressions, c)
		if err != nil {
			panic(err)
		}

		if len(cxAtomicOp1.Operator.Outputs) > 0 {
			// TODO: We're only considering one output per operator.
			/// Not because of practicality, but because multiple returns in CX are currently buggy anyway.
			arg.Offset += cxAtomicOp1.Operator.Outputs[0].TotalSize
		}
	}
}

func GetFunctionSet(fnNames []string) (fns []*cxast.CXFunction) {
	for _, fnName := range fnNames {
		fn := cxast.Natives[cxast.OpCodes[fnName]]
		if fn == nil {
			panic("standard library function not found.")
		}

		fns = append(fns, fn)
	}
	return fns
}

func GenerateSampleProgram(t *testing.T, withLiteral bool) *cxast.CXProgram {
	var cxProgram *cxast.CXProgram

	// Needed for AddNativeExpressionToFunction
	// because of dependency on cxast.OpNames
	cxparsingcompletor.InitCXCore()
	cxProgram = cxast.MakeProgram()

	err := astapi.AddEmptyPackage(cxProgram, "main")
	if err != nil {
		t.Errorf("want no error, got %v", err)
	}

	err = astapi.AddEmptyFunctionToPackage(cxProgram, "main", "TestFunction")
	if err != nil {
		t.Errorf("want no error, got %v", err)
	}

	err = astapi.AddNativeInputToFunction(cxProgram, "main", "TestFunction", "inputOne", types.Code(types.I32))
	if err != nil {
		t.Errorf("want no error, got %v", err)
	}

	err = astapi.AddNativeOutputToFunction(cxProgram, "main", "TestFunction", "outputOne", types.Code(types.I32))
	if err != nil {
		t.Errorf("want no error, got %v", err)
	}
	functionSetNames := []string{"i32.add", "i32.mul", "i32.sub", "i32.eq", "i32.uneq", "i32.gt", "i32.gteq", "i32.lt", "i32.lteq", "bool.not", "bool.or", "bool.and", "bool.uneq", "bool.eq", "i32.neg", "i32.abs", "i32.bitand", "i32.bitor", "i32.bitxor", "i32.bitclear", "i32.bitshl", "i32.bitshr", "i32.max", "i32.min", "i32.rand"}
	fns := GetFunctionSet(functionSetNames)

	fn, _ := cxProgram.GetFunction("TestFunction", "main")
	pkg, _ := cxProgram.GetPackage("main")
	GenerateRandomExpressions(cxProgram, fn, pkg, fns, 30)

	if withLiteral {
		buf := new(bytes.Buffer)
		var num int32 = 5
		binary.Write(buf, binary.LittleEndian, num)
		err = astapi.AddLiteralInputToExpression(cxProgram, "main", "TestFunction", buf.Bytes(), types.I32, 2)
		if err != nil {
			t.Errorf("want no error, got %v", err)
		}

	}

	cxProgram.PrintProgram()
	return cxProgram
}
