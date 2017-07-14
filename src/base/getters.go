package base

import (
	"fmt"
	"errors"
)

func (cxt *cxContext) GetCurrentModule () (*cxModule, error) {
	if cxt.CurrentModule != nil {
		return cxt.CurrentModule, nil
	} else {
		return nil, errors.New("Current module is nil")
	}
	
}

func (cxt *cxContext) GetCurrentStruct () (*cxStruct, error) {
	if cxt.CurrentModule != nil &&
		cxt.CurrentModule.CurrentStruct != nil {
		return cxt.CurrentModule.CurrentStruct, nil
	} else {
		return nil, errors.New("Current module or struct is nil")
	}
	
}

func (cxt *cxContext) GetCurrentFunction () (*cxFunction, error) {
	if cxt.CurrentModule != nil &&
		cxt.CurrentModule.CurrentFunction != nil {
		return cxt.CurrentModule.CurrentFunction, nil
	} else {
		return nil, errors.New("Current module or function is nil")
	}
	
}

func (mod *cxModule) GetCurrentFunction () (*cxFunction, error) {
	if mod.CurrentFunction != nil {
		return mod.CurrentFunction, nil
	} else {
		return nil, errors.New("Current function is nil")
	}
}

func (cxt *cxContext) GetCurrentExpression () (*cxExpression, error) {
	if cxt.CurrentModule != nil &&
		cxt.CurrentModule.CurrentFunction != nil &&
		cxt.CurrentModule.CurrentFunction.CurrentExpression != nil {
		return cxt.CurrentModule.CurrentFunction.CurrentExpression, nil
	} else {
		return nil, errors.New("Current module, function or expression is nil")
	}
}

// no, we're always going to return something
// if nil, we return the first expression, unless it's empty
func (fn *cxFunction) GetCurrentExpression () (*cxExpression, error) {
	if fn.CurrentExpression != nil {
		return fn.CurrentExpression, nil
	} else if fn.Expressions != nil {
		return fn.Expressions[0], nil
	} else {
		return nil, errors.New("Current expression is nil")
	}
}

func (cxt *cxContext) GetCurrentDefinitions () (map[string]*cxDefinition, error) {
	mod, err := cxt.GetCurrentModule()

	if err == nil {
		return mod.GetCurrentDefinitions()
	} else {
		return nil, err
	}
}

func (mod *cxModule) GetCurrentDefinitions () (map[string]*cxDefinition, error) {
	return mod.GetDefinitions()
}

func (mod *cxModule) GetDefinitions () (map[string]*cxDefinition, error) {
	if mod.Definitions != nil {
		return mod.Definitions, nil
	} else {
		return nil, errors.New("Definitions array is nil")
	}
}

func (cxt *cxContext) GetDefinition (name string) (*cxDefinition, error) {
	if mod, err := cxt.GetCurrentModule(); err == nil {
		found := mod.Definitions[name]
		if found == nil {
			return nil, errors.New("Definition not found")
		} else {
			return found, nil
		}
	} else {
		return nil, err
	}
}

func (strct *cxStruct) GetFields() ([]*cxField, error) {
	if strct.Fields != nil {
		return strct.Fields, nil
	} else {
		return nil, errors.New("Structure has no fields")
	}
}

func (mod *cxModule) GetFunctions() ([]*cxFunction, error) {
	// going from map to slice
	if mod.Functions != nil {
		funcs := make([]*cxFunction, len(mod.Functions))
		i := 0
		for _, v := range mod.Functions {
			funcs[i] = v
			i++
		}
		return funcs, nil
	} else {
		return nil, errors.New("Module has no functions")
	}
}

func (cxt *cxContext) GetModule(modName string) (*cxModule, error) {
	if cxt.Modules != nil {
		return cxt.Modules[modName], nil
	} else {
		return nil, errors.New(fmt.Sprintf("Module '%s'", modName))
	}
}

func (cxt *cxContext) GetFunction(fnName string, modName string) (*cxFunction, error) {
	if cxt.Modules != nil && cxt.Modules[modName].Functions != nil {
		return cxt.Modules[modName].Functions[fnName], nil
	} else {
		return nil, errors.New(fmt.Sprintf("Function '%s' not found in module '%s'", fnName, modName))
	}
}

func (fn *cxFunction) GetExpressions() ([]*cxExpression, error) {
	if fn.Expressions != nil {
		return fn.Expressions, nil
	} else {
		return nil, errors.New("Function has no expressions")
	}
}

func (fn *cxFunction) GetExpression(line int) (*cxExpression, error) {
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

func (expr *cxExpression) GetArguments() ([]*cxArgument, error) {
	if expr.Arguments != nil {
		return expr.Arguments, nil
	} else {
		return nil, errors.New("Expression has no arguments")
	}
}
