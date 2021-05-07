package evolve

import (
	"math/rand"

	copier "github.com/jinzhu/copier"
	cxast "github.com/skycoin/cx/cx/ast"
)

// Codes associated to each of the crossover functions.
const (
	CrossoverSinglePoint = iota // Default
)

// getCrossoverFn returns the crossover function associated to `crossoverCode`.
func (pop *Population) getCrossoverFn() func(*cxast.CXFunction, *cxast.CXFunction) (*cxast.CXFunction, *cxast.CXFunction) {
	crossoverCode := pop.CrossoverMethod
	switch crossoverCode {
	case CrossoverSinglePoint:
		return singlePointCrossover
	}

	// Non-existant crossover code.
	return nil
}

// singlePointCrossover ...
func singlePointCrossover(parent1, parent2 *cxast.CXFunction) (*cxast.CXFunction, *cxast.CXFunction) {
	var child1, child2 cxast.CXFunction

	cutPointOption := len(parent1.Expressions)
	if cutPointOption > len(parent2.Expressions) {
		cutPointOption = len(parent2.Expressions)
	}
	cutPoint := rand.Intn(cutPointOption)

	err := copier.Copy(&child1, *parent1)
	if err != nil {
		panic(err)
	}
	// reprint.FromTo(parent1, &child1)

	// Replacing reference to slice.
	child1.Expressions = make([]*cxast.CXExpression, len(child1.Expressions))

	// It's okay to keep the same references to expressions, though.
	// We only want to be handling a different slice of `*CXExpression`s.
	for i, expr := range parent1.Expressions {
		child1.Expressions[i] = expr
	}

	err = copier.Copy(&child2, *parent2)
	if err != nil {
		panic(err)
	}
	// reprint.FromTo(parent2, &child2)

	// Replacing expressions as we did for `child1`.
	child2.Expressions = make([]*cxast.CXExpression, len(child2.Expressions))

	for i, expr := range parent2.Expressions {
		child2.Expressions[i] = expr
	}

	for c := 0; c < cutPoint; c++ {
		child1.Expressions[c] = parent2.Expressions[c]
	}

	for c := 0; c < cutPoint; c++ {
		child2.Expressions[c] = parent1.Expressions[c]
	}

	return &child1, &child2
}
