package type_checks

import (
	"os"
	"regexp"

	"github.com/skycoin/cx/cmd/declaration_extractor"
	"github.com/skycoin/cx/cx/ast"
	cxinit "github.com/skycoin/cx/cx/init"
	"github.com/skycoin/cx/cxparser/actions"
)

// Parse Structs
// - takes in structs from cx/cmd/declaration_extractor
// - adds structs to AST
func ParseStructs(structs []declaration_extractor.StructDeclaration) error {

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

		structCX := ast.MakeStruct(strct.StructName)
		structCX.Package = ast.CXPackageIndex(pkg.Index)

		pkg = pkg.AddStruct(actions.AST, structCX)

		src, err := os.ReadFile(strct.FileID)
		if err != nil {
			return err
		}

		var structFields []*ast.CXArgument

		for _, strctFieldDec := range strct.StructFields {

			structFieldLine := src[strct.StartOffset : strctFieldDec.StartOffset+strctFieldDec.Length]
			reStructField := regexp.MustCompile(`\w+\s+(\w+)`)
			structField := reStructField.FindSubmatch(structFieldLine)

			structFieldArg := ast.MakeArgument(strctFieldDec.StructFieldName, strct.FileID, strct.LineNumber)
			structFieldArg = structFieldArg.SetPackage(pkg)

			structFieldSpecifier := actions.DeclarationSpecifiersBasic(TypesMap[string(structField[1])])

			structFieldSpecifier.Name = structFieldArg.Name
			structFieldSpecifier.Package = structFieldArg.Package

			structFields = append(structFields, structFieldSpecifier)
		}

		actions.DeclareStruct(actions.AST, strct.StructName, structFields)

	}
	return nil
}
