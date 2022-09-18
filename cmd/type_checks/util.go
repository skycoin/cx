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

	typeString, file, line, parameterDeclaration, err := ParseTypeSpecifier(parameterDeclarationTokens[2], fileName, lineno, parameterDeclaration)

	if typeString != nil && file != "" && line != 0 && err != nil {
		return nil, err
	}

	parameterDeclaration.Name = declarator.Name
	parameterDeclaration.Package = declarator.Package

	return parameterDeclaration, nil
}

func ParseTypeSpecifier(typeString []byte, fileName string, lineno int, declarationSpecifier *ast.CXArgument) ([]byte, string, int, *ast.CXArgument, error) {
	reTypeSpecifier := regexp.MustCompile(`(\*){0,1}\s*((?:\[(\d*)\])){0,1}\s*(\w*){0,1}`)
	typeSpecifierTokens := reTypeSpecifier.FindSubmatch(typeString)
	typeSpecifierTokensIdx := reTypeSpecifier.FindSubmatchIndex(typeString)

	if typeString == nil || typeSpecifierTokensIdx[1] == 0 {
		return nil, "", 0, declarationSpecifier, nil
	}

	if val, ok := TypesMap[string(typeSpecifierTokens[4])]; ok {
		fmt.Print("works", typeSpecifierTokens[4])
		return ParseTypeSpecifier(typeString[typeSpecifierTokensIdx[0]:typeSpecifierTokensIdx[1]-(typeSpecifierTokensIdx[9]-typeSpecifierTokensIdx[8])], fileName, lineno, actions.DeclarationSpecifiersBasic(val))
	}

	if bytes.Contains(typeSpecifierTokens[4], []byte(".")) {
		tokens := strings.Split(string(typeSpecifierTokens[4]), ".")

		if val, ok := TypesMap[tokens[0]]; ok {
			return ParseTypeSpecifier(typeString[typeSpecifierTokensIdx[0]:typeSpecifierTokensIdx[1]-(typeSpecifierTokensIdx[9]-typeSpecifierTokensIdx[8])], fileName, lineno, actions.DeclarationSpecifiersStruct(actions.AST, tokens[1], val.Name(), true, fileName, lineno))
		}
		return ParseTypeSpecifier(typeString[typeSpecifierTokensIdx[0]:typeSpecifierTokensIdx[1]-(typeSpecifierTokensIdx[9]-typeSpecifierTokensIdx[8])], fileName, lineno, actions.DeclarationSpecifiersStruct(actions.AST, tokens[1], tokens[0], true, fileName, lineno))
	} else {

	}

	// if typeSpecifierTokens[4] != nil {
	// 	fmt.Print(string(typeSpecifierTokens[4]), "works")
	// 	return ParseTypeSpecifier(typeString[typeSpecifierTokensIdx[0]:typeSpecifierTokensIdx[1]-(typeSpecifierTokensIdx[9]-typeSpecifierTokensIdx[8])], fileName, lineno, actions.DeclarationSpecifiersStruct(actions.AST, string(typeSpecifierTokens[4]), "", false, fileName, lineno))
	// }

	if typeSpecifierTokens[2] != nil && typeSpecifierTokens[3] != nil {
		byteToInt, err := strconv.Atoi(string(typeSpecifierTokens[3]))
		if err != nil {
			return typeString, fileName, lineno, declarationSpecifier, err
		}

		return ParseTypeSpecifier(typeString[typeSpecifierTokensIdx[0]:typeSpecifierTokensIdx[1]-(typeSpecifierTokensIdx[7]-typeSpecifierTokensIdx[6])], fileName, lineno, actions.DeclarationSpecifiers(declarationSpecifier, types.Cast_sint_to_sptr([]int{byteToInt}), constants.DECL_ARRAY))
	}

	if typeSpecifierTokens[2] != nil && typeSpecifierTokens == nil {
		return ParseTypeSpecifier(typeString[typeSpecifierTokensIdx[0]:typeSpecifierTokensIdx[1]-(typeSpecifierTokensIdx[7]-typeSpecifierTokensIdx[6])], fileName, lineno, actions.DeclarationSpecifiers(declarationSpecifier, []types.Pointer{0}, constants.DECL_SLICE))
	}

	if typeSpecifierTokens[1] != nil {
		fmt.Print(string(typeSpecifierTokens[0]))
		return nil, "", 0, actions.DeclarationSpecifiers(declarationSpecifier, []types.Pointer{0}, constants.DECL_POINTER), nil
	}

	return ParseTypeSpecifier(typeString[typeSpecifierTokensIdx[0]:typeSpecifierTokensIdx[1]-(typeSpecifierTokensIdx[9]-typeSpecifierTokensIdx[8])], fileName, lineno, actions.DeclarationSpecifiersStruct(actions.AST, string(typeSpecifierTokens[4]), "", false, fileName, lineno))

}
