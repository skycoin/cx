package ast

import (
	"fmt"
)

func stackValueHeader(fileName string, fileLine int) string {
	return fmt.Sprintf("%s:%d", fileName, fileLine)
}

// PrintStack ...
func (cxprogram *CXProgram) PrintStack() {
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
			fmt.Printf("\t%s : %s() : %s\n", stackValueHeader(inp.ArgDetails.FileName, inp.ArgDetails.FileLine), op.Name, GetPrintableValue(fp, inp))

			dupNames = append(dupNames, inp.ArgDetails.Package.Name+inp.ArgDetails.Name)
		}

		for _, out := range op.Outputs {
			fmt.Println("ProgramOutput")
			fmt.Printf("\t%s : %s() : %s\n", stackValueHeader(out.ArgDetails.FileName, out.ArgDetails.FileLine), op.Name, GetPrintableValue(fp, out))

			dupNames = append(dupNames, out.ArgDetails.Package.Name+out.ArgDetails.Name)
		}

		// fmt.Println("Expressions")
		exprs := ""
		for _, expr := range op.Expressions {
			for _, inp := range expr.Inputs {
				if inp.ArgDetails.Name == "" || expr.Operator == nil {
					continue
				}
				var dup bool
				for _, name := range dupNames {
					if name == inp.ArgDetails.Package.Name+inp.ArgDetails.Name {
						dup = true
						break
					}
				}
				if dup {
					continue
				}

				// fmt.Println("\t", inp.Name, "\t", ":", "\t", GetPrintableValue(fp, inp))
				// exprs += fmt.Sprintln("\t", stackValueHeader(inp.FileName, inp.FileLine), "\t", ":", "\t", GetPrintableValue(fp, inp))

				exprs += fmt.Sprintf("\t%s : %s() : %s\n", stackValueHeader(inp.ArgDetails.FileName, inp.ArgDetails.FileLine), ExprOpName(expr), GetPrintableValue(fp, inp))

				dupNames = append(dupNames, inp.ArgDetails.Package.Name+inp.ArgDetails.Name)
			}

			for _, out := range expr.Outputs {
				if out.ArgDetails.Name == "" || expr.Operator == nil {
					continue
				}
				var dup bool
				for _, name := range dupNames {
					if name == out.ArgDetails.Package.Name+out.ArgDetails.Name {
						dup = true
						break
					}
				}
				if dup {
					continue
				}

				// fmt.Println("\t", out.Name, "\t", ":", "\t", GetPrintableValue(fp, out))
				// exprs += fmt.Sprintln("\t", stackValueHeader(out.FileName, out.FileLine), ":", GetPrintableValue(fp, out))

				exprs += fmt.Sprintf("\t%s : %s() : %s\n", stackValueHeader(out.ArgDetails.FileName, out.ArgDetails.FileLine), ExprOpName(expr), GetPrintableValue(fp, out))

				dupNames = append(dupNames, out.ArgDetails.Package.Name+out.ArgDetails.Name)
			}
		}

		if len(exprs) > 0 {
			fmt.Println("Expressions\n", exprs)
		}
	}
}

// TODO: Deprecate
func ExprOpName(expr *CXExpression) string {
	if expr.Operator.IsBuiltin {
		return OpNames[expr.Operator.OpCode]
	}
	return expr.Operator.Name

}
