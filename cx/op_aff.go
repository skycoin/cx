package base

import (
	"fmt"
	"github.com/skycoin/skycoin/src/cipher/encoder"
)

func GetInferActions (inp *CXArgument, fp int) []string {
	inpOffset := GetFinalOffset(fp, inp)

	var off int32
	encoder.DeserializeAtomic(PROGRAM.Memory[inpOffset : inpOffset + TYPE_POINTER_SIZE], &off)

	var l int32
	_l := PROGRAM.Memory[off + OBJECT_HEADER_SIZE : off + OBJECT_HEADER_SIZE + SLICE_HEADER_SIZE]
	encoder.DeserializeAtomic(_l[:4], &l)

	result := make([]string, l)

	// for c := int(l); c > 0; c-- {
	for c := 0; c < int(l); c++ {
		var elOff int32
		// encoder.DeserializeAtomic(PROGRAM.Memory[int(off) + OBJECT_HEADER_SIZE + SLICE_HEADER_SIZE + (c - 1) * TYPE_POINTER_SIZE : int(off) + OBJECT_HEADER_SIZE + SLICE_HEADER_SIZE + c * STR_HEADER_SIZE], &elOff)
		encoder.DeserializeAtomic(PROGRAM.Memory[int(off) + OBJECT_HEADER_SIZE + SLICE_HEADER_SIZE + c * TYPE_POINTER_SIZE : int(off) + OBJECT_HEADER_SIZE + SLICE_HEADER_SIZE + (c + 1) * STR_HEADER_SIZE], &elOff)

		var size int32
		encoder.DeserializeAtomic(PROGRAM.Memory[elOff : elOff + STR_HEADER_SIZE], &size)

		var res string
		encoder.DeserializeRaw(PROGRAM.Memory[elOff : elOff + STR_HEADER_SIZE + size], &res)

		// result[int(l) - c] = res
		result[c] = res
	}

	return result
}

func op_aff_print (expr *CXExpression, fp int) {
	inp1 := expr.Inputs[0]
	fmt.Println(GetInferActions(inp1, fp))
}

func affExpr (expr *CXExpression) {
	
}

func op_aff_on (expr *CXExpression, fp int) {
	inp1, inp2 := expr.Inputs[0], expr.Inputs[1]

	_ = inp1

	// inp1Offset := GetFinalOffset(fp, inp1)
	// inp2Offset := GetFinalOffset(fp, inp2)

	prevPkg := PROGRAM.CurrentPackage
	prevFn := prevPkg.CurrentFunction
	prevExpr := prevFn.CurrentExpression

	var tgtPkg *CXPackage = prevPkg
	var tgtFn *CXFunction = prevFn
	var tgtExpr *CXExpression = prevExpr

	_ = tgtExpr

	// processing the target
	var tgtCmd string
	for _, aff := range GetInferActions(inp2, fp) {
		switch aff {
		case "pkg":
			tgtCmd = "pkg"
		case "fn":
			tgtCmd = "fn"
		case "expr":
			tgtCmd = "expr"
		default:
			switch tgtCmd {
			case "pkg":
				if pkg, err := PROGRAM.GetPackage(aff); err == nil {
					tgtPkg = pkg
				} else {
					panic(err)
				}
			case "fn":
				if fn, err := tgtPkg.GetFunction(aff); err == nil {
					tgtFn = fn
				} else {
					panic(err)
				}
			case "expr":
				if expr, err := tgtFn.GetExpressionByLabel(aff); err == nil {
					tgtExpr = expr
				} else {
					panic(err)
				}
			}
		}
	}

	// var affPkg *CXPackage = prevPkg
	// var affFn *CXFunction = prevFn
	// var affExpr *CXExpression = prevExpr
	
	// processing the affordances
	var affCmd string
	var results []string
	for _, aff := range GetInferActions(inp1, fp) {
		switch aff {
		case "expr":
			affCmd = "expr"
		case "program":
			affCmd = "program"
			// do it in here
		default:
			switch affCmd {
			case "pkg":
				if pkg, err := PROGRAM.GetPackage(aff); err == nil {
					// affPkg = pkg
					_ = pkg
				} else {
					panic(err)
				}
			case "fn":
				if fn, err := tgtPkg.GetFunction(aff); err == nil {
					// affFn = fn
					_ = fn
				} else {
					panic(err)
				}
			case "expr":
				if expr, err := tgtFn.GetExpressionByLabel(aff); err == nil {
					// affExpr = expr
					_ = expr
					switch tgtCmd {
					case "pkg":
					case "fn":
					case "expr":
						results = append(results, "expr-expr")
					}
				} else {
					panic(err)
				}
			}
		}
	}

	// returning to previous state
	PROGRAM.CurrentPackage = prevPkg
	PROGRAM.CurrentPackage.CurrentFunction = prevFn
	PROGRAM.CurrentPackage.CurrentFunction.CurrentExpression = prevExpr

	fmt.Println(results)
}

func op_aff_of (expr *CXExpression, fp int) {
	
}

func op_aff_inform (expr *CXExpression, fp int) {
	
}

func op_aff_request (expr *CXExpression, fp int) {
	
}

func CallAffPredicate (fn *CXFunction, predValue []byte) byte {
	prevCall := &PROGRAM.CallStack[PROGRAM.CallCounter]
	
	PROGRAM.CallCounter++
	newCall := &PROGRAM.CallStack[PROGRAM.CallCounter]
	newCall.Operator = fn
	newCall.Line = 0
	newCall.FramePointer = PROGRAM.StackPointer
	PROGRAM.StackPointer += newCall.Operator.Size

	newFP := newCall.FramePointer

	// wiping next mem frame (removing garbage)
	for c := 0; c < fn.Size; c++ {
		PROGRAM.Memory[newFP+c] = 0
	}

	// sending value to predicate function
	WriteMemory(
		GetFinalOffset(newFP, newCall.Operator.Inputs[0]),
		predValue)

	prevCC := PROGRAM.CallCounter
	for true {
		call := &PROGRAM.CallStack[PROGRAM.CallCounter]
		call.ccall(PROGRAM)
		if PROGRAM.CallCounter < prevCC {
			break
		}
	}

	prevCall.Line--

	return ReadMemory(GetFinalOffset(
		newCall.FramePointer,
		newCall.Operator.Outputs[0]),
		newCall.Operator.Outputs[0])[0]
}

// This might not make sense, as we can use normal programming to create conditions on values
// func QueryValue (fn *CXFunction, argOffsetB []byte, affOffset *int) {
// 	for c := 0; c <= PROGRAM.CallCounter; c++ {
// 		inFP := 0
// 		op := PROGRAM.CallStack[c].Operator
		
// 		for _, expr := range op.Expressions {
// 			if expr.Operator == nil {
// 				for _, out := range expr.Outputs {
// 					if fn.Inputs[0].Type == out.Type && out.Name != "" {
// 						res := CallAffPredicate(fn, PROGRAM.Memory[inFP + out.Offset : inFP + out.Offset + out.TotalSize])

// 						if res == 1 {
// 							*affOffset = WriteToSlice(*affOffset, argOffsetB)

// 							outNameB := encoder.Serialize(out.Name)
// 							outNameOffset := AllocateSeq(len(outNameB))
// 							WriteMemory(outNameOffset, outNameB)
// 							outNameOffsetB := encoder.SerializeAtomic(int32(outNameOffset))

// 							*affOffset = WriteToSlice(*affOffset, outNameOffsetB)
// 						}
// 					}
// 				}
// 			}
// 		}

// 		inFP += op.Size
// 	}
// }

func QueryArgument (fn *CXFunction, argOffsetB []byte, affOffset *int) {
	
}

func QueryExpressions (fn *CXFunction, expr *CXExpression, exprOffsetB []byte, affOffset *int) {
	for _, ex := range expr.Function.Expressions {
		if ex.Operator == nil || ex.Label == "" {
			// then it's a variable declaration
			// or it's a non-labelled expression
			continue
		}

		var opNameB []byte
		if ex.Operator.IsNative {
			opNameB = encoder.Serialize(OpNames[ex.Operator.OpCode])
		} else {
			opNameB = encoder.Serialize(ex.Operator.Name)
		}

		opNameOffset := AllocateSeq(len(opNameB))
		WriteMemory(opNameOffset, opNameB)
		opNameOffsetB := encoder.SerializeAtomic(int32(opNameOffset))

		res := CallAffPredicate(fn, opNameOffsetB)

		if res == 1 {
			*affOffset = WriteToSlice(*affOffset, exprOffsetB)

			lblNameB := encoder.Serialize(ex.Label)
			lblNameOffset := AllocateSeq(len(lblNameB))
			WriteMemory(lblNameOffset, lblNameB)
			lblNameOffsetB := encoder.SerializeAtomic(int32(lblNameOffset))

			*affOffset = WriteToSlice(*affOffset, lblNameOffsetB)
		}
	}
}

func getSignatureSlice (params []*CXArgument) int {
	var sliceOffset int
	for _, param := range params {
		
		var typOffset int
		if param.CustomType != nil {
			// then it's custom type
			typOffset = WriteObjectRetOff(encoder.Serialize(param.CustomType.Package.Name + "." + param.CustomType.Name))
		} else {
			// then it's native type
			typOffset = WriteObjectRetOff(encoder.Serialize(TypeNames[param.Type]))
		}
		
		sliceOffset = WriteToSlice(sliceOffset, encoder.SerializeAtomic(int32(typOffset)))
	}

	return sliceOffset
}

func QueryFunction (fn *CXFunction, expr *CXExpression, fnOffsetB []byte, affOffset *int) {
	for _, f := range expr.Package.Functions {
		if f.Name == SYS_INIT_FUNC {
			continue
		}
		
		var opNameB []byte
		if f.IsNative {
			opNameB = encoder.Serialize(OpNames[f.OpCode])
		} else {
			opNameB = encoder.Serialize(f.Name)
		}

		opNameOffsetB := encoder.SerializeAtomic(int32(WriteObjectRetOff(opNameB)))
		
		inpSigOffset := getSignatureSlice(f.Inputs)
		outSigOffset := getSignatureSlice(f.Outputs)

		fnOffset := AllocateSeq(OBJECT_HEADER_SIZE + STR_SIZE + TYPE_POINTER_SIZE + TYPE_POINTER_SIZE)
		// Name
		WriteMemory(fnOffset + OBJECT_HEADER_SIZE, opNameOffsetB)
		// InputSignature
		WriteMemory(fnOffset + OBJECT_HEADER_SIZE + TYPE_POINTER_SIZE, encoder.SerializeAtomic(int32(inpSigOffset)))
		// OutputSignature
		WriteMemory(fnOffset + OBJECT_HEADER_SIZE + TYPE_POINTER_SIZE + TYPE_POINTER_SIZE, encoder.SerializeAtomic(int32(outSigOffset)))

		val := PROGRAM.Memory[fnOffset + OBJECT_HEADER_SIZE : fnOffset + OBJECT_HEADER_SIZE + STR_SIZE + TYPE_POINTER_SIZE + TYPE_POINTER_SIZE]
		res := CallAffPredicate(fn, val)
		
		if res == 1 {
			*affOffset = WriteToSlice(*affOffset, fnOffsetB)
			*affOffset = WriteToSlice(*affOffset, opNameOffsetB)
		}
	}
}

func QueryCaller (fn *CXFunction, expr *CXExpression, callerOffsetB []byte, affOffset *int) {
	if PROGRAM.CallCounter == 0 {
		// then it's entry point
		return
	}

	call := PROGRAM.CallStack[PROGRAM.CallCounter - 1]

	var opNameB []byte
	if call.Operator.IsNative {
		opNameB = encoder.Serialize(OpNames[call.Operator.OpCode])
	} else {
		opNameB = encoder.Serialize(call.Operator.Package.Name + "." + call.Operator.Name)
	}


	opNameOffsetB := encoder.SerializeAtomic(int32(WriteObjectRetOff(opNameB)))

	callOffset := AllocateSeq(OBJECT_HEADER_SIZE + STR_SIZE + I32_SIZE)
	// FnName
	WriteMemory(callOffset + OBJECT_HEADER_SIZE, opNameOffsetB)
	// FnSize
	WriteMemory(callOffset + OBJECT_HEADER_SIZE + STR_SIZE, encoder.SerializeAtomic(int32(call.Operator.Size)))

	res := CallAffPredicate(fn, PROGRAM.Memory[callOffset + OBJECT_HEADER_SIZE : callOffset + OBJECT_HEADER_SIZE + STR_SIZE + I32_SIZE])

	if res == 1 {
		*affOffset = WriteToSlice(*affOffset, callerOffsetB)
	}
}

func QueryProgram (fn *CXFunction, expr *CXExpression, prgrmOffsetB []byte, affOffset *int) {
	prgrmOffset := AllocateSeq(OBJECT_HEADER_SIZE + I32_SIZE + I64_SIZE)
	// Callcounter
	WriteMemory(prgrmOffset + OBJECT_HEADER_SIZE, encoder.SerializeAtomic(int32(PROGRAM.CallCounter)))
	// HeapUsed
	WriteMemory(prgrmOffset + OBJECT_HEADER_SIZE + I32_SIZE, encoder.Serialize(int64(PROGRAM.HeapPointer)))
	res := CallAffPredicate(fn, PROGRAM.Memory[prgrmOffset + OBJECT_HEADER_SIZE : prgrmOffset + OBJECT_HEADER_SIZE + I32_SIZE + I64_SIZE])

	if res == 1 {
		*affOffset = WriteToSlice(*affOffset, prgrmOffsetB)
	}
}

func op_aff_query (expr *CXExpression, fp int) {
	inp1, out1 := expr.Inputs[0], expr.Outputs[0]

	out1Offset := GetFinalOffset(fp, out1)

	var affOffset int
	
	var cmd string
	for _, rule := range GetInferActions(inp1, fp) {
		switch rule {
		case "filter":
			cmd = "filter"
		case "sort":
			cmd = "sort"
		default:
			switch cmd {
			case "filter":
				if fn, err := inp1.Package.GetFunction(rule); err == nil {

					// arg keyword
					argB := encoder.Serialize("arg")
					argOffset := AllocateSeq(len(argB))
					WriteMemory(argOffset, argB)
					argOffsetB := encoder.SerializeAtomic(int32(argOffset))

					// expr keyword
					exprB := encoder.Serialize("expr")
					exprOffset := AllocateSeq(len(exprB))
					WriteMemory(exprOffset, exprB)
					exprOffsetB := encoder.SerializeAtomic(int32(exprOffset))
					
					// fn keyword
					fnB := encoder.Serialize("fn")
					fnOffset := AllocateSeq(len(fnB))
					WriteMemory(fnOffset, fnB)
					fnOffsetB := encoder.SerializeAtomic(int32(fnOffset))

					// caller keyword
					callerB := encoder.Serialize("caller")
					callerOffset := AllocateSeq(len(callerB))
					WriteMemory(callerOffset, callerB)
					callerOffsetB := encoder.SerializeAtomic(int32(callerOffset))

					// program keyword
					prgrmB := encoder.Serialize("program")
					prgrmOffset := AllocateSeq(len(prgrmB))
					WriteMemory(prgrmOffset, prgrmB)
					prgrmOffsetB := encoder.SerializeAtomic(int32(prgrmOffset))

					predInp := fn.Inputs[0]
					
					if predInp.Type == TYPE_CUSTOM {
						if predInp.CustomType != nil {
							switch predInp.CustomType.Name {
							case "Argument":
								QueryArgument(fn, argOffsetB, &affOffset)
							case "Expression":
								QueryExpressions(fn, expr, exprOffsetB, &affOffset)
							case "Function":
								QueryFunction(fn, expr, fnOffsetB, &affOffset)
							case "Caller":
								QueryCaller(fn, expr, callerOffsetB, &affOffset)
							case "Program":
								QueryProgram(fn, expr, prgrmOffsetB, &affOffset)
							}
						}
					}

				} else {
					panic(err)
				}
			case "sort":
				
			}
		}
	}

	WriteMemory(out1Offset, FromI32(int32(affOffset)))
}
