package ast

import (
	"github.com/skycoin/cx/cx/constants"
)

var (
	// OpNames ...
	OpNames = map[int]string{}

	// OpCodes ...
	OpCodes = map[string]int{}

	// Versions ...
	//OpVersions = map[int]int{}
)

var (
	OpcodeHandlers    []OpcodeHandler
	OpcodeHandlers_V2 []OpcodeHandler_V2
)

const (
	OPERATOR_COUNT         = constants.END_OF_OPERATORS - constants.START_OF_OPERATORS + 1
	OPERATOR_HANDLER_COUNT = constants.TYPE_COUNT * OPERATOR_COUNT
)

// OpcodeHandler ...
type OpcodeHandler func(expr *CXExpression, fp int)
type OpcodeHandler_V2 func(inputs []CXValue, outputs []CXValue)

//Todo: Rename Natives
//Todo: What is an operator?
var (
	// Natives ...
	Natives   = map[int]*CXFunction{}
	Operators []*CXFunction
)

