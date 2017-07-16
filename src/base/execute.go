package base

import (
	"fmt"
	"errors"
	"strings"
	"github.com/skycoin/skycoin/src/cipher/encoder"
)

// might not be needed, check to delete
func argsToDefs (args []*cxArgument, names []string) (map[string]*cxDefinition, error) {
	if len(names) == len(args) {
		defs := make(map[string]*cxDefinition, 0)
		for i, arg := range args {
			defs[names[i]] = &cxDefinition{
				Name: names[i],
				Typ: arg.Typ,
				Value: arg.Value,
			}
		}
		return defs, nil
	} else {
		return nil, errors.New("Not enough definition names provided")
	}
}





// 1. Check if main module and function exists
// 2. Make root call, with null returnaddress
// 3. Just call() this root call
//
// 4. The call() method should get the first expression


func PrintCallStackOld (callStack []*cxCall) {
	lastCall := callStack[len(callStack) - 1]
	pluses := strings.Repeat("   ", len(callStack) - 1)
	minuses := strings.Repeat("---", len(callStack) - 1)

	if lastCall.Line < len(lastCall.Operator.Expressions) {
		fmt.Printf("%sEntering function: '%s', Line#: %d \n",
			pluses,
			lastCall.Operator.Name,
			lastCall.Line)
	} else {
		fmt.Printf("%sNow exiting '%s'\n",
			minuses,
			lastCall.Operator.Name)
	}

	fmt.Printf("%sState:\n", pluses)
	for _, v := range lastCall.State {
		var val int32
		encoder.DeserializeAtomic(*v.Value, &val)
		fmt.Printf("%s\t'%s': %d\n",
			pluses,
			v.Name,
			val)
	}
}

func PrintCallStack (callStack []*cxCall) {
	for i, call := range callStack {
		tabs := strings.Repeat("\t", i)
		fmt.Printf("%s%s %d, ", tabs, call.Operator.Name, call.Line)

		lenState := len(call.State)
		idx := 0
		for _, def := range call.State {
			var val int32
			encoder.DeserializeAtomic(*def.Value, &val)
			if idx == lenState - 1 {
				fmt.Printf("%s: %d", def.Name, val)
			} else {
				fmt.Printf("%s: %d, ", def.Name, val)
			}
			idx++
		}
		fmt.Println()
	}
	fmt.Println()
}

func callsEqual (call1, call2 *cxCall) bool {
	if call1.Line != call2.Line ||
		len(call1.State) != len(call2.State) ||
		call1.Operator != call2.Operator ||
		call1.ReturnAddress != call2.ReturnAddress ||
		call1.Module != call2.Module {
		return false
	}

	for k, v := range call1.State {
		if call2.State[k] != v {
			return false
		}
	}

	return true
}

func saveStep (call *cxCall) {
	lenCallStack := len(call.Context.CallStack.Calls)
	newStep := MakeCallStack(lenCallStack)

	if len(call.Context.Steps) < 1 {
		// First call, copy everything
		for i, call := range call.Context.CallStack.Calls {
			newStep.Calls[i] = MakeCallCopy(call, call.Module, call.Context)
		}

		call.Context.Steps = append(call.Context.Steps, newStep)
		return
	}
	
	lastStep := call.Context.Steps[len(call.Context.Steps) - 1]
	lenLastStep := len(lastStep.Calls)
	
	smallerLen := 0
	if lenLastStep < lenCallStack {
		smallerLen = lenLastStep
	} else {
		smallerLen = lenCallStack
	}
	
	// Everytime a call changes, we need to make a hard copy of it
	// If the call doesn't change, we keep saving a pointer to it

	for i, call := range call.Context.CallStack.Calls[:smallerLen] {
		if callsEqual(call, lastStep.Calls[i]) {
			// if they are equal
			// append reference
			newStep.Calls[i] = lastStep.Calls[i]
		} else {
			newStep.Calls[i] = MakeCallCopy(call, call.Module, call.Context)
		}
	}

	// sizes can be different. if this is the case, we hard copy the rest
	for i, call := range call.Context.CallStack.Calls[smallerLen:] {
		newStep.Calls[i + smallerLen] = MakeCallCopy(call, call.Module, call.Context)
	}
	
	call.Context.Steps = append(call.Context.Steps, newStep)
	return
}

// It "un-runs" a program
func (cxt *cxContext) Reset() {
	cxt.CallStack = MakeCallStack(0)
	cxt.Steps = make([]*cxCallStack, 0)
	//cxt.ProgramSteps = nil
}

func (cxt *cxContext) ResetTo(stepNumber int) {
	// if no steps, we do nothing. the program will run from step 0
	if len(cxt.Steps) > 0 {
		if stepNumber > len(cxt.Steps) {
			stepNumber = len(cxt.Steps) - 1
		}
		reqStep := cxt.Steps[stepNumber]

		newStep := MakeCallStack(len(reqStep.Calls))
		
		var lastCall *cxCall
		for j, call := range reqStep.Calls {
			newCall := MakeCallCopy(call, call.Module, call.Context)
			newCall.ReturnAddress = lastCall
			lastCall = newCall
			newStep.Calls[j] = newCall
		}

		cxt.CallStack = newStep
		cxt.Steps = cxt.Steps[:stepNumber]
	}
}

func (cxt *cxContext) Run (withDebug bool, nCalls int) {
	var callCounter int = 0
	// we are going to do this if the CallStack is empty
	if cxt.CallStack != nil && len(cxt.CallStack.Calls) > 0 {
		// we resume the program
		lastCall := cxt.CallStack.Calls[len(cxt.CallStack.Calls) - 1]
		
		lastCall.call(withDebug, nCalls, callCounter)
	} else {
		// initialization and checking
		if mod, err := cxt.SelectModule("main"); err == nil {
			if fn, err := mod.SelectFunction("main"); err == nil {
				// main function
				state := make(map[string]*cxDefinition)
				mainCall := MakeCall(fn, state, nil, mod, mod.Context)
				
				cxt.CallStack.Calls = append(cxt.CallStack.Calls, mainCall)
				// saveStep(mainCall)
				// if withDebug {
				// 	PrintCallStack(cxt.CallStack)
				// }

				//cal := MakeCall(fn, state, mainCall)
				mainCall.call(withDebug, nCalls, callCounter)
			}
			
		} else {
			fmt.Println(err)
		}
	}
}

func (call *cxCall) call(withDebug bool, nCalls, callCounter int) {
	//  add a counter here to pause
	if nCalls > 0 && callCounter >= nCalls {
		return
	}
	callCounter++

	saveStep(call)
	if withDebug {
		PrintCallStack(call.Context.CallStack.Calls)
	}
	
	if call.Line >= len(call.Operator.Expressions) {
		if call.ReturnAddress != nil {
			// popping the stack
			call.Context.CallStack.Calls = call.Context.CallStack.Calls[:len(call.Context.CallStack.Calls) - 1]
			outName := call.Operator.Output.Name

			// this one is for returning result
			returnName := call.ReturnAddress.Operator.Expressions[call.ReturnAddress.Line - 1].OutputName
			output := call.State[outName]
			if output == nil {
				outName := call.Operator.Expressions[len(call.Operator.Expressions) - 1].OutputName
				output = call.State[outName]
			}

			if output != nil {
				call.ReturnAddress.State[returnName] = MakeDefinition(returnName, output.Value, output.Typ)
			}
			
			call.ReturnAddress.call(withDebug, nCalls, callCounter)
			return
		}
	} else {
		fn := call.Operator
		
		globals := call.Module.Definitions
		state := call.State

		if expr, err := fn.GetExpression(call.Line); err == nil {
			// getting arguments
			outName := expr.OutputName
			argsRefs, _ := expr.GetArguments()
			argsCopy := make([]*cxArgument, len(argsRefs))
			argNames := make([]string, len(argsRefs))

			for i, inp := range expr.Operator.Inputs {
				argNames[i] = inp.Name
			}
			
			// we are modifying by reference, we need to make copies
			for i := 0; i < len(argsRefs); i++ {
				if argsRefs[i].Typ.Name == "ident" {
					lookingFor := string(*argsRefs[i].Value)

					local := state[lookingFor]
					global := globals[lookingFor]

					if (local == nil && global == nil) {
						panic(fmt.Sprintf("'%s' is undefined", lookingFor))
					}

					// giving priority to local var
					if local != nil {
						argsCopy[i] = MakeArgument(local.Value, local.Typ)
					} else {
						argsCopy[i] = MakeArgument(global.Value, global.Typ)
					}
				} else {
					argsCopy[i] = argsRefs[i]
				}
			}

			// checking if native or not
			switch expr.Operator.Name {
			case "addI32":
				value := addI32(argsCopy[0], argsCopy[1])
				call.State[outName] = MakeDefinition(outName, value.Value, value.Typ)
				call.Line++
				call.call(withDebug, nCalls, callCounter)
				return
			case "mulI32":
				value := mulI32(argsCopy[0], argsCopy[1])
				call.State[outName] = MakeDefinition(outName, value.Value, value.Typ)
				call.Line++
				call.call(withDebug, nCalls, callCounter)
				return
			case "subI32":
				value := subI32(argsCopy[0], argsCopy[1])
				call.State[outName] = MakeDefinition(outName, value.Value, value.Typ)
				call.Line++
				call.call(withDebug, nCalls, callCounter)
				return
			default: // not native function
				call.Line++ // once the subcall finishes, call next line of the
				if argDefs, err := argsToDefs(argsCopy, argNames); err == nil {
					subcall := MakeCall(expr.Operator, argDefs, call, call.Module, call.Context)
					call.Context.CallStack.Calls = append(call.Context.CallStack.Calls, subcall)
					// debugging
					// saveStep(call)
					// if withDebug {
					// 	PrintCallStack(call.Context.CallStack)
					// }
					subcall.call(withDebug, nCalls, callCounter)
					return
				} else {
					fmt.Println(err)
				}
			}
		} else {
			fmt.Println(err)
		}
	}
}
