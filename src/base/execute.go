package base

import (
	// "fmt"
	// "errors"
	"math/rand"
	"time"
	// "github.com/skycoin/skycoin/src/cipher/encoder"
)

func (prgrm *CXProgram) Run () error {
	rand.Seed(time.Now().UTC().UnixNano())

	if mod, err := prgrm.SelectModule(MAIN_MOD); err == nil {
		if fn, err := mod.SelectFunction(MAIN_FUNC); err == nil {
			// main function
			mainCall := MakeCall(fn, nil, mod, mod.Program)
			
			// initializing program resources
			//prgrm.CallStack = append(prgrm.CallStack, mainCall)
			prgrm.CallStack[0] = mainCall
			prgrm.Stacks = append(prgrm.Stacks, MakeStack(1024))

			var lastCall CXCall
			var err error

			for !prgrm.Terminated {
				lastCall = prgrm.CallStack[prgrm.CallCounter]
				err = lastCall.call(prgrm)
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

var isTesting bool
var isErrorPresent bool

func execNative (prgrm *CXProgram) {
	call := &prgrm.CallStack[prgrm.CallCounter]
	stack := &prgrm.Stacks[0]
	expr := call.Operator.Expressions[call.Line]
	opCode := expr.Operator.OpCode
	fp := call.FramePointer
	
	switch opCode {
	case OP_IDENTITY:
		opIdentity(expr, stack, fp)
	case OP_ADD:
		opAdd(expr, stack, fp)
	case OP_SUB:
	case OP_MUL:
	case OP_DIV:
	case OP_ABS:
	case OP_MOD:
	case OP_POW:
	case OP_COS:
	case OP_SIN:
	case OP_BITAND:
	case OP_BITOR:
	case OP_BITXOR:
	case OP_BITCLEAR:
	case OP_BITSHL:
	case OP_BITSHR:
	case OP_PRINT:
	case OP_MAKE:
	case OP_READ:
	case OP_WRITE:
	case OP_LEN:
	case OP_CONCAT:
	case OP_APPEND:
	case OP_COPY:
	case OP_CAST:
	case OP_EQ:
	case OP_UNEQ:
	case OP_LT:
	case OP_GT:
	case OP_LTEQ:
	case OP_GTEQ:
	case OP_RAND:
	case OP_AND:
	case OP_OR:
	case OP_NOT:
	case OP_SLEEP:
	case OP_HALT:
	case OP_GOTO:
	case OP_REMCX:
	case OP_ADDCX:
	case OP_QUERY:
	case OP_EXECUTE:
	case OP_INDEX:
	case OP_NAME:
	case OP_EVOLVE:
	case OP_TEST_START:
	case OP_TEST_STOP:
	case OP_TEST_ERROR:
	case OP_TEST:
	case OP_TIME_UNIX:
	case OP_TIME_UNIXMILLI:
	case OP_TIME_UNIXNANO:
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
		// copying the outputs to the previous stack frame
		returnAddr := prgrm.CallStack[prgrm.CallCounter]
		returnOp := returnAddr.Operator
		returnLine := returnAddr.Line
		returnFP := returnAddr.FramePointer
		fp := call.FramePointer
		
		expr := returnOp.Expressions[returnLine]
		outOffset := 0
		for _, out := range expr.Outputs {
			// copy byte by byte to the previous stack frame
			for c := 0; c < out.Size; c++ {
				prgrm.Stacks[0].Stack[returnFP + out.Offset + c] =
					prgrm.Stacks[0].Stack[fp + expr.Operator.InputsSize + outOffset + c]
			}
			outOffset += out.Size
		}
		
		// return the stack pointer to its previous state
		prgrm.Stacks[0].StackPointer = call.FramePointer
		// we'll now execute the next command
		prgrm.CallStack[prgrm.CallCounter].Line++
		// calling the actual command
		prgrm.CallStack[prgrm.CallCounter].call(prgrm)
	} else {
		/*
                  continue with call operator's execution
                */
		fn := call.Operator
		expr := fn.Expressions[call.Line]
		/*
                  preparing inputs and outputs in the stack for the next function call
                */
		sp := prgrm.Stacks[0].StackPointer
		fp := call.FramePointer
		tmp := sp
		for _, inp := range expr.Inputs {
			// we write the input values for the next frame
			size := inp.Size
			offset := inp.Offset
			for c := 0; c < size; c++ {
				// we copy each byte outside of current frame
				prgrm.Stacks[0].Stack[tmp+c] = prgrm.Stacks[0].Stack[fp+offset+c]
			}
			tmp += size
		}
		for _, out := range expr.Outputs {
			// we make room to receive the outputs
			tmp += out.Size
		}
		// the stack pointer is moved to create room for the next call
		prgrm.Stacks[0].StackPointer += fn.Size
		// if it's a native, then we just process the arguments with execNative
		if expr.Operator.IsNative {
			execNative(prgrm)
		} else {
			/*
                          It was not a native, so we need to create another call
                          with the current expression's operator
                        */
			// once the subcall finishes, call next line
			call.Line++
			// we're going to use the next call in the callstack
			prgrm.CallCounter++
			newCall := prgrm.CallStack[prgrm.CallCounter]
			// setting the new call
			newCall.Operator = expr.Operator
			newCall.Line = 0
			newCall.FramePointer = prgrm.Stacks[0].StackPointer
		}
	}
	return nil
}
