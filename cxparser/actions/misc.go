package actions

import (
	"github.com/skycoin/cx/cx/ast"
	"github.com/skycoin/cx/cx/constants"
	"github.com/skycoin/cx/cx/types"
)

// hasDeclSpec determines if an argument has certain declaration specifier.
func hasDeclSpec(arg *ast.CXArgument, spec int) bool {
	found := false
	for _, s := range arg.DeclarationSpecifiers {
		if s == spec {
			found = true
		}
	}
	return found
}

// hasDerefOp determines if an argument has certain dereference operation
func hasDerefOp(arg *ast.CXArgument, spec int) bool {
	found := false
	for _, s := range arg.DereferenceOperations {
		if s == spec {
			found = true
		}
	}
	return found
}

// WritePrimary writes those bytes to prgrm.Data and returns a CXArg.
func WritePrimary(prgrm *ast.CXProgram, typeCode types.Code, byts []byte, isSlice bool) *ast.CXArgument {
	pkg, err := prgrm.GetCurrentPackage()
	if err != nil {
		panic(err)
	}

	arg := ast.MakeArgument("", CurrentFile, LineNo)
	arg.SetType(typeCode)
	arg.Package = ast.CXPackageIndex(pkg.Index)

	size := types.Cast_int_to_ptr(len(byts))

	arg.Size = typeCode.Size()
	arg.Offset = prgrm.Data.Size + prgrm.Data.StartsAt

	if arg.Type == types.STR || arg.Type == types.AFF {
		arg.Size = types.POINTER_SIZE
		if isSlice == false {

			types.Write_ptr(byts, 0, arg.Offset)
		}
	}

	// A CX program allocates min(INIT_HEAP_SIZE, MAX_HEAP_SIZE) bytes
	// after the stack segment. These bytes are used to allocate the data segment
	// at compile time. If the data segment is bigger than min(INIT_HEAP_SIZE, MAX_HEAP_SIZE),
	// we'll start appending the bytes to prgrm.Memory.
	// After compilation, we calculate how many bytes we need to add to have a heap segment
	// equal to `minHeapSize()` that is allocated after the data segment.
	memSize := types.Cast_int_to_ptr(len(prgrm.Memory))
	if (size + prgrm.Data.Size + prgrm.Data.StartsAt) > memSize {
		var i types.Pointer
		// First we need to fill the remaining free bytes in
		// the current `prgrm.Memory` slice.
		for i = types.Pointer(0); i < memSize-prgrm.Data.Size+prgrm.Data.StartsAt; i++ {
			prgrm.Memory[prgrm.Data.Size+prgrm.Data.StartsAt+i] = byts[i]
		}
		// Then we append the bytes that didn't fit.
		prgrm.Memory = append(prgrm.Memory, byts[i:]...)
	} else {
		for i, byt := range byts {
			prgrm.Memory[prgrm.Data.Size+prgrm.Data.StartsAt+types.Cast_int_to_ptr(i)] = byt
		}
	}
	prgrm.Data.Size += size
	prgrm.Heap.StartsAt = prgrm.Data.Size + prgrm.Data.StartsAt

	argIdx := prgrm.AddCXArgInArray(arg)
	argReturn := prgrm.GetCXArgFromArray(argIdx)

	return argReturn
}

// WritePrimary writes those bytes to prgrm.Data and returns a []CXExpression.
func WritePrimaryExprs(prgrm *ast.CXProgram, typeCode types.Code, byts []byte, isSlice bool) []ast.CXExpression {
	pkg, err := prgrm.GetCurrentPackage()
	if err != nil {
		panic(err)
	}

	arg := WritePrimary(prgrm, typeCode, byts, isSlice)

	expr := ast.MakeAtomicOperatorExpression(prgrm, nil)
	prgrm.CXAtomicOps[expr.Index].Package = ast.CXPackageIndex(pkg.Index)

	typeSig := ast.GetCXTypeSignatureRepresentationOfCXArg(prgrm, arg)
	typeSigIdx := prgrm.AddCXTypeSignatureInArray(typeSig)
	prgrm.CXAtomicOps[expr.Index].AddOutput(prgrm, typeSigIdx)

	return []ast.CXExpression{*expr}
}

func TotalLength(lengths []types.Pointer) types.Pointer {
	total := types.Pointer(1)
	for _, i := range lengths {
		total *= i
	}

	return total
}

func AffordanceStructs(prgrm *ast.CXProgram, pkg *ast.CXPackage, currentFile string, lineNo int) {
	// Argument type
	argStrct := ast.MakeStruct("Argument")
	// argStrct.Size = cxcore.GetArgSize(cxcore.TYPE_STR) + cxcore.GetArgSize(cxcore.TYPE_STR)

	argFldName := ast.MakeField("Name", types.STR, "", 0)
	argFldIndex := ast.MakeField("Index", types.I32, "", 0)
	argFldType := ast.MakeField("Type", types.STR, "", 0)

	argStrct.AddField(prgrm, argFldName)
	argStrct.AddField(prgrm, argFldIndex)
	argStrct.AddField(prgrm, argFldType)

	pkg.AddStruct(prgrm, argStrct)

	// Expression type
	exprStrct := ast.MakeStruct("Expression")
	// exprStrct.Size = cxcore.GetArgSize(cxcore.TYPE_STR)

	exprFldOperator := ast.MakeField("Operator", types.STR, "", 0)

	exprStrct.AddField(prgrm, exprFldOperator)

	pkg.AddStruct(prgrm, exprStrct)

	// Function type
	fnStrct := ast.MakeStruct("Function")
	// fnStrct.Size = cxcore.GetArgSize(cxcore.TYPE_STR) + cxcore.GetArgSize(cxcore.TYPE_STR) + cxcore.GetArgSize(cxcore.TYPE_STR)

	fnFldName := ast.MakeField("Name", types.STR, "", 0)

	fnFldInpSig := ast.MakeField("InputSignature", types.STR, "", 0)
	fnFldInpSig.Size = types.STR.Size()
	fnFldInpSig = DeclarationSpecifiers(fnFldInpSig, []types.Pointer{0}, constants.DECL_SLICE)

	fnFldOutSig := ast.MakeField("OutputSignature", types.STR, "", 0)
	fnFldOutSig.Size = types.STR.Size()
	fnFldOutSig = DeclarationSpecifiers(fnFldOutSig, []types.Pointer{0}, constants.DECL_SLICE)

	fnStrct.AddField(prgrm, fnFldName)
	fnStrct.AddField(prgrm, fnFldInpSig)

	fnStrct.AddField(prgrm, fnFldOutSig)

	pkg.AddStruct(prgrm, fnStrct)

	// Structure type
	strctStrct := ast.MakeStruct("Structure")
	// strctStrct.Size = cxcore.GetArgSize(cxcore.TYPE_STR)

	strctFldName := ast.MakeField("Name", types.STR, "", 0)

	strctStrct.AddField(prgrm, strctFldName)

	pkg.AddStruct(prgrm, strctStrct)

	// Package type
	pkgStrct := ast.MakeStruct("Structure")
	// pkgStrct.Size = cxcore.GetArgSize(cxcore.TYPE_STR)

	pkgFldName := ast.MakeField("Name", types.STR, "", 0)

	pkgStrct.AddField(prgrm, pkgFldName)

	pkg.AddStruct(prgrm, pkgStrct)

	// Caller type
	callStrct := ast.MakeStruct("Caller")
	// callStrct.Size = cxcore.GetArgSize(cxcore.TYPE_STR) + cxcore.GetArgSize(cxcore.TYPE_I32)

	callFldFnName := ast.MakeField("FnName", types.STR, "", 0)
	callFldFnSize := ast.MakeField("FnSize", types.I32, "", 0)

	callStrct.AddField(prgrm, callFldFnName)
	callStrct.AddField(prgrm, callFldFnSize)

	pkg.AddStruct(prgrm, callStrct)

	// Program type
	prgrmStrct := ast.MakeStruct("Program")
	// prgrmStrct.Size = cxcore.GetArgSize(cxcore.TYPE_I32) + cxcore.GetArgSize(cxcore.TYPE_I64)

	prgrmFldCallCounter := ast.MakeField("CallCounter", types.I32, "", 0)
	prgrmFldFreeHeap := ast.MakeField("HeapUsed", types.I64, "", 0)

	// prgrmFldCaller := cxcore.MakeField("Caller", cxcore.TYPE_STRUCT, "", 0)
	strctPkg, err := prgrm.GetPackageFromArray(callStrct.Package)
	if err != nil {
		panic(err)
	}
	prgrmFldCaller := DeclarationSpecifiersStruct(prgrm, callStrct.Name, strctPkg.Name, false, currentFile, lineNo)
	prgrmFldCaller.Name = "Caller"

	prgrmStrct.AddField(prgrm, prgrmFldCallCounter)
	prgrmStrct.AddField(prgrm, prgrmFldFreeHeap)
	prgrmStrct.AddField(prgrm, prgrmFldCaller)

	pkg.AddStruct(prgrm, prgrmStrct)
}

// PrimaryIdentifier creates an identifier expression with an output name of 'ident'.
func PrimaryIdentifier(prgrm *ast.CXProgram, ident string) []ast.CXExpression {
	pkg, err := prgrm.GetCurrentPackage()
	if err != nil {
		panic(err)
	}

	arg := ast.MakeArgument(ident, CurrentFile, LineNo) // fix: line numbers in errors sometimes report +1 or -1. Issue #195
	arg.SetType(types.IDENTIFIER)
	arg.Name = ident
	arg.Package = ast.CXPackageIndex(pkg.Index)

	expr := ast.MakeAtomicOperatorExpression(prgrm, nil)
	expression, err := prgrm.GetCXAtomicOp(expr.Index)
	if err != nil {
		panic(err)
	}

	typeSig := ast.GetCXTypeSignatureRepresentationOfCXArg(prgrm, arg)
	typeSigIdx := prgrm.AddCXTypeSignatureInArray(typeSig)
	expression.AddOutput(prgrm, typeSigIdx)
	expression.Package = ast.CXPackageIndex(pkg.Index)
	return []ast.CXExpression{*expr}
}

// IsAllArgsBasicTypes checks if all the input arguments in an expressions are of basic type.
func IsAllArgsBasicTypes(prgrm *ast.CXProgram, expr *ast.CXExpression) bool {
	expression, err := prgrm.GetCXAtomicOp(expr.Index)
	if err != nil {
		panic(err)
	}

	for _, inputIdx := range expression.GetInputs(prgrm) {
		input := prgrm.GetCXTypeSignatureFromArray(inputIdx)

		var inpType types.Code

		var inp *ast.CXArgument = &ast.CXArgument{}
		if input.Type == ast.TYPE_CXARGUMENT_DEPRECATE {
			inp = prgrm.GetCXArgFromArray(ast.CXArgumentIndex(input.Meta))

			inpType = inp.Type
			if inp.Type == types.POINTER {
				inpType = inp.PointerTargetType
			}
		} else {
			inpType = input.GetType(prgrm)
		}

		// TODO: Check why STR is considered as basic type.
		if !inpType.IsPrimitive() && inpType != types.STR {
			return false
		}
	}
	return true
}
