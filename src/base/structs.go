package base

type cxAction string
const (
	acAddArgument cxAction = "AddArgument"
	acAddDefinition cxAction = "AddDefinition"
	acAddExpression cxAction = "AddExpression"
	acAddInput cxAction = "AddInput"
	acAddOutput cxAction = "AddOutput"
)

type cxType struct {
	Name string
	Affordances []*cxAffordance
}

type cxField struct {
	Name string
	Typ *cxType
	Affordances []*cxAffordance
}

type cxStruct struct {
	Name string
	Fields []*cxField
	Affordances []*cxAffordance
}

/*
  Context
*/

type cxContext struct {
	Modules []*cxModule
	CurrentModule *cxModule
	CurrentFunction *cxFunction
	Scopes [][]*cxDefinition
	Affordances []*cxAffordance
}

/*
  Functions
*/

type cxParameter struct {
	Name string
	Typ *cxType
	Affordances []*cxAffordance
}

type cxArgument struct {
	Typ *cxType
	Value *[]byte
	Affordances []*cxAffordance
}

type cxExpression struct {
	Fn *cxFunction
	Arguments []*cxArgument
	Affordances []*cxAffordance
}

type cxFunction struct {
	Name string
	Inputs []*cxParameter
	Output *cxParameter
	Expressions []*cxExpression
	Affordances []*cxAffordance
}

/*
  Modules
*/

type cxDefinition struct {
	Name string
	Typ *cxType
	Value *[]byte
	Affordances []*cxAffordance
}

type cxModule struct {
	Name string
	Imports map[string]*cxModule
	Functions map[string]*cxFunction
	Structs map[string]*cxStruct
	Definitions map[string]*cxDefinition
	Affordances []*cxAffordance
}

/*
  Affordances
*/

type cxAffordance struct {
	Action cxAction
	Input interface{}
}
