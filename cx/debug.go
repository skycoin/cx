package cxcore

import (
	"bytes"
	"fmt"
	"os"
	"text/tabwriter"
)

// Debug is just a wrapper for `fmt.Println`. Its purpose is to bring
// the developer a quick way to `grep` all the instances of `Debug` in
// the source files to delete them after debugging.
func Debug(args ...interface{}) {
	fmt.Println(args...)
}

// PrintStack prints the stack trace of a CX program.
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

// DebugHeap prints the symbols that are acting as pointers in a CX program at certain point during the execution of the program along with the addresses they are pointing. Additionally, a list of the objects in the heap is printed, which shows their address in the heap, if they are marked as alive or as dead by the garbage collector, the address where they used to live after a garbage collector call, the full size of the object, the object itself as a slice of bytes and the pointers that are pointing to that object.
func DebugHeap() {
	// symsToAddrs will hold a list of symbols that are pointing to an address.
	symsToAddrs := make(map[int32][]string)

	// Processing global variables. Adding the address they are pointing to.
	for _, pkg := range PROGRAM.Packages {
		for _, glbl := range pkg.Globals {
			if glbl.IsPointer || glbl.IsSlice {
				heapOffset := Deserialize_i32(PROGRAM.Memory[glbl.Offset : glbl.Offset+TYPE_POINTER_SIZE])

				symsToAddrs[heapOffset] = append(symsToAddrs[heapOffset], glbl.Name)
			}
		}
	}

	// Processing local variables in every active function call in the `CallStack`.
	// Adding the address they are pointing to.
	var fp int
	for c := 0; c <= PROGRAM.CallCounter; c++ {
		op := PROGRAM.CallStack[c].Operator

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

			if ptr.Offset < PROGRAM.StackSize {
				offset += fp
			}

			heapOffset := Deserialize_i32(PROGRAM.Memory[offset : offset+TYPE_POINTER_SIZE])

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

	for c := PROGRAM.HeapStartsAt + NULL_HEAP_ADDRESS_OFFSET; c < PROGRAM.HeapStartsAt+PROGRAM.HeapPointer; {
		objSize := Deserialize_i32(PROGRAM.Memory[c+MARK_SIZE+FORWARDING_ADDRESS_SIZE : c+MARK_SIZE+FORWARDING_ADDRESS_SIZE+OBJECT_SIZE])

		// Setting a limit size for the object to be printed if the object is too large.
		// We don't want to print obscenely large objects to standard output.
		printObjSize := objSize
		if objSize > 50 {
			printObjSize = 50
		}

		var addrB [4]byte
		WriteMemI32(addrB[:], 0, int32(c))

		fmt.Fprintln(w, "Addr:\t", addrB, "\tMark:\t", PROGRAM.Memory[c:c+MARK_SIZE], "\tFwd:\t", PROGRAM.Memory[c+MARK_SIZE:c+MARK_SIZE+FORWARDING_ADDRESS_SIZE], "\tSize:\t", objSize, "\tObj:\t", PROGRAM.Memory[c+OBJECT_HEADER_SIZE:c+int(printObjSize)], "\tPtrs:", symsToAddrs[int32(c)])

		c += int(objSize)
	}

	// Just a newline.
	fmt.Fprintln(w)
	w.Flush()
}

// PrintProgram prints the abstract syntax tree of a CX program in a
// human-readable format.
func (prgrm *CXProgram) PrintProgram() {
	fmt.Println("Program")

	var currentFunction *CXFunction
	var currentPackage *CXPackage

	// Saving current program state because PrintProgram uses SelectXXX.
	// If we don't do this, calling `:dp` in a REPL will always switch the
	// user to the last function in the last package in the `CXProgram`
	// structure.
	if pkg, err := prgrm.GetCurrentPackage(); err == nil {
		currentPackage = pkg
	}

	if fn, err := prgrm.GetCurrentFunction(); err == nil {
		currentFunction = fn
	}

	printPackages(prgrm)

	// Restoring a program's state (what package and function were
	// selected.)
	if currentPackage != nil {
		_, err := prgrm.SelectPackage(currentPackage.Name)
		if err != nil {
			panic(err)
		}
	}
	if currentFunction != nil {
		_, err := prgrm.SelectFunction(currentFunction.Name)
		if err != nil {
			panic(err)
		}
	}

	prgrm.CurrentPackage = currentPackage
	if currentPackage != nil {
		currentPackage.CurrentFunction = currentFunction
	}
}

// printPackages is an auxiliary function for `PrintProgram`. It starts the
// process of printing the abstract syntax tree of a CX program.
func printPackages(prgrm *CXProgram) {
	// We need to declare the counter outside so we can
	// ignore the increments from core or stdlib packages.
	var i int
	for _, pkg := range prgrm.Packages {
		if IsCorePackage(pkg.Name) {
			continue
		}

		fmt.Printf("%d.- Package: %s\n", i, pkg.Name)

		printImports(pkg)
		printGlobals(pkg)
		printStructs(pkg)
		printFunctions(pkg)

		i++
	}
}

// printFunctions is an auxiliary function for `printProgram`. It prints all the
// functions defined in `pkg`.
func printFunctions(pkg *CXPackage) {
	if len(pkg.Functions) > 0 {
		fmt.Println("\tFunctions")
	}

	// We need to declare the counter outside so we can
	// ignore the increment from the `*init` function.
	var j int
	for _, fn := range pkg.Functions {
		if fn.Name == SYS_INIT_FUNC {
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

		fmt.Printf("\t\t%d.- Function: %s (%s) (%s)\n",
			j, fn.Name, inps.String(), outs.String())

		for k, expr := range fn.Expressions {
			var inps bytes.Buffer
			var outs bytes.Buffer
			var opName string
			var lbl string

			// Adding label in case a `goto` statement was used for the expression.
			if expr.Label != "" {
				lbl = " <<" + expr.Label + ">>"
			} else {
				lbl = ""
			}

			// Determining operator's name.
			if expr.Operator != nil {
				if expr.Operator.IsNative {
					opName = OpNames[expr.Operator.OpCode]
				} else {
					opName = expr.Operator.Name
				}
			}

			getFormattedParam(expr.Inputs, pkg, &inps)
			getFormattedParam(expr.Outputs, pkg, &outs)

			if expr.Operator != nil {
				assignOp := ""
				if outs.Len() > 0 {
					assignOp = " = "
				}
				fmt.Printf("\t\t\t%d.- Expression%s: %s%s%s(%s)\n",
					k,
					lbl,
					outs.String(),
					assignOp,
					opName,
					inps.String(),
				)
			} else {
				// Then it's a variable declaration. These are represented
				// by expressions without operators that only have outputs.
				if len(expr.Outputs) > 0 {
					out := expr.Outputs[len(expr.Outputs)-1]

					fmt.Printf("\t\t\t%d.- Declaration%s: %s %s\n",
						k,
						lbl,
						expr.Outputs[0].Name,
						GetFormattedType(out))
				}
			}
		}

		j++
	}
}

// printStructs is an auxiliary function for `printProgram`. It prints all the
// structures defined in `pkg`.
func printStructs(pkg *CXPackage) {
	if len(pkg.Structs) > 0 {
		fmt.Println("\tStructs")
	}

	for j, strct := range pkg.Structs {
		fmt.Printf("\t\t%d.- Struct: %s\n", j, strct.Name)

		for k, fld := range strct.Fields {
			fmt.Printf("\t\t\t%d.- Field: %s %s\n",
				k, fld.Name, GetFormattedType(fld))
		}
	}
}

// getFormattedParam is an auxiliary function for `PrintProgram`. It formats the
// name of a `CXExpression`'s input and output parameters (`CXArgument`s). Examples
// of these formattings are "pkg.foo[0]", "&*foo.field1". The result is written to
// `buf`.
func getFormattedParam(params []*CXArgument, pkg *CXPackage, buf *bytes.Buffer) {
	for i, param := range params {
		elt := GetAssignmentElement(param)

		// Checking if this argument comes from an imported package.
		externalPkg := false
		if pkg != param.Package {
			externalPkg = true
		}

		if i == len(params)-1 {
			buf.WriteString(fmt.Sprintf("%s %s", GetFormattedName(param, externalPkg), GetFormattedType(elt)))
		} else {
			buf.WriteString(fmt.Sprintf("%s %s, ", GetFormattedName(param, externalPkg), GetFormattedType(elt)))
		}
	}
}

// printImports is an auxiliary function for `printProgram`. It prints all the
// imported packages of `pkg`.
func printImports(pkg *CXPackage) {
	if len(pkg.Imports) > 0 {
		fmt.Println("\tImports")
	}

	for j, imp := range pkg.Imports {
		fmt.Printf("\t\t%d.- Import: %s\n", j, imp.Name)
	}
}

// printGlobals is an auxiliary function for `printProgram`. It prints all the
// global variables of `pkg`.
func printGlobals(pkg *CXPackage) {
	if len(pkg.Globals) > 0 {
		fmt.Println("\tGlobals")
	}

	for j, v := range pkg.Globals {
		fmt.Printf("\t\t%d.- Global: %s %s\n", j, v.Name, GetFormattedType(v))
	}
}

// GetPrintableValue tries to print the value (in human readable format) represented by `arg`.
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
