package base

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
}

/*
  Context
*/

type cxContext struct {
	Modules []*cxModule
	CurrentModule *cxModule
	CurrentFunction *cxFunction
	Scopes [][]*cxDefinition
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
}

/*
  Modules
*/

type cxDefinition struct {
	Name string
	Typ *cxType
	Value *[]byte
}

type cxModule struct {
	Name string
	Imports map[string]*cxModule
	Functions map[string]*cxFunction
	Structs map[string]*cxStruct
	Definitions map[string]*cxDefinition
}
