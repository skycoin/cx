package cxcore

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/skycoin/skycoin/src/cipher/encoder"
)

// It "un-runs" a program
// func (prgrm *CXProgram) Reset() {
// 	prgrm.CallStack = MakeCallStack(0)
// 	prgrm.Steps = make([][]CXCall, 0)
// 	prgrm.Outputs = make([]*CXArgument, 0)
// 	//prgrm.ProgramSteps = nil
// }

// UnRun ...
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

// ToCall ...
func (prgrm *CXProgram) ToCall() *CXExpression {
	for c := prgrm.CallCounter - 1; c >= 0; c-- {
		if prgrm.CallStack[c].Line+1 >= len(prgrm.CallStack[c].Operator.Expressions) {
			// then it'll also return from this function call; continue
			continue
		}
		return prgrm.CallStack[c].Operator.Expressions[prgrm.CallStack[c].Line+1]
		// prgrm.CallStack[c].Operator.Expressions[prgrm.CallStack[prgrm.CallCounter-1].Line + 1]
	}
	// error
	return &CXExpression{Operator: MakeFunction("", "", -1)}
	// panic("")
}

// Run ...
func (prgrm *CXProgram) Run(untilEnd bool, nCalls *int, untilCall int) error {
	defer RuntimeError()
	var err error

    var inputs []CXValue
    var outputs []CXValue
	for !prgrm.Terminated && (untilEnd || *nCalls != 0) && prgrm.CallCounter > untilCall {
		call := &prgrm.CallStack[prgrm.CallCounter]

		// checking if enough memory in stack
		if prgrm.StackPointer > STACK_SIZE {
			panic(STACK_OVERFLOW_ERROR)
		}

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

			fmt.Printf("in:%s, expr#:%d, calling:%s()\n", inName, call.Line+1, toCallName)
			*nCalls--
		}

		err = call.ccall(prgrm, &inputs, &outputs)
		if err != nil {
			return err
		}
	}

	return nil
}

// minHeapSize determines what's the minimum heap size that a CX program
// needs to have based on INIT_HEAP_SIZE, MAX_HEAP_SIZE and NULL_HEAP_ADDRESS_OFFSET.
func minHeapSize() int {
	minHeapSize := INIT_HEAP_SIZE
	if MAX_HEAP_SIZE < INIT_HEAP_SIZE {
		// Then MAX_HEAP_SIZE overrides INIT_HEAP_SIZE's value.
		minHeapSize = MAX_HEAP_SIZE
	}
	if minHeapSize < NULL_HEAP_ADDRESS_OFFSET {
		// Then the user is trying to allocate too little heap memory.
		// We need at least NULL_HEAP_ADDRESS_OFFSET bytes for `nil`.
		minHeapSize = NULL_HEAP_ADDRESS_OFFSET
	}

	return minHeapSize
}

// EnsureHeap ensures that `prgrm` has `minHeapSize()`
// bytes allocated after the data segment.
func (prgrm *CXProgram) EnsureHeap() {
	currHeapSize := len(prgrm.Memory) - prgrm.HeapStartsAt
	minHeapSize := minHeapSize()
	if currHeapSize < minHeapSize {
		prgrm.Memory = append(prgrm.Memory, make([]byte, minHeapSize-currHeapSize)...)
	}
}

// RunCompiled ...
func (prgrm *CXProgram) RunCompiled(nCalls int, args []string) error {
	_, err := prgrm.SelectProgram()
	if err != nil {
		panic(err)
	}
	prgrm.EnsureHeap()
	rand.Seed(time.Now().UTC().UnixNano())

	var untilEnd bool
	if nCalls == 0 {
		untilEnd = true
	}
	mod, err := prgrm.SelectPackage(MAIN_PKG)
	if err == nil {
		// initializing program resources
		// prgrm.Stacks = append(prgrm.Stacks, MakeStack(1024))

        var inputs []CXValue
        var outputs []CXValue
		if prgrm.CallStack[0].Operator == nil {
			// then the program is just starting and we need to run the SYS_INIT_FUNC
			if fn, err := mod.SelectFunction(SYS_INIT_FUNC); err == nil {
				// *init function
				mainCall := MakeCall(fn)
				prgrm.CallStack[0] = mainCall
				prgrm.StackPointer = fn.Size

				var err error

				for !prgrm.Terminated {
					call := &prgrm.CallStack[prgrm.CallCounter]
                    err = call.ccall(prgrm, &inputs, &outputs)
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
				mainCall := MakeCall(fn)
				mainCall.FramePointer = prgrm.StackPointer
				// initializing program resources
				prgrm.CallStack[0] = mainCall

				// prgrm.Stacks = append(prgrm.Stacks, MakeStack(1024))
				prgrm.StackPointer += fn.Size

				// feeding os.Args
				if osPkg, err := PROGRAM.SelectPackage(OS_PKG); err == nil {
					argsOffset := 0
					if osGbl, err := osPkg.GetGlobal(OS_ARGS); err == nil {
						for _, arg := range args {
							argBytes := encoder.Serialize(arg)
							argOffset := AllocateSeq(len(argBytes) + OBJECT_HEADER_SIZE)

							var header = make([]byte, OBJECT_HEADER_SIZE)
							WriteMemI32(header, 5, int32(encoder.Size(arg)+OBJECT_HEADER_SIZE))
							obj := append(header, argBytes...)

							WriteMemory(argOffset, obj)

							var argOffsetBytes [4]byte
							WriteMemI32(argOffsetBytes[:], 0, int32(argOffset))
							argsOffset = WriteToSlice(argsOffset, argOffsetBytes[:])
						}
						WriteI32(GetFinalOffset(0, osGbl), int32(argsOffset))
					}
				}
				prgrm.Terminated = false
			}

			if err = prgrm.Run(untilEnd, &nCalls, -1); err != nil {
				return err
			}

			if prgrm.Terminated {
				prgrm.Terminated = false
				prgrm.CallCounter = 0
				prgrm.CallStack[0].Operator = nil
			}

			// debugging memory
			// if len(prgrm.Memory) < 2000 {
			// 	fmt.Println("prgrm.Memory", prgrm.Memory)
			// }

			return err
		}
		return err

	}
	return err

}

func (call *CXCall) ccall(prgrm *CXProgram, globalInputs *[]CXValue, globalOutputs *[]CXValue) error {
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
			// prgrm.CallStack[prgrm.CallCounter].ccall(prgrm)
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
			// wiping this declaration's memory (removing garbage)
			newCall := &prgrm.CallStack[prgrm.CallCounter]
			newFP := newCall.FramePointer
			size := GetSize(expr.Outputs[0])
			for c := 0; c < size; c++ {
				prgrm.Memory[newFP+expr.Outputs[0].Offset+c] = 0
			}
			call.Line++
		} else {
			switch expr.Operator.Version {
				case 1: // old version
					opcodeHandlers[expr.Operator.OpCode](expr, call.FramePointer)
					call.Line++
				case 2: // new version
					fp := call.FramePointer;

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


		            argIndex := 0;
					for inputIndex := 0; inputIndex < inputCount; inputIndex++ {
                        input := inputs[inputIndex]
                        offset := GetFinalOffset(fp, input)
                        value := &inputValues[inputIndex]
                        value.Arg = input
                        value.Used = -1
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
                        value.Used = -1
                        value.Offset = offset
                        value.Type = output.Type
                        value.FramePointer = fp
                        value.Expr = expr
                        argIndex++
					}

					opcodeHandlers_V2[expr.Operator.OpCode](inputValues, outputValues)

					for inputIndex := 0; inputIndex < inputCount; inputIndex++ { // TODO: remove in release builds
						if inputValues[inputIndex].Used != int8(inputs[inputIndex].Type) { // TODO: remove cast
							panic(fmt.Sprintf("Input value not used for opcode: '%s', param #%d. Expected type %d, '%s', used type %d, '%s'.",
							 	OpNames[expr.Operator.OpCode],
                                inputIndex + 1,
							 	inputs[inputIndex].Type, TypeNames[inputs[inputIndex].Type],
								inputValues[inputIndex].Used, TypeNames[int(inputValues[inputIndex].Used)]))
						}
					}

					for outputIndex := 0; outputIndex < outputCount; outputIndex++ { // TODO: remove in release builds
						if outputValues[outputIndex].Used != int8(outputs[outputIndex].Type) { // TODO: remove cast
							panic(fmt.Sprintf("Output value not used for opcode: '%s', param #%d. Expected type %d, '%s', used type %d '%s'.",
							 	OpNames[expr.Operator.OpCode],
                                outputIndex + 1,
							 	outputs[outputIndex].Type, TypeNames[outputs[outputIndex].Type],
								outputValues[outputIndex].Used, TypeNames[int(outputValues[outputIndex].Used)]))
						}
					}

        			call.Line++
        		default:
					/*
					   It was not a native, so we need to create another call
					   with the current expression's operator
					*/
					// we're going to use the next call in the callstack
					prgrm.CallCounter++
					if prgrm.CallCounter >= CALLSTACK_SIZE {
						panic(STACK_OVERFLOW_ERROR)
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
					if prgrm.StackPointer > STACK_SIZE {
						panic(STACK_OVERFLOW_ERROR)
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
						if inp.PassBy == PASSBY_REFERENCE {
							// If we're referencing an inner element, like an element of a slice (&slc[0])
							// or a field of a struct (&struct.fld) we no longer need to add
							// the OBJECT_HEADER_SIZE to the offset
							if inp.IsInnerReference {
								finalOffset -= OBJECT_HEADER_SIZE
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
							// newFP + newCall.Operator.Inputs[i].Offset,
							// GetFinalOffset(prgrm.Memory, newFP, newCall.Operator.Inputs[i], MEM_WRITE),
							byts)
					}
			}
		}
	}
	return nil
}

// Callback ...
func (prgrm *CXProgram) Callback(fn *CXFunction, inputs [][]byte) (outputs [][]byte) {
	line := prgrm.CallStack[prgrm.CallCounter].Line
	previousCall := prgrm.CallCounter
	prgrm.CallCounter++
	newCall := &prgrm.CallStack[prgrm.CallCounter]
	newCall.Operator = fn
	newCall.Line = 0
	newCall.FramePointer = prgrm.StackPointer
	prgrm.StackPointer += newCall.Operator.Size
	newFP := newCall.FramePointer

	// wiping next mem frame (removing garbage)
	for c := 0; c < fn.Size; c++ {
		prgrm.Memory[newFP+c] = 0
	}

	for i, inp := range inputs {
		WriteMemory(GetFinalOffset(newFP, newCall.Operator.Inputs[i]), inp)
	}

	var nCalls = 0
	if err := prgrm.Run(true, &nCalls, previousCall); err != nil {
		os.Exit(CX_INTERNAL_ERROR)
	}

	prgrm.CallCounter = previousCall
	prgrm.CallStack[prgrm.CallCounter].Line = line

	for _, out := range fn.Outputs {
		// Making a copy of the bytes, so if we modify the bytes being held by `outputs`
		// we don't modify the program memory.
		mem := ReadMemory(GetFinalOffset(newFP, out), out)
		cop := make([]byte, len(mem))
		copy(cop, mem)
		outputs = append(outputs, cop)
	}
	return outputs
}
