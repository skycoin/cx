package evolve

import (
	cxast "github.com/skycoin/cx/cx/ast"
	cxconstants "github.com/skycoin/cx/cx/constants"
)

type Population struct {
	Individuals      []*cxast.CXProgram
	PopulationSize   int
	ExpressionsCount int
	Iterations       int
	TargetError      float64
	FunctionToEvolve *cxast.CXFunction
	FunctionSet      []*cxast.CXFunction
	EvaluationMethod int
	CrossoverMethod  int
	MutationMethod   int
	InputSignature   []string
	OutputSignature  []string
	Inputs           [][]byte
	Outputs          [][]byte
}

// MakePopulation creates a `Population` with a number of `Individuals` equal to `populationSize`.
func MakePopulation(populationSize int) *Population {
	var pop Population
	pop.PopulationSize = populationSize
	pop.Individuals = make([]*cxast.CXProgram, populationSize)
	return &pop
}

// InitIndividuals initializes the `Individuals` in a `Population` using a `CXProgram` which works as a template. In the end, all the `Individuals` in `pop` will be exact copies of `initPrgrm` (but not pointers to the same object).
func (pop *Population) InitIndividuals(initPrgrm *cxast.CXProgram) {
	// Serializing root CX program to create copies of it.
	sPrgrm := cxast.SerializeCXProgramV2(initPrgrm, true)

	for i := 0; i < len(pop.Individuals); i++ {
		pop.Individuals[i] = cxast.DeserializeCXProgramV2(sPrgrm)
	}
}

// InitFunctionSet gathers the functions contained in `prgrm` named by `fnNames`.
func (pop *Population) InitFunctionSet(fnNames []string) {
	pop.FunctionSet = GetFunctionSet(fnNames)
}

// InitFunctionsToEvolve initializes the `FunctionToEvolve` in each of the individuals in a `Population`, so each individual has a `FunctionToEvolve` with a random set of expressions.
func (pop *Population) InitFunctionsToEvolve(fnName string) {
	prgrm := pop.Individuals[0]
	fnToEvolve, err := prgrm.GetFunction(fnName, cxconstants.MAIN_PKG)
	if err != nil {
		panic(err)
	}
	pop.FunctionToEvolve = fnToEvolve

	numExprs := pop.ExpressionsCount
	fns := pop.FunctionSet

	for i := 0; i < len(pop.Individuals); i++ {
		// Initialize solution with random expressions.
		initSolution(pop.Individuals[i], fnToEvolve, fns, numExprs)
		adaptSolution(pop.Individuals[i], fnToEvolve)
		resetPrgrm(pop.Individuals[i])
	}
}

// SetInputs sets the `Inputs` of a `Population` to `inputs`.
func (pop *Population) SetInputs(inputs [][]byte) {
	pop.Inputs = inputs
}

// SetOutputs sets the `Outputs` of a `Population` to `outputs`.
func (pop *Population) SetOutputs(outputs [][]byte) {
	pop.Outputs = outputs
}

// SetTargetError sets the `TargetError` of a `Population` to `targetError`.
func (pop *Population) SetTargetError(targetError float64) {
	pop.TargetError = targetError
}

// SetIterations sets the `Iterations` of a `Population` to `iter`.
func (pop *Population) SetIterations(iter int) {
	pop.Iterations = iter
}

// SetExpressionsCount sets the `ExpressionsCount` of a `Population` to `exprCount`.
func (pop *Population) SetExpressionsCount(exprCount int) {
	pop.ExpressionsCount = exprCount
}
