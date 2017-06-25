package base

import (
	
)

func MakeModule (name string) *cxModule {
	return &cxModule{
		Name: name,
		Definitions: make(map[string]*cxDefinition),
		Imports: make(map[string]*cxModule),
		Functions: make(map[string]*cxFunction),
		Structs: make(map[string]*cxStruct),
	}
}

func MakeDefinition (name string, value *[]byte, typ *cxType) *cxDefinition {
	return &cxDefinition{Name: name, Typ: typ, Value: value}
}

func MakeField (name string, typ *cxType) *cxField {
	return &cxField{Name: name, Typ: typ}
}

func MakeStruct (name string) *cxStruct {
	return &cxStruct{Name: name}
}

func MakeContext () {
	
}

func MakeParameter(name string, typ *cxType) *cxParameter {
	return &cxParameter{Name: name,
		Typ: typ}
}

func MakeExpression (fn *cxFunction) *cxExpression {
	return &cxExpression{Fn: fn}
}

func MakeArgument(value *[]byte, typ *cxType) *cxArgument {
	return &cxArgument{Typ: typ, Value: value}
}

func MakeTypes(names []string) []*cxType {
	types := make([]*cxType, len(names))

	for i := 0; i < len(names); i++ {
		types = append(types, MakeType(names[i]))
	}
	
	return types
}

func MakeType(name string) *cxType {
	return &cxType{Name: name}
}

func MakeFunction(name string) *cxFunction {
	return &cxFunction{Name: name}
}
