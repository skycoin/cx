package globals

import (
	"github.com/skycoin/cx/cx/ast"
	"os"
)


//PROGRAM GLOBALS SHOULD NOT BE IN THIS FILE
// global reference to our program
var PROGRAM *ast.CXProgram //Why do we have global?

// Initializing `CXProgram` structure where packages, structs, functions and
// global variables that belong to core packages are stored.
func init() {
	prgrm := ast.CXProgram{Packages: make([]*ast.CXPackage, 0)}
	PROGRAM = &prgrm
}

// Var
var (
	HeapOffset    int
	GenSymCounter int
)

//Path is only used by os module and only to get working directory
//Path and working directory should not be hard coded into program struct (etc, when serialized)
//Working directory is property of executable and can be retrieved with golang library
var CxProgramPath string = ""

var CXPATH = os.Getenv("CXPATH") + "/"
var BINPATH = CXPATH + "bin/" // TODO @evanlinjin: Not used.
var PKGPATH = CXPATH + "pkg/" // TODO @evanlinjin: Not used.
var SRCPATH = CXPATH + "src/"
