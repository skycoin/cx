package base

import (
	// "bytes"
	// "errors"
	"fmt"
	"github.com/skycoin/skycoin/src/cipher/encoder"
	// "io/ioutil"
	"math/rand"
	// "runtime"
	"time"
)

func callsEqual(call1, call2 *CXCall) bool {
	if call1.Line != call2.Line ||
		len(call1.State) != len(call2.State) ||
		call1.Operator != call2.Operator ||
		call1.ReturnAddress != call2.ReturnAddress ||
		call1.Package != call2.Package {
		return false
	}

	for k, v := range call1.State {
		if call2.State[k] != v {
			return false
		}
	}

	return true
}

// It "un-runs" a program
func (prgrm *CXProgram) Reset() {
	prgrm.CallStack = MakeCallStack(0)
	prgrm.Steps = make([][]CXCall, 0)
	prgrm.Outputs = make([]*CXArgument, 0)
	//prgrm.ProgramSteps = nil
}

func (prgrm *CXProgram) UnRun(nCalls int) {
	if nCalls >= 0 || prgrm.CallCounter < 0 {
		return
	}

	call := &prgrm.CallStack[prgrm.CallCounter]
	
	for c := nCalls; c < 0; c++ {
		if call.Line >= c {
			// then we stay in this call counter
			call.Line += c
			c -= c
		} else {
			
			if prgrm.CallCounter == 0 {
				call.Line = 0
				return
			}
			c += call.Line
			call.Line = 0
			prgrm.CallCounter--
			call = &prgrm.CallStack[prgrm.CallCounter]
		}
	}
}

var isTesting bool
var isErrorPresent bool

func (prgrm *CXProgram) ToCall () *CXExpression {
	for c := prgrm.CallCounter - 1; c >= 0; c-- {
		if prgrm.CallStack[c].Line + 1 >= len(prgrm.CallStack[c].Operator.Expressions) {
			// then it'll also return from this function call; continue
			continue
		}
		return prgrm.CallStack[c].Operator.Expressions[prgrm.CallStack[c].Line + 1]
		// prgrm.CallStack[c].Operator.Expressions[prgrm.CallStack[prgrm.CallCounter-1].Line + 1]
	}
	// error
	return &CXExpression{Operator: MakeFunction("")}
	// panic("")
}

func (prgrm *CXProgram) RunCompiled(nCalls int) error {
	PROGRAM = prgrm
	// prgrm.PrintProgram()
	rand.Seed(time.Now().UTC().UnixNano())

	var untilEnd bool
	if nCalls == 0 {
		untilEnd = true
	}
	
	if mod, err := prgrm.SelectPackage(MAIN_PKG); err == nil {
		// initializing program resources
		// prgrm.Stacks = append(prgrm.Stacks, MakeStack(1024))

		if prgrm.CallStack[0].Operator == nil {
			// then the program is just starting and we need to run the SYS_INIT_FUNC
			if fn, err := mod.SelectFunction(SYS_INIT_FUNC); err == nil {
				// *init function
				mainCall := MakeCall(fn, nil, nil, mod, mod.Program)
				prgrm.CallStack[0] = mainCall
				prgrm.StackPointer = fn.Size

				var err error

				for !prgrm.Terminated {
					call := &prgrm.CallStack[prgrm.CallCounter]
					err = call.ccall(prgrm)
					if err != nil {
						return err
					}
				}
				// we reset call state
				prgrm.Terminated = false
				prgrm.CallCounter = 0
				prgrm.CallStack[0].Operator = nil
			} else {
				return err
			}
		}

		if fn, err := mod.SelectFunction(MAIN_FUNC); err == nil {
			if len(fn.Expressions) < 1 {
				return nil
			}

			if prgrm.CallStack[0].Operator == nil {
				// main function
				mainCall := MakeCall(fn, nil, nil, mod, mod.Program)
				// initializing program resources
				prgrm.CallStack[0] = mainCall
				// prgrm.Stacks = append(prgrm.Stacks, MakeStack(1024))
				prgrm.StackPointer = fn.Size

				prgrm.Terminated = false
			}

			var err error

			for !prgrm.Terminated && (untilEnd || nCalls != 0) {
				call := &prgrm.CallStack[prgrm.CallCounter]

				if !untilEnd {
					var inName string
					var toCallName string
					var toCall *CXExpression

					if call.Line >= call.Operator.Length && prgrm.CallCounter == 0 {
						prgrm.Terminated = true
						prgrm.CallStack[0].Operator = nil
						prgrm.CallCounter = 0
						fmt.Println("in:terminated")
						return err
					}

					if call.Line >= call.Operator.Length && prgrm.CallCounter != 0 {
						toCall = prgrm.ToCall()
						// toCall = prgrm.CallStack[prgrm.CallCounter-1].Operator.Expressions[prgrm.CallStack[prgrm.CallCounter-1].Line + 1]
						inName = prgrm.CallStack[prgrm.CallCounter-1].Operator.Name
					} else {
						toCall = call.Operator.Expressions[call.Line]
						inName = call.Operator.Name
					}

					if toCall.Operator == nil {
						// then it's a declaration
						toCallName = "declaration"
					} else if toCall.Operator.IsNative {
						toCallName = OpNames[toCall.Operator.OpCode]
					} else {
						if toCall.Operator.Name != "" {
							toCallName = toCall.Operator.Package.Name + "." + toCall.Operator.Name
						} else {
							// then it's the end of the program got from nested function calls
							prgrm.Terminated = true
							prgrm.CallStack[0].Operator = nil
							prgrm.CallCounter = 0
							fmt.Println("in:terminated")
							return err
						}
					}
					
					fmt.Printf("in:%s, expr#:%d, calling:%s()\n", inName, call.Line + 1, toCallName)
					
					nCalls--
				}
				
				err = call.ccall(prgrm)
				if err != nil {
					return err
				}
			}

			if prgrm.Terminated {
				prgrm.Terminated = false
				prgrm.CallCounter = 0
				prgrm.CallStack[0].Operator = nil
			}

			// debugging memory
			if len(prgrm.Memory) < 2000 {
				fmt.Println("prgrm.Memory", prgrm.Memory)
			}

			return err
		} else {
			return err
		}
	} else {
		return err
	}
}

func (call *CXCall) ccall(prgrm *CXProgram) error {
	// CX is still single-threaded, so only one stack
	if call.Line >= call.Operator.Length {
		/*
		   popping the stack
		*/
		// going back to the previous call
		prgrm.CallCounter--
		if prgrm.CallCounter < 0 {
			// then the program finished
			prgrm.Terminated = true
		} else {
			// copying the outputs to the previous stack frame
			returnAddr := &prgrm.CallStack[prgrm.CallCounter]
			returnOp := returnAddr.Operator
			returnLine := returnAddr.Line
			returnFP := returnAddr.FramePointer
			fp := call.FramePointer

			
			expr := returnOp.Expressions[returnLine]
			for i, out := range expr.Outputs {
				WriteMemory(
					GetFinalOffset(returnFP, out),
					ReadMemory(
						GetFinalOffset(fp, call.Operator.Outputs[i]),
						call.Operator.Outputs[i]))
			}

			// return the stack pointer to its previous state
			prgrm.StackPointer = call.FramePointer
			// we'll now execute the next command
			prgrm.CallStack[prgrm.CallCounter].Line++
			// calling the actual command
			prgrm.CallStack[prgrm.CallCounter].ccall(prgrm)
		}
	} else {
		/*
		   continue with call operator's execution
		*/
		fn := call.Operator
		expr := fn.Expressions[call.Line]
		// if it's a native, then we just process the arguments with execNative
		if expr.Operator == nil {
			// then it's a declaration
			call.Line++
		} else if expr.Operator.IsNative {
			execNative(prgrm)
			call.Line++
		} else {
			/*
			   It was not a native, so we need to create another call
			   with the current expression's operator
			*/
			// we're going to use the next call in the callstack
			prgrm.CallCounter++
			newCall := &prgrm.CallStack[prgrm.CallCounter]
			// setting the new call
			newCall.Operator = expr.Operator
			newCall.Line = 0
			newCall.FramePointer = prgrm.StackPointer
			// the stack pointer is moved to create room for the next call
			// prgrm.MemoryPointer += fn.Size
			prgrm.StackPointer += newCall.Operator.Size

			fp := call.FramePointer
			newFP := newCall.FramePointer

			// wiping next stack frame (removing garbage)
			for c := 0; c < expr.Operator.Size; c++ {
				prgrm.Memory[newFP+c] = 0
			}

			for i, inp := range expr.Inputs {
				var byts []byte
				// finalOffset := inp.Offset
				finalOffset := GetFinalOffset(fp, inp)
				// finalOffset := fp + inp.Offset

				// if inp.Indexes != nil {
				// 	finalOffset = GetFinalOffset(&prgrm.Stacks[0], fp, inp)
				// }
				if inp.PassBy == PASSBY_REFERENCE {
					byts = encoder.Serialize(int32(finalOffset))
				} else {
					byts = prgrm.Memory[finalOffset : finalOffset+inp.TotalSize]
				}

				// writing inputs to new stack frame
				WriteMemory(
					GetFinalOffset(newFP, newCall.Operator.Inputs[i]),
					// newFP + newCall.Operator.Inputs[i].Offset,
					// GetFinalOffset(prgrm.Memory, newFP, newCall.Operator.Inputs[i], MEM_WRITE),
					byts)
			}
		}
	}
	return nil
}
