package type_checks

/*
// Parse Function Headers
// - takes in funcs from cx/cmd/declaration_extractor
// - adds func headers to AST
func ParseFuncHeaders(funcs []declaration_extractor.FuncDeclaration) {

	// Make program
	if actions.AST == nil {
		actions.AST = cxinit.MakeProgram()
	}

	for _, fun := range funcs {

		pkg, err := actions.AST.GetPackage(fun.PackageID)

		if err != nil {

			newPkg := ast.MakePackage(fun.PackageID)
			pkgIdx := actions.AST.AddPackage(newPkg)
			newPkg, err = actions.AST.GetPackageFromArray(pkgIdx)

			if err != nil {
				// error handling
			}

			pkg = newPkg

		}


}
*/
