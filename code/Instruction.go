package code

/*
	ReadOperands returns operands for Instructions.
*/
func ReadOperands(def *Definition, ins Instructions) ([]int, int) {

	operands := make([]int, len(def.OperandWidths))

	offset := 0

	for i, width := range def.OperandWidths {

		switch width {

		case 2:

			operands[i] = int(ReadUint16(ins[offset:]))

		case 1:
			operands[i] = int(ReadUint8(ins[offset:]))
		}

		offset += width
	}

	return operands, offset
}
