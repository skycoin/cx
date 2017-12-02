package base

import (
	"fmt"
	"github.com/skycoin/skycoin/src/cipher/encoder"
)

/*
  Context
*/

type sIndex struct {
	ContextOffset int32
	NamesOffset int32
	ValuesOffset int32
	ModulesOffset int32
	DefinitionsOffset int32
	ImportsOffset int32
	FunctionsOffset int32
	StructsOffset int32
	FieldsOffset int32
	TypesOffset int32
	ParametersOffset int32
	ExpressionsOffset int32
	ArgumentsOffset int32
	CallsOffset int32
	OutputNamesOffset int32
}

type sProgram struct {
	ModulesOffset int32
	ModulesSize int32
	CurrentModuleOffset int32
	CallStackOffset int32
	CallStackSize int32
	Terminated int32
	StepsOffset int32
	StepsSize int32
}

type sCall struct {
	OperatorOffset int32
	Line int32
	StateOffset int32
	StateSize int32
	ReturnAddressOffset int32
	//ContextOffset // this might not be a problem, because the context will always be at byte 0
	ModuleOffset int32
}

/*
  Modules
*/

type sModule struct {
	NameOffset int32
	NameSize int32
	ImportsOffset int32
	ImportsSize int32
	FunctionsOffset int32
	FunctionsSize int32
	StructsOffset int32
	StructsSize int32
	DefinitionsOffset int32
	DefinitionsSize int32
	
	CurrentFunctionOffset int32
	CurrentStructOffset int32
}

type sDefinition struct {
	NameOffset int32
	NameSize int32
	TypeOffset int32
	TypeSize int32
	ValueOffset int32
	ValueSize int32

	ModuleOffset int32
}

/*
  Structs
*/

type sStruct struct {
	NameOffset int32
	NameSize int32
	FieldsOffset int32
	FieldsSize int32
	ModuleOffset int32
}

type sField struct {
	NameOffset int32
	NameSize int32
	TypeOffset int32
	TypeSize int32
}

// type sType struct {
// 	NameOffset int32
// 	NameSize int32
// }

/*
  Functions
*/

type sFunction struct {
	NameOffset int32
	NameSize int32
	InputsOffset int32
	InputsSize int32
	OutputsOffset int32
	OutputsSize int32
	ExpressionsOffset int32
	ExpressionsSize int32
	CurrentExpressionOffset int32
	ModuleOffset int32
}

type sParameter struct {
	NameOffset int32
	NameSize int32
	TypeOffset int32
	TypeSize int32
}

type sExpression struct {
	OperatorOffset int32
	ArgumentsOffset int32
	ArgumentsSize int32
	OutputNamesOffset int32 // these are CXDefinition
	OutputNamesSize int32
	Line int32
	FileLine int32
	TagOffset int32
	TagSize int32
	FunctionOffset int32
	ModuleOffset int32
}

type sArgument struct {
	TypeOffset int32
	TypeSize int32
	ValueOffset int32
	ValueSize int32
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

func serializeImports (imps []*CXModule, sModsMap *map[string]int, sImps *[]byte, sImpsCounter *int) (offset, size int32) {
	if imps != nil && len(imps) > 0 {
		offset = int32(*sImpsCounter)
		size = int32(len(imps))

		for _, imp := range imps {
			// we only need the index of the imported module
			*sImps = append(*sImps, encoder.SerializeAtomic(int32((*sModsMap)[imp.Name]))...)
			*sImpsCounter++
		}
	}
	return offset, size
}

func Serialize (cxt *CXProgram) *[]byte {
	// we will be appending the bytes here
	serialized := make([]byte, 0)
	
	sNames := make([]byte, 0)
	sNamesCounter := 0
	sNamesMap := make(map[string]int, 0)
	
	sValues := make([]byte, 0)
	sValuesCounter := 0
	// we don't use a map for values, as these can change, can be huge
	
	sMods := make([]byte, 0)
	sModsCounter := 0
	sModsMap := make(map[string]int, 0)

	sDefs := make([]byte, 0)
	sDefsCounter := 0

	sImps := make([]byte, 0)
	sImpsCounter := 0
	//sImpsMap not needed because they are modules

	sFns := make([]byte, 0)
	sFnsCounter := 0
	//sFnsMap := make(map[string]int, 0)
	sFnsMap := make([]string, 0)

	sStrcts := make([]byte, 0)
	sStrctsCounter := 0
	//sStrctsMap := make(map[string]int, 0)
	// we don't need the map because other elements don't reference structs

	sFlds := make([]byte, 0)
	sFldsCounter := 0

	sParams := make([]byte, 0)
	sParamsCounter := 0

	sExprs := make([]byte, 0)
	sExprsCounter := 0

	sArgs := make([]byte, 0)
	sArgsCounter := 0

	sCalls := make([]byte, 0)
	sCallsCounter := 0

	sOutNames := make([]byte, 0)
	sOutNamesCounter := 0
	
	// context
	sPrgrm := &sProgram{}

	sPrgrm.ModulesOffset = int32(sModsCounter)
	sPrgrm.ModulesSize = int32(len(cxt.Modules))

	//sPrgrm.CurrentModuleOffset =

	// adding function names first so expressions can reference their operators
	// module names too
	cxtModules := make([]*CXModule, 0)
	modFunctions := make([][]*CXFunction, len(cxt.Modules))
	modCounter := 0
	for _, mod := range cxt.Modules {
		sModsMap[mod.Name] = sModsCounter
		cxtModules = append(cxtModules, mod)
		sModsCounter++
		for _, fn := range mod.Functions {
			fnName := fmt.Sprintf("%s.%s", mod.Name, fn.Name)
			sFnsMap = append(sFnsMap, fnName)
			// converting mod's functions map to array
			modFunctions[modCounter] = append(modFunctions[modCounter], fn)
		}
		modCounter++
	}
	// resetting counters used above
	sModsCounter = 0


	// Modules
	for i, mod := range cxtModules {
		sMod := sModule{}

		// module name
		sMod.NameOffset, sMod.NameSize = serializeName(mod.Name, &sNamesMap, &sNames, &sNamesCounter)

		// serializing imports
		sMod.ImportsOffset, sMod.ImportsSize = serializeImports(mod.Imports, &sModsMap, &sImps, &sImpsCounter)

		// functions
		if modFunctions[i] != nil && len(modFunctions[i]) > 0 {
			sMod.FunctionsOffset = int32(sFnsCounter)
			sMod.FunctionsSize = int32(len(mod.Functions))

			for _, fn := range modFunctions[i] {
				sFn := sFunction{}

				// function name
				sFn.NameOffset, sFn.NameSize = serializeName(fn.Name, &sNamesMap, &sNames, &sNamesCounter)

				// inputs
				if fn.Inputs != nil && len(fn.Inputs) > 0 {
					sFn.InputsOffset = int32(sParamsCounter)
					sFn.InputsSize = int32(len(fn.Inputs))
					
					for _, inp := range fn.Inputs {
						sParam := sParameter{}

						// input name
						sParam.NameOffset, sParam.NameSize = serializeName(inp.Name, &sNamesMap, &sNames, &sNamesCounter)

						// input type
						sParam.TypeOffset, sParam.TypeSize = serializeName(inp.Typ, &sNamesMap, &sNames, &sNamesCounter)

						// save the sParam
						sParams = append(sParams, encoder.Serialize(sParam)...)
						sParamsCounter++
					}
				} else {
					sFn.InputsOffset = -1 // nil; fn does not have inputs
					sFn.InputsSize = -1
				}

				// outputs
				if fn.Outputs != nil && len(fn.Outputs) > 0 {
					sFn.OutputsOffset = int32(sParamsCounter)
					sFn.OutputsSize = int32(len(fn.Outputs))
					
					for _, out := range fn.Outputs {
						sParam := sParameter{}

						// output name
						sParam.NameOffset, sParam.NameSize = serializeName(out.Name, &sNamesMap, &sNames, &sNamesCounter)

						// output type
						sParam.TypeOffset, sParam.TypeSize = serializeName(out.Typ, &sNamesMap, &sNames, &sNamesCounter)

						// save the sParam
						sParams = append(sParams, encoder.Serialize(sParam)...)
						sParamsCounter++
					}
				} else {
					sFn.OutputsOffset = -1 // nil; fn does not have outputs
					sFn.OutputsSize = -1
				}

				// expressions
				if fn.Expressions != nil && len(fn.Expressions) > 0 {
					sFn.ExpressionsOffset = int32(sExprsCounter)
					sFn.ExpressionsSize = int32(len(fn.Expressions))
					
					for _, expr := range fn.Expressions {
						sExpr := sExpression{}
						opName := fmt.Sprintf("%s.%s", expr.Operator.Module.Name, expr.Operator.Name)

						// operator

						// looking for the function's offset
						opOffset := -1
						for i, fn := range sFnsMap {
							if opName == fn {
								opOffset = i
								break
							}
						}

						//if offset, ok := sFnsMap[opName]; ok {
						if opOffset >= 0 {
							sExpr.OperatorOffset = int32(opOffset)
						} else {
							panic(fmt.Sprintf("Expression's operator (%s) not found in sFnsMap", opName))
						}

						// output names
						if expr.OutputNames != nil && len(expr.OutputNames) > 0 {
							sExpr.OutputNamesOffset = int32(sOutNamesCounter)
							sExpr.OutputNamesSize = int32(len(expr.OutputNames))

							for _, outName := range expr.OutputNames {
								//sOutName := sOutputName{}
								sOutName := &sDefinition{}

								// outputName name
								sOutName.NameOffset, sOutName.NameSize = serializeName(outName.Name, &sNamesMap, &sNames, &sNamesCounter)
								// outputName type
								sOutName.TypeOffset, sOutName.TypeSize = serializeName(outName.Typ, &sNamesMap, &sNames, &sNamesCounter)
								// outputName value
								sOutName.ValueOffset, sOutName.ValueSize = serializeValue(outName.Value, &sValues, &sValuesCounter)

								// outputName Module
								sOutName.ModuleOffset =
									int32(sModsMap[outName.Module.Name])


								// saving the output name
								sOutNames = append(sOutNames, encoder.Serialize(sOutName)...)
								sOutNamesCounter++
							}
						}
						
						// arguments
						if expr.Arguments != nil && len(expr.Arguments) > 0 {
							sExpr.ArgumentsOffset = int32(sArgsCounter)
							sExpr.ArgumentsSize = int32(len(expr.Arguments))

							for _, arg := range expr.Arguments {
								sArg := sArgument{}

								// argument type
								sArg.TypeOffset, sArg.TypeSize = serializeName(arg.Typ, &sNamesMap, &sNames, &sNamesCounter)

								// argument value
								sArg.ValueOffset, sArg.ValueSize = serializeValue(arg.Value, &sValues, &sValuesCounter)

								// save the sArg
								sArgs = append(sArgs, encoder.Serialize(sArg)...)
								sArgsCounter++
							}
						}

						// line
						sExpr.Line = int32(expr.Line)

						// file line
						sExpr.FileLine = int32(expr.FileLine)

						// tag
						sExpr.TagOffset, sExpr.TagSize = serializeName(expr.Tag, &sNamesMap, &sNames, &sNamesCounter)


						// function
						fnOffset := 0
						for i, fnName := range sFnsMap {
							//fmt.Println(fnName)
							if fnName == fmt.Sprintf("%s.%s", mod.Name, fn.Name) {
								fnOffset = i
								break
							}
						}
						
						if fnOffset >= 0 {
							sExpr.FunctionOffset = int32(fnOffset)
						} else {
							panic(fmt.Sprintf("Function '%s' not found in sFnsMap", fn.Name))
						}
						
						// module
						sExpr.ModuleOffset = int32(sModsMap[expr.Module.Name])

						// save the expression
						sExprs = append(sExprs, encoder.Serialize(sExpr)...)
						// also checking if this expr is the fn's CurrentExpression
						if fn.CurrentExpression == expr {
							sFn.CurrentExpressionOffset = int32(sExprsCounter)
						}
						sExprsCounter++
					}
				}

				// module
				sFn.ModuleOffset = int32(sModsMap[fn.Module.Name])

				if mod.CurrentFunction == fn {
					sMod.CurrentFunctionOffset = int32(sFnsCounter)
				}

				// save the function
				sFns = append(sFns, encoder.Serialize(sFn)...)
				sFnsCounter++
			}
		}

		// structs
		if mod.Structs != nil && len(mod.Structs) > 0 {
			sMod.StructsOffset = int32(sStrctsCounter)
			sMod.StructsSize = int32(len(mod.Structs))

			for _, strct := range mod.Structs {
				sStrct := sStruct{}

				// struct name
				sStrct.NameOffset, sStrct.NameSize = serializeName(strct.Name, &sNamesMap, &sNames, &sNamesCounter)

				// fields
				if strct.Fields != nil && len(strct.Fields) > 0 {
					sStrct.FieldsOffset = int32(sFldsCounter)
					sStrct.FieldsSize = int32(len(strct.Fields))
					
					for _, fld := range strct.Fields {
						sFld := sField{}

						// field name
						sFld.NameOffset, sFld.NameSize = serializeName(fld.Name, &sNamesMap, &sNames, &sNamesCounter)

						// field type
						sFld.TypeOffset, sFld.TypeSize = serializeName(fld.Typ, &sNamesMap, &sNames, &sNamesCounter)

						// save the field
						sFlds = append(sFlds, encoder.Serialize(sFld)...)
						sFldsCounter++
					}
				}

				// module
				sStrct.ModuleOffset = int32(sModsMap[strct.Module.Name])

				if mod.CurrentStruct == strct {
					sMod.CurrentStructOffset = int32(sStrctsCounter)
				}

				// save the struct
				sStrcts = append(sStrcts, encoder.Serialize(sStrct)...)
				sStrctsCounter++
			}
		} else {
			sMod.CurrentStructOffset = int32(-1)
		}

		// definitions
		if mod.Definitions != nil && len(mod.Definitions) > 0 {
			sMod.DefinitionsOffset = int32(sDefsCounter)
			sMod.DefinitionsSize = int32(len(mod.Definitions))
			for _, def := range mod.Definitions {
				sDef := &sDefinition{}

				// definition name
				sDef.NameOffset, sDef.NameSize = serializeName(def.Name, &sNamesMap, &sNames, &sNamesCounter)

				// definition type
				sDef.TypeOffset, sDef.TypeSize = serializeName(def.Typ, &sNamesMap, &sNames, &sNamesCounter)

				// definition value
				sDef.ValueOffset, sDef.ValueSize = serializeValue(def.Value, &sValues, &sValuesCounter)

				sDef.ModuleOffset = int32(sModsMap[def.Module.Name])

				// save the definition
				sDefs = append(sDefs, encoder.Serialize(sDef)...)
				sDefsCounter++
			}
		}

		if cxt.CurrentModule == mod {
			sPrgrm.CurrentModuleOffset = int32(sModsCounter)
		}

		// save the mod
		sMods = append(sMods, encoder.Serialize(sMod)...)
		sModsCounter++
	}

	// Program terminated flag
	if cxt.Terminated {
		sPrgrm.Terminated = int32(1)
	} else {
		sPrgrm.Terminated = int32(0)
	}
	
	// Call stack
	sPrgrm.CallStackOffset = int32(sCallsCounter)
	sPrgrm.CallStackSize = int32(len(cxt.CallStack.Calls))
	lastCallOffset := int32(-1)
	for _, call := range cxt.CallStack.Calls {
		sCall := sCall{}

		// Operator
		opName := fmt.Sprintf("%s.%s", call.Operator.Module.Name, call.Operator.Name)

		// looking for the function's offset
		opOffset := -1
		for i, fn := range sFnsMap {
			//fmt.Printf("%s == %s\n", opName, fn)
			if opName == fn {
				opOffset = i
				break
			}
		}

		//if offset, ok := sFnsMap[opName]; ok {
		if opOffset >= 0 {
			sCall.OperatorOffset = int32(opOffset)
		} else {
			panic(fmt.Sprintf("Expression's operator (%s) not found in sFnsMap", opName))
		}

		// Line
		sCall.Line = int32(call.Line)

		// State
		if call.State != nil && len(call.State) > 0 {
			sCall.StateOffset = int32(sDefsCounter)
			sCall.StateSize = int32(len(call.State))
			for _, def := range call.State {
				sDef := &sDefinition{}

				// state definition name
				sDef.NameOffset, sDef.NameSize = serializeName(def.Name, &sNamesMap, &sNames, &sNamesCounter)

				// state definition type
				sDef.TypeOffset, sDef.TypeSize = serializeName(def.Typ, &sNamesMap, &sNames, &sNamesCounter)

				// state definition value
				sDef.ValueOffset, sDef.ValueSize = serializeValue(def.Value, &sValues, &sValuesCounter)

				sDef.ModuleOffset = int32(sModsMap[def.Module.Name])

				// save the definition
				sDefs = append(sDefs, encoder.Serialize(sDef)...)
				sDefsCounter++
			}
		}

		// Return address
		if lastCallOffset >= 0 {
			sCall.ReturnAddressOffset = lastCallOffset
		} else {
			sCall.ReturnAddressOffset = int32(-1) // nil
		}

		// Module
		sCall.ModuleOffset = int32(sModsMap[call.Module.Name])

		// save the call
		sCalls = append(sCalls, encoder.Serialize(sCall)...)
		lastCallOffset = int32(sCallsCounter)
		sCallsCounter++
	}

	// whole program
	sIdx := sIndex{}
	sIdx.ContextOffset = int32(encoder.Size(sIdx))
	sIdx.NamesOffset = sIdx.ContextOffset + int32(encoder.Size(sPrgrm))
	sIdx.ValuesOffset = sIdx.NamesOffset + int32(encoder.Size(sNames))
	sIdx.ModulesOffset = sIdx.ValuesOffset + int32(encoder.Size(sValues))
	sIdx.DefinitionsOffset = sIdx.ModulesOffset + int32(encoder.Size(sMods))
	sIdx.ImportsOffset = sIdx.DefinitionsOffset + int32(encoder.Size(sDefs))
	sIdx.FunctionsOffset = sIdx.ImportsOffset + int32(encoder.Size(sImps))
	sIdx.StructsOffset = sIdx.FunctionsOffset + int32(encoder.Size(sFns))
	sIdx.FieldsOffset = sIdx.StructsOffset + int32(encoder.Size(sStrcts))
	
	sIdx.ParametersOffset = sIdx.FieldsOffset + int32(encoder.Size(sFlds))
	
	sIdx.ExpressionsOffset = sIdx.ParametersOffset + int32(encoder.Size(sParams))
	sIdx.ArgumentsOffset = sIdx.ExpressionsOffset + int32(encoder.Size(sExprs))
	sIdx.CallsOffset = sIdx.ArgumentsOffset + int32(encoder.Size(sArgs))
	sIdx.OutputNamesOffset = sIdx.CallsOffset + int32(encoder.Size(sCalls))

	serialized = append(serialized, encoder.Serialize(sIdx)...)
	serialized = append(serialized, encoder.Serialize(sPrgrm)...)
	serialized = append(serialized, encoder.Serialize(sNames)...)
	serialized = append(serialized, encoder.Serialize(sValues)...)
	serialized = append(serialized, encoder.Serialize(sMods)...)
	serialized = append(serialized, encoder.Serialize(sDefs)...)
	serialized = append(serialized, encoder.Serialize(sImps)...)
	serialized = append(serialized, encoder.Serialize(sFns)...)
	serialized = append(serialized, encoder.Serialize(sStrcts)...)
	serialized = append(serialized, encoder.Serialize(sFlds)...)
	//serialized = append(serialized, encoder.Serialize(sTyps)...)
	serialized = append(serialized, encoder.Serialize(sParams)...)
	serialized = append(serialized, encoder.Serialize(sExprs)...)
	serialized = append(serialized, encoder.Serialize(sArgs)...)
	serialized = append(serialized, encoder.Serialize(sCalls)...)
	serialized = append(serialized, encoder.Serialize(sOutNames)...)
	
	return &serialized
}

func Deserialize (prgrm *[]byte) *CXProgram {
	cxt := CXProgram{}

	// First we deserialize the sIndex, as it contains the offsets of everything
	var dsIdx sIndex
	sIdx := (*prgrm)[:encoder.Size(sIndex{})]
	encoder.DeserializeRaw(sIdx, &dsIdx)

	// // Context
	var dsPrgrm sProgram
	sPrgrm := (*prgrm)[dsIdx.ContextOffset:dsIdx.NamesOffset]
	encoder.DeserializeRaw(sPrgrm, &dsPrgrm)
	
	// Names
	var dsNames []byte
	sNames := (*prgrm)[dsIdx.NamesOffset:dsIdx.ValuesOffset]
	encoder.DeserializeRaw(sNames, &dsNames)

	// Values
	var dsValues []byte
	sValues := (*prgrm)[dsIdx.ValuesOffset:dsIdx.ModulesOffset]
	encoder.DeserializeRaw(sValues, &dsValues)

	// Modules
	var dsMods []byte
	sMods := (*prgrm)[dsIdx.ModulesOffset:dsIdx.DefinitionsOffset]
	encoder.DeserializeRaw(sMods, &dsMods)

	// Definitions
	var dsDefs []byte
	sDefs := (*prgrm)[dsIdx.DefinitionsOffset:dsIdx.ImportsOffset]
	encoder.DeserializeRaw(sDefs, &dsDefs)

	// Imports
	var dsImps []byte
	sImps := (*prgrm)[dsIdx.ImportsOffset:dsIdx.FunctionsOffset]
	encoder.DeserializeRaw(sImps, &dsImps)

	// Functions
	var dsFns []byte
	sFns := (*prgrm)[dsIdx.FunctionsOffset:dsIdx.StructsOffset]
	encoder.DeserializeRaw(sFns, &dsFns)

	// Structs
	var dsStrcts []byte
	sStrcts := (*prgrm)[dsIdx.StructsOffset:dsIdx.FieldsOffset]
	encoder.DeserializeRaw(sStrcts, &dsStrcts)

	// Fields
	var dsFlds []byte
	sFlds := (*prgrm)[dsIdx.FieldsOffset:dsIdx.ParametersOffset]
	encoder.DeserializeRaw(sFlds, &dsFlds)

	// Parameters (Inputs & Outputs)
	var dsParams []byte
	sParams := (*prgrm)[dsIdx.ParametersOffset:dsIdx.ExpressionsOffset]
	encoder.DeserializeRaw(sParams, &dsParams)

	// Expressions
	var dsExprs []byte
	sExprs := (*prgrm)[dsIdx.ExpressionsOffset:dsIdx.ArgumentsOffset]
	encoder.DeserializeRaw(sExprs, &dsExprs)

	// Arguments
	var dsArgs []byte
	sArgs := (*prgrm)[dsIdx.ArgumentsOffset:dsIdx.CallsOffset]
	encoder.DeserializeRaw(sArgs, &dsArgs)

	// Calls
	var dsCalls []byte
	sCalls := (*prgrm)[dsIdx.CallsOffset:dsIdx.OutputNamesOffset]
	encoder.DeserializeRaw(sCalls, &dsCalls)

	// Output names
	var dsOutNames []byte
	sOutNames := (*prgrm)[dsIdx.OutputNamesOffset:]
	encoder.DeserializeRaw(sOutNames, &dsOutNames)

	/*
	   Deserializing elements
        */

	// Initializing CXModules for referencing modules as imports
	// Also initializing CXFunctions for referencing functions as expression operators
	mods := make([]*CXModule, 0)
	fns := make([]*CXFunction, 0)
	modSize := encoder.Size(sModule{})
	for i := 0; i < int(dsPrgrm.ModulesSize); i++ {
		mod := CXModule{}
		
		var dsMod sModule
		sMod := dsMods[i*modSize:(i+1)*modSize]
		encoder.DeserializeRaw(sMod, &dsMod)

		// Module Name
		var dsModName []byte
		sModName := dsNames[dsMod.NameOffset:dsMod.NameOffset + dsMod.NameSize]
		encoder.DeserializeRaw(sModName, &dsModName)

		//fmt.Println(string(dsModName))

		// Appending module with name attached
		mod.Name = string(dsModName)
		mods = append(mods, &mod)

		// Functions
		fnSize := encoder.Size(sFunction{})
		fnsOffset := int(dsMod.FunctionsOffset) * fnSize
		for i := 0; i < int(dsMod.FunctionsSize); i++ {
			fn := CXFunction{}

			var dsFn sFunction
			sFn := dsFns[fnsOffset + i*fnSize : fnsOffset + (i+1)*fnSize]
			encoder.DeserializeRaw(sFn, &dsFn)
			
			// Function Name
			var dsName []byte
			sName := dsNames[dsFn.NameOffset : dsFn.NameOffset + dsFn.NameSize]
			encoder.DeserializeRaw(sName, &dsName)

			fn.Name = string(dsName)
			fn.Module = &mod

			// Adding function's FQN and reference to map
			//fns[fmt.Sprintf("%s.%s", string(dsModName), string(dsName))] = &fn
			fns = append(fns, &fn)
		}
	}
	
	// Modules
	modsOffset := int(dsPrgrm.ModulesOffset) * modSize
	for i := 0; i < int(dsPrgrm.ModulesSize); i++ {
		//mod := CXModule{}
		
		var dsMod sModule
		sMod := dsMods[i*modSize:(i+1)*modSize]
		encoder.DeserializeRaw(sMod, &dsMod)

		// Name
		var dsModName []byte
		sModName := dsNames[dsMod.NameOffset:dsMod.NameOffset + dsMod.NameSize]
		encoder.DeserializeRaw(sModName, &dsModName)

		// Getting CXModule
		//mod := MakeModule(string(dsModName))
		var mod *CXModule
		for _, m := range mods {
			if m.Name == string(dsModName) {
				mod = m
				break
			}
		}

		//mod := mods[string(dsModName)]

		// Adding current module to context
		if int(dsPrgrm.CurrentModuleOffset) * modSize == modsOffset + i*modSize {
			cxt.CurrentModule = mod
		}

		// Imports (this []byte is holding module offsets, not sModules)
		//imps := make(map[string]*CXModule, 0)
		imps := make([]*CXModule, 0)
		impSize := encoder.Size(int32(0))
		impsOffset := int(dsMod.ImportsOffset)

		for i := 0; i < int(dsMod.ImportsSize); i++ {
			// Import (module) offset
			var dsModOffset int32
			sModOffset := dsImps[impsOffset + i*impSize : impsOffset + (i+1)*impSize]
			//encoder.DeserializeAtomic(sModOffset, &dsModOffset)
			encoder.DeserializeRaw(sModOffset, &dsModOffset)

			// Imported module
			var dsMod sModule
			sMod := dsMods[dsModOffset*int32(modSize) : dsModOffset*int32(modSize) + int32(modSize)]
			encoder.DeserializeRaw(sMod, &dsMod)

			// Imported module name
			var dsName []byte
			sName := dsNames[dsMod.NameOffset:dsMod.NameOffset + dsMod.NameSize]
			encoder.DeserializeRaw(sName, &dsName)

			//imps[string(dsName)] = mods[string(dsName)]

			for _, mod := range mods {
				if mod.Name == string(dsName) {
					imps = append(imps, mod)
					break
				}
			}
		}
		
		//Functions
		fnSize := encoder.Size(sFunction{})
		fnsOffset := int(dsMod.FunctionsOffset) * fnSize
		// fns contains ALL the functions. we need to do a subset
		modFns := make([]*CXFunction, 0)
		for i := 0; i < int(dsMod.FunctionsSize); i++ {
			fn := &CXFunction{}
			
			var dsFn sFunction
			sFn := dsFns[fnsOffset + i*fnSize : fnsOffset + (i+1)*fnSize]
			encoder.DeserializeRaw(sFn, &dsFn)

			// Name
			var dsName []byte
			sName := dsNames[dsFn.NameOffset : dsFn.NameOffset + dsFn.NameSize]
			encoder.DeserializeRaw(sName, &dsName)

			for _, f := range fns {
				if f.Name == string(dsName) && f.Module.Name == string(dsModName) {
					fn = f
					break
				}
			}

			//fn := fns[fmt.Sprintf("%s.%s", string(dsModName), string(dsName))]

			// Adding current function to module
			if int(dsMod.CurrentFunctionOffset) * fnSize == fnsOffset + i*fnSize {
				mod.CurrentFunction = fn
			}

			// Inputs
			var inps []*CXParameter
			paramSize := encoder.Size(sParameter{})
			inpsOffset := int(dsFn.InputsOffset) * paramSize
			for i := 0; i < int(dsFn.InputsSize); i++ {
				inp := CXParameter{}

				var dsParam sParameter
				sParam := dsParams[inpsOffset + i*paramSize : inpsOffset + (i+1)*paramSize]
				encoder.DeserializeRaw(sParam, &dsParam)

				// Name
				var dsName []byte
				sName := dsNames[dsParam.NameOffset : dsParam.NameOffset + dsParam.NameSize]
				encoder.DeserializeRaw(sName, &dsName)

				// Type
				var dsTypName []byte
				sTypName := dsNames[dsParam.TypeOffset : dsParam.TypeOffset + dsParam.TypeSize]
				encoder.DeserializeRaw(sTypName, &dsTypName)
				
				inp.Name = string(dsName)
				inp.Typ = string(dsTypName)

				// Appending final input
				inps = append(inps, &inp)
			}

			// Outputs
			var outs []*CXParameter
			outsOffset := int(dsFn.OutputsOffset) * paramSize
			for i := 0; i < int(dsFn.OutputsSize); i++ {
				out := CXParameter{}

				var dsParam sParameter
				sParam := dsParams[outsOffset + i*paramSize : outsOffset + (i+1)*paramSize]
				encoder.DeserializeRaw(sParam, &dsParam)

				// Name
				var dsName []byte
				sName := dsNames[dsParam.NameOffset : dsParam.NameOffset + dsParam.NameSize]
				encoder.DeserializeRaw(sName, &dsName)


				// Type
				var dsTypName []byte
				sTypName := dsNames[dsParam.TypeOffset : dsParam.TypeOffset + dsParam.TypeSize]
				encoder.DeserializeRaw(sTypName, &dsTypName)
				
				out.Name = string(dsName)
				out.Typ = string(dsTypName)

				// Appending final output
				outs = append(outs, &out)
			}

			// Expressions
			var exprs []*CXExpression
			exprSize := encoder.Size(sExpression{})
			exprsOffset := int(dsFn.ExpressionsOffset) * exprSize

			// Current expression
			var dsCurrExpr sExpression
			sCurrExpr := dsExprs[int(dsFn.CurrentExpressionOffset)*exprSize : int(dsFn.CurrentExpressionOffset)*exprSize + exprSize]
			encoder.DeserializeRaw(sCurrExpr, &dsCurrExpr)
			
			for i := 0; i < int(dsFn.ExpressionsSize); i++ {
				expr := CXExpression{}
				
				var dsExpr sExpression
				sExpr := dsExprs[exprsOffset + i*exprSize : exprsOffset + (i+1)*exprSize]
				encoder.DeserializeRaw(sExpr, &dsExpr)

				// Adding current expression to function
				if int(dsFn.CurrentExpressionOffset) * exprSize == exprsOffset + i*exprSize {
					fn.CurrentExpression = &expr
				}

				// Operator
				opSize := int32(encoder.Size(sFunction{}))

				var dsOp sFunction
				sOp := dsFns[dsExpr.OperatorOffset*opSize : dsExpr.OperatorOffset*opSize + opSize]
				encoder.DeserializeRaw(sOp, &dsOp)

				// Operator's name
				var dsOpName []byte
				sOpName := dsNames[dsOp.NameOffset : dsOp.NameOffset + dsOp.NameSize]
				encoder.DeserializeRaw(sOpName, &dsOpName)
				
				// Arguments
				var args []*CXArgument
				argSize := encoder.Size(sArgument{})
				argsOffset := int(dsExpr.ArgumentsOffset) * argSize
				for i := 0; i < int(dsExpr.ArgumentsSize); i++ {
					arg := CXArgument{}

					var dsArg sArgument
					sArg := dsArgs[argsOffset + i*argSize : argsOffset + (i+1)*argSize]
					encoder.DeserializeRaw(sArg, &dsArg)


					// Argument type
					var dsTypName []byte
					sTypName := dsNames[dsArg.TypeOffset : dsArg.TypeOffset + dsArg.TypeSize]
					encoder.DeserializeRaw(sTypName, &dsTypName)
					
					// Argument value
					var dsValue []byte
					sVal := dsValues[dsArg.ValueOffset : dsArg.ValueOffset + dsArg.ValueSize]
					encoder.DeserializeRaw(sVal, &dsValue)

					arg.Typ = string(dsTypName)
					arg.Value = &dsValue
					// arg.Offset = -1
					// arg.Size = -1

					// Appending final argument
					args = append(args, &arg)
				}

				//fmt.Println(int(dsExpr.OutputNamesSize))

				// Expression output names
				var outNames []*CXDefinition
				outNameSize := encoder.Size(sDefinition{})
				outNamesOffset := int(dsExpr.OutputNamesOffset) * outNameSize
				for i := 0; i < int(dsExpr.OutputNamesSize); i++ {
					outName := CXDefinition{}

					var dsOutName sDefinition
					// fmt.Println(dsExpr.OutputNamesOffset)
					// fmt.Println(len(dsOutNames))
					// fmt.Println(outNamesOffset)
					sOutName := dsOutNames[outNamesOffset + i*outNameSize : outNamesOffset + (i+1)*outNameSize]
					encoder.DeserializeRaw(sOutName, &dsOutName)

					// outName name
					var dsName []byte
					sName := dsNames[dsOutName.NameOffset : dsOutName.NameOffset + dsOutName.NameSize]
					encoder.DeserializeRaw(sName, &dsName)

					// outName type
					var dsTypName []byte
					sTypName := dsNames[dsOutName.TypeOffset : dsOutName.TypeOffset + dsOutName.TypeSize]
					encoder.DeserializeRaw(sTypName, &dsTypName)

					// outName value
					var dsValue []byte
					sVal := dsValues[dsOutName.ValueOffset : dsOutName.ValueOffset + dsOutName.ValueSize]
					encoder.DeserializeRaw(sVal, &dsValue)

					outName.Name = string(dsName)
					outName.Typ = string(dsTypName)
					outName.Value = &dsValue
					outName.Module = mod
					outName.Context = &cxt
					// outName.Offset = -1
					// outName.Size = -1
					
					// Appending final outName
					outNames = append(outNames, &outName)
				}

				// expression tag
				var dsTag []byte
				sTag := dsNames[dsExpr.TagOffset : dsExpr.TagOffset + dsExpr.TagSize]
				encoder.DeserializeRaw(sTag, &dsTag)

				for _, fn := range fns {
					if fn.Name == string(dsOpName) {
						expr.Operator = fn
						break
					}
				}
				
				//expr.Operator = fns[fmt.Sprintf("%s.%s", dsModName, dsOpName)]
				expr.Arguments = args
				expr.OutputNames = outNames
				expr.Line = int(dsExpr.Line)
				expr.Tag = string(dsTag)
				expr.FileLine = int(dsExpr.FileLine)
				expr.Function = fn
				expr.Context = &cxt
				expr.Module = mod
				
				// Appending final expression
				exprs = append(exprs, &expr)
			}

			// Constructing final function
			fn.Name = string(dsName)
			fn.Inputs = inps
			fn.Outputs = outs
			fn.Expressions = exprs
			// Current expression was added in the expression's loop
			fn.Module = mod

			// Appending final function to modFns
			//modFns[string(dsName)] = fn
			modFns = append(modFns, fn)
		}

		// Structs
		strcts := make([]*CXStruct, 0)
		strctSize := encoder.Size(sStruct{})
		strctsOffset := int(dsMod.StructsOffset) * strctSize

		// // Current struct
		// var dsCurrStrct sStruct
		// sCurrStrct := dsStrcts[int(dsMod.CurrentStructOffset)*strctSize : int(dsMod.CurrentStructOffset)*strctSize + strctSize]
		// encoder.DeserializeRaw(sCurrStrct, &dsCurrStrct)
		
		for i := 0; i < int(dsMod.StructsSize); i++ {
			strct := CXStruct{}
			
			var dsStrct sStruct
			sStrct := dsStrcts[strctsOffset + i*strctSize : strctsOffset + (i+1)*strctSize]
			encoder.DeserializeRaw(sStrct, &dsStrct)

			// Adding current struct to module
			if int(dsMod.CurrentStructOffset) * strctSize == strctsOffset + i*strctSize {
				mod.CurrentStruct = &strct
			}


			// Struct name
			var dsName []byte
			sName := dsNames[dsStrct.NameOffset : dsStrct.NameOffset + dsStrct.NameSize]
			encoder.DeserializeRaw(sName, &dsName)

			// Struct fields
			var flds []*CXField
			fldSize := encoder.Size(sField{})
			fldsOffset := int(dsStrct.FieldsOffset) * fldSize
			for i := 0; i < int(dsStrct.FieldsSize); i++ {
				fld := CXField{}

				var dsFld sField
				sFld := dsFlds[fldsOffset + i*fldSize : fldsOffset + (i+1)*fldSize]
				encoder.DeserializeRaw(sFld, &dsFld)

				// Field name
				var dsName []byte
				sName := dsNames[dsFld.NameOffset : dsFld.NameOffset + dsFld.NameSize]
				encoder.DeserializeRaw(sName, &dsName)

				// Field type
				var dsTypName []byte
				sTypName := dsNames[dsFld.TypeOffset : dsFld.TypeOffset + dsFld.TypeSize]
				encoder.DeserializeRaw(sTypName, &dsTypName)

				fld.Name = string(dsName)
				fld.Typ = string(dsTypName)

				// Appending final field
				flds = append(flds, &fld)
			}


			
			strct.Name = string(dsName)
			strct.Fields = flds
			strct.Module = mod
			strct.Context = &cxt

			// Appending final struct
			//strcts[string(dsName)] = &strct
			strcts = append(strcts, &strct)
		}

		// Definitions
		defs := make([]*CXDefinition, 0)
		defSize := encoder.Size(sDefinition{})
		defsOffset := int(dsMod.DefinitionsOffset) * defSize
		for i := 0; i < int(dsMod.DefinitionsSize); i++ {
			def := CXDefinition{}
			
			var dsDef sDefinition
			sDef := dsDefs[defsOffset + i*defSize : defsOffset + (i+1)*defSize]
			encoder.DeserializeRaw(sDef, &dsDef)

			// Definition name
			var dsName []byte
			sName := dsNames[dsDef.NameOffset : dsDef.NameOffset + dsDef.NameSize]
			encoder.DeserializeRaw(sName, &dsName)

			// Definition type
			var dsTypName []byte
			sTypName := dsNames[dsDef.TypeOffset : dsDef.TypeOffset + dsDef.TypeSize]
			encoder.DeserializeRaw(sTypName, &dsTypName)

			// Definition value
			var dsValue []byte
			sVal := dsValues[dsDef.ValueOffset : dsDef.ValueOffset + dsDef.ValueSize]
			encoder.DeserializeRaw(sVal, &dsValue)

			def.Name = string(dsName)
			def.Typ = string(dsTypName)
			def.Value = &dsValue
			def.Module = mod
			def.Context = &cxt
			// def.Offset = -1
			// def.Size = -1

			// Appending final definition
			//defs[string(dsName)] = &def
			defs = append(defs, &def)
		}

		mod.Name = string(dsModName)
		mod.Imports = imps
		mod.Functions = modFns
		mod.Structs = strcts
		mod.Definitions = defs
		mod.Context = &cxt
	}

	// Call stack
	calls := make([]*CXCall, 0)
	callSize := encoder.Size(sCall{})
	
	// this will always be 0. I'll leave it in the struct for consistency (like modulesoffset)
	//callsOffset := int(dsPrgrm.CallStackOffset) * callSize

	var lastCall *CXCall
	for i := 0; i < int(dsPrgrm.CallStackSize); i++ {
		call := CXCall{}

		var dsCall sCall
		sCal := dsCalls[i*callSize:(i+1)*callSize]
		encoder.DeserializeRaw(sCal, &dsCall)

		// Call's module
		var dsMod sModule
		modSize := encoder.Size(sModule{})
		sMod := dsMods[int(dsCall.ModuleOffset)*modSize : int(dsCall.ModuleOffset)*modSize + modSize]
		encoder.DeserializeRaw(sMod, &dsMod)

		// Call's module name
		var dsModName []byte
		sModName := dsNames[dsMod.NameOffset:dsMod.NameOffset + dsMod.NameSize]
		encoder.DeserializeRaw(sModName, &dsModName)

		//mod := mods[string(dsModName)]
		var mod *CXModule
		for _, m := range mods {
			if m.Name == string(dsModName) {
				mod = m
				break
			}
		}

		// Call's operator
		opSize := int32(encoder.Size(sFunction{}))

		var dsOp sFunction
		sOp := dsFns[dsCall.OperatorOffset*opSize : dsCall.OperatorOffset*opSize + opSize]
		encoder.DeserializeRaw(sOp, &dsOp)

		// Operator's name
		var dsOpName []byte
		sOpName := dsNames[dsOp.NameOffset : dsOp.NameOffset + dsOp.NameSize]
		encoder.DeserializeRaw(sOpName, &dsOpName)

		// Call's Operator Module
		var dsOpMod sModule
		sOpMod := dsMods[int(dsOp.ModuleOffset)*modSize : int(dsOp.ModuleOffset)*modSize + modSize]
		encoder.DeserializeRaw(sOpMod, &dsOpMod)

		// Call's Operator Module Name
		var dsOpModName []byte
		sOpModName := dsNames[dsOpMod.NameOffset:dsOpMod.NameOffset + dsOpMod.NameSize]
		encoder.DeserializeRaw(sOpModName, &dsOpModName)

		for _, mod := range mods {
			for _, fn := range mod.Functions {
				if fn.Name == string(dsOpName) && mod.Name == string(dsOpModName) {
					call.Operator = fn
				}
			}
		}

		// State
		defs := make([]*CXDefinition, 0)
		defSize := encoder.Size(sDefinition{})
		defsOffset := int(dsCall.StateOffset) * defSize
		for i := 0; i < int(dsCall.StateSize); i++ {
			def := CXDefinition{}
			
			var dsDef sDefinition
			sDef := dsDefs[defsOffset + i*defSize : defsOffset + (i+1)*defSize]
			encoder.DeserializeRaw(sDef, &dsDef)

			// Definition name
			var dsName []byte
			sName := dsNames[dsDef.NameOffset : dsDef.NameOffset + dsDef.NameSize]
			encoder.DeserializeRaw(sName, &dsName)

			// Definition type
			var dsTypName []byte
			sTypName := dsNames[dsDef.TypeOffset : dsDef.TypeOffset + dsDef.TypeSize]
			encoder.DeserializeRaw(sTypName, &dsTypName)

			// Definition value
			var dsValue []byte
			sVal := dsValues[dsDef.ValueOffset : dsDef.ValueOffset + dsDef.ValueSize]
			encoder.DeserializeRaw(sVal, &dsValue)

			def.Name = string(dsName)
			def.Typ = string(dsTypName)
			def.Value = &dsValue
			def.Module = mod
			def.Context = &cxt

			// Appending final definition
			defs = append(defs, &def)
		}

		call.Line = int(dsCall.Line)
		call.State = defs
		call.ReturnAddress = lastCall
		call.Module = mod
		call.Context = &cxt

		lastCall = &call

		// Appending final call
		calls = append(calls, &call)
	}

	//fmt.Println(calls)

	callStack := CXCallStack{}
	callStack.Calls = calls

	if dsPrgrm.Terminated > 0 {
		cxt.Terminated = true
	} else {
		cxt.Terminated = false
	}
	
	cxt.Modules = mods
	cxt.CallStack = &callStack
	cxt.Steps = make([]*CXCallStack, 0)

	return &cxt
}
