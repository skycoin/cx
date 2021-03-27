package ast

import "github.com/skycoin/cx/cx"

// ToString returns the abstract syntax tree of a CX program in a
// string format.
func (cxprogram *CXProgram) ToString() string {
	var ast string
	ast += "Program\n" //why is top line "Program" ???

	var currentFunction *CXFunction
	var currentPackage *CXPackage

	currentPackage, err := cxprogram.GetCurrentPackage()

	if err != nil {
		panic("CXProgram.ToString(): error, currentPackage is nil")
	}

	currentFunction, _ = cxprogram.GetCurrentFunction()
	currentPackage.CurrentFunction = currentFunction

	cxcore.BuildStrPackages(cxprogram, &ast) //what does this do?

	return ast
}