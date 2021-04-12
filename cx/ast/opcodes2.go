package ast

import (
	"github.com/skycoin/cx/cx/constants"
)

// GetOpCodeCount returns an op code that is available for usage on the CX standard library.
/*
func GetOpCodeCount() int {
	return len(OpcodeHandlers)
}
*/

//TODO: WHAT IS AN "OPERATOR"
func IsOperator(opCode int) bool {
	return opCode > constants.START_OF_OPERATORS && opCode < constants.END_OF_OPERATORS
}

func IsArithmeticOperator(opCode int) bool {
	return opCode > constants.START_OF_ARITHMETIC_OPERATORS && opCode < constants.END_OF_ARITHMETIC_OPERATORS
}

func IsComparisonOperator(opCode int) bool {
	return opCode > constants.START_OF_COMPARISON_OPERATORS && opCode < constants.END_OF_COMPARISON_OPERATORS
}

func GetTypedOperatorOffset(typeCode int, opCode int) int {
	return typeCode*OPERATOR_COUNT + opCode - constants.START_OF_OPERATORS - 1
}

func GetTypedOperator(typeCode int, opCode int) *CXFunction {
	return Operators[GetTypedOperatorOffset(typeCode, opCode)]
}

