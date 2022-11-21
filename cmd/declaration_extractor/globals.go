package declaration_extractor

import (
	"bufio"
	"bytes"
	"fmt"
	"path/filepath"
	"regexp"
)

// All units for offset/length are in counted in bytes

type GlobalDeclaration struct {
	PackageID          string // name of package
	FileID             string // name of file
	StartOffset        int    // offset with in the file starting from 'var' of file
	Length             int    // length of entire declaration i.e. 'var [name] [type]'
	LineNumber         int    // line number of declaration
	GlobalVariableName string // name of variable being declared
}

func ExtractGlobals(source []byte, fileName string) ([]GlobalDeclaration, error) {

	var GlobalDeclarationsArray []GlobalDeclaration
	var pkg string

	//Regexs
	rePkg := regexp.MustCompile(`^(?:.+\s+|\s*)package(?:\s+[\S\s]+|\s*)$`)
	rePkgName := regexp.MustCompile(`package\s+([_a-zA-Z][_a-zA-Z0-9]*)`)
	reGlobal := regexp.MustCompile(`^(?:.+\s+|\s*)var(?:\s+[\S\s]+|\s*)$`)
	reGlobalName := regexp.MustCompile(`var\s+([_a-zA-Z]\w*)\s+\*{0,1}\s*(?:\[(?:[1-9]\d+|[0-9]){0,1}\]){0,1}\s*[_a-zA-Z]\w*(?:\.[_a-zA-Z]\w*)*(?:\s*\=\s*[\s\S]+\S+){0,1}`)
	reBodyOpen := regexp.MustCompile("{")
	reBodyClose := regexp.MustCompile("}")

	reader := bytes.NewReader(source)
	scanner := bufio.NewScanner(reader)
	scanner.Split(scanLinesWithLineTerminator) // set scanner SplitFunc to custom ScanLines func at line 55

	var currentOffset int // offset of current line
	var lineno int        // line number
	var inBlock int       // in Body { } depth

	for scanner.Scan() {

		line := scanner.Bytes()
		lineno++

		// Package declaration extraction
		if rePkg.FindIndex(line) != nil {

			matchPkg := rePkgName.FindSubmatch(line)

			if matchPkg == nil || !bytes.Equal(matchPkg[0], bytes.TrimSpace(line)) {
				return GlobalDeclarationsArray, fmt.Errorf("%v:%v: syntax error: package declaration", filepath.Base(fileName), lineno)
			}

			pkg = string(matchPkg[1])

		}

		// if {  is found increment body depth
		if locs := reBodyOpen.FindAllIndex(line, -1); locs != nil {
			inBlock++
		}

		// if } is found decrement body depth
		if locs := reBodyClose.FindAllIndex(line, -1); locs != nil {
			inBlock--
		}

		// if match is found and body depth is 0
		if reGlobal.FindIndex(line) != nil {

			matchGlobal := reGlobalName.FindSubmatch(line)
			matchGlobalIdx := reGlobalName.FindIndex(line)

			if matchGlobal == nil || !bytes.Equal(matchGlobal[0], bytes.TrimSpace(line)) {
				return GlobalDeclarationsArray, fmt.Errorf("%v:%v: syntax error: global declaration", filepath.Base(fileName), lineno)
			}

			if inBlock == 0 {

				var tmp GlobalDeclaration

				tmp.PackageID = pkg
				tmp.FileID = fileName

				tmp.StartOffset = matchGlobalIdx[0] + currentOffset // offset is match index + current line offset
				tmp.Length = matchGlobalIdx[1] - matchGlobalIdx[0]
				tmp.LineNumber = lineno

				// gets the name directly with submatch index + current line offset
				tmp.GlobalVariableName = string(matchGlobal[1])

				GlobalDeclarationsArray = append(GlobalDeclarationsArray, tmp)
			}

		}

		currentOffset += len(line) // increments the currentOffset by line len
	}

	return GlobalDeclarationsArray, nil

}
