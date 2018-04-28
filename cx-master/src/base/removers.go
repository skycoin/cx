package base

import (

)

func (cxt *CXProgram) RemoveModule (modName string) {
	lenMods := len(cxt.Packages)
	for i, mod := range cxt.Packages {
		if mod.Name == modName {
			if i == lenMods - 1 {
				cxt.Packages = cxt.Packages[:len(cxt.Packages) - 1]
			} else {
				cxt.Packages = append(cxt.Packages[:i], cxt.Packages[i+1:]...)
			}
			break
		}
	}
}

func (mod *CXPackage) RemoveDefinition (defName string) {
	lenDefs := len(mod.Definitions)
	for i, def := range mod.Definitions {
		if def.Name == defName {
			if i == lenDefs - 1 {
				mod.Definitions = mod.Definitions[:len(mod.Definitions) - 1]
			} else {
				mod.Definitions = append(mod.Definitions[:i], mod.Definitions[i+1:]...)
			}
			break
		}
	}
}

func (mod *CXPackage) RemoveFunction (fnName string) {
	lenFns := len(mod.Functions)
	for i, fn := range mod.Functions {
		if fn.Name == fnName {
			if i == lenFns - 1 {
				mod.Functions = mod.Functions[:len(mod.Functions) - 1]
			} else {
				mod.Functions = append(mod.Functions[:i], mod.Functions[i+1:]...)
			}
			break
		}
	}
}

func (mod *CXPackage) RemoveStruct (strctName string) {
	lenStrcts := len(mod.Structs)
	for i, strct := range mod.Structs {
		if strct.Name == strctName {
			if i == lenStrcts - 1 {
				mod.Structs = mod.Structs[:len(mod.Structs) - 1]
			} else {
				mod.Structs = append(mod.Structs[:i], mod.Structs[i+1:]...)
			}
			break
		}
	}
}

func (mod *CXPackage) RemoveImport (impName string) {
	lenImps := len(mod.Imports)
	for i, imp := range mod.Imports {
		if imp.Name == impName {
			if i == lenImps - 1 {
				mod.Imports = mod.Imports[:len(mod.Imports) - 1]
			} else {
				mod.Imports = append(mod.Imports[:i], mod.Imports[i+1:]...)
			}
			break
		}
	}
}

func (strct *CXStruct) RemoveField (fldName string) {
	if len(strct.Fields) > 0 {
		lenFlds := len(strct.Fields)
		for i, fld := range strct.Fields {
			if fld.Name == fldName {
				if i == lenFlds - 1 {
					strct.Fields = strct.Fields[:len(strct.Fields) - 1]
				} else {
					strct.Fields = append(strct.Fields[:i], strct.Fields[i+1:]...)
				}
				break
			}
		}
	}
}

func (fn *CXFunction) RemoveExpression (line int) {
	if len(fn.Expressions) > 0 {
		lenExprs := len(fn.Expressions)
		if line >= lenExprs - 1 || line < 0 {
			fn.Expressions = fn.Expressions[:len(fn.Expressions) - 1]
		} else {
			fn.Expressions = append(fn.Expressions[:line], fn.Expressions[line+1:]...)
		}
		for i, expr := range fn.Expressions {
			expr.Line = i
		}
	}
}

func (fn *CXFunction) RemoveInput (inpName string) {
	if len(fn.Inputs) > 0 {
		lenInps := len(fn.Inputs)
		for i, inp := range fn.Inputs {
			if inp.Name == inpName {
				if i == lenInps {
					fn.Inputs = fn.Inputs[:len(fn.Inputs) - 1]
				} else {
					fn.Inputs = append(fn.Inputs[:i], fn.Inputs[i+1:]...)
				}
				break
			}
		}
	}
}

func (fn *CXFunction) RemoveOutput (outName string) {
	if len(fn.Outputs) > 0 {
		lenOuts := len(fn.Outputs)
		for i, out := range fn.Outputs {
			if out.Name == outName {
				if i == lenOuts {
					fn.Outputs = fn.Outputs[:len(fn.Outputs) - 1]
				} else {
					fn.Outputs = append(fn.Outputs[:i], fn.Outputs[i+1:]...)
				}
				break
			}
		}
	}
}

func (expr *CXExpression) RemoveArgument () {
	if len(expr.Arguments) > 0 {
		expr.Arguments = expr.Arguments[:len(expr.Arguments) - 1]
	}
}

func (expr *CXExpression) RemoveOutputName () {
	if len(expr.OutputNames) > 0 {
		expr.OutputNames = expr.OutputNames[:len(expr.OutputNames) - 1]
	}
}
