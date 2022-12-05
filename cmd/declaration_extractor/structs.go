package declaration_extractor

import (
	"bufio"
	"bytes"
	"fmt"
	"path/filepath"
	"regexp"
)

type StructDeclaration struct {
	PackageID    string         // name of package
	FileID       string         // name of file
	StartOffset  int            // offset with in the file starting from 'type'
	Length       int            // length of entire declaration i.e. 'type [name] [type]'
	LineNumber   int            // line number of declaration
	StructName   string         // name of struct being declared
	StructFields []*StructField // array of fields
}

type StructField struct {
	StartOffset     int    // offset with in the file starting from [name]
	Length          int    // length of entire declaration i.e. '[name] [type]'
	LineNumber      int    // line number of declaration
	StructFieldName string // name of variable being declared
}

func ExtractStructs(source []byte, fileName string) ([]StructDeclaration, error) {

	var StructDeclarationsArray []StructDeclaration
	var pkg string

	// Regexes
	reNotSpace := regexp.MustCompile(`\S+`)
	reName := regexp.MustCompile(`[_a-zA-Z]\w*`)
	// reStruct := regexp.MustCompile(`type\s+[_a-zA-Z][_a-zA-Z0-9]*\s+struct`)
	// reStructHeader := regexp.MustCompile(`type\s+([_a-zA-Z][_a-zA-Z0-9]*)\s+struct\s*{`)
	reRightBrace := regexp.MustCompile("}")
	reStructField := regexp.MustCompile(`([_a-zA-Z][_a-zA-Z0-9]*)\s+\*{0,1}\s*(?:\[(?:[1-9]\d+|[0-9]){0,1}\]\*{0,1}){0,1}\s*[_a-zA-Z]\w*(?:\.[_a-zA-Z]\w*)*`)

	reader := bytes.NewReader(source)
	scanner := bufio.NewScanner(reader)
	scanner.Split(scanLinesWithLineTerminator) // set scanner SplitFunc to custom ScanLines func at line 55

	var currentOffset int                   // offset of current line
	var lineno int                          // line number
	var inBlock bool                        // inBlock
	var structDeclaration StructDeclaration // temporary variable for Struct Declaration
	var structFieldsArray []*StructField    // struct fields

	//Reading code line by line
	for scanner.Scan() {
		line := scanner.Bytes()
		lineno++

		tokens := bytes.Fields(line)
		// Package declaration extraction
		if contains(tokens, []byte("package")) {
			if len(tokens) != 2 {
				return StructDeclarationsArray, fmt.Errorf("%v:%v: syntax error: package declaration", filepath.Base(fileName), lineno)
			}
			name := reName.Find(tokens[1])
			if name == nil || len(name) != len(tokens[1]) {
				return StructDeclarationsArray, fmt.Errorf("%v:%v: syntax error: package declaration", filepath.Base(fileName), lineno)

			}
			pkg = string(name)
		}

		// if struct declaration is found
		// i.e. type [name] [type]
		if len(tokens) >= 3 && contains(tokens, []byte("type")) && contains(tokens, []byte("struct")) {

			name := reName.Find(tokens[1])
			if name == nil || len(name) != len(tokens[1]) {
				return StructDeclarationsArray, fmt.Errorf("%v:%v: syntax error: struct declaration", filepath.Base(fileName), lineno)
			}

			if pkg == "" {
				return StructDeclarationsArray, fmt.Errorf("%v:%v: no package declared for global declaration", filepath.Base(fileName), lineno)
			}

			structDeclaration.PackageID = pkg
			structDeclaration.FileID = fileName

			structDeclaration.StartOffset = currentOffset // offset is current line offset + match index
			structDeclaration.Length = len(line)
			structDeclaration.StructName = string(name)

			structDeclaration.LineNumber = lineno

			inBlock = true

		}

		if match := reRightBrace.FindIndex(line); match != nil && inBlock {

			inBlock = false
			structDeclaration.StructFields = structFieldsArray
			StructDeclarationsArray = append(StructDeclarationsArray, structDeclaration)
			structFieldsArray = []*StructField{}

		}

		if inBlock && structDeclaration.LineNumber < lineno {

			var structField StructField
			matchStructField := reStructField.FindSubmatch(line)
			matchStructFieldIdx := reStructField.FindSubmatchIndex(line)

			if reNotSpace.Find(line) != nil && (matchStructField == nil || !bytes.Equal(matchStructField[0], bytes.TrimSpace(line))) {
				return StructDeclarationsArray, fmt.Errorf("%v:%v: syntax error:struct field", filepath.Base(fileName), lineno)
			}

			if matchStructField != nil {
				structField.StartOffset = matchStructFieldIdx[0] + currentOffset
				structField.Length = matchStructFieldIdx[1] - matchStructFieldIdx[0]
				structField.LineNumber = lineno
				structField.StructFieldName = string(matchStructField[1])
				structFieldsArray = append(structFieldsArray, &structField)
			}
		}

		currentOffset += len(line) // increments the currentOffset by line len
	}

	return StructDeclarationsArray, nil
}
