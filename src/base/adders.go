package base

import (
	
)

func (mod *cxModule) AddDefinition(def *cxDefinition) *cxModule {
	mod.Definitions[def.Name] = def
	return mod
}

func (mod *cxModule) AddFunction(fn *cxFunction) *cxModule {
	mod.Functions[fn.Name] = fn
	return mod
}

func (mod *cxModule) AddStruct(strct *cxStruct) *cxModule {
	mod.Structs[strct.Name] = strct
	return mod
}

func (mod *cxModule) AddImport(imp *cxModule) *cxModule {
	mod.Imports[imp.Name] = imp
	return mod
}

func (strct *cxStruct) AddField(field *cxField) *cxStruct {
	strct.Fields = append(strct.Fields, field)
	return strct
}

func (fn *cxFunction) AddExpression (expr *cxExpression) *cxFunction {
	fn.Expressions = append(fn.Expressions, expr)
	return fn
}

func (fn *cxFunction) AddInput(param *cxParameter) *cxFunction {
	fn.Inputs = append(fn.Inputs, param)
	return fn
}

func (fn *cxFunction) AddOutput(param *cxParameter) *cxFunction {
	fn.Output = param
	return fn
}

func (expr *cxExpression) AddArgument(arg *cxArgument) *cxExpression {
	expr.Arguments = append(expr.Arguments, arg)
	return expr
}


// Adders for Affordances

func (param *cxParameter) AddName (name string) *cxParameter {
	param.Name = name
	return param
}

func (param *cxParameter) AddType (name string) *cxParameter {
	param.Name = name
	return param
}
