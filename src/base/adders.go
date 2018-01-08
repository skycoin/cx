package base

import (
	//"github.com/skycoin/skycoin/src/cipher/encoder"
)

func (prgrm *CXProgram) AddModule (mod *CXModule) *CXProgram {
	mod.Program = prgrm
	found := false
	for _, md := range prgrm.Modules {
		if md.Name == mod.Name {
			prgrm.CurrentModule = md
			found = true
			break
		}
	}
	if !found {
		prgrm.Modules = append(prgrm.Modules, mod)
		prgrm.CurrentModule = mod
	}
	return prgrm
}

func (mod *CXModule) AddGlobal (def *CXArgument) *CXModule {
	def.Program = mod.Program
	def.Module = mod
	found := false
	for i, df := range mod.Globals {
		if df.Name == def.Name {
			mod.Globals[i] = def
			found = true
			break
		}
	}
	if !found {
		mod.Globals = append(mod.Globals, def)
	}
	return mod
}

func (mod *CXModule) AddFunction (fn *CXFunction) *CXModule {
	fn.Program = mod.Program
	fn.Module = mod
	
	found := false
	for i, f := range mod.Functions {
		if f.Name == fn.Name {
			mod.Functions[i].Name = fn.Name
			mod.Functions[i].Inputs = fn.Inputs
			mod.Functions[i].Outputs = fn.Outputs
			mod.Functions[i].Expressions = fn.Expressions
			mod.Functions[i].CurrentExpression = fn.CurrentExpression
			mod.Functions[i].Module = fn.Module
			mod.Functions[i].Program = fn.Program
			mod.CurrentFunction = mod.Functions[i]
			found = true
			break
		}
	}
	if !found {
		mod.Functions = append(mod.Functions, fn)
		mod.CurrentFunction = fn
	}

	return mod
}

func (mod *CXModule) AddStruct (strct *CXStruct) *CXModule {
	prgrm := mod.Program
	strct.Program = prgrm
	strct.Module = mod
	mod.CurrentStruct = strct
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

func (strct *CXStruct) AddField (fld *CXArgument) *CXStruct {
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
	expr.Program = fn.Program
	expr.Module = fn.Module
	expr.Function = fn
	fn.Expressions = append(fn.Expressions, expr)
	fn.CurrentExpression = expr
	return fn
}

func (fn *CXFunction) AddInput (param *CXArgument) *CXFunction {
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

func (fn *CXFunction) AddOutput (param *CXArgument) *CXFunction {
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

func (expr *CXExpression) AddOutput (param *CXArgument) *CXExpression {
	found := false
	for _, out := range expr.Outputs {
		if out.Name == param.Name {
			found = true
			break
		}
	}
	if !found {
		expr.Outputs = append(expr.Outputs, param)
	}

	return expr
}
