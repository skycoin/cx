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

	for c := int(l); c > 0; c-- {
		var elOff int32
		encoder.DeserializeAtomic(PROGRAM.Memory[int(off) + OBJECT_HEADER_SIZE + SLICE_HEADER_SIZE + (c - 1) * TYPE_POINTER_SIZE : int(off) + OBJECT_HEADER_SIZE + SLICE_HEADER_SIZE + c * STR_HEADER_SIZE], &elOff)

		var size int32
		encoder.DeserializeAtomic(PROGRAM.Memory[elOff : elOff + STR_HEADER_SIZE], &size)

		var res string
		encoder.DeserializeRaw(PROGRAM.Memory[elOff : elOff + STR_HEADER_SIZE + size], &res)

		result[int(l) - c] = res
	}

	return result
}

func op_aff_print (expr *CXExpression, fp int) {
	inp1 := expr.Inputs[0]
	fmt.Println(GetInferActions(inp1, fp))
}

func op_aff_on (expr *CXExpression, fp int) {
	
}

func op_aff_of (expr *CXExpression, fp int) {
	
}

func op_aff_inform (expr *CXExpression, fp int) {
	
}

func op_aff_request (expr *CXExpression, fp int) {
	
}

func op_aff_query (expr *CXExpression, fp int) {
	inp1, out1 := expr.Inputs[0], expr.Outputs[0]

	inp1Offset := GetFinalOffset(fp, inp1)
	out1Offset := GetFinalOffset(fp, out1)

	var sliceHeader []byte
	
	var len1 int32
	sliceHeader = PROGRAM.Memory[inp1Offset-SLICE_HEADER_SIZE : inp1Offset]
	encoder.DeserializeAtomic(sliceHeader[:4], &len1)

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
					inFP := 0

					var affOffset int

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

					predInp := fn.Inputs[0]
					
					if predInp.Type == TYPE_CUSTOM {
						if predInp.CustomType != nil {
							switch predInp.CustomType.Name {
							case "Argument":
								// get all possible values
								for c := 0; c <= PROGRAM.CallCounter; c++ {
									op := PROGRAM.CallStack[c].Operator
									
									for _, expr := range op.Expressions {
										if expr.Operator == nil {
											for _, out := range expr.Outputs {
												if fn.Inputs[0].Type == out.Type && out.Name != "" {
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
														PROGRAM.Memory[inFP + out.Offset : inFP + out.Offset + out.TotalSize])

													prevCC := PROGRAM.CallCounter
													for true {
														call := &PROGRAM.CallStack[PROGRAM.CallCounter]
														call.ccall(PROGRAM)
														if PROGRAM.CallCounter < prevCC {
															break
														}
													}

													prevCall.Line--

													if ReadMemory(GetFinalOffset(
														newCall.FramePointer,
														newCall.Operator.Outputs[0]),
														newCall.Operator.Outputs[0])[0] == 1 {
															affOffset = WriteToSlice(affOffset, argOffsetB)

															outNameB := encoder.Serialize(out.Name)
															outNameOffset := AllocateSeq(len(outNameB))
															WriteMemory(outNameOffset, outNameB)
															outNameOffsetB := encoder.SerializeAtomic(int32(outNameOffset))

															affOffset = WriteToSlice(affOffset, outNameOffsetB)
														}
												}
											}
										}
									}

									inFP += op.Size
								}
							case "Expression":
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
										opNameOffsetB)

									prevCC := PROGRAM.CallCounter
									for true {
										call := &PROGRAM.CallStack[PROGRAM.CallCounter]
										call.ccall(PROGRAM)
										if PROGRAM.CallCounter < prevCC {
											break
										}
									}

									prevCall.Line--

									if ReadMemory(GetFinalOffset(
										newCall.FramePointer,
										newCall.Operator.Outputs[0]),
										newCall.Operator.Outputs[0])[0] == 1 {
											affOffset = WriteToSlice(affOffset, exprOffsetB)

											lblNameB := encoder.Serialize(ex.Label)
											lblNameOffset := AllocateSeq(len(lblNameB))
											WriteMemory(lblNameOffset, lblNameB)
											lblNameOffsetB := encoder.SerializeAtomic(int32(lblNameOffset))

											affOffset = WriteToSlice(affOffset, lblNameOffsetB)
										}
								}
							}
						}
					}

					WriteMemory(out1Offset, FromI32(int32(affOffset)))
				} else {
					panic(err)
				}
			case "sort":
				
			}
		}
	}
}
