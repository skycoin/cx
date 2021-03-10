package cxcore

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime/debug"
	"strconv"
	"strings"

	"github.com/skycoin/skycoin/src/cipher/encoder"
)

// ExprOpName ...
func ExprOpName(expr *CXExpression) string {
	if expr.Operator.IsNative {
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

// getFormattedDerefs is an auxiliary function for `GetFormattedName`. This
// function formats indexing and pointer dereferences associated to `arg`.
func getFormattedDerefs(arg *CXArgument, includePkg bool) string {
	name := ""
	// Checking if we should include `arg`'s package name.
	if includePkg {
		name = fmt.Sprintf("%s.%s", arg.Package.Name, arg.Name)
	} else {
		name = arg.Name
	}

	// If it's a literal, just override the name with LITERAL_PLACEHOLDER.
	if arg.Name == "" {
		name = LITERAL_PLACEHOLDER
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
		if idx.Offset > PROGRAM.StackSize {
			// Then it's a literal.
			idxI32 := Deserialize_i32(PROGRAM.Memory[idx.Offset : idx.Offset+TYPE_POINTER_SIZE])
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
func GetFormattedName(arg *CXArgument, includePkg bool) string {
	// Getting formatted name which does not include fields.
	name := getFormattedDerefs(arg, includePkg)

	// Adding as suffixes all the fields.
	for _, fld := range arg.Fields {
		name = fmt.Sprintf("%s.%s", name, getFormattedDerefs(fld, includePkg))
	}

	// Checking if we're referencing `arg`.
	if arg.PassBy == PASSBY_REFERENCE {
		name = "&" + name
	}

	return name
}

// formatParameters returns a string containing a list of the formatted types of
// each of `params`, enclosed in parethesis. This function is used only when
// formatting functions as first-class objects.
func formatParameters(params []*CXArgument) string {
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
func GetFormattedType(arg *CXArgument) string {
	typ := ""
	elt := GetAssignmentElement(arg)

	// this is used to know what arg.Lengths index to use
	// used for cases like [5]*[3]i32, where we jump to another decl spec
	arrDeclCount := len(arg.Lengths) - 1
	// looping declaration specifiers
	for _, spec := range elt.DeclarationSpecifiers {
		switch spec {
		case DECL_POINTER:
			typ = "*" + typ
		case DECL_DEREF:
			typ = typ[1:]
		case DECL_ARRAY:
			typ = fmt.Sprintf("[%d]%s", arg.Lengths[arrDeclCount], typ)
			arrDeclCount--
		case DECL_SLICE:
			typ = "[]" + typ
		case DECL_INDEXING:
		default:
			// base type
			if elt.CustomType != nil {
				// then it's custom type
				typ += elt.CustomType.Name
			} else {
				// then it's basic type
				typ += TypeNames[elt.Type]

				// If it's a function, let's add the inputs and outputs.
				if elt.Type == TYPE_FUNC {
					if elt.IsLocalDeclaration {
						// Then it's a local variable, which can be assigned to a
						// lambda function, for example.
						typ += formatParameters(elt.Inputs)
						typ += formatParameters(elt.Outputs)
					} else {
						// Then it refers to a named function defined in a package.
						pkg, err := PROGRAM.GetPackage(arg.Package.Name)
						if err != nil {
							println(CompilationError(elt.FileName, elt.FileLine), err.Error())
							os.Exit(CX_COMPILATION_ERROR)
						}

						fn, err := pkg.GetFunction(elt.Name)
						if err == nil {
							// println(CompilationError(elt.FileName, elt.FileLine), err.Error())
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
func SignatureStringOfStruct(s *CXStruct) string {
	fields := ""
	for _, f := range s.Fields {
		fields += fmt.Sprintf(" %s %s;", f.Name, GetFormattedType(f))
	}

	return fmt.Sprintf("%s struct {%s }", s.Name, fields)
}

// SignatureStringOfFunction returns the signature string of a function.
func SignatureStringOfFunction(pkg *CXPackage, f *CXFunction) string {
	var ins bytes.Buffer
	var outs bytes.Buffer
	getFormattedParam(f.Inputs, pkg, &ins)
	getFormattedParam(f.Outputs, pkg, &outs)

	return fmt.Sprintf("func %s(%s) (%s)",
		f.Name, ins.String(), outs.String())
}

// CheckArithmeticOp ...
func CheckArithmeticOp(expr *CXExpression) bool {
	if expr.Operator.IsNative {
		switch expr.Operator.OpCode {
		case OP_I32_MUL, OP_I32_DIV, OP_I32_MOD, OP_I32_ADD,
			OP_I32_SUB, OP_I32_NEG, OP_I32_BITSHL, OP_I32_BITSHR, OP_I32_LT,
			OP_I32_GT, OP_I32_LTEQ, OP_I32_GTEQ, OP_I32_EQ, OP_I32_UNEQ,
			OP_I32_BITAND, OP_I32_BITXOR, OP_I32_BITOR, OP_STR_EQ:
			return true
		}
	}
	return false
}

// IsCorePackage ...
func IsCorePackage(ident string) bool {
	for _, core := range CorePackages {
		if core == ident {
			return true
		}
	}
	return false
}

// IsTempVar ...
func IsTempVar(name string) bool {
	if len(name) >= len(LOCAL_PREFIX) && name[:len(LOCAL_PREFIX)] == LOCAL_PREFIX {
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

// GetArgSize ...
func GetArgSize(typ int) int {
	switch typ {
	case TYPE_BOOL, TYPE_I8, TYPE_UI8:
		return 1
	case TYPE_I16, TYPE_UI16:
		return 2
	case TYPE_STR, TYPE_I32, TYPE_UI32, TYPE_F32, TYPE_AFF:
		return 4
	case TYPE_I64, TYPE_UI64, TYPE_F64:
		return 8
	default:
		return 4
		//return -1 // should be panic
		//panic(CX_INTERNAL_ERROR)
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
func GetAssignmentElement(arg *CXArgument) *CXArgument {
	if len(arg.Fields) > 0 {
		return arg.Fields[len(arg.Fields)-1]
	}
	return arg

}

// IsValidSliceIndex ...
func IsValidSliceIndex(offset int, index int, sizeofElement int) bool {
	sliceLen := GetSliceLen(int32(offset))
	bytesLen := sliceLen * int32(sizeofElement)
	index -= OBJECT_HEADER_SIZE + SLICE_HEADER_SIZE + offset

	if index >= 0 && index < int(bytesLen) && (index%sizeofElement) == 0 {
		return true
	}
	return false
}

// GetPointerOffset ...
func GetPointerOffset(pointer int32) int32 {
	return Deserialize_i32(PROGRAM.Memory[pointer : pointer+TYPE_POINTER_SIZE])
}

// GetSliceOffset ...
func GetSliceOffset(fp int, arg *CXArgument) int32 {
	element := GetAssignmentElement(arg)
	if element.IsSlice {
		return GetPointerOffset(int32(GetFinalOffset(fp, arg)))
	}

	return -1
}

// GetObjectHeader ...
func GetObjectHeader(offset int32) []byte {
	return PROGRAM.Memory[offset : offset+OBJECT_HEADER_SIZE]
}

// GetSliceHeader ...
func GetSliceHeader(offset int32) []byte {
	return PROGRAM.Memory[offset+OBJECT_HEADER_SIZE : offset+OBJECT_HEADER_SIZE+SLICE_HEADER_SIZE]
}

// GetSliceLen ...
func GetSliceLen(offset int32) int32 {
	sliceHeader := GetSliceHeader(offset)
	return Deserialize_i32(sliceHeader[4:8])
}

// GetSlice ...
func GetSlice(offset int32, sizeofElement int) []byte {
	if offset > 0 {
		sliceLen := GetSliceLen(offset)
		if sliceLen > 0 {
			dataOffset := offset + OBJECT_HEADER_SIZE + SLICE_HEADER_SIZE - 4
			dataLen := 4 + sliceLen*int32(sizeofElement)
			return PROGRAM.Memory[dataOffset : dataOffset+dataLen]
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
		panic(CX_RUNTIME_SLICE_INDEX_OUT_OF_RANGE) // TODO : should use uint32
	}

	var outputSliceHeader []byte
	var outputSliceCap int32

	if outputSliceOffset > 0 {
		outputSliceHeader = GetSliceHeader(outputSliceOffset)
		outputSliceCap = Deserialize_i32(outputSliceHeader[0:4])
	}

	var newLen = count
	var newCap = outputSliceCap
	if newLen > newCap {
		if newCap <= 0 {
			newCap = newLen
		} else {
			newCap *= 2
		}
		var outputObjectSize = OBJECT_HEADER_SIZE + SLICE_HEADER_SIZE + newCap*int32(sizeofElement)
		outputSliceOffset = int32(AllocateSeq(int(outputObjectSize)))
		WriteMemI32(GetObjectHeader(outputSliceOffset)[5:9], 0, outputObjectSize)

		outputSliceHeader = GetSliceHeader(outputSliceOffset)
		WriteMemI32(outputSliceHeader[0:4], 0, newCap)
		WriteMemI32(outputSliceHeader[4:8], 0, newLen)
	}

	return int(outputSliceOffset)
}

// SliceResize ...
func SliceResize(fp int, out *CXArgument, inp *CXArgument, count int32, sizeofElement int) int {
	outputSliceOffset := GetSliceOffset(fp, out)

	outputSliceOffset = int32(SliceResizeEx(outputSliceOffset, count, sizeofElement))

	SliceCopy(fp, outputSliceOffset, inp, count, sizeofElement)

	return int(outputSliceOffset)
}

// SliceCopyEx does the logic required by `SliceCopy`. It is separated because some other functions might have access to the offsets of the slices, but not the `CXArgument`s.
func SliceCopyEx(outputSliceOffset int32, inputSliceOffset int32, count int32, sizeofElement int) {
	if count < 0 {
		panic(CX_RUNTIME_SLICE_INDEX_OUT_OF_RANGE) // TODO : should use uint32
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
func SliceCopy(fp int, outputSliceOffset int32, inp *CXArgument, count int32, sizeofElement int) {
	inputSliceOffset := GetSliceOffset(fp, inp)
	SliceCopyEx(outputSliceOffset, inputSliceOffset, count, sizeofElement)
}

// SliceAppendResize prepares a slice to be able to store a new object of length `sizeofElement`. It checks if the slice needs to be relocated in memory, and if it is needed it relocates it and a new `outputSliceOffset` is calculated for the new slice.
func SliceAppendResize(fp int, out *CXArgument, inp *CXArgument, sizeofElement int) int32 {
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
func SliceInsert(fp int, out *CXArgument, inp *CXArgument, index int32, object []byte) int {
	inputSliceOffset := GetSliceOffset(fp, inp)
	// outputSliceOffset := GetSliceOffset(fp, out)

	var inputSliceLen int32
	if inputSliceOffset != 0 {
		inputSliceLen = GetSliceLen(inputSliceOffset)
	}

	if index < 0 || index > inputSliceLen {
		panic(CX_RUNTIME_SLICE_INDEX_OUT_OF_RANGE)
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
func SliceRemove(fp int, out *CXArgument, inp *CXArgument, index int32, sizeofElement int32) int {
	inputSliceOffset := GetSliceOffset(fp, inp)
	outputSliceOffset := GetSliceOffset(fp, out)

	var inputSliceLen int32
	if inputSliceOffset != 0 {
		inputSliceLen = GetSliceLen(inputSliceOffset)
	}

	if index < 0 || index >= inputSliceLen {
		panic(CX_RUNTIME_SLICE_INDEX_OUT_OF_RANGE)
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

// refactoring reuse in WriteObject and WriteObjectRetOff
func writeObj(obj []byte) int {
	// QUARENTINED: Check if `newwriteObj` can supersede this `writeObj`.
	// Especially check usage on CX chains.
	size := len(obj) + OBJECT_HEADER_SIZE
	sizeB := encoder.SerializeAtomic(int32(size))
	// heapOffset := AllocateSeq(size + OBJECT_HEADER_SIZE)
	heapOffset := AllocateSeq(size)

	// var finalObj = make([]byte, OBJECT_HEADER_SIZE+size)
	var finalObj = make([]byte, size)

	for c := OBJECT_GC_HEADER_SIZE; c < OBJECT_HEADER_SIZE; c++ {
		finalObj[c] = sizeB[c-OBJECT_GC_HEADER_SIZE]
	}
	// for c := OBJECT_HEADER_SIZE; c < size+OBJECT_HEADER_SIZE; c++ {
	for c := OBJECT_HEADER_SIZE; c < size; c++ {
		finalObj[c] = obj[c-OBJECT_HEADER_SIZE]
	}

	WriteMemory(heapOffset, finalObj)
	return heapOffset
}

// NewWriteObj refactoring reuse in WriteObject and WriteObjectRetOff
func NewWriteObj(obj []byte) int {
	// 2dbug introduces this new version of `writeObj`. It is unknown to me
	// (amherag) at the moment if it is safe to replace `writeObj` with
	// this version. Leaving `writeObj` in quarentine.
	size := len(obj)
	heapOffset := AllocateSeq(size + OBJECT_HEADER_SIZE)
	var finalObj = make([]byte, OBJECT_HEADER_SIZE+size)

	WriteMemI32(finalObj, OBJECT_GC_HEADER_SIZE, int32(size))
	for c := OBJECT_HEADER_SIZE; c < size+OBJECT_HEADER_SIZE; c++ {
		finalObj[c] = obj[c-OBJECT_HEADER_SIZE]
	}

	WriteMemory(heapOffset, finalObj)
	return heapOffset
}

// WriteObject ...
func WriteObject(out1Offset int, obj []byte) {
	WriteI32(out1Offset, int32(NewWriteObj(obj)))
}

// WriteObjectRetOff ...
func WriteObjectRetOff(obj []byte) int {
	return NewWriteObj(obj)
}

// ErrorHeader ...
func ErrorHeader(currentFile string, lineNo int) string {
	return "error: " + currentFile + ":" + strconv.FormatInt(int64(lineNo), 10)
}

// CompilationError is a helper function that concatenates the `currentFile` and `lineNo` data to a error header and returns the full error string.
func CompilationError(currentFile string, lineNo int) string {
	FoundCompileErrors = true
	return ErrorHeader(currentFile, lineNo)
}

// ErrorString ...
func ErrorString(code int) string {
	if str, found := ErrorStrings[code]; found {
		return str
	}
	return ErrorStrings[CX_RUNTIME_ERROR]
}

func errorCode(r interface{}) int {
	switch v := r.(type) {
	case int:
		return int(v)
	default:
		return CX_RUNTIME_ERROR
	}
}

func runtimeErrorInfo(r interface{}, printStack bool, defaultError int) {
	call := PROGRAM.CallStack[PROGRAM.CallCounter]
	expr := call.Operator.Expressions[call.Line]
	code := errorCode(r)
	if code == CX_RUNTIME_ERROR {
		code = defaultError
	}

	fmt.Printf("%s, %s, %v", ErrorHeader(expr.FileName, expr.FileLine), ErrorString(code), r)

	if printStack {
		PROGRAM.PrintStack()
	}

	if DBG_GOLANG_STACK_TRACE {
		debug.PrintStack()
	}

	os.Exit(code)
}

// RuntimeError ...
func RuntimeError() {
	if r := recover(); r != nil {
		switch r {
		case STACK_OVERFLOW_ERROR:
			call := PROGRAM.CallStack[PROGRAM.CallCounter]
			if PROGRAM.CallCounter > 0 {
				PROGRAM.CallCounter--
				PROGRAM.StackPointer = call.FramePointer
				runtimeErrorInfo(r, true, CX_RUNTIME_STACK_OVERFLOW_ERROR)
			} else {
				// error at entry point
				runtimeErrorInfo(r, false, CX_RUNTIME_STACK_OVERFLOW_ERROR)
			}
		case HEAP_EXHAUSTED_ERROR:
			runtimeErrorInfo(r, true, CX_RUNTIME_HEAP_EXHAUSTED_ERROR)
		default:
			runtimeErrorInfo(r, true, CX_RUNTIME_ERROR)
		}
	}
}

func getNonCollectionValue(fp int, arg, elt *CXArgument, typ string) string {
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
			os.Exit(CX_COMPILATION_ERROR)
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
					os.Exit(CX_COMPILATION_ERROR)
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
func IsPointer(sym *CXArgument) bool {
	// There's no need to add global variables in `fn.ListOfPointers` as we can access them easily through `CXPackage.Globals`
	// TODO: We could still pre-compute a list of candidates for globals.
	if sym.Offset >= PROGRAM.StackSize && sym.Name != "" {
		return false
	}
	// NOTE: Strings are considered as `IsPointer`s by the runtime.
	// if (sym.IsPointer || sym.IsSlice) && sym.Name != "" {
	// 	return true
	// }
	if (sym.IsPointer || sym.IsSlice) && sym.Name != "" && len(sym.Fields) == 0 {
		return true
	}
	if sym.Type == TYPE_STR && sym.Name != "" && len(sym.Fields) == 0 {
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

// WriteStringObj writes `str` to the heap as an object and returns its absolute offset.
func WriteStringObj(str string) int {
	strB := encoder.Serialize(str)
	return NewWriteObj(strB)
}

// ReadStringFromObject reads the string located at offset `off`.
func ReadStringFromObject(off int32) string {
	var plusOff int32
	if int(off) > PROGRAM.HeapStartsAt {
		// Found in heap segment.
		plusOff += OBJECT_HEADER_SIZE
	}

	size := Deserialize_i32(PROGRAM.Memory[off+plusOff : off+plusOff+STR_HEADER_SIZE])

	str := ""
	_, err := encoder.DeserializeRaw(PROGRAM.Memory[off+plusOff:off+plusOff+STR_HEADER_SIZE+size], &str)
	if err != nil {
		panic(err)
	}
	return str
}
