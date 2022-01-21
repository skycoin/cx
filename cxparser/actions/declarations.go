package actions

import (
	"fmt"
	"os"

	constants2 "github.com/skycoin/cx/cxparser/constants"

	"github.com/skycoin/cx/cx/ast"
	"github.com/skycoin/cx/cx/constants"
	globals2 "github.com/skycoin/cx/cx/globals"
	"github.com/skycoin/cx/cx/types"
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
func DeclareGlobal(prgrm *ast.CXProgram, declarator *ast.CXArgument, declarationSpecifiers *ast.CXArgument,
	initializer []*ast.CXExpression, doesInitialize bool) {
	pkg, err := prgrm.GetCurrentPackage()
	if err != nil {
		panic(err)
	}

	DeclareGlobalInPackage(prgrm, pkg, declarator, declarationSpecifiers, initializer, doesInitialize)
}

// DeclareGlobalInPackage creates a global variable in a specified package
//
// If `doesInitialize` is true, then `initializer` is used to initialize the
// new variable.
//
func DeclareGlobalInPackage(prgrm *ast.CXProgram, pkg *ast.CXPackage,
	declarator *ast.CXArgument, declaration_specifiers *ast.CXArgument,
	initializer []*ast.CXExpression, doesInitialize bool) {
	declaration_specifiers.Package = ast.CXPackageIndex(pkg.Index)

	// Treat the name a bit different whether it's defined already or not.
	if glbl, err := pkg.GetGlobal(declarator.Name); err == nil {
		// The name is already defined.

		if glbl.Offset < 0 || glbl.Size == 0 || glbl.TotalSize == 0 {
			// then it was only added a reference to the symbol
			var offExpr []*ast.CXExpression
			if declaration_specifiers.IsSlice { // TODO:PTR move branch in WritePrimary
				offExpr = WritePrimary(prgrm, declaration_specifiers.Type,
					make([]byte, types.POINTER_SIZE), true)
			} else {
				offExpr = WritePrimary(prgrm, declaration_specifiers.Type,
					make([]byte, declaration_specifiers.TotalSize), false)
			}

			offExprAtomicOp, err := prgrm.GetCXAtomicOpFromExpressions(offExpr, 0)
			if err != nil {
				panic(err)
			}

			glbl.Offset = offExprAtomicOp.Outputs[0].Offset
			glbl.PassBy = offExprAtomicOp.Outputs[0].PassBy
			// glbl.Package = offExpr[0].ProgramOutput[0].Package
		}

		// Checking if something is supposed to be initialized
		// and if `initializer` actually contains something.
		// If `initializer` is `nil`, this means that an expression
		// equivalent to nil was used, such as `[]i32{}`.
		if doesInitialize && initializer != nil {
			initializerAtomicOp, err := prgrm.GetCXAtomicOpFromExpressions(initializer, len(initializer)-1)
			if err != nil {
				panic(err)
			}

			// then we just re-assign offsets
			if initializerAtomicOp.Operator == nil {
				// then it's a literal
				declaration_specifiers.Name = glbl.Name
				declaration_specifiers.Offset = glbl.Offset
				declaration_specifiers.PassBy = glbl.PassBy
				declaration_specifiers.Package = glbl.Package

				*glbl = *declaration_specifiers

				initializerAtomicOp.AddInput(initializerAtomicOp.Outputs[0])
				initializerAtomicOp.Outputs = nil
				initializerAtomicOp.AddOutput(glbl)
				initializerAtomicOp.Operator = ast.Natives[constants.OP_IDENTITY]
				initializerAtomicOp.Package = glbl.Package

				//add intialization statements, to array
				prgrm.SysInitExprs = append(prgrm.SysInitExprs, initializer...)
			} else {
				// then it's an expression
				declaration_specifiers.Name = glbl.Name
				declaration_specifiers.Offset = glbl.Offset
				declaration_specifiers.PassBy = glbl.PassBy
				declaration_specifiers.Package = glbl.Package

				*glbl = *declaration_specifiers

				if initializer[len(initializer)-1].IsStructLiteral() {
					index := prgrm.AddCXAtomicOp(&ast.CXAtomicOperator{Outputs: []*ast.CXArgument{glbl}})
					initializer = StructLiteralAssignment(prgrm,
						[]*ast.CXExpression{
							{
								Index: index,
								Type:  ast.CX_ATOMIC_OPERATOR,
							},
						},
						initializer,
					)
				} else {
					initializerAtomicOp.Outputs = nil
					initializerAtomicOp.AddOutput(glbl)
				}
				//add intialization statements, to array
				prgrm.SysInitExprs = append(prgrm.SysInitExprs, initializer...)
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
		var offExpr []*ast.CXExpression
		if declaration_specifiers.IsSlice { // TODO:PTR move branch in WritePrimary
			offExpr = WritePrimary(prgrm, declaration_specifiers.Type, make([]byte, types.POINTER_SIZE), true)
		} else {
			offExpr = WritePrimary(prgrm, declaration_specifiers.Type, make([]byte, declaration_specifiers.TotalSize), false)
		}

		// Checking if something is supposed to be initialized
		// and if `initializer` actually contains something.
		// If `initializer` is `nil`, this means that an expression
		// equivalent to nil was used, such as `[]i32{}`.
		if doesInitialize && initializer != nil {
			initializerAtomicOp, err := prgrm.GetCXAtomicOpFromExpressions(initializer, len(initializer)-1)
			if err != nil {
				panic(err)
			}

			offExprAtomicOp, err := prgrm.GetCXAtomicOpFromExpressions(offExpr, 0)
			if err != nil {
				panic(err)
			}

			if initializerAtomicOp.Operator == nil {
				// then it's a literal

				declaration_specifiers.Name = declarator.Name
				declaration_specifiers.ArgDetails.FileLine = declarator.ArgDetails.FileLine
				declaration_specifiers.Offset = offExprAtomicOp.Outputs[0].Offset
				declaration_specifiers.Size = offExprAtomicOp.Outputs[0].Size
				declaration_specifiers.TotalSize = offExprAtomicOp.Outputs[0].TotalSize
				declaration_specifiers.Package = ast.CXPackageIndex(pkg.Index)

				initializerAtomicOp.Operator = ast.Natives[constants.OP_IDENTITY]
				initializerAtomicOp.AddInput(initializerAtomicOp.Outputs[0])
				initializerAtomicOp.Outputs = nil
				initializerAtomicOp.AddOutput(declaration_specifiers)

				pkg.AddGlobal(declaration_specifiers)
				//add intialization statements, to array
				prgrm.SysInitExprs = append(prgrm.SysInitExprs, initializer...)
			} else {
				// then it's an expression
				declaration_specifiers.Name = declarator.Name
				declaration_specifiers.ArgDetails.FileLine = declarator.ArgDetails.FileLine
				declaration_specifiers.Offset = offExprAtomicOp.Outputs[0].Offset
				declaration_specifiers.Size = offExprAtomicOp.Outputs[0].Size
				declaration_specifiers.TotalSize = offExprAtomicOp.Outputs[0].TotalSize
				declaration_specifiers.Package = ast.CXPackageIndex(pkg.Index)

				if initializer[len(initializer)-1].IsStructLiteral() {
					index := prgrm.AddCXAtomicOp(&ast.CXAtomicOperator{Outputs: []*ast.CXArgument{declaration_specifiers}})
					initializer = StructLiteralAssignment(prgrm,
						[]*ast.CXExpression{
							{
								Index: index,
								Type:  ast.CX_ATOMIC_OPERATOR,
							},
						},
						initializer,
					)
				} else {
					initializerAtomicOp.Outputs = nil
					initializerAtomicOp.AddOutput(declaration_specifiers)
				}

				pkg.AddGlobal(declaration_specifiers)
				//add intialization statements, to array
				prgrm.SysInitExprs = append(prgrm.SysInitExprs, initializer...)
			}
		} else {
			// offExpr := WritePrimary(declaration_specifiers.Type, make([]byte, declaration_specifiers.Size), false)
			// exprOut := expr[0].ProgramOutput[0]

			offExprAtomicOp, err := prgrm.GetCXAtomicOpFromExpressions(offExpr, 0)
			if err != nil {
				panic(err)
			}

			declaration_specifiers.Name = declarator.Name
			declaration_specifiers.ArgDetails.FileLine = declarator.ArgDetails.FileLine
			declaration_specifiers.Offset = offExprAtomicOp.Outputs[0].Offset
			declaration_specifiers.Size = offExprAtomicOp.Outputs[0].Size
			declaration_specifiers.TotalSize = offExprAtomicOp.Outputs[0].TotalSize
			declaration_specifiers.Package = ast.CXPackageIndex(pkg.Index)

			pkg.AddGlobal(declaration_specifiers)
		}
	}
}

// DeclareStruct takes a name of a struct and a slice of fields representing
// the members and adds the struct to the package.
//
func DeclareStruct(prgrm *ast.CXProgram, ident string, strctFlds []*ast.CXArgument) {
	// Make sure we are inside a package.
	pkg, err := prgrm.GetCurrentPackage()
	if err != nil {
		// FIXME: Should give a relevant error message
		panic(err)
	}

	// Make sure a struct with the same name is not yet defined.
	strct, err := prgrm.GetStruct(ident, pkg.Name)
	if err != nil {
		// FIXME: Should give a relevant error message
		panic(err)
	}

	strct.Fields = nil
	strct.Size = 0
	for _, fld := range strctFlds {
		if _, err := strct.GetField(fld.Name); err == nil {
			println(ast.CompilationError(fld.ArgDetails.FileName, fld.ArgDetails.FileLine), "Multiply defined struct field:", fld.Name)
		} else {
			strct.AddField(fld)
		}
	}
}

// DeclarePackage() switches the current package in the program.
//
func DeclarePackage(prgrm *ast.CXProgram, ident string) {
	// Add a new package to the program if it's not previously defined.
	if _, err := prgrm.GetPackage(ident); err != nil {
		pkg := ast.MakePackage(ident)
		prgrm.AddPackage(pkg)
	}

	prgrm.SelectPackage(ident)
}

// DeclareImport()
//
func DeclareImport(prgrm *ast.CXProgram, name string, currentFile string, lineNo int) {
	// Make sure we are inside a package
	pkg, err := prgrm.GetCurrentPackage()
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
	if _, err := pkg.GetImport(prgrm, ident); err == nil {
		return
	}

	// If the package is already defined in the program, just add it to
	// the importing package.
	if imp, err := prgrm.GetPackage(ident); err == nil {
		pkg.AddImport(prgrm, imp)
		return
	}

	// All packages are read during the first pass of the compilation.  So
	// if we get here during the 2nd pass, it's either a core package or
	// something is panic-level wrong.
	if constants2.IsCorePackage(ident) {
		imp := ast.MakePackage(ident)
		impIdx := prgrm.AddPackage(imp)
		newImp, err := prgrm.GetPackageFromArray(impIdx)
		if err != nil {
			panic(err)
		}
		pkg.AddImport(prgrm, newImp)

		prgrm.CurrentPackage = ast.CXPackageIndex(pkg.Index)

		if ident == "aff" {
			AffordanceStructs(prgrm, newImp, currentFile, lineNo)
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
func DeclareLocal(prgrm *ast.CXProgram, declarator *ast.CXArgument, declarationSpecifiers *ast.CXArgument,
	initializer []*ast.CXExpression, doesInitialize bool) []*ast.CXExpression {
	if globals2.FoundCompileErrors {
		return nil
	}

	declarationSpecifiers.IsLocalDeclaration = true

	pkg, err := prgrm.GetCurrentPackage()
	if err != nil {
		panic(err)
	}

	declCXLine := ast.MakeCXLineExpression(prgrm, CurrentFile, LineNo, LineStr)
	// Declaration expression to handle the inline initialization.
	// For example, `var foo i32 = 11` needs to be divided into two expressions:
	// one that declares `foo`, and another that assigns 11 to `foo`
	decl := ast.MakeAtomicOperatorExpression(prgrm, nil)
	cxAtomicOp, _, _, err := prgrm.GetOperation(decl)
	if err != nil {
		panic(err)
	}
	cxAtomicOp.Package = ast.CXPackageIndex(pkg.Index)

	declarationSpecifiers.Name = declarator.Name
	declarationSpecifiers.ArgDetails.FileLine = declarator.ArgDetails.FileLine
	declarationSpecifiers.Package = ast.CXPackageIndex(pkg.Index)
	declarationSpecifiers.PreviouslyDeclared = true
	cxAtomicOp.AddOutput(declarationSpecifiers)

	// Checking if something is supposed to be initialized
	// and if `initializer` actually contains something.
	// If `initializer` is `nil`, this means that an expression
	// equivalent to nil was used, such as `[]i32{}`.
	if doesInitialize && initializer != nil {
		initializerAtomicOp, err := prgrm.GetCXAtomicOpFromExpressions(initializer, len(initializer)-1)
		if err != nil {
			panic(err)
		}

		// THEN it's a literal, e.g. var foo i32 = 10;
		// ELSE it's an expression with an operator
		if initializerAtomicOp.Operator == nil {

			exprCXLine := ast.MakeCXLineExpression(prgrm, CurrentFile, LineNo, LineStr)
			// we need to create an expression that links the initializer expressions
			// with the declared variable
			expr := ast.MakeAtomicOperatorExpression(prgrm, ast.Natives[constants.OP_IDENTITY])
			cxExprAtomicOp, _, _, err := prgrm.GetOperation(expr)
			if err != nil {
				panic(err)
			}
			cxExprAtomicOp.Package = ast.CXPackageIndex(pkg.Index)

			initOut := initializerAtomicOp.Outputs[0]
			// CX checks the output of an expression to determine if it's being passed
			// by value or by reference, so we copy this property from the initializer's
			// output, in case of something like var foo *i32 = &bar
			declarationSpecifiers.PassBy = initOut.PassBy

			cxExprAtomicOp.AddOutput(declarationSpecifiers)
			cxExprAtomicOp.AddInput(initOut)

			initializer[len(initializer)-1] = exprCXLine
			initializer = append(initializer, expr)

			return append([]*ast.CXExpression{declCXLine, decl}, initializer...)
		} else {
			expr := initializer[len(initializer)-1]
			cxExprAtomicOp, _, _, err := prgrm.GetOperation(expr)
			if err != nil {
				panic(err)
			}

			// THEN the expression has outputs created from the result of
			// handling a dot notation initializer, and it needs to be replaced
			// ELSE we simply add it using `AddOutput`
			if len(cxExprAtomicOp.Outputs) > 0 {
				cxExprAtomicOp.Outputs = []*ast.CXArgument{declarationSpecifiers}
			} else {
				cxExprAtomicOp.AddOutput(declarationSpecifiers)
			}

			return append([]*ast.CXExpression{declCXLine, decl}, initializer...)
		}
	} else {
		exprCXLine := ast.MakeCXLineExpression(prgrm, CurrentFile, LineNo, LineStr)
		// There is no initialization.
		expr := ast.MakeAtomicOperatorExpression(prgrm, nil)
		cxAtomicOp, _, _, err := prgrm.GetOperation(expr)
		if err != nil {
			panic(err)
		}
		cxAtomicOp.Package = ast.CXPackageIndex(pkg.Index)

		declarationSpecifiers.Name = declarator.Name
		declarationSpecifiers.ArgDetails.FileLine = declarator.ArgDetails.FileLine
		declarationSpecifiers.Package = ast.CXPackageIndex(pkg.Index)
		declarationSpecifiers.PreviouslyDeclared = true
		cxAtomicOp.AddOutput(declarationSpecifiers)

		return []*ast.CXExpression{exprCXLine, expr}
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
func DeclarationSpecifiers(declSpec *ast.CXArgument, arrayLengths []types.Pointer, opTyp int) *ast.CXArgument {
	switch opTyp {
	case constants.DECL_POINTER:
		declSpec.DeclarationSpecifiers = append(declSpec.DeclarationSpecifiers, constants.DECL_POINTER)

		declSpec.Size = types.POINTER_SIZE
		declSpec.TotalSize = types.POINTER_SIZE
		// declSpec.IndirectionLevels++

		if declSpec.Type == types.STR || declSpec.StructType != nil {
			declSpec.PointerTargetType = declSpec.Type
			declSpec.Type = types.POINTER
		}
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
		// arg.IsReference = true
		// arg.IsArray = true
		arg.PassBy = constants.PASSBY_REFERENCE

		arg.Lengths = append([]types.Pointer{0}, arg.Lengths...)
		// arg.Lengths = arrayLengths
		// arg.TotalSize = arg.Size
		// arg.Size = cxcore.TYPE_POINTER_SIZE
		arg.TotalSize = types.POINTER_SIZE

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
func DeclarationSpecifiersBasic(typeCode types.Code) *ast.CXArgument {
	arg := ast.MakeArgument("", CurrentFile, LineNo)
	arg.AddType(typeCode)
	if typeCode == types.AFF {
		// equivalent to slice of strings
		return DeclarationSpecifiers(arg, []types.Pointer{0}, constants.DECL_SLICE)
	}

	return DeclarationSpecifiers(arg, []types.Pointer{0}, constants.DECL_BASIC)
}

// DeclarationSpecifiersStruct() declares a struct
func DeclarationSpecifiersStruct(prgrm *ast.CXProgram, ident string, pkgName string,
	isExternal bool, currentFile string, lineNo int) *ast.CXArgument {
	pkg, err := prgrm.GetCurrentPackage()
	if err != nil {
		panic(err)
	}

	if isExternal {
		// custom type in an imported package
		imp, err := pkg.GetImport(prgrm, pkgName)
		if err != nil {
			panic(err)
		}

		strct, err := prgrm.GetStruct(ident, imp.Name)
		if err != nil {
			println(ast.CompilationError(currentFile, lineNo), err.Error())
			return nil
		}

		arg := ast.MakeArgument("", currentFile, lineNo)
		arg.Type = types.STRUCT
		arg.StructType = strct
		arg.Size = strct.Size
		arg.TotalSize = strct.Size

		arg.Package = ast.CXPackageIndex(pkg.Index)
		arg.DeclarationSpecifiers = append(arg.DeclarationSpecifiers, constants.DECL_STRUCT)

		return arg
	} else {
		// custom type in the current package
		strct, err := prgrm.GetStruct(ident, pkg.Name)
		if err != nil {
			println(ast.CompilationError(currentFile, lineNo), err.Error())
			return nil
		}

		arg := ast.MakeArgument("", currentFile, lineNo)
		arg.Type = types.STRUCT
		arg.DeclarationSpecifiers = append(arg.DeclarationSpecifiers, constants.DECL_STRUCT)
		arg.StructType = strct
		arg.Size = strct.Size
		arg.TotalSize = strct.Size
		arg.Package = ast.CXPackageIndex(pkg.Index)

		return arg
	}
}
