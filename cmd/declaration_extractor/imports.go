package declaration_extractor

import (
	"bufio"
	"bytes"
	"fmt"
	"path/filepath"
)

type ImportDeclaration struct {
	PackageID  string // name of package
	FileID     string // name of file
	LineNumber int    // line number of declaration
	ImportName string // name of import being declared
}

func ExtractImports(source []byte, fileName string) ([]ImportDeclaration, error) {

	var ImportDeclarationsArray []ImportDeclaration
	var pkg string

	reader := bytes.NewReader(source)
	scanner := bufio.NewScanner(reader)

	var lineno int // line number

	for scanner.Scan() {

		line := scanner.Bytes()
		lineno++

		tokens := bytes.Fields(line)
		// Package declaration extraction
		if contains(tokens, []byte("package")) {
			if len(tokens) != 2 {
				return ImportDeclarationsArray, fmt.Errorf("%v:%v: syntax error: package declaration", filepath.Base(fileName), lineno)
			}
			name := reName.Find(tokens[1])
			if name == nil || len(name) != len(tokens[1]) {
				return ImportDeclarationsArray, fmt.Errorf("%v:%v: syntax error: package declaration", filepath.Base(fileName), lineno)

			}
			pkg = string(name)
		}

		// Extract Import
		if contains(tokens, []byte("import")) {

			if pkg == "" {
				return ImportDeclarationsArray, fmt.Errorf("%v:%v: no package declared for global declaration", filepath.Base(fileName), lineno)
			}

			var tmp ImportDeclaration

			tmp.PackageID = pkg
			tmp.FileID = fileName
			tmp.LineNumber = lineno

			tmp.ImportName = string(tokens[1][1 : len(tokens[1])-1])

			ImportDeclarationsArray = append(ImportDeclarationsArray, tmp)
		}
	}

	return ImportDeclarationsArray, nil
}
