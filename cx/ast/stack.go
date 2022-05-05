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
		for _, inpIdx := range opInputs {
			inp := cxprogram.GetCXArgFromArray(inpIdx)
			fmt.Println("ProgramInput")
			fmt.Printf("\t%s : %s() : %s\n", stackValueHeader(inp.ArgDetails.FileName, inp.ArgDetails.FileLine), op.Name, GetPrintableValue(cxprogram, fp, inp))

			inpPkg, err := cxprogram.GetPackageFromArray(inp.Package)
			if err != nil {
				panic(err)
			}
			dupNames = append(dupNames, inpPkg.Name+inp.Name)
		}

		for _, outIdx := range op.Outputs {
			out := cxprogram.GetCXArgFromArray(outIdx)
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
			for _, inpIdx := range cxAtomicOp.Inputs {
				inp := cxprogram.GetCXArgFromArray(inpIdx)
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

			for _, outIdx := range cxAtomicOp.Outputs {
				out := cxprogram.GetCXArgFromArray(outIdx)
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
