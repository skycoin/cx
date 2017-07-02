package base

import (
	"fmt"
)

var counter int = 0
func MakeGenSym (name string) string {
	gensym := fmt.Sprintf("%s%d", name, counter)
	counter++
	
	return gensym
}

func MakeContext () *cxContext {
	return &cxContext{}
}

func MakeModule (name string) *cxModule {
	return &cxModule{
		Name: name,
		Definitions: make(map[string]*cxDefinition, 0),
		Imports: make(map[string]*cxModule, 0),
		Functions: make(map[string]*cxFunction, 0),
		Structs: make(map[string]*cxStruct, 0),
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

func MakeParameter(name string, typ *cxType) *cxParameter {
	return &cxParameter{Name: name,
		Typ: typ}
}

func MakeExpression (fn *cxFunction) *cxExpression {
	return &cxExpression{Operator: fn}
}

func MakeArgument(value *[]byte, typ *cxType) *cxArgument {
	return &cxArgument{Typ: typ, Value: value}
}

func MakeType(name string) *cxType {
	return &cxType{Name: name}
}

func MakeFunction(name string) *cxFunction {
	return &cxFunction{Name: name}
}
