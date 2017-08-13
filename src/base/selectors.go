package base

import (
	"fmt"
	"errors"
)

func (cxt *CXContext) SelectModule (name string) (*CXModule, error) {
	prgrmStep := &CXProgramStep{
		Action: func(cxt *CXContext) {
			cxt.SelectModule(name)
		},
	}
	saveProgramStep(prgrmStep, cxt)
	
	var found *CXModule
	for _, mod := range cxt.Modules {
		if mod.Name == name {
			cxt.CurrentModule = mod
			found = mod
		}
	}

	if found == nil {
		return nil, errors.New(fmt.Sprintf("Module '%s' does not exist", name))
	}

	return found, nil
}

func (cxt *CXContext) SelectFunction (name string) (*CXFunction, error) {
	prgrmStep := &CXProgramStep{
		Action: func(cxt *CXContext) {
			cxt.SelectFunction(name)
		},
	}
	saveProgramStep(prgrmStep, cxt)
	
	mod, err := cxt.GetCurrentModule()
	if err == nil {
		return mod.SelectFunction(name)
	} else {
		return nil, err
	}
}

func (mod *CXModule) SelectFunction (name string) (*CXFunction, error) {
	prgrmStep := &CXProgramStep{
		Action: func(cxt *CXContext) {
			if mod, err := cxt.GetCurrentModule(); err == nil {
				mod.SelectFunction(name)
			}
		},
	}
	saveProgramStep(prgrmStep, mod.Context)
	
	var found *CXFunction
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

func (cxt *CXContext) SelectStruct (name string) (*CXStruct, error) {
	prgrmStep := &CXProgramStep{
		Action: func(cxt *CXContext) {
			cxt.SelectStruct(name)
		},
	}
	saveProgramStep(prgrmStep, cxt)
	
	mod, err := cxt.GetCurrentModule()
	if err == nil {
		return mod.SelectStruct(name)
	} else {
		return nil, err
	}
}

func (mod *CXModule) SelectStruct (name string) (*CXStruct, error) {
	prgrmStep := &CXProgramStep{
		Action: func(cxt *CXContext) {
			if mod, err := cxt.GetCurrentModule(); err == nil {
				mod.SelectStruct(name)
			}
		},
	}
	saveProgramStep(prgrmStep, mod.Context)

	var found *CXStruct
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

func (cxt *CXContext) SelectExpression (line int) (*CXExpression, error) {
	prgrmStep := &CXProgramStep{
		Action: func(cxt *CXContext) {
			cxt.SelectExpression(line)
		},
	}
	saveProgramStep(prgrmStep, cxt)

	mod, err := cxt.GetCurrentModule()
	if err == nil {
		return mod.SelectExpression(line)
	} else {
		return nil, err
	}
}

func (mod *CXModule) SelectExpression (line int) (*CXExpression, error) {
	prgrmStep := &CXProgramStep{
		Action: func(cxt *CXContext) {
			if mod, err := cxt.GetCurrentModule(); err == nil {
				mod.SelectExpression(line)
			}
		},
	}
	saveProgramStep(prgrmStep, mod.Context)
	fn, err := mod.GetCurrentFunction()
	if err == nil {
		return fn.SelectExpression(line)
	} else {
		return nil, err
	}
}

func (fn *CXFunction) SelectExpression (line int) (*CXExpression, error) {
	prgrmStep := &CXProgramStep{
		Action: func(cxt *CXContext) {
			if mod, err := cxt.GetCurrentModule(); err == nil {
				if fn, err := mod.GetCurrentFunction(); err == nil {
					fn.SelectExpression(line)
				}
			}
		},
	}
	saveProgramStep(prgrmStep, fn.Context)
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
