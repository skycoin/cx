package base

import (
	"fmt"
	"errors"
)

func (cxt *CXProgram) SelectModule (name string) (*CXPackage, error) {
	// prgrmStep := &CXProgramStep{
	// 	Action: func(cxt *CXProgram) {
	// 		cxt.SelectModule(name)
	// 	},
	// }
	// saveProgramStep(prgrmStep, cxt)
	
	var found *CXPackage
	for _, mod := range cxt.Packages {
		if mod.Name == name {
			cxt.CurrentPackage = mod
			found = mod
		}
	}

	if found == nil {
		return nil, errors.New(fmt.Sprintf("Module '%s' does not exist", name))
	}

	return found, nil
}

func (cxt *CXProgram) SelectFunction (name string) (*CXFunction, error) {
	// prgrmStep := &CXProgramStep{
	// 	Action: func(cxt *CXProgram) {
	// 		cxt.SelectFunction(name)
	// 	},
	// }
	// saveProgramStep(prgrmStep, cxt)
	
	mod, err := cxt.GetCurrentPackage()
	if err == nil {
		return mod.SelectFunction(name)
	} else {
		return nil, err
	}
}

func (mod *CXPackage) SelectFunction (name string) (*CXFunction, error) {
	// prgrmStep := &CXProgramStep{
	// 	Action: func(cxt *CXProgram) {
	// 		if mod, err := cxt.GetCurrentPackage(); err == nil {
	// 			mod.SelectFunction(name)
	// 		}
	// 	},
	// }
	// saveProgramStep(prgrmStep, mod.Context)
	
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

func (cxt *CXProgram) SelectStruct (name string) (*CXStruct, error) {
	// prgrmStep := &CXProgramStep{
	// 	Action: func(cxt *CXProgram) {
	// 		cxt.SelectStruct(name)
	// 	},
	// }
	// saveProgramStep(prgrmStep, cxt)
	
	mod, err := cxt.GetCurrentPackage()
	if err == nil {
		return mod.SelectStruct(name)
	} else {
		return nil, err
	}
}

func (mod *CXPackage) SelectStruct (name string) (*CXStruct, error) {
	// prgrmStep := &CXProgramStep{
	// 	Action: func(cxt *CXProgram) {
	// 		if mod, err := cxt.GetCurrentPackage(); err == nil {
	// 			mod.SelectStruct(name)
	// 		}
	// 	},
	// }
	// saveProgramStep(prgrmStep, mod.Context)

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

func (cxt *CXProgram) SelectExpression (line int) (*CXExpression, error) {
	// prgrmStep := &CXProgramStep{
	// 	Action: func(cxt *CXProgram) {
	// 		cxt.SelectExpression(line)
	// 	},
	// }
	// saveProgramStep(prgrmStep, cxt)

	mod, err := cxt.GetCurrentPackage()
	if err == nil {
		return mod.SelectExpression(line)
	} else {
		return nil, err
	}
}

func (mod *CXPackage) SelectExpression (line int) (*CXExpression, error) {
	// prgrmStep := &CXProgramStep{
	// 	Action: func(cxt *CXProgram) {
	// 		if mod, err := cxt.GetCurrentPackage(); err == nil {
	// 			mod.SelectExpression(line)
	// 		}
	// 	},
	// }
	// saveProgramStep(prgrmStep, mod.Context)
	fn, err := mod.GetCurrentFunction()
	if err == nil {
		return fn.SelectExpression(line)
	} else {
		return nil, err
	}
}

func (fn *CXFunction) SelectExpression (line int) (*CXExpression, error) {
	// prgrmStep := &CXProgramStep{
	// 	Action: func(cxt *CXProgram) {
	// 		if mod, err := cxt.GetCurrentPackage(); err == nil {
	// 			if fn, err := mod.GetCurrentFunction(); err == nil {
	// 				fn.SelectExpression(line)
	// 			}
	// 		}
	// 	},
	// }
	// saveProgramStep(prgrmStep, fn.Context)
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
