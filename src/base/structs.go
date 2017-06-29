package base

// type cxAction string
// const (
// 	// adders
// 	acAddArgument cxAction = "AddArgument"
// 	acAddModule cxAction = "AddModule"
// 	acAddDefinition cxAction = "AddDefinition"
// 	acAddFunction cxAction = "AddFunction"
// 	acAddStruct cxAction = "AddStruct"
// 	acAddImport cxAction = "AddImport"
// 	acAddExpression cxAction = "AddExpression"
// 	acAddInput cxAction = "AddInput"
// 	acAddOutput cxAction = "AddOutput"

// 	// selectors
// 	acSelectModule cxAction = "SelectModule"
// 	acSelectFunction cxAction = "SelectFunction"
// 	acSelectStruct cxAction = "SelectStruct"
// 	//acSelectExpression cxAction = "SelectExpression"
// 	// How do we handle arguments in expressions??
// 	// We would need to show every single combination with definitions
// 	//// Maybe this is the correct way
// 	//// If this is the case, then we don't need to select expressions
// )

type cxType struct {
	Name string
}

type cxField struct {
	Name string
	Typ *cxType
}

type cxStruct struct {
	Name string
	Fields []*cxField

	Context *cxContext
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
	Fn *cxFunction
	Arguments []*cxArgument
}

type cxFunction struct {
	Name string
	Inputs []*cxParameter
	Output *cxParameter
	Expressions []*cxExpression

	Context *cxContext
}

/*
  Modules
*/

type cxDefinition struct {
	Name string
	Typ *cxType
	Value *[]byte

	Context *cxContext
}

type cxModule struct {
	Name string
	Imports map[string]*cxModule
	Functions map[string]*cxFunction
	Structs map[string]*cxStruct
	Definitions map[string]*cxDefinition

	Types map[string]*cxType
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

