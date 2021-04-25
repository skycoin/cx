package actions

import (
	"fmt"
	"os"

	constants2 "github.com/skycoin/cx/cxparser/constants"

	"github.com/skycoin/cx/cx/ast"
	"github.com/skycoin/cx/cx/constants"
	globals2 "github.com/skycoin/cx/cx/globals"
	"github.com/skycoin/cx/cxparser/globals"
)

// DeclareGlobal creates a global variable in the current package.
//
// If `doesInitialize` is true, then `initializer` is used to initialize the
// new variable. This function is a wrapper around DeclareGlobalInPackage()
// which does the real work.
//
// FIXME: This function should be merged with DeclareGlobalInPackage.
//        Just use pkg=nil to indicate that CurrentPackage should be used.
//
func DeclareGlobal(declarator *ast.CXArgument, declarationSpecifiers *ast.CXArgument,
	initializer []*ast.CXExpression, doesInitialize bool) {
	pkg, err := AST.GetCurrentPackage()
	if err != nil {
		panic(err)
	}

	DeclareGlobalInPackage(pkg, declarator, declarationSpecifiers, initializer, doesInitialize)
}

// DeclareGlobalInPackage creates a global variable in a specified package
//
// If `doesInitialize` is true, then `initializer` is used to initialize the
// new variable.
//
func DeclareGlobalInPackage(pkg *ast.CXPackage,
	declarator *ast.CXArgument, declaration_specifiers *ast.CXArgument,
	initializer []*ast.CXExpression, doesInitialize bool) {
	declaration_specifiers.ArgDetails.Package = pkg

	// Treat the name a bit different whether it's defined already or not.
	if glbl, err := pkg.GetGlobal(declarator.ArgDetails.Name); err == nil {
		// The name is already defined.

		if glbl.Offset < 0 || glbl.Size == 0 || glbl.TotalSize == 0 {
			// then it was only added a reference to the symbol
			var offExpr []*ast.CXExpression
			if declaration_specifiers.IsSlice {
				offExpr = WritePrimary(declaration_specifiers.Type,
					make([]byte, declaration_specifiers.Size), true)
			} else {
				offExpr = WritePrimary(declaration_specifiers.Type,
					make([]byte, declaration_specifiers.TotalSize), true)
			}

			glbl.Offset = offExpr[0].Outputs[0].Offset
			glbl.PassBy = offExpr[0].Outputs[0].PassBy
			// glbl.Package = offExpr[0].ProgramOutput[0].Package
		}

		// Checking if something is supposed to be initialized
		// and if `initializer` actually contains something.
		// If `initializer` is `nil`, this means that an expression
		// equivalent to nil was used, such as `[]i32{}`.
		if doesInitialize && initializer != nil {
			// then we just re-assign offsets
			if initializer[len(initializer)-1].Operator == nil {
				// then it's a literal
				declaration_specifiers.ArgDetails.Name = glbl.ArgDetails.Name
				declaration_specifiers.Offset = glbl.Offset
				declaration_specifiers.PassBy = glbl.PassBy
				declaration_specifiers.ArgDetails.Package = glbl.ArgDetails.Package

				*glbl = *declaration_specifiers

				initializer[len(initializer)-1].AddInput(initializer[len(initializer)-1].Outputs[0])
				initializer[len(initializer)-1].Outputs = nil
				initializer[len(initializer)-1].AddOutput(glbl)
				initializer[len(initializer)-1].Operator = ast.Natives[constants.OP_IDENTITY]
				initializer[len(initializer)-1].Package = glbl.ArgDetails.Package

				//add intialization statements, to array
				globals.SysInitExprs = append(globals.SysInitExprs, initializer...)
			} else {
				// then it's an expression
				declaration_specifiers.ArgDetails.Name = glbl.ArgDetails.Name
				declaration_specifiers.Offset = glbl.Offset
				declaration_specifiers.PassBy = glbl.PassBy
				declaration_specifiers.ArgDetails.Package = glbl.ArgDetails.Package

				*glbl = *declaration_specifiers

				if initializer[len(initializer)-1].IsStructLiteral() {
					initializer = StructLiteralAssignment([]*ast.CXExpression{&ast.CXExpression{Outputs: []*ast.CXArgument{glbl}}}, initializer)
				} else {
					initializer[len(initializer)-1].Outputs = nil
					initializer[len(initializer)-1].AddOutput(glbl)
				}
				//add intialization statements, to array
				globals.SysInitExprs = append(globals.SysInitExprs, initializer...)
			}
		} else {
			// we keep the last value for now
			declaration_specifiers.ArgDetails.Name = glbl.ArgDetails.Name
			declaration_specifiers.Offset = glbl.Offset
			declaration_specifiers.PassBy = glbl.PassBy
			declaration_specifiers.ArgDetails.Package = glbl.ArgDetails.Package
			*glbl = *declaration_specifiers
		}
	} else {
		// then it hasn't been defined
		var offExpr []*ast.CXExpression
		if declaration_specifiers.IsSlice {
			offExpr = WritePrimary(declaration_specifiers.Type, make([]byte, declaration_specifiers.Size), true)
		} else {
			offExpr = WritePrimary(declaration_specifiers.Type, make([]byte, declaration_specifiers.TotalSize), true)
		}

		// Checking if something is supposed to be initialized
		// and if `initializer` actually contains something.
		// If `initializer` is `nil`, this means that an expression
		// equivalent to nil was used, such as `[]i32{}`.
		if doesInitialize && initializer != nil {
			if initializer[len(initializer)-1].Operator == nil {
				// then it's a literal

				declaration_specifiers.ArgDetails.Name = declarator.ArgDetails.Name
				declaration_specifiers.ArgDetails.FileLine = declarator.ArgDetails.FileLine
				declaration_specifiers.Offset = offExpr[0].Outputs[0].Offset
				declaration_specifiers.Size = offExpr[0].Outputs[0].Size
				declaration_specifiers.TotalSize = offExpr[0].Outputs[0].TotalSize
				declaration_specifiers.ArgDetails.Package = pkg

				initializer[len(initializer)-1].Operator = ast.Natives[constants.OP_IDENTITY]
				initializer[len(initializer)-1].AddInput(initializer[len(initializer)-1].Outputs[0])
				initializer[len(initializer)-1].Outputs = nil
				initializer[len(initializer)-1].AddOutput(declaration_specifiers)

				pkg.AddGlobal(declaration_specifiers)
				//add intialization statements, to array
				globals.SysInitExprs = append(globals.SysInitExprs, initializer...)
			} else {
				// then it's an expression
				declaration_specifiers.ArgDetails.Name = declarator.ArgDetails.Name
				declaration_specifiers.ArgDetails.FileLine = declarator.ArgDetails.FileLine
				declaration_specifiers.Offset = offExpr[0].Outputs[0].Offset
				declaration_specifiers.Size = offExpr[0].Outputs[0].Size
				declaration_specifiers.TotalSize = offExpr[0].Outputs[0].TotalSize
				declaration_specifiers.ArgDetails.Package = pkg

				if initializer[len(initializer)-1].IsStructLiteral() {
					initializer = StructLiteralAssignment([]*ast.CXExpression{&ast.CXExpression{Outputs: []*ast.CXArgument{declaration_specifiers}}}, initializer)
				} else {
					initializer[len(initializer)-1].Outputs = nil
					initializer[len(initializer)-1].AddOutput(declaration_specifiers)
				}

				pkg.AddGlobal(declaration_specifiers)
				//add intialization statements, to array
				globals.SysInitExprs = append(globals.SysInitExprs, initializer...)
			}
		} else {
			// offExpr := WritePrimary(declaration_specifiers.Type, make([]byte, declaration_specifiers.Size), true)
			// exprOut := expr[0].ProgramOutput[0]

			declaration_specifiers.ArgDetails.Name = declarator.ArgDetails.Name
			declaration_specifiers.ArgDetails.FileLine = declarator.ArgDetails.FileLine
			declaration_specifiers.Offset = offExpr[0].Outputs[0].Offset
			declaration_specifiers.Size = offExpr[0].Outputs[0].Size
			declaration_specifiers.TotalSize = offExpr[0].Outputs[0].TotalSize
			declaration_specifiers.ArgDetails.Package = pkg

			pkg.AddGlobal(declaration_specifiers)
		}
	}
}

// DeclareStruct takes a name of a struct and a slice of fields representing
// the members and adds the struct to the package.
//
func DeclareStruct(ident string, strctFlds []*ast.CXArgument) {
	// Make sure we are inside a package.
	pkg, err := AST.GetCurrentPackage()
	if err != nil {
		// FIXME: Should give a relevant error message
		panic(err)
	}

	// Make sure a struct with the same name is not yet defined.
	strct, err := AST.GetStruct(ident, pkg.Name)
	if err != nil {
		// FIXME: Should give a relevant error message
		panic(err)
	}

	strct.Fields = nil
	strct.Size = 0
	for _, fld := range strctFlds {
		if _, err := strct.GetField(fld.ArgDetails.Name); err == nil {
			println(ast.CompilationError(fld.ArgDetails.FileName, fld.ArgDetails.FileLine), "Multiply defined struct field:", fld.ArgDetails.Name)
		} else {
			strct.AddField(fld)
		}
	}
}

// DeclarePackage() switches the current package in the program.
//
func DeclarePackage(ident string) {
	// Add a new package to the program if it's not previously defined.
	if _, err := AST.GetPackage(ident); err != nil {
		pkg := ast.MakePackage(ident)
		AST.AddPackage(pkg)
	}

	AST.SelectPackage(ident)
}

// DeclareImport()
//
func DeclareImport(name string, currentFile string, lineNo int) {
	// Make sure we are inside a package
	pkg, err := AST.GetCurrentPackage()
	if err != nil {
		// FIXME: Should give a relevant error message
		panic(err)
	}

	// Checking if it's a package in the CX workspace by trying to find a
	// slash (/) in the name.
	// We start backwards and we stop if we find a slash.
	hasSlash := false
	c := len(name) - 1
	for ; c >= 0; c-- {
		if name[c] == '/' {
			hasSlash = true
			break
		}
	}
	ident := ""
	// If the `name` has a slash, then we need to strip
	// everything behind the slash and the slash itself.
	if hasSlash {
		ident = name[c+1:]
	} else {
		ident = name
	}

	// If the package is already imported, then there is nothing more to be done.
	if _, err := pkg.GetImport(ident); err == nil {
		return
	}

	// If the package is already defined in the program, just add it to
	// the importing package.
	if imp, err := AST.GetPackage(ident); err == nil {
		pkg.AddImport(imp)
		return
	}

	// All packages are read during the first pass of the compilation.  So
	// if we get here during the 2nd pass, it's either a core package or
	// something is panic-level wrong.
	if constants2.IsCorePackage(ident) {
		imp := ast.MakePackage(ident)
		pkg.AddImport(imp)
		AST.AddPackage(imp)
		AST.CurrentPackage = pkg

		if ident == "aff" {
			AffordanceStructs(imp, currentFile, lineNo)
		}
	} else {
		// This should never happen.
		println(ast.CompilationError(currentFile, lineNo), fmt.Sprintf("unkown error when trying to read package '%s'", ident))
		os.Exit(constants.CX_COMPILATION_ERROR)
	}
}

// DeclareLocal() creates a local variable inside a function.
// If `doesInitialize` is true, then `initializer` contains the initial values
// of the variable(s).
//
// Returns a list of expressions that contains the initialization, if any.
//
func DeclareLocal(declarator *ast.CXArgument, declarationSpecifiers *ast.CXArgument,
	initializer []*ast.CXExpression, doesInitialize bool) []*ast.CXExpression {
	if globals2.FoundCompileErrors {
		return nil
	}

	declarationSpecifiers.IsLocalDeclaration = true

	pkg, err := AST.GetCurrentPackage()
	if err != nil {
		panic(err)
	}

	// Declaration expression to handle the inline initialization.
	// For example, `var foo i32 = 11` needs to be divided into two expressions:
	// one that declares `foo`, and another that assigns 11 to `foo`
	decl := ast.MakeExpression(nil, declarator.ArgDetails.FileName, declarator.ArgDetails.FileLine)
	decl.Package = pkg

	declarationSpecifiers.ArgDetails.Name = declarator.ArgDetails.Name
	declarationSpecifiers.ArgDetails.FileLine = declarator.ArgDetails.FileLine
	declarationSpecifiers.ArgDetails.Package = pkg
	declarationSpecifiers.PreviouslyDeclared = true
	decl.AddOutput(declarationSpecifiers)

	// Checking if something is supposed to be initialized
	// and if `initializer` actually contains something.
	// If `initializer` is `nil`, this means that an expression
	// equivalent to nil was used, such as `[]i32{}`.
	if doesInitialize && initializer != nil {
		// THEN it's a literal, e.g. var foo i32 = 10;
		// ELSE it's an expression with an operator
		if initializer[len(initializer)-1].Operator == nil {
			// we need to create an expression that links the initializer expressions
			// with the declared variable
			expr := ast.MakeExpression(ast.Natives[constants.OP_IDENTITY], CurrentFile, LineNo)
			expr.Package = pkg

			initOut := initializer[len(initializer)-1].Outputs[0]
			// CX checks the output of an expression to determine if it's being passed
			// by value or by reference, so we copy this property from the initializer's
			// output, in case of something like var foo *i32 = &bar
			declarationSpecifiers.PassBy = initOut.PassBy

			expr.AddOutput(declarationSpecifiers)
			expr.AddInput(initOut)

			initializer[len(initializer)-1] = expr

			return append([]*ast.CXExpression{decl}, initializer...)
		} else {
			expr := initializer[len(initializer)-1]

			// THEN the expression has outputs created from the result of
			// handling a dot notation initializer, and it needs to be replaced
			// ELSE we simply add it using `AddOutput`
			if len(expr.Outputs) > 0 {
				expr.Outputs = []*ast.CXArgument{declarationSpecifiers}
			} else {
				expr.AddOutput(declarationSpecifiers)
			}

			return append([]*ast.CXExpression{decl}, initializer...)
		}
	} else {
		// There is no initialization.
		expr := ast.MakeExpression(nil, declarator.ArgDetails.FileName, declarator.ArgDetails.FileLine)
		expr.Package = pkg

		declarationSpecifiers.ArgDetails.Name = declarator.ArgDetails.Name
		declarationSpecifiers.ArgDetails.FileLine = declarator.ArgDetails.FileLine
		declarationSpecifiers.ArgDetails.Package = pkg
		declarationSpecifiers.PreviouslyDeclared = true
		expr.AddOutput(declarationSpecifiers)

		return []*ast.CXExpression{expr}
	}
}

// DeclarationSpecifiers is called to build a type of a variable or parameter.
//
// It is called repeatedly while the type is parsed.
//
//   declSpec:     The incoming type
//   arrayLengths: The lengths of the array if `opTyp` = cxcore.DECL_ARRAY
//   opTyp:        The type of modification to `declSpec` (array of, pointer to, ...)
//
// Returns the new type build from `declSpec` and `opTyp`.
//
func DeclarationSpecifiers(declSpec *ast.CXArgument, arrayLengths []int, opTyp int) *ast.CXArgument {
	switch opTyp {
	case constants.DECL_POINTER:
		declSpec.DeclarationSpecifiers = append(declSpec.DeclarationSpecifiers, constants.DECL_POINTER)
		// if !declSpec.IsPointer {
		declSpec.IsPointer = true
		declSpec.Size = constants.TYPE_POINTER_SIZE
		declSpec.TotalSize = constants.TYPE_POINTER_SIZE
		declSpec.IndirectionLevels++
		// }
		// else {
		// pointer := declSpec

		// for c := declSpec.IndirectionLevels; c > 1; c-- {
		// 	pointer.IndirectionLevels = c
		// 	pointer.IsPointer = true
		// }

		// declSpec.IndirectionLevels++

		// pointer.Size = constants.TYPE_POINTER_SIZE
		// pointer.TotalSize = constants.TYPE_POINTER_SIZE
		// }

		return declSpec
	case constants.DECL_ARRAY:
		for range arrayLengths {
			declSpec.DeclarationSpecifiers = append(declSpec.DeclarationSpecifiers, constants.DECL_ARRAY)
		}
		arg := declSpec
		// arg.IsArray = true
		arg.Lengths = arrayLengths
		arg.TotalSize = arg.Size * TotalLength(arg.Lengths)

		return arg
	case constants.DECL_SLICE:
		// for range arrayLengths {
		// 	declSpec.DeclarationSpecifiers = append(declSpec.DeclarationSpecifiers, cxcore.DECL_SLICE)
		// }

		arg := declSpec

		arg.DeclarationSpecifiers = append(arg.DeclarationSpecifiers, constants.DECL_SLICE)

		arg.IsSlice = true
		arg.IsReference = true
		// arg.IsArray = true
		arg.PassBy = constants.PASSBY_REFERENCE

		arg.Lengths = append([]int{0}, arg.Lengths...)
		// arg.Lengths = arrayLengths
		// arg.TotalSize = arg.Size
		// arg.Size = cxcore.TYPE_POINTER_SIZE
		arg.TotalSize = constants.TYPE_POINTER_SIZE

		return arg
	case constants.DECL_BASIC:
		arg := declSpec
		// arg.DeclarationSpecifiers = append(arg.DeclarationSpecifiers, cxcore.DECL_BASIC)
		arg.TotalSize = arg.Size
		return arg
	case constants.DECL_FUNC:
		// Creating this case if additional operations are needed in the
		// future.
		return declSpec
	}

	return nil
}

// DeclarationSpecifiersBasic() returns a type specifier created from one of the builtin types.
//
func DeclarationSpecifiersBasic(typ int) *ast.CXArgument {
	arg := ast.MakeArgument("", CurrentFile, LineNo)
	arg.AddType(constants.TypeNames[typ])
	arg.Type = typ

	arg.Size = constants.GetArgSize(typ)

	if typ == constants.TYPE_AFF {
		// equivalent to slice of strings
		return DeclarationSpecifiers(arg, []int{0}, constants.DECL_SLICE)
	}

	return DeclarationSpecifiers(arg, []int{0}, constants.DECL_BASIC)
}

// DeclarationSpecifiersStruct() declares a struct
func DeclarationSpecifiersStruct(ident string, pkgName string,
	isExternal bool, currentFile string, lineNo int) *ast.CXArgument {
	pkg, err := AST.GetCurrentPackage()
	if err != nil {
		panic(err)
	}

	if isExternal {
		// custom type in an imported package
		imp, err := pkg.GetImport(pkgName)
		if err != nil {
			panic(err)
		}

		strct, err := AST.GetStruct(ident, imp.Name)
		if err != nil {
			println(ast.CompilationError(currentFile, lineNo), err.Error())
			return nil
		}

		arg := ast.MakeArgument("", currentFile, lineNo)
		arg.Type = constants.TYPE_CUSTOM
		arg.CustomType = strct
		arg.Size = strct.Size
		arg.TotalSize = strct.Size

		arg.ArgDetails.Package = pkg
		arg.DeclarationSpecifiers = append(arg.DeclarationSpecifiers, constants.DECL_STRUCT)

		return arg
	} else {
		// custom type in the current package
		strct, err := AST.GetStruct(ident, pkg.Name)
		if err != nil {
			println(ast.CompilationError(currentFile, lineNo), err.Error())
			return nil
		}

		arg := ast.MakeArgument("", currentFile, lineNo)
		arg.Type = constants.TYPE_CUSTOM
		arg.DeclarationSpecifiers = append(arg.DeclarationSpecifiers, constants.DECL_STRUCT)
		arg.CustomType = strct
		arg.Size = strct.Size
		arg.TotalSize = strct.Size
		arg.ArgDetails.Package = pkg

		return arg
	}
}
