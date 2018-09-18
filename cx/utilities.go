package base

import (
	// "os"
	// "path/filepath"
	"bytes"
	// "errors"
	"fmt"
	"github.com/skycoin/skycoin/src/cipher/encoder"
	// "math/rand"
	// "regexp"
	"strings"
	// "time"
)

func Debug (args ...interface{}) {
	fmt.Println(args...)
}

func sameFields(flds1 []*CXArgument, flds2 []*CXArgument) bool {
	allSame := true

	if len(flds1) != len(flds2) {
		allSame = false
	} else {
		for i, fld := range flds1 {
			if flds2[i].Name != fld.Name {
				allSame = false
			}
		}
	}

	return allSame
}

func getArrayChunkSizes(size int, lens []int) []int {
	var result []int

	for c := len(lens) - 1; c >= 0; c-- {
		if len(result) > 0 {
			result = append([]int{lens[c] * result[len(result)-1]}, result...)
		} else {
			// first one to add
			result = append(result, lens[c]*size)
		}
	}

	return result
}

func IsArray(typ string) bool {
	if len(typ) > 2 && typ[:2] == "[]" {
		return true
	}
	return false
}
func IsStructInstance(typ string, mod *CXPackage) bool {
	if _, err := mod.Program.GetStruct(typ, mod.Name); err == nil {
		return true
	} else {
		return false
	}
}
// func IsLocal(identName string, call *CXCall) bool {
// 	for _, def := range call.State {
// 		if def.Name == identName {
// 			return true
// 		}
// 	}
// 	return false
// }
func IsGlobal(identName string, mod *CXPackage) bool {
	for _, def := range mod.Globals {
		if def.Name == identName {
			return true
		}
	}
	for _, imp := range mod.Imports {
		for _, def := range imp.Globals {
			if def.Name == identName {
				return true
			}
		}
	}
	return false
}

// It returns true if the operator receives undefined types as input parameters but also an operator that needs to mimic its input's type. For example, == should not return its input type, as it is always going to return a boolean
func IsUndOp (fn *CXFunction) bool {
	res := false
	switch fn.OpCode {
		case
		OP_UND_BITAND,
		OP_UND_BITXOR,
		OP_UND_BITOR,
		OP_UND_BITCLEAR,
		OP_UND_MUL,
		OP_UND_DIV,
		OP_UND_MOD,
		OP_UND_ADD,
		OP_UND_SUB,
		OP_UND_BITSHL, OP_UND_BITSHR:
		res = true
	}
	
	return res
}

func (prgrm *CXProgram) PrintStack() {
	fmt.Println()
	fmt.Println("===Stack===")

	fp := 0

	for c := 0; c <= prgrm.CallCounter; c++ {
		op := prgrm.CallStack[c].Operator

		var dupNames []string

		fmt.Println(">>>", op.Name, "()")

		for _, inp := range op.Inputs {
			fmt.Println("Inputs")
			fmt.Println("\t", inp.Name, "\t", ":", "\t", prgrm.Memory[inp.Offset:inp.Offset+inp.TotalSize])

			dupNames = append(dupNames, inp.Package.Name+inp.Name)
		}

		for _, out := range op.Outputs {
			fmt.Println("Outputs")
			fmt.Println("\t", out.Name, "\t", ":", "\t", prgrm.Memory[out.Offset:out.Offset+out.TotalSize])

			dupNames = append(dupNames, out.Package.Name+out.Name)
		}

		fmt.Println("Expressions")

		for _, expr := range op.Expressions {
			for _, inp := range expr.Inputs {
				if inp.Name == "" || expr.Operator == nil {
					continue
				}
				var dup bool
				for _, name := range dupNames {
					if name == inp.Package.Name+inp.Name {
						dup = true
						break
					}
				}
				if dup {
					continue
				}

				fmt.Println("\t", inp.Name, "\t", ":", "\t", prgrm.Memory[inp.Offset:inp.Offset+inp.TotalSize])

				dupNames = append(dupNames, inp.Package.Name+inp.Name)
			}

			for _, out := range expr.Outputs {
				if out.Name == "" || expr.Operator == nil {
					continue
				}
				var dup bool
				for _, name := range dupNames {
					if name == out.Package.Name+out.Name {
						dup = true
						break
					}
				}
				if dup {
					continue
				}

				fmt.Println("\t", out.Name, "\t", ":", "\t", prgrm.Memory[out.Offset:out.Offset+out.TotalSize])

				dupNames = append(dupNames, out.Package.Name+out.Name)
			}
		}

		fp += op.Size
	}
	fmt.Println()
}

func PrintCallStack(callStack []CXCall) {
	for i, call := range callStack {
		tabs := strings.Repeat("___", i)
		if tabs == "" {
			//fmt.Printf("%sfn:%s ln:%d, \tlocals: ", tabs, call.Operator.Name, call.Line)
			fmt.Printf("%sfn:%s ln:%d", tabs, call.Operator.Name, call.Line)
		} else {
			//fmt.Printf("↓%sfn:%s ln:%d, \tlocals: ", tabs, call.Operator.Name, call.Line)
			fmt.Printf("↓%sfn:%s ln:%d", tabs, call.Operator.Name, call.Line)
		}

		fmt.Println()
	}
}

func (prgrm *CXProgram) PrintProgram() {
	fmt.Println("Program")

	var currentFunction *CXFunction
	var currentPackage *CXPackage

	_ = currentFunction
	_ = currentPackage

	// saving current program state because PrintProgram uses SelectXXX
	if pkg, err := prgrm.GetCurrentPackage(); err == nil {
		currentPackage = pkg
	}

	if fn, err := prgrm.GetCurrentFunction(); err == nil {
		currentFunction = fn
	}

	i := 0
	for _, mod := range prgrm.Packages {
		if mod.Name == CORE_MODULE || mod.Name == "glfw" || mod.Name == "gl" || mod.Name == "gltext" {
			continue
		}

		fmt.Printf("%d.- Package: %s\n", i, mod.Name)

		if len(mod.Imports) > 0 {
			fmt.Println("\tImports")
		}

		j := 0
		for _, imp := range mod.Imports {
			fmt.Printf("\t\t%d.- Import: %s\n", j, imp.Name)
			j++
		}

		if len(mod.Globals) > 0 {
			fmt.Println("\tGlobals")
		}

		j = 0
		for _, v := range mod.Globals {
			var arrayStr string
			if v.IsArray {
				for _, l := range v.Lengths {
					arrayStr += fmt.Sprintf("[%d]", l)
				}
			}
			fmt.Printf("\t\t%d.- Global: %s %s%s\n", j, v.Name, arrayStr, TypeNames[v.Type])
			j++
		}

		if len(mod.Structs) > 0 {
			fmt.Println("\tStructs")
		}

		j = 0
		for _, strct := range mod.Structs {
			fmt.Printf("\t\t%d.- Struct: %s\n", j, strct.Name)

			for k, fld := range strct.Fields {
				fmt.Printf("\t\t\t%d.- Field: %s %s\n",
					k, fld.Name, TypeNames[fld.Type])
			}

			j++
		}

		if len(mod.Functions) > 0 {
			fmt.Println("\tFunctions")
		}

		j = 0
		for _, fn := range mod.Functions {
			mod.SelectFunction(fn.Name)

			var inps bytes.Buffer
			for i, inp := range fn.Inputs {
				if i == len(fn.Inputs)-1 {
					inps.WriteString(fmt.Sprintf("%s %s", inp.Name, TypeNames[inp.Type]))
				} else {
					inps.WriteString(fmt.Sprintf("%s %s, ", inp.Name, TypeNames[inp.Type]))
				}
			}

			var outs bytes.Buffer
			for i, out := range fn.Outputs {
				if i == len(fn.Outputs)-1 {
					outs.WriteString(fmt.Sprintf("%s %s", out.Name, TypeNames[out.Type]))
				} else {
					outs.WriteString(fmt.Sprintf("%s %s, ", out.Name, TypeNames[out.Type]))
				}
			}

			fmt.Printf("\t\t%d.- Function: %s (%s) (%s)\n",
				j, fn.Name, inps.String(), outs.String())

			k := 0
			for _, expr := range fn.Expressions {
				if expr.Operator == nil {
					continue
				}
				//Arguments
				var args bytes.Buffer

				for i, arg := range expr.Inputs {
					var name string
					var dat []byte

					if arg.Offset > STACK_SIZE {
						dat = prgrm.Memory[arg.Offset : arg.Offset+arg.Size]
					} else {
						name = arg.Name
					}

					if dat != nil {
						switch TypeNames[arg.Type] {
						case "str":
							encoder.DeserializeRaw(dat, &name)
							name = "\"" + name + "\""
						case "i32":
							var i32 int32
							encoder.DeserializeAtomic(dat, &i32)
							name = fmt.Sprintf("%v", i32)
						case "i64":
							var i64 int64
							encoder.DeserializeRaw(dat, &i64)
							name = fmt.Sprintf("%v", i64)
						case "f32":
							var f32 float32
							encoder.DeserializeRaw(dat, &f32)
							name = fmt.Sprintf("%v", f32)
						case "f64":
							var f64 float64
							encoder.DeserializeRaw(dat, &f64)
							name = fmt.Sprintf("%v", f64)
						case "bool":
							var b bool
							encoder.DeserializeRaw(dat, &b)
							name = fmt.Sprintf("%v", b)
						case "byte":
							var b bool
							encoder.DeserializeRaw(dat, &b)
							name = fmt.Sprintf("%v", b)
						}
					}

					if arg.Name != "" {
						name = arg.Name
					}

					var arrayStr string
					if arg.IsArray {
						for _, l := range arg.Lengths {
							arrayStr += fmt.Sprintf("[%d]", l)
						}
					}
					
					if i == len(expr.Inputs)-1 {
						
						args.WriteString(fmt.Sprintf("%s %s%s", name, arrayStr, TypeNames[arg.Type]))
						// args.WriteString(TypeNames[arg.Type])
					} else {
						args.WriteString(fmt.Sprintf("%s %s%s, ", name, arrayStr, TypeNames[arg.Type]))
						// args.WriteString(TypeNames[arg.Type] + ", ")
					}
				}

				var opName string
				if expr.Operator.IsNative {
					opName = OpNames[expr.Operator.OpCode]
				} else {
					opName = expr.Operator.Name
				}

				if len(expr.Outputs) > 0 {
					var outNames bytes.Buffer
					for i, outName := range expr.Outputs {
						if i == len(expr.Outputs)-1 {
							outNames.WriteString(fmt.Sprintf("%s %s", outName.Name, TypeNames[outName.Type]))
						} else {
							outNames.WriteString(fmt.Sprintf("%s %s", outName.Name, TypeNames[outName.Type]))
						}
					}

					var lbl string
					if expr.Label != "" {
						lbl = " <<" + expr.Label + ">>"
					} else {
						lbl = ""
					}

					if expr.Operator != nil {
						fmt.Printf("\t\t\t%d.- Expression%s: %s = %s(%s)\n",
							k,
							lbl,
							outNames.String(),
							opName,
							args.String(),
						)
					}

				} else {
					var lbl string

					if expr.Label != "" {
						lbl = " <<" + expr.Label + ">>"
					} else {
						lbl = ""
					}

					fmt.Printf("\t\t\t%d.- Expression%s: %s(%s)\n",
						k,
						lbl,
						opName,
						args.String(),
					)
				}
				k++
			}
			j++
		}
		i++
	}

	if currentPackage != nil {
		prgrm.SelectPackage(currentPackage.Name)
	}
	if currentFunction != nil {
		prgrm.SelectFunction(currentFunction.Name)
	}

	prgrm.CurrentPackage = currentPackage
	currentPackage.CurrentFunction = currentFunction
}

// this function adds the roots (pointers) for some GC algorithms
func AddPointer(fn *CXFunction, sym *CXArgument) {
	if sym.IsPointer {
		var found bool
		for _, ptr := range fn.ListOfPointers {
			if sym.Name == ptr.Name {
				found = true
				break
			}
		}
		if !found {
			fn.ListOfPointers = append(fn.ListOfPointers, sym)
		}
	}
}

func CheckArithmeticOp(expr *CXExpression) bool {
	if expr.Operator.IsNative {
		switch expr.Operator.OpCode {
		case OP_I32_MUL, OP_I32_DIV, OP_I32_MOD, OP_I32_ADD,
			OP_I32_SUB, OP_I32_BITSHL, OP_I32_BITSHR, OP_I32_LT,
			OP_I32_GT, OP_I32_LTEQ, OP_I32_GTEQ, OP_I32_EQ, OP_I32_UNEQ,
			OP_I32_BITAND, OP_I32_BITXOR, OP_I32_BITOR, OP_STR_EQ:
			return true
		}
	}
	return false
}

func IsCorePackage (ident string) bool {
	for _, core := range CorePackages {
		if core == ident {
			return true
		}
	}
	return false
}

func SetCorrectArithmeticOp(expr *CXExpression) {
	if expr.Operator == nil || len(expr.Outputs) < 1 {
		return
	}
	op := expr.Operator
	typ := expr.Outputs[0].Type

	if CheckArithmeticOp(expr) {
		// if !CheckSameNativeType(expr) {
		// 	panic("wrong types")
		// }
		switch op.OpCode {
		case OP_I32_MUL:
			switch typ {
			case TYPE_I32:
			case TYPE_I64:
				expr.Operator = Natives[OP_I64_MUL]
			case TYPE_F32:
				expr.Operator = Natives[OP_F32_MUL]
			case TYPE_F64:
				expr.Operator = Natives[OP_F64_MUL]
			}
		case OP_I32_DIV:
			switch typ {
			case TYPE_I32:
			case TYPE_I64:
				expr.Operator = Natives[OP_I64_DIV]
			case TYPE_F32:
				expr.Operator = Natives[OP_F32_DIV]
			case TYPE_F64:
				expr.Operator = Natives[OP_F64_DIV]
			}
		case OP_I32_MOD:
			switch typ {
			case TYPE_I32:
			case TYPE_I64:
				expr.Operator = Natives[OP_I64_MOD]
			}

		case OP_I32_ADD:
			switch typ {
			case TYPE_I32:
			case TYPE_I64:
				expr.Operator = Natives[OP_I64_ADD]
			case TYPE_F32:
				expr.Operator = Natives[OP_F32_ADD]
			case TYPE_F64:
				expr.Operator = Natives[OP_F64_ADD]
			}
		case OP_I32_SUB:
			switch typ {
			case TYPE_I32:
			case TYPE_I64:
				expr.Operator = Natives[OP_I64_ADD]
			case TYPE_F32:
				expr.Operator = Natives[OP_F32_ADD]
			case TYPE_F64:
				expr.Operator = Natives[OP_F64_ADD]
			}

		case OP_I32_BITSHL:
			switch typ {
			case TYPE_I32:
			case TYPE_I64:
				expr.Operator = Natives[OP_I64_BITSHL]
			}
		case OP_I32_BITSHR:
			switch typ {
			case TYPE_I32:
			case TYPE_I64:
				expr.Operator = Natives[OP_I64_BITSHR]
			}

		case OP_I32_LT:
			switch typ {
			case TYPE_I32:
			case TYPE_I64:
				expr.Operator = Natives[OP_I64_LT]
			case TYPE_F32:
				expr.Operator = Natives[OP_F32_LT]
			case TYPE_F64:
				expr.Operator = Natives[OP_F64_LT]
			}
		case OP_I32_GT:
			switch typ {
			case TYPE_I32:
			case TYPE_I64:
				expr.Operator = Natives[OP_I64_GT]
			case TYPE_F32:
				expr.Operator = Natives[OP_F32_GT]
			case TYPE_F64:
				expr.Operator = Natives[OP_F64_GT]
			}
		case OP_I32_LTEQ:
			switch typ {
			case TYPE_I32:
			case TYPE_I64:
				expr.Operator = Natives[OP_I64_LTEQ]
			case TYPE_F32:
				expr.Operator = Natives[OP_F32_LTEQ]
			case TYPE_F64:
				expr.Operator = Natives[OP_F64_LTEQ]
			}
		case OP_I32_GTEQ:
			switch typ {
			case TYPE_I32:
			case TYPE_I64:
				expr.Operator = Natives[OP_I64_GTEQ]
			case TYPE_F32:
				expr.Operator = Natives[OP_F32_GTEQ]
			case TYPE_F64:
				expr.Operator = Natives[OP_F64_GTEQ]
			}

		case OP_I32_EQ:
			switch typ {
			case TYPE_I32:
			case TYPE_I64:
				expr.Operator = Natives[OP_I64_EQ]
			case TYPE_F32:
				expr.Operator = Natives[OP_F32_EQ]
			case TYPE_F64:
				expr.Operator = Natives[OP_F64_EQ]
			}
		case OP_I32_UNEQ:
			switch typ {
			case TYPE_I32:
			case TYPE_I64:
				expr.Operator = Natives[OP_I64_UNEQ]
			case TYPE_F32:
				expr.Operator = Natives[OP_F32_UNEQ]
			case TYPE_F64:
				expr.Operator = Natives[OP_F64_UNEQ]
			}

		case OP_I32_BITAND:
			switch typ {
			case TYPE_I32:
			case TYPE_I64:
				expr.Operator = Natives[OP_I64_BITAND]
			}

		case OP_I32_BITXOR:
			switch typ {
			case TYPE_I32:
			case TYPE_I64:
				expr.Operator = Natives[OP_I64_BITXOR]
			}

		case OP_I32_BITOR:
			switch typ {
			case TYPE_I32:
			case TYPE_I64:
				expr.Operator = Natives[OP_I64_BITOR]
			}
		}
	}
}

func GetArgSize(typ int) int {
	switch typ {
	case TYPE_BOOL, TYPE_BYTE:
		return 1
	case TYPE_STR, TYPE_I32, TYPE_F32, TYPE_AFF:
		return 4
	case TYPE_I64, TYPE_F64:
		return 8
	default:
		return 4
	}
}

func MakeMultiDimArray(atomicSize int, lengths []int) []byte {
	var result []byte

	fstDLen := lengths[len(lengths)-1]

	sLen := encoder.SerializeAtomic(int32(fstDLen))

	byts := append(sLen, make([]byte, fstDLen*atomicSize)...)
	result = byts

	if len(lengths) > 1 {
		// -2 to ignore the first dimension
		for c := len(lengths) - 2; c >= 0; c-- {
			lenB := encoder.SerializeAtomic(int32(lengths[c]))

			var tmp []byte

			for i := 0; i < lengths[c]; i++ {
				tmp = append(tmp, result...)
			}

			result = append(lenB, tmp...)
		}
	}

	return result
}

func checkForEscapedChars (str string) []byte {
	var res []byte
	var lenStr int = len(str)
	for c := 0; c < len(str); c++ {
		var nextCh byte
		ch := str[c]
		if c < lenStr - 1{
			nextCh = str[c+1]
		}
		if ch == '\\' {
			switch nextCh {
			case '%':
				c++
				res = append(res, nextCh)
				continue
			case 'n':
				c++
				res = append(res, '\n')
				continue
			default:
				res = append(res, ch)
				continue
			}

		} else {
			res = append(res, ch)
		}
	}

	return res
}

func GetAssignmentElement (arg *CXArgument) *CXArgument {
	if len(arg.Fields) > 0 {
		return arg.Fields[len(arg.Fields) - 1]
	} else {
		return arg
	}
}

func WriteObject (out1Offset int, byts []byte) {
	size := encoder.Serialize(int32(len(byts)))
	heapOffset := AllocateSeq(len(byts) + OBJECT_HEADER_SIZE)
	
	var header []byte = make([]byte, OBJECT_HEADER_SIZE, OBJECT_HEADER_SIZE)
	for c := 5; c < OBJECT_HEADER_SIZE; c++ {
		header[c] = size[c-5]
	}

	obj := append(header, byts...)

	WriteMemory(heapOffset, obj)

	off := encoder.SerializeAtomic(int32(heapOffset + OBJECT_HEADER_SIZE))

	WriteMemory(out1Offset, off)
}
