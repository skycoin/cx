package base

import (

)

func saveProgramStep (prgrmStep *CXProgramStep, cxt *CXContext) {
	cxt.ProgramSteps = append(cxt.ProgramSteps, prgrmStep)
}

func (cxt *CXContext) AddModule (mod *CXModule) *CXContext {
	stepMod := MakeModuleCopy(mod, cxt)
	prgrmStep := &CXProgramStep{
		Action: func(cxt *CXContext) {
			newStep := MakeModuleCopy(stepMod, cxt)
			cxt.CurrentModule = newStep
			cxt.Modules[newStep.Name] = newStep
		},
	}
	saveProgramStep(prgrmStep, cxt)
	
	mod.Context = cxt
	cxt.CurrentModule = mod
	cxt.Modules[mod.Name] = mod
	return cxt
}

func (mod *CXModule) AddDefinition (def *CXDefinition) *CXModule {
	// identParts := getIdentParts(def.Name)
	// // we're ignoring nested structs for now
	// if len(identParts) == 2 {
		
	// }
	
	stepDef := MakeDefinitionCopy(def, mod, mod.Context)
	prgrmStep := &CXProgramStep{
		Action: func(cxt *CXContext) {
			if mod, err := cxt.GetCurrentModule(); err == nil {
				newDef := MakeDefinitionCopy(stepDef, mod, cxt)
				mod.Definitions[newDef.Name] = newDef
			}
		},
	}
	saveProgramStep(prgrmStep, mod.Context)
	
	def.Context = mod.Context
	def.Module = mod
	mod.Definitions[def.Name] = def
	return mod
}

func (mod *CXModule) AddFunction (fn *CXFunction) *CXModule {
	stepFn := MakeFunctionCopy(fn, mod, mod.Context)
	prgrmStep := &CXProgramStep{
		Action: func(cxt *CXContext) {
			if mod, err := cxt.GetCurrentModule(); err == nil {
				newFn := MakeFunctionCopy(stepFn, mod, cxt)
				mod.CurrentFunction = newFn
				mod.Functions[newFn.Name] = newFn
			}
		},
	}
	saveProgramStep(prgrmStep, mod.Context)
	
	fn.Context = mod.Context
	fn.Module = mod
	mod.CurrentFunction = fn
	mod.Functions[fn.Name] = fn
	return mod
}

func (mod *CXModule) AddStruct (strct *CXStruct) *CXModule {
	stepStrct := MakeStructCopy(strct, mod, mod.Context)
	prgrmStep := &CXProgramStep{
		Action: func(cxt *CXContext) {
			if mod, err := cxt.GetCurrentModule(); err == nil {
				newStrct := MakeStructCopy(stepStrct, mod, cxt)
				mod.CurrentStruct = newStrct
				mod.Structs[newStrct.Name] = newStrct
			}
		},
	}
	saveProgramStep(prgrmStep, mod.Context)
	
	cxt := mod.Context
	strct.Context = cxt
	strct.Module = mod
	mod.CurrentStruct = strct
	mod.Structs[strct.Name] = strct
	return mod
}

func (mod *CXModule) AddImport (imp *CXModule) *CXModule {
	prgrmStep := &CXProgramStep{
		Action: func(cxt *CXContext) {
			if mod, err := cxt.GetCurrentModule(); err == nil {
				mod.Imports[imp.Name] = cxt.Modules[imp.Name]
			}
		},
	}
	saveProgramStep(prgrmStep, mod.Context)
	
	mod.Imports[imp.Name] = imp
	return mod
}

func (strct *CXStruct) AddField (fld *CXField) *CXStruct {
	stepFld := MakeFieldCopy(fld)
	prgrmStep := &CXProgramStep{
		Action: func(cxt *CXContext) {
			newFld := MakeFieldCopy(stepFld)
			if strct, err := cxt.GetCurrentStruct(); err == nil {
				strct.Fields = append(strct.Fields, newFld)
			}
		},
	}
	saveProgramStep(prgrmStep, strct.Context)
	
	strct.Fields = append(strct.Fields, fld)
	return strct
}

func (fn *CXFunction) AddExpression (expr *CXExpression) *CXFunction {
	stepExpr := MakeExpressionCopy(expr, fn, fn.Module, fn.Context)
	prgrmStep := &CXProgramStep{
		Action: func(cxt *CXContext) {
			if mod, err := cxt.GetCurrentModule(); err == nil {
				if fn, err := cxt.GetCurrentFunction(); err == nil {
					newExpr := MakeExpressionCopy(stepExpr, fn, mod, cxt)
					var exprOperator *CXFunction
					for _, mod := range cxt.Modules {
						for _, fn := range mod.Functions {
							if fn.Name == expr.Operator.Name &&
								mod.Name == expr.Operator.Module.Name {
								exprOperator = fn
							}
						}
					}
					if exprOperator != nil {
						newExpr.Operator = exprOperator
					} else {
						panic("AddExpression: Expression operator not found when creating program step")
					}
					
					fn.CurrentExpression = newExpr
					fn.Expressions = append(fn.Expressions, newExpr)
				}
			}
		},
	}
	saveProgramStep(prgrmStep, fn.Context)
	
	expr.Context = fn.Context
	expr.Module = fn.Module
	expr.Function = fn
	expr.Line = len(fn.Expressions)
	fn.Expressions = append(fn.Expressions, expr)
	fn.CurrentExpression = expr
	return fn
}

func (fn *CXFunction) AddInput (param *CXParameter) *CXFunction {
	stepParam := MakeParameterCopy(param)
	prgrmStep := &CXProgramStep{
		Action: func(cxt *CXContext) {
			if fn, err := cxt.GetCurrentFunction(); err == nil {
				fn.Inputs = append(fn.Inputs, MakeParameterCopy(stepParam))
			}
		},
	}
	saveProgramStep(prgrmStep, fn.Context)
	
	fn.Inputs = append(fn.Inputs, param)
	return fn
}

func (fn *CXFunction) AddOutput (param *CXParameter) *CXFunction {
	stepParam := MakeParameterCopy(param)
	prgrmStep := &CXProgramStep{
		Action: func(cxt *CXContext) {
			if fn, err := cxt.GetCurrentFunction(); err == nil {
				fn.Outputs = append(fn.Outputs, MakeParameterCopy(stepParam))
			}
		},
	}
	saveProgramStep(prgrmStep, fn.Context)

	fn.Outputs = append(fn.Outputs, param)
	return fn
}

func (expr *CXExpression) AddArgument (arg *CXArgument) *CXExpression {
	stepArg := MakeArgumentCopy(arg)
	prgrmStep := &CXProgramStep{
		Action: func(cxt *CXContext) {
			if expr, err := cxt.GetCurrentExpression(); err == nil {
				expr.Arguments = append(expr.Arguments, MakeArgumentCopy(stepArg))
			}
		},
	}
	saveProgramStep(prgrmStep,
		expr.Context)

	expr.Arguments = append(expr.Arguments, arg)
	return expr
}

func (expr *CXExpression) AddOutputName (outName string) *CXExpression {
	expr.OutputNames = append(expr.OutputNames, outName)
	return expr
}
