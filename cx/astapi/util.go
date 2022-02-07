package astapi

import (
	"errors"

	cxast "github.com/skycoin/cx/cx/ast"
)

func FindPackage(cxprogram *cxast.CXProgram, packageName string) (*cxast.CXPackage, error) {
	for _, pkgIdx := range cxprogram.Packages {
		pkg, err := cxprogram.GetPackageFromArray(pkgIdx)
		if err != nil {
			panic(err)
		}

		if pkg.Name == packageName {
			return pkg, nil
		}
	}
	return nil, errors.New("package not found")
}

func FindFunction(cxprogram *cxast.CXProgram, functionName string) (*cxast.CXFunction, error) {
	for _, pkgIdx := range cxprogram.Packages {
		pkg, err := cxprogram.GetPackageFromArray(pkgIdx)
		if err != nil {
			panic(err)
		}

		for _, fnIdx := range pkg.Functions {
			fn, err := cxprogram.GetFunctionFromArray(fnIdx)
			if err != nil {
				return nil, err
			}

			if fn.Name == functionName {
				return fn, nil
			}
		}
	}
	return nil, errors.New("function not found")
}
