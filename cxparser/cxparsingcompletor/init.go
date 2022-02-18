package parsingcompletor

import (
	"github.com/skycoin/cx/cx/ast"
	cxinit "github.com/skycoin/cx/cx/init"
	"github.com/skycoin/cx/cx/opcodes"
	"github.com/skycoin/cx/cxparser/actions"
)

var Initialized bool

func init() {
	InitCXCore()
}

func InitCXCore() {
	if !Initialized {
		if actions.AST == nil {
			actions.AST = ast.MakeProgram()
		}

		opcodes.RegisterOpcodes(actions.AST)
		cxinit.RegisterPackages(actions.AST)

		Initialized = true
	}
}
