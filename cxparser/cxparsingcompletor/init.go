package parsingcompletor

import (
	cxinit "github.com/skycoin/cx/cx/init"
	"github.com/skycoin/cx/cx/opcodes"
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
