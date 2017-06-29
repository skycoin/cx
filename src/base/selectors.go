package base

import (
	//"github.com/skycoin/skycoin/src/cipher/encoder"
	//"fmt"
)



// This might be a good approach
// We are going to be stating on what module, function, struct, etc affordances must be occuring
// Affordances can also be SelectModule, let's do this now


func (cxt *cxContext) SelectModule (name string) *cxModule {
	var found *cxModule
	for _, mod := range cxt.Modules {
		if mod.Name == name {
			cxt.CurrentModule = mod
			found = mod
		}
	}

	if found == nil {
		// if not found, we return last current module
		// this is not the desired behaviour. we need to throw an error
		// the same applies for the other selectors
		found = cxt.CurrentModule
	}

	return found
}

func (mod *cxModule) SelectFunction (name string) *cxFunction {
	var found *cxFunction
	for _, fn := range mod.Functions {
		if fn.Name == name {
			mod.CurrentFunction = fn
			found = fn
		}
	}

	if found == nil {
		found = mod.CurrentFunction
	}

	return found
}

// Does this mean that these structures are the only ones
// that can have affordances??
// No, this means that these are "root" nodes that can have other structures
// which can have further affordances
// hmm, wait
// we can't have affordances (the ones we are focusing at the moment. we are going to have other types of affordances, like remove, change) of: cxType, cxField
// yes, this means these structures are the ones that can have the current type of
// affordances (adders)
// we won't have adders on cxTypes, cxFields, etc

func (mod *cxModule) SelectStruct (name string) *cxStruct {
	var found *cxStruct
	for _, strct := range mod.Structs {
		if strct.Name == name {
			mod.CurrentStruct = strct
			found = strct
		}
	}

	if found == nil {
		found = mod.CurrentStruct
	}

	return found
}
