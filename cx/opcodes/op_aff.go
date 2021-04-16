package opcodes

import (
	"fmt"
	"strconv"

	"github.com/skycoin/cx/cx/ast"
	"github.com/skycoin/cx/cx/constants"
	"github.com/skycoin/cx/cx/helper"
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
func GetInferActions(inp *ast.CXArgument, fp int) []string {
	inpOffset := ast.GetFinalOffset(fp, inp)

	off := helper.Deserialize_i32(ast.PROGRAM.Memory[inpOffset : inpOffset+constants.TYPE_POINTER_SIZE])

	l := helper.Deserialize_i32(ast.GetSliceHeader(ast.GetSliceOffset(fp, inp))[4:8])

	result := make([]string, l)

	// for c := int(l); c > 0; c-- {
	for c := 0; c < int(l); c++ {
		// elof := Deserialize_i32(PROGRAM.Memory[int(off) + OBJECT_HEADER_SIZE + SLICE_HEADER_SIZE + (c - 1) * TYPE_POINTER_SIZE : int(off) + OBJECT_HEADER_SIZE + SLICE_HEADER_SIZE + c * STR_HEADER_SIZE])
		elOff := helper.Deserialize_i32(ast.PROGRAM.Memory[int(off)+constants.OBJECT_HEADER_SIZE+constants.SLICE_HEADER_SIZE+c*constants.TYPE_POINTER_SIZE : int(off)+constants.OBJECT_HEADER_SIZE+constants.SLICE_HEADER_SIZE+(c+1)*constants.STR_HEADER_SIZE])
		// size := Deserialize_i32(PROGRAM.Memory[elOff : elOff+STR_HEADER_SIZE])
		// var res string
		// _, err := encoder.DeserializeRaw(PROGRAM.Memory[elOff:elOff+STR_HEADER_SIZE+size], &res)
		// if err != nil {
		// 	panic(err)
		// }

		// result[int(l) - c] = res
		result[c] = ast.ReadStringFromObject(elOff)
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
	for c := 0; c < fn.Size; c++ {
		ast.PROGRAM.Memory[newFP+c] = 0
	}

	// sending value to predicate function
	ast.WriteMemory(
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

	return ast.ReadMemory(ast.GetFinalOffset(
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
func queryParam(fn *ast.CXFunction, args []*ast.CXArgument, exprLbl string, argOffsetB []byte, affOffset *int) {
	for i, arg := range args {

		var typOffset int
		elt := ast.GetAssignmentElement(arg)
		if elt.CustomType != nil {
			// then it's custom type
			// typOffset = WriteObjectRetOff(encoder.Serialize(elt.CustomType.Package.Name + "." + elt.CustomType.Name))
			typOffset = ast.WriteStringData(elt.CustomType.Package.Name + "." + elt.CustomType.Name)
		} else {
			// then it's native type
			// typOffset = WriteObjectRetOff(encoder.Serialize(TypeNames[elt.Type]))
			typOffset = ast.WriteStringData(constants.TypeNames[elt.Type])
		}

		// Name
		// argNameB := encoder.Serialize(arg.Name)
		// argNameOffset := int32(WriteObjectRetOff(argNameB))
		argNameOffset := ast.WriteStringData(arg.ArgDetails.Name)

		argOffset := ast.AllocateSeq(constants.OBJECT_HEADER_SIZE + constants.STR_SIZE + constants.I32_SIZE + constants.STR_SIZE)
		ast.WriteI32(argOffset+constants.OBJECT_HEADER_SIZE, int32(argNameOffset))

		// Index
		ast.WriteI32(argOffset+constants.OBJECT_HEADER_SIZE+constants.STR_SIZE, int32(i))

		// Type
		ast.WriteI32(argOffset+constants.OBJECT_HEADER_SIZE+constants.STR_SIZE+constants.I32_SIZE, int32(typOffset))

		res := CallAffPredicate(fn, ast.PROGRAM.Memory[argOffset+constants.OBJECT_HEADER_SIZE:argOffset+constants.OBJECT_HEADER_SIZE+constants.STR_SIZE+constants.I32_SIZE+constants.STR_SIZE])

		if res == 1 {
			*affOffset = ast.WriteToSlice(*affOffset, argOffsetB)

			// affNameB := encoder.Serialize(fmt.Sprintf("%s.%d", exprLbl, i))
			// affNameOffset := AllocateSeq(len(affNameB))
			affNameOffset := ast.WriteStringData(fmt.Sprintf("%s.%d", exprLbl, i))
			// WriteMemory(affNameOffset, affNameB)

			var affNameOffsetBytes [4]byte
			ast.WriteMemI32(affNameOffsetBytes[:], 0, int32(affNameOffset))
			*affOffset = ast.WriteToSlice(*affOffset, affNameOffsetBytes[:])
		}
	}
}

// QueryArgument ...
func QueryArgument(fn *ast.CXFunction, expr *ast.CXExpression, argOffsetB []byte, affOffset *int) {
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
func QueryExpressions(fn *ast.CXFunction, expr *ast.CXExpression, exprOffsetB []byte, affOffset *int) {
	for _, ex := range expr.Function.Expressions {
		if ex.Operator == nil || ex.Label == "" {
			// then it's a variable declaration
			// or it's a non-labeled expression
			continue
		}

		// var opNameB []byte
		opNameOffset := 0
		if ex.Operator.IsBuiltin {
			// opNameB = encoder.Serialize(OpNames[ex.Operator.OpCode])
			opNameOffset = ast.WriteStringData(ast.OpNames[ex.Operator.OpCode])
		} else {
			// opNameB = encoder.Serialize(ex.Operator.Name)
			opNameOffset = ast.WriteStringData(ex.Operator.Name)
		}

		// opNameOffset := AllocateSeq(len(opNameB))
		// WriteMemory(opNameOffset, opNameB)
		var opNameOffsetB [4]byte
		ast.WriteMemI32(opNameOffsetB[:], 0, int32(opNameOffset))
		res := CallAffPredicate(fn, opNameOffsetB[:])

		if res == 1 {
			*affOffset = ast.WriteToSlice(*affOffset, exprOffsetB)

			// lblNameB := encoder.Serialize(ex.Label)
			// lblNameOffset := AllocateSeq(len(lblNameB))
			lblNameOffset := ast.WriteStringData(ex.Label)
			// WriteMemory(lblNameOffset, lblNameB)
			var lblNameOffsetB [4]byte
			ast.WriteMemI32(lblNameOffsetB[:], 0, int32(lblNameOffset))
			*affOffset = ast.WriteToSlice(*affOffset, lblNameOffsetB[:])
		}
	}
}

func getSignatureSlice(params []*ast.CXArgument) int {
	var sliceOffset int
	for _, param := range params {

		var typOffset int
		if param.CustomType != nil {
			// then it's custom type
			// typOffset = WriteObjectRetOff(encoder.Serialize(param.CustomType.Package.Name + "." + param.CustomType.Name))
			typOffset = ast.WriteStringData(param.CustomType.Package.Name + "." + param.CustomType.Name)
		} else {
			// then it's native type
			// typOffset = WriteObjectRetOff(encoder.Serialize(TypeNames[param.Type]))
			typOffset = ast.WriteStringData(constants.TypeNames[param.Type])
		}

		var typOffsetB [4]byte
		ast.WriteMemI32(typOffsetB[:], 0, int32(typOffset))
		sliceOffset = ast.WriteToSlice(sliceOffset, typOffsetB[:])
	}

	return sliceOffset
}

// Helper function for QueryStructure. Used to query all the structs in a particular package
func queryStructsInPackage(fn *ast.CXFunction, strctOffsetB []byte, affOffset *int, pkg *ast.CXPackage) {
	for _, f := range pkg.Structs {
		// strctNameB := encoder.Serialize(f.Name)

		// strctNameOffset := WriteObjectRetOff(strctNameB)
		strctNameOffset := ast.WriteStringData(f.Name)
		var strctNameOffsetB [4]byte
		ast.WriteMemI32(strctNameOffsetB[:], 0, int32(strctNameOffset))

		strctOffset := ast.AllocateSeq(constants.OBJECT_HEADER_SIZE + constants.STR_SIZE)
		// Name
		ast.WriteMemory(strctOffset+constants.OBJECT_HEADER_SIZE, strctNameOffsetB[:])

		val := ast.PROGRAM.Memory[strctOffset+constants.OBJECT_HEADER_SIZE : strctOffset+constants.OBJECT_HEADER_SIZE+constants.STR_SIZE]
		res := CallAffPredicate(fn, val)

		if res == 1 {
			*affOffset = ast.WriteToSlice(*affOffset, strctOffsetB)
			*affOffset = ast.WriteToSlice(*affOffset, strctNameOffsetB[:])
		}
	}
}

// QueryStructure ...
func QueryStructure(fn *ast.CXFunction, expr *ast.CXExpression, strctOffsetB []byte, affOffset *int) {
	queryStructsInPackage(fn, strctOffsetB, affOffset, expr.Package)
	for _, imp := range expr.Package.Imports {
		queryStructsInPackage(fn, strctOffsetB, affOffset, imp)
	}
}

// QueryFunction ...
func QueryFunction(fn *ast.CXFunction, expr *ast.CXExpression, fnOffsetB []byte, affOffset *int) {
	for _, f := range expr.Package.Functions {
		if f.Name == constants.SYS_INIT_FUNC {
			continue
		}

		// var opNameB []byte
		opNameOffset := 0
		if f.IsBuiltin {
			// opNameB = encoder.Serialize(OpNames[f.OpCode])
			opNameOffset = ast.WriteStringData(ast.OpNames[f.OpCode])
		} else {
			// opNameB = encoder.Serialize(f.Name)
			opNameOffset = ast.WriteStringData(f.Name)
		}

		var opNameOffsetB [4]byte
		// WriteMemI32(opNameOffsetB[:], 0, int32(WriteObjectRetOff(opNameB)))
		ast.WriteMemI32(opNameOffsetB[:], 0, int32(opNameOffset))

		inpSigOffset := getSignatureSlice(f.Inputs)
		outSigOffset := getSignatureSlice(f.Outputs)

		fnOffset := ast.AllocateSeq(constants.OBJECT_HEADER_SIZE + constants.STR_SIZE + constants.TYPE_POINTER_SIZE + constants.TYPE_POINTER_SIZE)
		// Name
		ast.WriteMemory(fnOffset+constants.OBJECT_HEADER_SIZE, opNameOffsetB[:])
		// InputSignature
		ast.WriteI32(fnOffset+constants.OBJECT_HEADER_SIZE+constants.TYPE_POINTER_SIZE, int32(inpSigOffset))
		// OutputSignature
		ast.WriteI32(fnOffset+constants.OBJECT_HEADER_SIZE+constants.TYPE_POINTER_SIZE+constants.TYPE_POINTER_SIZE, int32(outSigOffset))

		val := ast.PROGRAM.Memory[fnOffset+constants.OBJECT_HEADER_SIZE : fnOffset+constants.OBJECT_HEADER_SIZE+constants.STR_SIZE+constants.TYPE_POINTER_SIZE+constants.TYPE_POINTER_SIZE]
		res := CallAffPredicate(fn, val)

		if res == 1 {
			*affOffset = ast.WriteToSlice(*affOffset, fnOffsetB)
			*affOffset = ast.WriteToSlice(*affOffset, opNameOffsetB[:])
		}
	}
}

// QueryCaller ...
func QueryCaller(fn *ast.CXFunction, expr *ast.CXExpression, callerOffsetB []byte, affOffset *int) {
	if ast.PROGRAM.CallCounter == 0 {
		// then it's entry point
		return
	}

	call := ast.PROGRAM.CallStack[ast.PROGRAM.CallCounter-1]

	// var opNameB []byte
	opNameOffset := 0
	if call.Operator.IsBuiltin {
		// opNameB = encoder.Serialize(OpNames[call.Operator.OpCode])
		opNameOffset = ast.WriteStringData(ast.OpNames[call.Operator.OpCode])
	} else {
		// opNameB = encoder.Serialize(call.Operator.Package.Name + "." + call.Operator.Name)
		opNameOffset = ast.WriteStringData(call.Operator.Package.Name + "." + call.Operator.Name)
	}

	callOffset := ast.AllocateSeq(constants.OBJECT_HEADER_SIZE + constants.STR_SIZE + constants.I32_SIZE)

	// FnName
	var opNameOffsetB [4]byte
	// WriteMemI32(opNameOffsetB[:], 0, int32(WriteObjectRetOff(opNameB)))
	ast.WriteMemI32(opNameOffsetB[:], 0, int32(opNameOffset))
	ast.WriteMemory(callOffset+constants.OBJECT_HEADER_SIZE, opNameOffsetB[:])

	// FnSize
	ast.WriteI32(callOffset+constants.OBJECT_HEADER_SIZE+constants.STR_SIZE, int32(call.Operator.Size))

	res := CallAffPredicate(fn, ast.PROGRAM.Memory[callOffset+constants.OBJECT_HEADER_SIZE:callOffset+constants.OBJECT_HEADER_SIZE+constants.STR_SIZE+constants.I32_SIZE])

	if res == 1 {
		*affOffset = ast.WriteToSlice(*affOffset, callerOffsetB)
	}
}

// QueryProgram ...
func QueryProgram(fn *ast.CXFunction, expr *ast.CXExpression, prgrmOffsetB []byte, affOffset *int) {
	prgrmOffset := ast.AllocateSeq(constants.OBJECT_HEADER_SIZE + constants.I32_SIZE + constants.I64_SIZE + constants.STR_SIZE + constants.I32_SIZE)
	// Callcounter
	ast.WriteI32(prgrmOffset+constants.OBJECT_HEADER_SIZE, int32(ast.PROGRAM.CallCounter))
	// HeapUsed
	ast.WriteI64(prgrmOffset+constants.OBJECT_HEADER_SIZE+constants.I32_SIZE, int64(ast.PROGRAM.HeapPointer))

	// Caller
	if ast.PROGRAM.CallCounter != 0 {
		// then it's not just entry point
		call := ast.PROGRAM.CallStack[ast.PROGRAM.CallCounter-1]

		// var opNameB []byte
		opNameOffset := 0
		if call.Operator.IsBuiltin {
			// opNameB = encoder.Serialize(OpNames[call.Operator.OpCode])
			opNameOffset = ast.WriteStringData(ast.OpNames[call.Operator.OpCode])
		} else {
			// opNameB = encoder.Serialize(call.Operator.Package.Name + "." + call.Operator.Name)
			opNameOffset = ast.WriteStringData(call.Operator.Package.Name + "." + call.Operator.Name)
		}

		// callOffset := AllocateSeq(OBJECT_HEADER_SIZE + STR_SIZE + I32_SIZE)
		// FnName
		var opNameOffsetB [4]byte
		// WriteMemI32(opNameOffsetB[:], 0, int32(WriteObjectRetOff(opNameB)))
		ast.WriteMemI32(opNameOffsetB[:], 0, int32(opNameOffset))
		ast.WriteMemory(prgrmOffset+constants.OBJECT_HEADER_SIZE+constants.I32_SIZE+constants.I64_SIZE, opNameOffsetB[:])
		// FnSize
		ast.WriteI32(prgrmOffset+constants.OBJECT_HEADER_SIZE+constants.I32_SIZE+constants.I64_SIZE+constants.STR_SIZE, int32(call.Operator.Size))

		// res := CallAffPredicate(fn, PROGRAM.Memory[callOffset + OBJECT_HEADER_SIZE : callOffset + OBJECT_HEADER_SIZE + STR_SIZE + I32_SIZE])

		// if res == 1 {
		// 	*affOffset = WriteToSlice(*affOffset, callerOffsetB)
		// }
	}

	res := CallAffPredicate(fn, ast.PROGRAM.Memory[prgrmOffset+constants.OBJECT_HEADER_SIZE:prgrmOffset+constants.OBJECT_HEADER_SIZE+constants.I32_SIZE+constants.I64_SIZE+constants.STR_SIZE+constants.I32_SIZE])

	if res == 1 {
		*affOffset = ast.WriteToSlice(*affOffset, prgrmOffsetB)
		*affOffset = ast.WriteToSlice(*affOffset, prgrmOffsetB)
	}
}

func getTarget(inp2 *ast.CXArgument, fp int, tgtElt *string, tgtArgType *string, tgtArgIndex *int,
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

func getAffordances(inp1 *ast.CXArgument, fp int,
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
	eltIdx := ast.ReadI32(fp, inp2)
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
	eltIdx := ast.ReadI32(fp, inp2)
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
				if fn, err := inp1.ArgDetails.Package.GetFunction(rule); err == nil {

					// arg keyword
					// argB := encoder.Serialize("arg")
					// argOffset := AllocateSeq(len(argB))
					// WriteMemory(argOffset, argB)
					argOffset := ast.WriteStringData("arg")
					var argOffsetB [4]byte
					ast.WriteMemI32(argOffsetB[:], 0, int32(argOffset))

					// expr keyword
					// exprB := encoder.Serialize("expr")
					// exprOffset := AllocateSeq(len(exprB))
					// WriteMemory(exprOffset, exprB)
					exprOffset := ast.WriteStringData("expr")
					var exprOffsetB [4]byte
					ast.WriteMemI32(exprOffsetB[:], 0, int32(exprOffset))

					// fn keyword
					// fnB := encoder.Serialize("fn")
					// fnOffset := AllocateSeq(len(fnB))
					// WriteMemory(fnOffset, fnB)
					fnOffset := ast.WriteStringData("fn")
					var fnOffsetB [4]byte
					ast.WriteMemI32(fnOffsetB[:], 0, int32(fnOffset))

					// strct keyword
					// strctB := encoder.Serialize("strct")
					// strctOffset := AllocateSeq(len(strctB))
					// WriteMemory(strctOffset, strctB)
					strctOffset := ast.WriteStringData("strct")
					var strctOffsetB [4]byte
					ast.WriteMemI32(strctOffsetB[:], 0, int32(strctOffset))

					// caller keyword
					// callerB := encoder.Serialize("caller")
					// callerOffset := AllocateSeq(len(callerB))
					// WriteMemory(callerOffset, callerB)
					callerOffset := ast.WriteStringData("caller")
					var callerOffsetB [4]byte
					ast.WriteMemI32(callerOffsetB[:], 0, int32(callerOffset))

					// program keyword
					// prgrmB := encoder.Serialize("prgrm")
					// prgrmOffset := AllocateSeq(len(prgrmB))
					// WriteMemory(prgrmOffset, prgrmB)
					prgrmOffset := ast.WriteStringData("prgrm")
					var prgrmOffsetB [4]byte
					ast.WriteMemI32(prgrmOffsetB[:], 0, int32(prgrmOffset))

					predInp := fn.Inputs[0]

					if predInp.Type == constants.TYPE_CUSTOM {
						if predInp.CustomType != nil {
							switch predInp.CustomType.Name {
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

	ast.WriteI32(out1Offset, int32(affOffset))
}
