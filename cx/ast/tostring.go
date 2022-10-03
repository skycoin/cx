package ast

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	cxpackages "github.com/skycoin/cx/cx/packages"

	"github.com/skycoin/cx/cx/constants"
	"github.com/skycoin/cx/cx/types"
	"github.com/skycoin/cx/cx/util"
)

// ToString returns the abstract syntax tree of a CX program in a
// string format.
func ToString(cxprogram *CXProgram) string {
	var ast3 string
	// ast3 += "Program\n" //why is top line "Program" ???

	BuildStrPackages(cxprogram, &ast3) //what does this do?

	return ast3
}

// buildStrImports is an auxiliary function for `toString`. It builds
// string representation all the imported packages of `pkg`.
func buildStrImports(prgrm *CXProgram, pkg *CXPackage, ast *string) {
	if len(pkg.Imports) > 0 {
		*ast += "\tImports\n"
	}

	count := 0
	for _, impIdx := range pkg.Imports {
		impPkg, err := prgrm.GetPackageFromArray(impIdx)
		if err != nil {
			panic(err)
		}
		*ast += fmt.Sprintf("\t\t%d.- Import: %s\n", count, impPkg.Name)
		count++
	}
}

// buildStrGlobals is an auxiliary function for `toString`. It builds
// string representation of all the global variables of `pkg`.
func buildStrGlobals(prgrm *CXProgram, pkg *CXPackage, ast *string) {
	if len(pkg.Globals.Fields) > 0 {
		*ast += "\tGlobals\n"
	}

	for idx, glblIdx := range pkg.Globals.Fields {
		glbl := prgrm.GetCXTypeSignatureFromArray(glblIdx)
		if glbl.Type == TYPE_CXARGUMENT_DEPRECATE {
			*ast += fmt.Sprintf("\t\t%d.- Global: %s %s\n", idx, prgrm.GetCXArg(CXArgumentIndex(glbl.Meta)).Name, GetFormattedType(prgrm, glbl))
		} else if glbl.Type == TYPE_ATOMIC {
			*ast += fmt.Sprintf("\t\t%d.- Global: %s %s\n", idx, glbl.Name, types.Code(glbl.Meta).Name())
		} else if glbl.Type == TYPE_POINTER_ATOMIC {
			*ast += fmt.Sprintf("\t\t%d.- Global: %s *%s\n", idx, glbl.Name, types.Code(glbl.Meta).Name())
		} else {
			panic("type is not known")
		}

	}
}

// buildStrStructs is an auxiliary function for `toString`. It builds
// string representation of all the structures defined in `pkg`.
func buildStrStructs(prgrm *CXProgram, pkg *CXPackage, ast *string) {
	if len(pkg.Structs) > 0 {
		*ast += "\tStructs\n"
	}

	count := 0
	for _, strctIdx := range pkg.Structs {
		strct := prgrm.CXStructs[strctIdx]
		*ast += fmt.Sprintf("\t\t%d.- Struct: %s\n", count, strct.Name)

		for k, typeSignatureIdx := range strct.Fields {
			typeSignature := prgrm.GetCXTypeSignatureFromArray(typeSignatureIdx)

			if typeSignature.Type == TYPE_CXARGUMENT_DEPRECATE {
				*ast += fmt.Sprintf("\t\t\t%d.- Field: %s %s\n",
					k, typeSignature.Name, GetFormattedType(prgrm, typeSignature))
			} else if typeSignature.Type == TYPE_ATOMIC {
				*ast += fmt.Sprintf("\t\t\t%d.- Field: %s %s\n",
					k, typeSignature.Name, types.Code(typeSignature.Meta).Name())
			} else if typeSignature.Type == TYPE_POINTER_ATOMIC {
				*ast += fmt.Sprintf("\t\t\t%d.- Field: %s *%s\n",
					k, typeSignature.Name, types.Code(typeSignature.Meta).Name())
			} else {
				panic("type is not known")
			}
		}

		count++
	}
}

// buildStrFunctions is an auxiliary function for `toString`. It builds
// string representation of all the functions defined in `pkg`.
func buildStrFunctions(prgrm *CXProgram, pkg *CXPackage, ast1 *string) {
	if len(pkg.Functions) > 0 {
		*ast1 += "\tFunctions\n"
	}

	// We need to declare the counter outside so we can
	// ignore the increment from the `*init` function.
	var j int
	for _, fnIdx := range pkg.Functions {
		fn := prgrm.GetFunctionFromArray(fnIdx)

		if fn.Name == constants.SYS_INIT_FUNC {
			continue
		}
		_, err := pkg.SelectFunction(prgrm, fn.Name)
		if err != nil {
			panic(err)
		}

		var inps bytes.Buffer
		var outs bytes.Buffer

		getFormattedParam(prgrm, fn.GetInputs(prgrm), pkg, &inps)
		getFormattedParam(prgrm, fn.GetOutputs(prgrm), pkg, &outs)

		*ast1 += fmt.Sprintf("\t\t%d.- Function: %s (%s) (%s)\n",
			j, fn.Name, inps.String(), outs.String())

		for k, expr := range fn.Expressions {
			var inps bytes.Buffer
			var outs bytes.Buffer
			var opName1 string
			var lbl string

			cxAtomicOp, err := prgrm.GetCXAtomicOp(expr.Index)
			if err != nil {
				panic(err)
			}

			cxAtomicOpOperator := prgrm.GetFunctionFromArray(cxAtomicOp.Operator)
			// Adding label in case a `goto` statement was used for the expression.
			if cxAtomicOp.Label != "" {
				lbl = " <<" + cxAtomicOp.Label + ">>"
			} else {
				lbl = ""
			}

			// Determining operator's name.
			if cxAtomicOpOperator != nil {
				if cxAtomicOpOperator.IsBuiltIn() {

					opName1 = OpNames[cxAtomicOpOperator.AtomicOPCode]
				} else {
					opName1 = cxAtomicOpOperator.Name
				}
			}

			getFormattedParam(prgrm, cxAtomicOp.GetInputs(prgrm), pkg, &inps)
			getFormattedParam(prgrm, cxAtomicOp.GetOutputs(prgrm), pkg, &outs)

			if expr.Type == CX_LINE {
				cxLine, _ := prgrm.GetCXLine(expr.Index)
				*ast1 += fmt.Sprintf("\t\t\t%d.- Line: %v: %s\n",
					k,
					cxLine.LineNumber,
					strings.TrimSpace(cxLine.LineStr))
			} else if cxAtomicOpOperator != nil {
				assignOp := ""
				if outs.Len() > 0 {
					assignOp = " = "
				}
				*ast1 += fmt.Sprintf("\t\t\t%d.- Expression%s: %s%s%s(%s)\n",
					k,
					lbl,
					outs.String(),
					assignOp,
					opName1,
					inps.String(),
				)
			} else {
				// Then it's a variable declaration. These are represented
				// by expressions without operators that only have outputs.
				cxAtomicOpOutputs := cxAtomicOp.GetOutputs(prgrm)
				cxAtomicOpOutputTypeSignature := prgrm.GetCXTypeSignatureFromArray(cxAtomicOpOutputs[len(cxAtomicOpOutputs)-1])
				if len(cxAtomicOpOutputs) > 0 {
					if cxAtomicOpOutputTypeSignature.Type == TYPE_CXARGUMENT_DEPRECATE {
						*ast1 += fmt.Sprintf("\t\t\t%d.- Declaration%s: %s %s\n",
							k,
							lbl,
							cxAtomicOpOutputTypeSignature.Name,
							GetFormattedType(prgrm, cxAtomicOpOutputTypeSignature))
					} else if cxAtomicOpOutputTypeSignature.Type == TYPE_ATOMIC {
						*ast1 += fmt.Sprintf("\t\t\t%d.- Declaration%s: %s %s\n",
							k,
							lbl,
							cxAtomicOpOutputTypeSignature.Name,
							types.Code(cxAtomicOpOutputTypeSignature.Meta).Name())
					} else if cxAtomicOpOutputTypeSignature.Type == TYPE_POINTER_ATOMIC {
						*ast1 += fmt.Sprintf("\t\t\t%d.- Declaration%s: %s *%s\n",
							k,
							lbl,
							cxAtomicOpOutputTypeSignature.Name,
							types.Code(cxAtomicOpOutputTypeSignature.Meta).Name())
					} else if cxAtomicOpOutputTypeSignature.Type == TYPE_SLICE_ATOMIC {
						*ast1 += fmt.Sprintf("\t\t\t%d.- Declaration%s: %s *%s\n",
							k,
							lbl,
							cxAtomicOpOutputTypeSignature.Name,
							types.Code(types.SLICE).Name())
					} else {
						panic("type is not known")
					}
				}
			}
		}

		j++
	}
}

// BuildStrPackages is an auxiliary function for `ToString`. It starts the
// process of building string format of the abstract syntax tree of a CX program.
func BuildStrPackages(prgrm *CXProgram, ast *string) {
	// We need to declare the counter outside so we can
	// ignore the increments from core or stdlib packages.
	var i int
	for _, pkgIdx := range prgrm.Packages {
		pkg, err := prgrm.GetPackageFromArray(pkgIdx)
		if err != nil {
			panic(err)
		}

		if cxpackages.IsDefaultPackage(pkg.Name) {
			continue
		}

		*ast += fmt.Sprintf("%d.- Package: %s\n", i, pkg.Name)

		buildStrImports(prgrm, pkg, ast)
		buildStrGlobals(prgrm, pkg, ast)
		buildStrStructs(prgrm, pkg, ast)
		buildStrFunctions(prgrm, pkg, ast)

		i++
	}
}

// getFormattedParam is an auxiliary function for `ToString`. It formats the
// name of a `CXExpression`'s input and output parameters (`CXArgument`s). Examples
// of these formattings are "pkg.foo[0]", "&*foo.field1". The result is written to
// `buf`.
func getFormattedParam(prgrm *CXProgram, paramTypeSigIdxs []CXTypeSignatureIndex, pkg *CXPackage, buf *bytes.Buffer) {
	for i, paramTypeSigIdx := range paramTypeSigIdxs {
		paramTypeSig := prgrm.GetCXTypeSignatureFromArray(paramTypeSigIdx)

		// Checking if this argument comes from an imported package.
		externalPkg := false
		if CXPackageIndex(pkg.Index) != paramTypeSig.Package {
			externalPkg = true
		}
		if paramTypeSig.Type == TYPE_CXARGUMENT_DEPRECATE {
			param := prgrm.GetCXArgFromArray(CXArgumentIndex(paramTypeSig.Meta))

			buf.WriteString(fmt.Sprintf("%s %s", GetFormattedName(prgrm, param, externalPkg, pkg), GetFormattedType(prgrm, paramTypeSig)))
		} else if paramTypeSig.Type == TYPE_ATOMIC {
			name := paramTypeSig.Name

			// If it's a literal, just override the name with LITERAL_PLACEHOLDER.
			if paramTypeSig.Name == "" {
				name = constants.LITERAL_PLACEHOLDER
			}

			// TODO: Check if external pkg and pkg name shown are correct
			if externalPkg {
				name = fmt.Sprintf("%s.%s", pkg.Name, name)
			}

			buf.WriteString(fmt.Sprintf("%s %s", name, types.Code(paramTypeSig.Meta).Name()))
		} else if paramTypeSig.Type == TYPE_POINTER_ATOMIC {

			name := paramTypeSig.Name

			// If it's a literal, just override the name with LITERAL_PLACEHOLDER.
			if paramTypeSig.Name == "" {
				name = constants.LITERAL_PLACEHOLDER
			}

			// TODO: Check if external pkg and pkg name shown are correct
			if externalPkg {
				name = fmt.Sprintf("%s.%s", pkg.Name, name)
			}

			buf.WriteString(fmt.Sprintf("%s *%s", name, types.Code(paramTypeSig.Meta).Name()))
		} else if paramTypeSig.Type == TYPE_SLICE_ATOMIC {

			name := paramTypeSig.Name

			// If it's a literal, just override the name with LITERAL_PLACEHOLDER.
			if paramTypeSig.Name == "" {
				name = constants.LITERAL_PLACEHOLDER
			}

			// TODO: Check if external pkg and pkg name shown are correct
			if externalPkg {
				name = fmt.Sprintf("%s.%s", pkg.Name, name)
			}

			buf.WriteString(fmt.Sprintf("%s *%s", name, types.Code(types.SLICE).Name()))
		} else {
			panic("type is not known")
		}

		if i != len(paramTypeSigIdxs)-1 {
			buf.WriteString(", ")
		}

	}
}

// SignatureStringOfFunction returns the signature string of a function.
func SignatureStringOfFunction(prgrm *CXProgram, pkg *CXPackage, f *CXFunction) string {
	var ins bytes.Buffer
	var outs bytes.Buffer

	getFormattedParam(prgrm, f.GetInputs(prgrm), pkg, &ins)
	getFormattedParam(prgrm, f.GetOutputs(prgrm), pkg, &outs)

	return fmt.Sprintf("func %s(%s) (%s)",
		f.Name, ins.String(), outs.String())
}

func getNonCollectionValue(prgrm *CXProgram, fp types.Pointer, arg, elt *CXArgument, typ string) string {
	if arg.IsPointer() {
		return fmt.Sprintf("%v", types.Read_ptr(prgrm.Memory, GetFinalOffset(prgrm, fp, elt, nil)))
	}
	if arg.IsSlice {
		return fmt.Sprintf("%v", types.GetSlice_byte(prgrm.Memory, GetFinalOffset(prgrm, fp, elt, nil), GetArgSize(prgrm, elt)))
	}
	switch typ {
	case "bool":
		return fmt.Sprintf("%v", types.Read_bool(prgrm.Memory, GetFinalOffset(prgrm, fp, elt, nil)))
	case "str":
		return fmt.Sprintf("%v", types.Read_str(prgrm.Memory, GetFinalOffset(prgrm, fp, elt, nil)))
	case "i8":
		return fmt.Sprintf("%v", types.Read_i8(prgrm.Memory, GetFinalOffset(prgrm, fp, elt, nil)))
	case "i16":
		return fmt.Sprintf("%v", types.Read_i16(prgrm.Memory, GetFinalOffset(prgrm, fp, elt, nil)))
	case "i32":
		return fmt.Sprintf("%v", types.Read_i32(prgrm.Memory, GetFinalOffset(prgrm, fp, elt, nil)))
	case "i64":
		return fmt.Sprintf("%v", types.Read_i64(prgrm.Memory, GetFinalOffset(prgrm, fp, elt, nil)))
	case "ui8":
		return fmt.Sprintf("%v", types.Read_ui8(prgrm.Memory, GetFinalOffset(prgrm, fp, elt, nil)))
	case "ui16":
		return fmt.Sprintf("%v", types.Read_ui16(prgrm.Memory, GetFinalOffset(prgrm, fp, elt, nil)))
	case "ui32":
		return fmt.Sprintf("%v", types.Read_ui32(prgrm.Memory, GetFinalOffset(prgrm, fp, elt, nil)))
	case "ui64":
		return fmt.Sprintf("%v", types.Read_ui64(prgrm.Memory, GetFinalOffset(prgrm, fp, elt, nil)))
	case "f32":
		return fmt.Sprintf("%v", types.Read_f32(prgrm.Memory, GetFinalOffset(prgrm, fp, elt, nil)))
	case "f64":
		return fmt.Sprintf("%v", types.Read_f64(prgrm.Memory, GetFinalOffset(prgrm, fp, elt, nil)))
	default:
		// then it's a struct
		var val string
		val = "{"
		// for _, fld := range elt.StructType.Fields {
		lFlds := len(elt.StructType.Fields)
		off := types.Pointer(0)
		for c := 0; c < lFlds; c++ {
			typeSignatureIdx := elt.StructType.Fields[c]
			typeSignature := prgrm.GetCXTypeSignatureFromArray(typeSignatureIdx)
			var fldTotalSize types.Pointer
			if typeSignature.Type == TYPE_CXARGUMENT_DEPRECATE {
				fldIdx := typeSignature.Meta
				fldTotalSize = GetArgSize(prgrm, &prgrm.CXArgs[fldIdx])
			} else if typeSignature.Type == TYPE_ATOMIC || typeSignature.Type == TYPE_POINTER_ATOMIC {
				fldTotalSize = types.Code(typeSignature.Type).Size()
			} else if typeSignature.Type == TYPE_SLICE_ATOMIC {
				fldTotalSize = typeSignature.GetSize(prgrm, false)
			} else {
				panic("type is not known")
			}

			if c == lFlds-1 {
				val += fmt.Sprintf("%s: %s", typeSignature.Name, GetPrintableValue(prgrm, fp+arg.Offset+off, typeSignature))
			} else {
				val += fmt.Sprintf("%s: %s, ", typeSignature.Name, GetPrintableValue(prgrm, fp+arg.Offset+off, typeSignature))
			}
			off += fldTotalSize
		}
		val += "}"
		return val
	}
}

func readValue(prgrm *CXProgram, typ string, offset types.Pointer) string {
	switch typ {
	case "bool":
		return fmt.Sprintf("%v", types.Read_bool(prgrm.Memory, offset))
	case "str":
		return fmt.Sprintf("%v", types.Read_str(prgrm.Memory, offset))
	case "i8":
		return fmt.Sprintf("%v", types.Read_i8(prgrm.Memory, offset))
	case "i16":
		return fmt.Sprintf("%v", types.Read_i16(prgrm.Memory, offset))
	case "i32":
		return fmt.Sprintf("%v", types.Read_i32(prgrm.Memory, offset))
	case "i64":
		return fmt.Sprintf("%v", types.Read_i64(prgrm.Memory, offset))
	case "ui8":
		return fmt.Sprintf("%v", types.Read_ui8(prgrm.Memory, offset))
	case "ui16":
		return fmt.Sprintf("%v", types.Read_ui16(prgrm.Memory, offset))
	case "ui32":
		return fmt.Sprintf("%v", types.Read_ui32(prgrm.Memory, offset))
	case "ui64":
		return fmt.Sprintf("%v", types.Read_ui64(prgrm.Memory, offset))
	case "f32":
		return fmt.Sprintf("%v", types.Read_f32(prgrm.Memory, offset))
	case "f64":
		return fmt.Sprintf("%v", types.Read_f64(prgrm.Memory, offset))
	}

	return ""
}

// ReadSliceElements ...
func ReadSliceElements(prgrm *CXProgram, fp types.Pointer, arg, elt *CXArgument, sliceData []byte, size types.Pointer, typ string) string {
	switch typ {
	case "bool":
		return fmt.Sprintf("%v", types.Read_bool(sliceData, 0))
	case "str":
		return fmt.Sprintf("%v", types.Read_str(prgrm.Memory, types.Read_ptr(sliceData, 0)))
	case "i8":
		return fmt.Sprintf("%v", types.Read_i8(sliceData, 0))
	case "i16":
		return fmt.Sprintf("%v", types.Read_i16(sliceData, 0))
	case "i32":
		return fmt.Sprintf("%v", types.Read_i32(sliceData, 0))
	case "i64":
		return fmt.Sprintf("%v", types.Read_i64(sliceData, 0))
	case "ui8":
		return fmt.Sprintf("%v", types.Read_ui8(sliceData, 0))
	case "ui16":
		return fmt.Sprintf("%v", types.Read_ui16(sliceData, 0))
	case "ui32":
		return fmt.Sprintf("%v", types.Read_ui32(sliceData, 0))
	case "ui64":
		return fmt.Sprintf("%v", types.Read_ui64(sliceData, 0))
	case "f32":
		return fmt.Sprintf("%v", types.Read_f32(sliceData, 0))
	case "f64":
		return fmt.Sprintf("%v", types.Read_f64(sliceData, 0))
	default:
		// then it's a struct
		var val string
		val = "{"
		// for _, fld := range elt.StructType.Fields {
		lFlds := len(elt.StructType.Fields)
		off := types.Pointer(0)
		for c := 0; c < lFlds; c++ {
			typeSignatureIdx := elt.StructType.Fields[c]
			typeSignature := prgrm.GetCXTypeSignatureFromArray(typeSignatureIdx)

			var fldTotalSize types.Pointer
			if typeSignature.Type == TYPE_CXARGUMENT_DEPRECATE {
				fldIdx := typeSignature.Meta
				fldTotalSize = GetArgSize(prgrm, &prgrm.CXArgs[fldIdx])
			} else if typeSignature.Type == TYPE_ATOMIC || typeSignature.Type == TYPE_POINTER_ATOMIC {
				fldTotalSize = types.Code(typeSignature.Type).Size()
			} else if typeSignature.Type == TYPE_SLICE_ATOMIC {
				fldTotalSize = typeSignature.GetSize(prgrm, false)
			} else {
				panic("type is not known")
			}

			if c == lFlds-1 {
				val += fmt.Sprintf("%s: %s", typeSignature.Name, GetPrintableValue(prgrm, fp+arg.Offset+off, typeSignature))
			} else {
				val += fmt.Sprintf("%s: %s, ", typeSignature.Name, GetPrintableValue(prgrm, fp+arg.Offset+off, typeSignature))
			}
			off += fldTotalSize
		}
		val += "}"
		return val
	}
}

func arrayPrinter(c types.Pointer, arrayLengths []types.Pointer, closeStr, openStr string) string {
	val := types.Pointer(1)

	for _, valLen := range arrayLengths {
		val *= valLen
	}

	if c%val == 0 {
		closeStr += "]"
		openStr += "["
	}

	if len(arrayLengths) > 1 {
		return arrayPrinter(c, arrayLengths[1:], closeStr, openStr)
	}

	return closeStr + " " + openStr
}

// GetPrintableValue ...
func GetPrintableValue(prgrm *CXProgram, fp types.Pointer, argTypeSig *CXTypeSignature) string {
	var arg, elt *CXArgument
	var typ string
	if argTypeSig.Type == TYPE_CXARGUMENT_DEPRECATE {
		arg = prgrm.GetCXArgFromArray(CXArgumentIndex(argTypeSig.Meta))

		elt = arg.GetAssignmentElement(prgrm)
		if elt.StructType != nil {
			// then it's struct type
			typ = elt.StructType.Name
		} else {
			// then it's native type
			typ = elt.Type.Name()
		}

		if len(elt.Lengths) > 0 {
			var val string
			if len(elt.Lengths) == 1 {
				val = "["

				if arg.IsSlice {
					// for slices
					sliceOffset := GetSliceOffset(prgrm, fp, argTypeSig)

					sliceData := GetSlice(prgrm, sliceOffset, elt.Size)
					if len(sliceData) != 0 {
						sliceLen := types.Read_ptr(sliceData, 0)
						for c := types.Pointer(0); c < sliceLen; c++ {
							if c == sliceLen-1 {
								val += ReadSliceElements(prgrm, sliceOffset+constants.SLICE_HEADER_SIZE+types.OBJECT_HEADER_SIZE+c*elt.Size, arg, elt, sliceData[types.POINTER_SIZE+c*elt.Size:], elt.Size, typ)
							} else {
								val += ReadSliceElements(prgrm, sliceOffset+constants.SLICE_HEADER_SIZE+types.OBJECT_HEADER_SIZE+c*elt.Size, arg, elt, sliceData[types.POINTER_SIZE+c*elt.Size:], elt.Size, typ) + ", "
							}

						}
					}

				} else {
					// for Arrays
					for c := types.Pointer(0); c < elt.Lengths[0]; c++ {
						val += getNonCollectionValue(prgrm, fp+c*elt.Size, arg, elt, typ)
						if c != elt.Lengths[0]-1 {
							val += ", "
						}
					}
				}

				val += "]"
			} else {
				// 5, 4, 1
				val = ""

				finalSize := types.Pointer(1)
				for _, l := range elt.Lengths {
					finalSize *= l
					val += "["
				}

				// adding first element because of formatting reasons
				val += getNonCollectionValue(prgrm, fp, arg, elt, typ)

				for c := types.Pointer(1); c < finalSize; c++ {
					val += arrayPrinter(c, elt.Lengths, "", "")

					val += getNonCollectionValue(prgrm, fp+c*elt.Size, arg, elt, typ)

				}

				for range elt.Lengths {
					val += "]"
				}
			}

			return val
		}
	} else if argTypeSig.Type == TYPE_ATOMIC {
		// TODO: improve this

		if argTypeSig.PassBy == constants.PASSBY_REFERENCE {
			return fmt.Sprintf("%v", GetFinalOffset(prgrm, fp, nil, argTypeSig))
		}

		typ = types.Code(argTypeSig.Meta).Name()

		return readValue(prgrm, typ, GetFinalOffset(prgrm, fp, nil, argTypeSig))
	} else if argTypeSig.Type == TYPE_POINTER_ATOMIC {
		// TODO: improve this
		return fmt.Sprintf("%v", types.Read_ptr(prgrm.Memory, GetFinalOffset(prgrm, fp, nil, argTypeSig)))
	} else if argTypeSig.Type == TYPE_ARRAY_ATOMIC {
		arrDetails := prgrm.GetCXTypeSignatureArrayFromArray(argTypeSig.Meta)
		var val string

		typ := types.Code(arrDetails.Type).Name()
		if len(arrDetails.Lengths) == 1 {
			val += "["

			// for Arrays
			for c := types.Pointer(0); c < arrDetails.Lengths[0]; c++ {
				val += readValue(prgrm, typ, GetFinalOffset(prgrm, fp+(c*types.Code(arrDetails.Type).Size()), nil, argTypeSig))
				if c != arrDetails.Lengths[0]-1 {
					val += ", "
				}
			}

			val += "]"

			return val
		} else {
			val = ""

			finalSize := types.Pointer(1)
			for _, l := range arrDetails.Lengths {
				finalSize *= l
				val += "["
			}

			// adding first element because of formatting reasons
			val += readValue(prgrm, typ, GetFinalOffset(prgrm, fp, nil, argTypeSig))

			for c := types.Pointer(1); c < finalSize; c++ {
				val += arrayPrinter(c, arrDetails.Lengths, "", "")
				val += readValue(prgrm, typ, GetFinalOffset(prgrm, fp+(c*types.Code(arrDetails.Type).Size()), nil, argTypeSig))
			}

			for range arrDetails.Lengths {
				val += "]"
			}

			return val
		}

	} else if argTypeSig.Type == TYPE_POINTER_ARRAY_ATOMIC {
		arrDetails := prgrm.GetCXTypeSignatureArrayFromArray(argTypeSig.Meta)
		var val string

		typ := types.Code(arrDetails.Type).Name()
		if len(arrDetails.Lengths) == 1 {
			val += "["

			// for Arrays
			for c := types.Pointer(0); c < arrDetails.Lengths[0]; c++ {
				val += readValue(prgrm, typ, GetFinalOffset(prgrm, fp+(c*types.Code(arrDetails.Type).Size()), nil, argTypeSig))

				if c != arrDetails.Lengths[0]-1 {
					val += ", "
				}
			}

			val += "]"

			return val
		} else {
			val = ""

			finalSize := types.Pointer(1)
			for _, l := range arrDetails.Lengths {
				finalSize *= l
				val += "["
			}

			// adding first element because of formatting reasons
			val += readValue(prgrm, typ, GetFinalOffset(prgrm, fp, nil, argTypeSig))

			for c := types.Pointer(1); c < finalSize; c++ {
				val += arrayPrinter(c, arrDetails.Lengths, "", "")
				val += readValue(prgrm, typ, GetFinalOffset(prgrm, fp+(c*types.Code(arrDetails.Type).Size()), nil, argTypeSig))
			}

			for range arrDetails.Lengths {
				val += "]"
			}

			return val
		}
	} else if argTypeSig.Type == TYPE_SLICE_ATOMIC {
		arrDetails := prgrm.GetCXTypeSignatureArrayFromArray(argTypeSig.Meta)
		var val string

		typ := types.Code(arrDetails.Type).Name()
		if len(arrDetails.Lengths) == 1 {
			val += "["

			// for slices
			sizeOfElement := types.Code(arrDetails.Type).Size()
			sliceOffset := GetSliceOffset(prgrm, fp, argTypeSig)
			sliceData := GetSlice(prgrm, sliceOffset, sizeOfElement)

			if len(sliceData) != 0 {
				sliceLen := types.Read_ptr(sliceData, 0)
				for c := types.Pointer(0); c < sliceLen; c++ {
					if c == sliceLen-1 {
						val += ReadSliceElements(prgrm, sliceOffset+constants.SLICE_HEADER_SIZE+types.OBJECT_HEADER_SIZE+c*sizeOfElement, arg, elt, sliceData[types.POINTER_SIZE+c*sizeOfElement:], sizeOfElement, typ)
					} else {
						val += ReadSliceElements(prgrm, sliceOffset+constants.SLICE_HEADER_SIZE+types.OBJECT_HEADER_SIZE+c*sizeOfElement, arg, elt, sliceData[types.POINTER_SIZE+c*sizeOfElement:], sizeOfElement, typ) + ", "
					}

				}
			}

			val += "]"

			return val
		} else {
			val = ""

			finalSize := types.Pointer(1)
			for _, l := range arrDetails.Lengths {
				finalSize *= l
				val += "["
			}

			// adding first element because of formatting reasons
			val += readValue(prgrm, typ, GetFinalOffset(prgrm, fp, nil, argTypeSig))

			for c := types.Pointer(1); c < finalSize; c++ {
				val += arrayPrinter(c, arrDetails.Lengths, "", "")
				val += readValue(prgrm, typ, GetFinalOffset(prgrm, fp+(c*types.Code(arrDetails.Type).Size()), nil, argTypeSig))
			}

			for range arrDetails.Lengths {
				val += "]"
			}

			return val
		}

	} else {
		panic("type is not known")
	}

	return getNonCollectionValue(prgrm, fp, arg, elt, typ)
}

// filePathWalkDir scans all the files in a directory. It will automatically
// scan each sub-directories in the directory. Code obtained from manigandand's
// post in https://stackoverflow.com/questions/14668850/list-directory-in-go.
func filePathWalkDir(root string) ([]string, error) {
	var files []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			files = append(files, path)
		}
		return err
	})
	return files, err
}

// ioReadDir reads the directory named by dirname and returns a list of
// directory entries sorted by filename. Code obtained from manigandand's
// post in https://stackoverflow.com/questions/14668850/list-directory-in-go.
func ioReadDir(root string) ([]string, error) {
	var files []string
	fileInfo, err := ioutil.ReadDir(root)
	if err != nil {
		return files, err
	}

	for _, file := range fileInfo {
		files = append(files, fmt.Sprintf("%s/%s", root, file.Name()))
	}
	return files, nil
}

// ParseArgsForCX parses the arguments and returns:
//  - []arguments
//  - []file pointers	open files
//  - []sting		filenames
func ParseArgsForCX(args []string, alsoSubdirs bool) (cxArgs []string, sourceCode []*os.File, fileNames []string) {
	skip := false // flag for skipping arg

	for _, arg := range args {

		// skip arg if skip flag is specified
		if skip {
			skip = false
			continue
		}

		// cli flags are either "--key=value" or "-key value"
		// we have to skip both cases
		if len(arg) > 1 && arg[0] == '-' {
			if !strings.Contains(arg, "=") {
				skip = true
			}
			continue
		}

		// cli cx flags are prefixed with "++"
		if len(arg) > 2 && arg[:2] == "++" {
			cxArgs = append(cxArgs, arg)
			continue
		}

		fi, err := util.CXStatFile(arg)
		if err != nil {
			println(fmt.Sprintf("%s: source file or library not found", arg))
			os.Exit(constants.CX_COMPILATION_ERROR)
		}

		switch mode := fi.Mode(); {
		case mode.IsDir():
			var fileList []string
			var err error

			// Checking if we want to check all subdirectories.
			if alsoSubdirs {
				fileList, err = filePathWalkDir(arg)
			} else {
				fileList, err = ioReadDir(arg)
				// fileList, err = filePathWalkDir(arg)
			}

			if err != nil {
				panic(err)
			}

			for _, path := range fileList {
				file, err := util.CXOpenFile(path)

				if err != nil {
					println(fmt.Sprintf("%s: source file or library not found", arg))
					os.Exit(constants.CX_COMPILATION_ERROR)
				}

				fiName := file.Name()
				fiNameLen := len(fiName)

				if fiNameLen > 2 && fiName[fiNameLen-3:] == ".cx" {
					// only loading .cx files
					sourceCode = append(sourceCode, file)
					fileNames = append(fileNames, fiName)
				}
			}
		case mode.IsRegular():
			file, err := util.CXOpenFile(arg)

			if err != nil {
				panic(err)
			}

			fileNames = append(fileNames, file.Name())
			sourceCode = append(sourceCode, file)
		}
	}

	return cxArgs, sourceCode, fileNames
}

// getFormattedDerefs is an auxiliary function for `GetFormattedName`. This
// function formats indexing and pointer dereferences associated to `arg`.
func getFormattedDerefs(prgrm *CXProgram, arg *CXArgument, includePkg bool, pkg *CXPackage) string {
	name := ""

	argPkg := pkg
	// Checking if we should include `arg`'s package name.
	if includePkg {
		name = fmt.Sprintf("%s.%s", argPkg.Name, arg.Name)
	} else {
		name = arg.Name
	}

	// If it's a literal, just override the name with LITERAL_PLACEHOLDER.
	if arg.Name == "" {
		name = constants.LITERAL_PLACEHOLDER
	}

	// Checking if we have indexing operations, e.g. foo[2][1]
	for _, idxIdx := range arg.Indexes {
		idx := prgrm.GetCXTypeSignatureFromArray(idxIdx)
		// Checking if the value is in data segment.
		// If this is the case, we can safely display it.
		idxValue := ""
		if idx.Offset > prgrm.Stack.Size {
			// Then it's a literal.
			idxI32 := types.Read_ptr(prgrm.Memory, idx.Offset)
			idxValue = fmt.Sprintf("%d", idxI32)
		} else {
			// Then let's just print the variable name.
			idxValue = idx.Name
		}

		name = fmt.Sprintf("%s[%s]", name, idxValue)
	}

	return name
}

// GetFormattedName reads `arg.DereferenceOperations` and builds a string that
// depicts how an argument is being accessed. Example outputs: "foo[3]",
// "**bar", "foo.bar[0]". If `includePkg` is `true`, the argument name will
// include the package name that contains it, such as in "pkg.foo".
func GetFormattedName(prgrm *CXProgram, arg *CXArgument, includePkg bool, pkg *CXPackage) string {
	// Getting formatted name which does not include fields.
	name := getFormattedDerefs(prgrm, arg, includePkg, pkg)

	// Adding as suffixes all the fields.
	for _, fldIdx := range arg.Fields {
		fld := prgrm.GetCXArgFromArray(fldIdx)
		name = fmt.Sprintf("%s.%s", name, getFormattedDerefs(prgrm, fld, includePkg, pkg))
	}

	// Checking if we're referencing `arg`.
	if arg.PassBy == constants.PASSBY_REFERENCE {
		name = "&" + name
	}

	return name
}

// formatParameters returns a string containing a list of the formatted types of
// each of `params`, enclosed in parethesis. This function is used only when
// formatting functions as first-class objects.
func formatParameters(prgrm *CXProgram, params []CXTypeSignatureIndex) string {
	types := "("
	for i, paramIdx := range params {
		param := prgrm.GetCXTypeSignatureFromArray(paramIdx)
		types += GetFormattedType(prgrm, param)
		if i != len(params)-1 {
			types += ", "
		}
	}
	types += ")"

	return types
}

// GetFormattedType builds a string with the CXGO type representation of `arg`.
func GetFormattedType(prgrm *CXProgram, typeSig *CXTypeSignature) string {
	typ := ""

	if typeSig.Type == TYPE_CXARGUMENT_DEPRECATE {
		expressionOutputArg := prgrm.GetCXArgFromArray(CXArgumentIndex(typeSig.Meta))
		typ = getFormattedType_CXArg(prgrm, expressionOutputArg)
	} else if typeSig.Type == TYPE_ATOMIC {
		typ = types.Code(typeSig.Meta).Name()

		if typeSig.PassBy == constants.PASSBY_REFERENCE {
			typ = "*" + typ
		}
	} else if typeSig.Type == TYPE_POINTER_ATOMIC {
		typ = types.Code(typeSig.Meta).Name()
		if !typeSig.IsDeref {
			typ = "*" + typ
		}
	} else if typeSig.Type == TYPE_ARRAY_ATOMIC {
		arrayData := prgrm.GetCXTypeSignatureArrayFromArray(typeSig.Meta)

		arrLen := len(arrayData.Lengths) - len(arrayData.Indexes)
		if arrLen != 0 {
			for _, len := range arrayData.Lengths[len(arrayData.Indexes):] {
				typ = fmt.Sprintf("[%d]%s", len, typ)
			}
		}

		typ += types.Code(arrayData.Type).Name()
		if typeSig.PassBy == constants.PASSBY_REFERENCE {
			typ = "*" + typ
		}
	} else if typeSig.Type == TYPE_POINTER_ARRAY_ATOMIC {
		arrayData := prgrm.GetCXTypeSignatureArrayFromArray(typeSig.Meta)

		arrLen := len(arrayData.Lengths) - len(arrayData.Indexes)
		if arrLen != 0 {
			for _, len := range arrayData.Lengths[len(arrayData.Indexes):] {
				typ = fmt.Sprintf("[%d]%s", len, typ)
			}
		}

		typ += types.Code(arrayData.Type).Name()
		if !typeSig.IsDeref {
			typ = "*" + typ
		}

	} else if typeSig.Type == TYPE_SLICE_ATOMIC {
		arrayData := prgrm.GetCXTypeSignatureArrayFromArray(typeSig.Meta)

		arrLen := len(arrayData.Lengths) - len(arrayData.Indexes)
		if arrLen != 0 {
			for i := 0; i < len(arrayData.Lengths[len(arrayData.Indexes):]); i++ {
				typ = fmt.Sprintf("[]%s", typ)
			}
		}

		typ += types.Code(arrayData.Type).Name()
		if typeSig.PassBy == constants.PASSBY_REFERENCE {
			typ = "*" + typ
		}
	} else {
		panic("type is not known")
	}

	return typ
}

func getFormattedType_CXArg(prgrm *CXProgram, arg *CXArgument) string {
	typ := ""
	elt := arg.GetAssignmentElement(prgrm)

	// this is used to know what arg.Lengths index to use
	// used for cases like [5]*[3]i32, where we jump to another decl spec
	arrDeclCount := len(arg.Lengths) - 1
	// looping declaration specifiers
	for _, spec := range elt.DeclarationSpecifiers {
		switch spec {
		case constants.DECL_POINTER:
			typ = "*" + typ
		case constants.DECL_DEREF:
			typ = typ[1:]
		case constants.DECL_ARRAY:
			typ = fmt.Sprintf("[%d]%s", arg.Lengths[arrDeclCount], typ)
			arrDeclCount--
		case constants.DECL_SLICE:
			typ = "[]" + typ
		case constants.DECL_INDEXING:
		default:
			// base type
			if elt.StructType != nil {
				// then it's struct type
				typ += elt.StructType.Name
			} else {
				if elt.Type == types.POINTER {
					typ += elt.PointerTargetType.Name()
				} else {
					// then it's basic type
					typ += elt.Type.Name()
				}

				// If it's a function, let's add the inputs and outputs.
				if elt.Type == types.FUNC {
					// if elt.IsLocalDeclaration {
					// Then it's a local variable, which can be assigned to a
					// lambda function, for example.
					// typ += formatParameters(prgrm, prgrm.ConvertIndexArgsToPointerArgs(elt.Inputs))
					// typ += formatParameters(prgrm, prgrm.ConvertIndexArgsToPointerArgs(elt.Outputs))
					// } else {
					// Then it refers to a named function defined in a package.
					pkg, err := prgrm.GetPackageFromArray(arg.Package)
					if err != nil {
						println(CompilationError(elt.ArgDetails.FileName, elt.ArgDetails.FileLine), err.Error())
						os.Exit(constants.CX_COMPILATION_ERROR)
					}

					fn, err := pkg.GetFunction(prgrm, elt.Name)
					if err == nil {
						// println(CompilationError(elt.FileName, elt.FileLine), err.ProgramError())
						// os.Exit(CX_COMPILATION_ERROR)
						// Adding list of inputs and outputs types.
						typ += formatParameters(prgrm, fn.GetInputs(prgrm))
						typ += formatParameters(prgrm, fn.GetOutputs(prgrm))
					}
					// }
				}
			}
		}
	}
	return typ

}

// SignatureStringOfStruct returns the signature string of a struct.
func SignatureStringOfStruct(prgrm *CXProgram, s *CXStruct) string {
	fields := ""
	for _, typeSignatureIdx := range s.Fields {
		typeSignature := prgrm.GetCXTypeSignatureFromArray(typeSignatureIdx)

		if typeSignature.Type == TYPE_ATOMIC {
			fields += fmt.Sprintf(" %s %s;", typeSignature.Name, types.Code(typeSignature.Meta).Name())
			continue
		} else if typeSignature.Type == TYPE_POINTER_ATOMIC {
			fields += fmt.Sprintf(" %s *%s;", typeSignature.Name, types.Code(typeSignature.Meta).Name())
			continue
		} else {
			panic("type is not known")
		}

		fldIdx := typeSignature.Meta
		fld := prgrm.CXArgs[fldIdx]
		fields += fmt.Sprintf(" %s %s;", fld.Name, GetFormattedType(prgrm, typeSignature))
	}

	return fmt.Sprintf("%s struct {%s }", s.Name, fields)
}

// GetArgSizeFromTypeName ...
func GetArgSizeFromTypeName(typeName string) int {
	switch typeName {
	case "bool", "i8", "ui8":
		return 1
	case "i16", "ui16":
		return 2
	case "str", "i32", "ui32", "f32", "aff":
		return 4
	case "i64", "ui64", "f64":
		return 8
	default:
		return 4
		// return -1
		// panic(CX_INTERNAL_ERROR)
	}
}
