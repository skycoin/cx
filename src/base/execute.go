package base

import (
	"fmt"
	// "errors"
	"math/rand"
	"time"
	// "github.com/skycoin/skycoin/src/cipher/encoder"
)

func (prgrm *CXProgram) Run () error {
	// prgrm.PrintProgram()
	rand.Seed(time.Now().UTC().UnixNano())

	if mod, err := prgrm.SelectPackage(MAIN_PKG); err == nil {
		if fn, err := mod.SelectFunction(MAIN_FUNC); err == nil {
			if len(fn.Expressions) < 1 {
				return nil
			}
			// main function
			mainCall := MakeCall(fn, nil, mod, mod.Program)
			
			// initializing program resources
			prgrm.CallStack[0] = mainCall
			prgrm.Stacks = append(prgrm.Stacks, MakeStack(1024))
			prgrm.Stacks[0].StackPointer = fn.Size

			var err error

			for !prgrm.Terminated {
				call := &prgrm.CallStack[prgrm.CallCounter]
				err = call.call(prgrm)
				if err != nil {
					return err
				}
			}
			return err
		} else {
			return err
		}
	} else {
		return err
	}
}

func (call *CXCall) call (prgrm *CXProgram) error {
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
			fmt.Println(prgrm.Stacks[0].Stack)
			// fmt.Println("prgrm.Data", prgrm.Data)
		} else {
			// copying the outputs to the previous stack frame
			returnAddr := &prgrm.CallStack[prgrm.CallCounter]
			returnOp := returnAddr.Operator
			returnLine := returnAddr.Line
			returnFP := returnAddr.FramePointer
			fp := call.FramePointer

			expr := returnOp.Expressions[returnLine]
			outOffset := 0
			for i, out := range expr.Outputs {
				// copy byte by byte to the previous stack frame
				for c := 0; c < out.TotalSize; c++ {
					prgrm.Stacks[0].Stack[returnFP + out.Offset + c] =
						prgrm.Stacks[0].Stack[fp + call.Operator.Outputs[i].Offset + c]
				}
				outOffset += out.TotalSize
			}

			// return the stack pointer to its previous state
			prgrm.Stacks[0].StackPointer = call.FramePointer
			// we'll now execute the next command
			prgrm.CallStack[prgrm.CallCounter].Line++
			// calling the actual command
			prgrm.CallStack[prgrm.CallCounter].call(prgrm)
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
			newCall.FramePointer = prgrm.Stacks[0].StackPointer
			// the stack pointer is moved to create room for the next call
			prgrm.Stacks[0].StackPointer += fn.Size

			fp := call.FramePointer
			newFP := newCall.FramePointer

			for i, inp := range expr.Inputs {
				var byts []byte
				finalOffset := inp.Offset
				if inp.Indexes != nil {
					finalOffset = GetFinalOffset(&prgrm.Stacks[0], fp, inp)
				}
				switch inp.MemoryType {
				case MEM_STACK:
					// byts = prgrm.Stacks[0].Stack[fp + inp.Offset : fp + inp.Offset + inp.TotalSize]
					byts = prgrm.Stacks[0].Stack[fp + finalOffset : fp + finalOffset + inp.TotalSize]
				case MEM_DATA:
					// byts = prgrm.Data[inp.Offset : inp.Offset + inp.TotalSize]
					byts = prgrm.Data[finalOffset : finalOffset + inp.TotalSize]
				default:
					panic("implement the other mem types")
				}
				// we copy the inputs for the next call
				for c := 0; c < inp.TotalSize; c++ {
					prgrm.Stacks[0].Stack[newFP + newCall.Operator.Inputs[i].Offset + c] = 
					byts[c]
				}
			}
		}
	}
	return nil
}
