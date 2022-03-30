package init

import (
	"github.com/skycoin/cx/cx/ast"
	"github.com/skycoin/cx/cx/constants"
	"github.com/skycoin/cx/cx/packages/cipher"
	"github.com/skycoin/cx/cx/packages/cxfx"
	"github.com/skycoin/cx/cx/packages/cxos"
	"github.com/skycoin/cx/cx/packages/tcp"

	// "github.com/skycoin/cx/cx/packages/http"

	"github.com/skycoin/cx/cx/packages/regexp"
)

func RegisterPackages(prgrm *ast.CXProgram) {
	cipher.RegisterPackage(prgrm)
	cxfx.RegisterPackage(prgrm)
	cxos.RegisterPackage(prgrm)
	// http.RegisterPackage()
	regexp.RegisterPackage(prgrm)
	tcp.RegisterPackage(prgrm)
}

// MakeProgram ...
func MakeProgram() *ast.CXProgram {
	minHeapSize := ast.MinHeapSize()
	newPrgrm := &ast.CXProgram{
		Packages:       make(map[string]ast.CXPackageIndex, 0),
		CXFunctions:    make([]ast.CXFunction, 0),
		CurrentPackage: -1,
		CallStack:      make([]ast.CXCall, constants.CALLSTACK_SIZE),
		Memory:         make([]byte, constants.STACK_SIZE+minHeapSize),
		Stack: ast.StackSegmentStruct{
			Size: constants.STACK_SIZE,
		},
		Data: ast.DataSegmentStruct{
			StartsAt: constants.STACK_SIZE,
		},
		Heap: ast.HeapSegmentStruct{
			Size:    minHeapSize,
			Pointer: constants.NULL_HEAP_ADDRESS_OFFSET, // We can start adding objects to the heap after the NULL (nil) bytes.
		},
		CXArgs: make([]ast.CXArgument, 0),
	}

	RegisterPackages(newPrgrm)

	return newPrgrm
}
