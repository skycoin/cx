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

const (
	OPERATOR_COUNT         = constants.END_OF_OPERATORS - constants.START_OF_OPERATORS + 1
	OPERATOR_HANDLER_COUNT = constants.TYPE_COUNT * OPERATOR_COUNT
)

// OpcodeHandler ...
type OpcodeHandler func(inputs []CXValue, outputs []CXValue)

//Todo: Rename Natives
//Todo: What is an operator?
var (
	// Natives ...
	Natives   = map[int]*CXFunction{}
	Operators []*CXFunction
)

