package declaration_extractor

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"os"
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

type StructDeclaration struct {
	PackageID   string // name of package
	FileID      string // name of file
	StartOffset int    // offset with in the file starting from 'type'
	Length      int    // length of entire declaration i.e. 'type [name] [type]'
	LineNumber  int    // line number of declaration
	StructName  string // name of struct being declared
}

type StructField struct {
	PackageID       string // name of package
	FileID          string // name of file
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

func ExtractPackages(source []byte) string {
	rePkgName := regexp.MustCompile(`(^|[\s])package[ \t]+([_a-zA-Z][_a-zA-Z0-9]*)`)

	// Only extracts from the first line
	firstLineTerminator := bytes.IndexByte(source, byte('\n'))

	// Gets the first line
	line := source[0:firstLineTerminator]

	// extract package name
	pkg := rePkgName.FindStringSubmatch(string(line))[2]

	return pkg
}

func ExtractGlobals(source []byte, fileName string, pkg string) ([]GlobalDeclaration, error) {

	var GlblDec []GlobalDeclaration
	var err error

	//Regexs
	reGlbl := regexp.MustCompile(`var\s([_a-zA-Z][_a-zA-Z0-9]*)\s+[\[_a-zA-Z][\]_a-zA-Z0-9]*`)
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

		// if {  is found increment body depth
		if locs := reBodyOpen.FindAllIndex(line, -1); locs != nil {
			inBlock++
		}

		// if } is found decrement body depth
		if locs := reBodyClose.FindAllIndex(line, -1); locs != nil {
			inBlock--
		}

		// if match is found and body depth is 0
		if match := reGlbl.FindSubmatchIndex(line); match != nil && inBlock == 0 {

			var tmp GlobalDeclaration

			tmp.PackageID = pkg
			tmp.FileID = fileName

			tmp.StartOffset = currentOffset + match[0] // offset is current line offset + match index
			tmp.Length = match[1] - match[0]
			tmp.LineNumber = lineno

			// gets the name directly with current line offset + submatch index
			tmp.GlobalVariableName = string(source[match[2]+currentOffset : match[3]+currentOffset])

			GlblDec = append(GlblDec, tmp)
		}

		currentOffset += len(line) // increments the currentOffset by line len
	}

	return GlblDec, err

}

func ExtractEnums(source []byte, fileName string, pkg string) ([]EnumDeclaration, error) {

	var EnumDec []EnumDeclaration
	var err error

	//Regexes
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
			EnumDec = append(EnumDec, tmp)
			Index++
		}

		currentOffset += len(line) // increments the currentOffset by line len
	}

	return EnumDec, err

}

func ExtractStructs(source []byte, fileName string, pkg string) ([]StructDeclaration, error) {

	var StrctDec []StructDeclaration
	var err error

	reStruct := regexp.MustCompile(`type\s+([_a-zA-Z][_a-zA-Z0-9]*)\s+[_a-zA-Z][_a-zA-Z0-9]*`)

	reader := bytes.NewReader(source)
	scanner := bufio.NewScanner(reader)
	scanner.Split(scanLinesWithLineTerminator) // set scanner SplitFunc to custom ScanLines func at line 55

	var currentOffset int // offset of current line
	var lineno int        // line number

	//Reading code line by line
	for scanner.Scan() {
		line := scanner.Bytes()
		lineno++

		// if struct declaration is found
		// i.e. type [name] [type]
		if match := reStruct.FindSubmatchIndex(line); match != nil {

			var tmp StructDeclaration

			tmp.PackageID = pkg
			tmp.FileID = fileName

			tmp.StartOffset = match[0] + currentOffset // offset is current line offset + match index
			tmp.Length = match[1] - match[0]
			tmp.StructName = string(source[match[2]+currentOffset : match[3]+currentOffset])

			tmp.LineNumber = lineno

			StrctDec = append(StrctDec, tmp)

		}
		currentOffset += len(line) // increments the currentOffset by line len
	}

	return StrctDec, err
}

func ExtractFuncs(source []byte, fileName string, pkg string) ([]FuncDeclaration, error) {

	var FuncDec []FuncDeclaration
	var err error

	reFunc := regexp.MustCompile(`func\s+([_a-zA-Z]\w*)\s*\(.*\)\s+\S+\w+|func\s+([_a-zA-Z]\w*)\s*\(.*\)`)

	reader := bytes.NewReader(source)
	scanner := bufio.NewScanner(reader)
	scanner.Split(scanLinesWithLineTerminator) // set scanner SplitFunc to custom ScanLines func at line 55

	var currentOffset int // offset of current line
	var lineno int        // line number

	//Reading code line by line
	for scanner.Scan() {
		line := scanner.Bytes()
		lineno++

		// if function declaration is found
		// i.e.  func [name] ([params]) ([returns])
		if match := reFunc.FindSubmatchIndex(line); match != nil {

			var tmp FuncDeclaration

			tmp.PackageID = pkg
			tmp.FileID = fileName

			tmp.StartOffset = match[0] + currentOffset // offset is current line offset + match index
			tmp.Length = match[1] - match[0]

			// If func has multiple or no returns
			// i.e. func [name] ([params]) ([returns]) or func [name] ([params])
			tmp.FuncName = string(source[match[4]+currentOffset : match[5]+currentOffset])

			// If func has one return
			// i.e. func [name] ([params]) [return]
			if match[2] != -1 {
				tmp.FuncName = string(source[match[2]+currentOffset : match[3]+currentOffset])
			}

			tmp.LineNumber = lineno

			FuncDec = append(FuncDec, tmp)

		}

		currentOffset += len(line) // increments the currentOffset by line len
	}

	return FuncDec, err
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
			pkg := ExtractPackages(replaceComments)

			wg.Add(4)

			go func(globalChannel chan<- []GlobalDeclaration, replaceComments []byte, fileName string, pkg string, wg *sync.WaitGroup) {

				defer wg.Done()

				globals, err := ExtractGlobals(replaceComments, fileName, pkg)

				if err != nil {
					errorChannel <- err
					return
				}

				globalChannel <- globals

			}(globalChannel, replaceComments, fileName, pkg, wg)

			go func(enumChannel chan<- []EnumDeclaration, replaceComments []byte, fileName string, pkg string, wg *sync.WaitGroup) {

				defer wg.Done()

				enums, err := ExtractEnums(replaceComments, fileName, pkg)

				if err != nil {
					errorChannel <- err
					return
				}

				enumChannel <- enums

			}(enumChannel, replaceComments, fileName, pkg, wg)

			go func(structChannel chan<- []StructDeclaration, replaceComments []byte, fileName string, pkg string, wg *sync.WaitGroup) {

				defer wg.Done()

				structs, err := ExtractStructs(replaceComments, fileName, pkg)

				if err != nil {
					errorChannel <- err
					return
				}

				structChannel <- structs

			}(structChannel, replaceComments, fileName, pkg, wg)

			go func(funcChannel chan<- []FuncDeclaration, replaceComments []byte, fileName string, pkg string, wg *sync.WaitGroup) {

				defer wg.Done()

				funcs, err := ExtractFuncs(replaceComments, fileName, pkg)

				if err != nil {
					errorChannel <- err
					return
				}

				funcChannel <- funcs

			}(funcChannel, replaceComments, fileName, pkg, wg)

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
