package declaration_extractor

import (
	"bufio"
	"bytes"
	"fmt"
	"path/filepath"
	"regexp"
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

	//Regexes
	rePkg := regexp.MustCompile("package")
	rePkgName := regexp.MustCompile(`package\s+([_a-zA-Z][_a-zA-Z0-9]*)`)
	reEnumInit := regexp.MustCompile(`const\s+\(`)
	rePrtsClose := regexp.MustCompile(`\)`)
	reEnumDec := regexp.MustCompile(`([_a-zA-Z][_a-zA-Z0-9]*)(?:\s+([_a-zA-Z]\w*(?:\.[_a-zA-Z]\w*)*)){0,1}(?:\s*\=\s*[\s\S]+\S+){0,1}`)

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

		// Package declaration extraction
		if rePkg.FindIndex(line) != nil {

			matchPkg := rePkgName.FindSubmatch(line)

			if matchPkg == nil || !bytes.Equal(matchPkg[0], bytes.TrimSpace(line)) {
				return EnumDeclarationsArray, fmt.Errorf("%v:%v: syntax error: package declaration", filepath.Base(fileName), lineno)
			}

			pkg = string(matchPkg[1])

		}

		// initialize enum, increment parenthesis depth and skip to next line
		// if const ( is found
		if locs := reEnumInit.FindAllIndex(line, -1); locs != nil {
			EnumInit = true
			currentOffset += len(line) // increments the currentOffset by line len
			continue
		}

		// if ) is found and enum intialized, decrement parenthesis depth
		if locs := rePrtsClose.FindAllIndex(line, -1); locs != nil && EnumInit {
			EnumInit = false
			Type = ""
			Index = 0
		}

		// if match is found and enum initialized and parenthesis depth is 1
		if enumDec := reEnumDec.FindSubmatch(line); enumDec != nil && EnumInit {

			if !bytes.Equal(enumDec[0], bytes.TrimSpace(line)) {
				return EnumDeclarationsArray, fmt.Errorf("%v:%v: syntax error: enum declaration", filepath.Base(fileName), lineno)
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
