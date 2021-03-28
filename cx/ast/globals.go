package ast

var InREPL bool = false

//PROGRAM GLOBALS SHOULD NOT BE IN THIS FILE
// global reference to our program
var PROGRAM *CXProgram //Why do we have global?

// Initializing `CXProgram` structure where packages, structs, functions and
// global variables that belong to core packages are stored.
func init() {
	prgrm := CXProgram{Packages: make([]*CXPackage, 0)}
	PROGRAM = &prgrm
}

