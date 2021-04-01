package astapi

import (
	cxast "github.com/skycoin/cx/cx/ast"
)

// GetPackagesNameList returns all names of all packages in cx program.
//
// The inputs are
// cxprogram - in the form of cxast.CXProgram.
//
// Example:
// We have this CX Program:
// 0.- Package: main
// 1.- Package: secondpackage
// 2.- Package: thirdpackage
//
// We use GetPackagesNameList(cxprogram).
// The Result will be the list of package names:
// []string{"main", "secondpackage", "thirdpackage"}
//
func GetPackagesNameList(cxprogram *cxast.CXProgram) (list []string) {
	for _, pkg := range cxprogram.Packages {
		list = append(list, pkg.Name)
	}
	return list
}

// AddEmptyPackage adds an empty package in cx program.
//
// The inputs are
// cxprogram - in the form of cxast.CXProgram.
// pkgName - the name of the package to be added.
//
// Example:
// We have an empty program.
// We use AddEmptyPackage(cxprogram, "main").
// The Result will be:
// 0.- Package: main
//
// Note the new main package added to the cx program.
func AddEmptyPackage(cxprogram *cxast.CXProgram, pkgName string) error {
	pkg := cxast.MakePackage(pkgName)
	cxprogram.AddPackage(pkg)
	return nil
}
