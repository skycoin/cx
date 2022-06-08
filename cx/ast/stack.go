package ast

import (
	"fmt"

	"github.com/skycoin/cx/cx/types"
)

func stackValueHeader(fileName string, fileLine int) string {
	return fmt.Sprintf("%s:%d", fileName, fileLine)
}

// PrintStack ...
func (cxprogram *CXProgram) PrintStack() {
	fmt.Println()
	fmt.Println("===Callstack===")

	// we're going backwards in the stack
	fp := cxprogram.Stack.Pointer

	for c := cxprogram.CallCounter; c != types.InvalidPointer; c-- {
		op := cxprogram.CallStack[c].Operator
		fp -= op.Size

		var dupNames []string

		fmt.Printf(">>> %s()\n", op.Name)

		opInputs := op.GetInputs(cxprogram)
		for _, input := range opInputs {
			var inp *CXArgument
			if input.Type == TYPE_CXARGUMENT_DEPRECATE {
				inp = cxprogram.GetCXArgFromArray(CXArgumentIndex(input.Meta))
			}

			fmt.Println("ProgramInput")
			fmt.Printf("\t%s : %s() : %s\n", stackValueHeader(inp.ArgDetails.FileName, inp.ArgDetails.FileLine), op.Name, GetPrintableValue(cxprogram, fp, inp))

			inpPkg, err := cxprogram.GetPackageFromArray(inp.Package)
			if err != nil {
				panic(err)
			}
			dupNames = append(dupNames, inpPkg.Name+inp.Name)
		}

		opOutputs := op.GetOutputs(cxprogram)
		for _, output := range opOutputs {
			var out *CXArgument
			if output.Type == TYPE_CXARGUMENT_DEPRECATE {
				out = cxprogram.GetCXArgFromArray(CXArgumentIndex(output.Meta))
			}
			fmt.Println("ProgramOutput")
			fmt.Printf("\t%s : %s() : %s\n", stackValueHeader(out.ArgDetails.FileName, out.ArgDetails.FileLine), op.Name, GetPrintableValue(cxprogram, fp, out))

			outPkg, err := cxprogram.GetPackageFromArray(out.Package)
			if err != nil {
				panic(err)
			}
			dupNames = append(dupNames, outPkg.Name+out.Name)
		}

		// fmt.Println("Expressions")
		exprs := ""
		for _, expr := range op.Expressions {
			cxAtomicOp, err := cxprogram.GetCXAtomicOp(expr.Index)
			if err != nil {
				panic(err)
			}

			cxAtomicOpOperator := cxprogram.GetFunctionFromArray(cxAtomicOp.Operator)
			for _, input := range cxAtomicOp.GetInputs(cxprogram) {
				var inp *CXArgument
				if input.Type == TYPE_CXARGUMENT_DEPRECATE {
					inp = cxprogram.GetCXArgFromArray(CXArgumentIndex(input.Meta))
				}

				if inp.Name == "" || cxAtomicOpOperator == nil {
					continue
				}

				inpPkg, err := cxprogram.GetPackageFromArray(inp.Package)
				if err != nil {
					panic(err)
				}
				var dup bool
				for _, name := range dupNames {
					if name == inpPkg.Name+inp.Name {
						dup = true
						break
					}
				}
				if dup {
					continue
				}

				// fmt.Println("\t", inp.Name, "\t", ":", "\t", GetPrintableValue(fp, inp))
				// exprs += fmt.Sprintln("\t", stackValueHeader(inp.FileName, inp.FileLine), "\t", ":", "\t", GetPrintableValue(fp, inp))

				exprs += fmt.Sprintf("\t%s : %s() : %s\n", stackValueHeader(inp.ArgDetails.FileName, inp.ArgDetails.FileLine), cxAtomicOp.GetOperatorName(cxprogram), GetPrintableValue(cxprogram, fp, inp))

				dupNames = append(dupNames, inpPkg.Name+inp.Name)
			}

			for _, output := range cxAtomicOp.GetOutputs(cxprogram) {
				var out *CXArgument
				if output.Type == TYPE_CXARGUMENT_DEPRECATE {
					out = cxprogram.GetCXArgFromArray(CXArgumentIndex(output.Meta))
				}

				if out.Name == "" || cxAtomicOpOperator == nil {
					continue
				}

				outPkg, err := cxprogram.GetPackageFromArray(out.Package)
				if err != nil {
					panic(err)
				}
				var dup bool
				for _, name := range dupNames {
					if name == outPkg.Name+out.Name {
						dup = true
						break
					}
				}
				if dup {
					continue
				}

				// fmt.Println("\t", out.Name, "\t", ":", "\t", GetPrintableValue(fp, out))
				// exprs += fmt.Sprintln("\t", stackValueHeader(out.FileName, out.FileLine), ":", GetPrintableValue(fp, out))

				exprs += fmt.Sprintf("\t%s : %s() : %s\n", stackValueHeader(out.ArgDetails.FileName, out.ArgDetails.FileLine), cxAtomicOp.GetOperatorName(cxprogram), GetPrintableValue(cxprogram, fp, out))

				dupNames = append(dupNames, outPkg.Name+out.Name)
			}
		}

		if len(exprs) > 0 {
			fmt.Println("Expressions\n", exprs)
		}
	}
}
