package cxcore

import (
	"bytes"
	"fmt"
	"github.com/skycoin/cx/cx/ast"
	"github.com/skycoin/cx/cx/constants"
	"github.com/skycoin/cx/cx/helper"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/tabwriter"
)

// Debug ...
func Debug(args ...interface{}) {
	fmt.Println(args...)
}

// GetType ...
func GetType(arg *ast.CXArgument) int {
    fieldCount := len(arg.Fields)
    if fieldCount > 0 {
        return GetType(arg.Fields[fieldCount - 1])
    }

    return arg.Type
}

// ExprOpName ...
func ExprOpName(expr *ast.CXExpression) string {
	if expr.Operator.IsAtomic {
		return OpNames[expr.Operator.OpCode]
	}
	return expr.Operator.Name

}

// func limitString (str string) string {
// 	if len(str) > 3
// }

func stackValueHeader(fileName string, fileLine int) string {
	return fmt.Sprintf("%s:%d", fileName, fileLine)
}

// PrintStack ...
func (cxprogram *ast.CXProgram) PrintStack() {
	fmt.Println()
	fmt.Println("===Callstack===")

	// we're going backwards in the stack
	fp := cxprogram.StackPointer

	for c := cxprogram.CallCounter; c >= 0; c-- {
		op := cxprogram.CallStack[c].Operator
		fp -= op.Size

		var dupNames []string

		fmt.Printf(">>> %s()\n", op.Name)

		for _, inp := range op.Inputs {
			fmt.Println("ProgramInput")
			fmt.Printf("\t%s : %s() : %s\n", stackValueHeader(inp.FileName, inp.FileLine), op.Name, GetPrintableValue(fp, inp))

			dupNames = append(dupNames, inp.Package.Name+inp.Name)
		}

		for _, out := range op.Outputs {
			fmt.Println("ProgramOutput")
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

// getFormattedDerefs is an auxiliary function for `GetFormattedName`. This
// function formats indexing and pointer dereferences associated to `arg`.
func getFormattedDerefs(arg *ast.CXArgument, includePkg bool) string {
	name := ""
	// Checking if we should include `arg`'s package name.
	if includePkg {
		name = fmt.Sprintf("%s.%s", arg.Package.Name, arg.Name)
	} else {
		name = arg.Name
	}

	// If it's a literal, just override the name with LITERAL_PLACEHOLDER.
	if arg.Name == "" {
		name = constants.LITERAL_PLACEHOLDER
	}

	// Checking if we got dereferences, e.g. **foo
	derefLevels := ""
	if arg.DereferenceLevels > 0 {
		for c := 0; c < arg.DereferenceLevels; c++ {
			derefLevels += "*"
		}
	}
	name = derefLevels + name

	// Checking if we have indexing operations, e.g. foo[2][1]
	for _, idx := range arg.Indexes {
		// Checking if the value is in data segment.
		// If this is the case, we can safely display it.
		idxValue := ""
		if idx.Offset > ast.PROGRAM.StackSize {
			// Then it's a literal.
			idxI32 := helper.Deserialize_i32(ast.PROGRAM.Memory[idx.Offset : idx.Offset+constants.TYPE_POINTER_SIZE])
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
func GetFormattedName(arg *ast.CXArgument, includePkg bool) string {
	// Getting formatted name which does not include fields.
	name := getFormattedDerefs(arg, includePkg)

	// Adding as suffixes all the fields.
	for _, fld := range arg.Fields {
		name = fmt.Sprintf("%s.%s", name, getFormattedDerefs(fld, includePkg))
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
func formatParameters(params []*ast.CXArgument) string {
	types := "("
	for i, param := range params {
		types += GetFormattedType(param)
		if i != len(params)-1 {
			types += ", "
		}
	}
	types += ")"

	return types
}

// GetFormattedType builds a string with the CXGO type representation of `arg`.
func GetFormattedType(arg *ast.CXArgument) string {
	typ := ""
	elt := GetAssignmentElement(arg)

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
			if elt.CustomType != nil {
				// then it's custom type
				typ += elt.CustomType.Name
			} else {
				// then it's basic type
				typ += constants.TypeNames[elt.Type]

				// If it's a function, let's add the inputs and outputs.
				if elt.Type == constants.TYPE_FUNC {
					if elt.IsLocalDeclaration {
						// Then it's a local variable, which can be assigned to a
						// lambda function, for example.
						typ += formatParameters(elt.Inputs)
						typ += formatParameters(elt.Outputs)
					} else {
						// Then it refers to a named function defined in a package.
						pkg, err := ast.PROGRAM.GetPackage(arg.Package.Name)
						if err != nil {
							println(ast.CompilationError(elt.FileName, elt.FileLine), err.Error())
							os.Exit(constants.CX_COMPILATION_ERROR)
						}

						fn, err := pkg.GetFunction(elt.Name)
						if err == nil {
							// println(CompilationError(elt.FileName, elt.FileLine), err.ProgramError())
							// os.Exit(CX_COMPILATION_ERROR)
							// Adding list of inputs and outputs types.
							typ += formatParameters(fn.Inputs)
							typ += formatParameters(fn.Outputs)
						}
					}
				}
			}
		}
	}

	return typ
}

// SignatureStringOfStruct returns the signature string of a struct.
func SignatureStringOfStruct(s *ast.CXStruct) string {
	fields := ""
	for _, f := range s.Fields {
		fields += fmt.Sprintf(" %s %s;", f.Name, GetFormattedType(f))
	}

	return fmt.Sprintf("%s struct {%s }", s.Name, fields)
}

// SignatureStringOfFunction returns the signature string of a function.
func SignatureStringOfFunction(pkg *ast.CXPackage, f *ast.CXFunction) string {
	var ins bytes.Buffer
	var outs bytes.Buffer
	getFormattedParam(f.Inputs, pkg, &ins)
	getFormattedParam(f.Outputs, pkg, &outs)

	return fmt.Sprintf("func %s(%s) (%s)",
		f.Name, ins.String(), outs.String())
}

// IsCorePackage ...
func IsCorePackage(ident string) bool {
	for _, core := range constants.CorePackages {
		if core == ident {
			return true
		}
	}
	return false
}

// IsTempVar ...
func IsTempVar(name string) bool {
	if len(name) >= len(constants.LOCAL_PREFIX) && name[:len(constants.LOCAL_PREFIX)] == constants.LOCAL_PREFIX {
		return true
	}
	return false
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

func checkForEscapedChars(str string) []byte {
	var res []byte
	var lenStr = int(len(str))
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

// GetAssignmentElement ...
func GetAssignmentElement(arg *ast.CXArgument) *ast.CXArgument {
	if len(arg.Fields) > 0 {
		return arg.Fields[len(arg.Fields)-1]
	}
	return arg

}

// IsValidSliceIndex ...
func IsValidSliceIndex(offset int, index int, sizeofElement int) bool {
	sliceLen := GetSliceLen(int32(offset))
	bytesLen := sliceLen * int32(sizeofElement)
	index -= constants.OBJECT_HEADER_SIZE + constants.SLICE_HEADER_SIZE + offset

	if index >= 0 && index < int(bytesLen) && (index%sizeofElement) == 0 {
		return true
	}
	return false
}

// GetPointerOffset ...
func GetPointerOffset(pointer int32) int32 {
	return helper.Deserialize_i32(ast.PROGRAM.Memory[pointer : pointer+constants.TYPE_POINTER_SIZE])
}

// GetSliceOffset ...
func GetSliceOffset(fp int, arg *ast.CXArgument) int32 {
	element := GetAssignmentElement(arg)
	if element.IsSlice {
		return GetPointerOffset(int32(GetFinalOffset(fp, arg)))
	}

	return -1
}

// GetObjectHeader ...
func GetObjectHeader(offset int32) []byte {
	return ast.PROGRAM.Memory[offset : offset+constants.OBJECT_HEADER_SIZE]
}

// GetSliceHeader ...
func GetSliceHeader(offset int32) []byte {
	return ast.PROGRAM.Memory[offset+constants.OBJECT_HEADER_SIZE : offset+constants.OBJECT_HEADER_SIZE+constants.SLICE_HEADER_SIZE]
}

// GetSliceLen ...
func GetSliceLen(offset int32) int32 {
	sliceHeader := GetSliceHeader(offset)
	return helper.Deserialize_i32(sliceHeader[4:8])
}

// GetSlice ...
func GetSlice(offset int32, sizeofElement int) []byte {
	if offset > 0 {
		sliceLen := GetSliceLen(offset)
		if sliceLen > 0 {
			dataOffset := offset + constants.OBJECT_HEADER_SIZE + constants.SLICE_HEADER_SIZE - 4
			dataLen := 4 + sliceLen*int32(sizeofElement)
			return ast.PROGRAM.Memory[dataOffset : dataOffset+dataLen]
		}
	}
	return nil
}

// GetSliceData ...
func GetSliceData(offset int32, sizeofElement int) []byte {
	if slice := GetSlice(offset, sizeofElement); slice != nil {
		return slice[4:]
	}
	return nil
}

// SliceResizeEx does the logic required by `SliceResize`. It is separated because some other functions might have access to the offsets of the slices, but not the `CXArgument`s.
func SliceResizeEx(outputSliceOffset int32, count int32, sizeofElement int) int {
	if count < 0 {
		panic(constants.CX_RUNTIME_SLICE_INDEX_OUT_OF_RANGE) // TODO : should use uint32
	}

	var outputSliceHeader []byte
	var outputSliceCap int32

	if outputSliceOffset > 0 {
		outputSliceHeader = GetSliceHeader(outputSliceOffset)
		outputSliceCap = helper.Deserialize_i32(outputSliceHeader[0:4])
	}

	var newLen = count
	var newCap = outputSliceCap
	if newLen > newCap {
		if newCap <= 0 {
			newCap = newLen
		} else {
			newCap *= 2
		}
		var outputObjectSize = constants.OBJECT_HEADER_SIZE + constants.SLICE_HEADER_SIZE + newCap*int32(sizeofElement)
		outputSliceOffset = int32(AllocateSeq(int(outputObjectSize)))
		WriteMemI32(GetObjectHeader(outputSliceOffset)[5:9], 0, outputObjectSize)

		outputSliceHeader = GetSliceHeader(outputSliceOffset)
		WriteMemI32(outputSliceHeader[0:4], 0, newCap)
		WriteMemI32(outputSliceHeader[4:8], 0, newLen)
	}

	return int(outputSliceOffset)
}

// SliceResize ...
func SliceResize(fp int, out *ast.CXArgument, inp *ast.CXArgument, count int32, sizeofElement int) int {
	outputSliceOffset := GetSliceOffset(fp, out)

	outputSliceOffset = int32(SliceResizeEx(outputSliceOffset, count, sizeofElement))

	SliceCopy(fp, outputSliceOffset, inp, count, sizeofElement)

	return int(outputSliceOffset)
}

// SliceCopyEx does the logic required by `SliceCopy`. It is separated because some other functions might have access to the offsets of the slices, but not the `CXArgument`s.
func SliceCopyEx(outputSliceOffset int32, inputSliceOffset int32, count int32, sizeofElement int) {
	if count < 0 {
		panic(constants.CX_RUNTIME_SLICE_INDEX_OUT_OF_RANGE) // TODO : should use uint32
	}

	var inputSliceLen int32
	if inputSliceOffset != 0 {
		inputSliceLen = GetSliceLen(inputSliceOffset)
	}

	if outputSliceOffset > 0 {
		outputSliceHeader := GetSliceHeader(outputSliceOffset)
		WriteMemI32(outputSliceHeader[4:8], 0, count)
		outputSliceData := GetSliceData(outputSliceOffset, sizeofElement)
		if (outputSliceOffset != inputSliceOffset) && inputSliceLen > 0 {
			copy(outputSliceData, GetSliceData(inputSliceOffset, sizeofElement))
		}
	}
}

// SliceCopy copies the contents from the slice located at `inputSliceOffset` to the slice located at `outputSliceOffset`.
func SliceCopy(fp int, outputSliceOffset int32, inp *ast.CXArgument, count int32, sizeofElement int) {
	inputSliceOffset := GetSliceOffset(fp, inp)
	SliceCopyEx(outputSliceOffset, inputSliceOffset, count, sizeofElement)
}

// SliceAppendResize prepares a slice to be able to store a new object of length `sizeofElement`. It checks if the slice needs to be relocated in memory, and if it is needed it relocates it and a new `outputSliceOffset` is calculated for the new slice.
func SliceAppendResize(fp int, out *ast.CXArgument, inp *ast.CXArgument, sizeofElement int) int32 {
	inputSliceOffset := GetSliceOffset(fp, inp)
	var inputSliceLen int32
	if inputSliceOffset != 0 {
		inputSliceLen = GetSliceLen(inputSliceOffset)
	}

	// TODO: Are we limited then to only one element for now? (because of that +1)
	outputSliceOffset := int32(SliceResize(fp, out, inp, inputSliceLen+1, sizeofElement))
	return outputSliceOffset
}

// SliceAppendWrite writes `object` to a slice that is guaranteed to be able to hold `object`, i.e. it had to be checked by `SliceAppendResize` first in case it needed to be resized.
func SliceAppendWrite(outputSliceOffset int32, object []byte, index int32) {
	sizeofElement := len(object)
	outputSliceData := GetSliceData(outputSliceOffset, sizeofElement)
	copy(outputSliceData[int(index)*sizeofElement:], object)
}

// SliceAppendWriteByte writes `object` to a slice that is guaranteed to be able to hold `object`, i.e. it had to be checked by `SliceAppendResize` first in case it needed to be resized.
func SliceAppendWriteByte(outputSliceOffset int32, object []byte, index int32) {
	outputSliceData := GetSliceData(outputSliceOffset, 1)
	copy(outputSliceData[int(index):], object)
}

// SliceInsert ...
func SliceInsert(fp int, out *ast.CXArgument, inp *ast.CXArgument, index int32, object []byte) int {
	inputSliceOffset := GetSliceOffset(fp, inp)
	// outputSliceOffset := GetSliceOffset(fp, out)

	var inputSliceLen int32
	if inputSliceOffset != 0 {
		inputSliceLen = GetSliceLen(inputSliceOffset)
	}

	if index < 0 || index > inputSliceLen {
		panic(constants.CX_RUNTIME_SLICE_INDEX_OUT_OF_RANGE)
	}

	var newLen = inputSliceLen + 1
	sizeofElement := len(object)
	outputSliceOffset := int32(SliceResize(fp, out, inp, newLen, sizeofElement))
	outputSliceData := GetSliceData(outputSliceOffset, sizeofElement)
	copy(outputSliceData[int(index+1)*sizeofElement:], outputSliceData[int(index)*sizeofElement:])
	copy(outputSliceData[int(index)*sizeofElement:], object)
	return int(outputSliceOffset)
}

// SliceRemove ...
func SliceRemove(fp int, out *ast.CXArgument, inp *ast.CXArgument, index int32, sizeofElement int32) int {
	inputSliceOffset := GetSliceOffset(fp, inp)
	outputSliceOffset := GetSliceOffset(fp, out)

	var inputSliceLen int32
	if inputSliceOffset != 0 {
		inputSliceLen = GetSliceLen(inputSliceOffset)
	}

	if index < 0 || index >= inputSliceLen {
		panic(constants.CX_RUNTIME_SLICE_INDEX_OUT_OF_RANGE)
	}

	outputSliceData := GetSliceData(outputSliceOffset, int(sizeofElement))
	copy(outputSliceData[index*sizeofElement:], outputSliceData[(index+1)*sizeofElement:])
	outputSliceOffset = int32(SliceResize(fp, out, inp, inputSliceLen-1, int(sizeofElement)))
	return int(outputSliceOffset)
}

// WriteToSlice is used to create slices in the backend, i.e. not by calling `append`
// in a CX program, but rather by the CX code itself. This function is used by
// affordances, serialization and to store OS input arguments.
func WriteToSlice(off int, inp []byte) int {
	// TODO: Check all these parses from/to int32/int.
	var inputSliceLen int32
	if off != 0 {
		inputSliceLen = GetSliceLen(int32(off))
	}

	inpLen := len(inp)
	// We first check if a resize is needed. If a resize occurred
	// the address of the new slice will be stored in `newOff` and will
	// be different to `off`.
	newOff := SliceResizeEx(int32(off), inputSliceLen+1, inpLen)

	// Copy the data from the old slice at `off` to `newOff`.
	SliceCopyEx(int32(newOff), int32(off), inputSliceLen+1, inpLen)

	// Write the new slice element `inp` to the slice located at `newOff`.
	SliceAppendWrite(int32(newOff), inp, inputSliceLen)
	return newOff

}

func getNonCollectionValue(fp int, arg, elt *ast.CXArgument, typ string) string {
	switch typ {
	case "bool":
		return fmt.Sprintf("%v", ReadBool(fp, elt))
	case "str":
		return fmt.Sprintf("%v", ReadStr(fp, elt))
	case "i8":
		return fmt.Sprintf("%v", ReadI8(fp, elt))
	case "i16":
		return fmt.Sprintf("%v", ReadI16(fp, elt))
	case "i32":
		return fmt.Sprintf("%v", ReadI32(fp, elt))
	case "i64":
		return fmt.Sprintf("%v", ReadI64(fp, elt))
	case "ui8":
		return fmt.Sprintf("%v", ReadUI8(fp, elt))
	case "ui16":
		return fmt.Sprintf("%v", ReadUI16(fp, elt))
	case "ui32":
		return fmt.Sprintf("%v", ReadUI32(fp, elt))
	case "ui64":
		return fmt.Sprintf("%v", ReadUI64(fp, elt))
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

// GetPrintableValue ...
func GetPrintableValue(fp int, arg *ast.CXArgument) string {
	var typ string
	elt := GetAssignmentElement(arg)
	if elt.CustomType != nil {
		// then it's custom type
		typ = elt.CustomType.Name
	} else {
		// then it's native type
		typ = constants.TypeNames[elt.Type]
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

// DebugHeap prints the symbols that are acting as pointers in a CX program at certain point during the execution of the program along with the addresses they are pointing. Additionally, a list of the objects in the heap is printed, which shows their address in the heap, if they are marked as alive or as dead by the garbage collector, the address where they used to live after a garbage collector call, the full size of the object, the object itself as a slice of bytes and the pointers that are pointing to that object.
func DebugHeap() {
	// symsToAddrs will hold a list of symbols that are pointing to an address.
	symsToAddrs := make(map[int32][]string)

	// Processing global variables. Adding the address they are pointing to.
	for _, pkg := range ast.PROGRAM.Packages {
		for _, glbl := range pkg.Globals {
			if glbl.IsPointer || glbl.IsSlice {
				heapOffset := helper.Deserialize_i32(ast.PROGRAM.Memory[glbl.Offset : glbl.Offset+constants.TYPE_POINTER_SIZE])

				symsToAddrs[heapOffset] = append(symsToAddrs[heapOffset], glbl.Name)
			}
		}
	}

	// Processing local variables in every active function call in the `CallStack`.
	// Adding the address they are pointing to.
	var fp int
	for c := 0; c <= ast.PROGRAM.CallCounter; c++ {
		op := ast.PROGRAM.CallStack[c].Operator

		// TODO: Some standard library functions "manually" add a function
		// call (callbacks) to `PRGRM.CallStack`. These functions do not have an
		// operator associated to them. This can be considered as a bug or as an
		// undesirable mechanic.
		// [2019-06-24 Mon 22:39] Actually, if the GC is triggered in the middle
		// of a callback, things will certainly break.
		if op == nil {
			continue
		}

		for _, ptr := range op.ListOfPointers {
			offset := ptr.Offset
			symName := ptr.Name
			if len(ptr.Fields) > 0 {
				fld := ptr.Fields[len(ptr.Fields)-1]
				offset += fld.Offset
				symName += "." + fld.Name
			}

			if ptr.Offset < ast.PROGRAM.StackSize {
				offset += fp
			}

			heapOffset := helper.Deserialize_i32(ast.PROGRAM.Memory[offset : offset+constants.TYPE_POINTER_SIZE])

			symsToAddrs[heapOffset] = append(symsToAddrs[heapOffset], symName)
		}

		fp += op.Size
	}

	// Printing all the details.
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, '.', 0)

	for off, symNames := range symsToAddrs {
		var addrB [4]byte
		WriteMemI32(addrB[:], 0, off)
		fmt.Fprintln(w, "Addr:\t", addrB, "\tPtr:\t", symNames)
	}

	// Just a newline.
	fmt.Fprintln(w)
	w.Flush()

	w = tabwriter.NewWriter(os.Stdout, 0, 0, 2, '.', 0)

	for c := ast.PROGRAM.HeapStartsAt + constants.NULL_HEAP_ADDRESS_OFFSET; c < ast.PROGRAM.HeapStartsAt+ast.PROGRAM.HeapPointer; {
		objSize := helper.Deserialize_i32(ast.PROGRAM.Memory[c+constants.MARK_SIZE+constants.FORWARDING_ADDRESS_SIZE : c+constants.MARK_SIZE+constants.FORWARDING_ADDRESS_SIZE+constants.OBJECT_SIZE])

		// Setting a limit size for the object to be printed if the object is too large.
		// We don't want to print obscenely large objects to standard output.
		printObjSize := objSize
		if objSize > 50 {
			printObjSize = 50
		}

		var addrB [4]byte
		WriteMemI32(addrB[:], 0, int32(c))

		fmt.Fprintln(w, "Addr:\t", addrB, "\tMark:\t", ast.PROGRAM.Memory[c:c+constants.MARK_SIZE], "\tFwd:\t", ast.PROGRAM.Memory[c+constants.MARK_SIZE:c+constants.MARK_SIZE+constants.FORWARDING_ADDRESS_SIZE], "\tSize:\t", objSize, "\tObj:\t", ast.PROGRAM.Memory[c+constants.OBJECT_HEADER_SIZE:c+int(printObjSize)], "\tPtrs:", symsToAddrs[int32(c)])

		c += int(objSize)
	}

	// Just a newline.
	fmt.Fprintln(w)
	w.Flush()
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

		fi, err := CXStatFile(arg)
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
				file, err := CXOpenFile(path)

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
			file, err := CXOpenFile(arg)

			if err != nil {
				panic(err)
			}

			fileNames = append(fileNames, file.Name())
			sourceCode = append(sourceCode, file)
		}
	}

	return cxArgs, sourceCode, fileNames
}

// IsPointer checks if `sym` is a candidate for the garbage collector to check.
// For example, if `sym` is a slice, the garbage collector will need to check
// if the slice on the heap needs to be relocated.
func IsPointer(sym *ast.CXArgument) bool {
	// There's no need to add global variables in `fn.ListOfPointers` as we can access them easily through `CXPackage.Globals`
	// TODO: We could still pre-compute a list of candidates for globals.
	if sym.Offset >= ast.PROGRAM.StackSize && sym.Name != "" {
		return false
	}
	// NOTE: Strings are considered as `IsPointer`s by the runtime.
	// if (sym.IsPointer || sym.IsSlice) && sym.Name != "" {
	// 	return true
	// }
	if (sym.IsPointer || sym.IsSlice) && sym.Name != "" && len(sym.Fields) == 0 {
		return true
	}
	if sym.Type == constants.TYPE_STR && sym.Name != "" && len(sym.Fields) == 0 {
		return true
	}
	// if (sym.Type == TYPE_STR && sym.Name != "") {
	// 	return true
	// }
	// If `sym` is a structure instance, we need to check if the last field
	// being access is a pointer candidate
	// if len(sym.Fields) > 0 {
	// 	return isPointer(sym.Fields[len(sym.Fields)-1])
	// }
	return false
}
