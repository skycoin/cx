package base

import (
	"fmt"
	"unsafe"
	. "github.com/skycoin/skycoin/src/cipher/encoder"
)

// We might need to reorder the fields to have a better ordering in the program byte array



/*
  Context
*/

type sContext struct {
	ModulesOffset int32
	ModulesSize int32
	CurrentModuleOffset int32
	CallStackOffset int32
	StepsOffset int32
	StepsSize int32
	// we can't serialize steps because of skycoin/encoder limitations
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
	ModuleSize int32
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
	OutputOffset int32
	ExpressionsOffset int32
	ExpressionsSize int32
	CurrentExpressionOffest int32
	ModuleOffset int32
}

// I'll try creating the structures inline
// func SMakeFunction (fn *cxFunction) *sFunction {
// 	return &sFunction{
// 		NameOffset: ,
// 		NameSize: ,
// 		InputsOffset 
// 		InputsSize 
// 		OutputOffset 
// 		ExpressionsOffset 
// 		ExpressionsSize 
// 		CurrentExpressionOffest 
// 		ModuleOffset 
// 	}
// }

type sParameter struct {
	NameOffset int32
	NameSize int32
	TypOffset int32
}

type sExpression struct {
	OperatorOffset int32
	ArgumentsOffset int32
	ArgumentsSize int32
	OutputNameOffset int32 // these are also going to the names byte array
	Line int32
	FunctionOffset int32
	ModuleOffset int32
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


func Testing () {
	test := sContext{
		ModulesOffset: 5,
		ModulesSize: 10,
		CurrentModuleOffset: 3,
		CallStackOffset: 4,
		StepsOffset: 313,
		StepsSize: 11,
		}
	fmt.Println("Testing size")
	fmt.Println(unsafe.Sizeof(test))
	// fmt.Println(unsafe.Sizeof(Serialize(test)))
	// fmt.Println(Size(test))
	// fmt.Println(Size(Serialize(test)))

	fmt.Println("Testing serializing byte arrays (what happens?)")
	fmt.Println(Size([]byte("hello")))
	fmt.Println(Size(Serialize([]byte("hello"))))
	fmt.Println(Serialize(Serialize([]byte("hello"))))
	fmt.Println(Serialize("hello"))
	fmt.Println([]byte("hello"))
}

/*
  Makers
*/

func SMakeContext (cxt *cxContext) *sContext {
	return &sContext{
		ModulesOffset: 24,
		ModulesSize: int32(len(cxt.Modules)),
		CurrentModuleOffset: -1,
		CallStackOffset: -1,
		StepsOffset: -1,
		StepsSize: -1,
	}
}

func SMakeArgument () {
	
}


// okay, once we create, for example, a name, how do we keep track of what was its index
// we can use the auxiliary var "modsIndexes"


func SMakeName (name string) *[]byte {
	byts := []byte(name)
	return &byts
}


// we first need this function
func SGetContext (prgrm *[]byte) *sContext {
	cxt := &sContext{}
	// The first 24 bytes are the context
	
	

	return cxt
}

// we're going to have makers, setters, and getters

func SSetContext (sCxt *sContext, prgrm *[]byte) *[]byte {
	// okay, I think we need to leave this for later
	foo := make([]byte, 24)

	return &foo
}

func SGetModules (cxt *[]byte, offset, size int) *[]byte {
	mods := make([]byte, 0)
	
	return &mods
}



// I don't find the point here with these getters
// Can't we just have one?



func SGetModule (mods *[]byte, offset, size int) *[]byte {
	mod := make([]byte, 0)

	return &mod
}

func SSetModules (cxt *cxContext, offset, size int) {
	// we need to think a bit more about the setters
}







func (cxt *cxContext) Serialize () *[]byte {
	// we will be appending the bytes here
	sCxt := make([]byte, 0)
	
	sNames := make([]byte, 0)
	sNamesCounter := 0
	sNamesMap := make(map[string]int, 0)
	
	sValues := make([]byte, 0)
	sValuesCounter := 0
	// we don't use a map for values, as these can change, can be huge, and other stuff
	
	sMods := make([]byte, 0)
	sModsCounter := 0

	sFns := make([]byte, 0)
	sFnsCounter := 0
	sFnsMap := make(map[string]int, 0)

	sTyps := make([]byte, 0)
	sTypsCounter := 0

	sParams := make([]byte, 0)
	sParamsCounter := 0

	sExprs := make([]byte, 0)
	sExprsCounter := 0

	sArgs := make([]byte, 0)
	sArgsCounter := 0
	
	// context
	mods := cxt.Modules


	//currentModule := cxt.CurrentModule // constant
	// we could send currentModule to the end, although I don't think it matters
	//callStack := cxt.CallStack // not constant
	//steps := cxt.Steps // not constant
	//programSteps := cxt.ProgramSteps // not constant

	// we serialize modules

	// adding function names first so expressions can reference their operators
	for _, mod := range mods {
		for _, fn := range mod.Functions {
			fnName := fmt.Sprintf("%s.%s", mod.Name, fn.Name)
			sFnsMap[fnName] = sFnsCounter
			sFnsCounter++
		}
	}


	for _, mod := range mods {
		defs := mod.Definitions
		fns := mod.Functions
		
		// for _, def := range defs {
		// 	sName := []byte(def.Name)
		// 	sNames = append(sNames, sName...)
		// }


		// we need to order the functions in a way that
		// their expressions aren't calling themselves (I think we can do this)
		// and they are not calling functions which
		// ahhh, wait, we might be able to
		// yeah, we need to add *their names* to the sFnsMap, and then we proceed


		for _, fn := range fns {
			sFn := sFunction{}

			// name
			if offset, ok := sNamesMap[fn.Name]; ok {
				sFn.NameOffset = int32(offset)
				sFn.NameSize = int32(Size(fn.Name))
			} else {
				sNames = append(sNames, Serialize(fn.Name)...)
				sFn.NameOffset = int32(sNamesCounter)
				sFn.NameSize = int32(Size(fn.Name))
				sNamesMap[fn.Name] = sNamesCounter
				sNamesCounter++
			}

			// inputs
			if fn.Inputs != nil && len(fn.Inputs) > 0 {
				sFn.InputsOffset = int32(sParamsCounter)
				sFn.InputsSize = int32(len(fn.Inputs))
				
				for _, inp := range fn.Inputs {
					sParam := sParameter{}

					// input name
					if offset, ok := sNamesMap[inp.Name]; ok {
						sParam.NameOffset = int32(offset)
						sParam.NameSize = int32(Size(inp.Name))
					} else {
						sNames = append(sNames, Serialize(inp.Name)...)
						sNamesMap[inp.Name] = sNamesCounter
						sParam.NameOffset = int32(sNamesCounter)
						sParam.NameSize = int32(Size(inp.Name))
						sNamesCounter++
						
					}

					// input type
					typ := inp.Typ
					sTyp := sType{}
					
					// input type name
					if offset, ok := sNamesMap[typ.Name]; ok {
						sTyp.NameOffset = int32(offset)
						sTyp.NameSize = int32(Size(typ.Name))
					} else {
						sNames = append(sNames, Serialize(typ.Name)...)
						sNamesMap[typ.Name] = sNamesCounter
						sTyp.NameOffset = int32(sNamesCounter)
						sTyp.NameSize = int32(Size(typ.Name))
						sNamesCounter++
					}

					// save the sTyp
					sTyps = append(sTyps, Serialize(sTyp)...)
					sParam.TypOffset = int32(sTypsCounter)
					sTypsCounter++

					// save the sParam
					sParams = append(sParams, Serialize(sParam)...)
					sParamsCounter++
				}
			} else {
				sFn.InputsOffset = -1 // nil; fn does not have inputs
				sFn.InputsSize = -1
			}

			// output
			if fn.Output != nil {
				sFn.OutputOffset = int32(sParamsCounter)
				sParam := sParameter{}

				// name
				if offset, ok := sNamesMap[fn.Output.Name]; ok {
					sParam.NameOffset = int32(offset)
					sParam.NameSize = int32(Size(fn.Output.Name))
				} else {
					sNames = append(sNames, Serialize(fn.Output.Name)...)
					sNamesMap[fn.Output.Name] = sNamesCounter
					sParam.NameOffset = int32(sNamesCounter)
					sParam.NameSize = int32(Size(fn.Output.Name))
					sNamesCounter++
				}

				// output type
				typ := fn.Output.Typ
				sTyp := sType{}
				
				// output type name
				if offset, ok := sNamesMap[typ.Name]; ok {
					sTyp.NameOffset = int32(offset)
					sTyp.NameSize = int32(Size(typ.Name))
				} else {
					sNames = append(sNames, Serialize(typ.Name)...)
					sNamesMap[typ.Name] = sNamesCounter
					sTyp.NameOffset = int32(sNamesCounter)
					sTyp.NameSize = int32(Size(typ.Name))
					sNamesCounter++
				}

				// save the sTyp
				sTyps = append(sTyps, Serialize(sTyp)...)
				sParam.TypOffset = int32(sTypsCounter)
				sTypsCounter++

				// save the sParam
				sParams = append(sParams, Serialize(sParam)...)
				sParamsCounter++
			} else {
				sFn.OutputOffset = -1 // nil; fn does not have an output
			}

			// expressions
			if fn.Expressions != nil && len(fn.Expressions) > 0 {
				sFn.ExpressionsOffset = int32(sExprsCounter)
				sFn.ExpressionsSize = int32(len(fn.Expressions))
				
				for _, expr := range fn.Expressions {
					sExpr := sExpression{}
					opName := fmt.Sprintf("%s.%s", expr.Operator.Module.Name, expr.Operator.Name)

					// operator
					if offset, ok := sFnsMap[opName]; ok {
						sExpr.OperatorOffset = int32(offset)
					} else {
						panic(fmt.Sprintf("Expression's operator (%s) not found in sFnsMap", opName))
					}

					// arguments
					for _, arg := range expr.Arguments {
						sArg := sArgument{}
						
						// arg type
						typ := arg.Typ
						sTyp := sType{}
						
						// argument type name
						if offset, ok := sNamesMap[typ.Name]; ok {
							sTyp.NameOffset = int32(offset)
							sTyp.NameSize = int32(Size(typ.Name))
						} else {
							sNames = append(sNames, Serialize(typ.Name)...)
							sNamesMap[typ.Name] = sNamesCounter
							sTyp.NameOffset = int32(sNamesCounter)
							sTyp.NameSize = int32(Size(typ.Name))
							sNamesCounter++
						}

						// save the argument sTyp
						sTyps = append(sTyps, Serialize(sTyp)...)
						sArg.TypOffset = int32(sTypsCounter)
						sTypsCounter++

						// argument value

						sValues = append(sValues, Serialize(*arg.Value)...)
						// we are serializing anyway or we won't get the correct size with encoder.Size()
						sArg.ValueOffset = int32(sValuesCounter)
						sArg.ValueSize = int32(Size(*arg.Value))
						sValuesCounter++

						// save the sArg
						sArgs = append(sArgs, Serialize(sArg)...)
						sArgsCounter++
					}
				}
			}
		}
	}
	
	// field

	return &sCxt
}
