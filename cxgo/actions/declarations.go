package actions

import (
	. "github.com/skycoin/cx/cx"
)

func DeclareGlobal(declarator *CXArgument, declaration_specifiers *CXArgument,
	initializer []*CXExpression, doesInitialize bool) {
	if pkg, err := PRGRM.GetCurrentPackage(); err == nil {
		DeclareGlobalInPackage(pkg, declarator, declaration_specifiers, initializer, doesInitialize)
	} else {
		panic(err)
	}
}
func DeclareGlobalInPackage(pkg *CXPackage, declarator *CXArgument, declaration_specifiers *CXArgument, initializer []*CXExpression, doesInitialize bool) {
	declaration_specifiers.Package = pkg

	if glbl, err := PRGRM.GetGlobal(declarator.Name); err == nil {
		// then it is already defined

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
		}

		if doesInitialize {
			// then we just re-assign offsets
			if initializer[len(initializer)-1].Operator == nil {
				// then it's a literal
				declaration_specifiers.Name = glbl.Name
				declaration_specifiers.Offset = glbl.Offset
				declaration_specifiers.PassBy = glbl.PassBy

				*glbl = *declaration_specifiers

				initializer[len(initializer)-1].AddInput(initializer[len(initializer)-1].Outputs[0])
				initializer[len(initializer)-1].Outputs = nil
				initializer[len(initializer)-1].AddOutput(glbl)
				initializer[len(initializer)-1].Operator = Natives[OP_IDENTITY]

				SysInitExprs = append(SysInitExprs, initializer...)
			} else {
				// then it's an expression
				declaration_specifiers.Name = glbl.Name
				declaration_specifiers.Offset = glbl.Offset
				declaration_specifiers.PassBy = glbl.PassBy

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
			declaration_specifiers.Offset = offExpr[0].Outputs[0].Offset
			declaration_specifiers.Size = offExpr[0].Outputs[0].Size
			declaration_specifiers.TotalSize = offExpr[0].Outputs[0].TotalSize
			declaration_specifiers.Package = pkg

			pkg.AddGlobal(declaration_specifiers)
		}
	}
}

func DeclareStruct(ident string, strctFlds []*CXArgument) {
	if pkg, err := PRGRM.GetCurrentPackage(); err == nil {
		if strct, err := PRGRM.GetStruct(ident, pkg.Name); err == nil {
			strct.Fields = nil
			strct.Size = 0

			// var size int
			for _, fld := range strctFlds {
				strct.AddField(fld)
				// size += fld.TotalSize
			}
			// strct.Size = size
		} else {
			panic(err)
		}
	} else {
		panic(err)
	}
}

func DeclarePackage(ident string) {
	if pkg, err := PRGRM.GetPackage(ident); err != nil {
		pkg := MakePackage(ident)
		// pkg.AddImport(pkg)
		PRGRM.AddPackage(pkg)
		PRGRM.SelectPackage(pkg.Name)
	} else {
		PRGRM.SelectPackage(pkg.Name)
	}
}

func DeclareImport(ident string, currentFile string, lineNo int) {
	if pkg, err := PRGRM.GetCurrentPackage(); err == nil {
		if _, err := pkg.GetImport(ident); err != nil {

			if imp, err := PRGRM.GetPackage(ident); err == nil {
				pkg.AddImport(imp)
			} else {
				// TODO look in the workspace
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
		}
	} else {
		panic(err)
	}
}

func DeclareLocal(declarator *CXArgument, declaration_specifiers *CXArgument, initializer []*CXExpression, doesInitialize bool) []*CXExpression {
	if FoundCompileErrors {
		return nil
	}
	if doesInitialize {
		declaration_specifiers.IsLocalDeclaration = true

		if pkg, err := PRGRM.GetCurrentPackage(); err == nil {
			if initializer[len(initializer)-1].Operator == nil {
				// then it's a literal, e.g. var foo i32 = 10;
				expr := MakeExpression(Natives[OP_IDENTITY], CurrentFile, LineNo)
				expr.Package = pkg

				declaration_specifiers.Name = declarator.Name
				declaration_specifiers.Package = pkg
				declaration_specifiers.PreviouslyDeclared = true

				expr.AddOutput(declaration_specifiers)
				expr.AddInput(initializer[len(initializer)-1].Outputs[0])

				return []*CXExpression{expr}
			} else {
				// then it's an expression (it has an operator)
				declaration_specifiers.Name = declarator.Name
				declaration_specifiers.Package = pkg
				declaration_specifiers.PreviouslyDeclared = true

				expr := initializer[len(initializer)-1]
				expr.AddOutput(declaration_specifiers)

				// exprs := $5
				// exprs = append(exprs, expr)

				return initializer
			}
		} else {
			panic(err)
		}
	} else {
		declaration_specifiers.IsLocalDeclaration = true

		// this will tell the runtime that it's just a declaration
		if pkg, err := PRGRM.GetCurrentPackage(); err == nil {
			expr := MakeExpression(nil, declarator.FileName, declarator.FileLine)
			expr.Package = pkg

			declaration_specifiers.Name = declarator.Name
			declaration_specifiers.Package = pkg
			declaration_specifiers.PreviouslyDeclared = true
			expr.AddOutput(declaration_specifiers)

			return []*CXExpression{expr}
		} else {
			panic(err)
		}
	}
}

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
		arg.DeclarationSpecifiers = append(arg.DeclarationSpecifiers, DECL_BASIC)
		arg.TotalSize = arg.Size
		return arg
	}

	return nil
}

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

func DeclarationSpecifiersStruct(ident string, pkgName string, isExternal bool, currentFile string, lineNo int) *CXArgument {
	if isExternal {
		// custom type in an imported package
		if pkg, err := PRGRM.GetCurrentPackage(); err == nil {
			if imp, err := pkg.GetImport(pkgName); err == nil {
				if strct, err := PRGRM.GetStruct(ident, imp.Name); err == nil {
					arg := MakeArgument("", currentFile, lineNo)
					arg.Type = TYPE_CUSTOM
					arg.CustomType = strct
					arg.Size = strct.Size
					arg.TotalSize = strct.Size

					arg.Package = pkg
					arg.DeclarationSpecifiers = append(arg.DeclarationSpecifiers, DECL_STRUCT)

					return arg
				} else {
					println(CompilationError(currentFile, lineNo), err.Error())
					return nil
				}
			} else {
				panic(err)
			}
		} else {
			panic(err)
		}
	} else {
		// custom type in the current package
		if pkg, err := PRGRM.GetCurrentPackage(); err == nil {
			if strct, err := PRGRM.GetStruct(ident, pkg.Name); err == nil {
				arg := MakeArgument("", currentFile, lineNo)
				arg.Type = TYPE_CUSTOM
				arg.DeclarationSpecifiers = append(arg.DeclarationSpecifiers, DECL_STRUCT)
				arg.CustomType = strct
				arg.Size = strct.Size
				arg.TotalSize = strct.Size
				arg.Package = pkg

				return arg
			} else {
				println(CompilationError(currentFile, lineNo), err.Error())
				return nil
			}
		} else {
			panic(err)
		}
	}
}
