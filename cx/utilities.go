package base

import (
	"bytes"
	"fmt"
	"os"
	"runtime/debug"
	"strconv"

	"github.com/skycoin/skycoin/src/cipher/encoder"
)

func Debug(args ...interface{}) {
	fmt.Println(args...)
}

// It returns true if the operator receives undefined types as input parameters but also an operator that needs to mimic its input's type. For example, == should not return its input type, as it is always going to return a boolean
func IsUndOp(fn *CXFunction) bool {
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

func ExprOpName(expr *CXExpression) string {
	if expr.Operator.IsNative {
		return OpNames[expr.Operator.OpCode]
	} else {
		return expr.Operator.Name
	}
}

// func limitString (str string) string {
// 	if len(str) > 3
// }

func stackValueHeader(fileName string, fileLine int) string {
	return fmt.Sprintf("%s:%d", fileName, fileLine)
}

func (prgrm *CXProgram) PrintStack() {
	fmt.Println()
	fmt.Println("===Callstack===")

	// we're going backwards in the stack
	fp := prgrm.StackPointer

	for c := prgrm.CallCounter; c >= 0; c-- {
		op := prgrm.CallStack[c].Operator
		fp -= op.Size

		var dupNames []string

		fmt.Printf(">>> %s()\n", op.Name)

		for _, inp := range op.Inputs {
			fmt.Println("Inputs")
			fmt.Printf("\t%s : %s() : %s\n", stackValueHeader(inp.FileName, inp.FileLine), op.Name, GetPrintableValue(fp, inp))

			dupNames = append(dupNames, inp.Package.Name+inp.Name)
		}

		for _, out := range op.Outputs {
			fmt.Println("Outputs")
			fmt.Printf("\t%s : %s() : %s\n", stackValueHeader(out.FileName, out.FileLine), op.Name, GetPrintableValue(fp, out))

			dupNames = append(dupNames, out.Package.Name+out.Name)
		}

		// fmt.Println("Expressions")
		exprs := ""
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

				// fmt.Println("\t", inp.Name, "\t", ":", "\t", GetPrintableValue(fp, inp))
				// exprs += fmt.Sprintln("\t", stackValueHeader(inp.FileName, inp.FileLine), "\t", ":", "\t", GetPrintableValue(fp, inp))

				exprs += fmt.Sprintf("\t%s : %s() : %s\n", stackValueHeader(inp.FileName, inp.FileLine), ExprOpName(expr), GetPrintableValue(fp, inp))

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

				// fmt.Println("\t", out.Name, "\t", ":", "\t", GetPrintableValue(fp, out))
				// exprs += fmt.Sprintln("\t", stackValueHeader(out.FileName, out.FileLine), ":", GetPrintableValue(fp, out))

				exprs += fmt.Sprintf("\t%s : %s() : %s\n", stackValueHeader(out.FileName, out.FileLine), ExprOpName(expr), GetPrintableValue(fp, out))

				dupNames = append(dupNames, out.Package.Name+out.Name)
			}
		}

		if len(exprs) > 0 {
			fmt.Println("Expressions\n", exprs)
		}
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
		if IsCorePackage(mod.Name) {
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
				var isPointer string
				if inp.IsPointer {
					isPointer = "*"
				}

				var arrayStr string
				if inp.IsArray {
					for _, l := range inp.Lengths {
						arrayStr += fmt.Sprintf("[%d]", l)
					}
				}

				var typeName string
				elt := GetAssignmentElement(inp)
				if elt.CustomType != nil {
					// then it's custom type
					typeName = elt.CustomType.Name
				} else {
					// then it's native type
					typeName = TypeNames[elt.Type]
				}

				if i == len(fn.Inputs)-1 {
					inps.WriteString(fmt.Sprintf("%s %s%s%s", inp.Name, isPointer, arrayStr, typeName))
				} else {
					inps.WriteString(fmt.Sprintf("%s %s%s%s, ", inp.Name, isPointer, arrayStr, typeName))
				}
			}

			var outs bytes.Buffer
			for i, out := range fn.Outputs {
				var isPointer string
				if out.IsPointer {
					isPointer = "*"
				}

				var arrayStr string
				if out.IsArray {
					for _, l := range out.Lengths {
						arrayStr += fmt.Sprintf("[%d]", l)
					}
				}

				var typeName string
				elt := GetAssignmentElement(out)
				if elt.CustomType != nil {
					// then it's custom type
					typeName = elt.CustomType.Name
				} else {
					// then it's native type
					typeName = TypeNames[elt.Type]
				}

				if i == len(fn.Outputs)-1 {
					outs.WriteString(fmt.Sprintf("%s %s%s%s", out.Name, isPointer, arrayStr, typeName))
				} else {
					outs.WriteString(fmt.Sprintf("%s %s%s%s, ", out.Name, isPointer, arrayStr, typeName))
				}
			}

			fmt.Printf("\t\t%d.- Function: %s (%s) (%s)\n",
				j, fn.Name, inps.String(), outs.String())

			k := 0
			for _, expr := range fn.Expressions {
				// if expr.Operator == nil {
				//      continue
				// }
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
						for _, fld := range arg.Fields {
							name += "." + fld.Name
						}
					}

					var derefLevels string
					if arg.DereferenceLevels > 0 {
						for c := 0; c < arg.DereferenceLevels; c++ {
							derefLevels += "*"
						}
					}

					var isReference string
					if arg.PassBy == PASSBY_REFERENCE {
						isReference = "&"
					}

					var arrayStr string
					if arg.IsArray {
						for _, l := range arg.Lengths {
							arrayStr += fmt.Sprintf("[%d]", l)
						}
					}

					var typeName string
					elt := GetAssignmentElement(arg)
					if elt.CustomType != nil {
						// then it's custom type
						typeName = elt.CustomType.Name
					} else {
						// then it's native type
						typeName = TypeNames[elt.Type]
					}

					if i == len(expr.Inputs)-1 {
						args.WriteString(fmt.Sprintf("%s%s%s %s%s", isReference, derefLevels, name, arrayStr, typeName))
					} else {
						args.WriteString(fmt.Sprintf("%s%s%s %s%s, ", isReference, derefLevels, name, arrayStr, typeName))
					}
				}

				var opName string
				if expr.Operator != nil {
					if expr.Operator.IsNative {
						opName = OpNames[expr.Operator.OpCode]
					} else {
						opName = expr.Operator.Name
					}
				}

				if len(expr.Outputs) > 0 {
					var outNames bytes.Buffer
					for i, outName := range expr.Outputs {
						out := GetAssignmentElement(outName)

						var derefLevels string
						if outName.DereferenceLevels > 0 {
							for c := 0; c < outName.DereferenceLevels; c++ {
								derefLevels += "*"
							}
						}

						var arrayStr string
						if outName.IsArray {
							for _, l := range outName.Lengths {
								arrayStr += fmt.Sprintf("[%d]", l)
							}
						}

						var typeName string
						if out.CustomType != nil {
							// then it's custom type
							typeName = out.CustomType.Name
						} else {
							// then it's native type
							typeName = TypeNames[out.Type]
						}

						fullName := outName.Name

						for _, fld := range outName.Fields {
							fullName += "." + fld.Name
						}

						if i == len(expr.Outputs)-1 {
							outNames.WriteString(fmt.Sprintf("%s%s%s %s", derefLevels, fullName, arrayStr, typeName))
						} else {
							outNames.WriteString(fmt.Sprintf("%s%s%s %s, ", derefLevels, fullName, arrayStr, typeName))
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
					} else {
						if len(expr.Outputs) > 0 {
							var typs string

							for i, out := range expr.Outputs {
								if GetAssignmentElement(out).CustomType != nil {
									// then it's custom type
									typs += GetAssignmentElement(out).CustomType.Name
								} else {
									// then it's native type
									typs += TypeNames[GetAssignmentElement(out).Type]
								}

								if i != len(expr.Outputs) {
									typs += ", "
								}
							}

							fmt.Printf("\t\t\t%d.- Declaration%s: %s %s\n",
								k,
								lbl,
								expr.Outputs[0].Name,
								typs)
						}
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

func IsCorePackage(ident string) bool {
	for _, core := range CorePackages {
		if core == ident {
			return true
		}
	}
	return false
}

func IsTempVar(name string) bool {
	if len(name) >= len(LOCAL_PREFIX) && name[:len(LOCAL_PREFIX)] == LOCAL_PREFIX {
		return true
	}
	return false
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

func checkForEscapedChars(str string) []byte {
	var res []byte
	var lenStr int = len(str)
	for c := 0; c < len(str); c++ {
		var nextCh byte
		ch := str[c]
		if c < lenStr-1 {
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

func GetAssignmentElement(arg *CXArgument) *CXArgument {
	if len(arg.Fields) > 0 {
		return arg.Fields[len(arg.Fields)-1]
	} else {
		return arg
	}
}

func WriteToSlice(off int, inp []byte) int {
	var heapOffset int

	if off == 0 {
		// then it's a new slice
		heapOffset = AllocateSeq(OBJECT_HEADER_SIZE + SLICE_HEADER_SIZE + len(inp))

		var header []byte = make([]byte, OBJECT_HEADER_SIZE)

		size := encoder.SerializeAtomic(int32(len(inp)) + SLICE_HEADER_SIZE)

		for c := 5; c < OBJECT_HEADER_SIZE; c++ {
			header[c] = size[c-5]
		}

		one := []byte{1, 0, 0, 0}

		// len == 1
		finalObj := append(header, one...)
		// cap == 1
		finalObj = append(finalObj, one...)
		finalObj = append(finalObj, inp...)

		WriteMemory(heapOffset, finalObj)
		return heapOffset
	} else {
		// then it already exists
		sliceHeader := PROGRAM.Memory[off+OBJECT_HEADER_SIZE : off+OBJECT_HEADER_SIZE+SLICE_HEADER_SIZE]

		var l int32
		var c int32

		encoder.DeserializeAtomic(sliceHeader[:4], &l)
		encoder.DeserializeAtomic(sliceHeader[4:], &c)

		if l >= c {
			// then we need to increase cap and relocate slice
			obj := PROGRAM.Memory[off+OBJECT_HEADER_SIZE+SLICE_HEADER_SIZE : int32(off)+OBJECT_HEADER_SIZE+SLICE_HEADER_SIZE+l*int32(len(inp))]

			l++
			c = c * 2

			heapOffset = AllocateSeq(int(c)*len(inp) + OBJECT_HEADER_SIZE + SLICE_HEADER_SIZE)

			size := encoder.SerializeAtomic(int32(int(c)*len(inp) + SLICE_HEADER_SIZE))

			var header []byte = make([]byte, OBJECT_HEADER_SIZE)
			for c := 5; c < OBJECT_HEADER_SIZE; c++ {
				header[c] = size[c-5]
			}

			lB := encoder.SerializeAtomic(l)
			cB := encoder.SerializeAtomic(c)

			finalObj := append(header, lB...)
			finalObj = append(finalObj, cB...)
			finalObj = append(finalObj, obj...)
			finalObj = append(finalObj, inp...)

			WriteMemory(heapOffset, finalObj)

			return heapOffset
		} else {
			// then we can simply write the element

			// updating the length
			newL := encoder.SerializeAtomic(l + int32(1))

			for i, byt := range newL {
				PROGRAM.Memory[int(off)+OBJECT_HEADER_SIZE+i] = byt
			}

			// write the obj
			for i, byt := range inp {
				PROGRAM.Memory[off+OBJECT_HEADER_SIZE+SLICE_HEADER_SIZE+int(l)*len(inp)+i] = byt
			}

			return off
		}
	}
}

// refactoring reuse in WriteObject and WriteObjectRetOff
func writeObj(obj []byte) int {
	size := len(obj)
	sizeB := encoder.SerializeAtomic(int32(size))
	// heapOffset := AllocateSeq(size + OBJECT_HEADER_SIZE + SLICE_HEADER_SIZE)
	heapOffset := AllocateSeq(size + OBJECT_HEADER_SIZE)

	var finalObj []byte = make([]byte, OBJECT_HEADER_SIZE+size)

	for c := OBJECT_GC_HEADER_SIZE; c < OBJECT_HEADER_SIZE; c++ {
		finalObj[c] = sizeB[c-OBJECT_GC_HEADER_SIZE]
	}
	for c := OBJECT_HEADER_SIZE; c < size+OBJECT_HEADER_SIZE; c++ {
		finalObj[c] = obj[c-OBJECT_HEADER_SIZE]
	}

	WriteMemory(heapOffset, finalObj)
	return heapOffset + OBJECT_HEADER_SIZE
}

func WriteObject(out1Offset int, obj []byte) {
	off := encoder.SerializeAtomic(int32(writeObj(obj)))

	WriteMemory(out1Offset, off)
}

func WriteObjectRetOff(obj []byte) int {
	return writeObj(obj)
}

func ErrorHeader(currentFile string, lineNo int) string {
	return "error: " + currentFile + ":" + strconv.FormatInt(int64(lineNo), 10)
}

func runtimeErrorInfo(r interface{}, printStack bool) {
	call := PROGRAM.CallStack[PROGRAM.CallCounter]
	expr := call.Operator.Expressions[call.Line]
	fmt.Println(ErrorHeader(expr.FileName, expr.FileLine), r)

	if printStack {
		PROGRAM.PrintStack()
	}

	if DBG_GOLANG_STACK_TRACE {
		debug.PrintStack()
	}

	os.Exit(3)
}

func RuntimeError() {
	if r := recover(); r != nil {
		switch r {
		case STACK_OVERFLOW_ERROR:
			call := PROGRAM.CallStack[PROGRAM.CallCounter]
			if PROGRAM.CallCounter > 0 {
				PROGRAM.CallCounter--
				PROGRAM.StackPointer = call.FramePointer
				runtimeErrorInfo(r, true)
			} else {
				// error at entry point
				runtimeErrorInfo(r, false)
			}
		default:
			runtimeErrorInfo(r, true)
		}
		os.Exit(CX_RUNTIME_ERROR)
	}
}

func getNonCollectionValue(fp int, arg, elt *CXArgument, typ string) string {
	switch typ {
	case "bool":
		return fmt.Sprintf("%v", ReadBool(fp, elt))
	case "byte":
		return fmt.Sprintf("%v", ReadByte(fp, elt))
	case "str":
		return fmt.Sprintf("%v", ReadStr(fp, elt))
	case "i32":
		return fmt.Sprintf("%v", ReadI32(fp, elt))
	case "i64":
		return fmt.Sprintf("%v", ReadI64(fp, elt))
	case "f32":
		return fmt.Sprintf("%v", ReadF32(fp, elt))
	case "f64":
		return fmt.Sprintf("%v", ReadF64(fp, elt))
	default:
		// then it's a struct
		var val string
		val = "{"
		// for _, fld := range elt.CustomType.Fields {
		lFlds := len(elt.CustomType.Fields)
		off := 0
		for c := 0; c < lFlds; c++ {
			fld := elt.CustomType.Fields[c]
			if c == lFlds-1 {
				val += fmt.Sprintf("%s: %s", fld.Name, GetPrintableValue(fp+arg.Offset+off, fld))
			} else {
				val += fmt.Sprintf("%s: %s, ", fld.Name, GetPrintableValue(fp+arg.Offset+off, fld))
			}
			off += fld.TotalSize
		}
		val += "}"
		return val
	}
}

func GetPrintableValue(fp int, arg *CXArgument) string {
	var typ string
	elt := GetAssignmentElement(arg)
	if elt.CustomType != nil {
		// then it's custom type
		typ = elt.CustomType.Name
	} else {
		// then it's native type
		typ = TypeNames[elt.Type]
	}

	if len(elt.Lengths) > 0 {
		var val string
		if len(elt.Lengths) == 1 {
			val = "["
			for c := 0; c < elt.Lengths[0]; c++ {
				if c == elt.Lengths[0]-1 {
					val += getNonCollectionValue(fp+c*elt.Size, arg, elt, typ)
				} else {
					val += getNonCollectionValue(fp+c*elt.Size, arg, elt, typ) + ", "
				}

			}
			val += "]"
		} else {
			// 5, 4, 1
			val = ""

			finalSize := 1
			for _, l := range elt.Lengths {
				finalSize *= l
			}

			lens := make([]int, len(elt.Lengths))
			copy(lens, elt.Lengths)

			for c := 0; c < len(lens); c++ {
				for i := 0; i < len(lens[c+1:]); i++ {
					lens[c] *= lens[c+i]
				}
			}

			for range lens {
				val += "["
			}

			// adding first element because of formatting reasons
			val += getNonCollectionValue(fp, arg, elt, typ)
			for c := 1; c < finalSize; c++ {
				closeCount := 0
				for _, l := range lens {
					if c%l == 0 && c != 0 {
						// val += "] ["
						closeCount++
					}
				}

				if closeCount > 0 {
					for c := 0; c < closeCount; c++ {
						val += "]"
					}
					val += " "
					for c := 0; c < closeCount; c++ {
						val += "["
					}

					val += getNonCollectionValue(fp+c*elt.Size, arg, elt, typ)
				} else {
					val += " " + getNonCollectionValue(fp+c*elt.Size, arg, elt, typ)
				}
			}
			for range lens {
				val += "]"
			}
		}

		return val
	}

	return getNonCollectionValue(fp, arg, elt, typ)
}
