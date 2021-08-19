package opcodes

import (
	"fmt"
	"strconv"

	"github.com/skycoin/cx/cx/ast"
	"github.com/skycoin/cx/cx/constants"
	"github.com/skycoin/cx/cx/types"
	// "github.com/skycoin/skycoin/src/cipher/encoder"
)

var onMessages = map[string]string{
	"arg-arg-input":    "Replace %s.Input.%d with %s",
	"arg-arg-output":   "Replace %s.Output.%d with %s",
	"arg-expr":         "Add ",
	"prgrm-arg-input":  "Print %s.Input.%d's value",
	"prgrm-arg-output": "Print %s.Output.%d's value",
}
var ofMessages = map[string]string{
	"arg-arg-input":    "Replace %[3]s with %[1]s.Input.%[2]d",
	"arg-arg-output":   "Replace %[3]s with %[1]s.Output.%[2]d",
	"strct-arg-input":  "Add %[1]s.Input.%[2]d as a new field of %s",
	"strct-arg-output": "Add %[1]s.Output.%[2]d as a new field of %s",
	"prgrm-arg-input":  "Print %[1]s.Input.%[2]d's value",
	"prgrm-arg-output": "Print %[1]s.Output.%[2]d's value",
}

// GetInferActions ...
func GetInferActions(inp *ast.CXArgument, fp types.Pointer) []string {
	inpOffset := ast.GetFinalOffset(fp, inp)

	off := types.Read_ptr(ast.PROGRAM.Memory, inpOffset)

	l := types.Read_ptr(ast.GetSliceHeader(ast.GetSliceOffset(fp, inp)), types.POINTER_SIZE)

	result := make([]string, l)

	for c := types.Cast_int_to_ptr(0); c < l; c++ {
		elOff := types.Read_ptr(ast.PROGRAM.Memory, off+types.OBJECT_HEADER_SIZE+constants.SLICE_HEADER_SIZE+c*types.POINTER_SIZE)
		result[c] = types.Read_str_data(ast.PROGRAM.Memory, elOff)
	}

	return result
}

func opAffPrint(inputs []ast.CXValue, outputs []ast.CXValue) {
	inp1 := inputs[0]
	fmt.Println(GetInferActions(inp1.Arg, inp1.FramePointer))
	// for _, aff := range GetInferActions(inp1, fp) {
	// 	fmt.Println(aff)
	// }
}

// CallAffPredicate ...
func CallAffPredicate(fn *ast.CXFunction, predValue []byte) byte {
	prevCall := &ast.PROGRAM.CallStack[ast.PROGRAM.CallCounter]

	ast.PROGRAM.CallCounter++
	newCall := &ast.PROGRAM.CallStack[ast.PROGRAM.CallCounter]
	newCall.Operator = fn
	newCall.Line = 0
	newCall.FramePointer = ast.PROGRAM.StackPointer
	ast.PROGRAM.StackPointer += newCall.Operator.Size

	newFP := newCall.FramePointer

	// wiping next mem frame (removing garbage)
	for c := types.Pointer(0); c < fn.Size; c++ {
		ast.PROGRAM.Memory[newFP+c] = 0
	}

	// sending value to predicate function
	types.WriteSlice_byte(
		ast.PROGRAM.Memory,
		ast.GetFinalOffset(newFP, newCall.Operator.Inputs[0]),
		predValue)

	var inputs []ast.CXValue
	var outputs []ast.CXValue
	prevCC := ast.PROGRAM.CallCounter
	for {
		call := &ast.PROGRAM.CallStack[ast.PROGRAM.CallCounter]
		err := call.Ccall(ast.PROGRAM, &inputs, &outputs)
		if err != nil {
			panic(err)
		}
		if ast.PROGRAM.CallCounter < prevCC {
			break
		}
	}

	prevCall.Line--

	return types.GetSlice_byte(ast.PROGRAM.Memory, ast.GetFinalOffset(
		newCall.FramePointer,
		newCall.Operator.Outputs[0]),
		ast.GetSize(newCall.Operator.Outputs[0]))[0]
}

// This might not make sense, as we can use normal programming to create conditions on values
// func QueryValue (fn *CXFunction, argOffsetB []byte, affOffset *int) {
// 	for c := 0; c <= PROGRAM.CallCounter; c++ {
// 		inFP := 0
// 		op := PROGRAM.CallStack[c].Operator

// 		for _, expr := range op.Expressions {
// 			if expr.Operator == nil {
// 				for _, out := range expr.ProgramOutput {
// 					if fn.ProgramInput[0].Type == out.Type && out.Name != "" {
// 						res := CallAffPredicate(fn, PROGRAM.Memory[inFP + out.Offset : inFP + out.Offset + out.TotalSize])

// 						if res == 1 {
// 							*affOffset = WriteToSlice(*affOffset, argOffsetB)

// 							outNameB := encoder.Serialize(out.Name)
// 							outNameOffset := AllocateSeq(len(outNameB))
// 							WriteMemory(outNameOffset, outNameB)

// 							var outNameOffsetB [4]byte
//							WriteMemI32(outNameOffsetB[:], 0, int32(outNameOffset))
// 							*affOffset = WriteToSlice(*affOffset, outNameOffsetB[:])
// 						}
// 					}
// 				}
// 			}
// 		}

// 		inFP += op.Size
// 	}
// }

// Used by QueryArgument to query inputs and then outputs from expressions.
func queryParam(fn *ast.CXFunction, args []*ast.CXArgument, exprLbl string, argOffsetB []byte, affOffset *types.Pointer) {
	for i, arg := range args {

		var typOffset types.Pointer
		elt := ast.GetAssignmentElement(arg)
		if elt.StructType != nil {
			// typOffset = WriteObjectRetOff(encoder.Serialize(elt.StructType.Package.Name + "." + elt.StructType.Name))
			typOffset = types.AllocWrite_str_data(ast.PROGRAM.Memory, elt.StructType.Package.Name+"."+elt.StructType.Name)
		} else {
			// then it's native type
			// typOffset = WriteObjectRetOff(encoder.Serialize(TypeNames[elt.Type]))
			typOffset = types.AllocWrite_str_data(ast.PROGRAM.Memory, elt.Type.Name())
		}

		// Name
		// argNameB := encoder.Serialize(arg.Name)
		// argNameOffset := int32(WriteObjectRetOff(argNameB))
		argNameOffset := types.AllocWrite_str_data(ast.PROGRAM.Memory, arg.ArgDetails.Name)

		argOffset := ast.AllocateSeq(types.OBJECT_HEADER_SIZE + types.STR_SIZE + types.I32_SIZE + types.STR_SIZE)
		types.Write_ptr(ast.PROGRAM.Memory, argOffset+types.OBJECT_HEADER_SIZE, argNameOffset)

		// Index
		types.Write_ptr(ast.PROGRAM.Memory, argOffset+types.OBJECT_HEADER_SIZE+types.STR_SIZE, types.Cast_int_to_ptr(i))

		// Type
		types.Write_ptr(ast.PROGRAM.Memory, argOffset+types.OBJECT_HEADER_SIZE+types.STR_SIZE+types.I32_SIZE, typOffset)

		res := CallAffPredicate(fn, ast.PROGRAM.Memory[argOffset+types.OBJECT_HEADER_SIZE:argOffset+types.OBJECT_HEADER_SIZE+types.STR_SIZE+types.I32_SIZE+types.STR_SIZE])

		if res == 1 {
			*affOffset = ast.WriteToSlice(*affOffset, argOffsetB)

			// affNameB := encoder.Serialize(fmt.Sprintf("%s.%d", exprLbl, i))
			// affNameOffset := AllocateSeq(len(affNameB))
			affNameOffset := types.AllocWrite_str_data(ast.PROGRAM.Memory, fmt.Sprintf("%s.%d", exprLbl, i))
			// WriteMemory(affNameOffset, affNameB)

			var affNameOffsetBytes [4]byte
			types.Write_ptr(affNameOffsetBytes[:], 0, affNameOffset)
			*affOffset = ast.WriteToSlice(*affOffset, affNameOffsetBytes[:])
		}
	}
}

// QueryArgument ...
func QueryArgument(fn *ast.CXFunction, expr *ast.CXExpression, argOffsetB []byte, affOffset *types.Pointer) {
	for _, ex := range expr.Function.Expressions {
		if ex.Label == "" {
			// it's a non-labeled expression
			continue
		}

		queryParam(fn, ex.Inputs, ex.Label+".Input", argOffsetB, affOffset)
		queryParam(fn, ex.Outputs, ex.Label+".Output", argOffsetB, affOffset)
	}
}

// QueryExpressions ...
func QueryExpressions(fn *ast.CXFunction, expr *ast.CXExpression, exprOffsetB []byte, affOffset *types.Pointer) {
	for _, ex := range expr.Function.Expressions {
		if ex.Operator == nil || ex.Label == "" {
			// then it's a variable declaration
			// or it's a non-labeled expression
			continue
		}

		// var opNameB []byte
		opNameOffset := types.Pointer(0)
		if ex.Operator.IsBuiltin {
			// opNameB = encoder.Serialize(OpNames[ex.Operator.OpCode])
			opNameOffset = types.AllocWrite_str_data(ast.PROGRAM.Memory, ast.OpNames[ex.Operator.OpCode])
		} else {
			// opNameB = encoder.Serialize(ex.Operator.Name)
			opNameOffset = types.AllocWrite_str_data(ast.PROGRAM.Memory, ex.Operator.Name)
		}

		// opNameOffset := AllocateSeq(len(opNameB))
		// WriteMemory(opNameOffset, opNameB)
		var opNameOffsetB [4]byte
		types.Write_ptr(opNameOffsetB[:], 0, opNameOffset)
		res := CallAffPredicate(fn, opNameOffsetB[:])

		if res == 1 {
			*affOffset = ast.WriteToSlice(*affOffset, exprOffsetB)

			// lblNameB := encoder.Serialize(ex.Label)
			// lblNameOffset := AllocateSeq(len(lblNameB))
			lblNameOffset := types.AllocWrite_str_data(ast.PROGRAM.Memory, ex.Label)
			// WriteMemory(lblNameOffset, lblNameB)
			var lblNameOffsetB [4]byte
			types.Write_ptr(lblNameOffsetB[:], 0, lblNameOffset)
			*affOffset = ast.WriteToSlice(*affOffset, lblNameOffsetB[:])
		}
	}
}

func getSignatureSlice(params []*ast.CXArgument) types.Pointer {
	var sliceOffset types.Pointer
	for _, param := range params {

		var typOffset types.Pointer
		if param.StructType != nil {
			// typOffset = WriteObjectRetOff(encoder.Serialize(param.StructType.Package.Name + "." + param.StructType.Name))
			typOffset = types.AllocWrite_str_data(ast.PROGRAM.Memory, param.StructType.Package.Name+"."+param.StructType.Name)
		} else {
			// then it's native type
			// typOffset = WriteObjectRetOff(encoder.Serialize(TypeNames[param.Type]))
			typOffset = types.AllocWrite_str_data(ast.PROGRAM.Memory, param.Type.Name())
		}

		var typOffsetB [types.POINTER_SIZE]byte
		types.Write_ptr(typOffsetB[:], 0, typOffset)
		sliceOffset = ast.WriteToSlice(sliceOffset, typOffsetB[:])
	}

	return sliceOffset
}

// Helper function for QueryStructure. Used to query all the structs in a particular package
func queryStructsInPackage(fn *ast.CXFunction, strctOffsetB []byte, affOffset *types.Pointer, pkg *ast.CXPackage) {
	for _, f := range pkg.Structs {
		// strctNameB := encoder.Serialize(f.Name)

		// strctNameOffset := WriteObjectRetOff(strctNameB)
		strctNameOffset := types.AllocWrite_str_data(ast.PROGRAM.Memory, f.Name)
		var strctNameOffsetB [types.POINTER_SIZE]byte
		types.Write_ptr(strctNameOffsetB[:], 0, strctNameOffset)

		strctOffset := ast.AllocateSeq(types.OBJECT_HEADER_SIZE + types.STR_SIZE)
		// Name
		types.WriteSlice_byte(ast.PROGRAM.Memory, strctOffset+types.OBJECT_HEADER_SIZE, strctNameOffsetB[:])

		val := ast.PROGRAM.Memory[strctOffset+types.OBJECT_HEADER_SIZE : strctOffset+types.OBJECT_HEADER_SIZE+types.STR_SIZE]
		res := CallAffPredicate(fn, val)

		if res == 1 {
			*affOffset = ast.WriteToSlice(*affOffset, strctOffsetB)
			*affOffset = ast.WriteToSlice(*affOffset, strctNameOffsetB[:])
		}
	}
}

// QueryStructure ...
func QueryStructure(fn *ast.CXFunction, expr *ast.CXExpression, strctOffsetB []byte, affOffset *types.Pointer) {
	queryStructsInPackage(fn, strctOffsetB, affOffset, expr.Package)
	for _, imp := range expr.Package.Imports {
		queryStructsInPackage(fn, strctOffsetB, affOffset, imp)
	}
}

// QueryFunction ...
func QueryFunction(fn *ast.CXFunction, expr *ast.CXExpression, fnOffsetB []byte, affOffset *types.Pointer) {
	for _, f := range expr.Package.Functions {
		if f.Name == constants.SYS_INIT_FUNC {
			continue
		}

		// var opNameB []byte
		opNameOffset := types.Pointer(0)
		if f.IsBuiltin {
			// opNameB = encoder.Serialize(OpNames[f.OpCode])
			opNameOffset = types.AllocWrite_str_data(ast.PROGRAM.Memory, ast.OpNames[f.OpCode])
		} else {
			// opNameB = encoder.Serialize(f.Name)
			opNameOffset = types.AllocWrite_str_data(ast.PROGRAM.Memory, f.Name)
		}

		var opNameOffsetB [types.POINTER_SIZE]byte
		// WriteMemI32(opNameOffsetB[:], 0, int32(WriteObjectRetOff(opNameB)))
		types.Write_ptr(opNameOffsetB[:], 0, opNameOffset)

		inpSigOffset := getSignatureSlice(f.Inputs)
		outSigOffset := getSignatureSlice(f.Outputs)

		fnOffset := ast.AllocateSeq(types.OBJECT_HEADER_SIZE + types.STR_SIZE + types.POINTER_SIZE + types.POINTER_SIZE)
		// Name
		types.WriteSlice_byte(ast.PROGRAM.Memory, fnOffset+types.OBJECT_HEADER_SIZE, opNameOffsetB[:])
		// InputSignature
		types.Write_ptr(ast.PROGRAM.Memory, fnOffset+types.OBJECT_HEADER_SIZE+types.POINTER_SIZE, inpSigOffset)
		// OutputSignature
		types.Write_ptr(ast.PROGRAM.Memory, fnOffset+types.OBJECT_HEADER_SIZE+types.POINTER_SIZE+types.POINTER_SIZE, outSigOffset)

		val := ast.PROGRAM.Memory[fnOffset+types.OBJECT_HEADER_SIZE : fnOffset+types.OBJECT_HEADER_SIZE+types.STR_SIZE+types.POINTER_SIZE+types.POINTER_SIZE]
		res := CallAffPredicate(fn, val)

		if res == 1 {
			*affOffset = ast.WriteToSlice(*affOffset, fnOffsetB)
			*affOffset = ast.WriteToSlice(*affOffset, opNameOffsetB[:])
		}
	}
}

// QueryCaller ...
func QueryCaller(fn *ast.CXFunction, expr *ast.CXExpression, callerOffsetB []byte, affOffset *types.Pointer) {
	if ast.PROGRAM.CallCounter == 0 {
		// then it's entry point
		return
	}

	call := ast.PROGRAM.CallStack[ast.PROGRAM.CallCounter-1]

	// var opNameB []byte
	opNameOffset := types.Pointer(0)
	if call.Operator.IsBuiltin {
		// opNameB = encoder.Serialize(OpNames[call.Operator.OpCode])
		opNameOffset = types.AllocWrite_str_data(ast.PROGRAM.Memory, ast.OpNames[call.Operator.OpCode])
	} else {
		// opNameB = encoder.Serialize(call.Operator.Package.Name + "." + call.Operator.Name)
		opNameOffset = types.AllocWrite_str_data(ast.PROGRAM.Memory, call.Operator.Package.Name+"."+call.Operator.Name)
	}

	callOffset := ast.AllocateSeq(types.OBJECT_HEADER_SIZE + types.STR_SIZE + types.I32_SIZE)

	// FnName
	var opNameOffsetB [4]byte
	// WriteMemI32(opNameOffsetB[:], 0, int32(WriteObjectRetOff(opNameB)))
	types.Write_ptr(opNameOffsetB[:], 0, opNameOffset)
	types.WriteSlice_byte(ast.PROGRAM.Memory, callOffset+types.OBJECT_HEADER_SIZE, opNameOffsetB[:])

	// FnSize
	types.Write_ptr(ast.PROGRAM.Memory, callOffset+types.OBJECT_HEADER_SIZE+types.STR_SIZE, call.Operator.Size)

	res := CallAffPredicate(fn, ast.PROGRAM.Memory[callOffset+types.OBJECT_HEADER_SIZE:callOffset+types.OBJECT_HEADER_SIZE+types.STR_SIZE+types.I32_SIZE])

	if res == 1 {
		*affOffset = ast.WriteToSlice(*affOffset, callerOffsetB)
	}
}

// QueryProgram ...
func QueryProgram(fn *ast.CXFunction, expr *ast.CXExpression, prgrmOffsetB []byte, affOffset *types.Pointer) {
	prgrmOffset := ast.AllocateSeq(types.OBJECT_HEADER_SIZE + types.I32_SIZE + types.I64_SIZE + types.STR_SIZE + types.I32_SIZE)
	// Callcounter
	types.Write_ptr(ast.PROGRAM.Memory, prgrmOffset+types.OBJECT_HEADER_SIZE, ast.PROGRAM.CallCounter)
	// HeapUsed
	types.Write_ptr(ast.PROGRAM.Memory, prgrmOffset+types.OBJECT_HEADER_SIZE+types.I32_SIZE, ast.PROGRAM.HeapPointer)

	// Caller
	if ast.PROGRAM.CallCounter != 0 {
		// then it's not just entry point
		call := ast.PROGRAM.CallStack[ast.PROGRAM.CallCounter-1]

		// var opNameB []byte
		opNameOffset := types.Pointer(0)
		if call.Operator.IsBuiltin {
			// opNameB = encoder.Serialize(OpNames[call.Operator.OpCode])
			opNameOffset = types.AllocWrite_str_data(ast.PROGRAM.Memory, ast.OpNames[call.Operator.OpCode])
		} else {
			// opNameB = encoder.Serialize(call.Operator.Package.Name + "." + call.Operator.Name)
			opNameOffset = types.AllocWrite_str_data(ast.PROGRAM.Memory, call.Operator.Package.Name+"."+call.Operator.Name)
		}

		// callOffset := AllocateSeq(OBJECT_HEADER_SIZE + STR_SIZE + I32_SIZE)
		// FnName
		var opNameOffsetB [4]byte
		// WriteMemI32(opNameOffsetB[:], 0, int32(WriteObjectRetOff(opNameB)))
		types.Write_ptr(opNameOffsetB[:], 0, opNameOffset)
		types.WriteSlice_byte(ast.PROGRAM.Memory, prgrmOffset+types.OBJECT_HEADER_SIZE+types.I32_SIZE+types.I64_SIZE, opNameOffsetB[:])
		// FnSize
		types.Write_ptr(ast.PROGRAM.Memory, prgrmOffset+types.OBJECT_HEADER_SIZE+types.I32_SIZE+types.I64_SIZE+types.STR_SIZE, call.Operator.Size)

		// res := CallAffPredicate(fn, PROGRAM.Memory[callOffset + OBJECT_HEADER_SIZE : callOffset + OBJECT_HEADER_SIZE + STR_SIZE + I32_SIZE])

		// if res == 1 {
		// 	*affOffset = WriteToSlice(*affOffset, callerOffsetB)
		// }
	}

	res := CallAffPredicate(fn, ast.PROGRAM.Memory[prgrmOffset+types.OBJECT_HEADER_SIZE:prgrmOffset+types.OBJECT_HEADER_SIZE+types.I32_SIZE+types.I64_SIZE+types.STR_SIZE+types.I32_SIZE])

	if res == 1 {
		*affOffset = ast.WriteToSlice(*affOffset, prgrmOffsetB)
		*affOffset = ast.WriteToSlice(*affOffset, prgrmOffsetB)
	}
}

func getTarget(inp2 *ast.CXArgument, fp types.Pointer, tgtElt *string, tgtArgType *string, tgtArgIndex *int,
	tgtPkg *ast.CXPackage, tgtFn *ast.CXFunction, tgtExpr *ast.CXExpression) {
	for _, aff := range GetInferActions(inp2, fp) {
		switch aff {
		case "prgrm":
			*tgtElt = "prgrm"
		case "Pkg":
			*tgtElt = "Pkg"
		case "strct":
			*tgtElt = "strct"
		case "fn":
			*tgtElt = "fn"
		case "expr":
			*tgtElt = "expr"
		case "rec":
			*tgtElt = "rec"
		case "inp":
			*tgtElt = "inp"
		case "out":
			*tgtElt = "out"
		default:
			switch *tgtElt {
			case "Pkg":
				if pkg, err := ast.PROGRAM.GetPackage(aff); err == nil {
					*tgtPkg = *pkg
				} else {
					panic(err)
				}
			case "fn":
				if fn, err := tgtPkg.GetFunction(aff); err == nil {
					*tgtFn = *fn
				} else {
					panic(err)
				}
			case "expr":
				if expr, err := tgtFn.GetExpressionByLabel(aff); err == nil {
					*tgtExpr = *expr
				} else {
					panic(err)
				}
			case "inp":
				*tgtArgType = "inp"
				*tgtElt = "arg"
				i, err := strconv.ParseInt(aff, 10, 32)

				*tgtArgIndex = int(i)

				if err != nil {
					panic(err)
				}
			case "out":
				*tgtArgType = "out"
				*tgtElt = "arg"
				i, err := strconv.ParseInt(aff, 10, 32)

				*tgtArgIndex = int(i)

				if err != nil {
					panic(err)
				}
			}
		}
	}
}

func getAffordances(inp1 *ast.CXArgument, fp types.Pointer,
	tgtElt string, tgtArgType string, tgtArgIndex int,
	tgtPkg *ast.CXPackage, tgtFn *ast.CXFunction, tgtExpr *ast.CXExpression,
	affMsgs map[string]string,
	affs *[]string) {
	var fltrElt string
	elts := GetInferActions(inp1, fp)
	// for _, elt := range elts {
	for c := 0; c < len(elts); c++ {
		elt := elts[c]
		switch elt {
		case "arg":
			fltrElt = "arg"
		case "expr":
			fltrElt = "expr"
		case "fn":
			fltrElt = "fn"
		case "strct":
			fltrElt = "strct"
		case "Pkg":
			fltrElt = "Pkg"
		case "prgrm":
			fltrElt = "prgrm"
			// skipping the extra "prgrm"
			c++
			switch tgtElt {
			case "arg":
				if tgtArgType == "inp" {
					if msg, ok := affMsgs["prgrm-arg-input"]; ok {
						*affs = append(*affs, fmt.Sprintf(msg, tgtExpr.Label, tgtArgIndex))
					}
				} else {
					if msg, ok := affMsgs["prgrm-arg-output"]; ok {
						*affs = append(*affs, fmt.Sprintf(msg, tgtExpr.Label, tgtArgIndex))
					}
				}
			case "prgrm":
				*affs = append(*affs, "Run program")
			}
		default:
			switch fltrElt {
			case "arg":
				switch tgtElt {
				case "arg":
					if tgtArgType == "inp" {
						if msg, ok := affMsgs["arg-arg-input"]; ok {
							*affs = append(*affs, fmt.Sprintf(msg, tgtExpr.Label, tgtArgIndex, elt))
						}
					} else {
						if msg, ok := affMsgs["arg-arg-output"]; ok {
							*affs = append(*affs, fmt.Sprintf(msg, tgtExpr.Label, tgtArgIndex, elt))
						}
					}
				case "prgrm":
					*affs = append(*affs, "Print FA's value")
				}
			case "expr":
				if expr, err := tgtFn.GetExpressionByLabel(elt); err == nil {
					_ = expr
					switch tgtElt {
					case "arg":
						*affs = append(*affs, "Replace TA by FE")
					case "fn":
						*affs = append(*affs, "Add FE to TF")
					case "prgrm":
						*affs = append(*affs, "Call FE")
					}
				} else {
					panic(err)
				}
			case "fn":
				if fn, err := tgtPkg.GetFunction(elt); err == nil {
					_ = fn
					switch tgtElt {
					case "arg":
						*affs = append(*affs, "Replace TA by a call to FF")
					case "expr":
						*affs = append(*affs, "Change TE's operator to FF")
					case "Pkg":
						// affs = append(affs, fmt.Sprintf("[%s.Operator = %s]", tgtExpr.Label, fn.Name))
						*affs = append(*affs, "Move FF to TP")
					case "prgrm":
						*affs = append(*affs, "Call FF")
					}
				} else {
					panic(err)
				}
			case "strct":
				switch tgtElt {
				case "arg":
					if msg, ok := affMsgs["arg-strct"]; ok {
						*affs = append(*affs, fmt.Sprintf(msg, tgtExpr.Label, tgtArgIndex, elt))
					}
					if tgtArgType == "inp" {
						if msg, ok := affMsgs["strct-arg-input"]; ok {
							*affs = append(*affs, fmt.Sprintf(msg, tgtExpr.Label, tgtArgIndex, elt))
						}
					} else {
						if msg, ok := affMsgs["strct-arg-output"]; ok {
							*affs = append(*affs, fmt.Sprintf(msg, tgtExpr.Label, tgtArgIndex, elt))
						}
					}
				case "fn":
					*affs = append(*affs, "Add or change TF's receiver to FS")
				case "Pkg":
					*affs = append(*affs, "Move FS to TP")
				}
			case "Pkg":
				if pkg, err := ast.PROGRAM.GetPackage(elt); err == nil {
					_ = pkg
					switch tgtElt {
					case "Pkg":
						*affs = append(*affs, "Make TP import FP")
					}
				} else {
					panic(err)
				}
				// case "prgrm":
				// 	switch tgtElt {
				// 	case "prgrm":
				// 		affs = append(affs, "Run program")
				// 	}
			}
		}
	}
}

func opAffOn(inputs []ast.CXValue, outputs []ast.CXValue) {
	inp1, inp2 := inputs[0].Arg, inputs[1].Arg

	prevPkg := ast.PROGRAM.CurrentPackage
	prevFn := prevPkg.CurrentFunction
	prevExpr := prevFn.CurrentExpression

	call := ast.PROGRAM.GetCurrentCall()
	expr := call.Operator.Expressions[call.Line]
	fp := inputs[0].FramePointer

	var tgtPkg = ast.CXPackage(*prevPkg)
	var tgtFn = ast.CXFunction(*expr.Function)
	var tgtExpr = ast.CXExpression(*prevExpr)

	// processing the target
	var tgtElt string
	var tgtArgType string
	var tgtArgIndex int

	getTarget(inp2, fp, &tgtElt, &tgtArgType, &tgtArgIndex, &tgtPkg, &tgtFn, &tgtExpr)

	// var affPkg *CXPackage = prevPkg
	// var affFn *CXFunction = prevFn
	// var affExpr *CXExpression = prevExpr

	// processing the affordances
	var affs []string
	getAffordances(inp1, fp, tgtElt, tgtArgType, tgtArgIndex, &tgtPkg, &tgtFn, &tgtExpr, onMessages, &affs)

	// returning to previous state
	ast.PROGRAM.CurrentPackage = prevPkg
	ast.PROGRAM.CurrentPackage.CurrentFunction = prevFn
	ast.PROGRAM.CurrentPackage.CurrentFunction.CurrentExpression = prevExpr

	for i, aff := range affs {
		fmt.Printf("%d - %s\n", i, aff)
	}
}

func opAffOf(inputs []ast.CXValue, outputs []ast.CXValue) {
	inp1, inp2 := inputs[0].Arg, inputs[1].Arg

	prevPkg := ast.PROGRAM.CurrentPackage
	prevFn := prevPkg.CurrentFunction
	prevExpr := prevFn.CurrentExpression

	call := ast.PROGRAM.GetCurrentCall()
	expr := call.Operator.Expressions[call.Line]
	fp := inputs[0].FramePointer

	var tgtPkg = ast.CXPackage(*expr.Package)
	var tgtFn = ast.CXFunction(*expr.Function)
	var tgtExpr = ast.CXExpression(*prevExpr)

	// processing the target
	var tgtElt string
	var tgtArgType string
	var tgtArgIndex int

	getTarget(inp2, fp, &tgtElt, &tgtArgType, &tgtArgIndex, &tgtPkg, &tgtFn, &tgtExpr)

	// processing the affordances
	var affs []string
	getAffordances(inp1, fp, tgtElt, tgtArgType, tgtArgIndex, &tgtPkg, &tgtFn, &tgtExpr, ofMessages, &affs)

	// returning to previous state
	ast.PROGRAM.CurrentPackage = prevPkg
	ast.PROGRAM.CurrentPackage.CurrentFunction = prevFn
	ast.PROGRAM.CurrentPackage.CurrentFunction.CurrentExpression = prevExpr

	for i, aff := range affs {
		fmt.Printf("%d - %s\n", i, aff)
	}
}

func readStrctAff(aff string, tgtPkg *ast.CXPackage) *ast.CXStruct {
	strct, err := tgtPkg.GetStruct(aff)
	if err != nil {
		panic(err)
	}

	return strct
}

func readArgAff(aff string, tgtFn *ast.CXFunction) *ast.CXArgument {
	var affExpr *ast.CXExpression
	var lIdx int
	var rIdx int
	var ch rune

	for _, ch = range aff {
		if ch == '.' {
			exprLbl := aff[lIdx:rIdx]
			for _, expr := range tgtFn.Expressions {
				if exprLbl == expr.Label {
					affExpr = expr
					rIdx++
					break
				}
			}
		}

		if affExpr != nil {
			break
		}

		rIdx++
	}

	lIdx = rIdx

	var argType string

	for _, ch = range aff[lIdx:] {
		if ch == '.' {
			argType = aff[lIdx:rIdx]
			rIdx++
			break
		}

		if argType != "" {
			break
		}

		rIdx++
	}

	lIdx = rIdx

	argIdx, err := strconv.ParseInt(aff[lIdx:], 10, 32)
	if err != nil {
		panic(err)
	}

	if argType == "Input" {
		return affExpr.Inputs[argIdx]
	}
	return affExpr.Outputs[argIdx]

}

func opAffInform(inputs []ast.CXValue, outputs []ast.CXValue) {
	inp1, inp2, inp3 := inputs[0].Arg, inputs[1].Arg, inputs[2].Arg

	call := ast.PROGRAM.GetCurrentCall()
	expr := call.Operator.Expressions[call.Line]
	fp := inputs[0].FramePointer

	prevPkg := ast.PROGRAM.CurrentPackage
	prevFn := prevPkg.CurrentFunction
	prevExpr := prevFn.CurrentExpression

	var tgtPkg = ast.CXPackage(*prevPkg)
	var tgtFn = ast.CXFunction(*expr.Function)
	var tgtExpr = ast.CXExpression(*prevExpr)

	// processing the target
	var tgtElt string
	var tgtArgType string
	var tgtArgIndex int

	getTarget(inp3, fp, &tgtElt, &tgtArgType, &tgtArgIndex, &tgtPkg, &tgtFn, &tgtExpr)

	elts := GetInferActions(inp1, fp)
	eltIdx := types.Read_ptr(ast.PROGRAM.Memory, ast.GetFinalOffset(fp, inp2))
	eltType := elts[eltIdx*2]
	elt := elts[eltIdx*2+1]

	switch eltType {
	case "arg":
		switch tgtElt {
		case "arg":
			if tgtArgType == "inp" {
				tgtExpr.Inputs[tgtArgIndex] = readArgAff(elt, &tgtFn)
			} else {
				tgtExpr.Outputs[tgtArgIndex] = readArgAff(elt, &tgtFn)
			}
		case "strct":

		case "prgrm":

		}
	case "expr":
		if expr, err := tgtFn.GetExpressionByLabel(elt); err == nil {
			_ = expr
			switch tgtElt {
			case "arg":

			case "fn":

			case "prgrm":

			}
		} else {
			panic(err)
		}
	case "fn":
		if fn, err := tgtPkg.GetFunction(elt); err == nil {
			_ = fn
			switch tgtElt {
			case "arg":

			case "expr":

			case "Pkg":

			case "prgrm":

			}
		} else {
			panic(err)
		}
	case "strct":
		switch tgtElt {
		case "arg":

		case "fn":

		case "Pkg":

		}
	case "Pkg":
		if pkg, err := ast.PROGRAM.GetPackage(elt); err == nil {
			_ = pkg
			switch tgtElt {
			case "Pkg":

			}
		} else {
			panic(err)
		}
		// case "prgrm":
		// 	switch tgtElt {
		// 	case "prgrm":
		// 		affs = append(affs, "Run program")
		// 	}
	}

	// returning to previous state
	ast.PROGRAM.CurrentPackage = prevPkg
	ast.PROGRAM.CurrentPackage.CurrentFunction = prevFn
	ast.PROGRAM.CurrentPackage.CurrentFunction.CurrentExpression = prevExpr
}

func opAffRequest(inputs []ast.CXValue, outputs []ast.CXValue) {
	inp1, inp2, inp3 := inputs[0].Arg, inputs[1].Arg, inputs[2].Arg

	call := ast.PROGRAM.GetCurrentCall()
	expr := call.Operator.Expressions[call.Line]
	fp := inputs[0].FramePointer

	prevPkg := ast.PROGRAM.CurrentPackage
	prevFn := prevPkg.CurrentFunction
	prevExpr := prevFn.CurrentExpression

	var tgtPkg = ast.CXPackage(*prevPkg)
	var tgtFn = ast.CXFunction(*expr.Function)
	var tgtExpr = ast.CXExpression(*prevExpr)

	// processing the target
	var tgtElt string
	var tgtArgType string
	var tgtArgIndex int

	getTarget(inp3, fp, &tgtElt, &tgtArgType, &tgtArgIndex, &tgtPkg, &tgtFn, &tgtExpr)

	// var affs []string

	elts := GetInferActions(inp1, fp)
	eltIdx := types.Read_ptr(ast.PROGRAM.Memory, ast.GetFinalOffset(fp, inp2))
	eltType := elts[eltIdx*2]
	elt := elts[eltIdx*2+1]

	switch eltType {
	case "arg":
		switch tgtElt {
		case "arg":
			if tgtArgType == "inp" {
				// tgtExpr.ProgramInput[tgtArgIndex] = readArgAff(elt, &tgtFn)
				*readArgAff(elt, &tgtFn) = *tgtExpr.Inputs[tgtArgIndex]
			} else {
				// tgtExpr.ProgramOutput[tgtArgIndex] = readArgAff(elt, &tgtFn)
				*readArgAff(elt, &tgtFn) = *tgtExpr.Outputs[tgtArgIndex]
			}
		case "strct":

		case "prgrm":
			fmt.Println(ast.GetPrintableValue(fp, readArgAff(elt, &tgtFn)))
		}
	case "expr":
		if expr, err := tgtFn.GetExpressionByLabel(elt); err == nil {
			_ = expr
			switch tgtElt {
			case "arg":

			case "fn":

			case "prgrm":

			}
		} else {
			panic(err)
		}
	case "fn":
		fn := ast.Natives[ast.OpCodes[elt]]
		if fn == nil {
			var err error
			fn, err = tgtPkg.GetFunction(elt)
			if err != nil {
				panic(err)
			}
		}
		_ = fn
		switch tgtElt {
		case "arg":

		case "expr":

		case "Pkg":

		case "prgrm":

		}
	case "strct":
		switch tgtElt {
		case "arg":
			if tgtArgType == "inp" {
				// tgtExpr.ProgramInput[tgtArgIndex] = readArgAff(elt, &tgtFn)
				readStrctAff(elt, &tgtPkg).AddField(tgtExpr.Inputs[tgtArgIndex])
			} else {
				// tgtExpr.ProgramOutput[tgtArgIndex] = readArgAff(elt, &tgtFn)
				readStrctAff(elt, &tgtPkg).AddField(tgtExpr.Outputs[tgtArgIndex])
			}
		case "fn":

		case "Pkg":

		}
	case "Pkg":
		if pkg, err := ast.PROGRAM.GetPackage(elt); err == nil {
			_ = pkg
			switch tgtElt {
			case "Pkg":

			}
		} else {
			panic(err)
		}
	case "prgrm":
		switch tgtElt {
		case "arg":
			if tgtArgType == "inp" {
				fmt.Println(ast.GetPrintableValue(fp, tgtExpr.Inputs[tgtArgIndex]))
			} else {
				fmt.Println(ast.GetPrintableValue(fp, tgtExpr.Outputs[tgtArgIndex]))
			}
		case "prgrm":
			// affs = append(affs, "Run program")
		}
	}

	// returning to previous state
	ast.PROGRAM.CurrentPackage = prevPkg
	ast.PROGRAM.CurrentPackage.CurrentFunction = prevFn
	ast.PROGRAM.CurrentPackage.CurrentFunction.CurrentExpression = prevExpr
}

func opAffQuery(inputs []ast.CXValue, outputs []ast.CXValue) {
	inp1, out1 := inputs[0].Arg, outputs[0].Arg

	call := ast.PROGRAM.GetCurrentCall()
	expr := call.Operator.Expressions[call.Line]
	fp := inputs[0].FramePointer

	out1Offset := ast.GetFinalOffset(fp, out1)

	var affOffset types.Pointer

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
				if fn, err := inp1.ArgDetails.Package.GetFunction(rule); err == nil {

					// arg keyword
					// argB := encoder.Serialize("arg")
					// argOffset := AllocateSeq(len(argB))
					// WriteMemory(argOffset, argB)
					argOffset := types.AllocWrite_str_data(ast.PROGRAM.Memory, "arg")
					var argOffsetB [4]byte
					types.Write_ptr(argOffsetB[:], 0, argOffset)

					// expr keyword
					// exprB := encoder.Serialize("expr")
					// exprOffset := AllocateSeq(len(exprB))
					// WriteMemory(exprOffset, exprB)
					exprOffset := types.AllocWrite_str_data(ast.PROGRAM.Memory, "expr")
					var exprOffsetB [4]byte
					types.Write_ptr(exprOffsetB[:], 0, exprOffset)

					// fn keyword
					// fnB := encoder.Serialize("fn")
					// fnOffset := AllocateSeq(len(fnB))
					// WriteMemory(fnOffset, fnB)
					fnOffset := types.AllocWrite_str_data(ast.PROGRAM.Memory, "fn")
					var fnOffsetB [4]byte
					types.Write_ptr(fnOffsetB[:], 0, fnOffset)

					// strct keyword
					// strctB := encoder.Serialize("strct")
					// strctOffset := AllocateSeq(len(strctB))
					// WriteMemory(strctOffset, strctB)
					strctOffset := types.AllocWrite_str_data(ast.PROGRAM.Memory, "strct")
					var strctOffsetB [4]byte
					types.Write_ptr(strctOffsetB[:], 0, strctOffset)

					// caller keyword
					// callerB := encoder.Serialize("caller")
					// callerOffset := AllocateSeq(len(callerB))
					// WriteMemory(callerOffset, callerB)
					callerOffset := types.AllocWrite_str_data(ast.PROGRAM.Memory, "caller")
					var callerOffsetB [4]byte
					types.Write_ptr(callerOffsetB[:], 0, callerOffset)

					// program keyword
					// prgrmB := encoder.Serialize("prgrm")
					// prgrmOffset := AllocateSeq(len(prgrmB))
					// WriteMemory(prgrmOffset, prgrmB)
					prgrmOffset := types.AllocWrite_str_data(ast.PROGRAM.Memory, "prgrm")
					var prgrmOffsetB [4]byte
					types.Write_ptr(prgrmOffsetB[:], 0, prgrmOffset)

					predInp := fn.Inputs[0]

					if predInp.Type == types.STRUCT {
						if predInp.StructType != nil {
							switch predInp.StructType.Name {
							case "Argument":
								QueryArgument(fn, expr, argOffsetB[:], &affOffset)
							case "Expression":
								QueryExpressions(fn, expr, exprOffsetB[:], &affOffset)
							case "Function":
								QueryFunction(fn, expr, fnOffsetB[:], &affOffset)
							case "Structure":
								QueryStructure(fn, expr, strctOffsetB[:], &affOffset)
							case "Caller":
								QueryCaller(fn, expr, callerOffsetB[:], &affOffset)
							case "Program":
								QueryProgram(fn, expr, prgrmOffsetB[:], &affOffset)
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

	types.Write_ptr(ast.PROGRAM.Memory, out1Offset, affOffset)
}
