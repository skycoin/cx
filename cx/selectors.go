package base

import (
	"errors"
	"fmt"
)

func (mod *CXPackage) SelectFunction(name string) (*CXFunction, error) {
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
		return nil, fmt.Errorf("function '%s' does not exist", name)
	}

	return found, nil
}

func (mod *CXPackage) SelectStruct(name string) (*CXStruct, error) {
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

func (mod *CXPackage) SelectExpression(line int) (*CXExpression, error) {
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

func (fn *CXFunction) SelectExpression(line int) (*CXExpression, error) {
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
