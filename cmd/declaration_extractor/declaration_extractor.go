package declaration_extractor

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"sync"
	"unicode"
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

type TypeDefinitionDeclaration struct {
	PackageID          string // name of package
	FileID             string // name of file
	StartOffset        int    // offset with in the file starting from 'type'
	Length             int    // length of entire declaration i.e. 'type [name] [data type]'
	LineNumber         int    // line number of declaration
	TypeDefinitionName string // name of function being declared
}

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

type FuncDeclaration struct {
	PackageID   string // name of package
	FileID      string // name of file
	StartOffset int    // offset with in the file starting from 'func'
	Length      int    // length of entire declaration i.e. 'func [name] ()' or 'func [name] ([params])' or 'func [name] ([params]) [returns]'
	LineNumber  int    // line number of declaration
	FuncName    string // name of function being declared
}

// Modified ScanLines to include "\r\n" or "\n" in line
func scanLinesWithLineTerminator(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	if i := bytes.IndexAny(data, "\r\n"); i >= 0 {
		advance = i + 1 // i + 1 includes the line termninator
		if data[i] == '\n' {
			// We have a line terminated by single newline.
			return advance, data[0:advance], nil
		}

		if len(data) > i+1 && data[i+1] == '\n' {
			advance += 1
		}
		return advance, data[0:advance], nil
	}
	// If we're at EOF, we have a final, non-terminated line. Return it.
	if atEOF {
		return len(data), data, nil
	}
	// Request more data.
	return 0, nil, nil
}

func ReplaceCommentsWithWhitespaces(source []byte) []byte {

	reComment := regexp.MustCompile(`//.*|/\*[\s\S]*?\*/|\"//.*\"`)

	comments := reComment.FindAllIndex(source, -1)

	// Loops through each character and replaces with whitespcace
	for i := range comments {
		for loc := comments[i][0]; loc < comments[i][1]; loc++ {
			if unicode.IsSpace(rune(source[loc])) {
				continue
			}
			source[loc] = byte(' ')
		}
	}

	return source
}

func ReplaceStringContentsWithWhitespaces(source []byte) ([]byte, error) {

	var sourceWithoutStringContents []byte
	sourceWithoutStringContents = append(sourceWithoutStringContents, source...)
	var inStdString bool
	var inRawString bool
	var lineno int

	for i, char := range source {

		if char == '\n' {
			lineno++
		}

		//if end of line and quote not terminated
		if char == '\n' && inStdString {
			return sourceWithoutStringContents, fmt.Errorf("%v: syntax error: quote not terminated", lineno)
		}

		if char == '"' && !inStdString && !inRawString {
			inStdString = true
			continue
		}

		if char == '"' && inStdString {
			inStdString = false
			continue
		}

		if char == '`' && !inRawString && !inStdString {
			inRawString = true
			continue
		}

		if char == '`' && inRawString {
			inRawString = false
			continue
		}

		if (inStdString || inRawString) && !unicode.IsSpace(rune(char)) {
			sourceWithoutStringContents[i] = byte(' ')
		}
	}

	return sourceWithoutStringContents, nil
}

func ExtractGlobals(source []byte, fileName string) ([]GlobalDeclaration, error) {

	var GlobalDeclarationsArray []GlobalDeclaration
	var pkg string

	//Regexs
	rePkg := regexp.MustCompile("package")
	rePkgName := regexp.MustCompile(`package\s+([_a-zA-Z][_a-zA-Z0-9]*)`)
	reGlobal := regexp.MustCompile("var")
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
		if reGlobal.FindAllIndex(line, -1) != nil {

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

func ExtractTypeDefinitions(source []byte, fileName string) ([]TypeDefinitionDeclaration, error) {

	var TypeDefinitionDeclarationsArray []TypeDefinitionDeclaration
	var pkg string

	// Regexes
	rePkg := regexp.MustCompile("package")
	rePkgName := regexp.MustCompile(`package\s+([_a-zA-Z][_a-zA-Z0-9]*)`)
	reType := regexp.MustCompile(`type`)
	reTypeDefinition := regexp.MustCompile(`type\s+([_a-zA-Z][_a-zA-Z0-9]*)\s+([_a-zA-Z]\w*(?:\.[_a-zA-Z]\w*)*)`)

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

			if matchPkg == nil || !bytes.Equal(matchPkg[0], bytes.TrimSpace(line)) {
				return TypeDefinitionDeclarationsArray, fmt.Errorf("%v:%v: syntax error: package declaration", filepath.Base(fileName), lineno)
			}

			pkg = string(matchPkg[1])

		}

		if reType.Find(line) != nil {

			typeDefinition := reTypeDefinition.FindSubmatch(line)
			typeDefinitionIdx := reTypeDefinition.FindSubmatchIndex(line)

			if typeDefinition == nil {
				return TypeDefinitionDeclarationsArray, fmt.Errorf("%v:%v: syntax error: type definition declaration", filepath.Base(fileName), lineno)
			}

			if !bytes.Contains(typeDefinition[2], []byte("struct")) {

				if !bytes.Equal(typeDefinition[0], bytes.TrimSpace(line)) {
					return TypeDefinitionDeclarationsArray, fmt.Errorf("%v:%v: syntax error: type definition declaration", filepath.Base(fileName), lineno)
				}

				var typeDefinitionDeclaration TypeDefinitionDeclaration
				typeDefinitionDeclaration.PackageID = pkg
				typeDefinitionDeclaration.FileID = fileName
				typeDefinitionDeclaration.StartOffset = typeDefinitionIdx[0] + currentOffset
				typeDefinitionDeclaration.Length = typeDefinitionIdx[1] - typeDefinitionIdx[0]
				typeDefinitionDeclaration.LineNumber = lineno
				typeDefinitionDeclaration.TypeDefinitionName = string(typeDefinition[1])

				TypeDefinitionDeclarationsArray = append(TypeDefinitionDeclarationsArray, typeDefinitionDeclaration)

			}
		}

		currentOffset += len(line) // increments the currentOffset by line len

	}

	return TypeDefinitionDeclarationsArray, nil

}

func ExtractStructs(source []byte, fileName string) ([]StructDeclaration, error) {

	var StructDeclarationsArray []StructDeclaration
	var pkg string

	// Regexes
	reNotSpace := regexp.MustCompile(`\S+`)
	rePkg := regexp.MustCompile("package")
	rePkgName := regexp.MustCompile(`package\s+([_a-zA-Z][_a-zA-Z0-9]*)`)
	reStruct := regexp.MustCompile(`type\s+[_a-zA-Z][_a-zA-Z0-9]*\s+struct`)
	reStructHeader := regexp.MustCompile(`type\s+([_a-zA-Z][_a-zA-Z0-9]*)\s+struct\s*{`)
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

		// Package declaration extraction
		if rePkg.FindIndex(line) != nil {

			matchPkg := rePkgName.FindSubmatch(line)

			if matchPkg == nil || !bytes.Equal(matchPkg[0], bytes.TrimSpace(line)) {
				return StructDeclarationsArray, fmt.Errorf("%v:%v: syntax error: package declaration", filepath.Base(fileName), lineno)
			}

			pkg = string(matchPkg[1])

		}

		// if struct declaration is found
		// i.e. type [name] [type]
		if strct := reStruct.FindIndex(line); strct != nil {

			structHeader := reStructHeader.FindSubmatch(line)
			if structHeader == nil || !bytes.Equal(structHeader[0], bytes.TrimSpace(line)) {
				return StructDeclarationsArray, fmt.Errorf("%v:%v: syntax error: struct declaration", filepath.Base(fileName), lineno)
			}

			structDeclaration.PackageID = pkg
			structDeclaration.FileID = fileName

			structDeclaration.StartOffset = strct[0] + currentOffset // offset is current line offset + match index
			structDeclaration.Length = strct[1] - strct[0]
			structDeclaration.StructName = string(structHeader[1])

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

func ExtractFuncs(source []byte, fileName string) ([]FuncDeclaration, error) {

	var FuncDeclarationsArray []FuncDeclaration
	var pkg string

	// Regexes
	rePkg := regexp.MustCompile("package")
	rePkgName := regexp.MustCompile(`package\s+([_a-zA-Z][_a-zA-Z0-9]*)`)
	reFunc := regexp.MustCompile(`func`)
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
	reFuncDec := regexp.MustCompile(`(func(?:(?:\s*\(\s*[_a-zA-Z]\w*\s+\*{0,1}\s*[_a-zA-Z]\w*(?:\.[_a-zA-Z]\w*)*\s*\)\s*)|\s+)([_a-zA-Z]\w*)(?:\s*\(\s*(?:(?:[_a-zA-Z]\w*\s+\*{0,1}\s*(?:\[(?:[1-9]\d+|[0-9]){0,1}\]\*{0,1}){0,1}\s*[_a-zA-Z]\w*(?:\.[_a-zA-Z]\w*)*\s*,\s*)+[_a-zA-Z]\w*\s+\*{0,1}\s*(?:\[(?:[1-9]\d+|[0-9]){0,1}\]\*{0,1}){0,1}\s*[_a-zA-Z]\w*(?:\.[_a-zA-Z]\w*)*|(?:[_a-zA-Z]\w*\s+\*{0,1}\s*(?:\[(?:[1-9]\d+|[0-9]){0,1}\]\*{0,1}){0,1}\s*[_a-zA-Z]\w*(?:\.[_a-zA-Z]\w*)*){0,1})\s*\)){1,2})(?:\s*{){0,1}`)

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

func ReDeclarationCheck(Glbl []GlobalDeclaration, Enum []EnumDeclaration, TypeDef []TypeDefinitionDeclaration, Strct []StructDeclaration, Func []FuncDeclaration) error {

	// Checks for the first declaration redeclared
	// in the order:
	// Global -> Enum -> Type Definition -> Struct -> Func

	for i := 0; i < len(Glbl); i++ {
		for j := i + 1; j < len(Glbl); j++ {
			if Glbl[i].GlobalVariableName == Glbl[j].GlobalVariableName && Glbl[i].PackageID == Glbl[j].PackageID {
				return fmt.Errorf("%v:%v: redeclaration error: global: %v", filepath.Base(Glbl[j].FileID), Glbl[j].LineNumber, Glbl[i].GlobalVariableName)
			}
		}
	}

	for i := 0; i < len(Enum); i++ {
		for j := i + 1; j < len(Enum); j++ {
			if Enum[i].EnumName == Enum[j].EnumName && Enum[i].PackageID == Enum[j].PackageID {
				return fmt.Errorf("%v:%v: redeclaration error: enum: %v", filepath.Base(Enum[j].FileID), Enum[j].LineNumber, Enum[i].EnumName)
			}
		}
	}

	for i := 0; i < len(TypeDef); i++ {
		for j := i + 1; j < len(TypeDef); j++ {
			if TypeDef[i].TypeDefinitionName == TypeDef[j].TypeDefinitionName && TypeDef[i].PackageID == TypeDef[j].PackageID {
				return fmt.Errorf("%v:%v: redeclaration error: type definition: %v", filepath.Base(TypeDef[j].FileID), TypeDef[j].LineNumber, TypeDef[i].TypeDefinitionName)
			}
		}
	}

	for i := 0; i < len(Strct); i++ {

		StructFields := Strct[i].StructFields
		for m := 0; m < len(StructFields); m++ {
			for n := m + 1; n < len(StructFields); n++ {
				if StructFields[m].StructFieldName == StructFields[n].StructFieldName {
					return fmt.Errorf("%v:%v: redeclaration error: struct field: %v", filepath.Base(Strct[i].FileID), StructFields[n].LineNumber, StructFields[n].StructFieldName)
				}
			}
		}

		for j := i + 1; j < len(Strct); j++ {
			if Strct[i].StructName == Strct[j].StructName && Strct[i].PackageID == Strct[j].PackageID {
				return fmt.Errorf("%v:%v: redeclaration error: struct: %v", filepath.Base(Strct[j].FileID), Strct[j].LineNumber, Strct[i].StructName)
			}
		}
	}

	for i := 0; i < len(Func); i++ {
		for j := i + 1; j < len(Func); j++ {
			if Func[i].FuncName == Func[j].FuncName && Func[i].PackageID == Func[j].PackageID {
				return fmt.Errorf("%v:%v: redeclaration error: func: %v", filepath.Base(Func[j].FileID), Func[j].LineNumber, Func[i].FuncName)
			}
		}
	}

	return nil
}

func GetDeclarations(source []byte, Glbls []GlobalDeclaration, Enums []EnumDeclaration, TypeDefs []TypeDefinitionDeclaration, Strcts []StructDeclaration, Funcs []FuncDeclaration) []string {

	var declarations []string

	for _, glbl := range Glbls {
		declarations = append(declarations, string(source[glbl.StartOffset:glbl.StartOffset+glbl.Length]))
	}

	for _, enum := range Enums {
		declarations = append(declarations, string(source[enum.StartOffset:enum.StartOffset+enum.Length]))
	}

	for _, typeDef := range TypeDefs {
		declarations = append(declarations, string(source[typeDef.StartOffset:typeDef.StartOffset+typeDef.Length]))
	}

	for _, strct := range Strcts {
		declarations = append(declarations, string(source[strct.StartOffset:strct.StartOffset+strct.Length]))
	}

	for _, fun := range Funcs {
		declarations = append(declarations, string(source[fun.StartOffset:fun.StartOffset+fun.Length]))
	}

	return declarations
}

func ExtractAllDeclarations(source []*os.File) ([]GlobalDeclaration, []EnumDeclaration, []TypeDefinitionDeclaration, []StructDeclaration, []FuncDeclaration, error) {

	//Variable declarations
	var Globals []GlobalDeclaration
	var Enums []EnumDeclaration
	var TypeDefinitions []TypeDefinitionDeclaration
	var Structs []StructDeclaration
	var Funcs []FuncDeclaration

	//Channel declarations
	globalChannel := make(chan []GlobalDeclaration, len(source))
	enumChannel := make(chan []EnumDeclaration, len(source))
	typeDefinitionChannel := make(chan []TypeDefinitionDeclaration, len(source))
	structChannel := make(chan []StructDeclaration, len(source))
	funcChannel := make(chan []FuncDeclaration, len(source))
	errorChannel := make(chan error, len(source))

	var wg sync.WaitGroup

	// concurrent extractions start
	for _, currentFile := range source {

		wg.Add(1)

		go func(currentFile *os.File, globalChannel chan<- []GlobalDeclaration, enumChannel chan<- []EnumDeclaration, typeDefinition chan<- []TypeDefinitionDeclaration, structChannel chan<- []StructDeclaration, funcChannel chan<- []FuncDeclaration, errorChannel chan<- error, wg *sync.WaitGroup) {

			defer wg.Done()

			srcBytes, err := io.ReadAll(currentFile)
			if err != nil {
				errorChannel <- err
				return
			}

			fileName := currentFile.Name()
			replaceComments := ReplaceCommentsWithWhitespaces(srcBytes)
			replaceStringContents, err := ReplaceStringContentsWithWhitespaces(replaceComments)
			if err != nil {
				errorChannel <- fmt.Errorf("%v:%v", filepath.Base(fileName), err)
				return
			}

			wg.Add(5)

			go func(globalChannel chan<- []GlobalDeclaration, replaceStringContents []byte, fileName string, wg *sync.WaitGroup) {

				defer wg.Done()

				globals, err := ExtractGlobals(replaceStringContents, fileName)

				if err != nil {
					errorChannel <- err
					return
				}

				globalChannel <- globals

			}(globalChannel, replaceStringContents, fileName, wg)

			go func(enumChannel chan<- []EnumDeclaration, replaceStringContents []byte, fileName string, wg *sync.WaitGroup) {

				defer wg.Done()

				enums, err := ExtractEnums(replaceStringContents, fileName)

				if err != nil {
					errorChannel <- err
					return
				}

				enumChannel <- enums

			}(enumChannel, replaceStringContents, fileName, wg)

			go func(typeDefinitionChannel chan<- []TypeDefinitionDeclaration, replaceStringContents []byte, fileName string, wg *sync.WaitGroup) {

				defer wg.Done()

				typeDefinitions, err := ExtractTypeDefinitions(replaceStringContents, fileName)

				if err != nil {
					errorChannel <- err
					return
				}

				typeDefinitionChannel <- typeDefinitions

			}(typeDefinitionChannel, replaceStringContents, fileName, wg)

			go func(structChannel chan<- []StructDeclaration, replaceStringContents []byte, fileName string, wg *sync.WaitGroup) {

				defer wg.Done()

				structs, err := ExtractStructs(replaceStringContents, fileName)

				if err != nil {
					errorChannel <- err
					return
				}

				structChannel <- structs

			}(structChannel, replaceStringContents, fileName, wg)

			go func(funcChannel chan<- []FuncDeclaration, replaceStringContents []byte, fileName string, wg *sync.WaitGroup) {

				defer wg.Done()

				funcs, err := ExtractFuncs(replaceStringContents, fileName)

				if err != nil {
					errorChannel <- err
					return
				}

				funcChannel <- funcs

			}(funcChannel, replaceStringContents, fileName, wg)

		}(currentFile, globalChannel, enumChannel, typeDefinitionChannel, structChannel, funcChannel, errorChannel, &wg)
	}

	wg.Wait()

	// Close all channels for reading
	close(globalChannel)
	close(enumChannel)
	close(typeDefinitionChannel)
	close(structChannel)
	close(funcChannel)
	close(errorChannel)

	//Read from channels concurrently
	wg.Add(5)

	go func() {

		for global := range globalChannel {
			Globals = append(Globals, global...)
		}

		wg.Done()

	}()

	go func() {

		for enum := range enumChannel {
			Enums = append(Enums, enum...)
		}

		wg.Done()
	}()

	go func() {

		for typeDef := range typeDefinitionChannel {
			TypeDefinitions = append(TypeDefinitions, typeDef...)
		}

		wg.Done()

	}()

	go func() {

		for strct := range structChannel {
			Structs = append(Structs, strct...)
		}

		wg.Done()

	}()

	go func() {

		for fun := range funcChannel {
			Funcs = append(Funcs, fun...)
		}

		wg.Done()

	}()

	wg.Wait()

	// there's an error, return values with first error
	if err := <-errorChannel; err != nil {
		return Globals, Enums, TypeDefinitions, Structs, Funcs, err
	}

	reDeclarationCheck := ReDeclarationCheck(Globals, Enums, TypeDefinitions, Structs, Funcs)

	// there's declaration redeclared return values with error
	if reDeclarationCheck != nil {
		return Globals, Enums, TypeDefinitions, Structs, Funcs, reDeclarationCheck
	}

	return Globals, Enums, TypeDefinitions, Structs, Funcs, nil
}
