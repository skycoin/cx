package base

import (
	"errors"
)

func (cxt *cxContext) SelectModule (name string) (*cxModule, error) {
	var found *cxModule
	for _, mod := range cxt.Modules {
		if mod.Name == name {
			cxt.CurrentModule = mod
			found = mod
		}
	}

	if found == nil {
		return nil, errors.New("Desired module does not exist")
	}

	return found, nil
}

func (cxt *cxContext) SelectFunction (name string) (*cxFunction, error) {
	mod, err := cxt.GetCurrentModule()
	if err == nil {
		return mod.SelectFunction(name)
	} else {
		return nil, err
	}
}

func (mod *cxModule) SelectFunction (name string) (*cxFunction, error) {
	var found *cxFunction
	for _, fn := range mod.Functions {
		if fn.Name == name {
			mod.CurrentFunction = fn
			found = fn
		}
	}

	if found == nil {
		return nil, errors.New("Desired function does not exist")
	}

	return found, nil
}

func (cxt *cxContext) SelectStruct (name string) (*cxStruct, error) {
	mod, err := cxt.GetCurrentModule()
	if err == nil {
		return mod.SelectStruct(name)
	} else {
		return nil, err
	}
}

func (mod *cxModule) SelectStruct (name string) (*cxStruct, error) {
	var found *cxStruct
	for _, strct := range mod.Structs {
		if strct.Name == name {
			mod.CurrentStruct = strct
			found = strct
		}
	}

	if found == nil {
		return nil, errors.New("Desired structure does not exist")
	}

	return found, nil
}

func (cxt *cxContext) SelectExpression (line int) (*cxExpression, error) {
	mod, err := cxt.GetCurrentModule()
	if err == nil {
		return mod.SelectExpression(line)
	} else {
		return nil, err
	}
}

func (mod *cxModule) SelectExpression (line int) (*cxExpression, error) {
	fn, err := mod.GetCurrentFunction()
	if err == nil {
		return fn.SelectExpression(line)
	} else {
		return nil, err
	}
}

func (fn *cxFunction) SelectExpression (line int) (*cxExpression, error) {
	if len(fn.Expressions) == 0 {
		return nil, errors.New("There are no expressions in this function")
	}
	
	if line >= len(fn.Expressions) {
		line = len(fn.Expressions) - 1
	}

	if line < 0 {
		line = 0
	}
	
	expr := fn.Expressions[line]
	fn.CurrentExpression = expr
	
	return expr, nil
}
