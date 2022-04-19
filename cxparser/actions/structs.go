package actions

import "github.com/skycoin/cx/cx/ast"

// DeclareStruct takes a name of a struct and a slice of fields representing
// the members and adds the struct to the package.
//
// Input arguments description:
// 	prgrm - a CXProgram that contains all the data and arrays of the program.
// 	structName - name of the struct to declare.
//  structFields - fields of the struct to be added.
func DeclareStruct(prgrm *ast.CXProgram, structName string, structFields []*ast.CXArgument) {
	// Make sure we are inside a package.
	pkg, err := prgrm.GetCurrentPackage()
	if err != nil {
		// FIXME: Should give a relevant error message
		panic(err)
	}

	// Make sure a struct with the same name is not yet defined.
	strct, err := prgrm.GetStruct(structName, pkg.Name)
	if err != nil {
		// FIXME: Should give a relevant error message
		panic(err)
	}

	strct.Fields = nil
	for _, field := range structFields {
		if _, err := strct.GetField(prgrm, field.Name); err == nil {
			println(ast.CompilationError(field.ArgDetails.FileName, field.ArgDetails.FileLine), "Multiply defined struct field:", field.Name)
		} else {
			strct.AddField(prgrm, ast.TYPE_COMPLEX, field, nil)
		}
	}
}
