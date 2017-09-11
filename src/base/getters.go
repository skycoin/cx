package base

import (
	"fmt"
	"errors"
)

func (cxt *CXProgram) GetCurrentModule () (*CXModule, error) {
	if cxt.CurrentModule != nil {
		return cxt.CurrentModule, nil
	} else {
		return nil, errors.New("Current module is nil")
	}
	
}

func (cxt *CXProgram) GetCurrentStruct () (*CXStruct, error) {
	if cxt.CurrentModule != nil &&
		cxt.CurrentModule.CurrentStruct != nil {
		return cxt.CurrentModule.CurrentStruct, nil
	} else {
		return nil, errors.New("Current module or struct is nil")
	}
	
}

func (mod *CXModule) GetCurrentStruct () (*CXStruct, error) {
	if mod.CurrentStruct != nil {
		return mod.CurrentStruct, nil
	} else {
		return nil, errors.New("Current struct is nil")
	}
	
}

func (cxt *CXProgram) GetCurrentFunction () (*CXFunction, error) {
	if cxt.CurrentModule != nil &&
		cxt.CurrentModule.CurrentFunction != nil {
		return cxt.CurrentModule.CurrentFunction, nil
	} else {
		return nil, errors.New("Current module or function is nil")
	}
	
}

func (mod *CXModule) GetCurrentFunction () (*CXFunction, error) {
	if mod.CurrentFunction != nil {
		return mod.CurrentFunction, nil
	} else {
		return nil, errors.New("Current function is nil")
	}
}

func (cxt *CXProgram) GetCurrentExpression () (*CXExpression, error) {
	if cxt.CurrentModule != nil &&
		cxt.CurrentModule.CurrentFunction != nil &&
		cxt.CurrentModule.CurrentFunction.CurrentExpression != nil {
		return cxt.CurrentModule.CurrentFunction.CurrentExpression, nil
	} else {
		return nil, errors.New("Current module, function or expression is nil")
	}
}

func (fn *CXFunction) GetCurrentExpression () (*CXExpression, error) {
	if fn.CurrentExpression != nil {
		return fn.CurrentExpression, nil
	} else if fn.Expressions != nil {
		return fn.Expressions[0], nil
	} else {
		return nil, errors.New("Current expression is nil")
	}
}

func (cxt *CXProgram) GetCurrentDefinitions () ([]*CXDefinition, error) {
	mod, err := cxt.GetCurrentModule()

	if err == nil {
		return mod.GetCurrentDefinitions()
	} else {
		return nil, err
	}
}

func (mod *CXModule) GetCurrentDefinitions () ([]*CXDefinition, error) {
	return mod.GetDefinitions()
}

func (mod *CXModule) GetDefinitions () ([]*CXDefinition, error) {
	if mod.Definitions != nil {
		return mod.Definitions, nil
	} else {
		return nil, errors.New("Definitions array is nil")
	}
}

func (cxt *CXProgram) GetDefinition (name string) (*CXDefinition, error) {
	if mod, err := cxt.GetCurrentModule(); err == nil {
		var found *CXDefinition
		for _, def := range mod.Definitions {
			if def.Name == name {
				found = def
				break
			}
		}
		
		if found == nil {
			return nil, errors.New(fmt.Sprintf("GetDefinition: definition '%s' not found", name))
		} else {
			return found, nil
		}
	} else {
		return nil, err
	}
}

func (strct *CXStruct) GetFields() ([]*CXField, error) {
	if strct.Fields != nil {
		return strct.Fields, nil
	} else {
		return nil, errors.New("Structure has no fields")
	}
}

func (mod *CXModule) GetFunctions() ([]*CXFunction, error) {
	// going from map to slice
	if mod.Functions != nil {
		// funcs := make([]*CXFunction, len(mod.Functions))
		// i := 0
		// for _, v := range mod.Functions {
		// 	funcs[i] = v
		// 	i++
		// }
		return mod.Functions, nil
	} else {
		return nil, errors.New("Module has no functions")
	}
}

func (cxt *CXProgram) GetModule (modName string) (*CXModule, error) {
	if cxt.Modules != nil {
		var found *CXModule
		for _, mod := range cxt.Modules {
			if modName == mod.Name {
				found = mod
				break
			}
		}
		return found, nil
	} else {
		return nil, errors.New(fmt.Sprintf("Module '%s'", modName))
	}
}

func (cxt *CXProgram) GetStruct (strctName string, modName string) (*CXStruct, error) {
	var foundMod *CXModule
	for _, mod := range cxt.Modules {
		if modName == mod.Name {
			foundMod = mod
			break
		}
	}
	var foundStrct *CXStruct
	for _, strct := range foundMod.Structs {
		if strct.Name == strctName {
			foundStrct = strct
			break
		}
	}

	if foundMod != nil && foundStrct != nil {
		return foundStrct, nil
	} else {
		return nil, errors.New(fmt.Sprintf("Strct '%s' not found in module '%s'", strctName, modName))
	}
	
	// if cxt.Modules != nil && cxt.Modules[modName] != nil && cxt.Modules[modName].Structs != nil && cxt.Modules[modName].Structs[strctName] != nil {
	// 	return cxt.Modules[modName].Structs[strctName], nil
	// } else {
	// 	return nil, errors.New(fmt.Sprintf("Strct '%s' not found in module '%s'", strctName, modName))
	// }
}

func (mod *CXModule) GetDefinition (defName string) (*CXDefinition, error) {
	var foundDef *CXDefinition
	for _, def := range mod.Definitions {
		if def.Name == defName {
			foundDef = def
			break
		}
	}

	if foundDef != nil {
		return foundDef, nil
	} else {
		return nil, errors.New(fmt.Sprintf("Definition '%s' not found in module '%s'", defName, mod.Name))
	}
}

func (cxt *CXProgram) GetFunction (fnName string, modName string) (*CXFunction, error) {
	for _, nativeFn := range NATIVE_FUNCTIONS {
		if fnName == nativeFn {
			modName = CORE_MODULE
			break
		}
	}
	
	var foundMod *CXModule
	for _, mod := range cxt.Modules {
		if modName == mod.Name {
			foundMod = mod
			break
		}
	}

	var foundFn *CXFunction
	for _, fn := range foundMod.Functions {
		if fn.Name == fnName {
			foundFn = fn
			break
		}
	}

	if foundMod != nil && foundFn != nil {
		return foundFn, nil
	} else {
		return nil, errors.New(fmt.Sprintf("Function '%s' not found in module '%s'", fnName, modName))
	}
	
	// if cxt.Modules != nil && cxt.Modules[modName] != nil && cxt.Modules[modName].Functions != nil && cxt.Modules[modName].Functions[fnName] != nil {
	// 	return cxt.Modules[modName].Functions[fnName], nil
	// } else {
	// 	return nil, errors.New(fmt.Sprintf("Function '%s' not found in module '%s'", fnName, modName))
	// }
}

func (fn *CXFunction) GetExpressions () ([]*CXExpression, error) {
	if fn.Expressions != nil {
		return fn.Expressions, nil
	} else {
		return nil, errors.New("Function has no expressions")
	}
}

func (fn *CXFunction) GetExpression (line int) (*CXExpression, error) {
	if fn.Expressions != nil {
		if line <= len(fn.Expressions) {
			return fn.Expressions[line], nil
		} else {
			return nil, errors.New("Expression line number exceeds number of expressions in function")
		}
		
	} else {
		return nil, errors.New("Function has no expressions")
	}
}

func (expr *CXExpression) GetArguments () ([]*CXArgument, error) {
	if expr.Arguments != nil {
		return expr.Arguments, nil
	} else {
		return nil, errors.New("Expression has no arguments")
	}
}

