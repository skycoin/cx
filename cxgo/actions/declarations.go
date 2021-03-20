package actions

import (
	"fmt"
	"os"

	"github.com/skycoin/cx/cx"
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
func DeclareGlobal(declarator *cxcore.CXArgument, declarationSpecifiers *cxcore.CXArgument,
	initializer []*cxcore.CXExpression, doesInitialize bool) {
	pkg, err := PRGRM.GetCurrentPackage()
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
func DeclareGlobalInPackage(pkg *cxcore.CXPackage,
	declarator *cxcore.CXArgument, declaration_specifiers *cxcore.CXArgument,
	initializer []*cxcore.CXExpression, doesInitialize bool) {
	declaration_specifiers.Package = pkg

	// Treat the name a bit different whether it's defined already or not.
	if glbl, err := pkg.GetGlobal(declarator.Name); err == nil {
		// The name is already defined.

		if glbl.Offset < 0 || glbl.Size == 0 || glbl.TotalSize == 0 {
			// then it was only added a reference to the symbol
			var offExpr []*cxcore.CXExpression
			if declaration_specifiers.IsSlice {
				offExpr = WritePrimary(declaration_specifiers.Type,
					make([]byte, declaration_specifiers.Size), true)
			} else {
				offExpr = WritePrimary(declaration_specifiers.Type,
					make([]byte, declaration_specifiers.TotalSize), true)
			}

			glbl.Offset = offExpr[0].Outputs[0].Offset
			glbl.PassBy = offExpr[0].Outputs[0].PassBy
			// glbl.Package = offExpr[0].Outputs[0].Package
		}

		// Checking if something is supposed to be initialized
		// and if `initializer` actually contains something.
		// If `initializer` is `nil`, this means that an expression
		// equivalent to nil was used, such as `[]i32{}`.
		if doesInitialize && initializer != nil {
			// then we just re-assign offsets
			if initializer[len(initializer)-1].Operator == nil {
				// then it's a literal
				declaration_specifiers.Name = glbl.Name
				declaration_specifiers.Offset = glbl.Offset
				declaration_specifiers.PassBy = glbl.PassBy
				declaration_specifiers.Package = glbl.Package

				*glbl = *declaration_specifiers

				initializer[len(initializer)-1].AddInput(initializer[len(initializer)-1].Outputs[0])
				initializer[len(initializer)-1].Outputs = nil
				initializer[len(initializer)-1].AddOutput(glbl)
				initializer[len(initializer)-1].Operator = cxcore.Natives[cxcore.OP_IDENTITY]
				initializer[len(initializer)-1].Package = glbl.Package

				SysInitExprs = append(SysInitExprs, initializer...)
			} else {
				// then it's an expression
				declaration_specifiers.Name = glbl.Name
				declaration_specifiers.Offset = glbl.Offset
				declaration_specifiers.PassBy = glbl.PassBy
				declaration_specifiers.Package = glbl.Package

				*glbl = *declaration_specifiers

				if initializer[len(initializer)-1].IsStructLiteral {
					initializer = StructLiteralAssignment([]*cxcore.CXExpression{&cxcore.CXExpression{Outputs: []*cxcore.CXArgument{glbl}}}, initializer)
				} else {
					initializer[len(initializer)-1].Outputs = nil
					initializer[len(initializer)-1].AddOutput(glbl)
				}

				SysInitExprs = append(SysInitExprs, initializer...)
			}
		} else {
			// we keep the last value for now
			declaration_specifiers.Name = glbl.Name
			declaration_specifiers.Offset = glbl.Offset
			declaration_specifiers.PassBy = glbl.PassBy
			declaration_specifiers.Package = glbl.Package
			*glbl = *declaration_specifiers
		}
	} else {
		// then it hasn't been defined
		var offExpr []*cxcore.CXExpression
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

				declaration_specifiers.Name = declarator.Name
				declaration_specifiers.FileLine = declarator.FileLine
				declaration_specifiers.Offset = offExpr[0].Outputs[0].Offset
				declaration_specifiers.Size = offExpr[0].Outputs[0].Size
				declaration_specifiers.TotalSize = offExpr[0].Outputs[0].TotalSize
				declaration_specifiers.Package = pkg

				initializer[len(initializer)-1].Operator = cxcore.Natives[cxcore.OP_IDENTITY]
				initializer[len(initializer)-1].AddInput(initializer[len(initializer)-1].Outputs[0])
				initializer[len(initializer)-1].Outputs = nil
				initializer[len(initializer)-1].AddOutput(declaration_specifiers)

				pkg.AddGlobal(declaration_specifiers)

				SysInitExprs = append(SysInitExprs, initializer...)
			} else {
				// then it's an expression
				declaration_specifiers.Name = declarator.Name
				declaration_specifiers.FileLine = declarator.FileLine
				declaration_specifiers.Offset = offExpr[0].Outputs[0].Offset
				declaration_specifiers.Size = offExpr[0].Outputs[0].Size
				declaration_specifiers.TotalSize = offExpr[0].Outputs[0].TotalSize
				declaration_specifiers.Package = pkg

				if initializer[len(initializer)-1].IsStructLiteral {
					initializer = StructLiteralAssignment([]*cxcore.CXExpression{&cxcore.CXExpression{Outputs: []*cxcore.CXArgument{declaration_specifiers}}}, initializer)
				} else {
					initializer[len(initializer)-1].Outputs = nil
					initializer[len(initializer)-1].AddOutput(declaration_specifiers)
				}

				pkg.AddGlobal(declaration_specifiers)
				SysInitExprs = append(SysInitExprs, initializer...)
			}
		} else {
			// offExpr := WritePrimary(declaration_specifiers.Type, make([]byte, declaration_specifiers.Size), true)
			// exprOut := expr[0].Outputs[0]

			declaration_specifiers.Name = declarator.Name
			declaration_specifiers.FileLine = declarator.FileLine
			declaration_specifiers.Offset = offExpr[0].Outputs[0].Offset
			declaration_specifiers.Size = offExpr[0].Outputs[0].Size
			declaration_specifiers.TotalSize = offExpr[0].Outputs[0].TotalSize
			declaration_specifiers.Package = pkg

			pkg.AddGlobal(declaration_specifiers)
		}
	}
}

// DeclareStruct takes a name of a struct and a slice of fields representing
// the members and adds the struct to the package.
//
func DeclareStruct(ident string, strctFlds []*cxcore.CXArgument) {
	// Make sure we are inside a package.
	pkg, err := PRGRM.GetCurrentPackage()
	if err != nil {
		// FIXME: Should give a relevant error message
		panic(err)
	}

	// Make sure a struct with the same name is not yet defined.
	strct, err := PRGRM.GetStruct(ident, pkg.Name)
	if err != nil {
		// FIXME: Should give a relevant error message
		panic(err)
	}

	strct.Fields = nil
	strct.Size = 0
	for _, fld := range strctFlds {
		if _, err := strct.GetField(fld.Name); err == nil {
			println(cxcore.CompilationError(fld.FileName, fld.FileLine), "Multiply defined struct field:", fld.Name)
		} else {
			strct.AddField(fld)
		}
	}
}

// DeclarePackage() switches the current package in the program.
//
func DeclarePackage(ident string) {
	// Add a new package to the program if it's not previously defined.
	if pkg, err := PRGRM.GetPackage(ident); err != nil {
		pkg = cxcore.MakePackage(ident)
		PRGRM.AddPackage(pkg)
	}

	PRGRM.SelectPackage(ident)
}

// DeclareImport()
//
func DeclareImport(name string, currentFile string, lineNo int) {
	// Make sure we are inside a package
	pkg, err := PRGRM.GetCurrentPackage()
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
	if imp, err := PRGRM.GetPackage(ident); err == nil {
		pkg.AddImport(imp)
		return
	}

	// All packages are read during the first pass of the compilation.  So
	// if we get here during the 2nd pass, it's either a core package or
	// something is panic-level wrong.
	if cxcore.IsCorePackage(ident) {
		imp := cxcore.MakePackage(ident)
		pkg.AddImport(imp)
		PRGRM.AddPackage(imp)
		PRGRM.CurrentPackage = pkg

		if ident == "aff" {
			AffordanceStructs(imp, currentFile, lineNo)
		}
	} else {
		// This should never happen.
		println(cxcore.CompilationError(currentFile, lineNo), fmt.Sprintf("unkown error when trying to read package '%s'", ident))
		os.Exit(cxcore.CX_COMPILATION_ERROR)
	}
}

// DeclareLocal() creates a local variable inside a function.
// If `doesInitialize` is true, then `initializer` contains the initial values
// of the variable(s).
//
// Returns a list of expressions that contains the initialization, if any.
//
func DeclareLocal(declarator *cxcore.CXArgument, declarationSpecifiers *cxcore.CXArgument,
	initializer []*cxcore.CXExpression, doesInitialize bool) []*cxcore.CXExpression {
	if cxcore.FoundCompileErrors {
		return nil
	}

	declarationSpecifiers.IsLocalDeclaration = true

	pkg, err := PRGRM.GetCurrentPackage()
	if err != nil {
		panic(err)
	}

	// Declaration expression to handle the inline initialization.
	// For example, `var foo i32 = 11` needs to be divided into two expressions:
	// one that declares `foo`, and another that assigns 11 to `foo`
	decl := cxcore.MakeExpression(nil, declarator.FileName, declarator.FileLine)
	decl.Package = pkg

	declarationSpecifiers.Name = declarator.Name
	declarationSpecifiers.FileLine = declarator.FileLine
	declarationSpecifiers.Package = pkg
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
			expr := cxcore.MakeExpression(cxcore.Natives[cxcore.OP_IDENTITY], CurrentFile, LineNo)
			expr.Package = pkg

			initOut := initializer[len(initializer)-1].Outputs[0]
			// CX checks the output of an expression to determine if it's being passed
			// by value or by reference, so we copy this property from the initializer's
			// output, in case of something like var foo *i32 = &bar
			declarationSpecifiers.PassBy = initOut.PassBy

			expr.AddOutput(declarationSpecifiers)
			expr.AddInput(initOut)

			initializer[len(initializer)-1] = expr

			return append([]*cxcore.CXExpression{decl}, initializer...)
		} else {
			expr := initializer[len(initializer)-1]

			// THEN the expression has outputs created from the result of
			// handling a dot notation initializer, and it needs to be replaced
			// ELSE we simply add it using `AddOutput`
			if len(expr.Outputs) > 0 {
				expr.Outputs = []*cxcore.CXArgument{declarationSpecifiers}
			} else {
				expr.AddOutput(declarationSpecifiers)
			}

			return append([]*cxcore.CXExpression{decl}, initializer...)
		}
	} else {
		// There is no initialization.
		expr := cxcore.MakeExpression(nil, declarator.FileName, declarator.FileLine)
		expr.Package = pkg

		declarationSpecifiers.Name = declarator.Name
		declarationSpecifiers.FileLine = declarator.FileLine
		declarationSpecifiers.Package = pkg
		declarationSpecifiers.PreviouslyDeclared = true
		expr.AddOutput(declarationSpecifiers)

		return []*cxcore.CXExpression{expr}
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
func DeclarationSpecifiers(declSpec *cxcore.CXArgument, arrayLengths []int, opTyp int) *cxcore.CXArgument {
	switch opTyp {
	case cxcore.DECL_POINTER:
		declSpec.DeclarationSpecifiers = append(declSpec.DeclarationSpecifiers, cxcore.DECL_POINTER)
		if !declSpec.IsPointer {
			declSpec.IsPointer = true
			declSpec.Size = cxcore.TYPE_POINTER_SIZE
			declSpec.TotalSize = cxcore.TYPE_POINTER_SIZE
			declSpec.IndirectionLevels++
		} else {
			pointer := declSpec

			for c := declSpec.IndirectionLevels - 1; c > 0; c-- {
				pointer.IndirectionLevels = c
				pointer.IsPointer = true
			}

			declSpec.IndirectionLevels++

			pointer.Size = cxcore.TYPE_POINTER_SIZE
			pointer.TotalSize = cxcore.TYPE_POINTER_SIZE
		}

		return declSpec
	case cxcore.DECL_ARRAY:
		for range arrayLengths {
			declSpec.DeclarationSpecifiers = append(declSpec.DeclarationSpecifiers, cxcore.DECL_ARRAY)
		}
		arg := declSpec
		arg.IsArray = true
		arg.Lengths = arrayLengths
		arg.TotalSize = arg.Size * TotalLength(arg.Lengths)

		return arg
	case cxcore.DECL_SLICE:
		// for range arrayLengths {
		// 	declSpec.DeclarationSpecifiers = append(declSpec.DeclarationSpecifiers, cxcore.DECL_SLICE)
		// }

		arg := declSpec

		arg.DeclarationSpecifiers = append(arg.DeclarationSpecifiers, cxcore.DECL_SLICE)

		arg.IsSlice = true
		arg.IsReference = true
		arg.IsArray = true
		arg.PassBy = cxcore.PASSBY_REFERENCE

		arg.Lengths = append([]int{0}, arg.Lengths...)
		// arg.Lengths = arrayLengths
		// arg.TotalSize = arg.Size
		// arg.Size = cxcore.TYPE_POINTER_SIZE
		arg.TotalSize = cxcore.TYPE_POINTER_SIZE

		return arg
	case cxcore.DECL_BASIC:
		arg := declSpec
		// arg.DeclarationSpecifiers = append(arg.DeclarationSpecifiers, cxcore.DECL_BASIC)
		arg.TotalSize = arg.Size
		return arg
	case cxcore.DECL_FUNC:
		// Creating this case if additional operations are needed in the
		// future.
		return declSpec
	}

	return nil
}

// DeclarationSpecifiersBasic() returns a type specifier created from one of the builtin types.
//
func DeclarationSpecifiersBasic(typ int) *cxcore.CXArgument {
	arg := cxcore.MakeArgument("", CurrentFile, LineNo)
	arg.AddType(cxcore.TypeNames[typ])
	arg.Type = typ

	arg.Size = cxcore.GetArgSize(typ)

	if typ == cxcore.TYPE_AFF {
		// equivalent to slice of strings
		return DeclarationSpecifiers(arg, []int{0}, cxcore.DECL_SLICE)
	}

	return DeclarationSpecifiers(arg, []int{0}, cxcore.DECL_BASIC)
}

// DeclarationSpecifiersStruct() declares a struct
func DeclarationSpecifiersStruct(ident string, pkgName string,
	isExternal bool, currentFile string, lineNo int) *cxcore.CXArgument {
	pkg, err := PRGRM.GetCurrentPackage()
	if err != nil {
		panic(err)
	}

	if isExternal {
		// custom type in an imported package
		imp, err := pkg.GetImport(pkgName)
		if err != nil {
			panic(err)
		}

		strct, err := PRGRM.GetStruct(ident, imp.Name)
		if err != nil {
			println(cxcore.CompilationError(currentFile, lineNo), err.Error())
			return nil
		}

		arg := cxcore.MakeArgument("", currentFile, lineNo)
		arg.Type = cxcore.TYPE_CUSTOM
		arg.CustomType = strct
		arg.Size = strct.Size
		arg.TotalSize = strct.Size

		arg.Package = pkg
		arg.DeclarationSpecifiers = append(arg.DeclarationSpecifiers, cxcore.DECL_STRUCT)

		return arg
	} else {
		// custom type in the current package
		strct, err := PRGRM.GetStruct(ident, pkg.Name)
		if err != nil {
			println(cxcore.CompilationError(currentFile, lineNo), err.Error())
			return nil
		}

		arg := cxcore.MakeArgument("", currentFile, lineNo)
		arg.Type = cxcore.TYPE_CUSTOM
		arg.DeclarationSpecifiers = append(arg.DeclarationSpecifiers, cxcore.DECL_STRUCT)
		arg.CustomType = strct
		arg.Size = strct.Size
		arg.TotalSize = strct.Size
		arg.Package = pkg

		return arg
	}
}
