package type_checker

import (
	"bytes"
	"fmt"
	"path/filepath"
	"regexp"
	"strconv"
	"unicode"

	"github.com/skycoin/cx/cmd/packageloader2/loader"
	"github.com/skycoin/cx/cx/ast"
	"github.com/skycoin/cx/cx/constants"
	"github.com/skycoin/cx/cx/types"
	"github.com/skycoin/cx/cxparser/actions"
)

func ParseParameterDeclaration(parameterString []byte, pkg *ast.CXPackage, fileName string, lineno int) (*ast.CXArgument, error) {
	var parameterDeclaration *ast.CXArgument
	reParameterDeclaration := regexp.MustCompile(`(\w+)((?:(?:\s*[\[\]\*\d]+|\s+)\w+(?:\.\w+)*))`)
	parameterDeclarationTokens := reParameterDeclaration.FindSubmatch(parameterString)

	// Check if the tokenized result is empty
	if parameterDeclarationTokens == nil || len(parameterDeclarationTokens[0]) == 0 {
		return nil, fmt.Errorf("%s: %d: parameter declaration error", filepath.Base(fileName), lineno)
	}

	// Set the declarator or the name of the param
	declarator := ast.MakeArgument("", fileName, lineno)
	declarator.SetType(types.UNDEFINED)
	declarator.Package = ast.CXPackageIndex(pkg.Index)
	declarator.Name = string(parameterDeclarationTokens[1])

	//Set the decalaration type
	parameterDeclaration, err := ParseDeclarationSpecifier(bytes.TrimSpace(parameterDeclarationTokens[2]), fileName, lineno, parameterDeclaration)
	if err != nil {
		return nil, err
	}

	// Merge both CXArgs
	parameterDeclaration.Name = declarator.Name
	parameterDeclaration.Package = declarator.Package

	return parameterDeclaration, nil
}

func ParseDeclarationSpecifier(declarationSpecifierByte []byte, fileName string, lineno int, declarationSpecifier *ast.CXArgument) (*ast.CXArgument, error) {
	// Base case if all parts are parsed
	if declarationSpecifierByte == nil || len(declarationSpecifierByte) == 0 {
		return declarationSpecifier, nil
	}

	// Checks last byte to determine what to parse
	lastByte := declarationSpecifierByte[len(declarationSpecifierByte)-1]

	if unicode.IsLetter(rune(lastByte)) || unicode.IsNumber(rune(lastByte)) || lastByte == '_' {

		reWords := regexp.MustCompile(`[\w\.]+`)
		words := reWords.FindAll(declarationSpecifierByte, -1)
		wordsIdx := reWords.FindAllIndex(declarationSpecifierByte, -1)
		newLastIdx := wordsIdx[len(wordsIdx)-1][0]

		dataType := words[len(words)-1]
		splitDataType := bytes.Split(dataType, []byte("."))

		newDeclarationSpecifierByte := declarationSpecifierByte[:newLastIdx]

		if len(splitDataType) == 1 {

			// Types like i32, str, aff, etc...
			if val, ok := TypesMap[string(splitDataType[0])]; ok {
				newDeclarationSpecifierArg := actions.DeclarationSpecifiersBasic(val)
				return ParseDeclarationSpecifier(newDeclarationSpecifierByte, fileName, lineno, newDeclarationSpecifierArg)
			}

			// Structs
			newDeclarationSpecifierArg := actions.DeclarationSpecifiersStruct(actions.AST, string(splitDataType[0]), "", false, fileName, lineno)
			return ParseDeclarationSpecifier(newDeclarationSpecifierByte, fileName, lineno, newDeclarationSpecifierArg)
		}

		// External types
		if val, ok := TypesMap[string(splitDataType[0])]; ok {
			newDeclarationSpecifierArg := actions.DeclarationSpecifiersStruct(actions.AST, string(splitDataType[1]), val.Name(), true, fileName, lineno)
			return ParseDeclarationSpecifier(newDeclarationSpecifierByte, fileName, lineno, newDeclarationSpecifierArg)
		}

		// External structs
		newDeclarationSpecifierArg := actions.DeclarationSpecifiersStruct(actions.AST, string(splitDataType[1]), string(splitDataType[0]), true, fileName, lineno)
		return ParseDeclarationSpecifier(newDeclarationSpecifierByte, fileName, lineno, newDeclarationSpecifierArg)
	}

	if lastByte == ']' {
		reBrackets := regexp.MustCompile(`\[\s*(\d*)\s*\]`)
		brackets := reBrackets.FindAllSubmatch(declarationSpecifierByte, -1)
		bracketsIdx := reBrackets.FindAllIndex(declarationSpecifierByte, -1)
		newLastIdx := bracketsIdx[len(bracketsIdx)-1][0]
		reNumber := regexp.MustCompile(`\d+`)
		number := reNumber.Find(brackets[len(brackets)-1][1])
		newDeclarationSpecifierByte := declarationSpecifierByte[:newLastIdx]

		// Arrays
		if number != nil {
			byteToInt, err := strconv.Atoi(string(number))
			if err != nil {
				return declarationSpecifier, err
			}

			declarationSpecifier.Lengths = append(declarationSpecifier.Lengths, types.Pointer(byteToInt))
			newDeclarationSpecifierArg := actions.DeclarationSpecifiers(declarationSpecifier, declarationSpecifier.Lengths, constants.DECL_ARRAY)
			return ParseDeclarationSpecifier(newDeclarationSpecifierByte, fileName, lineno, newDeclarationSpecifierArg)
		}

		// Slices
		newDeclarationSpecifierArg := actions.DeclarationSpecifiers(declarationSpecifier, []types.Pointer{0}, constants.DECL_SLICE)
		return ParseDeclarationSpecifier(newDeclarationSpecifierByte, fileName, lineno, newDeclarationSpecifierArg)
	}

	// Pointer
	if lastByte == '*' {
		newLastIdx := bytes.LastIndex(declarationSpecifierByte, []byte("*"))
		newDeclarationSpecifierByte := declarationSpecifierByte[:newLastIdx]
		newDeclarationSpecifierArg := actions.DeclarationSpecifiers(declarationSpecifier, []types.Pointer{0}, constants.DECL_POINTER)
		return ParseDeclarationSpecifier(newDeclarationSpecifierByte, fileName, lineno, newDeclarationSpecifierArg)

	}

	// If bytes don't match any of the cases
	return nil, fmt.Errorf("%v: %d: declaration specifier error", fileName, lineno)
}

// Finds the SourceBytes from the files array
func GetSourceBytes(files []*loader.File, fileName string) ([]byte, error) {
	for _, file := range files {
		if file.FileName == fileName {
			return file.Content, nil
		}
	}

	return nil, fmt.Errorf("%s not found", fileName)
}
