package base

import (
	
)

func (cxt *cxContext) AddModule (mod *cxModule) *cxContext {
	mod.Context = cxt
	cxt.CurrentModule = mod
	cxt.Modules = append(cxt.Modules, mod)
	return cxt
}

func (mod *cxModule) AddDefinition (def *cxDefinition) *cxModule {
	mod.Definitions[def.Name] = def
	return mod
}

func (mod *cxModule) AddFunction (fn *cxFunction) *cxModule {
	mod.Functions[fn.Name] = fn
	return mod
}

func (mod *cxModule) AddStruct (strct *cxStruct) *cxModule {
	mod.Structs[strct.Name] = strct
	return mod
}

func (mod *cxModule) AddImport (imp *cxModule) *cxModule {
	mod.Imports[imp.Name] = imp
	return mod
}

func (strct *cxStruct) AddField (field *cxField) *cxStruct {
	strct.Fields = append(strct.Fields, field)
	return strct
}

func (fn *cxFunction) AddExpression (expr *cxExpression) *cxFunction {
	fn.Expressions = append(fn.Expressions, expr)
	return fn
}

func (fn *cxFunction) AddInput (param *cxParameter) *cxFunction {
	fn.Inputs = append(fn.Inputs, param)
	return fn
}

func (fn *cxFunction) AddOutput (param *cxParameter) *cxFunction {
	fn.Output = param
	return fn
}

func (expr *cxExpression) AddArgument (arg *cxArgument) *cxExpression {
	expr.Arguments = append(expr.Arguments, arg)
	return expr
}



// We would use these for empty interface{} inputs

// Module -> GetAffordances -> AddDefinition
// mod.ApplyAffordance()
// ApplyAffordance is a synonym for an adder
//// which adds an empty interface object, in this case &cxDefinition{}
// ApplyAffordance returns the mod itself
// We can then call ApplyAffordance again
//// whih adds a name or a type
// We can

//Problem:

// Variadic!


// We could have an "AlterAffordance" in case we want to 


// Leaf Adders for Affordances
// maybe these should be private?
func (param *cxParameter) addName (name string) *cxParameter {
	param.Name = name
	return param
}

func (def *cxDefinition) addName (name string) *cxDefinition {
	def.Name = name
	return def
}

func (def *cxDefinition) addValue (value *[]byte) *cxDefinition {
	def.Value = value
	return def
}

// Necessary leaf adders (for each cxStructure that needs it):
// AddName
// AddValue

// This one is most likely not needed, as an ApplyAffordance can
// automatically select one type, i.e. the list of types is fixed
func (param *cxParameter) AddType (name string) *cxParameter {
	param.Name = name
	return param
}
