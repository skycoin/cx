package base

type cxType struct {
	Name string
}
// used for affordances (and maybe other stuff)
var basicTypes = []string{"i32"}

type cxField struct {
	Name string
	Typ *cxType
}

type cxStruct struct {
	Name string
	Fields []*cxField

	Context *cxContext
	Module *cxModule
}

/*
  Context
*/

type cxContext struct {
	Modules []*cxModule
	CurrentModule *cxModule
	//Scopes [][]*cxDefinition
}

/*
  Functions
*/

type cxParameter struct {
	Name string
	Typ *cxType
}

type cxArgument struct {
	Typ *cxType
	Value *[]byte
}

type cxExpression struct {
	Operator *cxFunction
	Arguments []*cxArgument

	Line int
	Function *cxFunction
	Module *cxModule
	Context *cxContext
}

type cxFunction struct {
	Name string
	Inputs []*cxParameter
	Output *cxParameter
	Expressions []*cxExpression

	CurrentExpression *cxExpression
	Module *cxModule
	Context *cxContext
}

/*
  Modules
*/

type cxDefinition struct {
	Name string
	Typ *cxType
	Value *[]byte

	Module *cxModule
	Context *cxContext
}

type cxModule struct {
	Name string
	Imports map[string]*cxModule
	Functions map[string]*cxFunction
	Structs map[string]*cxStruct
	Definitions map[string]*cxDefinition

	CurrentFunction *cxFunction
	CurrentStruct *cxStruct
	Context *cxContext
}

/*
  Affordances
*/

type cxAffordance struct {
	//Typ *cxType // determines how to cast action // might not be needed actually, let's see
	Description string
	Action func() // func()*cxContext, func()*cxModule, etc
	//Object interface{} // to what object are we going to apply the action
	// might not be needed either, as this is a closure having everything needed
	// e.g. *cxModule{}, *cxFunction{}, etc.
}

