package declaration_extractor

import (
	"bufio"
	"bytes"
	"fmt"
	"path/filepath"
	"regexp"
)

type FuncDeclaration struct {
	PackageID   string // name of package
	FileID      string // name of file
	StartOffset int    // offset with in the file starting from 'func'
	Length      int    // length of entire declaration i.e. 'func [name] ()' or 'func [name] ([params])' or 'func [name] ([params]) [returns]'
	LineNumber  int    // line number of declaration
	FuncName    string // name of function being declared
}

func ExtractFuncs(source []byte, fileName string) ([]FuncDeclaration, error) {

	var FuncDeclarationsArray []FuncDeclaration
	var pkg string

	// Regexes
	rePkg := regexp.MustCompile(`^(?:.+\s+|\s*)package(?:\s+[\S\s]+|\s*)$`)
	rePkgName := regexp.MustCompile(`package\s+([_a-zA-Z][_a-zA-Z0-9]*)`)
	reFunc := regexp.MustCompile(`^(?:.+\s+|\s*)func(?:\s+[\S\s]+|\([\S\s]+|\s*)$`)
	reNotSpace := regexp.MustCompile(`\S+`)

	// Func Declaration regex for name extraction and syntax checking
	// Components:
	// func - func keyword
	// (?:\s*\(\s*[_a-zA-Z]\w*\s+\*{0,1}\s*[_a-zA-Z]\w*(?:\.[_a-zA-Z]\w*)*\s*\)\s*)|\s+) -  [space/no space] [([reciever object]) [name] | [space]]
	// ([_a-zA-Z]\w*) - name of func
	// (?:\s*\(\s*(?:(?:[_a-zA-Z]\w*\s+\*{0,1}\s*(?:\[(?:[1-9]\d+|[0-9]){0,1}\]\*{0,1}){0,1}\s*[_a-zA-Z]\w*(?:\.[_a-zA-Z]\w*)*\s*,\s*)+[_a-zA-Z]\w*\s+\*{0,1}\s*(?:\[(?:[1-9]\d+|[0-9]){0,1}\]\*{0,1}){0,1}\s*[_a-zA-Z]\w*(?:\.[_a-zA-Z]\w*)*|(?:[_a-zA-Z]\w*\s+\*{0,1}\s*(?:\[(?:[1-9]\d+|[0-9]){0,1}\]\*{0,1}){0,1}\s*[_a-zA-Z]\w*(?:\.[_a-zA-Z]\w*)*){0,1})\s*\)){1,2} - [[space/no space] ([params])]{1,2}
	//		(?:(?:[_a-zA-Z]\w*\s+\*{0,1}\s*(?:\[(?:[1-9]\d+|[0-9]){0,1}\]\*{0,1}){0,1}\s*[_a-zA-Z]\w*(?:\.[_a-zA-Z]\w*)*\s*,\s*)+[_a-zA-Z]\w*\s+\*{0,1}\s*(?:\[(?:[1-9]\d+|[0-9]){0,1}\]\*{0,1}){0,1}\s*[_a-zA-Z]\w*(?:\.[_a-zA-Z]\w*)*|(?:[_a-zA-Z]\w*\s+\*{0,1}\s*(?:\[(?:[1-9]\d+|[0-9]){0,1}\]\*{0,1}){0,1}\s*[_a-zA-Z]\w*(?:\.[_a-zA-Z]\w*)*){0,1}) - [[param name] [data type] [,]]{0,1} [param name] [data type] | [param name] [data type]
	// 			(?:[_a-zA-Z]\w*\s+\*{0,1}\s*(?:\[(?:[1-9]\d+|[0-9]){0,1}\]\*{0,1}){0,1}\s*[_a-zA-Z]\w*(?:\.[_a-zA-Z]\w*)* - [param name] [*]{0,1} [\[[1-9][0-9]+|[0-9]\][*]{0,1}]{0,1} [word] [[.][word]]*
	//
	// First, finds the func keyword
	// Second, finds out whether the function has a receiver object or not and extracts the func name
	// Third, finds whether there's one or two pairs of parenthesis
	// Forth, finds whether there's one or more params in the parenthesis
	reFuncDec := regexp.MustCompile(`(func(?:(?:\s*\(\s*[_a-zA-Z]\w*\s+\*{0,1}\s*[_a-zA-Z]\w*(?:\.[_a-zA-Z]\w*)*\s*\)\s*)|\s+)([_a-zA-Z]\w*)(?:\s*\(\s*(?:(?:[_a-zA-Z]\w*(?:\s*\*{0,1}\s*(?:\[(?:[1-9]\d+|[0-9]){0,1}\]\*{0,1}){0,1}\s*|\s+)[_a-zA-Z]\w*(?:\.[_a-zA-Z]\w*)*\s*,\s*)+[_a-zA-Z]\w*(?:\s*\*{0,1}\s*(?:\[(?:[1-9]\d+|[0-9]){0,1}\]\*{0,1}){0,1}|\s+)\s*[_a-zA-Z]\w*(?:\.[_a-zA-Z]\w*)*|(?:[_a-zA-Z]\w*(?:\s*\*{0,1}\s*(?:\[(?:[1-9]\d+|[0-9]){0,1}\]\*{0,1}){0,1}\s*|\s+)[_a-zA-Z]\w*(?:\.[_a-zA-Z]\w*)*){0,1})\s*\)){1,2})(?:\s*{){0,1}`)

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

			if matchPkg == nil || !bytes.Equal(matchPkg[0], bytes.TrimSpace(line)) && reNotSpace.Find(line) != nil {
				return FuncDeclarationsArray, fmt.Errorf("%v:%v: syntax error: package declaration", filepath.Base(fileName), lineno)
			}

			pkg = string(matchPkg[1])

		}

		if match := reFunc.FindIndex(line); match != nil {

			funcBytes := reFuncDec.FindSubmatch(line)
			funcIdx := reFuncDec.FindSubmatchIndex(line)

			if funcBytes == nil || !bytes.Equal(funcBytes[0], bytes.TrimSpace(line)) {
				return FuncDeclarationsArray, fmt.Errorf("%v:%v: syntax error: func declaration", filepath.Base(fileName), lineno)
			}

			var funcDeclaration FuncDeclaration
			funcDeclaration.PackageID = pkg
			funcDeclaration.FileID = fileName
			funcDeclaration.StartOffset = funcIdx[2] + currentOffset
			funcDeclaration.LineNumber = lineno
			funcDeclaration.FuncName = string(funcBytes[2])
			funcDeclaration.Length = funcIdx[3] - funcIdx[2]

			FuncDeclarationsArray = append(FuncDeclarationsArray, funcDeclaration)

		}

		currentOffset += len(line) // increments the currentOffset by line len
	}

	return FuncDeclarationsArray, nil
}
