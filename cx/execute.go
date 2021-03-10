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

		err = call.ccall(prgrm)
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

var globalInputs []CXValue
var globalOutputs []CXValue
var globalArgs []*CXArgument
var globalOffsets []uint64

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
					if inputCount > len(globalInputs) {
						globalInputs = make([]CXValue, inputCount)
					}
					inputValues := globalInputs[:inputCount]


			        outputs := expr.Outputs
			        outputCount := len(outputs)
					if outputCount > len(globalOutputs) {
						globalOutputs = make([]CXValue, outputCount)
					}
					outputValues := globalOutputs[:outputCount]

					argCount := inputCount + outputCount
					if argCount > len(globalArgs) {
						globalArgs = make([]*CXArgument, argCount)
						globalOffsets = make([]uint64, argCount)
					}

				 	argIndex := 0;					
					for inputIndex := 0; inputIndex < inputCount; inputIndex++ {
						globalArgs[argIndex] = inputs[inputIndex]
						inputValues[inputIndex].Used = -1
						argIndex++
					}

					for outputIndex := 0; outputIndex < outputCount; outputIndex++ {
						globalArgs[argIndex] = outputs[outputIndex]
						outputValues[outputIndex].Used = -1
						argIndex++
					}

					for argIndex := 0; argIndex < argCount; argIndex++ {
						var arg = globalArgs[argIndex]
						finalOffset := arg.Offset
						if finalOffset < PROGRAM.StackSize {
							// Then it's in the stack, not in data or heap and we need to consider the frame pointer.
							finalOffset += fp
						}

						CalculateDereferences(arg, &finalOffset, fp) // TODO: hidden cost
						for _, fld := range arg.Fields {
							// elt = fld
							finalOffset += fld.Offset
							CalculateDereferences(fld, &finalOffset, fp) // TODO: hidden cost
						}

						globalOffsets[argIndex] = uint64(finalOffset)
					}

					for inputIndex := 0; inputIndex < inputCount; inputIndex++ {
						offset := int(globalOffsets[inputIndex]) // TODO: remove cast
						input := inputs[inputIndex]
						inputValues[inputIndex].Type = int8(input.Type) // TODO: remove cast
						size := GetSize(input) // TODO: hidden cost
						memory := PROGRAM.Memory[offset : offset+size]

						switch input.Type {
							case TYPE_I8:
								inputValues[inputIndex].Value_i8 = Deserialize_i8(memory) // TODO: Read_i8
							case TYPE_I16:
								inputValues[inputIndex].Value_i16 = Deserialize_i16(memory) // TODO: Read_i16
							case TYPE_I32:
								inputValues[inputIndex].Value_i32 = Deserialize_i32(memory) // TODO: Read_i32
							case TYPE_I64:
								inputValues[inputIndex].Value_i64 = Deserialize_i64(memory) // TODO: Read_i64
							case TYPE_UI8:
								inputValues[inputIndex].Value_ui8 = Deserialize_ui8(memory) // TODO: Read_ui8
							case TYPE_UI16:
								inputValues[inputIndex].Value_ui16 = Deserialize_ui16(memory) // TODO: Read_ui16
							case TYPE_UI32:
								inputValues[inputIndex].Value_ui32 = Deserialize_ui32(memory) // TODO: Read_ui32
							case TYPE_UI64:
								inputValues[inputIndex].Value_ui64 = Deserialize_ui64(memory) // TODO: Read_ui64
							case TYPE_F32:
								inputValues[inputIndex].Value_f32 = Deserialize_f32(memory) // TODO: Read_f32
							case TYPE_F64:
								inputValues[inputIndex].Value_f64 = Deserialize_f64(memory) // TODO: Read_f64
							case TYPE_BOOL:
								inputValues[inputIndex].Value_bool = Deserialize_bool(memory) // TODO: Read_bool
							case TYPE_STR:
								inputValues[inputIndex].Value_str = ReadStrFromOffset(offset, input) // TODO: Read_str
							default:
								panic(fmt.Sprintf("Unhandled type : %d", input.Type))
						}							
					}

					opcodeHandlers_V2[expr.Operator.OpCode](inputValues, outputValues)

					for inputIndex := 0; inputIndex < inputCount; inputIndex++ { // TODO: remove in release builds
						if inputValues[inputIndex].Used != int8(inputs[inputIndex].Type) { // TODO: remove cast
							panic("Input value not used in opcode")
						}
					}

					for outputIndex := 0; outputIndex < outputCount; outputIndex++ { // TODO: remove in release builds
						if outputValues[outputIndex].Used != int8(outputs[outputIndex].Type) { // TODO: remove cast
							panic("Output value not used in opcode")
						}
					}

			        for outputIndex := 0; outputIndex < outputCount; outputIndex++ {
			        	offset := int(globalOffsets[inputCount + outputIndex]) // TODO: remove cast
			        	output := outputs[outputIndex]
			        	switch output.Type {
							case TYPE_I8:
								WriteI8(offset, outputValues[outputIndex].Value_i8)
							case TYPE_I16:
								WriteI16(offset, outputValues[outputIndex].Value_i16)
							case TYPE_I32:
								WriteI32(offset, outputValues[outputIndex].Value_i32)
							case TYPE_I64:
								WriteI64(offset, outputValues[outputIndex].Value_i64)
							case TYPE_UI8:
								WriteUI8(offset, outputValues[outputIndex].Value_ui8)
							case TYPE_UI16:
								WriteUI16(offset, outputValues[outputIndex].Value_ui16)
							case TYPE_UI32:
								WriteUI32(offset, outputValues[outputIndex].Value_ui32)
							case TYPE_UI64:
								WriteUI64(offset, outputValues[outputIndex].Value_ui64)
							case TYPE_F32:
								WriteF32(offset, outputValues[outputIndex].Value_f32)
							case TYPE_F64:
								WriteF64(offset, outputValues[outputIndex].Value_f64)
							case TYPE_BOOL:
								WriteBool(offset, outputValues[outputIndex].Value_bool)
							case TYPE_STR:
								WriteObject(offset, encoder.Serialize(outputValues[outputIndex].Value_str))
							default:
								panic(fmt.Sprintf("Unhandled type : %d", output.Type))
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
