package base

import (
	//"fmt"
	//"unsafe"
	//"bytes"
	//"github.com/skycoin/skycoin/src/cipher/encoder"
)

/*
  Context
*/

type sProgram struct {
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

type sContext struct {
	ModulesOffset int32
	ModulesSize int32
	CurrentModuleOffset int32
	CallStackOffset int32
	CallStackSize int32
	StepsOffset int32
	StepsSize int32
	// we can't serialize steps because of skycoin/encoder limitations; it can't serialize funcs
	//ProgramStepsOffset int32
	//ProgramStepsSize int32
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

type sProgramStep struct {
	// can we serialize funcs with skycoin/encoder?
	// no, we can't
	// serialized programs will lose their ProgramSteps
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
	TypOffset int32
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
	TypOffset int32
}

type sType struct {
	NameOffset int32
	NameSize int32
}

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
	TypOffset int32
}

type sExpression struct {
	OperatorOffset int32
	ArgumentsOffset int32
	ArgumentsSize int32
	OutputNamesOffset int32
	OutputNamesSize int32
	Line int32
	FunctionOffset int32
	ModuleOffset int32
}

// used by sExpression's OutputNamesOffset and OutputNamesSize
type sOutputName struct {
	NameOffset int32
	NameSize int32
}

type sArgument struct {
	TypOffset int32
	ValueOffset int32
	ValueSize int32
}

/*
  Affordances

  Affordances must not be serialized
*/

// MANY blocks of code could become functions to decrease bloat
// func Serialize (cxt *CXProgram) *[]byte {
// 	// we will be appending the bytes here
// 	serialized := make([]byte, 0)
	
// 	sNames := make([]byte, 0)
// 	sNamesCounter := 0
// 	sNamesMap := make(map[string]int, 0)
	
// 	sValues := make([]byte, 0)
// 	sValuesCounter := 0
// 	// we don't use a map for values, as these can change, can be huge, and other stuff
	
// 	sMods := make([]byte, 0)
// 	sModsCounter := 0
// 	sModsMap := make(map[string]int, 0)

// 	sDefs := make([]byte, 0)
// 	sDefsCounter := 0

// 	sImps := make([]byte, 0)
// 	sImpsCounter := 0
// 	//sImpsMap not needed because they are modules

// 	sFns := make([]byte, 0)
// 	sFnsCounter := 0
// 	//sFnsMap := make(map[string]int, 0)
// 	sFnsMap := make([]string, 0)

// 	sStrcts := make([]byte, 0)
// 	sStrctsCounter := 0
// 	//sStrctsMap := make(map[string]int, 0)
// 	// we don't need the map because other elements don't reference structs

// 	sFlds := make([]byte, 0)
// 	sFldsCounter := 0

// 	sTyps := make([]byte, 0)
// 	sTypsCounter := 0

// 	sParams := make([]byte, 0)
// 	sParamsCounter := 0

// 	sExprs := make([]byte, 0)
// 	sExprsCounter := 0

// 	sArgs := make([]byte, 0)
// 	sArgsCounter := 0

// 	sCalls := make([]byte, 0)
// 	sCallsCounter := 0

// 	sOutNames := make([]byte, 0)
// 	sOutNamesCounter := 0
	
// 	// context
// 	sCxt := &sContext{}

// 	sCxt.ModulesOffset = int32(sModsCounter)
// 	sCxt.ModulesSize = int32(len(cxt.Modules))

// 	//sCxt.CurrentModuleOffset =

// 	// adding function names first so expressions can reference their operators
// 	// module names too
// 	cxtModules := make([]*CXModule, 0)
// 	modFunctions := make([][]*CXFunction, len(cxt.Modules))
// 	modCounter := 0
// 	for _, mod := range cxt.Modules {
// 		sModsMap[mod.Name] = sModsCounter
// 		cxtModules = append(cxtModules, mod)
// 		sModsCounter++
// 		for _, fn := range mod.Functions {
// 			fnName := fmt.Sprintf("%s.%s", mod.Name, fn.Name)
// 			sFnsMap = append(sFnsMap, fnName)
// 			// converting mod's functions map to array
// 			modFunctions[modCounter] = append(modFunctions[modCounter], fn)
// 		}
// 		modCounter++
// 	}
// 	// resetting counters used above
// 	sModsCounter = 0


// 	// Modules
// 	for i, mod := range cxtModules {
// 		sMod := sModule{}

// 		// name
// 		if offset, ok := sNamesMap[mod.Name]; ok {
// 			sMod.NameOffset = int32(offset)
// 			sMod.NameSize = int32(encoder.Size(mod.Name))
// 		} else {
// 			sNames = append(sNames, encoder.Serialize(mod.Name)...)
// 			sMod.NameOffset = int32(sNamesCounter)
// 			sMod.NameSize = int32(encoder.Size(mod.Name))
// 			sNamesMap[mod.Name] = sNamesCounter
// 			sNamesCounter = sNamesCounter + int(sMod.NameSize)
// 		}

// 		// imports
// 		if mod.Imports != nil && len(mod.Imports) > 0 {
// 			sMod.ImportsOffset = int32(sImpsCounter)
// 			sMod.ImportsSize = int32(len(mod.Imports))

// 			for _, imp := range mod.Imports {
// 				// we only need the index of the imported module
// 				sImps = append(sImps, encoder.SerializeAtomic(int32(sModsMap[imp.Name]))...)
// 				sImpsCounter++
// 			}
// 		}

// 		// functions
// 		if modFunctions[i] != nil && len(modFunctions[i]) > 0 {
// 			sMod.FunctionsOffset = int32(sFnsCounter)
// 			sMod.FunctionsSize = int32(len(mod.Functions))

// 			for _, fn := range modFunctions[i] {
// 				sFn := sFunction{}

// 				// name
// 				if offset, ok := sNamesMap[fn.Name]; ok {
// 					sFn.NameOffset = int32(offset)
// 					sFn.NameSize = int32(encoder.Size(fn.Name))
// 				} else {
// 					sNames = append(sNames, encoder.Serialize(fn.Name)...)
// 					sFn.NameOffset = int32(sNamesCounter)
// 					sFn.NameSize = int32(encoder.Size(fn.Name))
// 					sNamesMap[fn.Name] = sNamesCounter
// 					sNamesCounter = sNamesCounter + int(sFn.NameSize)
// 				}

// 				// inputs
// 				if fn.Inputs != nil && len(fn.Inputs) > 0 {
// 					sFn.InputsOffset = int32(sParamsCounter)
// 					sFn.InputsSize = int32(len(fn.Inputs))
					
// 					for _, inp := range fn.Inputs {
// 						sParam := sParameter{}

// 						// input name
// 						if offset, ok := sNamesMap[inp.Name]; ok {
// 							sParam.NameOffset = int32(offset)
// 							sParam.NameSize = int32(encoder.Size(inp.Name))
// 						} else {
// 							sNames = append(sNames, encoder.Serialize(inp.Name)...)
// 							sNamesMap[inp.Name] = sNamesCounter
// 							sParam.NameOffset = int32(sNamesCounter)
// 							sParam.NameSize = int32(encoder.Size(inp.Name))
// 							sNamesCounter = sNamesCounter + int(sParam.NameSize)
							
// 						}

// 						// input type
// 						typ := inp.Typ
// 						sTyp := sType{}
						
// 						// input type name
// 						if offset, ok := sNamesMap[typ.Name]; ok {
// 							sTyp.NameOffset = int32(offset)
// 							sTyp.NameSize = int32(encoder.Size(typ.Name))
// 						} else {
// 							sNames = append(sNames, encoder.Serialize(typ.Name)...)
// 							sNamesMap[typ.Name] = sNamesCounter
// 							sTyp.NameOffset = int32(sNamesCounter)
// 							sTyp.NameSize = int32(encoder.Size(typ.Name))
// 							sNamesCounter = sNamesCounter + int(sTyp.NameSize)
// 						}

// 						// save the sTyp
// 						sTyps = append(sTyps, encoder.Serialize(sTyp)...)
// 						sParam.TypOffset = int32(sTypsCounter)
// 						sTypsCounter++

// 						// save the sParam
// 						sParams = append(sParams, encoder.Serialize(sParam)...)
// 						sParamsCounter++
// 					}
// 				} else {
// 					sFn.InputsOffset = -1 // nil; fn does not have inputs
// 					sFn.InputsSize = -1
// 				}

// 				// outputs
// 				if fn.Outputs != nil && len(fn.Outputs) > 0 {
// 					sFn.OutputsOffset = int32(sParamsCounter)
// 					sFn.OutputsSize = int32(len(fn.Outputs))
					
// 					for _, out := range fn.Outputs {
// 						sParam := sParameter{}

// 						// output name
// 						if offset, ok := sNamesMap[out.Name]; ok {
// 							sParam.NameOffset = int32(offset)
// 							sParam.NameSize = int32(encoder.Size(out.Name))
// 						} else {
// 							sNames = append(sNames, encoder.Serialize(out.Name)...)
// 							sNamesMap[out.Name] = sNamesCounter
// 							sParam.NameOffset = int32(sNamesCounter)
// 							sParam.NameSize = int32(encoder.Size(out.Name))
// 							sNamesCounter = sNamesCounter + int(sParam.NameSize)
							
// 						}

// 						// output type
// 						typ := out.Typ
// 						sTyp := sType{}
						
// 						// output type name
// 						if offset, ok := sNamesMap[typ.Name]; ok {
// 							sTyp.NameOffset = int32(offset)
// 							sTyp.NameSize = int32(encoder.Size(typ.Name))
// 						} else {
// 							sNames = append(sNames, encoder.Serialize(typ.Name)...)
// 							sNamesMap[typ.Name] = sNamesCounter
// 							sTyp.NameOffset = int32(sNamesCounter)
// 							sTyp.NameSize = int32(encoder.Size(typ.Name))
// 							sNamesCounter = sNamesCounter + int(sTyp.NameSize)
// 						}

// 						// save the sTyp
// 						sTyps = append(sTyps, encoder.Serialize(sTyp)...)
// 						sParam.TypOffset = int32(sTypsCounter)
// 						sTypsCounter++

// 						// save the sParam
// 						sParams = append(sParams, encoder.Serialize(sParam)...)
// 						sParamsCounter++
// 					}
// 				} else {
// 					sFn.OutputsOffset = -1 // nil; fn does not have outputs
// 					sFn.OutputsSize = -1
// 				}

// 				// // output
// 				// if fn.Output != nil {
// 				// 	sFn.OutputOffset = int32(sParamsCounter)
// 				// 	sParam := sParameter{}

// 				// 	// name
// 				// 	if offset, ok := sNamesMap[fn.Output.Name]; ok {
// 				// 		sParam.NameOffset = int32(offset)
// 				// 		sParam.NameSize = int32(encoder.Size(fn.Output.Name))
// 				// 	} else {
// 				// 		sNames = append(sNames, encoder.Serialize(fn.Output.Name)...)
// 				// 		sNamesMap[fn.Output.Name] = sNamesCounter
// 				// 		sParam.NameOffset = int32(sNamesCounter)
// 				// 		sParam.NameSize = int32(encoder.Size(fn.Output.Name))
// 				// 		sNamesCounter = sNamesCounter + int(sParam.NameSize)
// 				// 	}

// 				// 	// output type
// 				// 	typ := fn.Output.Typ
// 				// 	sTyp := sType{}
					
// 				// 	// output type name
// 				// 	if offset, ok := sNamesMap[typ.Name]; ok {
// 				// 		sTyp.NameOffset = int32(offset)
// 				// 		sTyp.NameSize = int32(encoder.Size(typ.Name))
// 				// 	} else {
// 				// 		sNames = append(sNames, encoder.Serialize(typ.Name)...)
// 				// 		sNamesMap[typ.Name] = sNamesCounter
// 				// 		sTyp.NameOffset = int32(sNamesCounter)
// 				// 		sTyp.NameSize = int32(encoder.Size(typ.Name))
// 				// 		sNamesCounter = sNamesCounter + int(sTyp.NameSize)
// 				// 	}

// 				// 	// save the sTyp
// 				// 	sTyps = append(sTyps, encoder.Serialize(sTyp)...)
// 				// 	sParam.TypOffset = int32(sTypsCounter)
// 				// 	sTypsCounter++

// 				// 	// save the sParam
// 				// 	sParams = append(sParams, encoder.Serialize(sParam)...)
// 				// 	sParamsCounter++
// 				// } else {
// 				// 	sFn.OutputOffset = -1 // nil; fn does not have an output
// 				// }

// 				// expressions
// 				if fn.Expressions != nil && len(fn.Expressions) > 0 {
// 					sFn.ExpressionsOffset = int32(sExprsCounter)
// 					sFn.ExpressionsSize = int32(len(fn.Expressions))
					
// 					for _, expr := range fn.Expressions {
// 						sExpr := sExpression{}
// 						opName := fmt.Sprintf("%s.%s", expr.Operator.Module.Name, expr.Operator.Name)

// 						// operator

// 						// looking for the function's offset
// 						opOffset := -1
// 						for i, fn := range sFnsMap {
// 							//fmt.Printf("%s == %s\n", opName, fn)
// 							if opName == fn {
// 								opOffset = i
// 								break
// 							}
// 						}

// 						//if offset, ok := sFnsMap[opName]; ok {
// 						if opOffset >= 0 {
// 							sExpr.OperatorOffset = int32(opOffset)
// 						} else {
// 							panic(fmt.Sprintf("Expression's operator (%s) not found in sFnsMap", opName))
// 						}

// 						// output names
// 						if expr.OutputNames != nil && len(expr.OutputNames) > 0 {
// 							sExpr.OutputNamesOffset = int32(sOutNamesCounter)
// 							sExpr.OutputNamesSize = int32(len(expr.OutputNames))

// 							for _, outName := range expr.OutputNames {
// 								sOutName := sOutputName{}

// 								// name
// 								if offset, ok := sNamesMap[outName]; ok {
// 									sOutName.NameOffset = int32(offset)
// 									sOutName.NameSize = int32(encoder.Size(outName))
// 								} else {
// 									sNames = append(sNames, encoder.Serialize(outName)...)
// 									sNamesMap[outName] = sNamesCounter
// 									sOutName.NameOffset = int32(sNamesCounter)
// 									sOutName.NameSize = int32(encoder.Size(outName))
// 									sNamesCounter = sNamesCounter + int(sOutName.NameSize)
// 								}
								
// 								// saving the output name
// 								sOutNames = append(sOutNames, encoder.Serialize(sOutName)...)
// 								sOutNamesCounter++
// 							}
// 						}
						
// 						// arguments
// 						if expr.Arguments != nil && len(expr.Arguments) > 0 {
// 							sExpr.ArgumentsOffset = int32(sArgsCounter)
// 							sExpr.ArgumentsSize = int32(len(expr.Arguments))

// 							for _, arg := range expr.Arguments {
// 								sArg := sArgument{}
								
// 								// arg type
// 								typ := arg.Typ
// 								sTyp := sType{}
								
// 								// argument type name
// 								if offset, ok := sNamesMap[typ.Name]; ok {
// 									sTyp.NameOffset = int32(offset)
// 									sTyp.NameSize = int32(encoder.Size(typ.Name))
// 								} else {
// 									sNames = append(sNames, encoder.Serialize(typ.Name)...)
// 									sNamesMap[typ.Name] = sNamesCounter
// 									sTyp.NameOffset = int32(sNamesCounter)
// 									sTyp.NameSize = int32(encoder.Size(typ.Name))
// 									sNamesCounter = sNamesCounter + int(sTyp.NameSize)
// 								}

// 								// save the argument sTyp
// 								sTyps = append(sTyps, encoder.Serialize(sTyp)...)
// 								sArg.TypOffset = int32(sTypsCounter)
// 								sTypsCounter++

// 								// argument value
// 								sValues = append(sValues, encoder.Serialize(*arg.Value)...)
// 								sArg.ValueOffset = int32(sValuesCounter)
// 								sArg.ValueSize = int32(encoder.Size(*arg.Value))
// 								sValuesCounter = sValuesCounter + int(sArg.ValueSize)

// 								// save the sArg
// 								sArgs = append(sArgs, encoder.Serialize(sArg)...)
// 								sArgsCounter++
// 							}
// 						}

// 						// output name
// 						// if offset, ok := sNamesMap[expr.OutputName]; ok {
// 						// 	sExpr.OutputNameOffset = int32(offset)
// 						// 	sExpr.OutputNameSize = int32(encoder.Size(expr.OutputName))
// 						// } else {
// 						// 	sNames = append(sNames, encoder.Serialize(expr.OutputName)...)
// 						// 	sNamesMap[expr.OutputName] = sNamesCounter
// 						// 	sExpr.OutputNameOffset = int32(sNamesCounter)
// 						// 	sExpr.OutputNameSize = int32(encoder.Size(expr.OutputName))
// 						// 	sNamesCounter = sNamesCounter + int(sExpr.OutputNameSize)
// 						// }

// 						// line
// 						sExpr.Line = int32(expr.Line)

// 						// function
// 						fnOffset := 0
// 						for i, fnName := range sFnsMap {
// 							//fmt.Println(fnName)
// 							if fnName == fmt.Sprintf("%s.%s", mod.Name, fn.Name) {
// 								fnOffset = i
// 								break
// 							}
// 						}
						
// 						//fmt.Printf("%s <===> %v\n", fmt.Sprintf("%s.%s", mod.Name, fn.Name), sFnsMap)
// 						if fnOffset >= 0 {
// 							sExpr.FunctionOffset = int32(fnOffset)
// 						} else {
// 							panic(fmt.Sprintf("Function '%s' not found in sFnsMap", fn.Name))
// 						}
						

// 						// module
// 						sExpr.ModuleOffset = int32(sModsMap[expr.Module.Name])

// 						// save the expression
// 						sExprs = append(sExprs, encoder.Serialize(sExpr)...)
// 						// also checking if this expr is the fn's CurrentExpression
// 						if fn.CurrentExpression == expr {
// 							sFn.CurrentExpressionOffset = int32(sExprsCounter)
// 						}
// 						sExprsCounter++
// 					}
// 				}

// 				// module
// 				sFn.ModuleOffset = int32(sModsMap[fn.Module.Name])

// 				if mod.CurrentFunction == fn {
// 					sMod.CurrentFunctionOffset = int32(sFnsCounter)
// 				}

// 				// save the function
// 				sFns = append(sFns, encoder.Serialize(sFn)...)
// 				sFnsCounter++
// 			}
// 		}

// 		// structs
// 		if mod.Structs != nil && len(mod.Structs) > 0 {
// 			sMod.StructsOffset = int32(sStrctsCounter)
// 			sMod.StructsSize = int32(len(mod.Structs))

// 			for _, strct := range mod.Structs {
// 				sStrct := sStruct{}

// 				// name
// 				if offset, ok := sNamesMap[strct.Name]; ok {
// 					sStrct.NameOffset = int32(offset)
// 					sStrct.NameSize = int32(encoder.Size(strct.Name))
// 				} else {
// 					sNames = append(sNames, encoder.Serialize(strct.Name)...)
// 					sStrct.NameOffset = int32(sNamesCounter)
// 					sStrct.NameSize = int32(encoder.Size(strct.Name))
// 					sNamesMap[strct.Name] = sNamesCounter
// 					sNamesCounter = sNamesCounter + int(sStrct.NameSize)
// 				}

// 				// fields
// 				if strct.Fields != nil && len(strct.Fields) > 0 {
// 					sStrct.FieldsOffset = int32(sFldsCounter)
// 					sStrct.FieldsSize = int32(len(strct.Fields))
					
// 					for _, fld := range strct.Fields {
// 						sFld := sField{}

// 						// name
// 						if offset, ok := sNamesMap[fld.Name]; ok {
// 							sFld.NameOffset = int32(offset)
// 							sFld.NameSize = int32(encoder.Size(fld.Name))
// 						} else {
// 							sNames = append(sNames, encoder.Serialize(fld.Name)...)
// 							sFld.NameOffset = int32(sNamesCounter)
// 							sFld.NameSize = int32(encoder.Size(fld.Name))
// 							sNamesMap[fld.Name] = sNamesCounter
// 							sNamesCounter = sNamesCounter + int(sFld.NameSize)
// 						}

// 						// type
// 						typ := fld.Typ
// 						sTyp := sType{}
						
// 						// field type name
// 						if offset, ok := sNamesMap[typ.Name]; ok {
// 							sTyp.NameOffset = int32(offset)
// 							sTyp.NameSize = int32(encoder.Size(typ.Name))
// 						} else {
// 							sNames = append(sNames, encoder.Serialize(typ.Name)...)
// 							sNamesMap[typ.Name] = sNamesCounter
// 							sTyp.NameOffset = int32(sNamesCounter)
// 							sTyp.NameSize = int32(encoder.Size(typ.Name))
// 							sNamesCounter = sNamesCounter + int(sTyp.NameSize)
// 						}

// 						// save the argument sTyp
// 						sTyps = append(sTyps, encoder.Serialize(sTyp)...)
// 						sFld.TypOffset = int32(sTypsCounter)
// 						sTypsCounter++

// 						// save the field
// 						sFlds = append(sFlds, encoder.Serialize(sFld)...)
// 						sFldsCounter++
// 					}
// 				}

// 				// module
// 				sStrct.ModuleOffset = int32(sModsMap[strct.Module.Name])

// 				if mod.CurrentStruct == strct {
// 					sMod.CurrentStructOffset = int32(sStrctsCounter)
// 				}

// 				// save the struct
// 				sStrcts = append(sStrcts, encoder.Serialize(sStrct)...)
// 				sStrctsCounter++
// 			}
// 		}

// 		// definitions
// 		if mod.Definitions != nil && len(mod.Definitions) > 0 {
// 			sMod.DefinitionsOffset = int32(sDefsCounter)
// 			sMod.DefinitionsSize = int32(len(mod.Definitions))
// 			for _, def := range mod.Definitions {
// 				sDef := &sDefinition{}

// 				// name
// 				if offset, ok := sNamesMap[def.Name]; ok {
// 					sDef.NameOffset = int32(offset)
// 					sDef.NameSize = int32(encoder.Size(def.Name))
// 				} else {
// 					sNames = append(sNames, encoder.Serialize(def.Name)...)
// 					sDef.NameOffset = int32(sNamesCounter)
// 					sDef.NameSize = int32(encoder.Size(def.Name))
// 					sNamesMap[def.Name] = sNamesCounter
// 					sNamesCounter = sNamesCounter + int(sDef.NameSize)
// 				}

// 				// type
// 				typ := def.Typ
// 				sTyp := sType{}
				
// 				// field type name
// 				if offset, ok := sNamesMap[typ.Name]; ok {
// 					sTyp.NameOffset = int32(offset)
// 					sTyp.NameSize = int32(encoder.Size(typ.Name))
// 				} else {
// 					sNames = append(sNames, encoder.Serialize(typ.Name)...)
// 					sNamesMap[typ.Name] = sNamesCounter
// 					sTyp.NameOffset = int32(sNamesCounter)
// 					sTyp.NameSize = int32(encoder.Size(typ.Name))
// 					sNamesCounter = sNamesCounter + int(sTyp.NameSize)
// 				}

// 				// save the definition sTyp
// 				sTyps = append(sTyps, encoder.Serialize(sTyp)...)
// 				sDef.TypOffset = int32(sTypsCounter)
// 				sTypsCounter++

// 				// value
// 				sValues = append(sValues, encoder.Serialize(*def.Value)...)
// 				sDef.ValueOffset = int32(sValuesCounter)
// 				sDef.ValueSize = int32(encoder.Size(*def.Value))
// 				sValuesCounter = sValuesCounter + int(sDef.ValueSize)

// 				sDef.ModuleOffset = int32(sModsMap[def.Module.Name])

// 				// save the definition
// 				sDefs = append(sDefs, encoder.Serialize(sDef)...)
// 				sDefsCounter++
// 			}
// 		}
		

// 		if cxt.CurrentModule == mod {
// 			sCxt.CurrentModuleOffset = int32(sModsCounter)
// 		}

// 		// save the mod
// 		sMods = append(sMods, encoder.Serialize(sMod)...)
// 		sModsCounter++
// 	}

// 	// Call stack
// 	sCxt.CallStackOffset = int32(sCallsCounter)
// 	sCxt.CallStackSize = int32(len(cxt.CallStack.Calls))
// 	lastCallOffset := int32(-1)
// 	for _, call := range cxt.CallStack.Calls {
// 		sCall := sCall{}

// 		// Operator
// 		opName := fmt.Sprintf("%s.%s", call.Operator.Module.Name, call.Operator.Name)

// 		// looking for the function's offset
// 		opOffset := -1
// 		for i, fn := range sFnsMap {
// 			//fmt.Printf("%s == %s\n", opName, fn)
// 			if opName == fn {
// 				opOffset = i
// 				break
// 			}
// 		}

// 		//if offset, ok := sFnsMap[opName]; ok {
// 		if opOffset >= 0 {
// 			sCall.OperatorOffset = int32(opOffset)
// 		} else {
// 			panic(fmt.Sprintf("Expression's operator (%s) not found in sFnsMap", opName))
// 		}

// 		// Line
// 		sCall.Line = int32(call.Line)

// 		// State
// 		if call.State != nil && len(call.State) > 0 {
// 			sCall.StateOffset = int32(sDefsCounter)
// 			sCall.StateSize = int32(len(call.State))
// 			for _, def := range call.State {
// 				sDef := &sDefinition{}

// 				// name
// 				if offset, ok := sNamesMap[def.Name]; ok {
// 					sDef.NameOffset = int32(offset)
// 					sDef.NameSize = int32(encoder.Size(def.Name))
// 				} else {
// 					sNames = append(sNames, encoder.Serialize(def.Name)...)
// 					sDef.NameOffset = int32(sNamesCounter)
// 					sDef.NameSize = int32(encoder.Size(def.Name))
// 					sNamesMap[def.Name] = sNamesCounter
// 					sNamesCounter = sNamesCounter + int(sDef.NameSize)
// 				}

// 				// type
// 				typ := def.Typ
// 				sTyp := sType{}
				
// 				// field type name
// 				if offset, ok := sNamesMap[typ.Name]; ok {
// 					sTyp.NameOffset = int32(offset)
// 					sTyp.NameSize = int32(encoder.Size(typ.Name))
// 				} else {
// 					sNames = append(sNames, encoder.Serialize(typ.Name)...)
// 					sNamesMap[typ.Name] = sNamesCounter
// 					sTyp.NameOffset = int32(sNamesCounter)
// 					sTyp.NameSize = int32(encoder.Size(typ.Name))
// 					sNamesCounter = sNamesCounter + int(sTyp.NameSize)
// 				}

// 				// save the definition sTyp
// 				sTyps = append(sTyps, encoder.Serialize(sTyp)...)
// 				sDef.TypOffset = int32(sTypsCounter)
// 				sTypsCounter++

// 				// value
// 				sValues = append(sValues, encoder.Serialize(*def.Value)...)
// 				sDef.ValueOffset = int32(sValuesCounter)
// 				sDef.ValueSize = int32(encoder.Size(*def.Value))
// 				sValuesCounter = sValuesCounter + int(sDef.ValueSize)

// 				sDef.ModuleOffset = int32(sModsMap[def.Module.Name])

// 				// save the definition
// 				sDefs = append(sDefs, encoder.Serialize(sDef)...)
// 				sDefsCounter++
// 			}
// 		}

// 		// Return address
// 		if lastCallOffset >= 0 {
// 			sCall.ReturnAddressOffset = lastCallOffset
// 		} else {
// 			sCall.ReturnAddressOffset = int32(-1) // nil
// 		}

// 		// Module
// 		sCall.ModuleOffset = int32(sModsMap[call.Module.Name])

// 		// save the call
// 		sCalls = append(sCalls, encoder.Serialize(sCall)...)
// 		lastCallOffset = int32(sCallsCounter)
// 		sCallsCounter++
		
// 	}
	
// 	// whole program
// 	sPrgrm := sProgram{}
// 	sPrgrm.ContextOffset = int32(encoder.Size(sPrgrm))
// 	sPrgrm.NamesOffset = sPrgrm.ContextOffset + int32(encoder.Size(sCxt))
// 	sPrgrm.ValuesOffset = sPrgrm.NamesOffset + int32(encoder.Size(sNames))
// 	sPrgrm.ModulesOffset = sPrgrm.ValuesOffset + int32(encoder.Size(sValues))
// 	sPrgrm.DefinitionsOffset = sPrgrm.ModulesOffset + int32(encoder.Size(sMods))
// 	sPrgrm.ImportsOffset = sPrgrm.DefinitionsOffset + int32(encoder.Size(sDefs))
// 	sPrgrm.FunctionsOffset = sPrgrm.ImportsOffset + int32(encoder.Size(sImps))
// 	sPrgrm.StructsOffset = sPrgrm.FunctionsOffset + int32(encoder.Size(sFns))
// 	sPrgrm.FieldsOffset = sPrgrm.StructsOffset + int32(encoder.Size(sStrcts))
// 	sPrgrm.TypesOffset = sPrgrm.FieldsOffset + int32(encoder.Size(sFlds))
// 	sPrgrm.ParametersOffset = sPrgrm.TypesOffset + int32(encoder.Size(sTyps))
// 	sPrgrm.ExpressionsOffset = sPrgrm.ParametersOffset + int32(encoder.Size(sParams))
// 	sPrgrm.ArgumentsOffset = sPrgrm.ExpressionsOffset + int32(encoder.Size(sExprs))
// 	sPrgrm.CallsOffset = sPrgrm.ArgumentsOffset + int32(encoder.Size(sArgs))
// 	sPrgrm.OutputNamesOffset = sPrgrm.CallsOffset + int32(encoder.Size(sCalls))

// 	serialized = append(serialized, encoder.Serialize(sPrgrm)...)
// 	serialized = append(serialized, encoder.Serialize(sCxt)...)
// 	serialized = append(serialized, encoder.Serialize(sNames)...)
// 	serialized = append(serialized, encoder.Serialize(sValues)...)
// 	serialized = append(serialized, encoder.Serialize(sMods)...)
// 	serialized = append(serialized, encoder.Serialize(sDefs)...)
// 	serialized = append(serialized, encoder.Serialize(sImps)...)
// 	serialized = append(serialized, encoder.Serialize(sFns)...)
// 	serialized = append(serialized, encoder.Serialize(sStrcts)...)
// 	serialized = append(serialized, encoder.Serialize(sFlds)...)
// 	serialized = append(serialized, encoder.Serialize(sTyps)...)
// 	serialized = append(serialized, encoder.Serialize(sParams)...)
// 	serialized = append(serialized, encoder.Serialize(sExprs)...)
// 	serialized = append(serialized, encoder.Serialize(sArgs)...)
// 	serialized = append(serialized, encoder.Serialize(sCalls)...)
// 	serialized = append(serialized, encoder.Serialize(sOutNames)...)
	
// 	return &serialized
// }

// Just as Serialize(), this function needs to be divided to be better organized and more maintainable
// func Deserialize (prgrm *[]byte) *CXProgram {
// 	cxt := CXProgram{}

// 	// First we deserialize the sProgram, as it contains the offsets of everything
// 	var dsPrgrm sProgram
// 	sPrgrm := (*prgrm)[:encoder.Size(sProgram{})]
// 	encoder.DeserializeRaw(sPrgrm, &dsPrgrm)

// 	// // Context
// 	var dsCxt sContext
// 	sCxt := (*prgrm)[dsPrgrm.ContextOffset:dsPrgrm.NamesOffset]
// 	encoder.DeserializeRaw(sCxt, &dsCxt)
	
// 	// Names
// 	var dsNames []byte
// 	sNames := (*prgrm)[dsPrgrm.NamesOffset:dsPrgrm.ValuesOffset]
// 	encoder.DeserializeRaw(sNames, &dsNames)

// 	// Values
// 	var dsValues []byte
// 	sValues := (*prgrm)[dsPrgrm.ValuesOffset:dsPrgrm.ModulesOffset]
// 	encoder.DeserializeRaw(sValues, &dsValues)

// 	// Modules
// 	var dsMods []byte
// 	sMods := (*prgrm)[dsPrgrm.ModulesOffset:dsPrgrm.DefinitionsOffset]
// 	encoder.DeserializeRaw(sMods, &dsMods)

// 	// Definitions
// 	var dsDefs []byte
// 	sDefs := (*prgrm)[dsPrgrm.DefinitionsOffset:dsPrgrm.ImportsOffset]
// 	encoder.DeserializeRaw(sDefs, &dsDefs)

// 	// Imports
// 	var dsImps []byte
// 	sImps := (*prgrm)[dsPrgrm.ImportsOffset:dsPrgrm.FunctionsOffset]
// 	encoder.DeserializeRaw(sImps, &dsImps)

// 	// Functions
// 	var dsFns []byte
// 	sFns := (*prgrm)[dsPrgrm.FunctionsOffset:dsPrgrm.StructsOffset]
// 	encoder.DeserializeRaw(sFns, &dsFns)

// 	// Structs
// 	var dsStrcts []byte
// 	sStrcts := (*prgrm)[dsPrgrm.StructsOffset:dsPrgrm.FieldsOffset]
// 	encoder.DeserializeRaw(sStrcts, &dsStrcts)

// 	// Fields
// 	var dsFlds []byte
// 	sFlds := (*prgrm)[dsPrgrm.FieldsOffset:dsPrgrm.TypesOffset]
// 	encoder.DeserializeRaw(sFlds, &dsFlds)

// 	// Types
// 	var dsTyps []byte
// 	sTyps := (*prgrm)[dsPrgrm.TypesOffset:dsPrgrm.ParametersOffset]
// 	encoder.DeserializeRaw(sTyps, &dsTyps)

// 	// Parameters (Inputs & Outputs)
// 	var dsParams []byte
// 	sParams := (*prgrm)[dsPrgrm.ParametersOffset:dsPrgrm.ExpressionsOffset]
// 	encoder.DeserializeRaw(sParams, &dsParams)

// 	// Expressions
// 	var dsExprs []byte
// 	sExprs := (*prgrm)[dsPrgrm.ExpressionsOffset:dsPrgrm.ArgumentsOffset]
// 	encoder.DeserializeRaw(sExprs, &dsExprs)

// 	// Arguments
// 	var dsArgs []byte
// 	sArgs := (*prgrm)[dsPrgrm.ArgumentsOffset:]
// 	encoder.DeserializeRaw(sArgs, &dsArgs)

// 	// Calls
// 	var dsCalls []byte
// 	sCalls := (*prgrm)[dsPrgrm.CallsOffset:]
// 	encoder.DeserializeRaw(sCalls, &dsCalls)

// 	// Output names
// 	var dsOutNames []byte
// 	sOutNames := (*prgrm)[dsPrgrm.OutputNamesOffset:]
// 	encoder.DeserializeRaw(sOutNames, &dsOutNames)


// 	/*
// 	   Deserializing elements
//         */

// 	// Initializing CXModules for referencing modules as imports
// 	// Also initializing CXFunctions for referencing functions as expression operators
// 	mods := make(map[string]*CXModule, 0)
// 	fns := make(map[string]*CXFunction, 0)
// 	modSize := encoder.Size(sModule{})
// 	for i := 0; i < int(dsCxt.ModulesSize); i++ {
// 		mod := CXModule{}
		
// 		var dsMod sModule
// 		sMod := dsMods[i*modSize:(i+1)*modSize]
// 		encoder.DeserializeRaw(sMod, &dsMod)

// 		// Module Name
// 		var dsModName []byte
// 		sModName := dsNames[dsMod.NameOffset:dsMod.NameOffset + dsMod.NameSize]
// 		encoder.DeserializeRaw(sModName, &dsModName)

// 		// Adding module name and reference to map
// 		mods[string(dsModName)] = &mod

// 		// Functions
// 		fnSize := encoder.Size(sFunction{})
// 		fnsOffset := int(dsMod.FunctionsOffset) * fnSize
// 		for i := 0; i < int(dsMod.FunctionsSize); i++ {
// 			fn := CXFunction{}

// 			var dsFn sFunction
// 			sFn := dsFns[fnsOffset + i*fnSize : fnsOffset + (i+1)*fnSize]
// 			encoder.DeserializeRaw(sFn, &dsFn)
			
// 			// Function Name
// 			var dsName []byte
// 			sName := dsNames[dsFn.NameOffset : dsFn.NameOffset + dsFn.NameSize]
// 			encoder.DeserializeRaw(sName, &dsName)

// 			// Adding function's FQDM and reference to map
// 			fns[fmt.Sprintf("%s.%s", string(dsModName), string(dsName))] = &fn
// 		}
// 	}
	
// 	// Modules
// 	modsOffset := int(dsCxt.ModulesOffset) * modSize
// 	for i := 0; i < int(dsCxt.ModulesSize); i++ {
// 		//mod := CXModule{}
		
// 		var dsMod sModule
// 		sMod := dsMods[i*modSize:(i+1)*modSize]
// 		encoder.DeserializeRaw(sMod, &dsMod)

// 		// Name
// 		var dsModName []byte
// 		sModName := dsNames[dsMod.NameOffset:dsMod.NameOffset + dsMod.NameSize]
// 		encoder.DeserializeRaw(sModName, &dsModName)

// 		// Getting CXModule
// 		mod := mods[string(dsModName)]

// 		// Adding current module to context
// 		if int(dsCxt.CurrentModuleOffset) * modSize == modsOffset + i*modSize {
// 			cxt.CurrentModule = mod
// 		}

// 		// Imports (this []byte is holding module offsets, not sModules)
// 		imps := make(map[string]*CXModule, 0)
// 		impSize := encoder.Size(int32(0))
// 		impsOffset := int(dsMod.ImportsOffset)

// 		//fmt.Printf("yeee %d\n", dsMod.ImportsOffset)
		
// 		for i := 0; i < int(dsMod.ImportsSize); i++ {
// 			// Import (module) offset
// 			var dsModOffset int32
// 			sModOffset := dsImps[impsOffset + i*impSize : impsOffset + (i+1)*impSize]
// 			//encoder.DeserializeAtomic(sModOffset, &dsModOffset)
// 			encoder.DeserializeRaw(sModOffset, &dsModOffset)

// 			// Imported module
// 			var dsMod sModule
// 			sMod := dsMods[dsModOffset*int32(modSize) : dsModOffset*int32(modSize) + int32(modSize)]
// 			encoder.DeserializeRaw(sMod, &dsMod)

// 			// Imported module name
// 			var dsName []byte
// 			sName := dsNames[dsMod.NameOffset:dsMod.NameOffset + dsMod.NameSize]
// 			encoder.DeserializeRaw(sName, &dsName)

// 			imps[string(dsName)] = mods[string(dsName)]
// 		}
		
// 		//Functions
// 		fnSize := encoder.Size(sFunction{})
// 		fnsOffset := int(dsMod.FunctionsOffset) * fnSize
// 		// fns contains ALL the functions. we need to do a subset
// 		modFns := make(map[string]*CXFunction, 0)
// 		for i := 0; i < int(dsMod.FunctionsSize); i++ {
// 			//fn := CXFunction{}
			
// 			var dsFn sFunction
// 			sFn := dsFns[fnsOffset + i*fnSize : fnsOffset + (i+1)*fnSize]
// 			encoder.DeserializeRaw(sFn, &dsFn)

// 			// Name
// 			var dsName []byte
// 			sName := dsNames[dsFn.NameOffset : dsFn.NameOffset + dsFn.NameSize]
// 			encoder.DeserializeRaw(sName, &dsName)

// 			fn := fns[fmt.Sprintf("%s.%s", string(dsModName), string(dsName))]

// 			// Adding current function to module
// 			if int(dsMod.CurrentFunctionOffset) * fnSize == fnsOffset + i*fnSize {
// 				mod.CurrentFunction = fn
// 			}

// 			// Inputs
// 			var inps []*CXParameter
// 			paramSize := encoder.Size(sParameter{})
// 			typSize := encoder.Size(sType{})
// 			inpsOffset := int(dsFn.InputsOffset) * paramSize
// 			for i := 0; i < int(dsFn.InputsSize); i++ {
// 				inp := CXParameter{}

// 				var dsParam sParameter
// 				sParam := dsParams[inpsOffset + i*paramSize : inpsOffset + (i+1)*paramSize]
// 				encoder.DeserializeRaw(sParam, &dsParam)

// 				// Name
// 				var dsName []byte
// 				sName := dsNames[dsParam.NameOffset : dsParam.NameOffset + dsParam.NameSize]
// 				encoder.DeserializeRaw(sName, &dsName)

// 				// Type
// 				typ := CXType{}
				
// 				var dsTyp sType
// 				sTyp := dsTyps[dsParam.TypOffset*int32(typSize) : dsParam.TypOffset*int32(typSize) + int32(typSize)]
// 				encoder.DeserializeRaw(sTyp, &dsTyp)

// 				// Type name
// 				var dsTypName []byte
// 				sTypName := dsNames[dsTyp.NameOffset : dsTyp.NameOffset + dsTyp.NameSize]
// 				encoder.DeserializeRaw(sTypName, &dsTypName)
// 				typ.Name = string(dsTypName)
				
// 				inp.Name = string(dsName)
// 				inp.Typ = &typ

// 				// Appending final input
// 				inps = append(inps, &inp)
// 			}

// 			// Outputs
// 			var outs []*CXParameter
// 			outsOffset := int(dsFn.OutputsOffset) * paramSize
// 			for i := 0; i < int(dsFn.OutputsSize); i++ {
// 				out := CXParameter{}

// 				var dsParam sParameter
// 				sParam := dsParams[outsOffset + i*paramSize : outsOffset + (i+1)*paramSize]
// 				encoder.DeserializeRaw(sParam, &dsParam)

// 				// Name
// 				var dsName []byte
// 				sName := dsNames[dsParam.NameOffset : dsParam.NameOffset + dsParam.NameSize]
// 				encoder.DeserializeRaw(sName, &dsName)

// 				// Type
// 				typ := CXType{}
				
// 				var dsTyp sType
// 				sTyp := dsTyps[dsParam.TypOffset*int32(typSize) : dsParam.TypOffset*int32(typSize) + int32(typSize)]
// 				encoder.DeserializeRaw(sTyp, &dsTyp)

// 				// Type name
// 				var dsTypName []byte
// 				sTypName := dsNames[dsTyp.NameOffset : dsTyp.NameOffset + dsTyp.NameSize]
// 				encoder.DeserializeRaw(sTypName, &dsTypName)
// 				typ.Name = string(dsTypName)
				
// 				out.Name = string(dsName)
// 				out.Typ = &typ

// 				// Appending final output
// 				outs = append(outs, &out)
// 			}

// 			// // Output
// 			// out := &CXParameter{}
// 			// outOffset := int(dsFn.OutputOffset) * paramSize

// 			// var dsOut sParameter
// 			// sOut := dsParams[outOffset : outOffset + paramSize]
// 			// encoder.DeserializeRaw(sOut, &dsOut)

// 			// // Output name
// 			// var dsOutName []byte
// 			// sOutName := dsNames[dsOut.NameOffset : dsOut.NameOffset + dsOut.NameSize]
// 			// encoder.DeserializeRaw(sOutName, &dsOutName)

// 			// // Output Type
// 			// typ := CXType{}
			
// 			// var dsTyp sType
// 			// sTyp := dsTyps[dsOut.TypOffset*int32(typSize) : dsOut.TypOffset*int32(typSize) + int32(typSize)]
// 			// encoder.DeserializeRaw(sTyp, &dsTyp)

// 			// // Type name
// 			// var dsTypName []byte
// 			// sTypName := dsNames[dsTyp.NameOffset : dsTyp.NameOffset + dsTyp.NameSize]
// 			// encoder.DeserializeRaw(sTypName, &dsTypName)
// 			// typ.Name = string(dsTypName)

// 			// out.Name = string(dsOutName)
// 			// out.Typ = &typ

// 			// Expressions
// 			var exprs []*CXExpression
// 			exprSize := encoder.Size(sExpression{})
// 			exprsOffset := int(dsFn.ExpressionsOffset) * exprSize

// 			// Current expression
// 			var dsCurrExpr sExpression
// 			sCurrExpr := dsExprs[int(dsFn.CurrentExpressionOffset)*exprSize : int(dsFn.CurrentExpressionOffset)*exprSize + exprSize]
// 			encoder.DeserializeRaw(sCurrExpr, &dsCurrExpr)
			
// 			for i := 0; i < int(dsFn.ExpressionsSize); i++ {
// 				expr := CXExpression{}
				
// 				var dsExpr sExpression
// 				sExpr := dsExprs[exprsOffset + i*exprSize : exprsOffset + (i+1)*exprSize]
// 				encoder.DeserializeRaw(sExpr, &dsExpr)

// 				// Adding current expression to function
// 				if int(dsFn.CurrentExpressionOffset) * exprSize == exprsOffset + i*exprSize {
// 					fn.CurrentExpression = &expr
// 				}

// 				// Operator
// 				opSize := int32(encoder.Size(sFunction{}))

// 				var dsOp sFunction
// 				sOp := dsFns[dsExpr.OperatorOffset*opSize : dsExpr.OperatorOffset*opSize + opSize]
// 				encoder.DeserializeRaw(sOp, &dsOp)

// 				// Operator's name
// 				var dsOpName []byte
// 				sOpName := dsNames[dsOp.NameOffset : dsOp.NameOffset + dsOp.NameSize]
// 				encoder.DeserializeRaw(sOpName, &dsOpName)
				
// 				// Arguments
// 				var args []*CXArgument
// 				argSize := encoder.Size(sArgument{})
// 				argsOffset := int(dsExpr.ArgumentsOffset) * argSize
// 				for i := 0; i < int(dsExpr.ArgumentsSize); i++ {
// 					arg := CXArgument{}

// 					var dsArg sArgument
// 					sArg := dsArgs[argsOffset + i*argSize : argsOffset + (i+1)*argSize]
// 					encoder.DeserializeRaw(sArg, &dsArg)

// 					// Argument type
// 					typ := CXType{}

// 					var dsTyp sType
// 					sTyp := dsTyps[dsArg.TypOffset*int32(typSize) : dsArg.TypOffset*int32(typSize) + int32(typSize)]
// 					encoder.DeserializeRaw(sTyp, &dsTyp)

// 					// Type name
// 					var dsTypName []byte
// 					sTypName := dsNames[dsTyp.NameOffset : dsTyp.NameOffset + dsTyp.NameSize]
// 					encoder.DeserializeRaw(sTypName, &dsTypName)
// 					typ.Name = string(dsTypName)

// 					// Argument value
// 					var dsValue []byte
// 					sVal := dsValues[dsArg.ValueOffset : dsArg.ValueOffset + dsArg.ValueSize]
// 					encoder.DeserializeRaw(sVal, &dsValue)

// 					arg.Typ = &typ
// 					arg.Value = &dsValue
// 					arg.Offset = -1
// 					arg.Size = -1

// 					// Appending final argument
// 					args = append(args, &arg)
// 				}


// 				// Expression output names
// 				var outNames []string
// 				outNameSize := encoder.Size(sOutputName{})
// 				outNamesOffset := int(dsExpr.OutputNamesOffset) * outNameSize
// 				for i := 0; i < int(dsExpr.OutputNamesSize); i++ {
// 					var outName string

// 					var dsOutName sOutputName
// 					sOutName := dsOutNames[outNamesOffset + i*outNameSize : outNamesOffset + (i+1)*outNameSize]
// 					encoder.DeserializeRaw(sOutName, &dsOutName)

// 					var dsName []byte
// 					sName := dsNames[dsOutName.NameOffset : dsOutName.NameOffset + dsOutName.NameSize]
// 					encoder.DeserializeRaw(sName, &dsName)
// 					outName = string(dsName)

// 					// Appending final output name
// 					outNames = append(outNames, outName)
// 				}

// 				// // Expression output name
// 				// var dsOutName []byte
// 				// sOutName := dsNames[dsExpr.OutputNameOffset : dsExpr.OutputNameOffset + dsExpr.OutputNameSize]
// 				// encoder.DeserializeRaw(sOutName, &dsOutName)

// 				expr.Operator = fns[fmt.Sprintf("%s.%s", dsModName, dsOpName)]
// 				expr.Arguments = args
// 				//expr.OutputName = string(dsOutName)
// 				expr.OutputNames = outNames
// 				expr.Line = int(dsExpr.Line)
// 				expr.Function = fn
// 				expr.Context = &cxt
// 				expr.Module = mod
				
// 				// Appending final expression
// 				exprs = append(exprs, &expr)
// 			}

			

// 			// Constructing final function
// 			fn.Name = string(dsName)
// 			fn.Inputs = inps
// 			fn.Outputs = outs
// 			fn.Expressions = exprs
// 			// Current expression was added in the expression's loop
// 			fn.Module = mod

// 			// Appending final function to modFns
// 			modFns[string(dsName)] = fn
// 		}

// 		// Structs
// 		strcts := make(map[string]*CXStruct, 0)
// 		strctSize := encoder.Size(sStruct{})
// 		strctsOffset := int(dsMod.StructsOffset) * strctSize

// 		// Current struct
// 		var dsCurrStrct sStruct
// 		sCurrStrct := dsStrcts[int(dsMod.CurrentStructOffset)*strctSize : int(dsMod.CurrentStructOffset)*strctSize + strctSize]
// 		encoder.DeserializeRaw(sCurrStrct, &dsCurrStrct)
		
// 		for i := 0; i < int(dsMod.StructsSize); i++ {
// 			strct := CXStruct{}
			
// 			var dsStrct sStruct
// 			sStrct := dsStrcts[strctsOffset + i*strctSize : strctsOffset + (i+1)*strctSize]
// 			encoder.DeserializeRaw(sStrct, &dsStrct)

// 			// Adding current struct to module
// 			if int(dsMod.CurrentStructOffset) * strctSize == strctsOffset + i*strctSize {
// 				mod.CurrentStruct = &strct
// 			}


// 			// Struct name
// 			var dsName []byte
// 			sName := dsNames[dsStrct.NameOffset : dsStrct.NameOffset + dsStrct.NameSize]
// 			encoder.DeserializeRaw(sName, &dsName)

// 			// Struct fields
// 			var flds []*CXField
// 			fldSize := encoder.Size(sField{})
// 			fldsOffset := int(dsStrct.FieldsOffset) * fldSize
// 			for i := 0; i < int(dsStrct.FieldsSize); i++ {
// 				fld := CXField{}

// 				var dsFld sField
// 				sFld := dsFlds[fldsOffset + i*fldSize : fldsOffset + (i+1)*fldSize]
// 				encoder.DeserializeRaw(sFld, &dsFld)

// 				// Field name
// 				var dsName []byte
// 				sName := dsNames[dsFld.NameOffset : dsFld.NameOffset + dsFld.NameSize]
// 				encoder.DeserializeRaw(sName, &dsName)

// 				// Field type
// 				typ := CXType{}

// 				var dsTyp sType
// 				typSize := encoder.Size(sType{})
// 				sTyp := dsTyps[dsFld.TypOffset*int32(typSize) : dsFld.TypOffset*int32(typSize) + int32(typSize)]
// 				encoder.DeserializeRaw(sTyp, &dsTyp)

// 				// Type name
// 				var dsTypName []byte
// 				sTypName := dsNames[dsTyp.NameOffset : dsTyp.NameOffset + dsTyp.NameSize]
// 				encoder.DeserializeRaw(sTypName, &dsTypName)
// 				typ.Name = string(dsTypName)
				
// 				fld.Name = string(dsName)
// 				fld.Typ = &typ

// 				// Appending final field
// 				flds = append(flds, &fld)
// 			}


			
// 			strct.Name = string(dsName)
// 			strct.Fields = flds
// 			strct.Module = mod
// 			strct.Context = &cxt

// 			// Appending final struct
// 			strcts[string(dsName)] = &strct
// 		}

// 		// Definitions
// 		defs := make(map[string]*CXDefinition, 0)
// 		defSize := encoder.Size(sDefinition{})
// 		defsOffset := int(dsMod.DefinitionsOffset) * defSize
// 		for i := 0; i < int(dsMod.DefinitionsSize); i++ {
// 			def := CXDefinition{}
			
// 			var dsDef sDefinition
// 			sDef := dsDefs[defsOffset + i*defSize : defsOffset + (i+1)*defSize]
// 			encoder.DeserializeRaw(sDef, &dsDef)

// 			// Definition name
// 			var dsName []byte
// 			sName := dsNames[dsDef.NameOffset : dsDef.NameOffset + dsDef.NameSize]
// 			encoder.DeserializeRaw(sName, &dsName)

// 			// Definition type
// 			typ := CXType{}

// 			var dsTyp sType
// 			typSize := encoder.Size(sType{})
// 			sTyp := dsTyps[dsDef.TypOffset*int32(typSize) : dsDef.TypOffset*int32(typSize) + int32(typSize)]
// 			encoder.DeserializeRaw(sTyp, &dsTyp)

// 			// Type name
// 			var dsTypName []byte
// 			sTypName := dsNames[dsTyp.NameOffset : dsTyp.NameOffset + dsTyp.NameSize]
// 			encoder.DeserializeRaw(sTypName, &dsTypName)
// 			typ.Name = string(dsTypName)

// 			// Definition value
// 			var dsValue []byte
// 			sVal := dsValues[dsDef.ValueOffset : dsDef.ValueOffset + dsDef.ValueSize]
// 			encoder.DeserializeRaw(sVal, &dsValue)

// 			def.Name = string(dsName)
// 			def.Typ = &typ
// 			def.Value = &dsValue
// 			def.Module = mod
// 			def.Context = &cxt

// 			// Appending final definition
// 			defs[string(dsName)] = &def
// 		}
		
// 		mod.Name = string(dsModName)
// 		mod.Imports = imps
// 		mod.Functions = modFns
// 		mod.Structs = strcts
// 		mod.Definitions = defs
// 		mod.Context = &cxt
// 		//fmt.Println(string(dsName))
// 	}

// 	// Call stack
// 	calls := make([]*CXCall, 0)
// 	callSize := encoder.Size(sCall{})
	
// 	// this will always be 0. I'll leave it in the struct for consistency (like modulesoffset)
// 	//callsOffset := int(dsCxt.CallStackOffset) * callSize

// 	var lastCall *CXCall
// 	for i := 0; i < int(dsCxt.CallStackSize); i++ {
// 		call := CXCall{}

// 		var dsCall sCall
// 		sCal := dsCalls[i*callSize:(i+1)*callSize]
// 		encoder.DeserializeRaw(sCal, &dsCall)

// 		// Call's module
// 		var dsMod sModule
// 		modSize := encoder.Size(sModule{})
// 		sMod := dsMods[int(dsCall.ModuleOffset)*modSize : int(dsCall.ModuleOffset)*modSize + modSize]
// 		encoder.DeserializeRaw(sMod, &dsMod)

// 		// Call's module name
// 		var dsModName []byte
// 		sModName := dsNames[dsMod.NameOffset:dsMod.NameOffset + dsMod.NameSize]
// 		encoder.DeserializeRaw(sModName, &dsModName)

// 		mod := mods[string(dsModName)]

// 		// Call's operator
// 		opSize := int32(encoder.Size(sFunction{}))

// 		var dsOp sFunction
// 		sOp := dsFns[dsCall.OperatorOffset*opSize : dsCall.OperatorOffset*opSize + opSize]
// 		encoder.DeserializeRaw(sOp, &dsOp)

// 		// Operator's name
// 		var dsOpName []byte
// 		sOpName := dsNames[dsOp.NameOffset : dsOp.NameOffset + dsOp.NameSize]
// 		encoder.DeserializeRaw(sOpName, &dsOpName)

// 		// Call's Operator Module
// 		var dsOpMod sModule
// 		sOpMod := dsMods[int(dsOp.ModuleOffset)*modSize : int(dsOp.ModuleOffset)*modSize + modSize]
// 		encoder.DeserializeRaw(sOpMod, &dsOpMod)

// 		// Call's Operator Module Name
// 		var dsOpModName []byte
// 		sOpModName := dsNames[dsOpMod.NameOffset:dsOpMod.NameOffset + dsOpMod.NameSize]
// 		encoder.DeserializeRaw(sOpModName, &dsOpModName)

// 		for _, mod := range mods {
// 			for _, fn := range mod.Functions {
// 				if fn.Name == string(dsOpName) && mod.Name == string(dsOpModName) {
// 					call.Operator = fn
// 				}
// 			}
// 		}

// 		// State
// 		defs := make(map[string]*CXDefinition, 0)
// 		defSize := encoder.Size(sDefinition{})
// 		defsOffset := int(dsCall.StateOffset) * defSize
// 		for i := 0; i < int(dsCall.StateSize); i++ {
// 			def := CXDefinition{}
			
// 			var dsDef sDefinition
// 			sDef := dsDefs[defsOffset + i*defSize : defsOffset + (i+1)*defSize]
// 			encoder.DeserializeRaw(sDef, &dsDef)

// 			// Definition name
// 			var dsName []byte
// 			sName := dsNames[dsDef.NameOffset : dsDef.NameOffset + dsDef.NameSize]
// 			encoder.DeserializeRaw(sName, &dsName)

// 			// Definition type
// 			typ := CXType{}

// 			var dsTyp sType
// 			typSize := encoder.Size(sType{})
// 			sTyp := dsTyps[dsDef.TypOffset*int32(typSize) : dsDef.TypOffset*int32(typSize) + int32(typSize)]
// 			encoder.DeserializeRaw(sTyp, &dsTyp)

// 			// Type name
// 			var dsTypName []byte
// 			sTypName := dsNames[dsTyp.NameOffset : dsTyp.NameOffset + dsTyp.NameSize]
// 			encoder.DeserializeRaw(sTypName, &dsTypName)
// 			typ.Name = string(dsTypName)

// 			// Definition value
// 			var dsValue []byte
// 			sVal := dsValues[dsDef.ValueOffset : dsDef.ValueOffset + dsDef.ValueSize]
// 			encoder.DeserializeRaw(sVal, &dsValue)

// 			def.Name = string(dsName)
// 			def.Typ = &typ
// 			def.Value = &dsValue
// 			def.Module = mod
// 			def.Context = &cxt

// 			// Appending final definition
// 			defs[string(dsName)] = &def
// 		}

// 		call.Line = int(dsCall.Line)
// 		call.State = defs
// 		call.ReturnAddress = lastCall
// 		call.Module = mod
// 		call.Context = &cxt

// 		lastCall = &call

// 		// Appending final call
// 		calls = append(calls, &call)
// 	}

// 	//fmt.Println(calls)

// 	callStack := CXCallStack{}
// 	callStack.Calls = calls

// 	cxt.Modules = mods
// 	cxt.CallStack = &callStack
// 	cxt.Steps = make([]*CXCallStack, 0)



// 	return &cxt
// }
