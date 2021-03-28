package ast

import "fmt"

//not used
func MakeNativeFunctionV2(opCode int, inputs []*CXArgument, outputs []*CXArgument) *CXFunction {
	fn := &CXFunction{
		IsAtomic: true,
		OpCode:   opCode,
		Version:  2,
	}

	offset := 0
	for _, inp := range inputs {
		inp.Offset = offset
		offset += GetSize(inp)
		fn.Inputs = append(fn.Inputs, inp)
	}
	for _, out := range outputs {
		fn.Outputs = append(fn.Outputs, out)
		out.Offset = offset
		offset += GetSize(out)
	}

	return fn
}

// Op ...
func Op_V2(code int, name string, handler OpcodeHandler_V2, inputs []*CXArgument, outputs []*CXArgument) {
	if code >= len(OpcodeHandlers_V2) {
		OpcodeHandlers_V2 = append(OpcodeHandlers_V2, make([]OpcodeHandler_V2, code+1)...)
	}
	if OpcodeHandlers_V2[code] != nil {
		panic(fmt.Sprintf("duplicate opcode %d : '%s' width '%s'.\n", code, name, OpNames[code]))
	}
	OpcodeHandlers_V2[code] = handler

	OpNames[code] = name
	OpCodes[name] = code
	//OpVersions[code] = 2

	if inputs == nil {
		inputs = []*CXArgument{}
	}
	if outputs == nil {
		outputs = []*CXArgument{}
	}
	Natives[code] = MakeNativeFunctionV2(code, inputs, outputs)
}


