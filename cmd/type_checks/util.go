package type_checks

import (
	"bytes"
	"regexp"
	"strconv"
	"strings"

	"github.com/skycoin/cx/cx/ast"
	"github.com/skycoin/cx/cx/constants"
	"github.com/skycoin/cx/cx/types"
	"github.com/skycoin/cx/cxparser/actions"
)

func ParseParameterDeclaration(parameterString string, fileName string, lineno int) (ast.CXArgument, error) {
	var parameterDeclaration ast.CXArgument
	// reParameterDeclaration := regexp.MustCompile(`(\w+)\s+(.+)`)

	return parameterDeclaration, nil
}

func ParseTypeSpecifier(typeString []byte, fileName string, lineno int, declarationSpecifier *ast.CXArgument) ([]byte, string, int, *ast.CXArgument, error) {
	reTypeSpecifier := regexp.MustCompile(`(\*){0,1}\s*((?:\[(\d*)\])){0,1}\s*(\w*){0,1}`)
	typeSpecifierTokens := reTypeSpecifier.FindSubmatch(typeString)
	typeSpecifierTokensIdx := reTypeSpecifier.FindSubmatchIndex(typeString)

	if typeString == nil {
		return []byte(""), "", 0, declarationSpecifier, nil
	}

	if val, ok := TypesMap[string(typeSpecifierTokens[4])]; ok {
		return ParseTypeSpecifier(typeString[typeSpecifierTokensIdx[0]:typeSpecifierTokensIdx[7]], fileName, lineno, actions.DeclarationSpecifiersBasic(val))
	}

	if bytes.Contains(typeSpecifierTokens[4], []byte(".")) {
		tokens := strings.Split(string(typeSpecifierTokens[4]), ".")

		if val, ok := TypesMap[tokens[0]]; ok {
			return ParseTypeSpecifier(typeString[typeSpecifierTokensIdx[0]:typeSpecifierTokensIdx[7]], fileName, lineno, actions.DeclarationSpecifiersStruct(actions.AST, tokens[1], val.Name(), true, fileName, lineno))
		}
		return ParseTypeSpecifier(typeString[typeSpecifierTokensIdx[0]:typeSpecifierTokensIdx[7]], fileName, lineno, actions.DeclarationSpecifiersStruct(actions.AST, tokens[1], tokens[0], true, fileName, lineno))
	}

	if typeSpecifierTokens[2] != nil && typeSpecifierTokens != nil {
		byteToInt, err := strconv.Atoi(string(typeSpecifierTokens[3]))
		if err != nil {
			return typeString, fileName, lineno, declarationSpecifier, err
		}
		return ParseTypeSpecifier(typeString[typeSpecifierTokensIdx[0]:typeSpecifierTokensIdx[5]], fileName, lineno, actions.DeclarationSpecifiers(declarationSpecifier, types.Cast_sint_to_sptr([]int{byteToInt}), constants.DECL_ARRAY))
	}

	if typeSpecifierTokens[2] != nil && typeSpecifierTokens == nil {
		return ParseTypeSpecifier(typeString[typeSpecifierTokensIdx[0]:typeSpecifierTokensIdx[5]], fileName, lineno, actions.DeclarationSpecifiers(declarationSpecifier, []types.Pointer{0}, constants.DECL_SLICE))
	}

	return ParseTypeSpecifier(typeString[typeSpecifierTokensIdx[0]:typeSpecifierTokensIdx[7]], fileName, lineno, actions.DeclarationSpecifiersStruct(actions.AST, string(typeSpecifierTokens[4]), "", false, fileName, lineno))
}
