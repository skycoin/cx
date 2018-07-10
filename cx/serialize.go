package base

import (
	"fmt"
	"github.com/skycoin/skycoin/src/cipher/encoder"
)

/*
  Program
*/

type sIndex struct {
	ProgramOffset int32
	NamesOffset int32
	ValuesOffset int32
	PackagesOffset int32
	GlobalsOffset int32
	ImportsOffset int32
	FunctionsOffset int32
	StructsOffset int32
	FieldsOffset int32
	ExpressionsOffset int32
	ArgumentsOffset int32
	CallsOffset int32
}

type sProgram struct {
	PackagesOffset 			int32
	PackagesSize 			int32
	CurrentPackageOffset 	int32
	
	InputsOffset 	int32
	InputsSize 		int32

	OutputsOffset 	int32
	OutputsSize 	int32
	
	CallStackOffset int32
	CallStackSize 	int32

	CallCounter 	int32

	StacksOffset 	int32
	StacksSize 		int32

	HeapOffset 		int32
	DataOffset 		int32

	Terminated 		int32

	PathOffset 		int32
	PathSize 		int32

	StepsOffset 	int32
	StepsSize 		int32
}

type sHeap struct {
	HeapOffset 		int32
	HeapSize		int32
	HeapPointer		int32
	ProgramOffset 	int32
}

type sStack struct {
	StackOffset 	int32
	StackSize 		int32
	StackPointer 	int32
	ProgramOffset 	int32
}

type sCall struct {
	OperatorOffset 		int32
	Line 				int32
	FramePointer 		int32

	StateOffset 		int32
	StateSize 			int32
	ReturnAddressOffset int32
	
	PackageOffset 		int32
	ProgramOffset 		int32
}

/*
  Packages
*/

type sPackage struct {
	NameOffset 			int32
	NameSize 			int32
	ImportsOffset 		int32
	ImportsSize 		int32
	FunctionsOffset 	int32
	FunctionsSize 		int32
	StructsOffset 		int32
	StructsSize 		int32
	GlobalsOffset 		int32
	GlobalsSize 		int32

	CurrentFunctionOffset 	int32
	CurrentStructOffset 	int32
	ProgramOffset 			int32
}

type sGlobal struct {
	NameOffset int32
	NameSize int32
	TypeOffset int32
	TypeSize int32
	ValueOffset int32
	ValueSize int32

	PackageOffset int32
}

/*
  Structs
*/

type sStruct struct {
	NameOffset 		int32
	NameSize 		int32
	FieldsOffset 	int32
	FieldsSize 		int32

	Size 			int32

	PackageOffset 	int32
	ProgramOffset 	int32
}

// type sField struct {
// 	NameOffset int32
// 	NameSize int32
// 	TypeOffset int32
// 	TypeSize int32
// }

// type sType struct {
// 	NameOffset int32
// 	NameSize int32
// }

/*
  Functions
*/

type sFunction struct {
	NameOffset 			int32
	NameSize 			int32
	InputsOffset 		int32
	InputsSize 			int32
	OutputsOffset 		int32
	OutputsSize 		int32
	ExpressionsOffset	int32
	ExpressionsSize 	int32
	Size 				int32
	Length 				int32

	ListOfPointersOffset 	int32
	ListOfPointersSize 		int32

	IsNative 	int32
	OpCode 		int32

	CurrentExpressionOffset int32
	PackageOffset 			int32
	ProgramOffset			int32
}

// type sParameter struct {
// 	NameOffset int32
// 	NameSize int32
// 	TypeOffset int32
// 	TypeSize int32
// }

type sExpression struct {
	OperatorOffset 		int32
	InputsOffset 		int32
	InputsSize 			int32
	OutputsOffset 	int32
	OutputsSize 	int32
	
	Line 			int32
	FileLine 		int32
	FileNameOffset 	int32
	FileNameSize 	int32
	
	LabelOffset int32
	LabelSize 	int32
	ThenLines 	int32
	ElseLines 	int32

	IsStructLiteral int32
	IsArrayLiteral 	int32

	FunctionOffset 	int32
	PackageOffset 	int32
	ProgramOffset 	int32
}

type sConstant struct {
	Type 		int32
	ValueOffset int32
	ValueSize 	int32
}

type sArgument struct {
	NameOffset 			int32
	NameSize 			int32
	TypeOffset 			int32
	CustomTypeOffset 	int32
	Size 				int32
	TotalSize 			int32
	PointeeSize 		int32

	MemoryRead 	int32
	MemoryWrite int32
	Offset 		int32
	HeapOffset 	int32
	
	IndirectionLevels 			int32
	DereferenceLevels 			int32
	PointeeOffset 				int32
	PointeeMemoryType 			int32
	DereferenceOperationsOffset int32
	DereferenceOperationsSize 	int32
	DereferenceSpecifiersOffset int32
	DereferenceSpecifiersSize 	int32

	IsArray 			int32
	IsArrayFirst 		int32
	IsPointer 			int32
	IsReference 		int32
	IsDeferenceFirst 	int32
	IsStruct 			int32
	IsField 			int32
	IsRest 				int32
	IsLocalDeclaration 	int32

	LengthsOffset 	int32
	LengthsSize 	int32
	IndexesOffset 	int32
	IndexesSize		int32
	FieldsOffset	int32
	FieldsSize		int32

	SynonymousToOffset 	int32
	SynonymousToSize 	int32

	PackageOffset int32
	ProgramOffset int32

	ValueOffset int32
	ValueSize 	int32
	TypOffset 	int32
	TypSize		int32
}

/*
  Affordances

  Affordances must not be serialized
*/

func serializeName (name string, sNamesMap *map[string]int, sNames *[]byte, sNamesCounter *int) (offset, size int32) {
	if off, ok := (*sNamesMap)[name]; ok {
		offset = int32(off)
		size = int32(encoder.Size(name))
	} else {
		offset = int32(*sNamesCounter)
		size = int32(encoder.Size(name))
		*sNames = append(*sNames, encoder.Serialize(name)...)
		(*sNamesMap)[name] = *sNamesCounter
		*sNamesCounter = *sNamesCounter + int(size)
	}
	return offset, size
}

func serializeValue (value, sValues *[]byte, sValuesCounter *int) (offset, size int32) {
	*sValues = append(*sValues, encoder.Serialize(*value)...)
	offset = int32(*sValuesCounter)
	size = int32(encoder.Size(*value))
	*sValuesCounter = *sValuesCounter + int(size)

	return offset, size
}

func serializeImports (imps []*CXPackage, sPacksMap *map[string]int, sImps *[]byte, sImpsCounter *int) (offset, size int32) {
	if imps != nil && len(imps) > 0 {
		offset = int32(*sImpsCounter)
		size = int32(len(imps))

		for _, imp := range imps {
			// we only need the index of the imported module
			*sImps = append(*sImps, encoder.SerializeAtomic(int32((*sPacksMap)[imp.Name]))...)
			*sImpsCounter++
		}
	}
	return offset, size
}

func Serialize (cxt *CXProgram) *[]byte {
	serialized := make([]byte, 0)

	sNames := make([]byte, 0)
	sNamesCounter := 0
	sNamesMap := make(map[string]int, 0)

	sValues := make([]byte, 0)
	sValuesCounter := 0

	sPacks := make([]byte, 0)
	sPacksCounter := 0
	sPacksMap := make(map[string]int, 0)

	sGlobals := make([]byte, 0)
	sGlobalsCounter := 0

	sImps := make([]byte, 0)
	sImpsCounter := 0

	sFns := make([]byte, 0)
	sFnsCounter := 0
	sFnsMap := make([]string, 0)

	sStrcts := make([]byte, 0)
	sStrctsCounter := 0

	sFlds := make([]byte, 0)
	sFldsCounter := 0

	sExprs := make([]byte, 0)
	sExprsCounter := 0

	sArgs := make([]byte, 0)
	sArgsCounter := 0

	sCalls := make([]byte, 0)
	sCallsCounter := 0

	sStacks := make([]byte, 0)
	sStacksCounter := 0
	
	sHeap := make([]byte, 0)
	sHeapCounter := 0

	sData := make([]byte, 0)
	sDataCOunter := 0

	// Program Serialize

	sPrgrm := &sProgram{}
	sPrgrm.PackagesOffset = int32(sPacksCounter)
	sPrgrm.PackagesSize = int32(len(cxt.Packages))

	cxtPackages := make([]*CXPackage, 0)
	packFunctions := make([][]*CXFunction, len(cxt.Packages))
	packCounter := 0

	// Program Packages
	for _, pack := range cxt.Packages {
		sPacksMap[pack.Name] = sPacksCounter
		cxtPackages = append(cxtPackages, pack)
		sPacksCounter++

		//Functions
		for _, fn := range pack.Functions {
			fnName := fmt.Sprintf("%s.%s", pack.Name, fn.Name)
			sFnsMap = append(sFnsMap, fnName)
			packFunctions[packCounter] = append(packFunctions[packCounter], fn)
		}
		packCounter++
	}

	sPacksCounter = 0

	// Packages Serialize
	for i, pack := range cxtPackages {
		sPack := sPackage{}

		// Pack's Name
		sPack.NameOffset, sPack.NameSize = serializeName(pack.Name, &sNamesMap, &sNames, &sNamesCounter)

		// Serialize Pack's Imports
		sPack.ImportsOffset, sPack.ImportsSize = serializeImports(pack.Imports, &sPacksMap, &sImps, &sImpsCounter)

		// Serialize Pack's Functions
		if packFunctions[i] != nil && len(packFunctions[i]) > 0{
			sPack.FunctionsOffset = int32(sFnsCounter)
			sPack.FunctionsSize = int32(len(pack.Functions))

			for _, fn := range packFunctions[i] {
				sFn := sFunction{}

				// Function's Name
				sFn.NameOffset, sFn.NameSize = serializeName(fn.Name, &sNamesMap, &sNames, &sNamesCounter)

				// Serialize Function's Inputs
				if fn.Inputs != nil && len(fn.Inputs) > 0 { 
					sFn.InputsOffset = int32(sArgsCounter)
					sFn.InputsSize = int32(len(fn.Inputs))

					for _, arg := range fn.Inputs {
						sInput := sArgument{}

						// Input's Name
						sInput.NameOffset, sInput.NameSize = serializeName(arg.Name, &sNamesMap, &sNames, &sNamesCounter)

						// Input Type
						sInput.TypOffset, sInput.TypSize = serializeName(arg.Typ, &sNamesMap, &sNames, &sNamesCounter)
						sInput.TypeOffset = int32(arg.Type)

						// Input values
						sInput.ValuesOffset, sInput.ValueSize = serializeValue(arg.Value, &sValues, &sValuesCounter)

						// Others Options
						sInput.Size = int32(arg.Size)
						sInput.TotalSize = int32(arg.Size)
						sInput.PointeeSize = int32(arg.PointeeSize)
						sInput.MemoryRead = int32(arg.MemoryRead)
						sInput.MemoryWrite = int32(arg.MemoryWrite)
						sInput.Offset = int32(arg.Offset)
						sInput.HeapOffset = int32(arg.HeapOffset)

						if arg.IsArray {
							sInput.IsArray = int32(1)
						} else {
							sInput.IsArray = int32(0)
						}

						if arg.IsArrayFirst {
							sInput.IsArrayFirst = int32(1)
						} else {
							sInput.IsArrayFirst = int32(0)
						}

						if arg.IsPointer {
							sInput.IsPointer = int32(1)
						} else {
							sInput.IsPointer = int32(0)
						}

						if arg.IsReference {
							sInput.IsReference = int32(1)
						} else {
							sInput.IsReference = int32(0)
						}

						if arg.IsDereferenceFirst {
							sInput.IsDereferenceFirst = int32(1)
						} else {
							sInput.IsDereferenceFirst = int32(0)
						}

						if arg.IsStruct {
							sInput.IsStruct = int32(1)
							

						} else {
							sInput.IsStruct = int32(0)
						}

						if arg.IsField {
							sInput.IsField = int32(1)
						} else {
							sInput.IsField = int32(0)
						}

						if arg.IsRest {
							sInput.IsRest = int32(1)
						} else {
							sInput.IsRest = int32(0)
						}

						if arg.IsLocalDeclaration {
							sInput.IsLocalDeclaration = int32(1)
						} else {
							sInput.IsLocalDeclaration = int32(0)
						}

						// Synonymous To...
						sInput.SynonymousToOffset, sInput.SynonymousToSize = serializeName(arg.SynonymousTo, &sNamesMap, &sNames, &sNamesCounter)

						// Package Offset
						sInput.PackageOffset = int32(sPacksMap[fn.Package.Name])

						// Saving Inputs
						sArgs = append(sArgs, encoder.Serialize(sInput)...)
						sArgsCounter++
					}
				} else {
					sFn.InputsOffset = -1
					sFn.InputsSize = -1
				}

				// Serialize Function's Outputs
				if fn.Outputs != nil && len(fn.Outputs) > 0 { 
					sFn.OutputsOffset = int32(sArgsCounter)
					sFn.OutputsSize = int32(len(fn.Outputs))

					for _, arg := range fn.Outputs {
						sOutput := sArgument{}

						// Output's Name
						sOutput.NameOffset, sOutput.NameSize = serializeName(arg.Name, &sNamesMap, &sNames, &sNamesCounter)

						// Output Type
						sOutput.TypOffset, sOutput.TypSize = serializeName(arg.Typ, &sNamesMap, &sNames, &sNamesCounter)
						sOutput.TypeOffset = int32(arg.Type)
						
						// Output values
						sOutput.ValuesOffset, sOutput.ValueSize = serializeValue(arg.Value, &sValues, &sValuesCounter)

						// Others Options
						sOutput.Size = int32(arg.Size)
						sOutput.TotalSize = int32(arg.Size)
						sOutput.PointeeSize = int32(arg.PointeeSize)
						sOutput.MemoryRead = int32(arg.MemoryRead)
						sOutput.MemoryWrite = int32(arg.MemoryWrite)
						sOutput.Offset = int32(arg.Offset)
						sOutput.HeapOffset = int32(arg.HeapOffset)

						if arg.IsArray {
							sOutput.IsArray = int32(1)
						} else {
							sOutput.IsArray = int32(0)
						}

						if arg.IsArrayFirst {
							sOutput.IsArrayFirst = int32(1)
						} else {
							sOutput.IsArrayFirst = int32(0)
						}

						if arg.IsPointer {
							sOutput.IsPointer = int32(1)
						} else {
							sOutput.IsPointer = int32(0)
						}

						if arg.IsReference {
							sOutput.IsReference = int32(1)
						} else {
							sOutput.IsReference = int32(0)
						}

						if arg.IsDereferenceFirst {
							sOutput.IsDereferenceFirst = int32(1)
						} else {
							sOutput.IsDereferenceFirst = int32(0)
						}

						if arg.IsStruct {
							sOutput.IsStruct = int32(1)
							//CustomType

						} else {
							sOutput.IsStruct = int32(0)
						}

						if arg.IsField {
							sOutput.IsField = int32(1)
						} else {
							sOutput.IsField = int32(0)
						}

						if arg.IsRest {
							sOutput.IsRest = int32(1)
						} else {
							sOutput.IsRest = int32(0)
						}

						if arg.IsLocalDeclaration {
							sOutput.IsLocalDeclaration = int32(1)
						} else {
							sOutput.IsLocalDeclaration = int32(0)
						}

						// Synonymous To...
						sOutput.SynonymousToOffset, sOutput.SynonymousToSize = serializeName(arg.SynonymousTo, &sNamesMap, &sNames, &sNamesCounter)

						// Package Offset
						sOutput.PackageOffset = int32(sPacksMap[fn.Package.Name])

						// Saving Inputs
						sArgs = append(sArgs, encoder.Serialize(sOutput)...)
						sArgsCounter++
					}
				} else {
					sFn.OutputsOffset = -1
					sFn.OutputsSize = -1
				}

				// Serialize Function's Expressions
				if fn.Expressions != nil && len(fn.Expressions) > 0 {
					sFn.ExpressionsOffset = int32(sExprsCounter)
					sFn.ExpressionsSize = int32(len(fn.Expressions))

					for _, expr := range fn.Expressions {
						sExpr := sExpression{}

						opName := fmt.Sprintf("%s.%s", expr.Operator.Package.Name, expr.Operator.Name)
						opOffset := -1
						for i, fn := range sFnsMap {
							if opName == fn {
								opOffset = i
								break
							}
						}

						//OperatorOffset
						if opOffset >= 0 {
							sExpr.OperatorOffset = int32(opOffset)
						} else {
							panic(fmt.Sprintf("Expression's operator (%s) not found in sFnsMap", opName))
						}

						// Expression's Inputs
						if expr.Inputs != nil && len(expr.Inputs) > 0 { 
							sExpr.InputsOffset = int32(sArgsCounter)
							sExpr.InputsSize = int32(len(expr.Inputs))
		
							for _, arg := range expr.Inputs {
								sInput := sArgument{}

								// Input's Name
								sInput.NameOffset, sInput.NameSize = serializeName(arg.Name, &sNamesMap, &sNames, &sNamesCounter)

								// Input Type
								sInput.TypOffset, sInput.TypSize = serializeName(arg.Typ, &sNamesMap, &sNames, &sNamesCounter)
								sInput.TypeOffset = int32(arg.Type)

								// Input values
								sInput.ValuesOffset, sInput.ValueSize = serializeValue(arg.Value, &sValues, &sValuesCounter)

								// Others Options
								sInput.Size = int32(arg.Size)
								sInput.TotalSize = int32(arg.Size)
								sInput.PointeeSize = int32(arg.PointeeSize)
								sInput.MemoryRead = int32(arg.MemoryRead)
								sInput.MemoryWrite = int32(arg.MemoryWrite)
								sInput.Offset = int32(arg.Offset)
								sInput.HeapOffset = int32(arg.HeapOffset)

								if arg.IsArray {
									sInput.IsArray = int32(1)
								} else {
									sInput.IsArray = int32(0)
								}

								if arg.IsArrayFirst {
									sInput.IsArrayFirst = int32(1)
								} else {
									sInput.IsArrayFirst = int32(0)
								}

								if arg.IsPointer {
									sInput.IsPointer = int32(1)
								} else {
									sInput.IsPointer = int32(0)
								}

								if arg.IsReference {
									sInput.IsReference = int32(1)
								} else {
									sInput.IsReference = int32(0)
								}

								if arg.IsDereferenceFirst {
									sInput.IsDereferenceFirst = int32(1)
								} else {
									sInput.IsDereferenceFirst = int32(0)
								}

								if arg.IsStruct {
									sInput.IsStruct = int32(1)
									//CustomType

								} else {
									sInput.IsStruct = int32(0)
								}

								if arg.IsField {
									sInput.IsField = int32(1)
								} else {
									sInput.IsField = int32(0)
								}

								if arg.IsRest {
									sInput.IsRest = int32(1)
								} else {
									sInput.IsRest = int32(0)
								}

								if arg.IsLocalDeclaration {
									sInput.IsLocalDeclaration = int32(1)
								} else {
									sInput.IsLocalDeclaration = int32(0)
								}

								// Synonymous To...
								sInput.SynonymousToOffset, sInput.SynonymousToSize = serializeName(arg.SynonymousTo, &sNamesMap, &sNames, &sNamesCounter)

								// Package Offset
								sInput.PackageOffset = int32(sPacksMap[fn.Package.Name])

								// Saving Inputs
								sArgs = append(sArgs, encoder.Serialize(sInput)...)
								sArgsCounter++
							}
						} else {
							sExpr.InputsOffset = -1
							sExpr.InputsSize = -1
						}

						// Expression's Outputs
						if expr.Outputs != nil && len(expr.Outputs) > 0 { 
							sExpr.OutputsOffset = int32(sArgsCounter)
							sExpr.OutputsSize = int32(len(expr.Outputs))
		
							for _, arg := range expr.Outputs {
								sOutput := sArgument{}

								// Output's Name
								sOutput.NameOffset, sOutput.NameSize = serializeName(arg.Name, &sNamesMap, &sNames, &sNamesCounter)

								// Output Type
								sOutput.TypOffset, sOutput.TypSize = serializeName(arg.Typ, &sNamesMap, &sNames, &sNamesCounter)
								sOutput.TypeOffset = int32(arg.Type)
								
								// Output values
								sOutput.ValuesOffset, sOutput.ValueSize = serializeValue(arg.Value, &sValues, &sValuesCounter)

								// Others Options
								sOutput.Size = int32(arg.Size)
								sOutput.TotalSize = int32(arg.Size)
								sOutput.PointeeSize = int32(arg.PointeeSize)
								sOutput.MemoryRead = int32(arg.MemoryRead)
								sOutput.MemoryWrite = int32(arg.MemoryWrite)
								sOutput.Offset = int32(arg.Offset)
								sOutput.HeapOffset = int32(arg.HeapOffset)

								if arg.IsArray {
									sOutput.IsArray = int32(1)
								} else {
									sOutput.IsArray = int32(0)
								}

								if arg.IsArrayFirst {
									sOutput.IsArrayFirst = int32(1)
								} else {
									sOutput.IsArrayFirst = int32(0)
								}

								if arg.IsPointer {
									sOutput.IsPointer = int32(1)
								} else {
									sOutput.IsPointer = int32(0)
								}

								if arg.IsReference {
									sOutput.IsReference = int32(1)
								} else {
									sOutput.IsReference = int32(0)
								}

								if arg.IsDereferenceFirst {
									sOutput.IsDereferenceFirst = int32(1)
								} else {
									sOutput.IsDereferenceFirst = int32(0)
								}

								if arg.IsStruct {
									sOutput.IsStruct = int32(1)
									//CustomType

								} else {
									sOutput.IsStruct = int32(0)
								}

								if arg.IsField {
									sOutput.IsField = int32(1)
								} else {
									sOutput.IsField = int32(0)
								}

								if arg.IsRest {
									sOutput.IsRest = int32(1)
								} else {
									sOutput.IsRest = int32(0)
								}

								if arg.IsLocalDeclaration {
									sOutput.IsLocalDeclaration = int32(1)
								} else {
									sOutput.IsLocalDeclaration = int32(0)
								}

								// Synonymous To...
								sOutput.SynonymousToOffset, sOutput.SynonymousToSize = serializeName(arg.SynonymousTo, &sNamesMap, &sNames, &sNamesCounter)

								// Package Offset
								sOutput.PackageOffset = int32(sPacksMap[fn.Package.Name])

								// Saving Inputs
								sArgs = append(sArgs, encoder.Serialize(sOutput)...)
								sArgsCounter++
							}
						} else {
							sExpr.OutputsOffset = -1
							sExpr.OutputsSize = -1
						}

						// Expression's FileLine and FileName
						sExpr.Line = int32(expr.Line)
						sExpr.FileLine = int32(expr.FileLine)

						sExpr.FileNameOffset, sExpr.FileNameSize = serializeName(expr.FileName, &sNamesMap, &sNames, &sNamesCounter)

						// Label Then Else

						sExpr.LabelOffset, sExpr.LabelSize = serializeName(expr.Label, &sNamesMap, &sNames, &sNamesCounter)
						sExpr.ThenLines = int32(expr.ThenLines)
						sExpr.ElseLines = int32(expr.ElseLines)


						// IsStrucLiteral and IsArrayLiteral flags
						if expr.IsStructLiteral {
							sExpr.IsStructLiteral = int32(1)
						} else {
							sExpr.IsStructLiteral = int32(0)
						}

						if expr.IsArrayLiteral {
							sExpr.IsArrayLiteral = int32(1)
						} else {
							sExpr.IsArrayLiteral = int32(0)
						}

						// Expression's Function
						fnOffset := 0
						for i, fnName := range sFnsMap {
							if fnName == fmt.Sprintf("%s.%s", pack.Name, fn.Name) {
								fnOffset = i
								break
							}
						}

						if fnOffset >= 0 {
							sExpr.FunctionOffset = int32(fnOffset)
						} else {
							panic(fmt.Sprintf("Function '%s' not found in sFnsMap", fn.Name))
						}

						// Expression's Package
						sExpr.PackageOffset = int32(sPacksMap[expr.Package.Name])

						// Saving Expression
						sExprs = append(sExprs, encoder.Serialize(sExpr)...)

						if fn.CurrentExpression == expr {
							sFn.CurrentExpressionOffset = int32(sExprsCounter)
						}
						sExprsCounter++
					}
				}

				// List of Pointers
				if fn.ListOfPointers != nil && len(fn.ListOfPointers) > 0 { 
					sFn.ListOfPointersOffset = int32(sArgsCounter)
					sFn.ListOfPointersSize = int32(len(fn.Outputs))

					for _, arg := range fn.ListOfPointers {
						sInput := sArgument{}

						// Input's Name
						sInput.NameOffset, sInput.NameSize = serializeName(arg.Name, &sNamesMap, &sNames, &sNamesCounter)

						// Input Type
						sInput.TypOffset, sInput.TypSize = serializeName(arg.Typ, &sNamesMap, &sNames, &sNamesCounter)
						sInput.TypeOffset = int32(arg.Type)

						// Input values
						sInput.ValuesOffset, sInput.ValueSize = serializeValue(arg.Value, &sValues, &sValuesCounter)

						// Others Options
						sInput.Size = int32(arg.Size)
						sInput.TotalSize = int32(arg.Size)
						sInput.PointeeSize = int32(arg.PointeeSize)
						sInput.MemoryRead = int32(arg.MemoryRead)
						sInput.MemoryWrite = int32(arg.MemoryWrite)
						sInput.Offset = int32(arg.Offset)
						sInput.HeapOffset = int32(arg.HeapOffset)

						if arg.IsArray {
							sInput.IsArray = int32(1)
						} else {
							sInput.IsArray = int32(0)
						}

						if arg.IsArrayFirst {
							sInput.IsArrayFirst = int32(1)
						} else {
							sInput.IsArrayFirst = int32(0)
						}

						if arg.IsPointer {
							sInput.IsPointer = int32(1)
						} else {
							sInput.IsPointer = int32(0)
						}

						if arg.IsReference {
							sInput.IsReference = int32(1)
						} else {
							sInput.IsReference = int32(0)
						}

						if arg.IsDereferenceFirst {
							sInput.IsDereferenceFirst = int32(1)
						} else {
							sInput.IsDereferenceFirst = int32(0)
						}

						if arg.IsStruct {
							sInput.IsStruct = int32(1)
							//CustomType

						} else {
							sInput.IsStruct = int32(0)
						}

						if arg.IsField {
							sInput.IsField = int32(1)
						} else {
							sInput.IsField = int32(0)
						}

						if arg.IsRest {
							sInput.IsRest = int32(1)
						} else {
							sInput.IsRest = int32(0)
						}

						if arg.IsLocalDeclaration {
							sInput.IsLocalDeclaration = int32(1)
						} else {
							sInput.IsLocalDeclaration = int32(0)
						}

						// Synonymous To...
						sInput.SynonymousToOffset, sInput.SynonymousToSize = serializeName(arg.SynonymousTo, &sNamesMap, &sNames, &sNamesCounter)

						// Package Offset
						sInput.PackageOffset = int32(sPacksMap[fn.Package.Name])

						// Saving Inputs
						sArgs = append(sArgs, encoder.Serialize(sInput)...)
						sArgsCounter++
					}
				} else {
					sFn.ListOfPointersOffset = -1
					sFn.ListOfPointersSize = -1
				}

				sFn.Size = int32(fn.Size)
				sFn.Length = int32(fn.Length)

				if fn.IsNative {
					sFn.IsNative = int32(1)
				} else {
					sFn.IsNative = int32(0)
				}

				sFn.OpCode = int32(fn.OpCode)

				sFn.PackageOffset = int32(sPacksMap[fn.Package.Name])

				if pack.CurrentFunction == fn {
					sPack.CurrentFunctionOffset = int32(sFnsCounter)
				}

				sFns = append(sFns, encoder.Serialize(sFn)...)
				sFnsCounter++
			}
		}

		// Serialize Pack's Strucs
		if pack.Structs != nil && len(pack.Structs) > 0 {
			sPack.StructsOffset = int32(sStrctsCounter)
			sPack.StructsSize = int32(len(pack.Structs))

			for _, strct := range pack.Structs {
				sStrct := sStruct{}
				
				sStrct.NameOffset, sStrct.NameSize = serializeName(strct.Name, &sNamesMap, &sNames, &sNamesCounter)

				if strct.Fields != nil && len(strct.Fields) > 0 {
					sStrct.FieldsOffset = int32(sFldsCounter)
					sStrct.FieldsSize = int32(len(strct.Fields))

					for _, fld := range strct.Fields {
						sFld := sArgument{}

						sFld.NameOffset, sFld.NameSize = serializeName(fld.Name, &sNamesMap, &sNames, &sNamesCounter)

						sFld.TypOffset, sFld.TypSize = serializeName(fld.Typ, &sNamesMap, &sNames, &sNamesCounter)

						sFld.TypeOffset = int32(fld.Type)

						sFlds = append(sFlds, encoder.Serialize(sFld)...)
						sFldsCounter++
					}
				}

				sStrct.Size = int32(strct.Size)

				sStrct.PackageOffset = int32(sPacksMap[strct.Package.Name])

				if pack.CurrentStruct == strct {
					sPack.CurrentStructOffset = int32(sStrctsCounter)
				}
				sStrcts = append(sStrcts, encoder.Serialize(sStrct)...)
				sStrctsCounter++
			}
		} else {
			sPack.CurrentStructOffset = int32(-1)
		}

		// Serialize Pack's Globals
		if pack.Globals != nil && len(pack.Globals) > 0 {
			sPack.GlobalsOffset = int32(sGlobalsCounter)
			sPack.GlobalsSize = int32(len(pack.Globals))
			for _, arg := range pack.Globals {
				sInput := sArgument{}

				// Input's Name
				sInput.NameOffset, sInput.NameSize = serializeName(arg.Name, &sNamesMap, &sNames, &sNamesCounter)

				// Input Type
				sInput.TypOffset, sInput.TypSize = serializeName(arg.Typ, &sNamesMap, &sNames, &sNamesCounter)
				sInput.TypeOffset = int32(arg.Type)

				// Input values
				sInput.ValuesOffset, sInput.ValueSize = serializeValue(arg.Value, &sValues, &sValuesCounter)

				// Others Options
				sInput.Size = int32(arg.Size)
				sInput.TotalSize = int32(arg.Size)
				sInput.PointeeSize = int32(arg.PointeeSize)
				sInput.MemoryRead = int32(arg.MemoryRead)
				sInput.MemoryWrite = int32(arg.MemoryWrite)
				sInput.Offset = int32(arg.Offset)
				sInput.HeapOffset = int32(arg.HeapOffset)

				if arg.IsArray {
					sInput.IsArray = int32(1)
				} else {
					sInput.IsArray = int32(0)
				}

				if arg.IsArrayFirst {
					sInput.IsArrayFirst = int32(1)
				} else {
					sInput.IsArrayFirst = int32(0)
				}

				if arg.IsPointer {
					sInput.IsPointer = int32(1)
				} else {
					sInput.IsPointer = int32(0)
				}

				if arg.IsReference {
					sInput.IsReference = int32(1)
				} else {
					sInput.IsReference = int32(0)
				}

				if arg.IsDereferenceFirst {
					sInput.IsDereferenceFirst = int32(1)
				} else {
					sInput.IsDereferenceFirst = int32(0)
				}

				if arg.IsStruct {
					sInput.IsStruct = int32(1)
					//CustomType

				} else {
					sInput.IsStruct = int32(0)
				}

				if arg.IsField {
					sInput.IsField = int32(1)
				} else {
					sInput.IsField = int32(0)
				}

				if arg.IsRest {
					sInput.IsRest = int32(1)
				} else {
					sInput.IsRest = int32(0)
				}

				if arg.IsLocalDeclaration {
					sInput.IsLocalDeclaration = int32(1)
				} else {
					sInput.IsLocalDeclaration = int32(0)
				}

				// Synonymous To...
				sInput.SynonymousToOffset, sInput.SynonymousToSize = serializeName(arg.SynonymousTo, &sNamesMap, &sNames, &sNamesCounter)

				// Package Offset
				sInput.PackageOffset = int32(sPacksMap[fn.Package.Name])

				// Saving Inputs
				sArgs = append(sArgs, encoder.Serialize(sInput)...)
				sArgsCounter++
			}
		}

		if cxt.CurrentPackage == pack {
			sPrgrm.CurrentPackageOffset = int32(sPacksCounter)
		}
		
		sPacks = append(sPacks, encoder.Serialize(sPack)...)
		sPacksCounter++
	}

	if cxt.Terminated {
		sPrgrm.Terminated = int32(1)
	} else {
		sPrgrm.Terminated = int32(0)
	}

	// Program's CallStack
	sPrgrm.CallStackOffset = int32(sCallsCounter)
	sPrgrm.CallStackSize = int32(len(cxt.CallStack))
	lastCallOffset := int32(-1)
	
	for _, call := range cxt.CallStack{
		sCll := sCall{}
		
		opName := fmt.Sprintf("%s.%s", call.Operator.Package.Name, call.Operator.Name)
		opOffset := -1
		for i, fn := range sFnsMap {
			if opName == fn {
				opOffset = i
				break
			}
		}

		if opOffset >= 0 {
			sCll.OperatorOffset = int32(opOffset)
		} else {
			panic(fmt.Sprintf("Expression's operator (%s) not found in sFnsMap", opName))
		}

		sCll.Line = int32(call.Line)
		sCll.FramePointer = int32(call.FramePointer)

		if call.State != nil && len(call.State) > 0 {
			sCll.StateOffset = int32(sArgsCounter)
			sCll.StateSize = int32(len(call.State))

			for _, arg := range call.State {
				sInput := sArgument{}

				// Input's Name
				sInput.NameOffset, sInput.NameSize = serializeName(arg.Name, &sNamesMap, &sNames, &sNamesCounter)

				// Input Type
				sInput.TypOffset, sInput.TypSize = serializeName(arg.Typ, &sNamesMap, &sNames, &sNamesCounter)
				sInput.TypeOffset = int32(arg.Type)

				// Input values
				sInput.ValuesOffset, sInput.ValueSize = serializeValue(arg.Value, &sValues, &sValuesCounter)

				// Others Options
				sInput.Size = int32(arg.Size)
				sInput.TotalSize = int32(arg.Size)
				sInput.PointeeSize = int32(arg.PointeeSize)
				sInput.MemoryRead = int32(arg.MemoryRead)
				sInput.MemoryWrite = int32(arg.MemoryWrite)
				sInput.Offset = int32(arg.Offset)
				sInput.HeapOffset = int32(arg.HeapOffset)

				if arg.IsArray {
					sInput.IsArray = int32(1)
				} else {
					sInput.IsArray = int32(0)
				}

				if arg.IsArrayFirst {
					sInput.IsArrayFirst = int32(1)
				} else {
					sInput.IsArrayFirst = int32(0)
				}

				if arg.IsPointer {
					sInput.IsPointer = int32(1)
				} else {
					sInput.IsPointer = int32(0)
				}

				if arg.IsReference {
					sInput.IsReference = int32(1)
				} else {
					sInput.IsReference = int32(0)
				}

				if arg.IsDereferenceFirst {
					sInput.IsDereferenceFirst = int32(1)
				} else {
					sInput.IsDereferenceFirst = int32(0)
				}

				if arg.IsStruct {
					sInput.IsStruct = int32(1)
					//CustomType

				} else {
					sInput.IsStruct = int32(0)
				}

				if arg.IsField {
					sInput.IsField = int32(1)
				} else {
					sInput.IsField = int32(0)
				}

				if arg.IsRest {
					sInput.IsRest = int32(1)
				} else {
					sInput.IsRest = int32(0)
				}

				if arg.IsLocalDeclaration {
					sInput.IsLocalDeclaration = int32(1)
				} else {
					sInput.IsLocalDeclaration = int32(0)
				}

				// Synonymous To...
				sInput.SynonymousToOffset, sInput.SynonymousToSize = serializeName(arg.SynonymousTo, &sNamesMap, &sNames, &sNamesCounter)

				// Package Offset
				sInput.PackageOffset = int32(sPacksMap[arg.Package.Name])

				// Saving Inputs
				sArgs = append(sArgs, encoder.Serialize(sInput)...)
				sArgsCounter++
			}
		}

		if lastCallOffset >= 0 {
			sCll.ReturnAddressOffset = lastCallOffset
		} else {
			sCll.ReturnAddressOffset = int32(-1) // nil
		}

		sCll.PackageOffset = int32(sPacksMap[call.Package.Name])

		sCalls = append(sCalls, encoder.Serialize(sCll)...)
		lastCallOffset = int32(sCallsCounter)
		sCallsCounter++
	}

	// Program's Stacks
	sPrgrm.StacksOffset = int32(sStacksCounter)
	sPrgrm.StacksSize = int32(len(cxt.Stacks))
	lastStackOffset := int32(-1)

	for i, stack := range cxt.Stacks {
		sStck := sStack{}

		sStck.StackOffset = int32(i)
		sStck.StackSize = int32(len(stack.Stack))

		sStck.StackPointer = int32(stack.StackPointer)

		lastStackOffset := int32(sStacksCounter)
		sStacksCounter++
	}

	// Program's Heap

	// sPrgrm.HeapOffset = int32(sHeapCounter)
	// lastHeapOffset := int32(-1)

	// sHp := sHeap{}
	// sHp.HeapOffset = int32(lastHeapOffset + 1)
	// sHp.HeapSize = int32(len(cxt.Heap.Heap))
	// sHp.HeapPointer = int32(cxt.Heap.HeapPointer)

	//Program Path

	sPrgrm.PathOffset, sPrgrm.PathSize = serializeName(cxt.Path, &sNamesMap, &sNames, &sNamesCounter)

	//Program's Steps

	sIdx := sIndex{}
	sIdx.ProgramOffset = int32(encoder.Size(sIdx))
	sIdx.NamesOffset = sIdx.ProgramOffset + int32(encoder.Size(sPrgrm))
	sIdx.ValuesOffset = sIdx.NamesOffset + int32(encoder.Size(sNames))
	sIdx.PackagesOffset = sIdx.ValuesOffset + int32(encoder.Size(sValues))
	sIdx.GlobalsOffset = sIdx.PackagesOffset + int32(encoder.Size(sPacks))
	sIdx.ImportsOffset = sIdx.GlobalsOffset + int32(encoder.Size(sGlobals))
	sIdx.FunctionsOffset = sIdx.ImportsOffset + int32(encoder.Size(sImps))
	sIdx.StructsOffset = sIdx.FunctionsOffset + int32(encoder.Size(sFns))
	sIdx.FieldsOffset = sIdx.StructsOffset + int32(encoder.Size(sStrcts))
	sIdx.ExpressionsOffset = sIdx.FieldsOffset + int32(encoder.Size(sFlds))
	sIdx.ArgumentsOffset = sIdx.ExpressionsOffset + int32(encoder.Size(sExprs))
	sIdx.CallsOffset = sIdx.ArgumentsOffset + int32(encoder.Size(sArgs))


	serialized = append(serialized, encoder.Serialize(sIdx)...)
	serialized = append(serialized, encoder.Serialize(sPrgrm)...)
	serialized = append(serialized, encoder.Serialize(sNames)...)
	serialized = append(serialized, encoder.Serialize(sValues)...)
	serialized = append(serialized, encoder.Serialize(sPacks)...)
	serialized = append(serialized, encoder.Serialize(sGlobals)...)
	serialized = append(serialized, encoder.Serialize(sImps)...)
	serialized = append(serialized, encoder.Serialize(sFns)...)
	serialized = append(serialized, encoder.Serialize(sStrcts)...)
	serialized = append(serialized, encoder.Serialize(sFlds)...)
	serialized = append(serialized, encoder.Serialize(sExprs)...)
	serialized = append(serialized, encoder.Serialize(sArgs)...)
	serialized = append(serialized, encoder.Serialize(sCalls)...)
	
	return &serialized
}


//Incomplete
func Deserialize (prgrm *[]byte) *CXProgram {
	cxt := CXProgram{}

	var dsIdx sIndex
	sIdx := (*prgrm)[:encoder.Size(sIndex{})]
	encoder.DeserializeRaw(sIdx, &dsIdx)

	var dsPrgrm sProgram
	sPrgrm := (*prgrm)[dsIdx.ProgramOffset:dsIdx.NamesOffset]
	encoder.DeserializeRaw(sPrgrm, &dsPrgrm)

	var dsNames []byte
	sNames := (*prgrm)[dsIdx.NamesOffset:dsIdx.ValuesOffset]
	encoder.DeserializeRaw(sNames, &dsNames)

	var dsValues []byte
	sValues := (*prgrm)[dsIdx.ValuesOffset:dsIdx.PackagesOffset]
	encoder.DeserializeRaw(sValues, &dsValues)

	var dsPacks []byte
	sPacks := (*prgrm)[dsIdx.PackageOffset: dsIdx.GlobalsOffset]
	encoder.DeserializeRaw(sPacks, &dsPacks)

	var dsGlobals []byte
	sGlobals := (*prgrm)[dsIdx.GlobalsOffset: dsIdx.ImportsOffset]
	encoder.DeserializeRaw(sGlobals, &dsGlobals)

	var dsImports []byte
	sImps := (*prgrm)[dsIdx.ImportsOffset: dsIdx.FunctionsOffset]
	encoder.DeserializeRaw(sImps, &dsImports)

	var dsFns []byte
	sFns := (*prgrm)[dsIdx.FunctionsOffset: dsIdx.StructsOffset]
	encoder.DeserializeRaw(sFns, &dsFns)

	var dsStrct []byte
	sStrct := (*prgrm)[dsIdx.StructsOffset: dsIdx.FieldsOffset]
	encoder.DeserializeRaw(sStrct, &dsStrct)

	var dsFlds []byte
	sFlds := (*prgrm)[dsIdx.FieldsOffset: dsIdx.ExpressionsOffset]
	encoder.DeserializeRaw(sFlds, &dsFlds)

	var dsExprs []byte
	sExprs := (*prgrm)[dsIdx.ExpressionsOffset: dsIdx.ArgumentsOffset]
	encoder.DeserializeRaw(sExprs, &dsExprs)

	var dsArgs []byte
	sArgs := (*prgrm)[dsIdx.ArgumentsOffset: dsIdx.CallsOffset]
	encoder.DeserializeRaw(sArgs, &dsArgs)

	var dsClls []byte
	sCll := (*prgrm)[dsIdx.CallsOffset:]
	encoder.DeserializeRaw(sCll, &dsClls)

	packs := make([]*CXPackage, 0)
	fns := make([]*CXFunction, 0)

	packSize := encoder.Size(sPackage{})
	for i:=0; i < int(dsPrgrm.PackagesSize); i++ {
		pack := CXPackage{}

		var dsPack sPackage
		sPack := dsPacks[i *packSize: (i + 1) * packSize]
		encoder.DeserializeRaw(sPack, &dsPack)

		var dsPackName []byte
		sPackName := dsNames[dsPack.NameOffset: dsPack.NameOffset + dsPack.NameSize]
		encoder.DeserializeRaw(sPackName, &dsPackName)

		pack.Name = string(dsPackName)
		packs = append(packs, &pack)

		fnSize := encoder.Size(sFunction{})
		fnsOffset := int(dsPack.FunctionsOffset) * fnSize
		for i := 0; i < int(dsPack.FunctionsSize); i++ {
			fn := CXFunction{}

			var dsFn sFunction
			sFn := dsFns[fnsOffset + ( i * fnSize): fnOffset + (i + 1) * fnSize]
			encoder.DeserializeRaw(sFn, &dsFn)

			var dsName []byte
			sName := dsNames[dsFn.NameOffset : dsFn.NameOffset + dsFn.NameSize]
			encoder.DeserializeRaw(sName, &dsName)

			fn.Name = string(dsName)
			fn.Module = &pack

			fns = append(fns, &fn)
		}
	}
	return &cxt
}
