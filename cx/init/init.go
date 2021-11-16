package init

import (
	"github.com/skycoin/cx/cx/ast"
	"github.com/skycoin/cx/cx/packages/cipher"
	"github.com/skycoin/cx/cx/packages/cxfx"
	"github.com/skycoin/cx/cx/packages/cxos"

	// "github.com/skycoin/cx/cx/packages/http"
	"github.com/skycoin/cx/cx/packages/regexp"
)

func RegisterPackages() {
	cipher.RegisterPackage(ast.PROGRAM)
	cxfx.RegisterPackage()
	cxos.RegisterPackage()
	// http.RegisterPackage()
	regexp.RegisterPackage(ast.PROGRAM)
}
