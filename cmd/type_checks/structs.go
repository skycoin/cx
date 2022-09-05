package type_checks

/*
// Parse Structs
// - takes in structs from cx/cmd/declaration_extractor
// - adds structs to AST
func ParseStructs(structs []declaration_extractor.StructDeclaration) {

	// Make porgram
	if actions.AST == nil {
		// fmt.Print(actions.AST, "\n")
		actions.AST = cxinit.MakeProgram()
		// fmt.Print("missing action API\n", &actions.AST, "\n")
	}

	for _, strct := range structs {

		pkg, err := actions.AST.GetPackage(strct.PackageID)

		if err != nil {

			newPkg := ast.MakePackage(strct.PackageID)
			pkgIdx := actions.AST.AddPackage(newPkg)
			newPkg, err = actions.AST.GetPackageFromArray(pkgIdx)

			if err != nil {
				// error handling
			}

			pkg = newPkg

		}

		structCX := ast.MakeStruct(strct.StructVariableName)
		structCX.Package = ast.CXPackageIndex(pkg.Index)

		pkg = pkg.AddStruct(actions.AST, structCX)

		file, err := os.Open(strct.FileID)

		if err != nil {
			// error handling
		}

		bracesOpen := regexp.MustCompile("{")
		bracesClose := regexp.MustCompile("}")

		srcBytes, err := io.ReadAll(file)

		strctDeclaration := srcBytes[strct.StartOffset : strct.StartOffset+strct.Length]

		if !bytes.Contains(strctDeclaration, []byte("struct")) {
			// not fields to be included
			fmt.Printf("syntax error: expecting type struct, line %v", strct.LineNumber)
			continue
		}

		var structFields []*ast.CXArgument


				structFieldArg := ast.MakeArgument(string(tokens[0]), strct.FileID, strct.LineNumber)
				structFieldArg = structFieldArg.SetPackage(pkg)

				structFieldSpecifier := actions.DeclarationSpecifiersBasic(primitiveTypesMap[string(tokens[1])])

				structFieldSpecifier.Name = structFieldArg.Name
				structFieldSpecifier.Package = structFieldArg.Package
				structFieldSpecifier.IsLocalDeclaration = true

				structFields = append(structFields, structFieldSpecifier)



		actions.DeclareStruct(actions.AST, strct.StructVariableName, structFields)

	}
}
*/
