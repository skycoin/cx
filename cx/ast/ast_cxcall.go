package ast

import (
	"github.com/skycoin/cx/cx/constants"
	"github.com/skycoin/cx/cx/types"
	//"fmt"
)

// CXCall ...
type CXCall struct {
	Operator     *CXFunction   // What CX function will be called when running this CXCall in the runtime
	Line         int           // What line in the CX function is currently being executed
	FramePointer types.Pointer // Where in the stack is this function call's local variables stored
}

func popStack(prgrm *CXProgram, call *CXCall) error {
	// going back to the previous call
	prgrm.CallCounter--
	if !prgrm.CallCounter.IsValid() {
		// then the program finished
		prgrm.Terminated = true
		return nil
	}

	// copying the outputs to the previous stack frame
	returnAddr := &prgrm.CallStack[prgrm.CallCounter]
	returnOp := returnAddr.Operator
	returnLine := returnAddr.Line
	returnFP := returnAddr.FramePointer
	fp := call.FramePointer

	expr := returnOp.Expressions[returnLine]

	cxAtomicOp, err := prgrm.GetCXAtomicOp(expr.Index)
	if err != nil {
		return err
	}

	cxAtomicOpOutputs := cxAtomicOp.GetOutputs(prgrm)
	lenOuts := len(cxAtomicOpOutputs)
	callOperatorOutputs := call.Operator.GetOutputs(prgrm)
	for i, outputIdx := range callOperatorOutputs {
		output := prgrm.GetCXTypeSignatureFromArray(outputIdx)
		// Continuing if there is no receiving variable available.
		if i >= lenOuts {
			continue
		}

		cxAtomicOpOutput := prgrm.GetCXTypeSignatureFromArray(cxAtomicOpOutputs[i])
		types.WriteSlice_byte(prgrm.Memory, GetFinalOffset(prgrm, returnFP, nil, cxAtomicOpOutput),
			types.GetSlice_byte(prgrm.Memory, GetFinalOffset(prgrm, fp, nil, output), output.GetSize(prgrm, false)))
	}

	// return the stack pointer to its previous state
	prgrm.Stack.Pointer = fp
	// we'll now execute the next command
	prgrm.CallStack[prgrm.CallCounter].Line++
	// calling the actual command
	// prgrm.CallStack[prgrm.CallCounter].Call(prgrm)
	return nil
}

func wipeDeclarationMemory(prgrm *CXProgram, expr *CXExpression) error {
	cxAtomicOp, err := prgrm.GetCXAtomicOp(expr.Index)
	if err != nil {
		return err
	}

	newCall := &prgrm.CallStack[prgrm.CallCounter]
	newFP := newCall.FramePointer
	cxAtomicOpOutputs := cxAtomicOp.GetOutputs(prgrm)
	cxAtomicOutputTypeSig := prgrm.GetCXTypeSignatureFromArray(cxAtomicOpOutputs[0])
	size := cxAtomicOutputTypeSig.GetSize(prgrm, false)
	var offset types.Pointer
	if cxAtomicOutputTypeSig.Type == TYPE_CXARGUMENT_DEPRECATE {
		offset = prgrm.CXArgs[cxAtomicOutputTypeSig.Meta].Offset
	} else if cxAtomicOutputTypeSig.Type == TYPE_ATOMIC || cxAtomicOutputTypeSig.Type == TYPE_POINTER_ATOMIC {
		offset = cxAtomicOutputTypeSig.Offset
	} else if cxAtomicOutputTypeSig.Type == TYPE_ARRAY_ATOMIC {
		offset = cxAtomicOutputTypeSig.Offset
	} else if cxAtomicOutputTypeSig.Type == TYPE_POINTER_ARRAY_ATOMIC {
		offset = cxAtomicOutputTypeSig.Offset
	}

	for c := types.Pointer(0); c < size; c++ {
		prgrm.Memory[newFP+offset+c] = 0
	}

	return nil
}

func processBuiltInOperators(prgrm *CXProgram, expr *CXExpression, globalInputs *[]CXValue, globalOutputs *[]CXValue, fp types.Pointer) error {
	cxAtomicOp, err := prgrm.GetCXAtomicOp(expr.Index)
	if err != nil {
		return err
	}

	cxAtomicOpOperator := prgrm.GetFunctionFromArray(cxAtomicOp.Operator)

	inputs := cxAtomicOp.GetInputs(prgrm)
	inputCount := len(inputs)
	if inputCount > len(*globalInputs) {
		*globalInputs = make([]CXValue, inputCount)
	}
	inputValues := (*globalInputs)[:inputCount]

	outputs := cxAtomicOp.GetOutputs(prgrm)
	outputCount := len(outputs)
	if outputCount > len(*globalOutputs) {
		*globalOutputs = make([]CXValue, outputCount)
	}
	outputValues := (*globalOutputs)[:outputCount]

	argIndex := 0
	for inputIndex := 0; inputIndex < inputCount; inputIndex++ {
		inputTypeSignature := prgrm.GetCXTypeSignatureFromArray(inputs[inputIndex])
		value := &inputValues[inputIndex]
		value.TypeSignature = inputTypeSignature
		value.Size = inputTypeSignature.GetSize(prgrm, false)
		offset := GetFinalOffset(prgrm, fp, nil, inputTypeSignature)
		value.Offset = offset

		var input *CXArgument = &CXArgument{}
		if inputTypeSignature.Type == TYPE_CXARGUMENT_DEPRECATE {
			input = prgrm.GetCXArgFromArray(CXArgumentIndex(inputTypeSignature.Meta))

			value.Type = input.Type
			if input.Type == types.POINTER {
				value.Type = input.PointerTargetType
			}
		} else if inputTypeSignature.Type == TYPE_ATOMIC || inputTypeSignature.Type == TYPE_POINTER_ATOMIC {
			value.Type = types.Code(inputTypeSignature.Meta)
		} else if inputTypeSignature.Type == TYPE_ARRAY_ATOMIC {
			arrDetails := prgrm.GetCXTypeSignatureArrayFromArray(inputTypeSignature.Meta)
			value.Type = types.Code(arrDetails.Type)
		} else if inputTypeSignature.Type == TYPE_POINTER_ARRAY_ATOMIC {
			arrDetails := prgrm.GetCXTypeSignatureArrayFromArray(inputTypeSignature.Meta)
			value.Type = types.Code(arrDetails.Type)
		}

		value.FramePointer = fp
		value.Expr = expr
		argIndex++
	}

	for outputIndex := 0; outputIndex < outputCount; outputIndex++ {
		outputTypeSignature := prgrm.GetCXTypeSignatureFromArray(outputs[outputIndex])

		value := &outputValues[outputIndex]
		value.TypeSignature = outputTypeSignature
		value.Size = outputTypeSignature.GetSize(prgrm, false)
		offset := GetFinalOffset(prgrm, fp, nil, outputTypeSignature)
		value.Offset = offset

		var output *CXArgument = &CXArgument{}
		if outputTypeSignature.Type == TYPE_CXARGUMENT_DEPRECATE {
			output = prgrm.GetCXArgFromArray(CXArgumentIndex(outputTypeSignature.Meta))

			value.Type = output.Type
			if output.Type == types.POINTER {
				value.Type = output.PointerTargetType
			}
		} else if outputTypeSignature.Type == TYPE_ATOMIC || outputTypeSignature.Type == TYPE_POINTER_ATOMIC {
			value.Type = types.Code(outputTypeSignature.Meta)
		} else if outputTypeSignature.Type == TYPE_ARRAY_ATOMIC {
			arrDetails := prgrm.GetCXTypeSignatureArrayFromArray(outputTypeSignature.Meta)
			value.Type = types.Code(arrDetails.Type)
		} else if outputTypeSignature.Type == TYPE_POINTER_ARRAY_ATOMIC {
			arrDetails := prgrm.GetCXTypeSignatureArrayFromArray(outputTypeSignature.Meta)
			value.Type = types.Code(arrDetails.Type)
		}

		value.FramePointer = fp
		value.Expr = expr
		argIndex++
	}

	OpcodeHandlers[cxAtomicOpOperator.AtomicOPCode](prgrm, inputValues, outputValues)

	return nil
}

func processNonAtomicOperators(prgrm *CXProgram, expr *CXExpression, fp types.Pointer) error {
	cxAtomicOp, err := prgrm.GetCXAtomicOp(expr.Index)
	if err != nil {
		return err
	}

	cxAtomicOpOperator := prgrm.GetFunctionFromArray(cxAtomicOp.Operator)
	//TODO: Is this only called for user defined functions?
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
	newCall.Operator = cxAtomicOpOperator
	newCall.Line = 0
	newCall.FramePointer = prgrm.Stack.Pointer
	// the stack pointer is moved to create room for the next call
	// prgrm.MemoryPointer += fn.Size
	prgrm.Stack.Pointer += newCall.Operator.Size

	// checking if enough memory in stack
	if prgrm.Stack.Pointer > constants.STACK_SIZE {
		panic(constants.STACK_OVERFLOW_ERROR)
	}

	newFP := newCall.FramePointer
	// wiping next stack frame (removing garbage)
	for c := types.Pointer(0); c < cxAtomicOpOperator.Size; c++ {
		prgrm.Memory[newFP+c] = 0
	}

	newCallOperatorInputs := newCall.Operator.GetInputs(prgrm)
	for i, inputIdx := range cxAtomicOp.GetInputs(prgrm) {
		input := prgrm.GetCXTypeSignatureFromArray(inputIdx)
		var inp *CXArgument = &CXArgument{}
		if input.Type == TYPE_CXARGUMENT_DEPRECATE {
			inp = prgrm.GetCXArgFromArray(CXArgumentIndex(input.Meta))
		}

		var byts []byte

		finalOffset := GetFinalOffset(prgrm, fp, nil, input)
		if input.PassBy == constants.PASSBY_REFERENCE || inp.PassBy == constants.PASSBY_REFERENCE {

			// If we're referencing an inner element, like an element of a slice (&slc[0])
			// or a field of a struct (&struct.fld) we no longer need to add
			// the OBJECT_HEADER_SIZE to the offset
			if inp.IsInnerReference {
				finalOffset -= types.OBJECT_HEADER_SIZE
			}
			var finalOffsetB [types.POINTER_SIZE]byte
			types.Write_ptr(finalOffsetB[:], 0, finalOffset)
			byts = finalOffsetB[:]

		} else {
			size := input.GetSize(prgrm, false)
			byts = prgrm.Memory[finalOffset : finalOffset+size]
		}

		newCallOperatorInputTypeSignature := prgrm.GetCXTypeSignatureFromArray(newCallOperatorInputs[i])
		// writing inputs to new stack frame
		types.WriteSlice_byte(
			prgrm.Memory,
			GetFinalOffset(prgrm, newFP, nil, newCallOperatorInputTypeSignature),
			// newFP + newCall.Operator.ProgramInput[i].Offset,
			// GetFinalOffset(prgrm.Memory, newFP, newCall.Operator.ProgramInput[i], MEM_WRITE),
			byts)
	}

	return nil
}

func (call *CXCall) Call(prgrm *CXProgram, globalInputs *[]CXValue, globalOutputs *[]CXValue) error {
	// CX is still single-threaded, so only one stack
	if call.Line >= call.Operator.LineCount {
		/*
		   popping the stack
		*/
		err := popStack(prgrm, call)
		if err != nil {
			return err
		}

		return nil
	}

	/*
	   continue with call operator's execution
	*/

	fn := call.Operator
	expr := fn.Expressions[call.Line]
	fp := call.FramePointer

	if expr.Type == CX_LINE {
		call.Line++
		return nil
	}

	cxAtomicOp, err := prgrm.GetCXAtomicOp(expr.Index)
	if err != nil {
		return err
	}
	cxAtomicOpOperator := prgrm.GetFunctionFromArray(cxAtomicOp.Operator)
	// if it's a native, then we just process the arguments with execNative
	//TODO: WHEN WOULD OPERATOR EVER BE NIL?
	if cxAtomicOpOperator == nil {
		// then it's a declaration
		// wiping this declaration's memory (removing garbage)
		err := wipeDeclarationMemory(prgrm, &expr)
		if err != nil {
			return err
		}

		call.Line++
	} else if cxAtomicOpOperator.IsBuiltIn() {
		//TODO: SLICES ARE NON ATOMIC
		err := processBuiltInOperators(prgrm, &expr, globalInputs, globalOutputs, fp)
		if err != nil {
			return err
		}

		call.Line++
	} else {
		//NON-ATOMIC OPERATOR
		err := processNonAtomicOperators(prgrm, &expr, fp)
		if err != nil {
			return err
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
