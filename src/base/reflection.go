package base

import (
	"github.com/skycoin/skycoin/src/cipher/encoder"
)

// These would be part of the "functions_of"

func (strct *cxStruct) GetFields() []*cxField {
	return strct.Fields
}

func (mod *cxModule) GetFunctions() []*cxFunction {
	funcs := make([]*cxFunction, len(mod.Functions))
	i := 0
	for _, v := range mod.Functions {
		funcs[i] = v
		i++
	}
	return funcs
}



// This is one idea
func (mod *cxModule) ApplyAffordance(aff *cxAffordance) *cxModule {
	return nil
}

func (mod *cxFunction) ApplyAffordance(aff *cxAffordance) *cxFunction {
	return nil
}

// all cxStructures could have a field Affordances
// this way: fn.PopulateAffordances().ApplyAffordance()
// populating a cxStructure with its affordances would be automatically, so we don't need a PopulateAffordances method

// everytime we call a maker or an adder, we can populate affordances for each object
// or maybe not, because other processes could add new affordances and these would need to be recalculated anyway
// the best approach would be to recalculate everytime we require them


// if I receive an affordance with a cxFunction, what does this mean? that I can call it on this object? Is this always going to be the case?
// for example, if I receive an expression for a function, what does this mean: I can add it to this function

//list of possible affordances:
// module: AddDefinition
// apparently, we should focus on adders, makers and already defined stuff



// how can we get the current module? only by using context?
//// There's always only one context object, so we could use this
//*****
//// Or every cxStructure could have a pointer to its module, an expression a pointer to its function, etc. this would contribute to reflection
//*****

// should GetAffordances() print the options to the user? i.e. should it work only with side effects?
// no, it should return the affordances, because: what if we want to create a function which takes a decision based on the affordances it returns
// again, we can bypass this by populating the cxStructure with the affordances. If we require a function which operates on the affordances of a cxStructure, we only access them via cxStructure.Affordances



// Now, how do we apply an affordance?
// ApplyAffordance should be on an object, i.e. func (fn *cxFunction) ApplyAffordance()
// And it should return the same object, and thus, we can call ApplyAffordance several times
///// How do we tell the method what affordance to appyl?
//////// The first solution is to use an index
//////// Another solution is to describe the affordance, but this doesn't make sense, because we're not supposed to know anything about an afordance, the programming language should tell us

// we need another method which is UpdateAffordanecs. GetAffordances would be part of the reflection functions, Populate and Apply would be part of the affordance.go file


// The idea is to call UpdateAffordances() and then access the struct's affordances field


// These are the scenarios:
// 1. The programmer wants to know the affordances: He will call GetAffordances to know, then call ApplyAffordance
// 2. The program itself wants to know: It will simply call ApplyAffordances


// Needs to be a private method
func (fn *cxFunction) updateAffordances() *cxFunction {
	num1 := encoder.SerializeAtomic(int32(25))
	
	// Next step: how do we determine these available actions
	/// Some are fixed: addargument, addinput, addoutput
	aff1 := &cxAffordance{Action: acAddArgument,
		// The values can vary, as well as type
		// How can we help the programmer make an argument
		// We'll think about that later I guess
		Input: MakeArgument(&num1, MakeType("i32"))}

	aff2 := &cxAffordance{Action: acAddInput,
		Input: MakeParameter("something1", MakeType("i32"))}

	aff3 := &cxAffordance{Action: acAddOutput,
		Input: MakeParameter("something2", MakeType("i32"))}

	fn.Affordances = append(fn.Affordances, aff1)
	fn.Affordances = append(fn.Affordances, aff2)
	fn.Affordances = append(fn.Affordances, aff3)
	
	return fn
	//AddExpression // This could be type, but name it Action
	//cxExpression

	// Action // These are fixed: "AddExpression", "AddField", "AddArgument"
	// Input
}





func (fn *cxFunction) GetAffordances() []*cxAffordance {
	// Can be parameter (input and output) and expression
	//fn.Affordances =
	return fn.Affordances
}

// Operators
// Functions
// int32 => add, sub, mul, div

// In the end, we can only return funcs or structs!!
/// We could return makers and adders
// I think we can stick with returning funcs

func (mod *cxModule) GetAffordances() {
	
}

func (typ *cxType) GetAffordances() {
	// if cxType.Name == "i32" {}

	// returns what cx functions are using this type
	/// we would need a list/map/array of all functions by type
	
}

func (strct *cxStruct) GetAffordances() {
	
}


