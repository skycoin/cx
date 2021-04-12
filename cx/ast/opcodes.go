package ast

import (
	"github.com/skycoin/cx/cx/constants"
)

var (
	// OpNames ...
	OpNames = map[int]string{}

	// OpCodes ...
	OpCodes = map[string]int{}
)

var (
	OpcodeHandlers    []OpcodeHandler
)

// OpcodeHandler ...
//TODO: make special op-code handler for 2 input, 1 output atomics
type OpcodeHandler func(inputs []CXValue, outputs []CXValue)

//TODO: Do atomic opcode handlers (not slices, not arrays)
//type AtomicOpcodeHandler func(input1 CXValue, input2 CXValue, output CXValue)

//TODO: RENAME THIS. The nameing is very bad
const (
	OPERATOR_COUNT         = constants.END_OF_OPERATORS - constants.START_OF_OPERATORS + 1
	OPERATOR_HANDLER_COUNT = constants.TYPE_COUNT * OPERATOR_COUNT
)

//Todo: Rename Natives
//Todo: What is an operator?
var (
	// Natives ...
	Natives   = map[int]*CXFunction{}
	Operators []*CXFunction
)

