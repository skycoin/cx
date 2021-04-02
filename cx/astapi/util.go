package astapi

import (
	"errors"

	cxast "github.com/skycoin/cx/cx/ast"
)

func FindPackage(cxprogram *cxast.CXProgram, packageName string) (*cxast.CXPackage, error) {
	for _, pkg := range cxprogram.Packages {
		if pkg.Name == packageName {
			return pkg, nil
		}
	}
	return nil, errors.New("package not found")
}

func FindFunction(cxprogram *cxast.CXProgram, functionName string) (*cxast.CXFunction, error) {
	for _, pkg := range cxprogram.Packages {
		for _, fn := range pkg.Functions {
			if fn.Name == functionName {
				return fn, nil
			}
		}
	}
	return nil, errors.New("function not found")
}
