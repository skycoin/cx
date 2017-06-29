package base

import (
	//"github.com/skycoin/skycoin/src/cipher/encoder"
	"fmt"
	"bytes"
)


// FUNCTION -> ListAffordances -> ApplyAffordance -> FUNCTION
// Affordance <=> Adder(fixed) -> Maker(variable: types, literals)


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


// In the end, we can only return funcs or structs!!
/// We could return makers and adders
// I think we can stick with returning funcs

// 1. We can keep the idea of adding the affordances to the object
// 2.

// func (mod *cxModule) ApplyAffordance(int idx) *cxModule {
	
// }

// func (cxt *cxContext) ShowAffordances () *cxContext {
// 	for idx, aff := range cxt.Affordances {
// 		fmt.Printf("%d%+v\n", idx, aff)
// 	}
	
// 	return cxt
// }

// Another possibility is to mix some ideas from the abandoned branch:
// Always apply affordances over contexts
// This isn't a terrible idea, as the programmer or any agent won't really know
/// on what object to ask affordances

// Are getters affordances too? (most likely yes)

//MakeModule("Something").ShowAffordances() // Purely user

//MakeModule("Something").GetAffordances(fn) // This could receive a function or something
// This function would be an "agent"
// This agent can automatically select an affordance for you
// This agent could decide on an affordance in other ways

// Notes about what I thought a few hours ago:
// DONE We take subtrees using selectors
// DONE Every node needs a reference to cxt (except: cxType, cxContext, cxAffordance)
// Select (subtree) -> Affordances? -> Apply
// ---- Let's create ContextAffordances(), ModuleAffordances()
// DONE No, let's make selectors return the selected object






// We can now: Select
// Every required structure now has a reference to context
////// Ah, let's quickly modify adders
///////// No, let's do this when needed
// Now, on to the structure for an affordance:
/// Options were: to return functions
//// Can we do this now?? (This would be the best)


// Every (*cxModule)GetAffordances would be eg: AddFunction () (cxModule)
// Functions are: func()()
// This adds an empty function (with gensym name) and sets current function to this
// Returns the module

// (*cxFunction)GetAffordances    :   AddExpression

// We need these structs: ModuleAffordance, ContextAffordance, etc. so they hold funcs
/// of different types
// An affordance struct would have a pointer to a function or ACTION (they don't receive args, always return module, context, function, etc.)
// They would also have a description:
/// AddExpression (+ DEF3 INP2)
/// AddExpression (+ INP1 INP7)
/// AddInput INP2 I32
/// AddOutput OUT1 I32


func concat (str1 string, str2 string) string {
	var buffer bytes.Buffer
	buffer.WriteString(str1)
	buffer.WriteString(str2)
	return buffer.String()
}

func PrintAffordances (affs []*cxAffordance) {
	for _, aff := range affs {
		fmt.Println(aff.Description)
	}
}

func (aff *cxAffordance) ApplyAffordance () {
	aff.Action()
}

func (cxt *cxContext) GetAffordances() []*cxAffordance {
	affs := make([]*cxAffordance, 0)

	gensym := MakeGenSym()
	affs = append(affs, &cxAffordance{
		//Typ: MakeType("ContextAffordance"),
		Description: concat("AddModule ", gensym),
		// actions will be closures
		Action: func() {
			mod := MakeModule(gensym)
			cxt.Modules = append(cxt.Modules, mod)
			cxt.CurrentModule = mod
		}})
	
	return affs
}

// func (cxt *cxContext) oldGetAffordances() *cxContext {
// 	affs := make([]*cxAffordance, 0)

// 	// We can add a module to the current context
// 	affs = append(affs, &cxAffordance{
// 		Action: acAddModule,
// 		Object: cxt,
// 		Input: &cxModule{Name: MakeGenSym()},
// 	})

// 	// fixed
// 	for _, module := range cxt.Modules {
// 		affs = append(affs, &cxAffordance{
// 			Action: acSelectModule,
// 			Object: cxt,
// 			Input: module.Name,
// 		})
// 	}
// 	cxt.Affordances = affs
// 	return cxt
// }

// func (mod *cxModule) getAffordances() {
// 	affs := make([]*cxAffordance, 0)

// 	// var
// 	affs = append(affs, &cxAffordance{
// 		Action: acAddDefinition,
// 		Input: &cxDefinition{},
// 		// We need to populate these when we apply the Affordance
// 	})

// 	// var
// 	affs = append(affs, &cxAffordance{
// 		Action: acAddFunction,
// 		Input: &cxFunction{},
// 	})


// 	// This will only get the affordances
// 	// We then apply the affordance

// 	// var
// 	affs = append(affs, &cxAffordance{
// 		Action: acAddStruct,
// 		Input: &cxStruct{},
// 	})

// 	// fixed
// 	// for _, module := range mod.Context.Modules {
// 	// 	affs = append(affs, &cxAffordance{
// 	// 		Action: acAddImport,
// 	// 		Input: module,
// 	// 	})
// 	// }

// 	mod.Affordances = affs
// }

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
