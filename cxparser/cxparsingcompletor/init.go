package parsingcompletor

import (
	cxinit "github.com/skycoin/cx/cx/init"
	"github.com/skycoin/cx/cxparser/actions"
)

var Initialized bool

func init() {
	InitCXCore()
}

func InitCXCore() {
	if !Initialized {
		if actions.AST == nil {
			actions.AST = cxinit.MakeProgram()
		}

		Initialized = true
	}
}
