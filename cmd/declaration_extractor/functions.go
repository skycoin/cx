package declaration_extractor

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
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
		if ContainsTokenByte(tokens, []byte("package")) {

			name := reName.Find(tokens[1])

			if len(tokens) != 2 || len(tokens[1]) != len(name) {
				return FuncDeclarationsArray, fmt.Errorf("%v:%v: syntax error: package declaration", filepath.Base(fileName), lineno)
			}

			pkg = string(name)

		}

		if ContainsTokenByte(tokens, []byte("func")) || ContainsTokenByte(tokens, []byte("func(")) {

			funcBytes := reFuncDec.FindSubmatch(line)
			funcIdx := reFuncDec.FindSubmatchIndex(line)

			if pkg == "" {
				return FuncDeclarationsArray, fmt.Errorf("%v:%v: syntax error: missing package", filepath.Base(fileName), lineno)
			}

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

func ExtractMethod(fun FuncDeclaration) (string, error) {

	file, err := os.Open(fun.FileID)
	if err != nil {
		return "", err
	}

	tmp := bytes.NewBuffer(nil)
	io.Copy(tmp, file)
	bytes := tmp.Bytes()

	reFuncMethod := regexp.MustCompile(`func\s*\(\s*\w+\s+(\w+)\s*\)`)
	funcMethod := reFuncMethod.FindSubmatch(bytes[fun.StartOffset : fun.StartOffset+fun.Length])
	if funcMethod == nil {
		return "", nil
	}
	return string(funcMethod[1]), nil
}
