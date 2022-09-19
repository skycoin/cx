package type_checks

import (
	"bytes"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/skycoin/cx/cx/ast"
	"github.com/skycoin/cx/cx/constants"
	"github.com/skycoin/cx/cx/types"
	"github.com/skycoin/cx/cxparser/actions"
)

func ParseParameterDeclaration(parameterString []byte, pkg *ast.CXPackage, fileName string, lineno int) (*ast.CXArgument, error) {
	var parameterDeclaration *ast.CXArgument
	reParameterDeclaration := regexp.MustCompile(`(\w+)\s+(.+)`)
	parameterDeclarationTokens := reParameterDeclaration.FindSubmatch(parameterString)

	declarator := ast.MakeArgument("", fileName, lineno)
	declarator.SetType(types.UNDEFINED)
	declarator.Package = ast.CXPackageIndex(pkg.Index)
	declarator.Name = string(parameterDeclarationTokens[1])

	parameterDeclaration, err := ParseDeclarationSpecifier(parameterDeclarationTokens[2], fileName, lineno, parameterDeclaration)

	if err != nil {
		return nil, err
	}

	parameterDeclaration.Name = declarator.Name
	parameterDeclaration.Package = declarator.Package

	return parameterDeclaration, nil
}

func ParseDeclarationSpecifier(declarationSpecifierByte []byte, fileName string, lineno int, declarationSpecifier *ast.CXArgument) (*ast.CXArgument, error) {
	reDeclarationSpecifier := regexp.MustCompile(`(\*){0,1}\s*((?:\[(\d*)\])){0,1}\s*([\w.]*)`)
	declarationSpecifierTokens := reDeclarationSpecifier.FindSubmatch(declarationSpecifierByte)
	declarationSpecifierTokensIdx := reDeclarationSpecifier.FindIndex(declarationSpecifierByte)

	// Base case if all parts are parsed
	if declarationSpecifierByte == nil || declarationSpecifierTokensIdx[1] == 0 {
		return declarationSpecifier, nil
	}

	// Types like i32, str, aff, etc...
	if val, ok := TypesMap[string(declarationSpecifierTokens[4])]; ok {
		newDeclarationSpecifierByte := declarationSpecifierByte[:declarationSpecifierTokensIdx[1]-len(declarationSpecifierTokens[4])]
		newDeclarationSpecifierArg := actions.DeclarationSpecifiersBasic(val)
		return ParseDeclarationSpecifier(newDeclarationSpecifierByte, fileName, lineno, newDeclarationSpecifierArg)
	}

	// External structs and types like myPackage.Animal
	if bytes.Contains(declarationSpecifierTokens[4], []byte(".")) {
		tokens := strings.Split(string(declarationSpecifierTokens[4]), ".")

		// External types
		if val, ok := TypesMap[tokens[0]]; ok {
			newDeclarationSpecifierByte := declarationSpecifierByte[:declarationSpecifierTokensIdx[1]-len(declarationSpecifierTokens[4])]
			newDeclarationSpecifierArg := actions.DeclarationSpecifiersStruct(actions.AST, tokens[1], val.Name(), true, fileName, lineno)
			return ParseDeclarationSpecifier(newDeclarationSpecifierByte, fileName, lineno, newDeclarationSpecifierArg)
		}

		// External structs
		newDeclarationSpecifierByte := declarationSpecifierByte[:declarationSpecifierTokensIdx[1]-len(declarationSpecifierTokens[4])]
		newDeclarationSpecifierArg := actions.DeclarationSpecifiersStruct(actions.AST, tokens[1], tokens[0], true, fileName, lineno)
		return ParseDeclarationSpecifier(newDeclarationSpecifierByte, fileName, lineno, newDeclarationSpecifierArg)
	}

	// Structs
	if len(declarationSpecifierTokens[4]) != 0 {
		fmt.Print("worksfsdfsdfdsfsdfsdf")

		newDeclarationSpecifierByte := declarationSpecifierByte[:declarationSpecifierTokensIdx[1]-len(declarationSpecifierTokens[4])]
		newDeclarationSpecifierArg := actions.DeclarationSpecifiersStruct(actions.AST, string(declarationSpecifierTokens[4]), "", false, fileName, lineno)
		return ParseDeclarationSpecifier(newDeclarationSpecifierByte, fileName, lineno, newDeclarationSpecifierArg)
	}

	// Arrays
	if declarationSpecifierTokens[2] != nil && len(declarationSpecifierTokens[3]) != 0 {
		byteToInt, err := strconv.Atoi(string(declarationSpecifierTokens[3]))
		if err != nil {
			return declarationSpecifier, err
		}

		newDeclarationSpecifierByte := declarationSpecifierByte[:declarationSpecifierTokensIdx[1]-len(declarationSpecifierTokens[2])]
		newDeclarationSpecifierArg := actions.DeclarationSpecifiers(declarationSpecifier, types.Cast_sint_to_sptr([]int{byteToInt}), constants.DECL_ARRAY)
		return ParseDeclarationSpecifier(newDeclarationSpecifierByte, fileName, lineno, newDeclarationSpecifierArg)
	}

	// Slices
	if declarationSpecifierTokens[2] != nil && len(declarationSpecifierTokens[3]) == 0 {
		newDeclarationSpecifierByte := declarationSpecifierByte[:declarationSpecifierTokensIdx[1]-len(declarationSpecifierTokens[2])]
		newDeclarationSpecifierArg := actions.DeclarationSpecifiers(declarationSpecifier, []types.Pointer{0}, constants.DECL_SLICE)
		return ParseDeclarationSpecifier(newDeclarationSpecifierByte, fileName, lineno, newDeclarationSpecifierArg)
	}

	// Pointer
	if declarationSpecifierTokens[1] != nil {
		return actions.DeclarationSpecifiers(declarationSpecifier, []types.Pointer{0}, constants.DECL_POINTER), nil
	}

	return nil, fmt.Errorf("%v: %d: declaration specifier error", fileName, lineno)
}
