package declaration_extractor

import (
	"bufio"
	"bytes"
	"fmt"
	"path/filepath"
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
		if locs := reLeftBrace.FindAllIndex(line, -1); locs != nil {
			inBlock++
		}

		// if } is found decrement body depth
		if locs := reRightBrace.FindAllIndex(line, -1); locs != nil {
			inBlock--
		}

		if inBlock == 0 {

			// if match is found and body depth is 0
			if reGlobal.FindIndex(line) != nil {

				if pkg == "" {
					return GlobalDeclarationsArray, fmt.Errorf("%v:%v: syntax error: missing package", filepath.Base(fileName), lineno)
				}

				matchGlobal := reGlobalName.FindSubmatch(line)
				matchGlobalIdx := reGlobalName.FindIndex(line)

				if matchGlobal == nil || !bytes.Equal(matchGlobal[0], bytes.TrimSpace(line)) {
					return GlobalDeclarationsArray, fmt.Errorf("%v:%v: syntax error: global declaration", filepath.Base(fileName), lineno)
				}

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
