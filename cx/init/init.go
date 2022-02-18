package init

import (
	"github.com/skycoin/cx/cx/ast"
	"github.com/skycoin/cx/cx/packages/cipher"
	"github.com/skycoin/cx/cx/packages/cxfx"
	"github.com/skycoin/cx/cx/packages/cxos"

	// "github.com/skycoin/cx/cx/packages/http"
	"github.com/skycoin/cx/cx/packages/regexp"
)

func RegisterPackages(prgrm *ast.CXProgram) {
	cipher.RegisterPackage(prgrm)
	cxfx.RegisterPackage(prgrm)
	cxos.RegisterPackage(prgrm)
	// http.RegisterPackage()
	regexp.RegisterPackage(prgrm)
}
