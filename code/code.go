package code

import (
	"encoding/binary"
	"fmt"
)

type Instructions []byte

type OpCode byte

const (
	OpConstant OpCode = iota
)

type Definition struct {
	Name          string
	OperandWidths []int
}

var definitions = map[OpCode]*Definition{
	OpConstant: {"OpConstant", []int{2}},
}

func Lookup(op byte) (*Definition, error) {

	def, ok := definitions[OpCode(op)]

	if !ok {
		return nil, fmt.Errorf("opcde %d undefined", op)
	}

	return def, nil
}

func Make(op OpCode, operands ...int) []byte {

	def, ok := definitions[op]

	if !ok {
		return []byte{}
	}

	instructionLen := 1

	for _, operandWidth := range def.OperandWidths {

		instructionLen = instructionLen + operandWidth
	}

	instruction := make([]byte, instructionLen)

	instruction[0] = byte(op)

	offset := 1

	for i, operand := range operands {

		width := def.OperandWidths[i]

		switch width {

		case 2:
			binary.BigEndian.PutUint16(instruction[offset:], uint16(uint16(operand)))

		}

		offset = offset + width
	}
	return instruction
}
