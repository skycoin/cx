package opcodes

import (
	"fmt"

	"github.com/skycoin/cx/cx/ast"
	"github.com/skycoin/cx/cx/globals"
	"github.com/skycoin/cx/cx/types"
)

// RegisterPackage registers a package on the CX standard library. This does not create a `CXPackage` structure,
// it only tells the CX runtime that `pkgName` will exist by the time a CX program is run.
//func RegisterPackage(pkgName string) {
//	constants.CorePackages = append(constants.CorePackages, pkgName)
//}

// GetOpCodeCount returns an op code that is available for usage on the CX standard library.
/*
func GetOpCodeCount() int {
	return len(OpcodeHandlers)
}
*/

// RegisterOpCode ...
func RegisterOpCode(code int, name string, handler ast.OpcodeHandler, inputs []*ast.CXArgument, outputs []*ast.CXArgument) {
	if code >= len(ast.OpcodeHandlers) {
		ast.OpcodeHandlers = append(ast.OpcodeHandlers, make([]ast.OpcodeHandler, code+1)...)
	}
	if ast.OpcodeHandlers[code] != nil {
		panic(fmt.Sprintf("duplicate opcode %d : '%s' width '%s'.\n", code, name, ast.OpNames[code]))
	}
	ast.OpcodeHandlers[code] = handler

	ast.OpNames[code] = name
	ast.OpCodes[name] = code
	//OpVersions[code] = 2

	if inputs == nil {
		inputs = []*ast.CXArgument{}
	}
	if outputs == nil {
		outputs = []*ast.CXArgument{}
	}
	ast.Natives[code] = MakeNativeFunction(code, inputs, outputs)
}

// RegisterFunction ...
func RegisterFunction(name string, handler ast.OpcodeHandler, inputs []*ast.CXArgument, outputs []*ast.CXArgument) {
	RegisterOpCode(globals.OpCodeSystemCounter, name, handler, inputs, outputs)
	globals.OpCodeSystemCounter++
}

// RegisterOperator ...
func RegisterOperator(name string, handler ast.OpcodeHandler, inputs []*ast.CXArgument, outputs []*ast.CXArgument, atomicType types.Code, operator int) {
	RegisterOpCode(globals.OpCodeSystemCounter, name, handler, inputs, outputs)
	native := ast.Natives[globals.OpCodeSystemCounter]
	ast.Operators[ast.GetTypedOperatorOffset(atomicType, operator)] = native
	globals.OpCodeSystemCounter++
}

// MakeNativeFunction ...
func MakeNativeFunction(opCode int, inputs []*ast.CXArgument, outputs []*ast.CXArgument) *ast.CXFunction {
	fn := &ast.CXFunction{
		AtomicOPCode: opCode,
		Index:        -1,
	}

	offset := types.Pointer(0)
	for _, inp := range inputs {
		inp.Offset = offset
		offset.Add(ast.GetSize(inp))
		fn.Inputs = append(fn.Inputs, inp)
	}
	for _, out := range outputs {
		fn.Outputs = append(fn.Outputs, out)
		out.Offset = offset
		offset.Add(ast.GetSize(out))
	}

	return fn
}

func opDebugPrintStack(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	prgrm.PrintStack()
}
