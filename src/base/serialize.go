package base

import (
	//"github.com/skycoin/skycoin/src/cipher/encoder"
)

// We might need to reorder the fields to have a better ordering in the program byte array



/*
  Context
*/

type sContext struct {
	ModulesOffset int
	ModulesSize int
	CurrentModuleOffset int
	CallStackOffset int
	CallStackSize int
	

	// Steps are variable in size
	// So this is going to be a problem
	// It's a problem because they hold an array of arrays of calls

	// A possible solution would be to create a struct cxCallStack
	// but we'd need to make changes to the code, so let's try to not
	
	StepsOffest int
	StepsSize int
	// we can't serialize steps because of skycoin/encoder limitations
	//ProgramStepsOffset int
	//ProgramStepsSize int
}

type sCall struct {
	OperatorOffset int
	Line int
	StateOffset int
	StateSize int
	ReturnAddressOffset int
	//ContextOffset // this might not be a problem, because the context will always be at byte 0
	ModuleOffset int
	ModuleSize int
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
	NameOffset int
	NameSize int
	ImportsOffset int
	ImportsSize int
	FunctionsOffset int
	FunctionsSize int
	StructsOffset int
	StructsSize int
	DefinitionsOffset int
	DefinitionsSize int
	CurrentFunctionOffset int
	CurrentStructOffset int
}

type sDefinition struct {
	NameOffset int
	NameSize int
	TypOffset int
	ValueOffset int
	ValueSize int
	ModuleOffset int
}

/*
  Structs
*/

type sStruct struct {
	NameOffset int
	NameSize int
	FieldsOffset int
	FieldsSize int
	ModuleOffset int
}

type sField struct {
	NameOffset int
	NameSize int
	TypOffset int
}

type sType struct {
	NameOffset int
	NameSize int
}

/*
  Functions
*/

type sFunction struct {
	NameOffset int
	NameSize int
	InputsOffset int
	InputsSize int
	OutputOffset int
	ExpressionsOffset int
	ExpressionsSize int
	CurrentExpressionOffest int
	ModuleOffset int
}

type sParameter struct {
	NameOffset int
	NameSize int
	TypOffset int
}

type sExpression struct {
	OperatorOffset int
	ArgumentsOffset int
	ArgumentsSize int
	OutputNameOffset int // these are also going to the names byte array
	Line int
	FunctionOffset int
	ModuleOffset int
}

type sArgument struct {
	TypOffset int
	ValueOffset int
	ValueSize int
}

/*
  Affordances

  Affordances must not be serialized
*/



// we first need this function
func sGetContext (prgrm *[]byte) *sContext {
	return &sContext{}
}

func sGetModules (cxt *[]byte, offset, size int) *[]byte {
	mods := make([]byte, 0)
	
	return &mods
}



// I don't find the point here with these getters
// Can't we just have one?



func sGetModule (mods *[]byte, offset, size int) *[]byte {
	mod := make([]byte, 0)

	return &mod
}

func sSetModules (cxt *cxContext, offset, size int) {
	// we need to think a bit more about the setters
}

// func (cxt *cxContext) Serialize () *[]byte {
// 	// we will be appending the bytes here
// 	sCxt := make([]byte, 0)
// 	sNames := make([]byte, 0)
// 	sValues := make([]byte, 0)
	
// 	sMods := make([]byte, 0)
	
// 	// context
// 	mods := cxt.Modules

// 	modsIndexes := make(map[string]int, 0)

// 	// serializing only the modules references
// 	for i, mod := range mods {
// 		// we can use a serialization and a map which temporarily holds the offset to each module
// 		modsIndexes[mod.Name] = i

// 		// now, we need a struct (or maybe we don't necessarily need a struct)
// 		// we need
// 	}

// 	currentModule := cxt.CurrentModule // constant
// 	// we could send currentModule to the end, although I don't think it matters
// 	callStack := cxt.CallStack // not constant
// 	steps := cxt.Steps // not constant
// 	programSteps := cxt.ProgramSteps // not constant

// 	// we serialize modules



// 	// we would need to start from the leaves

// 	// definitions
// 	for _, mod := range mods {
// 		defs := mod.Definitions
// 		for _, def := range defs {
// 			sName := []byte(def.Name)
// 			names = append(names, sName...)
// 			// at this ponit we got that name, but how can we retrieve it when we decide to deserialize
// 			//def.Typ.Name
// 			// these are going to be repeated a lot of times
// 			// so they're just like these:
// 			//def.Module
// 			//def.Context
// 			// we'd need reference to these
// 			// we can just reference the reference

			
// 		}
// 	}

// 	// arguments
	
// 	// inputs
// 	// outputs
	
// 	// field

// 	return &program
// }
