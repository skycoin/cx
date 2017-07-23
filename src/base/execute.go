package base

import (
	"fmt"
	"errors"
	"strings"
	"github.com/skycoin/skycoin/src/cipher/encoder"
)

func argsToDefs (args []*cxArgument, names []string, mod *cxModule, cxt *cxContext) (map[string]*cxDefinition, error) {
	if len(names) == len(args) {
		defs := make(map[string]*cxDefinition, 0)
		for i, arg := range args {
			defs[names[i]] = &cxDefinition{
				Name: names[i],
				Typ: arg.Typ,
				Value: arg.Value,
				Module: mod,
				Context: cxt,
			}
		}
		return defs, nil
	} else {
		return nil, errors.New("Not enough definition names provided")
	}
}

func PrintCallStack (callStack []*cxCall) {
	for i, call := range callStack {
		tabs := strings.Repeat("\t", i)
		fmt.Printf("%s%s %d, ", tabs, call.Operator.Name, call.Line)

		lenState := len(call.State)
		idx := 0
		for _, def := range call.State {
			var valI32 int32
			var valI64 int64
			switch def.Typ.Name {
			case "i32":
				encoder.DeserializeAtomic(*def.Value, &valI32)
				if idx == lenState - 1 {
					fmt.Printf("%s: %d", def.Name, valI32)
				} else {
					fmt.Printf("%s: %d, ", def.Name, valI32)
				}
			case "i64":
				encoder.DeserializeAtomic(*def.Value, &valI64)
				if idx == lenState - 1 {
					fmt.Printf("%s: %d", def.Name, valI64)
				} else {
					fmt.Printf("%s: %d, ", def.Name, valI64)
				}
			case "byte":
				if idx == lenState - 1 {
					fmt.Printf("%s: %d", def.Name, (*def.Value)[0])
				} else {
					fmt.Printf("%s: %d, ", def.Name, (*def.Value)[0])
				}
			case "[]byte":
				if idx == lenState - 1 {
					fmt.Printf("%s: %v", def.Name, (*def.Value))
				} else {
					fmt.Printf("%s: %v, ", def.Name, (*def.Value))
				}
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

func beforeDot (str string ) (beforeDot string, afterDot string) {
	beforeDot = ""
	afterDot = ""
	beforeDotCounter := 0
	foundDot := false
	for _, letter := range str {
		if letter != '.' {
			beforeDot = concat(beforeDot, string(letter))
		} else {
			foundDot = true
			break
		}
		beforeDotCounter++
	}
	if foundDot {
		afterDot = str[beforeDotCounter + 1:] // ignore the dot
	}
	
	return beforeDot, afterDot
}

func getIdentParts (str string) []string {
	var identParts []string
	before, after := beforeDot(str) // => "Math", "myPoint.x"
	identParts = append(identParts, before)
	for after != "" {
		before, after = beforeDot(after)
		identParts = append(identParts, before)
	}
	return identParts
}

func (cxt *cxContext) Compile () *cxContext {
	heap := *cxt.Heap
	heapCounter := 0

	for _, mod := range cxt.Modules {
		for _, fn := range mod.Functions {

			// using struct fields
			for _, expr := range fn.Expressions {
				for argIdx, arg := range expr.Arguments {
					if arg.Typ.Name == "ident" {
						identParts := getIdentParts(string(*arg.Value))

						var def *cxDefinition
						
						if len(identParts) == 1 {
							// it's a current module's definition
							// TODO: compile later
						} else {
							if identMod, ok := cxt.Modules[identParts[0]]; ok {
								// identParts[1] will always be a definition if before is a module
								if len(identParts) == 2 {
									// then we are referring to the struct itself or a normal ident
									// TODO: compile later
								} else {
									// we're referring to a struct field then
									if len(identParts) > 3 {
										// nested structs
										// TODO: compile later
									} else {
										def = identMod.Definitions[concat(identParts[1], ".", identParts[2])]
										
										defSize := len(*def.Value)
										heapArg := &cxArgument{
											Typ: def.Typ,
											//Value: arg.Value, //irrelevant now
											Offset: heapCounter,
											Size: defSize,
										}
										// replacing argument
										arg = heapArg

										heap = append(heap, *def.Value...)
										heapCounter = heapCounter + defSize
									}
								}
							} else {
								if len(identParts) > 2 {
									fmt.Println("identParts > 2")
									// nested structs
									// TODO: compile later
								} else {
									def = mod.Definitions[concat(identParts[0], ".", identParts[1])]

									//defSize := encoder.Size(*def.Value)
									defSize := len(*def.Value)
									heapArg := &cxArgument{
										Typ: def.Typ,
										//Value: arg.Value, //irrelevant now
										Offset: heapCounter,
										Size: defSize,
									}
									// replacing argument
									expr.Arguments[argIdx] = heapArg

									//heap = append(heap, encoder.Serialize(*def.Value)...)
									heap = append(heap, *def.Value...)
									heapCounter = heapCounter + defSize
								}
							}
						}

						// checking if first part is a module
						
					}
				}
			}
		}
	}
	cxt.Heap = &heap
	return cxt
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
		if len(call.Operator.Expressions) < 1 {
			panic(fmt.Sprintf("Calling function without expressions '%s'", call.Operator.Name))
		}
		// popping the stack
		call.Context.CallStack.Calls = call.Context.CallStack.Calls[:len(call.Context.CallStack.Calls) - 1]
		outName := call.Operator.Output.Name

		// this one is for returning result
		output := call.State[outName]
		if output == nil {
			outName := call.Operator.Expressions[len(call.Operator.Expressions) - 1].OutputName
			output = call.State[outName]
		}

		// checking if output var has the same type as the required output
		if output.Typ.Name != call.Operator.Output.Typ.Name {
			panic(fmt.Sprintf("output var '%s' is of type '%s'; function '%s' requires output of type '%s'",
				output.Name, output.Typ.Name, call.Operator.Name, call.Operator.Output.Typ.Name))
		}

		if call.ReturnAddress != nil {
			returnName := call.ReturnAddress.Operator.Expressions[call.ReturnAddress.Line - 1].OutputName
			
			if output != nil {
				def := MakeDefinition(returnName, output.Value, output.Typ)
				def.Module = call.Module
				def.Context = call.Context
				call.ReturnAddress.State[returnName] = def
			} else {
				panic(fmt.Sprintf("Function '%s' couldn't return anything", call.Operator.Name))
			}

			call.ReturnAddress.call(withDebug, nCalls, callCounter)
		} else {
			// no return address. should only be for main
			call.Context.Output = output
			//fmt.Printf("\nProgram's output:\n%v\n", output.Value)
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
				if argsRefs[i].Offset > -1 {
					offset := argsRefs[i].Offset
					size := argsRefs[i].Size
					//var val []byte
					//encoder.DeserializeRaw((*call.Context.Heap)[offset:offset+size], &val)
					//argsRefs[i].Value = &val
					val := (*call.Context.Heap)[offset:offset+size]
					argsRefs[i].Value = &val
				}
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

				// checking if arguments types match with expressions required types
				if expr.Operator.Inputs[i].Typ.Name != argsCopy[i].Typ.Name {
					//panic(fmt.Sprintf("%s, line #%d: wrong argument type", expr.Operator.Name, call.Line))
					panic(fmt.Sprintf("%s, line #%d: %s argument #%d is type '%s'; expected type '%s'",
						fn.Name, call.Line, expr.Operator.Name, i, argsCopy[i].Typ.Name, expr.Operator.Inputs[i].Typ.Name))
				}
			}

			// checking if native or not
			var value *cxArgument
			
			switch expr.Operator.Name {
			case "addI32":
				value = addI32(argsCopy[0], argsCopy[1])
			case "mulI32":
				value = mulI32(argsCopy[0], argsCopy[1])
			case "subI32":
				value = subI32(argsCopy[0], argsCopy[1])
			case "divI32":
				value = divI32(argsCopy[0], argsCopy[1])
			case "addI64":
				value = addI64(argsCopy[0], argsCopy[1])
			case "mulI64":
				value = mulI64(argsCopy[0], argsCopy[1])
			case "subI64":
				value = subI64(argsCopy[0], argsCopy[1])
			case "divI64":
				value = divI64(argsCopy[0], argsCopy[1])
			case "readAByte":
				value = readAByte(argsCopy[0], argsCopy[1])
			case "writeAByte":
				value = writeAByte(argsCopy[0], argsCopy[1], argsCopy[2])
			case "":
			}
			if value != nil {
				// operator was a native function
				def := MakeDefinition(outName, value.Value, value.Typ)
				def.Module = call.Module
				def.Context = call.Context
				
				call.State[outName] = def
				call.Line++
				call.call(withDebug, nCalls, callCounter)
			} else {
				// operator was not a native function
				call.Line++ // once the subcall finishes, call next line of the
				if argDefs, err := argsToDefs(argsCopy, argNames, call.Module, call.Context); err == nil {
					subcall := MakeCall(expr.Operator, argDefs, call, call.Module, call.Context)
					call.Context.CallStack.Calls = append(call.Context.CallStack.Calls, subcall)
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
