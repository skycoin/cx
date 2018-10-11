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

func op_aff_execute (expr *CXExpression, fp int) {
	
}

func op_aff_query (expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]

	inp1Offset := GetFinalOffset(fp, inp1)
	inp2Offset := GetFinalOffset(fp, inp2)
	out1Offset := GetFinalOffset(fp, out1)

	_ = inp2Offset

	_ = out1

	var sliceHeader []byte
	
	var len1 int32
	sliceHeader = PROGRAM.Memory[inp1Offset-SLICE_HEADER_SIZE : inp1Offset]
	encoder.DeserializeAtomic(sliceHeader[:4], &len1)

	var cmd string
	for _, rule := range GetInferActions(inp2, fp) {
		switch rule {
		case "filter":
			cmd = "filter"
		case "sort":
			cmd = "sort"
		default:
			switch cmd {
			case "filter":
				if fn, err := inp2.Package.GetFunction(rule); err == nil {
					inFP := 0

					var affOffset int

					argB := encoder.Serialize("arg")
					argOffset := AllocateSeq(len(argB))
					WriteMemory(argOffset, argB)
					argOffsetB := encoder.SerializeAtomic(int32(argOffset))
					
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

					WriteMemory(out1Offset, FromI32(int32(affOffset)))
				} else {
					panic(err)
				}
			case "sort":
				
			}
		}
	}
}
