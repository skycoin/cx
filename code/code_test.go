package code

import "testing"

func TestMake(t *testing.T) {

	tests := []struct {
		op       OpCode
		operands []int
		expected []byte
	}{
		{OpConstant, []int{65534}, []byte{byte(OpConstant), 255, 254}},
	}

	for _, test := range tests {
		instruction := Make(test.op, test.operands...)

		if len(instruction) != len(test.expected) {
			t.Errorf("instuction has wrong length. want=%d, got=%d", len(test.expected), len(instruction))
		}

		for i, bt := range test.expected {

			if instruction[i] != test.expected[i] {
				t.Errorf("wrong byte at pos %d. want=%d, got=%d", i, bt, instruction[i])
			}
		}
	}
}
