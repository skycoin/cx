package declaration_extractor

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
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

	var sourceWithoutComments []byte = source

	reComment := regexp.MustCompile(`//.*|/\*[\s\S]*?\*/|\"//.*\"`)

	comments := reComment.FindAllIndex(source, -1)

	// Loops through each character and replaces with whitespcace
	for i := range comments {
		for loc := comments[i][0]; loc < comments[i][1]; loc++ {
			if unicode.IsSpace(rune(sourceWithoutComments[loc])) {
				continue
			}
			sourceWithoutComments[loc] = byte(' ')
		}
	}

	return sourceWithoutComments
}

func ExtractGlobals(source []byte, fileName string) ([]GlobalDeclaration, error) {

	var GlobalDeclarationsArray []GlobalDeclaration
	var pkg string

	//Regexs
	rePkg := regexp.MustCompile("package")
	rePkgName := regexp.MustCompile(`(^|[\s])package\s+([_a-zA-Z][_a-zA-Z0-9]*)`)
	reGlobal := regexp.MustCompile("var")
	reGlobalName := regexp.MustCompile(`var\s([_a-zA-Z][_a-zA-Z0-9]*)\s+[\[_a-zA-Z][\]_a-zA-Z0-9]*`)
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
			tokens := bytes.Fields(line)

			if !bytes.Equal(tokens[0], []byte("package")) {
				col := bytes.IndexAny(line, string(tokens[0]))

				return GlobalDeclarationsArray, fmt.Errorf("%d:%d %s", lineno, col, "syntax error: unexpected IDENTIFIER")
			}

			if len(tokens) == 1 {
				col := bytes.LastIndex(line, tokens[0])

				return GlobalDeclarationsArray, fmt.Errorf("%d:%d %s", lineno, col, "syntax error: unexpected IDENTIFIER")
			}

			if len(tokens) > 2 {
				col := bytes.IndexAny(line, string(tokens[2]))

				return GlobalDeclarationsArray, fmt.Errorf("%d:%d %s", lineno, col, "syntax error: unexpected IDENTIFIER")
			}

			if match := rePkgName.FindStringSubmatch(string(line)); match != nil {
				pkg = match[2]
			}

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

			tokens := bytes.Fields(line)

			if !bytes.Equal(tokens[0], []byte("var")) {
				col := bytes.Index(line, tokens[0])
				return GlobalDeclarationsArray, fmt.Errorf("%d:%d %s", lineno, col, "syntax error: unexpected IDENTIFIER")
			}

			if match := reGlobalName.FindSubmatchIndex(line); match != nil && inBlock == 0 {

				var tmp GlobalDeclaration

				tmp.PackageID = pkg
				tmp.FileID = fileName

				tmp.StartOffset = match[0] + currentOffset // offset is match index + current line offset
				tmp.Length = match[1] - match[0]
				tmp.LineNumber = lineno

				// gets the name directly with submatch index + current line offset
				tmp.GlobalVariableName = string(source[match[2]+currentOffset : match[3]+currentOffset])

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
	rePkgName := regexp.MustCompile(`(^|[\s])package\s+([_a-zA-Z][_a-zA-Z0-9]*)`)
	reEnumInit := regexp.MustCompile(`const\s+\(`)
	rePrtsOpen := regexp.MustCompile(`\(`)
	rePrtsClose := regexp.MustCompile(`\)`)
	reEnumDec := regexp.MustCompile(`([_a-zA-Z][_a-zA-Z0-9]*)\s+([_a-zA-Z][_a-zA-Z0-9]*)|([_a-zA-Z][_a-zA-Z0-9]*)`)

	reader := bytes.NewReader(source)
	scanner := bufio.NewScanner(reader)
	scanner.Split(scanLinesWithLineTerminator) // set scanner SplitFunc to custom ScanLines func at line 55

	var EnumInit bool     // is in a enum declaration
	var inPrts int        // parenthesis depth
	var Type string       // type for later enum declaration
	var Index int         // index for enum declaration
	var currentOffset int // offset of current line
	var lineno int        // line number

	for scanner.Scan() {

		line := scanner.Bytes()
		lineno++

		// Package declaration extraction
		if rePkg.FindIndex(line) != nil {
			tokens := bytes.Fields(line)

			if !bytes.Equal(tokens[0], []byte("package")) {
				col := bytes.IndexAny(line, string(tokens[0]))

				return EnumDeclarationsArray, fmt.Errorf("%d:%d %s", lineno, col, "syntax error: unexpected IDENTIFIER")
			}

			if len(tokens) == 1 {
				col := bytes.LastIndex(line, tokens[0])

				return EnumDeclarationsArray, fmt.Errorf("%d:%d %s", lineno, col, "syntax error: unexpected IDENTIFIER")
			}

			if len(tokens) > 2 {
				col := bytes.IndexAny(line, string(tokens[2]))

				return EnumDeclarationsArray, fmt.Errorf("%d:%d %s", lineno, col, "syntax error: unexpected IDENTIFIER")
			}

			if match := rePkgName.FindStringSubmatch(string(line)); match != nil {
				pkg = match[2]
			}

		}

		// initialize enum, increment parenthesis depth and skip to next line
		// if const ( is found
		if locs := reEnumInit.FindAllIndex(line, -1); locs != nil {
			EnumInit = true
			inPrts++
			currentOffset += len(line) // increments the currentOffset by line len
			continue
		}

		// if ( is found and enum initialized, increment parenthesis depth
		if locs := rePrtsOpen.FindAllIndex(line, -1); locs != nil && EnumInit {
			inPrts++
		}

		// if ) is found and enum intialized, decrement parenthesis depth
		if locs := rePrtsClose.FindAllIndex(line, -1); locs != nil && EnumInit {
			inPrts--
		}

		// if parenthesis depth is 0, reset all enum related variables
		if inPrts == 0 {
			EnumInit = false
			Type = ""
			Index = 0
		}

		// if match is found and enum initialized and parenthesis depth is 1
		if match := reEnumDec.FindSubmatchIndex(line); match != nil && inPrts == 1 && EnumInit {

			var tmp EnumDeclaration

			tmp.PackageID = pkg
			tmp.FileID = fileName

			tmp.StartOffset = match[0] + currentOffset // offset is current line offset + match index
			tmp.Length = match[1] - match[0]
			tmp.LineNumber = lineno

			tmp.EnumName = string(source[match[6]+currentOffset : match[7]+currentOffset])

			// If there is type declaration for Enum declaration
			// i.e. [enum] [type] = iota
			// set the type to type in declaration
			if match[2] != -1 {
				Type = string(source[match[4]+currentOffset : match[5]+currentOffset])
				tmp.EnumName = string(source[match[2]+currentOffset : match[3]+currentOffset])
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

func ExtractStructs(source []byte, fileName string) ([]StructDeclaration, error) {

	var StructDeclarationsArray []StructDeclaration
	var pkg string

	// Package
	rePkg := regexp.MustCompile("package")
	rePkgName := regexp.MustCompile(`(^|[\s])package\s+([_a-zA-Z][_a-zA-Z0-9]*)`)
	reStructHeader := regexp.MustCompile(`type\s+([_a-zA-Z][_a-zA-Z0-9]*)\s+struct`)
	reLeftBrace := regexp.MustCompile("{")
	reRightBrace := regexp.MustCompile("}")
	reStructField := regexp.MustCompile(`([_a-zA-Z][_a-zA-Z0-9]+)\s+[_a-zA-Z][_a-zA-Z0-9]+`)

	reader := bytes.NewReader(source)
	scanner := bufio.NewScanner(reader)
	scanner.Split(scanLinesWithLineTerminator) // set scanner SplitFunc to custom ScanLines func at line 55

	var currentOffset int                   // offset of current line
	var lineno int                          // line number
	var inBlock int                         // inBlock
	var structDeclaration StructDeclaration // temporary variable for Struct Declaration
	var structFieldsArray []*StructField    // struct fields

	//Reading code line by line
	for scanner.Scan() {
		line := scanner.Bytes()
		lineno++

		// Package declaration extraction
		if rePkg.FindIndex(line) != nil {
			tokens := bytes.Fields(line)

			if !bytes.Equal(tokens[0], []byte("package")) {
				col := bytes.IndexAny(line, string(tokens[0]))

				return StructDeclarationsArray, fmt.Errorf("%d:%d %s", lineno, col, "syntax error: unexpected IDENTIFIER")
			}

			if len(tokens) == 1 {
				col := bytes.LastIndex(line, tokens[0])

				return StructDeclarationsArray, fmt.Errorf("%d:%d %s", lineno, col, "syntax error: unexpected IDENTIFIER")
			}

			if len(tokens) > 2 {
				col := bytes.IndexAny(line, string(tokens[2]))

				return StructDeclarationsArray, fmt.Errorf("%d:%d %s", lineno, col, "syntax error: unexpected IDENTIFIER")
			}

			if match := rePkgName.FindStringSubmatch(string(line)); match != nil {
				pkg = match[2]
			}

		}

		// if struct declaration is found
		// i.e. type [name] [type]
		if match := reStructHeader.FindSubmatchIndex(line); match != nil && reLeftBrace.FindIndex(line)[0] > match[0] {

			structDeclaration.PackageID = pkg
			structDeclaration.FileID = fileName

			structDeclaration.StartOffset = match[0] + currentOffset // offset is current line offset + match index
			structDeclaration.Length = match[1] - match[0]
			structDeclaration.StructName = string(source[match[2]+currentOffset : match[3]+currentOffset])

			structDeclaration.LineNumber = lineno

			inBlock++

		}

		if match := reRightBrace.FindIndex(line); match != nil && inBlock == 1 {

			inBlock--
			structDeclaration.StructFields = structFieldsArray
			StructDeclarationsArray = append(StructDeclarationsArray, structDeclaration)
			structFieldsArray = []*StructField{}
		}

		if inBlock == 1 && structDeclaration.LineNumber < lineno {

			tokens := strings.Fields(string(line))

			if len(tokens) == 1 {
				return StructDeclarationsArray, errors.New("missing type")
			}

			if len(tokens) > 2 {
				return StructDeclarationsArray, errors.New("unexpected token")
			}

			var structField StructField

			if len(tokens) == 2 {

				match := reStructField.FindSubmatchIndex(line)

				structField.StartOffset = match[0] + currentOffset
				structField.Length = match[1] - match[0]
				structField.LineNumber = lineno
				structField.StructFieldName = string(source[match[2]+currentOffset : match[3]+currentOffset])
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
	rePkgName := regexp.MustCompile(`(^|[\s])package\s+([_a-zA-Z][_a-zA-Z0-9]*)`)
	reFunc := regexp.MustCompile(`func`)

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
			tokens := bytes.Fields(line)

			if !bytes.Equal(tokens[0], []byte("package")) {
				col := bytes.IndexAny(line, string(tokens[0]))

				return FuncDeclarationsArray, fmt.Errorf("%d:%d %s", lineno, col, "syntax error: unexpected IDENTIFIER")
			}

			if len(tokens) == 1 {
				col := bytes.LastIndex(line, tokens[0])

				return FuncDeclarationsArray, fmt.Errorf("%d:%d %s", lineno, col, "syntax error: unexpected IDENTIFIER")
			}

			if len(tokens) > 2 {
				col := bytes.IndexAny(line, string(tokens[2]))

				return FuncDeclarationsArray, fmt.Errorf("%d:%d %s", lineno, col, "syntax error: unexpected IDENTIFIER")
			}

			if match := rePkgName.FindStringSubmatch(string(line)); match != nil {
				pkg = match[2]
			}

		}

		if match := reFunc.FindIndex(line); match != nil {

			var funcDeclaration FuncDeclaration
			funcDeclaration.PackageID = pkg
			funcDeclaration.FileID = fileName
			funcDeclaration.StartOffset = match[0] + currentOffset
			funcDeclaration.LineNumber = lineno

			reFuncRegular := regexp.MustCompile(`func\s+([_a-zA-Z][_a-zA-Z0-9]*)`)
			reFuncMethod := regexp.MustCompile(`func\s*\(\s*[_a-zA-Z][_a-zA-Z0-9]*\s+\*{0,1}\s*[_a-zA-Z][_a-zA-Z0-9]*\s*\)\s*([_a-zA-Z][_a-zA-Z0-9]*)`)

			funcRegular := reFuncRegular.FindSubmatchIndex(line)
			funcMethod := reFuncMethod.FindSubmatchIndex(line)
			var funcNameIdx int

			if funcRegular == nil && funcMethod == nil {
				return FuncDeclarationsArray, fmt.Errorf("func err")
			}

			if funcRegular != nil {
				funcDeclaration.FuncName = string(line[funcRegular[2]:funcRegular[3]])
				funcNameIdx = funcRegular[1]
			}

			if funcMethod != nil {
				funcDeclaration.FuncName = string(line[funcMethod[2]:funcMethod[3]])
				funcNameIdx = funcMethod[1]
			}

			reInParen := regexp.MustCompile(`\(([\w\s\*\,]*)\)`)

			inParens := reInParen.FindAllSubmatchIndex(line[funcNameIdx:], -1)

			if inParens == nil || len(inParens) > 2 {
				return FuncDeclarationsArray, fmt.Errorf("parenthesis error")
			}

			for _, inParen := range inParens {

				inParenByte := line[inParen[2]+funcNameIdx : inParen[3]+funcNameIdx]

				reNotEmptyParen := regexp.MustCompile(`\S+`)

				if reNotEmptyParen.FindIndex(inParenByte) == nil {
					continue
				}

				params := bytes.Split(inParenByte, []byte(","))

				for _, param := range params {

					removeSpace := bytes.TrimSpace(param)

					reParam := regexp.MustCompile(`[_a-zA-Z][_a-zA-Z0-9]*\s+\*{0,1}\s*[_a-zA-Z][_a-zA-Z0-9]*`)

					if len(reParam.Find(removeSpace)) != len(removeSpace) {
						return FuncDeclarationsArray, fmt.Errorf("param err")
					}

				}

			}

			lastIndex := bytes.LastIndex(line, []byte(")"))

			funcDeclaration.Length = lastIndex - match[0] + 1

			FuncDeclarationsArray = append(FuncDeclarationsArray, funcDeclaration)

		}

		currentOffset += len(line) // increments the currentOffset by line len
	}

	return FuncDeclarationsArray, nil
}

func ReDeclarationCheck(Glbl []GlobalDeclaration, Enum []EnumDeclaration, Strct []StructDeclaration, Func []FuncDeclaration) error {

	var err error

	// Checks for the first declaration redeclared
	// in the order:
	// Globals -> Enums -> Struct -> Func

	for i := 0; i < len(Glbl); i++ {
		for j := i + 1; j < len(Glbl); j++ {
			if Glbl[i].GlobalVariableName == Glbl[j].GlobalVariableName {
				err = errors.New("global redeclared")
				return err
			}
		}
	}

	for i := 0; i < len(Enum); i++ {
		for j := i + 1; j < len(Enum); j++ {
			if Enum[i].EnumName == Enum[j].EnumName {
				err = errors.New("enum redeclared")
				return err
			}
		}
	}

	for i := 0; i < len(Strct); i++ {
		for j := i + 1; j < len(Strct); j++ {
			if Strct[i].StructName == Strct[j].StructName {
				err = errors.New("struct redeclared")
				return err
			}
		}
	}

	for _, structDeclaration := range Strct {
		StructFields := structDeclaration.StructFields
		for i := 0; i < len(StructFields); i++ {
			for j := i + 1; j < len(StructFields); j++ {
				if StructFields[i].StructFieldName == StructFields[j].StructFieldName {
					err = errors.New("struct field redeclared")
					return err
				}
			}
		}
	}

	for i := 0; i < len(Func); i++ {
		for j := i + 1; j < len(Func); j++ {
			if Func[i].FuncName == Func[j].FuncName {
				err = errors.New("func redeclared")
				return err
			}
		}
	}

	err = nil
	return err
}

func GetDeclarations(source []byte, Glbl []GlobalDeclaration, Enum []EnumDeclaration, Strct []StructDeclaration, Func []FuncDeclaration) []string {

	var declarations []string

	for i := 0; i < len(Glbl); i++ {
		declarations = append(declarations, string(source[Glbl[i].StartOffset:Glbl[i].StartOffset+Glbl[i].Length]))
	}

	for i := 0; i < len(Enum); i++ {
		declarations = append(declarations, string(source[Enum[i].StartOffset:Enum[i].StartOffset+Enum[i].Length]))
	}

	for i := 0; i < len(Strct); i++ {
		declarations = append(declarations, string(source[Strct[i].StartOffset:Strct[i].StartOffset+Strct[i].Length]))
	}

	for i := 0; i < len(Func); i++ {
		declarations = append(declarations, string(source[Func[i].StartOffset:Func[i].StartOffset+Func[i].Length]))
	}

	return declarations
}

func ExtractAllDeclarations(source []*os.File) ([]GlobalDeclaration, []EnumDeclaration, []StructDeclaration, []FuncDeclaration, error) {

	//Variable declarations
	var Globals []GlobalDeclaration
	var Enums []EnumDeclaration
	var Structs []StructDeclaration
	var Funcs []FuncDeclaration

	//Channel declarations
	globalChannel := make(chan []GlobalDeclaration, len(source))
	enumChannel := make(chan []EnumDeclaration, len(source))
	structChannel := make(chan []StructDeclaration, len(source))
	funcChannel := make(chan []FuncDeclaration, len(source))
	errorChannel := make(chan error, len(source))

	var wg sync.WaitGroup

	// concurrent extractions start
	for _, currentFile := range source {

		wg.Add(1)

		go func(currentFile *os.File, globalChannel chan<- []GlobalDeclaration, enumChannel chan<- []EnumDeclaration, structChannel chan<- []StructDeclaration, funcChannel chan<- []FuncDeclaration, errorChannel chan<- error, wg *sync.WaitGroup) {

			defer wg.Done()

			srcBytes, err := io.ReadAll(currentFile)
			if err != nil {
				errorChannel <- err
				return
			}

			fileName := currentFile.Name()
			replaceComments := ReplaceCommentsWithWhitespaces(srcBytes)

			wg.Add(4)

			go func(globalChannel chan<- []GlobalDeclaration, replaceComments []byte, fileName string, wg *sync.WaitGroup) {

				defer wg.Done()

				globals, err := ExtractGlobals(replaceComments, fileName)

				if err != nil {
					errorChannel <- err
					return
				}

				globalChannel <- globals

			}(globalChannel, replaceComments, fileName, wg)

			go func(enumChannel chan<- []EnumDeclaration, replaceComments []byte, fileName string, wg *sync.WaitGroup) {

				defer wg.Done()

				enums, err := ExtractEnums(replaceComments, fileName)

				if err != nil {
					errorChannel <- err
					return
				}

				enumChannel <- enums

			}(enumChannel, replaceComments, fileName, wg)

			go func(structChannel chan<- []StructDeclaration, replaceComments []byte, fileName string, wg *sync.WaitGroup) {

				defer wg.Done()

				structs, err := ExtractStructs(replaceComments, fileName)

				if err != nil {
					errorChannel <- err
					return
				}

				structChannel <- structs

			}(structChannel, replaceComments, fileName, wg)

			go func(funcChannel chan<- []FuncDeclaration, replaceComments []byte, fileName string, wg *sync.WaitGroup) {

				defer wg.Done()

				funcs, err := ExtractFuncs(replaceComments, fileName)

				if err != nil {
					errorChannel <- err
					return
				}

				funcChannel <- funcs

			}(funcChannel, replaceComments, fileName, wg)

		}(currentFile, globalChannel, enumChannel, structChannel, funcChannel, errorChannel, &wg)
	}

	wg.Wait()

	// Close all channels for reading
	close(globalChannel)
	close(enumChannel)
	close(structChannel)
	close(funcChannel)
	close(errorChannel)

	//Read from channels concurrently
	wg.Add(4)

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
		return Globals, Enums, Structs, Funcs, err
	}

	reDeclarationCheck := ReDeclarationCheck(Globals, Enums, Structs, Funcs)

	// there's declaration redeclared return values with error
	if reDeclarationCheck != nil {
		return Globals, Enums, Structs, Funcs, reDeclarationCheck
	}

	return Globals, Enums, Structs, Funcs, nil
}
