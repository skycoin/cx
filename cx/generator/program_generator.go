package generator

import (
	"bytes"
	"encoding/binary"
	"math/rand"
	"strconv"
	"testing"

	"github.com/skycoin/cx/cx/types"

	"github.com/jinzhu/copier"
	cxast "github.com/skycoin/cx/cx/ast"
	"github.com/skycoin/cx/cx/astapi"
	cxconstants "github.com/skycoin/cx/cx/constants"
	cxinit "github.com/skycoin/cx/cx/init"
	cxparsingcompletor "github.com/skycoin/cx/cxparser/cxparsingcompletor"
)

func GenerateRandomExpressions(prgrm *cxast.CXProgram, inputFn *cxast.CXFunction, inputPkg *cxast.CXPackage, fns []*cxast.CXNativeFunction, numExprs int) {
	preExistingExpressions := len(inputFn.Expressions)
	// Checking if we need to add more expressions.
	for i := 0; i < numExprs-preExistingExpressions; i++ {
		op := getRandFn(fns)

		inputFnOutputs := inputFn.GetOutputs(prgrm)
		// Last expression output must be the same as function output.
		if i == (numExprs-preExistingExpressions)-1 && len(op.Outputs) > 0 && len(inputFnOutputs) > 0 {
			for len(op.Outputs) == 0 || op.Outputs[0].Type != prgrm.GetCXArgFromArray(inputFnOutputs[0]).Type {
				op = getRandFn(fns)
			}
		}

		exprCXLine := cxast.MakeCXLineExpression(prgrm, "", -1, "")
		expr := cxast.MakeAtomicOperatorExpression(prgrm, op)
		cxAtomicOp, err := prgrm.GetCXAtomicOp(expr.Index)
		if err != nil {
			panic(err)
		}

		cxAtomicOp.Package = cxast.CXPackageIndex(inputPkg.Index)
		cxAtomicOp.Function = cxast.CXFunctionIndex(inputFn.Index)
		for c := 0; c < len(op.Inputs); c++ {
			inpIdx := prgrm.AddCXArgInArray(getRandInp(prgrm, expr))
			cxAtomicOp.Inputs = append(cxAtomicOp.Inputs, inpIdx)
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
		inputFn.Expressions = append(inputFn.Expressions, *exprCXLine, *expr)
		prgrm.CXFunctions[inputFn.Index] = *inputFn
		// Adding last expression, so output must be fn's output.
		if i == numExprs-preExistingExpressions-1 {
			cxAtomicOp.Outputs = append(cxAtomicOp.Outputs, inputFnOutputs[0])
		} else {
			for c := 0; c < len(op.Outputs); c++ {
				outIdx := prgrm.AddCXArgInArray(getRandOut(prgrm, expr))
				cxAtomicOp.Outputs = append(cxAtomicOp.Outputs, outIdx)
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

func getRandFn(fnSet []*cxast.CXNativeFunction) *cxast.CXNativeFunction {
	return fnSet[rand.Intn(len(fnSet))]
}

func calcFnSize(prgrm *cxast.CXProgram, fn *cxast.CXFunction) (size types.Pointer) {
	for _, argIdx := range fn.GetInputs(prgrm) {
		arg := prgrm.GetCXArgFromArray(argIdx)
		size += arg.TotalSize
	}

	fnOutputs := fn.GetOutputs(prgrm)
	for _, argIdx := range fnOutputs {
		arg := prgrm.GetCXArgFromArray(argIdx)
		size += arg.TotalSize
	}
	for _, expr := range fn.Expressions {
		if expr.Type == cxast.CX_LINE {
			continue
		}
		cxAtomicOp, err := prgrm.GetCXAtomicOp(expr.Index)
		if err != nil {
			panic(err)
		}

		cxAtomicOpOperator := prgrm.GetFunctionFromArray(cxAtomicOp.Operator)
		cxAtomicOpOperatorOutputs := cxAtomicOpOperator.GetOutputs(prgrm)
		// TODO: We're only considering one output per operator.
		/// Not because of practicality, but because multiple returns in CX are currently buggy anyway.
		if len(cxAtomicOpOperatorOutputs) > 0 {
			size += prgrm.GetCXArgFromArray(cxAtomicOpOperatorOutputs[0]).TotalSize
		}
	}

	return size
}

func getRandInp(prgrm *cxast.CXProgram, expr *cxast.CXExpression) *cxast.CXArgument {
	var rndExprIdx int
	var argToCopy *cxast.CXArgument
	var arg cxast.CXArgument

	cxAtomicOp, err := prgrm.GetCXAtomicOp(expr.Index)
	if err != nil {
		panic(err)
	}

	cxAtomicOpOperator := prgrm.GetFunctionFromArray(cxAtomicOp.Operator)
	cxAtomicOpOperatorInputs := cxAtomicOpOperator.GetInputs(prgrm)
	cxAtomicOpFunction := prgrm.GetFunctionFromArray(cxAtomicOp.Function)
	cxAtomicOpFunctionInputs := cxAtomicOpFunction.GetInputs(prgrm)
	// Find available arg options.
	optionsFromInputs, optionsFromExpressions := findArgOptions(prgrm, expr, prgrm.GetCXArgFromArray(cxAtomicOpOperatorInputs[0]).Type)
	lengthOfOptions := len(optionsFromInputs) + len(optionsFromExpressions)

	// if no options available or if operator is jump, add new i32_LT expression.
	if lengthOfOptions == 0 || cxAtomicOpOperator.AtomicOPCode == cxconstants.OP_JMP {
		// TODO: improve process when there's OP_JMP
		return addNewExpression(prgrm, expr, cxast.OpCodes["i32.lt"])
	}

	rndExprIdx = rand.Intn(lengthOfOptions)
	gotOptionsFromFunctionInputs := rndExprIdx < len(optionsFromInputs)

	if gotOptionsFromFunctionInputs {
		argToCopy = prgrm.GetCXArgFromArray(cxAtomicOpFunctionInputs[optionsFromInputs[rndExprIdx]])
	} else {
		rndExprIdx -= len(optionsFromInputs)
		cxAtomicOp1, err := prgrm.GetCXAtomicOpFromExpressions(cxAtomicOpFunction.Expressions, optionsFromExpressions[rndExprIdx])
		if err != nil {
			panic(err)
		}
		cxAtomicOp1Operator := prgrm.GetFunctionFromArray(cxAtomicOp1.Operator)
		cxAtomicOp1OperatorOutputs := cxAtomicOp1Operator.GetOutputs(prgrm)
		argToCopy = prgrm.GetCXArgFromArray(cxAtomicOp1OperatorOutputs[0])
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
	var argToAdd cxast.CXArgumentIndex

	exprCXAtomicOp, err := prgrm.GetCXAtomicOp(expr.Index)
	if err != nil {
		panic(err)
	}

	exprCXAtomicOpFunction := prgrm.GetFunctionFromArray(exprCXAtomicOp.Function)
	exprCXAtomicOpFunctionInputs := exprCXAtomicOpFunction.GetInputs(prgrm)

	newExprCXLine := cxast.MakeCXLineExpression(prgrm, "", -1, "")
	newExpr := cxast.MakeAtomicOperatorExpression(prgrm, cxast.Natives[expressionType])
	newExprCXAtomicOp, _, _, err := prgrm.GetOperation(newExpr)
	if err != nil {
		panic(err)
	}
	newExprCXAtomicOpOperator := prgrm.GetFunctionFromArray(newExprCXAtomicOp.Operator)
	newExprCXAtomicOpOperatorInputs := newExprCXAtomicOpOperator.GetInputs(prgrm)
	newExprCXAtomicOpOperator.Name = cxast.OpNames[expressionType]

	// Add expression's inputs
	for i := 0; i < 2; i++ {
		optionsFromInputs, optionsFromExpressions := findArgOptions(prgrm, expr, prgrm.GetCXArgFromArray(newExprCXAtomicOpOperatorInputs[0]).Type)
		rndExprIdx = rand.Intn(len(optionsFromInputs) + len(optionsFromExpressions))
		if rndExprIdx < len(optionsFromInputs) {
			argToAdd = exprCXAtomicOpFunctionInputs[optionsFromInputs[rndExprIdx]]
		} else {
			rndExprIdx -= len(optionsFromInputs)
			argToAddCXAtomicOp, err := prgrm.GetCXAtomicOpFromExpressions(exprCXAtomicOpFunction.Expressions, optionsFromExpressions[rndExprIdx])
			if err != nil {
				panic(err)
			}
			argToAdd = argToAddCXAtomicOp.Outputs[0]
		}
		newExprCXAtomicOp.AddInput(prgrm, argToAdd)
	}

	// Add expression's output
	argOutName := strconv.Itoa(len(exprCXAtomicOpFunction.Expressions))
	argOut := cxast.MakeField(argOutName, types.BOOL, "", -1)
	argOut.SetType(types.Code(types.BOOL))
	argOut.Package = exprCXAtomicOpFunction.Package
	argOutIdx := prgrm.AddCXArgInArray(argOut)
	newExprCXAtomicOp.AddOutput(prgrm, argOutIdx)
	exprCXAtomicOpFunction.AddExpression(prgrm, newExprCXLine)
	exprCXAtomicOpFunction.AddExpression(prgrm, newExpr)

	prgrm.CXFunctions[exprCXAtomicOpFunction.Index] = *exprCXAtomicOpFunction

	determineExpressionOffset(prgrm, argOut, expr, len(exprCXAtomicOpFunction.Expressions))

	return argOut
}

func findArgOptions(prgrm *cxast.CXProgram, expr *cxast.CXExpression, argTypeToFind types.Code) ([]int, []int) {
	var optionsFromInputs []int
	var optionsFromExpressions []int

	cxAtomicOp, err := prgrm.GetCXAtomicOp(expr.Index)
	if err != nil {
		panic(err)
	}

	cxAtomicOpFunction := prgrm.GetFunctionFromArray(cxAtomicOp.Function)
	cxAtomicOpFunctionInputs := cxAtomicOpFunction.GetInputs(prgrm)
	// loop in inputs
	for i, inpIdx := range cxAtomicOpFunctionInputs {
		inp := prgrm.GetCXArgFromArray(inpIdx)
		if inp.Type == argTypeToFind && inp.Name != "" {
			// add index to options from inputs
			optionsFromInputs = append(optionsFromInputs, i)
		}
	}

	// loop in expression outputs
	for i, exp := range cxAtomicOpFunction.Expressions {
		expCXAtomicOp, _, _, err := prgrm.GetOperation(&exp)
		if err != nil {
			panic(err)
		}

		if len(expCXAtomicOp.Outputs) > 0 && prgrm.GetCXArgFromArray(expCXAtomicOp.Outputs[0]).Type == argTypeToFind && prgrm.GetCXArgFromArray(expCXAtomicOp.Outputs[0]).Name != "" {
			// add index to options from inputs
			optionsFromExpressions = append(optionsFromExpressions, i)
		}
	}
	return optionsFromInputs, optionsFromExpressions
}

func getRandOut(prgrm *cxast.CXProgram, expr *cxast.CXExpression) *cxast.CXArgument {
	var arg cxast.CXArgument
	var optionsFromExpressions []int

	cxAtomicOp, err := prgrm.GetCXAtomicOp(expr.Index)
	if err != nil {
		panic(err)
	}
	cxAtomicOpOperator := prgrm.GetFunctionFromArray(cxAtomicOp.Operator)
	cxAtomicOpOperatorOutputs := cxAtomicOpOperator.GetOutputs(prgrm)
	cxAtomicOpFunction := prgrm.GetFunctionFromArray(cxAtomicOp.Function)

	for i, exp := range cxAtomicOpFunction.Expressions {
		if exp.Type == cxast.CX_LINE {
			continue
		}
		expCXAtomicOp, _, _, err := prgrm.GetOperation(&exp)
		if err != nil {
			panic(err)
		}
		expCXAtomicOpOperator := prgrm.GetFunctionFromArray(expCXAtomicOp.Operator)
		expCXAtomicOpOperatorOutputs := expCXAtomicOpOperator.GetOutputs(prgrm)
		if len(expCXAtomicOpOperatorOutputs) > 0 && prgrm.GetCXArgFromArray(expCXAtomicOpOperatorOutputs[0]).Type == prgrm.GetCXArgFromArray(cxAtomicOpOperatorOutputs[0]).Type {
			optionsFromExpressions = append(optionsFromExpressions, i)
		}
	}

	rndExprIdx := rand.Intn(len(optionsFromExpressions))

	copyCXAtomicOp, err := prgrm.GetCXAtomicOpFromExpressions(cxAtomicOpFunction.Expressions, optionsFromExpressions[rndExprIdx])
	if err != nil {
		panic(err)
	}
	copyCXAtomicOpOperator := prgrm.GetFunctionFromArray(copyCXAtomicOp.Operator)
	copyCXAtomicOpOperatorOutputs := copyCXAtomicOpOperator.GetOutputs(prgrm)
	// Making a copy of the argument
	err = copier.Copy(&arg, copyCXAtomicOpOperatorOutputs[0])
	if err != nil {
		panic(err)
	}

	determineExpressionOffset(prgrm, &arg, expr, optionsFromExpressions[rndExprIdx])
	arg.Package = cxAtomicOpFunction.Package
	arg.Name = strconv.Itoa(optionsFromExpressions[rndExprIdx])

	return &arg
}

func determineExpressionOffset(prgrm *cxast.CXProgram, arg *cxast.CXArgument, expr *cxast.CXExpression, indexOfSelectedOption int) {
	cxAtomicOp, err := prgrm.GetCXAtomicOp(expr.Index)
	if err != nil {
		panic(err)
	}

	cxAtomicOpFunction := prgrm.GetFunctionFromArray(cxAtomicOp.Function)
	cxAtomicOpFunctionInputs := cxAtomicOpFunction.GetInputs(prgrm)
	cxAtomicOpFunctionOutputs := cxAtomicOpFunction.GetOutputs(prgrm)
	// Determining the offset where the expression should be writing to.
	for c := 0; c < len(cxAtomicOpFunctionInputs); c++ {
		arg.Offset += prgrm.GetCXArgFromArray(cxAtomicOpFunctionInputs[c]).TotalSize
	}
	for c := 0; c < len(cxAtomicOpFunctionOutputs); c++ {
		arg.Offset += prgrm.GetCXArgFromArray(cxAtomicOpFunctionOutputs[c]).TotalSize
	}
	for c := 0; c < indexOfSelectedOption; c++ {
		cxAtomicOp1, err := prgrm.GetCXAtomicOpFromExpressions(cxAtomicOpFunction.Expressions, c)
		if err != nil {
			panic(err)
		}
		cxAtomicOp1Operator := prgrm.GetFunctionFromArray(cxAtomicOp1.Operator)
		cxAtomicOp1OperatorOutputs := cxAtomicOp1Operator.GetOutputs(prgrm)
		if len(cxAtomicOp1OperatorOutputs) > 0 {
			// TODO: We're only considering one output per operator.
			/// Not because of practicality, but because multiple returns in CX are currently buggy anyway.
			arg.Offset += prgrm.GetCXArgFromArray(cxAtomicOp1OperatorOutputs[0]).TotalSize
		}
	}
}

func GetFunctionSet(fnNames []string) (fns []*cxast.CXNativeFunction) {
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
	cxProgram = cxinit.MakeProgram()

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
