package actions

import (
	. "github.com/skycoin/cx/cx"
)

// DeclareGlobal() creates a global variable in the current package.
//
// If `doesInitialize` is true, then `initializer` is used to initialize the
// new variable. This function is a wrapper around DeclareGlobalInPackage()
// which does the real work.
//
// FIXME: This function should be merged with DeclareGlobalInPackage.
//        Just use pkg=nil to indicate that CurrentPackage should be used.
//
func DeclareGlobal(declarator *CXArgument, declaration_specifiers *CXArgument,
	           initializer []*CXExpression, doesInitialize bool) {
	pkg, err := PRGRM.GetCurrentPackage()
	if err != nil {
		panic(err)
	}

	DeclareGlobalInPackage(pkg, declarator, declaration_specifiers, initializer, doesInitialize)
}

// DeclareGlobalInPackage() creates a global variable in a specified package
//
// If `doesInitialize` is true, then `initializer` is used to initialize the
// new variable.
//
func DeclareGlobalInPackage(pkg *CXPackage,
			    declarator *CXArgument, declaration_specifiers *CXArgument,
			    initializer []*CXExpression, doesInitialize bool) {
	declaration_specifiers.Package = pkg

	// Treat the name a bit different whether it's defined already or not.
	if glbl, err := pkg.GetGlobal(declarator.Name); err == nil {
		// The name is already defined.

		if glbl.Offset < 0 || glbl.Size == 0 || glbl.TotalSize == 0 {
			// then it was only added a reference to the symbol
			var offExpr []*CXExpression
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

		if doesInitialize {
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
				initializer[len(initializer)-1].Operator = Natives[OP_IDENTITY]
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
					initializer = StructLiteralAssignment([]*CXExpression{&CXExpression{Outputs: []*CXArgument{glbl}}}, initializer)
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
		var offExpr []*CXExpression
		if declaration_specifiers.IsSlice {
			offExpr = WritePrimary(declaration_specifiers.Type, make([]byte, declaration_specifiers.Size), true)
		} else {
			offExpr = WritePrimary(declaration_specifiers.Type, make([]byte, declaration_specifiers.TotalSize), true)
		}
		if doesInitialize {
			if initializer[len(initializer)-1].Operator == nil {
				// then it's a literal

				declaration_specifiers.Name = declarator.Name
				declaration_specifiers.FileLine = declarator.FileLine
				declaration_specifiers.Offset = offExpr[0].Outputs[0].Offset
				declaration_specifiers.Size = offExpr[0].Outputs[0].Size
				declaration_specifiers.TotalSize = offExpr[0].Outputs[0].TotalSize
				declaration_specifiers.Package = pkg

				initializer[len(initializer)-1].Operator = Natives[OP_IDENTITY]
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
					initializer = StructLiteralAssignment([]*CXExpression{&CXExpression{Outputs: []*CXArgument{declaration_specifiers}}}, initializer)
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
func DeclareStruct(ident string, strctFlds []*CXArgument) {
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
			println(CompilationError(fld.FileName, fld.FileLine), "Multiply defined struct field:", fld.Name)
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
		pkg = MakePackage(ident)
		PRGRM.AddPackage(pkg)
	}

	PRGRM.SelectPackage(ident)
}

// DeclareImport()
//
func DeclareImport(ident string, currentFile string, lineNo int) {
        // Make sure we are inside a package
	pkg, err := PRGRM.GetCurrentPackage()
	if err != nil {
		// FIXME: Should give a relevant error message
		panic(err)
	}

	// If the package is already imported, then there is nothing more to be done.
	if _, err := pkg.GetImport(ident); err == nil {
		return;
	}

	// If the package is already defined in the program, just add it to
	// the program.
	if imp, err := PRGRM.GetPackage(ident); err == nil {
		pkg.AddImport(imp)
		return
	}

	// All packages are read during the first pass of the compilation.  So
	// if we get here during the 2nd pass, it's either a core package or
	// something is wrong.
	if IsCorePackage(ident) {
		imp := MakePackage(ident)
		pkg.AddImport(imp)
		PRGRM.AddPackage(imp)
		PRGRM.CurrentPackage = pkg

		if ident == "aff" {
			AffordanceStructs(imp, currentFile, lineNo)
		}
	} else {
		println(CompilationError(currentFile, lineNo), err.Error())
	}
}

// DeclareLocal() creates a local variable inside a function.
// If `doesInitialize` is true, then `initializer` contains the initial values
// of the variable(s).
//
// Returns a list of expressions that contains the initialization, if any.
//
func DeclareLocal(declarator *CXArgument, declaration_specifiers *CXArgument,
     	          initializer []*CXExpression, doesInitialize bool) []*CXExpression {
	if FoundCompileErrors {
		return nil
	}

	declaration_specifiers.IsLocalDeclaration = true

        pkg, err := PRGRM.GetCurrentPackage()
	if err != nil {
		panic(err)
	}

	if doesInitialize {
		// THEN it's a literal, e.g. var foo i32 = 10;
		// ELSE it's an expression with an operator
		if initializer[len(initializer)-1].Operator == nil {
			// we need to create an expression that links the initializer expressions
			// with the declared variable
			expr := MakeExpression(Natives[OP_IDENTITY], CurrentFile, LineNo)
			expr.Package = pkg

			declaration_specifiers.Name = declarator.Name
			declaration_specifiers.FileLine = declarator.FileLine
			declaration_specifiers.Package = pkg
			declaration_specifiers.PreviouslyDeclared = true

			initOut := initializer[len(initializer)-1].Outputs[0]
			// CX checks the output of an expression to determine if it's being passed
			// by value or by reference, so we copy this property from the initializer's
			// output, in case of something like var foo *i32 = &bar
			declaration_specifiers.PassBy = initOut.PassBy

			expr.AddOutput(declaration_specifiers)
			expr.AddInput(initOut)

			initializer[len(initializer)-1] = expr

			return initializer
		} else {
			expr := initializer[len(initializer)-1]

			declaration_specifiers.Name = declarator.Name
			declaration_specifiers.FileLine = declarator.FileLine
			declaration_specifiers.Package = pkg
			declaration_specifiers.PreviouslyDeclared = true
				
			// THEN the expression has outputs created from the result of
			// handling a dot notation initializer, and it needs to be replaced
			// ELSE we simply add it using `AddOutput`
			if len(expr.Outputs) > 0 {
				expr.Outputs = []*CXArgument{declaration_specifiers}
			} else {
				expr.AddOutput(declaration_specifiers)
			}

			return initializer
		}
	} else {
		// There is no initialization.
		expr := MakeExpression(nil, declarator.FileName, declarator.FileLine)
		expr.Package = pkg

		declaration_specifiers.Name = declarator.Name
		declaration_specifiers.FileLine = declarator.FileLine
		declaration_specifiers.Package = pkg
		declaration_specifiers.PreviouslyDeclared = true
		expr.AddOutput(declaration_specifiers)

		return []*CXExpression{expr}
	}
}

// DeclarationSpecifiers() is called to build a type of a variable or parameter.
//
// It is called repeatedly while the type is parsed.
//
//   declSpec:  The incoming type
//   arraySize: The size of the array if `opTyp` = DECL_ARRAY
//   opTyp:     The type of modification to `declSpec` (array of, pointer to, ...)
//
// Returns the new type build from `declSpec` and `opTyp`.
//
func DeclarationSpecifiers(declSpec *CXArgument, arraySize int, opTyp int) *CXArgument {
	switch opTyp {
	case DECL_POINTER:
		declSpec.DeclarationSpecifiers = append(declSpec.DeclarationSpecifiers, DECL_POINTER)
		if !declSpec.IsPointer {
			declSpec.IsPointer = true
			declSpec.Size = TYPE_POINTER_SIZE
			declSpec.TotalSize = TYPE_POINTER_SIZE
			declSpec.IndirectionLevels++
		} else {
			pointer := declSpec

			for c := declSpec.IndirectionLevels - 1; c > 0; c-- {
				pointer.IndirectionLevels = c
				pointer.IsPointer = true
			}

			declSpec.IndirectionLevels++

			pointer.Size = TYPE_POINTER_SIZE
			pointer.TotalSize = TYPE_POINTER_SIZE
		}

		return declSpec
	case DECL_ARRAY:
		declSpec.DeclarationSpecifiers = append(declSpec.DeclarationSpecifiers, DECL_ARRAY)
		arg := declSpec
		arg.IsArray = true
		arg.Lengths = append([]int{arraySize}, arg.Lengths...)
		arg.TotalSize = arg.Size * TotalLength(arg.Lengths)

		return arg
	case DECL_SLICE:
		declSpec.DeclarationSpecifiers = append(declSpec.DeclarationSpecifiers, DECL_SLICE)

		arg := declSpec
		arg.IsSlice = true
		arg.IsReference = true
		arg.IsArray = true
		arg.PassBy = PASSBY_REFERENCE

		arg.Lengths = append([]int{0}, arg.Lengths...)
		arg.TotalSize = arg.Size
		arg.Size = TYPE_POINTER_SIZE

		return arg
	case DECL_BASIC:
		arg := declSpec
		// arg.DeclarationSpecifiers = append(arg.DeclarationSpecifiers, DECL_BASIC)
		arg.TotalSize = arg.Size
		return arg
	}

	return nil
}

// DeclarationSpecifiersBasic() returns a type specifier created from one of the builtin types.
//
func DeclarationSpecifiersBasic(typ int) *CXArgument {
	arg := MakeArgument("", CurrentFile, LineNo)
	arg.AddType(TypeNames[typ])
	arg.Type = typ

	arg.Size = GetArgSize(typ)

	if typ == TYPE_AFF {
		// equivalent to slice of strings
		return DeclarationSpecifiers(arg, 0, DECL_SLICE)
	}

	return DeclarationSpecifiers(arg, 0, DECL_BASIC)
}

// DeclarationSpecifiersStruct() declares a struct
func DeclarationSpecifiersStruct(ident string, pkgName string,
				 isExternal bool, currentFile string, lineNo int) *CXArgument {
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
			println(CompilationError(currentFile, lineNo), err.Error())
			return nil
		}

		arg := MakeArgument("", currentFile, lineNo)
		arg.Type = TYPE_CUSTOM
		arg.CustomType = strct
		arg.Size = strct.Size
		arg.TotalSize = strct.Size

		arg.Package = pkg
		arg.DeclarationSpecifiers = append(arg.DeclarationSpecifiers, DECL_STRUCT)

		return arg
	} else {
		// custom type in the current package
		strct, err := PRGRM.GetStruct(ident, pkg.Name)
		if err != nil {
			println(CompilationError(currentFile, lineNo), err.Error())
			return nil
		}

		arg := MakeArgument("", currentFile, lineNo)
		arg.Type = TYPE_CUSTOM
		arg.DeclarationSpecifiers = append(arg.DeclarationSpecifiers, DECL_STRUCT)
		arg.CustomType = strct
		arg.Size = strct.Size
		arg.TotalSize = strct.Size
		arg.Package = pkg

		return arg
	}
}
