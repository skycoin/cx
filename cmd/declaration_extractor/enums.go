package declaration_extractor

import (
	"bufio"
	"bytes"
	"fmt"
	"path/filepath"
)

type EnumDeclaration struct {
	PackageID   string // name of package
	FileID      string // name of file
	StartOffset int    // offset with in the file starting from '[name of enum]' of file
	Length      int    // length of entire declaration i.e. '[name] [type]' or '[name]'
	LineNumber  int    // line number of declaration
	Type        string // type of enum
	Value       int    // value of enum
	EnumName    string // name of enum being declared
}

func ExtractEnums(source []byte, fileName string) ([]EnumDeclaration, error) {

	var EnumDeclarationsArray []EnumDeclaration
	var pkg string

	reader := bytes.NewReader(source)
	scanner := bufio.NewScanner(reader)
	scanner.Split(scanLinesWithLineTerminator) // set scanner SplitFunc to custom ScanLines func at line 55

	var EnumInit bool     // is in a enum declaration
	var Type string       // type for later enum declaration
	var Index int         // index for enum declaration
	var currentOffset int // offset of current line
	var lineno int        // line number

	for scanner.Scan() {

		line := scanner.Bytes()
		lineno++

		tokens := bytes.Fields(line)
		// Package declaration extraction
		if contains(tokens, []byte("package")) {
			if len(tokens) != 2 {
				return EnumDeclarationsArray, fmt.Errorf("%v:%v: syntax error: package declaration", filepath.Base(fileName), lineno)
			}
			name := reName.Find(tokens[1])
			if name == nil || len(name) != len(tokens[1]) {
				return EnumDeclarationsArray, fmt.Errorf("%v:%v: syntax error: package declaration", filepath.Base(fileName), lineno)

			}
			pkg = string(name)
		}

		// initialize enum, increment parenthesis depth and skip to next line
		// if const ( is found
		if locs := reEnumInit.FindIndex(line); locs != nil {
			EnumInit = true
			currentOffset += len(line) // increments the currentOffset by line len
			continue
		}

		// if ) is found and enum intialized, decrement parenthesis depth
		if contains(tokens, []byte(")")) && EnumInit {
			EnumInit = false
			Type = ""
			Index = 0
		}

		// if match is found and enum initialized and parenthesis depth is 1
		if len(tokens) > 0 && EnumInit {

			enumDec := reEnumDec.FindSubmatch(line)
			if !bytes.Equal(enumDec[0], bytes.TrimSpace(line)) {
				return EnumDeclarationsArray, fmt.Errorf("%v:%v: syntax error: enum declaration", filepath.Base(fileName), lineno)
			}

			if pkg == "" {
				return EnumDeclarationsArray, fmt.Errorf("%v:%v: no package declared for enum declaration", filepath.Base(fileName), lineno)
			}

			enumDecIdx := reEnumDec.FindIndex(line)

			var tmp EnumDeclaration

			tmp.PackageID = pkg
			tmp.FileID = fileName

			tmp.StartOffset = enumDecIdx[0] + currentOffset // offset is current line offset + match index
			tmp.Length = enumDecIdx[1] - enumDecIdx[0]
			tmp.LineNumber = lineno

			tmp.EnumName = string(enumDec[1])

			// If there is type declaration for Enum declaration
			// i.e. [enum] [type] = iota
			// set the type to type in declaration
			if enumDec[2] != nil {
				Type = string(enumDec[2])
			}

			// otherwise set it to previous type
			tmp.Type = Type

			tmp.Value = Index
			EnumDeclarationsArray = append(EnumDeclarationsArray, tmp)
			Index++
		}

		currentOffset += len(line) // increments the currentOffset by line len
	}

	return EnumDeclarationsArray, nil

}
