package base

import (
	"fmt"
	"errors"
	"math/rand"
	"time"
	"github.com/skycoin/skycoin/src/cipher/encoder"
)

func (prgrm *CXProgram) Run () error {
	rand.Seed(time.Now().UTC().UnixNano())
	
}

func (prgrm *CXProgram) Run (nCalls int) error {
	rand.Seed(time.Now().UTC().UnixNano())

	if prgrm.CallStack != nil && len(prgrm.CallStack) > 0 {
		// we resume the program
		var lastCall *CXCall
		var err error

		var untilEnd = false
		if nCalls < 1 {
			nCalls = 1 // so the for loop executes
			untilEnd = true
		}

		for !prgrm.Terminated && nCalls > 0 {
			lastCall = prgrm.CallStack[len(prgrm.CallStack) - 1]
			err = lastCall.call(1, callCounter, prgrm)
			if err != nil {
				return err
			}
			if !untilEnd {
				nCalls = nCalls - 1
			}
		}
	} else {
		// initialization and checking
		if mod, err := prgrm.SelectModule(MAIN_MOD); err == nil {
			if fn, err := mod.SelectFunction(MAIN_FUNC); err == nil {
				// main function
				state := make([]*CXArgument, 0, 20)
				mainCall := MakeCall(fn, state, nil, mod, mod.Program)
				
				prgrm.CallStack = append(prgrm.CallStack, mainCall)

				var lastCall *CXCall
				var err error

				var untilEnd = false
				if nCalls < 1 {
					nCalls = 1 // so the for loop executes
					untilEnd = true
				}
				
				for !prgrm.Terminated && nCalls > 0 {
					lastCall = prgrm.CallStack[len(prgrm.CallStack) - 1]
					err = lastCall.call(1, callCounter, prgrm)
					if err != nil {
						return err
					}
					if !untilEnd {
						nCalls = nCalls - 1
					}
				}
				return err
			}
		} else {
			return err
		}
	}
	return nil
}

var isTesting bool
var isErrorPresent bool

func (call *CXCall) call (nCalls, callCounter int, prgrm *CXProgram) error {
	//  add a counter here to pause
	if nCalls > 0 && callCounter >= nCalls {
		return nil
	}
	callCounter++

	// exceptions
	var exc bool
	var excError error

	if call.Line >= len(call.Operator.Expressions) || call.Line < 0 {
		// popping the stack
		prgrm.CallStack = prgrm.CallStack[:len(prgrm.CallStack) - 1]
		numOutputs := len(call.Operator.Outputs)
		for i, out := range call.Operator.Outputs {
			found := true
			for _, def := range call.State {
				/////////// throw error if output was not defined, or handle outputs from last expression
				if out.Name == def.Name {
					if call.ReturnAddress != nil {
						retName := call.ReturnAddress.Operator.Expressions[call.ReturnAddress.Line - 1].OutputNames[i].Name

						found := false
						for _, retDef := range call.ReturnAddress.State {
							if retDef.Name == retName {
								retDef.Value = def.Value
								found = true
								break
							}
						}
						if !found {
							def.Name = retName
							call.ReturnAddress.State = append(call.ReturnAddress.State, def)
						}

						found = true
						// break
						if i == numOutputs {
							return call.ReturnAddress.call(nCalls, callCounter)
						}
					} else {
						// no return address. should only be for main
						prgrm.Terminated = true
						prgrm.Outputs = append(prgrm.Outputs, def)
					}
				}
			}

			// this isn't complete yet
			if !found {
				return errors.New(fmt.Sprintf("'%s' output(s) not specified", call.Operator.Name))
			}
		}

		if call.ReturnAddress != nil {
			return call.ReturnAddress.call(nCalls, callCounter)
		} else {
			// no return address. should only be for main
			prgrm.Terminated = true
			//prgrm.Outputs = append(prgrm.Outputs, def)
		}
	} else {
		fn := call.Operator
		
		if expr, err := fn.GetExpression(call.Line); err == nil {
			
			// getting arguments
			argsRefs, _ := expr.GetArguments()

			argsCopy := make([]*CXArgument, len(argsRefs))
			//argNames := make([]string, len(argsRefs))

			if len(argsRefs) != len(expr.Operator.Inputs) {
				
				if len(argsRefs) == 1 {
					return errors.New(fmt.Sprintf("%s: %d: %s: expected %d arguments; %d was provided",
						expr.FileName, expr.FileLine, expr.Operator.Name, len(expr.Operator.Inputs), len(argsRefs)))
				} else {
					return errors.New(fmt.Sprintf("%s: %d: %s: expected %d arguments; %d were provided",
						expr.FileName, expr.FileLine, expr.Operator.Name, len(expr.Operator.Inputs), len(argsRefs)))
				}
			}
			
			// we don't want to modify by reference, we need to make copies
			for i := 0; i < len(argsRefs); i++ {
				
				if argsRefs[i].Typ == "ident" || argsRefs[i].Typ == "$ident" || argsRefs[i].Typ == "*ident" {
					var lookingFor string
					encoder.DeserializeRaw(*argsRefs[i].Value, &lookingFor)
					if arg, err := resolveIdent(lookingFor, call); err == nil {

						argsCopy[i] = arg
						if argsRefs[i].Typ == "$ident" {
							argsCopy[i].Typ = argsCopy[i].Typ[1:]
						}
						if argsRefs[i].Typ == "*ident" {
							argsCopy[i].Typ = "*" + argsCopy[i].Typ
						}
					} else {
						return errors.New(fmt.Sprintf("%s: %d: %s", expr.FileName, expr.FileLine, err.Error()))
					}
				} else {
					argsCopy[i] = argsRefs[i]
				}

				// checking if arguments types match with expressions required types
				if len(expr.Operator.Inputs) > 0 &&
					expr.Operator.Inputs[i].Typ !=
					argsCopy[i].Typ {
					return errors.New(fmt.Sprintf("%s: %d: %s: argument %d is type '%s'; expected type '%s'\n",
						expr.FileName, expr.FileLine, expr.Operator.Name, i+1, argsCopy[i].Typ, expr.Operator.Inputs[i].Typ))
				}
			}

			var opName string
			if expr.Operator != nil {
				opName = expr.Operator.Name
			} else {
				opName = "id" // return the same
			}

			isNative := false
			if _, ok := NATIVE_FUNCTIONS[opName]; ok {
				isNative = true
			}

			// check if struct array function
			if isNative {
				checkNative(opName, expr, call, &argsCopy, &exc, &excError)
				if exc && isTesting {
					isErrorPresent = true
				}
				if exc && !isTesting {
					fmt.Println()
					fmt.Println("Call's State:")
					for _, def := range call.State {
						isBasic := false
						for _, basic := range BASIC_TYPES {
							if basic == def.Typ {
								isBasic = true
								break
							}
						}

						if len(def.Name) > len(NON_ASSIGN_PREFIX) && def.Name[:len(NON_ASSIGN_PREFIX)] != NON_ASSIGN_PREFIX {
							if isBasic {
								fmt.Printf("%s:\t\t%s\n", def.Name, PrintValue(def.Name, def.Value, def.Typ, prgrm))
							} else {
								fmt.Println(def.Name)
								PrintValue(def.Name, def.Value, def.Typ, prgrm)
							}
						}
					}
					fmt.Println()
					fmt.Printf("%s() Arguments:\n", expr.Operator.Name)
					for i, arg := range argsCopy {
						fmt.Printf("%d: %s\n", i, PrintValue("", arg.Value, arg.Typ, prgrm))
					}
					fmt.Println()
					return excError
				}
				
				call.Line++
				return call.call(nCalls, callCounter)
			} else {
				// operator was not a native function
				if exc && isTesting {
					isErrorPresent = true
					//fmt.Println(excError)
				}
				if exc && !isTesting {
					fmt.Println()
					fmt.Println("Call's State:")
					for _, def := range call.State {
						fmt.Printf("%s:\t\t%s\n", def.Name, PrintValue(def.Name, def.Value, def.Typ, prgrm))
					}
					fmt.Println()
					fmt.Printf("%s() Arguments:\n", expr.Operator.Name)
					for i, arg := range argsCopy {
						fmt.Printf("%d: %s\n", i, PrintValue("", arg.Value, arg.Typ, prgrm))
					}
					fmt.Println()
					return excError
				}
				
				call.Line++ // once the subcall finishes, call next line
				if argDefs, err := argsToDefs(argsCopy, expr.Operator.Inputs, expr.Operator.Outputs, call.Module, prgrm); err == nil {
					subcall := MakeCall(expr.Operator, argDefs, call, call.Module, prgrm)

					prgrm.CallStack.Calls = append(prgrm.CallStack.Calls, subcall)
					return subcall.call(nCalls, callCounter)
				} else {
					fmt.Println(err)
				}
			}
		} else {
			fmt.Println(err)
		}
	}
	return nil
}
