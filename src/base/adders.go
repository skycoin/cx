package base

import (

)

func saveProgramStep (prgrmStep *CXProgramStep, cxt *CXProgram) {
	cxt.ProgramSteps = append(cxt.ProgramSteps, prgrmStep)
}

func (cxt *CXProgram) AddModule (mod *CXModule) *CXProgram {
	//stepMod := MakeModuleCopy(mod, cxt)
	// prgrmStep := &CXProgramStep{
	// 	Action: func(cxt *CXProgram) {
	// 		newMod := MakeModuleCopy(stepMod, cxt)
	// 		cxt.CurrentModule = newMod
	// 		//cxt.Modules[newMod.Name] = newMod
	// 		found := false
	// 		for i, md := range cxt.Modules {
	// 			if md.Name == newMod.Name {
	// 				cxt.Modules[i] = newMod
	// 				found = true
	// 				break
	// 			}
	// 		}
	// 		if !found {
	// 			cxt.Modules = append(cxt.Modules, newMod)
	// 		}
	// 	},
	// }


	
	// if mod.Name != CORE_MODULE {
	// 	prgrmStep := &CXProgramStep{
	// 		Action: func () {
	// 			cxt.RemoveModule(mod.Name)
	// 		},
	// 	}
	// 	saveProgramStep(prgrmStep, cxt)
	// }
	
	mod.Context = cxt
	//cxt.Modules[mod.Name] = mod
	found := false
	for _, md := range cxt.Modules {
		if md.Name == mod.Name {
			cxt.CurrentModule = md
			found = true
			break
		}
	}
	if !found {
		cxt.Modules = append(cxt.Modules, mod)
		cxt.CurrentModule = mod
	}
	return cxt
}

func (mod *CXModule) AddDefinition (def *CXDefinition) *CXModule {
	//stepDef := MakeDefinitionCopy(def, mod, mod.Context)
	// prgrmStep := &CXProgramStep{
	// 	Action: func(cxt *CXProgram) {
	// 		if mod, err := cxt.GetCurrentModule(); err == nil {
	// 			newDef := MakeDefinitionCopy(stepDef, mod, cxt)
	// 			//mod.Definitions[newDef.Name] = newDef
	// 			found := false
	// 			for i, df := range mod.Definitions {
	// 				if df.Name == newDef.Name {
	// 					mod.Definitions[i] = newDef
	// 					found = true
	// 					break
	// 				}
	// 			}
	// 			if !found {
	// 				mod.Definitions = append(mod.Definitions, newDef)
	// 			}
	// 		}
	// 	},
	// }

	// prgrmStep := &CXProgramStep{
	// 	Action: func () {
	// 		mod.RemoveDefinition(def.Name)
	// 	},
	// }
	// saveProgramStep(prgrmStep, mod.Context)
	
	def.Context = mod.Context
	def.Module = mod
	//mod.Definitions[def.Name] = def
	found := false
	for i, df := range mod.Definitions {
		if df.Name == def.Name {
			mod.Definitions[i] = def
			found = true
			break
		}
	}
	if !found {
		mod.Definitions = append(mod.Definitions, def)
	}
	return mod
}

func (mod *CXModule) AddFunction (fn *CXFunction) *CXModule {
	// stepFn := MakeFunctionCopy(fn, mod, mod.Context)
	// prgrmStep := &CXProgramStep{
	// 	Action: func(cxt *CXProgram) {
	// 		if mod, err := cxt.GetCurrentModule(); err == nil {
	// 			newFn := MakeFunctionCopy(stepFn, mod, cxt)
	// 			mod.CurrentFunction = newFn
	// 			//mod.Functions[newFn.Name] = newFn
	// 			found := false
	// 			for i, fn := range mod.Functions {
	// 				if fn.Name == newFn.Name {
	// 					mod.Functions[i] = newFn
	// 					found = true
	// 					break
	// 				}
	// 			}
	// 			if !found {
	// 				mod.Functions = append(mod.Functions, newFn)
	// 			}
	// 		}
	// 	},
	// }
	
	// isNative := false
	// for _, native := range NATIVE_FUNCTIONS {
	// 	if native == fn.Name {
	// 		isNative = true
	// 	}
	// }
	// if !isNative {
	// 	prgrmStep := &CXProgramStep{
	// 		Action: func () {
	// 			mod.RemoveFunction(fn.Name)
	// 		},
	// 	}
	// 	saveProgramStep(prgrmStep, mod.Context)
	// }
	
	fn.Context = mod.Context
	fn.Module = mod
	mod.CurrentFunction = fn
	//mod.Functions[fn.Name] = fn
	found := false
	for i, f := range mod.Functions {
		if f.Name == fn.Name {
			mod.Functions[i] = fn
			found = true
			break
		}
	}
	if !found {
		mod.Functions = append(mod.Functions, fn)
	}
	return mod
}

func (mod *CXModule) AddStruct (strct *CXStruct) *CXModule {
	// stepStrct := MakeStructCopy(strct, mod, mod.Context)
	// prgrmStep := &CXProgramStep{
	// 	Action: func(cxt *CXProgram) {
	// 		if mod, err := cxt.GetCurrentModule(); err == nil {
	// 			newStrct := MakeStructCopy(stepStrct, mod, cxt)
	// 			mod.CurrentStruct = newStrct
	// 			//mod.Structs[newStrct.Name] = newStrct
	// 			found := false
	// 			for i, strct := range mod.Structs {
	// 				if strct.Name == newStrct.Name {
	// 					mod.Structs[i] = newStrct
	// 					found = true
	// 					break
	// 				}
	// 			}
	// 			if !found {
	// 				mod.Structs = append(mod.Structs, newStrct)
	// 			}
	// 		}
	// 	},
	// }
	// saveProgramStep(prgrmStep, mod.Context)

	// prgrmStep := &CXProgramStep{
	// 	Action: func () {
	// 		mod.RemoveStruct(strct.Name)
	// 	},
	// }
	// saveProgramStep(prgrmStep, mod.Context)

	cxt := mod.Context
	strct.Context = cxt
	strct.Module = mod
	mod.CurrentStruct = strct
	//mod.Structs[strct.Name] = strct
	found := false
	for i, s := range mod.Structs {
		if s.Name == strct.Name {
			mod.Structs[i] = strct
			found = true
			break
		}
	}
	if !found {
		mod.Structs = append(mod.Structs, strct)
	}
	return mod
}

func (mod *CXModule) AddImport (imp *CXModule) *CXModule {
	// prgrmStep := &CXProgramStep{
	// 	Action: func(cxt *CXProgram) {
	// 		if mod, err := cxt.GetCurrentModule(); err == nil {
	// 			//mod.Imports[imp.Name] = cxt.Modules[imp.Name]
	// 			found := false
	// 			for _, im := range mod.Imports {
	// 				if im.Name == imp.Name {
	// 					found = true
	// 					break
	// 				}
	// 			}
	// 			if !found {
	// 				for _, md := range cxt.Modules {
	// 					if md.Name == imp.Name {
	// 						mod.Imports = append(mod.Imports, md)
	// 					}
	// 				}
	// 			}
	// 		}
	// 	},
	// }
	// saveProgramStep(prgrmStep, mod.Context)

	// prgrmStep := &CXProgramStep{
	// 	Action: func () {
	// 		mod.RemoveImport(imp.Name)
	// 	},
	// }
	// saveProgramStep(prgrmStep, mod.Context)
	
	//mod.Imports[imp.Name] = imp
	found := false
	for _, im := range mod.Imports {
		if im.Name == imp.Name {
			found = true
			break
		}
	}
	if !found {
		mod.Imports = append(mod.Imports, imp)
	}
	
	return mod
}

func (mod *CXModule) AddObject (obj string) *CXModule {
	mod.Objects = append(mod.Objects, obj)

	return mod
}

func (mod *CXModule) AddClauses (clauses string) *CXModule {
	mod.Clauses = clauses
	return mod
}

func (mod *CXModule) AddQuery (query string) *CXModule {
	mod.Query = query
	return mod
}

func (strct *CXStruct) AddField (fld *CXField) *CXStruct {
	// stepFld := MakeFieldCopy(fld)
	// prgrmStep := &CXProgramStep{
	// 	Action: func(cxt *CXProgram) {
	// 		newFld := MakeFieldCopy(stepFld)
	// 		if strct, err := cxt.GetCurrentStruct(); err == nil {
	// 			strct.Fields = append(strct.Fields, newFld)
	// 		}
	// 	},
	// }
	// saveProgramStep(prgrmStep, strct.Context)

	// prgrmStep := &CXProgramStep{
	// 	Action: func () {
	// 		strct.RemoveField(fld.Name)
	// 	},
	// }
	// saveProgramStep(prgrmStep, strct.Context)

	found := false
	for _, fl := range strct.Fields {
		if fl.Name == fld.Name {
			found = true
			break
		}
	}
	if !found {
		strct.Fields = append(strct.Fields, fld)
	}
	return strct
}

func (fn *CXFunction) AddExpression (expr *CXExpression) *CXFunction {
	// stepExpr := MakeExpressionCopy(expr, fn, fn.Module, fn.Context)
	// prgrmStep := &CXProgramStep{
	// 	Action: func(cxt *CXProgram) {
	// 		if mod, err := cxt.GetCurrentModule(); err == nil {
	// 			if fn, err := cxt.GetCurrentFunction(); err == nil {
	// 				newExpr := MakeExpressionCopy(stepExpr, fn, mod, cxt)
	// 				var exprOperator *CXFunction
	// 				for _, mod := range cxt.Modules {
	// 					for _, fn := range mod.Functions {
	// 						if fn.Name == expr.Operator.Name &&
	// 							mod.Name == expr.Operator.Module.Name {
	// 							exprOperator = fn
	// 						}
	// 					}
	// 				}
	// 				if exprOperator != nil {
	// 					newExpr.Operator = exprOperator
	// 				} else {
	// 					panic("AddExpression: Expression operator not found when creating program step")
	// 				}
					
	// 				fn.CurrentExpression = newExpr
	// 				fn.Expressions = append(fn.Expressions, newExpr)
	// 			}
	// 		}
	// 	},
	// }
	// saveProgramStep(prgrmStep, fn.Context)

	// prgrmStep := &CXProgramStep{
	// 	Action: func () {
	// 		fn.RemoveExpression()
	// 	},
	// }
	// saveProgramStep(prgrmStep, fn.Context)
	
	expr.Context = fn.Context
	expr.Module = fn.Module
	expr.Function = fn
	expr.Line = len(fn.Expressions)
	fn.Expressions = append(fn.Expressions, expr)
	fn.CurrentExpression = expr
	return fn
}

func (fn *CXFunction) AddInput (param *CXParameter) *CXFunction {
	// stepParam := MakeParameterCopy(param)
	// prgrmStep := &CXProgramStep{
	// 	Action: func(cxt *CXProgram) {
	// 		if fn, err := cxt.GetCurrentFunction(); err == nil {
	// 			fn.Inputs = append(fn.Inputs, MakeParameterCopy(stepParam))
	// 		}
	// 	},
	// }
	// saveProgramStep(prgrmStep, fn.Context)

	// prgrmStep := &CXProgramStep{
	// 	Action: func () {
	// 		fn.RemoveInput(param.Name)
	// 	},
	// }
	// saveProgramStep(prgrmStep, fn.Context)

	found := false
	for _, inp := range fn.Inputs {
		if inp.Name == param.Name {
			found = true
			break
		}
	}
	if !found {
		fn.Inputs = append(fn.Inputs, param)
	}
	
	return fn
}

func (fn *CXFunction) AddOutput (param *CXParameter) *CXFunction {
	// stepParam := MakeParameterCopy(param)
	// prgrmStep := &CXProgramStep{
	// 	Action: func(cxt *CXProgram) {
	// 		if fn, err := cxt.GetCurrentFunction(); err == nil {
	// 			fn.Outputs = append(fn.Outputs, MakeParameterCopy(stepParam))
	// 		}
	// 	},
	// }
	// saveProgramStep(prgrmStep, fn.Context)

	// prgrmStep := &CXProgramStep{
	// 	Action: func () {
	// 		fn.RemoveOutput(param.Name)
	// 	},
	// }
	// saveProgramStep(prgrmStep, fn.Context)

	found := false
	for _, out := range fn.Outputs {
		if out.Name == param.Name {
			found = true
			break
		}
	}
	if !found {
		fn.Outputs = append(fn.Outputs, param)
	}

	return fn
}

func (expr *CXExpression) AddArgument (arg *CXArgument) *CXExpression {
	// stepArg := MakeArgumentCopy(arg)
	// prgrmStep := &CXProgramStep{
	// 	Action: func(cxt *CXProgram) {
	// 		if expr, err := cxt.GetCurrentExpression(); err == nil {
	// 			expr.Arguments = append(expr.Arguments, MakeArgumentCopy(stepArg))
	// 		}
	// 	},
	// }
	// saveProgramStep(prgrmStep, expr.Context)

	// prgrmStep := &CXProgramStep{
	// 	Action: func () {
	// 		expr.RemoveArgument()
	// 	},
	// }
	// saveProgramStep(prgrmStep, expr.Context)

	expr.Arguments = append(expr.Arguments, arg)
	return expr
}

func (expr *CXExpression) AddOutputName (outName string) *CXExpression {
	// // was going to create the addstep procedure, but all of them are wrong now (because of removers)
	// stepOutName := MakeArgumentCopy(arg)
	// prgrmStep := &CXProgramStep{
	// 	Action: func(cxt *CXProgram) {
	// 		if expr, err := cxt.GetCurrentExpression(); err == nil {
	// 			expr.Arguments = append(expr.Arguments, MakeArgumentCopy(stepArg))
	// 		}
	// 	},
	// }
	// saveProgramStep(prgrmStep, expr.Context)

	// prgrmStep := &CXProgramStep{
	// 	Action: func () {
	// 		expr.RemoveOutputName()
	// 	},
	// }
	// saveProgramStep(prgrmStep, expr.Context)

	// The error raises because everything in CX needs to return something
	// The outputName is created and is waiting for an output
	// A solution would be to add a nil output in case there aren't any

	if len(expr.Operator.Outputs) > 0 {
		nextOutIdx := len(expr.OutputNames)
		outDef := MakeDefinition(
			outName,
			MakeDefaultValue(expr.Operator.Outputs[nextOutIdx].Typ.Name),
			expr.Operator.Outputs[nextOutIdx].Typ)
		expr.OutputNames = append(expr.OutputNames, outDef)
	}
	
	return expr
}

func (expr *CXExpression) AddTag (tag string) *CXExpression {
	expr.Tag = tag
	return expr
}

func (mod *CXModule) AddTag (tag string) *CXModule {
	mod.Tag = tag
	return mod
}

func (strct *CXStruct) AddTag (tag string) *CXStruct {
	strct.Tag = tag
	return strct
}

func (fn *CXFunction) AddTag (tag string) *CXFunction {
	fn.Tag = tag
	return fn
}
