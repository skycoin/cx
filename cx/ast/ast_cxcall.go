package ast

import (
	"github.com/skycoin/cx/cx/constants"
)

// CXCall ...
type CXCall struct {
	Operator     *CXFunction // What CX function will be called when running this CXCall in the runtime
	Line         int         // What line in the CX function is currently being executed
	FramePointer int         // Where in the stack is this function call's local variables stored
}

//function is only called once and by affordances
//is a function on CXCal, not PROGRAM
func (call *CXCall) Ccall(prgrm *CXProgram, globalInputs *[]CXValue, globalOutputs *[]CXValue) error {
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

			lenOuts := len(expr.Outputs)
			for i, out := range call.Operator.Outputs {
				// Continuing if there is no receiving variable available.
				if i >= lenOuts {
					continue
				}
				WriteMemory(
					GetFinalOffset(returnFP, expr.Outputs[i]),
					ReadMemory(
						GetFinalOffset(fp, out),
						out))
			}

			// return the stack pointer to its previous state
			prgrm.StackPointer = call.FramePointer
			// we'll now execute the next command
			prgrm.CallStack[prgrm.CallCounter].Line++
			// calling the actual command
			// prgrm.CallStack[prgrm.CallCounter].Ccall(prgrm)
		}
	} else {
		/*
		   continue with call operator's execution
		*/
		fn := call.Operator
		expr := fn.Expressions[call.Line]

		// if it's a native, then we just process the arguments with execNative

		//TODO: WHEN WOULD OPERATOR EVER BE NIL?
		if expr.Operator == nil {
			// then it's a declaration
			// wiping this declaration's memory (removing garbage)
			newCall := &prgrm.CallStack[prgrm.CallCounter]
			newFP := newCall.FramePointer
			size := GetSize(expr.Outputs[0])
			for c := 0; c < size; c++ {
				prgrm.Memory[newFP+expr.Outputs[0].Offset+c] = 0
			}
			call.Line++
		} else if expr.Operator.IsBuiltin {
			//TODO: SLICES ARE NON ATOMIC

			fp := call.FramePointer
			if IsOperator(expr.Operator.OpCode) {
				// TODO: resolve this at compile time
				atomicType := GetType(expr.Inputs[0])
				expr.Operator = GetTypedOperator(atomicType, expr.Operator.OpCode)
			}
			inputs := expr.Inputs
			inputCount := len(inputs)
			if inputCount > len(*globalInputs) {
				*globalInputs = make([]CXValue, inputCount)
			}
			inputValues := (*globalInputs)[:inputCount]

			outputs := expr.Outputs
			outputCount := len(outputs)
			if outputCount > len(*globalOutputs) {
				*globalOutputs = make([]CXValue, outputCount)
			}
			outputValues := (*globalOutputs)[:outputCount]

			argIndex := 0
			for inputIndex := 0; inputIndex < inputCount; inputIndex++ {
				input := inputs[inputIndex]
				offset := GetFinalOffset(fp, input)
				value := &inputValues[inputIndex]
				value.Arg = input
				//value.Used = -1
				value.Offset = offset
				value.Type = input.Type
				value.FramePointer = fp
				value.Expr = expr
				value.memory = PROGRAM.Memory[offset : offset+GetSize(input)]
				argIndex++
			}

			for outputIndex := 0; outputIndex < outputCount; outputIndex++ {
				output := outputs[outputIndex]
				offset := GetFinalOffset(fp, output)
				value := &outputValues[outputIndex]
				value.Arg = output
				//value.Used = -1
				value.Offset = offset
				value.Type = output.Type
				value.FramePointer = fp
				value.Expr = expr
				argIndex++
			}

			OpcodeHandlers[expr.Operator.OpCode](inputValues, outputValues)

			//inputValues[inputIndex].Used
			//WTF is ".Used"
			/*
				for inputIndex := 0; inputIndex < inputCount; inputIndex++ { // TODO: remove in release builds
					if inputValues[inputIndex].Used != int8(inputs[inputIndex].Type) { // TODO: remove cast
						panic(fmt.Sprintf("Input value not used for opcode: '%s', param #%d. Expected type %d, '%s', used type %d, '%s'.",
							OpNames[expr.Operator.OpCode],
							inputIndex+1,
							inputs[inputIndex].Type, constants.TypeNames[inputs[inputIndex].Type],
							inputValues[inputIndex].Used, constants.TypeNames[int(inputValues[inputIndex].Used)]))
					}
				}
			*/

			/*
				for outputIndex := 0; outputIndex < outputCount; outputIndex++ { // TODO: remove in release builds
					if outputValues[outputIndex].Used != int8(outputs[outputIndex].Type) { // TODO: remove cast
						panic(fmt.Sprintf("Output value not used for opcode: '%s', param #%d. Expected type %d, '%s', used type %d '%s'.",
							OpNames[expr.Operator.OpCode],
							outputIndex+1,
							outputs[outputIndex].Type, constants.TypeNames[outputs[outputIndex].Type],
							outputValues[outputIndex].Used, constants.TypeNames[int(outputValues[outputIndex].Used)]))
					}
				}
			*/

			call.Line++
		} else { //NON-ATOMIC OPERATOR
			//TODO: Is this only called for user defined functions?

			//panic("BULLSHIT")
			/*
			   It was not a native, so we need to create another call
			   with the current expression's operator
			*/
			// we're going to use the next call in the callstack
			prgrm.CallCounter++
			if prgrm.CallCounter >= constants.CALLSTACK_SIZE {
				panic(constants.STACK_OVERFLOW_ERROR)
			}
			newCall := &prgrm.CallStack[prgrm.CallCounter]
			// setting the new call
			newCall.Operator = expr.Operator
			newCall.Line = 0
			newCall.FramePointer = prgrm.StackPointer
			// the stack pointer is moved to create room for the next call
			// prgrm.MemoryPointer += fn.Size
			prgrm.StackPointer += newCall.Operator.Size

			// checking if enough memory in stack
			if prgrm.StackPointer > constants.STACK_SIZE {
				panic(constants.STACK_OVERFLOW_ERROR)
			}

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
				if inp.PassBy == constants.PASSBY_REFERENCE {
					// If we're referencing an inner element, like an element of a slice (&slc[0])
					// or a field of a struct (&struct.fld) we no longer need to add
					// the OBJECT_HEADER_SIZE to the offset
					if inp.IsInnerReference {
						finalOffset -= constants.OBJECT_HEADER_SIZE
					}
					var finalOffsetB [4]byte
					WriteMemI32(finalOffsetB[:], 0, int32(finalOffset))
					byts = finalOffsetB[:]
				} else {
					size := GetSize(inp)
					byts = prgrm.Memory[finalOffset : finalOffset+size]
				}

				// writing inputs to new stack frame
				WriteMemory(
					GetFinalOffset(newFP, newCall.Operator.Inputs[i]),
					// newFP + newCall.Operator.ProgramInput[i].Offset,
					// GetFinalOffset(prgrm.Memory, newFP, newCall.Operator.ProgramInput[i], MEM_WRITE),
					byts)
			}
		}
	}
	return nil
}

//prgrm.CallStack = MakeCallStack(0)
func MakeCallStack(size int) []CXCall {
	return make([]CXCall, 0)
	// return &CXCallStack{
	// 	Calls: make([]*CXCall, size),
	// }
}
