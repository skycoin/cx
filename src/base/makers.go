package base

import (
	"fmt"
	"github.com/satori/go.uuid"
)

//var counter int = 0
func MakeGenSym (name string) string {
	u1 := uuid.NewV4()
	for i := 0; i < len(u1); i++ {
		if u1[i] == '-' {
			u1[i] = '_'
		}
	}
	gensym := fmt.Sprintf("%s_%s", name, u1)
	//counter++
	
	return gensym
}

func MakeContext () *CXContext {
	heap := make([]byte, 0)
	newContext := &CXContext{
		Modules: make(map[string]*CXModule, 0),
		CallStack: MakeCallStack(0),
		Heap: &heap,
		Steps: make([]*CXCallStack, 0)}
	return newContext
}

// after implementing structure stepping, we need to change this
// we'll need to make hard copies of everything

//at the moment we need to make a hard copy of the modules
// we'll need it later too, so let's do it

// complete
func MakeParameterCopy (param *CXParameter) *CXParameter {
	return &CXParameter{
		Name: param.Name,
		Typ: MakeType(param.Typ.Name),
	}
}

// complete
func MakeArgumentCopy (arg *CXArgument) *CXArgument {
	value := *arg.Value
	return &CXArgument{
		Typ: MakeType(arg.Typ.Name),
		Value: &value,
		Offset: arg.Offset,
		Size: arg.Size,
	}
}

// complete
func MakeExpressionCopy (expr *CXExpression, fn *CXFunction, mod *CXModule, cxt *CXContext) *CXExpression {
	argsCopy := make([]*CXArgument, len(expr.Arguments))
	for i, arg := range expr.Arguments {
		argsCopy[i] = MakeArgumentCopy(arg)
	}
	return &CXExpression{
		Operator: expr.Operator,
		Arguments: argsCopy,
		OutputNames: expr.OutputNames,
		Line: expr.Line,
		Function: fn,
		Module: mod,
		Context: cxt,
	}
}

func MakeFunctionCopy (fn *CXFunction, mod *CXModule, cxt *CXContext) *CXFunction {
	newFn := &CXFunction{}
	inputsCopy := make([]*CXParameter, len(fn.Inputs))
	outputsCopy := make([]*CXParameter, len(fn.Outputs))
	exprsCopy := make([]*CXExpression, len(fn.Expressions))
	for i, inp := range fn.Inputs {
		inputsCopy[i] = MakeParameterCopy(inp)
	}
	for i, out := range fn.Outputs {
		outputsCopy[i] = MakeParameterCopy(out)
	}
	
	for i, expr := range fn.Expressions {
		exprsCopy[i] = MakeExpressionCopy(expr, newFn, mod, cxt)
	}
	
	newFn.Name = fn.Name
	newFn.Inputs = inputsCopy
	newFn.Outputs = outputsCopy
	
	// if fn.Output != nil {
	// 	newFn.Output = MakeParameterCopy(fn.Output)
	// }
	newFn.Expressions = exprsCopy
	if len(exprsCopy) > 0 {
		newFn.CurrentExpression = exprsCopy[len(exprsCopy) - 1]
	}
	newFn.Module = mod
	newFn.Context = cxt

	return newFn
}

func MakeFieldCopy (fld *CXField) *CXField {
	return &CXField{
		Name: fld.Name,
		Typ: MakeType(fld.Typ.Name),
	}
}

func MakeStructCopy (strct *CXStruct, mod *CXModule, cxt *CXContext) *CXStruct {
	fldsCopy := make([]*CXField, len(strct.Fields))
	for i, fld := range strct.Fields {
		fldsCopy[i] = MakeFieldCopy(fld)
	}
	return &CXStruct{
		Name: strct.Name,
		Fields: fldsCopy,
		Module: mod,
		Context: cxt,
	}
}

func MakeDefinitionCopy (def *CXDefinition, mod *CXModule, cxt *CXContext) *CXDefinition {
	valCopy := *def.Value
	return &CXDefinition{
		Name: def.Name,
		Typ: MakeType(def.Typ.Name),
		Value: &valCopy,
		Offset: def.Offset,
		Size: def.Size,
		Module: mod,
		Context: cxt,
	}
}

func MakeModuleCopy (mod *CXModule, cxt *CXContext) *CXModule {
	newMod := &CXModule{Context: cxt}
	fnsCopy := make(map[string]*CXFunction, len(mod.Functions))
	strctsCopy := make(map[string]*CXStruct, len(mod.Structs))
	defsCopy := make(map[string]*CXDefinition, len(mod.Definitions))
	
	for k, fn := range mod.Functions {
		fnsCopy[k] = MakeFunctionCopy(fn, newMod, cxt)
	}
	for k, strct := range mod.Structs {
		strctsCopy[k] = MakeStructCopy(strct, newMod, cxt)
	}
	for k, def := range mod.Definitions {
		defsCopy[k] = MakeDefinitionCopy(def, newMod, cxt)
	}

	// Setting current function in copy
	for _, fn := range fnsCopy {
		if fn.Name == mod.CurrentFunction.Name {
			newMod.CurrentFunction = fn
		}
	}

	newMod.Name = mod.Name
	newMod.Imports = mod.Imports
	newMod.Functions = fnsCopy
	newMod.Structs = strctsCopy
	newMod.Definitions = defsCopy
	newMod.Context = cxt
	
	return newMod
}

func MakeCallCopy (call *CXCall, mod *CXModule, cxt *CXContext) *CXCall {
	stateCopy := make(map[string]*CXDefinition)
	for k, v := range call.State {
		//var valueCopy []byte = *v.Value
		//stateCopy[k] = MakeDefinition(v.Name, &valueCopy, MakeType(v.Typ.Name))
		stateCopy[k] = MakeDefinitionCopy(v, mod, cxt)
	}
	return &CXCall{
		Operator: call.Operator,
		Line: call.Line,
		State: stateCopy,
		ReturnAddress: call.ReturnAddress,
		Module: mod,
		Context: cxt,
	}
}

func MakeCallStack (size int) *CXCallStack {
	return &CXCallStack{
		Calls: make([]*CXCall, size),
	}
}

func MakeContextCopy (cxt *CXContext, stepNumber int) *CXContext {
	newContext := &CXContext{}

	newContext.Heap = cxt.Heap

	modsCopy := make(map[string]*CXModule, len(cxt.Modules))
	if stepNumber >= len(cxt.Steps) || stepNumber < 0 {
		stepNumber = len(cxt.Steps) - 1
	}
	
	for k, mod := range cxt.Modules {
		modsCopy[k] = MakeModuleCopy(mod, newContext)
	}

	// Setting current module in copy
	for _, mod := range modsCopy {
		if mod.Name == cxt.CurrentModule.Name {
			newContext.CurrentModule = mod
		}
	}

	newContext.Modules = modsCopy

	// Making imports copies
	for _, mod := range modsCopy {
		for impKey, _ := range mod.Imports {
			mod.Imports[impKey] =  modsCopy[impKey]
		}
	}

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

	if len(cxt.Steps) > 0 {
		reqStep := cxt.Steps[stepNumber]
		newStep := MakeCallStack(len(reqStep.Calls))
		
		var lastCall *CXCall
		for j, call := range reqStep.Calls {
			var callModule *CXModule
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
			lastCall = newCall
			newStep.Calls[j] = newCall
		}
		
		newContext.CallStack = newStep
		//newContext.CallStack = reqStep
		//newContext.Steps = append(newContext.Steps, newStep)
		newContext.Steps = nil
	}
	
	return newContext
}

func MakeModule (name string) *CXModule {
	return &CXModule{
		Name: name,
		Definitions: make(map[string]*CXDefinition, 0),
		Imports: make(map[string]*CXModule, 0),
		Functions: make(map[string]*CXFunction, 0),
		Structs: make(map[string]*CXStruct, 0),
		//Instances: make(map[string]*CXInstance, 0),
	}
}

// func MakeInstance (name string, typ *CXStruct) *CXInstance {
// 	// the instance definitions need to be initialized
// 	defs := make([]*CXDefinition, 0)
// 	for _, fld := range typ.Fields {
// 		byteArray := make([]byte, 0)
// 		defs = append(defs, &CXDefinition{
// 			Name: fld.Name,
// 			Typ: fld.Typ,
// 			Value: &byteArray,
// 		})
// 	}
	
// 	return &CXInstance{
// 		Name: name,
// 		Typ: MakeType(typ.Name),
// 		Definitions: defs,
// 	}
// }

func MakeDefinition (name string, value *[]byte, typ *CXType) *CXDefinition {
	return &CXDefinition{
		Name: name,
		Typ: typ,
		Value: value,
		Offset: -1,
		Size: -1,}
}

func MakeField (name string, typ *CXType) *CXField {
	return &CXField{Name: name, Typ: typ}
}

func MakeStruct (name string) *CXStruct {
	return &CXStruct{Name: name}
}

func MakeParameter (name string, typ *CXType) *CXParameter {
	return &CXParameter{Name: name,
		Typ: typ}
}

func MakeExpression (fn *CXFunction) *CXExpression {
	return &CXExpression{Operator: fn}
}

func MakeArgument (value *[]byte, typ *CXType) *CXArgument {
	return &CXArgument{
		Typ: typ,
		Value: value,
		Offset: -1,
		Size: -1,
	}
}

func MakeType (name string) *CXType {
	return &CXType{
		Name: name,
	}
}

func MakeFunction (name string) *CXFunction {
	return &CXFunction{Name: name}
}

func MakeValue (value string) *[]byte {
	byts := []byte(value)
	return &byts
}

func MakeCall (op *CXFunction, state map[string]*CXDefinition, ret *CXCall, mod *CXModule, cxt *CXContext) *CXCall {
	return &CXCall{
		Operator: op,
		Line: 0,
		State: state,
		ReturnAddress: ret,
		Module: mod,
		Context: cxt,}
}

func MakeAffordance (desc string, action func()) *CXAffordance {
	return &CXAffordance{
		Description: desc,
		Action: action,
	}
}
