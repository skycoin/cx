package parsingcompletor

import (
	"github.com/skycoin/cx/cx/opcodes"
	cxinit "github.com/skycoin/cx/cx/init"
)

var Initialized bool

func init() {
	InitCXCore()
}

func InitCXCore() {
	if !Initialized {
		opcodes.RegisterOpcodes()
		cxinit.RegisterPackages()
		Initialized = true
	}
}
