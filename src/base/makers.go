package base

import (
	"fmt"
)

var counter int = 0
func MakeGenSym (name string) string {
	gensym := fmt.Sprintf("%s%d", name, counter)
	counter++
	
	return gensym
}

func MakeContext () *cxContext {
	return &cxContext{
		Modules: make(map[string]*cxModule, 0),
		CallStack: make([]*cxCall, 0),
		Steps: make([][]*cxCall, 0)}
}

// after implementing structure stepping, we need to change this
// we'll need to make hard copies of everything

//at the moment we need to make a hard copy of the modules
// we'll need it later too, so let's do it

// complete
func MakeParameterCopy (param *cxParameter) *cxParameter {
	return &cxParameter{
		Name: param.Name,
		Typ: MakeType(param.Typ.Name),
	}
}

// complete
func MakeArgumentCopy (arg *cxArgument) *cxArgument {
	value := *arg.Value
	return &cxArgument{
		Typ: MakeType(arg.Typ.Name),
		Value: &value,
	}
}

// should the operator be the same as the copied expression?

// complete
func MakeExpressionCopy (expr *cxExpression, fn *cxFunction, mod *cxModule, cxt *cxContext) *cxExpression {
	argsCopy := make([]*cxArgument, len(expr.Arguments))
	for i, arg := range expr.Arguments {
		argsCopy[i] = MakeArgumentCopy(arg)
	}
	return &cxExpression{
		Operator: expr.Operator,
		Arguments: argsCopy,
		OutputName: expr.OutputName,
		Line: expr.Line,
		Function: fn,
		Module: mod,
		Context: cxt,
	}
}

func MakeFunctionCopy (fn *cxFunction, mod *cxModule, cxt *cxContext) *cxFunction {
	newFn := &cxFunction{}
	inputsCopy := make([]*cxParameter, len(fn.Inputs))
	exprsCopy := make([]*cxExpression, len(fn.Expressions))
	for i, inp := range fn.Inputs {
		inputsCopy[i] = MakeParameterCopy(inp)
	}
	for i, expr := range fn.Expressions {
		exprsCopy[i] = MakeExpressionCopy(expr, newFn, mod, cxt)
	}

	newFn.Name = fn.Name
	newFn.Inputs = inputsCopy
	if fn.Output != nil {
		newFn.Output = MakeParameterCopy(fn.Output)
	}
	newFn.Expressions = exprsCopy
	if fn.Expressions != nil {
		newFn.CurrentExpression = exprsCopy[len(exprsCopy) - 1]
	}
	newFn.Module = mod
	newFn.Context = cxt

	return newFn
}

func MakeFieldCopy (fld *cxField) *cxField {
	return &cxField{
		Name: fld.Name,
		Typ: MakeType(fld.Typ.Name),
	}
}

func MakeStructCopy (strct *cxStruct, mod *cxModule, cxt *cxContext) *cxStruct {
	fldsCopy := make([]*cxField, len(strct.Fields))
	for i, fld := range strct.Fields {
		fldsCopy[i] = MakeFieldCopy(fld)
	}
	return &cxStruct{
		Name: strct.Name,
		Fields: fldsCopy,
		Module: mod,
		Context: cxt,
	}
}

func MakeDefinitionCopy (def *cxDefinition, mod *cxModule, cxt *cxContext) *cxDefinition {
	valCopy := *def.Value
	return &cxDefinition{
		Name: def.Name,
		Typ: MakeType(def.Typ.Name),
		Value: &valCopy,
		Module: mod,
		Context: cxt,
	}
}

func MakeModuleCopy (mod *cxModule, cxt *cxContext) *cxModule {
	newMod := &cxModule{Context: cxt}
	fnsCopy := make(map[string]*cxFunction, len(mod.Functions))
	strctsCopy := make(map[string]*cxStruct, len(mod.Structs))
	defsCopy := make(map[string]*cxDefinition, len(mod.Definitions))
	
	for k, fn := range mod.Functions {
		fnsCopy[k] = MakeFunctionCopy(fn, newMod, cxt)
	}
	for k, strct := range mod.Structs {
		strctsCopy[k] = MakeStructCopy(strct, newMod, cxt)
	}
	for k, def := range mod.Definitions {
		defsCopy[k] = MakeDefinitionCopy(def, newMod, cxt)
	}

	newMod.Name = mod.Name
	newMod.Imports = mod.Imports
	newMod.Functions = fnsCopy
	newMod.Structs = strctsCopy
	newMod.Definitions = defsCopy
	newMod.Context = cxt
	
	return newMod
}

func MakeCallCopy (call *cxCall, mod *cxModule, cxt *cxContext) *cxCall {
	stateCopy := make(map[string]*cxDefinition)
	for k, v := range call.State {
		//var valueCopy []byte = *v.Value
		//stateCopy[k] = MakeDefinition(v.Name, &valueCopy, MakeType(v.Typ.Name))
		stateCopy[k] = MakeDefinitionCopy(v, mod, cxt)
	}
	return &cxCall{
		Operator: call.Operator,
		Line: call.Line,
		State: stateCopy,
		ReturnAddress: call.ReturnAddress,
		Module: mod,
		Context: cxt,
	}
}

func MakeContextCopy (cxt *cxContext, stepNumber int) *cxContext {
	newContext := &cxContext{}
	modsCopy := make(map[string]*cxModule, len(cxt.Modules))
	if stepNumber > len(cxt.Steps) || stepNumber < 0 {
		stepNumber = len(cxt.Steps) - 1
	}
	
	for k, mod := range cxt.Modules {
		modsCopy[k] = MakeModuleCopy(mod, newContext)
	}

	// Making imports copies
	for _, mod := range modsCopy {
		for impKey, _ := range mod.Imports {
			mod.Imports[impKey] =  modsCopy[impKey]
		}
	}

	newContext.Modules = modsCopy
	
	// Making expressions/operators
	for _, mod := range modsCopy {
		for _, fn := range mod.Functions {
			for _, expr := range fn.Expressions {
				if op, err := newContext.GetFunction(expr.Operator.Name, expr.Module.Name); err == nil {
					expr.Operator = op
				}
			}
		}
	}


	reqStep := cxt.Steps[stepNumber]

	newStep := make([]*cxCall, len(reqStep))
	var lastCall *cxCall
	for j, call := range reqStep {
		var callModule *cxModule
		for _, mod := range modsCopy {
			if call.Module.Name == mod.Name {
				callModule = mod
			}
		}
		
		newCall := MakeCallCopy(call, callModule, newContext)
		if callOp, err := newContext.GetFunction(call.Operator.Name, call.Operator.Module.Name); err == nil {
			newCall.Operator = callOp
		}
		newCall.ReturnAddress = lastCall
		fmt.Printf("**%p\n", newCall.Context)
		lastCall = newCall
		newStep[j] = newCall
	}

	
	newContext.CallStack = newStep
	//newContext.CallStack = reqStep
	//newContext.Steps = append(newContext.Steps, newStep)
	newContext.Steps = nil
	
	return newContext
}

func MakeContextCopy2 (cxt *cxContext, numberSteps int) *cxContext {
	if numberSteps > len(cxt.Steps) || numberSteps < 0 {
		numberSteps = len(cxt.Steps)
	}
	
	newContext := &cxContext{
		Modules: cxt.Modules,
		CurrentModule: cxt.CurrentModule,
	}

	for _, mod := range newContext.Modules {
		newMod := MakeModule(mod.Name)
		newMod.Context = newContext
		for _, fn := range mod.Functions {
			fn.Context = newContext
		}
		for _, strct := range mod.Structs {
			strct.Context = newContext
		}
		for _, def := range mod.Definitions {
			def.Context = newContext
		}
	}

	reqStep := cxt.Steps[numberSteps - 1]

	newStep := make([]*cxCall, len(reqStep))
	var lastCall *cxCall
	for j, call := range reqStep {
		newCall := MakeCallCopy(call, call.Module, newContext)
		//newCall.Context = newContext
		newCall.ReturnAddress = lastCall
		lastCall = newCall
		//saveStep(newCall)
		newStep[j] = newCall
	}

	newContext.Steps = nil
	newContext.CallStack = newStep
	
	//newSteps[i] = newStep

	//newSteps := make([][]*cxCall, numberSteps)
	// for i, step := range cxt.Steps {
	// 	if i >= numberSteps {
	// 		break
	// 	}
	// 	//newStep := make([]*cxCall, len(step))

	// 	// only the last step matters for now
		

	// 	var lastCall *cxCall
	// 	for _, call := range step {
	// 		newCall := MakeCallCopy(call)
	// 		newCall.Context = newContext
	// 		newCall.ReturnAddress = lastCall
	// 		lastCall = newCall
	// 		saveStep(newCall)
	// 		//newStep[j] = newCall
	// 	}
	// 	//newSteps[i] = newStep
	// }

	//fmt.Println(newSteps)
	
	//newContext.CallStack = newSteps[len(newSteps) - 1]
	//newContext.CallStack = newContext.Steps[len(newContext.Steps) - 1]
	//newContext.Steps = newSteps

	// for _, step := range newSteps {
	// 	fmt.Printf("%p\n", step[len(step) - 1].Context)
	// }
	
	return newContext
}

// func MakeModuleCopy (mod *cxModule) *cxModule {
// 	newMod := MakeModule(mod.Name)
// 	defCopies := make([]*cxDefinition, len(mod.Definitions))

// 	i := 0
// 	for _, def := range mod.Definitions {
// 		newDef := MakeDefinition(def.Name, def.Value, def.Typ)
// 		defCopies[i] = newDef
// 		i++
// 	}

// 	i = 0
// 	for _, fn := range mod.Functions {
// 		newFn := MakeFunction(fn.Name)
// 		newFn.Inputs = fn.Inputs
// 		newFn.Output = fn.Output
// 		newFn.Expressions = fn.Expressions
// 		newFn.CurrentExpression = fn.CurrentExpression
// 		//newFn.Context = 
		
// 	}

// 	// Import copies must be made in makecontextcopy I think
// 	// At the moment we can leave the references to the original imports
// 	newMod.Imports = mod.Imports
// 	// for _, imp := range mod.Imports {
// 	// 	MakeModule(def.Name, def.Value, def.Typ)
// 	// }

	
// }

func MakeModule (name string) *cxModule {
	return &cxModule{
		Name: name,
		Definitions: make(map[string]*cxDefinition, 0),
		Imports: make(map[string]*cxModule, 0),
		Functions: make(map[string]*cxFunction, 0),
		Structs: make(map[string]*cxStruct, 0),
	}
}

func MakeDefinition (name string, value *[]byte, typ *cxType) *cxDefinition {
	return &cxDefinition{Name: name, Typ: typ, Value: value}
}

func MakeField (name string, typ *cxType) *cxField {
	return &cxField{Name: name, Typ: typ}
}

func MakeStruct (name string) *cxStruct {
	return &cxStruct{Name: name}
}

func MakeParameter (name string, typ *cxType) *cxParameter {
	return &cxParameter{Name: name,
		Typ: typ}
}

func MakeExpression (outputName string, fn *cxFunction) *cxExpression {
	return &cxExpression{Operator: fn, OutputName: outputName}
}

func MakeArgument (value *[]byte, typ *cxType) *cxArgument {
	return &cxArgument{Typ: typ, Value: value}
}

func MakeType (name string) *cxType {
	return &cxType{Name: name}
}

func MakeFunction (name string) *cxFunction {
	return &cxFunction{Name: name}
}

func MakeValue (value string) *[]byte {
	byts := []byte(value)
	return &byts
}

func MakeCall (op *cxFunction, state map[string]*cxDefinition, ret *cxCall, mod *cxModule, cxt *cxContext) *cxCall {
	return &cxCall{
		Operator: op,
		Line: 0,
		State: state,
		ReturnAddress: ret,
		Module: mod,
		Context: cxt,}
}
