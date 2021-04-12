package parsingcompletor

import (
	"github.com/skycoin/cx/cx/opcodes"
)

var Initialized bool

func init() {
	InitCXCore()
}

func InitCXCore() {
	if !Initialized {
		opcodes.LoadOpCodeTables()
		Initialized = true
	}
}
