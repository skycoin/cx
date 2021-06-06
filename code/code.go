package code

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

/*
	Instructions represents cx instruction as opcode.
*/
type Instructions []byte

/*
	String represents Instructions into string.
*/
func (ins Instructions) String() string {

	var out bytes.Buffer

	i := 0

	for i < len(ins) {

		def, err := Lookup(ins[i])

		if err != nil {
			fmt.Fprintf(&out, "ERROR: %s\n", err)
			continue
		}

		operands, read := ReadOperands(def, ins[i+1:])

		fmt.Fprintf(&out, "%04d %s\n", i, ins.fmtInstruction(def, operands))

		i += 1 + read
	}

	return out.String()
}

/*
	fmtInstruction print instruction into string format.
*/
func (ins Instructions) fmtInstruction(def *Definition, operands []int) string {

	operandCount := len(def.OperandWidths)

	if len(operands) != operandCount {
		return fmt.Sprintf("ERROR: operand len %d does not match defined %d\n",
			len(operands), operandCount)
	}

	switch operandCount {
	case 0:
		return def.Name
	case 1:
		return fmt.Sprintf("%s %d", def.Name, operands[0])
	case 2:
		return fmt.Sprintf("%s %d %d", def.Name, operands[0], operands[1])
	}

	return fmt.Sprintf("ERROR: unhandled operandCount for %s\n", def.Name)
}

/*
 Lookup return Definition base on op byte.
*/
func Lookup(op byte) (*Definition, error) {

	def, ok := definitions[Opcode(op)]

	if !ok {
		return nil, fmt.Errorf("opcode %d undefined", op)
	}

	return def, nil
}

/*
	Make creates instruction for cx.
*/
func Make(op Opcode, operands ...int) []byte {

	def, ok := definitions[op]

	if !ok {
		return []byte{}
	}

	instructionLen := 1

	for _, w := range def.OperandWidths {
		instructionLen += w
	}

	instruction := make([]byte, instructionLen)
	instruction[0] = byte(op)

	offset := 1

	for i, o := range operands {

		width := def.OperandWidths[i]

		switch width {

		case 2:
			binary.BigEndian.PutUint16(instruction[offset:], uint16(o))

		case 1:
			instruction[offset] = byte(o)
		}
		offset += width
	}

	return instruction
}

/*
	ReadUint8 returns uint8
*/
func ReadUint8(ins Instructions) uint8 {

	return uint8(ins[0])
}

/*
	ReadUint16 returns uint16
*/
func ReadUint16(ins Instructions) uint16 {
	return binary.BigEndian.Uint16(ins)
}
