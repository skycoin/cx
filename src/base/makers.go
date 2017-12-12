package base

import (
	"fmt"
	//"github.com/satori/go.uuid"
	"github.com/skycoin/skycoin/src/cipher/encoder"
)

var genSymCounter int = 0
func MakeGenSym (name string) string {
	gensym := fmt.Sprintf("%s_%d", name, genSymCounter)
	genSymCounter++
	
	return gensym
}

func MakeContext () *CXProgram {
	//heap := make([]byte, 0)
	newContext := &CXProgram{
		Modules: make([]*CXModule, 0),
		CallStack: MakeCallStack(0),
		//Heap: &heap,
		Steps: make([]*CXCallStack, 0)}
	return newContext
}

func MakeParameterCopy (param *CXParameter) *CXParameter {
	return &CXParameter{
		Name: param.Name,
		Typ: param.Typ,
	}
}

func MakeArgumentCopy (arg *CXArgument) *CXArgument {
	value := *arg.Value
	return &CXArgument{
		Typ: arg.Typ,
		Value: &value,
	}
}

func MakeExpressionCopy (expr *CXExpression, fn *CXFunction, mod *CXModule, cxt *CXProgram) *CXExpression {
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

func MakeFunctionCopy (fn *CXFunction, mod *CXModule, cxt *CXProgram) *CXFunction {
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
		Typ: fld.Typ,
	}
}

func MakeStructCopy (strct *CXStruct, mod *CXModule, cxt *CXProgram) *CXStruct {
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

func MakeDefinitionCopy (def *CXDefinition, mod *CXModule, cxt *CXProgram) *CXDefinition {
	valCopy := *def.Value
	return &CXDefinition{
		Name: def.Name,
		Typ: def.Typ,
		Value: &valCopy,
		Module: mod,
		Context: cxt,
	}
}

func MakeModuleCopy (mod *CXModule, cxt *CXProgram) *CXModule {
	newMod := &CXModule{Context: cxt}
	fnsCopy := make([]*CXFunction, len(mod.Functions))
	strctsCopy := make([]*CXStruct, len(mod.Structs))
	defsCopy := make([]*CXDefinition, len(mod.Definitions))
	
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

func MakeCallCopy (call *CXCall, mod *CXModule, cxt *CXProgram) *CXCall {
	stateCopy := make([]*CXDefinition, len(call.State))
	for k, v := range call.State {
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

func MakeContextCopy (cxt *CXProgram, stepNumber int) *CXProgram {
	newContext := &CXProgram{}

	modsCopy := make([]*CXModule, len(cxt.Modules))
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
		newContext.Steps = nil
	}
	
	return newContext
}

func MakeModule (name string) *CXModule {
	return &CXModule{
		Name: name,
		Definitions: make([]*CXDefinition, 0, 10),
		Imports: make([]*CXModule, 0),
		Functions: make([]*CXFunction, 0, 10),
		Structs: make([]*CXStruct, 0),
	}
}

func MakeDefinition (name string, value *[]byte, typ string) *CXDefinition {
	return &CXDefinition{
		Name: name,
		Typ: typ,
		Value: value,
	}
}

func MakeField (name string, typ string) *CXField {
	return &CXField{Name: name, Typ: typ}
}

func MakeFieldFromParameter (param *CXParameter) *CXField {
	return &CXField{Name: param.Name, Typ: param.Typ}
}

// Used only for native types
func MakeDefaultValue (typName string) *[]byte {
	var zeroVal []byte
	switch typName {
	case "str": zeroVal = encoder.Serialize("")
	case "bool": zeroVal = encoder.Serialize(int32(0))
	case "byte": zeroVal = []byte{byte(0)}
	case "i32": zeroVal = encoder.Serialize(int32(0))
	case "i64": zeroVal = encoder.Serialize(int64(0))
	case "f32": zeroVal = encoder.Serialize(float32(0))
	case "f64": zeroVal = encoder.Serialize(float64(0))
	case "[]byte": zeroVal = encoder.Serialize([]byte{})
	case "[]bool": zeroVal = encoder.Serialize([]int32{})
	case "[]str": zeroVal = encoder.Serialize([]string{})
	case "[]i32": zeroVal = encoder.Serialize([]int32{})
	case "[]i64": zeroVal = encoder.Serialize([]int64{})
	case "[]f32": zeroVal = encoder.Serialize([]float32{})
	case "[]f64": zeroVal = encoder.Serialize([]float64{})
	}
	return &zeroVal
}

func MakeStruct (name string) *CXStruct {
	return &CXStruct{Name: name}
}

func MakeParameter (name string, typ string) *CXParameter {
	return &CXParameter{Name: name,
		Typ: typ}
}

func MakeExpression (op *CXFunction) *CXExpression {
	return &CXExpression{Operator: op}
}

// var argPool = sync.Pool{
// 	New: func() interface{} {
// 		return &CXArgument{}
// 	},
// }

func MakeArgument (value *[]byte, typ string) *CXArgument {
	// arg := argPool.Get().(*CXArgument)
	// arg.Typ = typ
	// arg.Value = value
	// return arg
	
	return &CXArgument{
		Typ: typ,
		Value: value,
	}
}

func MakeFunction (name string) *CXFunction {
	return &CXFunction{Name: name}
}

func MakeValue (value string) *[]byte {
	byts := encoder.Serialize(value)
	return &byts
}

func MakeCall (op *CXFunction, state []*CXDefinition, ret *CXCall, mod *CXModule, cxt *CXProgram) *CXCall {
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

func MakeIdentityOpName (typeName string) string {
	switch typeName {
	case "str":
		return "str.id"
	case "bool":
		return "bool.id"
	case "byte":
		return "byte.id"
	case "i32":
		return "i32.id"
	case "i64":
		return "i64.id"
	case "f32":
		return "f32.id"
	case "f64":
		return "f64.id"
	case "[]bool":
		return "[]bool.id"
	case "[]byte":
		return "[]byte.id"
	case "[]str":
		return "[]str.id"
	case "[]i32":
		return "[]i32.id"
	case "[]i64":
		return "[]i64.id"
	case "[]f32":
		return "[]f32.id"
	case "[]f64":
		return "[]f64.id"
	default:
		return ""
	}
}
