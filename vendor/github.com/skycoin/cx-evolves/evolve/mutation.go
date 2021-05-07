package evolve

import (
	"math/rand"

	cxast "github.com/skycoin/cx/cx/ast"
)

// Codes associated to each of the mutation functions.
const (
	MutationRandom = iota // Default
	MutationMirror
	MutationBitFlip
)

// getCrossoverFn returns the crossover function associated to `mutationCode`.
// func (pop *Population) getMutationFn(mutationCode int) func(*cxast.CXFunction) {
// 	switch mutationCode {
// 	case MutationRandom:
// 		return randomMutation
// 	case MutationMirror:
// 		return mirrorMutation
// 	case MutationBitFlip:
// 		return bitFlipMutation
// 	}
// }

// mirrorMutation swaps a gene (*CXExpression) from fn.Expressions (our genome) in a mirror-like manner.
// func mirrorMutation(fn *cxcore.CXFunction) {
// 	randIdx := rand.Intn(len(fn.Expressions))
// 	tmpExpr := fn.Expressions[randIdx]
// 	mirrorIdx := len(fn.Expressions) - randIdx - 1
// 	fn.Expressions[randIdx] = fn.Expressions[mirrorIdx]
// 	fn.Expressions[mirrorIdx] = tmpExpr
// }

func randomMutation(pop *Population, sPrgrm []byte) {
	fnToEvolve := pop.FunctionToEvolve
	numExprs := pop.ExpressionsCount
	fns := pop.FunctionSet
	randIdx := rand.Intn(len(pop.Individuals))
	pop.Individuals[randIdx] = cxast.DeserializeCXProgramV2(sPrgrm)
	initSolution(pop.Individuals[randIdx], fnToEvolve, fns, numExprs)
	adaptSolution(pop.Individuals[randIdx], fnToEvolve)
	resetPrgrm(pop.Individuals[randIdx])
}

// func bitflipMutation(fn *cxcore.CXFunction, fnBag []*cxcore.CXFunction) {
// 	rndExprIdx := rand.Intn(len(fn.Expressions))
// 	rndFn := getRandFn(fnBag)

// 	expr := cxcore.MakeExpression(rndFn, "", -1)
// 	expr.Package = fn.Package
// 	expr.Inputs = fn.Expressions[rndExprIdx].Inputs
// 	expr.Outputs = fn.Expressions[rndExprIdx].Outputs

// 	exprs := make([]*cxcore.CXExpression, len(fn.Expressions))
// 	for i, ex := range fn.Expressions {
// 		if i == rndExprIdx {
// 			exprs[i] = expr
// 		} else {
// 			exprs[i] = ex
// 		}
// 	}

// 	// fn.Expressions[rndExprIdx] = expr
// 	fn.Expressions = exprs
// }
