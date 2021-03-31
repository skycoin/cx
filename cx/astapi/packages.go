package astapi

import (
	cxast "github.com/skycoin/cx/cx/ast"
)

func GetPackagesNameList(cxprogram *cxast.CXProgram) (list []string) {
	for _, pkg := range cxprogram.Packages {
		list = append(list, pkg.Name)
	}
	return list
}

func AddEmptyPackage(cxprogram *cxast.CXProgram, pkgName string) error {
	pkg := cxast.MakePackage(pkgName)
	cxprogram.AddPackage(pkg)
	return nil
}
