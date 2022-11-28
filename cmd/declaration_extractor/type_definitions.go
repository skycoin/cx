package declaration_extractor

import (
	"bufio"
	"bytes"
	"fmt"
	"path/filepath"
	"regexp"
)

type TypeDefinitionDeclaration struct {
	PackageID          string // name of package
	FileID             string // name of file
	StartOffset        int    // offset with in the file starting from 'type'
	Length             int    // length of entire declaration i.e. 'type [name] [data type]'
	LineNumber         int    // line number of declaration
	TypeDefinitionName string // name of function being declared
}

func ExtractTypeDefinitions(source []byte, fileName string) ([]TypeDefinitionDeclaration, error) {
	var TypeDefinitionDeclarationsArray []TypeDefinitionDeclaration
	var pkg string

	// Regexes
	rePkg := regexp.MustCompile(`^(?:.+\s+|\s*)package(?:\s+[\S\s]+|\s*)$`)
	rePkgName := regexp.MustCompile(`package\s+([_a-zA-Z][_a-zA-Z0-9]*)`)
	reType := regexp.MustCompile(`^(?:.+\s+|\s*)type(?:\s+[\S\s]+|\s*)$`)
	reTypeDefinition := regexp.MustCompile(`type\s+([_a-zA-Z][_a-zA-Z0-9]*)\s+([_a-zA-Z]\w*(?:\.[_a-zA-Z]\w*)*)`)

	reader := bytes.NewReader(source)
	scanner := bufio.NewScanner(reader)
	scanner.Split(scanLinesWithLineTerminator) // set scanner SplitFunc to custom ScanLines func at line 55

	var currentOffset int // offset of current line
	var lineno int        // line number

	//Reading code line by line
	for scanner.Scan() {
		line := scanner.Bytes()
		lineno++

		// Package declaration extraction
		if rePkg.FindIndex(line) != nil {

			matchPkg := rePkgName.FindSubmatch(line)

			if matchPkg == nil || !bytes.Equal(matchPkg[0], bytes.TrimSpace(line)) {
				return TypeDefinitionDeclarationsArray, fmt.Errorf("%v:%v: syntax error: package declaration", filepath.Base(fileName), lineno)
			}

			pkg = string(matchPkg[1])

		}

		if reType.Find(line) != nil {

			typeDefinition := reTypeDefinition.FindSubmatch(line)
			typeDefinitionIdx := reTypeDefinition.FindSubmatchIndex(line)

			if typeDefinition == nil {
				return TypeDefinitionDeclarationsArray, fmt.Errorf("%v:%v: syntax error: type definition declaration", filepath.Base(fileName), lineno)
			}

			if !bytes.Contains(typeDefinition[2], []byte("struct")) {

				if !bytes.Equal(typeDefinition[0], bytes.TrimSpace(line)) {
					return TypeDefinitionDeclarationsArray, fmt.Errorf("%v:%v: syntax error: type definition declaration", filepath.Base(fileName), lineno)
				}

				var typeDefinitionDeclaration TypeDefinitionDeclaration
				typeDefinitionDeclaration.PackageID = pkg
				typeDefinitionDeclaration.FileID = fileName
				typeDefinitionDeclaration.StartOffset = typeDefinitionIdx[0] + currentOffset
				typeDefinitionDeclaration.Length = typeDefinitionIdx[1] - typeDefinitionIdx[0]
				typeDefinitionDeclaration.LineNumber = lineno
				typeDefinitionDeclaration.TypeDefinitionName = string(typeDefinition[1])

				TypeDefinitionDeclarationsArray = append(TypeDefinitionDeclarationsArray, typeDefinitionDeclaration)

			}
		}

		currentOffset += len(line) // increments the currentOffset by line len

	}

	return TypeDefinitionDeclarationsArray, nil

}
