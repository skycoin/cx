package cxcore

import (
	"fmt"
	"strconv"
	"math"
	"math/rand"
	"github.com/jinzhu/copier"

	// "github.com/SkycoinProject/skycoin/src/cipher/encoder"
)

// adaptMainFn removes the main function from the main
// package. Then it creates a new main function that will contain a call
// to the solution function.
func (prgrm *CXProgram) adaptSolution(solution []string) {
	solutionName := solution[len(solution)-1]
	// Ensuring that main pkg exists.
	var mainPkg *CXPackage
	mainPkg, err := prgrm.GetPackage(MAIN_PKG)
	if err != nil {
		panic(err)
	}
	
	// mainFn, err := mainPkg.GetFunction(MAIN_FUNC)
	// if err != nil {
	// 	panic(err)
	// }

	mainFn := MakeFunction(MAIN_FUNC, "", -1)
	mainFn.Package = mainPkg
	for i, fn := range mainPkg.Functions {
		if fn.Name == MAIN_FUNC {
			mainPkg.Functions[i] = mainFn
			break
		}
	}

	mainFn.Expressions = nil
	mainFn.Inputs = nil
	mainFn.Outputs = nil
	
	// idx := -1
	// for i, fn := range pkg.Functions {
	// 	if fn.Name == MAIN_FUNC {
	// 		idx = i
	// 		break
	// 	}
	// }
	// _ = idx
	// // Removing main function.
	// pkg.Functions = append(pkg.Functions[:idx], pkg.Functions[idx+1:]...)

	// var solPkg *CXPackage
	// solPkg, err = prgrm.GetCurrentPackage()
	// if err != nil {
	// 	panic(err)
	// }
	// pkg, err = prgrm.GetPackage(MAIN_PKG)
	// if err != nil {
	// 	panic(err)
	// }
	// mainFn := MakeFunction(MAIN_FUNC, "", -1)
	// pkg.AddFunction(mainFn)

	

	mainInp := MakeArgument("inp", "", -1).AddType("f64")
	mainInp.Package = mainPkg
	mainOut := MakeArgument("out", "", -1).AddType("f64")
	mainOut.Package = mainPkg
	mainOut.Offset += mainInp.TotalSize
	mainFn.AddInput(mainInp)
	mainFn.AddOutput(mainOut)

	var sol *CXFunction
	sol, err = mainPkg.GetFunction(solutionName)
	if err != nil {
		panic(err)
	}

	// // We need to replace the solution function, as it is a pointer.
	// /// It'd maintain a reference to it.
	// newSol := MakeFunction(solutionName, sol.FileName, sol.FileLine)
	// // Inputs and outputs can be the same. We only need their offsets.
	// newSol.Inputs = sol.Inputs
	// newSol.Outputs = sol.Outputs
	// newSol.Package = solPkg

	// // We'll need to replace the pointer, not only the object.
	/// So we need the solution index to replace it from solPkg.Functions.
	// solFnIdx := -1
	// for i, fn := range solPkg.Functions {
	// 	if fn.Name == solutionName {
	// 		solFnIdx = i
	// 		break
	// 	}
	// }
	// solPkg.Functions[solFnIdx] = newSol

	expr := MakeExpression(sol, "", -1)
	expr.Package = mainPkg
	expr.AddOutput(mainOut)
	expr.AddInput(mainInp)

	// prnt := MakeExpression(Natives[OpCodes["f64.print"]], "", -1)
	// prnt.Package = mainPkg
	// prnt.AddInput(mainOut)
	
	
	mainFn.AddExpression(expr)
	// mainFn.AddExpression(prnt)
	mainFn.Length = 1
	// TODO: Assuming one input and one output of f64 type.
	mainFn.Size = calcFnSize(mainFn)
}

func getFnBag(prgrm *CXProgram, fnBag []string) (fns []*CXFunction) {
	pkgName := ""
	fnName := ""
	for i, name := range fnBag {
		if name == "pkg" {
			pkgName = fnBag[i+1]
		}
		if name == "fn" {
			fnName = fnBag[i+1]
		}
		if pkgName != "" && fnName != "" {
			// Then it's a standard library function, like i32.add.
			var fn *CXFunction
			if pkgName == STDLIB_PKG {
				fn = Natives[OpCodes[fnName]]
				if fn == nil {
					panic("standard library function not found.")
				}
			} else {
				var err error
				fn, err = prgrm.GetFunction(fnName, pkgName)
				if err != nil {
					panic(err)
				}
			}
			
			fns = append(fns, fn)
			pkgName = ""
			fnName = ""
		}
	}
	return fns
}

func getRandFn(fnBag []*CXFunction) *CXFunction {
	return fnBag[rand.Intn(len(fnBag))]
}

func getFnArgs(fn *CXFunction) (args []*CXArgument) {
	for _, arg := range fn.Inputs {
		args = append(args, arg)
	}

	for _, arg := range fn.Outputs {
		args = append(args, arg)
	}

	for _, expr := range fn.Expressions {
		for _, arg := range expr.Inputs {
			args = append(args, arg)
		}

		for _, arg := range expr.Outputs {
			args = append(args, arg)
		}
	}

	return args
}

func calcFnSize(fn *CXFunction) (size int) {
	for _, arg := range fn.Inputs {
		size += arg.TotalSize
	}
	for _, arg := range fn.Outputs {
		size += arg.TotalSize
	}
	for _, expr := range fn.Expressions {
		// TODO: We're only considering one output per operator.
		/// Not because of practicality, but because multiple returns in CX are currently buggy anyway.
		if len(expr.Operator.Outputs) > 0 {
			size += expr.Operator.Outputs[0].TotalSize
		}
	}

	return size
}

func getRandInp(fn *CXFunction) *CXArgument {
	var arg CXArgument
	// Unlike getRandOut, we need to also consider the function inputs.
	rndExprIdx := rand.Intn(len(fn.Inputs) + len(fn.Expressions))
	// Then we're returning one of fn.Inputs as the input argument.
	if rndExprIdx < len(fn.Inputs) {
		// Making a copy of the operator.
		// Inputs should have already a compiled offset.
		err := copier.Copy(&arg, fn.Inputs[rndExprIdx])
		if err != nil {
			panic(err)
		}
		arg.Package = fn.Package
		return &arg
	}
	// It was not a function input.
	// We need to subtract the number of inputs to rndExprIdx.
	rndExprIdx -= len(fn.Inputs)
	// Making a copy of the argument
	err := copier.Copy(&arg, fn.Expressions[rndExprIdx].Operator.Outputs[0])
	if err != nil {
		panic(err)
	}
	// Determining the offset where the expression should be writing to.
	for c := 0; c < len(fn.Inputs); c++ {
		arg.Offset += fn.Inputs[c].TotalSize
	}
	for c := 0; c < len(fn.Outputs); c++ {
		arg.Offset += fn.Outputs[c].TotalSize
	}
	for c := 0; c < rndExprIdx; c++ {
		// TODO: We're only considering one output per operator.
		/// Not because of practicality, but because multiple returns in CX are currently buggy anyway.
		arg.Offset += fn.Expressions[c].Operator.Outputs[0].TotalSize
	}

	arg.Package = fn.Package
	arg.Name = strconv.Itoa(rndExprIdx)
	return &arg
}

func getRandOut(fn *CXFunction) *CXArgument {
	var arg CXArgument
	rndExprIdx := rand.Intn(len(fn.Expressions))
	// Making a copy of the argument
	err := copier.Copy(&arg, fn.Expressions[rndExprIdx].Operator.Outputs[0])
	if err != nil {
		panic(err)
	}
	// Determining the offset where the expression should be writing to.
	for c := 0; c < len(fn.Inputs); c++ {
		arg.Offset += fn.Inputs[c].TotalSize
	}
	for c := 0; c < len(fn.Outputs); c++ {
		arg.Offset += fn.Outputs[c].TotalSize
	}
	for c := 0; c < rndExprIdx; c++ {
		// TODO: We're only considering one output per operator.
		/// Not because of practicality, but because multiple returns in CX are currently buggy anyway.
		arg.Offset += fn.Expressions[c].Operator.Outputs[0].TotalSize
	}

	arg.Package = fn.Package
	arg.Name = strconv.Itoa(rndExprIdx)
	return &arg
}

func (fn *CXFunction) mutateFn(fnBag []*CXFunction) {
	rndExprIdx := rand.Intn(len(fn.Expressions))
	rndFn := getRandFn(fnBag)

	expr := MakeExpression(rndFn, "", -1)
	expr.Package = fn.Package
	expr.Inputs = fn.Expressions[rndExprIdx].Inputs
	expr.Outputs = fn.Expressions[rndExprIdx].Outputs

	exprs := make([]*CXExpression, len(fn.Expressions))
	for i, ex := range fn.Expressions {
		if i == rndExprIdx {
			exprs[i] = expr
		} else {
			exprs[i] = ex
		}
	}

	// fn.Expressions[rndExprIdx] = expr
	fn.Expressions = exprs
}

func crossover(parent1, parent2 *CXFunction) (*CXFunction, *CXFunction) {
	var child1, child2 CXFunction

	cutPoint := rand.Intn(len(parent1.Expressions))

	err := copier.Copy(&child1, *parent1)
	if err != nil {
		panic(err)
	}

	err = copier.Copy(&child2, *parent2)
	if err != nil {
		panic(err)
	}

	for c := 0; c < cutPoint; c++ {
		child1.Expressions[c] = parent2.Expressions[c]
	}

	for c := 0; c < cutPoint; c++ {
		child2.Expressions[c] = parent1.Expressions[c]
	}

	return &child1, &child2
}

func (prgrm *CXProgram) initSolution(solution []string, fns []*CXFunction, numExprs int) {
	solutionName := solution[len(solution)-1]

	pkg, err := prgrm.GetPackage(MAIN_PKG)
	if err != nil {
		panic(err)
	}

	var newPkg CXPackage
	copier.Copy(&newPkg, *pkg)
	pkgs := make([]*CXPackage, len(prgrm.Packages))
	for i, _ := range pkgs {
		pkgs[i] = prgrm.Packages[i]
	}
	prgrm.Packages = pkgs

	for i, pkg := range prgrm.Packages {
		if pkg.Name == MAIN_PKG {
			prgrm.Packages[i] = &newPkg
			break
		}
	}
	
	fn, err := prgrm.GetFunction(solutionName, MAIN_PKG)
	if err != nil {
		panic(err)
	}

	var newFn CXFunction
	newFn.Name = fn.Name
	newFn.Inputs = fn.Inputs
	newFn.Outputs = fn.Outputs
	newFn.Package = fn.Package
	// copier.Copy(&newFn, *fn)

	tmpFns := make([]*CXFunction, len(newPkg.Functions))
	for i, _ := range tmpFns {
		tmpFns[i] = newPkg.Functions[i]
	}
	newPkg.Functions = tmpFns

	for i, fn := range newPkg.Functions {
		if fn.Name == solutionName {
			newPkg.Functions[i] = &newFn
			break
		}
	}
	
	preExistingExpressions := len(newFn.Expressions)
	// Checking if we need to add more expressions.
	for i := 0; i < numExprs - preExistingExpressions; i++ {
		op := getRandFn(fns)
		expr := MakeExpression(op, "", -1)
		for c := 0; c < len(op.Inputs); c++ {
			expr.Inputs = append(expr.Inputs, getRandInp(&newFn))
		}
		// We need to add the expression at this point, so we
		// can consider this expression's output as a
		// possibility to assign stuff.
		newFn.Expressions = append(newFn.Expressions, expr)
		// Adding last expression, so output must be fn's output.
		if i == numExprs - preExistingExpressions - 1 {
			expr.Outputs = append(expr.Outputs, newFn.Outputs[0])
		} else {
			for c := 0; c < len(op.Outputs); c++ {
				expr.Outputs = append(expr.Outputs, getRandOut(&newFn))
			}
		}
	}
	newFn.Size = calcFnSize(&newFn)
	newFn.Length = numExprs
}

func (prgrm *CXProgram) injectMainInput(inp []byte) {
	// TODO: Assuming F64 output and one F64 input.
	for c := 0; c < len(inp); c++ {
		prgrm.Memory[c] = inp[c]
	}
}

func (prgrm *CXProgram) extractMainOutput() []byte {
	// TODO: Assuming F64 output and one F64 input.
	return prgrm.Memory[8:16]
}

func (prgrm *CXProgram) resetPrgrm() {
	prgrm.CallCounter = 0
	prgrm.StackPointer = 0
	prgrm.CallStack = make([]CXCall, CALLSTACK_SIZE)
	minHeapSize := minHeapSize()
	prgrm.Memory = make([]byte, STACK_SIZE+minHeapSize)
}

func mae(real, sim []float64) float64 {
	var sum float64
	for c := 0; c < len(real); c++ {
		sum += math.Abs(real[c] - sim[c])
	}
	return sum / float64(len(real))
}

func evalInd(ind *CXProgram, fp int, inputs *CXArgument, outputs *CXArgument) float64 {
	inps := ReadSliceBytes(fp, inputs, inputs.Type)
	outs := ReadSliceBytes(fp, outputs, outputs.Type)
	
	var tmp *CXProgram
	tmp = PROGRAM
	PROGRAM = ind

	// ind.PrintProgram()

	// TODO: We're calculating the error in here.
	/// Migrate to functions when we have other fitness functions.
	var sum float64
	for c := 0; c < len(inps); c += inputs.Size {
		ind.injectMainInput(inps[c:c+inputs.Size])
		ind.RunCompiled(0, nil)
		simOut := mustDeserializeF64(ind.extractMainOutput())
		realOut := mustDeserializeF64(outs[c:c+inputs.Size])
		sum += math.Abs(simOut - realOut)
	}

	PROGRAM = tmp
	
	return sum / float64((len(inps) / inputs.Size))
}

func getLowErrorIdx(errors []float64, chance float32) int {
	idx := 0
	lowest := errors[0]
	for i, err := range errors {
		if err < lowest && rand.Float32() < chance {
			lowest = err
			idx = i
		}
	}
	return idx
}

func getHighErrorIdx(errors []float64, chance float32) int {
	idx := 0
	highest := errors[0]
	for i, err := range errors {
		if err > highest && rand.Float32() < chance {
			highest = err
			idx = i
		}
	}
	return idx
}

func (ind *CXProgram) replaceSolution(solution []string, sol *CXFunction) {
	solutionName := solution[len(solution)-1]
	mainPkg, err := ind.GetPackage(MAIN_PKG)
	if err != nil {
		panic(err)
	}
	for i, fn := range mainPkg.Functions {
		if fn.Name == solutionName {
			mainPkg.Functions[i] = sol
		}
	}
	mainFn, err := mainPkg.GetFunction(MAIN_FUNC)
	if err != nil {
		panic(err)
	}
	mainFn.Expressions[0].Operator = sol
}

func opEvolve(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	inp1, inp2, inp3, inp4, inp5, inp6, inp7, inp8 := expr.Inputs[0], expr.Inputs[1], expr.Inputs[2], expr.Inputs[3], expr.Inputs[4], expr.Inputs[5], expr.Inputs[6], expr.Inputs[7]

	solution := GetInferActions(inp1, fp)

	numExprs := int(ReadI32(fp, inp5))
	numIter := ReadI32(fp, inp6)
	numPop := ReadI32(fp, inp7)
	eps := ReadF64(fp, inp8)

	fnBag := GetInferActions(inp2, fp)
	inps := ReadSliceBytes(fp, inp3, inp3.Type)
	outs := ReadSliceBytes(fp, inp4, inp4.Type)

	fns := getFnBag(prgrm, fnBag)

	_ = inps
	_ = outs

	// Initializing population.
	pop := make([]CXProgram, numPop)
	for i, _ := range pop {
		err := copier.Copy(&pop[i], prgrm)
		if err != nil {
			panic(err)
		}
		// Initialize solution with random expressions.
		pop[i].initSolution(solution, fns, numExprs)
		pop[i].adaptSolution(solution)
		
		pop[i].resetPrgrm()
	}

	// Evaluating all solutions.
	errors := make([]float64, numPop)
	for i, _ := range pop {
		errors[i] = evalInd(&pop[i], fp, inp3, inp4)
	}

	fmt.Printf("errors: %v\n", errors)

	// Crossover.
	for c := 0; c < int(numIter); c++ {
		pop1Idx := getLowErrorIdx(errors, 0.2)
		// pop1Idx := rand.Intn(int(numPop))
		pop2Idx := getLowErrorIdx(errors, 0.2)
		// pop2Idx := rand.Intn(int(numPop))
		dead1Idx := getHighErrorIdx(errors, 0.2)
		// dead1Idx := rand.Intn(int(numPop))
		dead2Idx := getHighErrorIdx(errors, 0.2)
		// dead2Idx := rand.Intn(int(numPop))
		// pop1 := pop[pop1Idx]
		// pop2 := pop[pop2Idx]

		pop1MainPkg, err := pop[pop1Idx].GetPackage(MAIN_PKG)
		if err != nil {
			panic(err)
		}
		parent1, err := pop1MainPkg.GetFunction(solution[len(solution)-1])
		if err != nil {
			panic(err)
		}

		pop2MainPkg, err := pop[pop2Idx].GetPackage(MAIN_PKG)
		if err != nil {
			panic(err)
		}
		parent2, err := pop2MainPkg.GetFunction(solution[len(solution)-1])
		if err != nil {
			panic(err)
		}

		child1, child2 := crossover(parent1, parent2)
		
		if rand.Float32() < 0.1 {
			child1.mutateFn(fns)
		}
		if rand.Float32() < 0.1 {
			child2.mutateFn(fns)
		}

		// pop[dead1Idx].PrintProgram()
		pop[dead1Idx].replaceSolution(solution, child1)
		// pop[dead1Idx].PrintProgram()
		pop[dead2Idx].replaceSolution(solution, child2)

		for i, _ := range pop {
			errors[i] = evalInd(&pop[i], fp, inp3, inp4)
			if errors[i] <= eps {
				fmt.Printf("Found solution. Bot #", i)
				fmt.Printf("errors: %v\n", errors)
				pop[i].PrintProgram()
				return
			}
		}

		fmt.Printf("errors: %v\n", errors)
	}
}
