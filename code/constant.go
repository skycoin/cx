package code

/*
	Opcode reprsents Opcode.
*/
type Opcode byte

const (
	OpConstant Opcode = iota

	OpAdd

	OpPop

	OpSub
	OpMul
	OpDiv

	OpTrue
	OpFalse

	OpEqual
	OpNotEqual
	OpGreaterThan

	OpMinus
	OpBang

	OpJumpNotTruthy
	OpJump

	OpNull

	OpGetGlobal
	OpSetGlobal

	OpArray
	OpHash
	OpIndex

	OpCall

	OpReturnValue
	OpReturn

	OpGetLocal
	OpSetLocal

	OpGetBuiltin

	OpClosure

	OpGetFree

	OpCurrentClosure
)

/*
	Definition represents opcode as name and  OperandWidths.
*/

type Definition struct {
	Name          string
	OperandWidths []int
}

var definitions = map[Opcode]*Definition{

	OpConstant: {"OpConstant", []int{2}},

	OpAdd: {"OpAdd", []int{}},

	OpPop: {"OpPop", []int{}},

	OpSub: {"OpSub", []int{}},
	OpMul: {"OpMul", []int{}},
	OpDiv: {"OpDiv", []int{}},

	OpTrue:  {"OpTrue", []int{}},
	OpFalse: {"OpFalse", []int{}},

	OpEqual:       {"OpEqual", []int{}},
	OpNotEqual:    {"OpNotEqual", []int{}},
	OpGreaterThan: {"OpGreaterThan", []int{}},

	OpMinus: {"OpMinus", []int{}},
	OpBang:  {"OpBang", []int{}},

	OpJumpNotTruthy: {"OpJumpNotTruthy", []int{2}},
	OpJump:          {"OpJump", []int{2}},

	OpNull: {"OpNull", []int{}},

	OpGetGlobal: {"OpGetGlobal", []int{2}},
	OpSetGlobal: {"OpSetGlobal", []int{2}},

	OpArray: {"OpArray", []int{2}},
	OpHash:  {"OpHash", []int{2}},
	OpIndex: {"OpIndex", []int{}},

	OpCall: {"OpCall", []int{1}},

	OpReturnValue: {"OpReturnValue", []int{}},
	OpReturn:      {"OpReturn", []int{}},

	OpGetLocal: {"OpGetLocal", []int{1}},
	OpSetLocal: {"OpSetLocal", []int{1}},

	OpGetBuiltin: {"OpGetBuiltin", []int{1}},

	OpClosure: {"OpClosure", []int{2, 1}},

	OpGetFree: {"OpGetFree", []int{1}},

	OpCurrentClosure: {"OpCurrentClosure", []int{}},
}
