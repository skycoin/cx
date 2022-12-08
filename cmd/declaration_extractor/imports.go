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

		// Package declaration extraction
		if rePkg.FindIndex(line) != nil {

			matchPkg := rePkgName.FindSubmatch(line)

			if matchPkg == nil || !bytes.Equal(matchPkg[0], bytes.TrimSpace(line)) {
				return ImportDeclarationsArray, fmt.Errorf("%v:%v: syntax error: package declaration", filepath.Base(fileName), lineno)
			}

			pkg = string(matchPkg[1])

		}

		// Extract Import
		checkLine := bytes.Split(line, []byte(" "))
		if bytes.Equal(checkLine[0], []byte("import")) {

			if pkg == "" {
				return ImportDeclarationsArray, fmt.Errorf("%v:%v: syntax error: missing package", filepath.Base(fileName), lineno)
			}

			var tmp ImportDeclaration

			tmp.PackageID = pkg
			tmp.FileID = fileName
			tmp.LineNumber = lineno

			tmp.ImportName = string(checkLine[1][1 : len(checkLine[1])-1])

			ImportDeclarationsArray = append(ImportDeclarationsArray, tmp)
		}
	}

	return ImportDeclarationsArray, nil
}
