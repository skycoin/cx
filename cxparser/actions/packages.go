package actions

import (
	"fmt"

	"github.com/skycoin/cx/cx/ast"
	cxpackages "github.com/skycoin/cx/cx/packages"
)

// DeclarePackage() switches the current package in the program.
//
// Input arguments description:
// 	prgrm - a CXProgram that contains all the data and arrays of the program.
// 	pkgName - name of the package to declare.
func DeclarePackage(prgrm *ast.CXProgram, pkgName string) {
	// Add a new package to the program if it's not previously defined.
	if _, err := prgrm.GetPackage(pkgName); err != nil {
		pkg := ast.MakePackage(pkgName)
		prgrm.AddPackage(pkg)
	}

	_, err := prgrm.SelectPackage(pkgName)
	if err != nil {
		panic(err)
	}
}

// DeclareImport()
//
// Input arguments description:
// 	prgrm - a CXProgram that contains all the data and arrays of the program.
// 	impName - name of the import to declare.
// 	currentFile - name of the current cx source code file.
// 	lineNo - the current line number from the source code.
func DeclareImport(prgrm *ast.CXProgram, impName string, currentFile string, lineNo int) {
	// Make sure we are inside a package
	pkg, err := prgrm.GetCurrentPackage()
	if err != nil {
		// FIXME: Should give a relevant error message
		panic(err)
	}

	// Checking if it's a package in the CX workspace by trying to find a
	// slash (/) in the impName.
	// We start backwards and we stop if we find a slash.
	hasSlash := false
	c := len(impName) - 1
	for ; c >= 0; c-- {
		if impName[c] == '/' {
			hasSlash = true
			break
		}
	}
	ident := ""
	// If the `impName` has a slash, then we need to strip
	// everything behind the slash and the slash itself.
	if hasSlash {
		ident = impName[c+1:]
	} else {
		ident = impName
	}

	// If the package is already imported, then there is nothing more to be done.
	if _, err := pkg.GetImport(prgrm, ident); err == nil {
		return
	}

	// If the package is already defined in the program, just add it to
	// the importing package.
	if imp, err := prgrm.GetPackage(ident); err == nil {
		pkg.AddImport(prgrm, imp)
		return
	}

	// All packages are read during the first pass of the compilation.  So
	// if we get here during the 2nd pass, it's either a core package or
	// something is panic-level wrong.
	if cxpackages.IsDefaultPackage(ident) {
		imp := ast.MakePackage(ident)
		impIdx := prgrm.AddPackage(imp)
		newImp, err := prgrm.GetPackageFromArray(impIdx)
		if err != nil {
			panic(err)
		}
		pkg.AddImport(prgrm, newImp)

		prgrm.CurrentPackage = ast.CXPackageIndex(pkg.Index)

		if ident == "aff" {
			AffordanceStructs(prgrm, newImp, currentFile, lineNo)
		}
	} else {
		// This should never happen.
		panic(fmt.Sprintf("%v: unkown error when trying to read package '%s'", ast.CompilationError(currentFile, lineNo), ident))
	}
}
