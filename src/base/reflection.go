package base

import (
	
)

// What should we return?

func (strct *cxStruct) GetFields() []*cxField {
	return strct.Fields
}

func (mod *cxModule) GetFunctions() []*cxFunction {
	funcs := make([]*cxFunction, len(mod.Functions))
	i := 0
	for _, v := range mod.Functions {
		funcs[i] = v
		i++
	}
	return funcs
}
