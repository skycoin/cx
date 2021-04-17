package ast

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	constants2 "github.com/skycoin/cx/cxparser/constants"

	"github.com/skycoin/cx/cx/constants"
	"github.com/skycoin/cx/cx/helper"
	"github.com/skycoin/cx/cx/util"
)

// ToString returns the abstract syntax tree of a CX program in a
// string format.
func ToString(cxprogram *CXProgram) string {
	var ast3 string
	ast3 += "Program\n" //why is top line "Program" ???

	var currentFunction *CXFunction
	var currentPackage *CXPackage

	currentPackage, err := cxprogram.GetCurrentPackage()

	if err != nil {
		panic("CXProgram.ToString(): error, currentPackage is nil")
	}

	currentFunction, _ = cxprogram.GetCurrentFunction()
	currentPackage.CurrentFunction = currentFunction

	BuildStrPackages(cxprogram, &ast3) //what does this do?

	return ast3
}

// buildStrImports is an auxiliary function for `toString`. It builds
// string representation all the imported packages of `pkg`.
func buildStrImports(pkg *CXPackage, ast *string) {
	if len(pkg.Imports) > 0 {
		*ast += "\tImports\n"
	}

	for j, imp := range pkg.Imports {
		*ast += fmt.Sprintf("\t\t%d.- Import: %s\n", j, imp.Name)
	}
}

// buildStrGlobals is an auxiliary function for `toString`. It builds
// string representation of all the global variables of `pkg`.
func buildStrGlobals(pkg *CXPackage, ast *string) {
	if len(pkg.Globals) > 0 {
		*ast += "\tGlobals\n"
	}

	for j, v := range pkg.Globals {
		*ast += fmt.Sprintf("\t\t%d.- Global: %s %s\n", j, v.ArgDetails.Name, GetFormattedType(v))
	}
}

// buildStrStructs is an auxiliary function for `toString`. It builds
// string representation of all the structures defined in `pkg`.
func buildStrStructs(pkg *CXPackage, ast *string) {
	if len(pkg.Structs) > 0 {
		*ast += "\tStructs\n"
	}

	for j, strct := range pkg.Structs {
		*ast += fmt.Sprintf("\t\t%d.- Struct: %s\n", j, strct.Name)

		for k, fld := range strct.Fields {
			*ast += fmt.Sprintf("\t\t\t%d.- Field: %s %s\n",
				k, fld.ArgDetails.Name, GetFormattedType(fld))
		}
	}
}

// buildStrFunctions is an auxiliary function for `toString`. It builds
// string representation of all the functions defined in `pkg`.
func buildStrFunctions(pkg *CXPackage, ast1 *string) {
	if len(pkg.Functions) > 0 {
		*ast1 += "\tFunctions\n"
	}

	// We need to declare the counter outside so we can
	// ignore the increment from the `*init` function.
	var j int
	for _, fn := range pkg.Functions {
		if fn.Name == constants.SYS_INIT_FUNC {
			continue
		}
		_, err := pkg.SelectFunction(fn.Name)
		if err != nil {
			panic(err)
		}

		var inps bytes.Buffer
		var outs bytes.Buffer
		getFormattedParam(fn.Inputs, pkg, &inps)
		getFormattedParam(fn.Outputs, pkg, &outs)

		*ast1 += fmt.Sprintf("\t\t%d.- Function: %s (%s) (%s)\n",
			j, fn.Name, inps.String(), outs.String())

		for k, expr := range fn.Expressions {
			var inps bytes.Buffer
			var outs bytes.Buffer
			var opName1 string
			var lbl string

			// Adding label in case a `goto` statement was used for the expression.
			if expr.Label != "" {
				lbl = " <<" + expr.Label + ">>"
			} else {
				lbl = ""
			}

			// Determining operator's name.
			if expr.Operator != nil {
				if expr.Operator.IsBuiltin {

					opName1 = OpNames[expr.Operator.OpCode]
				} else {
					opName1 = expr.Operator.Name
				}
			}

			getFormattedParam(expr.Inputs, pkg, &inps)
			getFormattedParam(expr.Outputs, pkg, &outs)

			if expr.Operator != nil {
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
				if len(expr.Outputs) > 0 {
					out := expr.Outputs[len(expr.Outputs)-1]

					*ast1 += fmt.Sprintf("\t\t\t%d.- Declaration%s: %s %s\n",
						k,
						lbl,
						expr.Outputs[0].ArgDetails.Name,
						GetFormattedType(out))
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
	for _, pkg := range prgrm.Packages {
		if constants2.IsCorePackage(pkg.Name) {
			continue
		}

		*ast += fmt.Sprintf("%d.- Package: %s\n", i, pkg.Name)

		buildStrImports(pkg, ast)
		buildStrGlobals(pkg, ast)
		buildStrStructs(pkg, ast)
		buildStrFunctions(pkg, ast)

		i++
	}
}

// getFormattedParam is an auxiliary function for `ToString`. It formats the
// name of a `CXExpression`'s input and output parameters (`CXArgument`s). Examples
// of these formattings are "pkg.foo[0]", "&*foo.field1". The result is written to
// `buf`.
func getFormattedParam(params []*CXArgument, pkg *CXPackage, buf *bytes.Buffer) {
	for i, param := range params {
		elt := GetAssignmentElement(param)

		// Checking if this argument comes from an imported package.
		externalPkg := false
		if pkg != param.ArgDetails.Package {
			externalPkg = true
		}

		if i == len(params)-1 {
			buf.WriteString(fmt.Sprintf("%s %s", GetFormattedName(param, externalPkg), GetFormattedType(elt)))
		} else {
			buf.WriteString(fmt.Sprintf("%s %s, ", GetFormattedName(param, externalPkg), GetFormattedType(elt)))
		}
	}
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

func getNonCollectionValue(fp int, arg, elt *CXArgument, typ string) string {
	if arg.IsPointer {
		return fmt.Sprintf("%v", ReadPtr(fp, elt))
	}
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
				val += fmt.Sprintf("%s: %s", fld.ArgDetails.Name, GetPrintableValue(fp+arg.Offset+off, fld))
			} else {
				val += fmt.Sprintf("%s: %s, ", fld.ArgDetails.Name, GetPrintableValue(fp+arg.Offset+off, fld))
			}
			off += fld.TotalSize
		}
		val += "}"
		return val
	}
}

// GetPrintableValue ...
func GetPrintableValue(fp int, arg *CXArgument) string {
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
					for i := 0; i < closeCount; i++ {
						val += "]"
					}
					val += " "
					for i := 0; i < closeCount; i++ {
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

// IsPointer checks if `sym` is a candidate for the garbage collector to check.
// For example, if `sym` is a slice, the garbage collector will need to check
// if the slice on the heap needs to be relocated.
func IsPointer(sym *CXArgument) bool {
	// There's no need to add global variables in `fn.ListOfPointers` as we can access them easily through `CXPackage.Globals`
	// TODO: We could still pre-compute a list of candidates for globals.
	if sym.Offset >= PROGRAM.StackSize && sym.ArgDetails.Name != "" {
		return false
	}
	// NOTE: Strings are considered as `IsPointer`s by the runtime.
	// if (sym.IsPointer || sym.IsSlice) && sym.ArgDetails.Name != "" {
	// 	return true
	// }
	if (sym.IsPointer || sym.IsSlice) && sym.ArgDetails.Name != "" && len(sym.Fields) == 0 {
		return true
	}
	if sym.Type == constants.TYPE_STR && sym.ArgDetails.Name != "" && len(sym.Fields) == 0 {
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

// getFormattedDerefs is an auxiliary function for `GetFormattedName`. This
// function formats indexing and pointer dereferences associated to `arg`.
func getFormattedDerefs(arg *CXArgument, includePkg bool) string {
	name := ""
	// Checking if we should include `arg`'s package name.
	if includePkg {
		name = fmt.Sprintf("%s.%s", arg.ArgDetails.Package.Name, arg.ArgDetails.Name)
	} else {
		name = arg.ArgDetails.Name
	}

	// If it's a literal, just override the name with LITERAL_PLACEHOLDER.
	if arg.ArgDetails.Name == "" {
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
		if idx.Offset > PROGRAM.StackSize {
			// Then it's a literal.
			idxI32 := helper.Deserialize_i32(PROGRAM.Memory[idx.Offset : idx.Offset+constants.TYPE_POINTER_SIZE])
			idxValue = fmt.Sprintf("%d", idxI32)
		} else {
			// Then let's just print the variable name.
			idxValue = idx.ArgDetails.Name
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
	if arg.PassBy == constants.PASSBY_REFERENCE {
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
						pkg, err := PROGRAM.GetPackage(arg.ArgDetails.Package.Name)
						if err != nil {
							println(CompilationError(elt.ArgDetails.FileName, elt.ArgDetails.FileLine), err.Error())
							os.Exit(constants.CX_COMPILATION_ERROR)
						}

						fn, err := pkg.GetFunction(elt.ArgDetails.Name)
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
func SignatureStringOfStruct(s *CXStruct) string {
	fields := ""
	for _, f := range s.Fields {
		fields += fmt.Sprintf(" %s %s;", f.ArgDetails.Name, GetFormattedType(f))
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
