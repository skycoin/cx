package declaration_extractor

import (
	"bufio"
	"bytes"
	"fmt"
	"path/filepath"
	"regexp"

	"github.com/skycoin/cx/cmd/packageloader2/loader"
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
	reName := regexp.MustCompile(`[_a-zA-Z]\w*`)
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

		tokens := bytes.Fields(line)
		// Package declaration extraction
		if contains(tokens, []byte("package")) {
			if len(tokens) != 2 {
				return FuncDeclarationsArray, fmt.Errorf("%v:%v: syntax error: package declaration", filepath.Base(fileName), lineno)
			}
			name := reName.Find(tokens[1])
			if name == nil || len(name) != len(tokens[1]) {
				return FuncDeclarationsArray, fmt.Errorf("%v:%v: syntax error: package declaration", filepath.Base(fileName), lineno)

			}
			pkg = string(name)
		}

		if contains(tokens, []byte("func")) {

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
			funcDeclaration.Length = len(line)

			FuncDeclarationsArray = append(FuncDeclarationsArray, funcDeclaration)

		}

		currentOffset += len(line) // increments the currentOffset by line len
	}

	return FuncDeclarationsArray, nil
}

func ExtractMethod(fun FuncDeclaration, files []*loader.File) (string, error) {

	bytes, err := GetSourceBytes(files, filepath.Base(fun.FileID))
	if err != nil {
		return "", err
	}

	reFuncMethod := regexp.MustCompile(`func\s*\(\s*\w+\s+(\w+)\s*\)`)
	funcMethod := reFuncMethod.FindSubmatch(bytes[fun.StartOffset : fun.StartOffset+fun.Length])
	if funcMethod == nil {
		return "", nil
	}
	return string(funcMethod[1]), nil
}
