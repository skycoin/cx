package globals

import (
	"github.com/skycoin/cx/cx/constants"
)

var (
	OpCodeSystemCounter int = constants.END_OF_NAMED_OPCODES
	OpCodeMap map[string]int
	OpCodeReverseMap map[int]string
)
// MakeGenSym ...

func RegisterOpCodeWithIndex(name string, id int) {
	if id >= OpCodeSystemCounter {
		OpCodeSystemCounter = id+1
	}

	//if OpCodeMap
}

func RegisterOpCode(name string) {


	return
}