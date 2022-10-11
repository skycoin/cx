package actions

import (
	"github.com/skycoin/cx/cx/ast"
	"github.com/skycoin/cx/cx/constants"
	"github.com/skycoin/cx/cx/globals"
	"github.com/skycoin/cx/cx/types"
)

func processGlobalInitializer(prgrm *ast.CXProgram, initializer []ast.CXExpression, glblIdx ast.CXTypeSignatureIndex) {
	initializerExpressionIdx := initializer[len(initializer)-1].Index
	initializerExpressionOperator := prgrm.GetFunctionFromArray(prgrm.CXAtomicOps[initializerExpressionIdx].Operator)
	// then we just re-assign offsets
	if initializerExpressionOperator == nil {
		// then it's a literal
		typeSigIdx := prgrm.CXAtomicOps[initializerExpressionIdx].GetOutputs(prgrm)[0]
		prgrm.CXAtomicOps[initializerExpressionIdx].AddInput(prgrm, typeSigIdx)
		prgrm.CXAtomicOps[initializerExpressionIdx].Outputs = nil

		prgrm.CXAtomicOps[initializerExpressionIdx].AddOutput(prgrm, glblIdx)
		opIdx := prgrm.AddNativeFunctionInArray(ast.Natives[constants.OP_IDENTITY])
		prgrm.CXAtomicOps[initializerExpressionIdx].Operator = opIdx

		//add intialization statements, to array
		prgrm.SysInitExprs = append(prgrm.SysInitExprs, initializer...)
	} else {
		// then it's an expression
		if initializer[len(initializer)-1].IsStructLiteral() {
			outputStruct := &ast.CXStruct{}
			outputStruct.AddField_TypeSignature(prgrm, glblIdx)
			index := prgrm.AddCXAtomicOp(&ast.CXAtomicOperator{Outputs: outputStruct, Operator: -1, Function: -1})
			initializer = StructLiteralAssignment(prgrm,
				[]ast.CXExpression{
					{
						Index: index,
						Type:  ast.CX_ATOMIC_OPERATOR,
					},
				},
				initializer,
			)
		} else {
			prgrm.CXAtomicOps[initializerExpressionIdx].Outputs = nil
			prgrm.CXAtomicOps[initializerExpressionIdx].AddOutput(prgrm, glblIdx)
		}

		//add intialization statements, to array
		prgrm.SysInitExprs = append(prgrm.SysInitExprs, initializer...)
	}
}

// DeclareGlobalInPackage creates a global variable in a specified package
//
// If `doesInitialize` is true, then `initializer` is used to initialize the
// new variable.
//
// Input arguments description:
// 	prgrm - a CXProgram that contains all the data and arrays of the program.
// 	pkg - the package where the global will belong.
// 	declarator - contains the name of the global var.
// 	declaration_specifiers - contains the type build of the global.
// 	initializer - if doesInitialize is true then this contains the initial
// 				  value of the global.
// 	doesInitialize - true if global is initialized or given a value.
func DeclareGlobalInPackage(prgrm *ast.CXProgram, pkg *ast.CXPackage,
	declarator *ast.CXArgument, declaration_specifiers *ast.CXArgument,
	initializer []ast.CXExpression, doesInitialize bool) {
	if pkg == nil {
		var err error
		pkg, err = prgrm.GetCurrentPackage()
		if err != nil {
			panic(err)
		}
	}
	declaration_specifiers.Package = ast.CXPackageIndex(pkg.Index)

	totalSize := declaration_specifiers.Size
	for _, decl := range declaration_specifiers.DeclarationSpecifiers {
		switch decl {
		case constants.DECL_POINTER, constants.DECL_SLICE:
			totalSize = types.POINTER_SIZE
		case constants.DECL_ARRAY:
			totalSize = declaration_specifiers.Size * TotalLength(declaration_specifiers.Lengths)
		}
	}

	// Treat the name a bit different whether it's defined already or not.
	if glbl, err := pkg.GetGlobal(prgrm, declarator.Name); err == nil {
		// then it was only added a reference to the symbol

		// if offset is not valid, assign new offset
		var glblArg *ast.CXArgument = &ast.CXArgument{}
		if glbl.Offset < 0 || !glbl.Offset.IsValid() {
			if declaration_specifiers.IsSlice { // TODO:PTR move branch in WritePrimary
				glblArg = WritePrimary(prgrm, declaration_specifiers.Type,
					make([]byte, types.POINTER_SIZE), true)
			} else {
				glblArg = WritePrimary(prgrm, declaration_specifiers.Type,
					make([]byte, totalSize), false)
			}

		}

		if ast.IsTypeAtomic(declaration_specifiers) {
			glbl.Type = ast.TYPE_ATOMIC
			glbl.Meta = int(declaration_specifiers.Type)
			if glbl.Offset < 0 || !glbl.Offset.IsValid() {
				glbl.Offset = glblArg.Offset
			}

			glbl.PassBy = declaration_specifiers.PassBy
		} else if ast.IsTypePointerAtomic(declaration_specifiers) {
			glbl.Type = ast.TYPE_POINTER_ATOMIC
			glbl.Meta = int(declaration_specifiers.Type)
			if glbl.Offset < 0 || !glbl.Offset.IsValid() {
				glbl.Offset = glblArg.Offset
			}

			glbl.PassBy = declaration_specifiers.PassBy
		} else if ast.IsTypeArrayAtomic(declaration_specifiers) {
			glbl.Type = ast.TYPE_ARRAY_ATOMIC

			typeSignatureForArray := &ast.CXTypeSignature_Array{
				Type:    int(declaration_specifiers.Type),
				Lengths: declaration_specifiers.Lengths,
				Indexes: declaration_specifiers.Indexes,
			}
			typeSignatureForArrayIdx := prgrm.AddCXTypeSignatureArrayInArray(typeSignatureForArray)

			glbl.Meta = typeSignatureForArrayIdx

			glbl.PassBy = declaration_specifiers.PassBy

			if glbl.Offset < 0 || !glbl.Offset.IsValid() {
				glbl.Offset = glblArg.Offset
			}
		} else if ast.IsTypePointerArrayAtomic(declaration_specifiers) {
			glbl.Type = ast.TYPE_POINTER_ARRAY_ATOMIC

			typeSignatureForArray := &ast.CXTypeSignature_Array{
				Type:    int(declaration_specifiers.Type),
				Lengths: declaration_specifiers.Lengths,
				Indexes: declaration_specifiers.Indexes,
			}

			typeSignatureForArrayIdx := prgrm.AddCXTypeSignatureArrayInArray(typeSignatureForArray)
			glbl.Meta = typeSignatureForArrayIdx
			glbl.PassBy = declaration_specifiers.PassBy

			if glbl.Offset < 0 || !glbl.Offset.IsValid() {
				glbl.Offset = glblArg.Offset
			}
		} else if ast.IsTypeSliceAtomic(declaration_specifiers) {
			// If slice atomic type, i.e. []i32, []f64, etc.

			glbl.Type = ast.TYPE_SLICE_ATOMIC

			typeSignatureForArray := &ast.CXTypeSignature_Array{
				Type:    int(declaration_specifiers.Type),
				Lengths: declaration_specifiers.Lengths,
				Indexes: declaration_specifiers.Indexes,
			}

			typeSignatureForArrayIdx := prgrm.AddCXTypeSignatureArrayInArray(typeSignatureForArray)
			glbl.Meta = typeSignatureForArrayIdx
			glbl.PassBy = declaration_specifiers.PassBy

			if glbl.Offset < 0 || !glbl.Offset.IsValid() {
				glbl.Offset = glblArg.Offset
			}
		} else if ast.IsTypePointerSliceAtomic(declaration_specifiers) {
			// If pointer slice atomic type, i.e. *[]i32, *[]f64, etc.

			glbl.Type = ast.TYPE_POINTER_SLICE_ATOMIC

			typeSignatureForArray := &ast.CXTypeSignature_Array{
				Type:    int(declaration_specifiers.Type),
				Lengths: declaration_specifiers.Lengths,
				Indexes: declaration_specifiers.Indexes,
			}

			typeSignatureForArrayIdx := prgrm.AddCXTypeSignatureArrayInArray(typeSignatureForArray)
			glbl.Meta = typeSignatureForArrayIdx
			glbl.PassBy = declaration_specifiers.PassBy

			if glbl.Offset < 0 || !glbl.Offset.IsValid() {
				glbl.Offset = glblArg.Offset
			}
		} else {
			// its a cxargument_deprecate type
			var glblIdx int
			var glblCXArg *ast.CXArgument = &ast.CXArgument{}
			if glbl.Type == ast.TYPE_CXARGUMENT_DEPRECATE {
				glblIdx = glbl.Meta
				glblCXArg = prgrm.GetCXArgFromArray(ast.CXArgumentIndex(glblIdx))
			} else {
				panic("type is not type cx argument deprecate\n\n")
			}

			if glbl.Offset < 0 || glblCXArg.Size == 0 || totalSize == 0 {
				prgrm.CXArgs[glblIdx].Offset = glblArg.Offset
				prgrm.CXArgs[glblIdx].PassBy = glblArg.PassBy
			}

			declaration_specifiers.Name = prgrm.CXArgs[glblIdx].Name
			declaration_specifiers.Offset = prgrm.CXArgs[glblIdx].Offset
			declaration_specifiers.PassBy = prgrm.CXArgs[glblIdx].PassBy
			declaration_specifiers.Package = prgrm.CXArgs[glblIdx].Package

			prgrm.CXArgs[glblIdx] = *declaration_specifiers
			prgrm.CXArgs[glblIdx].Index = glblIdx

		}

		// Checking if something is supposed to be initialized
		// and if `initializer` actually contains something.
		// If `initializer` is `nil`, this means that an expression
		// equivalent to nil was used, such as `[]i32{}`.
		if doesInitialize && initializer != nil {
			processGlobalInitializer(prgrm, initializer, glbl.Index)
		}
	} else {
		// then it hasn't been defined
		var glblArg *ast.CXArgument = &ast.CXArgument{}
		if declaration_specifiers.IsSlice { // TODO:PTR move branch in WritePrimary
			glblArg = WritePrimary(prgrm, declaration_specifiers.Type, make([]byte, types.POINTER_SIZE), true)
		} else {
			glblArg = WritePrimary(prgrm, declaration_specifiers.Type, make([]byte, totalSize), false)
		}

		declaration_specifiers.Name = declarator.Name
		declaration_specifiers.ArgDetails.FileLine = declarator.ArgDetails.FileLine
		declaration_specifiers.Offset = glblArg.Offset
		declaration_specifiers.Size = glblArg.Size
		declaration_specifiers.Package = ast.CXPackageIndex(pkg.Index)

		typeSignature := ast.GetCXTypeSignatureRepresentationOfCXArg(prgrm, declaration_specifiers)
		typeSignatureIdx := prgrm.AddCXTypeSignatureInArray(typeSignature)
		pkg.AddGlobal_TypeSignature(prgrm, typeSignatureIdx)

		// Checking if something is supposed to be initialized
		// and if `initializer` actually contains something.
		// If `initializer` is `nil`, this means that an expression
		// equivalent to nil was used, such as `[]i32{}`.
		if doesInitialize && initializer != nil {
			processGlobalInitializer(prgrm, initializer, typeSignatureIdx)
		}
	}

}

func processLocalInitialization(prgrm *ast.CXProgram, pkg *ast.CXPackage, initializer *[]ast.CXExpression, localVarIdx ast.CXTypeSignatureIndex, declSpecIdx ast.CXArgumentIndex) {
	initializerExpression, err := prgrm.GetCXAtomicOpFromExpressions(*initializer, len(*initializer)-1)
	if err != nil {
		panic(err)
	}

	initializerExpressionOperator := prgrm.GetFunctionFromArray(initializerExpression.Operator)
	if initializerExpressionOperator == nil {
		// THEN it's a literal, e.g. var foo i32 = 10;

		exprCXLine := ast.MakeCXLineExpression(prgrm, CurrentFile, LineNo, LineStr)

		// we need to create an expression that links the initializer expressions
		// with the declared variable
		expr := ast.MakeAtomicOperatorExpression(prgrm, ast.Natives[constants.OP_IDENTITY])
		cxExprAtomicOpIdx := expr.Index
		prgrm.CXAtomicOps[cxExprAtomicOpIdx].Package = ast.CXPackageIndex(pkg.Index)

		var initOutTypeSigIdx ast.CXTypeSignatureIndex
		initializerExpressionOutputTypeSig := prgrm.GetCXTypeSignatureFromArray(initializerExpression.GetOutputs(prgrm)[0])
		if initializerExpressionOutputTypeSig.Type == ast.TYPE_CXARGUMENT_DEPRECATE {
			initializerExpressionOutputArg := prgrm.GetCXArgFromArray(ast.CXArgumentIndex(initializerExpressionOutputTypeSig.Meta))
			initOut := initializerExpressionOutputArg
			initOutIdx := ast.CXArgumentIndex(initializerExpressionOutputTypeSig.Meta)

			// CX checks the output of an expression to determine if it's being passed
			// by value or by reference, so we copy this property from the initializer's
			// output, in case of something like var foo *i32 = &bar
			prgrm.CXArgs[declSpecIdx].PassBy = initOut.PassBy

			typeSig := ast.GetCXTypeSignatureRepresentationOfCXArg(prgrm, prgrm.GetCXArgFromArray(initOutIdx))
			initOutTypeSigIdx = prgrm.AddCXTypeSignatureInArray(typeSig)
		} else if initializerExpressionOutputTypeSig.Type == ast.TYPE_ATOMIC || initializerExpressionOutputTypeSig.Type == ast.TYPE_POINTER_ATOMIC {
			initOutTypeSigIdx = initializerExpression.GetOutputs(prgrm)[0]
		} else {
			panic("type is not known")
		}

		prgrm.CXAtomicOps[cxExprAtomicOpIdx].AddOutput(prgrm, localVarIdx)
		prgrm.CXAtomicOps[cxExprAtomicOpIdx].AddInput(prgrm, initOutTypeSigIdx)

		(*initializer)[len(*initializer)-1] = *exprCXLine
		(*initializer) = append(*initializer, *expr)

	} else {
		// ELSE it's an expression with an operator

		expr := (*initializer)[len(*initializer)-1]
		cxExprAtomicOp, err := prgrm.GetCXAtomicOp(expr.Index)
		if err != nil {
			panic(err)
		}

		// THEN the expression has outputs created from the result of
		// handling a dot notation initializer, and it needs to be replaced
		// ELSE we simply add it using `AddOutput`
		if len(cxExprAtomicOp.GetOutputs(prgrm)) > 0 {
			// declSpecIdx := prgrm.AddCXArgInArray(declarationSpecifiers)
			cxExprAtomicOp.Outputs.Fields = nil

		}
		cxExprAtomicOp.AddOutput(prgrm, localVarIdx)
	}
}

// DeclareLocal() creates a local variable inside a function.
//
// Returns a list of expressions that contains the initialization, if any.
//
// Input arguments description:
// 	prgrm - a CXProgram that contains all the data and arrays of the program.
// 	declarator - contains the name of the var.
// 	declaration_specifiers - contains the type build of the var.
// 	initializer - if doesInitialize is true then this contains the initial
// 				  value of the var.
// 	doesInitialize - true if var is initialized or given a value.
func DeclareLocal(prgrm *ast.CXProgram, declarator *ast.CXArgument, declarationSpecifiers *ast.CXArgument,
	initializer []ast.CXExpression, doesInitialize bool) []ast.CXExpression {
	if globals.FoundCompileErrors {
		return nil
	}

	currFn, err := prgrm.GetCurrentFunction()
	if err != nil {
		// TODO: improve error handling
		panic("error getting current function")
	}
	err = currFn.AddLocalVariableName(declarator.Name)
	if err != nil {
		panic("error adding local variable")
	}

	pkg, err := prgrm.GetCurrentPackage()
	if err != nil {
		panic(err)
	}

	declCXLine := ast.MakeCXLineExpression(prgrm, CurrentFile, LineNo, LineStr)
	// Declaration expression to handle the inline initialization.
	// For example, `var foo i32 = 11` needs to be divided into two expressions:
	// one that declares `foo`, and another that assigns 11 to `foo`
	decl := ast.MakeAtomicOperatorExpression(prgrm, nil)
	expressionIdx := decl.Index
	prgrm.CXAtomicOps[expressionIdx].Package = ast.CXPackageIndex(pkg.Index)

	declarationSpecifiers.Name = declarator.Name
	declarationSpecifiers.ArgDetails.FileLine = declarator.ArgDetails.FileLine
	declarationSpecifiers.Package = ast.CXPackageIndex(pkg.Index)
	declarationSpecifiers.PreviouslyDeclared = true
	declSpecIdx := prgrm.AddCXArgInArray(declarationSpecifiers)

	localVarTypeSig := ast.GetCXTypeSignatureRepresentationOfCXArg(prgrm, prgrm.GetCXArgFromArray(declSpecIdx))
	localVarTypeSigIdx := prgrm.AddCXTypeSignatureInArray(localVarTypeSig)
	prgrm.CXAtomicOps[expressionIdx].AddOutput(prgrm, localVarTypeSigIdx)

	// Checking if something is supposed to be initialized
	// and if `initializer` actually contains something.
	// If `initializer` is `nil`, this means that an expression
	// equivalent to nil was used, such as `[]i32{}`.
	if doesInitialize && initializer != nil {
		processLocalInitialization(prgrm, pkg, &initializer, localVarTypeSigIdx, declSpecIdx)

		return append([]ast.CXExpression{*declCXLine, *decl}, initializer...)
	} else {
		// There is no initialization.
		expr := ast.MakeAtomicOperatorExpression(prgrm, nil)
		cxAtomicOpIdx := expr.Index
		prgrm.CXAtomicOps[cxAtomicOpIdx].Package = ast.CXPackageIndex(pkg.Index)
		prgrm.CXAtomicOps[cxAtomicOpIdx].AddOutput(prgrm, localVarTypeSigIdx)

		return []ast.CXExpression{*declCXLine, *expr}
	}

}
