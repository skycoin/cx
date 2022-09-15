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
		for _, inputIdx := range opInputs {
			input := cxprogram.GetCXTypeSignatureFromArray(inputIdx)

			var inp *CXArgument = &CXArgument{ArgDetails: &CXArgumentDebug{}}
			if input.Type == TYPE_CXARGUMENT_DEPRECATE {
				inp = cxprogram.GetCXArgFromArray(CXArgumentIndex(input.Meta))
			} else if input.Type == TYPE_ATOMIC || input.Type == TYPE_POINTER_ATOMIC {
				inp = &CXArgument{ArgDetails: &CXArgumentDebug{}}
			} else {
				panic("type is not known")
			}

			fmt.Println("ProgramInput")
			fmt.Printf("\t%s : %s() : %s\n", stackValueHeader(inp.ArgDetails.FileName, inp.ArgDetails.FileLine), op.Name, GetPrintableValue(cxprogram, fp, input))

			inpPkg, err := cxprogram.GetPackageFromArray(input.Package)
			if err != nil {
				panic(err)
			}
			dupNames = append(dupNames, inpPkg.Name+input.Name)
		}

		opOutputs := op.GetOutputs(cxprogram)
		for _, outputIdx := range opOutputs {
			output := cxprogram.GetCXTypeSignatureFromArray(outputIdx)

			var out *CXArgument = &CXArgument{ArgDetails: &CXArgumentDebug{}}
			if output.Type == TYPE_CXARGUMENT_DEPRECATE {
				out = cxprogram.GetCXArgFromArray(CXArgumentIndex(output.Meta))
			} else if output.Type == TYPE_ATOMIC || output.Type == TYPE_POINTER_ATOMIC {
				out = &CXArgument{ArgDetails: &CXArgumentDebug{}}
			} else {
				panic("type is not known")
			}

			fmt.Println("ProgramOutput")
			fmt.Printf("\t%s : %s() : %s\n", stackValueHeader(out.ArgDetails.FileName, out.ArgDetails.FileLine), op.Name, GetPrintableValue(cxprogram, fp, output))

			outPkg, err := cxprogram.GetPackageFromArray(output.Package)
			if err != nil {
				panic(err)
			}
			dupNames = append(dupNames, outPkg.Name+output.Name)
		}

		// fmt.Println("Expressions")
		exprs := ""
		for _, expr := range op.Expressions {
			cxAtomicOp, err := cxprogram.GetCXAtomicOp(expr.Index)
			if err != nil {
				panic(err)
			}

			cxAtomicOpOperator := cxprogram.GetFunctionFromArray(cxAtomicOp.Operator)
			for _, inputIdx := range cxAtomicOp.GetInputs(cxprogram) {

				input := cxprogram.GetCXTypeSignatureFromArray(inputIdx)

				var inp *CXArgument = &CXArgument{ArgDetails: &CXArgumentDebug{}}
				if input.Type == TYPE_CXARGUMENT_DEPRECATE {
					inp = cxprogram.GetCXArgFromArray(CXArgumentIndex(input.Meta))
				} else if input.Type == TYPE_ATOMIC || input.Type == TYPE_POINTER_ATOMIC {
					// do nothing
				} else if input.Type == TYPE_ARRAY_ATOMIC {
					// do nothing
				}

				if input.Name == "" || cxAtomicOpOperator == nil {
					continue
				}

				inpPkg, err := cxprogram.GetPackageFromArray(input.Package)
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
				exprs += fmt.Sprintf("\t%s : %s() : %s\n", stackValueHeader(inp.ArgDetails.FileName, inp.ArgDetails.FileLine), cxAtomicOp.GetOperatorName(cxprogram), GetPrintableValue(cxprogram, fp, input))

				dupNames = append(dupNames, inpPkg.Name+inp.Name)
			}

			for _, outputIdx := range cxAtomicOp.GetOutputs(cxprogram) {
				output := cxprogram.GetCXTypeSignatureFromArray(outputIdx)

				var out *CXArgument = &CXArgument{ArgDetails: &CXArgumentDebug{}}
				if output.Type == TYPE_CXARGUMENT_DEPRECATE {
					out = cxprogram.GetCXArgFromArray(CXArgumentIndex(output.Meta))

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

					exprs += fmt.Sprintf("\t%s : %s() : %s\n", stackValueHeader(out.ArgDetails.FileName, out.ArgDetails.FileLine), cxAtomicOp.GetOperatorName(cxprogram), GetPrintableValue(cxprogram, fp, output))

					dupNames = append(dupNames, outPkg.Name+out.Name)
				} else if output.Type == TYPE_ATOMIC || output.Type == TYPE_POINTER_ATOMIC || output.Type == TYPE_ARRAY_ATOMIC {
					if output.Name == "" || cxAtomicOpOperator == nil {
						continue
					}

					outputPkg, err := cxprogram.GetPackageFromArray(output.Package)
					if err != nil {
						panic(err)
					}

					var dup bool
					for _, name := range dupNames {
						if name == outputPkg.Name+output.Name {
							dup = true
							break
						}
					}
					if dup {
						continue
					}

					// TODO: Make GetPrintableValue() and other functions receive CXTypeSignature instead of CXArgument
					exprs += fmt.Sprintf("\t%s : %s() : %s\n", stackValueHeader(out.ArgDetails.FileName, out.ArgDetails.FileLine), cxAtomicOp.GetOperatorName(cxprogram), GetPrintableValue(cxprogram, fp, output))

					dupNames = append(dupNames, outputPkg.Name+output.Name)

				} else {
					panic("type is not known")
				}

			}
		}

		if len(exprs) > 0 {
			fmt.Println("Expressions\n", exprs)
		}
	}
}
