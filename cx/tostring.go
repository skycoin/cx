package cxcore

import (
	"bytes"
	"fmt"
	"github.com/skycoin/cx/cx/ast"
	"github.com/skycoin/cx/cx/constants"
)

// buildStrImports is an auxiliary function for `toString`. It builds
// string representation all the imported packages of `pkg`.
func buildStrImports(pkg *ast.CXPackage, ast *string) {
	if len(pkg.Imports) > 0 {
		*ast += "\tImports\n"
	}

	for j, imp := range pkg.Imports {
		*ast += fmt.Sprintf("\t\t%d.- Import: %s\n", j, imp.Name)
	}
}

// buildStrGlobals is an auxiliary function for `toString`. It builds
// string representation of all the global variables of `pkg`.
func buildStrGlobals(pkg *ast.CXPackage, ast *string) {
	if len(pkg.Globals) > 0 {
		*ast += "\tGlobals\n"
	}

	for j, v := range pkg.Globals {
		*ast += fmt.Sprintf("\t\t%d.- Global: %s %s\n", j, v.Name, GetFormattedType(v))
	}
}

// buildStrStructs is an auxiliary function for `toString`. It builds
// string representation of all the structures defined in `pkg`.
func buildStrStructs(pkg *ast.CXPackage, ast *string) {
	if len(pkg.Structs) > 0 {
		*ast += "\tStructs\n"
	}

	for j, strct := range pkg.Structs {
		*ast += fmt.Sprintf("\t\t%d.- Struct: %s\n", j, strct.Name)

		for k, fld := range strct.Fields {
			*ast += fmt.Sprintf("\t\t\t%d.- Field: %s %s\n",
				k, fld.Name, GetFormattedType(fld))
		}
	}
}

// buildStrFunctions is an auxiliary function for `toString`. It builds
// string representation of all the functions defined in `pkg`.
func buildStrFunctions(pkg *ast.CXPackage, ast *string) {
	if len(pkg.Functions) > 0 {
		*ast += "\tFunctions\n"
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

		*ast += fmt.Sprintf("\t\t%d.- Function: %s (%s) (%s)\n",
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
				if expr.Operator.IsAtomic {
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
				*ast += fmt.Sprintf("\t\t\t%d.- Expression%s: %s%s%s(%s)\n",
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

					*ast += fmt.Sprintf("\t\t\t%d.- Declaration%s: %s %s\n",
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

// BuildStrPackages is an auxiliary function for `ToString`. It starts the
// process of building string format of the abstract syntax tree of a CX program.
func BuildStrPackages(prgrm *ast.CXProgram, ast *string) {
	// We need to declare the counter outside so we can
	// ignore the increments from core or stdlib packages.
	var i int
	for _, pkg := range prgrm.Packages {
		if IsCorePackage(pkg.Name) {
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
func getFormattedParam(params []*ast.CXArgument, pkg *ast.CXPackage, buf *bytes.Buffer) {
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

