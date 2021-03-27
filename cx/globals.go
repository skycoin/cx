package cxcore

import (
	"os"
)

// Initializing `CXProgram` structure where packages, structs, functions and
// global variables that belong to core packages are stored.
func init() {
	prgrm := CXProgram{Packages: make([]*CXPackage, 0)}
	PROGRAM = &prgrm
}

var InREPL bool = false
var FoundCompileErrors bool

const DBG_GOLANG_STACK_TRACE = true

//PROGRAM GLOBALS SHOULD NOT BE IN THIS FILE
// global reference to our program
var PROGRAM *CXProgram //Why do we have global?

var CXPATH = os.Getenv("CXPATH") + "/"
var BINPATH = CXPATH + "bin/" // TODO @evanlinjin: Not used.
var PKGPATH = CXPATH + "pkg/" // TODO @evanlinjin: Not used.
var SRCPATH = CXPATH + "src/"
